package fsmount

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Files-com/files-sdk-go/v3/fsmount/events"
)

const (
	transferProgressMinBytes    int64 = 128 * 1024
	transferProgressMinInterval       = 500 * time.Millisecond
	transferredBytesUnchanged   int64 = -1
)

type transferReporter struct {
	publisher events.EventPublisher

	id         string
	direction  events.TransferDirection
	localPath  string
	remotePath string
	size       int64
	startedAt  time.Time

	transferredBytes  int64
	lastProgressAt    time.Time
	lastProgressBytes int64
	progressEmitted   bool
	terminalEmitted   bool

	mu sync.Mutex
}

func (fs *RemoteFs) newTransferReporter(direction events.TransferDirection, path string, size int64) *transferReporter {
	localPath, remotePath := fs.paths(path)
	return fs.newTransferReporterForPaths(direction, localPath, remotePath, size)
}

func (fs *RemoteFs) newTransferReporterForPaths(direction events.TransferDirection, localPath string, remotePath string, size int64) *transferReporter {
	if size < 0 {
		size = 0
	}
	startedAt := time.Now()
	return &transferReporter{
		publisher:      fs.eventPublisher(),
		id:             fs.nextTransferID(direction),
		direction:      direction,
		localPath:      localPath,
		remotePath:     remotePath,
		size:           size,
		startedAt:      startedAt,
		lastProgressAt: startedAt,
	}
}

func (fs *RemoteFs) eventPublisher() events.EventPublisher {
	if fs.events != nil {
		return fs.events
	}
	return &events.NoOpEventPublisher{}
}

func (fs *RemoteFs) nextTransferID(direction events.TransferDirection) string {
	seq := atomic.AddUint64(&fs.transferSeq, 1)
	return fmt.Sprintf("fsmount-%s-%d-%d", direction, time.Now().UnixNano(), seq)
}

func (r *transferReporter) Queued() {
	r.publish(events.TransferStatusQueued, 0, time.Time{}, "")
}

func (r *transferReporter) Progress(delta int64) {
	if delta == 0 {
		return
	}

	now := time.Now()
	r.mu.Lock()
	if r.terminalEmitted {
		r.mu.Unlock()
		return
	}

	r.transferredBytes += delta
	r.normalizeTransferredLocked()
	if delta < 0 {
		if r.lastProgressBytes > r.transferredBytes {
			r.lastProgressBytes = r.transferredBytes
		}
		r.mu.Unlock()
		return
	}

	shouldPublish := !r.progressEmitted ||
		r.transferredBytes-r.lastProgressBytes >= transferProgressMinBytes ||
		now.Sub(r.lastProgressAt) >= transferProgressMinInterval
	if !shouldPublish {
		r.mu.Unlock()
		return
	}

	r.progressEmitted = true
	r.lastProgressAt = now
	r.lastProgressBytes = r.transferredBytes
	event := r.eventLocked(events.TransferStatusTransferring, time.Time{}, "")
	r.mu.Unlock()

	r.publisher.Publish(event)
}

func (r *transferReporter) Complete(transferredBytes int64) {
	r.terminal(events.TransferStatusComplete, transferredBytes, "")
}

func (r *transferReporter) Error(err error, transferredBytes int64) {
	if err == nil {
		return
	}

	status := events.TransferStatusErrored
	if errors.Is(err, context.Canceled) {
		status = events.TransferStatusCanceled
	}
	r.terminal(status, transferredBytes, err.Error())
}

func (r *transferReporter) terminal(status events.TransferStatus, transferredBytes int64, message string) {
	r.mu.Lock()
	if r.terminalEmitted {
		r.mu.Unlock()
		return
	}
	if transferredBytes >= 0 {
		r.transferredBytes = transferredBytes
	}
	r.normalizeTransferredLocked()
	if r.size == 0 && r.transferredBytes > 0 {
		r.size = r.transferredBytes
	}
	r.terminalEmitted = true
	event := r.eventLocked(status, time.Now(), message)
	r.mu.Unlock()

	r.publisher.Publish(event)
}

func (r *transferReporter) publish(status events.TransferStatus, transferredBytes int64, endedAt time.Time, message string) {
	r.mu.Lock()
	if transferredBytes >= 0 {
		r.transferredBytes = transferredBytes
	}
	r.normalizeTransferredLocked()
	event := r.eventLocked(status, endedAt, message)
	r.mu.Unlock()

	r.publisher.Publish(event)
}

func (r *transferReporter) eventLocked(status events.TransferStatus, endedAt time.Time, message string) events.TransferEvent {
	return events.TransferEvent{
		ID:               r.id,
		Direction:        r.direction,
		LocalPath:        r.localPath,
		RemotePath:       r.remotePath,
		Size:             r.size,
		TransferredBytes: r.transferredBytes,
		Status:           status,
		StartedAt:        r.startedAt,
		EndedAt:          endedAt,
		Error:            message,
	}
}

func (r *transferReporter) normalizeTransferredLocked() {
	if r.transferredBytes < 0 {
		r.transferredBytes = 0
	}
	if r.size > 0 && r.transferredBytes > r.size {
		r.transferredBytes = r.size
	}
}
