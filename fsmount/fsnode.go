package fsmount

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

type FSWriter interface {
	writeFile(path string, reader io.Reader, mtime time.Time, fh uint64)
}

type fsNode struct {
	path        string
	downloadUri string
	info        fsNodeInfo

	cacheTTL time.Duration
	logger   lib.LeveledLogger

	// infoExpires is the time when the node info is no longer within the cache
	// window.
	infoExpires time.Time

	// the set of paths that are children of this node.
	childPaths map[string]struct{}

	// childPathsExpires is the time when the child paths need to be rebuilt.
	childPathsExpires time.Time

	// childPathsMutex is used to synchronize access to childPaths and childPathsExpires.
	childPathsMutex sync.Mutex

	// coordinates and caches out of order writes to the remote file system
	// until they can be written in the correct order.
	writer *orderedPipe

	// Used to prevent creation of multiple writers for the same node.
	writeMu sync.Mutex

	// writerOwner is the handle id that opened the writer
	writerOwner uint64

	// Used to prevent simultaneous lock/unlock operations.
	lockMutex sync.Mutex

	// Used to prevent changes while calling status type methods like isWriterOpen, isLocked, etc.
	statusMu sync.Mutex

	// upload is the active upload for this node, if any. It is nil if there is no active upload.
	upload *activeUpload

	// uploadMu is used to synchronize access to the upload field.
	uploadMu sync.Mutex
}

var (
	// timeZero is a zero value for time.Time used to indicate that a time has not been set.
	timeZero = time.Time{}
)

