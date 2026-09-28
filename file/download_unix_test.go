//go:build darwin || linux

package file

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_DownloadToFileClosesDestinationAndResponseBody(t *testing.T) {
	// A finalizer could otherwise close a leaked descriptor before it is observed.
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	payload := bytes.Repeat([]byte{0, 0xff, '\r', '\n'}, 8192)
	truncated := payload[:len(payload)/4]
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(len(payload)))
		if r.URL.Path == "/truncated" {
			_, _ = w.Write(truncated)
			return
		}
		_, _ = w.Write(payload)
	}))
	t.Cleanup(server.Close)
	errRejected := errors.New("rejected by caller option")

	tests := []struct {
		name      string
		path      string
		opts      []files_sdk.RequestResponseOption
		wantErr   error
		wantBytes []byte
	}{
		{name: "complete", path: "/download", wantBytes: payload},
		{name: "truncated body", path: "/truncated", wantErr: io.ErrUnexpectedEOF, wantBytes: truncated},
		{
			name:      "request option error",
			path:      "/download",
			opts:      []files_sdk.RequestResponseOption{files_sdk.RequestOption(func(*http.Request) error { return errRejected })},
			wantErr:   errRejected,
			wantBytes: []byte{},
		},
		{
			name:      "response option error before copy",
			path:      "/download",
			opts:      []files_sdk.RequestResponseOption{files_sdk.ResponseOption(func(*http.Response) error { return errRejected })},
			wantErr:   errRejected,
			wantBytes: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transport := &closeRecordingTransport{RoundTripper: http.DefaultTransport.(*http.Transport).Clone()}
			t.Cleanup(transport.RoundTripper.(*http.Transport).CloseIdleConnections)
			client := &Client{Config: files_sdk.Config{EndpointOverride: server.URL}.Init().SetCustomClient(&http.Client{Transport: transport})}
			destination := filepath.Join(t.TempDir(), "download.bin")

			_, err := client.DownloadToFile(
				files_sdk.FileDownloadParams{File: files_sdk.File{Path: "remote.bin", DownloadUri: server.URL + tt.path}},
				destination,
				tt.opts...,
			)

			require.ErrorIs(t, err, tt.wantErr)
			assert.Empty(t, openDescriptorsFor(t, destination), "destination is still open")
			assert.Zero(t, transport.unclosedBodies(), "response body is still open")
			contents, err := os.ReadFile(destination)
			require.NoError(t, err)
			assert.Equal(t, tt.wantBytes, contents)
		})
	}
}

// openDescriptorsFor lists this process's open descriptors that refer to path.
func openDescriptorsFor(t *testing.T, path string) []string {
	t.Helper()
	info, err := os.Stat(path)
	require.NoError(t, err)
	target := info.Sys().(*syscall.Stat_t)
	fdDir, err := os.Open("/dev/fd")
	require.NoError(t, err)
	names, err := fdDir.Readdirnames(-1)
	require.NoError(t, fdDir.Close())
	require.NoError(t, err)

	var open []string
	for _, name := range names {
		fd, err := strconv.Atoi(name)
		if err != nil {
			continue
		}
		var stat syscall.Stat_t
		if syscall.Fstat(fd, &stat) == nil && stat.Dev == target.Dev && stat.Ino == target.Ino {
			open = append(open, name)
		}
	}
	return open
}

type closeRecordingTransport struct {
	http.RoundTripper
	mu     sync.Mutex
	bodies []*closeRecordingBody
}

func (t *closeRecordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.RoundTripper.RoundTrip(req)
	if err == nil {
		body := &closeRecordingBody{ReadCloser: res.Body}
		res.Body = body
		t.mu.Lock()
		t.bodies = append(t.bodies, body)
		t.mu.Unlock()
	}
	return res, err
}

func (t *closeRecordingTransport) unclosedBodies() (unclosed int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, body := range t.bodies {
		if !body.closed.Load() {
			unclosed++
		}
	}
	return unclosed
}

type closeRecordingBody struct {
	io.ReadCloser
	closed atomic.Bool
}

func (b *closeRecordingBody) Close() error {
	b.closed.Store(true)
	return b.ReadCloser.Close()
}
