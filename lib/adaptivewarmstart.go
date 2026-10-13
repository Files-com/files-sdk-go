package lib

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync/atomic"
	"time"
)

// Warm start (AdaptiveConcurrencyConfig.WarmStart): the cautious startup of the adaptive upload reference
// (verification/adaptive-upload-spec.md, "Warm-start transitions"). A manager that would start at InitialTarget,
// enough parts at once to flood a slow link, starts at one attempt instead and raises its target only as measured
// windows earn each raise:
//
//   - A warm window closes once it has an answer for each attempt of the target, and has lasted at least as long as
//     those answers took on average. A window is judged only if attempts sending filled the target for at least 75%
//     of it; otherwise too few parts were ready, or a shared limiter held them, and it teaches nothing.
//   - The first window at one attempt warms its connection; the second is the baseline.
//   - A raise from n adds ceil(n/2) attempts, or one after a rest. The raised target is judged by the first window
//     opened after the attempts sending when it first filled have ended (by their sending sequence, not by how long
//     they took), or once ThroughputWindow answers since the raise have shown it, whichever comes first.
//   - A raise from n by d is kept when its window delivered at least baseline * (1 + 0.5*d/n), and its queue estimate,
//     average sending attempts * (1 - fastest median seconds per byte judged at n / this window's median), is under
//     0.5*d. Otherwise the target goes back to n, and the warm start rests: ThroughputHoldWindows and then one more
//     window of ThroughputWindow answers, after which it adds one attempt. A baseline whose requests queue rests too.
//   - Judged windows hold at least two answers, whatever ThroughputWindow, and more where their answers' times vary
//     enough to hide the gain or queue judged (warmWindowAnswersLocked): a sampling rule of thumb, not a confidence
//     level. Only the first window, which is never judged, and one already showing a fast link close sooner.
//   - Any judged window delivering ThroughputProbeFloorRate ends the warm start at InitialTarget at once: the link
//     is fast. So does a raise reaching InitialTarget. From there everything else in the config applies as it would
//     have from InitialTarget; the warm start's measurements do not carry into it.
//   - A failure, or back-pressure, while warming rests where its shrink leaves the target, down to one attempt.
//     Only attempts admitted since the target last fell lower it again, so a burst of failures counts once.
//
// It measures attempts sending: those admitted with WaitForAttemptWithContext from when AttemptSending reports them,
// to when the sample or DoneNeutralAttempt that ends them gives back their number. An attempt admitted with Wait or
// WaitWithContext counts as sending from its admission to its end; its raise can then only be judged once
// ThroughputWindow answers have followed it.

const (
	adaptiveWarmStartTarget = 1
	// The share of what each attempt delivered before a raise that each attempt it added must deliver, and of the
	// attempts it added that may queue, for the raise to be kept.
	adaptiveWarmGainPerAttempt = 0.5
	// How many spreads of its queue estimate the point between no queue and a queue that keeps no raise should lie
	// from either, in sizing a judged window (warmWindowAnswersLocked).
	adaptiveWarmWindowMargin = 2.0
	// A window is judged only if its target was sending for at least this share of it.
	adaptiveWarmJudgedTimeAtTarget = 0.75
	// Windows of ThroughputWindow answers a rest takes after its hold windows before it adds an attempt.
	adaptiveWarmGrowthIntervalWindows = 1
	// The fewest answers a window the warm start judges may hold.
	adaptiveWarmMinJudgedAnswers = 2
)

