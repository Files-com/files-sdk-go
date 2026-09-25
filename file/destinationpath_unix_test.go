//go:build !windows

package file

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"testing/fstest"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadFolder_acceptsBackslashInFileNameOnPosix(t *testing.T) {
	setup := NewTestSetup()
	defer setup.TearDown()
	setup.MapFS["reports"] = mapDir("reports")
	setup.MapFS["reports/entry"] = mapFileAt("reports/entry", `reports/a\b.txt`, 10)
	destination := filepath.Join(setup.tempDir, "downloads")
	setup.DownloaderParams = DownloaderParams{RemotePath: "reports/", LocalPath: destination + string(os.PathSeparator)}

	job := setup.Call()

	assert.Equal(t, status.Complete, findStatus(t, job, `reports/a\b.txt`).Status())
	assert.FileExists(t, filepath.Join(destination, `a\b.txt`))
}

// A 244-byte file name is longer than a temporary name can hold as it is, so
// it is staged under a shortened name. Finalizing it through an external temp
// directory must still deliver the file under its own name.
func TestClient_Downloader_externalTempPathFinalizesLongFileName(t *testing.T) {
	root := t.TempDir()
	temp := t.TempDir()
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	client := server.Client()
	name := strings.Repeat("a", 240) + ".txt"
	server.MockFiles[name] = mockFile{SizeTrust: TrustedSizeValue, File: files_sdk.File{Size: 100}}

	job := client.Downloader(DownloaderParams{RemotePath: name, LocalPath: root + string(os.PathSeparator), TempPath: temp})
	job.Start()
	job.Wait()

	require.Len(t, job.Statuses, 1)
	require.NoError(t, job.Statuses[0].Err())
	assert.FileExists(t, filepath.Join(root, name))
	entries, err := os.ReadDir(temp)
	require.NoError(t, err)
	assert.Empty(t, entries)
}

// Some test environments run as root or ignore mode bits. Skip only when a
// write probe proves that they cannot exercise directory write permissions.
func makeDownloadRootReadOnly(t *testing.T, root string) {
	t.Helper()
	require.NoError(t, os.Chmod(root, 0555))
	t.Cleanup(func() { _ = os.Chmod(root, 0755) })
	probe := filepath.Join(root, "write-probe")
	err := os.WriteFile(probe, nil, 0600)
	if err == nil {
		require.NoError(t, os.Remove(probe))
		t.Skip("directory write permissions are not enforced")
	}
	require.ErrorIs(t, err, fs.ErrPermission)
}

// A fresh download is created the way files are usually created, with 0666
// before the process umask, so a group-writable destination folder keeps
// working the way it always has. Umask 0 makes the whole mode observable.
func TestClient_Downloader_freshFileKeepsTheUsualCreationMode(t *testing.T) {
	previous := syscall.Umask(0)
	t.Cleanup(func() { syscall.Umask(previous) })
	root := t.TempDir()
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	mockFolder(server, "reports", map[string][]byte{"report.csv": []byte("report bytes")})

	job := server.Client().Downloader(DownloaderParams{RemotePath: "reports/report.csv", LocalPath: filepath.Join(root, "report.csv")})
	job.Start()
	job.Wait()

	require.NoError(t, job.Statuses[0].Err())
	info, err := os.Stat(filepath.Join(root, "report.csv"))
	require.NoError(t, err)
	assert.Equal(t, fs.FileMode(0666), info.Mode().Perm())
}

func TestClient_Downloader_writableChildOfReadOnlyRoot(t *testing.T) {
	for _, externalTemp := range []bool{false, true} {
		name := "default_temp"
		if externalTemp {
			name = "external_temp"
		}
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			child := filepath.Join(root, "sub")
			require.NoError(t, os.Mkdir(child, 0755))
			// Include a long basename and an existing destination so this also
			// checks that staging does not append to the final name or content.
			filename := strings.Repeat("a", 240) + ".txt"
			final := filepath.Join(child, filename)
			require.NoError(t, os.WriteFile(final, []byte("old content"), 0600))
			makeDownloadRootReadOnly(t, root)
			server := (&MockAPIServer{T: t}).Do()
			defer server.Shutdown()
			server.MockFiles["folder"] = mockFile{File: files_sdk.File{Path: "folder", DisplayName: "folder", Type: "directory"}}
			server.MockFiles["folder/sub"] = mockFile{File: files_sdk.File{Path: "folder/sub", DisplayName: "sub", Type: "directory"}}
			remote := "folder/sub/" + filename
			data := []byte("new content")
			server.MockFiles[remote] = mockFile{SizeTrust: TrustedSizeValue, File: files_sdk.File{Path: remote, DisplayName: filename, Type: "file", Size: int64(len(data))}, Data: data}
			params := DownloaderParams{RemotePath: "folder/", LocalPath: root}
			if externalTemp {
				params.TempPath = t.TempDir()
			}
			job := server.Client().Downloader(params)
			job.Start()
			job.Wait()

			for _, s := range job.Statuses {
				require.NoError(t, s.Err(), s.RemotePath())
			}
			assert.Equal(t, status.Complete, findStatus(t, job, remote).Status())
			actual, err := os.ReadFile(final)
			require.NoError(t, err)
			assert.Equal(t, data, actual)
			entries, err := os.ReadDir(child)
			require.NoError(t, err)
			require.Len(t, entries, 1, "the final parent must not keep a staging file")
			assert.Equal(t, filename, entries[0].Name())
			if externalTemp {
				entries, err := os.ReadDir(params.TempPath)
				require.NoError(t, err)
				assert.Empty(t, entries, "the external temp and its macOS package must be removed")
			}
		})
	}
}

