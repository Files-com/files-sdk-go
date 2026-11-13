package file

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/require"
)

type versionDownloadRequest struct{ uri, ranges, match, unmodifiedSince string }

// boundStageCheckpoint is fixture data from a settled response: the stage's
// bytes belong to this exact source version and destination.
func boundStageCheckpoint(config files_sdk.Config, remote, local, tmp, etag string, total, prefix int64) *JobDownloadCheckpoint {
	return &JobDownloadCheckpoint{Bindings: []DownloadBinding{{
		Source:     downloadSource(&File{FS: &FS{Config: config}}),
		RemotePath: remote,
		LocalPath:  checkpointLocalPath(local),
		TmpPath:    checkpointLocalPath(tmp),
		ETag:       etag,
		Total:      total,
		Prefix:     prefix,
	}}}
}

func waitForDownloadTestJob(t *testing.T, ctx context.Context, job *Job) {
	t.Helper()
	settled := make(chan struct{})
	go func() { job.Wait(); close(settled) }()
	select {
	case <-settled:
		require.NoError(t, ctx.Err(), "download test exceeded its deadline")
	case <-ctx.Done():
		job.Cancel()
		select {
		case <-settled:
		case <-time.After(5 * time.Second):
			t.Fatal("canceled download test did not settle within five seconds")
		}
		t.Fatalf("download test deadline exceeded: %v", ctx.Err())
	}
}

type versionDownloadFixture struct {
	beforeStart                                   func()
	mu                                            sync.Mutex
	url                                           string
	data                                          []byte
	requests                                      []versionDownloadRequest
	failed, expire, expired, changed, noValidator bool
	alwaysShort                                   bool
	ignoresRange                                  bool
	lastModified, laterLastModified               string
	ignoreDateCondition, expireSecond             bool
	omitLaterDate                                 bool
	duplicateDate                                 bool
	clientDateCheck, localDateCheck               *bool
}

func newVersionDownloadFixture(t *testing.T) *versionDownloadFixture {
	f := &versionDownloadFixture{data: bytes.Repeat([]byte("a varied download payload\x00\xff"), (34<<20)/26+1)}
	server := httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(server.Close)
	f.url = server.URL
	return f
}

func (f *versionDownloadFixture) serve(w http.ResponseWriter, req *http.Request) {
	f.mu.Lock()
	if strings.HasPrefix(req.URL.Path, "/api/") {
		uri := f.url + "/transfer"
		if f.expired {
			uri = f.url + "/fresh"
		}
		f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"path": "remote.bin", "display_name": "remote.bin", "type": "file", "size": len(f.data), "download_uri": uri})
		return
	}
	ranges := req.Header.Get("Range")
	f.requests = append(f.requests, versionDownloadRequest{uri: req.URL.Path, ranges: ranges, match: req.Header.Get("If-Match"), unmodifiedSince: req.Header.Get("If-Unmodified-Since")})
	if f.expireSecond && len(f.requests) == 2 {
		f.expired = true
	}
	if f.ignoresRange {
		ranges = ""
	}
	if f.expired && req.URL.Path != "/fresh" {
		f.mu.Unlock()
		w.WriteHeader(403)
		return
	}
	start, end := int64(0), int64(len(f.data)-1)
	if ranges != "" {
		_, _ = fmt.Sscanf(ranges, "bytes=%d-%d", &start, &end)
		end = min(end, int64(len(f.data)-1))
	}
	drop := ranges != "" && !f.noValidator && (f.alwaysShort || start == 0 && !f.failed)
	if drop {
		f.failed = true
		f.expired = f.expire
	}
	etag := `"original"`
	if f.changed && start > 0 {
		etag = `"replacement"`
	}
	noValidator := f.noValidator
	modified := f.lastModified
	if len(f.requests) > 1 && f.laterLastModified != "" {
		modified = f.laterLastModified
	}
	if len(f.requests) > 1 && f.omitLaterDate {
		modified = ""
	}
	condition, conditionErr := http.ParseTime(req.Header.Get("If-Unmodified-Since"))
	responseDate, responseDateErr := http.ParseTime(modified)
	rejectDate := !f.ignoreDateCondition && req.Header.Get("If-Match") == "" && conditionErr == nil && responseDateErr == nil && responseDate.After(condition)
	f.mu.Unlock()
	if modified != "" {
		w.Header().Set("Last-Modified", modified)
		if f.duplicateDate {
			w.Header().Add("Last-Modified", modified)
		}
	}
	if rejectDate {
		w.WriteHeader(http.StatusPreconditionFailed)
		return
	}
	if !noValidator {
		w.Header().Set("ETag", etag)
	}
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	if ranges != "" {
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(f.data)))
		w.WriteHeader(206)
	}
	if drop {
		end = start + (1 << 20) - 1
	}
	_, _ = w.Write(f.data[start : end+1])
}

