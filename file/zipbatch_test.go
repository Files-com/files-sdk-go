package file

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZipBatchIncrementFilename(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		want    string
		wantErr bool
	}{
		{name: "file.txt", count: 1, want: "file (1).txt"},
		{name: "archive.tar.gz", count: 2, want: "archive.tar (2).gz"},
		{name: "README", count: 1, want: "README (1)"},
		{name: ".env", count: 1, wantErr: true},
		{name: "name.", count: 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := zipBatchIncrementFilename(tt.name, tt.count)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestZipBatchEntryMapDetectsAmbiguity(t *testing.T) {
	_, err := zipBatchEntryMap([]*DownloadStatus{
		{remotePath: "root/a (1).txt"},
		{remotePath: "root/a.txt"},
		{remotePath: "root/z/a.txt"},
	})
	require.Error(t, err)
}

func TestZipBatchParse422Paths(t *testing.T) {
	body := []byte(`{"error":"Paths does not match existing file at: root/missing.txt","errors":[{"error":"Paths cannot be downloaded at: root/locked.txt"}]}`)
	assert.Equal(t, map[string]struct{}{
		"root/missing.txt": {},
		"root/locked.txt":  {},
	}, zipBatchParse422Paths(body))
}

func TestZipBatchMockRejectsInvalidPathsParam(t *testing.T) {
	server := newZipBatchMockServer(t, nil, nil)
	defer server.Shutdown()

	for _, body := range []string{
		`{"paths":null,"encoded_paths":[]}`,
		`{"paths":"nope","encoded_paths":[]}`,
	} {
		res, err := http.Post(server.Server.URL+"/api/rest/v1/zip_downloads", "application/json", strings.NewReader(body))
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
		require.NoError(t, res.Body.Close())
	}
}

func TestZipBatchParamsDefaults(t *testing.T) {
	params := (ZipBatchParams{}).withDefaults()

	assert.Equal(t, int64(128*1024), params.EligibleSize)
	assert.Equal(t, 500, params.MinFiles)
	assert.Equal(t, 32, params.MaxFiles)
	assert.Equal(t, 0, params.BatchSize)
	assert.Equal(t, int64(25*1024*1024), params.MaxBytes)
	assert.Equal(t, 64, params.ConcurrentBatches)
	assert.Equal(t, defaultZipBatchMinAdvantage, params.MinAdvantage)
	assert.Equal(t, zipBatchTuning.DefaultReprobeInterval, params.ReprobeInterval)
	assert.Equal(t, ZipBatchExtractionStream, params.Extraction)
}

func TestZipBatchParamsNegativeValuesUseDefaults(t *testing.T) {
	params := (ZipBatchParams{
		EligibleSize:      -1,
		MinFiles:          -1,
		MaxFiles:          -1,
		BatchSize:         -1,
		MaxBytes:          -1,
		ConcurrentBatches: -1,
		ReprobeInterval:   -1,
	}).withDefaults()

	assert.Equal(t, int64(defaultZipBatchEligibleSize), params.EligibleSize)
	assert.Equal(t, defaultZipBatchMinFiles, params.MinFiles)
	assert.Equal(t, defaultZipBatchMaxFiles, params.MaxFiles)
	assert.Equal(t, 0, params.BatchSize)
	assert.Equal(t, int64(defaultZipBatchMaxBytes), params.MaxBytes)
	assert.Equal(t, defaultZipBatchConcurrentBatches, params.ConcurrentBatches)
	assert.Equal(t, defaultZipBatchMinAdvantage, params.MinAdvantage)
	assert.Equal(t, time.Duration(-1), params.ReprobeInterval)
}

func TestZipBatchEffectiveConcurrentBatchesCapsToFileSlots(t *testing.T) {
	assert.Equal(t, 37, zipBatchEffectiveConcurrentBatches(64, 50))
	assert.Equal(t, 1, zipBatchEffectiveConcurrentBatches(64, 1))
	assert.Equal(t, 2, zipBatchEffectiveConcurrentBatches(2, 50))
}

func TestZipBatchDownloaderNegativeConcurrencyUsesDefault(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/a.txt": "alpha",
		"batch/b.txt": "bravo",
	}, nil)
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{MinFiles: 1, ConcurrentBatches: -1, MinAdvantage: -1}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	assert.NotEmpty(t, server.ZipCreateRequests)
}

func TestZipBatchDownloaderHappyPath(t *testing.T) {
	forEachZipBatchExtractionMode(t, func(t *testing.T, extraction ZipBatchExtractionMode) {
		root := t.TempDir()
		mtime := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
		server := newZipBatchMockServer(t, map[string]string{
			"batch/a.txt": "alpha",
			"batch/b.txt": "bravo",
		}, func(path string, file files_sdk.File) files_sdk.File {
			file.ProvidedMtime = &mtime
			return file
		})
		defer server.Shutdown()

		job := runZipBatchDownload(t, server, "batch", root, zipBatchParamsForMode(extraction, ZipBatchParams{MinFiles: 2}), true, RetryPolicy{})

		requireZipBatchJobClean(t, job)
		assertZipBatchStats(t, job, ZipBatchStatsSnapshot{
			BatchesDispatched: 1,
			BatchFiles:        2,
			CreateRequests:    1,
			StreamAttempts:    1,
			CleanFinalized:    2,
		})
		assert.Equal(t, [][]string{{"batch/a.txt", "batch/b.txt"}}, server.ZipCreateRequests)
		assert.Equal(t, "alpha", readLocalFile(t, root, "batch", "a.txt"))
		assert.Equal(t, "bravo", readLocalFile(t, root, "batch", "b.txt"))
		stat, err := os.Stat(filepath.Join(root, "batch", "a.txt"))
		require.NoError(t, err)
		assert.True(t, stat.ModTime().Equal(mtime.Local()))
	})
}

