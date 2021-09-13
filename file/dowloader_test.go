package file

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/Files-com/files-sdk-go/v2/file/status"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/stretchr/testify/assert"
)

type MockerDownloader struct {
	Calls       []files_sdk.FileDownloadParams
	ReturnError error
	*sync.Mutex
	files []Entity
}

func (m *MockerDownloader) Download(_ context.Context, p files_sdk.FileDownloadParams) (files_sdk.File, error) {
	m.Lock()
	m.Calls = append(m.Calls, p)
	p.OnDownload(&http.Response{ContentLength: 100})
	p.Writer.Write([]byte("one hundred bytes-----------------------------------------------------------------------------------"))
	m.Unlock()
	return files_sdk.File{}, m.ReturnError
}

func (m *MockerDownloader) Index(_ context.Context, fileQueue chan Entity, _ string) int {
	for _, file := range m.files {
		fileQueue <- file
	}

	return len(m.files)
}

type ReporterCall struct {
	status.File
	err error
}

type TestSetup struct {
	files         []Entity
	reporterCalls []ReporterCall
	fstest.MapFS
	DownloadFolderParams
	rootDestination string
	tempDir         string
}

func NewTestSetup() *TestSetup {
	t := &TestSetup{}
	t.MapFS = make(fstest.MapFS)
	t.TempDir()
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
		events[s] = callback
	}

	for _, s := range status.Excluded {
		events[s] = callback
	}

	return events
}

func (setup *TestSetup) TempDir() error {
	var err error
	setup.tempDir, err = ioutil.TempDir("", "test")

	return err
}

func (setup *TestSetup) TearDown() error {
	return os.RemoveAll(setup.tempDir)
}

