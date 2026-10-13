package file

import (
	"context"
	"time"

	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

// adaptiveFileAdmissionRecheckInterval bounds how long admission waits before
// re-reading the learned part-concurrency target, which can grow without any
// file completing.
const adaptiveFileAdmissionRecheckInterval = time.Second

// adaptivePartTargets reports the learned part-concurrency target for the
// adaptive work this admission gate is protecting. ok is false until a shared
// manager exists.
type adaptivePartTargets interface {
	admissionTarget() (target int, ok bool)
}

// adaptiveStartupTargets also reports the target with its warm start's phase and
// the part attempts waiting for it (Job.uploadV2AdmissionStartup).
type adaptiveStartupTargets interface {
	admissionStartup() (lib.AdaptiveConcurrencyStartupState, bool)
}

// fileAdmissionPool is the subset of the FilesManager pool the admission gate
// reads. The pool's own Wait remains the static admission ceiling.
type fileAdmissionPool interface {
	RunningCount() int
	WaitForADoneWithContext(ctx context.Context) bool
}

func adaptiveFileAdmissionInitialTarget() int {
	return manager.AdaptiveFileAdmissionFloor
}

// waitForAdaptiveFileAdmission keeps the number of in-flight files near the
// learned part-concurrency target so a constrained connection runs a few
// files to completion instead of starving every admitted file. It blocks
// while the in-flight count is at or above max(floor, learned target) and
// returns false only when ctx ends. Before any shared adaptive manager
// exists the caller-provided initial target is used so the static pool is not
// filled before the adaptive manager has a chance to learn.
func waitForAdaptiveFileAdmission(ctx context.Context, pool fileAdmissionPool, targets adaptivePartTargets, initialTarget int) bool {
	return waitForSuppliedFileAdmission(ctx, pool, targets, initialTarget, nil)
}

// waitForSuppliedFileAdmission is waitForAdaptiveFileAdmission, with supply, if
// not nil, admitting more files while a warm start lacks parts to send.
func waitForSuppliedFileAdmission(ctx context.Context, pool fileAdmissionPool, targets adaptivePartTargets, initialTarget int, supply *adaptiveFileSupply) bool {
	for {
		limit, ok := adaptiveFileAdmissionLimit(targets, initialTarget)
		running := pool.RunningCount()
		if !ok || running < limit {
			return ctx.Err() == nil
		}
		wait := adaptiveFileAdmissionRecheckInterval
		if supply != nil {
			supplied, recheck := supply.limit(targets, running, limit)
			if running < supplied {
				return ctx.Err() == nil
			}
			if recheck > 0 {
				wait = recheck
			}
		}
		waitCtx, cancel := context.WithTimeout(ctx, wait)
		pool.WaitForADoneWithContext(waitCtx)
		cancel()
		if ctx.Err() != nil {
			return false
		}
	}
}

// adaptiveFileSupplyRecheck is how often the files a warm-starting job admits
// past its target's limit may change, by one file.
const adaptiveFileSupplyRecheck = 100 * time.Millisecond

// adaptiveFileSupply feeds a warming target when one-part files spend time beginning or finishing. It adjusts supply
// by one file per 100 ms based on queued part demand. Warm raises carry prepared capacity in proportion, including
// the final raise; later policy raises use paced growth. The job's file cap always bounds supply.
type adaptiveFileSupply struct {
	fileCap   int
	grown     int
	changedAt time.Time
	// The target and whether a warm start under way set it, read together, when
	// it last looked; 0 before it saw one.
	grownTarget int
	grownWarm   bool
	now         func() time.Time
}

func newAdaptiveFileSupply(pool fileAdmissionPool) *adaptiveFileSupply {
	fileCap := manager.AdaptiveUploadV2ConcurrentFiles
	if maxer, ok := pool.(interface{ Max() int }); ok && maxer.Max() > 0 {
		fileCap = maxer.Max()
	}
	return &adaptiveFileSupply{fileCap: fileCap, changedAt: time.Now(), now: time.Now}
}

// limit returns how many files the job may have in progress while it waits at
// targetLimit, and how soon to look again. It grows past targetLimit only for a
// job whose manager warm-started.
func (s *adaptiveFileSupply) limit(targets adaptivePartTargets, running int, targetLimit int) (int, time.Duration) {
	startup, ok := targets.(adaptiveStartupTargets)
	if !ok {
		return targetLimit, 0
	}
	state, ok := startup.admissionStartup()
	if !ok || !state.WarmStarted {
		return targetLimit, 0
	}
	now := s.now()
	if s.grownTarget == 0 {
		// Files grown before the first look fed the warm start's first target.
		s.grownTarget = 1
		s.grownWarm = true
	}
	if s.grown > 0 && state.Target > s.grownTarget && (s.grownWarm || state.Warming) {
		s.grown = min(s.fileCap, (s.grown*state.Target+s.grownTarget-1)/s.grownTarget)
	}
	s.grownTarget = state.Target
	s.grownWarm = state.Warming
	tooFew := state.Waiting < state.Target/4+1
	if now.Sub(s.changedAt) >= adaptiveFileSupplyRecheck && ((tooFew && running >= s.grown) || (!tooFew && s.grown > targetLimit)) {
		if tooFew {
			s.grown = min(s.fileCap, running+1)
		} else {
			s.grown--
		}
		s.changedAt = now
	}
	recheck := adaptiveFileSupplyRecheck
	if left := adaptiveFileSupplyRecheck - now.Sub(s.changedAt); tooFew && left > 0 {
		recheck = left
	}
	return max(targetLimit, s.grown), recheck
}

// adaptiveFileAdmissionLimit allows the learned target plus growth headroom.
// The headroom (target/4 + 1) always covers the controller's next growth step
// (max(growStep, sqrt(target)), and t/4+1 >= sqrt(t) for all t), so workloads
// whose queued part demand equals the admitted file count — single-part files
// — still present enough excess demand for a growth probe to demonstrate a
// throughput gain. Admission equal to the target would starve the probe and
// pin ramp-up on fast connections.
func adaptiveFileAdmissionLimit(targets adaptivePartTargets, initialTarget int) (int, bool) {
	target, ok := targets.admissionTarget()
	if !ok {
		if initialTarget <= 0 {
			return 0, false
		}
		target = initialTarget
	}
	return max(target+target/4+1, manager.AdaptiveFileAdmissionFloor), true
}