func TestZipBatchDownloaderDisabledDoesNotCreateZipDownloads(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/a.txt": "alpha",
		"batch/b.txt": "bravo",
	}, nil)
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{Disabled: true, MinFiles: 1}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	assert.Empty(t, server.ZipCreateRequests)
	assert.NotEmpty(t, server.TrackRequest["/download/:download_id"])
}

func TestZipBatchDownloaderDryRunDoesNotCreateZipDownloads(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/a.txt": "alpha",
		"batch/b.txt": "bravo",
	}, nil)
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		DryRun:      true,
		RetryPolicy: RetryPolicy{},
		ZipBatch:    ZipBatchParams{MinFiles: 1},
	})

	requireZipBatchJobClean(t, job)
	assert.Empty(t, server.ZipCreateRequests)
	assert.NoFileExists(t, filepath.Join(root, "batch", "a.txt"))
}

func TestZipBatchDownloaderSingleStreamManagerDoesNotCreateZipDownloads(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/a.txt": "alpha",
		"batch/b.txt": "bravo",
	}, nil)
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.Build(2, 2, true),
		RetryPolicy: RetryPolicy{},
		ZipBatch:    ZipBatchParams{MinFiles: 1},
	})

	requireZipBatchJobClean(t, job)
	assert.Empty(t, server.ZipCreateRequests)
	assert.NotEmpty(t, server.TrackRequest["/download/:download_id"])
}

func TestZipBatchDownloaderDuplicateBasenames(t *testing.T) {
	forEachZipBatchExtractionMode(t, func(t *testing.T, extraction ZipBatchExtractionMode) {
		root := t.TempDir()
		server := newZipBatchMockServer(t, map[string]string{
			"batch/a/same.txt": "first",
			"batch/b/same.txt": "second",
		}, nil)
		defer server.Shutdown()

		job := runZipBatchDownload(t, server, "batch", root, zipBatchParamsForMode(extraction, ZipBatchParams{MinFiles: 2}), false, RetryPolicy{})

		requireZipBatchJobClean(t, job)
		assert.Equal(t, "first", readLocalFile(t, root, "batch", "a", "same.txt"))
		assert.Equal(t, "second", readLocalFile(t, root, "batch", "b", "same.txt"))
		assert.Equal(t, [][]string{{"batch/a/same.txt", "batch/b/same.txt"}}, server.ZipCreateRequests)
	})
}

func TestZipBatchDownloaderCorruptionSalvagesPrefixAndRetriesRemaining(t *testing.T) {
	forEachZipBatchExtractionMode(t, func(t *testing.T, extraction ZipBatchExtractionMode) {
		root := t.TempDir()
		server := newZipBatchMockServer(t, map[string]string{
			"batch/a.txt": "alpha",
			"batch/b.txt": "bravo",
			"batch/c.txt": "charlie",
		}, nil)
		server.ZipDownloadMutator = func(attempt int, _ []string, body []byte) []byte {
			if attempt != 1 {
				return body
			}
			return corruptAfterFirstZipEntry(body)
		}
		defer server.Shutdown()

		job := runZipBatchDownload(t, server, "batch", root, zipBatchParamsForMode(extraction, ZipBatchParams{MinFiles: 2}), false, RetryPolicy{RetryCount: 1})

		requireZipBatchJobClean(t, job)
		want := ZipBatchStatsSnapshot{
			BatchesDispatched: 1,
			BatchFiles:        3,
			CreateRequests:    2,
			StreamAttempts:    2,
			StreamFailures:    1,
			CleanFinalized:    2,
			SalvageFinalized:  1,
		}
		if extraction == ZipBatchExtractionStream {
			want.CleanFinalized = 3
			want.SalvageFinalized = 0
		}
		assertZipBatchStats(t, job, want)
		require.Len(t, server.ZipCreateRequests, 2)
		assert.Equal(t, []string{"batch/a.txt", "batch/b.txt", "batch/c.txt"}, server.ZipCreateRequests[0])
		assert.Equal(t, []string{"batch/b.txt", "batch/c.txt"}, server.ZipCreateRequests[1])
		assert.Equal(t, "alpha", readLocalFile(t, root, "batch", "a.txt"))
		assert.Equal(t, "bravo", readLocalFile(t, root, "batch", "b.txt"))
		assert.Equal(t, "charlie", readLocalFile(t, root, "batch", "c.txt"))
	})
}

func TestZipBatchDownloaderStreamRetryResetsPartialProgress(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/a.txt": strings.Repeat("abcdef", 512),
		"batch/b.txt": "bravo",
	}, nil)
	var mutated bool
	server.ZipDownloadMutator = func(attempt int, _ []string, body []byte) []byte {
		if attempt != 1 {
			return body
		}
		corrupt, ok := corruptFirstZipDataDescriptorCRC(body)
		mutated = ok
		if !ok {
			return body
		}
		return corrupt
	}
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{Extraction: ZipBatchExtractionStream, MinFiles: 1, MinAdvantage: -1}, false, RetryPolicy{RetryCount: 1})

	require.True(t, mutated, "test zip should use data descriptors")
	requireZipBatchJobClean(t, job)
	assert.Equal(t, job.TotalBytes(status.Included...), job.TransferBytes(status.Included...))
	assert.Equal(t, "bravo", readLocalFile(t, root, "batch", "b.txt"))
}

