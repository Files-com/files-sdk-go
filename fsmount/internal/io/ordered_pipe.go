// Package io provides I/O utilities for Files.com FUSE mount.
package io

import (
	"errors"
	"fmt"
	"io"
	"os"
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

	// in is the side of the io.Pipe that receives streamed data when Close() is called.
	// A reference to the io.PipeWriter is kept here so that it can be closed
	// when the OrderedPipe is closed.
	in *io.PipeWriter

	// file is the temporary file used to store writes at specific offsets.
	// It allows reading from the OrderedPipe while writes are still in progress.
	file *os.File

	// offset is the current write offset in the file.
	offset int64

	// hasWrites tracks whether any WriteAt calls have occurred (excludes initial content).
	hasWrites bool

	// initialContentSize is the size of content loaded before any writes.
	initialContentSize int64

	// cacheWriter is an optional callback that writes data to cache on every WriteAt.
	cacheWriter CacheWriter

	// mu is used to synchronize access to offset, hasWrites, and file operations.
	mu sync.Mutex

	// closeOnce ensures that the close operation is only performed once.
	closeOnce sync.Once
	closeErr  error
	closed    bool

	// writingDone is closed when FinishedWriting is called, signaling to Read that
	// no more writes will occur.
	writingDone chan struct{}
	writingOnce sync.Once
}

// OrderedPipeOption defines a function type for configuring OrderedPipe.
type OrderedPipeOption func(*OrderedPipe)

// CacheWriter is a function that writes data to cache at a specific offset.
type CacheWriter func(data []byte, offset int64) (int, error)

// WithLogger sets a custom logger for OrderedPipe.
func WithLogger(logger log.Logger) OrderedPipeOption {
	return func(op *OrderedPipe) {
		op.logger = logger
	}
}

// WithCacheWriter sets a callback that writes data to cache on every WriteAt call.
// This allows real-time cache updates as writes occur, avoiding a separate copy step.
func WithCacheWriter(cacheWriter CacheWriter) OrderedPipeOption {
	return func(op *OrderedPipe) {
		op.cacheWriter = cacheWriter
	}
}

// WithInitialContent copies existing file content into the OrderedPipe's temp file
// before any writes occur. This allows partial updates to preserve unmodified data.
func WithInitialContent(reader io.Reader) OrderedPipeOption {
	return func(op *OrderedPipe) {
		if reader != nil && op.file != nil {
			n, err := io.Copy(op.file, reader)
			if err != nil {
				op.logger.Error("OrderedPipe: failed to copy initial content: %v", err)
				return
			}
			op.initialContentSize = n
			op.offset = n
			op.logger.Debug("OrderedPipe: loaded %d bytes of initial content for %v", n, op.ident)

			// Seek back to start for subsequent operations
			if _, err := op.file.Seek(0, io.SeekStart); err != nil {
				op.logger.Error("OrderedPipe: failed to seek after initial content: %v", err)
			}
		}
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
		Out:         pipeReader,
		logger:      &log.NoOpLogger{},
		ident:       ident,
		in:          pipeWriter,
		file:        file,
		offset:      0,
		hasWrites:   false,
		writingDone: make(chan struct{}),
	}
	for _, opt := range opts {
		opt(op)
	}
	return op, nil
}