func (f *versionDownloadFixture) client() *Client {
	return &Client{Config: files_sdk.Config{EndpointOverride: f.url, DisableDirectTransfers: true, CheckDownloadModificationTime: f.clientDateCheck}.Init()}
}

func (f *versionDownloadFixture) run(t *testing.T, adaptive bool) (*Job, string) {
	dest := filepath.Join(t.TempDir(), "out.bin")
	job := f.client().Downloader(DownloaderParams{RemotePath: "remote.bin", LocalPath: dest, Manager: manager.Build(2, 1, false), AdaptiveConcurrency: adaptive, CheckDownloadModificationTime: f.localDateCheck})
	if f.beforeStart != nil {
		f.beforeStart()
	}
	job.Start()
	job.Wait()
	require.Len(t, job.Statuses, 1)
	return job, dest
}

func TestDownloaderResumesTheFirstUnreceivedByteWithItsVersion(t *testing.T) {
	for _, adaptive := range []bool{false, true} {
		for _, renew := range []bool{false, true} {
			t.Run(fmt.Sprintf("adaptive=%v/renew=%v", adaptive, renew), func(t *testing.T) {
				f := newVersionDownloadFixture(t)
				f.expire = renew
				f.lastModified = "Wed, 21 Oct 2015 07:28:00 GMT"
				f.laterLastModified = "Wed, 21 Oct 2015 07:29:00 GMT"
				job, dest := f.run(t, adaptive)
				require.NoError(t, job.Statuses[0].Err())
				out, err := os.ReadFile(dest)
				require.NoError(t, err)
				require.Equal(t, sha256.Sum256(f.data), sha256.Sum256(out))
				f.mu.Lock()
				requests := append([]versionDownloadRequest(nil), f.requests...)
				f.mu.Unlock()
				resumed := false
				for _, req := range requests {
					if req.match != "" {
						require.Empty(t, req.unmodifiedSince, "If-Match takes precedence over the date guard")
					}
					if strings.HasPrefix(req.ranges, "bytes=1048576-") && req.match == `"original"` && (!renew || req.uri == "/fresh") {
						resumed = true
					}
				}
				require.True(t, resumed, "requests: %v", requests)
			})
		}
	}
}

func TestDownloaderRejectsNewerDateBeforePublishingAWholeResponse(t *testing.T) {
	for _, adaptive := range []bool{false, true} {
		for _, ignoreCondition := range []bool{false, true} {
			t.Run(fmt.Sprintf("adaptive=%v/ignores-condition=%v", adaptive, ignoreCondition), func(t *testing.T) {
				f := newVersionDownloadFixture(t)
				f.noValidator = true
				f.lastModified = "Wed, 21 Oct 2015 07:28:00 GMT"
				f.laterLastModified = "Wed, 21 Oct 2015 07:28:01 GMT"
				f.ignoreDateCondition = ignoreCondition
				job, dest := f.run(t, adaptive)
				require.ErrorIs(t, job.Statuses[0].Err(), files_sdk.ErrDownloadSourceChanged)
				_, err := os.Stat(dest)
				require.True(t, os.IsNotExist(err))
				f.mu.Lock()
				defer f.mu.Unlock()
				require.Len(t, f.requests, 2, "a timestamp does not qualify the source for parallel ranges")
				require.Empty(t, f.requests[0].unmodifiedSince)
				require.Empty(t, f.requests[1].ranges)
				require.Equal(t, f.lastModified, f.requests[1].unmodifiedSince)
			})
		}
	}
}

