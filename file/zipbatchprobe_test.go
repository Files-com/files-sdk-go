package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZipBatchProbeMedianWindow(t *testing.T) {
	got := zipBatchProbeMedianWindow([]zipBatchProbeRateWindow{
		{rate: 100, completions: 10, duration: 100 * time.Millisecond},
		{rate: 10, completions: 1, duration: 100 * time.Millisecond},
		{rate: 100, completions: 10, duration: 100 * time.Millisecond},
	})

	assert.Equal(t, 100.0, got.rate)
	assert.Equal(t, 21, got.completions)
	assert.Equal(t, 300*time.Millisecond, got.duration)
}

func TestZipBatchDownloaderDefaultsEngageEndToEnd(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(defaultZipBatchMinFiles), nil)
	server.DownloadDelay = func(string) time.Duration { return 2 * time.Millisecond }
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	assert.NotEmpty(t, server.ZipCreateRequests)
	assert.Greater(t, job.ZipBatchStats().BatchesDispatched, int64(0))
}

func TestZipBatchProbeChoosesZip(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(512), nil)
	server.DownloadDelay = func(path string) time.Duration {
		var idx int
		_, _ = fmt.Sscanf(filepath.Base(path), "file-%03d.txt", &idx)
		return 200*time.Millisecond + time.Duration(idx%64)*10*time.Millisecond
	}
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.New(64, 1, 1),
		RetryPolicy: RetryPolicy{},
		ZipBatch: ZipBatchParams{
			MinFiles:          32,
			MaxFiles:          16,
			BatchSize:         16,
			ConcurrentBatches: 8,
		},
	})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, string(zipBatchProbeDecisionCommitted), stats.ProbeDecision)
	assert.Equal(t, int64(1), stats.ProbeBaselineWindows)
	assert.Greater(t, stats.ProbeZipRateMilli, stats.ProbePerFileRateMilli)
	assert.Equal(t, int64(0), stats.Reprobes)
	assert.Greater(t, stats.BatchesDispatched, int64(1))
	assert.Equal(t, "file-511", readLocalFile(t, root, "batch", "file-511.txt"))
}

func TestZipBatchProbeLargeJobUsesThreeBaselineWindows(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	zipBatchTuning.ProbePhaseAMinCompletions = 5
	zipBatchTuning.ProbePhaseAMinDuration = 50 * time.Millisecond
	zipBatchTuning.ProbePhaseAWindow = 50 * time.Millisecond
	zipBatchTuning.ProbePhaseAHardCap = time.Second
	zipBatchTuning.ProbePhaseAWindow2After = 4
	zipBatchTuning.ProbePhaseAWindow3After = 12
	zipBatchTuning.ProbePhaseBMinCompletions = 4
	zipBatchTuning.ProbePhaseBWarmup = 10 * time.Millisecond
	zipBatchTuning.ProbePhaseBWindow = 50 * time.Millisecond
	zipBatchTuning.ProbePhaseBEndMinWindow = 10 * time.Millisecond
	zipBatchTuning.ProbePhaseBHardCap = time.Second
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(120), nil)
	var directRequests atomic.Int64
	server.DownloadDelay = func(string) time.Duration {
		n := directRequests.Add(1)
		if n >= 11 && n <= 15 {
			return 50 * time.Millisecond
		}
		return 5 * time.Millisecond
	}
	server.ZipDownloadDelay = func(int, []string) time.Duration { return 100 * time.Millisecond }
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.New(1, 1, 1),
		RetryPolicy: RetryPolicy{},
		ZipBatch: ZipBatchParams{
			MinFiles:          16,
			MaxFiles:          1,
			BatchSize:         1,
			ConcurrentBatches: 1,
			ReprobeInterval:   -1,
		},
	})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, int64(3), stats.ProbeBaselineWindows)
	assert.Equal(t, string(zipBatchProbeDecisionDissolved), stats.ProbeDecision)
	assert.Greater(t, stats.ProbePerFileRateMilli, stats.ProbeZipRateMilli)
	assert.Equal(t, "file-119", readLocalFile(t, root, "batch", "file-119.txt"))
}

func TestZipBatchProbeChoosesPerFile(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(512), nil)
	server.DownloadDelay = func(string) time.Duration { return time.Millisecond }
	server.ZipDownloadChunkDelay = func(int, []string) time.Duration { return 30 * time.Millisecond }
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.New(16, 1, 1),
		RetryPolicy: RetryPolicy{},
		ZipBatch: ZipBatchParams{
			MinFiles:          32,
			MaxFiles:          8,
			BatchSize:         8,
			ConcurrentBatches: 8,
		},
	})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, string(zipBatchProbeDecisionDissolved), stats.ProbeDecision)
	assert.Greater(t, stats.ProbePerFileRateMilli, stats.ProbeZipRateMilli)
	assert.Greater(t, stats.ProbeZipFiles, int64(0))
	assert.Equal(t, "file-511", readLocalFile(t, root, "batch", "file-511.txt"))
}

