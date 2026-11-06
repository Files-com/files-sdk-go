package file

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/stretchr/testify/require"
)

func TestPublicDownloadAcceptsANilHTTPClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = io.WriteString(w, "downloaded") }))
	defer server.Close()
	config := files_sdk.Config{DisableDirectTransfers: true}.Init()
	config.HTTPClient = nil
	var output []byte
	_, err := (&Client{Config: config}).Download(files_sdk.FileDownloadParams{File: files_sdk.File{Path: "file", DownloadUri: server.URL}}, files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
		defer body.Close()
		var err error
		output, err = io.ReadAll(body)
		return err
	}))
	require.NoError(t, err)
	require.Equal(t, "downloaded", string(output))
	require.Nil(t, config.HTTPClient)
}

func TestDownloadDoesNotRenewAnExpiryErrorReturnedAfterTheBodyCallbackStarts(t *testing.T) {
	var url string
	var metadata atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			metadata.Add(1)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"path": "file", "download_uri": url + "/data"})
			return
		}
		_, _ = io.WriteString(w, "body bytes")
	}))
	defer server.Close()
	url = server.URL
	client := &Client{Config: files_sdk.Config{EndpointOverride: url, DisableDirectTransfers: true}.Init()}
	var calls int
	_, err := client.Download(files_sdk.FileDownloadParams{File: files_sdk.File{Path: "file", DownloadUri: url + "/data"}}, files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
		calls++
		_, readErr := io.ReadAll(body)
		closeErr := body.Close()
		require.NoError(t, errors.Join(readErr, closeErr))
		// A nested SDK operation can return this cause to a user's callback.
		return errDownloadURLExpired
	}))
	require.ErrorIs(t, err, errDownloadURLExpired)
	require.Equal(t, 1, calls)
	require.Zero(t, metadata.Load())
}

func TestSiblingPublicDownloadsReuseTheirJobsDirectConnection(t *testing.T) {
	var mu sync.Mutex
	var peers []string
	info := testDirectTransferInfo(t, "/download", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		peers = append(peers, r.RemoteAddr)
		mu.Unlock()
		_, _ = io.WriteString(w, "direct")
	}))
	var proxyCalls atomic.Int32
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { proxyCalls.Add(1); _, _ = io.WriteString(w, "proxy") }))
	defer proxy.Close()
	job := (&Job{}).Init()
	ctx := job.withDirectTransferClientCache(context.Background())
	defer job.Finish()
	client := &Client{Config: files_sdk.Config{EndpointOverride: proxy.URL, Environment: files_sdk.Development}.Init()}
	for _, path := range []string{"first.bin", "second.bin"} {
		var output []byte
		_, err := client.Download(files_sdk.FileDownloadParams{File: files_sdk.File{Path: path, DownloadUri: proxy.URL + "/" + path, DirectConnectionInfo: info}}, files_sdk.WithContext(ctx), files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
			defer body.Close()
			var err error
			output, err = io.ReadAll(body)
			return err
		}))
		require.NoError(t, err)
		require.Equal(t, "direct", string(output))
	}
	mu.Lock()
	defer mu.Unlock()
	require.Len(t, peers, 2)
	require.Equal(t, peers[0], peers[1])
	require.Zero(t, proxyCalls.Load())
}

func TestPublicFileRangeAccessorsResumeBoundResponses(t *testing.T) {
	for _, operation := range []string{"read-at", "reader-range"} {
		t.Run(operation, func(t *testing.T) {
			f := newVersionDownloadFixture(t)
			remote := (&File{File: &files_sdk.File{Path: "remote.bin", DownloadUri: f.url + "/transfer"}, FS: (&FS{Context: context.Background()}).Init(f.client().Config, false)}).Init()
			const size = 6 << 20
			var out []byte
			if operation == "read-at" {
				out = make([]byte, size)
				n, err := remote.ReadAt(out, 0)
				require.NoError(t, err)
				require.Equal(t, size, n)
			} else {
				reader, err := remote.ReaderRange(0, size-1)
				require.NoError(t, err)
				out, err = io.ReadAll(reader)
				require.NoError(t, err)
				require.NoError(t, reader.Close())
			}
			require.Equal(t, f.data[:size], out)
			f.mu.Lock()
			defer f.mu.Unlock()
			require.Len(t, f.requests, 2)
			require.Equal(t, "bytes=1048576-6291455", f.requests[1].ranges)
			require.Equal(t, `"original"`, f.requests[1].match)
		})
	}
}

