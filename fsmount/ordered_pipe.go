package fsmount

import (
	"fmt"
	"io"
	"slices"
	"sync"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

// orderedPipe coordinates writing to the remote file system through an
// io.Pipe. Incoming write operations are buffered and sorted by offset before
// being written to the remote file system.
type orderedPipe struct {
	logger lib.LeveledLogger

	// path is the path of the file being written to. Only used for logging and debugging.
	path string

	// in is the side of the io.Pipe streaming data from the host
	// to the Filescomfs.Write function.
	in *io.PipeWriter

	// out is the side of the io.Pipe streaming data to the SDK
	// uploader.
	out *io.PipeReader

	// offset is the current write offset in the file.
	offset int64

	// bufCache holds data chunks that have been written out of order.
	bufCache map[int64][]byte

	// cacheMu is used to synchronize access to bufCache.
	cacheMu *sync.Mutex

	// closeOnce ensures that the close operation is only performed once.
	closeOnce sync.Once

	completedCond *sync.Cond
	handle        uint64
}

func newOrderedPipe(path string, handle uint64, logger lib.LeveledLogger) *orderedPipe {
	pipeReader, pipeWriter := io.Pipe()
	return &orderedPipe{
		logger:        logger,
		path:          path,
		in:            pipeWriter,
		out:           pipeReader,
		offset:        0,
		bufCache:      make(map[int64][]byte),
		cacheMu:       &sync.Mutex{},
		handle:        handle,
		completedCond: sync.NewCond(&sync.Mutex{}),
	}
}

func (w *orderedPipe) writeAt(buff []byte, offset int64) (n int, err error) {
	w.cacheMu.Lock()
	if offset < w.offset {
		// This happens on Windows when a write operation is paused. It writes a 56 byte buffer at
		// offset 0. It's unclear how to handle this to properly resume the write.
		w.cacheMu.Unlock()
		return 0, fmt.Errorf("write at offset %d is less than current offset %d", offset, w.offset)
	}

	if offset > w.offset {
		// Sometimes parts come in out of order. We need to cache them until it's time to write them.
		w.logger.Trace("Write: path=%v, offset %d is greater than write offset %d, caching", w.path, offset, w.offset)

		// TODO: Allow for configuring the cache size.
		w.bufCache[offset] = slices.Clone(buff)

		// Return that we wrote the full buffer, otherwise fuse will eventually fail the write.
		w.cacheMu.Unlock()
		return len(buff), nil
	}

	n, err = w.in.Write(buff)
	if err != nil {
		w.cacheMu.Unlock()
		return n, err
	}

	// update the offset so we know how many total bytes have been written
	w.offset += int64(n)

	w.logger.Trace("Write: path=%v, wrote %d bytes, new write offset is %d", w.path, n, w.offset)

	// check to see if there is a part in the cache at the new offset, and if there is
	// recurse
	if part, ok := w.bufCache[w.offset]; ok {
		partOffset := w.offset
		w.cacheMu.Unlock()
		l, err := w.writeAt(part, partOffset)
		if err != nil {
			return 0, err
		}

		w.logger.Trace("Write: path=%v, wrote %d bytes from cache, new write offset is %d", w.path, l, w.offset)

		// TODO: this might be better before calling writeAt, otherwise parts are not removed from the cache
		// until the call to writeAt returns. If there are multiple levels of recursion, the cache could grow
		// unbounded.
		delete(w.bufCache, partOffset)
	} else {
		w.cacheMu.Unlock()
	}

	return n, err
}

func (w *orderedPipe) close() {
	w.closeOnce.Do(func() {
		if w.in != nil {
			w.in.Close()
			w.in = nil
		}
	})
}

func (w *orderedPipe) done() {
	w.close()

	w.completedCond.L.Lock()
	defer w.completedCond.L.Unlock()
	w.completedCond.Broadcast()
	w.completedCond = nil
}

func (w *orderedPipe) waitForCompletion() {
	if w.completedCond != nil {
		w.completedCond.L.Lock()
		defer w.completedCond.L.Unlock()
		w.completedCond.Wait()
	}
}
