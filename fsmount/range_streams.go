//go:build linux || windows

package fsmount

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
)

var (
	errRangeStreamsClosed  = errors.New("range downloads are closed")
	errRangeVersionChanged = errors.New("remote file version changed")
)

type rangeCacheStore interface {
	InvalidateRanges(path string) error
	RangeMetadata(path string, meta cache.EntryMetadata) (cache.EntryMetadata, bool, error)
	Pin(path string)
	Unpin(path string)
	MissingRanges(path string, meta cache.EntryMetadata, requested cache.ByteRange) ([]cache.ByteRange, error)
	WriteRange(path string, meta cache.EntryMetadata, buff []byte, ofst int64) (int, error)
	WriteCompleteRange(path string, meta cache.EntryMetadata, buff []byte, ofst int64) (int, error)
	ReadRange(path string, meta cache.EntryMetadata, buff []byte, ofst int64) (int, error)
}

type rangeSource interface {
	DownloadRange(ctx context.Context, path string, meta cache.EntryMetadata, requested cache.ByteRange) (remoteRangeResponse, error)
}

type completeRangeSource interface {
	DownloadComplete(ctx context.Context, path string, meta cache.EntryMetadata, requested cache.ByteRange) (remoteRangeResponse, error)
}

type rangeFlusher interface {
	FlushRanges(path string, meta cache.EntryMetadata) error
}

// rangeStreams owns the remote downloads of one mounted drive.
//
// A stream is one HTTP response that starts at some offset and writes toward
// the end of its range as data arrives. Reads never download anything
// themselves: they wait until a stream brings the bytes they need. A stream
// lives while a reader is using it, pauses when it gets too far ahead of its
// reader, and stops a short grace period after its last reader leaves.
type rangeStreams struct {
	store            rangeCacheStore
	source           rangeSource
	policy           streamPolicy
	logger           log.Logger
	newReporter      func(path string, meta cache.EntryMetadata) *transferReporter
	versionChanged   func(path string)
	responseAccepted func(path, uri string)

	mu        sync.Mutex
	files     map[string]*streamFile
	wholeOnly map[string]cache.EntryMetadata
	running   int
	slotFreed chan struct{}
	// demandQueued counts streams with a waiting read that are queued for a
	// slot. Read-ahead is not admitted while any are.
	demandQueued int
	// recentRate is the download speed of recent streams, used for a stream
	// too new to have measured its own.
	recentRate float64
	nextID     uint64
	closed     bool

	diagnostics rangeStreamDiagnostics
}

// streamPolicy holds the rules' tunable values. See defaultStreamPolicy.
type streamPolicy struct {
	// sparse is false when sparse range reads are disabled. Every download is
	// then one complete-file response from byte zero, as before sparse reads.
	sparse bool

	probeSize        int64
	smallFileSize    int64
	reuseWithin      time.Duration
	minReuseDistance int64
	reorderWindow    int64
	streamingAfter   int64
	minLead          int64
	maxLead          int64
	leadWindow       time.Duration
	readerGrace      time.Duration
	pauseTimeout     time.Duration
	retryInitial     time.Duration
	retryMax         time.Duration
	retryGiveUp      time.Duration
	slots            int
	aheadSlots       int
	streamsPerFile   int
	chunkSize        int
}

func defaultStreamPolicy(sparse bool, cacheCapacity int64) streamPolicy {
	policy := streamPolicy{
		sparse:           sparse,
		probeSize:        1 << 20,
		smallFileSize:    2 << 20,
		reuseWithin:      time.Second,
		minReuseDistance: 1 << 20,
		reorderWindow:    8 << 20,
		streamingAfter:   8 << 20,
		minLead:          8 << 20,
		leadWindow:       5 * time.Second,
		readerGrace:      5 * time.Second,
		pauseTimeout:     30 * time.Second,
		retryInitial:     250 * time.Millisecond,
		retryMax:         4 * time.Second,
		retryGiveUp:      30 * time.Second,
		slots:            downloadOpLimit,
		aheadSlots:       downloadOpLimit - 1,
		streamsPerFile:   3,
		chunkSize:        cacheWriteSize,
	}
	if cacheCapacity > 0 {
		policy.maxLead = max(policy.minLead, cacheCapacity/4)
	}
	return policy
}

type streamKind int

