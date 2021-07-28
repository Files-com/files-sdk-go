package file

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/file/status"
)

type DownloadStatus struct {
	file files_sdk.File
	status.Status
	DownloadedBytes int64
	destination     string
	cancel          context.CancelFunc
	runStats        *status.Job
}

func (r *DownloadStatus) SetStatus(s status.Status) {
	if s.Valid() && r.Invalid() {
		return
	}

	r.Status = s
}

func (r DownloadStatus) Destination() string {
	return r.destination
}

func (r DownloadStatus) TransferBytes() int64 {
	return r.DownloadedBytes
}

func (r DownloadStatus) File() files_sdk.File {
	return r.file
}

func (r DownloadStatus) Job() status.Job {
	return *r.runStats
}

func (r DownloadStatus) Id() string {
	return r.Job().Id + ":" + r.File().Path
}

func (r *DownloadStatus) incrementDownloadedBytes(b int64) {
	r.DownloadedBytes += b
}

func (r DownloadStatus) Cancel() {
	if r.cancel != nil {
		r.cancel()
		r.Status = status.Canceled
	}
}
