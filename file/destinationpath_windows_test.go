//go:build windows

package file

import (
	"os"
	"path/filepath"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Reserved device names without an extension are rejected on every Windows
// version. A reserved name with an extension, such as NUL.txt, is a valid file
// name since Windows 11 and is decided by the running system, so it is not
// asserted here.
func TestDownloadFolder_rejectsWindowsPathForms(t *testing.T) {
	cases := map[string]string{
		"backslash traversal":   `reports/..\..\outside.txt`,
		"drive prefix":          `reports/C:outside.txt`,
		"drive directory":       `reports/E:/dropped2.txt`,
		"alternate data stream": `reports/file.txt:stream`,
		"UNC path":              `reports/\\server\share\outside.txt`,
		"device name":           `reports/CON`,
		"lowercase device name": `reports/nul`,
		"serial port name":      `reports/COM1`,
		"console handle name":   `reports/CONIN$`,
	}
	for name, serverPath := range cases {
		t.Run(name, func(t *testing.T) {
			setup := NewTestSetup()
			defer setup.TearDown()
			setup.MapFS["reports"] = mapDir("reports")
			setup.MapFS["reports/entry"] = mapFileAt("reports/entry", serverPath, 10)
			destination := filepath.Join(setup.tempDir, "downloads")
			setup.DownloaderParams = DownloaderParams{
				RemotePath:  "reports/",
				LocalPath:   destination + string(os.PathSeparator),
				RetryPolicy: RetryPolicy{Type: RetryUnfinished, RetryCount: 2},
			}

			job := setup.Call()

			rejected := findStatus(t, job, serverPath)
			assert.Equal(t, status.Errored, rejected.Status())
			require.ErrorContains(t, rejected.Err(), "not a valid local path")
			var classified interface{ ErrorType() string }
			require.ErrorAs(t, ToStatusFile(rejected).Err, &classified)
			assert.Equal(t, "invalid_path", classified.ErrorType())
			assert.Zero(t, rejected.StatusChanges().Count(status.Retrying))
			entries, err := os.ReadDir(destination)
			require.NoError(t, err)
			assert.Empty(t, entries, "nothing may be written for the rejected path")
			assert.True(t, job.All(status.Ended...))
		})
	}
}

// With no LocalPath, a download is written to the working directory under the
// remote file's or folder's name. A name with a drive prefix must not become a
// path on that drive or lose its prefix, and a stream name must not open
// another file. The drive prefix is the working directory's own drive, whose
// current directory is the working directory, so a download that escaped
// would replace victim.txt or add an entry beside it.
func TestClient_Downloader_rejectsDefaultLocalNameWindowsReserves(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	drive := filepath.VolumeName(root)
	require.Len(t, drive, 2, "the working directory must be on a drive letter")
	cases := map[string]struct {
		remotePath string
		mock       func(server *MockAPIServer)
	}{
		"file with a drive prefix at the root": {drive + "victim.txt", func(server *MockAPIServer) {
			server.MockFiles[drive+"victim.txt"] = mockFile{
				SizeTrust: TrustedSizeValue,
				Data:      []byte("replacement"),
				File:      files_sdk.File{Path: drive + "victim.txt", DisplayName: drive + "victim.txt", Type: "file", Size: int64(len("replacement"))},
			}
		}},
		"file with a drive prefix in a folder": {"reports/" + drive + "victim.txt", func(server *MockAPIServer) {
			mockFolder(server, "reports", map[string][]byte{drive + "victim.txt": []byte("replacement")})
		}},
		"file with a stream name": {"reports/victim.txt::$DATA", func(server *MockAPIServer) {
			mockFolder(server, "reports", map[string][]byte{"victim.txt::$DATA": []byte("replacement")})
		}},
		"folder with a drive prefix": {"reports/" + drive + "stuff", func(server *MockAPIServer) {
			mockFolder(server, "reports/"+drive+"stuff", map[string][]byte{"victim.txt": []byte("replacement")})
		}},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.WriteFile("victim.txt", []byte("original"), 0644))
			server := (&MockAPIServer{T: t}).Do()
			defer server.Shutdown()
			c.mock(server)

			job := server.Client().Downloader(DownloaderParams{
				RemotePath:  c.remotePath,
				RetryPolicy: RetryPolicy{Type: RetryUnfinished, RetryCount: 2},
			})
			job.Start()
			job.Wait()

			rejected := findStatus(t, job, c.remotePath)
			assert.Equal(t, status.Errored, rejected.Status())
			var classified interface{ ErrorType() string }
			require.ErrorAs(t, ToStatusFile(rejected).Err, &classified)
			assert.Equal(t, "invalid_path", classified.ErrorType())
			assert.Zero(t, rejected.StatusChanges().Count(status.Retrying))
			content, err := os.ReadFile("victim.txt")
			require.NoError(t, err)
			assert.Equal(t, "original", string(content), "the download must not replace another local file")
			entries, err := os.ReadDir(".")
			require.NoError(t, err)
			assert.Len(t, entries, 1, "nothing may be written beside victim.txt")
		})
	}
}