const (
	// streamProbe fetches a short fixed range for a first, possibly tiny, read.
	streamProbe streamKind = iota
	// streamRange runs toward the end of the file at the pace of its reader.
	streamRange
	// streamWhole is one complete-file response. It never pauses and is never
	// stopped for lack of readers, because it cannot resume where it left off.
	streamWhole
)

func (k streamKind) String() string {
	switch k {
	case streamProbe:
		return "probe"
	case streamRange:
		return "range"
	default:
		return "whole"
	}
}

// streamFile is the download state of one remote file version.
type streamFile struct {
	path    string
	meta    cache.EntryMetadata
	ctx     context.Context
	cancel  context.CancelFunc
	err     error
	changed chan struct{}
	// stopped closes once stopFile has finished cleaning up after err.
	stopped  chan struct{}
	stopping bool
	streams  map[uint64]*rangeStream
	readers  map[uint64]*streamReader
	session  *streamSession

	// wholeFile is set when the server cannot return ranges that can be proven
	// to belong to one version. The file then uses one complete download.
	wholeFile bool
	// unvalidated is set once bytes without a validator are cached. Nothing
	// proves another response holds the same file, so later responses are
	// checked against those bytes.
	unvalidated bool

	// publishMu keeps cancellation from returning while a cache write for this
	// version is in progress.
	publishMu sync.Mutex
}

// streamReader tracks one open file handle's reading pattern.
type streamReader struct {
	stream  *rangeStream
	lastEnd int64
	hasRead bool
	// sequential counts contiguous bytes read since the last jump. Once it
	// reaches streamingAfter the reader counts as streaming for good, and its
	// misses start range streams directly instead of probes.
	sequential int64
	streaming  bool
	// waiting counts this handle's reads waiting now, and requestedEnd is the
	// furthest byte they asked for. Concurrent out-of-order reads of one
	// sequential reader share its stream.
	waiting      int
	requestedEnd int64
}

type rangeStream struct {
	id     uint64
	file   *streamFile
	kind   streamKind
	start  int64
	end    int64
	pos    int64
	ctx    context.Context
	cancel context.CancelFunc

	waiters  int
	readers  int
	furthest int64
	samples  []streamSample
	rate     float64
	hasSlot  bool
	paused   bool
	followed bool
	done     bool
	err      error
	wake     chan struct{}
	orphan   *time.Timer

	// verifyCached is set on a response that repeats bytes cached without a
	// validator. Each repeated byte must match the cache.
	verifyCached bool

	requested   cache.ByteRange
	transferred int64
	// firstByteAt and firstBytes let rate exclude the wait for the response.
	firstByteAt time.Time
	firstBytes  int64
}

type streamSample struct {
	at       time.Time
	furthest int64
}

// streamSession reports one burst of downloads for a file as one transfer row.
type streamSession struct {
	reporter    *transferReporter
	transferred int64
	wholeFile   bool
	err         error
}

// hydrationReader identifies waits made on behalf of an edit rather than a handle.
const hydrationReader = ^uint64(0)

var errNoAheadSlot = errors.New("no download slot for read-ahead")

func newRangeStreams(store rangeCacheStore, source rangeSource, policy streamPolicy) *rangeStreams {
	return &rangeStreams{
		store:       store,
		source:      source,
		policy:      policy,
		logger:      &log.NoOpLogger{},
		files:       make(map[string]*streamFile),
		wholeOnly:   make(map[string]cache.EntryMetadata),
		slotFreed:   make(chan struct{}),
		diagnostics: newRangeStreamDiagnostics(),
	}
}

// Wait returns once requested is cached for meta. It joins a download that
// will reach requested soon, or starts one. reader identifies the file handle.
func (s *rangeStreams) Wait(ctx context.Context, reader uint64, path string, meta cache.EntryMetadata, requested cache.ByteRange) error {
	return s.WaitAfterJoining(ctx, reader, path, meta, requested, nil)
}

