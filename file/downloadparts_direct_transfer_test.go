package file

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"os"
	"sync"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/stretchr/testify/require"
)

func TestDownloadPartsReusesDirectTransferClient(t *testing.T) {
	const size = int64(11 * 1024 * 1024)
	caPEM, _ := testDirectTransferCertificate(t, testDirectCredentialServerName)
	state := &downloadPartsDirectClientState{}
	fileInfo := Info{
		File:      files_sdk.File{DisplayName: "download.bin", Size: size},
		sizeTrust: TrustedSizeValue,
	}
	remote := &downloadPartsDirectClientFile{
		state: state,
		info:  fileInfo,
		config: files_sdk.Config{
			Environment: files_sdk.Development,
		}.Init(),
		directInfo: files_sdk.DirectConnectionInfo{
			Version:    1,
			ServerName: testDirectCredentialServerName,
			Addresses:  []string{"tcp://127.0.0.1:4001"},
			DirectUri:  "https://" + testDirectCredentialServerName + "/downloads",
			CaPem:      string(caPEM),
		},
	}
	destination, err := os.CreateTemp(t.TempDir(), "download-*")
	require.NoError(t, err)

	download := (&DownloadParts{}).Init(
		remote,
		fileInfo,
		manager.ConcurrencyManager{}.New(2),
		destination,
		remote.config,
		0,
	)
	require.NoError(t, download.Run(context.Background()))

	state.mu.Lock()
	defer state.mu.Unlock()
	require.Greater(t, len(state.clients), 1)
	for _, client := range state.clients[1:] {
		require.Same(t, state.clients[0], client)
	}
}

type downloadPartsDirectClientState struct {
	mu      sync.Mutex
	clients []*http.Client
}

type downloadPartsDirectClientFile struct {
	state      *downloadPartsDirectClientState
	info       Info
	config     files_sdk.Config
	directInfo files_sdk.DirectConnectionInfo
	ctx        context.Context
}

func (f *downloadPartsDirectClientFile) WithContext(ctx context.Context) fs.File {
	copy := *f
	copy.ctx = ctx
	return &copy
}

func (f *downloadPartsDirectClientFile) ReaderRange(off, end int64) (io.ReadCloser, error) {
	_, client, err := files_sdk.DirectTransferRetryableClient(f.ctx, f.config, f.directInfo)
	if err != nil {
		return nil, err
	}
	f.state.mu.Lock()
	f.state.clients = append(f.state.clients, client.HTTPClient)
	f.state.mu.Unlock()
	return io.NopCloser(io.LimitReader(downloadPartsZeroReader{}, end-off+1)), nil
}

func (f *downloadPartsDirectClientFile) Stat() (fs.FileInfo, error) { return f.info, nil }
func (f *downloadPartsDirectClientFile) Read([]byte) (int, error)   { return 0, io.EOF }
func (f *downloadPartsDirectClientFile) Close() error               { return nil }

type downloadPartsZeroReader struct{}

func (downloadPartsZeroReader) Read(p []byte) (int, error) {
	clear(p)
	return len(p), nil
}
