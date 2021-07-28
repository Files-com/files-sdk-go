package status

import (
	"sync"

	filesSDK "github.com/Files-com/files-sdk-go"
	"github.com/chilts/sid"
)

type Job struct {
	statuses      []Report
	statusesMutex *sync.Mutex
	ended         int
	Id            string
}

func (r Job) Init() *Job {
	r.statusesMutex = &sync.Mutex{}
	if r.Id == "" {
		r.Id = sid.IdBase64()
	}
	return &r
}

func (r Job) Count() int {
	return len(r.statuses)
}

func (r *Job) Add(report Report) {
	r.statusesMutex.Lock()
	r.statuses = append(r.statuses, report)
	r.statusesMutex.Unlock()
}

func (r Job) TotalBytes() int64 {
	var total int64
	r.statusesMutex.Lock()
	for _, s := range r.statuses {
		if s.Valid() {
			total += s.File().Size
		}
	}
	r.statusesMutex.Unlock()
	return total
}

func (r Job) TransferBytes() int64 {
	var transferBytes int64
	r.statusesMutex.Lock()
	for _, s := range r.statuses {
		if s.Valid() {
			transferBytes += s.TransferBytes()
		}
	}
	r.statusesMutex.Unlock()
	return transferBytes
}

func (r Job) AllEnded() bool {
	allEnded := true
	r.statusesMutex.Lock()
	for _, s := range r.statuses {
		if !s.Ended() {
			allEnded = false
		}
	}
	r.statusesMutex.Unlock()
	return allEnded
}

func (r Job) Files() []filesSDK.File {
	var files []filesSDK.File
	r.statusesMutex.Lock()
	for _, s := range r.statuses {
		files = append(files, s.File())
	}
	r.statusesMutex.Unlock()
	return files
}