func TestZipBatchDownloaderRetriesExhaustedFallsBackPerFile(t *testing.T) {
	forEachZipBatchExtractionMode(t, func(t *testing.T, extraction ZipBatchExtractionMode) {
		root := t.TempDir()
		server := newZipBatchMockServer(t, map[string]string{
			"batch/a.txt": "alpha",
			"batch/b.txt": "bravo",
		}, nil)
		server.ZipDownloadMutator = func(_ int, _ []string, body []byte) []byte {
			return corruptAfterFirstZipEntry(body)
		}
		defer server.Shutdown()

		job := runZipBatchDownload(t, server, "batch", root, zipBatchParamsForMode(extraction, ZipBatchParams{MinFiles: 2}), false, RetryPolicy{})

		requireZipBatchJobClean(t, job)
		want := ZipBatchStatsSnapshot{
			BatchesDispatched:        1,
			BatchFiles:               2,
			CreateRequests:           1,
			StreamAttempts:           1,
			StreamFailures:           1,
			SalvageFinalized:         1,
			FallbackRetriesExhausted: 1,
		}
		if extraction == ZipBatchExtractionStream {
			want.CleanFinalized = 1
			want.SalvageFinalized = 0
		}
		assertZipBatchStats(t, job, want)
		assert.Equal(t, "alpha", readLocalFile(t, root, "batch", "a.txt"))
		assert.Equal(t, "bravo", readLocalFile(t, root, "batch", "b.txt"))
		assert.NotEmpty(t, server.TrackRequest["/download/:download_id"])
	})
}

func TestZipBatchDownloaderStoredDescriptorIsTripwire(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/a.txt": "alpha",
	}, nil)
	server.ZipDownloadMutator = func(_ int, _ []string, _ []byte) []byte {
		return storedDescriptorZip("a.txt", []byte("alpha"))
	}
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{Extraction: ZipBatchExtractionStream, MinFiles: 1, MinAdvantage: -1}, false, RetryPolicy{RetryCount: 2})

	requireZipBatchJobClean(t, job)
	assertZipBatchStats(t, job, ZipBatchStatsSnapshot{
		BatchesDispatched: 1,
		BatchFiles:        1,
		CreateRequests:    1,
		StreamAttempts:    1,
		FallbackTripwire:  1,
	})
	assert.Equal(t, "alpha", readLocalFile(t, root, "batch", "a.txt"))
}

func TestZipBatchDownloaderWrongSizeEntryFallsBackAsTripwire(t *testing.T) {
	forEachZipBatchExtractionMode(t, func(t *testing.T, extraction ZipBatchExtractionMode) {
		root := t.TempDir()
		server := newZipBatchMockServer(t, map[string]string{
			"batch/a.txt": "alp",
		}, nil)
		server.ZipDownloadMutator = func(_ int, _ []string, _ []byte) []byte {
			return zipWithEntry("a.txt", []byte("alpha"))
		}
		defer server.Shutdown()

		job := runZipBatchDownload(t, server, "batch", root, zipBatchParamsForMode(extraction, ZipBatchParams{MinFiles: 1}), false, RetryPolicy{})

		requireZipBatchJobClean(t, job)
		assertZipBatchStats(t, job, ZipBatchStatsSnapshot{
			BatchesDispatched: 1,
			BatchFiles:        1,
			CreateRequests:    1,
			StreamAttempts:    1,
			FallbackTripwire:  1,
		})
		assert.Equal(t, "alp", readLocalFile(t, root, "batch", "a.txt"))
	})
}

func TestZipBatchDownloader422Shrink(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/a.txt": "alpha",
		"batch/b.txt": "bravo",
		"batch/c.txt": "charlie",
	}, nil)
	var once sync.Once
	server.ZipCreateResponse = func(paths []string) (int, any, bool) {
		handled := false
		once.Do(func() {
			handled = true
		})
		if !handled {
			return 0, nil, false
		}
		return http.StatusUnprocessableEntity, files_sdk.ResponseError{
			ErrorMessage: "Paths does not match existing file at: batch/b.txt",
		}, true
	}
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{MinFiles: 2, MinAdvantage: -1}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	require.Len(t, server.ZipCreateRequests, 2)
	assert.Equal(t, []string{"batch/a.txt", "batch/b.txt", "batch/c.txt"}, server.ZipCreateRequests[0])
	assert.Equal(t, []string{"batch/a.txt", "batch/c.txt"}, server.ZipCreateRequests[1])
	assert.Equal(t, "bravo", readLocalFile(t, root, "batch", "b.txt"))
}