func TestDownloadDateGuardKeepsTheFirstValidDateAcrossRanges(t *testing.T) {
	for _, tc := range []struct {
		name, initial, later, condition string
		omitLater, duplicate            bool
	}{
		{name: "stable", initial: "Wed, 21 Oct 2015 07:28:00 GMT", condition: "Wed, 21 Oct 2015 07:28:00 GMT"},
		{name: "earlier", initial: "Wed, 21 Oct 2015 07:28:00 GMT", later: "Wed, 21 Oct 2015 07:27:59 GMT", condition: "Wed, 21 Oct 2015 07:28:00 GMT"},
		{name: "obsolete-http-date", initial: "Wednesday, 21-Oct-15 07:28:00 GMT", later: "Wed, 21 Oct 2015 07:28:00 GMT", condition: "Wed, 21 Oct 2015 07:28:00 GMT"},
		{name: "missing"},
		{name: "duplicate", initial: "Wed, 21 Oct 2015 07:28:00 GMT", duplicate: true},
		{name: "invalid", initial: "not an HTTP date"},
		{name: "later-missing", initial: "Wed, 21 Oct 2015 07:28:00 GMT", condition: "Wed, 21 Oct 2015 07:28:00 GMT", omitLater: true},
		{name: "later-invalid", initial: "Wed, 21 Oct 2015 07:28:00 GMT", later: "not an HTTP date", condition: "Wed, 21 Oct 2015 07:28:00 GMT"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newVersionDownloadFixture(t)
			f.data, f.noValidator = []byte("first part then second part"), true
			f.lastModified, f.laterLastModified = tc.initial, tc.later
			f.omitLaterDate, f.duplicateDate = tc.omitLater, tc.duplicate
			remote := (&File{File: &files_sdk.File{Path: "remote.bin", DownloadUri: f.url + "/transfer"}, FS: (&FS{Context: context.Background()}).Init(f.client().Config, false)}).Init()
			for offset := int64(0); offset < 3; offset++ {
				reader, err := remote.ReaderRange(offset, offset)
				require.NoError(t, err)
				data, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.NoError(t, reader.Close())
				require.Equal(t, f.data[offset:offset+1], data)
			}
			f.mu.Lock()
			defer f.mu.Unlock()
			for _, req := range f.requests[1:] {
				require.Equal(t, tc.condition, req.unmodifiedSince)
				require.Empty(t, req.match)
			}
		})
	}
}

func TestDownloadDateGuardSurvivesURLRenewal(t *testing.T) {
	f := newVersionDownloadFixture(t)
	f.data, f.noValidator = []byte("date-bound range data"), true
	f.lastModified, f.expireSecond = "Wed, 21 Oct 2015 07:28:00 GMT", true
	remote := (&File{File: &files_sdk.File{Path: "remote.bin", DownloadUri: f.url + "/transfer"}, FS: (&FS{Context: context.Background()}).Init(f.client().Config, false)}).Init()
	for offset := int64(0); offset < 2; offset++ {
		reader, err := remote.ReaderRange(offset, offset)
		require.NoError(t, err)
		_, err = io.Copy(io.Discard, reader)
		require.NoError(t, err)
		require.NoError(t, reader.Close())
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	require.Len(t, f.requests, 3)
	require.Equal(t, "/fresh", f.requests[2].uri)
	for _, req := range f.requests[1:] {
		require.Equal(t, f.lastModified, req.unmodifiedSince)
	}
}

func TestDownloaderModificationTimeCheckOptions(t *testing.T) {
	for _, tc := range []struct {
		name          string
		client, local *bool
		checks        bool
		mutate        bool
	}{
		{name: "explicit-true", client: lib.Bool(true), checks: true},
		{name: "client-disabled", client: lib.Bool(false)},
		{name: "job-disabled", client: lib.Bool(true), local: lib.Bool(false)},
		{name: "job-enabled", client: lib.Bool(false), local: lib.Bool(true), checks: true},
		{name: "client-snapshot", client: lib.Bool(true), checks: true, mutate: true},
		{name: "job-snapshot", local: lib.Bool(true), checks: true, mutate: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newVersionDownloadFixture(t)
			f.noValidator, f.ignoreDateCondition = true, true
			f.lastModified, f.laterLastModified = "Wed, 21 Oct 2015 07:28:00 GMT", "Wed, 21 Oct 2015 07:28:01 GMT"
			f.clientDateCheck, f.localDateCheck = tc.client, tc.local
			if tc.mutate {
				f.beforeStart = func() {
					if tc.local != nil {
						*tc.local = false
					} else {
						*tc.client = false
					}
				}
			}
			job, dest := f.run(t, true)
			if tc.checks {
				require.ErrorIs(t, job.Statuses[0].Err(), files_sdk.ErrDownloadSourceChanged)
			} else {
				require.NoError(t, job.Statuses[0].Err())
				data, err := os.ReadFile(dest)
				require.NoError(t, err)
				require.Equal(t, f.data, data)
			}
			f.mu.Lock()
			defer f.mu.Unlock()
			condition := ""
			if tc.checks {
				condition = f.lastModified
			}
			require.Len(t, f.requests, 2)
			require.Equal(t, condition, f.requests[1].unmodifiedSince)
			require.Empty(t, f.requests[1].ranges, "date opt-out does not authorize parallel ranges")
		})
	}
}

