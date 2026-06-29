package file

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/ignore"
	"github.com/Files-com/files-sdk-go/v3/lib"
	gitignore "github.com/sabhiram/go-gitignore"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ReporterCall struct {
	JobFile
	err error
}

type TestSetup struct {
	files         []Entity
	reporterCalls []ReporterCall
	fstest.MapFS
	DownloaderParams
	rootDestination string
	tempDir         string
	files_sdk.Config
}

func NewTestSetup() *TestSetup {
	t := &TestSetup{Config: files_sdk.Config{}.Init()}
	t.MapFS = make(fstest.MapFS)
	err := t.TempDir()
	if err != nil {
		panic(err)
	}
	return t
}

func (setup *TestSetup) Reporter() EventsReporter {
	m := sync.Mutex{}

	callback := func(status JobFile) {
		m.Lock()
		setup.reporterCalls = append(setup.reporterCalls, ReporterCall{JobFile: status})
		m.Unlock()
	}

	return CreateFileEvents(callback, append(status.Excluded, status.Included...)...)
}

func (setup *TestSetup) TempDir() error {
	var err error
	setup.tempDir, err = os.MkdirTemp("", "test")

	return err
}

func (setup *TestSetup) TearDown() error {
	return os.RemoveAll(setup.tempDir)
}

func (setup *TestSetup) Call() *Job {
	setup.DownloaderParams.config = setup.Config
	job := downloader(
		context.Background(),
		setup.MapFS,
		setup.DownloaderParams,
	)

	job.Start()
	job.Wait()
	return job
}

func (setup *TestSetup) RootDestination() string {
	if setup.rootDestination != "" && setup.rootDestination[len(setup.rootDestination)-1:] == string(os.PathSeparator) {
		return filepath.Join(setup.tempDir, setup.rootDestination) + string(os.PathSeparator)
	}

	return filepath.Join(setup.tempDir, setup.rootDestination)
}