type adaptiveWarmStart struct {
	started  bool
	warming  bool
	resting  bool
	settling bool
	// Whether it rested since it last kept a raise, so that its next raise adds one attempt.
	restedSinceRaise bool
	measuring        bool
	promotionTarget  int
	end              string

	// Windows judged at the start target, the throughput of the window it last rose after (level) at from, the
	// fastest median seconds per byte judged at the current target (levelLowest) and at the target a raise under
	// judgment came from (fromLowest), and the queue estimate last judged.
	windows     int
	level       float64
	from        int
	levelLowest float64
	fromLowest  float64
	queue       float64
	// Over the windows judged, the squares of each answer's seconds per byte less its window's mean, relative to
	// that mean, and their degrees of freedom: what sizes judged windows.
	variationSquares float64
	variationAnswers int

	// A raised target, until judged: the sending sequence of the attempt that first filled it (settleLast, 0 until
	// one did), how many of the attempts sending then still are, when the last of them ended, and the answers since
	// the raise.
	settleLast    uint64
	settleLeft    int
	settledAt     time.Time
	settleAnswers int
	// A rest: its hold windows left, then its windows toward the next raise.
	holdWindows         int
	windowsTowardGrowth int

	// The first attempt number admitted since the target last fell (newEpochLocked).
	epochStart uint64
	// Attempts sending, by number, with their sending sequence; attempts admitted without a number, which count as
	// sending until they end; and all of them.
	sending    map[uint64]uint64
	sent       uint64
	unnumbered int
	inFlight   int

	lastEvent           time.Time
	windowStart         time.Time
	closedWindowStart   time.Time
	windowAttempts      int
	windowBytes         int64
	windowSeconds       float64
	windowPerByte       []float64
	closedPerByte       []float64
	lastMedian          float64
	lastAverageInFlight float64
	windowAtTarget      time.Duration
	windowInFlight      float64

	raiseTotal    int
	settleTotal   int
	unjudgedTotal int
	// The answers in the window last closed, its last judgment, and its last raise taken back.
	closedAnswers  int
	lastDecision   AdaptiveConcurrencyWarmDecision
	lastWithdrawal AdaptiveConcurrencyWarmDecision
}

// AdaptiveConcurrencyWarmDecision is a judgment of the warm start (AdaptiveConcurrencySnapshot.WarmStartLastDecision,
// WarmStartLastWithdrawal): what it decided, from which target to which, and the window it decided on. It is the last
// such judgment when the snapshot was taken, not a record of every one. Decision is empty for none, or:
//
//   - "raised": a judged window earned a raise, or ended a rest; "raised_to_start": that raise reached InitialTarget,
//     ending the warm start; "fast_link": the window reached ThroughputProbeFloorRate, ending it at InitialTarget.
//   - "withdrawn": a raised target did not earn its keep, so it went back to To; Throughput below Required, or
//     QueueEstimate at or above QueueLimit, or both, say why.
//   - "queued": the baseline, or a rest's last window, queued: it rests again instead of raising.
//   - "failure", "back_pressure": a failed or throttled attempt lowered the target from From to To.
type AdaptiveConcurrencyWarmDecision struct {
	Decision string
	From     int
	To       int
	// Answers and Throughput (bytes per second) of the window judged; Required, the throughput a raise judged needed.
	Answers    int
	Throughput float64
	Required   float64
	// The window's queue estimate against the fastest window judged at the target compared with, and the estimate at
	// which a raise is not kept: half the attempts added, or that the next raise would add.
	QueueEstimate float64
	QueueLimit    float64
}

// String formats a decision for logs, or "" for none.
func (d AdaptiveConcurrencyWarmDecision) String() string {
	if d.Decision == "" {
		return ""
	}
	return fmt.Sprintf("%s %d->%d answers=%d throughput=%.0f required=%.0f queue=%.2f queue_limit=%.2f",
		d.Decision, d.From, d.To, d.Answers, d.Throughput, d.Required, d.QueueEstimate, d.QueueLimit)
}

// warmPhaseLocked names where the warm start is: "baseline" until its first raise, "settling" while a raise awaits
// its judgment, "resting", or "" when it is not warming.
func (a *AdaptiveConcurrencyManager) warmPhaseLocked() string {
	switch w := &a.warm; {
	case !w.warming:
		return ""
	case w.resting:
		return "resting"
	case w.settling:
		return "settling"
	default:
		return "baseline"
	}
}

