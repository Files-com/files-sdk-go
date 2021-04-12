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
	client.Config.Debug = lib.Bool(false)
	client.HttpClient = httpClient

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Filesapi-Key")
		return nil
	})
	return &client, r, nil
}

func deletePath(client *Client, path string) {
	_, err := client.Delete(files_sdk.FileDeleteParams{Path: path})
	responseError, ok := err.(files_sdk.ResponseError)
	if ok && responseError.Type == "not-found" {
	} else if ok && responseError.Type == "processing-failure/folder-not-empty" {
		_, err = client.Delete(files_sdk.FileDeleteParams{Path: path, Recursive: lib.Bool(true)})
		responseError, ok = err.(files_sdk.ResponseError)
		if ok && responseError.Type == "not-found" {

		} else if ok {
			panic(err)
		}
	} else if ok {
		panic(err)
	}
}

func buildScenario(base string, client *Client) {
	folderClient := folder.Client{Config: client.Config}

	folderClient.Create(files_sdk.FolderCreateParams{Path: base})
	folderClient.Create(files_sdk.FolderCreateParams{Path: filepath.Join(base, "nested_1")})
	folderClient.Create(files_sdk.FolderCreateParams{Path: filepath.Join(base, "nested_1", "nested_2")})
	folderClient.Create(files_sdk.FolderCreateParams{Path: filepath.Join(base, "nested_1", "nested_2", "nested_3")})

	client.Upload(strings.NewReader("testing 3"), files_sdk.FileActionBeginUploadParams{Path: filepath.Join(base, "nested_1", "nested_2", "3.text")}, &UploadProgress{})
	client.Upload(strings.NewReader("testing 3"), files_sdk.FileActionBeginUploadParams{Path: filepath.Join(base, "nested_1", "nested_2", "nested_3", "4.text")}, &UploadProgress{})

}

func runDownloadScenario(path string, destination string, client *Client) ([]string, error) {
	var results []string
	err := client.DownloadFolder(
		files_sdk.FolderListForParams{Path: path},
		destination,
		func(incDownloadedBytes int64, file files_sdk.File, destination string, err error, message string, _ int) {
			if message != "" {
				results = append(results, message)
			}
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

	return results, err
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

	deletePath(client, "golib")
}

func TestClient_UploadFolder_Dot(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_Dot")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	resultsMapMutex := sync.RWMutex{}
	results := make(map[string][]int64)

	_, err = client.UploadFolder(
		&UploadParams{
			Source:      ".",
			Destination: "go-from-dot",
			ProgressReporter: func(source string, file files_sdk.File, newBytesCount int64, batchStats UploadBatchStats, err error) {
				resultsMapMutex.Lock()
				results[file.Path] = append(results[file.Path], newBytesCount)
				resultsMapMutex.Unlock()
			},
		})
	assert.NoError(err)

	assert.Contains(results, "go-from-dot/fixtures/TestClient_UploadFolder.yaml")
	assert.Contains(results, "go-from-dot/client_test.go")
	assert.Contains(results, "go-from-dot/client.go")
	assert.Contains(results, "go-from-dot/download.go")
	assert.Contains(results, "go-from-dot/upload.go")

	deletePath(client, "go-from-dot")
}

func TestClient_UploadFolder_Relative(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_Relative")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	resultsMapMutex := sync.RWMutex{}
	results := make(map[string][]int64)

	_, err = client.UploadFolder(
		&UploadParams{
			Source:      "fixtures",
			Destination: "file-fixtures",
			ProgressReporter: func(source string, file files_sdk.File, newBytesCount int64, batchStats UploadBatchStats, err error) {
				resultsMapMutex.Lock()
				results[file.Path] = append(results[file.Path], newBytesCount)
				resultsMapMutex.Unlock()
			},
		})
	assert.NoError(err)

	assert.Contains(results, "file-fixtures/TestClient_UploadFolder.yaml")

	deletePath(client, "file-fixtures")
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

func TestClient_UploadFile_To_Existing_Dir(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile_To_Existing_Dir")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	folderClient := folder.Client{Config: client.Config}
	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: "docs"})
	defer deletePath(client, "docs")

	assert.NoError(err)
	uploadPath := "../LICENSE"
	_, err = client.UploadFile(&UploadParams{Source: uploadPath, Destination: "docs"})
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
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "docs/LICENSE"}, downloadPath)
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

func TestClient_UploadFile_To_NonExistingPath(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile_To_NonExistingPath")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	deletePath(client, "taco")
	uploadPath := "../LICENSE"
	_, err = client.UploadFile(&UploadParams{Source: uploadPath, Destination: "taco"})
	defer deletePath(client, "taco")
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
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "taco"}, downloadPath)
	assert.NoError(err)

	assert.Equal("taco", file.DisplayName, "because the docs did not exist as a folder it becomes the file")

	downloadData, err := ioutil.ReadFile(downloadPath)
	assert.NoError(err)
	localData, err := ioutil.ReadFile(uploadPath)
	if err != nil {
		panic(err)
	}
	assert.Equal(string(downloadData), string(localData))
	defer os.Remove(tempFile.Name())
}

