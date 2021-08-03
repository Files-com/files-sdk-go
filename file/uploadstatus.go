package file

import (
	"context"

	"github.com/Files-com/files-sdk-go/file/status"

	files_sdk "github.com/Files-com/files-sdk-go"
)

type UploadStatus struct {
	files_sdk.File
	status.Status
	context.CancelFunc
	Job           *status.Job
	LocalPath     string
	RemotePath    string
	UploadedBytes int64
	Sync          bool
}

func (u UploadStatus) ToStatusFile() status.File {
	return status.File{
		TransferBytes: u.UploadedBytes,
		File:          u.File,
		Cancel:        u.CancelFunc,
		LocalPath:     u.LocalPath,
		RemotePath:    u.RemotePath,
		Id:            u.Id(),
		Job:           u.Job,
		Status:        u.Status,
		Sync:          u.Sync,
	}
}

func (u *UploadStatus) SetStatus(s status.Status) {
	if s.Valid() && u.Invalid() {
		return
	}

	u.Status = s
}

func (u UploadStatus) Id() string {
	return u.Job.Id + ":" + u.File.Path
}

func (u *UploadStatus) incrementDownloadedBytes(b int64) {
	u.UploadedBytes += b
}