func TestZipBatchDownloaderMissingEntryFallsBack(t *testing.T) {
	forEachZipBatchExtractionMode(t, func(t *testing.T, extraction ZipBatchExtractionMode) {
		root := t.TempDir()
		server := newZipBatchMockServer(t, map[string]string{
			"batch/a.txt": "alpha",
			"batch/b.txt": "bravo",
		}, nil)
		server.ZipOmitPaths = map[string]bool{"batch/b.txt": true}
		defer server.Shutdown()

		job := runZipBatchDownload(t, server, "batch", root, zipBatchParamsForMode(extraction, ZipBatchParams{MinFiles: 2}), false, RetryPolicy{})

		requireZipBatchJobClean(t, job)
		assert.Equal(t, "alpha", readLocalFile(t, root, "batch", "a.txt"))
		assert.Equal(t, "bravo", readLocalFile(t, root, "batch", "b.txt"))
		assert.NotEmpty(t, server.TrackRequest["/download/:download_id"])
	})
}

func TestZipBatchDownloaderSyncSkipsUpToDateFiles(t *testing.T) {
	forEachZipBatchExtractionMode(t, func(t *testing.T, extraction ZipBatchExtractionMode) {
		root := t.TempDir()
		require.NoError(t, os.MkdirAll(filepath.Join(root, "batch"), 0755))
		require.NoError(t, os.WriteFile(filepath.Join(root, "batch", "a.txt"), []byte("alpha"), 0644))
		server := newZipBatchMockServer(t, map[string]string{
			"batch/a.txt": "alpha",
			"batch/b.txt": "bravo",
		}, nil)
		defer server.Shutdown()

		job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
			RemotePath:    "batch",
			LocalPath:     root + string(os.PathSeparator),
			Sync:          true,
			RetryPolicy:   RetryPolicy{},
			ZipBatch:      zipBatchParamsForMode(extraction, ZipBatchParams{MinFiles: 1}),
			PreserveTimes: false,
		})

		requireZipBatchJobClean(t, job)
		assert.Equal(t, int64(2), job.ZipBatchStats().BatchFiles)
		assert.Equal(t, [][]string{{"batch/b.txt"}}, server.ZipCreateRequests)
		assert.Equal(t, status.Skipped, statusForPath(job, "batch/a.txt"))
	})
}

func TestZipBatchDownloaderPrepareRunsAtBatchDispatch(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "batch"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "batch", "a.txt"), []byte("alpha"), 0644))
	server := newZipBatchMockServer(t, map[string]string{"batch/a.txt": "alpha"}, nil)
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:    "batch",
		LocalPath:     root + string(os.PathSeparator),
		Sync:          true,
		RetryPolicy:   RetryPolicy{},
		ZipBatch:      ZipBatchParams{MinFiles: 1, BatchSize: 1, MinAdvantage: -1},
		PreserveTimes: false,
	})

	requireZipBatchJobClean(t, job)
	assert.Equal(t, status.Skipped, statusForPath(job, "batch/a.txt"))
	assert.Equal(t, int64(0), job.TransferBytes())
	assert.Equal(t, int64(1), job.ZipBatchStats().BatchFiles)
	assert.Empty(t, server.ZipCreateRequests)
	assert.Empty(t, server.DownloadRequests)
}

func TestZipBatchDownloaderMinFilesDissolves(t *testing.T) {
	forEachZipBatchExtractionMode(t, func(t *testing.T, extraction ZipBatchExtractionMode) {
		root := t.TempDir()
		server := newZipBatchMockServer(t, map[string]string{
			"batch/a.txt": "alpha",
			"batch/b.txt": "bravo",
		}, nil)
		defer server.Shutdown()

		job := runZipBatchDownload(t, server, "batch", root, zipBatchParamsForMode(extraction, ZipBatchParams{MinFiles: 3}), false, RetryPolicy{})

		requireZipBatchJobClean(t, job)
		assertZipBatchStats(t, job, ZipBatchStatsSnapshot{
			BatchesDissolved: 1,
			DissolvedFiles:   2,
		})
		assert.Empty(t, server.ZipCreateRequests)
		assert.Equal(t, "alpha", readLocalFile(t, root, "batch", "a.txt"))
		assert.Equal(t, "bravo", readLocalFile(t, root, "batch", "b.txt"))
	})
}

func TestZipBatchDownloaderMinFilesGateBlocksSmallJobs(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(20), nil)
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{MinFiles: 500, MaxFiles: 16}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	assertZipBatchStats(t, job, ZipBatchStatsSnapshot{
		BatchesDissolved: 1,
		DissolvedFiles:   20,
	})
	assert.Empty(t, server.ZipCreateRequests)
	assert.NotEmpty(t, server.TrackRequest["/download/:download_id"])
}

func TestZipBatchDownloaderMinFilesGateIncludesAccumulatedFiles(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, zipBatchTestFiles(20), nil)
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{MinFiles: 17, MaxFiles: 16, MinAdvantage: -1}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	assertZipBatchStats(t, job, ZipBatchStatsSnapshot{
		BatchesDispatched: 2,
		BatchFiles:        20,
		CreateRequests:    2,
		StreamAttempts:    2,
		CleanFinalized:    20,
	})
	assert.ElementsMatch(t, []int{4, 16}, zipBatchRequestLengths(server.ZipCreateRequests))
	assert.ElementsMatch(t, zipBatchTestPaths(20), flattenZipBatchRequests(server.ZipCreateRequests))
}