// WaitAfterJoining is Wait for a caller holding a lock that must cover the
// choice of download but not the wait for it. joined runs once, as soon as the
// read has joined or started the download that brings its bytes, or before
// returning if it never does. A download joined that way that is stopped later
// fails the read instead of starting over, so nothing downloads for the file
// the caller looked at once a change to it has stopped its downloads.
func (s *rangeStreams) WaitAfterJoining(ctx context.Context, reader uint64, path string, meta cache.EntryMetadata, requested cache.ByteRange, joined func()) error {
	release := func() {}
	if joined != nil {
		release = sync.OnceFunc(joined)
	}
	defer release()
	if err := validateStreamRequest(path, meta, requested); err != nil {
		return err
	}
	if requested.Start == requested.End {
		return nil
	}
	// previous and waitedOn are the file state and stream this read last
	// waited on. If either failed, the read reports that error instead of
	// starting over, even when it woke up before the failure.
	var previous *streamFile
	var waitedOn *rangeStream
	for {
		s.mu.Lock()
		file, retired, err := s.fileLocked(path, meta)
		if retired != nil {
			// The listing reports a new version. Its readers get an error, and
			// the next lookup starts clean state for the new version.
			s.mu.Unlock()
			s.stopFile(retired, errRangeVersionChanged, stopDiscard)
			continue
		}
		if err != nil {
			s.mu.Unlock()
			if file != nil {
				// A stopped file reports its error after its cleanup finishes.
				<-file.stopped
			}
			return err
		}
		if previous != nil && previous != file && previous.err != nil {
			stopped := previous
			s.mu.Unlock()
			<-stopped.stopped
			return stopped.err
		}
		previous = file
		// Take the wake channel before looking at the cache, so a chunk
		// published after the lookup cannot be missed.
		notify := file.changed
		s.mu.Unlock()

		gaps, err := s.store.MissingRanges(path, meta, cache.ByteRange{Start: requested.Start, End: meta.Size})
		if err != nil {
			return err
		}
		if len(gaps) == 0 || gaps[0].Start >= requested.End {
			s.NoteRead(reader, path, meta, requested)
			return nil
		}

		s.mu.Lock()
		if s.files[path] != file || file.err != nil {
			s.mu.Unlock()
			continue
		}
		if waitedOn != nil && waitedOn.done && waitedOn.err != nil {
			err := waitedOn.err
			s.mu.Unlock()
			return err
		}
		select {
		case <-notify:
			// Something was published after the lookup, so gaps may be stale.
			s.mu.Unlock()
			continue
		default:
		}
		stream := s.streamForLocked(file, reader, gaps[0], requested)
		waitedOn = stream
		stream.waiters++
		s.attachLocked(file, reader, stream)
		state := file.readers[reader]
		state.waiting++
		state.requestedEnd = max(state.requestedEnd, requested.End)
		s.advanceLocked(stream, requested.End)
		s.wakeLocked(stream)
		s.mu.Unlock()
		release()

		select {
		case <-ctx.Done():
			s.mu.Lock()
			s.stopWaitingLocked(file, reader, stream)
			s.mu.Unlock()
			return ctx.Err()
		case <-notify:
		}

		s.mu.Lock()
		s.stopWaitingLocked(file, reader, stream)
		failed := stream.done && stream.err != nil
		streamErr, fileErr := stream.err, file.err
		s.mu.Unlock()
		if fileErr != nil {
			// Return only after the stopped file's cache cleanup, so the next
			// read cannot find bytes the error was about.
			select {
			case <-file.stopped:
			case <-ctx.Done():
			}
			return fileErr
		}
		if failed {
			return streamErr
		}
	}
}

func (s *rangeStreams) stopWaitingLocked(file *streamFile, reader uint64, stream *rangeStream) {
	stream.waiters--
	if state := file.readers[reader]; state != nil && state.waiting > 0 {
		state.waiting--
	}
	if stream.waiters == 0 && stream.hasSlot {
		// The stream is read-ahead now, so a read queued for a slot may take it.
		s.slotsChangedLocked()
	}
	s.maybeOrphanLocked(stream)
}

// EnsureComplete waits until every byte of meta is cached, for an edit baseline.
func (s *rangeStreams) EnsureComplete(ctx context.Context, path string, meta cache.EntryMetadata) error {
	err := s.Wait(ctx, hydrationReader, path, meta, cache.ByteRange{Start: 0, End: meta.Size})
	s.ReleaseReader(path, hydrationReader)
	return err
}

// StartComplete starts a complete-file download without waiting for it. With
// sparse reads disabled, a read that finds a partial entry uses it to finish
// the file the way a normal disabled-mode download would.
func (s *rangeStreams) StartComplete(path string, meta cache.EntryMetadata) {
	if err := validateStreamRequest(path, meta, cache.ByteRange{}); err != nil || meta.Size == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	file, retired, err := s.fileLocked(path, meta)
	if err != nil || retired != nil {
		return
	}
	for _, stream := range file.streams {
		if stream.kind == streamWhole && stream.active() {
			return
		}
	}
	s.startStreamLocked(file, streamWhole, 0, meta.Size)
}

