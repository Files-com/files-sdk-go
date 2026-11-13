//go:build linux || windows

package fsmount

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

func TestSDKRemoteBackendDownloadRange(t *testing.T) {
	payload := []byte("0123456789")
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Range"); got != "bytes=10-19" {
			t.Errorf("Range header = %q, want bytes=10-19", got)
		}
		if got := request.Header.Get("X-Test-Request"); got != "preserved" {
			t.Errorf("X-Test-Request = %q, want preserved", got)
		}
		response.Header().Set("Content-Range", "bytes 10-19/100")
		response.Header().Set("Content-Length", "10")
		response.WriteHeader(http.StatusPartialContent)
		_, _ = response.Write(payload)
	}))
	defer server.Close()

	backend := newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true})
	headers := &http.Header{}
	headers.Set("X-Test-Request", "preserved")
	requested := cache.ByteRange{Start: 10, End: 20}
	result, err := backend.downloadRange(
		files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.bin", Size: 100, DownloadUri: server.URL}},
		requested,
		files_sdk.RequestHeadersOption(headers),
	)
	if err != nil {
		t.Fatalf("downloadRange failed: %v", err)
	}
	defer result.Body.Close()
	if result.Returned != requested {
		t.Fatalf("returned range = %#v, want %#v", result.Returned, requested)
	}
	if !result.Partial || result.TotalSize != 100 {
		t.Fatalf("result = partial %t total %d, want true and 100", result.Partial, result.TotalSize)
	}
	data, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("reading range body failed: %v", err)
	}
	if !bytes.Equal(data, payload) {
		t.Fatalf("range body = %q, want %q", data, payload)
	}
}

func TestSDKRemoteBackendDownloadRangeAcceptsCompleteFallback(t *testing.T) {
	payload := []byte("complete response body")
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.WriteHeader(http.StatusOK)
		_, _ = response.Write(payload)
	}))
	defer server.Close()

	backend := newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true})
	result, err := backend.downloadRange(
		files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.bin", Size: int64(len(payload)), DownloadUri: server.URL}},
		cache.ByteRange{Start: 10, End: 15},
	)
	if err != nil {
		t.Fatalf("downloadRange failed: %v", err)
	}
	defer result.Body.Close()
	if result.Partial {
		t.Fatal("complete fallback was marked partial")
	}
	wantReturned := cache.ByteRange{Start: 0, End: int64(len(payload))}
	if result.Returned != wantReturned || result.TotalSize != int64(len(payload)) {
		t.Fatalf("complete fallback = returned %#v total %d, want %#v and %d", result.Returned, result.TotalSize, wantReturned, len(payload))
	}
	data, err := io.ReadAll(result.Body)
	if err != nil || !bytes.Equal(data, payload) {
		t.Fatalf("complete body = %q, %v; want %q, nil", data, err, payload)
	}
}

func TestRemoteFsRangeSourceCompleteDownloadOmitsRangeHeader(t *testing.T) {
	payload := []byte("ordinary complete response body")
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Range"); got != "" {
			t.Errorf("Range header = %q, want empty", got)
		}
		response.Header().Set("Content-Length", fmt.Sprint(len(payload)))
		response.WriteHeader(http.StatusOK)
		_, _ = response.Write(payload)
	}))
	defer server.Close()

	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()
	path := "/ordinary-complete.bin"
	modTime := time.Now().Add(-time.Minute).Round(0)
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(payload)),
		modTime:      modTime,
		creationTime: modTime,
	})
	node.setDownloadURI(server.URL)
	fs.backend = newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true})

	meta := cacheEntryMetadata(path, int64(len(payload)), modTime)
	complete := cache.ByteRange{Start: 0, End: int64(len(payload))}
	result, err := (&remoteFsRangeSource{fs: fs}).DownloadComplete(context.Background(), path, meta, complete)
	if err != nil {
		t.Fatalf("DownloadComplete failed: %v", err)
	}
	defer result.Body.Close()
	data, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("reading complete body failed: %v", err)
	}
	if !bytes.Equal(data, payload) {
		t.Fatalf("complete body = %q, want %q", data, payload)
	}
	if result.Partial || result.Returned != complete || result.TotalSize != int64(len(payload)) {
		t.Fatalf("complete response = partial %t returned %#v total %d", result.Partial, result.Returned, result.TotalSize)
	}
}