func Test_downloadFolder_ending_in_slash(t *testing.T) {
	setup := NewTestSetup()
	setup.MapFS["some-path"] = &fstest.MapFile{
		Data:    nil,
		Mode:    fs.ModeDir,
		ModTime: time.Time{},
		Sys:     files_sdk.File{DisplayName: "some-path", Path: "some-path", Type: "directory"},
	}

	setup.MapFS["some-path/taco.png"] = &fstest.MapFile{
		Data:    make([]byte, 100),
		Mode:    fs.ModePerm,
		ModTime: time.Time{},
		Sys:     files_sdk.File{DisplayName: "taco.png", Path: "some-path/taco.png", Type: "file", Size: 100},
	}

	setup.DownloaderParams = DownloaderParams{RemotePath: "some-path", EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
	setup.rootDestination = "some-path/"
	setup.Call()

	assert.Equal(t, 2, setup.reporterCalls[0].Job.Count())
	assert.Equal(t, 5, len(setup.reporterCalls))

	expectedStatus := []status.Status{status.Queued, status.Queued, status.FolderCreated, status.Downloading, status.Complete}
	var actualStatus []status.Status
	for _, call := range setup.reporterCalls {
		actualStatus = append(actualStatus, call.Status)
		assert.NoError(t, call.err)
		switch call.Status {
		case status.FolderCreated:
			assert.Equal(t, "some-path", call.File.Path)
		case status.Complete:
			assert.Equal(t, "some-path/taco.png", call.File.Path)
		}
	}
	assert.ElementsMatch(t, expectedStatus, actualStatus)

	assert.Equal(t, int64(0), setup.reporterCalls[0].TransferBytes)
	assert.Equal(t, true, setup.reporterCalls[0].Job.All(status.Ended...))
	assert.Equal(t, int64(100), setup.reporterCalls[0].Job.TransferBytes())
	assert.Equal(t, int64(100), setup.reporterCalls[0].Job.TotalBytes())

	assert.NoError(t, setup.TearDown())
}

func TestClient_Downloader(t *testing.T) {
	t.Run("small file with size", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["small-file-with-size.txt"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 1999},
		}
		job := client.Downloader(DownloaderParams{RemotePath: "small-file-with-size.txt", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		f, err := os.Open(filepath.Join(root, "small-file-with-size.txt"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, int64(1999), stat.Size())
	})

	t.Run("large file with size", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["large-file-with-size.txt"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 19999999},
		}
		job := client.Downloader(DownloaderParams{RemotePath: "large-file-with-size.txt", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		f, err := os.Open(filepath.Join(root, "large-file-with-size.txt"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, int64(19999999), stat.Size())
	})

	t.Run("large file with size with max concurrent connections of 1", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["large-file-with-size.txt"] = mockFile{
			SizeTrust:      TrustedSizeValue,
			File:           files_sdk.File{Size: 1024 * 1024 * 100},
			MaxConnections: 1,
		}
		m := manager.Build(1, 1)
		job := client.Downloader(DownloaderParams{RemotePath: "large-file-with-size.txt", LocalPath: root + "/", Manager: m})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		f, err := os.Open(filepath.Join(root, "large-file-with-size.txt"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, int64(1024*1024*100), stat.Size())
		assert.Len(t, server.TrackRequest["/download/:download_id"], 1)
	})

	t.Run("large file with size with max concurrent connections of 1", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["large-file-with-size.txt"] = mockFile{
			SizeTrust:      TrustedSizeValue,
			File:           files_sdk.File{Size: 1024 * 1024 * 50},
			MaxConnections: 1,
		}
		m := manager.Build(1, 1)
		job := client.Downloader(DownloaderParams{RemotePath: "large-file-with-size.txt", LocalPath: root + "/", Manager: m})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		f, err := os.Open(filepath.Join(root, "large-file-with-size.txt"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, int64(1024*1024*50), stat.Size())
		assert.Len(t, server.TrackRequest["/download/:download_id"], 1)
	})

	t.Run("large file with size DownloadFilesAsSingleStream", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["large-file-with-size.txt"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 1024 * 1024 * 50},
		}
		m := manager.Build(10, 1, true)
		job := client.Downloader(DownloaderParams{RemotePath: "large-file-with-size.txt", LocalPath: root + "/", Manager: m})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		f, err := os.Open(filepath.Join(root, "large-file-with-size.txt"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, int64(1024*1024*50), stat.Size())
		assert.Len(t, server.TrackRequest["/download/:download_id"], 1)
	})

	t.Run("large file with no size", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["large-file-with-no-size.txt"] = mockFile{
			SizeTrust: UntrustedSizeValue,
			File:      files_sdk.File{Size: 19999999},
		}

		job := client.Downloader(DownloaderParams{RemotePath: "large-file-with-no-size.txt", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		f, err := os.Open(filepath.Join(root, "large-file-with-no-size.txt"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, int64(19999999), stat.Size())
	})

	t.Run("large file with no size - extra parts are canceled", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		realSize := int64((1024 * 1024 * 5) - 256)
		server.MockFiles["large-file-with-no-size.txt"] = mockFile{
			SizeTrust: UntrustedSizeValue,
			File:      files_sdk.File{Size: 1024 * 1024 * 100},
			RealSize:  &realSize,
		}

		job := client.Downloader(DownloaderParams{RemotePath: "large-file-with-no-size.txt", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		f, err := os.Open(filepath.Join(root, "large-file-with-no-size.txt"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, realSize, stat.Size())
	})

	t.Run("large file with no size - client does not receive all bytes server reported to send", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		serverBytesSent := int64((1024 * 1024 * 5) + 256)
		server.MockFiles["large-file-with-no-size.txt"] = mockFile{
			SizeTrust:       UntrustedSizeValue,
			File:            files_sdk.File{Size: 1024 * 1024 * 15},
			ServerBytesSent: &serverBytesSent,
		}

		job := client.Downloader(DownloaderParams{RemotePath: "large-file-with-no-size.txt", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.EqualError(t, job.Statuses[0].Err(), `received size did not match server send size
expected 5243136 bytes sent 5242880 received`)
		_, err := os.Open(filepath.Join(root, "large-file-with-no-size.txt"))
		require.Error(t, err)
	})

	t.Run("large file with no size - client received more bytes than server reported to send", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		serverBytesSent := int64(1024 * 1024 * 4)
		server.MockFiles["large-file-with-no-size.txt"] = mockFile{
			SizeTrust:       UntrustedSizeValue,
			File:            files_sdk.File{Size: 1024 * 1024 * 15},
			ServerBytesSent: &serverBytesSent,
		}

		job := client.Downloader(DownloaderParams{RemotePath: "large-file-with-no-size.txt", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.EqualError(t, job.Statuses[0].Err(), `received size did not match server send size
expected 4194304 bytes sent 5242880 received`)
		_, err := os.Open(filepath.Join(root, "large-file-with-no-size.txt"))
		require.Error(t, err)
	})

	t.Run("large file with no size - when sever has invalid request status", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		serverBytesSent := int64(1024 * 1024 * 4)
		server.MockFiles["large-file-with-no-size.txt"] = mockFile{
			SizeTrust:          UntrustedSizeValue,
			File:               files_sdk.File{Size: 1024 * 1024 * 15},
			ServerBytesSent:    &serverBytesSent,
			ForceRequestStatus: "started",
		}

		job := client.Downloader(DownloaderParams{RemotePath: "large-file-with-no-size.txt", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		f, err := os.Open(filepath.Join(root, "large-file-with-no-size.txt"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, int64(1024*1024*15), stat.Size())
	})

	t.Run("large file with no size - when sever has failed request status", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["large-file-with-no-size.txt"] = mockFile{
			SizeTrust:           UntrustedSizeValue,
			File:                files_sdk.File{Size: 1024 * 1024 * 15},
			ForceRequestStatus:  "failed",
			ForceRequestMessage: "problem",
		}
		var events []JobFile
		eventReporter := CreateFileEvents(
			func(file JobFile) {
				events = append(events, file)
			},
			status.Included...,
		)

		job := client.Downloader(DownloaderParams{RemotePath: "large-file-with-no-size.txt", LocalPath: root + "/", EventsReporter: eventReporter})
		transferBytes := []string{"zero"}
		wait := make(chan bool)
		go func() {
			for {
				select {
				case <-job.Finished.C:
					wait <- true
					return
				default:
					bytes := job.TransferBytes()
					if bytes > 0 && transferBytes[len(transferBytes)-1] == "zero" {
						transferBytes = append(transferBytes, "bytes")
					}
					if bytes == 0 && transferBytes[len(transferBytes)-1] != "zero" {
						transferBytes = append(transferBytes, "zero")
					}
				}
			}
		}()

		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		assert.Error(t, job.Statuses[0].Err(), `received size did not match server send size
expected 4194304 bytes sent 5242880 received`)
		assert.Equal(t, []int64{0, 32768, 0}, lo.Map[JobFile, int64](events, func(item JobFile, index int) int64 { return item.TransferBytes }))
		assert.Equal(t, []string{"queued", "downloading", "errored"}, lo.Map[JobFile, string](events, func(item JobFile, index int) string { return item.StatusName }))
		<-wait
		assert.GreaterOrEqual(t, lo.Count[string](transferBytes, "zero"), 2, "After error transfer bytes are set to zero")
		assert.GreaterOrEqual(t, lo.Count[string](transferBytes, "bytes"), 2, "After error transfer bytes are set to zero")
	})

	t.Run("large file with bad size info real size is bigger", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		realSize := int64(20000000)
		server.MockFiles["file-with-mismatch-size-bigger"] = mockFile{
			SizeTrust: UntrustedSizeValue,
			File:      files_sdk.File{Size: 19999999},
			RealSize:  &realSize,
		}

		job := client.Downloader(DownloaderParams{RemotePath: "file-with-mismatch-size-bigger", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		require.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		f, err := os.Open(filepath.Join(root, "file-with-mismatch-size-bigger"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, int64(20000000), stat.Size())
	})

	t.Run("large file with bad size info real size is smaller", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		realSize := int64(19999999)
		server.MockFiles["file-with-mismatch-size-smaller"] = mockFile{
			SizeTrust: UntrustedSizeValue,
			File:      files_sdk.File{Size: 20000000},
			RealSize:  &realSize,
		}

		job := client.Downloader(DownloaderParams{RemotePath: "file-with-mismatch-size-smaller", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		require.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		f, err := os.Open(filepath.Join(root, "file-with-mismatch-size-smaller"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, int64(19999999), stat.Size())
	})

	multipleFiles := func(relativeRoot string, t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles[filepath.Join(relativeRoot, "file1")] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 6},
		}
		server.MockFiles[filepath.Join(relativeRoot, "file2")] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 1024 * 1024},
		}
		server.MockFiles[filepath.Join(relativeRoot, "file3")] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 1024 * 1024 * 2},
		}
		server.MockFiles[filepath.Join(relativeRoot, "file4")] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 1024 * 1024 * 10},
		}
		server.MockFiles[filepath.Join(relativeRoot, "file5")] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 100},
		}
		if relativeRoot != "" {
			server.MockFiles[relativeRoot] = mockFile{
				File: files_sdk.File{Type: "directory"},
			}
		}

		job := client.Downloader(DownloaderParams{RemotePath: relativeRoot, LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 6)
		require.NoError(t, job.Statuses[0].Err())

		for k, v := range server.MockFiles {
			f, err := os.Open(filepath.Join(root, k))
			require.NoError(t, err)
			stat, err := f.Stat()
			require.NoError(t, err)
			if !stat.IsDir() {
				assert.Equal(t, v.Size, stat.Size())
			}
		}
	}

	t.Run("list folder from a path", func(t *testing.T) {
		multipleFiles("a-root", t)
	})

	t.Run("multiple files from root", func(t *testing.T) {
		multipleFiles("", t)
	})

	t.Run("PreserveTimes with mtime", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		mtime := time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC).Truncate(time.Millisecond)
		server.MockFiles["small-file-with-size.txt"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 1999, Mtime: &mtime},
		}
		job := client.Downloader(DownloaderParams{RemotePath: "small-file-with-size.txt", LocalPath: root + "/", PreserveTimes: true})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		f, err := os.Open(filepath.Join(root, "small-file-with-size.txt"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, int64(1999), stat.Size())
		assert.Equal(t, mtime, stat.ModTime().UTC())
	})

	t.Run("PreserveTimes with providedMtime", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		providedMtime := time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC).Truncate(time.Millisecond)
		server.MockFiles["small-file-with-size.txt"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 1999, Mtime: lib.Time(time.Now()), ProvidedMtime: &providedMtime},
		}
		job := client.Downloader(DownloaderParams{RemotePath: "small-file-with-size.txt", LocalPath: root + "/", PreserveTimes: true})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		f, err := os.Open(filepath.Join(root, "small-file-with-size.txt"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, int64(1999), stat.Size())
		assert.Equal(t, providedMtime, stat.ModTime().UTC())
	})

	t.Run("sync already downloaded", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["taco.png"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 100},
		}
		taco, err := os.Create(filepath.Join(root, "taco.png"))
		assert.NoError(t, err)
		_, err = taco.Write(make([]byte, 100))
		require.NoError(t, err)
		require.NoError(t, taco.Close())
		job := client.Downloader(DownloaderParams{Sync: true, RemotePath: "taco.png", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		assert.Equal(t, status.Skipped, job.Statuses[0].Status())
	})

	t.Run("sync does not exist locally", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["taco.png"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 100},
		}
		job := client.Downloader(DownloaderParams{Sync: true, RemotePath: "taco.png", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		assert.Equal(t, status.Complete, job.Statuses[0].Status())
	})

	t.Run("sync is out of date locally by size", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["taco.png"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 100},
		}
		taco, err := os.Create(filepath.Join(root, "taco.png"))
		assert.NoError(t, err)
		require.NoError(t, taco.Close())
		job := client.Downloader(DownloaderParams{Sync: true, RemotePath: "taco.png", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		assert.Equal(t, status.Complete, job.Statuses[0].Status())
	})

	t.Run("no overwrite file exists", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["taco.png"] = mockFile{
			File: files_sdk.File{},
		}
		taco, err := os.Create(filepath.Join(root, "taco.png"))
		assert.NoError(t, err)
		require.NoError(t, taco.Close())
		job := client.Downloader(DownloaderParams{NoOverwrite: true, RemotePath: "taco.png", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		assert.Equal(t, status.FileExists, job.Statuses[0].Status())
	})

	t.Run("no overwrite file does not exists", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["taco.png"] = mockFile{
			File: files_sdk.File{},
		}
		job := client.Downloader(DownloaderParams{NoOverwrite: true, RemotePath: "taco.png", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		assert.Equal(t, status.Complete, job.Statuses[0].Status())
		f, err := os.Open(filepath.Join(root, "taco.png"))
		require.NoError(t, err)
		stat, err := f.Stat()
		require.NoError(t, err)
		assert.Equal(t, int64(0), stat.Size())
	})

	t.Run("local directory is privileged", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["taco.png"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 100},
		}

		require.NoError(t, os.Mkdir(filepath.Join(root, "restricted"), 0000))

		t.Cleanup(func() {
			require.NoError(t, os.Chmod(filepath.Join(root, "restricted"), 0777))
		})

		job := client.Downloader(DownloaderParams{Sync: true, RemotePath: "taco.png", LocalPath: filepath.Join(root, "restricted") + string(os.PathSeparator)})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.True(t, os.IsPermission(job.Statuses[0].Err()))
		assert.Equal(t, status.Errored, job.Statuses[0].Status())
	})

	t.Run("local path is invalid", func(t *testing.T) {
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["taco.png"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 100},
		}

		job := client.Downloader(DownloaderParams{Sync: true, RemotePath: "taco.png", LocalPath: "invalid\000path"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.Error(t, job.Statuses[0].Err())
		require.Contains(t, job.Statuses[0].Err().Error(), "invalid argument")
		assert.Equal(t, status.Errored, job.Statuses[0].Status())
	})

	t.Run("with a temp path", func(t *testing.T) {
		root := t.TempDir()
		temp := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["taco.png"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 100},
		}

		job := client.Downloader(DownloaderParams{RemotePath: "taco.png", LocalPath: root, TempPath: temp})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
	})

	t.Run("with a temp path and a privileged local directory", func(t *testing.T) {
		root := t.TempDir()
		temp := t.TempDir()
		restricted := filepath.Join(root, "restricted")
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["taco.png"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 100},
		}

		require.NoError(t, os.Mkdir(restricted, 0000))

		t.Cleanup(func() {
			require.NoError(t, os.Chmod(restricted, 0777))
		})

		job := client.Downloader(DownloaderParams{RemotePath: "taco.png", LocalPath: restricted, TempPath: temp})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.True(t, os.IsPermission(job.Statuses[0].Err()))
		assert.Equal(t, status.Errored, job.Statuses[0].Status())
		_, err := os.Open(filepath.Join(temp, "taco.png.download"))
		require.Error(t, err)
	})
}