// NoteRead records that reader received read. Streams use it to pace
// themselves, and a sequential reader's probe grows into a range stream.
func (s *rangeStreams) NoteRead(reader uint64, path string, meta cache.EntryMetadata, read cache.ByteRange) {
	if !s.policy.sparse || read.Start == read.End {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	file := s.files[path]
	if file == nil || file.err != nil || !file.meta.Matches(meta) {
		return
	}
	state := file.readers[reader]
	if state == nil {
		state = &streamReader{}
		file.readers[reader] = state
	}
	s.recordReadLocked(state, read)

	stream := state.stream
	if stream == nil || !stream.active() || !rangesTouch(cache.ByteRange{Start: stream.start, End: stream.end}, read) {
		stream = s.streamServingLocked(file, read.Start)
		if stream == nil {
			s.detachLocked(file, reader)
			return
		}
		s.attachLocked(file, reader, stream)
	}
	s.advanceLocked(stream, read.End)
	s.followProbeLocked(stream)
}

// followProbeLocked continues a probe with a range stream once its reader has
// read past the probe's middle, so a sequential reader does not stop at the
// probe's end. It runs on reads and again when the probe finishes.
func (s *rangeStreams) followProbeLocked(probe *rangeStream) {
	file := probe.file
	if probe.kind != streamProbe || probe.followed || probe.readers == 0 || probe.end >= file.meta.Size ||
		probe.furthest <= probe.start+(probe.end-probe.start)/2 || s.streamReachingLocked(file, nil, probe.end) != nil {
		return
	}
	probe.followed = true
	follow := s.startStreamLocked(file, streamRange, probe.end, file.meta.Size)
	follow.furthest = probe.furthest
	for reader, state := range file.readers {
		if state.stream == probe {
			s.attachLocked(file, reader, follow)
		}
	}
}

// ReleaseReader detaches a closed handle from its stream.
func (s *rangeStreams) ReleaseReader(path string, reader uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	file := s.files[path]
	if file == nil {
		return
	}
	s.detachLocked(file, reader)
	delete(file.readers, reader)
	s.cleanupFileLocked(file)
}

// CancelPath stops the downloads for path and waits for any cache write in
// progress. Nothing downloaded before the call is published afterward.
func (s *rangeStreams) CancelPath(path string) {
	s.mu.Lock()
	file := s.files[path]
	s.mu.Unlock()
	if file != nil {
		s.stopFile(file, context.Canceled, stopDiscard)
	}
}

// CancelAll stops the downloads active when it is called.
func (s *rangeStreams) CancelAll() {
	for _, file := range s.snapshotFiles() {
		s.stopFile(file, context.Canceled, stopDiscard)
	}
}

// Close stops every download, keeps what already arrived, and rejects new work.
func (s *rangeStreams) Close() {
	s.mu.Lock()
	s.closed = true
	s.mu.Unlock()
	for _, file := range s.snapshotFiles() {
		s.stopFile(file, errRangeStreamsClosed, stopFlush)
	}
}

func (s *rangeStreams) snapshotFiles() []*streamFile {
	s.mu.Lock()
	defer s.mu.Unlock()
	files := make([]*streamFile, 0, len(s.files))
	for _, file := range s.files {
		files = append(files, file)
	}
	return files
}

// fileLocked returns the state for the listing version meta. When the path has
// state for a different version, it returns that state as retired so the caller
// can stop it outside the lock and try again. A file being stopped is returned
// with its error.
func (s *rangeStreams) fileLocked(path string, meta cache.EntryMetadata) (*streamFile, *streamFile, error) {
	if s.closed {
		return nil, nil, errRangeStreamsClosed
	}
	if file := s.files[path]; file != nil {
		if file.err != nil {
			return file, nil, file.err
		}
		if !file.meta.Matches(meta) {
			if file.meta.ModTime.After(meta.ModTime) {
				// The caller's listing is older than the one in use.
				return nil, nil, errRangeVersionChanged
			}
			return nil, file, nil
		}
		return file, nil, nil
	}

	stored, found, err := s.store.RangeMetadata(path, meta)
	if err != nil {
		return nil, nil, err
	}
	if found {
		meta.ETag = stored.ETag
	}
	ctx, cancel := context.WithCancel(context.Background())
	file := &streamFile{
		path:    path,
		meta:    meta,
		ctx:     ctx,
		cancel:  cancel,
		changed: make(chan struct{}),
		stopped: make(chan struct{}),
		streams: make(map[uint64]*rangeStream),
		readers: make(map[uint64]*streamReader),
	}
	if wholeOnly, ok := s.wholeOnly[path]; ok && wholeOnly.Matches(meta) {
		file.wholeFile = true
	}
	// A handle may already have read bytes cached without a validator, so the
	// next response has to repeat them exactly.
	file.unvalidated = found && len(stored.Ranges) > 0 && stored.ETag == ""
	s.files[path] = file
	return file, nil, nil
}

// streamForLocked returns the stream a read of requested, whose first missing
// bytes are gap, should wait on.
func (s *rangeStreams) streamForLocked(file *streamFile, reader uint64, gap, requested cache.ByteRange) *rangeStream {
	if !s.policy.sparse || file.wholeFile || file.meta.Size <= s.policy.smallFileSize {
		for _, stream := range file.streams {
			if stream.kind == streamWhole && stream.active() {
				return stream
			}
		}
		return s.startStreamLocked(file, streamWhole, 0, file.meta.Size)
	}

	if stream := s.streamReachingLocked(file, file.readers[reader], gap.Start); stream != nil {
		return stream
	}

	kind := streamProbe
	if state := file.readers[reader]; reader == hydrationReader || state != nil && (state.streaming || s.continuesLocked(state, gap.Start)) {
		kind = streamRange
	}
	end := gap.End
	if kind == streamProbe {
		// Cover the whole first read, so a large read is one request.
		end = min(end, gap.Start+max(s.policy.probeSize, 2*(requested.End-gap.Start)))
	}
	return s.startStreamLocked(file, kind, gap.Start, end)
}

// streamReachingLocked returns a stream that will write offset soon: within
// about policy.reuseWithin at its speed, anywhere inside a short probe, or
// within the reorder window of reads the same reader is already waiting on.
func (s *rangeStreams) streamReachingLocked(file *streamFile, reader *streamReader, offset int64) *rangeStream {
	var best *rangeStream
	for _, stream := range file.streams {
		if !stream.active() || offset < stream.pos || offset >= stream.end {
			continue
		}
		soon := offset-stream.pos <= s.reuseDistanceLocked(stream) || stream.kind == streamProbe
		concurrent := reader != nil && reader.stream == stream && reader.waiting > 0 && offset <= reader.requestedEnd+s.policy.reorderWindow
		if !soon && !concurrent {
			continue
		}
		if best == nil || stream.pos > best.pos {
			best = stream
		}
	}
	return best
}

// streamServingLocked returns the stream whose written range contains offset.
func (s *rangeStreams) streamServingLocked(file *streamFile, offset int64) *rangeStream {
	for _, stream := range file.streams {
		if stream.active() && stream.start <= offset && offset < stream.end && offset <= stream.pos+s.reuseDistanceLocked(stream) {
			return stream
		}
	}
	return nil
}

// reuseDistanceLocked is how far ahead of a stream a read can be and still
// wait for it. A stream paused ahead of its reader resumes at full speed once
// a read waits, so the drive's recent speed counts too.
func (s *rangeStreams) reuseDistanceLocked(stream *rangeStream) int64 {
	rate := max(stream.rate, s.recentRate)
	return max(s.policy.minReuseDistance, int64(rate*s.policy.reuseWithin.Seconds()))
}

func (s *rangeStreams) continuesLocked(state *streamReader, offset int64) bool {
	return state.hasRead && offset >= state.lastEnd-s.policy.reorderWindow && offset <= state.lastEnd+s.policy.minReuseDistance
}

func (s *rangeStreams) recordReadLocked(state *streamReader, read cache.ByteRange) {
	if s.continuesLocked(state, read.Start) {
		if read.End > state.lastEnd {
			state.sequential += read.End - state.lastEnd
			state.lastEnd = read.End
		}
	} else {
		state.sequential = read.End - read.Start
		state.lastEnd = read.End
	}
	state.hasRead = true
	if state.sequential >= s.policy.streamingAfter {
		state.streaming = true
	}
}

// startStreamLocked creates a stream for [start, end) and starts its download.
// A newer stream takes over from an older one it overlaps, so the older one is
// trimmed to end where this one starts.
func (s *rangeStreams) startStreamLocked(file *streamFile, kind streamKind, start, end int64) *rangeStream {
	if kind != streamWhole {
		for _, other := range file.streams {
			if !other.active() || other.kind == streamWhole {
				continue
			}
			if other.pos < start && start < other.end {
				other.end = start
			}
			if start < other.start && other.start < end {
				end = other.start
			}
		}
		s.limitStreamsPerFileLocked(file)
	}
	s.nextID++
	ctx, cancel := context.WithCancel(file.ctx)
	stream := &rangeStream{
		id:        s.nextID,
		file:      file,
		kind:      kind,
		start:     start,
		end:       end,
		pos:       start,
		ctx:       ctx,
		cancel:    cancel,
		furthest:  start,
		wake:      make(chan struct{}, 1),
		requested: cache.ByteRange{Start: start, End: end},
	}
	file.streams[stream.id] = stream
	go s.runStream(stream)
	return stream
}

// limitStreamsPerFileLocked makes room for one more stream by stopping the
// least useful one nobody is waiting for.
func (s *rangeStreams) limitStreamsPerFileLocked(file *streamFile) {
	active := 0
	var victim *rangeStream
	for _, stream := range file.streams {
		if !stream.active() || stream.kind == streamWhole {
			continue
		}
		active++
		if stream.waiters == 0 && (victim == nil || lessUseful(stream, victim)) {
			victim = stream
		}
	}
	if active >= s.policy.streamsPerFile && victim != nil {
		victim.cancel()
	}
}

// lessUseful prefers streams without readers, then paused ones, then the one
// furthest ahead of its reader.
func lessUseful(left, right *rangeStream) bool {
	if (left.readers == 0) != (right.readers == 0) {
		return left.readers == 0
	}
	if left.paused != right.paused {
		return left.paused
	}
	return left.pos-left.furthest > right.pos-right.furthest
}

func (s *rangeStreams) attachLocked(file *streamFile, reader uint64, stream *rangeStream) {
	state := file.readers[reader]
	if state == nil {
		state = &streamReader{}
		file.readers[reader] = state
	}
	if state.stream == stream {
		return
	}
	if previous := state.stream; previous != nil {
		previous.readers--
		s.maybeOrphanLocked(previous)
	}
	state.stream = stream
	stream.readers++
	if stream.orphan != nil {
		stream.orphan.Stop()
		stream.orphan = nil
	}
}

func (s *rangeStreams) detachLocked(file *streamFile, reader uint64) {
	state := file.readers[reader]
	if state == nil || state.stream == nil {
		return
	}
	previous := state.stream
	state.stream = nil
	previous.readers--
	s.maybeOrphanLocked(previous)
}

// maybeOrphanLocked stops a stream readerGrace after it loses its last reader.
func (s *rangeStreams) maybeOrphanLocked(stream *rangeStream) {
	if stream.done || stream.kind == streamWhole || stream.readers > 0 || stream.waiters > 0 || stream.orphan != nil {
		return
	}
	var timer *time.Timer
	timer = time.AfterFunc(s.policy.readerGrace, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if stream.orphan != timer {
			return
		}
		stream.orphan = nil
		if stream.readers == 0 && stream.waiters == 0 {
			stream.cancel()
		}
	})
	stream.orphan = timer
}

