//go:build darwin

package file

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// On macOS an external temporary file lives in a temporary download folder
// inside the temp directory. Replacing that folder with a link to another tree
// between the write and the finalization must not move that tree's file into
// the destination.
func TestClient_Downloader_externalTempFolderReplacedByLinkDoesNotMoveOtherTree(t *testing.T) {
	root := t.TempDir()
	temp := t.TempDir()
	other := t.TempDir()
	decoy := filepath.Join(other, "file.txt")
	require.NoError(t, os.WriteFile(decoy, []byte("decoy"), 0644))
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	client := server.Client()
	server.MockFiles["file.txt"] = mockFile{SizeTrust: TrustedSizeValue, File: files_sdk.File{Size: 19999999}}

	job := client.Downloader(DownloaderParams{RemotePath: "file.txt", LocalPath: root + string(os.PathSeparator), TempPath: temp})
	var swap sync.Once
	job.RegisterFileEvent(func(JobFile) {
		swap.Do(func() {
			folder := filepath.Join(temp, tmpDownloadElement("file.txt", ""))
			assert.NoError(t, os.Rename(folder, filepath.Join(temp, "moved-folder")))
			assert.NoError(t, os.Symlink(other, folder))
		})
	}, status.Downloading)
	job.Start()
	job.Wait()

	require.Len(t, job.Statuses, 1)
	assert.Equal(t, status.Errored, job.Statuses[0].Status())
	assert.Error(t, job.Statuses[0].Err())
	content, err := os.ReadFile(decoy)
	require.NoError(t, err, "the other tree's file must stay where it is")
	assert.Equal(t, "decoy", string(content))
	assert.NoFileExists(t, filepath.Join(root, "file.txt"))
}
