package file

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/require"
)

type testDirectTransferDownloadSuppressor struct {
	disabled atomic.Bool
	disables atomic.Int32
	reason   atomic.Value
}

func (s *testDirectTransferDownloadSuppressor) directTransferDownloadAttemptAllowed() bool {
	return !s.disabled.Load()
}

func (s *testDirectTransferDownloadSuppressor) disableDirectTransferDownload(reason string, err error) {
	if s.disabled.CompareAndSwap(false, true) {
		s.reason.Store(reason)
		s.disables.Add(1)
	}
}

func TestDownloadV2DirectAttemptsAllowedUntilDisabled(t *testing.T) {
	engine := &downloadV2Engine{}

	require.True(t, engine.directTransferDownloadAttemptAllowed())
	require.True(t, engine.directTransferDownloadAttemptAllowed())

	engine.disableDirectTransferDownload("direct_request_failed", nil)
	require.False(t, engine.directTransferDownloadAttemptAllowed())
}

func TestDownloadV2DirectFailureSuppressesLaterRanges(t *testing.T) {
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
	suppressor := &testDirectTransferDownloadSuppressor{}
	ctx := context.WithValue(context.Background(), directTransferDownloadSuppressorContextKey{}, suppressor)
	client := &Client{Config: files_sdk.Config{
		Logger: log.New(io.Discard, "", 0),
	}.Init()}
	closeBody := files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
		return body.Close()
	})

	_, err := client.Download(files_sdk.FileDownloadParams{File: sdkFile}, files_sdk.WithContext(ctx), closeBody)
	require.NoError(t, err)
	require.True(t, suppressor.disabled.Load())
	require.Equal(t, int32(1), suppressor.disables.Load())

	_, err = client.Download(files_sdk.FileDownloadParams{File: sdkFile}, files_sdk.WithContext(ctx), closeBody)
	require.NoError(t, err)
	require.Equal(t, int32(1), suppressor.disables.Load())
	require.Equal(t, int32(2), proxyRequests.Load())
}

func TestDownloadDirectTransfersDisabledUsesProxy(t *testing.T) {
	proxyRequests := atomic.Int32{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyRequests.Add(1)
		_, _ = w.Write([]byte("ok"))
	}))
	t.Cleanup(server.Close)

	suppressor := &testDirectTransferDownloadSuppressor{}
	ctx := context.WithValue(context.Background(), directTransferDownloadSuppressorContextKey{}, suppressor)
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
	require.Zero(t, suppressor.disables.Load())
	require.Equal(t, int32(1), proxyRequests.Load())
}

func TestDownloadV2DirectBodyReadFailureSuppressesLaterRanges(t *testing.T) {
	suppressor := &testDirectTransferDownloadSuppressor{}
	require.True(t, suppressor.directTransferDownloadAttemptAllowed())

	response := &http.Response{Body: &downloadV2DirectFailingBody{readErr: io.ErrUnexpectedEOF, closeErr: io.ErrClosedPipe}}
	_, err := files_sdk.BuildResponse(response, directTransferDownloadFailureOptions(suppressor, files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
		buf := make([]byte, 16)
		n, readErr := body.Read(buf)
		require.Equal(t, 7, n)
		require.ErrorIs(t, readErr, io.ErrUnexpectedEOF)
		require.ErrorIs(t, body.Close(), io.ErrClosedPipe)
		return nil
	}))...)

	require.NoError(t, err)
	require.True(t, suppressor.disabled.Load())
	require.Equal(t, int32(1), suppressor.disables.Load())
	require.Equal(t, "direct_body_read_failed", suppressor.reason.Load())
	require.False(t, suppressor.directTransferDownloadAttemptAllowed())
}

func TestDownloadV2DirectBodyEOFDoesNotSuppressLaterRanges(t *testing.T) {
	suppressor := &testDirectTransferDownloadSuppressor{}
	require.True(t, suppressor.directTransferDownloadAttemptAllowed())

	response := &http.Response{Body: &downloadV2DirectFailingBody{readErr: io.EOF}}
	_, err := files_sdk.BuildResponse(response, directTransferDownloadFailureOptions(suppressor, files_sdk.ResponseBodyOption(func(body io.ReadCloser) error {
		buf := make([]byte, 16)
		n, readErr := body.Read(buf)
		require.Equal(t, 7, n)
		require.ErrorIs(t, readErr, io.EOF)
		require.NoError(t, body.Close())
		return nil
	}))...)

	require.NoError(t, err)
	require.False(t, suppressor.disabled.Load())
	require.Equal(t, int32(0), suppressor.disables.Load())
	require.True(t, suppressor.directTransferDownloadAttemptAllowed())
}

type downloadV2DirectFailingBody struct {
	readErr  error
	closeErr error
	read     atomic.Bool
}

func (b *downloadV2DirectFailingBody) Read(p []byte) (int, error) {
	if !b.read.CompareAndSwap(false, true) {
		return 0, io.EOF
	}
	return copy(p, "partial"), b.readErr
}

func (b *downloadV2DirectFailingBody) Close() error {
	return b.closeErr
}