func TestDownloadV2PreallocatedTempFileWriteAt(t *testing.T) {
	size := int64(20 * 1024 * 1024)
	source := bytes.Repeat([]byte("a"), int(size))
	ranger := &downloadV2TestRangeFile{
		data:        source,
		downloadURI: "https://bucket.s3.us-east-1.amazonaws.com/native.bin?X-Amz-Signature=test",
		info: Info{File: files_sdk.File{
			DisplayName: "native.bin",
			Path:        "native.bin",
			Type:        "file",
			Size:        size,
			Crc32:       "crc32",
		}, sizeTrust: TrustedSizeValue},
	}
	tmpPath := filepath.Join(t.TempDir(), "native.bin.download")
	reportStatus := downloadV2TestStatus(ranger, ranger.info, DownloaderParams{
		AdaptiveConcurrency: true,
		Manager:             manager.Build(2, 1),
	}, tmpPath)

	used, finalSize, err := runDownloadV2IfSupported(context.Background(), reportStatus, ranger.info, tmpPath, 0)
	require.True(t, used)
	require.NoError(t, err)
	assert.Equal(t, size, finalSize)
	assert.Equal(t, size, reportStatus.TransferBytes())
	assert.ElementsMatch(t, []downloadV2TestRange{{off: 0, end: 16*1024*1024 - 1}, {off: 16 * 1024 * 1024, end: size - 1}}, ranger.Ranges())

	written, err := os.ReadFile(tmpPath)
	require.NoError(t, err)
	assert.Equal(t, source, written)
}