func TestZipBatchDownloaderDynamicBatchSizeGrowth(t *testing.T) {
	root := t.TempDir()
	fileCount := 260
	server := newZipBatchMockServer(t, zipBatchTestFiles(fileCount), nil)
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{
		MinFiles:          1,
		MaxFiles:          64,
		ConcurrentBatches: 4,
		MinAdvantage:      -1,
	}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	lengths := zipBatchRequestLengths(server.ZipCreateRequests)
	assert.ElementsMatch(t, zipBatchExpectedBatchLengths(fileCount, 1, 64, 4, 0), lengths)
	assert.Contains(t, lengths, 16)
	assert.True(t, zipBatchAnyLengthOver(lengths, 16))
	assert.False(t, zipBatchAnyLengthOver(lengths, 64))
}

func TestZipBatchDownloaderPinnedBatchSize(t *testing.T) {
	root := t.TempDir()
	fileCount := 45
	server := newZipBatchMockServer(t, zipBatchTestFiles(fileCount), nil)
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{
		MinFiles:          1,
		MaxFiles:          64,
		BatchSize:         20,
		ConcurrentBatches: 4,
		MinAdvantage:      -1,
	}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	assert.ElementsMatch(t, []int{20, 20, 5}, zipBatchRequestLengths(server.ZipCreateRequests))
}

func TestZipBatchDownloaderMaxFilesFlushes(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/a.txt": "alpha",
		"batch/b.txt": "bravo",
		"batch/c.txt": "charlie",
		"batch/d.txt": "delta",
		"batch/e.txt": "echo",
	}, nil)
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{MinFiles: 1, MaxFiles: 2, MinAdvantage: -1}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	require.Len(t, server.ZipCreateRequests, 3)
	assert.ElementsMatch(t, []int{1, 2, 2}, zipBatchRequestLengths(server.ZipCreateRequests))
	assert.ElementsMatch(t, []string{"batch/a.txt", "batch/b.txt", "batch/c.txt", "batch/d.txt", "batch/e.txt"}, flattenZipBatchRequests(server.ZipCreateRequests))
}

func TestZipBatchDownloaderMaxBytesFlushes(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/a.txt": "aaaaaa",
		"batch/b.txt": "bbbbbb",
		"batch/c.txt": "cccccc",
	}, nil)
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{MinFiles: 1, MaxBytes: 10, MinAdvantage: -1}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	require.Len(t, server.ZipCreateRequests, 3)
	assert.ElementsMatch(t, []int{1, 1, 1}, zipBatchRequestLengths(server.ZipCreateRequests))
	assert.ElementsMatch(t, []string{"batch/a.txt", "batch/b.txt", "batch/c.txt"}, flattenZipBatchRequests(server.ZipCreateRequests))
}

func TestZipBatchDownloaderEligibleSizeFallsBackPerFile(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/small.txt": "tiny",
		"batch/large.txt": "larger",
	}, nil)
	defer server.Shutdown()

	job := runZipBatchDownload(t, server, "batch", root, ZipBatchParams{EligibleSize: 5, MinFiles: 1, MinAdvantage: -1}, false, RetryPolicy{})

	requireZipBatchJobClean(t, job)
	assert.Equal(t, [][]string{{"batch/small.txt"}}, server.ZipCreateRequests)
	assert.NotEmpty(t, server.TrackRequest["/download/:download_id"])
	assert.Equal(t, "tiny", readLocalFile(t, root, "batch", "small.txt"))
	assert.Equal(t, "larger", readLocalFile(t, root, "batch", "large.txt"))
}

func TestZipBatchDownloaderPrepareHandledFilesExcludedBeforeCreate(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "batch"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "batch", "existing.txt"), []byte("local"), 0644))
	server := newZipBatchMockServer(t, map[string]string{
		"batch/existing.txt": "remote",
		"batch/new.txt":      "new",
	}, nil)
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:    "batch",
		LocalPath:     root + string(os.PathSeparator),
		NoOverwrite:   true,
		RetryPolicy:   RetryPolicy{},
		ZipBatch:      ZipBatchParams{MinFiles: 2, MaxFiles: 1, MinAdvantage: -1},
		PreserveTimes: false,
	})

	requireZipBatchJobClean(t, job)
	assert.Equal(t, int64(2), job.ZipBatchStats().BatchFiles)
	assert.Equal(t, [][]string{{"batch/new.txt"}}, server.ZipCreateRequests)
	assert.Equal(t, status.FileExists, statusForPath(job, "batch/existing.txt"))
	assert.Equal(t, "new", readLocalFile(t, root, "batch", "new.txt"))
}

func TestZipBatchDownloaderFallbackEnqueuesAfterBatchAdmissionRelease(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/a.txt": "alpha",
		"batch/b.txt": "bravo",
	}, nil)
	server.ZipDownloadMutator = func(_ int, _ []string, body []byte) []byte {
		return corruptAfterFirstZipEntry(body)
	}
	defer server.Shutdown()

	job := runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:    "batch",
		LocalPath:     root + string(os.PathSeparator),
		Manager:       manager.New(1, 1, 1),
		RetryPolicy:   RetryPolicy{},
		ZipBatch:      ZipBatchParams{MinFiles: 2, MaxFiles: 2, ConcurrentBatches: 1, MinAdvantage: -1},
		PreserveTimes: false,
	})

	requireZipBatchJobClean(t, job)
	assert.NotEmpty(t, server.TrackRequest["/download/:download_id"])
	assert.Equal(t, "alpha", readLocalFile(t, root, "batch", "a.txt"))
	assert.Equal(t, "bravo", readLocalFile(t, root, "batch", "b.txt"))
}

