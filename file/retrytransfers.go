package file

import (
	"context"
	"strings"
	"time"

	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"
	"github.com/bradfitz/iter"
)

type RetryPolicy struct {
	Type       RetryPolicyType
	RetryCount int
	Backoff    int
}

func (p RetryPolicy) WaitSec(retry int) time.Duration {
	if p.Backoff == 0 {
		p.Backoff = 2
	}
	return time.Second * time.Duration(p.Backoff*retry)
}

type RetryPolicyType string

const (
	RetryAll                    = RetryPolicyType("RetryAll")
	RetryUnfinished             = RetryPolicyType("RetryUnfinished")
	RetryErroredIfSomeCompleted = RetryPolicyType("RetryErroredIfSomeCompleted")
)

func RetryByPolicy(ctx context.Context, job *status.Job, policy RetryPolicy, signalEvents bool) {
	switch policy.Type {
	case RetryAll:
		RetryByStatus(ctx, job, signalEvents, policy, status.Included...)
	case RetryUnfinished:
		RetryByStatus(ctx, job, signalEvents, policy, append(status.Running, []status.Status{status.Errored, status.Canceled}...)...)
	case RetryErroredIfSomeCompleted:
		retryErroredIfSomeCompleted(ctx, job, signalEvents, policy)
	}
}

func RetryByStatus(ctx context.Context, job *status.Job, signalEvents bool, policy RetryPolicy, s ...status.Status) {
	for i := range iter.N(policy.RetryCount) {
		switch job.Direction {
		case direction.DownloadType:
			retryDownload(ctx, job, signalEvents, s)
		case direction.UploadType:
			retryUpload(ctx, job, signalEvents, s)
		default:
			panic("invalid direction")
		}
		if len(job.Sub(s...).Statuses) > 0 {
			job.Logger.Printf("retry (%v): backing off %v sec", i, policy.WaitSec(i))
			time.Sleep(policy.WaitSec(i))
		} else {
			return
		}
	}
}

func retryUpload(ctx context.Context, job *status.Job, signalEvents bool, s []status.Status) {
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
}

func retryDownload(ctx context.Context, job *status.Job, signalEvents bool, s []status.Status) {
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
}

func retryErroredIfSomeCompleted(ctx context.Context, job *status.Job, signalEvents bool, policy RetryPolicy) {
	lastFailure := time.Time{}
	for _, s := range job.Sub(status.Errored).Statuses {
		if lastFailure.Before(s.LastByte()) {
			lastFailure = s.LastByte()
		}
	}
	if lastFailure.IsZero() {
		return
	}
	RetryByPolicy(ctx, job, RetryPolicy{Type: RetryUnfinished, RetryCount: policy.RetryCount}, signalEvents)
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

	files := job.Sub(s...).Statuses

	var types []string
	for _, st := range s {
		types = append(types, st.String())
	}
	job.Logger.Printf("retrying %v files (%v)", strings.Join(types, ", "), len(files))

	for _, file := range files {
		count += 1
		go enqueue(file, jobCtx)
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