func TestDownloadV2KeepsStatusQueuedUntilAdaptivePartSlot(t *testing.T) {
	size := int64(20 * 1024 * 1024)
	source := bytes.Repeat([]byte("a"), int(size))
	ranger := &downloadV2TestRangeFile{
		data:        source,
		downloadURI: "https://bucket.s3.us-east-1.amazonaws.com/native.bin?X-Amz-Signature=test",
		info: Info{File: files_sdk.File{
			DisplayName: "native.bin",
			Path:        "native.bin",
			Type:        "file",
			Size:        size,
			Crc32:       "crc32",
		}, sizeTrust: TrustedSizeValue},
	}
	tmpPath := filepath.Join(t.TempDir(), "native.bin.download")
	params := DownloaderParams{
		AdaptiveConcurrency: true,
		Manager:             manager.Build(2, 1),
	}
	reportStatus := downloadV2TestStatus(ranger, ranger.info, params, tmpPath)
	downloadingEvents := 0
	reportStatus.Job().RegisterFileEvent(func(JobFile) {
		downloadingEvents++
	}, status.Downloading)
	file, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_RDWR, 0644)
	require.NoError(t, err)

	blockedManager := lib.NewAdaptiveConcurrencyManagerWithConfig(lib.AdaptiveConcurrencyConfig{
		MaxConcurrency: 1,
		InitialTarget:  1,
		MinTarget:      1,
	})
	blockedManager.Wait()
	released := false
	done := make(chan error, 1)
	defer func() {
		if !released {
			released = true
			blockedManager.DoneNeutral()
			select {
			case <-done:
			case <-time.After(time.Second):
			}
		}
	}()

	engine := newDownloadV2Engine(reportStatus, ranger, file, downloadV2TargetS3, size, 0, downloadV2KnownSizePartSize(downloadV2TargetS3, size), params)
	engine.manager = blockedManager
	go func() {
		done <- engine.Run(context.Background())
	}()

	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, status.Queued, reportStatus.Status())
	assert.Empty(t, ranger.Ranges())

	released = true
	blockedManager.DoneNeutral()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for download v2 to finish")
	}
	assert.Equal(t, 1, downloadingEvents)
	assert.Equal(t, size, reportStatus.TransferBytes())
	assert.NotEmpty(t, ranger.Ranges())
}

