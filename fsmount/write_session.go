package fsmount

import (
	"errors"
	"os"
	"sync"
	"time"
)

type writeSession struct {
	path            string
	workingCopyPath string
	workingCopy     *os.File
	baselineSize    int64
	currentSize     int64
	mtime           time.Time
	mtimeExplicit   bool
	hydrated        bool
	dirty           bool
	uploading       bool
	finalizing      bool
	lastUploadErr   error

	handles map[uint64]struct{}
	upload  *activeUpload

	mu   sync.Mutex
	cond *sync.Cond
}

func newWriteSession(path string, mtime time.Time) (*writeSession, error) {
	f, err := os.CreateTemp("", "filescomfs-write-*")
	if err != nil {
		return nil, err
	}

	session := &writeSession{
		path:            path,
		workingCopyPath: f.Name(),
		workingCopy:     f,
		mtime:           mtime,
		handles:         make(map[uint64]struct{}),
	}
	session.cond = sync.NewCond(&session.mu)
	return session, nil
}

func (s *writeSession) closeAndRemoveWorkingCopy() error {
	s.mu.Lock()
	path := s.workingCopyPath
	f := s.workingCopy
	s.workingCopy = nil
	s.mu.Unlock()

	var errs []error
	if f != nil {
		errs = append(errs, f.Close())
	}
	if path != "" {
		errs = append(errs, os.Remove(path))
	}

	return errors.Join(errs...)
}

func (s *writeSession) addHandle(fh uint64) {
	if fh == ^uint64(0) {
		return
	}
	s.mu.Lock()
	s.handles[fh] = struct{}{}
	s.mu.Unlock()
}

func (s *writeSession) removeHandle(fh uint64) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.handles, fh)
	return len(s.handles)
}

func (s *writeSession) hasHandle(fh uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.handles[fh]
	return ok
}

func (s *writeSession) snapshot() writeSessionSnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	return writeSessionSnapshot{
		path:          s.path,
		currentSize:   s.currentSize,
		mtime:         s.mtime,
		hydrated:      s.hydrated,
		dirty:         s.dirty,
		uploading:     s.uploading,
		finalizing:    s.finalizing,
		lastUploadErr: s.lastUploadErr,
		handleCount:   len(s.handles),
	}
}

type writeSessionSnapshot struct {
	path          string
	currentSize   int64
	mtime         time.Time
	hydrated      bool
	dirty         bool
	uploading     bool
	finalizing    bool
	lastUploadErr error
	handleCount   int
}
