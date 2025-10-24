// Package io provides I/O utilities for Files.com FUSE mount.
package io

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"sync"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
)

// OrderedPipe provides a wrapper around an io.Pipe that allows out-of-order writes
// to be read in the correct order from the Out [io.PipeReader].
//
// In addition to the Out io.PipeReader, OrderedPipe implements [io.ReaderAt] to allow reading
// correctly ordered data at specific offsets. The io.ReaderAt implementation only supports reading
// data that has already been sorted, returning 0 bytes if attempting to read beyond offset of the
// currently sorted data.
//
// In the Files.com use case, this allows the FUSE file system to read data from the OrderedPipe
// before an upload is finalized.
type OrderedPipe struct {
	// Out is an io.Reader that returns the data written to the OrderedPipe in the correct order.
	Out *io.PipeReader

	// ident is used to provide context about the data being written for logging and debugging.
	ident string

	// logger is used for logging within the OrderedPipe.
	logger log.Logger

	// in is the side of the io.Pipe streaming data from the host
	// to the Filescomfs.Write function. The OrderedPipe writes to "writers",
	// but a reference to the io.PipeWriter is kept here so that it can be closed
	// when the OrderedPipe is closed.
	in *io.PipeWriter

	// writers is used to write data to both the pipe and a temporary file
	// simultaneously.
	writers io.Writer

	// file allows reading from the temporary file created to allow the OrderedPipe
	// to service Read requests while writes are still in progress.
	file *os.File

	// offset is the current write offset in the file.
	offset int64

	// bufCache holds data chunks that have been written out of order.
	bufCache map[int64][]byte

	// cacheMu is used to synchronize access to bufCache.
	cacheMu sync.Mutex

	// closeOnce ensures that the close operation is only performed once.
	closeOnce sync.Once
	closeErr  error
	closed    bool
}

// OrderedPipeOption defines a function type for configuring OrderedPipe.
type OrderedPipeOption func(*OrderedPipe)

// WithLogger sets a custom logger for OrderedPipe.
func WithLogger(logger log.Logger) OrderedPipeOption {
	return func(op *OrderedPipe) {
		op.logger = logger
	}
}

// NewOrderedPipe creates a new OrderedPipe with options.
func NewOrderedPipe(ident string, opts ...OrderedPipeOption) (*OrderedPipe, error) {
	pipeReader, pipeWriter := io.Pipe()
	file, err := os.CreateTemp("", "fsio-ordered-pipe")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file for ordered pipe: %w", err)
	}
	op := &OrderedPipe{
		Out:      pipeReader,
		logger:   &log.NoOpLogger{},
		ident:    ident,
		in:       pipeWriter,
		writers:  io.MultiWriter(pipeWriter, file),
		file:     file,
		offset:   0,
		bufCache: make(map[int64][]byte),
	}
	for _, opt := range opts {
		opt(op)
	}
	return op, nil
}

