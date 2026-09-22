package file

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeStage writes content to tmp the way a transfer does, and reports the
// identity that transfer would have recorded for it.
func writeStage(t *testing.T, tmp destinationPath, content []byte) fs.FileInfo {
	t.Helper()
	out, err := tmp.create()
	require.NoError(t, err)
	defer out.Close()
	_, err = out.Write(content)
	require.NoError(t, err)
	identity, err := out.Stat()
	require.NoError(t, err)
	return identity
}

// Publishing a download must deliver the bytes this transfer wrote. A
// filesystem can reach one file through more than one name, so another
// transfer can put its own file at the temporary path while this one runs;
// that file must never be delivered under this file's name.
func Test_finalizeTmpDownload_refusesToPublishAFileItDidNotWrite(t *testing.T) {
	dir := t.TempDir()
	final := explicitDestination(filepath.Join(dir, "victim.bin"))
	tmp, err := tmpDownloadPath(final, "")
	require.NoError(t, err)

	mine := writeStage(t, tmp, []byte("mine"))

	// Another file takes the temporary path, as a rename through a second name
	// for it would do.
	impostor := filepath.Join(dir, "impostor.bin")
	require.NoError(t, os.WriteFile(impostor, []byte("theirs"), 0600))
	require.NoError(t, os.Rename(impostor, tmp.String()))

	err = finalizeTmpDownload(t.Context(), tmp, final, mine)

	require.Error(t, err)
	assert.ErrorIs(t, err, errTempDownloadReplaced)
	assert.NoFileExists(t, final.String(), "a file this transfer did not write must not be published")

	// What the downloader does next, and what a later run then finds: nothing
	// left over, so the download stages under its usual name again rather than
	// under an alternate one from then on.
	require.NoError(t, removeTmpDownload(tmp))
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Empty(t, entries, "a refused download must leave nothing behind")
	next, err := tmpDownloadPath(final, "")
	require.NoError(t, err)
	assert.Equal(t, tmp, next, "the next attempt must stage under the usual name")
}

// The same contract through the real downloader: the temporary file is replaced
// while the transfer is running, the way a second transfer finalizing its own
// file onto that path does, so this transfer must not report success with the
// other file's bytes.
func TestClient_Downloader_doesNotPublishAStageReplacedMidTransfer(t *testing.T) {
	root := t.TempDir()
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	server.MockFiles["victim.bin"] = mockFile{SizeTrust: TrustedSizeValue, File: files_sdk.File{Size: 19999999}}

	local := filepath.Join(root, "victim.bin")
	planted := bytes.Repeat([]byte("EVIL"), 64)
	job := server.Client().Downloader(DownloaderParams{RemotePath: "victim.bin", LocalPath: local})
	var once sync.Once
	var swapErr error
	swapped := false
	job.RegisterFileEvent(func(JobFile) {
		once.Do(func() {
			swapErr = replaceStage(explicitDestination(local), planted)
			swapped = swapErr == nil
		})
	}, status.Downloading)
	job.Start()
	job.Wait()

	require.True(t, swapped, "the temporary file could not be replaced, so this checks nothing: %v", swapErr)
	require.Len(t, job.Statuses, 1)
	assert.NotEqual(t, status.Complete, job.Statuses[0].Status(), "a replaced temporary file must not be reported as a completed download")
	if delivered, err := os.ReadFile(local); err == nil {
		assert.NotEqual(t, planted, delivered, "another file's bytes were published as the requested file")
	}
	entries, err := os.ReadDir(root)
	require.NoError(t, err)
	assert.Empty(t, entries, "a refused download must leave nothing behind")
}

// A finished download is claimed under another name just before it is
// published. If the process stops between the two, what it leaves behind is
// still a temporary download of this client, so a sync running over the folder
// afterwards must not upload it, whatever rules the caller passes. That holds
// for a temporary file the caller named explicitly as well (a checkpoint's
// TmpPath), whose own folder is not in the reserved namespace.
func TestClient_Uploader_neverUploadsAClaimedDownload(t *testing.T) {
	for _, stage := range []struct {
		name string
		tmp  func(t *testing.T, root string) destinationPath
	}{
		{"the usual temporary file", func(t *testing.T, root string) destinationPath {
			tmp, err := tmpDownloadPath(explicitDestination(filepath.Join(root, "victim.bin")), "")
			require.NoError(t, err)
			return tmp
		}},
		{"an explicitly named temporary file from an earlier version", func(t *testing.T, root string) destinationPath {
			folder := filepath.Join(root, "legacy.download")
			require.NoError(t, os.Mkdir(folder, 0700))
			return explicitTmpDownload(filepath.Join(folder, "victim.bin"))
		}},
	} {
		t.Run(stage.name, func(t *testing.T) {
			root := t.TempDir()
			tmp := stage.tmp(t, root)
			mine := writeStage(t, tmp, []byte("complete but not yet published"))
			claimed, err := claimTmpDownload(tmp, mine)
			require.NoError(t, err)
			excluded, err := localPathReachesTempDownload(claimed.String())
			require.NoError(t, err)
			assert.True(t, excluded, "%v must be in the reserved namespace", claimed)
			require.NoError(t, os.WriteFile(filepath.Join(root, "report.csv"), []byte("an ordinary file"), 0600))

			server := (&MockAPIServer{T: t}).Do()
			defer server.Shutdown()
			job := server.Client().Uploader(UploaderParams{
				LocalPath:  root + string(os.PathSeparator),
				RemotePath: "reports",
				Manager:    manager.Build(1, 1),
				Ignore:     []string{"!.~files-cli*"},
				Include:    []string{"*"},
			})
			job.Start()
			job.Wait()

			assert.Equal(t, status.Complete, findStatus(t, job, "reports/report.csv").Status())
			var claimedUpload IFile
			for _, s := range job.Statuses {
				if s.LocalPath() == claimed.String() {
					claimedUpload = s
				}
				if s.Status().Is(status.Complete) {
					assert.Equalf(t, "reports/report.csv", s.RemotePath(), "%v must not be uploaded", s.LocalPath())
				}
			}
			require.NotNil(t, claimedUpload, "the claimed download must have been considered")
			assert.Equal(t, status.Ignored, claimedUpload.Status())
			for path := range server.MockFiles {
				assert.NotContains(t, path, filepath.Base(claimed.name), "the claimed download reached the server as %v", path)
			}
		})
	}
}

// replaceStage puts a different file at the temporary download path of final,
// through the same rooted rename a second transfer finalizing its own download
// onto that path would use. A plain rename is refused on Windows while the
// transfer holds the file open; the rooted one a transfer actually uses is not.
func replaceStage(final destinationPath, content []byte) error {
	stage, ok := existingTmpDownloadPath(final, "")
	if !ok {
		return errors.New("the transfer has no temporary file yet")
	}
	dir := filepath.Dir(stage.String())
	root, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer root.Close()
	const other = "other.bin"
	if err := root.WriteFile(other, content, 0600); err != nil {
		return err
	}
	return root.Rename(other, filepath.Base(stage.String()))
}
