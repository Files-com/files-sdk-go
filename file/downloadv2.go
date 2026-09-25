package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	goFs "io/fs"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"golang.org/x/sync/singleflight"
)

const (
	downloadV2CopyBufferSize   = 1024 * 1024
	downloadV2ProgressBatch    = 1024 * 1024
	downloadV2RetryAttempts    = 3
	downloadV2OutputPreallocAt = "preallocated_temp_file_write_at"
)

const (
	downloadV2TargetS3      = TransferV2TargetS3
	downloadV2TargetDefault = TransferV2TargetDefault
	downloadV2TargetDirect  = TransferV2TargetDirect
)

type downloadV2Part struct {
	number int
	off    int64
	len    int64
}

type downloadV2PartResult struct {
	part         downloadV2Part
	bytes        int64
	duration     time.Duration
	statusCode   int
	backPressure bool
	retryAfter   time.Duration
	err          error
}

type downloadV2Engine struct {
	reportStatus *DownloadStatus
	ranger       ReaderRange
	file         *os.File
	manager      *lib.AdaptiveConcurrencyManager
	target       TransferV2TargetClass
	startOffset  int64
	partSize     int64
	// planReady is closed once the first range response has established the
	// transfer size and the remaining ranges have been planned from it.
	planReady chan struct{}

	mu sync.Mutex
	// totalSize starts as the listing/download-request metadata size, which can
	// be stale, and becomes the transfer size reported by the first range
	// response. Every later response must agree with it.
	totalSize              int64
	transferSizeKnown      bool
	parts                  []downloadV2Part
	completed              map[int64]int64
	contiguous             int64
	directTransferDisabled bool
}

// downloadV2Range describes the bytes a range response covers and the transfer
// size it reports.
type downloadV2Range struct {
	off, end int64
	total    int64
	trust    SizeTrust
}

// downloadV2TerminalError marks a part failure that retrying the same range
// cannot fix: the response contradicted the plan, the transfer size could not
// be trusted, or the output file could not be prepared.
type downloadV2TerminalError struct{ err error }

func (e downloadV2TerminalError) Error() string { return e.err.Error() }
func (e downloadV2TerminalError) Unwrap() error { return e.err }

var (
	errDownloadV2UntrustedTransferSize = errors.New("download v2 range response did not report a trusted transfer size")
	errDownloadSizeConflict            = errors.New("download size conflict")
)

// downloadV2TerminalPartError reports errors that must fail the attempt instead
// of being retried per range: a rejected download request (source changed), an
// unsatisfiable range, or a response that contradicts the plan.
func downloadV2TerminalPartError(err error) bool {
	var terminal downloadV2TerminalError
	return errors.As(err, &terminal) ||
		downloadSourceChanged(err) ||
		downloadV2StatusCode(err) == http.StatusRequestedRangeNotSatisfiable
}

type downloadV2AdaptiveManagerCacheKey struct {
	target         TransferV2TargetClass
	maxConcurrency int
	tuning         UploadV2Tuning
}

type downloadV2URIProvider interface {
	downloadV2URI(context.Context) (string, error)
}

type downloadV2SharedAdaptiveManagerRegistry struct {
	adaptiveManagerRegistry[downloadV2AdaptiveManagerCacheKey]
}

var (
	downloadV2CopyBufferPool = sync.Pool{
		New: func() any {
			buf := make([]byte, downloadV2CopyBufferSize)
			return &buf
		},
	}
	downloadV2URIRefreshGroup        singleflight.Group
	downloadV2SharedAdaptiveManagers downloadV2SharedAdaptiveManagerRegistry
)

// runDownloadV2IfSupported runs the adaptive engine when the file qualifies.
// When it declines (used is false) planStat is the FileInfo the fallback engine
// must plan from: the metadata stat, or an empty plan when the first range
// response showed the metadata size cannot be right.
func runDownloadV2IfSupported(ctx context.Context, reportStatus *DownloadStatus, remoteStat goFs.FileInfo, tmp destinationPath, startOffset int64) (used bool, finalSize int64, planStat goFs.FileInfo, err error) {
	params, _ := reportStatus.Job().Params.(DownloaderParams)
	if !params.AdaptiveConcurrency {
		return false, 0, remoteStat, nil
	}
	if params.AdaptiveDownloadV2TuningSet {
		if err := params.AdaptiveDownloadV2Tuning.validate(); err != nil {
			return true, 0, remoteStat, err
		}
	}
	target, totalSize, partSize, ok := downloadV2PlanIfSupported(ctx, reportStatus, remoteStat, startOffset)
	if !ok {
		return false, 0, remoteStat, nil
	}

	file, err := tmp.openFile(os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return true, 0, remoteStat, err
	}
	if err := reportStatus.recordTmpDownloadIdentity(file); err != nil {
		file.Close()
		return true, 0, remoteStat, err
	}
	ranger := reportStatus.fsFile.(ReaderRange)
	engine := newDownloadV2Engine(reportStatus, ranger, file, target, totalSize, startOffset, partSize, params)
	err = engine.Run(ctx)
	if err != nil && engine.ContiguousSize() == startOffset {
		// Nothing was written, so the fallback engine can take over cleanly.
		switch {
		case errors.Is(err, errDownloadV2UntrustedTransferSize):
			// The transfer does not report its size, so range responses cannot
			// validate the plan; the fallback engine verifies unknown-size
			// downloads through the download request status.
			reportStatus.Job().Config.LogPath(reportStatus.RemotePath(), map[string]interface{}{
				"message": "download v2 transfer size untrusted; using fallback download",
				"error":   err,
			})
			return false, 0, remoteStat, nil
		case startOffset == 0 && downloadV2StatusCode(err) == http.StatusRequestedRangeNotSatisfiable:
			// No byte exists at offset 0, so the transfer is probably empty
			// despite the metadata size. A full download confirms that against
			// its Content-Length instead of trusting the 416.
			reportStatus.Job().Config.LogPath(reportStatus.RemotePath(), map[string]interface{}{
				"message": "download v2 first range unsatisfiable; confirming transfer size with fallback download",
				"error":   err,
			})
			return false, 0, downloadV2EmptyPlanStat(remoteStat), nil
		}
	}
	if err != nil {
		return true, engine.ContiguousSize(), remoteStat, err
	}
	return true, engine.FinalSize(), remoteStat, nil
}

