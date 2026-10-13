package lib

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// warmStartTestConfig is the scaled S3 upload policy's shape (file.uploadV2AdaptiveConcurrencyConfigWithInitial):
// a start of 150, raised only through measured warm windows.
func warmStartTestConfig(initialTarget int) AdaptiveConcurrencyConfig {
	return AdaptiveConcurrencyConfig{
		MaxConcurrency:                         1024,
		InitialTarget:                          initialTarget,
		MinTarget:                              8,
		GrowthCeiling:                          150,
		GrowthCeilingProbeBytes:                64 << 30,
		GrowthCeilingProbeRate:                 96 << 20,
		ThroughputFloor:                        50,
		GrowEvery:                              16,
		GrowStep:                               4,
		SqrtGrowth:                             true,
		FailureShrinkPercent:                   35,
		BackPressureShrinkPercent:              10,
		ThroughputWindow:                       32,
		ThroughputMinGainPercent:               1,
		ThroughputShrinkPercent:                8,
		ThroughputHoldWindows:                  1,
		ThroughputProbeMinWindows:              2,
		ThroughputProbeFloor:                   150,
		ThroughputProbeFloorRate:               96 << 20,
		ThroughputProbePlateauTarget:           200,
		ThroughputProbeMinGainPerTargetPercent: 0.15,
		ThroughputProbeLossTolerancePercent:    2,
		LatencyFloor:                           50,
		LatencyShrinkPercent:                   8,
		LatencyQueueHigh:                       160,
		LatencyGrowthQueueHigh:                 96,
		WarmStart:                              true,
	}
}

const warmTestPartBytes = 8 << 20

type warmSimEvent struct {
	at       time.Time
	answer   bool
	attempt  uint64
	duration time.Duration
	failed   bool
}

// warmSim admits every attempt the manager allows, on a clock of its own: each is sent after sendDelay and answered
// after duration, given how many attempts were sending as it started. It drives the manager only through its
// exported admission, sending and sample methods.
type warmSim struct {
	t         *testing.T
	m         *AdaptiveConcurrencyManager
	now       time.Time
	events    []warmSimEvent
	sending   int
	sent      int
	answered  int
	duration  func(attempt uint64, sending int) time.Duration
	sendDelay func(attempt uint64) time.Duration
	fails     func(attempt uint64) bool
	// peakBeforeEnd is the most attempts sending at once while the warm start was under way.
	peakBeforeEnd int
}

func newWarmSim(t *testing.T, config AdaptiveConcurrencyConfig, duration func(attempt uint64, sending int) time.Duration) *warmSim {
	s := &warmSim{t: t, now: time.Unix(1_000_000, 0), duration: duration}
	s.m = NewAdaptiveConcurrencyManagerWithConfig(config)
	s.m.clock = func() time.Time { return s.now }
	return s
}

func (s *warmSim) admit() {
	for {
		s.m.mu.Lock()
		ok := s.m.canAcquireLocked()
		s.m.mu.Unlock()
		if !ok {
			return
		}
		attempt, admitted := s.m.WaitForAttemptWithContext(s.t.Context())
		require.True(s.t, admitted)
		delay := time.Duration(0)
		if s.sendDelay != nil {
			delay = s.sendDelay(attempt)
		}
		s.events = append(s.events, warmSimEvent{at: s.now.Add(delay), attempt: attempt})
	}
}

// step admits what it can and plays the next event; it returns false once nothing is left to play.
func (s *warmSim) step() bool {
	s.admit()
	if len(s.events) == 0 {
		return false
	}
	sort.SliceStable(s.events, func(i, j int) bool {
		if !s.events[i].at.Equal(s.events[j].at) {
			return s.events[i].at.Before(s.events[j].at)
		}
		return s.events[i].answer && !s.events[j].answer
	})
	event := s.events[0]
	s.events = s.events[1:]
	s.now = event.at
	if !event.answer {
		s.m.AttemptSending(event.attempt)
		s.sending++
		s.sent++
		if s.m.Snapshot().Warming && s.sending > s.peakBeforeEnd {
			s.peakBeforeEnd = s.sending
		}
		d := s.duration(event.attempt, s.sending)
		failed := s.fails != nil && s.fails(event.attempt)
		s.events = append(s.events, warmSimEvent{at: s.now.Add(d), answer: true, attempt: event.attempt, duration: d, failed: failed})
		return true
	}
	s.sending--
	s.answered++
	s.m.DoneWithSample(AdaptiveConcurrencySample{Success: !event.failed, Duration: event.duration, Bytes: warmTestPartBytes, Attempt: event.attempt})
	return true
}

