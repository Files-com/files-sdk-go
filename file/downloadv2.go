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
	totalSize    int64
	startOffset  int64
	partSize     int64
	parts        []downloadV2Part

	mu         sync.Mutex
	completed  map[int64]int64
	contiguous int64
}

type downloadV2AdaptiveManagerCacheKey struct {
	target         TransferV2TargetClass
	maxConcurrency int
}

type downloadV2URIProvider interface {
	downloadV2URI(context.Context) (string, error)
}

type downloadV2SharedAdaptiveManagerRegistry struct {
	mu       sync.Mutex
	managers map[downloadV2AdaptiveManagerCacheKey]*lib.AdaptiveConcurrencyManager
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

func runDownloadV2IfSupported(ctx context.Context, reportStatus *DownloadStatus, remoteStat goFs.FileInfo, tmpName string, startOffset int64) (bool, int64, error) {
	params, _ := reportStatus.Job().Params.(DownloaderParams)
	if !params.AdaptiveConcurrency {
		return false, 0, nil
	}
	if reportStatus.Job().Manager.FilePartsManager.DownloadFilesAsSingleStream {
		return false, 0, nil
	}
	ranger, ok := reportStatus.fsFile.(ReaderRange)
	if !ok {
		return false, 0, nil
	}
	totalSize := remoteStat.Size()
	if totalSize < 0 || startOffset < 0 || startOffset > totalSize {
		return false, 0, nil
	}
	if untrusted, ok := remoteStat.(UntrustedSize); ok && untrusted.SizeTrust() == UntrustedSizeValue {
		return false, 0, nil
	}
	if totalSize-startOffset <= downloadV2SmallFileFallbackSize() {
		return false, 0, nil
	}
	target, ok := classifyDownloadV2Target(ctx, remoteStat, ranger, params.AdaptiveDownloadV2TargetClassifier)
	if !ok {
		return false, 0, nil
	}
	partSize := downloadV2KnownSizePartSize(target, totalSize)
	if totalSize-startOffset <= partSize {
		return false, 0, nil
	}

	file, err := os.OpenFile(tmpName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return true, 0, err
	}
	engine := newDownloadV2Engine(reportStatus, ranger, file, target, totalSize, startOffset, partSize, params)
	if err := engine.Run(ctx); err != nil {
		return true, engine.ContiguousSize(), err
	}
	return true, engine.FinalSize(), nil
}

func newDownloadV2Engine(reportStatus *DownloadStatus, ranger ReaderRange, file *os.File, target TransferV2TargetClass, totalSize int64, startOffset int64, partSize int64, params DownloaderParams) *downloadV2Engine {
	maxConcurrency := downloadV2MaxConcurrency(reportStatus.Job(), params)
	manager := lib.NewAdaptiveConcurrencyManagerWithConfig(downloadV2AdaptiveConcurrencyConfig(target, maxConcurrency, totalSize, partSize))
	if params.Manager == nil || params.AdaptiveConcurrencyUseSDKDefaultCaps {
		manager = downloadV2SharedAdaptiveManagers.get(target, maxConcurrency, totalSize, partSize)
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
		parts:        downloadV2BuildParts(startOffset, totalSize, partSize),
		completed:    make(map[int64]int64),
		contiguous:   startOffset,
	}
}

func (e *downloadV2Engine) Run(parentCtx context.Context) (err error) {
	defer func() {
		if err != nil {
			_ = e.file.Truncate(e.ContiguousSize())
		}
		closeErr := e.file.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if err = downloadV2PreallocateFile(e.file, e.totalSize); err != nil {
		return err
	}
	e.logStart()
	if e.totalSize == 0 || e.startOffset == e.totalSize {
		return nil
	}

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	results := make(chan downloadV2PartResult, max(1, e.manager.Max()))
	var wg sync.WaitGroup
	go func() {
		defer func() {
			wg.Wait()
			close(results)
		}()
		started := false
		for _, part := range e.parts {
			if ctx.Err() != nil {
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
		}
	}()

	for result := range results {
		if result.err != nil {
			cancel()
			if err == nil {
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
	if e.ContiguousSize() != e.totalSize {
		return fmt.Errorf("download v2 wrote non-contiguous file. expected: %v, actual: %v", e.totalSize, e.ContiguousSize())
	}
	e.logFinish()
	return nil
}

func (e *downloadV2Engine) FinalSize() int64 {
	return e.totalSize
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

func (e *downloadV2Engine) downloadPartWithRetry(ctx context.Context, part downloadV2Part) downloadV2PartResult {
	var result downloadV2PartResult
	for attempt := 1; attempt <= downloadV2RetryAttempts; attempt++ {
		result = e.downloadPart(ctx, part)
		if result.err == nil {
			return result
		}
		if ctx.Err() != nil {
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

func (e *downloadV2Engine) downloadPart(ctx context.Context, part downloadV2Part) downloadV2PartResult {
	start := time.Now()
	result := downloadV2PartResult{part: part}
	reader, err := downloadV2ReaderRange(ctx, e.ranger, part.off, part.off+part.len-1)
	if err != nil {
		result.err = err
		result.duration = time.Since(start)
		result.statusCode = downloadV2StatusCode(err)
		result.backPressure = downloadV2BackPressureStatus(result.statusCode)
		return result
	}

	progress := newDownloadV2ProgressBatcher(func(delta int64) {
		e.reportStatus.Job().UpdateStatusWithBytes(status.Downloading, e.reportStatus, delta)
	})
	written, copyErr := downloadV2CopyAt(e.file, part.off, part.len, reader, progress.Add)
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
	if written != part.len {
		progress.Subtract(written)
		result.err = fmt.Errorf("download v2 part size mismatch for part %v. expected: %v, actual: %v", part.number, part.len, written)
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

func downloadV2MaxConcurrency(job *Job, params DownloaderParams) int {
	if params.Manager != nil && !params.AdaptiveConcurrencyUseSDKDefaultCaps {
		return max(1, job.Manager.FilePartsManager.Max())
	}
	return manager.EffectiveAdaptiveDownloadV2ConcurrentFileParts()
}

func downloadV2AdaptiveConcurrencyConfig(target TransferV2TargetClass, maxConcurrency int, totalSize int64, partSize int64) lib.AdaptiveConcurrencyConfig {
	maxConcurrency = max(1, maxConcurrency)
	switch target {
	case downloadV2TargetS3:
		plan := uploadV2PartPlan{target: uploadV2TargetS3, totalSize: &totalSize, partSize: partSize, mode: "download_v2_known_size"}
		initial := uploadV2InitialConcurrencyForPlan(plan, maxConcurrency, UploadV2Tuning{})
		return uploadV2AdaptiveConcurrencyConfigWithInitial(plan, maxConcurrency, initial, UploadV2Tuning{})
	default:
		return lib.AdaptiveConcurrencyConfig{
			MaxConcurrency:            maxConcurrency,
			InitialTarget:             min(15, maxConcurrency),
			MinTarget:                 1,
			GrowEvery:                 8,
			GrowStep:                  1,
			FailureShrinkPercent:      50,
			BackPressureShrinkPercent: 35,
			BackPressurePause:         500 * time.Millisecond,
		}
	}
}

func (s *downloadV2SharedAdaptiveManagerRegistry) get(target TransferV2TargetClass, maxConcurrency int, totalSize int64, partSize int64) *lib.AdaptiveConcurrencyManager {
	key := downloadV2AdaptiveManagerCacheKey{target: target, maxConcurrency: maxConcurrency}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.managers == nil {
		s.managers = make(map[downloadV2AdaptiveManagerCacheKey]*lib.AdaptiveConcurrencyManager)
	}
	if manager := s.managers[key]; manager != nil {
		return manager
	}
	manager := lib.NewAdaptiveConcurrencyManagerWithConfig(downloadV2SharedAdaptiveConcurrencyConfig(target, maxConcurrency, totalSize, partSize))
	s.managers[key] = manager
	return manager
}

func downloadV2SharedAdaptiveConcurrencyConfig(target TransferV2TargetClass, maxConcurrency int, totalSize int64, partSize int64) lib.AdaptiveConcurrencyConfig {
	if target == downloadV2TargetS3 {
		plan := uploadV2PartPlan{target: uploadV2TargetS3, totalSize: &totalSize, partSize: partSize, mode: "download_v2_known_size"}
		return uploadV2SharedAdaptiveConcurrencyConfig(plan, maxConcurrency, UploadV2Tuning{})
	}
	return downloadV2AdaptiveConcurrencyConfig(target, maxConcurrency, totalSize, partSize)
}

func classifyDownloadV2Target(ctx context.Context, _ goFs.FileInfo, ranger ReaderRange, classifier DownloadV2TargetClassifier) (TransferV2TargetClass, bool) {
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

func downloadV2ReaderRange(ctx context.Context, ranger ReaderRange, off int64, end int64) (io.ReadCloser, error) {
	if withContext, ok := ranger.(lib.FileWithContext); ok {
		ranger = withContext.WithContext(ctx).(ReaderRange)
	}
	if file, ok := ranger.(*File); ok {
		return file.downloadV2ReaderRange(ctx, off, end, false)
	}
	return ranger.ReaderRange(off, end)
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

func (f *File) downloadV2ReaderRange(ctx context.Context, off int64, end int64, refreshed bool) (io.ReadCloser, error) {
	downloadURI, err := f.downloadV2URI(ctx)
	if err != nil {
		return nil, err
	}
	f.fileMutex.Lock()
	fileCopy := *f.File
	f.fileMutex.Unlock()

	var body io.ReadCloser
	headers := &http.Header{}
	headers.Set("Range", fmt.Sprintf("bytes=%v-%v", off, end))
	_, err = (&Client{Config: f.Config}).Download(
		files_sdk.FileDownloadParams{File: fileCopy},
		files_sdk.WithContext(ctx),
		files_sdk.RequestHeadersOption(headers),
		files_sdk.ResponseOption(func(response *http.Response) error {
			size, trust := parseSize(response)
			maxConnections := parseMaxConnections(response)
			downloadRequestID := response.Header.Get("X-Files-Download-Request-Id")
			f.fileMutex.Lock()
			f.MaxConnections = maxConnections
			f.downloadRequestId = downloadRequestID
			f.Size = size
			f.SizeTrust = trust
			f.fileMutex.Unlock()
			if err := lib.ResponseErrors(response, lib.IsStatus(http.StatusForbidden), files_sdk.APIError(), lib.NotStatus(http.StatusPartialContent)); err != nil {
				return &goFs.PathError{Path: f.File.Path, Err: err, Op: "downloadV2ReaderRange"}
			}
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
		return nil, err
	}
	if body == nil {
		return nil, &goFs.PathError{Path: f.File.Path, Err: errors.New("missing download response body"), Op: "downloadV2ReaderRange"}
	}
	return body, nil
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
