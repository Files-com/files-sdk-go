package file

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"

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
		*setup.Config,
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

	setup.DownloaderParams = DownloaderParams{RemotePath: "/some-path", EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
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

func TestClient_Downloader_path_spec(t *testing.T) {
	t.Run("files", func(t *testing.T) {
		for _, tt := range pathSpec() {
			t.Run(tt.name, func(t *testing.T) {
				srcFs := make(fstest.MapFS)
				filesDest, err := os.MkdirTemp("", "files-dest")
				assert.NoError(t, err)

				for _, e := range tt.src {
					fileType := "file"
					mode := fs.ModePerm
					if e.dir {
						fileType = "directory"
						mode = fs.ModeDir
					}
					_, displayName := filepath.Split(e.path)
					srcFs[e.path] = &fstest.MapFile{
						Data: nil,
						Mode: mode,
						Sys: files_sdk.File{
							DisplayName: displayName,
							Path:        e.path,
							Type:        fileType,
						},
					}
				}
				for _, e := range tt.dest {
					if !e.preexisting {
						continue
					}
					if e.dir {
						err = os.MkdirAll(filepath.Join(filesDest, e.path), 0750)
					} else {
						_, err = os.Create(filepath.Join(filesDest, e.path))
					}
					assert.NoError(t, err)
				}
				params := DownloaderParams{
					LocalPath:  strings.Join([]string{filesDest, tt.args.dest}, string(os.PathSeparator)),
					RemotePath: tt.args.src,
				}
				if tt.args.dest == "" {
					params.LocalPath = ""
				}
				t.Logf("RemotePath: %v, LocalPath: %v", params.RemotePath, params.LocalPath)
				originalDir, err := os.Getwd()
				require.NoError(t, err)
				err = os.Chdir(filesDest)
				require.NoError(t, err)
				job := downloader(context.Background(), srcFs, params, files_sdk.Config{})

				job.Start()
				job.Wait()
				err = os.Chdir(originalDir)
				require.NoError(t, err)
				for _, e := range tt.dest {
					fileInfo, err := os.Stat(filepath.Join(filesDest, e.path))
					require.NoError(t, err, e.path)
					assert.Equal(t, e.dir, fileInfo.IsDir(), e.path)
				}

				assert.NoError(t, os.RemoveAll(filesDest))
			})
		}
	})
}

func Test_downloadFolder_more_than_one_file(t *testing.T) {
	setup := NewTestSetup()
	setup.MapFS["some-path"] = &fstest.MapFile{
		Data: nil,
		Mode: fs.ModeDir,
		Sys: files_sdk.File{
			DisplayName:   "some-path",
			Path:          "some-path",
			Type:          "directory",
			ProvidedMtime: lib.Time(time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)),
			Mtime:         lib.Time(time.Now()),
		},
	}

	setup.MapFS["some-path/taco.png"] = &fstest.MapFile{
		Data: make([]byte, 100),
		Mode: fs.ModePerm,
		Sys: files_sdk.File{
			DisplayName: "taco.png",
			Path:        "some-path/taco.png",
			Type:        "file",
			Size:        100,
			Mtime:       lib.Time(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
		},
	}

	setup.MapFS["some-path/pizza.png"] = &fstest.MapFile{
		Data: make([]byte, 102),
		Mode: fs.ModePerm,
		Sys: files_sdk.File{
			DisplayName:   "pizza.png",
			Path:          "some-path/pizza.png",
			Type:          "file",
			Size:          102,
			ProvidedMtime: lib.Time(time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)),
			Mtime:         lib.Time(time.Now()),
		},
	}
	setup.DownloaderParams = DownloaderParams{
		RemotePath:     "some-path/",
		EventsReporter: setup.Reporter(),
		LocalPath:      setup.RootDestination(),
		PreserveTimes:  true,
	}
	setup.rootDestination = "some-path"

	job := setup.Call()
	job.Wait()
	statuses := make(map[string]int)
	var paths []string
	for _, call := range setup.reporterCalls {
		i, ok := statuses[call.Status.Name]
		if ok {
			statuses[call.Status.Name] = i + 1
		} else {
			statuses[call.Status.Name] = 1
		}
		paths = append(paths, call.File.Path)
	}
	t.Log("it goes through all statuses")
	{
		assert.Equal(t, 2, setup.reporterCalls[0].Job.Count())
		assert.Equal(t, map[string]int{"complete": 2, "downloading": 2, "queued": 2}, statuses)
		assert.Equal(t, 6, len(setup.reporterCalls))
	}

	t.Log("it uses Mtime")
	{
		stat, err := os.Stat(filepath.Join(setup.tempDir, "taco.png"))
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC), stat.ModTime().UTC())
		assert.Contains(t, paths, "some-path/taco.png")
		assert.Equal(t, int64(0), setup.reporterCalls[0].TransferBytes)
	}

	t.Log("it uses ProvidedMtime")
	{
		stat, err := os.Stat(filepath.Join(setup.tempDir, "pizza.png"))
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC), stat.ModTime().UTC())
		assert.Contains(t, paths, "some-path/pizza.png")
		assert.Equal(t, int64(0), setup.reporterCalls[0].TransferBytes)
	}

	t.Log("it all ends with correct bytes transferred")
	{
		assert.Equal(t, true, job.All(status.Ended...))
		assert.Equal(t, int64(202), job.TransferBytes())
		assert.Equal(t, int64(202), job.TotalBytes())
	}

	assert.NoError(t, setup.TearDown())
}

