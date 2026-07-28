package file

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StatusFile struct {
	file JobFile
	status.Changes
}

func (f StatusFile) StatusChanges() status.Changes {
	return f.Changes
}

func (f StatusFile) SetStatus(status status.Status, _ error) {
	f.file.Status = status
}

func (f StatusFile) TransferBytes() int64 {
	return f.file.TransferBytes
}

func (f StatusFile) IncrementTransferBytes(int64) {}

func (f StatusFile) File() files_sdk.File {
	return f.file.File
}

func (f StatusFile) LocalPath() string {
	return f.file.LocalPath
}

func (f StatusFile) RemotePath() string {
	return f.file.RemotePath
}

func (f StatusFile) Status() status.Status {
	return f.file.Status
}

func (f StatusFile) Err() error {
	return f.file.Err
}

func (f StatusFile) Job() *Job {
	return f.file.Job
}

func (f StatusFile) Id() string {
	return f.file.Id
}

func (f StatusFile) Size() int64 {
	return f.file.Size
}

func (f StatusFile) EndedAt() time.Time {
	return time.Time{}
}

func (f StatusFile) StartedAt() time.Time {
	return time.Time{}
}

func TestJob_TransferRate(t *testing.T) {
	assert := assert.New(t)
	job := (&Job{Logger: lib.NullLogger{}}).Init()
	job.Timer.Start()
	file := StatusFile{file: JobFile{Status: status.Downloading}}
	job.Add(file)
	job.UpdateStatusWithBytes(status.Uploading, file, 1000)
	assert.InDelta(int64(200), job.TransferRate(), 100)
	time.Sleep(1 * time.Second)
	assert.InDelta(int64(200), job.TransferRate(), 100)
	assert.False(job.Idle(), "Nothing has happened recently so rate is zero")
}

func TestJobDirectTransferClientCacheIsSharedAcrossFiles(t *testing.T) {
	job := (&Job{}).Init()
	jobContext := job.withDirectTransferClientCache(context.Background())
	defer job.Finish()

	firstFileContext, closeFirst := files_sdk.WithDirectTransferClientCache(jobContext)
	defer closeFirst()
	secondFileContext, closeSecond := files_sdk.WithDirectTransferClientCache(jobContext)
	defer closeSecond()

	caPEM, _ := testDirectTransferCertificate(t, testDirectCredentialServerName)
	directInfo := files_sdk.DirectConnectionInfo{
		Version:    1,
		ServerName: testDirectCredentialServerName,
		Addresses:  []string{"tcp://127.0.0.1:4001"},
		DirectUri:  "https://" + testDirectCredentialServerName + "/uploads",
		CaPem:      string(caPEM),
	}
	config := files_sdk.Config{Environment: files_sdk.Development}.Init()
	_, firstClient, err := files_sdk.DirectTransferRetryableClient(firstFileContext, config, directInfo)
	require.NoError(t, err)
	_, secondClient, err := files_sdk.DirectTransferRetryableClient(secondFileContext, config, directInfo)
	require.NoError(t, err)
	require.Same(t, firstClient.HTTPClient, secondClient.HTTPClient)
}

func TestJob_ETA(t *testing.T) {
	assert := assert.New(t)
	job := (&Job{Logger: lib.NullLogger{}}).Init()
	job.Timer.Start()
	file := StatusFile{
		file: JobFile{
			Status: status.Downloading,
			Size:   10000,
		},
	}
	job.Add(file)
	job.UpdateStatusWithBytes(status.Uploading, file, 1000)
	time.Sleep(1 * time.Second)
	assert.InDelta(50000, job.ETA().Milliseconds(), 100)
}

func TestJob_ElapsedTime(t *testing.T) {
	assert := assert.New(t)
	job := (&Job{}).Init()
	job.Start()

	file := StatusFile{
		file: JobFile{
			TransferBytes: 1000,
			Status:        status.Complete,
			Size:          10000,
		},
	}
	file.file.TransferBytes = +5000
	time.Sleep(1 * time.Second)
	file.file.TransferBytes = +5000
	time.Sleep(1 * time.Second)
	file.file.Status = status.Complete
	job.Stop()

	job.Add(file)
	assert.InDelta(2000, job.ElapsedTime().Milliseconds(), 100)
}

