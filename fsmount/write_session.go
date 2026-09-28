//go:build linux || windows

package fsmount

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib/privatefile"
)

type writeSession struct {
	path             string
	workingCopyPath  string
	workingCopy      *os.File
	ownedDir         string
	baselineSize     int64
	currentSize      int64
	mtime            time.Time
	mtimeExplicit    bool
	hydrated         bool
	dirty            bool
	uploading        bool
	finalizing       bool
	remoteCommitted  bool
	uploadPathFixed  bool
	renaming         bool
	mutationCount    int
	clearAfterRename bool
	lastUploadErr    error

	handles map[uint64]struct{}
	upload  *activeUpload

	mu   sync.Mutex
	cond *sync.Cond
}

// newWriteSession creates the working copy of the file at path inside workingDir, a
// directory only the current user can read (the mount's private storage). The working copy
// holds the file's content while it is written and uploaded, so it must never be readable
// by other users. When workingDir is empty, as it is for a RemoteFs built without a mount,
// a private directory is created for this session and removed with the working copy.
func newWriteSession(workingDir string, path string, mtime time.Time) (*writeSession, error) {
	var ownedDir string
	if workingDir == "" {
		if err := privatefile.CheckContainer(os.TempDir()); err != nil {
			return nil, fmt.Errorf("temporary directory cannot hold a private working copy: %w", err)
		}
		dir, err := privatefile.MkdirTemp("", "filescomfs-write-*")
		if err != nil {
			return nil, err
		}
		workingDir, ownedDir = dir, dir
	}
	f, err := os.CreateTemp(workingDir, "filescomfs-write-*")
	if err != nil {
		if ownedDir != "" {
			_ = os.Remove(ownedDir)
		}
		return nil, err
	}

	session := &writeSession{
		path:            path,
		workingCopyPath: f.Name(),
		workingCopy:     f,
		ownedDir:        ownedDir,
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
	if s.ownedDir != "" {
		errs = append(errs, os.Remove(s.ownedDir))
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

func (s *writeSession) endMutation() {
	s.mu.Lock()
	s.mutationCount--
	s.cond.Broadcast()
	s.mu.Unlock()
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