func (s *warmSim) runUntil(done func(AdaptiveConcurrencySnapshot) bool, maxAnswers int) AdaptiveConcurrencySnapshot {
	for s.answered < maxAnswers && !done(s.m.Snapshot()) {
		if !s.step() {
			break
		}
	}
	return s.m.Snapshot()
}

// A link where each request carries 25 MB/s on its own and the link has room for many: a part of 8 MiB takes 40 ms
// plus its transfer, and the link stops adding throughput only far above where the warm start reaches 96 MiB/s.
func fastLink(_ uint64, _ int) time.Duration {
	return 40*time.Millisecond + time.Duration(float64(warmTestPartBytes)/25e6*float64(time.Second))
}

// A home link of 6 MB/s that its requests share: each request added slows every request as much as it adds.
func saturatedLink(_ uint64, sending int) time.Duration {
	return 30*time.Millisecond + time.Duration(float64(warmTestPartBytes)*float64(sending)/6e6*float64(time.Second))
}

func TestAdaptiveConcurrencyWarmStartPromotesToItsStartOnAFastLink(t *testing.T) {
	sim := newWarmSim(t, warmStartTestConfig(150), fastLink)
	require.Equal(t, 1, sim.m.Target())

	snapshot := sim.runUntil(func(s AdaptiveConcurrencySnapshot) bool { return !s.Warming }, 200)

	assert.True(t, snapshot.WarmStarted)
	assert.False(t, snapshot.Warming)
	assert.Equal(t, "fast_link", snapshot.WarmStartEnd)
	assert.Equal(t, "fast_link", snapshot.WarmStartLastDecision.Decision)
	assert.Equal(t, 150, snapshot.WarmStartLastDecision.To)
	assert.GreaterOrEqual(t, snapshot.WarmStartLastDecision.Throughput, float64(96<<20), "the window it promoted on")
	assert.Empty(t, snapshot.WarmStartPhase)
	assert.Equal(t, 150, snapshot.Target, "a measured fast link goes straight to the policy's start")
	assert.Zero(t, snapshot.WarmStartSettleTotal, "no raise on a fast link is taken back")
	assert.LessOrEqual(t, sim.answered, 30, "the first answers that show the fast link decide")
	assert.LessOrEqual(t, sim.peakBeforeEnd, 6, "until then, only the attempts earned were sent")
}

func TestAdaptiveConcurrencyWarmStartKeepsASaturatedLinkNearOneAttempt(t *testing.T) {
	sim := newWarmSim(t, warmStartTestConfig(150), saturatedLink)

	snapshot := sim.runUntil(func(AdaptiveConcurrencySnapshot) bool { return false }, 300)

	assert.True(t, snapshot.Warming, "a link that adding requests does not speed up never looks fast")
	assert.LessOrEqual(t, snapshot.PeakRunning, 2, "at most one attempt more than the link could use was tried")
	assert.Positive(t, snapshot.WarmStartSettleTotal, "a raise that only queued was taken back")
	assert.LessOrEqual(t, snapshot.Target, 2)

	// The same policy started at its start sends 150 requests at once over that link.
	config := warmStartTestConfig(150)
	config.WarmStart = false
	aggressive := newWarmSim(t, config, saturatedLink)
	aggressive.step()
	assert.Equal(t, 150, aggressive.m.Snapshot().Running)
}

