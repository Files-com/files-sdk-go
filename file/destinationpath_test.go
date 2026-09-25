package file

import (
	"bytes"
	"io/fs"
	"net/http"
	"os"
	pathpkg "path"
	"path/filepath"
	"testing"
	"testing/fstest"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mapDir(path string) *fstest.MapFile {
	return &fstest.MapFile{
		Mode: fs.ModeDir,
		Sys:  files_sdk.File{DisplayName: pathpkg.Base(path), Path: path, Type: "directory"},
	}
}

func mapFile(path string, size int) *fstest.MapFile {
	return mapFileAt(path, path, size)
}

// mapFileAt lists a file under walkPath while its server response reports
// serverPath.
func mapFileAt(walkPath, serverPath string, size int) *fstest.MapFile {
	return &fstest.MapFile{
		Data: make([]byte, size),
		Mode: fs.ModePerm,
		Sys:  files_sdk.File{DisplayName: pathpkg.Base(walkPath), Path: serverPath, Type: "file", Size: int64(size)},
	}
}

func findStatus(t *testing.T, job *Job, remotePath string) IFile {
	t.Helper()
	for _, s := range job.Statuses {
		if s.RemotePath() == remotePath {
			return s
		}
	}
	t.Fatalf("no status for %q", remotePath)
	return nil
}

func TestClient_Downloader_rejectsServerPathOutsideDestination(t *testing.T) {
	for _, policy := range []RetryPolicyType{"no retry", RetryUnfinished, RetryAll} {
		t.Run(string(policy), func(t *testing.T) {
			root := t.TempDir()
			server := (&MockAPIServer{T: t}).Do()
			defer server.Shutdown()
			client := server.Client()
			server.MockFiles["reports"] = mockFile{File: files_sdk.File{Path: "reports", DisplayName: "reports", Type: "directory"}}
			server.MockFiles["reports/inside.txt"] = mockFile{
				SizeTrust: TrustedSizeValue,
				File:      files_sdk.File{Path: "reports/inside.txt", DisplayName: "inside.txt", Type: "file", Size: 10},
			}
			// The listing names this entry alias.txt, so the walk looks up
			// reports/alias.txt, whose metadata reports a path outside the folder.
			server.MockFiles["reports/real.txt"] = mockFile{
				SizeTrust: TrustedSizeValue,
				File:      files_sdk.File{Path: "reports/real.txt", DisplayName: "alias.txt", Type: "file", Size: 10},
			}
			server.MockRoute("/api/rest/v1/file_actions/metadata/reports/alias.txt", func(c *gin.Context, _ interface{}) bool {
				c.JSON(http.StatusOK, files_sdk.File{Path: "reports/../../outside.txt", DisplayName: "alias.txt", Type: "file", Size: 10})
				return true
			})
			destination := filepath.Join(root, "a", "b", "downloads")

			job := client.Downloader(DownloaderParams{
				RemotePath:  "reports/",
				LocalPath:   destination + string(os.PathSeparator),
				RetryPolicy: RetryPolicy{Type: policy, RetryCount: 2},
			})
			job.Start()
			job.Wait()

			rejected := findStatus(t, job, "reports/../../outside.txt")
			assert.Equal(t, status.Errored, rejected.Status())
			require.ErrorContains(t, rejected.Err(), `download: server path "reports/../../outside.txt"`)
			var classified interface {
				ErrorType() string
				PublicError() string
			}
			require.ErrorAs(t, ToStatusFile(rejected).Err, &classified)
			assert.Equal(t, "invalid_path", classified.ErrorType())
			assert.Equal(t, "Cannot download this file because its path is not valid on this computer.", classified.PublicError())
			assert.Zero(t, rejected.StatusChanges().Count(status.Retrying), "an invalid destination cannot be fixed by retrying")
			assert.NoFileExists(t, filepath.Join(root, "a", "outside.txt"))
			assert.Equal(t, status.Complete, findStatus(t, job, "reports/inside.txt").Status())
			assert.FileExists(t, filepath.Join(destination, "inside.txt"))
		})
	}
}

func TestClient_Downloader_keepsFolderListedInDifferentCaseInsideDestination(t *testing.T) {
	root := t.TempDir()
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	client := server.Client()
	server.MockFiles["Reports"] = mockFile{File: files_sdk.File{Path: "Reports", DisplayName: "Reports", Type: "directory"}}
	server.MockFiles["reports/x.txt"] = mockFile{
		SizeTrust: TrustedSizeValue,
		File:      files_sdk.File{Path: "reports/x.txt", DisplayName: "x.txt", Type: "file", Size: 10},
	}
	destination := filepath.Join(root, "downloads")

	job := client.Downloader(DownloaderParams{RemotePath: "Reports/", LocalPath: destination + string(os.PathSeparator)})
	job.Start()
	job.Wait()

	for _, s := range job.Statuses {
		require.NoError(t, s.Err(), s.RemotePath())
	}
	assert.FileExists(t, filepath.Join(destination, "x.txt"))
	assert.NoDirExists(t, filepath.Join(root, "reports"), "the folder must not become a sibling of the destination")
}

func TestDownloadFolder_acceptsServerPathWithLeadingSlash(t *testing.T) {
	setup := NewTestSetup()
	defer setup.TearDown()
	setup.MapFS["reports"] = mapDir("reports")
	setup.MapFS["reports/file.txt"] = mapFileAt("reports/file.txt", "/reports/file.txt", 10)
	setup.DownloaderParams = DownloaderParams{RemotePath: "reports", LocalPath: setup.tempDir + string(os.PathSeparator)}

	job := setup.Call()

	assert.Equal(t, status.Complete, findStatus(t, job, "/reports/file.txt").Status())
	assert.FileExists(t, filepath.Join(setup.tempDir, "reports", "file.txt"))
}

func TestDownloadFolder_doesNotWriteThroughLinkLeavingDestination(t *testing.T) {
	setup := NewTestSetup()
	defer setup.TearDown()
	setup.MapFS["folder"] = mapDir("folder")
	setup.MapFS["folder/sub"] = mapDir("folder/sub")
	setup.MapFS["folder/sub/file.txt"] = mapFile("folder/sub/file.txt", 100)
	other := t.TempDir()
	destination := filepath.Join(setup.tempDir, "folder")
	require.NoError(t, os.MkdirAll(destination, 0755))
	if err := os.Symlink(other, filepath.Join(destination, "sub")); err != nil {
		t.Skipf("symlinks are not available: %v", err)
	}
	setup.DownloaderParams = DownloaderParams{RemotePath: "folder", LocalPath: setup.tempDir + string(os.PathSeparator)}

	job := setup.Call()

	entries, err := os.ReadDir(other)
	require.NoError(t, err)
	assert.Empty(t, entries, "nothing may be written through the link")
	file := findStatus(t, job, "folder/sub/file.txt")
	assert.Equal(t, status.Errored, file.Status())
	assert.Error(t, file.Err())
	assert.Equal(t, status.Errored, findStatus(t, job, "folder/sub").Status())
	assert.True(t, job.All(status.Ended...))
}

func TestDownloadFolder_externalTempPathFinalizesNestedFiles(t *testing.T) {
	setup := NewTestSetup()
	defer setup.TearDown()
	setup.MapFS["folder"] = mapDir("folder")
	setup.MapFS["folder/sub"] = mapDir("folder/sub")
	setup.MapFS["folder/sub/file.txt"] = mapFile("folder/sub/file.txt", 100)
	temp := t.TempDir()
	setup.DownloaderParams = DownloaderParams{RemotePath: "folder", LocalPath: setup.tempDir + string(os.PathSeparator), TempPath: temp}

	job := setup.Call()

	for _, s := range job.Statuses {
		require.NoError(t, s.Err(), s.RemotePath())
	}
	stat, err := os.Stat(filepath.Join(setup.tempDir, "folder", "sub", "file.txt"))
	require.NoError(t, err)
	assert.Equal(t, int64(100), stat.Size())
	tempEntries, err := os.ReadDir(temp)
	require.NoError(t, err)
	assert.Empty(t, tempEntries, "the temp directory must not keep the temporary file")
	destinationEntries, err := os.ReadDir(filepath.Join(setup.tempDir, "folder"))
	require.NoError(t, err)
	require.Len(t, destinationEntries, 1, "the destination must not keep a placeholder")
	assert.Equal(t, "sub", destinationEntries[0].Name())
}

func TestClient_Downloader_rejectsListingNameThatRedirectsWalk(t *testing.T) {
	root := t.TempDir()
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	client := server.Client()
	server.MockFiles["reports"] = mockFile{File: files_sdk.File{Path: "reports", DisplayName: "reports", Type: "directory"}}
	server.MockFiles["reports/evil"] = mockFile{File: files_sdk.File{Path: "reports/evil", DisplayName: "../secrets", Type: "directory"}}
	server.MockFiles["secrets"] = mockFile{File: files_sdk.File{Path: "secrets", DisplayName: "secrets", Type: "directory"}}
	server.MockFiles["secrets/key.txt"] = mockFile{
		SizeTrust: TrustedSizeValue,
		File:      files_sdk.File{Path: "secrets/key.txt", DisplayName: "key.txt", Type: "file", Size: 5},
		Data:      []byte("shhh!"),
	}

	job := client.Downloader(DownloaderParams{RemotePath: "reports", LocalPath: root + string(os.PathSeparator)})
	job.Start()
	job.Wait()

	assert.NoFileExists(t, filepath.Join(root, "secrets", "key.txt"))
	assert.NotContains(t, server.DownloadRequests, "secrets/key.txt")
	errored, ok := job.Find(status.Errored)
	require.True(t, ok, "the malformed listing must be reported")
	assert.ErrorContains(t, errored.Err(), `invalid name "../secrets"`)
	assert.True(t, job.All(status.Ended...))
}

func TestClient_Downloader_resumesExplicitTmpPathOutsideDestination(t *testing.T) {
	root := t.TempDir()
	elsewhere := t.TempDir()
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	client := server.Client()
	fileSize := int64(19999999)
	server.MockFiles["file.txt"] = mockFile{SizeTrust: TrustedSizeValue, File: files_sdk.File{Size: fileSize}}
	partial := bytes.Repeat([]byte("x"), int(fileSize/2))
	tmpPath, err := tmpDownloadPath(explicitDestination(filepath.Join(elsewhere, "file.txt")), "")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(tmpPath.String(), partial, 0644))
	// A checkpoint names the file as the paused run reported it.
	paused := parkStage(t, tmpPath)
	localPath := filepath.Join(root, "file.txt")

	job := client.Downloader(DownloaderParams{RemotePath: "file.txt", LocalPath: localPath, ResumeTmpPath: paused.String()})
	job.Start()
	job.Wait()

	require.Len(t, job.Statuses, 1)
	require.NoError(t, job.Statuses[0].Err())
	written, err := os.ReadFile(localPath)
	require.NoError(t, err)
	require.Len(t, written, int(fileSize))
	assert.Equal(t, partial, written[:len(partial)], "the partial data in the explicit temp file must be kept")
	entries, err := os.ReadDir(elsewhere)
	require.NoError(t, err)
	assert.Empty(t, entries, "the explicit temp file must be cleaned up after finalization")
}