// downloadV2EmptyPlanStat plans the fallback as an empty file so it runs as one
// full download, which validates the body against Content-Length whatever size
// the transfer turns out to have.
func downloadV2EmptyPlanStat(remoteStat goFs.FileInfo) goFs.FileInfo {
	file, _ := remoteStat.Sys().(files_sdk.File)
	if file.DisplayName == "" {
		file.DisplayName = remoteStat.Name()
	}
	file.Size = 0
	return Info{File: file, sizeTrust: NullSizeTrust}
}

func downloadV2PlanIfSupported(ctx context.Context, reportStatus *DownloadStatus, remoteStat goFs.FileInfo, startOffset int64) (TransferV2TargetClass, int64, int64, bool) {
	params, _ := reportStatus.Job().Params.(DownloaderParams)
	if !params.AdaptiveConcurrency {
		return "", 0, 0, false
	}
	if reportStatus.Job().Manager.FilePartsManager.DownloadFilesAsSingleStream {
		return "", 0, 0, false
	}
	ranger, ok := reportStatus.fsFile.(ReaderRange)
	if !ok {
		return "", 0, 0, false
	}
	totalSize := remoteStat.Size()
	if totalSize < 0 || startOffset < 0 || startOffset > totalSize {
		return "", 0, 0, false
	}
	if untrusted, ok := remoteStat.(UntrustedSize); ok && untrusted.SizeTrust() == UntrustedSizeValue {
		return "", 0, 0, false
	}
	if totalSize-startOffset <= downloadV2SmallFileFallbackSize() {
		return "", 0, 0, false
	}
	target, ok := classifyDownloadV2Target(ctx, remoteStat, ranger, params.AdaptiveDownloadV2TargetClassifier)
	if !ok {
		return "", 0, 0, false
	}
	partSize := downloadV2KnownSizePartSize(target, totalSize)
	if totalSize-startOffset <= partSize {
		return "", 0, 0, false
	}
	return target, totalSize, partSize, true
}

func shouldGateAdaptiveDownloadAdmission(reportStatus *DownloadStatus, adaptiveAdmission bool) bool {
	if !adaptiveAdmission ||
		reportStatus == nil ||
		reportStatus.error != nil ||
		reportStatus.fsFile == nil ||
		reportStatus.FileInfo == nil ||
		reportStatus.dryRun ||
		reportStatus.FileInfo.IsDir() {
		return false
	}
	job := reportStatus.Job()
	if job == nil || job.Manager == nil || job.Manager.FilePartsManager.DownloadFilesAsSingleStream {
		return false
	}
	if _, ok := reportStatus.fsFile.(ReaderRange); !ok {
		return false
	}
	totalSize := reportStatus.FileInfo.Size()
	startOffset := downloadV2AdmissionStartOffset(reportStatus)
	if totalSize < 0 || startOffset < 0 || startOffset > totalSize {
		return false
	}
	if untrusted, ok := reportStatus.FileInfo.(UntrustedSize); ok && untrusted.SizeTrust() == UntrustedSizeValue {
		return false
	}
	remaining := totalSize - startOffset
	if remaining <= downloadV2SmallFileFallbackSize() {
		return false
	}
	return remaining > downloadV2AdmissionMinPartSize(totalSize)
}

func downloadV2AdmissionStartOffset(reportStatus *DownloadStatus) int64 {
	if reportStatus == nil || reportStatus.FileInfo == nil {
		return 0
	}
	tmp, ok := reportStatus.tmpDownloadToResume()
	if !ok {
		return 0
	}
	fi, err := tmp.stat()
	if err != nil {
		return 0
	}
	if fi.Size() > reportStatus.FileInfo.Size() {
		return 0
	}
	return fi.Size()
}

