package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lib/direction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sizeFixtureMiB = int64(1024 * 1024)

// sizeFixture is an API and transfer server whose listing and download-request
// metadata can disagree with the bytes it serves. Transfer URLs are pinned to
// the source size they were issued for, like Agent size pinning: once the
// source changes size, every range of an old URL answers 409
// download_source_changed and only a new download request succeeds.
type sizeFixture struct {
	mu  sync.Mutex
	url string

	metadataSize int64 // size reported by listing and download-request JSON
	size         int64 // size of the bytes served
	version      int   // content generation; changes with the source

	lifecycles   int      // download URLs issued
	rejected     []string // path and range of every 409
	responses    int      // transfer responses selected since the last change
	servedRanges []string // Range header of every completed transfer response

	// metadataTracksSource makes the JSON metadata follow the current size;
	// otherwise it keeps reporting metadataSize.
	metadataTracksSource bool
	// changeAfterResponse mutates the source (size += changeDelta) after this
	// response snapshots its bytes and headers. Later requests see the new
	// source even while that response streams. With repeatChange it happens
	// again in every later lifecycle; otherwise only once.
	changeAfterResponse int
	changeDelta         int64
	repeatChange        bool
	// changeBeforeFirstTransfer mutates the source before the first transfer
	// request, so the first URL is rejected before any bytes are served.
	changeBeforeFirstTransfer bool
	changes                   int
	// conflictTotalOnResponse makes that selected range response report a
	// Content-Range total that contradicts the size it serves.
	conflictTotalOnResponse int
	// unsatisfiedContentRangeOnResponse makes that 206 carry "bytes */total", a
	// Content-Range that states no byte coverage.
	unsatisfiedContentRangeOnResponse int
	shortBody                         bool // send half of the promised bytes
	// rejectAs selects the 409 body: "" for the typed download_source_changed
	// JSON, "text" for plain text that merely mentions the token, "other" for a
	// different typed JSON error.
	rejectAs string
}

func newSizeFixture(t *testing.T, metadataSize, size int64) *sizeFixture {
	f := &sizeFixture{metadataSize: metadataSize, size: size}
	server := httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(server.Close)
	f.url = server.URL
	return f
}

func sizeFixtureByte(version int, offset int64) byte {
	return byte(offset%251) ^ byte(version*31+7)
}