func TestRemoteFsDateGuardSurvivesTheRangeToWholeFallback(t *testing.T) {
	const baseline = "Wed, 21 Oct 2015 07:28:00 GMT"
	for _, tc := range []struct {
		name, initial, later                string
		ignoreCondition, changed, duplicate bool
		clientCheck, mountCheck             *bool
		disabled, mutate                    bool
	}{
		{name: "stable", initial: baseline, later: baseline},
		{name: "earlier", initial: baseline, later: "Wed, 21 Oct 2015 07:27:59 GMT"},
		{name: "missing"},
		{name: "duplicate", initial: baseline, later: baseline, duplicate: true},
		{name: "invalid", initial: "not an HTTP date", later: "not an HTTP date"},
		{name: "later-missing", initial: baseline},
		{name: "later-invalid", initial: baseline, later: "not an HTTP date"},
		{name: "newer-precondition", initial: baseline, later: "Wed, 21 Oct 2015 07:28:01 GMT", changed: true},
		{name: "newer-ignored-condition", initial: baseline, later: "Wed, 21 Oct 2015 07:28:01 GMT", ignoreCondition: true, changed: true},
		{name: "client-disabled", initial: baseline, later: "Wed, 21 Oct 2015 07:28:01 GMT", clientCheck: lib.Bool(false), disabled: true},
		{name: "mount-disabled", initial: baseline, later: "Wed, 21 Oct 2015 07:28:01 GMT", clientCheck: lib.Bool(true), mountCheck: lib.Bool(false), disabled: true},
		{name: "mount-enabled", initial: baseline, later: "Wed, 21 Oct 2015 07:28:01 GMT", clientCheck: lib.Bool(false), mountCheck: lib.Bool(true), changed: true},
		{name: "mount-snapshot", initial: baseline, later: "Wed, 21 Oct 2015 07:28:01 GMT", clientCheck: lib.Bool(true), mutate: true, changed: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			payload := bytes.Repeat([]byte("a"), 2<<20)
			var mu sync.Mutex
			var requests []http.Header
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				requests = append(requests, r.Header.Clone())
				count := len(requests)
				mu.Unlock()
				modified := tc.initial
				if count > 1 {
					modified = tc.later
				}
				w.Header().Set("Last-Modified", modified)
				if tc.duplicate {
					w.Header().Add("Last-Modified", modified)
				}
				if tc.changed && !tc.ignoreCondition && r.Header.Get("If-Unmodified-Since") == baseline {
					w.WriteHeader(http.StatusPreconditionFailed)
					return
				}
				start, end := 0, len(payload)-1
				if ranges := r.Header.Get("Range"); ranges != "" {
					_, _ = fmt.Sscanf(ranges, "bytes=%d-%d", &start, &end)
					w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(payload)))
					w.Header().Set("Content-Length", fmt.Sprint(end-start+1))
					w.WriteHeader(http.StatusPartialContent)
				} else {
					w.Header().Set("Content-Length", fmt.Sprint(len(payload)))
				}
				_, _ = w.Write(payload[start : end+1])
			}))
			defer server.Close()
			fs, vfs, _ := newTestRemoteFs(t)
			defer vfs.destroy()
			path := "/dated.bin"
			modTime := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)
			node := vfs.getOrCreate(path, nodeTypeFile)
			node.updateInfo(fsNodeInfo{nodeType: nodeTypeFile, size: int64(len(payload)), modTime: modTime, creationTime: modTime})
			node.setDownloadURI(server.URL)
			store, err := disk.NewDiskCache(t.TempDir(), disk.WithCapacityBytes(8<<30))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(store.StopMaintenance)
			config := files_sdk.Config{DisableDirectTransfers: true, CheckDownloadModificationTime: tc.clientCheck, FeatureFlags: map[string]bool{files_sdk.FeatureFlagSparseRangeReads: true}}
			fs, err = newRemoteFs(MountParams{Config: &config, CheckDownloadModificationTime: tc.mountCheck, IgnorePatterns: []string{}}, vfs, fs.log, store)
			if err != nil {
				t.Fatal(err)
			}
			if tc.mutate {
				*tc.clientCheck = false
			}
			fs.backend = newRangeTestSDKBackend(config)
			streams := fs.rangeDownloads()
			t.Cleanup(streams.Close)
			streams.policy = testStreamPolicy()
			streams.policy.smallFileSize = 0
			meta := cacheEntryMetadata(path, int64(len(payload)), modTime)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = streams.Wait(ctx, 1, path, meta, cache.ByteRange{Start: 0, End: 4096})
			if tc.changed {
				if !errors.Is(err, errRangeVersionChanged) {
					t.Fatalf("download error = %v, want changed source", err)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if got := readThrough(t, store, streams, meta, int64(len(payload))-4096, 4096); !bytes.Equal(got, payload[len(payload)-4096:]) {
					t.Fatal("whole response returned the wrong bytes")
				}
			}
			streams.Close()
			mu.Lock()
			defer mu.Unlock()
			if len(requests) != 2 || requests[0].Get("Range") == "" || requests[1].Get("Range") != "" {
				t.Fatalf("requests = %v, want one range probe followed by a whole response", requests)
			}
			condition := ""
			if tc.initial == baseline && !tc.duplicate && !tc.disabled {
				condition = baseline
			}
			if got := requests[1].Get("If-Unmodified-Since"); got != condition {
				t.Fatalf("If-Unmodified-Since = %q, want %q", got, condition)
			}
		})
	}
}