func TestMoveToReadOnlyRootPreservesSourceOnFailure(t *testing.T) {
	for _, failure := range []string{"destination_is_directory", "parent_links_outside", "paused"} {
		t.Run(failure, func(t *testing.T) {
			root := t.TempDir()
			child := filepath.Join(root, "sub")
			outside := t.TempDir()
			if failure == "parent_links_outside" {
				require.NoError(t, os.Symlink(outside, child))
			} else {
				require.NoError(t, os.Mkdir(child, 0755))
			}
			final := destinationPath{dir: root, name: filepath.Join("sub", "file.txt")}
			if failure == "destination_is_directory" {
				require.NoError(t, os.Mkdir(final.String(), 0755))
				require.NoError(t, os.WriteFile(filepath.Join(final.String(), "keep.txt"), []byte("keep"), 0600))
			} else {
				require.NoError(t, os.WriteFile(final.String(), []byte("keep"), 0600))
			}
			makeDownloadRootReadOnly(t, root)
			tmp, err := tmpDownloadPath(final, t.TempDir())
			require.NoError(t, err)
			require.NoError(t, os.WriteFile(tmp.String(), []byte("completed download"), 0640))
			ctx, cancel := context.WithCancelCause(context.Background())
			defer cancel(nil)
			if failure == "paused" {
				cancel(ErrJobPaused)
			}

			err = tmp.moveTo(ctx, final)

			require.Error(t, err)
			if failure == "paused" {
				assert.ErrorIs(t, err, ErrJobPaused)
			}
			data, err := os.ReadFile(tmp.String())
			require.NoError(t, err)
			assert.Equal(t, "completed download", string(data))
			keep := final.String()
			if failure == "destination_is_directory" {
				keep = filepath.Join(keep, "keep.txt")
			}
			data, err = os.ReadFile(keep)
			require.NoError(t, err)
			assert.Equal(t, "keep", string(data))
			entries, err := os.ReadDir(child)
			require.NoError(t, err)
			require.Len(t, entries, 1, "a failed finalization must not keep a staging file")
			assert.Equal(t, "file.txt", entries[0].Name())
		})
	}
}

// Cancel at a real filesystem boundary without depending on copy speed or
// goroutine scheduling. The context is passed directly to the transfer item,
// because the public Downloader wraps its caller's context in a child context.
type interruptStagedDownloadContext struct {
	context.Context
	cancel context.CancelCauseFunc
	cause  error
	dir    string
}

func (c *interruptStagedDownloadContext) Err() error {
	entries, _ := os.ReadDir(c.dir)
	for _, entry := range entries {
		if entry.Name() == "file.txt" {
			continue
		}
		if info, err := entry.Info(); err == nil && info.Size() > 0 {
			c.cancel(c.cause)
		}
	}
	return c.Context.Err()
}

