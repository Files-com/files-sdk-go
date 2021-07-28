package file

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/Files-com/files-sdk-go/file/status"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/stretchr/testify/assert"
)

type MockerDownloader struct {
	Calls       []files_sdk.FileDownloadParams
	ReturnError error
	*sync.Mutex
}

func (m *MockerDownloader) Download(_ context.Context, p files_sdk.FileDownloadParams) (files_sdk.File, error) {
	m.Lock()
	m.Calls = append(m.Calls, p)
	p.OnDownload(&http.Response{ContentLength: 100})
	p.Writer.Write([]byte("one hundred bytes-----------------------------------------------------------------------------------"))
	m.Unlock()
	return files_sdk.File{}, m.ReturnError
}

type ReporterCall struct {
	status.Report
	err error
}

type TestSetup struct {
	files         []Entity
	reporterCalls []ReporterCall
	downloader    *MockerDownloader
	DownloadFolderParams
	rootDestination string
	tempDir         string
}

func NewTestSetup() *TestSetup {
	t := &TestSetup{}
	t.downloader = &MockerDownloader{Mutex: &sync.Mutex{}}
	t.TempDir()
	return t
}

func (setup *TestSetup) Reporter() Reporter {
	return func(status status.Report, err error) {
		setup.reporterCalls = append(setup.reporterCalls, ReporterCall{Report: status, err: err})
	}
}

func (setup *TestSetup) TempDir() error {
	var err error
	setup.tempDir, err = ioutil.TempDir("", "test")

	return err
}

func (setup *TestSetup) TearDown() error {
	return os.RemoveAll(setup.tempDir)
}