// A link where each request carries 3 MB/s and the link 40 MB/s: throughput rises with every attempt added up to 13,
// and never shows a fast link. Each kept raise adds half the target again, so the warm start reaches 12 attempts in
// five raises (1, 2, 3, 5, 8, 12) rather than eleven. Raises past what the link carries are taken back, and after
// each rest it tries one attempt more: it stays near the 13 the link can use.
func TestAdaptiveConcurrencyWarmStartRaisesByHalfAgainWhileAttemptsPay(t *testing.T) {
	mediumLink := func(_ uint64, sending int) time.Duration {
		rate := min(3e6, 40e6/float64(sending))
		return 30*time.Millisecond + time.Duration(float64(warmTestPartBytes)/rate*float64(time.Second))
	}
	sim := newWarmSim(t, warmStartTestConfig(150), mediumLink)

	snapshot := sim.runUntil(func(s AdaptiveConcurrencySnapshot) bool { return s.Target >= 12 }, 400)

	assert.Equal(t, 12, snapshot.Target)
	assert.Equal(t, 5, snapshot.WarmStartRaiseTotal)
	assert.Zero(t, snapshot.WarmStartSettleTotal)

	snapshot = sim.runUntil(func(AdaptiveConcurrencySnapshot) bool { return false }, 600)
	assert.True(t, snapshot.Warming)
	assert.Positive(t, snapshot.WarmStartSettleTotal)
	assert.GreaterOrEqual(t, snapshot.Target, 12)
	assert.LessOrEqual(t, snapshot.Target, 14)
}

// warmScript lets each attempt of a warm start take a time of its own and wait before it is sent, and records the
// manager's target each time it changes.
type warmScript struct {
	*warmSim
	targets []int
}

func newWarmScript(t *testing.T, config AdaptiveConcurrencyConfig, duration func(attempt uint64, sending int) time.Duration, sendDelay func(attempt uint64) time.Duration) *warmScript {
	s := &warmScript{warmSim: newWarmSim(t, config, duration)}
	s.sendDelay = sendDelay
	s.targets = []int{s.m.Target()}
	return s
}

func (s *warmScript) run(maxAnswers int) {
	for s.answered < maxAnswers && s.step() {
		if target := s.m.Target(); target != s.targets[len(s.targets)-1] {
			s.targets = append(s.targets, target)
		}
	}
}

// Attempts 1 to 3 run at one attempt: the warming window, then the baseline's two answers. The baseline's attempts
// wait a quarter of their time to be sent, so its window delivers only 80% of what its answers' times show. From
// attempt 4 on, the raised target of two answers in raisedTime each.
func baselineThenRaise(baseline time.Duration, raisedTime time.Duration) (func(uint64, int) time.Duration, func(uint64) time.Duration) {
	duration := func(attempt uint64, _ int) time.Duration {
		if attempt <= 3 {
			return baseline
		}
		return raisedTime
	}
	sendDelay := func(attempt uint64) time.Duration {
		if attempt == 2 || attempt == 3 {
			return baseline / 4
		}
		return 0
	}
	return duration, sendDelay
}

func TestAdaptiveConcurrencyWarmStartTakesBackARaiseWhoseRequestsQueueWhateverItsGain(t *testing.T) {
	// The baseline delivers 0.8 parts per answer time; two attempts each taking 1.4 times as long deliver 1.43 of
	// them, more than the 1.5 times the baseline required, but every byte took 40% longer than at one attempt: the
	// second attempt mostly waited behind the first (a queue estimate of 0.57 of the one attempt added).
	duration, sendDelay := baselineThenRaise(time.Second, 1400*time.Millisecond)
	script := newWarmScript(t, warmStartTestConfig(150), duration, sendDelay)

	script.run(8)

	assert.Equal(t, []int{1, 2, 1}, script.targets)
	snapshot := script.m.Snapshot()
	assert.Equal(t, 1, snapshot.WarmStartSettleTotal)
	// The withdrawal says why: the gain was there, the queue was not acceptable.
	withdrawal := snapshot.WarmStartLastWithdrawal
	assert.Equal(t, AdaptiveConcurrencyWarmDecision{Decision: "withdrawn", From: 2, To: 1, Answers: 2}, AdaptiveConcurrencyWarmDecision{Decision: withdrawal.Decision, From: withdrawal.From, To: withdrawal.To, Answers: withdrawal.Answers})
	assert.GreaterOrEqual(t, withdrawal.Throughput, withdrawal.Required)
	assert.GreaterOrEqual(t, withdrawal.QueueEstimate, withdrawal.QueueLimit)
	assert.Equal(t, 0.5, withdrawal.QueueLimit)
	assert.Equal(t, "resting", snapshot.WarmStartPhase)

	// Requests that take 1.2 times as long queue less than half an attempt (0.33): the raise is kept, and the next
	// one tried.
	duration, sendDelay = baselineThenRaise(time.Second, 1200*time.Millisecond)
	kept := newWarmScript(t, warmStartTestConfig(150), duration, sendDelay)
	kept.run(8)
	assert.Equal(t, []int{1, 2, 3}, kept.targets)
	assert.Zero(t, kept.m.Snapshot().WarmStartSettleTotal)
}

