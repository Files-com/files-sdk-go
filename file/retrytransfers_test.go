package file

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

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
		job.RetryPolicy = string(RetryAll)
		events := status.EventsReporter{}
		var retryingEvents []status.Status
		events[status.Retrying] = func(file status.File) {
			retryingEvents = append(retryingEvents, file.Status)
		}
		job.SetEventsReporter(events)
		RetryTransfers(context.Background(), job)
		assert.Equal(true, job.All(status.Complete))
		assert.Equal(3, len(retryingEvents), "sets a retrying status before starting")
	})

	buildDownloadTest(func(job *status.Job) {
		job.RetryPolicy = string(RetryErrored)
		RetryTransfers(context.Background(), job)
		assert.Equal(false, job.All(status.Complete))
		assert.Equal(2, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
	})

	buildDownloadTest(func(job *status.Job) {
		job.RetryPolicy = string(RetryErrored)
		job.Stopped = true
		RetryTransfers(context.Background(), job)
		assert.Equal(1, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
		assert.Equal(1, job.Count(status.Errored), "job was stopped so no retrying")
	})

	buildDownloadTest(func(job *status.Job) {
		job.RetryPolicy = string(RetryErroredIfSomeCompleted)
		RetryTransfers(context.Background(), job)
		assert.Equal(false, job.All(status.Complete))
		assert.Equal(2, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
	})

	buildDownloadTest(func(job *status.Job) {
		job.RetryPolicy = string(RetryErroredIfSomeCompleted)
		job.Statuses[0].(*DownloadStatus).lastByte = time.Time{}
		job.Statuses[1].(*DownloadStatus).lastByte = time.Time{}
		job.Statuses[2].(*DownloadStatus).lastByte = time.Time{}
		RetryTransfers(context.Background(), job)
		assert.Equal(1, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
		assert.Equal(1, job.Count(status.Errored), "is not retried because completed happened before this error not after")
	})

	buildUploadTest(func(job *status.Job, _ *MockUploader) {
		job.RetryPolicy = string(RetryAll)
		RetryTransfers(context.Background(), job)
		assert.Equal(true, job.All(status.Complete))
	})

	buildUploadTest(func(job *status.Job, _ *MockUploader) {
		job.RetryPolicy = string(RetryErrored)
		RetryTransfers(context.Background(), job)
		assert.Equal(false, job.All(status.Complete))
		assert.Equal(2, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
		assert.Equal(false, job.EndTime.IsZero(), "does change the original time")
	})

	buildUploadTest(func(job *status.Job, _ *MockUploader) {
		job.StartTime = time.Now().AddDate(0, -1, 0)
		startTime := job.StartTime
		job.EndTime = time.Now().AddDate(0, -1, 0)
		endTime := job.EndTime
		job.RetryPolicy = string(RetryErrored)
		uploadStatusErrored := job.Statuses[0].(*UploadStatus)
		uploadStatusErrored.Status = status.Complete
		RetryTransfers(context.Background(), job)
		assert.Equal(false, job.All(status.Complete))
		assert.Equal(2, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
		assert.Equal(startTime, job.StartTime, "does not change the original time")
		assert.Equal(endTime, job.EndTime, "does not change the original time")
	})

	buildUploadTest(func(job *status.Job, _ *MockUploader) {
		job.RetryPolicy = string(RetryErroredIfSomeCompleted)
		uploadStatusErrored := job.Statuses[0].(*UploadStatus)
		uploadStatusComplete := job.Statuses[1].(*UploadStatus)
		uploadStatusComplete.lastByte = uploadStatusErrored.lastByte
		RetryTransfers(context.Background(), job)
		assert.Equal(false, job.All(status.Complete))
		assert.Equal(1, job.Count(status.Complete))
		assert.Equal(1, job.Count(status.Queued))
		assert.Equal(1, job.Count(status.Errored))
		assert.Equal(true, job.EndTime.IsZero(), "does not change the original time")
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
