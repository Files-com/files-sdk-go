package file

import (
	"context"
	"errors"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/directory"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

const (
	uploadV2S3ThroughputProbeFloor                   = 150
	uploadV2S3ThroughputProbePlateau                 = 200
	uploadV2S3ThroughputProbeFloorRate               = 96 * 1024 * 1024
	uploadV2S3ThroughputProbeMinGainPerTargetPercent = 0.15
	uploadV2S3ThroughputProbeLossTolerancePercent    = 2
	uploadV2S3MaxConcurrency                         = 1024
	uploadV2S3GrowthCeiling                          = 150
	uploadV2S3GrowthCeilingProbeBytes                = 64 * uploadV2GiB
	uploadV2S3GrowthCeilingProbeSuccesses            = 0
	uploadV2S3GrowthCeilingProbeRate                 = 96 * 1024 * 1024
	uploadV2S3InitialConcurrency                     = 150
	uploadV2S3AdaptiveFloor                          = 50
	uploadV2S3LatencyQueueHigh                       = 160
	uploadV2S3LatencyGrowthQueueHigh                 = 96
	uploadV2S3WorkloadTargetPartMultiplier           = 8
	uploadV2S3WorkloadMinPartSizeMiB                 = 8
	uploadV2WorkloadScanWait                         = 250 * time.Millisecond
	uploadV2WorkloadScanPoll                         = 10 * time.Millisecond
	uploadV2DefaultReadyRunwayParts                  = 4
	uploadV2DefaultReadyRunwayBytes                  = 256 * uploadV2MiB
)

// UploadV2Tuning overrides V2 upload defaults for diagnostics and benchmark
// tuning. Zero values keep the built-in defaults.
type UploadV2Tuning struct {
	// S3InitialTarget is the starting adaptive concurrency target for S3 uploads.
	S3InitialTarget int
	// S3AdaptiveFloor is the lowest adaptive target the S3 controller should shrink toward.
	S3AdaptiveFloor int
	// S3GrowEvery is the number of successful part samples required before a normal growth step.
	S3GrowEvery int
	// S3GrowStep is the number of additional concurrent part slots added during normal growth.
	S3GrowStep int
	// S3ThroughputWindow is the number of recent part samples used for throughput decisions.
	S3ThroughputWindow int
	// S3ThroughputMinGainPercent is the required throughput improvement for normal growth.
	S3ThroughputMinGainPercent int
	// S3ThroughputProbeMinWindows is the number of missed probe windows before backing off.
	S3ThroughputProbeMinWindows int
	// S3ThroughputProbeFloor is the fast-link target floor used when measured throughput is high.
	S3ThroughputProbeFloor int
	// S3ThroughputProbeFloorRateBytesPerSecond is the measured throughput required to use the probe floor.
	S3ThroughputProbeFloorRateBytesPerSecond int64
	// S3ThroughputProbePlateau is the target used as the first high-throughput probe plateau.
	S3ThroughputProbePlateau int
	// S3ThroughputShrinkPercent is the percent to shrink the target after throughput regression.
	S3ThroughputShrinkPercent int
	// S3ThroughputHoldWindows is the number of throughput windows to hold after a shrink.
	S3ThroughputHoldWindows int
	// S3ThroughputProbeMinGainPerTargetPercent is the minimum gain required per extra target above the plateau.
	S3ThroughputProbeMinGainPerTargetPercent float64
	// S3ThroughputProbeLossTolerancePercent is the tolerated throughput loss while probing above the plateau.
	S3ThroughputProbeLossTolerancePercent int
	// S3GrowthCeiling is the soft concurrency ceiling before large-workload probing unlocks higher targets.
	S3GrowthCeiling int
	// S3GrowthCeilingProbeBytes is the workload size required before probing above the soft ceiling.
	S3GrowthCeilingProbeBytes int64
	// S3GrowthCeilingProbeSuccesses is the successful part count required before probing above the soft ceiling.
	S3GrowthCeilingProbeSuccesses int
	// S3GrowthCeilingProbeRateBytesPerSecond is the throughput required before probing above the soft ceiling.
	S3GrowthCeilingProbeRateBytesPerSecond int64
	// S3LatencyQueueHigh is the observed queue/latency threshold that triggers backoff.
	S3LatencyQueueHigh float64
	// S3LatencyGrowthQueueHigh is the queue/latency threshold that suppresses further growth.
	S3LatencyGrowthQueueHigh float64
	// S3PartSizeMiB forces the known-size S3 part size in MiB. Zero uses the planner.
	S3PartSizeMiB int64
	// S3WorkloadBytes overrides the aggregate upload job size used by the workload-aware planner.
	S3WorkloadBytes int64
	// S3WorkloadTargetPartMultiplier sets desired planned parts per initial target for workload sizing.
	S3WorkloadTargetPartMultiplier int
	// S3WorkloadMinPartSizeMiB sets the minimum workload-tuned S3 part size in MiB.
	S3WorkloadMinPartSizeMiB int64
	// S3WorkloadScanWaitMillis is how long a job may wait for directory scanning before sizing from estimates.
	S3WorkloadScanWaitMillis int
}

// UploadWithV2 enables opt-in upload v2 behavior for this upload.
func UploadWithV2() UploadOption {
	return func(params uploadIO) (uploadIO, error) {
		params.uploadV2 = true
		return params, nil
	}
}

func uploadWithV2SDKDefaultCaps() UploadOption {
	return func(params uploadIO) (uploadIO, error) {
		params.uploadV2UseSDKDefaultCaps = true
		return params, nil
	}
}

// UploadWithV2ReadyRunway configures how many V2 upload parts may be prepared
// ahead of admitted HTTP concurrency. A parts value of 0 disables the runway.
// A bytes value of 0 leaves queued runway bytes uncapped.
func UploadWithV2ReadyRunway(parts int, bytes int64) UploadOption {
	return func(params uploadIO) (uploadIO, error) {
		if parts < 0 {
			return params, errors.New("upload v2 ready runway parts must be greater than or equal to zero")
		}
		if bytes < 0 {
			return params, errors.New("upload v2 ready runway bytes must be greater than or equal to zero")
		}
		params.uploadV2ReadyRunway = uploadV2ReadyRunwayConfig{
			parts:      parts,
			bytes:      bytes,
			configured: true,
		}
		return params, nil
	}
}

// UploadWithV2Tuning applies diagnostic tuning overrides to upload V2.
func UploadWithV2Tuning(tuning UploadV2Tuning) UploadOption {
	return func(params uploadIO) (uploadIO, error) {
		if err := tuning.validate(); err != nil {
			return params, err
		}
		params.uploadV2Tuning = tuning
		return params, nil
	}
}

func (t UploadV2Tuning) validate() error {
	if t.S3InitialTarget < 0 {
		return errors.New("upload v2 tuning s3 initial target must be greater than or equal to zero")
	}
	if t.S3AdaptiveFloor < 0 {
		return errors.New("upload v2 tuning s3 adaptive floor must be greater than or equal to zero")
	}
	if t.S3GrowEvery < 0 {
		return errors.New("upload v2 tuning s3 grow every must be greater than or equal to zero")
	}
	if t.S3GrowStep < 0 {
		return errors.New("upload v2 tuning s3 grow step must be greater than or equal to zero")
	}
	if t.S3ThroughputWindow < 0 {
		return errors.New("upload v2 tuning s3 throughput window must be greater than or equal to zero")
	}
	if t.S3ThroughputMinGainPercent < 0 {
		return errors.New("upload v2 tuning s3 throughput min gain percent must be greater than or equal to zero")
	}
	if t.S3ThroughputProbeMinWindows < 0 {
		return errors.New("upload v2 tuning s3 throughput probe min windows must be greater than or equal to zero")
	}
	if t.S3ThroughputProbeFloor < 0 {
		return errors.New("upload v2 tuning s3 throughput probe floor must be greater than or equal to zero")
	}
	if t.S3ThroughputProbeFloorRateBytesPerSecond < 0 {
		return errors.New("upload v2 tuning s3 throughput probe floor rate must be greater than or equal to zero")
	}
	if t.S3ThroughputProbePlateau < 0 {
		return errors.New("upload v2 tuning s3 throughput probe plateau must be greater than or equal to zero")
	}
	if t.S3ThroughputShrinkPercent < 0 || t.S3ThroughputShrinkPercent > 100 {
		return errors.New("upload v2 tuning s3 throughput shrink percent must be between zero and one hundred")
	}
	if t.S3ThroughputHoldWindows < 0 {
		return errors.New("upload v2 tuning s3 throughput hold windows must be greater than or equal to zero")
	}
	if t.S3ThroughputProbeMinGainPerTargetPercent < 0 {
		return errors.New("upload v2 tuning s3 throughput probe min gain per target percent must be greater than or equal to zero")
	}
	if t.S3ThroughputProbeLossTolerancePercent < 0 {
		return errors.New("upload v2 tuning s3 throughput probe loss tolerance percent must be greater than or equal to zero")
	}
	if t.S3GrowthCeiling < 0 {
		return errors.New("upload v2 tuning s3 growth ceiling must be greater than or equal to zero")
	}
	if t.S3GrowthCeilingProbeBytes < 0 {
		return errors.New("upload v2 tuning s3 growth ceiling probe bytes must be greater than or equal to zero")
	}
	if t.S3GrowthCeilingProbeSuccesses < 0 {
		return errors.New("upload v2 tuning s3 growth ceiling probe successes must be greater than or equal to zero")
	}
	if t.S3GrowthCeilingProbeRateBytesPerSecond < 0 {
		return errors.New("upload v2 tuning s3 growth ceiling probe rate must be greater than or equal to zero")
	}
	if t.S3LatencyQueueHigh < 0 {
		return errors.New("upload v2 tuning s3 latency queue high must be greater than or equal to zero")
	}
	if t.S3LatencyGrowthQueueHigh < 0 {
		return errors.New("upload v2 tuning s3 latency growth queue high must be greater than or equal to zero")
	}
	if t.S3PartSizeMiB < 0 {
		return errors.New("upload v2 tuning s3 part size mib must be greater than or equal to zero")
	}
	if t.S3WorkloadBytes < 0 {
		return errors.New("upload v2 tuning s3 workload bytes must be greater than or equal to zero")
	}
	if t.S3WorkloadTargetPartMultiplier < 0 {
		return errors.New("upload v2 tuning s3 workload target part multiplier must be greater than or equal to zero")
	}
	if t.S3WorkloadMinPartSizeMiB < 0 {
		return errors.New("upload v2 tuning s3 workload min part size mib must be greater than or equal to zero")
	}
	if t.S3WorkloadScanWaitMillis < 0 {
		return errors.New("upload v2 tuning s3 workload scan wait millis must be greater than or equal to zero")
	}
	return nil
}

func (t UploadV2Tuning) managerTuning() UploadV2Tuning {
	t.S3PartSizeMiB = 0
	t.S3WorkloadBytes = 0
	t.S3WorkloadTargetPartMultiplier = 0
	t.S3WorkloadMinPartSizeMiB = 0
	t.S3WorkloadScanWaitMillis = 0
	return t
}

type uploadV2AdaptiveManagerProvider func(uploadV2PartPlan, int, UploadV2Tuning) *lib.AdaptiveConcurrencyManager

func uploadWithV2AdaptiveManagerProvider(provider uploadV2AdaptiveManagerProvider) UploadOption {
	return func(params uploadIO) (uploadIO, error) {
		params.uploadV2ManagerProvider = provider
		return params, nil
	}
}

type uploadV2HTTPClientProvider func(*Client, uploadV2PartPlan, int, int) (*Client, uploadV2HTTPClientLimits, bool)

func uploadWithV2HTTPClientProvider(provider uploadV2HTTPClientProvider) UploadOption {
	return func(params uploadIO) (uploadIO, error) {
		params.uploadV2HTTPClientProvider = provider
		return params, nil
	}
}

type uploadV2HTTPClientCacheKey struct {
	target              uploadV2TargetClass
	maxConnsPerHost     int
	maxIdleConnsPerHost int
}

type uploadV2AdaptiveManagerCacheKey struct {
	target         uploadV2TargetClass
	maxConcurrency int
	tuning         UploadV2Tuning
}

type uploadV2HTTPClientCacheEntry struct {
	client *Client
	limits uploadV2HTTPClientLimits
}

type uploadV2ReadyRunwayConfig struct {
	parts      int
	bytes      int64
	configured bool
}

func (c uploadV2ReadyRunwayConfig) resolved() uploadV2ReadyRunwayConfig {
	if c.configured {
		return c
	}
	return uploadV2ReadyRunwayConfig{
		parts: uploadV2DefaultReadyRunwayParts,
		bytes: uploadV2DefaultReadyRunwayBytes,
	}
}

type uploadV2SharedAdaptiveManagerRegistry struct {
	mu       sync.Mutex
	managers map[uploadV2AdaptiveManagerCacheKey]*lib.AdaptiveConcurrencyManager
}

var uploadV2SharedAdaptiveManagers uploadV2SharedAdaptiveManagerRegistry

func (s *uploadV2SharedAdaptiveManagerRegistry) get(plan uploadV2PartPlan, maxConcurrency int, tuning UploadV2Tuning) *lib.AdaptiveConcurrencyManager {
	tuning = tuning.managerTuning()
	key := uploadV2AdaptiveManagerCacheKey{
		target:         plan.target,
		maxConcurrency: maxConcurrency,
		tuning:         tuning,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.managers == nil {
		s.managers = make(map[uploadV2AdaptiveManagerCacheKey]*lib.AdaptiveConcurrencyManager)
	}
	if manager := s.managers[key]; manager != nil {
		return manager
	}
	manager := lib.NewAdaptiveConcurrencyManagerWithConfig(uploadV2SharedAdaptiveConcurrencyConfig(plan, maxConcurrency, tuning))
	s.managers[key] = manager
	return manager
}

func (r *Job) uploadV2AdaptiveManager(plan uploadV2PartPlan, maxConcurrency int, tuning UploadV2Tuning) *lib.AdaptiveConcurrencyManager {
	return uploadV2SharedAdaptiveManagers.get(plan, maxConcurrency, tuning)
}

func (r *Job) uploadV2WorkloadBytes(currentFileSize int64, tuning UploadV2Tuning) int64 {
	if r == nil {
		return max(currentFileSize, 0)
	}
	if !r.uploadV2WorkloadScanComplete(tuning) {
		return 0
	}
	var total int64
	r.statusesMutex.RLock()
	for _, file := range r.Statuses {
		if file.Status().Any(status.Valid...) {
			total += file.Size()
		}
	}
	r.statusesMutex.RUnlock()
	if total < currentFileSize {
		total = currentFileSize
	}
	return max(total, 0)
}

func (r *Job) uploadV2WorkloadScanComplete(tuning UploadV2Tuning) bool {
	if r == nil || r.Type == directory.File || r.EndScanning == nil || r.EndScanning.Called() {
		return true
	}
	scanWait := uploadV2WorkloadScanWait
	if tuning.S3WorkloadScanWaitMillis > 0 {
		scanWait = time.Duration(tuning.S3WorkloadScanWaitMillis) * time.Millisecond
	}
	if scanWait <= 0 {
		return r.EndScanning.Called()
	}
	deadline := time.Now().Add(scanWait)
	for {
		if r.EndScanning.Called() {
			return true
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return r.EndScanning.Called()
		}
		sleep := uploadV2WorkloadScanPoll
		if sleep <= 0 || sleep > remaining {
			sleep = remaining
		}
		time.Sleep(sleep)
	}
}

func (r *Job) uploadV2HTTPClient(client *Client, plan uploadV2PartPlan, maxConnsPerHost int, maxIdleConnsPerHost int) (*Client, uploadV2HTTPClientLimits, bool) {
	if client == nil {
		return nil, uploadV2HTTPClientLimits{}, false
	}
	key := uploadV2HTTPClientCacheKey{
		target:              plan.target,
		maxConnsPerHost:     maxConnsPerHost,
		maxIdleConnsPerHost: maxIdleConnsPerHost,
	}
	r.adaptiveUploadV2Mu.Lock()
	defer r.adaptiveUploadV2Mu.Unlock()
	if r.adaptiveUploadV2Clients == nil {
		r.adaptiveUploadV2Clients = make(map[uploadV2HTTPClientCacheKey]uploadV2HTTPClientCacheEntry)
	}
	if entry, ok := r.adaptiveUploadV2Clients[key]; ok {
		return entry.client, entry.limits, true
	}
	adjustedClient, limits := configuredUploadV2HTTPClient(client, maxConnsPerHost, maxIdleConnsPerHost)
	if adjustedClient == nil || !limits.adjusted {
		return adjustedClient, limits, false
	}
	entry := uploadV2HTTPClientCacheEntry{client: adjustedClient, limits: limits}
	r.adaptiveUploadV2Clients[key] = entry
	return entry.client, entry.limits, true
}

func (u *uploadIO) runUploadV2(ctx context.Context) (UploadResumable, error, bool) {
	if !u.uploadV2Enabled() {
		return UploadResumable{}, nil, false
	}
	plan, ok, reason := u.newUploadV2PartPlanForUpload()
	if !ok {
		u.logUploadV2Fallback(reason)
		return UploadResumable{}, nil, false
	}

	u.uploadV2 = true
	engine := newUploadV2Engine(u, plan)
	if reason := engine.resumeResetReason(); reason != "" {
		u.logUploadV2(map[string]any{
			"timestamp": time.Now(),
			"event":     "upload v2 resume reset",
			"reason":    reason,
		})
		u.Parts = Parts{}
		var defaultSize int64
		if u.Size != nil {
			defaultSize = *u.Size
		}
		part, err := u.startUpload(ctx, files_sdk.FileBeginUploadParams{
			Size:         defaultSize,
			Path:         u.Path,
			MkdirParents: u.MkdirParents,
		})
		if err != nil {
			return u.UploadResumable(), err, true
		}
		u.FileUploadPart = part
		u.FileUploadPart.Path = u.Path
		plan, ok, reason = u.newUploadV2PartPlanForUpload()
		if !ok {
			u.logUploadV2Fallback(reason)
			return UploadResumable{}, nil, false
		}
		engine = newUploadV2Engine(u, plan)
	}
	resumable, err := engine.run(ctx)
	return resumable, err, true
}

func (u *uploadIO) uploadV2Enabled() bool {
	return u.uploadV2
}

func (u *uploadIO) logUploadV2(attrs map[string]any) {
	if u.Client == nil {
		return
	}
	u.LogPath(u.Path, attrs)
}

func (u *uploadIO) logUploadV2Fallback(reason string) {
	u.logUploadV2(map[string]any{
		"timestamp": time.Now(),
		"event":     "upload v2 fallback",
		"reason":    reason,
	})
	u.uploadV2 = false
}

func (u *uploadIO) logUploadV2Complete(err error, bytesWritten int64) {
	if !u.uploadV2 {
		return
	}
	attrs := map[string]any{
		"timestamp":     time.Now(),
		"event":         "upload v2 complete",
		"bytes_written": bytesWritten,
		"success":       err == nil,
	}
	if err != nil {
		attrs["error"] = err.Error()
	}
	if snapshotter, ok := u.Manager.(interface {
		Snapshot() lib.AdaptiveConcurrencySnapshot
	}); ok {
		snapshot := snapshotter.Snapshot()
		attrs["adaptive_final_target"] = snapshot.Target
		attrs["adaptive_max_target"] = snapshot.Max
		attrs["adaptive_growth_ceiling"] = snapshot.GrowthCeiling
		attrs["adaptive_growth_ceiling_unlocked"] = snapshot.GrowthCeilingUnlocked
		attrs["adaptive_growth_ceiling_probe_success_threshold"] = snapshot.GrowthCeilingProbeSuccessThreshold
		attrs["adaptive_peak_target"] = snapshot.PeakTarget
		attrs["adaptive_peak_running"] = snapshot.PeakRunning
		attrs["adaptive_success_total"] = snapshot.SuccessTotal
		attrs["adaptive_failure_total"] = snapshot.FailureTotal
		attrs["adaptive_grow_total"] = snapshot.GrowTotal
		attrs["adaptive_shrink_total"] = snapshot.ShrinkTotal
		attrs["adaptive_back_pressure_total"] = snapshot.BackPressureTotal
		attrs["adaptive_retry_after_total"] = snapshot.RetryAfterTotal
		attrs["adaptive_throughput_backoff_total"] = snapshot.ThroughputBackoffTotal
		attrs["adaptive_throughput_probe_miss_total"] = snapshot.ThroughputProbeMissTotal
		attrs["adaptive_throughput_probe_efficiency_miss_total"] = snapshot.ThroughputProbeEfficiencyMissTotal
		attrs["adaptive_latency_backoff_total"] = snapshot.LatencyBackoffTotal
		attrs["adaptive_latency_growth_suppression_total"] = snapshot.LatencyGrowthSuppressionTotal
		attrs["adaptive_bytes_total"] = snapshot.BytesTotal
		attrs["adaptive_average_duration_ms"] = snapshot.AverageDuration.Milliseconds()
		attrs["adaptive_last_throughput_bytes_per_second"] = snapshot.LastThroughputBytesPerSecond
		attrs["adaptive_best_throughput_bytes_per_second"] = snapshot.BestThroughputBytesPerSecond
		attrs["adaptive_last_throughput_probe_gain_percent"] = snapshot.LastThroughputProbeGainPercent
		attrs["adaptive_last_throughput_probe_target_delta"] = snapshot.LastThroughputProbeTargetDelta
		attrs["adaptive_last_throughput_probe_gain_per_target_percent"] = snapshot.LastThroughputProbeGainPerTargetPercent
		attrs["adaptive_last_queue_estimate"] = snapshot.LastQueueEstimate
		attrs["adaptive_min_duration_per_byte"] = snapshot.MinDurationPerByte
		attrs["adaptive_last_duration_per_byte"] = snapshot.LastDurationPerByte
	}
	u.logUploadV2(attrs)
}

func (u *uploadIO) uploadV2MaxConcurrency() int {
	if u.managerSet && !u.uploadV2UseSDKDefaultCaps {
		if maxer, ok := u.Manager.(interface{ Max() int }); ok {
			return maxer.Max()
		}
	}

	switch classifyUploadV2Target(u.FileUploadPart) {
	case uploadV2TargetS3:
		return uploadV2S3MaxConcurrency
	case uploadV2TargetAgent:
		return 128
	case uploadV2TargetFIW:
		return 192
	default:
		return 64
	}
}

func (u *uploadIO) uploadV2InitialConcurrency(maxConcurrency int) int {
	plan, ok, _ := u.newUploadV2PartPlanForUpload()
	if !ok {
		return min(maxConcurrency, 8)
	}
	return uploadV2InitialConcurrencyForPlan(plan, maxConcurrency, u.uploadV2Tuning)
}

func uploadV2AdaptiveConcurrencyConfig(plan uploadV2PartPlan, maxConcurrency int) lib.AdaptiveConcurrencyConfig {
	return uploadV2AdaptiveConcurrencyConfigWithInitial(plan, maxConcurrency, uploadV2InitialConcurrencyForPlan(plan, maxConcurrency, UploadV2Tuning{}), UploadV2Tuning{})
}

func uploadV2SharedAdaptiveConcurrencyConfig(plan uploadV2PartPlan, maxConcurrency int, tuning UploadV2Tuning) lib.AdaptiveConcurrencyConfig {
	return uploadV2AdaptiveConcurrencyConfigWithInitial(plan, maxConcurrency, uploadV2InitialConcurrencyForSharedPlan(plan, maxConcurrency, tuning), tuning)
}

func uploadV2AdaptiveConcurrencyConfigWithInitial(plan uploadV2PartPlan, maxConcurrency int, initial int, tuning UploadV2Tuning) lib.AdaptiveConcurrencyConfig {
	config := lib.AdaptiveConcurrencyConfig{
		MaxConcurrency:            maxConcurrency,
		InitialTarget:             initial,
		GrowEvery:                 8,
		GrowStep:                  1,
		FailureShrinkPercent:      50,
		BackPressureShrinkPercent: 50,
		BackPressurePause:         time.Second,
	}
	switch plan.target {
	case uploadV2TargetS3:
		config.MinTarget = min(8, maxConcurrency)
		config.ThroughputFloor = min(uploadV2S3AdaptiveFloor, initial)
		config.GrowEvery = 16
		config.GrowStep = 4
		config.SqrtGrowth = true
		config.FailureShrinkPercent = 35
		config.BackPressureShrinkPercent = 10
		config.BackPressurePause = 0
		config.ThroughputWindow = 32
		config.ThroughputMinGainPercent = 1
		config.ThroughputShrinkPercent = 8
		config.ThroughputHoldWindows = 1
		config.ThroughputProbeMinWindows = 2
		config.ThroughputProbeFloor = min(uploadV2S3ThroughputProbeFloor, maxConcurrency)
		config.ThroughputProbeFloorRate = uploadV2S3ThroughputProbeFloorRate
		config.ThroughputProbePlateauTarget = min(uploadV2S3ThroughputProbePlateau, maxConcurrency)
		config.ThroughputProbeMinGainPerTargetPercent = uploadV2S3ThroughputProbeMinGainPerTargetPercent
		config.ThroughputProbeLossTolerancePercent = uploadV2S3ThroughputProbeLossTolerancePercent
		config.GrowthCeiling = min(uploadV2S3GrowthCeiling, maxConcurrency)
		config.GrowthCeilingProbeBytes = uploadV2S3GrowthCeilingProbeBytes
		config.GrowthCeilingProbeSuccesses = uploadV2S3GrowthCeilingProbeSuccesses
		config.GrowthCeilingProbeRate = uploadV2S3GrowthCeilingProbeRate
		config.LatencyFloor = min(uploadV2S3AdaptiveFloor, initial)
		config.LatencyShrinkPercent = 8
		config.LatencyQueueHigh = uploadV2S3LatencyQueueHigh
		config.LatencyGrowthQueueHigh = uploadV2S3LatencyGrowthQueueHigh
		applyS3UploadV2Tuning(&config, tuning, maxConcurrency, initial)
	case uploadV2TargetFIW:
		config.MinTarget = min(4, maxConcurrency)
		config.GrowEvery = 4
		config.GrowStep = min(2, max(1, maxConcurrency/32))
		config.FailureShrinkPercent = 50
		config.BackPressureShrinkPercent = 25
		config.BackPressurePause = 250 * time.Millisecond
	case uploadV2TargetAgent:
		config.MinTarget = min(2, maxConcurrency)
		config.GrowEvery = 6
		config.GrowStep = 1
		config.FailureShrinkPercent = 50
		config.BackPressureShrinkPercent = 35
		config.BackPressurePause = 500 * time.Millisecond
	}
	return config
}

func applyS3UploadV2Tuning(config *lib.AdaptiveConcurrencyConfig, tuning UploadV2Tuning, maxConcurrency int, initial int) {
	if tuning.S3AdaptiveFloor > 0 {
		floor := min(tuning.S3AdaptiveFloor, maxConcurrency)
		config.ThroughputFloor = min(floor, initial)
		config.LatencyFloor = min(floor, initial)
	}
	if tuning.S3GrowEvery > 0 {
		config.GrowEvery = tuning.S3GrowEvery
	}
	if tuning.S3GrowStep > 0 {
		config.GrowStep = tuning.S3GrowStep
	}
	if tuning.S3ThroughputWindow > 0 {
		config.ThroughputWindow = tuning.S3ThroughputWindow
	}
	if tuning.S3ThroughputMinGainPercent > 0 {
		config.ThroughputMinGainPercent = tuning.S3ThroughputMinGainPercent
	}
	if tuning.S3ThroughputProbeMinWindows > 0 {
		config.ThroughputProbeMinWindows = tuning.S3ThroughputProbeMinWindows
	}
	if tuning.S3ThroughputProbeFloor > 0 {
		config.ThroughputProbeFloor = min(tuning.S3ThroughputProbeFloor, maxConcurrency)
	}
	if tuning.S3ThroughputProbeFloorRateBytesPerSecond > 0 {
		config.ThroughputProbeFloorRate = float64(tuning.S3ThroughputProbeFloorRateBytesPerSecond)
	}
	if tuning.S3ThroughputProbePlateau > 0 {
		config.ThroughputProbePlateauTarget = min(tuning.S3ThroughputProbePlateau, maxConcurrency)
	}
	if tuning.S3ThroughputShrinkPercent > 0 {
		config.ThroughputShrinkPercent = tuning.S3ThroughputShrinkPercent
	}
	if tuning.S3ThroughputHoldWindows > 0 {
		config.ThroughputHoldWindows = tuning.S3ThroughputHoldWindows
	}
	if tuning.S3ThroughputProbeMinGainPerTargetPercent > 0 {
		config.ThroughputProbeMinGainPerTargetPercent = tuning.S3ThroughputProbeMinGainPerTargetPercent
	}
	if tuning.S3ThroughputProbeLossTolerancePercent > 0 {
		config.ThroughputProbeLossTolerancePercent = tuning.S3ThroughputProbeLossTolerancePercent
	}
	if tuning.S3GrowthCeiling > 0 {
		config.GrowthCeiling = min(tuning.S3GrowthCeiling, maxConcurrency)
	}
	if tuning.S3GrowthCeilingProbeBytes > 0 {
		config.GrowthCeilingProbeBytes = tuning.S3GrowthCeilingProbeBytes
	}
	if tuning.S3GrowthCeilingProbeSuccesses > 0 {
		config.GrowthCeilingProbeSuccesses = tuning.S3GrowthCeilingProbeSuccesses
	}
	if tuning.S3GrowthCeilingProbeRateBytesPerSecond > 0 {
		config.GrowthCeilingProbeRate = float64(tuning.S3GrowthCeilingProbeRateBytesPerSecond)
	}
	if tuning.S3LatencyQueueHigh > 0 {
		config.LatencyQueueHigh = tuning.S3LatencyQueueHigh
	}
	if tuning.S3LatencyGrowthQueueHigh > 0 {
		config.LatencyGrowthQueueHigh = tuning.S3LatencyGrowthQueueHigh
	}
}

func uploadV2InitialConcurrencyForPlan(plan uploadV2PartPlan, maxConcurrency int, tuning UploadV2Tuning) int {
	initial := uploadV2InitialConcurrencyForSharedPlan(plan, maxConcurrency, tuning)
	switch plan.target {
	case uploadV2TargetS3:
		if plan.totalSize != nil && plan.partSize > 0 {
			partCount := int(ceilDiv(*plan.totalSize, plan.partSize))
			if partCount > 0 && partCount < initial {
				initial = partCount
			}
		}
		return max(1, initial)
	default:
		return initial
	}
}

func uploadV2InitialConcurrencyForSharedPlan(plan uploadV2PartPlan, maxConcurrency int, tuning UploadV2Tuning) int {
	maxConcurrency = max(maxConcurrency, 1)
	switch plan.target {
	case uploadV2TargetS3:
		if tuning.S3InitialTarget > 0 {
			return min(maxConcurrency, tuning.S3InitialTarget)
		}
		return min(maxConcurrency, uploadV2S3InitialConcurrency)
	default:
		return min(8, maxConcurrency)
	}
}