// AdaptiveConcurrencyStartupState is a manager's target read together with the phase of its warm start, its start
// (InitialTarget, where a warm start ends), and the attempts waiting for admission (StartupState).
type AdaptiveConcurrencyStartupState struct {
	Target      int
	StartTarget int
	Warming     bool
	WarmStarted bool
	Waiting     int
}

// WaitForAttemptWithContext waits like WaitWithContext, and returns the number of the attempt it admitted, from one,
// which its AttemptSending and the sample or DoneNeutralAttempt that ends it give back.
func (a *AdaptiveConcurrencyManager) WaitForAttemptWithContext(ctx context.Context) (uint64, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.waitForCapacityLocked(ctx) {
		return 0, false
	}
	a.acquireLocked()
	return a.attempts, true
}

// AttemptSending reports that an admitted attempt is being sent, past any other limiter it waited for.
func (a *AdaptiveConcurrencyManager) AttemptSending(attempt uint64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	w := &a.warm
	if !w.warming || attempt == 0 {
		return
	}
	a.advanceWarmLocked(a.now())
	w.sent++
	w.sending[attempt] = w.sent
	w.inFlight++
	if w.settling && w.settleLast == 0 && w.inFlight >= a.target {
		w.settleLast = w.sent
		w.settleLeft = w.inFlight
	}
}

// DoneNeutralAttempt is DoneNeutral for a numbered attempt.
func (a *AdaptiveConcurrencyManager) DoneNeutralAttempt(attempt uint64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.endAttemptLocked(a.now(), attempt)
	a.releaseLocked()
}

// StartupState returns the target with the warm start's phase, read together, so that a target the warm start set
// is never read with the phase after it ended, and the attempts waiting for admission.
func (a *AdaptiveConcurrencyManager) StartupState() AdaptiveConcurrencyStartupState {
	a.mu.Lock()
	defer a.mu.Unlock()
	return AdaptiveConcurrencyStartupState{
		Target:      a.target,
		StartTarget: a.startTarget,
		Warming:     a.warm.warming,
		WarmStarted: a.warm.started,
		Waiting:     a.waiters,
	}
}

func (s *AdaptiveConcurrencySubWorker) WaitForAttemptWithContext(ctx context.Context) (uint64, bool) {
	attempt, ok := s.parent.WaitForAttemptWithContext(ctx)
	if ok {
		s.mu.Lock()
		s.wg.Add(1)
		atomic.AddInt32(&s.runningCount, 1)
		s.mu.Unlock()
	}
	return attempt, ok
}

func (s *AdaptiveConcurrencySubWorker) AttemptSending(attempt uint64) {
	s.parent.AttemptSending(attempt)
}

func (s *AdaptiveConcurrencySubWorker) DoneNeutralAttempt(attempt uint64) {
	s.done(func() {
		s.parent.DoneNeutralAttempt(attempt)
	})
}

func (a *AdaptiveConcurrencyManager) startWarm() {
	a.warm = adaptiveWarmStart{
		started:         true,
		warming:         true,
		promotionTarget: a.target,
		sending:         make(map[uint64]uint64),
		epochStart:      1,
	}
	a.target = adaptiveWarmStartTarget
	a.peakTarget = a.target
}

// advanceWarmLocked adds the time since the last event to the current window: time with the target sending, and
// attempts sending over time. The first admitted attempt opens the first window.
func (a *AdaptiveConcurrencyManager) advanceWarmLocked(now time.Time) {
	w := &a.warm
	if !w.warming {
		return
	}
	if !w.measuring {
		w.measuring = true
		w.lastEvent = now
		a.openWarmWindowLocked(now)
		return
	}
	if !now.After(w.lastEvent) {
		return
	}
	from := w.lastEvent
	w.lastEvent = now
	if w.windowStart.After(from) {
		from = w.windowStart
	}
	if !now.After(from) {
		return
	}
	elapsed := now.Sub(from)
	if w.inFlight >= a.target {
		w.windowAtTarget += elapsed
	}
	w.windowInFlight += float64(w.inFlight) * elapsed.Seconds()
}

