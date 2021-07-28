package file

import (
	"context"

	"github.com/Files-com/files-sdk-go/file/status"

	files_sdk "github.com/Files-com/files-sdk-go"
)

type UploadStatus struct {
	file files_sdk.File
	status.Status
	destination   string
	Source        string
	cancel        context.CancelFunc
	UploadedBytes int64
	job           *status.Job
}

func (u *UploadStatus) SetStatus(s status.Status) {
	if s.Valid() && u.Invalid() {
		return
	}

	u.Status = s
}

func (u UploadStatus) Cancel() {
	if u.cancel != nil {
		u.cancel()
		u.Status = status.Canceled
	}
}

func (u UploadStatus) TransferBytes() int64 {
	return u.UploadedBytes
}

func (u UploadStatus) File() files_sdk.File {
	return u.file
}

func (u UploadStatus) Destination() string {
	return u.destination
}

func (u UploadStatus) Job() status.Job {
	return *u.job
}

func (u UploadStatus) Id() string {
	return u.Job().Id + ":" + u.File().Path
}

func (u *UploadStatus) incrementDownloadedBytes(b int64) {
	u.UploadedBytes += b
}
