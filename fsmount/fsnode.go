package fsmount

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

type FSWriter interface {
	writeFile(path string, reader io.Reader, mtime time.Time)
}

type fsNode struct {
	path         string
	downloadUri  string
	readerHandle uint64
	info         fsNodeInfo

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

	// Used to prevent simultaneous lock/unlock operations.
	lockMutex sync.Mutex

	// Used to prevent changes while calling status type methods like isWriterOpen, isLocked, etc.
	statusMu sync.Mutex
}

func (n *fsNode) String() string {
	return fmt.Sprintf("path: %s, uri: %s, h: %d", n.path, n.downloadUri, n.readerHandle)
}

type fsNodeInfo struct {
	nodeType     nodeType
	size         int64
	creationTime time.Time
	modTime      time.Time
	lockOwner    string
}

func (n *fsNode) updateInfo(info fsNodeInfo) {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	if n.info.size != info.size {
		n.downloadUri = ""
	}

	n.info = info
	n.infoExpires = time.Now().Add(n.cacheTTL)
	n.childPathsExpires = time.Time{} // Force a rebuild of child paths (if we're a directory).
}

func (n *fsNode) updateSize(size int64) {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	if n.info.size != size {
		n.downloadUri = ""
	}

	n.info.size = size
	n.infoExpires = time.Now().Add(n.cacheTTL)
}

func (n *fsNode) updateChildPaths(buildChildPaths func(string) (map[string]struct{}, error)) (err error) {
	n.childPathsMutex.Lock()
	defer n.childPathsMutex.Unlock()

	if !n.childPathsExpired() {
		return
	}

	childPaths, err := buildChildPaths(n.path)
	if err != nil {
		return err
	}

	n.childPaths = childPaths
	n.childPathsExpires = time.Now().Add(n.cacheTTL)
	return
}

func (n *fsNode) infoExpired() bool {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	return n.infoExpires.IsZero() || n.infoExpires.Before(time.Now())
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
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	return n.writer != nil
}

func (n *fsNode) openWriter(fsWriter FSWriter, fh uint64) {
	n.writeMu.Lock()
	defer n.writeMu.Unlock()
	// not wrapped in a sync.Once because a node could be cached, written to, and then
	// the writer could be closed, and then the node is accessed again, so we need to
	// be able to open a new writer for the same node.
	if n.writer == nil {
		n.logger.Debug("openWriter from node: %v, ptr: %p, fh: %v", n.String(), n, fh)
		pipe := newOrderedPipe(n.path, fh, n.logger)
		n.writer = pipe
		n.downloadUri = ""

		go func() {
			defer pipe.done()
			fsWriter.writeFile(n.path, pipe.out, n.info.modTime)
		}()
	}
}

func (n *fsNode) closeWriterByHandle(handle uint64) bool {
	if n.writer != nil && n.writer.handle == handle {
		n.closeWriter(true)
		return true
	}

	return false
}

func (n *fsNode) closeWriter(wait bool) {
	if n.writer != nil {
		n.writer.close()
		if wait {
			n.writer.waitForCompletion()
		}
		n.writer = nil
	}
}