// sendingUnnumberedLocked counts an attempt just admitted without a number as sending.
func (a *AdaptiveConcurrencyManager) sendingUnnumberedLocked() {
	w := &a.warm
	if !w.warming {
		return
	}
	a.advanceWarmLocked(a.now())
	w.unnumbered++
	w.inFlight++
}

// endAttemptLocked ends a sending attempt's part in the warm start's measurements; attempt 0 ends one admitted
// without a number.
func (a *AdaptiveConcurrencyManager) endAttemptLocked(now time.Time, attempt uint64) {
	w := &a.warm
	if attempt == 0 {
		if w.unnumbered > 0 {
			a.advanceWarmLocked(now)
			w.unnumbered--
			w.inFlight--
		}
		return
	}
	if len(w.sending) == 0 {
		return
	}
	sequence, ok := w.sending[attempt]
	if !ok {
		return
	}
	a.advanceWarmLocked(now)
	delete(w.sending, attempt)
	w.inFlight--
	if w.settling && sequence <= w.settleLast && w.settleLeft > 0 {
		w.settleLeft--
		if w.settleLeft == 0 {
			w.settledAt = now
		}
	}
}

func (a *AdaptiveConcurrencyManager) recordWarmLocked(now time.Time, sample AdaptiveConcurrencySample, backPressure bool) {
	a.advanceWarmLocked(now)
	w := &a.warm
	if sample.Success {
		a.successTotal++
	} else {
		a.failureTotal++
	}
	if sample.Success && !backPressure {
		if w.settling {
			w.settleAnswers++
		}
		a.addToWarmWindowLocked(now, sample)
		return
	}
	// The warm start rests from where the failure leaves the target. Attempts admitted before the target last fell
	// lower it no further.
	a.restWarmLocked()
	current := sample.Attempt == 0 || sample.Attempt >= w.epochStart
	from := a.target
	if backPressure {
		if current && a.lowerWarmLocked(a.backPressureShrink) {
			a.shrinkTotal++
			w.lastDecision = AdaptiveConcurrencyWarmDecision{Decision: "back_pressure", From: from, To: a.target}
		}
		// Measuring starts again once any pause has ended.
		a.openWarmWindowLocked(now)
		return
	}
	if current && a.lowerWarmLocked(a.failureShrink) {
		a.shrinkTotal++
		w.lastDecision = AdaptiveConcurrencyWarmDecision{Decision: "failure", From: from, To: a.target}
		a.openWarmWindowLocked(now)
	}
}

func (a *AdaptiveConcurrencyManager) lowerWarmLocked(percent int) bool {
	if percent <= 0 || a.target <= adaptiveWarmStartTarget {
		return false
	}
	next := a.target * (100 - percent) / 100
	if next >= a.target {
		next = a.target - 1
	}
	a.target = max(adaptiveWarmStartTarget, next)
	a.newWarmEpochLocked()
	return true
}

func (a *AdaptiveConcurrencyManager) newWarmEpochLocked() {
	a.warm.epochStart = a.attempts + 1
}

func (a *AdaptiveConcurrencyManager) openWarmWindowLocked(now time.Time) {
	w := &a.warm
	w.windowStart = now
	if a.pauseUntil.After(now) {
		w.windowStart = a.pauseUntil
	}
	w.windowAttempts = 0
	w.windowBytes = 0
	w.windowSeconds = 0
	w.windowPerByte = w.windowPerByte[:0]
	w.windowAtTarget = 0
	w.windowInFlight = 0
}