func TestJob_TotalBytes(t *testing.T) {
	assert := assert.New(t)
	job := (&Job{}).Init()
	file := StatusFile{}
	file.file.Status = status.Complete
	file.file.Size = 10000
	job.Add(file)
	job.Add(file)
	job.Add(file)
	assert.Equal(int64(30000), job.TotalBytes())
}

func TestJobSetManagerUsesSharedDefault(t *testing.T) {
	first := (&Job{}).Init()
	second := (&Job{}).Init()

	first.SetManager(nil)
	second.SetManager(nil)

	assert.Same(t, first.Manager, second.Manager)
}

func TestJobSetManagerKeepsExplicitManagerIsolated(t *testing.T) {
	job := (&Job{}).Init()
	explicit := manager.New(2, 3, 4)

	job.SetManager(explicit)

	assert.Same(t, explicit, job.Manager)
}

func TestJob_RemainingBytes(t *testing.T) {
	assert := assert.New(t)
	job := (&Job{}).Init()
	file := StatusFile{}
	file.file.Status = status.Complete
	file.file.Size = 10000
	job.Add(file)
	file.file.TransferBytes = +5000
	job.Add(file)
	file.file.TransferBytes = +5000
	job.Add(file)
	assert.Equal(int64(20000), job.RemainingBytes())
}

func TestJob_Count(t *testing.T) {
	assert := assert.New(t)
	job := (&Job{}).Init()
	file := StatusFile{}
	file.file.Status = status.Complete
	file.file.Size = 10000
	job.Add(file)
	job.Add(file)
	file.file.Status = status.Queued
	job.Add(file)
	file.file.Status = status.Errored
	job.Add(file)
	assert.Equal(3, job.Count(status.Ended...))
	assert.Equal(4, job.Count())
}

func TestJob_Sub(t *testing.T) {
	assert := assert.New(t)
	job := (&Job{}).Init()
	file := StatusFile{}
	file.file.Status = status.Complete
	file.file.Size = 10000
	job.Add(file)
	job.Add(file)
	file.file.Status = status.Queued
	job.Add(file)
	file.file.Status = status.Errored
	job.Add(file)
	file.file.Status = status.Ignored
	job.Add(file)
	file.file.Status = status.Skipped
	job.Add(file)
	assert.Equal(3, job.Sub(status.Valid...).Count())
	assert.Equal(2, job.Sub(status.Valid...).Count(status.Ended...))
	assert.Equal(2, job.Sub(status.Excluded...).Count())
	assert.Equal(1, job.Sub(status.Excluded...).Count(status.Skipped))
	assert.Equal(6, job.Count())
	assert.True(job.Any(status.Skipped), "A job should be skipped")
	assert.True(job.Any(status.Ignored), "A job should be ignored")
	assert.False(job.Any(status.Running...), "No jobs should be running")
}

func TestJob_Percentage(t *testing.T) {
	assert := assert.New(t)
	job := (&Job{}).Init()
	file := StatusFile{}
	file.file.Status = status.Complete
	file.file.Size = 100
	file.file.TransferBytes = 100
	job.Add(file)
	job.Add(file)
	file.file.TransferBytes = 0
	file.file.Status = status.Queued
	job.Add(file)
	file.file.TransferBytes = 50
	file.file.Status = status.Errored
	job.Add(file)
	file.file.TransferBytes = 0
	file.file.Status = status.Ignored
	job.Add(file)
	file.file.TransferBytes = 0
	file.file.Status = status.Skipped
	job.Add(file)
	assert.Equal(62, job.Percentage(status.Included...)) // 250 / 400

	job = (&Job{}).Init()
	file = StatusFile{}
	file.file.Status = status.Canceled
	file.file.Size = 100
	file.file.TransferBytes = 1
	job.Add(file)
	assert.Equal(1, job.Percentage(status.Included...))
}

