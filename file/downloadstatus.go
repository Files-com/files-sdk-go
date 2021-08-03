package file

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/file/status"
)

type DownloadStatus struct {
	files_sdk.File
	status.Status
	DownloadedBytes int64
	LocalPath       string
	RemotePath      string
	Sync            bool
	context.CancelFunc
	*status.Job
}

func (u DownloadStatus) ToStatusFile() status.File {
	return status.File{
		TransferBytes: u.DownloadedBytes,
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

func (r *DownloadStatus) SetStatus(s status.Status) {
	if s.Valid() && r.Invalid() {
		return
	}

	r.Status = s
}

func (r DownloadStatus) Id() string {
	return r.Job.Id + ":" + r.File.Path
}

func (r *DownloadStatus) incrementDownloadedBytes(b int64) {
	r.DownloadedBytes += b
}
