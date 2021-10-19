package file

import (
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v2/file/status"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
)

type UploadStatus struct {
	file          files_sdk.File
	status        status.Status
	job           *status.Job
	localPath     string
	remotePath    string
	uploadedBytes int64
	Sync          bool
	lastByte      time.Time
	error
	Uploader
	Parts
	files_sdk.FileUploadPart
	Mutex *sync.RWMutex
}

func (u *UploadStatus) Job() *status.Job {
	return u.job
}

func (u *UploadStatus) TransferBytes() int64 {
	u.Mutex.RLock()
	defer u.Mutex.RUnlock()
	return u.uploadedBytes
}

func (u *UploadStatus) File() files_sdk.File {
	return u.file
}

func (u *UploadStatus) LocalPath() string {
	return u.localPath
}

func (u *UploadStatus) RemotePath() string {
	return u.remotePath
}

func (u *UploadStatus) Status() status.Status {
	u.Mutex.RLock()
	defer u.Mutex.RUnlock()
	return u.status
}

func (u *UploadStatus) LastByte() time.Time {
	u.Mutex.RLock()
	defer u.Mutex.RUnlock()
	return u.lastByte
}

func (u *UploadStatus) Err() error {
	u.Mutex.RLock()
	defer u.Mutex.RUnlock()
	return u.error
}

func (u *UploadStatus) SetStatus(s status.Status, err error) {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	var setError bool
	u.status, setError = status.SetStatus(u.status, s, err)
	if setError {
		u.error = err
	}

	if s.Is(status.Retrying) {
		u.uploadedBytes = 0
		u.lastByte = time.Time{}
	}
}

func (u *UploadStatus) Id() string {
	return u.job.Id + ":" + u.file.Path
}

func (u *UploadStatus) incrementDownloadedBytes(b int64) {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	u.lastByte = time.Now()
	u.uploadedBytes += b
}

func (u *UploadStatus) SetUploadedBytes(b int64) {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	u.uploadedBytes = b
}
