package file

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/ignore"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lib/direction"
	"github.com/Files-com/files-sdk-go/v3/lib/timer"
	"github.com/stretchr/testify/assert"
)

func TestRetryTransfers(t *testing.T) {
	assert := assert.New(t)
	t.Run("downloads", func(t *testing.T) {
		t.Run("RetryAll RetryCount 1", func(t *testing.T) {
			buildDownloadTest(func(job *Job) {
				events := EventsReporter{}
				eventsMutex := sync.Mutex{}
				var retryingEvents []status.Status
				job.RegisterFileEvent(func(file JobFile) {
					eventsMutex.Lock()
					defer eventsMutex.Unlock()
					retryingEvents = append(retryingEvents, file.Status)
				}, status.Retrying)
				job.SetEventsReporter(events)
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryAll, RetryCount: 1}, false)
				assert.True(job.All(status.Complete))
				assert.Equal(3, len(retryingEvents), "sets a retrying status before starting")
			})
		})

		t.Run("RetryUnfinished RetryCount 1", func(t *testing.T) {
			buildDownloadTest(func(job *Job) {
				job.Start(false)
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryUnfinished, RetryCount: 1}, false)
				assert.False(job.All(status.Complete))
				assert.Equal(2, job.Count(status.Complete))
				assert.Equal(1, job.Count(status.Queued))
				assert.True(job.Started.Called())
				assert.False(job.Scanning.Called())
				assert.False(job.EndScanning.Called())
				assert.False(job.Finished.Called())
				job.Finish()
			})
		})

		t.Run("RetryUnfinished RetryCount 1 signalEvents true", func(t *testing.T) {
			buildDownloadTest(func(job *Job) {
				job.Start()
				job.Finish() // Finish already called, this happens in the desktop app
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryUnfinished, RetryCount: 1}, true)
				assert.False(job.All(status.Complete))
				assert.Equal(2, job.Count(status.Complete))
				assert.Equal(1, job.Count(status.Queued))

				assert.True(job.Started.Called())
				assert.True(job.Scanning.Called())
				assert.True(job.EndScanning.Called())
				assert.True(job.Finished.Called())
			})
		})
	})
	t.Run("uploads", func(t *testing.T) {
		t.Run("RetryAll RetryCount 1", func(t *testing.T) {
			buildUploadTest(func(job *Job, _ *MockUploader) {
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryAll, RetryCount: 1}, false)
				assert.True(job.All(status.Complete))
			})
		})

		t.Run("RetryUnfinished RetryCount 1", func(t *testing.T) {
			buildUploadTest(func(job *Job, _ *MockUploader) {
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryUnfinished, RetryCount: 1}, false)
				assert.False(job.All(status.Complete))
				assert.Equal(2, job.Count(status.Complete))
				assert.Equal(1, job.Count(status.Queued))
				assert.False(job.Running(), "does change the original time")
			})
		})

		t.Run("RetryUnfinished RetryCount 2", func(t *testing.T) {
			buildUploadTest(func(job *Job, uploader *MockUploader) {
				uploader.uploadError = fmt.Errorf("bad things")
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryUnfinished, RetryCount: 2, Backoff: -1}, false)
				assert.False(job.All(status.Complete))
				assert.Equal(1, job.Count(status.Complete))
				assert.Equal(2, job.Count(status.Errored))
				assert.False(job.Running(), "does change the original time")
			}, status.Errored, status.Errored, status.Complete)
		})

		t.Run("RetryUnfinished RetryCount 1", func(t *testing.T) {
			buildUploadTest(func(job *Job, _ *MockUploader) {
				startTime := time.Now().AddDate(0, -1, 0)
				endTime := time.Now().AddDate(0, -1, 0)
				job.Timer.Runs = timer.Runs{timer.Run{
					Start:  startTime,
					Finish: endTime,
				}}
				uploadStatusErrored := job.Statuses[0].(*UploadStatus)
				uploadStatusErrored.status = status.Complete
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryUnfinished, RetryCount: 1}, false)
				assert.False(job.All(status.Complete))
				assert.Equal(2, job.Count(status.Complete))
				assert.Equal(1, job.Count(status.Queued))
				assert.Equal(startTime, job.StartTime(), "does not change the original time")
				assert.Equal(endTime, job.FinishTime(), "does not change the original time")
			})
		})
	})
}

func buildDownloadTest(test func(*Job)) {
	job := (&Job{Direction: direction.DownloadType, Manager: manager.Default(), Config: files_sdk.Config{}.Init(), Logger: lib.NullLogger{}}).Init()
	temps := make([]string, 3)
	statuses := []status.Status{status.Errored, status.Complete, status.Queued}
	tmpDir, err := ioutil.TempDir(os.TempDir(), "client_test")
	if err != nil {
		log.Fatal(err)
	}
	for i, s := range statuses {
		tempFile, err := ioutil.TempFile(tmpDir, fmt.Sprintf("%v.txt", i))
		if err != nil {
			panic(err)
		}
		_, err = tempFile.Write([]byte("taco"))
		if err != nil {
			panic(err)
		}
		temps = append(temps, tempFile.Name())
		localFile, err := os.Open(tempFile.Name())
		if err != nil {
			panic(err)
		}
		tempFile.Close()
		job.Add(&DownloadStatus{Mutex: &sync.RWMutex{}, localPath: tempFile.Name(), fsFile: localFile, status: s, job: job, file: files_sdk.File{DisplayName: fmt.Sprintf("%v.txt", i)}})
	}

	test(job)

	for _, tmp := range temps {
		os.Remove(tmp)
	}
	os.RemoveAll(tmpDir)
}

func buildUploadTest(test func(*Job, *MockUploader), statuses ...status.Status) {
	job := (&Job{Direction: direction.UploadType, Manager: manager.Default(), Params: UploaderParams{}, Config: files_sdk.Config{}.Init(), Logger: lib.NullLogger{}}).Init()
	job.Ignore, _ = ignore.New()
	var temps []string
	if len(statuses) == 0 {
		statuses = []status.Status{status.Errored, status.Complete, status.Queued}
	}
	uploader := &MockUploader{}
	tmpDir, err := ioutil.TempDir(os.TempDir(), "retrytransfers")
	if err != nil {
		log.Fatal(err)
	}
	for i, s := range statuses {
		tempFile, err := ioutil.TempFile(tmpDir, fmt.Sprintf("%v.txt", i))
		if err != nil {
			panic(err)
		}
		_, err = tempFile.Write([]byte("taco"))
		if err != nil {
			panic(err)
		}
		temps = append(temps, tempFile.Name())
		tempFile.Close()
		job.Add(&UploadStatus{Mutex: &sync.RWMutex{}, Uploader: uploader, localPath: tempFile.Name(), status: s, job: job, file: files_sdk.File{DisplayName: fmt.Sprintf("%v.txt", i), Mtime: lib.Time(time.Now())}})
	}

	test(job, uploader)

	for _, tmp := range temps {
		os.Remove(tmp)
	}
	if err := os.RemoveAll(tmpDir); err != nil {
		panic(err)
	}
}