func TestDownloadV2RequiresExplicitAdaptiveConcurrency(t *testing.T) {
	size := int64(20 * 1024 * 1024)
	ranger := &downloadV2TestRangeFile{
		data:        make([]byte, size),
		downloadURI: "https://bucket.s3.us-east-1.amazonaws.com/native.bin?X-Amz-Signature=test",
		info: Info{File: files_sdk.File{
			DisplayName: "native.bin",
			Path:        "native.bin",
			Type:        "file",
			Size:        size,
			Crc32:       "crc32",
		}, sizeTrust: TrustedSizeValue},
	}
	tmpPath := filepath.Join(t.TempDir(), "native.bin.download")
	reportStatus := downloadV2TestStatus(ranger, ranger.info, DownloaderParams{
		Manager: manager.Build(2, 1),
	}, tmpPath)

	used, _, err := runDownloadV2IfSupported(context.Background(), reportStatus, ranger.info, tmpPath, 0)

	require.NoError(t, err)
	assert.False(t, used)
	assert.Empty(t, ranger.Ranges())
}

func TestDownloadV2FallsBackForUntrustedSize(t *testing.T) {
	size := int64(20 * 1024 * 1024)
	ranger := &downloadV2TestRangeFile{
		data: make([]byte, size),
		info: Info{File: files_sdk.File{
			DisplayName: "remote-mount.bin",
			Path:        "remote-mount.bin",
			Type:        "file",
			Size:        size,
			Crc32:       "crc32",
		}, sizeTrust: UntrustedSizeValue},
	}
	tmpPath := filepath.Join(t.TempDir(), "remote-mount.bin.download")
	reportStatus := downloadV2TestStatus(ranger, ranger.info, DownloaderParams{
		AdaptiveConcurrency: true,
		Manager:             manager.Build(2, 1),
	}, tmpPath)

	used, _, err := runDownloadV2IfSupported(context.Background(), reportStatus, ranger.info, tmpPath, 0)
	require.NoError(t, err)
	assert.False(t, used)
	assert.Empty(t, ranger.Ranges())
}

func TestDownloadV2TruncatesFailedPreallocatedTempFileToContiguousPrefix(t *testing.T) {
	size := int64(20 * 1024 * 1024)
	ranger := &downloadV2TestRangeFile{
		data:               bytes.Repeat([]byte("b"), int(size)),
		downloadURI:        "https://bucket.s3.us-east-1.amazonaws.com/native.bin?X-Amz-Signature=test",
		failAfterOffset:    16 * 1024 * 1024,
		failAfterReadBytes: 1024 * 1024,
		info: Info{File: files_sdk.File{
			DisplayName: "native.bin",
			Path:        "native.bin",
			Type:        "file",
			Size:        size,
			Crc32:       "crc32",
		}, sizeTrust: TrustedSizeValue},
	}
	tmpPath := filepath.Join(t.TempDir(), "native.bin.download")
	reportStatus := downloadV2TestStatus(ranger, ranger.info, DownloaderParams{
		AdaptiveConcurrency: true,
		Manager:             manager.Build(1, 1),
	}, tmpPath)

	used, finalSize, err := runDownloadV2IfSupported(context.Background(), reportStatus, ranger.info, tmpPath, 0)
	require.True(t, used)
	require.Error(t, err)
	assert.Equal(t, int64(16*1024*1024), finalSize)
	stat, statErr := os.Stat(tmpPath)
	require.NoError(t, statErr)
	assert.Equal(t, int64(16*1024*1024), stat.Size())
}

func TestDownloadV2UsesDefaultTargetForGenericNonS3DownloadURIWithCrc32(t *testing.T) {
	size := int64(20 * 1024 * 1024)
	ranger := &downloadV2TestRangeFile{
		data:        make([]byte, size),
		downloadURI: "https://files.example.com/download/native.bin",
		info: Info{File: files_sdk.File{
			DisplayName: "native.bin",
			Path:        "native.bin",
			Type:        "file",
			Size:        size,
			Crc32:       "crc32",
		}, sizeTrust: TrustedSizeValue},
	}
	tmpPath := filepath.Join(t.TempDir(), "native.bin.download")
	reportStatus := downloadV2TestStatus(ranger, ranger.info, DownloaderParams{
		AdaptiveConcurrency: true,
		Manager:             manager.Build(2, 1),
	}, tmpPath)

	used, finalSize, err := runDownloadV2IfSupported(context.Background(), reportStatus, ranger.info, tmpPath, 0)

	require.NoError(t, err)
	assert.True(t, used)
	assert.Equal(t, size, finalSize)
	assert.NotEmpty(t, ranger.Ranges())
}

func TestDownloadV2FallsBackForSinglePartS3Download(t *testing.T) {
	size := int64(9)
	ranger := &downloadV2TestRangeFile{
		data:        []byte("123456789"),
		downloadURI: "https://bucket.s3.us-east-1.amazonaws.com/tiny.bin?X-Amz-Signature=test",
		info: Info{File: files_sdk.File{
			DisplayName: "tiny.bin",
			Path:        "tiny.bin",
			Type:        "file",
			Size:        size,
			Crc32:       "crc32",
		}, sizeTrust: TrustedSizeValue},
	}
	tmpPath := filepath.Join(t.TempDir(), "tiny.bin.download")
	reportStatus := downloadV2TestStatus(ranger, ranger.info, DownloaderParams{
		AdaptiveConcurrency: true,
		Manager:             manager.Build(2, 1),
	}, tmpPath)

	used, _, err := runDownloadV2IfSupported(context.Background(), reportStatus, ranger.info, tmpPath, 0)

	require.NoError(t, err)
	assert.False(t, used)
	assert.Empty(t, ranger.Ranges())
}

