package file

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/directory"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lib/direction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncAfterActionsValidateRejectsDeleteSourceFilesWithMoveSource(t *testing.T) {
	err := SyncAfterActions{DeleteSourceFiles: true, MoveSource: "archive"}.Validate()
	require.EqualError(t, err, "delete source files and move source cannot both be enabled")
}

func TestRegisterSyncAfterActionsDryRunSkipsActions(t *testing.T) {
	tmpDir := t.TempDir()
	sourcePath := filepath.Join(tmpDir, "source.txt")
	movePath := filepath.Join(tmpDir, "moved.txt")
	require.NoError(t, os.WriteFile(sourcePath, []byte("contents"), 0600))

	job, fileStatus := syncAfterActionUploadJob(sourcePath)
	job.Type = directory.File
	logs := make([]status.Log, 0)

	registerSyncAfterActions(
		job,
		SyncAfterActions{
			DeleteSourceFiles:        true,
			DeleteSourceEmptyFolders: true,
			MoveSource:               movePath,
			Log: func(log status.Log, err error) {
				logs = append(logs, log)
			},
		},
		true,
		files_sdk.Config{}.Init(),
	)

	// dryRun=true skips registration entirely; later status updates cannot run cleanup actions.
	job.UpdateStatus(status.Complete, fileStatus, nil)

	assert.FileExists(t, sourcePath)
	assert.NoFileExists(t, movePath)
	assert.Empty(t, logs)
}

func TestSyncAfterActionsRunWithCallerEventsReporter(t *testing.T) {
	sourcePath := filepath.Join(t.TempDir(), "source.txt")
	require.NoError(t, os.WriteFile(sourcePath, []byte("contents"), 0600))

	job, fileStatus := syncAfterActionUploadJob(sourcePath)
	logs := make([]status.Log, 0)
	callerReporterCalled := false

	registerSyncAfterActions(
		job,
		SyncAfterActions{
			DeleteSourceFiles: true,
			Log: func(log status.Log, err error) {
				require.NoError(t, err)
				logs = append(logs, log)
			},
		},
		false,
		files_sdk.Config{}.Init(),
	)
	job.SetEventsReporter(CreateFileEvents(func(file JobFile) {
		callerReporterCalled = true
	}, status.Complete))

	job.UpdateStatus(status.Complete, fileStatus, nil)

	assert.True(t, callerReporterCalled)
	assert.NoFileExists(t, sourcePath)
	require.Len(t, logs, 1)
	assert.Equal(t, "delete source", logs[0].Action)
}

func TestRegisterSyncAfterActionsDeleteSource(t *testing.T) {
	tmpDir := t.TempDir()
	sourcePath := filepath.Join(tmpDir, "source.txt")
	require.NoError(t, os.WriteFile(sourcePath, []byte("contents"), 0600))

	job, fileStatus := syncAfterActionUploadJob(sourcePath)
	logs := make([]status.Log, 0)

	registerSyncAfterActions(
		job,
		SyncAfterActions{
			DeleteSourceFiles: true,
			Log: func(log status.Log, err error) {
				require.NoError(t, err)
				logs = append(logs, log)
			},
		},
		false,
		files_sdk.Config{}.Init(),
	)

	job.UpdateStatus(status.Complete, fileStatus, nil)

	assert.NoFileExists(t, sourcePath)
	require.Len(t, logs, 1)
	assert.Equal(t, "delete source", logs[0].Action)
}

func TestRegisterSyncAfterActionsInvalidOptionsLogAndSkipActions(t *testing.T) {
	tmpDir := t.TempDir()
	sourcePath := filepath.Join(tmpDir, "source.txt")
	movePath := filepath.Join(tmpDir, "moved.txt")
	require.NoError(t, os.WriteFile(sourcePath, []byte("contents"), 0600))

	job, fileStatus := syncAfterActionUploadJob(sourcePath)
	var logErr error

	registerSyncAfterActions(
		job,
		SyncAfterActions{
			DeleteSourceFiles: true,
			MoveSource:        movePath,
			Log: func(log status.Log, err error) {
				assert.Equal(t, "sync after actions", log.Action)
				logErr = err
			},
		},
		false,
		files_sdk.Config{}.Init(),
	)

	job.UpdateStatus(status.Complete, fileStatus, nil)

	require.EqualError(t, logErr, "delete source files and move source cannot both be enabled")
	assert.FileExists(t, sourcePath)
	assert.NoFileExists(t, movePath)
}