// WriteAt writes data to the OrderedPipe at the specified offset.
// The underlying os.File handles out-of-order writes automatically.
func (w *OrderedPipe) WriteAt(buff []byte, offset int64) (n int, err error) {
	if w.closed {
		return 0, errors.New("write to closed OrderedPipe")
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	// Mark that actual writes have occurred
	w.hasWrites = true

	// Write directly to the temp file at the specified offset
	// os.File.WriteAt handles sparse files and out-of-order writes
	n, err = w.file.WriteAt(buff, offset)
	if err != nil {
		return n, fmt.Errorf("failed to write to temp file: %w", err)
	}

	// If a cache writer is configured, also write to cache
	if w.cacheWriter != nil {
		if cn, cerr := w.cacheWriter(buff[:n], offset); cerr != nil || cn != n {
			w.logger.Debug("OrderedPipe.WriteAt: cache write failed for ident=%v at offset %d: %v", w.ident, offset, cerr)
			// Continue anyway - cache write failure shouldn't fail the main write
		}
	}

	// Update offset to track the highest written position
	endOffset := offset + int64(n)
	if endOffset > w.offset {
		w.offset = endOffset
	}

	w.logger.Trace("OrderedPipe.WriteAt: ident=%v, wrote %d bytes at offset %d, current max offset is %d",
		w.ident, n, offset, w.offset)

	return n, nil
}

// ReadAt reads data from the OrderedPipe at the specified offset.
// Returns data that has been written to the temp file.
func (w *OrderedPipe) ReadAt(buff []byte, offset int64) (n int) {
	if w.closed {
		w.logger.Error("OrderedPipe.ReadAt: read from closed OrderedPipe for ident=%v, at offset %d", w.ident, offset)
		return 0
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file == nil {
		w.logger.Error("OrderedPipe.ReadAt: file is nil for ident=%v, at offset %d", w.ident, offset)
		return 0
	}

	w.logger.Trace("OrderedPipe.ReadAt: attempting to read data for ident=%v, at offset %d", w.ident, offset)

	// Read directly from temp file - it contains all written data
	n, err := w.file.ReadAt(buff, offset)
	if err != nil && err != io.EOF {
		w.logger.Error("OrderedPipe.ReadAt: error reading data for ident=%v, at offset %d: %v", w.ident, offset, err)
		return 0
	}
	return n
}

// Offset returns the current maximum write offset of the OrderedPipe.
func (w *OrderedPipe) Offset() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.offset
}

// HasWrites reports whether any WriteAt calls have occurred (excludes initial content).
func (w *OrderedPipe) HasWrites() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.hasWrites
}

// streamToPipe copies the temp file content to the pipe writer for upload.
// This provides backpressure so uploads don't appear to complete instantly.
func (w *OrderedPipe) streamToPipe() {
	w.mu.Lock()
	if w.file == nil || w.in == nil {
		w.mu.Unlock()
		return
	}

	// Seek to start of file
	if _, err := w.file.Seek(0, io.SeekStart); err != nil {
		w.logger.Error("OrderedPipe.streamToPipe: failed to seek to start: %v", err)
		w.mu.Unlock()
		_ = w.in.CloseWithError(err)
		return
	}
	w.mu.Unlock()

	// Copy file to pipe (this provides backpressure)
	_, err := io.Copy(w.in, w.file)
	if err != nil {
		w.logger.Error("OrderedPipe.streamToPipe: error copying to pipe: %v", err)
		_ = w.in.CloseWithError(err)
	} else {
		_ = w.in.Close()
	}
	w.logger.Debug("OrderedPipe.streamToPipe: finished streaming for ident=%v", w.ident)
}

// Close closes the OrderedPipe and performs cleanup:
//   - Signals that no more writes will occur
//   - Streams any pending data to the pipe
//   - Closes the underlying io.PipeWriter
//   - Closes and removes the temporary file
//   - Returns any errors encountered (a closed OrderedPipe cannot be reused)
func (w *OrderedPipe) Close() error {
	// Signal that no more writes will occur and start streaming (once only)
	w.writingOnce.Do(func() {
		w.logger.Debug("OrderedPipe.Close: starting stream to pipe for ident=%v", w.ident)

		// Start goroutine to stream file content to pipe
		go func() {
			w.streamToPipe()
			// Signal that streaming is complete
			close(w.writingDone)
		}()
	})

	// Wait for streaming to complete
	<-w.writingDone

	w.closeOnce.Do(func() {
		var fErr error
		var rmErr error

		if w.file != nil {
			if fErr = w.file.Close(); fErr != nil {
				w.logger.Error("OrderedPipe.Close: error closing file for ident=%v: %v", w.ident, fErr)
			}
			if rmErr = os.Remove(w.file.Name()); rmErr != nil {
				w.logger.Error("OrderedPipe.Close: error removing temporary file for ident=%v: %v", w.ident, rmErr)
			}
			w.file = nil
		}
		w.closeErr = errors.Join(fErr, rmErr)
		w.closed = true
	})
	return w.closeErr
}