func downloadV2AdmissionMinPartSize(totalSize int64) int64 {
	return min(
		downloadV2KnownSizePartSize(downloadV2TargetS3, totalSize),
		downloadV2KnownSizePartSize(downloadV2TargetDefault, totalSize),
	)
}

func newDownloadV2Engine(reportStatus *DownloadStatus, ranger ReaderRange, file *os.File, target TransferV2TargetClass, totalSize int64, startOffset int64, partSize int64, params DownloaderParams) *downloadV2Engine {
	maxConcurrency := downloadV2MaxConcurrency(reportStatus.Job(), params, target)
	tuning := params.AdaptiveDownloadV2Tuning
	if !params.AdaptiveDownloadV2TuningSet {
		tuning = UploadV2Tuning{}
	}
	manager := lib.NewAdaptiveConcurrencyManagerWithConfig(downloadV2AdaptiveConcurrencyConfig(target, maxConcurrency, totalSize, partSize, tuning))
	if params.Manager == nil || params.AdaptiveConcurrencyUseSDKDefaultCaps {
		manager = reportStatus.Job().downloadV2AdaptiveManager(target, maxConcurrency, totalSize, partSize, tuning)
	}
	return &downloadV2Engine{
		reportStatus: reportStatus,
		ranger:       ranger,
		file:         file,
		manager:      manager,
		target:       target,
		totalSize:    totalSize,
		startOffset:  startOffset,
		partSize:     partSize,
		planReady:    make(chan struct{}),
		parts:        downloadV2BuildParts(startOffset, totalSize, partSize),
		completed:    make(map[int64]int64),
		contiguous:   startOffset,
	}
}

func (e *downloadV2Engine) Run(parentCtx context.Context) (err error) {
	defer func() {
		if err != nil {
			// The file was preallocated to the transfer size and parts land at
			// their offsets, so it must end at the received prefix before a
			// resume can continue from its length. When that fails the caller
			// is told, so the file is not kept for a resume.
			if trimErr := truncateTmpDownload(e.file, e.ContiguousSize()); trimErr != nil {
				err = errors.Join(err, tmpDownloadNotTrimmedError{cause: trimErr})
			}
		}
		closeErr := e.file.Close()
		if err == nil {
			err = closeErr
		}
	}()

	e.logStart()
	if len(e.parts) == 0 {
		return downloadV2PreallocateFile(e.file, e.totalSize)
	}

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	ctx, closeDirectClients := files_sdk.WithDirectTransferClientCache(ctx)
	defer closeDirectClients()
	ctx = context.WithValue(ctx, directTransferDownloadSuppressorContextKey{}, e)

	results := make(chan downloadV2PartResult, max(1, e.manager.Max()))
	var wg sync.WaitGroup
	go func() {
		defer func() {
			wg.Wait()
			close(results)
		}()
		started := false
		for i := 0; ; i++ {
			part, ok := e.part(i)
			if !ok || ctx.Err() != nil {
				return
			}
			if !e.manager.WaitWithContext(ctx) {
				return
			}
			if !started {
				e.reportStatus.Job().UpdateStatusWithBytes(status.Downloading, e.reportStatus, 0)
				started = true
			}
			wg.Add(1)
			go func(part downloadV2Part) {
				defer wg.Done()
				result := e.downloadPartWithRetry(ctx, part)
				e.manager.DoneWithSample(result.sample())
				results <- result
			}(part)
			// The first range response establishes the transfer size; the
			// remaining ranges are planned from it, not from listing metadata.
			if i == 0 && !e.waitForPlan(ctx) {
				return
			}
		}
	}()

	for result := range results {
		if result.err != nil {
			cancel()
			// A source change rejects the whole download request; report it over
			// the cancellation errors of the other ranges.
			if err == nil || (downloadSourceChanged(result.err) && !downloadSourceChanged(err)) {
				err = result.err
			}
			continue
		}
		e.markComplete(result.part, result.bytes)
	}
	if err != nil {
		return err
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if e.ContiguousSize() != e.FinalSize() {
		return fmt.Errorf("download v2 wrote non-contiguous file. expected: %v, actual: %v", e.FinalSize(), e.ContiguousSize())
	}
	e.logFinish()
	return nil
}

// FinalSize is the transfer size: the size reported by the range responses once
// the first one has been accepted, otherwise the planned metadata size.
func (e *downloadV2Engine) FinalSize() int64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.totalSize
}

func (e *downloadV2Engine) part(i int) (downloadV2Part, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if i < len(e.parts) {
		return e.parts[i], true
	}
	return downloadV2Part{}, false
}

func (e *downloadV2Engine) waitForPlan(ctx context.Context) bool {
	select {
	case <-e.planReady:
		return true
	case <-ctx.Done():
		return false
	}
}