func TestRegisterSyncAfterActionsMoveSource(t *testing.T) {
	tmpDir := t.TempDir()
	sourcePath := filepath.Join(tmpDir, "source.txt")
	movePath := filepath.Join(tmpDir, "archive", "moved.txt")
	require.NoError(t, os.WriteFile(sourcePath, []byte("contents"), 0600))

	job, fileStatus := syncAfterActionUploadJob(sourcePath)
	logs := make([]status.Log, 0)

	registerSyncAfterActions(
		job,
		SyncAfterActions{
			MoveSource: movePath,
			Log: func(log status.Log, err error) {
				require.NoError(t, err)
				logs = append(logs, log)
			},
		},
		false,
		files_sdk.Config{}.Init(),
	)

	job.UpdateStatus(status.Complete, fileStatus, nil)

	assert.NoFileExists(t, sourcePath)
	assert.FileExists(t, movePath)
	require.Len(t, logs, 1)
	assert.Equal(t, "move source", logs[0].Action)
	assert.Equal(t, movePath, logs[0].Path)
}

func TestRegisterSyncAfterActionsDownloadDeleteSource(t *testing.T) {
	var deletedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		deletedPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	localPath := filepath.Join(t.TempDir(), "source.txt")
	job, fileStatus := syncAfterActionDownloadJob(localPath, "source.txt")
	logs := make([]status.Log, 0)

	registerSyncAfterActions(
		job,
		SyncAfterActions{
			DeleteSourceFiles: true,
			Log: func(log status.Log, err error) {
				require.NoError(t, err)
				logs = append(logs, log)
			},
		},
		false,
		files_sdk.Config{EndpointOverride: server.URL}.Init(),
	)

	job.UpdateStatus(status.Complete, fileStatus, nil)

	assert.Equal(t, "/api/rest/v1/files/source.txt", deletedPath)
	require.Len(t, logs, 1)
	assert.Equal(t, "delete source", logs[0].Action)
	assert.Equal(t, "source.txt", logs[0].Path)
}

func TestRegisterSyncAfterActionsDownloadDeleteSourceEmptyFoldersUsesNonRecursiveDelete(t *testing.T) {
	var deletedPath string
	var recursiveParam string
	listRequests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/rest/v1/folders/source-folder":
			listRequests++
			w.Write([]byte("[]"))
		case r.Method == http.MethodDelete && r.URL.Path == "/api/rest/v1/files/source-folder":
			deletedPath = r.URL.Path
			recursiveParam = r.URL.Query().Get("recursive")
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	}))
	defer server.Close()

	job, _ := syncAfterActionDownloadJob(filepath.Join(t.TempDir(), "source"), "source-folder")
	logCh := make(chan status.Log, 1)

	actions := SyncAfterActions{
		DeleteSourceEmptyFolders: true,
		Log: func(log status.Log, err error) {
			require.NoError(t, err)
			logCh <- log
		},
	}
	config := files_sdk.Config{EndpointOverride: server.URL}.Init()
	runSyncAfterEmptyFolders(job, actions, false, config)

	log := requireSyncActionLog(t, logCh)
	assert.Equal(t, 2, listRequests)
	assert.Equal(t, "/api/rest/v1/files/source-folder", deletedPath)
	assert.Empty(t, recursiveParam)
	assert.Equal(t, "delete source folder", log.Action)
	assert.Equal(t, "source-folder", log.Path)
}

func TestRegisterSyncAfterActionsDownloadDeleteSourceEmptyFoldersWalksDepthFirst(t *testing.T) {
	requests := make([]string, 0)
	rootListRequests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		requests = append(requests, r.Method+" "+r.URL.Path)
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/rest/v1/folders/source-folder":
			rootListRequests++
			if rootListRequests == 1 {
				w.Write([]byte(`[{"path":"source-folder/empty-child","display_name":"empty-child","type":"directory"}]`))
				return
			}
			w.Write([]byte("[]"))
		case r.Method == http.MethodGet && r.URL.Path == "/api/rest/v1/folders/source-folder/empty-child":
			w.Write([]byte("[]"))
		case r.Method == http.MethodDelete && r.URL.Path == "/api/rest/v1/files/source-folder/empty-child":
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodDelete && r.URL.Path == "/api/rest/v1/files/source-folder":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	}))
	defer server.Close()

	job, _ := syncAfterActionDownloadJob(filepath.Join(t.TempDir(), "source"), "source-folder")
	logCh := make(chan status.Log, 1)

	actions := SyncAfterActions{
		DeleteSourceEmptyFolders: true,
		Log: func(log status.Log, err error) {
			require.NoError(t, err)
			logCh <- log
		},
	}
	config := files_sdk.Config{EndpointOverride: server.URL}.Init()
	runSyncAfterEmptyFolders(job, actions, false, config)

	requireSyncActionLog(t, logCh)
	assert.Equal(t, []string{
		"GET /api/rest/v1/folders/source-folder",
		"GET /api/rest/v1/folders/source-folder/empty-child",
		"GET /api/rest/v1/folders/source-folder/empty-child",
		"DELETE /api/rest/v1/files/source-folder/empty-child",
		"GET /api/rest/v1/folders/source-folder",
		"DELETE /api/rest/v1/files/source-folder",
	}, requests)
}

