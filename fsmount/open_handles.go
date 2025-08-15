package fsmount

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

// OpenHandles tracks currently-open FUSE file handles and allocates fresh IDs.
type OpenHandles struct {
	mu      sync.RWMutex
	entries map[uint64]*fileHandle
	ticker  *time.Ticker
	lib.LeveledLogger
}

const (
	// max int64 as a uint64 - 9,223,372,036,854,775,807
	handleMask = uint64(math.MaxInt64)

	// sweeptInterval is the interval at which the upload sweeper runs to cancel
	// uploads that have not been written to in a long time.
	sweeptInterval = 10 * time.Minute
)

var (
	// Mask with MaxInt64 so int64(h) is always >= 0 if anyone miscasts it.
	next uint64 = 1
)

// NewOpenHandles initializes a new OpenHandles instance.
func NewOpenHandles(logger lib.LeveledLogger) *OpenHandles {
	oh := &OpenHandles{
		entries:       make(map[uint64]*fileHandle),
		LeveledLogger: logger,
	}
	go oh.startUploadSweeper()
	return oh
}

// Open creates a new ID and stores the handle.
// Never hold h.mu while allocating the ID to avoid lock-order issues.
func (h *OpenHandles) Open(n *fsNode, fl FuseFlags) (id uint64, fh *fileHandle) {
	fh = &fileHandle{node: n, FuseFlags: fl, writtenAt: time.Now(), readAt: time.Now()}
	for {
		id = (atomic.AddUint64(&next, 1)) & handleMask
		fh.id = id
		h.mu.Lock()
		if _, clash := h.entries[id]; !clash {
			h.entries[id] = fh
			h.mu.Unlock()
			return id, fh
		}
		h.mu.Unlock()
		// On the vanishingly small chance of a clash, try the next value.
	}
}

// Lookup finds a fileHandle and node without holding the lock during use.
func (h *OpenHandles) Lookup(id uint64) (*fileHandle, *fsNode, bool) {
	h.mu.RLock()
	fh, ok := h.entries[id]
	h.mu.RUnlock()
	return fh, fh.node, ok
}

// Release removes a fileHandle by id. Safe to call even if the handle
// was already closed.
func (h *OpenHandles) Release(id uint64) (*fileHandle, bool) {
	h.mu.Lock()
	fh, ok := h.entries[id]
	if ok {
		delete(h.entries, id)
	}
	h.mu.Unlock()
	return fh, ok
}

// OpenHandles returns a slice of open file handles, excluding the root handle.
func (h *OpenHandles) OpenHandles() []*fileHandle {
	h.mu.RLock()
	defer h.mu.RUnlock()
	seen := make(map[string]struct{})
	handles := make([]*fileHandle, 0, len(h.entries))
	for _, fh := range h.entries {
		if fh.node != nil && fh.node.path != "/" {
			if _, exists := seen[fh.node.path]; !exists {
				seen[fh.node.path] = struct{}{}
				handles = append(handles, fh)
			}
		}
	}
	return handles
}

// ExtendOpenHandleTtls extends the TTL of all open handles.
// This is useful to keep the file handles alive while the OS is still using them.
func (h *OpenHandles) ExtendOpenHandleTtls() {
	handles := h.OpenHandles()
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, handle := range handles {
		handle.node.extendTtl()
	}
}

func (h *OpenHandles) stopUploadSweeper() {
	if h.ticker != nil {
		h.ticker.Stop()
	}
}

func (h *OpenHandles) startUploadSweeper() {
	h.ticker = time.NewTicker(sweeptInterval)
	for range h.ticker.C {
		h.mu.Lock()
		for id, handle := range h.entries {
			// if the handle is not a write operation or has no cancelUpload function, skip it
			if !handle.isWriteOp() || handle.cancelUpload == nil {
				continue
			}

			// if the handle has never been written to, and 24 hours have passed, cancel the upload
			// and remove it from the entries
			if handle.bytesWritten.Load() == 0 && time.Since(handle.writtenAt) > 24*time.Hour {
				h.Debug("Upload sweeper: cancelling upload for handle %d, never written to. Created at: %v", id, handle.writtenAt)
				handle.cancelUpload()
				handle.cancelUpload = nil
				delete(h.entries, id)
				continue
			}
			// if the handle has been written to, and 1 hour has passed since the last write, cancel the upload
			// and remove it from the entries
			if handle.bytesWritten.Load() > 0 && time.Since(handle.writtenAt) > 1*time.Hour {
				h.Debug("Upload sweeper: cancelling upload for handle %d, last written at: %v", id, handle.writtenAt)
				handle.cancelUpload()
				handle.cancelUpload = nil
				delete(h.entries, id)
			}
		}
		h.mu.Unlock()
	}
}