func TestDownloadDateGuardSnapshotsTheReadHandleSetting(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprint(enabled), func(t *testing.T) {
			f := newVersionDownloadFixture(t)
			f.data, f.noValidator, f.ignoreDateCondition = []byte("date guard"), true, true
			f.lastModified, f.laterLastModified = "Wed, 21 Oct 2015 07:28:00 GMT", "Wed, 21 Oct 2015 07:28:01 GMT"
			setting := enabled
			f.clientDateCheck = &setting
			remote := (&File{File: &files_sdk.File{Path: "remote.bin", DownloadUri: f.url + "/transfer"}, FS: (&FS{Context: context.Background()}).Init(f.client().Config, false)}).Init()
			reader, err := remote.ReaderRange(0, 0)
			require.NoError(t, err)
			_, err = io.Copy(io.Discard, reader)
			require.NoError(t, err)
			require.NoError(t, reader.Close())
			setting = !enabled
			reader, err = remote.ReaderRange(1, 1)
			if enabled {
				require.ErrorIs(t, err, files_sdk.ErrDownloadSourceChanged)
			} else {
				require.NoError(t, err)
				_, err = io.Copy(io.Discard, reader)
				require.NoError(t, err)
				require.NoError(t, reader.Close())
			}
			f.mu.Lock()
			defer f.mu.Unlock()
			condition := ""
			if enabled {
				condition = f.lastModified
			}
			require.Equal(t, condition, f.requests[1].unmodifiedSince)
		})
	}
}

func TestDownloadDateOptOutKeepsStrongVersionAndHTTPErrors(t *testing.T) {
	f := newVersionDownloadFixture(t)
	f.data, f.failed, f.clientDateCheck = []byte("version-bound"), true, lib.Bool(false)
	remote := (&File{File: &files_sdk.File{Path: "remote.bin", DownloadUri: f.url + "/transfer"}, FS: (&FS{Context: context.Background()}).Init(f.client().Config, false)}).Init()
	reader, err := remote.ReaderRange(0, 0)
	require.NoError(t, err)
	_, err = io.Copy(io.Discard, reader)
	require.NoError(t, err)
	require.NoError(t, reader.Close())
	f.mu.Lock()
	f.changed = true
	f.mu.Unlock()
	_, err = remote.ReaderRange(1, 1)
	require.ErrorIs(t, err, files_sdk.ErrDownloadSourceChanged)
	f.mu.Lock()
	require.Equal(t, `"original"`, f.requests[1].match)
	require.Empty(t, f.requests[1].unmodifiedSince)
	f.mu.Unlock()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusPreconditionFailed) }))
	defer server.Close()
	remote = (&File{File: &files_sdk.File{Path: "remote.bin", DownloadUri: server.URL}, FS: (&FS{Context: context.Background()}).Init(f.client().Config, false)}).Init()
	_, err = remote.ReaderRange(0, 0)
	require.ErrorIs(t, err, files_sdk.ErrDownloadSourceChanged, "date opt-out does not make HTTP412 successful")
}

func TestDownloaderRejectsSameSizeRangeReplacement(t *testing.T) {
	for _, adaptive := range []bool{false, true} {
		t.Run(fmt.Sprint(adaptive), func(t *testing.T) {
			f := newVersionDownloadFixture(t)
			f.changed = true
			job, dest := f.run(t, adaptive)
			require.ErrorIs(t, job.Statuses[0].Err(), files_sdk.ErrDownloadSourceChanged)
			_, err := os.Stat(dest)
			require.True(t, os.IsNotExist(err))
		})
	}
}

func TestVersionBoundRangeContinuationUsesOneFiniteAttemptAllowance(t *testing.T) {
	for _, attempts := range []int{downloadV2RetryAttempts, legacyDownloadRangeAttempts} {
		t.Run(fmt.Sprint(attempts), func(t *testing.T) {
			f := newVersionDownloadFixture(t)
			f.alwaysShort = true
			remote := (&File{File: &files_sdk.File{Path: "remote.bin", DownloadUri: f.url + "/transfer"}, FS: (&FS{Context: context.Background()}).Init(f.client().Config, false)}).Init()
			reader, _, err := remote.openDownloadRange(context.Background(), 0, (6<<20)-1, attempts)
			require.NoError(t, err)
			_, err = io.Copy(io.Discard, reader)
			require.ErrorIs(t, err, io.ErrUnexpectedEOF)
			_ = reader.Close()
			f.mu.Lock()
			defer f.mu.Unlock()
			require.Len(t, f.requests, attempts)
			for index, request := range f.requests {
				require.Equal(t, fmt.Sprintf("bytes=%d-%d", index*(1<<20), (6<<20)-1), request.ranges)
			}
		})
	}
}

