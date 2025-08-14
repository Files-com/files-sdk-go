package fsmount

import (
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

	// Used to prevent simultaneous lock/unlock operations.
	lockMutex sync.Mutex

	// Used to prevent changes while calling status type methods like isWriterOpen, isLocked, etc.
	statusMu sync.Mutex
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

func (n *fsNode) String() string {
	uri := truncate(n.downloadUri, 20) // truncate for readability in logs
	return fmt.Sprintf("fsNode{path: %s, uri: %s, info: %v, expires: %v}", n.path, uri, n.info, n.infoExpires)
}

type fsNodeInfo struct {
	nodeType     nodeType
	size         int64
	creationTime time.Time
	modTime      time.Time
	lockOwner    string
}

func (n fsNodeInfo) String() string {
	return fmt.Sprintf("fsNodeInfo{type: %v, size: %d, created: %v, modified: %v, lockOwner: %s}",
		n.nodeType, n.size, n.creationTime, n.modTime, n.lockOwner)
}

func (n *fsNode) updateInfo(info fsNodeInfo) {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	if n.info.size != info.size {
		n.downloadUri = ""
	}

	n.info = info
	n.extendTtl()
	// Force a rebuild of child paths (if the current node is a directory).
	n.childPathsExpires = time.Time{}
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

func (n *fsNode) incrementSize(size int64) {
	n.statusMu.Lock()
	defer n.statusMu.Unlock()
	n.info.size += size
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
	return err
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

func (n *fsNode) openWriter(fsWriter FSWriter, fh uint64) error {
	n.writeMu.Lock()
	defer n.writeMu.Unlock()
	// not wrapped in a sync.Once because a node could be cached, written to, and then
	// the writer could be closed. If a subsequent request needs to write to this node
	// again, it needs to be able to open a new writer.
	if n.writer == nil {
		n.logger.Debug("openWriter from node: %v, ptr: %p, fh: %v", n.String(), n, fh)
		pipe, err := newOrderedPipe(n.path, fh, n.logger)
		if err != nil {
			return fmt.Errorf("failed to open writer: %v", err)
		}
		n.writer = pipe
		n.downloadUri = ""

		go func() {
			fsWriter.writeFile(n.path, pipe.out, n.info.modTime, fh)
		}()
	}
	return nil
}
