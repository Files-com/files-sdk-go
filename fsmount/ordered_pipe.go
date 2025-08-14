package fsmount

import (
	"fmt"
	"io"
	"os"
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
	// to the Filescomfs.Write function. The orderedPipe writes to "writers",
	// but a reference to the io.PipeWriter is kept here so that it can be closed
	// when the orderedPipe is closed.
	in *io.PipeWriter

	// out is the side of the io.Pipe streaming data to the SDK
	// uploader.
	out *io.PipeReader

	// writers is used to write data to both the pipe and a temporary file
	// simultaneously.
	writers io.Writer

	// file allows reading from the temporary file created to allow the orderedPipe
	// to service Read requests while writes are still in progress.
	file *os.File

	// offset is the current write offset in the file.
	offset int64

	// bufCache holds data chunks that have been written out of order.
	bufCache map[int64][]byte

	// cacheMu is used to synchronize access to bufCache.
	cacheMu *sync.Mutex

	// closeOnce ensures that the close operation is only performed once.
	closeOnce sync.Once
	closeErr  error

	completedCond *sync.Cond
	handle        uint64
}

func newOrderedPipe(path string, handle uint64, logger lib.LeveledLogger) (*orderedPipe, error) {
	pipeReader, pipeWriter := io.Pipe()
	// Create a temporary file to allow reading while writes are still in progress.
	file, err := os.CreateTemp("", "Filescomfs-ordered-pipe")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file for ordered pipe: %w", err)
	}
	return &orderedPipe{
		logger:        logger,
		path:          path,
		in:            pipeWriter,
		out:           pipeReader,
		writers:       io.MultiWriter(pipeWriter, file),
		file:          file,
		offset:        0,
		bufCache:      make(map[int64][]byte),
		cacheMu:       &sync.Mutex{},
		handle:        handle,
		completedCond: sync.NewCond(&sync.Mutex{}),
	}, nil
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
		// Sometimes parts come in out of order. Those parts need to be cached until it's time to write them.
		w.logger.Trace("Write: path=%v, offset %d is greater than write offset %d, caching", w.path, offset, w.offset)

		// TODO: Allow for configuring the cache size.
		w.bufCache[offset] = slices.Clone(buff)

		// Return that the full buffer was written, otherwise fuse will eventually fail the write.
		w.cacheMu.Unlock()
		return len(buff), nil
	}

	n, err = w.writers.Write(buff)
	if err != nil {
		w.cacheMu.Unlock()
		return n, err
	}

	// update the offset so to maintain a record of how many total bytes have been written
	w.offset += int64(n)

	w.logger.Trace("Write: path=%v, wrote %d bytes, new write offset is %d", w.path, n, w.offset)

	// check to see if there is a part in the cache at the new offset, and if there is
	// recurse
	if part, ok := w.bufCache[w.offset]; ok {
		partOffset := w.offset
		w.cacheMu.Unlock() // explicit unlock to allow recursion without deadlock
		l, err := w.writeAt(part, partOffset)
		if err != nil {
			return 0, err
		}

		w.logger.Trace("Write: path=%v, wrote %d bytes from cache, new write offset is %d", w.path, l, w.offset)

		// TODO: this might be better before calling writeAt, otherwise parts are not removed from the cache
		// until the call to writeAt returns. If there are multiple levels of recursion, the cache could grow
		// unbounded.
		w.cacheMu.Lock()
		delete(w.bufCache, partOffset)
		w.cacheMu.Unlock()
	} else {
		w.cacheMu.Unlock()
	}

	return n, err
}

func (w *orderedPipe) readAt(buff []byte, offset int64) (n int) {
	w.cacheMu.Lock()
	defer w.cacheMu.Unlock()

	if w.file == nil {
		w.logger.Error("readAt: file is nil for path=%v, at offset %d", w.path, offset)
		return 0
	}

	w.logger.Trace("readAt: attempting to read data for path=%v, at offset %d", w.path, offset)
	if offset > w.offset {
		w.logger.Trace("readAt: path=%v, offset %d is greater than write offset %d, returning 0 bytes", w.path, offset, w.offset)
		return 0
	}

	n, err := w.file.ReadAt(buff, offset)
	if err != nil && err != io.EOF {
		w.logger.Error("readAt: error reading data for path=%v, at offset %d: %v", w.path, offset, err)
		return 0
	}
	return n
}

func (w *orderedPipe) close() error {
	w.closeOnce.Do(func() {
		if w.in != nil {
			w.closeErr = w.in.Close()
			w.in = nil
		}

		if w.file != nil {
			if err := w.file.Close(); err != nil {
				w.logger.Error("done: error closing file for path=%v: %v", w.path, err)
			}
			if err := os.Remove(w.file.Name()); err != nil {
				w.logger.Error("done: error removing temporary file for path=%v: %v", w.path, err)
			}
			w.file = nil
		}
	})
	return w.closeErr
}