func TestClient_UploadFile_To_NonExistingPath_WithSlash(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile_To_NonExistingPath_WithSlash")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	assert.NoError(err)
	uploadPath := "../LICENSE"
	deletePath(client, "docs")
	_, err = client.UploadFile(&UploadParams{Source: uploadPath, Destination: "docs/"})
	defer deletePath(client, "docs")
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
	file, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "docs/LICENSE"}, downloadPath)
	assert.NoError(err)

	assert.Equal("file", file.Type)
	assert.Equal("LICENSE", file.DisplayName)

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

	buildScenario("TestClient_DownloadFolder", client)

	assert := assert.New(t)
	folderClient := folder.Client{Config: client.Config}

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

	results, err := runDownloadScenario("TestClient_DownloadFolder", "download", client)

	assert.NoError(err)

	var expected []string
	expected = append(expected, "9 bytes TestClient_DownloadFolder/2.text => download/2.text")
	expected = append(expected, "9 bytes TestClient_DownloadFolder/1.text => download/1.text")
	expected = append(expected, "9 bytes TestClient_DownloadFolder/nested_1/nested_2/3.text => download/nested_1/nested_2/3.text")
	assert.Subset(results, expected)
	os.RemoveAll("download")
}

func TestClient_DownloadFolder_Smart(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_Smart")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	buildScenario("TestClient_DownloadFolder_Smart", client)

	assert := assert.New(t)

	results, err := runDownloadScenario(filepath.Join("TestClient_DownloadFolder_Smart", "nested_1", "nested_2", "3.text"), "download/", client)

	assert.NoError(err)

	var expected []string
	expected = append(expected, "9 bytes TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text => download/3.text")
	assert.Subset(results, expected)

	results2, err := runDownloadScenario(filepath.Join("TestClient_DownloadFolder_Smart", "nested_1", "nested_2"), "download", client)

	assert.NoError(err)

	var expected2 []string
	expected = append(expected2, "9 bytes TestClient_DownloadFolder_Smart/nested_1/nested_2/3.text => download/3.text")
	expected = append(expected2, "9 bytes TestClient_DownloadFolder_Smart/nested_1/nested_2/nested_3/4.text => download/nested_3/4.text")
	assert.Subset(results2, expected2)

	os.RemoveAll("download")
}

func TestClient_DownloadFolder_file_to_file(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_file_to_file")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	buildScenario("TestClient_DownloadFolder_file_to_file", client)
	assert := assert.New(t)

	results, err := runDownloadScenario(filepath.Join("TestClient_DownloadFolder_file_to_file", "nested_1", "nested_2", "3.text"), "3.text", client)
	assert.NoError(err)

	var expected []string
	expected = append(expected, "9 bytes TestClient_DownloadFolder_file_to_file/nested_1/nested_2/3.text => 3.text")
	assert.Subset(results, expected)

	os.RemoveAll("3.text")
}

func TestClient_DownloadFolder_file_to_implicit(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_file_to_implicit")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	buildScenario("file_to_implicit", client)
	assert := assert.New(t)
	results, err := runDownloadScenario(filepath.Join("file_to_implicit", "nested_1", "nested_2", "3.text"), "", client)
	assert.NoError(err)

	var expected []string
	expected = append(expected, "9 bytes file_to_implicit/nested_1/nested_2/3.text => 3.text")
	assert.Subset(results, expected)

	os.RemoveAll("3.text")
}

func TestClient_DownloadFolder_file_only(t *testing.T) {
	client, r, err := CreateClient("TestClient_DownloadFolder_file_only")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	client.Upload(strings.NewReader("hello"), files_sdk.FileActionBeginUploadParams{Path: filepath.Join("i am at the root.text")}, &UploadProgress{})
	assert := assert.New(t)
	results, err := runDownloadScenario("i am at the root.text", "", client)
	assert.NoError(err)

	var expected []string
	expected = append(expected, "5 bytes i am at the root.text => i am at the root.text")
	assert.Subset(results, expected)

	os.RemoveAll("i am at the root.text")
}

func TestClient_DownloadToFile_No_files(t *testing.T) {
	assert := assert.New(t)
	client, r, err := CreateClient("TestClient_DownloadToFile_No_files")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	folderClient := folder.Client{Config: client.Config}

	folderClient.Create(files_sdk.FolderCreateParams{Path: "empty folder"})
	results, err := runDownloadScenario("empty folder", "", client)
	assert.NoError(err)
	var expected []string
	expected = append(expected, "No files to download")
	assert.Subset(results, expected)
}