func TestZipBatchDownloaderLeavesFileAdmissionForLargeFiles(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/a.txt":     "alpha",
		"batch/b.txt":     "bravo",
		"batch/large.bin": strings.Repeat("x", defaultZipBatchEligibleSize+1),
	}, nil)
	blocked := make(chan struct{})
	release := make(chan struct{})
	var blockOnce sync.Once
	var releaseOnce sync.Once
	server.ZipDownloadBefore = func(_ int, _ []string) {
		blockOnce.Do(func() { close(blocked) })
		<-release
	}
	defer func() {
		releaseOnce.Do(func() { close(release) })
		server.Shutdown()
	}()

	client := server.Client()
	job := client.Downloader(DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		Manager:     manager.New(2, 1, 1),
		RetryPolicy: RetryPolicy{},
		ZipBatch:    ZipBatchParams{MinFiles: 1, BatchSize: 2, ConcurrentBatches: 64, MinAdvantage: -1},
	})
	done := make(chan struct{})
	go func() {
		job.Start()
		job.Wait()
		close(done)
	}()

	select {
	case <-blocked:
	case <-time.After(2 * time.Second):
		t.Fatal("zip download did not start")
	}
	require.Eventually(t, func() bool {
		return zipBatchTrackedRequestCount(server, "/download/:download_id") > 0
	}, 2*time.Second, 10*time.Millisecond)
	releaseOnce.Do(func() { close(release) })
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("download did not finish")
	}
	requireZipBatchJobClean(t, job)
}

func TestZipBatchDownloaderStreamModeDoesNotCreateSpoolFile(t *testing.T) {
	root := t.TempDir()
	server := newZipBatchMockServer(t, map[string]string{
		"batch/a.txt": "alpha",
		"batch/b.txt": "bravo",
	}, nil)
	blocked := make(chan struct{})
	release := make(chan struct{})
	var blockOnce sync.Once
	var releaseOnce sync.Once
	server.ZipDownloadBefore = func(_ int, _ []string) {
		blockOnce.Do(func() { close(blocked) })
		<-release
	}
	defer func() {
		releaseOnce.Do(func() { close(release) })
		server.Shutdown()
	}()

	client := server.Client()
	job := client.Downloader(DownloaderParams{
		RemotePath:  "batch",
		LocalPath:   root + string(os.PathSeparator),
		RetryPolicy: RetryPolicy{},
		ZipBatch:    ZipBatchParams{Extraction: ZipBatchExtractionStream, MinFiles: 2, MinAdvantage: -1},
	})
	done := make(chan struct{})
	go func() {
		job.Start()
		job.Wait()
		close(done)
	}()

	select {
	case <-blocked:
	case <-time.After(2 * time.Second):
		t.Fatal("zip download did not start")
	}
	assertNoZipBatchSpoolFiles(t, root)
	releaseOnce.Do(func() { close(release) })
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("download did not finish")
	}

	requireZipBatchJobClean(t, job)
	assertNoZipBatchSpoolFiles(t, root)
}

func withZipBatchProbeTestWindows(t *testing.T) {
	t.Helper()
	oldPhaseAMinCompletions := zipBatchTuning.ProbePhaseAMinCompletions
	oldPhaseAMinDuration := zipBatchTuning.ProbePhaseAMinDuration
	oldPhaseAHardCap := zipBatchTuning.ProbePhaseAHardCap
	oldPhaseAWindow := zipBatchTuning.ProbePhaseAWindow
	oldPhaseAWindow2After := zipBatchTuning.ProbePhaseAWindow2After
	oldPhaseAWindow3After := zipBatchTuning.ProbePhaseAWindow3After
	oldPhaseBMinCompletions := zipBatchTuning.ProbePhaseBMinCompletions
	oldPhaseBWarmup := zipBatchTuning.ProbePhaseBWarmup
	oldPhaseBWindow := zipBatchTuning.ProbePhaseBWindow
	oldPhaseBEndMinWindow := zipBatchTuning.ProbePhaseBEndMinWindow
	oldPhaseBHardCap := zipBatchTuning.ProbePhaseBHardCap
	oldDefaultReprobeInterval := zipBatchTuning.DefaultReprobeInterval
	oldMaxReprobeInterval := zipBatchTuning.MaxReprobeInterval

	zipBatchTuning.ProbePhaseAMinCompletions = 16
	zipBatchTuning.ProbePhaseAMinDuration = 20 * time.Millisecond
	zipBatchTuning.ProbePhaseAHardCap = 250 * time.Millisecond
	zipBatchTuning.ProbePhaseAWindow = 10 * time.Millisecond
	zipBatchTuning.ProbePhaseAWindow2After = 2000
	zipBatchTuning.ProbePhaseAWindow3After = 5000
	zipBatchTuning.ProbePhaseBMinCompletions = 16
	zipBatchTuning.ProbePhaseBWarmup = time.Millisecond
	zipBatchTuning.ProbePhaseBWindow = 10 * time.Millisecond
	zipBatchTuning.ProbePhaseBEndMinWindow = 2 * time.Millisecond
	zipBatchTuning.ProbePhaseBHardCap = 250 * time.Millisecond
	zipBatchTuning.DefaultReprobeInterval = 20 * time.Millisecond
	zipBatchTuning.MaxReprobeInterval = 80 * time.Millisecond

	t.Cleanup(func() {
		zipBatchTuning.ProbePhaseAMinCompletions = oldPhaseAMinCompletions
		zipBatchTuning.ProbePhaseAMinDuration = oldPhaseAMinDuration
		zipBatchTuning.ProbePhaseAHardCap = oldPhaseAHardCap
		zipBatchTuning.ProbePhaseAWindow = oldPhaseAWindow
		zipBatchTuning.ProbePhaseAWindow2After = oldPhaseAWindow2After
		zipBatchTuning.ProbePhaseAWindow3After = oldPhaseAWindow3After
		zipBatchTuning.ProbePhaseBMinCompletions = oldPhaseBMinCompletions
		zipBatchTuning.ProbePhaseBWarmup = oldPhaseBWarmup
		zipBatchTuning.ProbePhaseBWindow = oldPhaseBWindow
		zipBatchTuning.ProbePhaseBEndMinWindow = oldPhaseBEndMinWindow
		zipBatchTuning.ProbePhaseBHardCap = oldPhaseBHardCap
		zipBatchTuning.DefaultReprobeInterval = oldDefaultReprobeInterval
		zipBatchTuning.MaxReprobeInterval = oldMaxReprobeInterval
	})
}

