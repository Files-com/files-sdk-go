package fsmount

import (
	"io"
	"slices"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

type FSWriter interface {
	writeFile(path string, reader io.Reader, mtime *time.Time)
}

type fsNode struct {
	fs                *virtualfs
	path              string
	downloadUri       string
	readerHandle      uint64
	info              fsNodeInfo
	infoExpires       *time.Time
	childPaths        map[string]struct{}
	childPathsExpires *time.Time
	childPathsMutex   sync.Mutex
	writer            *fsNodeWriter
}

type fsNodeInfo struct {
	dir          bool
	size         int64
	creationTime *time.Time
	modTime      time.Time
}

func (n *fsNode) updateInfo(info fsNodeInfo) {
	n.info = info
	n.infoExpires = lib.Ptr(time.Now().Add(n.fs.cacheTTL))
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
	n.childPathsExpires = lib.Ptr(time.Now().Add(n.fs.cacheTTL))
	return
}

func (n *fsNode) updateSize(size int64) {
	n.info.size = size
	n.infoExpires = lib.Ptr(time.Now().Add(n.fs.cacheTTL))
}

func (n *fsNode) infoExpired() bool {
	return n.infoExpires == nil || n.infoExpires.Before(time.Now())
}

func (n *fsNode) childPathsExpired() bool {
	return n.childPathsExpires == nil || n.childPathsExpires.Before(time.Now())
}

func (n *fsNode) openWriter(fsWriter FSWriter, handle uint64) {
	if n.writer == nil {
		nodeWriter := newFsNodeWriter(n, handle)
		n.writer = nodeWriter
		n.downloadUri = ""

		go func() {
			defer nodeWriter.done()
			fsWriter.writeFile(n.path, nodeWriter.out, &n.info.modTime)
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

type fsNodeWriter struct {
	*fsNode
	in        *io.PipeWriter
	out       *io.PipeReader
	handle    uint64
	offset    int64
	completed *sync.Cond
	partCache map[int64][]byte
}

func newFsNodeWriter(node *fsNode, handle uint64) *fsNodeWriter {
	pipeReader, pipeWriter := io.Pipe()
	return &fsNodeWriter{
		fsNode:    node,
		in:        pipeWriter,
		out:       pipeReader,
		handle:    handle,
		completed: sync.NewCond(&sync.Mutex{}),
		partCache: make(map[int64][]byte),
	}
}

func (w *fsNodeWriter) writeAt(buff []byte, offset int64) (n int, err error) {
	if offset < w.offset {
		// This happens on Windows when a write operation is paused. It writes a 56 byte buffer at
		// offset 0. It's unclear how to handle this to properly resume the write.
		w.fs.Trace("Write: path=%v, offset %d is less than write offset %d, closing writer", w.path, offset, w.offset)
		w.closeWriter(true)
		return len(buff), nil
	}

	if offset > w.offset {
		// Sometimes parts come in out of order. We need to cache them until it's time to write them.
		w.fs.Trace("Write: path=%v, offset %d is greater than write offset %d, caching", w.path, offset, w.offset)
		// TODO: Allow for configuring the cache size.
		w.partCache[offset] = slices.Clone(buff)
		// Return that we wrote the full buffer, otherwise fuse will eventually fail the write.
		return len(buff), nil
	}

	n, err = w.in.Write(buff)
	if err != nil {
		return
	}

	w.offset += int64(n)
	w.updateSize(w.offset)

	w.fs.Trace("Write: path=%v, wrote %d bytes, new write offset is %d", w.path, n, w.offset)

	if part, ok := w.partCache[w.offset]; ok {
		partOffset := w.offset
		l, err := w.writeAt(part, partOffset)
		if err != nil {
			return 0, err
		}

		w.fs.Trace("Write: path=%v, wrote %d bytes from cache, new write offset is %d", w.path, l, w.offset)

		delete(w.partCache, partOffset)
	}

	return
}

func (w *fsNodeWriter) close() {
	if w.in != nil {
		w.in.Close()
		w.in = nil
	}
}

func (w *fsNodeWriter) done() {
	if w.out != nil {
		w.out.Close()
		w.out = nil

		w.completed.L.Lock()
		defer w.completed.L.Unlock()
		w.completed.Broadcast()
		w.completed = nil
	}
}

func (w *fsNodeWriter) waitForCompletion() {
	if w.completed != nil {
		w.completed.L.Lock()
		defer w.completed.L.Unlock()
		w.completed.Wait()
	}
}