// acceptRange validates a range response against the plan and returns the
// number of bytes the part must copy. The first accepted response establishes
// the transfer size, because the metadata the plan came from can be stale,
// re-plans the remaining ranges from it and preallocates the output; later
// ranges are only admitted once planReady closes, so nothing is written before
// that. Every later response must report the same size and cover exactly the
// requested range.
func (e *downloadV2Engine) acceptRange(part downloadV2Part, response downloadV2Range) (int64, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	conflict := func(format string, args ...any) (int64, error) {
		return 0, downloadV2TerminalError{fmt.Errorf("%w: part %v: %s", errDownloadSizeConflict, part.number, fmt.Sprintf(format, args...))}
	}
	if response.trust == UntrustedSizeValue {
		if !e.transferSizeKnown {
			return 0, downloadV2TerminalError{errDownloadV2UntrustedTransferSize}
		}
		return conflict("range response did not report the transfer size %v", e.totalSize)
	}
	if e.transferSizeKnown && response.total != e.totalSize {
		return conflict("response reported total %v, transfer size is %v", response.total, e.totalSize)
	}
	expectedEnd := min(part.off+part.len-1, response.total-1)
	if response.total <= part.off || response.off != part.off || response.end != expectedEnd {
		return conflict("requested bytes %v-%v of %v, response covered bytes %v-%v", part.off, expectedEnd, response.total, response.off, response.end)
	}
	if !e.transferSizeKnown {
		// planReady admits further ranges only once the output is preallocated.
		// On failure the scheduler stays parked and the terminal result cancels
		// the attempt, which releases waitForPlan.
		if err := downloadV2PreallocateFile(e.file, response.total); err != nil {
			return 0, downloadV2TerminalError{err}
		}
		e.totalSize = response.total
		e.transferSizeKnown = true
		e.parts = downloadV2BuildParts(e.startOffset, e.totalSize, e.partSize)
		close(e.planReady)
	}
	return response.end - response.off + 1, nil
}

func (e *downloadV2Engine) ContiguousSize() int64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.contiguous
}

func (e *downloadV2Engine) markComplete(part downloadV2Part, bytes int64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.completed[part.off] = bytes
	for {
		bytes, ok := e.completed[e.contiguous]
		if !ok {
			return
		}
		delete(e.completed, e.contiguous)
		e.contiguous += bytes
	}
}

func (e *downloadV2Engine) directTransferDownloadAttemptAllowed() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return !e.directTransferDisabled
}

func (e *downloadV2Engine) disableDirectTransferDownload(reason string, err error) {
	e.mu.Lock()
	if e.directTransferDisabled {
		e.mu.Unlock()
		return
	}
	e.directTransferDisabled = true
	e.mu.Unlock()

	if e.reportStatus == nil || e.reportStatus.Job() == nil {
		return
	}
	attrs := map[string]interface{}{
		"message":   "direct download disabled; using proxy URL for remaining Download V2 ranges",
		"direction": "download",
		"reason":    reason,
	}
	if err != nil {
		attrs["error"] = uploadRetryLogError(err)
	}
	e.reportStatus.Job().Config.LogPath(e.reportStatus.RemotePath(), attrs)
}

func (e *downloadV2Engine) downloadPartWithRetry(ctx context.Context, part downloadV2Part) downloadV2PartResult {
	var result downloadV2PartResult
	for attempt := 1; attempt <= downloadV2RetryAttempts; attempt++ {
		result = e.downloadPart(ctx, part)
		if result.err == nil || ctx.Err() != nil || downloadV2TerminalPartError(result.err) {
			return result
		}
		e.reportStatus.Job().Config.LogPath(e.reportStatus.RemotePath(), map[string]interface{}{
			"message":     "download v2 part retry",
			"part_number": part.number,
			"part_offset": part.off,
			"part_size":   part.len,
			"attempt":     attempt,
			"error":       result.err,
		})
	}
	return result
}

func (e *downloadV2Engine) downloadPart(ctx context.Context, part downloadV2Part) (result downloadV2PartResult) {
	start := time.Now()
	result = downloadV2PartResult{part: part}
	backpressure := &directTransferBackpressure{}
	ctx = context.WithValue(ctx, directTransferBackpressureContextKey{}, backpressure)
	defer func() {
		if backpressure.seen {
			result.statusCode = http.StatusTooManyRequests
			result.backPressure = true
			result.retryAfter = backpressure.retryAfter
		}
	}()
	reader, response, err := downloadV2ReaderRange(ctx, e.ranger, part.off, part.off+part.len-1)
	if err != nil {
		result.err = err
		result.duration = time.Since(start)
		result.statusCode = downloadV2StatusCode(err)
		result.backPressure = downloadV2BackPressureStatus(result.statusCode)
		return result
	}
	expected, err := e.acceptRange(part, response)
	if err != nil {
		_ = reader.Close()
		result.err = err
		result.duration = time.Since(start)
		return result
	}

	progress := newDownloadV2ProgressBatcher(func(delta int64) {
		e.reportStatus.Job().UpdateStatusWithBytes(status.Downloading, e.reportStatus, delta)
	})
	written, copyErr := downloadV2CopyAt(e.file, part.off, expected, reader, progress.Add)
	closeErr := reader.Close()
	progress.Flush()
	result.bytes = written
	result.duration = time.Since(start)
	if copyErr != nil {
		progress.Subtract(written)
		result.err = copyErr
		return result
	}
	if closeErr != nil {
		progress.Subtract(written)
		result.err = closeErr
		return result
	}
	if written != expected {
		progress.Subtract(written)
		result.err = fmt.Errorf("download v2 part size mismatch for part %v. expected: %v, actual: %v", part.number, expected, written)
		return result
	}
	return result
}

