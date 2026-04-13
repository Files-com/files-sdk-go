package fsmount

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

type fsNode struct {
	path         string
	downloadUri  string
	info         fsNodeInfo
	writeSession *writeSession

	// pendingVisible keeps a locally-created file visible in directory listings
	// until the remote listing confirms its existence or the upload fails.
	pendingVisible bool

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

	// Used to synchronize write-session state on the node.
	writeMu sync.Mutex

	// Used to prevent simultaneous lock/unlock operations.
	lockMutex sync.Mutex

	// Used to prevent changes while calling status type methods like isWriterOpen, isLocked, etc.
	statusMu sync.Mutex
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
	doneClosed   bool
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
	uid          uint32
	gid          uint32
}

func (n fsNodeInfo) String() string {
	return fmt.Sprintf("fsNodeInfo{type: %v, size: %d, created: %v, modified: %v, lockOwner: '%s', uid: %d, gid: %d}",
		n.nodeType, n.size, n.creationTime, n.modTime, n.lockOwner, n.uid, n.gid)
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
	if info.uid == 0 {
		info.uid = n.info.uid
	}
	if info.gid == 0 {
		info.gid = n.info.gid
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

func (n *fsNode) updateSizeAtLeast(size int64) {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	if size > n.info.size {
		n.downloadUri = ""
		n.info.size = size
	}
	n.extendTtl()
}

func (n *fsNode) setLockOwner(owner string) {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	n.info.lockOwner = owner
	n.extendTtl()
}

func (n *fsNode) setOwner(uid uint32, gid uint32) {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	n.info.uid = uid
	n.info.gid = gid
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

func (n *fsNode) setPendingVisible() {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	n.pendingVisible = true
}

func (n *fsNode) clearPendingVisible() {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	n.pendingVisible = false
}

func (n *fsNode) isPendingVisible() bool {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	return n.pendingVisible
}

func (n *fsNode) childPathsExpired() bool {
	return n.childPathsExpires.IsZero() || n.childPathsExpires.Before(time.Now())
}

func (n *fsNode) isLocked() bool {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	return n.info.lockOwner != ""
}

func (n *fsNode) pathAndRef() (string, string) {
	if path, ref, ok := n.writeSessionPathAndRef(); ok {
		return path, ref
	}
	return "", ""
}

func (n *fsNode) hasActiveWriteSession() bool {
	n.writeMu.Lock()
	defer n.writeMu.Unlock()
	return n.writeSession != nil
}

func (n *fsNode) getWriteSession() *writeSession {
	n.writeMu.Lock()
	defer n.writeMu.Unlock()
	return n.writeSession
}

func (n *fsNode) ensureWriteSession(initialPath string) (*writeSession, bool, error) {
	n.writeMu.Lock()
	defer n.writeMu.Unlock()

	if n.writeSession != nil {
		return n.writeSession, false, nil
	}

	session, err := newWriteSession(initialPath, n.info.modTime)
	if err != nil {
		return nil, false, err
	}

	session.baselineSize = n.info.size
	session.currentSize = n.info.size
	n.writeSession = session
	n.downloadUri = ""
	return session, true, nil
}

func (n *fsNode) clearWriteSession() error {
	n.writeMu.Lock()
	session := n.writeSession
	n.writeSession = nil
	n.writeMu.Unlock()
	if session == nil {
		return nil
	}
	return session.closeAndRemoveWorkingCopy()
}

func (n *fsNode) updateWriteSessionPath(path string) {
	n.writeMu.Lock()
	if n.writeSession != nil {
		n.writeSession.mu.Lock()
		n.writeSession.path = path
		n.writeSession.mu.Unlock()
	}
	n.writeMu.Unlock()
}

func (n *fsNode) readFromWriteSession(buff []byte, ofst int64) (int, bool, error) {
	n.writeMu.Lock()
	session := n.writeSession
	n.writeMu.Unlock()
	if session == nil {
		return 0, false, nil
	}

	session.mu.Lock()
	size := session.currentSize
	file := session.workingCopy
	session.mu.Unlock()

	if file == nil {
		return 0, true, io.ErrClosedPipe
	}
	if ofst < 0 {
		ofst = 0
	}
	if ofst >= size {
		return 0, true, nil
	}

	maxRead := min(int64(len(buff)), size-ofst)
	nread, err := file.ReadAt(buff[:maxRead], ofst)
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, true, err
	}
	return nread, true, nil
}

func (n *fsNode) writeSessionUploadStats() (string, int64, time.Time, bool) {
	n.writeMu.Lock()
	session := n.writeSession
	n.writeMu.Unlock()
	if session == nil {
		return "", 0, timeZero, false
	}

	session.mu.Lock()
	upload := session.upload
	session.mu.Unlock()
	if upload == nil {
		return "", 0, timeZero, true
	}
	ref, bytes, last := upload.stats()
	return ref, bytes, last, true
}

func (n *fsNode) writeSessionCancelUpload() bool {
	n.writeMu.Lock()
	session := n.writeSession
	n.writeMu.Unlock()
	if session == nil {
		return false
	}

	session.mu.Lock()
	upload := session.upload
	session.upload = nil
	session.uploading = false
	session.finalizing = false
	session.cond.Broadcast()
	session.mu.Unlock()

	if upload != nil {
		upload.cancelUpload()
	}
	return true
}

func (n *fsNode) writeSessionWaitForFinalize(stallTimeout time.Duration) error {
	n.writeMu.Lock()
	session := n.writeSession
	n.writeMu.Unlock()
	if session == nil {
		return nil
	}

	checkInterval := max(stallTimeout/3, time.Second)
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		session.mu.Lock()
		if session.lastUploadErr != nil {
			err := session.lastUploadErr
			session.mu.Unlock()
			return err
		}
		if !session.uploading && !session.finalizing {
			session.mu.Unlock()
			return nil
		}
		upload := session.upload
		session.mu.Unlock()

		if upload == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		select {
		case <-upload.done:
		case <-ticker.C:
			_, _, lastActivity := upload.stats()
			if time.Since(lastActivity) > stallTimeout {
				return context.DeadlineExceeded
			}
		}
	}
}

func (n *fsNode) writeSessionPathAndRef() (string, string, bool) {
	n.writeMu.Lock()
	session := n.writeSession
	n.writeMu.Unlock()
	if session == nil {
		return "", "", false
	}

	session.mu.Lock()
	defer session.mu.Unlock()
	var ref string
	if session.upload != nil {
		ref, _, _ = session.upload.stats()
	}
	return session.path, ref, true
}

func (n *fsNode) writeSessionCaptureRef(ref string) bool {
	n.writeMu.Lock()
	session := n.writeSession
	n.writeMu.Unlock()
	if session == nil {
		return false
	}

	session.mu.Lock()
	upload := session.upload
	session.mu.Unlock()
	if upload == nil {
		return true
	}
	upload.captureRef(ref)
	return true
}

func (n *fsNode) writeSessionStartUpload(cancel context.CancelFunc) (*activeUpload, bool) {
	n.writeMu.Lock()
	session := n.writeSession
	n.writeMu.Unlock()
	if session == nil {
		return nil, false
	}

	session.mu.Lock()
	defer session.mu.Unlock()
	upload := &activeUpload{
		path:         session.path,
		startedAt:    time.Now(),
		cancel:       cancel,
		done:         make(chan struct{}),
		lastActivity: time.Now(),
	}
	session.upload = upload
	session.uploading = true
	session.finalizing = true
	session.cond.Broadcast()
	return upload, true
}

func (n *fsNode) writeSessionFinishUpload(size int64, err error) bool {
	n.writeMu.Lock()
	session := n.writeSession
	n.writeMu.Unlock()
	if session == nil {
		return false
	}

	session.mu.Lock()
	upload := session.upload
	if err == nil {
		session.currentSize = size
		session.dirty = false
		session.lastUploadErr = nil
	} else {
		session.lastUploadErr = err
	}
	session.uploading = false
	session.finalizing = false
	session.upload = nil
	session.cond.Broadcast()
	session.mu.Unlock()

	if upload != nil {
		upload.closeDone()
	}
	return true
}

func (n *fsNode) writeSessionRecordProgress(delta int64) bool {
	n.writeMu.Lock()
	session := n.writeSession
	n.writeMu.Unlock()
	if session == nil {
		return false
	}

	session.mu.Lock()
	upload := session.upload
	session.mu.Unlock()
	if upload != nil {
		upload.recordProgress(delta)
	}
	return true
}

func (n *fsNode) poisonedWriteSessionErr() error {
	n.writeMu.Lock()
	session := n.writeSession
	n.writeMu.Unlock()
	if session == nil {
		return nil
	}

	session.mu.Lock()
	defer session.mu.Unlock()
	return session.lastUploadErr
}

func (n *fsNode) ensureWriteSessionHydrated(hydrate func(*writeSession) error) error {
	n.writeMu.Lock()
	session := n.writeSession
	n.writeMu.Unlock()
	if session == nil {
		return fmt.Errorf("write session missing for %s", n.path)
	}

	session.mu.Lock()
	if session.lastUploadErr != nil {
		err := session.lastUploadErr
		session.mu.Unlock()
		return err
	}
	if session.hydrated {
		session.mu.Unlock()
		return nil
	}
	session.mu.Unlock()

	if err := hydrate(session); err != nil {
		return err
	}

	session.mu.Lock()
	session.hydrated = true
	session.cond.Broadcast()
	session.mu.Unlock()
	return nil
}

func (n *fsNode) captureRef(ref string) {
	if n.writeSessionCaptureRef(ref) {
		return
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

// markDeleted clears the remote-file identity from the node so surviving open
// handles behave like a new file rather than a stale reference to a deleted object.
func (n *fsNode) markDeleted() {
	n.updateSize(0)
	n.clearDownloadURI()
	n.expireInfo()
}

func (n *fsNode) uploadActive() bool {
	_, _, _, ok := n.writeSessionUploadStats()
	return ok
}

// waitForUploadWithProgressTimeout blocks until the upload completes or stalls.
// Unlike waitForUploadIfFinalizing, the stall deadline resets whenever upload
// progress is recorded, so a large but actively progressing upload never times
// out prematurely. Returns an error only when no progress has been observed for
// stallTimeout.
func (n *fsNode) waitForUploadWithProgressTimeout(stallTimeout time.Duration) error {
	return n.writeSessionWaitForFinalize(stallTimeout)
}

// uploadStats safely returns (ref, bytes, lastActivity) for logging/metrics.
func (n *fsNode) uploadStats() (string, int64, time.Time) {
	ref, bytes, last, _ := n.writeSessionUploadStats()
	return ref, bytes, last
}

func (n *fsNode) cancelUpload() {
	if n.writeSessionCancelUpload() {
		_ = n.clearWriteSession()
	}
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
	if !u.doneClosed && u.done != nil {
		close(u.done)
		u.doneClosed = true
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
