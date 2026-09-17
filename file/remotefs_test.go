package file

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestFileReadUsesStatusOnlyWithRequestID(t *testing.T) {
	for _, requestID := range []string{"", "request-id"} {
		t.Run("request_id="+requestID, func(t *testing.T) {
			var downloadRequests, statusRequests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.URL.Path {
				case "/download":
					downloadRequests.Add(1)
					w.Header().Set("X-Files-Download-Request-Id", requestID)
					w.WriteHeader(http.StatusBadRequest)
					_, _ = io.WriteString(w, `{"error":"original download failure"}`)
				case "/download/request-id":
					statusRequests.Add(1)
					_, _ = io.WriteString(w, `{"error":"detailed status failure","data":{"status":"failed"}}`)
				default:
					t.Errorf("unexpected request: %s", r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			t.Cleanup(server.Close)
			config := files_sdk.Config{Logger: log.New(io.Discard, "", 0)}.Init()
			f := (&File{
				File: &files_sdk.File{Path: "/download.bin", DownloadUri: server.URL + "/download"},
				FS:   (&FS{Context: context.Background()}).Init(config, false),
			}).Init()

			_, err := f.Read(make([]byte, 1))

			require.Error(t, err)
			require.EqualValues(t, 1, downloadRequests.Load())
			if requestID == "" {
				require.ErrorContains(t, err, "original download failure")
				require.Zero(t, statusRequests.Load())
			} else {
				require.EqualError(t, err, "detailed status failure")
				require.EqualValues(t, 1, statusRequests.Load())
			}
		})
	}
}

func TestFS_Open(t *testing.T) {
	client, r, err := CreateClient("TestFS_Open")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	client.Upload(
		UploadWithReader(strings.NewReader("testing 3")),
		UploadWithSize(9),
		UploadWithDestinationPath(filepath.Join("remotefs_test", "1.text")),
	)

	fs := (&FS{}).Init(client.Config, true)
	fs = fs.WithContext(context.TODO()).(*FS)
	f, err := fs.Open("remotefs_test")
	assert.NoError(err)
	rf, ok := f.(*ReadDirFile)
	assert.True(ok)
	entry, err := rf.ReadDir(0)
	assert.NoError(err)
	assert.Equal(1, len(entry))
	assert.False(entry[0].IsDir())
	info, err := entry[0].Info()
	assert.NoError(err)
	assert.Equal("1.text", info.Name())
	fsFile, ok := entry[0].(*File)
	assert.True(ok)
	buf := make([]byte, 8)
	_, err = fsFile.Read(buf)
	assert.NoError(err)
	err = fsFile.Close()
	assert.NoError(err)
	assert.Equal("testing ", string(buf))

	buf = make([]byte, 9)
	_, err = fsFile.Read(buf)
	assert.NoError(err)
	err = fsFile.Close()
	assert.NoError(err)

	assert.Equal("testing 3", string(buf))

	_, err = fs.ReadDir(".")
	require.NoError(t, err)

	_, err = fs.Open(".")
	require.NoError(t, err)
}