func (setup *TestSetup) Call() *status.Job {
	job := downloader(
		context.Background(),
		setup.MapFS,
		setup.DownloadFolderParams,
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
	assert := assert.New(t)
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
		Sys:     files_sdk.File{DisplayName: "taco.png", Path: "some-path/taco.png", Type: "file"},
	}

	setup.DownloadFolderParams = DownloadFolderParams{RemotePath: "some-path", EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
	setup.rootDestination = "some-path/"
	setup.Call()

	assert.Equal(1, setup.reporterCalls[0].Job.Count())
	assert.Equal(3, len(setup.reporterCalls))
	assert.Equal(status.Queued, setup.reporterCalls[0].Status)
	assert.Equal(status.Downloading, setup.reporterCalls[1].Status)
	assert.Equal(status.Complete, setup.reporterCalls[2].Status)
	assert.NoError(setup.reporterCalls[2].err)
	assert.Equal("some-path/taco.png", setup.reporterCalls[0].File.Path)
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes)
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes)
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes)

	assert.Equal(true, setup.reporterCalls[0].Job.All(status.Ended...))
	assert.Equal(int64(100), setup.reporterCalls[0].Job.TransferBytes())
	assert.Equal(int64(100), setup.reporterCalls[0].Job.TotalBytes())

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_more_than_one_file(t *testing.T) {
	assert := assert.New(t)
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
		Sys:     files_sdk.File{DisplayName: "taco.png", Path: "some-path/taco.png", Type: "file"},
	}

	setup.MapFS["some-path/pizza.png"] = &fstest.MapFile{
		Data:    make([]byte, 102),
		Mode:    fs.ModePerm,
		ModTime: time.Time{},
		Sys:     files_sdk.File{DisplayName: "pizza.png", Path: "some-path/pizza.png", Type: "file"},
	}
	setup.DownloadFolderParams = DownloadFolderParams{RemotePath: "some-path", EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
	setup.rootDestination = "some-path"

	setup.Call()

	statuses := make(map[string]int)
	for _, call := range setup.reporterCalls {
		i, ok := statuses[call.Status.Name]
		if ok {
			statuses[call.Status.Name] = i + 1
		} else {
			statuses[call.Status.Name] = 1
		}
	}
	assert.Equal(2, setup.reporterCalls[0].Job.Count())
	assert.Equal(map[string]int{"complete": 2, "downloading": 2, "queued": 2}, statuses)
	assert.Equal(6, len(setup.reporterCalls))
	assert.Contains([]string{setup.reporterCalls[0].File.Path, setup.reporterCalls[1].File.Path}, "some-path/taco.png")
	assert.Contains([]string{setup.reporterCalls[0].File.Path, setup.reporterCalls[1].File.Path}, "some-path/pizza.png")
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes)

	assert.Equal(true, setup.reporterCalls[1].Job.All(status.Ended...))
	assert.Equal(int64(202), setup.reporterCalls[1].Job.TransferBytes())
	assert.Equal(int64(202), setup.reporterCalls[1].Job.TotalBytes())

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_sync_already_downloaded(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()

	setup.MapFS["taco.png"] = &fstest.MapFile{
		Data:    make([]byte, 100),
		Mode:    fs.ModePerm,
		ModTime: time.Now().AddDate(0, -1, 0),
		Sys:     files_sdk.File{DisplayName: "taco.png", Path: "taco.png", Type: "file"},
	}
	setup.DownloadFolderParams = DownloadFolderParams{RemotePath: "", Sync: true, EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
	setup.rootDestination = ""
	_, err := os.Create(filepath.Join(setup.tempDir, "taco.png"))
	assert.NoError(err)
	setup.Call()

	assert.Equal(2, len(setup.reporterCalls))
	assert.NoError(setup.reporterCalls[0].err)
	assert.Equal(status.Queued, setup.reporterCalls[0].Status)
	assert.Equal(status.Skipped, setup.reporterCalls[1].Status)
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes)

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_sync_not_already_downloaded(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.MapFS["taco.png"] = &fstest.MapFile{
		Data:    make([]byte, 100),
		ModTime: time.Now().AddDate(0, 1, 0),
		Sys:     files_sdk.File{DisplayName: "taco.png", Path: "taco.png", Type: "file"},
	}
	setup.DownloadFolderParams = DownloadFolderParams{RemotePath: "", Sync: true, EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
	setup.rootDestination = ""
	_, err := os.Create(filepath.Join(setup.tempDir, "taco.png"))
	assert.NoError(err)
	setup.Call()

	assert.Equal(3, len(setup.reporterCalls))
	assert.Equal("taco.png", setup.reporterCalls[0].File.Path)
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes)
	assert.Equal(status.Queued, setup.reporterCalls[0].Status)
	assert.Equal(status.Downloading, setup.reporterCalls[1].Status)
	assert.Equal(status.Complete, setup.reporterCalls[2].Status)

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_download_file(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.MapFS["some-path/taco.png"] = &fstest.MapFile{
		Data:    make([]byte, 100),
		Mode:    fs.ModePerm,
		ModTime: time.Time{},
		Sys:     files_sdk.File{DisplayName: "taco.png", Path: "some-path/taco.png", Type: "file"},
	}
	setup.DownloadFolderParams = DownloadFolderParams{RemotePath: "some-path", EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
	setup.rootDestination = "taco.png"

	job := setup.Call()

	assert.Equal(1, job.Count())
	assert.Equal(3, len(setup.reporterCalls))
	assert.Equal("some-path/taco.png", setup.reporterCalls[0].File.Path)
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes)

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_OnDownload(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.MapFS["some-path/taco.png"] = &fstest.MapFile{
		Data:    make([]byte, 100),
		ModTime: time.Time{},
		Sys:     files_sdk.File{DisplayName: "taco.png", Path: "some-path/taco.png", Type: "file", Size: 99},
	}
	setup.DownloadFolderParams = DownloadFolderParams{RemotePath: "some-path", EventsReporter: setup.Reporter(), LocalPath: setup.RootDestination()}
	setup.rootDestination = "taco.png"

	setup.Call()

	assert.Equal(3, len(setup.reporterCalls))

	assert.Equal(int64(100), setup.reporterCalls[2].File.Size, "Updates with real file size")

	assert.NoError(setup.TearDown())
}

func Test_tmpDownloadPath(t *testing.T) {
	assert := assert.New(t)

	path := tmpDownloadPath("you-wont-find-me")

	assert.Equal("you-wont-find-me.download", path)
	file, err := os.Create("find-me.download")
	defer func() {
		os.Remove(file.Name())
	}()
	if err != nil {
		panic(err)
	}
	file.Write([]byte("hello"))
	err = file.Close()
	if err != nil {
		panic(err)
	}
	path = tmpDownloadPath("find-me")

	assert.Equal(fmt.Sprintf("%v (1)", file.Name()), path, "it increments a number")
}