func TestClient_Downloader_reportsListingEntryWithMalformedPath(t *testing.T) {
	for _, policy := range []RetryPolicyType{"no retry", RetryUnfinished, RetryAll} {
		t.Run(string(policy), func(t *testing.T) {
			root := t.TempDir()
			server := (&MockAPIServer{T: t}).Do()
			defer server.Shutdown()
			client := server.Client()
			server.MockFiles["reports"] = mockFile{File: files_sdk.File{Path: "reports", DisplayName: "reports", Type: "directory"}}
			server.MockFiles["reports/inside.txt"] = mockFile{
				SizeTrust: TrustedSizeValue,
				File:      files_sdk.File{Path: "reports/inside.txt", DisplayName: "inside.txt", Type: "file", Size: 10},
			}
			// The listing of reports returns an entry whose path leaves the folder.
			server.MockFiles["reports/escape.txt"] = mockFile{
				SizeTrust: TrustedSizeValue,
				File:      files_sdk.File{Path: "reports/../../outside.txt", DisplayName: "escape.txt", Type: "file", Size: 10},
			}
			destination := filepath.Join(root, "a", "b", "downloads")

			job := client.Downloader(DownloaderParams{
				RemotePath:  "reports/",
				LocalPath:   destination + string(os.PathSeparator),
				RetryPolicy: RetryPolicy{Type: policy, RetryCount: 2},
			})
			job.Start()
			job.Wait()

			errored, ok := job.Find(status.Errored)
			require.True(t, ok, "the bad listing entry must be reported")
			require.ErrorContains(t, errored.Err(), `listing entry "reports/../../outside.txt" has a malformed path`)
			retries := 0
			if policy != "no retry" {
				retries = 2
			}
			assert.Equal(t, retries, errored.StatusChanges().Count(status.Retrying), "the error must survive every retry")
			assert.NoFileExists(t, filepath.Join(root, "a", "outside.txt"))
			assert.NotContains(t, server.DownloadRequests, "reports/../../outside.txt")
			assert.True(t, job.All(status.Ended...))
		})
	}
}