func TestRegisterSyncAfterActionsDownloadDeleteSourceEmptyFoldersCleansChildFoldersWhenParentHasFile(t *testing.T) {
	deletedPaths := make([]string, 0)
	nestedListRequests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/rest/v1/folders/source-folder":
			w.Write([]byte(`[
				{"path":"source-folder/keep.txt","display_name":"keep.txt","type":"file"},
				{"path":"source-folder/empty-child","display_name":"empty-child","type":"directory"}
			]`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/rest/v1/folders/source-folder/empty-child":
			nestedListRequests++
			w.Write([]byte("[]"))
		case r.Method == http.MethodDelete:
			deletedPaths = append(deletedPaths, r.URL.Path)
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	}))
	defer server.Close()

	job, _ := syncAfterActionDownloadJob(filepath.Join(t.TempDir(), "source"), "source-folder")
	logCh := make(chan status.Log, 1)

	actions := SyncAfterActions{
		DeleteSourceEmptyFolders: true,
		Log: func(log status.Log, err error) {
			require.NoError(t, err)
			logCh <- log
		},
	}
	config := files_sdk.Config{EndpointOverride: server.URL}.Init()
	runSyncAfterEmptyFolders(job, actions, false, config)

	requireSyncActionLog(t, logCh)
	assert.Equal(t, 2, nestedListRequests)
	assert.Equal(t, []string{"/api/rest/v1/files/source-folder/empty-child"}, deletedPaths)
}

func TestRegisterSyncAfterActionsDeleteSourceEmptyFoldersOnFinish(t *testing.T) {
	sourceDir := filepath.Join(t.TempDir(), "source")
	sourcePath := filepath.Join(sourceDir, "done.txt")
	require.NoError(t, os.MkdirAll(sourceDir, 0750))

	job, _ := syncAfterActionUploadJob(sourcePath)
	logCh := make(chan status.Log, 1)

	actions := SyncAfterActions{
		DeleteSourceEmptyFolders: true,
		Log: func(log status.Log, err error) {
			require.NoError(t, err)
			logCh <- log
		},
	}
	runSyncAfterEmptyFolders(job, actions, false, files_sdk.Config{}.Init())

	log := requireSyncActionLog(t, logCh)
	assert.Equal(t, "delete source folder", log.Action)
	assert.Equal(t, sourceDir, log.Path)
	assert.NoDirExists(t, sourceDir)
}

func TestSyncAfterEmptyFoldersRunBeforeJobFinish(t *testing.T) {
	sourceDir := filepath.Join(t.TempDir(), "source")
	require.NoError(t, os.MkdirAll(sourceDir, 0750))

	job, fileStatus := syncAfterActionUploadJob(sourceDir)
	job.Type = directory.Dir
	logCh := make(chan status.Log, 1)
	var logErr error
	finishedWhenLogged := true
	actions := SyncAfterActions{
		DeleteSourceEmptyFolders: true,
		Log: func(log status.Log, err error) {
			logErr = err
			finishedWhenLogged = job.Finished.Called()
			logCh <- log
		},
	}

	onComplete := make(chan *UploadStatus, 1)
	WaitTellFinished(job, onComplete, func() {
		runSyncAfterEmptyFolders(job, actions, false, files_sdk.Config{}.Init())
	})

	job.Start()
	job.EndScan()
	onComplete <- fileStatus
	job.Wait()

	log := requireSyncActionLog(t, logCh)
	require.NoError(t, logErr)
	assert.Equal(t, "delete source folder", log.Action)
	assert.False(t, finishedWhenLogged)
	assert.True(t, job.Finished.Called())
	assert.NoDirExists(t, sourceDir)
}

func TestRegisterSyncAfterActionsDeleteSourceEmptyFoldersForFolderSourceStaysInSourceRoot(t *testing.T) {
	parentDir := t.TempDir()
	sourceDir := filepath.Join(parentDir, "source")
	siblingDir := filepath.Join(parentDir, "sibling")
	require.NoError(t, os.MkdirAll(sourceDir, 0750))
	require.NoError(t, os.MkdirAll(siblingDir, 0750))

	job, _ := syncAfterActionUploadJob(sourceDir)
	job.Type = directory.Dir
	logCh := make(chan status.Log, 1)

	actions := SyncAfterActions{
		DeleteSourceEmptyFolders: true,
		Log: func(log status.Log, err error) {
			require.NoError(t, err)
			logCh <- log
		},
	}
	runSyncAfterEmptyFolders(job, actions, false, files_sdk.Config{}.Init())

	log := requireSyncActionLog(t, logCh)
	assert.Equal(t, "delete source folder", log.Action)
	assert.Equal(t, sourceDir, log.Path)
	assert.NoDirExists(t, sourceDir)
	assert.DirExists(t, siblingDir)
}