// addToWarmWindowLocked adds an accepted answer to the window and closes it when it holds an answer for each attempt
// of the target, or ThroughputWindow of them while resting, and has lasted as long as they took on average. A window
// that would be judged also waits for the answers warmWindowAnswersLocked asks for, unless it already shows a fast
// link.
func (a *AdaptiveConcurrencyManager) addToWarmWindowLocked(now time.Time, sample AdaptiveConcurrencySample) {
	w := &a.warm
	if !now.After(w.windowStart) {
		return
	}
	w.windowAttempts++
	w.windowBytes += sample.Bytes
	w.windowSeconds += sample.Duration.Seconds()
	if sample.Bytes > 0 {
		w.windowPerByte = append(w.windowPerByte, sample.Duration.Seconds()/float64(sample.Bytes))
	}
	elapsed := now.Sub(w.windowStart).Seconds()
	need := a.target
	if w.resting {
		need = max(adaptiveWarmMinJudgedAnswers, a.throughputWindow, a.target)
	}
	if w.windowAttempts < need || elapsed < w.windowSeconds/float64(w.windowAttempts) {
		return
	}
	if !w.resting && w.windowAttempts < a.warmWindowAnswersLocked() && float64(w.windowBytes)/elapsed < a.throughputProbeFloorRate {
		return
	}
	a.closeWarmWindowLocked(now)
}

func (a *AdaptiveConcurrencyManager) closeWarmWindowLocked(now time.Time) {
	w := &a.warm
	length := now.Sub(w.windowStart)
	throughput := float64(w.windowBytes) / length.Seconds()
	if w.windowBytes > 0 && w.windowSeconds > 0 {
		w.lastAverageInFlight = w.windowInFlight / length.Seconds()
	}
	judged := float64(w.windowAtTarget) >= adaptiveWarmJudgedTimeAtTarget*float64(length)
	anyBytes := w.windowBytes > 0
	w.closedAnswers = w.windowAttempts
	w.closedWindowStart = w.windowStart
	w.windowPerByte, w.closedPerByte = w.closedPerByte, w.windowPerByte
	w.lastMedian = adaptiveMedian(w.closedPerByte)
	a.openWarmWindowLocked(now)
	if !anyBytes {
		return
	}
	a.lastThroughput = throughput
	if !judged {
		w.unjudgedTotal++
		return
	}
	a.judgeWarmLocked(now, throughput)
}

// judgeWarmLocked judges a window of the warm start (see the top of this file).
func (a *AdaptiveConcurrencyManager) judgeWarmLocked(now time.Time, throughput float64) {
	w := &a.warm
	decision := AdaptiveConcurrencyWarmDecision{From: a.target, To: a.target, Answers: w.closedAnswers, Throughput: throughput}
	if throughput >= a.throughputProbeFloorRate {
		decision.Decision = "fast_link"
		decision.To = w.promotionTarget
		w.lastDecision = decision
		a.endWarmLocked("fast_link")
		return
	}
	// queued judges the window against reference, for a raise of attempts, into the decision.
	queued := func(reference float64, attempts int) bool {
		result := a.warmQueuedLocked(reference, attempts)
		decision.QueueEstimate = w.queue
		decision.QueueLimit = adaptiveWarmGainPerAttempt * float64(attempts)
		return result
	}
	switch {
	case w.resting:
		a.warmJudgedAtTargetLocked()
		if w.holdWindows > 0 {
			w.holdWindows--
			return
		}
		w.windowsTowardGrowth++
		if w.windowsTowardGrowth < adaptiveWarmGrowthIntervalWindows {
			return
		}
		w.resting = false
		if queued(w.levelLowest, a.warmStepLocked()) {
			decision.Decision = "queued"
			w.lastDecision = decision
			a.restWarmLocked()
			return
		}
	case !w.settling:
		// The start target's first window only warms the connection it opened.
		w.windows++
		if w.windows < 2 {
			return
		}
		a.warmJudgedAtTargetLocked()
		if queued(w.levelLowest, a.warmStepLocked()) {
			decision.Decision = "queued"
			w.lastDecision = decision
			a.restWarmLocked()
			return
		}
	case !a.warmSettledWindowLocked():
		return
	default:
		a.warmCountVariationLocked()
		added := a.target - w.from
		decision.Required = w.level * (1 + adaptiveWarmGainPerAttempt*float64(added)/float64(w.from))
		isQueued := queued(w.fromLowest, added)
		if throughput < decision.Required || isQueued {
			a.target = w.from
			a.shrinkTotal++
			w.settleTotal++
			decision.Decision = "withdrawn"
			decision.To = a.target
			w.lastDecision = decision
			w.lastWithdrawal = decision
			a.newWarmEpochLocked()
			a.openWarmWindowLocked(now)
			a.restWarmLocked()
			return
		}
		// Kept: this window is the level the next raise is judged against.
		w.restedSinceRaise = false
		w.levelLowest = w.lastMedian
	}
	w.level = throughput
	w.from = a.target
	a.raiseWarmLocked(now)
	decision.Decision = "raised"
	if !w.warming {
		decision.Decision = "raised_to_start"
	}
	decision.To = a.target
	w.lastDecision = decision
}

