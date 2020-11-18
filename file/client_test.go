package file

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/lib"
	recorder "github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func CreateClient(fixture string) (*Client, *recorder.Recorder, error) {
	client := Client{}
	r, err := recorder.New(filepath.Join("fixtures", fixture))
	if err != nil {
		return &client, r, err
	}

	httpClient := &http.Client{
		Transport: r,
	}
	client.Config.Debug = lib.Bool(true)
	client.HttpClient = httpClient
	return &client, r, nil
}

func TestClient_UploadFolder(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)

	firstRun := true
	var results []string
	files, err := client.UploadFolder("../../go/lib", lib.String("lib"), func(source string, file files_sdk.File, largestSize int, largestFilePath int, totalUploads int, err error) {
		if firstRun {
			results = append(results, fmt.Sprint("Number of files/directories ", totalUploads))
			firstRun = false
		}
		if err != nil {
			results = append(results, fmt.Sprint(file.Path, err))
		} else {
			results = append(results, fmt.Sprint(
				fmt.Sprintf("%-"+strconv.Itoa(len(strconv.Itoa(largestSize)))+"d bytes", file.Size),
				fmt.Sprintf("%-"+strconv.Itoa(largestFilePath)+"s => %s", source, file.Path),
			))
		}
	})
	if err != nil {
		panic(err)
	}

	assert.Equal(len(files), 9)
	var expected []string
	expected = append(expected, "Number of files/directories 9")
	expected = append(expected, "68   bytes../../go/lib/string.go        => lib/string.go")
	expected = append(expected, "691  bytes../../go/lib/required_test.go => lib/required_test.go")
	expected = append(expected, "1087 bytes../../go/lib/required.go      => lib/required.go")
	expected = append(expected, "332  bytes../../go/lib/query.go         => lib/query.go")
	expected = append(expected, "231  bytes../../go/lib/export_params.go => lib/export_params.go")
	expected = append(expected, "68   bytes../../go/lib/string.go        => lib/string.go")
	expected = append(expected, "3141 bytes../../go/lib/iter.go          => lib/iter.go")
	expected = append(expected, "1593 bytes../../go/lib/iter_test.go     => lib/iter_test.go")
	expected = append(expected, "75   bytes../../go/lib/interface.go     => lib/interface.go")
	assert.Subset(results, expected)
}

func TestClient_UploadFile(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFile")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	uploadPath := "../LICENSE"
	_, err = client.UploadFile(uploadPath, nil)
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
	if err != nil {
		panic(err)
	}

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

func TestClient_UploadFolder_as_file(t *testing.T) {
	client, r, err := CreateClient("TestClient_UploadFolder_as_file")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)

	uploadPath := "../LICENSE"
	_, err = client.UploadFolder(uploadPath, nil)
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
	if err != nil {
		panic(err)
	}

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