func (f *sizeFixture) serve(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.Contains(r.URL.Path, "/api/rest/v1/file_actions/metadata/"):
		f.writeMetadata(w, false)
	case strings.Contains(r.URL.Path, "/api/rest/v1/files/"):
		f.writeMetadata(w, true)
	case strings.HasPrefix(r.URL.Path, "/transfer/"):
		f.serveTransfer(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"type":"not-found","error":"not found","http-code":404}`)
	}
}

func (f *sizeFixture) writeMetadata(w http.ResponseWriter, withDownloadURI bool) {
	f.mu.Lock()
	size := f.metadataSize
	if f.metadataTracksSource {
		size = f.size
	}
	metadata := map[string]any{"path": "fixture.bin", "display_name": "fixture.bin", "type": "file", "size": size}
	if withDownloadURI {
		f.lifecycles++
		metadata["download_uri"] = fmt.Sprintf("%s/transfer/%d?sz=%d", f.url, f.lifecycles, f.size)
	}
	f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(metadata)
}

func (f *sizeFixture) serveTransfer(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	if f.changeBeforeFirstTransfer && f.changes == 0 {
		f.changeSource()
	}
	pinned, _ := strconv.ParseInt(r.URL.Query().Get("sz"), 10, 64)
	if pinned != f.size {
		f.rejected = append(f.rejected, r.URL.Path+" "+r.Header.Get("Range"))
		rejectAs := f.rejectAs
		f.mu.Unlock()
		switch rejectAs {
		case "text":
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, "download_source_changed mentioned in a plain text conflict")
		case "other":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, `{"type":"processing-failure/conflict","error":"download_source_changed mentioned in another error","http-code":409}`)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, `{"type":"download_source_changed","error":"The source file changed while it was being downloaded. Request a new download URL and restart from the beginning.","http-code":409}`)
		}
		return
	}
	start, end := int64(0), f.size-1
	code := http.StatusOK
	if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
		spec := strings.TrimPrefix(rangeHeader, "bytes=")
		startSpec, endSpec, _ := strings.Cut(spec, "-")
		start, _ = strconv.ParseInt(startSpec, 10, 64)
		if endSpec != "" {
			end, _ = strconv.ParseInt(endSpec, 10, 64)
		}
		if start >= f.size {
			w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", f.size))
			f.mu.Unlock()
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}
		end = min(end, f.size-1)
		code = http.StatusPartialContent
	}
	transfer := r.Method != http.MethodHead
	if transfer {
		f.responses++
	}
	if code == http.StatusPartialContent {
		total := f.size
		if transfer && f.responses == f.conflictTotalOnResponse {
			total += sizeFixtureMiB
		}
		if transfer && f.responses == f.unsatisfiedContentRangeOnResponse {
			w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", total))
		} else {
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, total))
		}
	}
	length := end - start + 1
	send := length
	if f.shortBody {
		send = length / 2
	}
	version := f.version
	if transfer && f.responses == f.changeAfterResponse && (f.changes == 0 || f.repeatChange) {
		f.changeSource()
	}
	f.mu.Unlock()
	// A blocked response body must not prevent another range from receiving
	// its headers. Each response keeps the source generation it started with.
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", strconv.FormatInt(length, 10))
	w.WriteHeader(code)
	if !transfer {
		return
	}
	buf := make([]byte, 64*1024)
	for written := int64(0); written < send; {
		n := min(int64(len(buf)), send-written)
		for i := range buf[:n] {
			buf[i] = sizeFixtureByte(version, start+written+int64(i))
		}
		if _, err := w.Write(buf[:n]); err != nil {
			return
		}
		written += n
	}
	if send != length {
		return // the connection closes short of Content-Length
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.servedRanges = append(f.servedRanges, r.Header.Get("Range"))
}

func (f *sizeFixture) changeSource() {
	f.size += f.changeDelta
	f.version++
	f.changes++
	f.responses = 0
}

type sizeFixtureRun struct {
	adaptive     bool
	singleStream bool
	retryCount   int
}

func (f *sizeFixture) download(t *testing.T, run sizeFixtureRun) (*Job, string) {
	return f.downloadTo(t, t.TempDir(), run)
}

func (f *sizeFixture) downloadTo(t *testing.T, dir string, run sizeFixtureRun) (*Job, string) {
	dest := filepath.Join(dir, "output.bin")
	config := files_sdk.Config{
		APIKey:                 "fixture-only-not-a-real-key",
		EndpointOverride:       f.url,
		DisableDirectTransfers: true,
		Logger:                 log.New(io.Discard, "", 0),
	}.Init()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	job := (&Client{Config: config}).Downloader(DownloaderParams{
		RemotePath:          "fixture.bin",
		LocalPath:           dest,
		Manager:             manager.Build(2, 1, run.singleStream),
		AdaptiveConcurrency: run.adaptive,
		RetryPolicy:         RetryPolicy{Type: RetryUnfinished, RetryCount: run.retryCount, Backoff: -1},
	}, files_sdk.WithContext(ctx))
	job.Start()
	job.Wait()
	require.Len(t, job.Statuses, 1)
	return job, dest
}

func (f *sizeFixture) requireExactOutput(t *testing.T, dest string) {
	t.Helper()
	f.mu.Lock()
	size, version := f.size, f.version
	f.mu.Unlock()
	data, err := os.ReadFile(dest)
	require.NoError(t, err)
	require.Equal(t, size, int64(len(data)), "output size")
	for i, b := range data {
		if b != sizeFixtureByte(version, int64(i)) {
			t.Fatalf("byte %d does not match the current source content", i)
		}
	}
}

func (f *sizeFixture) stats() (lifecycles int, rejected []string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.lifecycles, append([]string(nil), f.rejected...)
}

func (f *sizeFixture) ranges() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.servedRanges...)
}

// The adaptive engine plans 16 MiB parts; the parts engine plans 10 MiB parts.
const adaptivePartBoundary = "bytes=16777216-"

func requireNoOutput(t *testing.T, dest string) {
	t.Helper()
	entries, err := os.ReadDir(filepath.Dir(dest))
	require.NoError(t, err)
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	assert.Empty(t, names, "failed download must not leave output or temp files")
}

func requireEachRangeRejectedOnce(t *testing.T, rejected []string) {
	t.Helper()
	seen := map[string]int{}
	for _, r := range rejected {
		seen[r]++
	}
	for r, n := range seen {
		assert.Equal(t, 1, n, "range %q was retried against a rejected download request", r)
	}
}

func TestDownloadStaleMetadataUsesTransferSize(t *testing.T) {
	adaptive := sizeFixtureRun{adaptive: true}
	singleStream := sizeFixtureRun{singleStream: true}
	cases := []struct {
		name     string
		metadata int64
		actual   int64
		run      sizeFixtureRun
		ranges   []string // exact transfer requests when they must be deterministic
	}{
		{"adaptive metadata smaller than transfer", 20 * sizeFixtureMiB, 21 * sizeFixtureMiB, adaptive,
			[]string{"bytes=0-16777215", "bytes=16777216-22020095"}},
		{"adaptive metadata larger than transfer", 40 * sizeFixtureMiB, 21 * sizeFixtureMiB, adaptive,
			[]string{"bytes=0-16777215", "bytes=16777216-22020095"}},
		{"adaptive metadata larger than transfer within one part", 40 * sizeFixtureMiB, 3 * sizeFixtureMiB, adaptive,
			[]string{"bytes=0-16777215"}},
		{"adaptive metadata zero for a small transfer", 0, 3000, adaptive, []string{""}},
		{"adaptive metadata larger than empty transfer", 40 * sizeFixtureMiB, 0, adaptive, []string{""}},
		{"single stream metadata larger than empty transfer", 20 * sizeFixtureMiB, 0, singleStream, []string{""}},
		{"single stream metadata larger than transfer", 20 * sizeFixtureMiB, sizeFixtureMiB, singleStream, []string{""}},
		{"single stream metadata smaller than transfer", sizeFixtureMiB, 3 * sizeFixtureMiB, singleStream, []string{""}},
		{"single stream metadata zero", 0, 3000, singleStream, []string{""}},
		{"empty file", 0, 0, adaptive, []string{""}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newSizeFixture(t, tc.metadata, tc.actual)
			job, dest := fixture.download(t, tc.run)

			require.NoError(t, job.Statuses[0].Err())
			assert.Equal(t, status.Complete, job.Statuses[0].Status())
			fixture.requireExactOutput(t, dest)
			assert.Equal(t, tc.actual, job.Statuses[0].Size(), "reported size must be the completed output size")
			assert.Equal(t, tc.actual, job.TransferBytes())
			lifecycles, rejected := fixture.stats()
			assert.Equal(t, 1, lifecycles, "a stable source needs one download request")
			assert.Empty(t, rejected)
			assert.ElementsMatch(t, tc.ranges, fixture.ranges(), "transfer requests must follow the transfer size, not the metadata")
		})
	}
}

func TestDownloadResumePrefixBeyondEmptyTransferRestarts(t *testing.T) {
	// A temp file left by an earlier attempt can be longer than the source is
	// now. Its bytes cannot be validated against the current source, so the
	// attempt fails and the retry starts over from offset 0.
	t.Run("retry budget 0 fails without output", func(t *testing.T) {
		fixture := newSizeFixture(t, 40*sizeFixtureMiB, 0)
		dir := t.TempDir()
		createPausedTmpFile(t, filepath.Join(dir, "output.bin"), 5*sizeFixtureMiB)
		job, dest := fixture.downloadTo(t, dir, sizeFixtureRun{adaptive: true})

		require.Error(t, job.Statuses[0].Err())
		requireNoOutput(t, dest)
	})
	t.Run("retry budget 1 completes empty", func(t *testing.T) {
		fixture := newSizeFixture(t, 40*sizeFixtureMiB, 0)
		dir := t.TempDir()
		createPausedTmpFile(t, filepath.Join(dir, "output.bin"), 5*sizeFixtureMiB)
		job, dest := fixture.downloadTo(t, dir, sizeFixtureRun{adaptive: true, retryCount: 1})

		require.NoError(t, job.Statuses[0].Err())
		fixture.requireExactOutput(t, dest)
		assert.Equal(t, int64(0), job.Statuses[0].Size())
		assert.Equal(t, int64(0), job.TransferBytes())
		assert.Equal(t, []string{""}, fixture.ranges(), "only the confirming full download is served")
	})
}

func TestDownloadPartsStaleMetadataConvergesOnRetry(t *testing.T) {
	// The non-adaptive part engine plans from metadata and fails safely when the
	// transfer reports a different size; the retry plans from the transfer size
	// learned from the responses instead of the same stale metadata.
	t.Run("retry budget 0 fails without output", func(t *testing.T) {
		fixture := newSizeFixture(t, 20*sizeFixtureMiB, 21*sizeFixtureMiB)
		job, dest := fixture.download(t, sizeFixtureRun{retryCount: 0})

		require.Error(t, job.Statuses[0].Err())
		requireNoOutput(t, dest)
	})
	t.Run("retry budget 1 completes with the transfer size", func(t *testing.T) {
		fixture := newSizeFixture(t, 20*sizeFixtureMiB, 21*sizeFixtureMiB)
		job, dest := fixture.download(t, sizeFixtureRun{retryCount: 1})

		require.NoError(t, job.Statuses[0].Err())
		fixture.requireExactOutput(t, dest)
		lifecycles, _ := fixture.stats()
		assert.Equal(t, 1, lifecycles, "a stable source keeps its download request across the retry")
	})
}

func TestDownloadAdaptiveConflictingRangeTotalsFail(t *testing.T) {
	fixture := newSizeFixture(t, 21*sizeFixtureMiB, 21*sizeFixtureMiB)
	fixture.conflictTotalOnResponse = 2
	job, dest := fixture.download(t, sizeFixtureRun{adaptive: true})

	err := job.Statuses[0].Err()
	require.ErrorIs(t, err, errDownloadSizeConflict)
	requireNoOutput(t, dest)
}

func TestDownloadAdaptiveRangeWithoutByteCoverage(t *testing.T) {
	// A 206 whose Content-Range does not state which bytes it carries cannot be
	// trusted for the plan: before any byte is written the fallback engine takes
	// over; afterwards it contradicts the established plan and fails safely.
	t.Run("first response declines to the fallback engine", func(t *testing.T) {
		fixture := newSizeFixture(t, 21*sizeFixtureMiB, 21*sizeFixtureMiB)
		fixture.unsatisfiedContentRangeOnResponse = 1
		job, dest := fixture.download(t, sizeFixtureRun{adaptive: true})

		require.NoError(t, job.Statuses[0].Err())
		fixture.requireExactOutput(t, dest)
		assert.NotContains(t, fixture.ranges(), adaptivePartBoundary+"22020095", "the adaptive plan must not continue from an untrusted first response")
	})
	t.Run("later response fails the attempt", func(t *testing.T) {
		fixture := newSizeFixture(t, 21*sizeFixtureMiB, 21*sizeFixtureMiB)
		fixture.unsatisfiedContentRangeOnResponse = 2
		job, dest := fixture.download(t, sizeFixtureRun{adaptive: true})

		require.ErrorIs(t, job.Statuses[0].Err(), errDownloadSizeConflict)
		requireNoOutput(t, dest)
	})
}

func TestDownloadV2RangeFromResponse(t *testing.T) {
	response := func(contentRange string) *http.Response {
		header := http.Header{}
		if contentRange != "" {
			header.Set("Content-Range", contentRange)
		}
		return &http.Response{StatusCode: http.StatusPartialContent, Header: header, ContentLength: 1000}
	}
	trusted := downloadV2Range{off: 0, end: 999, total: 64000, trust: TrustedSizeValue}
	untrusted := downloadV2Range{off: 0, end: 999, trust: UntrustedSizeValue}
	for _, tc := range []struct {
		contentRange string
		want         downloadV2Range
	}{
		{"bytes 0-999/64000", trusted},
		{"0-999/64000", trusted},
		{"bytes */64000", untrusted},
		{"bytes 0-999/*", untrusted},
		{"items 0-999/64000", untrusted},
		{"", untrusted},
	} {
		assert.Equal(t, tc.want, downloadV2RangeFromResponse(response(tc.contentRange), 0, 999), tc.contentRange)
	}
}

func TestDownloadShortBodiesFail(t *testing.T) {
	for _, tc := range []struct {
		name string
		run  sizeFixtureRun
	}{
		{"adaptive range", sizeFixtureRun{adaptive: true}},
		{"parts range", sizeFixtureRun{}},
		{"single stream", sizeFixtureRun{singleStream: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newSizeFixture(t, 21*sizeFixtureMiB, 21*sizeFixtureMiB)
			fixture.shortBody = true
			job, dest := fixture.download(t, tc.run)

			require.Error(t, job.Statuses[0].Err())
			requireNoOutput(t, dest)
		})
	}
}

func TestDownloadSourceChangedRestartsFromNewDownloadRequest(t *testing.T) {
	changing := func(t *testing.T) *sizeFixture {
		fixture := newSizeFixture(t, 21*sizeFixtureMiB, 21*sizeFixtureMiB)
		fixture.changeAfterResponse = 1
		fixture.changeDelta = sizeFixtureMiB
		return fixture
	}

	t.Run("adaptive retry budget 0 fails safely without retrying ranges", func(t *testing.T) {
		fixture := changing(t)
		job, dest := fixture.download(t, sizeFixtureRun{adaptive: true})

		err := job.Statuses[0].Err()
		require.ErrorIs(t, err, files_sdk.ErrDownloadSourceChanged)
		requireNoOutput(t, dest)
		lifecycles, rejected := fixture.stats()
		assert.Equal(t, 1, lifecycles)
		require.NotEmpty(t, rejected)
		requireEachRangeRejectedOnce(t, rejected)
	})

	t.Run("adaptive retry budget 1 completes from a new download request despite stale metadata", func(t *testing.T) {
		fixture := changing(t)
		job, dest := fixture.download(t, sizeFixtureRun{adaptive: true, retryCount: 1})

		require.NoError(t, job.Statuses[0].Err())
		fixture.requireExactOutput(t, dest)
		assert.Equal(t, 22*sizeFixtureMiB, job.Statuses[0].Size())
		lifecycles, rejected := fixture.stats()
		assert.Equal(t, 2, lifecycles, "the retry must use a new download request")
		requireEachRangeRejectedOnce(t, rejected)
		assert.Contains(t, fixture.ranges(), adaptivePartBoundary+"23068671", "the adaptive engine planned the retry from the transfer size")
	})

	t.Run("adaptive source emptied after first response completes empty from a new download request", func(t *testing.T) {
		fixture := changing(t)
		fixture.changeDelta = -21 * sizeFixtureMiB
		job, dest := fixture.download(t, sizeFixtureRun{adaptive: true, retryCount: 1})

		require.NoError(t, job.Statuses[0].Err())
		fixture.requireExactOutput(t, dest)
		assert.Equal(t, int64(0), job.Statuses[0].Size())
		assert.Equal(t, int64(0), job.TransferBytes())
		lifecycles, rejected := fixture.stats()
		assert.Equal(t, 2, lifecycles)
		requireEachRangeRejectedOnce(t, rejected)
	})

	t.Run("adaptive completes when metadata tracks the source", func(t *testing.T) {
		fixture := changing(t)
		fixture.metadataTracksSource = true
		job, dest := fixture.download(t, sizeFixtureRun{adaptive: true, retryCount: 1})

		require.NoError(t, job.Statuses[0].Err())
		fixture.requireExactOutput(t, dest)
		lifecycles, _ := fixture.stats()
		assert.Equal(t, 2, lifecycles)
	})

	t.Run("parts engine restarts then converges on the transfer size", func(t *testing.T) {
		fixture := changing(t)
		job, dest := fixture.download(t, sizeFixtureRun{retryCount: 2})

		require.NoError(t, job.Statuses[0].Err())
		fixture.requireExactOutput(t, dest)
		lifecycles, rejected := fixture.stats()
		assert.Equal(t, 2, lifecycles, "only the source change replaces the download request")
		requireEachRangeRejectedOnce(t, rejected)
	})

	t.Run("single stream rejected before bytes restarts", func(t *testing.T) {
		fixture := newSizeFixture(t, 21*sizeFixtureMiB, 21*sizeFixtureMiB)
		fixture.changeBeforeFirstTransfer = true
		fixture.changeDelta = -20 * sizeFixtureMiB
		job, dest := fixture.download(t, sizeFixtureRun{singleStream: true, retryCount: 1})

		require.NoError(t, job.Statuses[0].Err())
		fixture.requireExactOutput(t, dest)
		lifecycles, rejected := fixture.stats()
		assert.Equal(t, 2, lifecycles)
		assert.Len(t, rejected, 1)
	})

	t.Run("single stream retry budget 0 fails safely", func(t *testing.T) {
		fixture := newSizeFixture(t, 21*sizeFixtureMiB, 21*sizeFixtureMiB)
		fixture.changeBeforeFirstTransfer = true
		fixture.changeDelta = -20 * sizeFixtureMiB
		job, dest := fixture.download(t, sizeFixtureRun{singleStream: true})

		require.ErrorIs(t, job.Statuses[0].Err(), files_sdk.ErrDownloadSourceChanged)
		requireNoOutput(t, dest)
	})

	t.Run("repeatedly changing source exhausts the retry budget", func(t *testing.T) {
		fixture := changing(t)
		fixture.repeatChange = true
		job, dest := fixture.download(t, sizeFixtureRun{adaptive: true, retryCount: 2})

		require.ErrorIs(t, job.Statuses[0].Err(), files_sdk.ErrDownloadSourceChanged)
		requireNoOutput(t, dest)
		lifecycles, rejected := fixture.stats()
		assert.Equal(t, 3, lifecycles, "one download request per attempt, no refresh loop")
		requireEachRangeRejectedOnce(t, rejected)
	})
}

func TestDownloadUntypedConflictsKeepTheDownloadRequest(t *testing.T) {
	// Only the typed download_source_changed error starts a new lifecycle; other
	// 409s, even ones mentioning the token, keep the ordinary retry behavior.
	for _, shape := range []string{"text", "other"} {
		t.Run(shape, func(t *testing.T) {
			fixture := newSizeFixture(t, 21*sizeFixtureMiB, 21*sizeFixtureMiB)
			fixture.changeAfterResponse = 1
			fixture.changeDelta = sizeFixtureMiB
			fixture.rejectAs = shape
			job, dest := fixture.download(t, sizeFixtureRun{adaptive: true, retryCount: 1})

			err := job.Statuses[0].Err()
			require.Error(t, err)
			assert.NotErrorIs(t, err, files_sdk.ErrDownloadSourceChanged)
			requireNoOutput(t, dest)
			lifecycles, rejected := fixture.stats()
			assert.Equal(t, 1, lifecycles, "an untyped conflict must not replace the download request")
			assert.NotEmpty(t, rejected)
		})
	}
}

func TestDownloadSourceChangedRequiresTypedError(t *testing.T) {
	typed := files_sdk.ResponseError{Type: string(files_sdk.ErrDownloadSourceChanged), ErrorMessage: "changed", HttpCode: http.StatusConflict}
	assert.True(t, downloadSourceChanged(&fs.PathError{Op: "ReaderRange", Path: "a", Err: typed}))
	assert.True(t, downloadSourceChanged(fmt.Errorf("part 2: %w", typed)))

	assert.False(t, downloadSourceChanged(nil))
	assert.False(t, downloadSourceChanged(files_sdk.ResponseError{Type: "processing-failure/conflict", ErrorMessage: "download_source_changed mentioned", HttpCode: http.StatusConflict}))
	plainText := lib.NotStatus(http.StatusPartialContent)(&http.Response{
		StatusCode:    http.StatusConflict,
		Header:        http.Header{"Content-Type": {"text/plain"}},
		ContentLength: -1,
		Body:          io.NopCloser(strings.NewReader("download_source_changed")),
	})
	require.Error(t, plainText)
	assert.False(t, downloadSourceChanged(plainText))
}

func TestDownloadV2UntrustedRangeResponseUsesFallback(t *testing.T) {
	root := t.TempDir()
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	client := server.Client()
	size := int64(20 * sizeFixtureMiB)
	server.MockFiles["untrusted.bin"] = mockFile{SizeTrust: UntrustedSizeValue, File: files_sdk.File{Size: size}}

	job := client.Downloader(DownloaderParams{
		RemotePath:          "untrusted.bin",
		LocalPath:           root + "/",
		Manager:             manager.Build(2, 1),
		AdaptiveConcurrency: true,
	})
	job.Start()
	job.Wait()

	require.Len(t, job.Statuses, 1)
	require.NoError(t, job.Statuses[0].Err())
	stat, err := os.Stat(filepath.Join(root, "untrusted.bin"))
	require.NoError(t, err)
	assert.Equal(t, size, stat.Size())
}

func TestDownloadV2CancelWhileWaitingForFirstRange(t *testing.T) {
	size := int64(20 * sizeFixtureMiB)
	ranger := &downloadV2BlockingRangeFile{
		info: Info{File: files_sdk.File{
			DisplayName: "native.bin",
			Path:        "native.bin",
			Type:        "file",
			Size:        size,
		}, sizeTrust: TrustedSizeValue},
		ctx: context.Background(),
	}
	tmpPath := filepath.Join(t.TempDir(), "native.bin.download")
	reportStatus := downloadV2TestStatus(ranger, ranger.info, DownloaderParams{
		AdaptiveConcurrency: true,
		Manager:             manager.Build(2, 1),
	}, tmpPath)
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(100*time.Millisecond, cancel)

	done := make(chan error, 1)
	go func() {
		_, _, _, err := runDownloadV2IfSupported(ctx, reportStatus, ranger.info, explicitDestination(tmpPath), 0)
		done <- err
	}()
	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(5 * time.Second):
		t.Fatal("download v2 did not stop after cancellation while waiting for the first range")
	}
}

func TestDownloadV2PreallocationFailureStopsBeforeAdmittingRanges(t *testing.T) {
	size := int64(20 * sizeFixtureMiB)
	ranger := &downloadV2TestRangeFile{
		data:        make([]byte, size),
		downloadURI: "https://files.example.com/download/native.bin",
		info: Info{File: files_sdk.File{
			DisplayName: "native.bin",
			Path:        "native.bin",
			Type:        "file",
			Size:        size,
		}, sizeTrust: TrustedSizeValue},
	}
	tmpPath := filepath.Join(t.TempDir(), "native.bin.download")
	params := DownloaderParams{AdaptiveConcurrency: true, Manager: manager.Build(2, 1)}
	reportStatus := downloadV2TestStatus(ranger, ranger.info, params, tmpPath)
	// A read-only output file makes preallocation fail after the first range
	// response has been accepted.
	require.NoError(t, os.WriteFile(tmpPath, nil, 0644))
	file, err := os.Open(tmpPath)
	require.NoError(t, err)
	engine := newDownloadV2Engine(reportStatus, ranger, file, downloadV2TargetDefault, size, 0, downloadV2KnownSizePartSize(downloadV2TargetDefault, size), params)

	done := make(chan error, 1)
	go func() { done <- engine.Run(context.Background()) }()
	select {
	case err := <-done:
		require.Error(t, err)
		var terminal downloadV2TerminalError
		assert.ErrorAs(t, err, &terminal)
	case <-time.After(5 * time.Second):
		t.Fatal("download v2 did not stop after the preallocation failure")
	}
	assert.Len(t, ranger.Ranges(), 1, "no further range may be admitted after a failed preallocation")
	stat, err := os.Stat(tmpPath)
	require.NoError(t, err)
	assert.Equal(t, int64(0), stat.Size())
}

// downloadV2BlockingRangeFile never answers a range request until its context
// ends, like a stalled first range.
type downloadV2BlockingRangeFile struct {
	info Info
	ctx  context.Context
}

func (f *downloadV2BlockingRangeFile) Stat() (fs.FileInfo, error) { return f.info, nil }
func (f *downloadV2BlockingRangeFile) Read([]byte) (int, error)   { return 0, io.EOF }
func (f *downloadV2BlockingRangeFile) Close() error               { return nil }

func (f *downloadV2BlockingRangeFile) downloadV2URI(context.Context) (string, error) {
	return "https://files.example.com/download/native.bin", nil
}

func (f *downloadV2BlockingRangeFile) WithContext(ctx context.Context) fs.File {
	return &downloadV2BlockingRangeFile{info: f.info, ctx: ctx}
}

func (f *downloadV2BlockingRangeFile) ReaderRange(int64, int64) (io.ReadCloser, error) {
	<-f.ctx.Done()
	return nil, f.ctx.Err()
}

func TestRetryByStatusStopsBackoffOnCancel(t *testing.T) {
	job := (&Job{Direction: direction.DownloadType, Manager: manager.Default(), Config: files_sdk.Config{}.Init(), Logger: lib.NullLogger{}}).Init()
	job.Add(&DownloadStatus{Mutex: &sync.RWMutex{}, status: status.Errored, error: errors.New("boom"), job: job, file: files_sdk.File{DisplayName: "a.txt"}})
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(100*time.Millisecond, cancel)

	start := time.Now()
	RetryByPolicy(ctx, job, RetryPolicy{Type: RetryUnfinished, RetryCount: 3, Backoff: 3}, false)

	assert.Less(t, time.Since(start), 2*time.Second, "backoff must stop when the context is canceled")
	assert.Equal(t, 1, job.Count(status.Errored))
}

func TestParseContentRange(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  contentRange
		ok    bool
	}{
		{"bytes 0-999/64000", contentRange{start: 0, end: 999, total: 64000, totalKnown: true}, true},
		{"0-999/64000", contentRange{start: 0, end: 999, total: 64000, totalKnown: true}, true},
		{"bytes 500-999/*", contentRange{start: 500, end: 999}, true},
		{"BYTES 0-999/64000", contentRange{start: 0, end: 999, total: 64000, totalKnown: true}, true},
		{"bytes */64000", contentRange{start: -1, end: -1, total: 64000, totalKnown: true}, true},
		{"items 0-999/64000", contentRange{}, false},
		{"bytes 0-999", contentRange{}, false},
		{"bytes 999-0/64000", contentRange{}, false},
		{"", contentRange{}, false},
	} {
		got, ok := parseContentRange(tc.value)
		assert.Equal(t, tc.ok, ok, tc.value)
		assert.Equal(t, tc.want, got, tc.value)
	}
}
