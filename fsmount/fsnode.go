package fsmount

import (
	"io"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

const (
	infoCacheTime = 1 * time.Second
)

type FSWriter interface {
	writeFile(path string, reader io.Reader, mtime *time.Time)
}

type fsNode struct {
	fs                *virtualfs
	path              string
	info              fsNodeInfo
	infoExpires       *time.Time
	writer            *io.PipeWriter
	reader            *io.PipeReader
	uploadCompleted   *sync.Cond
	writeOffset       int64
	partCache         map[int64][]byte
	childPaths        map[string]struct{}
	childPathsExpires *time.Time
}

type fsNodeInfo struct {
	dir          bool
	size         int64
	creationTime *time.Time
	modTime      time.Time
}

func (n *fsNode) updateInfo(info fsNodeInfo) {
	n.info = info
	n.infoExpires = lib.Ptr(time.Now().Add(infoCacheTime))
}

func (n *fsNode) updateChildPaths(childPaths map[string]struct{}) {
	n.childPaths = childPaths
	n.childPathsExpires = lib.Ptr(time.Now().Add(infoCacheTime))
}

func (n *fsNode) infoExpired() bool {
	return n.infoExpires == nil || n.infoExpires.Before(time.Now())
}

func (n *fsNode) childPathsExpired() bool {
	return n.childPathsExpires == nil || n.childPathsExpires.Before(time.Now())
}

func (n *fsNode) openWriter(writer FSWriter) {
	if n.writer == nil {
		n.reader, n.writer = io.Pipe()
		n.partCache = make(map[int64][]byte)
		n.uploadCompleted = sync.NewCond(&sync.Mutex{})

		go func() {
			defer n.closeReader()
			writer.writeFile(n.path, n.reader, &n.info.modTime)
		}()
	}
}

func (n *fsNode) write(buff []byte) (int, error) {
	l, err := n.writer.Write(buff)
	if err != nil {
		return 0, err
	}

	// Remove the part from the cache. No-op if it's not in the cache.
	delete(n.partCache, n.writeOffset)

	n.writeOffset += int64(l)
	n.info.size = n.writeOffset
	n.infoExpires = lib.Ptr(time.Now().Add(infoCacheTime))

	return l, nil
}

func (n *fsNode) closeWriter() {
	if n.writer != nil {
		n.writer.Close()
		n.writer = nil
		n.writeOffset = 0
		n.partCache = nil
	}
}

func (n *fsNode) closeReader() {
	if n.reader != nil {
		n.reader.Close()
		n.reader = nil

		n.uploadCompleted.L.Lock()
		defer n.uploadCompleted.L.Unlock()
		n.uploadCompleted.Broadcast()
		n.uploadCompleted = nil
	}
}

func (n *fsNode) waitForUploadCompletion() {
	if n.uploadCompleted != nil {
		n.uploadCompleted.L.Lock()
		defer n.uploadCompleted.L.Unlock()
		n.uploadCompleted.Wait()
	}
}
