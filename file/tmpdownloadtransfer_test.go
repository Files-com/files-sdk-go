package file

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockFolder registers a remote folder holding the given files.
func mockFolder(server *MockAPIServer, folder string, files map[string][]byte) {
	server.MockFiles[folder] = mockFile{File: files_sdk.File{Path: folder, DisplayName: filepath.Base(folder), Type: "directory"}}
	for name, data := range files {
		path := folder + "/" + name
		server.MockFiles[path] = mockFile{
			SizeTrust: TrustedSizeValue,
			Data:      data,
			File:      files_sdk.File{Path: path, DisplayName: filepath.Base(name), Type: "file", Size: int64(len(data))},
		}
	}
}

// downloadOne downloads one remote file under its own name into root, the way
// the CLI configures a sync: no default rules, so an ordinary name ending in
// ".download" is transferred.
func downloadOne(t *testing.T, server *MockAPIServer, root string, name string) {
	t.Helper()
	job := server.Client().Downloader(DownloaderParams{
		RemotePath: "reports/" + name,
		LocalPath:  filepath.Join(root, name),
		Ignore:     []string{},
	})
	job.Start()
	job.Wait()
	require.NoError(t, job.Statuses[0].Err())
}

// The reported defect, end to end: a remote file named after another file's
// unfinished download must not be delivered as that file, and must still be
// delivered correctly as itself. The two files are the same size, so a wrongly
// adopted temporary file cannot be rejected for being longer than the file
// being downloaded — that is the case that used to finish a download from
// someone else's bytes without asking the server for any content.
func TestClient_Downloader_deliversASiblingNamedLikeATemporaryDownloadItsOwnBytes(t *testing.T) {
	const name = "quarterly.csv"
	report := bytes.Repeat([]byte("report bytes;"), 5000)
	sibling := bytes.Repeat([]byte("SIBLING BYTE;"), 5000)
	require.Len(t, sibling, len(report))

	assertEachGotItsOwn := func(t *testing.T, root string) {
		t.Helper()
		delivered, err := os.ReadFile(filepath.Join(root, name))
		require.NoError(t, err)
		assert.Equal(t, report, delivered, "%v received another file's bytes", name)
		delivered, err = os.ReadFile(filepath.Join(root, name+".download"))
		require.NoError(t, err)
		assert.Equal(t, sibling, delivered, "%v received another file's bytes", name+".download")
	}

	// Each order is two transfers one after the other, so the first one's
	// result is already on disk when the second one runs.
	for _, order := range [][]string{
		{name + ".download", name},
		{name, name + ".download"},
	} {
		t.Run(order[0]+" first", func(t *testing.T) {
			root := t.TempDir()
			server := (&MockAPIServer{T: t}).Do()
			defer server.Shutdown()
			mockFolder(server, "reports", map[string][]byte{name: report, name + ".download": sibling})

			for _, transfer := range order {
				downloadOne(t, server, root, transfer)
			}

			assertEachGotItsOwn(t, root)
		})
	}

	t.Run("both in one folder download", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		mockFolder(server, "reports", map[string][]byte{name: report, name + ".download": sibling})

		job := server.Client().Downloader(DownloaderParams{
			RemotePath: "reports",
			LocalPath:  root + string(os.PathSeparator),
			Ignore:     []string{},
		})
		job.Start()
		job.Wait()

		assert.Equal(t, status.Complete, findStatus(t, job, "reports/"+name).Status())
		assert.Equal(t, status.Complete, findStatus(t, job, "reports/"+name+".download").Status())
		assertEachGotItsOwn(t, filepath.Join(root, "reports"))
	})
}