func TestReadAtRejectsWrongCoverageBeforeChangingTheBuffer(t *testing.T) {
	for _, etag := range []string{"", `"v1"`} {
		t.Run(etag, func(t *testing.T) {
			body := &ownedHTTPResponseBody{Reader: strings.NewReader("abc")}
			header := http.Header{"Content-Range": []string{"bytes 0-2/6"}}
			if etag != "" {
				header.Set("ETag", etag)
			}
			config := ownedResponseConfig(body, 206, header)
			file := (&File{File: &files_sdk.File{Path: "file", DownloadUri: "https://download.test/file"}, FS: (&FS{Context: context.Background()}).Init(config, false)}).Init()
			out := []byte("XYZ")
			n, err := file.ReadAt(out, 3)
			require.Error(t, err)
			require.Zero(t, n)
			require.Equal(t, []byte("XYZ"), out)
			require.Equal(t, 1, body.closes)
		})
	}
}

func TestDownloaderFallsBackWhenTheInitialRangeIsIgnored(t *testing.T) {
	f := newVersionDownloadFixture(t)
	f.ignoresRange = true
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

func TestDirectRangeWithoutProofSelectsAnOwnedWholeResponse(t *testing.T) {
	const size = 20 << 20
	data := bytes.Repeat([]byte("v"), size)
	var mu sync.Mutex
	var ranges []string
	info := testDirectTransferInfo(t, "/download", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		ranges = append(ranges, r.Header.Get("Range"))
		mu.Unlock()
		http.ServeContent(w, r, "remote.bin", time.Time{}, bytes.NewReader(data))
	}))
	var url string
	var proxy atomic.Int32
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			proxy.Add(1)
			_, _ = io.WriteString(w, "proxy")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(files_sdk.File{Path: "remote.bin", DisplayName: "remote.bin", Type: "file", Size: size, DownloadUri: url + "/proxy", DirectConnectionInfo: info})
	}))
	defer api.Close()
	url = api.URL
	client := &Client{Config: files_sdk.Config{EndpointOverride: url, Environment: files_sdk.Development}.Init()}
	dest := filepath.Join(t.TempDir(), "out.bin")
	job := client.Downloader(DownloaderParams{RemotePath: "remote.bin", LocalPath: dest, AdaptiveConcurrency: true, Manager: manager.Build(2, 1)})
	job.Start()
	job.Wait()
	require.NoError(t, job.Statuses[0].Err())
	out, err := os.ReadFile(dest)
	require.NoError(t, err)
	require.Equal(t, sha256.Sum256(data), sha256.Sum256(out))
	require.Zero(t, proxy.Load())
	mu.Lock()
	defer mu.Unlock()
	require.Len(t, ranges, 2)
	require.NotEmpty(t, ranges[0])
	require.Empty(t, ranges[1])
}

type downloadCloseFailure struct {
	io.ReadCloser
	err error
}

func (b downloadCloseFailure) Close() error { return errors.Join(b.ReadCloser.Close(), b.err) }

func TestWholeResponseFallbackDoesNotDiscardAResponseCloseFailure(t *testing.T) {
	for _, adaptive := range []bool{false, true} {
		for _, failure := range []struct {
			name string
			err  error
		}{
			{"close error", errors.New("closing unused range failed")},
			{"close EOF", io.EOF},
			{"close cancellation", fmt.Errorf("closing canceled response: %w", context.Canceled)},
		} {
			t.Run(fmt.Sprintf("adaptive=%v/%s", adaptive, failure.name), func(t *testing.T) {
				f := newVersionDownloadFixture(t)
				f.noValidator = true
				client := f.client()
				base := client.Config.HTTPClient.Transport
				httpClient := *client.Config.HTTPClient
				httpClient.Transport = ownedHTTPResponseTransport(func(req *http.Request) (*http.Response, error) {
					response, err := base.RoundTrip(req)
					if response != nil && response.StatusCode == 206 {
						response.Body = downloadCloseFailure{ReadCloser: response.Body, err: failure.err}
					}
					return response, err
				})
				client.Config = client.Config.SetCustomClient(&httpClient)
				dest := filepath.Join(t.TempDir(), "out.bin")
				var reportedErr error
				events := CreateFileEvents(func(file JobFile) { reportedErr = file.Err }, status.Errored, status.Canceled)
				job := client.Downloader(DownloaderParams{RemotePath: "remote.bin", LocalPath: dest, Manager: manager.Build(2, 1), AdaptiveConcurrency: adaptive, EventsReporter: events})
				job.Start()
				job.Wait()
				require.ErrorIs(t, job.Statuses[0].Err(), failure.err)
				require.ErrorIs(t, reportedErr, failure.err)
				_, err := os.Stat(dest)
				require.True(t, os.IsNotExist(err))
				f.mu.Lock()
				defer f.mu.Unlock()
				require.Len(t, f.requests, 1)
				require.NotEmpty(t, f.requests[0].ranges)
			})
		}
	}
}