type activeUpload struct {
	path         string
	startedAt    time.Time
	cancel       context.CancelFunc
	done         chan struct{}
	ref          string
	bytesWritten int64
	lastActivity time.Time
	// Guards fields inside activeUpload
	mu sync.Mutex
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

type fsNodeInfo struct {
	nodeType     nodeType
	size         int64
	creationTime time.Time
	modTime      time.Time
	lockOwner    string
}

func (n fsNodeInfo) String() string {
	return fmt.Sprintf("fsNodeInfo{type: %v, size: %d, created: %v, modified: %v, lockOwner: '%s'}",
		n.nodeType, n.size, n.creationTime, n.modTime, n.lockOwner)
}

func (n *fsNode) String() string {
	uri := truncate(n.downloadUri, 20) // truncate for readability in logs
	return fmt.Sprintf("fsNode{path: %s, uri: %s, info: %v, expires: %v}", n.path, uri, n.info, n.infoExpires)
}

func (n *fsNode) updateInfo(info fsNodeInfo) {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	if n.info.size != info.size {
		n.downloadUri = ""
	}

	n.info = info
	// Force a rebuild of child paths (if the current node is a directory).
	n.childPathsExpires = timeZero
	if info.nodeType == nodeTypeFile {
		n.extendTtl()
	}
}

func (n *fsNode) extendTtl() {
	n.infoExpires = time.Now().Add(n.cacheTTL)
}

func (n *fsNode) updateSize(size int64) {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	if n.info.size != size {
		n.downloadUri = ""
	}

	n.info.size = size
	n.extendTtl()
}

func (n *fsNode) updateChildPaths(buildChildPaths func(string) (map[string]struct{}, error)) (err error) {
	n.childPathsMutex.Lock()
	defer n.childPathsMutex.Unlock()

	if !n.childPathsExpired() {
		return err
	}

	childPaths, err := buildChildPaths(n.path)
	if err != nil {
		return err
	}

	n.childPaths = childPaths
	n.childPathsExpires = time.Now().Add(n.cacheTTL)
	n.extendTtl()
	return err
}

func (n *fsNode) infoExpired() bool {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	return n.infoExpires.IsZero() || n.infoExpires.Before(time.Now())
}

func (n *fsNode) expireInfo() {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	n.infoExpires = timeZero
	n.downloadUri = ""
}

func (n *fsNode) childPathsExpired() bool {
	return n.childPathsExpires.IsZero() || n.childPathsExpires.Before(time.Now())
}

func (n *fsNode) isLocked() bool {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	return n.info.lockOwner != ""
}

func (n *fsNode) isWriterOpen() bool {
	return n.writerIsOpen()
}

func (n *fsNode) openWriter(fsWriter FSWriter, fh uint64) error {
	n.writeMu.Lock()
	defer n.writeMu.Unlock()
	if n.writer == nil {
		n.logger.Debug("openWriter from node: %v, fh: %v", n.String(), fh)
		pipe, err := newOrderedPipe(n.path, n.logger)
		if err != nil {
			return fmt.Errorf("failed to open writer: %v", err)
		}
		n.writer = pipe
		n.writerOwner = fh
		n.downloadUri = ""
		go func() {
			fsWriter.writeFile(n.path, pipe.out, n.info.modTime, fh)
		}()
	}
	return nil
}

// closeWriter closes the writer if it is open and sets it to nil.
func (n *fsNode) closeWriter() error {
	n.writeMu.Lock()
	defer n.writeMu.Unlock()
	if n.writer != nil {
		n.logger.Debug("closeWriter from node: %s", n.String())
		err := n.writer.close()
		n.writer = nil
		return err
	}
	return nil
}

func (n *fsNode) recordProgress(delta int64) {
	n.uploadMu.Lock()
	defer n.uploadMu.Unlock()
	if n.upload == nil {
		return
	}
	n.upload.recordProgress(delta)
}

func (n *fsNode) startUpload(path string, cancel context.CancelFunc) (int, error) {
	n.uploadMu.Lock()
	defer n.uploadMu.Unlock()
	n.upload = &activeUpload{
		path:         path,
		startedAt:    time.Now(),
		cancel:       cancel,
		done:         make(chan struct{}),
		lastActivity: time.Now(),
	}
	return 0, nil
}

func (n *fsNode) closeUpload(size int64) {
	n.updateSize(size)
	n.uploadMu.Lock()
	defer n.uploadMu.Unlock()
	if n.upload != nil {
		n.upload.closeDone()
		n.upload = nil
	}
}

func (n *fsNode) pathAndRef() (string, string) {
	n.uploadMu.Lock()
	defer n.uploadMu.Unlock()
	var ref string
	if n.upload != nil {
		ref = n.upload.ref
	}
	return n.path, ref
}

func (n *fsNode) captureRef(ref string) {
	n.uploadMu.Lock()
	defer n.uploadMu.Unlock()
	if n.upload != nil {
		n.upload.captureRef(ref)
	}
}

func (n *fsNode) setDownloadURI(uri string) {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	n.downloadUri = uri
}

func (n *fsNode) clearDownloadURI() {
	n.setDownloadURI("")
}

// writerIsOpen reports if a writer is present.
func (n *fsNode) writerIsOpen() bool {
	n.writeMu.Lock()
	defer n.writeMu.Unlock()
	return n.writer != nil
}

// writerSnapshot returns a snapshot of writer pointer, owner, and committed status.
// committed := writer.Offset() > 0
func (n *fsNode) writerSnapshot() (w *orderedPipe, owner uint64, committed bool) {
	n.writeMu.Lock()
	w = n.writer
	owner = n.writerOwner
	if w != nil && w.Offset() > 0 {
		committed = true
	}
	n.writeMu.Unlock()
	return w, owner, committed
}

// adoptWriterIfUncommitted sets writerOwner to fh if writer exists and Offset()==0.
// Returns (adopted, committed).
func (n *fsNode) adoptWriterIfUncommitted(fh uint64) (bool, bool) {
	n.writeMu.Lock()
	defer n.writeMu.Unlock()
	if n.writer == nil {
		return false, false
	}
	if n.writer.Offset() == 0 {
		n.writerOwner = fh
		return true, false
	}
	return false, true
}

// readFromWriter reads from the in-flight writer without exposing locks.
func (n *fsNode) readFromWriter(buff []byte, ofst int64) int {
	n.writeMu.Lock()
	w := n.writer
	n.writeMu.Unlock()
	if w == nil {
		return 0
	}
	return w.readAt(buff, ofst)
}

// waitForUploadIfFinalizing blocks until finalize completes when appropriate.
// It never blocks for an unwritten upload.
func (n *fsNode) waitForUploadIfFinalizing(ctx context.Context) error {
	n.uploadMu.Lock()
	w, _, committed := n.writerSnapshot()
	var done <-chan struct{}
	if n.upload != nil && (w == nil || committed) {
		done = n.upload.done
	}
	n.uploadMu.Unlock()
	if done == nil {
		return nil
	}
	if ctx == nil {
		<-done
		return nil
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// uploadStats safely returns (ref, bytes, lastActivity) for logging/metrics.
func (n *fsNode) uploadStats() (string, int64, time.Time) {
	n.uploadMu.Lock()
	defer n.uploadMu.Unlock()
	if n.upload == nil {
		return "", 0, timeZero
	}
	return n.upload.stats()
}

// Lock order: uploadMu -> writeMu (maintain across call paths to avoid deadlocks).
func (n *fsNode) cancelUpload() {
	n.uploadMu.Lock()
	defer n.uploadMu.Unlock()
	if n.upload == nil {
		return
	}
	n.upload.cancelUpload()
	n.upload = nil
	n.closeWriter()
}

func (u *activeUpload) cancelUpload() {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.cancel != nil {
		u.cancel()
	}
	u.cancel = nil
	u.lastActivity = time.Now()
}

func (u *activeUpload) captureRef(ref string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.ref = ref
	u.lastActivity = time.Now()
}

func (u *activeUpload) closeDone() {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.done != nil {
		close(u.done)
		u.done = nil
	}
}

func (u *activeUpload) recordProgress(delta int64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.bytesWritten += delta
	u.lastActivity = time.Now()
}

func (u *activeUpload) stats() (string, int64, time.Time) {
	u.mu.Lock()
	ref := u.ref
	bytes := u.bytesWritten
	last := u.lastActivity
	u.mu.Unlock()
	return ref, bytes, last
}