func TestSDKRemoteBackendDownloadRangeRejectsInvalidResponses(t *testing.T) {
	tests := []struct {
		name          string
		status        int
		contentRange  string
		contentLen    string
		wantErr       error
		unknownLength bool
	}{
		{name: "missing content range", status: http.StatusPartialContent, contentLen: "10", wantErr: errMalformedRangeResponse},
		{name: "wrong returned range", status: http.StatusPartialContent, contentRange: "bytes 11-20/100", contentLen: "10", wantErr: errMalformedRangeResponse},
		{name: "wrong content length", status: http.StatusPartialContent, contentRange: "bytes 10-19/100", contentLen: "9", wantErr: errMalformedRangeResponse},
		{name: "unexpected success status", status: http.StatusNoContent, wantErr: errMalformedRangeResponse},
		{name: "complete response without length", status: http.StatusOK, unknownLength: true, wantErr: errMalformedRangeResponse},
		{name: "precondition failed", status: http.StatusPreconditionFailed, wantErr: errRangeVersionChanged},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				if test.contentRange != "" {
					response.Header().Set("Content-Range", test.contentRange)
				}
				if test.contentLen != "" {
					response.Header().Set("Content-Length", test.contentLen)
				}
				response.WriteHeader(test.status)
				if test.unknownLength {
					response.(http.Flusher).Flush()
				}
			}))
			defer server.Close()

			transport := &rangeCloseTransport{RoundTripper: http.DefaultTransport.(*http.Transport).Clone()}
			t.Cleanup(transport.RoundTripper.(*http.Transport).CloseIdleConnections)
			backend := newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true}.Init().SetCustomClient(&http.Client{Transport: transport}))
			knownSize := int64(100)
			if test.unknownLength {
				knownSize = -1
			}
			_, err := backend.downloadRange(
				files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.bin", Size: knownSize, DownloadUri: server.URL}},
				cache.ByteRange{Start: 10, End: 20},
			)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("downloadRange error = %v, want %v", err, test.wantErr)
			}
			transport.assertClosedOnce(t)
		})
	}
}

func TestSDKRemoteBackendDownloadRangeRejectsShortAndLongBodies(t *testing.T) {
	tests := []struct {
		name    string
		write   func(http.ResponseWriter)
		wantErr error
	}{
		{
			name: "short",
			write: func(response http.ResponseWriter) {
				response.Header().Set("Content-Range", "bytes 0-9/100")
				response.Header().Set("Content-Length", "10")
				response.WriteHeader(http.StatusPartialContent)
				_, _ = response.Write([]byte("short"))
			},
			wantErr: errShortRangeResponse,
		},
		{
			name: "long chunked body",
			write: func(response http.ResponseWriter) {
				response.Header().Set("Content-Range", "bytes 0-9/100")
				response.WriteHeader(http.StatusPartialContent)
				response.(http.Flusher).Flush()
				_, _ = response.Write([]byte("0123456789x"))
			},
			wantErr: errLongRangeResponse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				test.write(response)
			}))
			defer server.Close()

			backend := newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true})
			result, err := backend.downloadRange(
				files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.bin", Size: 100, DownloadUri: server.URL}},
				cache.ByteRange{Start: 0, End: 10},
			)
			if err != nil {
				t.Fatalf("downloadRange failed before reading body: %v", err)
			}
			_, readErr := io.ReadAll(result.Body)
			_ = result.Body.Close()
			if !errors.Is(readErr, test.wantErr) {
				t.Fatalf("body error = %v, want %v", readErr, test.wantErr)
			}
		})
	}
}