func Test_downloadFolder_sync_already_downloaded(t *testing.T) {
	setup := NewTestSetup()

	setup.MapFS["taco.png"] = &fstest.MapFile{
		Data: make([]byte, 100),
		Mode: fs.ModePerm,
		Sys:  files_sdk.File{DisplayName: "taco.png", Path: "taco.png", Type: "file", Size: 100},
	}
	setup.DownloaderParams = DownloaderParams{RemotePath: "", Sync: true, EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
	setup.rootDestination = ""
	taco, err := os.Create(filepath.Join(setup.tempDir, "taco.png"))
	assert.NoError(t, err)
	_, err = taco.Write(make([]byte, 100))
	require.NoError(t, err)
	setup.Call()

	assert.Equal(t, 2, len(setup.reporterCalls))
	assert.NoError(t, setup.reporterCalls[0].err)
	assert.Equal(t, status.Queued, setup.reporterCalls[0].Status)
	assert.Equal(t, status.Skipped, setup.reporterCalls[1].Status)
	assert.Equal(t, int64(0), setup.reporterCalls[0].TransferBytes)

	assert.NoError(t, setup.TearDown())
}

func Test_downloadFolder_sync_not_already_downloaded(t *testing.T) {
	setup := NewTestSetup()
	setup.MapFS["taco.png"] = &fstest.MapFile{
		Data:    make([]byte, 100),
		ModTime: time.Now().AddDate(0, 1, 0),
		Sys:     files_sdk.File{DisplayName: "taco.png", Path: "taco.png", Type: "file", Size: 100},
	}
	setup.DownloaderParams = DownloaderParams{RemotePath: "", Sync: true, EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
	setup.rootDestination = ""
	_, err := os.Create(filepath.Join(setup.tempDir, "taco.png"))
	assert.NoError(t, err)
	setup.Call()

	assert.Equal(t, 3, len(setup.reporterCalls))
	assert.Equal(t, "taco.png", setup.reporterCalls[0].File.Path)
	assert.Equal(t, int64(0), setup.reporterCalls[0].TransferBytes)
	assert.Equal(t, status.Queued, setup.reporterCalls[0].Status)
	assert.Equal(t, status.Downloading, setup.reporterCalls[1].Status)
	assert.Equal(t, status.Complete, setup.reporterCalls[2].Status)

	assert.NoError(t, setup.TearDown())
}

func Test_downloadFolder_sync_local_does_not_exist(t *testing.T) {
	setup := NewTestSetup()
	setup.MapFS["taco.png"] = &fstest.MapFile{
		Data:    make([]byte, 100),
		ModTime: time.Now().AddDate(0, 1, 0),
		Sys:     files_sdk.File{DisplayName: "taco.png", Path: "taco.png", Type: "file", Size: 100},
	}
	setup.DownloaderParams = DownloaderParams{RemotePath: "", Sync: true, EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
	setup.rootDestination = ""
	setup.Call()

	assert.Equal(t, 3, len(setup.reporterCalls))
	assert.Equal(t, "taco.png", setup.reporterCalls[0].File.Path)
	assert.Equal(t, int64(0), setup.reporterCalls[0].TransferBytes)
	assert.Equal(t, status.Queued, setup.reporterCalls[0].Status)
	assert.Equal(t, status.Downloading, setup.reporterCalls[1].Status)
	assert.Equal(t, status.Complete, setup.reporterCalls[2].Status)

	assert.NoError(t, setup.TearDown())
}

func Test_downloadFolder_download_file(t *testing.T) {
	setup := NewTestSetup()
	setup.MapFS["some-path/taco.png"] = &fstest.MapFile{
		Data:    make([]byte, 100),
		Mode:    fs.ModePerm,
		ModTime: time.Time{},
		Sys:     files_sdk.File{DisplayName: "taco.png", Path: "some-path/taco.png", Type: "file", Size: 100},
	}
	setup.DownloaderParams = DownloaderParams{RemotePath: "some-path", EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
	setup.rootDestination = "taco.png"

	job := setup.Call()

	assert.Equal(t, 1, job.Count())
	assert.Equal(t, 3, len(setup.reporterCalls))
	assert.Equal(t, "some-path/taco.png", setup.reporterCalls[0].File.Path)
	assert.Equal(t, int64(0), setup.reporterCalls[0].TransferBytes)

	assert.NoError(t, setup.TearDown())
}

func Test_downloadFolder_OnDownload(t *testing.T) {
	setup := NewTestSetup()
	setup.MapFS["some-path/taco.png"] = &fstest.MapFile{
		Data:    make([]byte, 100),
		ModTime: time.Time{},
		Sys:     files_sdk.File{DisplayName: "taco.png", Path: "some-path/taco.png", Type: "file", Size: 100},
	}
	setup.DownloaderParams = DownloaderParams{RemotePath: "some-path", EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
	setup.rootDestination = "taco.png"

	setup.Call()

	assert.Equal(t, 3, len(setup.reporterCalls))

	assert.Equal(t, int64(100), setup.reporterCalls[2].File.Size, "Updates with real file size")

	assert.NoError(t, setup.TearDown())
}