func TestAdaptiveConcurrencyWarmStartTakesBackARaiseThatAddsTooLittle(t *testing.T) {
	// Without waiting, two attempts taking 1.3 times as long deliver 1.54 times the baseline's throughput, with a
	// queue estimate of 0.46 of the attempt added: kept.
	duration := func(attempt uint64, _ int) time.Duration {
		if attempt <= 3 {
			return time.Second
		}
		return 1300 * time.Millisecond
	}
	kept := newWarmScript(t, warmStartTestConfig(150), duration, nil)
	kept.run(8)
	assert.Equal(t, []int{1, 2, 3}, kept.targets)

	// When the raised target's attempts each wait a fifth of their time to be sent, still judged, the same requests
	// deliver only 1.23 times the baseline, less than the added attempt's half share: taken back.
	slowSupply := newWarmScript(t, warmStartTestConfig(150), duration, func(attempt uint64) time.Duration {
		if attempt > 3 {
			return 330 * time.Millisecond
		}
		return 0
	})
	slowSupply.run(10)
	assert.Equal(t, []int{1, 2, 1}, slowSupply.targets)
	// The gain fell short; the queue alone would have kept it.
	withdrawal := slowSupply.m.Snapshot().WarmStartLastWithdrawal
	assert.Equal(t, "withdrawn", withdrawal.Decision)
	assert.Less(t, withdrawal.Throughput, withdrawal.Required)
	assert.Less(t, withdrawal.QueueEstimate, withdrawal.QueueLimit)
}

// With a custom ThroughputWindow of one, judged windows still hold two answers each. At 8 MiB/s, far from a fast
// link, the first raise waits for the warming window's answer and two baseline answers, not one; a raise taken back
// (attempts 4 and 5, twice as slow, add nothing) rests for two windows of two answers before trying again.
func TestAdaptiveConcurrencyWarmStartJudgesNoWindowOfOneAnswer(t *testing.T) {
	config := warmStartTestConfig(150)
	config.ThroughputWindow = 1
	duration := func(attempt uint64, _ int) time.Duration {
		if attempt == 4 || attempt == 5 {
			return 2 * time.Second
		}
		return time.Second
	}
	script := newWarmScript(t, config, duration, nil)
	var answersAt []int
	for script.answered < 20 && len(script.targets) < 4 && script.step() {
		if target := script.m.Target(); target != script.targets[len(script.targets)-1] {
			script.targets = append(script.targets, target)
			answersAt = append(answersAt, script.answered)
		}
	}

	assert.Equal(t, []int{1, 2, 1, 2}, script.targets)
	assert.Equal(t, []int{3, 5, 9}, answersAt, "raised after 1+2 answers, taken back, raised again after 2+2")
	withdrawal := script.m.Snapshot().WarmStartLastWithdrawal
	assert.Equal(t, "withdrawn", withdrawal.Decision)
	assert.GreaterOrEqual(t, withdrawal.Answers, 2)
	assert.Less(t, withdrawal.Throughput, withdrawal.Required)
}