func TestJob_Called(t *testing.T) {
	assert := assert.New(t)
	job := (&Job{}).Init()

	job.Start()

	assert.True(job.Started.Called())
	assert.False(job.Finished.Called())
	job.Finish()
	assert.True(job.Finished.Called())
	job.ClearCalled()
	assert.False(job.Started.Called())
	assert.False(job.Finished.Called())
}

func TestJobWaitsForSettledCancellation(t *testing.T) {
	job := (&Job{Logger: lib.NullLogger{}}).Init()
	jobCtx := job.WithContext(context.Background())
	file := &UploadStatus{
		file:      files_sdk.File{Path: "file.txt"},
		status:    status.Uploading,
		job:       job,
		localPath: "file.txt",
		UploadResumable: UploadResumable{
			FileUploadPart: files_sdk.FileUploadPart{Ref: "pending-ref"},
		},
		Mutex: &sync.RWMutex{},
	}
	job.Add(file)
	job.Start(true)
	job.Scan()

	onComplete := make(chan *UploadStatus)
	WaitTellFinished(job, onComplete, func() {})
	job.EndScan()
	job.Cancel()
	assert.ErrorIs(t, jobCtx.Err(), context.Canceled)

	waitDone := make(chan struct{})
	go func() {
		job.Wait()
		close(waitDone)
	}()

	select {
	case <-waitDone:
		t.Fatal("Wait returned before completion accounting settled")
	case <-time.After(50 * time.Millisecond):
	}

	checkpoint := job.UploadCheckpoint()
	assert.Empty(t, checkpoint.CompletedPaths)
	assert.Contains(t, checkpoint.PendingParts, "file.txt")

	job.UpdateStatus(status.Complete, file, nil)
	select {
	case onComplete <- file:
	case <-time.After(time.Second):
		t.Fatal("completion accounting did not accept the final file")
	}
	select {
	case <-waitDone:
	case <-time.After(time.Second):
		t.Fatal("Wait did not return after completion accounting settled")
	}

	assert.True(t, job.Finished.Called())
	assert.Equal(t, 1, job.Count(status.Complete))
	checkpoint = job.UploadCheckpoint()
	assert.Contains(t, checkpoint.CompletedPaths, "file.txt")
	assert.NotContains(t, checkpoint.PendingParts, "file.txt")
}

func TestJobCancelBeforeStart(t *testing.T) {
	job := (&Job{}).Init()

	require.NotPanics(t, job.Cancel)
	jobCtx := job.WithContext(context.Background())
	assert.ErrorIs(t, jobCtx.Err(), context.Canceled)

	waitDone := make(chan struct{})
	go func() {
		job.Wait()
		close(waitDone)
	}()
	select {
	case <-waitDone:
	case <-time.After(time.Second):
		t.Fatal("Wait blocked for a canceled job that never started")
	}

	assert.True(t, job.Canceled.Called())
	assert.False(t, job.Started.Called())
	assert.False(t, job.Finished.Called())
}

func TestJobWaitThenCancelBeforeStart(t *testing.T) {
	job := (&Job{}).Init()
	waitDone := make(chan struct{})
	go func() {
		job.Wait()
		close(waitDone)
	}()

	select {
	case <-waitDone:
		t.Fatal("Wait returned before cancellation")
	case <-time.After(50 * time.Millisecond):
	}
	job.Cancel()

	select {
	case <-waitDone:
	case <-time.After(time.Second):
		t.Fatal("Wait blocked after an unstarted job was canceled")
	}
}

func TestJobStatusFromPausedUploadIsCanceled(t *testing.T) {
	job := (&Job{Logger: lib.NullLogger{}}).Init()
	file := &UploadStatus{
		job:    job,
		status: status.Uploading,
		Mutex:  &sync.RWMutex{},
	}

	job.StatusFromError(file, fmt.Errorf("upload interrupted: %w", ErrJobPaused))

	assert.Equal(t, status.Canceled, file.Status())
	assert.NoError(t, file.Err())
}
