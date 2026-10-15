package file

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/stretchr/testify/require"
)

func TestDownloadV2DirectFailureMovesRemainingRangesToProxy(t *testing.T) {
	size := int64(20 * 1024 * 1024)
	source := bytes.Repeat([]byte("direct-range"), int(size)/len("direct-range")+1)[:size]
	directRequests := atomic.Int32{}
	// The first range fails before the plan admits the second, so the whole
	// transfer must reach the direct endpoint exactly once.
	info := testDirectTransferInfo(t, "/downloads?jwt=direct-token", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		directRequests.Add(1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	proxyRequests := atomic.Int32{}
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyRequests.Add(1)
		http.ServeContent(w, r, "download.bin", time.Time{}, bytes.NewReader(source))
	}))
	t.Cleanup(proxy.Close)

	config := files_sdk.Config{Environment: files_sdk.Development, Logger: log.New(io.Discard, "", 0)}.Init()
	sdkFile := files_sdk.File{
		DisplayName:          "download.bin",
		Path:                 "download.bin",
		Type:                 "file",
		Size:                 size,
		DownloadUri:          proxy.URL + "/download.bin?jwt=proxy-token",
		DirectConnectionInfo: info,
	}
	ranger := (&File{File: &sdkFile, FS: (&FS{Context: context.Background()}).Init(config, false)}).Init()
	fileInfo := Info{File: sdkFile, sizeTrust: TrustedSizeValue}
	tmpPath := filepath.Join(t.TempDir(), "download.bin.download")
	reportStatus := downloadV2TestStatus(ranger, fileInfo, DownloaderParams{
		AdaptiveConcurrency: true,
		Manager:             manager.Build(15, 1),
	}, tmpPath)

	used, finalSize, _, err := runDownloadV2IfSupported(context.Background(), reportStatus, fileInfo, explicitDestination(tmpPath), 0)

	require.True(t, used)
	require.NoError(t, err)
	require.Equal(t, size, finalSize)
	require.Equal(t, int32(1), directRequests.Load(), "a direct failure moves the rest of the transfer to the proxy")
	require.Equal(t, int32(2), proxyRequests.Load())
	written, err := os.ReadFile(tmpPath)
	require.NoError(t, err)
	require.True(t, bytes.Equal(source, written))
}

func TestDownloadSharedDirectSessionSendsLaterRequestsToProxy(t *testing.T) {
	proxyRequests := atomic.Int32{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyRequests.Add(1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	t.Cleanup(server.Close)

	sdkFile := files_sdk.File{
		Path:                 "/download.bin",
		DownloadUri:          server.URL + "/download.bin?jwt=proxy-token",
		DirectConnectionInfo: testDirectConnectionInfo("/downloads?jwt=direct-token"),
	}
	session := newDirectDownloadSession(nil)
	t.Cleanup(session.Close)
	ctx := context.WithValue(context.Background(), directTransferSessionContextKey{}, session)
	client := &Client{Config: files_sdk.Config{
		Logger: log.New(io.Discard, "", 0),
	}.Init()}
	closeBody := files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
		return body.Close()
	})

	_, err := client.Download(files_sdk.FileDownloadParams{File: sdkFile}, files_sdk.WithContext(ctx), closeBody)
	require.NoError(t, err)
	_, err = client.Download(files_sdk.FileDownloadParams{File: sdkFile}, files_sdk.WithContext(ctx), closeBody)
	require.NoError(t, err)

	stats := session.Stats()
	require.Equal(t, int64(1), stats.Attempts)
	require.Equal(t, "direct_request_failed", stats.DisabledReason)
	require.Equal(t, int32(2), proxyRequests.Load())
}

func TestDownloadDirectTransfersDisabledUsesProxy(t *testing.T) {
	proxyRequests := atomic.Int32{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyRequests.Add(1)
		_, _ = w.Write([]byte("ok"))
	}))
	t.Cleanup(server.Close)

	session := newDirectDownloadSession(nil)
	t.Cleanup(session.Close)
	ctx := context.WithValue(context.Background(), directTransferSessionContextKey{}, session)
	client := &Client{Config: files_sdk.Config{
		DisableDirectTransfers: true,
		Logger:                 log.New(io.Discard, "", 0),
	}.Init()}
	_, err := client.Download(files_sdk.FileDownloadParams{File: files_sdk.File{
		Path:                 "/download.bin",
		DownloadUri:          server.URL + "/download.bin?jwt=proxy-token",
		DirectConnectionInfo: testDirectConnectionInfo("/downloads?jwt=direct-token"),
	}}, files_sdk.WithContext(ctx), files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
		return body.Close()
	}))

	require.NoError(t, err)
	require.Zero(t, session.Stats().Attempts)
	require.Equal(t, int32(1), proxyRequests.Load())
}