func TestAdaptiveConcurrencyWarmStartJudgesARaiseOnlyOnceTheAttemptsThatFilledItHaveEnded(t *testing.T) {
	// Attempt 5, one of the two that first fill the raised target, takes 30 seconds; every other answer comes in a
	// second. Windows it overlaps judge nothing, so the raise is not judged until it ends, 25 answers later.
	duration := func(attempt uint64, _ int) time.Duration {
		if attempt == 5 {
			return 30 * time.Second
		}
		return time.Second
	}
	script := newWarmScript(t, warmStartTestConfig(150), duration, nil)

	for script.answered < 30 && script.now.Before(time.Unix(1_000_000, 0).Add(29*time.Second)) && script.step() {
	}
	snapshot := script.m.Snapshot()
	require.Equal(t, 2, snapshot.Target)
	assert.Equal(t, 1, snapshot.WarmStartRaiseTotal, "the raise to two waits for its slow attempt")

	script.run(60)
	assert.Greater(t, script.m.Snapshot().WarmStartRaiseTotal, 1, "once it ended, the raise was judged and kept")
}

func TestAdaptiveConcurrencyWarmStartLowersOnceForFailuresOfOneEpochDownToOneAttempt(t *testing.T) {
	sim := newWarmSim(t, warmStartTestConfig(150), func(uint64, int) time.Duration { return time.Second })
	snapshot := sim.runUntil(func(s AdaptiveConcurrencySnapshot) bool { return s.Target >= 5 }, 100)
	require.Equal(t, 5, snapshot.Target)
	require.True(t, snapshot.Warming)
	for sim.sending < 5 {
		require.True(t, sim.step())
	}

	// Every attempt sending now fails: they began in one epoch, so the target falls once, by 35%.
	var attempts []uint64
	for _, event := range sim.events {
		if event.answer {
			attempts = append(attempts, event.attempt)
		}
	}
	require.Len(t, attempts, 5)
	for _, attempt := range attempts {
		sim.m.DoneWithSample(AdaptiveConcurrencySample{Success: false, Attempt: attempt})
	}
	snapshot = sim.m.Snapshot()
	assert.Equal(t, 3, snapshot.Target)
	assert.True(t, snapshot.Warming)

	// Fresh attempts that fail lower it again, but never below one.
	for i := 0; i < 6; i++ {
		attempt, ok := sim.m.WaitForAttemptWithContext(t.Context())
		require.True(t, ok)
		sim.m.DoneWithSample(AdaptiveConcurrencySample{Success: false, Attempt: attempt})
	}
	assert.Equal(t, 1, sim.m.Target())
}

func TestAdaptiveConcurrencyWarmStartJudgesNoWindowItsTargetWasNotSending(t *testing.T) {
	// Fast answers, but each attempt waits as long to be sent as it then takes: the target sends half the time, so
	// no window is judged, and the warm start neither grows nor sees a fast link.
	sim := newWarmSim(t, warmStartTestConfig(150), func(uint64, int) time.Duration { return 40 * time.Millisecond })
	sim.sendDelay = func(uint64) time.Duration { return 40 * time.Millisecond }

	snapshot := sim.runUntil(func(AdaptiveConcurrencySnapshot) bool { return false }, 50)

	assert.True(t, snapshot.Warming)
	assert.Equal(t, 1, snapshot.Target)
	assert.Positive(t, snapshot.WarmStartUnjudgedWindowTotal)
}

func TestAdaptiveConcurrencyWarmStartEndsWithinAStartBelowItsDefault(t *testing.T) {
	// A caller cap of six makes six the policy's start: a fast link promotes there, and no further.
	config := warmStartTestConfig(6)
	config.MaxConcurrency = 6
	sim := newWarmSim(t, config, fastLink)

	snapshot := sim.runUntil(func(s AdaptiveConcurrencySnapshot) bool { return !s.Warming }, 200)

	assert.False(t, snapshot.Warming)
	assert.Equal(t, 6, snapshot.Target)
	assert.LessOrEqual(t, snapshot.PeakRunning, 6)
}

func TestAdaptiveConcurrencyWarmStartNeedsAStartAboveOne(t *testing.T) {
	config := warmStartTestConfig(1)
	manager := NewAdaptiveConcurrencyManagerWithConfig(config)
	assert.False(t, manager.Snapshot().WarmStarted)
	assert.Equal(t, 1, manager.Target())
}
