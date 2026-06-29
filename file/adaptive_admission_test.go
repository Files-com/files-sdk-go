package file

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

type stubPartTargets struct {
	target int
	ok     bool
}

func (s stubPartTargets) admissionTarget() (int, bool) {
	return s.target, s.ok
}

func TestAdaptiveFileAdmissionLimitNoSignal(t *testing.T) {
	if _, ok := adaptiveFileAdmissionLimit(stubPartTargets{}, 0); ok {
		t.Fatal("expected no limit before any shared adaptive manager exists")
	}
}

func TestAdaptiveFileAdmissionLimitUsesInitialTargetBeforeSignal(t *testing.T) {
	limit, ok := adaptiveFileAdmissionLimit(stubPartTargets{}, 4)
	if !ok || limit != 6 {
		t.Fatalf("limit = %d, %t, want 6, true", limit, ok)
	}
}

func TestAdaptiveFileAdmissionLimitUsesTargetWithGrowthHeadroom(t *testing.T) {
	limit, ok := adaptiveFileAdmissionLimit(stubPartTargets{target: 12, ok: true}, 4)
	if !ok || limit != 16 {
		t.Fatalf("limit = %d, %t, want 16, true", limit, ok)
	}
}

func TestAdaptiveFileAdmissionLimitAppliesFloor(t *testing.T) {
	limit, ok := adaptiveFileAdmissionLimit(stubPartTargets{target: 1, ok: true}, 4)
	if !ok || limit != manager.AdaptiveFileAdmissionFloor {
		t.Fatalf("limit = %d, %t, want %d, true", limit, ok, manager.AdaptiveFileAdmissionFloor)
	}
}

func TestWaitForAdaptiveFileAdmissionNoSignalAdmits(t *testing.T) {
	pool := lib.NewConstrainedWorkGroup(8)
	for range 8 {
		pool.Wait()
	}
	defer pool.WaitAllDone()
	defer func() {
		for range 8 {
			pool.Done()
		}
	}()

	if !waitForAdaptiveFileAdmission(context.Background(), pool, stubPartTargets{}, 0) {
		t.Fatal("expected admission without a learned signal")
	}
}

