package file

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

// The simulations drive the real AdaptiveConcurrencyManager — with the real
// production S3 adaptive config — through the real admission gate using
// single-part files, the workload where queued part demand equals the
// admitted file count. That is the worst case for the gate/controller
// feedback loop. They guard the two properties this feature trades between:
//
//   - Fast network: gated admission must not inhibit the controller's
//     ramp-up relative to an ungated baseline. (Production keeps throughput
//     probes dormant above 96 MiB/s below target 150, so growth is
//     success-driven and admission follows it.)
//   - Slow network: gated admission must stay near the learned target
//     instead of filling the static file pool, so a few files run to
//     completion instead of many files starving.

type simTargets struct {
	acm *lib.AdaptiveConcurrencyManager
}

func (s simTargets) admissionTarget() (int, bool) {
	return s.acm.Target(), true
}

type admissionSimResult struct {
	finalTarget  int
	peakAdmitted int
}

type admissionSimProfile struct {
	gated     bool
	ceiling   int
	acm       *lib.AdaptiveConcurrencyManager
	partBytes int64
	// perPartDuration receives the number of parts currently in flight so a
	// profile can model a shared-bandwidth (slow) or uncongested (fast) link.
	perPartDuration func(concurrent int) time.Duration
	deadline        time.Duration
	stop            func(target, admitted int) bool
}

// runAdmissionSimulation models a many-file transfer of single-part files
// through the same admission path production uses.
func runAdmissionSimulation(t *testing.T, profile admissionSimProfile) admissionSimResult {
	t.Helper()

	pool := lib.NewConstrainedWorkGroup(profile.ceiling)
	targets := simTargets{acm: profile.acm}
	ctx, cancel := context.WithTimeout(context.Background(), profile.deadline)
	defer cancel()

	var inFlightParts atomic.Int64
	var peakAdmitted atomic.Int64
	var wg sync.WaitGroup

	for {
		if profile.stop(profile.acm.Target(), pool.RunningCount()) {
			break
		}
		if profile.gated && !waitForAdaptiveFileAdmission(ctx, pool, targets, manager.AdaptiveFileAdmissionFloor) {
			break
		}
		if ctx.Err() != nil {
			break
		}
		if !pool.WaitWithContext(ctx) {
			break
		}
		if admitted := int64(pool.RunningCount()); admitted > peakAdmitted.Load() {
			peakAdmitted.Store(admitted)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer pool.Done()
			if !profile.acm.WaitWithContext(ctx) {
				return
			}
			concurrent := int(inFlightParts.Add(1))
			duration := profile.perPartDuration(concurrent)
			timer := time.NewTimer(duration)
			select {
			case <-timer.C:
				inFlightParts.Add(-1)
				profile.acm.DoneWithSample(lib.AdaptiveConcurrencySample{
					Success:  true,
					Duration: duration,
					Bytes:    profile.partBytes,
				})
			case <-ctx.Done():
				timer.Stop()
				inFlightParts.Add(-1)
				profile.acm.DoneNeutral()
			}
		}()
	}

	cancel()
	wg.Wait()
	return admissionSimResult{finalTarget: profile.acm.Target(), peakAdmitted: int(peakAdmitted.Load())}
}

func simProductionS3ACM(initial int) *lib.AdaptiveConcurrencyManager {
	plan := uploadV2PartPlan{target: uploadV2TargetS3}
	return lib.NewAdaptiveConcurrencyManagerWithConfig(
		uploadV2AdaptiveConcurrencyConfigWithInitial(plan, 1024, initial, UploadV2Tuning{}),
	)
}

func TestAdmissionSimulationFastNetworkRampIsNotInhibited(t *testing.T) {
	if testing.Short() {
		t.Skip("simulation test")
	}
	// Uncongested link: constant part duration, throughput scales with
	// concurrency and stays far above the 96 MiB/s probe-floor rate.
	fast := func(int) time.Duration { return 2 * time.Millisecond }
	const rampGoal = 64

	gated := runAdmissionSimulation(t, admissionSimProfile{
		gated:           true,
		ceiling:         128,
		acm:             simProductionS3ACM(8),
		partBytes:       1 << 20,
		perPartDuration: fast,
		deadline:        20 * time.Second,
		stop:            func(target, _ int) bool { return target >= rampGoal },
	})

	ungated := runAdmissionSimulation(t, admissionSimProfile{
		gated:           false,
		ceiling:         128,
		acm:             simProductionS3ACM(8),
		partBytes:       1 << 20,
		perPartDuration: fast,
		deadline:        20 * time.Second,
		stop:            func(target, _ int) bool { return target >= rampGoal },
	})

	if ungated.finalTarget < rampGoal {
		t.Fatalf("baseline did not ramp; simulation profile is invalid: ungated target = %d, want >= %d", ungated.finalTarget, rampGoal)
	}
	if gated.finalTarget < rampGoal {
		t.Fatalf("gated admission pinned ramp-up: target = %d, ungated reached %d", gated.finalTarget, ungated.finalTarget)
	}
	if gated.peakAdmitted <= manager.AdaptiveFileAdmissionFloor {
		t.Fatalf("admission never followed the growing target: peak admitted = %d", gated.peakAdmitted)
	}
}

func TestAdmissionSimulationSlowNetworkStaysNearTarget(t *testing.T) {
	if testing.Short() {
		t.Skip("simulation test")
	}
	// Shared-bandwidth link: every concurrent part slows all others, total
	// throughput is flat (~5 MiB/s, far below the 96 MiB/s probe-floor rate),
	// so throughput probes engage and refuse growth.
	slow := func(concurrent int) time.Duration {
		return time.Duration(concurrent) * 12 * time.Millisecond
	}
	timed := func(deadline time.Time) func(int, int) bool {
		return func(int, int) bool { return time.Now().After(deadline) }
	}

	gatedACM := simProductionS3ACM(8)
	gated := runAdmissionSimulation(t, admissionSimProfile{
		gated:           true,
		ceiling:         50,
		acm:             gatedACM,
		partBytes:       64 << 10,
		perPartDuration: slow,
		deadline:        15 * time.Second,
		stop:            timed(time.Now().Add(3 * time.Second)),
	})

	ungated := runAdmissionSimulation(t, admissionSimProfile{
		gated:           false,
		ceiling:         50,
		acm:             simProductionS3ACM(8),
		partBytes:       64 << 10,
		perPartDuration: slow,
		deadline:        15 * time.Second,
		stop:            timed(time.Now().Add(3 * time.Second)),
	})

	if ungated.peakAdmitted < 40 {
		t.Fatalf("ungated baseline did not fill the static pool; simulation profile is invalid: peak admitted = %d", ungated.peakAdmitted)
	}
	// Allow slack above the final-target limit for transient growth probes.
	gatedLimit, _ := adaptiveFileAdmissionLimit(stubPartTargets{target: gatedACM.Target(), ok: true}, manager.AdaptiveFileAdmissionFloor)
	if gated.peakAdmitted > gatedLimit+16 {
		t.Fatalf("gated admission overshot the learned target: peak admitted = %d, learned-target limit = %d", gated.peakAdmitted, gatedLimit)
	}
	if gated.peakAdmitted >= ungated.peakAdmitted {
		t.Fatalf("gate provided no admission reduction: gated peak = %d, ungated peak = %d", gated.peakAdmitted, ungated.peakAdmitted)
	}
}