// advanceLocked records how far a stream's reader has read.
func (s *rangeStreams) advanceLocked(stream *rangeStream, end int64) {
	if end <= stream.furthest {
		return
	}
	stream.furthest = end
	now := time.Now()
	if n := len(stream.samples); n == 0 || now.Sub(stream.samples[n-1].at) >= 100*time.Millisecond {
		stream.samples = append(stream.samples, streamSample{at: now, furthest: end})
	} else {
		stream.samples[n-1].furthest = end
	}
	cutoff := now.Add(-s.policy.leadWindow)
	for len(stream.samples) > 1 && stream.samples[1].at.Before(cutoff) {
		stream.samples = stream.samples[1:]
	}
	s.wakeLocked(stream)
}

// leadLimitLocked is how far a stream may run ahead of its furthest reader:
// the larger of minLead and twice what the reader consumed in leadWindow.
func (s *rangeStreams) leadLimitLocked(stream *rangeStream) int64 {
	var consumed int64
	if len(stream.samples) > 0 {
		consumed = stream.furthest - stream.samples[0].furthest
	}
	limit := max(s.policy.minLead, 2*consumed)
	if s.policy.maxLead > 0 {
		limit = min(limit, s.policy.maxLead)
	}
	return limit
}

func (s *rangeStreams) shouldPauseLocked(stream *rangeStream) bool {
	return stream.kind != streamWhole && stream.waiters == 0 && stream.pos-stream.furthest > s.leadLimitLocked(stream)
}