func newZipBatchMockServer(t *testing.T, files map[string]string, edit func(string, files_sdk.File) files_sdk.File) *MockAPIServer {
	server := (&MockAPIServer{T: t}).Do()
	server.MockFiles["batch"] = mockFile{File: files_sdk.File{Path: "batch", DisplayName: "batch", Type: "directory"}}
	dirs := map[string]struct{}{"batch": {}}
	for path, contents := range files {
		dir := filepath.ToSlash(filepath.Dir(path))
		for dir != "." && dir != "" {
			if _, ok := dirs[dir]; !ok {
				dirs[dir] = struct{}{}
				server.MockFiles[dir] = mockFile{File: files_sdk.File{Path: dir, DisplayName: filepath.Base(dir), Type: "directory"}}
			}
			next := filepath.ToSlash(filepath.Dir(dir))
			if next == dir {
				break
			}
			dir = next
		}
		file := files_sdk.File{Path: path, DisplayName: filepath.Base(path), Type: "file", Size: int64(len(contents))}
		if edit != nil {
			file = edit(path, file)
		}
		server.MockFiles[path] = mockFile{File: file, Data: []byte(contents), SizeTrust: TrustedSizeValue}
	}
	return server
}

func forEachZipBatchExtractionMode(t *testing.T, run func(*testing.T, ZipBatchExtractionMode)) {
	t.Helper()
	for _, extraction := range []ZipBatchExtractionMode{ZipBatchExtractionSpool, ZipBatchExtractionStream} {
		t.Run(string(extraction), func(t *testing.T) {
			run(t, extraction)
		})
	}
}

func zipBatchParamsForMode(extraction ZipBatchExtractionMode, params ZipBatchParams) ZipBatchParams {
	params.Extraction = extraction
	if params.MinAdvantage == 0 {
		params.MinAdvantage = -1
	}
	return params
}

func runZipBatchDownload(t *testing.T, server *MockAPIServer, remotePath, localRoot string, zipParams ZipBatchParams, preserveTimes bool, retryPolicy RetryPolicy) *Job {
	return runZipBatchDownloadWithParams(t, server, DownloaderParams{
		RemotePath:    remotePath,
		LocalPath:     localRoot + string(os.PathSeparator),
		RetryPolicy:   retryPolicy,
		ZipBatch:      zipParams,
		PreserveTimes: preserveTimes,
	})
}

func runZipBatchDownloadWithParams(t *testing.T, server *MockAPIServer, params DownloaderParams) *Job {
	t.Helper()
	client := server.Client()
	job := client.Downloader(params)
	job.Start()
	job.Wait()
	return job
}

func requireZipBatchJobClean(t *testing.T, job *Job) {
	t.Helper()
	require.True(t, job.All(status.Ended...), "all statuses should be ended")
	for _, st := range job.Statuses {
		require.NoError(t, st.Err(), st.RemotePath())
	}
}

func assertZipBatchStats(t *testing.T, job *Job, want ZipBatchStatsSnapshot) {
	t.Helper()
	if want.ProbeDecision == "" {
		want.ProbeDecision = string(zipBatchProbeDecisionNone)
	}
	assert.Equal(t, want, job.ZipBatchStats())
}

func readLocalFile(t *testing.T, root string, parts ...string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(append([]string{root}, parts...)...))
	require.NoError(t, err)
	return string(data)
}

func statusForPath(job *Job, remotePath string) status.Status {
	for _, st := range job.Statuses {
		if st.RemotePath() == remotePath {
			return st.Status()
		}
	}
	return status.Null
}

func zipBatchTestFiles(count int) map[string]string {
	files := make(map[string]string, count)
	for i := 0; i < count; i++ {
		files[fmt.Sprintf("batch/file-%03d.txt", i)] = fmt.Sprintf("file-%03d", i)
	}
	return files
}

func zipBatchTestPaths(count int) []string {
	paths := make([]string, 0, count)
	for i := 0; i < count; i++ {
		paths = append(paths, fmt.Sprintf("batch/file-%03d.txt", i))
	}
	return paths
}