// warmJudgedAtTargetLocked counts a window judged at the target: it may be the fastest there, which the next raise
// is compared with, and its answers' times count toward their variation.
func (a *AdaptiveConcurrencyManager) warmJudgedAtTargetLocked() {
	w := &a.warm
	if w.lastMedian > 0 && (w.levelLowest <= 0 || w.lastMedian < w.levelLowest) {
		w.levelLowest = w.lastMedian
	}
	a.warmCountVariationLocked()
}

func (a *AdaptiveConcurrencyManager) warmCountVariationLocked() {
	w := &a.warm
	if len(w.closedPerByte) >= 2 {
		w.variationSquares += adaptiveRelativeSquares(w.closedPerByte)
		w.variationAnswers += len(w.closedPerByte) - 1
	}
}

// warmQueuedLocked reports whether the window just closed queued: its estimate, against reference, the fastest
// window judged at the target compared with (0 for none, never queued), is the work of half the given attempts or
// more. A link that added attempts could not use makes them and those already sending wait for each other, whether
// or not it also delivered more for another reason, such as a slow connection finishing.
func (a *AdaptiveConcurrencyManager) warmQueuedLocked(reference float64, attempts int) bool {
	w := &a.warm
	w.queue = 0
	if reference > 0 && w.lastMedian > 0 {
		w.queue = w.lastAverageInFlight * (1 - reference/w.lastMedian)
	}
	return w.queue >= adaptiveWarmGainPerAttempt*float64(attempts)
}

// warmStepLocked is how many attempts the next raise adds: half the target, rounded up, or one after a rest.
func (a *AdaptiveConcurrencyManager) warmStepLocked() int {
	if a.warm.restedSinceRaise {
		return 1
	}
	return max(1, (a.target+1)/2)
}

// warmWindowAnswersLocked is how many answers a window the warm start will judge needs: at least two, and, from the
// variation of answers' times so far, enough that 2 * variation * (raised * margin / (0.5 * added))^2 answers
// separate the queue that keeps no raise from none, within ThroughputWindow. Answers' times are correlated and
// compared by medians, so this only sizes the sample; it promises no confidence level.
func (a *AdaptiveConcurrencyManager) warmWindowAnswersLocked() int {
	w := &a.warm
	judged := w.windows > 0
	if w.settling {
		judged = w.settleLast > 0 && w.settleLeft == 0 && !w.windowStart.Before(w.settledAt)
	}
	if !judged {
		return 0
	}
	squares := w.variationSquares
	freedom := w.variationAnswers
	if len(w.windowPerByte) >= 2 {
		squares += adaptiveRelativeSquares(w.windowPerByte)
		freedom += len(w.windowPerByte) - 1
	}
	added := a.warmStepLocked()
	raised := a.target + added
	if w.settling {
		added = a.target - w.from
		raised = a.target
	}
	needed := 2.0
	if freedom > 0 {
		ratio := float64(raised) * adaptiveWarmWindowMargin / (adaptiveWarmGainPerAttempt * float64(added))
		needed = math.Max(2, math.Ceil(2*adaptiveSmallSampleVariation(squares, freedom)*ratio*ratio))
	}
	return max(adaptiveWarmMinJudgedAnswers, int(math.Min(float64(a.throughputWindow), needed)))
}

