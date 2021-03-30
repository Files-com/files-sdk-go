package file

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/folder"
	"github.com/Files-com/files-sdk-go/lib"
	"github.com/dnaeon/go-vcr/cassette"
	recorder "github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func CreateClient(fixture string) (*Client, *recorder.Recorder, error) {
	client := Client{}
	var r *recorder.Recorder
	var err error
	if os.Getenv("GITLAB") != "" {
		fmt.Println("using ModeReplaying")
		r, err = recorder.NewAsMode(filepath.Join("fixtures", fixture), recorder.ModeReplaying, nil)
	} else {
		r, err = recorder.New(filepath.Join("fixtures", fixture))
	}
	if err != nil {
		return &client, r, err
	}

	httpClient := &http.Client{
		Transport: r,
	}
	client.Config.Debug = lib.Bool(true)
	client.HttpClient = httpClient

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Filesapi-Key")
		return nil
	})
	return &client, r, nil
}

func TestClient_UploadFolder(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	resultsMapMutex := sync.RWMutex{}
	results := make(map[string][]int64)

	_, err = client.UploadFolder(
		&UploadParams{
			Source:      "../lib",
			Destination: "golib",
			ProgressReporter: func(source string, file files_sdk.File, newBytesCount int64, batchStats UploadBatchStats, err error) {
				resultsMapMutex.Lock()
				results[file.Path] = append(results[file.Path], newBytesCount)
				resultsMapMutex.Unlock()
			},
		})
	assert.NoError(err)

	assert.Equal(10, len(results))
	assert.Contains(results, "golib/bool.go")
	assert.Contains(results, "golib/export_params.go")
	assert.Contains(results, "golib/interface.go")
	assert.Contains(results, "golib/iter.go")
	assert.Contains(results, "golib/string.go")
	assert.Contains(results, "golib/required_test.go")
	assert.Contains(results, "golib/required.go")
	assert.Contains(results, "golib/query.go")
	assert.Contains(results, "golib/progresswriter.go")
	assert.Contains(results, "golib/iter_test.go")
}

func TestClient_UploadFile(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	uploadPath := "../LICENSE"
	_, err = client.UploadFile(&UploadParams{Source: uploadPath})
	assert.NoError(err)
	_, err1 := os.Stat("../tmp")
	if os.IsNotExist(err1) {
		os.Mkdir("../tmp", 0700)
	}
	tempFile, err := ioutil.TempFile("../tmp", "LICENSE")
	if err != nil {
		panic(err)
	}
	downloadPath, err := filepath.Abs(filepath.Dir(tempFile.Name()))
	if err != nil {
		panic(err)
	}
	downloadPath = path.Join(downloadPath, tempFile.Name())
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "LICENSE"}, downloadPath)
	assert.NoError(err)

	assert.Equal(file.DisplayName, "LICENSE")

	downloadData, err := ioutil.ReadFile(downloadPath)
	if err != nil {
		panic(err)
	}
	localData, err := ioutil.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_UploadFolder_as_file2(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_as_file2")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	uploadPath := "../LICENSE"
	_, err = client.UploadFolder(&UploadParams{Source: uploadPath})
	if err != nil {
		panic(err)
	}
	_, err1 := os.Stat("../tmp")
	if os.IsNotExist(err1) {
		os.Mkdir("../tmp", 0700)
	}
	tempFile, err := ioutil.TempFile("../tmp", "LICENSE")
	if err != nil {
		panic(err)
	}
	downloadPath, err := filepath.Abs(filepath.Dir(tempFile.Name()))
	if err != nil {
		panic(err)
	}
	downloadPath = path.Join(downloadPath, tempFile.Name())
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "LICENSE"}, downloadPath)
	assert.NoError(err)

	assert.Equal(file.DisplayName, "LICENSE")

	downloadData, err := ioutil.ReadFile(downloadPath)
	if err != nil {
		panic(err)
	}
	localData, err := ioutil.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_DownloadFolder(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	folderClient := folder.Client{Config: client.Config}

	folderClient.Create(files_sdk.FolderCreateParams{Path: "TestClient_DownloadFolder"})
	folderClient.Create(files_sdk.FolderCreateParams{Path: filepath.Join("TestClient_DownloadFolder", "nested_1")})
	folderClient.Create(files_sdk.FolderCreateParams{Path: filepath.Join("TestClient_DownloadFolder", "nested_1", "nested_2")})

	client.Upload(strings.NewReader("testing 1"), filepath.Join("TestClient_DownloadFolder", "1.text"), &UploadProgress{})
	client.Upload(strings.NewReader("testing 2"), filepath.Join("TestClient_DownloadFolder", "2.text"), &UploadProgress{})
	client.Upload(strings.NewReader("testing 3"), filepath.Join("TestClient_DownloadFolder", "nested_1", "nested_2", "3.text"), &UploadProgress{})

	assert := assert.New(t)

	it, err := folderClient.ListFor(files_sdk.FolderListForParams{
		PerPage: 1,
		Path:    "TestClient_DownloadFolder",
	})

	assert.NoError(err)
	folders := files_sdk.FolderCollection{}
	for it.Next() {
		folders = append(folders, it.Folder())
	}

	assert.Len(folders, 3, "something is wrong with cursor")

	var results []string
	err = client.DownloadFolder(
		files_sdk.FolderListForParams{Path: "./TestClient_DownloadFolder"},
		"download",
		func(incDownloadedBytes int64, file files_sdk.File, destination string, err error, message string) {
			if err != nil {
				results = append(results, fmt.Sprint(file.Path, err))
			} else {
				results = append(results, fmt.Sprint(
					fmt.Sprintf("%d bytes ", incDownloadedBytes),
					fmt.Sprintf("%s => %s", file.Path, destination),
				))
			}
		},
	)

	assert.NoError(err)

	var expected []string
	expected = append(expected, "9 bytes TestClient_DownloadFolder/2.text => download/TestClient_DownloadFolder/2.text")
	expected = append(expected, "9 bytes TestClient_DownloadFolder/1.text => download/TestClient_DownloadFolder/1.text")
	expected = append(expected, "9 bytes TestClient_DownloadFolder/nested_1/nested_2/3.text => download/TestClient_DownloadFolder/nested_1/nested_2/3.text")
	assert.Equal(6, len(results))
	assert.Subset(results, expected)
	os.RemoveAll("download")
}