func TestDownloaderUsesAWholeResponseWithoutAVersionProof(t *testing.T) {
	f := newVersionDownloadFixture(t)
	f.noValidator = true
	job, dest := f.run(t, true)
	require.NoError(t, job.Statuses[0].Err())
	out, err := os.ReadFile(dest)
	require.NoError(t, err)
	require.Equal(t, sha256.Sum256(f.data), sha256.Sum256(out))
	f.mu.Lock()
	defer f.mu.Unlock()
	require.Len(t, f.requests, 2)
	require.NotEmpty(t, f.requests[0].ranges)
	require.Empty(t, f.requests[1].ranges)
}

func TestDownloadRejectsFrontendBeforeRetryRedirectAndBodyCallbacks(t *testing.T) {
	for _, code := range []int{200, 302, 503} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			body := &ownedHTTPResponseBody{Reader: strings.NewReader("frontend page")}
			config := ownedResponseConfig(body, code, http.Header{"X-Files-Frontend-App": nil, "Location": []string{"https://download.test/other"}})
			var callbacks, requests int
			original := config.Client.HTTPClient.Transport
			config.Client.HTTPClient.Transport = ownedHTTPResponseTransport(func(req *http.Request) (*http.Response, error) { requests++; return original.RoundTrip(req) })
			config.Client.CheckRetry = func(_ context.Context, res *http.Response, _ error) (bool, error) {
				if res != nil {
					_, _ = io.ReadAll(res.Body)
				}
				return true, nil
			}
			config.Client.ErrorHandler = func(_ *http.Response, _ error, _ int) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("masked"))}, nil
			}
			_, err := (&Client{Config: config}).Download(files_sdk.FileDownloadParams{File: files_sdk.File{Path: "file", DownloadUri: "https://download.test/file"}}, files_sdk.ResponseBodyOption(func(io.ReadCloser) error { callbacks++; return nil }))
			require.ErrorIs(t, err, errDownloadFrontend)
			require.Equal(t, 0, callbacks)
			require.Equal(t, 1, requests)
			require.Equal(t, 1, body.closes)
		})
	}
}

func TestDownloadRenewsWhenRetryAfterCrossesExpiry(t *testing.T) {
	var stale, fresh, metadata, prepare atomic.Int64
	var url string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/api/"):
			metadata.Add(1)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"path": "file", "download_uri": url + "/fresh"})
		case r.URL.Path == "/stale":
			stale.Add(1)
			w.Header().Set("Retry-After", "3")
			w.WriteHeader(503)
		default:
			fresh.Add(1)
			_, _ = io.WriteString(w, "file bytes")
		}
	}))
	defer server.Close()
	url = server.URL
	config := files_sdk.Config{EndpointOverride: url, DisableDirectTransfers: true}.Init().SetCustomClient(server.Client())
	config.PrepareRetry = func(*http.Request) error { prepare.Add(1); return nil }
	client, transport := config.Client, config.HTTPClient.Transport
	expires := time.Now().Add(2 * time.Second).UTC().Format("20060102T150405Z")
	var output []byte
	_, err := (&Client{Config: config}).Download(files_sdk.FileDownloadParams{File: files_sdk.File{Path: "file", DownloadUri: url + "/stale?X-Files-Date=" + expires}}, files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
		defer body.Close()
		var err error
		output, err = io.ReadAll(body)
		return err
	}))
	require.NoError(t, err)
	require.Equal(t, "file bytes", string(output))
	require.EqualValues(t, 1, stale.Load())
	require.EqualValues(t, 1, fresh.Load())
	require.EqualValues(t, 1, metadata.Load())
	require.EqualValues(t, 1, prepare.Load())
	require.Same(t, client, config.Client)
	require.Equal(t, transport, config.HTTPClient.Transport)
}