func TestRegisterSyncAfterActionsSkipsDeleteSourceEmptyFoldersWhenJobErrored(t *testing.T) {
	sourceDir := filepath.Join(t.TempDir(), "source")
	sourcePath := filepath.Join(sourceDir, "done.txt")
	require.NoError(t, os.MkdirAll(sourceDir, 0750))

	job, fileStatus := syncAfterActionUploadJob(sourcePath)
	logCh := make(chan status.Log, 1)

	actions := SyncAfterActions{
		DeleteSourceEmptyFolders: true,
		Log: func(log status.Log, err error) {
			logCh <- log
		},
	}

	job.UpdateStatus(status.Errored, fileStatus, errors.New("failed transfer"))
	runSyncAfterEmptyFolders(job, actions, false, files_sdk.Config{}.Init())

	assertNoSyncActionLog(t, logCh)
	assert.DirExists(t, sourceDir)
}

func TestRegisterSyncAfterActionsCancelSkipsDeleteSourceEmptyFolders(t *testing.T) {
	sourceDir := filepath.Join(t.TempDir(), "source")
	sourcePath := filepath.Join(sourceDir, "done.txt")
	require.NoError(t, os.MkdirAll(sourceDir, 0750))

	job, _ := syncAfterActionUploadJob(sourcePath)
	job.WithContext(context.Background())
	logCh := make(chan status.Log, 1)

	actions := SyncAfterActions{
		DeleteSourceEmptyFolders: true,
		Log: func(log status.Log, err error) {
			logCh <- log
		},
	}

	job.Cancel()
	runSyncAfterEmptyFolders(job, actions, false, files_sdk.Config{}.Init())

	assertNoSyncActionLog(t, logCh)
	assert.DirExists(t, sourceDir)
}

func TestRegisterSyncAfterActionsCancelSkipsPerFileActions(t *testing.T) {
	sourcePath := filepath.Join(t.TempDir(), "source.txt")
	require.NoError(t, os.WriteFile(sourcePath, []byte("contents"), 0600))

	job, fileStatus := syncAfterActionUploadJob(sourcePath)
	job.WithContext(context.Background())
	logCh := make(chan status.Log, 1)

	registerSyncAfterActions(
		job,
		SyncAfterActions{
			DeleteSourceFiles: true,
			Log: func(log status.Log, err error) {
				logCh <- log
			},
		},
		false,
		files_sdk.Config{}.Init(),
	)

	job.Cancel()
	job.UpdateStatus(status.Complete, fileStatus, nil)

	assertNoSyncActionLog(t, logCh)
	assert.FileExists(t, sourcePath)
}

func syncAfterActionUploadJob(localPath string) (*Job, *UploadStatus) {
	job := (&Job{
		Config:     files_sdk.Config{}.Init(),
		Direction:  direction.UploadType,
		Logger:     lib.NullLogger{},
		LocalPath:  localPath,
		RemotePath: "source.txt",
		Type:       directory.File,
	}).Init()

	fileStatus := &UploadStatus{
		file:       files_sdk.File{Path: "source.txt"},
		status:     status.Queued,
		job:        job,
		localPath:  localPath,
		remotePath: "source.txt",
		Mutex:      &sync.RWMutex{},
	}
	job.Add(fileStatus)
	return job, fileStatus
}

func syncAfterActionDownloadJob(localPath string, remotePath string) (*Job, *DownloadStatus) {
	job := (&Job{
		Config:     files_sdk.Config{}.Init(),
		Direction:  direction.DownloadType,
		Logger:     lib.NullLogger{},
		LocalPath:  localPath,
		RemotePath: remotePath,
		Type:       directory.File,
	}).Init()

	fileStatus := &DownloadStatus{
		file:       files_sdk.File{Path: remotePath},
		status:     status.Queued,
		job:        job,
		localPath:  localPath,
		remotePath: remotePath,
		Mutex:      &sync.RWMutex{},
	}
	job.Add(fileStatus)
	return job, fileStatus
}

func requireSyncActionLog(t *testing.T, ch <-chan status.Log) status.Log {
	t.Helper()
	select {
	case log := <-ch:
		return log
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for sync action log")
		return status.Log{}
	}
}

func assertNoSyncActionLog(t *testing.T, ch <-chan status.Log) {
	t.Helper()
	select {
	case log := <-ch:
		t.Fatalf("unexpected sync action log: %+v", log)
	case <-time.After(100 * time.Millisecond):
	}
}
