package file

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/stretchr/testify/assert"
)

type MockerDownloader struct {
	Calls       []files_sdk.FileDownloadParams
	ReturnError error
}

func (m *MockerDownloader) Download(p files_sdk.FileDownloadParams) (files_sdk.File, error) {
	m.Calls = append(m.Calls, p)
	return files_sdk.File{}, m.ReturnError
}

type ReporterCall struct {
	incDownloadedBytes int64
	file               files_sdk.File
	destination        string
	err                error
	onlyMessage        string
	totalFiles         int
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
	t.downloader = &MockerDownloader{}
	t.TempDir()
	return t
}

func (setup *TestSetup) TempDir() error {
	var err error
	setup.tempDir, err = ioutil.TempDir("", "test")

	return err
}

func (setup *TestSetup) TearDown() error {
	return os.RemoveAll(setup.tempDir)
}

func (setup *TestSetup) Call() error {
	return downloadFolder(
		setup.files,
		setup.downloader,
		setup.DownloadFolderParams,
		setup.RootDestination(),
		func(incDownloadedBytes int64, file files_sdk.File, destination string, err error, onlyMessage string, totalFiles int) {
			setup.reporterCalls = append(setup.reporterCalls, ReporterCall{incDownloadedBytes, file, destination, err, onlyMessage, totalFiles})
		},
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
	setup.files = append(setup.files, Entity{file: files_sdk.Folder{Path: "some-path/taco.png", Size: 100, Type: "file"}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: "some-path"}}
	setup.rootDestination = "some-path/"

	assert.NoError(setup.Call())

	assert.Equal(1, len(setup.reporterCalls))
	assert.Equal("some-path/taco.png", setup.reporterCalls[0].file.Path)
	assert.Equal(int64(0), setup.reporterCalls[0].incDownloadedBytes)
	assert.Equal(1, len(setup.downloader.Calls))
	assert.Equal("some-path/taco.png", setup.downloader.Calls[0].Path)

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_more_than_one_file(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.Folder{Path: "some-path/taco.png", Size: 100, Type: "file"}})
	setup.files = append(setup.files, Entity{file: files_sdk.Folder{Path: "some-path/pizza.png", Size: 102, Type: "file"}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: "some-path"}}
	setup.rootDestination = "some-path"

	assert.NoError(setup.Call())

	assert.Equal(2, len(setup.reporterCalls))
	assert.ElementsMatch([]string{"some-path/taco.png", "some-path/pizza.png"}, []string{setup.reporterCalls[0].file.Path, setup.reporterCalls[1].file.Path})
	assert.Equal(int64(0), setup.reporterCalls[0].incDownloadedBytes)
	assert.Equal(2, len(setup.downloader.Calls))
	assert.ElementsMatch([]string{"some-path/taco.png", "some-path/pizza.png"}, []string{setup.downloader.Calls[0].Path, setup.downloader.Calls[1].Path})

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_sync_already_downloaded(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.Folder{Path: "taco.png", Size: 100, Type: "file", Mtime: time.Now().AddDate(0, -1, 0)}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: ""}, Sync: true}
	setup.rootDestination = ""
	_, err := os.Create(filepath.Join(setup.tempDir, "taco.png"))
	assert.NoError(err)
	assert.NoError(setup.Call())

	assert.Equal(1, len(setup.reporterCalls))
	assert.Equal("", setup.reporterCalls[0].file.Path)
	assert.Equal("No files to download", setup.reporterCalls[0].onlyMessage)
	assert.Equal(int64(0), setup.reporterCalls[0].incDownloadedBytes)
	assert.Equal(0, len(setup.downloader.Calls))

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_sync_not_already_downloaded(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.Folder{Path: "taco.png", Size: 100, Type: "file", Mtime: time.Now().AddDate(0, 1, 0)}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: ""}, Sync: true}
	setup.rootDestination = ""
	_, err := os.Create(filepath.Join(setup.tempDir, "taco.png"))
	assert.NoError(err)
	assert.NoError(setup.Call())

	assert.Equal(1, len(setup.reporterCalls))
	assert.Equal("taco.png", setup.reporterCalls[0].file.Path)
	assert.Equal("", setup.reporterCalls[0].onlyMessage)
	assert.Equal(int64(0), setup.reporterCalls[0].incDownloadedBytes)
	assert.Equal(1, len(setup.downloader.Calls))
	assert.Equal("taco.png", setup.downloader.Calls[0].Path)

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_Entity_error(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.Folder{Path: "taco.png", Size: 100, Type: "file"}, error: fmt.Errorf("something Happened")})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: "some-path"}}
	setup.rootDestination = "some-path/"

	assert.Error(setup.Call(), "something Happened")

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_download_file(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.Folder{Path: "taco.png", Size: 100, Type: "file"}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: "some-path"}}
	setup.rootDestination = "taco.png"

	assert.NoError(setup.Call())

	assert.Equal(1, len(setup.reporterCalls))
	assert.Equal("taco.png", setup.reporterCalls[0].file.Path)
	assert.Equal("", setup.reporterCalls[0].onlyMessage)
	assert.Equal(int64(0), setup.reporterCalls[0].incDownloadedBytes)
	assert.Equal(1, len(setup.downloader.Calls))
	assert.Equal("taco.png", setup.downloader.Calls[0].Path)

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_OnDownload(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.Folder{Path: "taco.png", Size: 100, Type: "file"}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: "some-path"}}
	setup.rootDestination = "taco.png"

	assert.NoError(setup.Call())

	assert.Equal(1, len(setup.downloader.Calls))
	response := &http.Response{ContentLength: 1000}
	assert.Equal(int64(100), setup.reporterCalls[0].file.Size)
	setup.downloader.Calls[0].OnDownload(response)
	assert.Equal(int64(1000), setup.reporterCalls[1].file.Size, "Updates with real file size")

	assert.NoError(setup.TearDown())
}

func Test_downloadFolder_Download_error(t *testing.T) {
	assert := assert.New(t)
	setup := NewTestSetup()
	setup.files = append(setup.files, Entity{file: files_sdk.Folder{Path: "taco.png", Size: 100, Type: "file"}})
	setup.DownloadFolderParams = DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: "some-path"}}
	setup.rootDestination = "taco.png"
	setup.downloader.ReturnError = fmt.Errorf("download error")

	assert.NoError(setup.Call())

	assert.Equal(2, len(setup.reporterCalls))
	assert.Equal("taco.png", setup.reporterCalls[0].file.Path)
	assert.Equal("download error", setup.reporterCalls[1].err.Error())

	assert.NoError(setup.TearDown())
}
