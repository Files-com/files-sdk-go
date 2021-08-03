package status

import (
	"context"
	"sync"

	"github.com/chilts/sid"

	filesSDK "github.com/Files-com/files-sdk-go"
)

type ToStatusFile interface {
	ToStatusFile() File
}

type Job struct {
	statuses      []ToStatusFile
	statusesMutex *sync.RWMutex
	Id            string
	context.CancelFunc
}

func (r Job) Init(id ...string) *Job {
	r.statusesMutex = &sync.RWMutex{}
	if len(id) == 1 {
		r.Id = id[0]
	}
	if r.Id == "" {
		r.Id = sid.IdBase64()
	}
	return &r
}

func (r Job) Cancel() {
	r.CancelFunc()
}

func (r Job) Job() Job {
	return r
}

func (r Job) Count() int {
	return len(r.statuses)
}

func (r *Job) Add(report ToStatusFile) {
	r.statusesMutex.Lock()
	r.statuses = append(r.statuses, report)
	r.statusesMutex.Unlock()
}

func (r Job) TotalBytes() int64 {
	var total int64
	r.statusesMutex.RLock()
	for _, s := range r.statuses {
		if s.ToStatusFile().Valid() {
			total += s.ToStatusFile().Size
		}
	}
	r.statusesMutex.RUnlock()
	return total
}

func (r Job) TransferBytes() int64 {
	var transferBytes int64
	r.statusesMutex.RLock()
	for _, s := range r.statuses {
		if s.ToStatusFile().Valid() {
			transferBytes += s.ToStatusFile().TransferBytes
		}
	}
	r.statusesMutex.RUnlock()
	return transferBytes
}

func (r Job) AllEnded() bool {
	allEnded := true
	r.statusesMutex.RLock()
	for _, s := range r.statuses {
		if !s.ToStatusFile().Ended() {
			allEnded = false
			break
		}
	}
	r.statusesMutex.RUnlock()
	return allEnded
}

func (r Job) Files() []filesSDK.File {
	var files []filesSDK.File
	r.statusesMutex.RLock()
	for _, s := range r.statuses {
		files = append(files, s.ToStatusFile().File)
	}
	r.statusesMutex.RUnlock()
	return files
}
