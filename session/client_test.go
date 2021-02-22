package session

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go"
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

func TestClient_Delete(t *testing.T) {
	client, r, err := CreateClient("Delete")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	os.Unsetenv("FILES_API_KEY")
	client.Config.SessionId = "9f799aff7f518514a0b6b5cfd1047e73dddd5cf5"
	_, err = client.Delete(files_sdk.SessionDeleteParams{})
	assert.Nil(err, "logout returns success")
}