func TestDownloadFrontendDirectResponseDoesNotFallBackToProxy(t *testing.T) {
	// An empty header must still be carried over HTTP.
	info := testDirectTransferInfo(t, "/download", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Files-Frontend-App", "")
		_, _ = io.WriteString(w, "frontend")
	}))
	var proxyCalls, callbacks int
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { proxyCalls++; _, _ = io.WriteString(w, "proxy") }))
	defer proxy.Close()
	client := &Client{Config: files_sdk.Config{EndpointOverride: proxy.URL, Environment: files_sdk.Development}.Init()}
	_, err := client.Download(files_sdk.FileDownloadParams{File: files_sdk.File{Path: "file", DownloadUri: proxy.URL + "/download", DirectConnectionInfo: info}}, files_sdk.ResponseBodyOption(func(io.ReadCloser) error { callbacks++; return nil }))
	require.ErrorIs(t, err, errDownloadFrontend)
	require.Zero(t, proxyCalls)
	require.Zero(t, callbacks)
}

func TestFileReadRenewalStopsAfterOneForbiddenResponseRefresh(t *testing.T) {
	for _, operation := range []string{"range", "read-at"} {
		t.Run(operation, func(t *testing.T) {
			var transfers, metadata int
			config := files_sdk.Config{EndpointOverride: "https://api.test", DisableDirectTransfers: true}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(req *http.Request) (*http.Response, error) {
				code, payload, contentType := 403, "forbidden", "text/plain"
				if req.URL.Host == "api.test" {
					metadata++
					code, payload, contentType = 200, `{"path":"file","download_uri":"https://download.test/fresh"}`, "application/json"
				} else {
					transfers++
				}
				return &http.Response{StatusCode: code, Header: http.Header{"Content-Type": []string{contentType}}, Body: io.NopCloser(strings.NewReader(payload)), ContentLength: int64(len(payload)), Request: req}, nil
			})})
			file := (&File{File: &files_sdk.File{Path: "file", DownloadUri: "https://download.test/stale"}, FS: (&FS{Context: context.Background()}).Init(config, false)}).Init()
			var err error
			if operation == "range" {
				_, err = file.ReaderRange(0, 2)
			} else {
				_, err = file.ReadAt(make([]byte, 3), 0)
			}
			require.Error(t, err)
			require.Equal(t, 2, transfers)
			require.Equal(t, 1, metadata)
		})
	}
}

// transferRefusal is how a transfer URL is refused: its status, content type and body.
type transferRefusal struct {
	name, contentType, body string
	status                  int
}

// How storage, and Files.com for a URL it serves, refuse a signed URL that has expired, whatever the URL's own expiry.
var expiredURLRefusals = []transferRefusal{
	{"s3", "application/xml", `<Error><Code>AccessDenied</Code><Message>Request has expired</Message></Error>`, http.StatusForbidden},
	{"expired token", "application/xml", `<Error><Code>ExpiredToken</Code><Message>The provided token has expired.</Message></Error>`, http.StatusBadRequest},
	{"azure", "application/xml", `<Error><Code>AuthenticationFailed</Code><Message>Server failed to authenticate the request.</Message><AuthenticationErrorDetail>Signature not valid in the specified time frame</AuthenticationErrorDetail></Error>`, http.StatusForbidden},
	{"download request", "application/json", `{"type":"download_request_expired","error":"Download request expired"}`, http.StatusBadRequest},
	{"upload request", "application/json", `{"type":"upload_request_expired","error":"Upload request expired"}`, http.StatusBadRequest},
}