func (r downloadV2PartResult) sample() lib.AdaptiveConcurrencySample {
	return lib.AdaptiveConcurrencySample{
		Success:      r.err == nil,
		Duration:     r.duration,
		Bytes:        r.bytes,
		StatusCode:   r.statusCode,
		BackPressure: r.backPressure,
		RetryAfter:   r.retryAfter,
	}
}

func downloadV2CopyAt(dst io.WriterAt, writeOff int64, expected int64, src io.Reader, progress func(int64)) (written int64, err error) {
	bufPtr := downloadV2CopyBufferPool.Get().(*[]byte)
	buf := *bufPtr
	defer downloadV2CopyBufferPool.Put(bufPtr)
	if expected < 0 {
		return 0, errors.New("negative expected download size")
	}
	for written < expected {
		readSize := len(buf)
		if remaining := expected - written; remaining < int64(readSize) {
			readSize = int(remaining)
		}
		nr, er := io.ReadFull(src, buf[:readSize])
		if nr > 0 {
			nw, ew := dst.WriteAt(buf[:nr], writeOff+written)
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("invalid write result")
				}
			}
			written += int64(nw)
			if progress != nil && nw > 0 {
				progress(int64(nw))
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

func downloadV2BuildParts(startOffset int64, totalSize int64, partSize int64) []downloadV2Part {
	if partSize <= 0 || startOffset >= totalSize {
		return nil
	}
	parts := make([]downloadV2Part, 0, int(ceilDiv(totalSize-startOffset, partSize)))
	for off, partNumber := startOffset, 1; off < totalSize; partNumber++ {
		length := partSize
		if remaining := totalSize - off; remaining < length {
			length = remaining
		}
		parts = append(parts, downloadV2Part{number: partNumber, off: off, len: length})
		off += length
	}
	return parts
}

func downloadV2KnownSizePartSize(target TransferV2TargetClass, totalSize int64) int64 {
	switch target {
	case downloadV2TargetS3:
		return s3KnownSizePreferredPartSize(totalSize)
	default:
		return defaultKnownSizePreferredPartSize(totalSize)
	}
}

func downloadV2SmallFileFallbackSize() int64 {
	return min(
		s3KnownSizePreferredPartSize(0),
		defaultKnownSizePreferredPartSize(0),
	)
}

func downloadV2MaxConcurrency(job *Job, params DownloaderParams, target TransferV2TargetClass) int {
	if params.Manager != nil && !params.AdaptiveConcurrencyUseSDKDefaultCaps {
		return max(1, job.Manager.FilePartsManager.Max())
	}
	maxConcurrency := adaptiveTransferMaxConcurrency(target)
	return manager.EffectiveAdaptiveDownloadV2ConcurrentFileParts(maxConcurrency)
}

func downloadV2AdaptiveConcurrencyConfig(target TransferV2TargetClass, maxConcurrency int, totalSize int64, partSize int64, tuning UploadV2Tuning) lib.AdaptiveConcurrencyConfig {
	maxConcurrency = max(1, maxConcurrency)
	switch target {
	case downloadV2TargetS3, downloadV2TargetDirect:
		plan := uploadV2PartPlan{target: target, totalSize: &totalSize, partSize: partSize, mode: "download_v2_known_size"}
		initial := uploadV2InitialConcurrencyForPlan(plan, maxConcurrency, tuning)
		return uploadV2AdaptiveConcurrencyConfigWithInitial(plan, maxConcurrency, initial, tuning)
	default:
		initial := min(AdaptiveDownloadDefaultTargetInitialTarget, maxConcurrency)
		if tuning.InitialTarget > 0 {
			initial = min(tuning.InitialTarget, maxConcurrency)
		}
		config := lib.AdaptiveConcurrencyConfig{
			MaxConcurrency:            maxConcurrency,
			InitialTarget:             initial,
			MinTarget:                 AdaptiveDownloadDefaultTargetMinTarget,
			GrowEvery:                 AdaptiveTransferDefaultTargetGrowEvery,
			GrowStep:                  AdaptiveTransferDefaultTargetGrowStep,
			FailureShrinkPercent:      AdaptiveTransferDefaultTargetFailureShrinkPercent,
			BackPressureShrinkPercent: AdaptiveDownloadDefaultTargetBackPressureShrinkPercent,
			BackPressurePause:         AdaptiveDownloadDefaultTargetBackPressurePause,
		}
		return config
	}
}

func (s *downloadV2SharedAdaptiveManagerRegistry) get(target TransferV2TargetClass, maxConcurrency int, totalSize int64, partSize int64, tuning UploadV2Tuning) *lib.AdaptiveConcurrencyManager {
	tuning = tuning.managerTuning()
	key := downloadV2AdaptiveManagerCacheKey{target: target, maxConcurrency: maxConcurrency, tuning: tuning}
	return s.managerFor(key, func() *lib.AdaptiveConcurrencyManager {
		return lib.NewAdaptiveConcurrencyManagerWithConfig(downloadV2SharedAdaptiveConcurrencyConfig(target, maxConcurrency, totalSize, partSize, tuning))
	})
}

type downloadV2JobAdmissionTargets struct {
	job *Job
}

func (t downloadV2JobAdmissionTargets) admissionTarget() (int, bool) {
	if t.job == nil {
		return 0, false
	}
	return t.job.downloadV2AdmissionTarget()
}

func (r *Job) downloadV2AdmissionTarget() (target int, ok bool) {
	r.adaptiveDownloadV2Mu.Lock()
	defer r.adaptiveDownloadV2Mu.Unlock()
	for _, manager := range r.adaptiveDownloadV2Managers {
		t := manager.Target()
		if !ok || t < target {
			target = t
		}
		ok = true
	}
	return target, ok
}

func (r *Job) downloadV2AdaptiveManager(target TransferV2TargetClass, maxConcurrency int, totalSize int64, partSize int64, tuning UploadV2Tuning) *lib.AdaptiveConcurrencyManager {
	tuning = tuning.managerTuning()
	key := downloadV2AdaptiveManagerCacheKey{target: target, maxConcurrency: maxConcurrency, tuning: tuning}
	manager := downloadV2SharedAdaptiveManagers.get(target, maxConcurrency, totalSize, partSize, tuning)
	r.adaptiveDownloadV2Mu.Lock()
	if r.adaptiveDownloadV2Managers == nil {
		r.adaptiveDownloadV2Managers = make(map[downloadV2AdaptiveManagerCacheKey]*lib.AdaptiveConcurrencyManager)
	}
	r.adaptiveDownloadV2Managers[key] = manager
	r.adaptiveDownloadV2Mu.Unlock()
	return manager
}

func downloadV2SharedAdaptiveConcurrencyConfig(target TransferV2TargetClass, maxConcurrency int, totalSize int64, partSize int64, tuning UploadV2Tuning) lib.AdaptiveConcurrencyConfig {
	if target == downloadV2TargetS3 || target == downloadV2TargetDirect {
		plan := uploadV2PartPlan{target: target, totalSize: &totalSize, partSize: partSize, mode: "download_v2_known_size"}
		return uploadV2SharedAdaptiveConcurrencyConfig(plan, maxConcurrency, tuning)
	}
	return downloadV2AdaptiveConcurrencyConfig(target, maxConcurrency, totalSize, partSize, tuning)
}

func classifyDownloadV2Target(ctx context.Context, _ goFs.FileInfo, ranger ReaderRange, classifier DownloadV2TargetClassifier) (TransferV2TargetClass, bool) {
	if classifier == nil && downloadV2DirectTarget(ranger) {
		return downloadV2TargetDirect, true
	}
	provider, ok := ranger.(downloadV2URIProvider)
	if !ok {
		return "", false
	}
	downloadURI, err := provider.downloadV2URI(ctx)
	if err != nil {
		return "", false
	}
	return classifyDownloadV2URI(downloadURI, classifier)
}

func downloadV2DirectTarget(ranger ReaderRange) bool {
	file, ok := ranger.(*File)
	if !ok {
		return false
	}
	file.fileMutex.Lock()
	info := file.File.DirectConnectionInfo
	file.fileMutex.Unlock()
	return files_sdk.DirectConnectionInfoPresent(info)
}

func classifyDownloadV2URI(downloadURI string, classifiers ...DownloadV2TargetClassifier) (TransferV2TargetClass, bool) {
	parsed, err := url.Parse(downloadURI)
	if err != nil {
		return "", false
	}
	if len(classifiers) > 0 && classifiers[0] != nil {
		return normalizeTransferV2TargetClass(classifiers[0](downloadURI)), true
	}
	host := strings.ToLower(parsed.Hostname())
	if isS3UploadHost(host) {
		return downloadV2TargetS3, true
	}
	return downloadV2TargetDefault, true
}

func downloadV2ReaderRange(ctx context.Context, ranger ReaderRange, off int64, end int64) (io.ReadCloser, downloadV2Range, error) {
	if withContext, ok := ranger.(lib.FileWithContext); ok {
		ranger = withContext.WithContext(ctx).(ReaderRange)
	}
	if file, ok := ranger.(*File); ok {
		return file.downloadV2ReaderRange(ctx, off, end, false)
	}
	reader, err := ranger.ReaderRange(off, end)
	if err != nil {
		return nil, downloadV2Range{}, err
	}
	return reader, downloadV2RangeFromStat(ranger, off, end), nil
}

// downloadV2RangeFromStat describes a range served by a generic ReaderRange,
// whose Stat reports the size it serves (the pattern DownloadParts relies on).
func downloadV2RangeFromStat(ranger ReaderRange, off int64, end int64) downloadV2Range {
	response := downloadV2Range{off: off, end: end, trust: TrustedSizeValue}
	info, err := ranger.Stat()
	if err != nil || info == nil {
		response.trust = UntrustedSizeValue
		return response
	}
	response.total = info.Size()
	if untrusted, ok := info.(UntrustedSize); ok && untrusted.SizeTrust() == UntrustedSizeValue {
		response.trust = UntrustedSizeValue
	}
	response.end = min(end, response.total-1)
	return response
}

// downloadV2RangeFromResponse describes the range a 206 response covers. Only a
// "bytes start-end/total" Content-Range states which bytes the body carries and
// the transfer size; anything else leaves the response untrusted for the plan.
func downloadV2RangeFromResponse(response *http.Response, off int64, end int64) downloadV2Range {
	served := downloadV2Range{off: off, end: end, trust: UntrustedSizeValue}
	contentRange, ok := parseContentRange(response.Header.Get("Content-Range"))
	if !ok || !contentRange.totalKnown || contentRange.start < 0 {
		return served
	}
	served.off, served.end = contentRange.start, contentRange.end
	served.total, served.trust = contentRange.total, TrustedSizeValue
	return served
}

func (f *File) downloadV2EnsureURI(ctx context.Context) error {
	f.fileMutex.Lock()
	if f.File.DownloadUri != "" {
		f.fileMutex.Unlock()
		return nil
	}
	f.fileMutex.Unlock()

	_, err, _ := downloadV2URIRefreshGroup.Do(f.downloadV2URIRefreshKey(), func() (any, error) {
		f.fileMutex.Lock()
		if f.File.DownloadUri != "" {
			f.fileMutex.Unlock()
			return nil, nil
		}
		current := *f.File
		f.fileMutex.Unlock()

		fileInfo, err := (&Client{Config: f.Config}).DownloadUri(files_sdk.FileDownloadParams{File: current}, files_sdk.WithContext(ctx))
		if err != nil {
			return nil, err
		}
		f.fileMutex.Lock()
		*f.File = fileInfo
		f.fileMutex.Unlock()
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (f *File) downloadV2URIRefreshKey() string {
	return fmt.Sprintf("%p", f.File)
}

func (f *File) downloadV2URI(ctx context.Context) (string, error) {
	if err := f.downloadV2EnsureURI(ctx); err != nil {
		return "", err
	}
	f.fileMutex.Lock()
	downloadURI := f.File.DownloadUri
	f.fileMutex.Unlock()
	return downloadURI, nil
}

func (f *File) downloadV2ReaderRange(ctx context.Context, off int64, end int64, refreshed bool) (io.ReadCloser, downloadV2Range, error) {
	downloadURI, err := f.downloadV2URI(ctx)
	if err != nil {
		return nil, downloadV2Range{}, err
	}
	f.fileMutex.Lock()
	fileCopy := *f.File
	f.fileMutex.Unlock()

	var body io.ReadCloser
	var served downloadV2Range
	headers := &http.Header{}
	headers.Set("Range", fmt.Sprintf("bytes=%v-%v", off, end))
	_, err = (&Client{Config: f.Config}).Download(
		files_sdk.FileDownloadParams{File: fileCopy},
		files_sdk.WithContext(ctx),
		files_sdk.RequestHeadersOption(headers),
		files_sdk.ResponseOption(func(response *http.Response) error {
			if err := lib.ResponseErrors(response, lib.IsStatus(http.StatusForbidden), files_sdk.APIError(), lib.NotStatus(http.StatusPartialContent)); err != nil {
				return &goFs.PathError{Path: f.File.Path, Err: err, Op: "downloadV2ReaderRange"}
			}
			// Only a successful range response may update the size state; an
			// error body's Content-Length is not a file size.
			served = downloadV2RangeFromResponse(response, off, end)
			maxConnections := parseMaxConnections(response)
			downloadRequestID := response.Header.Get("X-Files-Download-Request-Id")
			f.fileMutex.Lock()
			f.MaxConnections = maxConnections
			f.downloadRequestId = downloadRequestID
			f.SizeTrust = served.trust
			if served.trust == TrustedSizeValue {
				f.Size = served.total
			}
			f.fileMutex.Unlock()
			body = response.Body
			return nil
		}),
	)
	if downloadRequestExpired(err) && !refreshed {
		f.Config.LogPath(f.File.Path, map[string]interface{}{"message": "downloadV2DownloadRequestExpired", "error": err})
		f.fileMutex.Lock()
		if f.File.DownloadUri == downloadURI {
			f.File.DownloadUri = ""
		}
		f.fileMutex.Unlock()
		return f.downloadV2ReaderRange(ctx, off, end, true)
	}
	if err != nil {
		return nil, downloadV2Range{}, err
	}
	if body == nil {
		return nil, downloadV2Range{}, &goFs.PathError{Path: f.File.Path, Err: errors.New("missing download response body"), Op: "downloadV2ReaderRange"}
	}
	return body, served, nil
}

func downloadV2StatusCode(err error) int {
	var responseErr lib.ResponseError
	if errors.As(err, &responseErr) {
		return responseErr.StatusCode
	}
	return 0
}

func downloadV2BackPressureStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode == http.StatusServiceUnavailable || statusCode == http.StatusGatewayTimeout
}

func (e *downloadV2Engine) logStart() {
	snapshot := e.manager.Snapshot()
	e.reportStatus.Job().Config.LogPath(e.reportStatus.RemotePath(), map[string]interface{}{
		"message":                    "download v2 start",
		"download_v2_enabled":        true,
		"download_v2_target":         e.target,
		"download_v2_output_mode":    downloadV2OutputPreallocAt,
		"download_v2_total_size":     e.totalSize,
		"download_v2_start_offset":   e.startOffset,
		"download_v2_part_size":      e.partSize,
		"download_v2_part_count":     len(e.parts),
		"download_v2_adaptive_max":   snapshot.Max,
		"download_v2_adaptive_start": snapshot.Target,
	})
}

func (e *downloadV2Engine) logFinish() {
	snapshot := e.manager.Snapshot()
	e.reportStatus.Job().Config.LogPath(e.reportStatus.RemotePath(), map[string]interface{}{
		"message":                                                            "download v2 finish",
		"download_v2_enabled":                                                true,
		"download_v2_target":                                                 e.target,
		"download_v2_output_mode":                                            downloadV2OutputPreallocAt,
		"download_v2_adaptive_target":                                        snapshot.Target,
		"download_v2_adaptive_peak_target":                                   snapshot.PeakTarget,
		"download_v2_adaptive_peak_running":                                  snapshot.PeakRunning,
		"download_v2_adaptive_success_total":                                 snapshot.SuccessTotal,
		"download_v2_adaptive_failure_total":                                 snapshot.FailureTotal,
		"download_v2_adaptive_grow_total":                                    snapshot.GrowTotal,
		"download_v2_adaptive_shrink_total":                                  snapshot.ShrinkTotal,
		"download_v2_adaptive_growth_ceiling":                                snapshot.GrowthCeiling,
		"download_v2_adaptive_growth_unlocked":                               snapshot.GrowthCeilingUnlocked,
		"download_v2_adaptive_back_pressure_total":                           snapshot.BackPressureTotal,
		"download_v2_adaptive_retry_after_total":                             snapshot.RetryAfterTotal,
		"download_v2_adaptive_throughput_backoff_total":                      snapshot.ThroughputBackoffTotal,
		"download_v2_adaptive_throughput_probe_miss_total":                   snapshot.ThroughputProbeMissTotal,
		"download_v2_adaptive_throughput_probe_efficiency_miss_total":        snapshot.ThroughputProbeEfficiencyMissTotal,
		"download_v2_adaptive_latency_backoff_total":                         snapshot.LatencyBackoffTotal,
		"download_v2_adaptive_latency_growth_suppression_total":              snapshot.LatencyGrowthSuppressionTotal,
		"download_v2_adaptive_bytes_total":                                   snapshot.BytesTotal,
		"download_v2_adaptive_average_duration_ms":                           snapshot.AverageDuration.Milliseconds(),
		"download_v2_adaptive_last_throughput_bps":                           snapshot.LastThroughputBytesPerSecond,
		"download_v2_adaptive_best_throughput_bps":                           snapshot.BestThroughputBytesPerSecond,
		"download_v2_adaptive_last_throughput_probe_gain_percent":            snapshot.LastThroughputProbeGainPercent,
		"download_v2_adaptive_last_throughput_probe_target_delta":            snapshot.LastThroughputProbeTargetDelta,
		"download_v2_adaptive_last_throughput_probe_gain_per_target_percent": snapshot.LastThroughputProbeGainPerTargetPercent,
		"download_v2_adaptive_last_queue_estimate":                           snapshot.LastQueueEstimate,
		"download_v2_adaptive_min_duration_per_byte":                         snapshot.MinDurationPerByte,
		"download_v2_adaptive_last_duration_per_byte":                        snapshot.LastDurationPerByte,
		"download_v2_contiguous_size":                                        e.ContiguousSize(),
		"download_v2_transfer_size":                                          e.FinalSize(),
	})
}

type downloadV2ProgressBatcher struct {
	progress func(int64)
	pending  int64
}

func newDownloadV2ProgressBatcher(progress func(int64)) *downloadV2ProgressBatcher {
	if progress == nil {
		progress = func(int64) {}
	}
	return &downloadV2ProgressBatcher{progress: progress}
}

func (b *downloadV2ProgressBatcher) Add(delta int64) {
	if delta == 0 {
		return
	}
	b.pending += delta
	if b.pending >= downloadV2ProgressBatch || b.pending <= -downloadV2ProgressBatch || delta < 0 {
		b.Flush()
	}
}

func (b *downloadV2ProgressBatcher) Subtract(delta int64) {
	if delta > 0 {
		b.Add(-delta)
	}
}

func (b *downloadV2ProgressBatcher) Flush() {
	if b.pending == 0 {
		return
	}
	b.progress(b.pending)
	b.pending = 0
}
