package file

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/Files-com/files-sdk-go/v2/file/manager"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/lib"
)

type ReporterCall struct {
	status.File
	err error
}

type TestSetup struct {
	files         []Entity
	reporterCalls []ReporterCall
	fstest.MapFS
	DownloaderParams
	rootDestination string
	tempDir         string
	*files_sdk.Config
}

func NewTestSetup() *TestSetup {
	t := &TestSetup{Config: &files_sdk.Config{}}
	t.MapFS = make(fstest.MapFS)
	err := t.TempDir()
	if err != nil {
		panic(err)
	}
	return t
}

func (setup *TestSetup) Reporter() status.EventsReporter {
	m := sync.Mutex{}
	events := make(status.EventsReporter)

	callback := func(status status.File) {
		m.Lock()
		setup.reporterCalls = append(setup.reporterCalls, ReporterCall{File: status})
		m.Unlock()
	}

	for _, s := range status.Included {
		events[s] = append(events[s], callback)
	}

	for _, s := range status.Excluded {
		events[s] = append(events[s], callback)
	}

	return events
}

func (setup *TestSetup) TempDir() error {
	var err error
	setup.tempDir, err = os.MkdirTemp("", "test")

	return err
}

func (setup *TestSetup) TearDown() error {
	return os.RemoveAll(setup.tempDir)
}

func (setup *TestSetup) Call() *status.Job {
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

	assert.Equal(t, 1, setup.reporterCalls[0].Job.Count())
	assert.Equal(t, 3, len(setup.reporterCalls))
	assert.Equal(t, status.Queued, setup.reporterCalls[0].Status)
	assert.Equal(t, status.Downloading, setup.reporterCalls[1].Status)
	assert.Equal(t, status.Complete, setup.reporterCalls[2].Status)
	assert.NoError(t, setup.reporterCalls[2].err)
	assert.Equal(t, "some-path/taco.png", setup.reporterCalls[0].File.Path)
	assert.Equal(t, int64(0), setup.reporterCalls[0].TransferBytes)

	assert.Equal(t, true, setup.reporterCalls[0].Job.All(status.Ended...))
	assert.Equal(t, int64(100), setup.reporterCalls[0].Job.TransferBytes())
	assert.Equal(t, int64(100), setup.reporterCalls[0].Job.TotalBytes())

	assert.NoError(t, setup.TearDown())
}

func Test_downloader_RemoteStartingSlash(t *testing.T) {
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
	setup.rootDestination = "some-path" + string(os.PathSeparator)
	setup.Call()

	fi, ok := setup.reporterCalls[0].Find(status.Errored)
	if ok {
		require.NoError(t, fi.Err())
	}
	assert.Equal(t, 1, setup.reporterCalls[0].Job.Count())
	assert.Equal(t, 3, len(setup.reporterCalls))
	assert.Equal(t, status.Queued, setup.reporterCalls[0].Status)
	assert.Equal(t, status.Downloading, setup.reporterCalls[1].Status)
	assert.Equal(t, status.Complete, setup.reporterCalls[2].Status)
	assert.NoError(t, setup.reporterCalls[2].err)
	assert.Equal(t, "some-path/taco.png", setup.reporterCalls[0].File.Path)
	assert.Equal(t, int64(0), setup.reporterCalls[0].TransferBytes)

	assert.Equal(t, true, setup.reporterCalls[0].Job.All(status.Ended...))
	assert.Equal(t, int64(100), setup.reporterCalls[0].Job.TransferBytes())
	assert.Equal(t, int64(100), setup.reporterCalls[0].Job.TotalBytes())

	assert.NoError(t, setup.TearDown())
}