// Nothing is ever written to a temporary download path, whatever ignore and
// include rules the caller passes, and names that only look similar still
// transfer.
func TestClient_Downloader_neverDownloadsIntoATemporaryDownloadPath(t *testing.T) {
	root := t.TempDir()
	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	report := bytes.Repeat([]byte("report bytes;"), 100)
	mockFolder(server, "reports", map[string][]byte{
		"report.csv":                           report,
		".~files-cli.report.csv.download":      []byte("planted"),
		".~files-cli-notes.download":           []byte("an ordinary file"),
		".~files-clients.download":             []byte("another ordinary file"),
		".~FILES-CLI.other.csv.download/x.csv": []byte("planted in a folder"),
		".~fileſ-cli.longs.csv.download":       []byte("planted under an alias"),
	})
	server.MockFiles["reports/.~FILES-CLI.other.csv.download"] = mockFile{File: files_sdk.File{
		Path: "reports/.~FILES-CLI.other.csv.download", DisplayName: ".~FILES-CLI.other.csv.download", Type: "directory",
	}}

	job := server.Client().Downloader(DownloaderParams{
		RemotePath: "reports",
		LocalPath:  root + string(os.PathSeparator),
		// A caller asking for exactly these names still does not get them.
		Ignore:  []string{"!.~files-cli*"},
		Include: []string{"*"},
	})
	job.Start()
	job.Wait()

	assert.Equal(t, status.Ignored, findStatus(t, job, "reports/.~files-cli.report.csv.download").Status())
	assert.Equal(t, status.Ignored, findStatus(t, job, "reports/.~FILES-CLI.other.csv.download/x.csv").Status())
	assert.Equal(t, status.Ignored, findStatus(t, job, "reports/.~fileſ-cli.longs.csv.download").Status())
	assert.NoFileExists(t, filepath.Join(root, "reports", ".~files-cli.report.csv.download"))
	assert.NoFileExists(t, filepath.Join(root, "reports", ".~FILES-CLI.other.csv.download", "x.csv"))
	assert.NoFileExists(t, filepath.Join(root, "reports", ".~fileſ-cli.longs.csv.download"))

	assert.Equal(t, status.Complete, findStatus(t, job, "reports/report.csv").Status())
	assert.Equal(t, status.Complete, findStatus(t, job, "reports/.~files-cli-notes.download").Status())
	assert.Equal(t, status.Complete, findStatus(t, job, "reports/.~files-clients.download").Status())
	delivered, err := os.ReadFile(filepath.Join(root, "reports", "report.csv"))
	require.NoError(t, err)
	assert.Equal(t, report, delivered)
}

// A temporary download is never uploaded either, so it cannot come back as a
// remote file able to take over the download it belongs to.
func TestClient_Uploader_neverUploadsATemporaryDownload(t *testing.T) {
	root := t.TempDir()
	staged := filepath.Join(root, ".~FILES-CLI.other.csv.download")
	require.NoError(t, os.Mkdir(staged, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(staged, "other.csv"), []byte("half of other.csv"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".~files-cli.report.csv.download"), []byte("half of report.csv"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "report.csv"), []byte("all of report.csv"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".~fileſ-cli.longs.csv.download"), []byte("half of longs.csv"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".~files-cli-notes.download"), []byte("an ordinary file"), 0600))

	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	job := server.Client().Uploader(UploaderParams{
		LocalPath:  root + string(os.PathSeparator),
		RemotePath: "reports",
		// One file at a time: the mock server's file map is not safe for
		// concurrent uploads, which is unrelated to what this checks.
		Manager: manager.Build(1, 1),
		// A caller asking for exactly these names still does not send them.
		Ignore:  []string{"!.~files-cli*"},
		Include: []string{"*"},
	})
	job.Start()
	job.Wait()

	assert.Equal(t, status.Ignored, findStatus(t, job, "reports/.~files-cli.report.csv.download").Status())
	assert.Equal(t, status.Ignored, findStatus(t, job, "reports/.~FILES-CLI.other.csv.download/other.csv").Status())
	assert.Equal(t, status.Ignored, findStatus(t, job, "reports/.~fileſ-cli.longs.csv.download").Status())
	assert.Equal(t, status.Complete, findStatus(t, job, "reports/report.csv").Status())
	assert.Equal(t, status.Complete, findStatus(t, job, "reports/.~files-cli-notes.download").Status())
}

