package file

import (
	"fmt"
	"os"
	"path/filepath"
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
	assert.Greater(t, stats.ProbeZipRateMilli, stats.ProbePerFileRateMilli)
	assert.Greater(t, stats.BatchesDispatched, int64(1))
	assert.Equal(t, "file-511", readLocalFile(t, root, "batch", "file-511.txt"))
}

func TestZipBatchProbeChoosesFasterSample(t *testing.T) {
	for _, tt := range []struct {
		name       string
		zipSpacing time.Duration
		decision   zipBatchProbeDecision
		mode       zipBatchMode
	}{
		{"zip", 100 * time.Microsecond, zipBatchProbeDecisionCommitted, zipBatchModeCommitted},
		{"per-file", 2 * time.Millisecond, zipBatchProbeDecisionDissolved, zipBatchModeDissolved},
	} {
		t.Run(tt.name, func(t *testing.T) {
			withZipBatchProbeTestWindows(t)
			batcher := &zipBatchDownloader{params: ZipBatchParams{}.withDefaults(), probe: newZipBatchProbe()}
			batcher.unlock()

			start := time.Unix(1, 0)
			batcher.probe.perFileFirstCompletion = start
			for elapsed := time.Duration(0); elapsed <= 20*time.Millisecond; elapsed += time.Millisecond {
				batcher.probe.perFileCompletions = append(batcher.probe.perFileCompletions, start.Add(elapsed))
			}
			batcher.maybeAdvanceProbeLocked(false, start.Add(20*time.Millisecond))
			batcher.applyProbeState(false)
			require.Equal(t, zipBatchModeProbePhaseB, batcher.mode)
			require.Equal(t, zipBatchProbeDecisionNone, batcher.currentProbeDecision())

			zipStart := start.Add(time.Second)
			batcher.probe.zipFirstCompletion = zipStart
			for elapsed := time.Duration(0); elapsed <= 40*time.Millisecond; elapsed += tt.zipSpacing {
				batcher.probe.zipCompletions = append(batcher.probe.zipCompletions, zipStart.Add(elapsed))
			}
			batcher.maybeAdvanceProbeLocked(false, zipStart.Add(40*time.Millisecond))
			batcher.applyProbeState(false)

			assert.Equal(t, tt.decision, batcher.currentProbeDecision())
			assert.Equal(t, tt.mode, batcher.mode)
			assert.Equal(t, 1, batcher.currentProbeAttempt())
			assert.Len(t, batcher.probe.perFileRateWindows, 1)
			if tt.decision == zipBatchProbeDecisionCommitted {
				assert.Greater(t, batcher.probe.zipRate, batcher.probe.perFileRate)
				assert.True(t, batcher.reprobeAfter.IsZero())
			} else {
				assert.Greater(t, batcher.probe.perFileRate, batcher.probe.zipRate)
				assert.False(t, batcher.reprobeAfter.IsZero())
			}
		})
	}
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
	// Measure sustained rates instead of a few completions in a 10 ms scheduling burst.
	zipBatchTuning.ProbePhaseAMinDuration = 100 * time.Millisecond
	zipBatchTuning.ProbePhaseAWindow = 100 * time.Millisecond
	zipBatchTuning.ProbePhaseBWindow = 100 * time.Millisecond
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
			// This test covers the initial choice; later attempts have their own tests.
			ReprobeInterval: -1,
		},
	})

	requireZipBatchJobClean(t, job)
	stats := job.ZipBatchStats()
	assert.Equal(t, string(zipBatchProbeDecisionDissolved), stats.ProbeDecision)
	assert.Greater(t, stats.ProbePerFileRateMilli, stats.ProbeZipRateMilli)
	assert.Greater(t, stats.ProbeZipFiles, int64(0))
	assert.Equal(t, int64(0), stats.Reprobes)
	assert.Equal(t, "file-511", readLocalFile(t, root, "batch", "file-511.txt"))
}