func TestWaitForAdaptiveFileAdmissionInitialTargetBlocksUntilFileFinishes(t *testing.T) {
	pool := lib.NewConstrainedWorkGroup(8)
	for range 6 {
		pool.Wait()
	}

	admitted := make(chan bool)
	go func() {
		admitted <- waitForAdaptiveFileAdmission(context.Background(), pool, stubPartTargets{}, 4)
	}()

	select {
	case <-admitted:
		t.Fatal("admission should block while the pool is at the initial target limit")
	case <-time.After(50 * time.Millisecond):
	}

	pool.Done()
	select {
	case ok := <-admitted:
		if !ok {
			t.Fatal("expected admission after a file finished")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("admission did not unblock after a file finished")
	}

	for range 5 {
		pool.Done()
	}
	pool.WaitAllDone()
}

func TestWaitForAdaptiveFileAdmissionBelowLimitAdmits(t *testing.T) {
	pool := lib.NewConstrainedWorkGroup(8)
	if !waitForAdaptiveFileAdmission(context.Background(), pool, stubPartTargets{target: 5, ok: true}, 4) {
		t.Fatal("expected admission below the learned limit")
	}
}

func TestWaitForAdaptiveFileAdmissionBlocksUntilFileFinishes(t *testing.T) {
	pool := lib.NewConstrainedWorkGroup(8)
	floor := manager.AdaptiveFileAdmissionFloor
	for range floor {
		pool.Wait()
	}

	admitted := make(chan bool)
	go func() {
		admitted <- waitForAdaptiveFileAdmission(context.Background(), pool, stubPartTargets{target: 1, ok: true}, 4)
	}()

	select {
	case <-admitted:
		t.Fatal("admission should block while the pool is at the limit")
	case <-time.After(50 * time.Millisecond):
	}

	pool.Done()
	select {
	case ok := <-admitted:
		if !ok {
			t.Fatal("expected admission after a file finished")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("admission did not unblock after a file finished")
	}

	for range floor - 1 {
		pool.Done()
	}
	pool.WaitAllDone()
}

func TestWaitForAdaptiveFileAdmissionReturnsFalseWhenCanceled(t *testing.T) {
	pool := lib.NewConstrainedWorkGroup(8)
	floor := manager.AdaptiveFileAdmissionFloor
	for range floor {
		pool.Wait()
	}

	ctx, cancel := context.WithCancel(context.Background())
	admitted := make(chan bool)
	go func() {
		admitted <- waitForAdaptiveFileAdmission(ctx, pool, stubPartTargets{target: 1, ok: true}, 4)
	}()

	cancel()
	select {
	case ok := <-admitted:
		if ok {
			t.Fatal("expected admission to fail after cancellation")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("admission did not return after cancellation")
	}

	for range floor {
		pool.Done()
	}
	pool.WaitAllDone()
}

func TestShouldGateAdaptiveUploadAdmissionSkipsNonTransferWork(t *testing.T) {
	directory := &UploadStatus{
		file: files_sdk.File{
			Type: "directory",
		},
		status: status.Queued,
		Mutex:  &sync.RWMutex{},
	}
	if shouldGateAdaptiveUploadAdmission(directory, true, nil, UploadV2Tuning{}) {
		t.Fatal("expected directories to bypass adaptive file admission")
	}

	file := &UploadStatus{
		file: files_sdk.File{
			Type: "file",
		},
		status: status.Queued,
		Mutex:  &sync.RWMutex{},
	}
	if !shouldGateAdaptiveUploadAdmission(file, true, nil, UploadV2Tuning{}) {
		t.Fatal("expected adaptive file uploads to use adaptive file admission")
	}
}

func TestUploadAdmissionInitialTargetBootstrapsBeforeSignal(t *testing.T) {
	job := (&Job{}).Init()
	limit, ok := adaptiveFileAdmissionLimit(uploadV2JobAdmissionTargets{job: job}, adaptiveFileAdmissionInitialTarget())
	if !ok || limit != 6 {
		t.Fatalf("limit = %d, %t, want 6, true", limit, ok)
	}
}

func TestJobFileAdmissionManagerCountsOnlyJobWork(t *testing.T) {
	shared := manager.New(8, 8, 8)
	jobA := (&Job{}).Init()
	jobA.SetManager(shared)
	jobB := (&Job{}).Init()
	jobB.SetManager(shared)

	for range 6 {
		jobA.fileAdmissionManager().Wait()
	}
	defer func() {
		for range 6 {
			jobA.fileAdmissionManager().Done()
		}
		jobA.fileAdmissionManager().WaitAllDone()
	}()

	if got := shared.FilesManager.RunningCount(); got != 6 {
		t.Fatalf("shared running count = %d, want 6", got)
	}
	if got := jobB.fileAdmissionManager().RunningCount(); got != 0 {
		t.Fatalf("job B running count = %d, want 0", got)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if !waitForAdaptiveFileAdmission(ctx, jobB.fileAdmissionManager(), stubPartTargets{}, 4) {
		t.Fatal("expected job B admission to ignore job A's running files")
	}
}

func TestShouldGateAdaptiveUploadAdmissionSkipsKnownUploadV2Fallback(t *testing.T) {
	file := &UploadStatus{
		file: files_sdk.File{
			Type: "file",
			Size: 64 * uploadV2MiB,
		},
		status: status.Queued,
		Mutex:  &sync.RWMutex{},
		UploadResumable: UploadResumable{FileUploadPart: files_sdk.FileUploadPart{
			Ref:           "upload-ref",
			ParallelParts: lib.Bool(false),
		}},
	}
	if shouldGateAdaptiveUploadAdmission(file, true, nil, UploadV2Tuning{}) {
		t.Fatal("expected known upload V2 fallbacks to bypass adaptive file admission")
	}

	file.UploadResumable.FileUploadPart.ParallelParts = lib.Bool(true)
	if !shouldGateAdaptiveUploadAdmission(file, true, nil, UploadV2Tuning{}) {
		t.Fatal("expected known upload V2-capable files to use adaptive file admission")
	}
}

func TestDownloadAdmissionTargetsSkipSmallFallback(t *testing.T) {
	size := downloadV2SmallFileFallbackSize()
	ranger := &downloadV2TestRangeFile{
		data:        make([]byte, size),
		downloadURI: "https://bucket.s3.us-east-1.amazonaws.com/small.bin?X-Amz-Signature=test",
		info: Info{File: files_sdk.File{
			DisplayName: "small.bin",
			Path:        "small.bin",
			Type:        "file",
			Size:        size,
			Crc32:       "crc32",
		}, sizeTrust: TrustedSizeValue},
	}
	reportStatus := downloadV2TestStatus(ranger, ranger.info, DownloaderParams{
		AdaptiveConcurrency: true,
	}, filepath.Join(t.TempDir(), "small.bin.download"))

	if shouldGateAdaptiveDownloadAdmission(reportStatus, true) {
		t.Fatal("expected small download fallback to bypass adaptive file admission")
	}
}

func TestDownloadAdmissionDoesNotFetchURI(t *testing.T) {
	size := 64 * uploadV2MiB
	ranger := &downloadV2TestRangeFile{
		data: make([]byte, size),
		info: Info{File: files_sdk.File{
			DisplayName: "large.bin",
			Path:        "large.bin",
			Type:        "file",
			Size:        int64(size),
			Crc32:       "crc32",
		}, sizeTrust: TrustedSizeValue},
	}
	reportStatus := downloadV2TestStatus(ranger, ranger.info, DownloaderParams{
		AdaptiveConcurrency: true,
	}, filepath.Join(t.TempDir(), "large.bin.download"))

	if !shouldGateAdaptiveDownloadAdmission(reportStatus, true) {
		t.Fatal("expected large download to use adaptive file admission")
	}
	if ranger.downloadURICalls != 0 {
		t.Fatalf("download URI calls = %d, want 0", ranger.downloadURICalls)
	}
}

func TestDownloadAdmissionSkipsNearCompleteResume(t *testing.T) {
	size := int64(64 * uploadV2MiB)
	tmpPath := filepath.Join(t.TempDir(), "large.bin.download")
	if err := os.WriteFile(tmpPath, make([]byte, size-downloadV2SmallFileFallbackSize()), 0644); err != nil {
		t.Fatal(err)
	}
	ranger := &downloadV2TestRangeFile{
		data: make([]byte, size),
		info: Info{File: files_sdk.File{
			DisplayName: "large.bin",
			Path:        "large.bin",
			Type:        "file",
			Size:        size,
			Crc32:       "crc32",
		}, sizeTrust: TrustedSizeValue},
	}
	reportStatus := downloadV2TestStatus(ranger, ranger.info, DownloaderParams{
		AdaptiveConcurrency: true,
	}, tmpPath)

	if shouldGateAdaptiveDownloadAdmission(reportStatus, true) {
		t.Fatal("expected near-complete resume fallback to bypass adaptive file admission")
	}
}

func TestUploadAdmissionTargetsUseJobManagers(t *testing.T) {
	fastJob := (&Job{}).Init()
	slowJob := (&Job{}).Init()
	fastJob.adaptiveUploadV2Managers = map[uploadV2AdaptiveManagerCacheKey]*lib.AdaptiveConcurrencyManager{
		{maxConcurrency: 64}: lib.NewAdaptiveConcurrencyManagerWithInitial(64, 9),
	}
	slowJob.adaptiveUploadV2Managers = map[uploadV2AdaptiveManagerCacheKey]*lib.AdaptiveConcurrencyManager{
		{maxConcurrency: 64}: lib.NewAdaptiveConcurrencyManagerWithInitial(64, 3),
	}

	target, ok := (uploadV2JobAdmissionTargets{job: fastJob}).admissionTarget()
	if !ok || target != 9 {
		t.Fatalf("fast job admissionTarget = %d, %t, want 9, true", target, ok)
	}
	target, ok = (uploadV2JobAdmissionTargets{job: slowJob}).admissionTarget()
	if !ok || target != 3 {
		t.Fatalf("slow job admissionTarget = %d, %t, want 3, true", target, ok)
	}

	emptyJob := (&Job{}).Init()
	if _, ok := (uploadV2JobAdmissionTargets{job: emptyJob}).admissionTarget(); ok {
		t.Fatal("expected no target from an empty job")
	}
}

func TestDownloadAdmissionTargetsUseJobManagers(t *testing.T) {
	fastJob := (&Job{}).Init()
	slowJob := (&Job{}).Init()
	fastJob.adaptiveDownloadV2Managers = map[downloadV2AdaptiveManagerCacheKey]*lib.AdaptiveConcurrencyManager{
		{target: downloadV2TargetS3, maxConcurrency: 64}: lib.NewAdaptiveConcurrencyManagerWithInitial(64, 7),
	}
	slowJob.adaptiveDownloadV2Managers = map[downloadV2AdaptiveManagerCacheKey]*lib.AdaptiveConcurrencyManager{
		{target: downloadV2TargetDefault, maxConcurrency: 64}: lib.NewAdaptiveConcurrencyManagerWithInitial(64, 3),
	}

	target, ok := (downloadV2JobAdmissionTargets{job: fastJob}).admissionTarget()
	if !ok || target != 7 {
		t.Fatalf("fast job admissionTarget = %d, %t, want 7, true", target, ok)
	}
	target, ok = (downloadV2JobAdmissionTargets{job: slowJob}).admissionTarget()
	if !ok || target != 3 {
		t.Fatalf("slow job admissionTarget = %d, %t, want 3, true", target, ok)
	}

	emptyJob := (&Job{}).Init()
	if _, ok := (downloadV2JobAdmissionTargets{job: emptyJob}).admissionTarget(); ok {
		t.Fatal("expected no target from an empty job")
	}
}
