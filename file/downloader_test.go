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
		assert.Equal(t, time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC).Truncate(time.Millisecond), stat.ModTime().UTC().Truncate(time.Millisecond))
		assert.Contains(t, paths, "some-path/taco.png")
		assert.Equal(t, int64(0), setup.reporterCalls[0].TransferBytes)
	}

	t.Log("it uses ProvidedMtime")
	{
		stat, err := os.Stat(filepath.Join(setup.tempDir, "pizza.png"))
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC).Truncate(time.Millisecond), stat.ModTime().UTC().Truncate(time.Millisecond))
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
	require.NoError(t, taco.Close())
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
	taco, err := os.Create(filepath.Join(setup.tempDir, "taco.png"))
	assert.NoError(t, err)
	assert.NoError(t, taco.Close())
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