func TestClient_Downloader(t *testing.T) {
	t.Run("small file with size", func(t *testing.T) {
		root := t.TempDir()
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["small-file-with-size.txt"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 1999},
		}
		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "small-file-with-size.txt", LocalPath: root + "/"})
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
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["large-file-with-size.txt"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 19999999},
		}
		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "large-file-with-size.txt", LocalPath: root + "/"})
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
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["large-file-with-size.txt"] = mockFile{
			SizeTrust:      TrustedSizeValue,
			File:           files_sdk.File{Size: 1024 * 1024 * 100},
			MaxConnections: 1,
		}
		m := manager.Build(1, 1)
		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "large-file-with-size.txt", LocalPath: root + "/", Manager: m})
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
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["large-file-with-size.txt"] = mockFile{
			SizeTrust:      TrustedSizeValue,
			File:           files_sdk.File{Size: 1024 * 1024 * 50},
			MaxConnections: 1,
		}
		m := manager.Build(1, 1)
		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "large-file-with-size.txt", LocalPath: root + "/", Manager: m})
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
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["large-file-with-size.txt"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 1024 * 1024 * 50},
		}
		m := manager.Build(10, 1, true)
		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "large-file-with-size.txt", LocalPath: root + "/", Manager: m})
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
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["large-file-with-no-size.txt"] = mockFile{
			SizeTrust: UntrustedSizeValue,
			File:      files_sdk.File{Size: 19999999},
		}

		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "large-file-with-no-size.txt", LocalPath: root + "/"})
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
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		realSize := int64((1024 * 1024 * 5) - 256)
		server.MockFiles["large-file-with-no-size.txt"] = mockFile{
			SizeTrust: UntrustedSizeValue,
			File:      files_sdk.File{Size: 1024 * 1024 * 100},
			RealSize:  &realSize,
		}

		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "large-file-with-no-size.txt", LocalPath: root + "/"})
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
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		serverBytesSent := int64((1024 * 1024 * 5) + 256)
		server.MockFiles["large-file-with-no-size.txt"] = mockFile{
			SizeTrust:       UntrustedSizeValue,
			File:            files_sdk.File{Size: 1024 * 1024 * 15},
			ServerBytesSent: &serverBytesSent,
		}

		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "large-file-with-no-size.txt", LocalPath: root + "/"})
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
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		serverBytesSent := int64(1024 * 1024 * 4)
		server.MockFiles["large-file-with-no-size.txt"] = mockFile{
			SizeTrust:       UntrustedSizeValue,
			File:            files_sdk.File{Size: 1024 * 1024 * 15},
			ServerBytesSent: &serverBytesSent,
		}

		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "large-file-with-no-size.txt", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.EqualError(t, job.Statuses[0].Err(), `received size did not match server send size
expected 4194304 bytes sent 5242880 received`)
		_, err := os.Open(filepath.Join(root, "large-file-with-no-size.txt"))
		require.Error(t, err)
	})

	t.Run("large file with bad size info real size is bigger", func(t *testing.T) {
		root := t.TempDir()
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		realSize := int64(20000000)
		server.MockFiles["file-with-mismatch-size-bigger"] = mockFile{
			SizeTrust: UntrustedSizeValue,
			File:      files_sdk.File{Size: 19999999},
			RealSize:  &realSize,
		}

		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "file-with-mismatch-size-bigger", LocalPath: root + "/"})
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
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		realSize := int64(19999999)
		server.MockFiles["file-with-mismatch-size-smaller"] = mockFile{
			SizeTrust: UntrustedSizeValue,
			File:      files_sdk.File{Size: 20000000},
			RealSize:  &realSize,
		}

		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "file-with-mismatch-size-smaller", LocalPath: root + "/"})
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
		server := FakeDownloadServer{T: t}.Do()
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

		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: relativeRoot, LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 5)
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
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		mtime := time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC).Truncate(time.Millisecond)
		server.MockFiles["small-file-with-size.txt"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 1999, Mtime: &mtime},
		}
		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "small-file-with-size.txt", LocalPath: root + "/", PreserveTimes: true})
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
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		providedMtime := time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC).Truncate(time.Millisecond)
		server.MockFiles["small-file-with-size.txt"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 1999, Mtime: lib.Time(time.Now()), ProvidedMtime: &providedMtime},
		}
		job := client.Downloader(context.Background(), DownloaderParams{RemotePath: "small-file-with-size.txt", LocalPath: root + "/", PreserveTimes: true})
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
		server := FakeDownloadServer{T: t}.Do()
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
		job := client.Downloader(context.Background(), DownloaderParams{Sync: true, RemotePath: "taco.png", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		assert.Equal(t, status.Skipped, job.Statuses[0].Status())
	})

	t.Run("sync does not exist locally", func(t *testing.T) {
		root := t.TempDir()
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["taco.png"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 100},
		}
		job := client.Downloader(context.Background(), DownloaderParams{Sync: true, RemotePath: "taco.png", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		assert.Equal(t, status.Complete, job.Statuses[0].Status())
	})

	t.Run("sync is out of date locally by size", func(t *testing.T) {
		root := t.TempDir()
		server := FakeDownloadServer{T: t}.Do()
		defer server.Shutdown()
		client := server.Client()
		server.MockFiles["taco.png"] = mockFile{
			SizeTrust: TrustedSizeValue,
			File:      files_sdk.File{Size: 100},
		}
		taco, err := os.Create(filepath.Join(root, "taco.png"))
		assert.NoError(t, err)
		require.NoError(t, taco.Close())
		job := client.Downloader(context.Background(), DownloaderParams{Sync: true, RemotePath: "taco.png", LocalPath: root + "/"})
		job.Start()
		job.Wait()
		assert.Len(t, job.Statuses, 1)
		require.NoError(t, job.Statuses[0].Err())
		assert.Equal(t, status.Complete, job.Statuses[0].Status())
	})
}

func TestDownload(t *testing.T) {
	mutex := &sync.Mutex{}
	t.Run("downloader", func(t *testing.T) {
		sourceFs := &FS{Context: context.Background()}
		destinationFs := lib.ReadWriteFs(lib.LocalFileSystem{})
		for _, tt := range lib.PathSpec(sourceFs.PathSeparator(), destinationFs.PathSeparator()) {
			t.Run(tt.Name, func(t *testing.T) {
				client, r, err := CreateClient(t.Name())
				if err != nil {
					t.Fatal(err)
				}
				config := client.Config
				sourceFs := (&FS{Context: context.Background()}).Init(config, false)
				lib.BuildPathSpecTest(t, mutex, tt, sourceFs, destinationFs, func(source, destination string) lib.Cmd {
					return &CmdRunner{
						run: func() *status.Job {
							return downloader(context.Background(), sourceFs, DownloaderParams{RemotePath: source, LocalPath: destination})
						},
						args: []string{source, destination},
					}
				})
				r.Stop()
			})
		}
	})
}

type CmdRunner struct {
	run    func() *status.Job
	stderr io.Writer
	stdout io.Writer
	args   []string
	*status.Job
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