func (s *rangeStreams) canResumeLocked(stream *rangeStream) bool {
	return stream.waiters > 0 || stream.pos-stream.furthest <= s.leadLimitLocked(stream)/2
}

func (s *rangeStreams) wakeLocked(stream *rangeStream) {
	if !stream.paused {
		return
	}
	select {
	case stream.wake <- struct{}{}:
	default:
	}
}

func (s *rangeStreams) signalLocked(file *streamFile) {
	close(file.changed)
	file.changed = make(chan struct{})
}

// cleanupFileLocked forgets a file with no downloads and no open readers. A
// stopped file is left for stopFile, which still has to finish its cleanup.
func (s *rangeStreams) cleanupFileLocked(file *streamFile) {
	if s.files[file.path] != file || file.err != nil || len(file.streams) != 0 || len(file.readers) != 0 {
		return
	}
	file.cancel()
	delete(s.files, file.path)
}

type stopAction int

const (
	stopDiscard stopAction = iota
	stopFlush
	stopInvalidate
)

// stopFile ends every download for file. It waits for a cache write already in
// progress, so nothing from these downloads is published after it returns.
// Reads waiting on the file are released once the cleanup is done.
func (s *rangeStreams) stopFile(file *streamFile, err error, action stopAction) {
	s.mu.Lock()
	if file.stopping {
		s.mu.Unlock()
		<-file.stopped
		return
	}
	if s.files[file.path] != file {
		s.mu.Unlock()
		return
	}
	file.stopping = true
	file.err = err
	file.cancel()
	completion := s.endSessionLocked(file, err)
	s.mu.Unlock()
	completion.publish()
	defer close(file.stopped)

	file.publishMu.Lock()
	defer file.publishMu.Unlock()
	switch action {
	case stopInvalidate:
		if invalidateErr := s.store.InvalidateRanges(file.path); invalidateErr != nil {
			s.logger.Error("RemoteFs: failed to invalidate changed cache entry %s: %v", file.path, invalidateErr)
		}
		if s.versionChanged != nil {
			s.versionChanged(file.path)
		}
	case stopFlush:
		s.flush(file)
	}

	s.mu.Lock()
	if s.files[file.path] == file {
		delete(s.files, file.path)
	}
	s.signalLocked(file)
	s.mu.Unlock()
}