func TestZipBatchReprobeCommitsAfterDissolve(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	zipBatchTuning.ProbePhaseBMinCompletions = 8
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(300), nil)
	server.DownloadDelay = func(string) time.Duration { return 10 * time.Millisecond }
	server.ZipDownloadChunkDelay = func(attempt int, _ []string) time.Duration {
		if attempt == 1 {
			return 3 * time.Millisecond
		}
		return 0
	}
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.New(4, 1, 1),
		RetryPolicy: RetryPolicy{},
		ZipBatch: ZipBatchParams{
			MinFiles:          16,
			MaxFiles:          8,
			BatchSize:         8,
			ConcurrentBatches: 1,
		},
	})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, string(zipBatchProbeDecisionCommitted), stats.ProbeDecision)
	assert.GreaterOrEqual(t, stats.Reprobes, int64(1))
	assert.Greater(t, len(server.ZipCreateRequests), 1)
	assert.Equal(t, "file-299", readLocalFile(t, root, "batch", "file-299.txt"))
}

func TestZipBatchReprobeBacksOffRepeatedDissolves(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	zipBatchTuning.ProbePhaseBMinCompletions = 4
	zipBatchTuning.DefaultReprobeInterval = 10 * time.Millisecond
	zipBatchTuning.MaxReprobeInterval = 40 * time.Millisecond
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(900), nil)
	server.DownloadDelay = func(string) time.Duration { return 2 * time.Millisecond }
	server.ZipDownloadChunkDelay = func(int, []string) time.Duration { return 4 * time.Millisecond }
	var zipStarts []time.Time
	var zipStartsMu sync.Mutex
	server.ZipDownloadBefore = func(int, []string) {
		zipStartsMu.Lock()
		zipStarts = append(zipStarts, time.Now())
		zipStartsMu.Unlock()
	}
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.New(2, 1, 1),
		RetryPolicy: RetryPolicy{},
		ZipBatch: ZipBatchParams{
			MinFiles:          16,
			MaxFiles:          4,
			BatchSize:         4,
			ConcurrentBatches: 1,
		},
	})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, string(zipBatchProbeDecisionDissolved), stats.ProbeDecision)
	assert.GreaterOrEqual(t, stats.Reprobes, int64(2))
	zipStartsMu.Lock()
	starts := append([]time.Time(nil), zipStarts...)
	zipStartsMu.Unlock()
	require.GreaterOrEqual(t, len(starts), 3)
	assert.GreaterOrEqual(t, starts[2].Sub(starts[1]), starts[1].Sub(starts[0]))
}

func TestZipBatchReprobeNegativeIntervalDisablesReprobe(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(180), nil)
	server.DownloadDelay = func(string) time.Duration { return time.Millisecond }
	server.ZipDownloadChunkDelay = func(int, []string) time.Duration { return 3 * time.Millisecond }
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.New(4, 1, 1),
		RetryPolicy: RetryPolicy{},
		ZipBatch: ZipBatchParams{
			MinFiles:          16,
			MaxFiles:          8,
			BatchSize:         8,
			ConcurrentBatches: 1,
			ReprobeInterval:   -1,
		},
	})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, string(zipBatchProbeDecisionDissolved), stats.ProbeDecision)
	assert.Equal(t, int64(0), stats.Reprobes)
}

func TestZipBatchProbeUsesSequentialPhases(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(180), nil)
	server.DownloadDelay = func(string) time.Duration { return 5 * time.Millisecond }
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.New(16, 1, 1),
		RetryPolicy: RetryPolicy{},
		ZipBatch:    ZipBatchParams{MinFiles: 32, MaxFiles: 16, BatchSize: 16, ConcurrentBatches: 4},
	})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.GreaterOrEqual(t, stats.ProbePerFileFiles, int64(zipBatchTuning.ProbePhaseAMinCompletions))
	assert.Greater(t, stats.ProbeZipFiles, int64(0))
	require.NotEmpty(t, server.ZipCreateRequests)
	assert.NotContains(t, flattenZipBatchRequests(server.ZipCreateRequests), "batch/file-000.txt")
	downloads := zipBatchDownloadRequests(server)
	for _, path := range []string{"batch/file-000.txt", "batch/file-001.txt", "batch/file-002.txt"} {
		assert.Contains(t, downloads, path)
	}
}