// A read renews its URL once when storage's refusal says that URL has expired. Any other refusal fails the read without
// a renewal, and a successful body that reads like a refusal is the file's content.
func TestFileReadRenewsOnlyWhenStorageSaysItsURLExpired(t *testing.T) {
	content := `<Error><Code>ExpiredToken</Code></Error>`
	typed := `{"type":"download_request_expired","error":"Download request expired"}`
	refusals := append([]transferRefusal{
		// Its first 4096 bytes parse as an expired type, but the whole body is not JSON.
		{"json cut short", "application/json", typed + strings.Repeat(" ", 4096-len(typed)) + `}`, http.StatusBadRequest},
		{name: "none"},
		{"expired word", "application/xml", `<Error><Code>InvalidArgument</Code><Message>X-Amz-Expires must be less than a week</Message></Error>`, http.StatusBadRequest},
		{"structured error", "application/json", `{"type":"bad-request/invalid-body","error":"Request has expired"}`, http.StatusBadRequest},
		// A JSON body says an expiry only by its type, whatever words it carries.
		{"json other type", "application/json", `{"type":"not_authorized","error":"Request has expired"}`, http.StatusBadRequest},
		{"json other data", "application/json", `{"type":"other","data":"<Code>ExpiredToken</Code>"}`, http.StatusBadRequest},
		{"json array", "application/json", `[{"Message":"Request has expired"}]`, http.StatusBadRequest},
		{"json string", "application/json", `"Request has expired"`, http.StatusBadRequest},
		{"malformed json object", "application/json", `{"error":"Request has expired"`, http.StatusBadRequest},
		{"malformed json array", "application/json", `[{"Message":"<Code>ExpiredToken</Code>"`, http.StatusBadRequest},
	}, expiredURLRefusals...)
	for _, operation := range []string{"read", "range"} {
		for _, refusal := range refusals {
			t.Run(operation+"/"+refusal.name, func(t *testing.T) {
				var stale, fresh, metadata int
				answer := func(req *http.Request, code int, contentType, payload string) *http.Response {
					header := http.Header{"Content-Type": []string{contentType}}
					if code == http.StatusPartialContent {
						header.Set("Content-Range", fmt.Sprintf("bytes 0-%d/%d", len(payload)-1, len(payload)))
						header.Set("ETag", `"v1"`)
					}
					return &http.Response{StatusCode: code, Header: header, Body: io.NopCloser(strings.NewReader(payload)), ContentLength: int64(len(payload)), Request: req}
				}
				config := files_sdk.Config{EndpointOverride: "https://api.test", DisableDirectTransfers: true}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(req *http.Request) (*http.Response, error) {
					switch {
					case req.URL.Host == "api.test":
						metadata++
						return answer(req, http.StatusOK, "application/json", `{"path":"file","download_uri":"https://download.test/fresh"}`), nil
					case req.URL.Path == "/stale":
						stale++
						if refusal.status != 0 {
							return answer(req, refusal.status, refusal.contentType, refusal.body), nil
						}
					default:
						fresh++
					}
					if operation == "range" {
						return answer(req, http.StatusPartialContent, "application/octet-stream", content), nil
					}
					return answer(req, http.StatusOK, "application/octet-stream", content), nil
				})})
				file := (&File{File: &files_sdk.File{Path: "file", DownloadUri: "https://download.test/stale"}, FS: (&FS{Context: context.Background()}).Init(config, false)}).Init()
				var got []byte
				var err error
				if operation == "range" {
					var reader io.ReadCloser
					if reader, err = file.ReaderRange(0, int64(len(content))-1); err == nil {
						got, err = io.ReadAll(reader)
						require.NoError(t, reader.Close())
					}
				} else {
					got, err = io.ReadAll(file)
				}
				switch {
				case slices.Contains(expiredURLRefusals, refusal):
					require.NoError(t, err)
					require.Equal(t, content, string(got))
					require.Equal(t, []int{1, 1, 1}, []int{stale, metadata, fresh})
				case refusal.status == 0:
					require.NoError(t, err)
					require.Equal(t, content, string(got))
					require.Equal(t, []int{1, 0, 0}, []int{stale, metadata, fresh})
				default:
					require.Error(t, err)
					require.Equal(t, []int{0, 0}, []int{metadata, fresh})
				}
			})
		}
	}
}

func TestDownloadNeverSendsAnUnchangedExpiredRenewal(t *testing.T) {
	uri := "https://download.test/file?X-Files-Date=20000101T000000Z"
	var metadata, transfers, callbacks int
	config := files_sdk.Config{EndpointOverride: "https://api.test", DisableDirectTransfers: true}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(req *http.Request) (*http.Response, error) {
		if req.URL.Host != "api.test" {
			transfers++
		}
		metadata++
		payload, _ := json.Marshal(map[string]any{"path": "file", "download_uri": uri})
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(bytes.NewReader(payload)), ContentLength: int64(len(payload)), Request: req}, nil
	})})
	_, err := (&Client{Config: config}).Download(files_sdk.FileDownloadParams{File: files_sdk.File{Path: "file", DownloadUri: uri}}, files_sdk.ResponseBodyOption(func(io.ReadCloser) error { callbacks++; return nil }))
	require.ErrorIs(t, err, errDownloadURLExpired)
	require.Equal(t, 1, metadata)
	require.Zero(t, transfers)
	require.Zero(t, callbacks)
}

