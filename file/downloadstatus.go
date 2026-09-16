package file

import (
	"io/fs"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/status"
)

type DownloadStatus struct {
	fsFile fs.File
	fs.FS
	fs.FileInfo
	file            files_sdk.File
	status          status.Status
	job             *Job
	DownloadedBytes int64
	localPath       string
	destination     destinationPath
	remotePath      string
	tempPath        string
	TmpPath         string
	Sync            bool
	NoOverwrite     bool
	endedAt         time.Time
	startedAt       time.Time
	Mutex           *sync.RWMutex
	PreserveTimes   bool
	error
	lastError error
	// indexErr records why the file can never be downloaded: its listing or
	// open failed, or its server path does not resolve inside the destination.
	// Retrying clears error, so this is checked separately to keep the file
	// Errored on every retry instead of downloading it.
	indexErr error
	dryRun   bool
	status.Changes
}

var _ IFile = &DownloadStatus{}

func (d *DownloadStatus) EndedAt() time.Time {
	return d.endedAt
}

func (d *DownloadStatus) StartedAt() time.Time {
	return d.startedAt
}

func (d *DownloadStatus) Size() (size int64) {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	if d.FileInfo != nil {
		size = d.FileInfo.Size()
	}

	if size <= 0 {
		size = d.file.Size
	}

	return
}

func (d *DownloadStatus) RecentError() error {
	if d.error != nil {
		return d.error
	}

	return d.lastError
}

func (d *DownloadStatus) SetStatus(s status.Status, err error) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	var setError bool
	d.status, setError = status.SetStatus(d.status, s, err)
	if setError {
		if d.error != nil {
			d.lastError = d.error
		}
		d.error = err
	}

	if s.Is(status.Downloading) && d.startedAt.IsZero() {
		d.startedAt = time.Now()
	}

	if s.Is(status.Retrying) {
		d.DownloadedBytes = 0
	}

	if s.Is(status.Ended...) {
		d.endedAt = time.Now()
	}

	d.Changes = append(d.Changes, status.Change{Status: d.status, Err: d.error, Time: time.Now()})
}

func (d *DownloadStatus) StatusChanges() status.Changes {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()

	return d.Changes
}

func (d *DownloadStatus) Id() string {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	return d.job.Id + ":" + d.file.Path
}

func (d *DownloadStatus) incrementDownloadedBytes(b int64) {
	d.Mutex.Lock()
	d.DownloadedBytes += b
	d.Mutex.Unlock()
}

func (d *DownloadStatus) IncrementTransferBytes(b int64) {
	d.incrementDownloadedBytes(b)
}

func (d *DownloadStatus) TransferBytes() int64 {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	return d.DownloadedBytes
}

func (d *DownloadStatus) File() files_sdk.File {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	return d.file
}

func (d *DownloadStatus) LocalPath() string {
	return d.localPath
}

func (d *DownloadStatus) RemotePath() string {
	return d.remotePath
}

func (d *DownloadStatus) Status() status.Status {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	return d.status
}

func (d *DownloadStatus) Err() error {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	return d.error
}

func (d *DownloadStatus) Job() *Job {
	return d.job
}

// SetFinalSize records the completed output size. The listing metadata the
// status was indexed from can be stale, so the completed output defines the
// size from here on, including a completed empty file.
func (d *DownloadStatus) SetFinalSize(written int64) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	d.DownloadedBytes = written
	d.file.Size = written
	d.FileInfo = Info{File: d.file, sizeTrust: TrustedSizeValue}
}