func TestZipBatchReprobeCommitsAfterDissolve(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	zipBatchTuning.ProbePhaseBMinCompletions = 8
	// Both probes compare sustained 100ms samples of each route, at rates the
	// server sets. Until the first ZIP download starts, four per-file
	// connections wait 10ms per file, while the first probe's two ZIP
	// streams send 32 bytes every 5ms, a few times slower, so the first
	// probe dissolves. From then on each per-file request waits 150ms, so no
	// 100ms window holds more than four per-file files, while later ZIP
	// streams are not paced; a re-probe commits once its baseline no longer
	// reaches back to the fast per-file phase. This still depends on the
	// machine keeping per-file requests well ahead of the paced ZIP stream
	// at first, and later finishing an unpaced 8-file ZIP batch well within
	// 100ms.
	zipBatchTuning.ProbePhaseAMinDuration = 100 * time.Millisecond
	zipBatchTuning.ProbePhaseAWindow = 100 * time.Millisecond
	zipBatchTuning.ProbePhaseBWindow = 100 * time.Millisecond
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(300), nil)
	var zipStarted atomic.Bool
	server.ZipDownloadBefore = func(int, []string) { zipStarted.Store(true) }
	server.DownloadDelay = func(string) time.Duration {
		if zipStarted.Load() {
			return 150 * time.Millisecond
		}
		return 10 * time.Millisecond
	}
	server.ZipDownloadChunkDelay = func(attempt int, _ []string) time.Duration {
		if attempt <= 2 {
			return 5 * time.Millisecond
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
}

// Each dissolved verdict schedules the next re-probe: the first after ReprobeInterval, each later one after twice the
// previous wait, up to MaxReprobeInterval. The re-probe starts exactly at that deadline when a per-file completion
// from the cooldown gives it a baseline.
func TestZipBatchReprobeCooldownDoublesToCapAndOpensAtDeadline(t *testing.T) {
	withZipBatchProbeTestWindows(t)
	zipBatchTuning.DefaultReprobeInterval = 10 * time.Millisecond
	zipBatchTuning.MaxReprobeInterval = 40 * time.Millisecond
	batcher := &zipBatchDownloader{params: ZipBatchParams{}.withDefaults(), probe: newZipBatchProbe()}
	batcher.unlock()

	for i, wait := range []time.Duration{10 * time.Millisecond, 20 * time.Millisecond, 40 * time.Millisecond, 40 * time.Millisecond} {
		attempt := i + 1
		batcher.probe.mu.Lock()
		batcher.probe.decision = zipBatchProbeDecisionDissolved
		batcher.probe.mu.Unlock()

		before := time.Now()
		batcher.applyProbeState(false)
		after := time.Now()
		require.Equal(t, zipBatchModeDissolved, batcher.mode)
		deadline := batcher.reprobeAfter
		require.GreaterOrEqual(t, deadline.Sub(before), wait, "attempt %d cooldown", attempt)
		require.LessOrEqual(t, deadline.Sub(after), wait, "attempt %d cooldown", attempt)

		batcher.probe.mu.Lock()
		batcher.probe.perFileCompletions = []time.Time{deadline.Add(-time.Millisecond)}
		batcher.probe.mu.Unlock()

		require.False(t, batcher.startReprobeIfDue(deadline.Add(-time.Nanosecond)), "attempt %d", attempt)
		require.Equal(t, attempt, batcher.currentProbeAttempt())
		require.Equal(t, zipBatchModeDissolved, batcher.mode)

		require.True(t, batcher.startReprobeIfDue(deadline), "attempt %d", attempt)
		require.Equal(t, attempt+1, batcher.currentProbeAttempt())
		require.Equal(t, zipBatchModeProbePhaseB, batcher.mode)
		require.Equal(t, zipBatchProbePhaseB, batcher.currentProbePhase())
		require.True(t, batcher.reprobeAfter.IsZero())
	}
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