func TestOldCompleteHTTPCheckpointRestartsFromZero(t *testing.T) {
	f := newVersionDownloadFixture(t)
	f.data = bytes.Repeat([]byte("new"), 1<<18)
	dest := explicitDestination(filepath.Join(t.TempDir(), "out.bin"))
	tmp, err := tmpDownloadPath(dest, "")
	require.NoError(t, err)
	file, err := tmp.create()
	require.NoError(t, err)
	_, err = file.Write(bytes.Repeat([]byte("old"), 1<<18))
	require.NoError(t, err)
	identity, err := file.Stat()
	require.NoError(t, err)
	require.NoError(t, file.Close())
	paused, err := parkTmpDownload(tmp, identity)
	require.NoError(t, err)
	job := f.client().Downloader(DownloaderParams{RemotePath: "remote.bin", LocalPath: dest.String(), ResumeTmpPath: paused.String(), PriorJobCheckpoint: &JobDownloadCheckpoint{}, Manager: manager.Build(1, 1, true)})
	job.Start()
	job.Wait()
	require.NoError(t, job.Statuses[0].Err())
	out, err := os.ReadFile(dest.String())
	require.NoError(t, err)
	require.Equal(t, f.data, out)
	f.mu.Lock()
	defer f.mu.Unlock()
	require.Len(t, f.requests, 1)
	require.Empty(t, f.requests[0].ranges)
}

func TestHTTPPauseCheckpointRetainsItsVersion(t *testing.T) {
	const size, prefix = 4 << 20, (1 << 20) + 12345
	data := bytes.Repeat([]byte("x"), size)
	holding := make(chan struct{})
	var once sync.Once
	var mu sync.Mutex
	var requests []versionDownloadRequest
	var url string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"path": "remote.bin", "display_name": "remote.bin", "type": "file", "size": size, "download_uri": url + "/transfer"})
			return
		}
		mu.Lock()
		requests = append(requests, versionDownloadRequest{ranges: r.Header.Get("Range"), match: r.Header.Get("If-Match")})
		hold := len(requests) == 1
		mu.Unlock()
		start := int64(0)
		end := int64(size - 1)
		if r.Header.Get("Range") != "" {
			_, _ = fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &start, &end)
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
		}
		w.Header().Set("ETag", `"paused"`)
		w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
		if r.Header.Get("Range") != "" {
			w.WriteHeader(206)
			_, _ = w.Write(data[start : end+1])
			return
		}
		if !hold {
			_, _ = w.Write(data)
			return
		}
		_, _ = w.Write(data[:prefix])
		w.(http.Flusher).Flush()
		once.Do(func() { close(holding) })
		<-r.Context().Done()
	}))
	defer server.Close()
	url = server.URL
	client := &Client{Config: files_sdk.Config{EndpointOverride: url, DisableDirectTransfers: true}.Init()}
	dest := filepath.Join(t.TempDir(), "out.bin")
	testCtx, stopTest := context.WithTimeout(context.Background(), 30*time.Second)
	defer stopTest()
	ctx, cancel := context.WithCancelCause(testCtx)
	defer cancel(nil)
	params := DownloaderParams{RemotePath: "remote.bin", LocalPath: dest, Manager: manager.Build(1, 1, true)}
	job := client.Downloader(params, files_sdk.WithContext(ctx))
	job.Start()
	select {
	case <-holding:
	case <-time.After(5 * time.Second):
		t.Fatal("download did not reach the pause boundary")
	}
	require.Eventually(t, func() bool { return job.TransferBytes() == prefix }, 5*time.Second, time.Millisecond)
	cancel(ErrJobPaused)
	waitForDownloadTestJob(t, testCtx, job)
	checkpoint := job.DownloadCheckpoint()
	require.Len(t, checkpoint.Bindings, 1)
	require.EqualValues(t, prefix, checkpoint.Bindings[0].Prefix)
	params.PriorJobCheckpoint = checkpoint
	require.NotNil(t, ToStatusFile(job.Statuses[0]).DownloadBinding)
	resumed := client.Downloader(params, files_sdk.WithContext(testCtx))
	resumed.Start()
	waitForDownloadTestJob(t, testCtx, resumed)
	require.Equal(t, 1, resumed.Count(status.Complete))
	out, err := os.ReadFile(dest)
	require.NoError(t, err)
	require.Equal(t, data, out)
	mu.Lock()
	defer mu.Unlock()
	require.Len(t, requests, 2)
	require.Equal(t, fmt.Sprintf("bytes=%d-%d", prefix, size-1), requests[1].ranges)
	require.Equal(t, `"paused"`, requests[1].match)
}