// warmSettledWindowLocked reports whether the window just closed may judge the raised target: it opened once every
// attempt sending when the target first filled had ended, or the target has had ThroughputWindow answers since the
// raise, so that a request that never answers does not hold the warm start forever.
func (a *AdaptiveConcurrencyManager) warmSettledWindowLocked() bool {
	w := &a.warm
	return (w.settleLast > 0 && w.settleLeft == 0 && !w.closedWindowStart.Before(w.settledAt)) ||
		w.settleAnswers >= max(a.throughputWindow, a.target)
}

func (a *AdaptiveConcurrencyManager) restWarmLocked() {
	w := &a.warm
	w.resting = true
	w.settling = false
	w.restedSinceRaise = true
	w.windows = 0
	w.level = 0
	w.levelLowest = 0
	w.holdWindows = a.throughputHold
	w.windowsTowardGrowth = 0
}

// raiseWarmLocked raises the target by warmStepLocked, judged against the fastest window at the target it rose from;
// reaching the policy's start ends the warm start there.
func (a *AdaptiveConcurrencyManager) raiseWarmLocked(now time.Time) {
	w := &a.warm
	step := a.warmStepLocked()
	if a.target+step >= w.promotionTarget {
		a.endWarmLocked("raised")
		return
	}
	a.target += step
	a.peakTarget = max(a.peakTarget, a.target)
	a.growTotal++
	w.raiseTotal++
	w.fromLowest = w.levelLowest
	w.levelLowest = 0
	w.settling = true
	w.settleLast = 0
	w.settleLeft = 0
	w.settleAnswers = 0
	a.openWarmWindowLocked(now)
}

// endWarmLocked ends the warm start at the policy's start (InitialTarget), from which the manager carries on as one
// that started there: nothing measured while warming carries into its growth, probes or backoff.
func (a *AdaptiveConcurrencyManager) endWarmLocked(end string) {
	w := &a.warm
	w.warming = false
	w.resting = false
	w.settling = false
	w.end = end
	w.windowPerByte = nil
	w.closedPerByte = nil
	a.target = w.promotionTarget
	a.peakTarget = max(a.peakTarget, a.target)
	a.growTotal++
	a.successes = 0
	a.probePending = false
	a.throughputProbeStartTarget = 0
	a.throughputProbeMisses = 0
	a.throughputHoldRemaining = 0
}

// adaptiveMedian returns the median of values, sorting them in place; zero for none.
func adaptiveMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sort.Float64s(values)
	middle := len(values) / 2
	if len(values)%2 == 1 {
		return values[middle]
	}
	return (values[middle-1] + values[middle]) / 2
}

// adaptiveRelativeSquares returns the squares of the values less their mean, relative to that mean.
func adaptiveRelativeSquares(values []float64) float64 {
	var mean float64
	for _, value := range values {
		mean += value / float64(len(values))
	}
	if mean <= 0 {
		return 0
	}
	var squares float64
	for _, value := range values {
		squares += (value - mean) * (value - mean)
	}
	return squares / (mean * mean)
}

// adaptiveSmallSampleVariation inflates the observed relative variance where few answers could understate it, by a
// Wilson-Hilferty-style allowance; a heuristic, not a confidence bound. One degree of freedom is unbounded.
func adaptiveSmallSampleVariation(squares float64, freedom int) float64 {
	if squares <= 0 {
		return 0
	}
	spread := 2.0 / (9 * float64(freedom))
	root := 1 - spread - adaptiveWarmWindowMargin*math.Sqrt(spread)
	if root <= 0 {
		return math.Inf(1)
	}
	return squares / (float64(freedom) * root * root * root)
}