func TestDownloadV2UsesDefaultDownloadURI(t *testing.T) {
	size := int64(20 * 1024 * 1024)
	source := bytes.Repeat([]byte("a"), int(size))
	ranger := &downloadV2TestRangeFile{
		data:        source,
		downloadURI: "https://app-us-east-1.files.com/download/native.bin?X-Files-Date=20260605T211830Z&X-Files-Expires=180&jwt=test",
		info: Info{File: files_sdk.File{
			DisplayName: "default.bin",
			Path:        "default.bin",
			Type:        "file",
			Size:        size,
		}, sizeTrust: TrustedSizeValue},
	}
	tmpPath := filepath.Join(t.TempDir(), "default.bin.download")
	reportStatus := downloadV2TestStatus(ranger, ranger.info, DownloaderParams{
		AdaptiveConcurrency: true,
		Manager:             manager.Build(15, 1),
	}, tmpPath)

	used, finalSize, err := runDownloadV2IfSupported(context.Background(), reportStatus, ranger.info, tmpPath, 0)
	require.True(t, used)
	require.NoError(t, err)
	assert.Equal(t, size, finalSize)
	assert.Equal(t, size, reportStatus.TransferBytes())
	assert.ElementsMatch(t, []downloadV2TestRange{{off: 0, end: 16*1024*1024 - 1}, {off: 16 * 1024 * 1024, end: size - 1}}, ranger.Ranges())

	written, err := os.ReadFile(tmpPath)
	require.NoError(t, err)
	assert.Equal(t, source, written)
}

func TestDownloadV2ClassifiesS3AndDefaultDownloadURLs(t *testing.T) {
	target, ok := classifyDownloadV2URI("https://s3.amazonaws.com/bucket/key?X-Amz-Signature=test")
	require.True(t, ok)
	assert.Equal(t, downloadV2TargetS3, target)

	target, ok = classifyDownloadV2URI("https://uploads.example.com/download/native.bin")
	require.True(t, ok)
	assert.Equal(t, downloadV2TargetDefault, target)

	target, ok = classifyDownloadV2URI("https://uploads.example.com/download/native.bin", func(string) TransferV2TargetClass {
		return "custom"
	})
	require.True(t, ok)
	assert.Equal(t, TransferV2TargetClass("custom"), target)

	_, ok = classifyDownloadV2URI("://bad")
	assert.False(t, ok)
}

func TestDownloadV2SharedManagerDoesNotStartAtTinyFilePartCount(t *testing.T) {
	registry := downloadV2SharedAdaptiveManagerRegistry{}
	adaptiveManager := registry.get(downloadV2TargetS3, manager.AdaptiveDownloadV2ConcurrentFileParts, 1024, 16*1024*1024)

	assert.Equal(t, uploadV2S3InitialConcurrency, adaptiveManager.Snapshot().Target)
}

func TestAdaptiveTransferStatsReportsSharedDownloadManagers(t *testing.T) {
	resetDownloadV2SharedAdaptiveManagersForTest()
	first := downloadV2SharedAdaptiveManagers.get(downloadV2TargetDefault, 4, 20*1024*1024, 16*1024*1024)
	second := downloadV2SharedAdaptiveManagers.get(downloadV2TargetDefault, 4, 20*1024*1024, 16*1024*1024)
	idle := downloadV2SharedAdaptiveManagers.get(downloadV2TargetDefault, 8, 20*1024*1024, 16*1024*1024)
	require.Same(t, first, second)
	require.NotSame(t, first, idle)
	first.Wait()
	t.Cleanup(first.Done)

	stats := AdaptiveTransferStats()

	assert.Equal(t, 0, stats.Upload.Active)
	assert.Equal(t, 0, stats.Upload.Max)
	assert.Equal(t, 1, stats.Download.Active)
	assert.Equal(t, 4, stats.Download.Max)
}

func TestDownloadV2DefaultManagerStartsAtLegacyRangeConcurrency(t *testing.T) {
	adaptiveManager := lib.NewAdaptiveConcurrencyManagerWithConfig(downloadV2AdaptiveConcurrencyConfig(
		downloadV2TargetDefault,
		manager.AdaptiveDownloadV2ConcurrentFileParts,
		20*1024*1024,
		16*1024*1024,
	))

	assert.Equal(t, 15, adaptiveManager.Snapshot().Target)
}

func TestDownloadV2CopyAtCoalescesShortReadsBeforeWriting(t *testing.T) {
	expected := int64(downloadV2CopyBufferSize*2 + 123)
	writer := &downloadV2RecordingWriterAt{}
	reader := &downloadV2ShortChunkReader{remaining: expected, chunkSize: 4 * 1024}

	written, err := downloadV2CopyAt(writer, 0, expected, reader, nil)
	require.NoError(t, err)
	assert.Equal(t, expected, written)
	assert.Equal(t, []int{downloadV2CopyBufferSize, downloadV2CopyBufferSize, 123}, writer.writeSizes)
}

func resetDownloadV2SharedAdaptiveManagersForTest() {
	downloadV2SharedAdaptiveManagers.resetForTest()
}

type downloadV2TestRange struct {
	off int64
	end int64
}

type downloadV2TestRangeFile struct {
	data               []byte
	downloadURI        string
	info               Info
	mu                 sync.Mutex
	ranges             []downloadV2TestRange
	downloadURICalls   int
	failAfterOffset    int64
	failAfterReadBytes int64
}

