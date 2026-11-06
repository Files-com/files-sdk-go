package test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

func CreateConfig(fixture string) (files_sdk.Config, *recorder.Recorder, error) {
	return CreateConfigWithConfig(fixture, files_sdk.Config{})
}

// CreateConfigWithConfig binds replay metadata normalization to the supplied
// configuration's API endpoint. Set a custom endpoint before calling this helper.
func CreateConfigWithConfig(fixture string, config files_sdk.Config) (files_sdk.Config, *recorder.Recorder, error) {
	config = config.Init()
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

	var timeline *replayTimeline
	var transport http.RoundTripper = r
	if r.Mode() == recorder.ModeReplaying {
		timeline, err = newReplayTimeline(filepath.Join("fixtures", fixture), time.Now())
		if err != nil {
			return files_sdk.Config{}, r, err
		}
		if timeline != nil {
			apiRoot, parseErr := url.Parse(config.RootPath())
			if parseErr != nil {
				return files_sdk.Config{}, r, parseErr
			}
			transport = replayTimeTransport{base: r, timeline: timeline, apiRoot: apiRoot}
		}
	}
	config = config.SetCustomClient(&http.Client{Transport: transport})

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Filesapi-Key")
		return nil
	})
	r.SetMatcher(func(r *http.Request, i cassette.Request) bool {
		if timeline != nil {
			i.URL = timeline.normalize(i.URL)
		}
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
