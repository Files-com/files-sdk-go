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

	"github.com/Files-com/files-sdk-go/v2/lib/timer"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/ignore"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"
	"github.com/stretchr/testify/assert"
)

func TestRetryTransfers(t *testing.T) {
	assert := assert.New(t)
	t.Run("downloads", func(t *testing.T) {
		t.Run("RetryAll RetryCount 1", func(t *testing.T) {
			buildDownloadTest(func(job *status.Job) {
				events := status.EventsReporter{}
				eventsMutex := sync.Mutex{}
				var retryingEvents []status.Status
				job.RegisterFileEvent(func(file status.File) {
					eventsMutex.Lock()
					defer eventsMutex.Unlock()
					retryingEvents = append(retryingEvents, file.Status)
				}, status.Retrying)
				job.SetEventsReporter(events)
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryAll, RetryCount: 1}, false)
				assert.Equal(true, job.All(status.Complete))
				assert.Equal(3, len(retryingEvents), "sets a retrying status before starting")
			})
		})

		t.Run("RetryUnfinished RetryCount 1", func(t *testing.T) {
			buildDownloadTest(func(job *status.Job) {
				job.Start(false)
				originalStartWhen := job.Started.When()
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryUnfinished, RetryCount: 1}, false)
				assert.Equal(false, job.All(status.Complete))
				assert.Equal(2, job.Count(status.Complete))
				assert.Equal(1, job.Count(status.Queued))
				assert.Equal(true, job.Started.Called())
				assert.Equal(originalStartWhen, job.Started.When())
				assert.Equal(false, job.Scanning.Called())
				assert.Equal(false, job.EndScanning.Called())
				assert.Equal(false, job.Finished.Called())
				job.Finish()
			})
		})

		t.Run("RetryUnfinished RetryCount 1 signalEvents true", func(t *testing.T) {
			buildDownloadTest(func(job *status.Job) {
				job.Start()
				job.Finish() // Finish already called, this happens in the desktop app
				originalFinishTime := job.Finished.When()
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryUnfinished, RetryCount: 1}, true)
				assert.Equal(false, job.All(status.Complete))
				assert.Equal(2, job.Count(status.Complete))
				assert.Equal(1, job.Count(status.Queued))
				assert.NotEqual(originalFinishTime, job.Finished.When())

				assert.Equal(true, job.Started.Called())
				assert.Equal(true, job.Scanning.Called())
				assert.Equal(true, job.EndScanning.Called())
				assert.Equal(true, job.Finished.Called())
			})
		})

		t.Run("RetryErroredIfSomeCompleted RetryCount 1", func(t *testing.T) {
			buildDownloadTest(func(job *status.Job) {
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryErroredIfSomeCompleted, RetryCount: 1}, false)
				assert.Equal(false, job.All(status.Complete))
				assert.Equal(2, job.Count(status.Complete))
				assert.Equal(1, job.Count(status.Queued))
			})
		})

		t.Run("RetryErroredIfSomeCompleted RetryCount 1", func(t *testing.T) {
			buildDownloadTest(func(job *status.Job) {
				job.Statuses[0].(*DownloadStatus).lastByte = time.Time{}
				job.Statuses[1].(*DownloadStatus).lastByte = time.Time{}
				job.Statuses[2].(*DownloadStatus).lastByte = time.Time{}
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryErroredIfSomeCompleted, RetryCount: 1}, false)
				assert.Equal(1, job.Count(status.Complete))
				assert.Equal(1, job.Count(status.Queued))
				assert.Equal(1, job.Count(status.Errored), "is not retried because completed happened before this error not after")
			})
		})
	})
	t.Run("uploads", func(t *testing.T) {
		t.Run("RetryAll RetryCount 1", func(t *testing.T) {
			buildUploadTest(func(job *status.Job, _ *MockUploader) {
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryAll, RetryCount: 1}, false)
				assert.Equal(true, job.All(status.Complete))
			})
		})

		t.Run("RetryUnfinished RetryCount 1", func(t *testing.T) {
			buildUploadTest(func(job *status.Job, _ *MockUploader) {
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryUnfinished, RetryCount: 1}, false)
				assert.Equal(false, job.All(status.Complete))
				assert.Equal(2, job.Count(status.Complete))
				assert.Equal(1, job.Count(status.Queued))
				assert.Equal(false, job.Running(), "does change the original time")
			})
		})

		t.Run("RetryUnfinished RetryCount 2", func(t *testing.T) {
			buildUploadTest(func(job *status.Job, uploader *MockUploader) {
				uploader.uploadError = fmt.Errorf("bad things")
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryUnfinished, RetryCount: 2, Backoff: -1}, false)
				assert.Equal(false, job.All(status.Complete))
				assert.Equal(1, job.Count(status.Complete))
				assert.Equal(2, job.Count(status.Errored))
				assert.Equal(false, job.Running(), "does change the original time")
			}, status.Errored, status.Errored, status.Complete)
		})

		t.Run("RetryUnfinished RetryCount 1", func(t *testing.T) {
			buildUploadTest(func(job *status.Job, _ *MockUploader) {
				startTime := time.Now().AddDate(0, -1, 0)
				endTime := time.Now().AddDate(0, -1, 0)
				job.Timer.Runs = timer.Runs{timer.Run{
					Start:  startTime,
					Finish: endTime,
				}}
				uploadStatusErrored := job.Statuses[0].(*UploadStatus)
				uploadStatusErrored.status = status.Complete
				RetryByPolicy(context.Background(), job, RetryPolicy{Type: RetryUnfinished, RetryCount: 1}, false)
				assert.Equal(false, job.All(status.Complete))
				assert.Equal(2, job.Count(status.Complete))
				assert.Equal(1, job.Count(status.Queued))
				assert.Equal(startTime, job.StartTime(), "does not change the original time")
				assert.Equal(endTime, job.FinishTime(), "does not change the original time")
			})
		})
	})
}

func buildDownloadTest(test func(*status.Job)) {
	job := status.Job{Direction: direction.DownloadType, Manager: manager.Default(), Config: files_sdk.Config{}, Logger: (&files_sdk.Config{}).Logger()}.Init()
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
		job.Add(&DownloadStatus{Mutex: &sync.RWMutex{}, lastByte: time.Now(), localPath: tempFile.Name(), fsFile: localFile, status: s, job: job, file: files_sdk.File{DisplayName: fmt.Sprintf("%v.txt", i)}})
	}

	test(job)

	for _, tmp := range temps {
		os.Remove(tmp)
	}
	os.RemoveAll(tmpDir)
}

func buildUploadTest(test func(*status.Job, *MockUploader), statuses ...status.Status) {
	job := status.Job{Direction: direction.UploadType, Manager: manager.Default(), Params: UploaderParams{}, Logger: (&files_sdk.Config{}).Logger()}.Init()
	job.GitIgnore, _ = ignore.New()
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
		job.Add(&UploadStatus{Mutex: &sync.RWMutex{}, Uploader: uploader, lastByte: time.Now(), localPath: tempFile.Name(), status: s, job: job, file: files_sdk.File{DisplayName: fmt.Sprintf("%v.txt", i)}})
	}

	test(job, uploader)

	for _, tmp := range temps {
		os.Remove(tmp)
	}
	if err := os.RemoveAll(tmpDir); err != nil {
		panic(err)
	}
}