func (f *downloadV2TestRangeFile) Stat() (fs.FileInfo, error) {
	return f.info, nil
}

func (f *downloadV2TestRangeFile) Read([]byte) (int, error) {
	return 0, io.EOF
}

func (f *downloadV2TestRangeFile) Close() error {
	return nil
}

func (f *downloadV2TestRangeFile) ReaderRange(off int64, end int64) (io.ReadCloser, error) {
	f.mu.Lock()
	f.ranges = append(f.ranges, downloadV2TestRange{off: off, end: end})
	f.mu.Unlock()
	reader := io.NewSectionReader(bytes.NewReader(f.data), off, end-off+1)
	if f.failAfterReadBytes > 0 && off >= f.failAfterOffset {
		return &downloadV2FailingReadCloser{reader: reader, remaining: f.failAfterReadBytes}, nil
	}
	return io.NopCloser(reader), nil
}

func (f *downloadV2TestRangeFile) downloadV2URI(context.Context) (string, error) {
	f.mu.Lock()
	f.downloadURICalls++
	f.mu.Unlock()
	return f.downloadURI, nil
}

func (f *downloadV2TestRangeFile) Ranges() []downloadV2TestRange {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]downloadV2TestRange(nil), f.ranges...)
}

type downloadV2FailingReadCloser struct {
	reader    io.Reader
	remaining int64
}

func (r *downloadV2FailingReadCloser) Read(p []byte) (int, error) {
	if r.remaining <= 0 {
		return 0, errors.New("forced range read failure")
	}
	if int64(len(p)) > r.remaining {
		p = p[:r.remaining]
	}
	n, err := r.reader.Read(p)
	r.remaining -= int64(n)
	if err != nil {
		return n, err
	}
	return n, nil
}

func (r *downloadV2FailingReadCloser) Close() error {
	return nil
}

type downloadV2ShortChunkReader struct {
	remaining int64
	chunkSize int
}

func (r *downloadV2ShortChunkReader) Read(p []byte) (int, error) {
	if r.remaining <= 0 {
		return 0, io.EOF
	}
	if len(p) > r.chunkSize {
		p = p[:r.chunkSize]
	}
	if int64(len(p)) > r.remaining {
		p = p[:r.remaining]
	}
	for i := range p {
		p[i] = 'x'
	}
	r.remaining -= int64(len(p))
	return len(p), nil
}

type downloadV2RecordingWriterAt struct {
	writeSizes []int
}

func (w *downloadV2RecordingWriterAt) WriteAt(p []byte, _ int64) (int, error) {
	w.writeSizes = append(w.writeSizes, len(p))
	return len(p), nil
}

func downloadV2TestStatus(file fs.File, info Info, params DownloaderParams, tmpPath string) *DownloadStatus {
	config := files_sdk.Config{}.Init()
	job := (&Job{
		Config:  config,
		Logger:  config.Logger,
		Params:  params,
		Manager: params.Manager,
	}).Init()
	job.SetManager(params.Manager)
	return &DownloadStatus{
		fsFile:     file,
		FileInfo:   info,
		file:       info.File,
		job:        job,
		localPath:  filepath.Join(filepath.Dir(tmpPath), info.Name()),
		remotePath: info.File.Path,
		status:     status.Queued,
		Mutex:      &sync.RWMutex{},
		TmpPath:    tmpPath,
	}
}

func TestDownload(t *testing.T) {
	mutex := &sync.Mutex{}
	t.Run("downloader", func(t *testing.T) {
		sourceFs := &FS{Context: context.Background()}
		destinationFs := lib.ReadWriteFs(lib.LocalFileSystem{})
		for _, tt := range lib.PathSpec(t, sourceFs.PathSeparator(), destinationFs.PathSeparator()) {
			t.Run(tt.Name, func(t *testing.T) {
				client, r, err := CreateClient(t.Name())
				if err != nil {
					t.Fatal(err)
				}
				config := client.Config
				sourceFs := (&FS{Context: context.Background()}).Init(config, false)
				lib.BuildPathSpecTest(t, mutex, tt, sourceFs, destinationFs, func(args lib.PathSpecArgs) lib.Cmd {
					return &CmdRunner{
						run: func() *Job {
							return downloader(context.Background(), sourceFs, DownloaderParams{config: config, RemotePath: args.Src, LocalPath: args.Dest, PreserveTimes: args.PreserveTimes})
						},
						args: []string{args.Src, args.Dest, "--times", fmt.Sprintf("%v", args.PreserveTimes)},
					}
				})
				r.Stop()
			})
		}
	})
}

func TestIgnoreDownload(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		ignorePattern  string
		includePattern string
		ignored        bool
	}{
		{
			name:    "nothing explicitly ignored or included",
			path:    "path/not/excluded/computers.pdf",
			ignored: false,
		},
		{
			name:          "ignore css files",
			path:          "path/excluded/main.css",
			ignorePattern: "*.css",
			ignored:       true,
		},
		{
			name:           "only include pdf files",
			path:           "path/excluded/computers.txt",
			includePattern: "*.pdf",
			ignored:        true,
		},
		{
			name:           "include directory path matches folder *.pdf",
			path:           "path/excluded/computers.pdf",
			includePattern: "*/excluded/*.pdf",
			ignored:        false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			var ignored, included *gitignore.GitIgnore
			ignored, err = ignore.New(tc.ignorePattern)
			if err != nil {
				t.Error(err)
			}
			if tc.includePattern != "" {
				included, err = ignore.New(tc.includePattern)
				if err != nil {
					t.Error(err)
				}
			}
			got := ignorePath(tc.path, ignored, included)
			if got != tc.ignored {
				t.Errorf("ignorePath(%q, %q, %q) = %v, want %v", tc.path, tc.ignorePattern, tc.includePattern, got, tc.ignored)
			}
		})
	}
}