func TestRejectedRangeCloseFailureStopsRenewalAndRetry(t *testing.T) {
	for _, code := range []int{http.StatusForbidden, http.StatusBadGateway} {
		for _, closeFails := range []bool{true, false} {
			t.Run(fmt.Sprintf("status=%d/close-fails=%v", code, closeFails), func(t *testing.T) {
				closeErr := errors.New("closing rejected range failed")
				rejected := &ownedHTTPResponseBody{Reader: strings.NewReader("rejected")}
				if closeFails {
					rejected.closeErr = closeErr
				}
				var storage, metadata int
				config := files_sdk.Config{EndpointOverride: "https://api.test", DisableDirectTransfers: true}.Init().SetCustomClient(&http.Client{Transport: ownedHTTPResponseTransport(func(req *http.Request) (*http.Response, error) {
					if req.URL.Host == "api.test" {
						metadata++
						payload := `{"path":"file","download_uri":"https://storage.test/fresh"}`
						return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(payload)), ContentLength: int64(len(payload)), Request: req}, nil
					}
					storage++
					if storage == 1 {
						return &http.Response{StatusCode: code, Status: fmt.Sprintf("%d %s", code, http.StatusText(code)), Header: http.Header{"Content-Type": []string{"text/plain"}}, Body: rejected, ContentLength: 8, Request: req}, nil
					}
					header := http.Header{}
					header.Set("ETag", `"same"`)
					header.Set("Content-Range", "bytes 0-2/3")
					return &http.Response{StatusCode: 206, Header: header, Body: io.NopCloser(strings.NewReader("new")), ContentLength: 3, Request: req}, nil
				})})
				// Let the shared reader, rather than the HTTP-status layer, own
				// any second attempt in this focused regression.
				config.RetryMax = 0
				remote := (&File{File: &files_sdk.File{Path: "file", DownloadUri: "https://storage.test/old"}, FS: (&FS{Context: context.Background()}).Init(config, false)}).Init()
				reader, _, err := remote.openDownloadRange(context.Background(), 0, 2, downloadV2RetryAttempts)
				if closeFails {
					require.ErrorIs(t, err, closeErr)
					require.Nil(t, reader)
					require.Equal(t, 1, storage)
					require.Zero(t, metadata)
				} else {
					require.NoError(t, err)
					data, err := io.ReadAll(reader)
					require.NoError(t, err)
					require.NoError(t, reader.Close())
					require.Equal(t, "new", string(data))
					require.Equal(t, 2, storage)
					if code == http.StatusForbidden {
						require.Equal(t, 1, metadata)
					} else {
						require.Zero(t, metadata)
					}
				}
				require.Equal(t, 1, rejected.closes)
			})
		}
	}
}

func TestDownloadV2RetainsAnOutputCloseFailureAfterPause(t *testing.T) {
	const size = int64(20 << 20)
	ranger := &downloadV2TestRangeFile{downloadURI: "https://download.test/file", info: Info{File: files_sdk.File{DisplayName: "file", Path: "file", Type: "file", Size: size}, sizeTrust: TrustedSizeValue}}
	path := filepath.Join(t.TempDir(), "out.download")
	file, err := os.Create(path)
	require.NoError(t, err)
	params := DownloaderParams{AdaptiveConcurrency: true, Manager: manager.Build(2, 1)}
	report := downloadV2TestStatus(ranger, ranger.info, params, path)
	engine := newDownloadV2Engine(report, ranger, file, downloadV2TargetDefault, size, 0, downloadV2KnownSizePartSize(downloadV2TargetDefault, size), params)
	previous := truncateTmpDownload
	truncateTmpDownload = func(target interface{ Truncate(int64) error }, size int64) error {
		if err := previous(target, size); err != nil {
			return err
		}
		// The final close must report an independently invalidated handle.
		return file.Close()
	}
	defer func() { truncateTmpDownload = previous }()
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(ErrJobPaused)
	err = engine.Run(ctx)
	require.ErrorIs(t, err, context.Canceled)
	require.ErrorIs(t, err, os.ErrClosed)
	require.False(t, stoppedForPauseOnly(err))
}

func TestPublicSingleResponseRangeWithoutETagRemainsSupported(t *testing.T) {
	f := newVersionDownloadFixture(t)
	f.noValidator = true
	remote := (&File{File: &files_sdk.File{Path: "remote.bin", DownloadUri: f.url + "/transfer"}, FS: (&FS{Context: context.Background()}).Init(f.client().Config, false)}).Init()
	reader, err := remote.ReaderRange(5, 20)
	require.NoError(t, err)
	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.NoError(t, reader.Close())
	require.Equal(t, f.data[5:21], data)
}
