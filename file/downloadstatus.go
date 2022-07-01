package file

import (
	"io/fs"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"github.com/Files-com/files-sdk-go/v2/file/status"
)

type DownloadStatus struct {
	fsFile          fs.File
	file            files_sdk.File
	status          status.Status
	job             *status.Job
	DownloadedBytes int64
	localPath       string
	remotePath      string
	Sync            bool
	lastByte        time.Time
	Mutex           *sync.RWMutex
	PreserveTimes   bool
	error
	lastError error
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

	if s.Is(status.Retrying) {
		d.DownloadedBytes = 0
		d.lastByte = time.Time{}
	}
}

func (d *DownloadStatus) Id() string {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	return d.job.Id + ":" + d.file.Path
}

func (d *DownloadStatus) incrementDownloadedBytes(b int64) {
	d.Mutex.Lock()
	d.lastByte = time.Now()
	d.DownloadedBytes += b
	d.Mutex.Unlock()
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

func (d *DownloadStatus) LastByte() time.Time {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	return d.lastByte
}

func (d *DownloadStatus) Err() error {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	return d.error
}

func (d *DownloadStatus) Job() *status.Job {
	return d.job
}

func (d *DownloadStatus) SetFinalSize(written int64) {
	d.Mutex.Lock()
	d.DownloadedBytes = written
	d.file.Size = written
	d.Mutex.Unlock()
}
