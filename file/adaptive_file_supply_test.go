package file

import (
	"testing"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
)

type stubStartupTargets struct {
	state lib.AdaptiveConcurrencyStartupState
}

func (s *stubStartupTargets) admissionTarget() (int, bool) {
	return s.state.Target, true
}

func (s *stubStartupTargets) admissionStartup() (lib.AdaptiveConcurrencyStartupState, bool) {
	return s.state, true
}

func newTestFileSupply(fileCap int) (*adaptiveFileSupply, *time.Time) {
	now := time.Unix(1_000_000, 0)
	supply := &adaptiveFileSupply{fileCap: fileCap, changedAt: now}
	supply.now = func() time.Time { return now }
	return supply, &now
}

// A warm start's files: while too few parts wait for its target, one more file each 100 ms; each raise of the warm
// start, its last included, as many more files at once as fed each part before; the ordinary policy's rises, a file
// at a time again; while enough parts wait, one fewer file each 100 ms, down to the target's limit.
func TestAdaptiveFileSupplyFeedsAWarmStartWithinTheFileCap(t *testing.T) {
	supply, now := newTestFileSupply(128)
	targets := &stubStartupTargets{state: lib.AdaptiveConcurrencyStartupState{Target: 1, StartTarget: 150, Warming: true, WarmStarted: true}}
	limit := func(running int) int {
		target, _ := targets.admissionTarget()
		targetLimit := max(target+target/4+1, 4)
		supplied, _ := supply.limit(targets, running, targetLimit)
		return supplied
	}

	*now = now.Add(50 * time.Millisecond)
	assert.Equal(t, 4, limit(4), "no sooner than 100 ms after it began")
	*now = now.Add(50 * time.Millisecond)
	assert.Equal(t, 5, limit(4))
	assert.Equal(t, 5, limit(5), "one file per 100 ms, however often it looks")
	*now = now.Add(100 * time.Millisecond)
	assert.Equal(t, 6, limit(5))

	targets.state.Target = 2
	assert.Equal(t, 12, limit(6), "a raise from one part to two needs twice the files")
	targets.state = lib.AdaptiveConcurrencyStartupState{Target: 3, StartTarget: 3, WarmStarted: true}
	assert.Equal(t, 18, limit(12), "so does the raise that ends the warm start")
	targets.state.Target = 4
	assert.Equal(t, 18, limit(18), "the policy's own rises are not carried")
	*now = now.Add(100 * time.Millisecond)
	assert.Equal(t, 19, limit(18))

	targets.state.Waiting = 2
	*now = now.Add(100 * time.Millisecond)
	assert.Equal(t, 18, limit(19), "enough parts waiting: one fewer per 100 ms")
	for range 20 {
		*now = now.Add(100 * time.Millisecond)
		limit(18)
	}
	assert.Equal(t, 6, limit(6), "down to the target's own limit")
}

func TestAdaptiveFileSupplyStaysWithinTheFileCap(t *testing.T) {
	supply, now := newTestFileSupply(10)
	targets := &stubStartupTargets{state: lib.AdaptiveConcurrencyStartupState{Target: 1, StartTarget: 150, Warming: true, WarmStarted: true}}
	for range 3 {
		*now = now.Add(100 * time.Millisecond)
		supply.limit(targets, 6, 4)
	}
	targets.state.Target = 5
	supplied, _ := supply.limit(targets, 7, 7)
	assert.Equal(t, 10, supplied)
}

// Without a warm start, under the aggressive startup or to another destination, files are admitted to the target's
// limit only, as before.
func TestAdaptiveFileSupplyLeavesUploadsWithoutAWarmStartAlone(t *testing.T) {
	supply, now := newTestFileSupply(128)
	targets := &stubStartupTargets{state: lib.AdaptiveConcurrencyStartupState{Target: 150, StartTarget: 150}}
	for range 5 {
		*now = now.Add(100 * time.Millisecond)
		supplied, _ := supply.limit(targets, 188, 188)
		assert.Equal(t, 188, supplied)
	}
}
