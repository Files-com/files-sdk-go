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
	RetryUnfinished             = RetryPolicy("RetryUnfinished")
	RetryErroredIfSomeCompleted = RetryPolicy("RetryErroredIfSomeCompleted")
)

func RetryByPolicy(ctx context.Context, job *status.Job, policy RetryPolicy, signalEvents bool) {
	switch policy {
	case RetryAll:
		RetryByStatus(ctx, job, signalEvents, status.Included...)
	case RetryUnfinished:
		RetryByStatus(ctx, job, signalEvents, append(status.Running, []status.Status{status.Errored, status.Canceled}...)...)
	case RetryErroredIfSomeCompleted:
		retryErroredIfSomeCompleted(ctx, job, signalEvents)
	}
}

func RetryByStatus(ctx context.Context, job *status.Job, signalEvents bool, s ...status.Status) {
	switch job.Direction {
	case direction.DownloadType:
		onComplete := make(chan *DownloadStatus)
		defer close(onComplete)
		enqueueByStatus(ctx, job, signalEvents,
			func(s status.IFile, jobCxt context.Context) {
				job.UpdateStatus(status.Retrying, s.(*DownloadStatus), nil)
				enqueueDownload(jobCxt, job, s.(*DownloadStatus), onComplete)
			}, func() {
				<-onComplete
			},
			s...,
		)
	case direction.UploadType:
		onComplete := make(chan *UploadStatus)
		defer close(onComplete)
		enqueueByStatus(ctx, job, signalEvents,
			func(s status.IFile, jobCxt context.Context) {
				job.UpdateStatus(status.Retrying, s.(*UploadStatus), nil)
				enqueueUpload(jobCxt, job, s.(*UploadStatus), onComplete)
			}, func() {
				<-onComplete
			},
			s...,
		)
	default:
		panic("invalid direction")
	}
}

func retryErroredIfSomeCompleted(ctx context.Context, job *status.Job, signalEvents bool) {
	lastFailure := time.Time{}
	for _, s := range job.Sub(status.Errored).Statuses {
		if lastFailure.Before(s.LastByte()) {
			lastFailure = s.LastByte()
		}
	}
	if lastFailure.IsZero() {
		return
	}
	for _, s := range job.Sub(status.Complete).Statuses {
		if lastFailure.Before(s.LastByte()) {
			RetryByPolicy(ctx, job, RetryUnfinished, signalEvents)
			return
		}
	}
}

func enqueueByStatus(ctx context.Context, job *status.Job, signalEvents bool, enqueue func(status.IFile, context.Context), onComplete func(), s ...status.Status) {
	if job.Count(s...) == 0 {
		return
	}
	jobCtx := job.WithContext(ctx)

	count := 0

	if signalEvents {
		job.ClearCalled()
		job.Start(false)
		job.Scan()
	}

	for _, s := range job.Sub(s...).Statuses {
		count += 1
		enqueue(s, jobCtx)
	}
	if signalEvents {
		job.EndScan()
	}
	for range iter.N(count) {
		onComplete()
	}
	if signalEvents {
		job.Finish()
	}
}