func TestZipBatchProbeEndOfWalkDissolvesMidProbe(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(20), nil)
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{MinFiles: 16, MaxFiles: 16, BatchSize: 16}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, string(zipBatchProbeDecisionDissolved), stats.ProbeDecision)
	assert.Empty(t, server.ZipCreateRequests)
}

func TestZipBatchProbePhaseBInsufficientCompletionsDissolves(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	zipBatchTuning.ProbePhaseBMinCompletions = 64
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(50), nil)
	server.DownloadDelay = func(string) time.Duration { return 5 * time.Millisecond }
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.New(4, 1, 1),
		RetryPolicy: RetryPolicy{},
		ZipBatch:    ZipBatchParams{MinFiles: 32, MaxFiles: 16, BatchSize: 16, ConcurrentBatches: 4},
	})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, string(zipBatchProbeDecisionDissolved), stats.ProbeDecision)
	assert.Greater(t, stats.ProbeZipFiles, int64(0))
	assert.Less(t, stats.ProbeZipFiles, int64(zipBatchTuning.ProbePhaseBMinCompletions))
}

func TestZipBatchProbeRampSkipIgnoresSlowFirstZipWave(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	zipBatchTuning.ProbePhaseBWarmup = 20 * time.Millisecond
	zipBatchTuning.ProbePhaseBWindow = 30 * time.Millisecond
	zipBatchTuning.ProbePhaseBHardCap = 500 * time.Millisecond
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(256), nil)
	server.DownloadDelay = func(string) time.Duration { return 50 * time.Millisecond }
	server.ZipDownloadDelay = func(attempt int, _ []string) time.Duration {
		if attempt <= 2 {
			return 300 * time.Millisecond
		}
		return 0
	}
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.New(16, 1, 1),
		RetryPolicy: RetryPolicy{},
		ZipBatch: ZipBatchParams{
			MinFiles:          32,
			MaxFiles:          16,
			BatchSize:         16,
			ConcurrentBatches: 2,
		},
	})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, string(zipBatchProbeDecisionCommitted), stats.ProbeDecision)
	assert.Greater(t, stats.ProbeZipRateMilli, int64(float64(stats.ProbePerFileRateMilli)*defaultZipBatchMinAdvantage))
}

func TestZipBatchProbeTrailingWindowIgnoresEarlyPerFileBurst(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	zipBatchTuning.ProbePhaseAMinDuration = 80 * time.Millisecond
	zipBatchTuning.ProbePhaseAWindow = 30 * time.Millisecond
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(512), nil)
	server.DownloadDelay = func(path string) time.Duration {
		var idx int
		_, _ = fmt.Sscanf(filepath.Base(path), "file-%03d.txt", &idx)
		if idx < 16 {
			return 0
		}
		return 20 * time.Millisecond
	}
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.New(16, 1, 1),
		RetryPolicy: RetryPolicy{},
		ZipBatch: ZipBatchParams{
			MinFiles:          32,
			MaxFiles:          16,
			BatchSize:         16,
			ConcurrentBatches: 4,
		},
	})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, string(zipBatchProbeDecisionCommitted), stats.ProbeDecision)
	assert.Less(t, stats.ProbePerFileRateMilli, int64(2_000_000))
}

func TestZipBatchMinAdvantageNegativeDisablesProbe(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(32), nil)
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{MinFiles: 16, MaxFiles: 16, BatchSize: 16, MinAdvantage: -1}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, string(zipBatchProbeDecisionNone), stats.ProbeDecision)
	assert.Equal(t, int64(0), stats.ProbeZipFiles)
	assert.Equal(t, int64(0), stats.ProbePerFileFiles)
	assert.Equal(t, int64(0), stats.Reprobes)
	assert.Empty(t, server.DownloadRequests)
	assert.NotEmpty(t, server.ZipCreateRequests)
}

func TestZipBatchCircuitBreakerDissolvesRemainder(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(80), nil)
	server.ZipDownloadMutator = func(_ int, _ []string, body []byte) []byte {
		return corruptAfterFirstZipEntry(body)
	}
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.New(1, 1, 1),
		RetryPolicy: RetryPolicy{},
		ZipBatch:    ZipBatchParams{MinFiles: 16, MaxFiles: 16, BatchSize: 16, ConcurrentBatches: 1, MinAdvantage: -1},
	})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, int64(1), stats.CircuitBreakerTrips)
	assert.Len(t, server.ZipCreateRequests, zipBatchCircuitBreakerBatches)
	assert.NotEmpty(t, server.DownloadRequests)
}