// A range refreshes its URL once when the refusal says the URL expired; any other refusal fails the range without a
// refresh.
func TestSDKRemoteBackendDownloadRangeRefreshesExpiredURLOnce(t *testing.T) {
	for name, test := range map[string]struct {
		status        int
		contentType   string
		body          string
		wantRefreshed bool
	}{
		"files expired":    {http.StatusForbidden, "application/json", `{"type":"download_request_expired","error":"expired"}`, true},
		"storage expired":  {http.StatusBadRequest, "application/xml", `<Error><Code>ExpiredToken</Code><Message>The provided token has expired.</Message></Error>`, true},
		"expired word":     {http.StatusBadRequest, "application/xml", `<Error><Code>InvalidArgument</Code><Message>X-Amz-Expires must be less than a week</Message></Error>`, false},
		"structured error": {http.StatusBadRequest, "application/json", `{"type":"bad-request/invalid-body","error":"Request has expired"}`, false},
	} {
		t.Run(name, func(t *testing.T) {
			var expiredCalls atomic.Int32
			var metadataCalls atomic.Int32
			var freshCalls atomic.Int32
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
				switch request.URL.Path {
				case "/expired":
					expiredCalls.Add(1)
					response.Header().Set("Content-Type", test.contentType)
					response.WriteHeader(test.status)
					_, _ = response.Write([]byte(test.body))
				case "/fresh":
					freshCalls.Add(1)
					response.Header().Set("Content-Range", "bytes 5-9/100")
					response.Header().Set("Content-Length", "5")
					response.WriteHeader(http.StatusPartialContent)
					_, _ = response.Write([]byte("fresh"))
				default:
					metadataCalls.Add(1)
					response.Header().Set("Content-Type", "application/json")
					_, _ = fmt.Fprintf(response, `{"path":"remote.bin","size":100,"download_uri":%q}`, server.URL+"/fresh")
				}
			}))
			defer server.Close()

			backend := newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true, EndpointOverride: server.URL})
			result, err := backend.downloadRange(
				files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.bin", Size: 100, DownloadUri: server.URL + "/expired"}},
				cache.ByteRange{Start: 5, End: 10},
			)
			if !test.wantRefreshed {
				if err == nil {
					_ = result.Body.Close()
					t.Fatal("downloadRange succeeded after a refusal that is not an expiry")
				}
				if expiredCalls.Load() != 1 || metadataCalls.Load() != 0 || freshCalls.Load() != 0 {
					t.Fatalf("calls = expired %d metadata %d fresh %d; want 1, 0, 0", expiredCalls.Load(), metadataCalls.Load(), freshCalls.Load())
				}
				return
			}
			if err != nil {
				t.Fatalf("downloadRange failed: %v", err)
			}
			data, readErr := io.ReadAll(result.Body)
			_ = result.Body.Close()
			if readErr != nil || string(data) != "fresh" {
				t.Fatalf("refreshed body = %q, %v; want fresh, nil", data, readErr)
			}
			if expiredCalls.Load() != 1 || metadataCalls.Load() != 1 || freshCalls.Load() != 1 {
				t.Fatalf("calls = expired %d metadata %d fresh %d; want 1 each", expiredCalls.Load(), metadataCalls.Load(), freshCalls.Load())
			}
		})
	}
}

func TestSDKRemoteBackendDownloadRangeRefreshesMetadataAfter416(t *testing.T) {
	var staleCalls atomic.Int32
	var metadataCalls atomic.Int32
	var freshCalls atomic.Int32
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/stale":
			staleCalls.Add(1)
			response.Header().Set("Content-Range", "bytes */100")
			response.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		case "/fresh":
			freshCalls.Add(1)
			response.Header().Set("Content-Range", "bytes 5-9/100")
			response.Header().Set("Content-Length", "5")
			response.WriteHeader(http.StatusPartialContent)
			_, _ = response.Write([]byte("fresh"))
		default:
			metadataCalls.Add(1)
			response.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(response, `{"path":"remote.bin","size":100,"download_uri":%q}`, server.URL+"/fresh")
		}
	}))
	defer server.Close()

	transport := &rangeCloseTransport{RoundTripper: http.DefaultTransport.(*http.Transport).Clone()}
	t.Cleanup(transport.RoundTripper.(*http.Transport).CloseIdleConnections)
	backend := newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true, EndpointOverride: server.URL}.Init().SetCustomClient(&http.Client{Transport: transport}))
	result, err := backend.downloadRange(
		files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.bin", Size: 100, DownloadUri: server.URL + "/stale"}},
		cache.ByteRange{Start: 5, End: 10},
	)
	if err != nil {
		t.Fatalf("downloadRange failed: %v", err)
	}
	data, readErr := io.ReadAll(result.Body)
	_ = result.Body.Close()
	if readErr != nil || string(data) != "fresh" {
		t.Fatalf("refreshed body = %q, %v; want fresh, nil", data, readErr)
	}
	transport.assertClosedOnce(t)
	if staleCalls.Load() != 1 || metadataCalls.Load() != 1 || freshCalls.Load() != 1 {
		t.Fatalf("calls = stale %d metadata %d fresh %d; want 1 each", staleCalls.Load(), metadataCalls.Load(), freshCalls.Load())
	}
}

