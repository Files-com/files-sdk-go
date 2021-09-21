package file

import (
	"io/fs"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"github.com/Files-com/files-sdk-go/v2/file/status"
)

type DownloadStatus struct {
	fsFile fs.File
	files_sdk.File
	status.Status
	*status.Job
	DownloadedBytes int64
	LocalPath       string
	RemotePath      string
	Sync            bool
	lastByte        time.Time
	error
}

func (d DownloadStatus) ToStatusFile() status.File {
	return status.File{
		TransferBytes: d.DownloadedBytes,
		File:          d.File,
		LocalPath:     d.LocalPath,
		RemotePath:    d.RemotePath,
		Id:            d.Id(),
		Job:           d.Job,
		Status:        d.Status,
		LastByte:      d.lastByte,
		Err:           d.error,
	}
}

func (d *DownloadStatus) SetStatus(s status.Status, err error) {
	var setError bool
	d.Status, setError = status.SetStatus(d.Status, s, err)
	if setError {
		d.error = err
	}

	if s.Is(status.Retrying) {
		d.DownloadedBytes = 0
		d.lastByte = time.Time{}
	}
}

func (d DownloadStatus) Id() string {
	return d.Job.Id + ":" + d.File.Path
}

func (d *DownloadStatus) incrementDownloadedBytes(b int64) {
	d.lastByte = time.Now()
	d.DownloadedBytes += b
}
