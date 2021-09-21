package file

import (
	"time"

	"github.com/Files-com/files-sdk-go/v2/file/status"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
)

type UploadStatus struct {
	files_sdk.File
	status.Status
	Job           *status.Job
	LocalPath     string
	RemotePath    string
	UploadedBytes int64
	Sync          bool
	lastByte      time.Time
	error
	Uploader
}

func (u *UploadStatus) ToStatusFile() status.File {
	return status.File{
		TransferBytes: u.UploadedBytes,
		File:          u.File,
		LocalPath:     u.LocalPath,
		RemotePath:    u.RemotePath,
		Id:            u.Id(),
		Job:           u.Job,
		Status:        u.Status,
		LastByte:      u.lastByte,
		Err:           u.error,
	}
}

func (u *UploadStatus) SetStatus(s status.Status, err error) {
	var setError bool
	u.Status, setError = status.SetStatus(u.Status, s, err)
	if setError {
		u.error = err
	}

	if s.Is(status.Retrying) {
		u.UploadedBytes = 0
		u.lastByte = time.Time{}
	}
}

func (u UploadStatus) Id() string {
	return u.Job.Id + ":" + u.File.Path
}

func (u *UploadStatus) incrementDownloadedBytes(b int64) {
	u.lastByte = time.Now()
	u.UploadedBytes += b
}