func (setup *TestSetup) Call() (status.Job, error) {
	return downloadFolder(
		context.Background(),
		setup.files,
		setup.downloader,
		setup.DownloadFolderParams,
	)
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
	setup.files = append(setup.files, Entity{file: files_sdk.File{Path: "some-path/taco.png", Size: 100, Type: "file"}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: "some-path"}, Reporter: setup.Reporter(), RootDestination: setup.RootDestination()}
	setup.rootDestination = "some-path/"
	_, err := setup.Call()
	assert.NoError(err)

	assert.Equal(1, setup.reporterCalls[0].Job().Count())
	assert.Equal(4, len(setup.reporterCalls))
	assert.Equal(status.Queued, setup.reporterCalls[0].Type())
	assert.Equal(status.Downloading, setup.reporterCalls[1].Type())
	assert.Equal(status.Downloading, setup.reporterCalls[2].Type())
	assert.Equal(status.Complete, setup.reporterCalls[3].Type())
	assert.Equal("some-path/taco.png", setup.reporterCalls[0].File().Path)
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes())
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes())
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes())
	assert.Equal(1, len(setup.downloader.Calls))
	assert.Equal("some-path/taco.png", setup.downloader.Calls[0].Path)

	assert.Equal(true, setup.reporterCalls[0].Job().AllEnded())
	assert.Equal(int64(100), setup.reporterCalls[0].Job().TransferBytes())
	assert.Equal(int64(100), setup.reporterCalls[0].Job().TotalBytes())

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_more_than_one_file(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.File{Path: "some-path/taco.png", Size: 100, Type: "file"}})
	setup.files = append(setup.files, Entity{file: files_sdk.File{Path: "some-path/pizza.png", Size: 102, Type: "file"}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: "some-path"}, Reporter: setup.Reporter(), RootDestination: setup.RootDestination()}
	setup.rootDestination = "some-path"

	_, err := setup.Call()
	assert.NoError(err)

	assert.Equal(2, setup.reporterCalls[0].Job().Count())
	assert.Equal(8, len(setup.reporterCalls))
	assert.ElementsMatch([]string{"some-path/taco.png", "some-path/pizza.png"}, []string{setup.reporterCalls[0].File().Path, setup.reporterCalls[1].File().Path})
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes())
	assert.Equal(2, len(setup.downloader.Calls))
	assert.ElementsMatch([]string{"some-path/taco.png", "some-path/pizza.png"}, []string{setup.downloader.Calls[0].Path, setup.downloader.Calls[1].Path})

	assert.Equal(true, setup.reporterCalls[1].Job().AllEnded())
	assert.Equal(int64(200), setup.reporterCalls[1].Job().TransferBytes())
	assert.Equal(int64(200), setup.reporterCalls[1].Job().TotalBytes())

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_sync_already_downloaded(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.File{Path: "taco.png", Size: 100, Type: "file", Mtime: time.Now().AddDate(0, -1, 0)}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: ""}, Sync: true, Reporter: setup.Reporter(), RootDestination: setup.RootDestination()}
	setup.rootDestination = ""
	_, err := os.Create(filepath.Join(setup.tempDir, "taco.png"))
	assert.NoError(err)
	_, err = setup.Call()
	assert.NoError(err)

	assert.Equal(2, len(setup.reporterCalls))
	assert.Equal(status.Queued, setup.reporterCalls[0].Type())
	assert.Equal(status.Skipped, setup.reporterCalls[1].Type())
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes())
	assert.Equal(0, len(setup.downloader.Calls))

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_sync_not_already_downloaded(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.File{Path: "taco.png", Size: 100, Type: "file", Mtime: time.Now().AddDate(0, 1, 0)}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: ""}, Sync: true, Reporter: setup.Reporter(), RootDestination: setup.RootDestination()}
	setup.rootDestination = ""
	_, err := os.Create(filepath.Join(setup.tempDir, "taco.png"))
	assert.NoError(err)
	_, err = setup.Call()
	assert.NoError(err)

	assert.Equal(4, len(setup.reporterCalls))
	assert.Equal("taco.png", setup.reporterCalls[0].File().Path)
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes())
	assert.Equal(status.Queued, setup.reporterCalls[0].Type())
	assert.Equal(status.Downloading, setup.reporterCalls[1].Type())
	assert.Equal(status.Downloading, setup.reporterCalls[2].Type())
	assert.Equal(status.Complete, setup.reporterCalls[3].Type())
	assert.Equal(1, len(setup.downloader.Calls))
	assert.Equal("taco.png", setup.downloader.Calls[0].Path)

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_Entity_error(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.File{Path: "taco.png", Size: 100, Type: "file"}, error: fmt.Errorf("something Happened")})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: "some-path"}, Reporter: setup.Reporter(), RootDestination: setup.RootDestination()}
	setup.rootDestination = "some-path/"
	job, err := setup.Call()
	assert.Equal(0, len(setup.reporterCalls))
	assert.Errorf(err, "something Happened")
	assert.Equal(0, job.Count())
	assert.Equal(true, job.AllEnded())
	assert.Equal(int64(0), job.TotalBytes())
	assert.Equal(int64(0), job.TransferBytes())

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_download_file(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.File{Path: "taco.png", Size: 100, Type: "file"}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: "some-path"}, Reporter: setup.Reporter(), RootDestination: setup.RootDestination()}
	setup.rootDestination = "taco.png"

	job, err := setup.Call()
	assert.NoError(err)

	assert.Equal(1, job.Count())
	assert.Equal(4, len(setup.reporterCalls))
	assert.Equal("taco.png", setup.reporterCalls[0].File().Path)
	assert.Equal(int64(0), setup.reporterCalls[0].TransferBytes())
	assert.Equal(1, len(setup.downloader.Calls))
	assert.Equal("taco.png", setup.downloader.Calls[0].Path)

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_OnDownload(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.File{Path: "taco.png", Size: 100, Type: "file"}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: "some-path"}, Reporter: setup.Reporter(), RootDestination: setup.RootDestination()}
	setup.rootDestination = "taco.png"

	_, err := setup.Call()
	assert.NoError(err)

	assert.Equal(4, len(setup.reporterCalls))
	assert.Equal(1, len(setup.downloader.Calls))

	assert.Equal(int64(100), setup.reporterCalls[3].File().Size, "Updates with real file size")

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_Download_error(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.File{Path: "taco.png", Size: 100, Type: "file"}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: "some-path"}, Reporter: setup.Reporter(), RootDestination: setup.RootDestination()}
	setup.rootDestination = "taco.png"
	setup.downloader.ReturnError = fmt.Errorf("download error")

	_, err := setup.Call()
	assert.NoError(err)

	assert.Equal(4, len(setup.reporterCalls))
	assert.Equal("taco.png", setup.reporterCalls[0].File().Path)
	assert.Contains(setup.reporterCalls[0].Destination(), "taco.png")
	var errStatus error
	for _, r := range setup.reporterCalls {
		if r.err != nil {
			errStatus = r.err
		}
	}
	assert.Errorf(errStatus, "download error")

	assert.NoError(setup.TearDown())
}