// A transfer maps each file to a path on the other side, and the mapping can
// drop the temporary name from one side: a temporary download folder chosen as
// the folder to transfer keeps its name only at the source, and an explicitly
// chosen output keeps it only at the destination. Each case below is stopped by
// one side of the check alone.
func TestTransferGateCoversBothSidesOfThePathMapping(t *testing.T) {
	const staged = ".~files-cli.victim.csv.download"

	t.Run("a temporary download file downloaded to an ordinary output", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		mockFolder(server, "reports", map[string][]byte{staged: []byte("planted")})

		job := server.Client().Downloader(DownloaderParams{
			RemotePath: "reports/" + staged,
			LocalPath:  filepath.Join(root, "ordinary.csv"),
			Ignore:     []string{},
		})
		job.Start()
		job.Wait()

		// The destination keeps none of the temporary name, so only the source
		// side of the check can stop this transfer.
		transfer := findStatus(t, job, "reports/"+staged)
		require.Equal(t, filepath.Join(root, "ordinary.csv"), transfer.LocalPath())
		assert.Equal(t, status.Ignored, transfer.Status())
		assert.NoFileExists(t, filepath.Join(root, "ordinary.csv"))
	})

	t.Run("an ordinary remote file downloaded to a temporary download path", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		mockFolder(server, "reports", map[string][]byte{"ordinary.csv": []byte("ordinary")})

		job := server.Client().Downloader(DownloaderParams{
			RemotePath: "reports/ordinary.csv",
			LocalPath:  filepath.Join(root, staged),
			Ignore:     []string{},
		})
		job.Start()
		job.Wait()

		// The source keeps none of the temporary name.
		assert.Equal(t, status.Ignored, findStatus(t, job, "reports/ordinary.csv").Status())
		assert.NoFileExists(t, filepath.Join(root, staged))
	})

	t.Run("a temporary download folder chosen as the folder to upload", func(t *testing.T) {
		root := t.TempDir()
		folder := filepath.Join(root, staged)
		require.NoError(t, os.Mkdir(folder, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(folder, "victim.csv"), []byte("half of victim.csv"), 0600))
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()

		job := server.Client().Uploader(UploaderParams{
			LocalPath:  folder + string(os.PathSeparator),
			RemotePath: "reports",
			Manager:    manager.Build(1, 1),
			Ignore:     []string{},
		})
		job.Start()
		job.Wait()

		// The destination keeps none of the temporary name.
		assert.Equal(t, status.Ignored, findStatus(t, job, "reports/victim.csv").Status())
	})

	t.Run("an ordinary local file uploaded into a temporary download folder", func(t *testing.T) {
		root := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(root, "ordinary.csv"), []byte("ordinary"), 0600))
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		mockFolder(server, staged, nil)

		job := server.Client().Uploader(UploaderParams{
			LocalPath:  filepath.Join(root, "ordinary.csv"),
			RemotePath: staged,
			Manager:    manager.Build(1, 1),
			Ignore:     []string{},
		})
		job.Start()
		job.Wait()

		// The source keeps none of the temporary name.
		assert.Equal(t, status.Ignored, findStatus(t, job, staged+"/ordinary.csv").Status())
	})
}

// A file already sitting where a client from before this naming would have
// staged the download must not be taken for the unfinished download of the file
// it is named after. It is the same size as the remote file, which is what used
// to finish the download from it without asking the server for any content.
func TestClient_Downloader_doesNotFinalizeAPlantedFileAsTheRequestedFile(t *testing.T) {
	root := t.TempDir()
	destination := filepath.Join(root, "reports")
	require.NoError(t, os.MkdirAll(destination, 0755))
	planted := []byte("planted bytes")
	genuine := []byte("genuine bytes")
	require.Len(t, genuine, len(planted))

	// The two shapes the previous naming took: a file on Linux and Windows, a
	// folder holding the file on macOS.
	flat := filepath.Join(destination, "flat.csv.download")
	require.NoError(t, os.WriteFile(flat, planted, 0600))
	folder := filepath.Join(destination, "folder.csv.download")
	require.NoError(t, os.Mkdir(folder, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(folder, "folder.csv"), planted, 0600))

	server := (&MockAPIServer{T: t}).Do()
	defer server.Shutdown()
	mockFolder(server, "reports", map[string][]byte{"flat.csv": genuine, "folder.csv": genuine})

	job := server.Client().Downloader(DownloaderParams{
		RemotePath: "reports",
		LocalPath:  root + string(os.PathSeparator),
		Ignore:     []string{},
	})
	job.Start()
	job.Wait()

	for _, name := range []string{"flat.csv", "folder.csv"} {
		assert.Equal(t, status.Complete, findStatus(t, job, "reports/"+name).Status())
		delivered, err := os.ReadFile(filepath.Join(destination, name))
		require.NoError(t, err)
		assert.Equalf(t, genuine, delivered, "%v was delivered a planted file's bytes", name)
	}

	// The leftovers of the previous naming are left alone rather than consumed.
	leftover, err := os.ReadFile(flat)
	require.NoError(t, err)
	assert.Equal(t, planted, leftover)
}