func zipBatchExpectedBatchLengths(total, minFiles, maxFiles, concurrentBatches, pinnedSize int) []int {
	var lengths []int
	pending := 0
	eligibleSeen := 0
	unlocked := false
	for i := 0; i < total; i++ {
		pending++
		eligibleSeen++
		if !unlocked && eligibleSeen >= minFiles {
			unlocked = true
		}
		for unlocked {
			size := zipBatchExpectedBatchSize(eligibleSeen, maxFiles, concurrentBatches, pinnedSize)
			if pending < size {
				break
			}
			lengths = append(lengths, size)
			pending -= size
		}
	}
	if unlocked && pending > 0 {
		lengths = append(lengths, pending)
	}
	return lengths
}

func zipBatchExpectedBatchSize(eligibleSeen, maxFiles, concurrentBatches, pinnedSize int) int {
	size := pinnedSize
	if size <= 0 {
		size = eligibleSeen / (2 * concurrentBatches)
		if size < defaultZipBatchDynamicBatchFloor {
			size = defaultZipBatchDynamicBatchFloor
		}
	}
	if size > maxFiles {
		size = maxFiles
	}
	if size < 1 {
		size = 1
	}
	return size
}

func zipBatchAnyLengthOver(lengths []int, limit int) bool {
	for _, length := range lengths {
		if length > limit {
			return true
		}
	}
	return false
}

func zipBatchTrackedRequestCount(server *MockAPIServer, route string) int {
	server.traceMutex.Lock()
	defer server.traceMutex.Unlock()
	return len(server.TrackRequest[route])
}

func zipBatchDownloadRequests(server *MockAPIServer) []string {
	server.traceMutex.Lock()
	defer server.traceMutex.Unlock()
	return append([]string(nil), server.DownloadRequests...)
}

func zipBatchRequestLengths(requests [][]string) []int {
	lengths := make([]int, 0, len(requests))
	for _, request := range requests {
		lengths = append(lengths, len(request))
	}
	return lengths
}

func flattenZipBatchRequests(requests [][]string) []string {
	var paths []string
	for _, request := range requests {
		paths = append(paths, request...)
	}
	return paths
}

func assertNoZipBatchSpoolFiles(t *testing.T, root string) {
	t.Helper()
	var paths []string
	require.NoError(t, filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(filepath.Base(path), ".zip-batch") {
			paths = append(paths, path)
		}
		return nil
	}))
	assert.Empty(t, paths)
}

func corruptAfterFirstZipEntry(body []byte) []byte {
	signature := []byte{0x50, 0x4b, 0x03, 0x04}
	next := bytes.Index(body[len(signature):], signature)
	if next < 0 {
		return append(body, []byte("9999999")...)
	}
	cut := len(signature) + next
	return append(append([]byte(nil), body[:cut]...), []byte("9999999")...)
}

func corruptFirstZipDataDescriptorCRC(body []byte) ([]byte, bool) {
	signature := []byte{0x50, 0x4b, 0x07, 0x08}
	offset := bytes.Index(body, signature)
	if offset < 0 || offset+4 >= len(body) {
		return nil, false
	}
	corrupt := append([]byte(nil), body...)
	corrupt[offset+4] ^= 0xff
	return corrupt, true
}

func zipWithEntry(name string, data []byte) []byte {
	out := bytes.NewBuffer(nil)
	writer := zip.NewWriter(out)
	entry, _ := writer.Create(name)
	_, _ = entry.Write(data)
	_ = writer.Close()
	return out.Bytes()
}

func storedDescriptorZip(name string, data []byte) []byte {
	out := bytes.NewBuffer(nil)
	writeUint32 := func(v uint32) { _ = binary.Write(out, binary.LittleEndian, v) }
	writeUint16 := func(v uint16) { _ = binary.Write(out, binary.LittleEndian, v) }
	crc := crc32.ChecksumIEEE(data)
	size := uint32(len(data))

	writeUint32(zipLocalFileHeaderSignature)
	writeUint16(20)
	writeUint16(zipFlagDataDescriptor)
	writeUint16(zipMethodStore)
	writeUint16(0)
	writeUint16(0)
	writeUint32(0)
	writeUint32(0)
	writeUint32(0)
	writeUint16(uint16(len(name)))
	writeUint16(0)
	out.WriteString(name)
	out.Write(data)
	writeUint32(zipDataDescriptorSignature)
	writeUint32(crc)
	writeUint32(size)
	writeUint32(size)

	centralOffset := out.Len()
	writeUint32(zipCentralDirectorySignature)
	writeUint16(20)
	writeUint16(20)
	writeUint16(zipFlagDataDescriptor)
	writeUint16(zipMethodStore)
	writeUint16(0)
	writeUint16(0)
	writeUint32(crc)
	writeUint32(size)
	writeUint32(size)
	writeUint16(uint16(len(name)))
	writeUint16(0)
	writeUint16(0)
	writeUint16(0)
	writeUint16(0)
	writeUint32(0)
	writeUint32(0)
	out.WriteString(name)
	centralSize := out.Len() - centralOffset

	writeUint32(zipEndOfCentralDirSignature)
	writeUint16(0)
	writeUint16(0)
	writeUint16(1)
	writeUint16(1)
	writeUint32(uint32(centralSize))
	writeUint32(uint32(centralOffset))
	writeUint16(0)
	return out.Bytes()
}
