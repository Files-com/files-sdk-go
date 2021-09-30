package file

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
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
	buildDownloadTest(func(job *status.Job) {
		events := status.EventsReporter{}
		var retryingEvents []status.Status
		events[status.Retrying] = func(file status.File) {
			retryingEvents = append(retryingEvents, file.Status)
		}
		job.SetEventsReporter(events)
		RetryByPolicy(context.Background(), job, RetryAll, false)
		assert.Equal(true, job.All(status.Complete))
		assert.Equal(3, len(retryingEvents), "sets a retrying status before starting")
	})

	buildDownloadTest(func(job *status.Job) {
		job.Start(false)
		originalStartWhen := job.Started.When
		RetryByPolicy(context.Background(), job, RetryUnfinished, false)
		assert.Equal(false, job.All(status.Complete))
		assert.Equal(2, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
		assert.Equal(true, job.Started.Called)
		assert.Equal(originalStartWhen, job.Started.When)
		assert.Equal(false, job.Scanning.Called)
		assert.Equal(false, job.EndScanning.Called)
		assert.Equal(false, job.Finished.Called)
		job.Finish()
	})

	buildDownloadTest(func(job *status.Job) {
		job.Start()
		job.Finish() // Finish already called, this happens in the desktop app
		originalFinishTime := job.Finished.When
		RetryByPolicy(context.Background(), job, RetryUnfinished, true)
		assert.Equal(false, job.All(status.Complete))
		assert.Equal(2, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
		assert.NotEqual(originalFinishTime, job.Finished.When)

		assert.Equal(true, job.Started.Called)
		assert.Equal(true, job.Scanning.Called)
		assert.Equal(true, job.EndScanning.Called)
		assert.Equal(true, job.Finished.Called)
	})

	buildDownloadTest(func(job *status.Job) {
		RetryByPolicy(context.Background(), job, RetryErroredIfSomeCompleted, false)
		assert.Equal(false, job.All(status.Complete))
		assert.Equal(2, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
	})

	buildDownloadTest(func(job *status.Job) {
		job.Statuses[0].(*DownloadStatus).lastByte = time.Time{}
		job.Statuses[1].(*DownloadStatus).lastByte = time.Time{}
		job.Statuses[2].(*DownloadStatus).lastByte = time.Time{}
		RetryByPolicy(context.Background(), job, RetryErroredIfSomeCompleted, false)
		assert.Equal(1, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
		assert.Equal(1, job.Count(status.Errored), "is not retried because completed happened before this error not after")
	})

	buildUploadTest(func(job *status.Job, _ *MockUploader) {
		RetryByPolicy(context.Background(), job, RetryAll, false)
		assert.Equal(true, job.All(status.Complete))
	})

	buildUploadTest(func(job *status.Job, _ *MockUploader) {
		RetryByPolicy(context.Background(), job, RetryUnfinished, false)
		assert.Equal(false, job.All(status.Complete))
		assert.Equal(2, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
		assert.Equal(false, job.Running(), "does change the original time")
	})

	buildUploadTest(func(job *status.Job, _ *MockUploader) {
		startTime := time.Now().AddDate(0, -1, 0)
		endTime := time.Now().AddDate(0, -1, 0)
		job.Timer.Runs = timer.Runs{timer.Run{
			Start:  startTime,
			Finish: endTime,
		}}
		uploadStatusErrored := job.Statuses[0].(*UploadStatus)
		uploadStatusErrored.Status = status.Complete
		RetryByPolicy(context.Background(), job, RetryUnfinished, false)
		assert.Equal(false, job.All(status.Complete))
		assert.Equal(2, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
		assert.Equal(startTime, job.StartTime(), "does not change the original time")
		assert.Equal(endTime, job.FinishTime(), "does not change the original time")
	})

	buildUploadTest(func(job *status.Job, _ *MockUploader) {
		uploadStatusErrored := job.Statuses[0].(*UploadStatus)
		uploadStatusComplete := job.Statuses[1].(*UploadStatus)
		uploadStatusComplete.lastByte = uploadStatusErrored.lastByte
		RetryByPolicy(context.Background(), job, RetryErroredIfSomeCompleted, false)
		assert.Equal(false, job.All(status.Complete))
		assert.Equal(1, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
		assert.Equal(1, job.Count(status.Errored))
		assert.Equal(false, job.Running(), "does not change the original time")
	})
}

func buildDownloadTest(test func(*status.Job)) {
	job := status.Job{Direction: direction.DownloadType, Manager: manager.Default()}.Init()
	temps := make([]string, 3)
	statuses := []status.Status{status.Errored, status.Complete, status.Queued}
	for i, s := range statuses {
		tempFile, err := ioutil.TempFile("../tmp", fmt.Sprintf("%v.txt", i))
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
		job.Add(&DownloadStatus{lastByte: time.Now(), LocalPath: tempFile.Name(), fsFile: localFile, Status: s, Job: job, File: files_sdk.File{DisplayName: fmt.Sprintf("%v.txt", i)}})
	}

	test(job)

	for _, tmp := range temps {
		os.Remove(tmp)
	}
}

func buildUploadTest(test func(*status.Job, *MockUploader)) {
	job := status.Job{Direction: direction.UploadType, Manager: manager.Default(), Params: UploadParams{}}.Init()
	job.GitIgnore, _ = ignore.New()
	temps := make([]string, 3)
	statuses := []status.Status{status.Errored, status.Complete, status.Queued}
	uploader := &MockUploader{}
	for i, s := range statuses {
		tempFile, err := ioutil.TempFile("../tmp", fmt.Sprintf("%v.txt", i))
		if err != nil {
			panic(err)
		}
		_, err = tempFile.Write([]byte("taco"))
		if err != nil {
			panic(err)
		}
		temps = append(temps, tempFile.Name())
		if err != nil {
			panic(err)
		}
		tempFile.Close()
		job.Add(&UploadStatus{Uploader: uploader, lastByte: time.Now(), LocalPath: tempFile.Name(), Status: s, Job: job, File: files_sdk.File{DisplayName: fmt.Sprintf("%v.txt", i)}})
	}

	test(job, uploader)

	for _, tmp := range temps {
		os.Remove(tmp)
	}
}
