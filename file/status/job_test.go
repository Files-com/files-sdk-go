package status

import (
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"github.com/stretchr/testify/assert"
)

type StatusFile struct {
	file File
}

func (f StatusFile) SetStatus(status Status, _ error) {
	f.file.Status = status
}

func (f StatusFile) TransferBytes() int64 {
	return f.file.TransferBytes
}

func (f StatusFile) File() files_sdk.File {
	return f.file.File
}

func (f StatusFile) LocalPath() string {
	return f.file.LocalPath
}

func (f StatusFile) RemotePath() string {
	return f.file.RemotePath
}

func (f StatusFile) Status() Status {
	return f.file.Status
}

func (f StatusFile) LastByte() time.Time {
	return f.file.LastByte
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

func TestJob_TransferRate(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	job.Timer.Start()
	file := StatusFile{file: File{LastByte: time.Now(), TransferBytes: 1000}}
	job.Add(file)
	time.Sleep(1 * time.Second)
	assert.InDelta(int64(1000), job.TransferRate(), 100)
	assert.Equal(false, job.Idle(), "Nothing has happened recently so rate is zero")

	time.Sleep(3500 * time.Millisecond)
	assert.Equal(true, job.Idle(), "Nothing has happened recently so rate is zero")
}

func TestJob_ETA(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	job.Timer.Start()
	file := StatusFile{
		file: File{
			TransferBytes: 1000,
			LastByte:      time.Now(),
			Status:        Downloading,
			File:          files_sdk.File{Size: 10000},
		},
	}
	job.Add(file)

	time.Sleep(1 * time.Second)
	assert.InDelta(9000, job.ETA().Milliseconds(), 100)
}

func TestJob_ElapsedTime(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	job.Timer.Start()

	file := StatusFile{
		file: File{
			TransferBytes: 1000,
			LastByte:      time.Now(),
			Status:        Complete,
			File:          files_sdk.File{Size: 10000},
		},
	}
	file.file.TransferBytes = +5000
	time.Sleep(1 * time.Second)
	file.file.TransferBytes = +5000
	time.Sleep(1 * time.Second)
	file.file.Status = Complete
	job.Timer.Stop()

	job.Add(file)
	assert.InDelta(2000, job.ElapsedTime().Milliseconds(), 100)
}

func TestJob_TotalBytes(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	file := StatusFile{}
	file.file.Status = Complete
	file.file.Size = 10000
	job.Add(file)
	job.Add(file)
	job.Add(file)
	assert.Equal(int64(30000), job.TotalBytes())
}

func TestJob_RemainingBytes(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	file := StatusFile{}
	file.file.Status = Complete
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
	job := Job{}.Init()
	file := StatusFile{}
	file.file.Status = Complete
	file.file.Size = 10000
	job.Add(file)
	job.Add(file)
	file.file.Status = Queued
	job.Add(file)
	file.file.Status = Errored
	job.Add(file)
	assert.Equal(3, job.Count(Ended...))
	assert.Equal(4, job.Count())
}

func TestJob_Sub(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	file := StatusFile{}
	file.file.Status = Complete
	file.file.Size = 10000
	job.Add(file)
	job.Add(file)
	file.file.Status = Queued
	job.Add(file)
	file.file.Status = Errored
	job.Add(file)
	file.file.Status = Ignored
	job.Add(file)
	file.file.Status = Skipped
	job.Add(file)
	assert.Equal(3, job.Sub(Valid...).Count())
	assert.Equal(2, job.Sub(Valid...).Count(Ended...))
	assert.Equal(2, job.Sub(Excluded...).Count())
	assert.Equal(1, job.Sub(Excluded...).Count(Skipped))
	assert.Equal(6, job.Count())
	assert.Equal(true, job.Any(Skipped))
	assert.Equal(true, job.Any(Ignored))
	assert.Equal(false, job.Any(Running...))
}

func TestJob_Percentage(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	file := StatusFile{}
	file.file.Status = Complete
	file.file.Size = 100
	file.file.TransferBytes = 100
	job.Add(file)
	job.Add(file)
	file.file.TransferBytes = 0
	file.file.Status = Queued
	job.Add(file)
	file.file.TransferBytes = 50
	file.file.Status = Errored
	job.Add(file)
	file.file.TransferBytes = 0
	file.file.Status = Ignored
	job.Add(file)
	file.file.TransferBytes = 0
	file.file.Status = Skipped
	job.Add(file)
	assert.Equal(62, job.Percentage(Included...)) // 250 / 400

	job = Job{}.Init()
	file = StatusFile{}
	file.file.Status = Canceled
	file.file.Size = 100
	file.file.TransferBytes = 1
	job.Add(file)
	assert.Equal(1, job.Percentage(Included...))
}

func TestJob_Called(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()

	job.Start()

	assert.Equal(true, job.Started.Called())
	assert.Equal(false, job.Finished.Called())
	job.Finish()
	assert.Equal(true, job.Finished.Called())
	job.ClearCalled()
	assert.Equal(false, job.Started.Called())
	assert.Equal(false, job.Finished.Called())
}
