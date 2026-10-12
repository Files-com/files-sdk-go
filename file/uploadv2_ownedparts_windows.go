package file

import (
	"errors"
	"io"
	"os"
	"sync/atomic"
)

// Per-part source handles. On Windows, Go serializes every read of one *os.File: ReadAt holds the file's read and write
// locks for the whole read, and on a synchronous handle also saves and restores the file pointer around it. The parts
// of an upload that read concurrently from the one file it opened therefore wait for each other. So each running part
// of an upload v2 whose file the SDK opened itself reads its range through a handle of its own to the same file
// (ownedPartSource), opened and identified by its handle as the upload's file before any of the part's data is read,
// and closed before the part's result is reported. Elsewhere ReadAt reads concurrently already, parts share the file,
// and uploadv2_ownedsource_notwindows.go stands in for all of this with nothing.

// uploadOwnedFile is the file the SDK opened for an upload, never one a caller passed in (uploadIO.ownedFile).
type uploadOwnedFile struct {
	file *os.File
}

// uploadWithOwnedFile makes file, a local file the SDK opened for this upload and closes after it, the upload's ReaderAt.
// Unlike UploadWithReaderAt, whose file stays its caller's, the upload may read it through handles of its own as well.
func uploadWithOwnedFile(file *os.File) UploadOption {
	return func(params uploadIO) (uploadIO, error) {
		params, err := UploadWithReaderAt(file)(params)
		params.ownedFile = uploadOwnedFile{file: file}
		return params, err
	}
}

// ownedSource returns the file the SDK opened for this upload, while it is still the upload's ReaderAt: a ReaderAt a
// caller gave in a later option makes the upload's source the caller's again.
func (u *uploadIO) ownedSource() (*os.File, bool) {
	if u.ownedFile.file == nil {
		return nil, false
	}
	if file, ok := u.readerAt.(*os.File); ok && file == u.ownedFile.file {
		return file, true
	}
	return nil, false
}

// uploadV2OwnedFile is the regular file the SDK opened for an upload v2 and its identity, when the upload's running parts
// read it through handles of their own. Its file is nil otherwise.
type uploadV2OwnedFile struct {
	file *os.File
	id   ownedFileIdentity
}

func newUploadV2OwnedFile(u *uploadIO) uploadV2OwnedFile {
	file, ok := u.ownedSource()
	if !ok {
		return uploadV2OwnedFile{}
	}
	if info, err := file.Stat(); err != nil || !info.Mode().IsRegular() {
		return uploadV2OwnedFile{}
	}
	id, ok := ownedFileIdentityOf(file)
	if !ok {
		return uploadV2OwnedFile{}
	}
	return uploadV2OwnedFile{file: file, id: id}
}

// readerFor returns what a part of readerAt reads its range from, and the part's source: for the file the SDK opened,
// one that gives the part a handle of its own once it runs (runPart); readerAt itself otherwise.
func (o uploadV2OwnedFile) readerFor(readerAt io.ReaderAt) (io.ReaderAt, uploadV2PartSource) {
	if o.file == nil {
		return readerAt, uploadV2PartSource{}
	}
	source := newOwnedPartSource(o.file)
	return source, uploadV2PartSource{owned: source}
}

// uploadV2PartSource is a part's ownedPartSource, if it reads the file the SDK opened.
type uploadV2PartSource struct {
	owned *ownedPartSource
}

// attachSource gives a part that reads the file the SDK opened a handle of its own, as its run starts (runPart): the
// part reads its range through it for every attempt, until its run ends (releaseSource), before either caller of runPart
// sends the part's result.
func (e *uploadV2Engine) attachSource(part *uploadV2Part) {
	if source := part.source.owned; source != nil {
		e.stats.recordOwnedPartHandle(source.attach(e.ownedFile.id))
	}
}

// releaseSource closes the part's handle of its own, if it had one, as its run ends (endRun).
func (p *uploadV2Part) releaseSource() {
	if p != nil && p.source.owned != nil {
		p.source.owned.release()
	}
}

// uploadV2OwnedPartStats counts the running parts that read through a handle of their own, and those that could not
// open or match one and read the upload's file instead.
type uploadV2OwnedPartStats struct {
	handles   int
	fallbacks int
}

func (s *uploadV2SchedulerStats) recordOwnedPartHandle(attached bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if attached {
		s.ownedParts.handles++
	} else {
		s.ownedParts.fallbacks++
	}
}

func (s *uploadV2OwnedPartStats) addLogAttrs(attrs map[string]any) {
	attrs["scheduler_owned_part_handles"] = s.handles
	attrs["scheduler_owned_part_handle_fallbacks"] = s.fallbacks
}

// ownedPartSource is the ReaderAt of one part of an upload whose file the SDK opened: the upload's file, or, while the
// part runs, a handle of its own to the same file.
type ownedPartSource struct {
	file    *os.File
	current atomic.Pointer[os.File]
}

func newOwnedPartSource(file *os.File) *ownedPartSource {
	source := &ownedPartSource{file: file}
	source.current.Store(file)
	return source
}

func (s *ownedPartSource) ReadAt(p []byte, off int64) (int, error) {
	file := s.current.Load()
	n, err := file.ReadAt(p, off)
	if file != s.file && errors.Is(err, os.ErrClosed) {
		// The part ended and closed its handle while its request body was still being read, as an HTTP transport may
		// after a request returns: that read goes on in the upload's file, as it would have without a handle of its own.
		m, err := s.file.ReadAt(p[n:], off+int64(n))
		return n + m, err
	}
	return n, err
}

// attach opens the part's own handle by the upload's file's name and reads through it, if the handle identifies the
// upload's file (id). It reports whether it did; otherwise the part reads the upload's file.
func (s *ownedPartSource) attach(id ownedFileIdentity) bool {
	own, err := os.Open(s.file.Name())
	if err != nil {
		return false
	}
	if ownID, ok := ownedFileIdentityOf(own); !ok || !sameOwnedFile(id, ownID) {
		_ = own.Close()
		return false
	}
	s.current.Store(own)
	return true
}

// release reads through the upload's file again, then closes the part's own handle, if it had one, which waits for a
// read of it under way to finish.
func (s *ownedPartSource) release() {
	if own := s.current.Swap(s.file); own != s.file {
		_ = own.Close()
	}
}