func TestDownloadPauseResume(t *testing.T) {
	t.Run("pause preserves temp file", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["file.txt"] = mockFile{SizeTrust: TrustedSizeValue, File: files_sdk.File{Size: 19999999}}
		ctx, cancel := context.WithCancelCause(context.Background())
		defer cancel(nil)
		job := client.Downloader(
			DownloaderParams{RemotePath: "file.txt", LocalPath: filepath.Join(root, "file.txt")},
			files_sdk.WithContext(ctx),
		)
		downloadStarted := make(chan struct{}, 1)
		job.RegisterFileEvent(func(f JobFile) {
			select {
			case downloadStarted <- struct{}{}:
			default:
			}
		}, status.Downloading)
		job.Start()
		select {
		case <-downloadStarted:
		case <-job.Finished.C:
		}
		cancel(ErrJobPaused)
		job.Wait()

		tmpPath := existingTmpDownloadPath(filepath.Join(root, "file.txt"), "")
		assert.NotEmpty(t, tmpPath, "temp file should be preserved on pause")
	})

	t.Run("cancel removes temp file", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["file.txt"] = mockFile{SizeTrust: TrustedSizeValue, File: files_sdk.File{Size: 19999999}}
		ctx, cancel := context.WithCancelCause(context.Background())
		defer cancel(nil)
		job := client.Downloader(
			DownloaderParams{RemotePath: "file.txt", LocalPath: filepath.Join(root, "file.txt")},
			files_sdk.WithContext(ctx),
		)
		downloadStarted := make(chan struct{}, 1)
		job.RegisterFileEvent(func(f JobFile) {
			select {
			case downloadStarted <- struct{}{}:
			default:
			}
		}, status.Downloading)
		job.Start()
		select {
		case <-downloadStarted:
		case <-job.Finished.C:
		}
		cancel(fmt.Errorf("job canceled"))
		job.Wait()

		tmpPath := existingTmpDownloadPath(filepath.Join(root, "file.txt"), "")
		assert.Empty(t, tmpPath, "temp file should be removed on normal cancel")
	})

	t.Run("resume skips completed paths", func(t *testing.T) {
		setup := NewTestSetup()
		defer setup.TearDown()
		setup.MapFS["folder"] = &fstest.MapFile{
			Mode: fs.ModeDir,
			Sys:  files_sdk.File{DisplayName: "folder", Path: "folder", Type: "directory"},
		}
		setup.MapFS["folder/a.txt"] = &fstest.MapFile{
			Data: make([]byte, 100),
			Mode: fs.ModePerm,
			Sys:  files_sdk.File{DisplayName: "a.txt", Path: "folder/a.txt", Type: "file", Size: 100},
		}
		setup.MapFS["folder/b.txt"] = &fstest.MapFile{
			Data: make([]byte, 100),
			Mode: fs.ModePerm,
			Sys:  files_sdk.File{DisplayName: "b.txt", Path: "folder/b.txt", Type: "file", Size: 100},
		}
		alreadyDone := filepath.Join(setup.tempDir, "folder", "a.txt")
		setup.DownloaderParams = DownloaderParams{
			RemotePath:         "folder",
			LocalPath:          setup.tempDir + "/",
			EventsReporter:     setup.Reporter(),
			PriorJobCheckpoint: &JobDownloadCheckpoint{CompletedPaths: []string{alreadyDone}},
		}
		setup.Call()

		var statuses []status.Status
		for _, c := range setup.reporterCalls {
			if c.LocalPath == alreadyDone {
				statuses = append(statuses, c.Status)
			}
		}
		assert.Contains(t, statuses, status.Skipped)
		assert.NotContains(t, statuses, status.Complete)
	})

	t.Run("resume starts from existing temp file offset", func(t *testing.T) {
		root := t.TempDir()
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()
		fileSize := int64(19999999)
		server.MockFiles["file.txt"] = mockFile{SizeTrust: TrustedSizeValue, File: files_sdk.File{Size: fileSize}}

		localPath := filepath.Join(root, "file.txt")
		createCanonicalTmpFile(t, localPath, fileSize/2)

		job := client.Downloader(DownloaderParams{RemotePath: "file.txt", LocalPath: localPath})
		job.Start()
		job.Wait()

		require.NoError(t, job.Statuses[0].Err())
		fi, err := os.Stat(localPath)
		require.NoError(t, err)
		assert.Equal(t, fileSize, fi.Size())
	})
}

func createCanonicalTmpFile(t *testing.T, localPath string, size int64) string {
	t.Helper()
	tmpPath, err := tmpDownloadPath(localPath, "")
	require.NoError(t, err)
	err = os.WriteFile(tmpPath, make([]byte, size), 0644)
	require.NoError(t, err)
	return tmpPath
}

type CmdRunner struct {
	run    func() *Job
	stderr io.Writer
	stdout io.Writer
	args   []string
	*Job
}

func (c *CmdRunner) Run() error {
	c.Job = c.run()
	c.Job.Start()
	c.Job.Wait()
	for _, f := range c.Job.Sub(status.Errored).Statuses {
		c.stderr.Write([]byte(f.Err().Error()))
	}
	return nil
}

func (c *CmdRunner) Args() []string {
	return c.args
}

func (c *CmdRunner) SetOut(w io.Writer) {
	c.stdout = w
}

func (c *CmdRunner) SetErr(stderr io.Writer) {
	c.stderr = stderr
}