func TestDownloadItem_interruptReadOnlyRootFinalization(t *testing.T) {
	for _, transport := range []string{"ordinary", "resumed", "zip_stream", "zip_spool"} {
		for _, interruption := range []string{"pause", "cancel"} {
			t.Run(transport+"/"+interruption, func(t *testing.T) {
				root := t.TempDir()
				child := filepath.Join(root, "sub")
				require.NoError(t, os.Mkdir(child, 0755))
				final := destinationPath{dir: root, name: filepath.Join("sub", "file.txt")}
				require.NoError(t, os.WriteFile(final.String(), []byte("old content"), 0600))
				makeDownloadRootReadOnly(t, root)
				temp := t.TempDir()
				data := bytes.Repeat([]byte("completed download\n"), 4096)
				const remote = "folder/sub/file.txt"
				input := fstest.MapFS{
					"folder":     mapDir("folder"),
					"folder/sub": mapDir("folder/sub"),
					remote:       mapFile(remote, len(data)),
				}
				input[remote].Data = data
				params := DownloaderParams{RemotePath: "folder/", LocalPath: root, TempPath: temp, config: files_sdk.Config{}.Init()}
				job := downloader(context.Background(), input, params)
				defer job.Cancel()
				file, err := input.Open(remote)
				require.NoError(t, err)
				defer file.Close()
				createIndexedStatus(Entity{File: file, FS: input}, params, job)
				s := job.Statuses[0].(*DownloadStatus)
				baseCtx, cancel := context.WithCancelCause(context.Background())
				defer cancel(nil)
				cause := context.Canceled
				if interruption == "pause" {
					cause = ErrJobPaused
				}
				ctx := &interruptStagedDownloadContext{Context: baseCtx, cancel: cancel, cause: cause, dir: child}

				if transport == "ordinary" || transport == "resumed" {
					if transport == "resumed" {
						// Also exercise the completed-temp resume branch, including
						// preservation of a resumed file's custom permissions.
						tmp, err := tmpDownloadPath(final, temp)
						require.NoError(t, err)
						require.NoError(t, os.WriteFile(tmp.String(), data, 0640))
						// Set the custom mode explicitly, independent of the process umask.
						require.NoError(t, os.Chmod(tmp.String(), 0640))
						// Only a paused temporary download is published without content.
						parkStage(t, tmp)
					}
					runDownloadFolderItem(ctx, s)
				} else {
					archive := zipWithEntry("file.txt", data)
					signaled := 0
					finalized := func(*DownloadStatus) { signaled++ }
					if transport == "zip_stream" {
						expected, mapErr := zipBatchEntryMap([]*DownloadStatus{s})
						require.NoError(t, mapErr)
						_, err = extractZipBatchStream(ctx, bytes.NewReader(archive), expected, finalized, nil)
					} else {
						spool := explicitDestination(filepath.Join(t.TempDir(), "archive.zip"))
						require.NoError(t, os.WriteFile(spool.String(), archive, 0600))
						batch := &zipBatchDownloader{job: job, ctx: ctx}
						_, err = batch.extractZipSpool(spool, []*DownloadStatus{s}, finalized)
					}
					require.NoError(t, err)
					assert.Equal(t, 1, signaled, "canceled ZIP entries still complete job accounting exactly once")
					assert.Zero(t, job.ZipBatchStats().CleanFinalized, "an interrupted copy is not a completed ZIP entry")
				}

				require.ErrorIs(t, context.Cause(ctx), cause, "the copy must be interrupted after staging data")
				assert.Equal(t, status.Canceled, s.Status())
				actual, err := os.ReadFile(final.String())
				require.NoError(t, err)
				assert.Equal(t, "old content", string(actual), "an interrupted copy must not replace the destination")
				entries, err := os.ReadDir(child)
				require.NoError(t, err)
				require.Len(t, entries, 1, "the interrupted copy must remove its staging file")
				if interruption == "cancel" {
					entries, err := os.ReadDir(temp)
					require.NoError(t, err)
					assert.Empty(t, entries, "normal cancel removes the original temp")
					return
				}

				tmp, exists := pausedTmpDownloadPath(final, temp)
				require.True(t, exists, "pause keeps the original completed download")
				info, err := tmp.stat()
				require.NoError(t, err)
				assert.Equal(t, int64(len(data)), info.Size())
				originalMode := info.Mode().Perm()
				if transport == "resumed" {
					assert.Equal(t, fs.FileMode(0640), originalMode, "pause preserves the temp file's permissions")
				}
				// A completed temp must be reused, even if the remote bytes have
				// since changed without changing the listed size.
				input[remote].Data = bytes.Repeat([]byte("x"), len(data))
				resumed := downloader(context.Background(), input, params)
				resumed.Start()
				resumed.Wait()
				assert.Equal(t, status.Complete, findStatus(t, resumed, remote).Status())
				info, err = os.Stat(final.String())
				require.NoError(t, err)
				assert.Equal(t, originalMode, info.Mode().Perm(), "copying preserves the downloaded file's permissions")
				actual, err = os.ReadFile(final.String())
				require.NoError(t, err)
				assert.Equal(t, data, actual, "resume finalizes the original completed bytes without redownloading")
				entries, err = os.ReadDir(temp)
				require.NoError(t, err)
				assert.Empty(t, entries)
			})
		}
	}
}