func (s *rangeStreams) flush(file *streamFile) {
	flusher, ok := s.store.(rangeFlusher)
	if !ok {
		return
	}
	s.mu.Lock()
	meta := file.meta
	s.mu.Unlock()
	if err := flusher.FlushRanges(file.path, meta); err != nil {
		s.logger.Debug("RemoteFs: flushing cached ranges for %s failed: %v", file.path, err)
	}
}

func validateStreamRequest(path string, meta cache.EntryMetadata, requested cache.ByteRange) error {
	if err := meta.Validate(); err != nil {
		return err
	}
	if path != meta.Path {
		return fmt.Errorf("%w: path %q does not match metadata path %q", cache.ErrInvalidEntryMetadata, path, meta.Path)
	}
	if err := requested.Validate(); err != nil {
		return err
	}
	if requested.End > meta.Size {
		return fmt.Errorf("%w: range [%d, %d) exceeds file size %d", cache.ErrInvalidByteRange, requested.Start, requested.End, meta.Size)
	}
	return nil
}

// active reports whether the stream is still downloading or about to.
func (stream *rangeStream) active() bool {
	return !stream.done && stream.ctx.Err() == nil
}

func byteRangeCount(ranges []cache.ByteRange) int64 {
	var count int64
	for _, byteRange := range ranges {
		count += byteRange.End - byteRange.Start
	}
	return count
}

// rangesTouch reports whether two ranges overlap or are adjacent.
func rangesTouch(left, right cache.ByteRange) bool {
	return left.Start <= right.End && right.Start <= left.End
}
