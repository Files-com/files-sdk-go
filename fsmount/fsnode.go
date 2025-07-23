package fsmount

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type FSWriter interface {
	writeFile(path string, reader io.Reader, mtime time.Time)
}

type fsNode struct {
	fs                *virtualfs
	path              string
	downloadUri       string
	readerHandle      uint64
	info              fsNodeInfo
	infoExpires       time.Time
	childPaths        map[string]struct{}
	childPathsExpires time.Time
	childPathsMutex   sync.Mutex
	writer            *orderedPipe
	writeMu           sync.Mutex
	lockMutex         sync.Mutex // Used to prevent simultaneous lock/unlock operations.
}

func (n *fsNode) String() string {
	return fmt.Sprintf("path: %s, uri: %s, h: %d", n.path, n.downloadUri, n.readerHandle)
}

type fsNodeInfo struct {
	dir          bool
	size         int64
	creationTime time.Time
	modTime      time.Time
	lockOwner    string
}

func (n *fsNode) updateInfo(info fsNodeInfo) {
	if n.info.size != info.size {
		n.downloadUri = ""
	}

	n.info = info
	n.infoExpires = time.Now().Add(n.fs.cacheTTL)
	n.childPathsExpires = time.Time{} // Force a rebuild of child paths (if we're a directory).
}

func (n *fsNode) updateSize(size int64) {
	if n.info.size != size {
		n.downloadUri = ""
	}

	n.info.size = size
	n.infoExpires = time.Now().Add(n.fs.cacheTTL)
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
	n.childPathsExpires = time.Now().Add(n.fs.cacheTTL)
	return
}

func (n *fsNode) infoExpired() bool {
	return n.infoExpires.IsZero() || n.infoExpires.Before(time.Now())
}

func (n *fsNode) childPathsExpired() bool {
	return n.childPathsExpires.IsZero() || n.childPathsExpires.Before(time.Now())
}

func (n *fsNode) isLocked() bool {
	return n.info.lockOwner != ""
}

func (n *fsNode) openWriter(fsWriter FSWriter, fh uint64) {
	if n.writer == nil {
		n.fs.Debug("openWriter from node: %v, ptr: %p, fh: %v", n.String(), n, fh)
		pipe := newOrderedPipe(n.path, fh, n.fs)
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

func (n *fsNode) isWriterOpen() bool {
	return n.writer != nil
}
