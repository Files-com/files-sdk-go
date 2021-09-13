package file

import (
	"context"
	"time"

	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"
	"github.com/bradfitz/iter"
)

type RetryPolicy string

const (
	RetryAll                    = RetryPolicy("RetryAll")
	RetryErrored                = RetryPolicy("RetryErrored")
	RetryErroredIfSomeCompleted = RetryPolicy("RetryErroredIfSomeCompleted")
)

func RetryTransfers(ctx context.Context, job *status.Job) {
	if job.Stopped {
		return
	}
	switch RetryPolicy(job.RetryPolicy) {
	case RetryAll:
		retryAll(ctx, job)
	case RetryErrored:
		retryErrored(ctx, job)
	case RetryErroredIfSomeCompleted:
		retryErroredIfSomeCompleted(ctx, job)
	}
}

func retryErroredIfSomeCompleted(ctx context.Context, job *status.Job) {
	lastFailure := time.Time{}
	for _, s := range job.Sub(status.Errored).Statuses {
		if lastFailure.Before(s.ToStatusFile().LastByte) {
			lastFailure = s.ToStatusFile().LastByte
		}
	}
	if lastFailure.IsZero() {
		return
	}
	for _, s := range job.Sub(status.Complete).Statuses {
		if lastFailure.Before(s.ToStatusFile().LastByte) {
			retryErrored(ctx, job)
			return
		}
	}
}

func retryAll(ctx context.Context, job *status.Job) {
	retryByStatus(ctx, job, status.Included...)
}

func retryErrored(ctx context.Context, job *status.Job) {
	retryByStatus(ctx, job, status.Errored)
}

func retryByStatus(ctx context.Context, job *status.Job, s ...status.Status) {
	switch job.Direction {
	case direction.DownloadType:
		onComplete := make(chan *DownloadStatus)
		enqueueByStatus(job,
			func(s status.ToStatusFile) {
				job.UpdateStatus(status.Retrying, s.(*DownloadStatus), nil)
				enqueueDownload(ctx, job, s.(*DownloadStatus), onComplete)
			}, func() {
				<-onComplete
			},
			s...,
		)
		close(onComplete)
	case direction.UploadType:
		onComplete := make(chan *UploadStatus)
		enqueueByStatus(job,
			func(s status.ToStatusFile) {
				job.UpdateStatus(status.Retrying, s.(*UploadStatus), nil)
				enqueueUpload(ctx, job, s.(*UploadStatus), onComplete)
			}, func() {
				<-onComplete
			},
			s...,
		)
		close(onComplete)
	default:
		panic("invalid direction")
	}
}

func enqueueByStatus(job *status.Job, enqueue func(s status.ToStatusFile), onComplete func(), s ...status.Status) {
	if job.Count(s...) == 0 {
		return
	}

	job.Reset()
	count := 0
	for _, s := range job.Sub(s...).Statuses {
		count += 1
		enqueue(s)
	}
	for range iter.N(count) {
		onComplete()
	}
	job.EndTime = time.Now()
}
