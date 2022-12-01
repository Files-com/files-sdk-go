package test

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

func CreateConfig(fixture string) (files_sdk.Config, *recorder.Recorder, error) {
	config := files_sdk.Config{}
	var r *recorder.Recorder
	var err error
	if os.Getenv("GITLAB") != "" {
		fmt.Println("using ModeReplaying")
		r, err = recorder.NewAsMode(filepath.Join("fixtures", fixture), recorder.ModeReplaying, nil)
	} else {
		r, err = recorder.New(filepath.Join("fixtures", fixture))
	}
	if err != nil {
		return config, r, err
	}

	httpClient := &http.Client{
		Transport: r,
	}
	config.SetHttpClient(httpClient)

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Filesapi-Key")
		return nil
	})
	return config, r, nil
}