// WriteAt writes data to the OrderedPipe at the specified offset. Out of order writes are cached
// until they can be written in order.
func (w *OrderedPipe) WriteAt(buff []byte, offset int64) (n int, err error) {
	if w.closed {
		return 0, errors.New("write to closed OrderedPipe")
	}

	// Lock the cache while checking and updating the write offset and cache.
	w.cacheMu.Lock()
	if offset < w.offset {
		// Windows editors sometimes replay a tiny header (~56 bytes) at 0 after pause/resume.
		// If so, accept and discard to keep the stream consistent.
		if offset == 0 && len(buff) <= 64 {
			w.logger.Debug("OrderedPipe.Write: rewind-at-zero %d bytes ignored for %v (current offset %d)", len(buff), w.ident, w.offset)
			w.cacheMu.Unlock()
			return len(buff), nil
		}
		w.cacheMu.Unlock()
		return 0, fmt.Errorf("write at offset %d is less than current offset %d", offset, w.offset)
	}

	if offset > w.offset {
		// Sometimes parts come in out of order. Those parts need to be cached until it's time to write them.
		w.logger.Trace("OrderedPipe.Write: ident=%v, offset %d is greater than write offset %d, caching", w.ident, offset, w.offset)

		// TODO: Allow for configuring the cache size.
		w.bufCache[offset] = slices.Clone(buff)

		// Return that the full buffer was written, otherwise fuse will eventually fail the write.
		w.cacheMu.Unlock()
		return len(buff), nil
	}

	// The current offset if the ordered pipe matches the requested offset for the current WriteAt
	// operation. Write the data to the underlying writers (pipe and temp file).
	n, err = w.writers.Write(buff)
	if err != nil {
		w.cacheMu.Unlock()
		return n, err
	}

	// update the offset so to maintain a record of how many total bytes have been written
	w.offset += int64(n)

	w.logger.Trace("OrderedPipe.Write: ident=%v, wrote %d bytes, new write offset is %d", w.ident, n, w.offset)

	// check to see if there is a part in the cache at the new offset, and if there is
	// recurse
	if part, ok := w.bufCache[w.offset]; ok {
		partOffset := w.offset
		w.cacheMu.Unlock() // explicit unlock to allow recursion without deadlock
		l, err := w.WriteAt(part, partOffset)
		if err != nil {
			return 0, err
		}

		w.logger.Trace("OrderedPipe.Write: ident=%v, wrote %d bytes from cache, new write offset is %d", w.ident, l, w.offset)

		// TODO: consider moving this before calling WriteAt, otherwise parts are not removed from
		// the cache until the all recursive calls to WriteAt return. If there are multiple levels
		// of recursion, the cache could grow unbounded. Maybe flush all cached parts to disk at a
		// certain size threshold, and check for cached parts on disk in addition to in-memory cache
		// when ordering writes.
		w.cacheMu.Lock()
		delete(w.bufCache, partOffset)
		w.cacheMu.Unlock()
	} else {
		w.cacheMu.Unlock()
	}

	return n, err
}

// ReadAt reads data from the OrderedPipe at the specified offset. It only supports reading data
// that has already been written in order. If attempting to read beyond the current write offset,
// it returns 0 bytes.
func (w *OrderedPipe) ReadAt(buff []byte, offset int64) (n int) {
	if w.closed {
		w.logger.Error("OrderedPipe.ReadAt: read from closed OrderedPipe for ident=%v, at offset %d", w.ident, offset)
		return 0
	}
	w.cacheMu.Lock()
	defer w.cacheMu.Unlock()

	if w.file == nil {
		w.logger.Error("OrderedPipe.ReadAt: file is nil for ident=%v, at offset %d", w.ident, offset)
		return 0
	}

	w.logger.Trace("OrderedPipe.ReadAt: attempting to read data for ident=%v, at offset %d", w.ident, offset)
	if offset > w.offset {
		w.logger.Trace("OrderedPipe.ReadAt: ident=%v, offset %d is greater than write offset %d, returning 0 bytes", w.ident, offset, w.offset)
		return 0
	}

	n, err := w.file.ReadAt(buff, offset)
	if err != nil && err != io.EOF {
		// log the error, but return 0. This can happen if reading while writing is in progress
		// but nothing has been written yet.
		w.logger.Error("OrderedPipe.ReadAt: error reading data for ident=%v, at offset %d: %v", w.ident, offset, err)
		return 0
	}
	return n
}

// Offset returns the current write offset of the OrderedPipe.
func (w *OrderedPipe) Offset() int64 {
	w.cacheMu.Lock()
	defer w.cacheMu.Unlock()
	return w.offset
}

// Close closes the OrderedPipe, including the underlying io.PipeWriter and temporary file. The
// temporary file is also removed.
//
// It returns any error encountered during the close operations. A closed OrderedPipe cannot be used
// again.
func (w *OrderedPipe) Close() error {
	w.closeOnce.Do(func() {
		var inErr error
		var fErr error
		var rmErr error
		if w.in != nil {
			inErr = w.in.Close()
			w.in = nil
		}

		if w.file != nil {
			if fErr = w.file.Close(); fErr != nil {
				w.logger.Error("OrderedPipe.Close: error closing file for ident=%v: %v", w.ident, fErr)
			}
			if rmErr = os.Remove(w.file.Name()); rmErr != nil {
				w.logger.Error("OrderedPipe.Close: error removing temporary file for ident=%v: %v", w.ident, rmErr)
			}
			w.file = nil
		}
		w.closeErr = errors.Join(inErr, fErr, rmErr)
		w.closed = true
	})
	return w.closeErr
}