func TestClient_Downloader_listsNestedListingEntriesWithTheirOwnFolder(t *testing.T) {
	root := t.TempDir()
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	client := server.Client()
	server.MockFiles["reports"] = mockFile{File: files_sdk.File{Path: "reports", DisplayName: "reports", Type: "directory"}}
	server.MockFiles["reports/sub"] = mockFile{File: files_sdk.File{Path: "reports/sub", DisplayName: "sub", Type: "directory"}}
	server.MockFiles["reports/sub/deep.txt"] = mockFile{
		SizeTrust: TrustedSizeValue,
		File:      files_sdk.File{Path: "reports/sub/deep.txt", DisplayName: "deep.txt", Type: "file", Size: 10},
	}
	// The listing of reports also returns the nested file, as the API can do.
	server.MockFiles["reports/nested-copy"] = mockFile{
		SizeTrust: TrustedSizeValue,
		File:      files_sdk.File{Path: "reports/sub/deep.txt", DisplayName: "deep.txt", Type: "file", Size: 10},
	}

	job := client.Downloader(DownloaderParams{RemotePath: "reports/", LocalPath: root + string(os.PathSeparator)})
	job.Start()
	job.Wait()

	for _, s := range job.Statuses {
		require.NoError(t, s.Err(), s.RemotePath())
	}
	listed := 0
	for _, s := range job.Statuses {
		if s.RemotePath() == "reports/sub/deep.txt" {
			listed++
		}
	}
	assert.Equal(t, 1, listed, "the nested file is listed once, with its own folder")
	assert.FileExists(t, filepath.Join(root, "sub", "deep.txt"))
}

func TestServerPathComparisonPreservesComponentBoundaries(t *testing.T) {
	for _, tc := range []struct {
		folder  string
		path    string
		matches bool
	}{
		{"q\u0301/カ", "q/か/Original.txt", true},
		{"Straße", "STRASSE/Original.txt", true},
		{"folder", "folder /Original.txt", false},
		{"folder", "folder2/Original.txt", false},
	} {
		t.Run(tc.path, func(t *testing.T) {
			child, listingErr := listingEntryIsChild(tc.folder, files_sdk.File{Path: tc.path, DisplayName: "Original.txt", Type: "file"})
			name, destinationErr := localNameBelow(tc.folder, tc.path)
			if tc.matches {
				require.NoError(t, listingErr)
				require.True(t, child)
				require.NoError(t, destinationErr)
				require.Equal(t, "Original.txt", name)
			} else {
				require.Error(t, listingErr)
				require.Error(t, destinationErr)
			}
		})
	}
}
