package status

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type StatusFile struct {
	File
}

func (f StatusFile) ToStatusFile() File {
	return f.File
}

func (f StatusFile) SetStatus(status Status, _ error) {
	f.File.Status = status
}

func TestJob_TransferRate(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	job.Timer.Start()
	file := StatusFile{}
	file.LastByte = time.Now()
	file.TransferBytes = 1000
	job.Add(file)
	time.Sleep(1 * time.Second)
	assert.InDelta(int64(1000), job.TransferRate(), 100)
	assert.Equal(false, job.Idle(), "Nothing has happened recently so rate is zero")

	time.Sleep(2 * time.Second)
	assert.Equal(true, job.Idle(), "Nothing has happened recently so rate is zero")
}

func TestJob_ETA(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	job.Timer.Start()
	file := StatusFile{}
	file.Status = Downloading
	file.Size = 10000
	file.LastByte = time.Now()
	file.TransferBytes = +1000
	job.Add(file)

	time.Sleep(1 * time.Second)
	assert.InDelta(9000, job.ETA().Milliseconds(), 100)
}

func TestJob_ElapsedTime(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	job.Timer.Start()

	file := StatusFile{}
	file.Status = Complete
	file.Size = 10000
	file.LastByte = time.Now()
	file.TransferBytes = +5000
	time.Sleep(1 * time.Second)
	file.TransferBytes = +5000
	time.Sleep(1 * time.Second)
	file.Status = Complete
	job.Timer.Stop()

	job.Add(file)
	assert.InDelta(2000, job.ElapsedTime().Milliseconds(), 100)
}

func TestJob_TotalBytes(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	file := StatusFile{}
	file.Status = Complete
	file.Size = 10000
	job.Add(file)
	job.Add(file)
	job.Add(file)
	assert.Equal(int64(30000), job.TotalBytes())
}

func TestJob_RemainingBytes(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	file := StatusFile{}
	file.Status = Complete
	file.Size = 10000
	job.Add(file)
	file.TransferBytes = +5000
	job.Add(file)
	file.TransferBytes = +5000
	job.Add(file)
	assert.Equal(int64(20000), job.RemainingBytes())
}

func TestJob_Count(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	file := StatusFile{}
	file.Status = Complete
	file.Size = 10000
	job.Add(file)
	job.Add(file)
	file.Status = Queued
	job.Add(file)
	file.Status = Errored
	job.Add(file)
	assert.Equal(3, job.Count(Ended...))
	assert.Equal(4, job.Count())
}

func TestJob_Sub(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()
	file := StatusFile{}
	file.Status = Complete
	file.Size = 10000
	job.Add(file)
	job.Add(file)
	file.Status = Queued
	job.Add(file)
	file.Status = Errored
	job.Add(file)
	file.Status = Ignored
	job.Add(file)
	file.Status = Skipped
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
	file.Status = Complete
	file.Size = 100
	file.TransferBytes = 100
	job.Add(file)
	job.Add(file)
	file.TransferBytes = 0
	file.Status = Queued
	job.Add(file)
	file.TransferBytes = 50
	file.Status = Errored
	job.Add(file)
	file.TransferBytes = 0
	file.Status = Ignored
	job.Add(file)
	file.TransferBytes = 0
	file.Status = Skipped
	job.Add(file)
	assert.Equal(62, job.Percentage(Included...)) // 250 / 400

	job = Job{}.Init()
	file = StatusFile{}
	file.Status = Canceled
	file.Size = 100
	file.TransferBytes = 1
	job.Add(file)
	assert.Equal(1, job.Percentage(Included...))
}

func TestJob_Called(t *testing.T) {
	assert := assert.New(t)
	job := Job{}.Init()

	job.Start()

	assert.Equal(true, job.Started.Called)
	assert.Equal(false, job.Finished.Called)
	job.Finish()
	assert.Equal(true, job.Finished.Called)
	job.ClearCalled()
	assert.Equal(false, job.Started.Called)
	assert.Equal(false, job.Finished.Called)
}
