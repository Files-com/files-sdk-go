package file

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/lib"
	"github.com/stretchr/testify/assert"
)

func Test_UploadFile(t *testing.T) {
	client := Client{}
	if client.GetAPIKey() == "" {
		return
	}
	assert := assert.New(t)

	client.Config.Debug = lib.Bool(true)
	uploadPath := "../LICENSE"
	_, err := client.UploadFile(uploadPath, nil)
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
