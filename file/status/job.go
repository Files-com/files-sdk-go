package status

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v2/lib/timer"

	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/chilts/sid"
	ignore "github.com/sabhiram/go-gitignore"

	"github.com/Files-com/files-sdk-go/v2/directory"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"

	filesSDK "github.com/Files-com/files-sdk-go/v2"
)

type EventsReporter map[Status]Reporter

type ToStatusFile interface {
	ToStatusFile() File
	SetStatus(Status, error)
}

type Job struct {
	Id string
	*timer.Timer
	Statuses      []ToStatusFile
	Scanning      bool
	Direction     direction.Type
	statusesMutex *sync.RWMutex
	LocalPath     string
	RemotePath    string
	Sync          bool
	Start         func()
	Canceled      bool
	context.CancelFunc
	Wait   func()
	Params interface{}
	Client interface{}
	EventsReporter
	directory.Type
	*manager.Manager
	RetryPolicy string
	*ignore.GitIgnore
}

func (r Job) Init() *Job {
	r.statusesMutex = &sync.RWMutex{}
	r.Id = sid.IdBase64()
	r.EventsReporter = make(map[Status]Reporter)
	r.Wait = func() {}
	r.Timer = timer.New()
	return &r
}

func (r *Job) SetManager(m *manager.Manager) {
	if m == nil {
		r.Manager = manager.Default()
	} else {
		r.Manager = m
	}
}

func (r *Job) SetEventsReporter(e EventsReporter) {
	if len(e) > 0 {
		r.EventsReporter = e
	}
}

func (r *Job) ClearStatuses() Job {
	newJob := *r
	newJob.Reset()
	newJob.Statuses = []ToStatusFile{}
	return newJob
}

func (r *Job) Reset() {
	r.Canceled = false
	r.Timer = timer.New()
}

func (r *Job) Cancel() {
	r.Canceled = true
	r.Timer.Stop()
	r.CancelFunc()
}

func (r *Job) Job() *Job {
	return r
}

func (r *Job) WithContext(ctx context.Context) context.Context {
	jobCtx, cancel := context.WithCancel(ctx)
	r.CancelFunc = cancel
	return jobCtx
}

func (r *Job) Events(event Status, callback Reporter) {
	r.EventsReporter[event] = callback
}

func (r *Job) UpdateStatus(status Status, file ToStatusFile, err error) {
	if err != nil && strings.Contains(err.Error(), "context canceled") {
		err = nil
		status = Canceled
	}
	if err != nil && errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}
	file.SetStatus(status, err)
	callback, ok := r.EventsReporter[status]
	if ok {
		callback(file.ToStatusFile())
	}
}

func (r *Job) Count(t ...Status) int {
	if len(t) == 0 {
		return len(r.Statuses)
	}
	var total int
	r.statusesMutex.RLock()
	for _, s := range r.Statuses {
		if s.ToStatusFile().Status.Any(t...) {
			total += 1
		}
	}
	r.statusesMutex.RUnlock()
	return total
}

func (r *Job) Add(report ToStatusFile) {
	r.statusesMutex.Lock()
	r.Statuses = append(r.Statuses, report)
	r.statusesMutex.Unlock()
}

func (r *Job) TotalBytes(t ...Status) int64 {
	var total int64
	r.statusesMutex.RLock()
	for _, s := range r.Statuses {
		if s.ToStatusFile().Status.Any(t...) {
			total += s.ToStatusFile().Size
		}
	}
	r.statusesMutex.RUnlock()
	return total
}

func (r *Job) RemainingBytes(t ...Status) int64 {
	return r.TotalBytes(t...) - r.TransferBytes(t...)
}

func (r *Job) TransferBytes(t ...Status) int64 {
	var transferBytes int64
	r.statusesMutex.RLock()
	for _, s := range r.Statuses {
		if s.ToStatusFile().Status.Any(t...) {
			transferBytes += s.ToStatusFile().TransferBytes
		}
	}
	r.statusesMutex.RUnlock()
	return transferBytes
}

func (r *Job) mostRecentBytes(t ...Status) (recent time.Time) {
	r.statusesMutex.RLock()
	for _, s := range r.Statuses {
		if !s.ToStatusFile().Status.Any(t...) {
			continue
		}
		if recent.IsZero() || recent.Before(s.ToStatusFile().LastByte) {
			recent = s.ToStatusFile().LastByte
		}
	}
	r.statusesMutex.RUnlock()
	return
}

func (r *Job) Idle(t ...Status) bool {
	return r.mostRecentBytes(t...).Before(time.Now().Add(time.Duration(-2) * time.Second))
}

func (r *Job) TransferRate(t ...Status) int64 {
	millisecondsSinceStart := time.Now().Sub(r.LastStart()).Milliseconds()
	bytesPerMilliseconds := float64(r.TransferBytes(t...)) / float64(millisecondsSinceStart)
	bytesPerSecond := bytesPerMilliseconds * float64(1000)

	if bytesPerSecond < 0 {
		return 0
	}
	return int64(bytesPerSecond)
}

func (r *Job) ETA(t ...Status) time.Duration {
	if r.TransferRate() == 0 {
		return 0
	}
	seconds := time.Duration(r.RemainingBytes(t...) / r.TransferRate(t...))
	eta := seconds * time.Second
	if eta < 0 {
		return 0
	}
	return eta
}

func (r *Job) ElapsedTime() time.Duration {
	return r.Elapsed()
}

func (r *Job) All(t ...Status) bool {
	allEnded := true
	r.statusesMutex.RLock()
	for _, s := range r.Statuses {
		if !s.ToStatusFile().Status.Any(t...) {
			allEnded = false
			break
		}
	}
	r.statusesMutex.RUnlock()
	return allEnded
}

func (r *Job) Any(t ...Status) (b bool) {
	r.statusesMutex.RLock()
	for _, s := range r.Statuses {
		if s.ToStatusFile().Status.Any(t...) {
			b = true
			break
		}
	}
	r.statusesMutex.RUnlock()
	return
}

func (r *Job) Find(t Status) (ToStatusFile, bool) {
	r.statusesMutex.RLock()
	defer r.statusesMutex.RUnlock()

	for _, s := range r.Statuses {
		if s.ToStatusFile().Status.Any(t) {
			return s, true
		}
	}

	return nil, false
}

func (r *Job) Sub(t ...Status) *Job {
	var sub []ToStatusFile
	r.statusesMutex.RLock()
	for _, s := range r.Statuses {
		if s.ToStatusFile().Status.Any(t...) {
			sub = append(sub, s)
		}
	}
	r.statusesMutex.RUnlock()
	newJob := *r
	newJob.Statuses = sub
	return &newJob
}

func (r *Job) Files() []filesSDK.File {
	var files []filesSDK.File
	r.statusesMutex.RLock()
	for _, s := range r.Statuses {
		files = append(files, s.ToStatusFile().File)
	}
	r.statusesMutex.RUnlock()
	return files
}

func (r *Job) Percentage(t ...Status) int {
	p := int((float64(r.TransferBytes(t...)) / float64(r.TotalBytes(t...))) * float64(100))
	if p < 0 {
		return 0
	}
	return p
}

func (r *Job) StatusFromError(s ToStatusFile, err error) {
	if r.Canceled {
		r.UpdateStatus(Canceled, s, nil)
	} else {
		r.UpdateStatus(Errored, s, err)
	}
}
