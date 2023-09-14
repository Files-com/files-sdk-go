package test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

func CreateConfig(fixture string) (files_sdk.Config, *recorder.Recorder, error) {
	var r *recorder.Recorder
	var err error
	if os.Getenv("GITLAB") != "" {
		fmt.Println("using ModeReplaying")
		r, err = recorder.NewAsMode(filepath.Join("fixtures", fixture), recorder.ModeReplaying, nil)
	} else {
		r, err = recorder.New(filepath.Join("fixtures", fixture))
	}
	if err != nil {
		return files_sdk.Config{}, r, err
	}

	config := files_sdk.Config{}.Init().SetCustomClient(&http.Client{
		Transport: r,
	})

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Filesapi-Key")
		return nil
	})
	r.SetMatcher(func(r *http.Request, i cassette.Request) bool {
		if cassette.DefaultMatcher(r, i) {
			if r.Body != nil {
				io.ReadAll(r.Body)
				r.Body.Close()
			}

			return true
		}
		return false
	})
	return config, r, nil
}
