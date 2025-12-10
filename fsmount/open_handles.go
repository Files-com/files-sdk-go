package fsmount

import (
	"fmt"
	"math"
	"os"
	"sync"
	"sync/atomic"
	"time"

	ff "github.com/Files-com/files-sdk-go/v3/fsmount/internal/flags"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

// OpenHandles tracks currently-open FUSE file handles and allocates fresh IDs.
type OpenHandles struct {
	mu         sync.RWMutex
	entries    map[uint64]*fileHandle
	ticker     *time.Ticker
	stopTicker chan struct{}
	log        lib.LeveledLogger
}

const (
	// max int64 as a uint64 - 9,223,372,036,854,775,807
	handleMask = uint64(math.MaxInt64)

	// sweepInterval is the interval at which the upload sweeper runs to cancel
	// uploads that have not been written to in a long time.
	sweepInterval = 10 * time.Minute

	// uploadIdleTimeout is the duration after which an upload that has not been
	// written to will be cancelled.
	uploadIdleTimeout = 1 * time.Hour

	// unopenedUploadTimeout is the duration after which an upload that has not
	// been written to at all will be cancelled.
	unopenedUploadTimeout = 24 * time.Hour
)

var (
	// Mask with MaxInt64 so int64(h) is always >= 0 if anyone miscasts it.
	next uint64 = 1

	ErrFileHandleInUse = fmt.Errorf("file handle ID already in use")
)

// NewOpenHandles initializes a new OpenHandles instance.
func NewOpenHandles(logger lib.LeveledLogger) *OpenHandles {
	oh := &OpenHandles{
		entries: make(map[uint64]*fileHandle),
		log:     logger,
	}
	go oh.startUploadSweeper()
	return oh
}

// Open creates a new ID and stores the handle.
// Never hold h.mu while allocating the ID to avoid lock-order issues.
func (h *OpenHandles) Open(node *fsNode, flags ff.FuseFlags) (id uint64, fh *fileHandle) {
	fh = &fileHandle{
		node:      node,
		FuseFlags: flags,
		readAt:    time.Time{},
	}
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

// OpenWithFile uses the given *os.File to create and store a new handle.
// Never hold h.mu while allocating the ID to avoid lock-order issues.
func (h *OpenHandles) OpenWithFile(node *fsNode, flags ff.FuseFlags, file *os.File) (id uint64, fh *fileHandle) {
	id, fh = h.Open(node, flags)
	fh.localFile = file
	return id, fh
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

// OpenDirectoryPaths returns unique FUSE paths for directories with open handles.
func (h *OpenHandles) OpenDirectoryPaths() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	paths := make([]string, 0, len(h.entries))
	seen := make(map[string]struct{})
	for _, fh := range h.entries {
		if fh == nil || fh.node == nil {
			continue
		}
		if fh.node.info.nodeType != nodeTypeDir {
			continue
		}
		if _, ok := seen[fh.node.path]; ok {
			continue
		}
		seen[fh.node.path] = struct{}{}
		paths = append(paths, fh.node.path)
	}
	return paths
}

// ExtendOpenHandleTtls extends the TTL of all open handles.
// This is useful to keep the file handles alive while the OS is still using them.
func (h *OpenHandles) ExtendOpenHandleTtls() {
	handles := h.OpenHandles()
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, handle := range handles {
		if handle.isWriteOp() {
			handle.node.extendTtl()
		}
	}
}

func (h *OpenHandles) Close() {
	h.stopUploadSweeper()
	h.mu.Lock()
	defer h.mu.Unlock()
	for id, fh := range h.entries {
		if fh.node != nil {
			fh.node.logger.Debug("Closing file handle %d for path %s", id, fh.node.path)
		} else {
			h.log.Debug("Closing file handle %d with no associated node", id)
		}
		if fh.localFile != nil {
			fh.localFile.Close()
			fh.localFile = nil
		}
		delete(h.entries, id)
	}
}

func (h *OpenHandles) stopUploadSweeper() {
	h.log.Debug("Stopping upload sweeper")
	if h.stopTicker != nil {
		close(h.stopTicker)
		h.stopTicker = nil
	}
	if h.ticker != nil {
		h.ticker.Stop()
		h.ticker = nil
	}
}

// startUploadSweeper periodically checks for uploads that have been idle
// for too long and cancels them.
func (h *OpenHandles) startUploadSweeper() {
	h.stopTicker = make(chan struct{})
	h.ticker = time.NewTicker(sweepInterval)
	for {
		select {
		case <-h.stopTicker:
			h.log.Debug("Upload sweeper stopped")
			return
		case <-h.ticker.C:
			h.mu.Lock()
			for _, handle := range h.entries {
				if handle == nil || handle.node == nil {
					continue
				}
				node := handle.node
				if node.upload == nil {
					continue
				}
				_, bytesWritten, lastActivity := node.uploadStats()
				idle := time.Since(lastActivity)
				h.log.Debug("Upload sweeper: checking upload for path %s, bytes written: %d, last activity: %v, idle time: %v", node.path, bytesWritten, lastActivity, idle)
				// Case A: upload opened but never wrote any bytes — allow a long grace period.
				if bytesWritten == 0 && idle > unopenedUploadTimeout {
					h.log.Debug("Upload sweeper: cancelling unopened upload for path %s, last activity: %v", node.path, lastActivity)
					node.cancelUpload()
					continue
				}
				// Case B: upload has written bytes but has been idle too long.
				if bytesWritten > 0 && idle > uploadIdleTimeout {
					h.log.Debug("Upload sweeper: cancelling idle upload for path %s, last activity: %v", node.path, lastActivity)
					node.cancelUpload()
				}
			}
			h.mu.Unlock()
		}
	}
}