func TestSDKRemoteBackendDownloadRangeReturns416AfterOneRefresh(t *testing.T) {
	var rangeCalls atomic.Int32
	var metadataCalls atomic.Int32
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/api/rest/v1/files/remote.bin" {
			metadataCalls.Add(1)
			response.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(response, `{"path":"remote.bin","size":100,"download_uri":%q}`, server.URL+"/range")
			return
		}
		rangeCalls.Add(1)
		response.Header().Set("Content-Range", "bytes */100")
		response.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
	}))
	defer server.Close()

	backend := newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true, EndpointOverride: server.URL})
	_, err := backend.downloadRange(
		files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.bin", Size: 100, DownloadUri: server.URL + "/range"}},
		cache.ByteRange{Start: 5, End: 10},
	)
	if !errors.Is(err, errRangeNotSatisfiable) {
		t.Fatalf("downloadRange error = %v, want errRangeNotSatisfiable", err)
	}
	if rangeCalls.Load() != 2 || metadataCalls.Load() != 1 {
		t.Fatalf("calls = range %d metadata %d; want 2 and 1", rangeCalls.Load(), metadataCalls.Load())
	}
}

func TestSDKRemoteBackendDownloadRangeHonorsCancellation(t *testing.T) {
	requestStarted := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Range", "bytes 0-9/100")
		response.Header().Set("Content-Length", "10")
		response.WriteHeader(http.StatusPartialContent)
		response.(http.Flusher).Flush()
		close(requestStarted)
		<-request.Context().Done()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	backend := newRangeTestSDKBackend(files_sdk.Config{DisableDirectTransfers: true})
	result, err := backend.downloadRange(
		files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.bin", Size: 100, DownloadUri: server.URL}},
		cache.ByteRange{Start: 0, End: 10},
		files_sdk.WithContext(ctx),
	)
	if err != nil {
		t.Fatalf("downloadRange failed: %v", err)
	}
	<-requestStarted
	cancel()
	_, readErr := io.ReadAll(result.Body)
	_ = result.Body.Close()
	if !errors.Is(readErr, context.Canceled) {
		t.Fatalf("body error = %v, want context.Canceled", readErr)
	}
}

func TestProviderRemoteBackendDownloadRangeIsUnsupported(t *testing.T) {
	backend := &providerRemoteBackend{}
	_, err := backend.downloadRange(files_sdk.FileDownloadParams{}, cache.ByteRange{Start: 0, End: 1})
	if !errors.Is(err, errRangeDownloadUnsupported) {
		t.Fatalf("downloadRange error = %v, want errRangeDownloadUnsupported", err)
	}
}

func newRangeTestSDKBackend(config files_sdk.Config) *sdkRemoteBackend {
	config = config.Init()
	return &sdkRemoteBackend{fileClient: &file.Client{Config: config}}
}

// Observe the real HTTP bodies while leaving request and response processing to the SDK.
type rangeCloseTransport struct {
	http.RoundTripper
	bodies []*rangeCloseBody
}

func (r *rangeCloseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	response, err := r.RoundTripper.RoundTrip(req)
	if err == nil {
		body := &rangeCloseBody{ReadCloser: response.Body}
		response.Body = body
		r.bodies = append(r.bodies, body)
	}
	return response, err
}
func (r *rangeCloseTransport) assertClosedOnce(t *testing.T) {
	t.Helper()
	if len(r.bodies) == 0 {
		t.Fatal("no response reached the range callback")
	}
	for _, body := range r.bodies {
		if got := body.closes.Load(); got != 1 {
			t.Errorf("raw response close calls = %d, want1", got)
		}
	}
}

type rangeCloseBody struct {
	io.ReadCloser
	closes atomic.Int32
}

func (b *rangeCloseBody) Close() error { b.closes.Add(1); return b.ReadCloser.Close() }
