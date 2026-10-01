//go:build linux || windows

package fsmount

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/lib"
)

var (
	// errStreamIdle stops a stream that stayed paused for policy.pauseTimeout.
	errStreamIdle = errors.New("range stream paused too long")
	// errRangeCacheWrite marks local cache failures, which are not retried.
	errRangeCacheWrite = errors.New("writing downloaded bytes to the cache failed")
)

func (s *rangeStreams) runStream(stream *rangeStream) {
	// A download can outlive the handle that started it. Keep its entry from
	// being evicted until the stream stops.
	s.store.Pin(stream.file.path)
	defer s.store.Unpin(stream.file.path)
	err := s.downloadStream(stream)
	s.finishStream(stream, err)
}

// downloadStream fetches the stream's range, resuming after transient errors.
// It returns nil when the stream stopped on purpose.
func (s *rangeStreams) downloadStream(stream *rangeStream) error {
	pending, err := s.trimToMissing(stream)
	if err != nil || !pending {
		return err
	}
	if err := s.acquireSlot(stream); err != nil {
		if errors.Is(err, errNoAheadSlot) || stream.ctx.Err() != nil {
			return nil
		}
		return err
	}
	defer s.releaseSlot(stream)

	backoff := s.policy.retryInitial
	var failingSince time.Time
	for {
		before := s.position(stream)
		err := s.fetchOnce(stream)
		if err == nil || errors.Is(err, errStreamIdle) || stream.ctx.Err() != nil {
			return nil
		}
		switch {
		case errors.Is(err, errRangeVersionChanged):
			return err
		case errors.Is(err, errRangeNotSatisfiable):
			// The range stayed unsatisfiable after a fresh link: the file shrank.
			return fmt.Errorf("%w: %w", errRangeVersionChanged, err)
		case errors.Is(err, errRangeValidatorUnavailable), errors.Is(err, cache.ErrSparseFilesUnsupported), errors.Is(err, errMalformedRangeResponse):
			if !s.useWholeFile(stream, err) {
				return err
			}
			continue
		case !transientDownloadError(err):
			return err
		}

		if s.position(stream) > before {
			failingSince = time.Time{}
			backoff = s.policy.retryInitial
		}
		if !s.hasReader(stream) {
			return nil
		}
		if failingSince.IsZero() {
			failingSince = time.Now()
		}
		if time.Since(failingSince) >= s.policy.retryGiveUp {
			return err
		}
		s.logger.Debug("RemoteFs: retrying range stream path=%v kind=%v position=%d after %v: %v", stream.file.path, stream.kind, s.position(stream), backoff, err)
		select {
		case <-stream.ctx.Done():
			return nil
		case <-time.After(backoff):
		}
		backoff = min(backoff*2, s.policy.retryMax)
	}
}

// trimToMissing starts the stream at its first uncached byte and ends it at
// the next cached one. It reports whether anything is left to download.
func (s *rangeStreams) trimToMissing(stream *rangeStream) (bool, error) {
	s.mu.Lock()
	meta := stream.file.meta
	requested := cache.ByteRange{Start: stream.start, End: stream.end}
	s.mu.Unlock()
	gaps, err := s.store.MissingRanges(stream.file.path, meta, requested)
	if err != nil || len(gaps) == 0 {
		return false, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if stream.kind != streamWhole {
		stream.start = max(stream.start, gaps[0].Start)
		stream.pos = max(stream.pos, stream.start)
		stream.end = min(stream.end, gaps[0].End)
		stream.furthest = max(stream.furthest, stream.start)
	}
	s.signalLocked(stream.file)
	return stream.pos < stream.end, nil
}

// acquireSlot waits for one of the drive's download slots. A read that is
// waiting takes the slot of a read-ahead stream when none is free; read-ahead
// that finds no free slot is dropped instead of queued.
func (s *rangeStreams) acquireSlot(stream *rangeStream) error {
	queued := false
	leaveQueue := func() {
		if queued {
			s.demandQueued--
			queued = false
		}
	}
	for {
		s.mu.Lock()
		if err := stream.ctx.Err(); err != nil {
			leaveQueue()
			s.mu.Unlock()
			return err
		}
		demand := stream.waiters > 0 || stream.kind == streamWhole
		aheadAllowed := s.demandQueued == 0 && s.aheadRunningLocked() < s.policy.aheadSlots
		if s.running < s.policy.slots && (demand || aheadAllowed) {
			leaveQueue()
			s.running++
			stream.hasSlot = true
			started := s.startSessionLocked(stream)
			s.mu.Unlock()
			// The event goes to the Desktop app, which can be slow to read it.
			if started != nil {
				started.Queued()
			}
			return nil
		}
		if !demand {
			leaveQueue()
			s.mu.Unlock()
			return errNoAheadSlot
		}
		if !queued {
			s.demandQueued++
			queued = true
		}
		if victim := s.preemptionVictimLocked(); victim != nil {
			victim.cancel()
		}
		freed := s.slotFreed
		s.mu.Unlock()

		select {
		case <-stream.ctx.Done():
			s.mu.Lock()
			leaveQueue()
			s.mu.Unlock()
			return stream.ctx.Err()
		case <-freed:
		}
	}
}

func (s *rangeStreams) releaseSlot(stream *rangeStream) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !stream.hasSlot {
		return
	}
	stream.hasSlot = false
	s.running--
	s.slotsChangedLocked()
}

// slotsChangedLocked wakes streams queued for a slot to look again.
func (s *rangeStreams) slotsChangedLocked() {
	close(s.slotFreed)
	s.slotFreed = make(chan struct{})
}

func (s *rangeStreams) aheadRunningLocked() int {
	count := 0
	for _, file := range s.files {
		for _, stream := range file.streams {
			if stream.hasSlot && stream.waiters == 0 && stream.kind != streamWhole {
				count++
			}
		}
	}
	return count
}

// preemptionVictimLocked picks the read-ahead stream to stop so a waiting read
// can have its slot. Complete-file downloads cannot resume and are never picked.
func (s *rangeStreams) preemptionVictimLocked() *rangeStream {
	var victim *rangeStream
	for _, file := range s.files {
		for _, stream := range file.streams {
			if !stream.hasSlot || stream.waiters > 0 || stream.kind == streamWhole || stream.ctx.Err() != nil {
				continue
			}
			if victim == nil || lessUseful(stream, victim) {
				victim = stream
			}
		}
	}
	return victim
}

// fetchOnce sends one request for the rest of the stream and writes what
// arrives into the cache.
func (s *rangeStreams) fetchOnce(stream *rangeStream) error {
	s.mu.Lock()
	meta := stream.file.meta
	kind := stream.kind
	requested := cache.ByteRange{Start: stream.pos, End: stream.end}
	if kind == streamWhole {
		requested = cache.ByteRange{Start: 0, End: meta.Size}
	}
	started := time.Now()
	id := s.startRequestLocked(requested)
	s.mu.Unlock()
	if requested.Start >= requested.End {
		s.finishRequest(id, requested, started, 0, nil)
		return nil
	}

	var response remoteRangeResponse
	var err error
	if kind == streamWhole {
		complete, ok := s.source.(completeRangeSource)
		if !ok {
			err = errors.New("range source does not support complete-file downloads")
		} else {
			response, err = complete.DownloadComplete(stream.ctx, stream.file.path, meta, requested)
		}
	} else {
		response, err = s.source.DownloadRange(stream.ctx, stream.file.path, meta, requested)
	}
	var transferred int64
	if err == nil {
		transferred, err = s.consume(stream, meta, requested, response)
	}
	if response.Body != nil {
		_ = response.Body.Close()
	}
	s.finishRequest(id, requested, started, transferred, err)
	s.logger.Debug(
		"RemoteFs: range stream path=%v kind=%v requested=[%d, %d) returned=[%d, %d) partial=%t transferred=%d duration=%v err=%v",
		stream.file.path, kind, requested.Start, requested.End, response.Returned.Start, response.Returned.End,
		response.Partial, transferred, time.Since(started), err,
	)
	return err
}

// consume validates a response and publishes its body chunk by chunk. It
// pauses between chunks while the stream is too far ahead of its reader.
func (s *rangeStreams) consume(stream *rangeStream, meta cache.EntryMetadata, requested cache.ByteRange, response remoteRangeResponse) (int64, error) {
	if response.Body == nil {
		return 0, fmt.Errorf("%w: range response has no body", errMalformedRangeResponse)
	}
	if response.TotalSize >= 0 && response.TotalSize != meta.Size {
		return 0, fmt.Errorf("%w: remote size %d does not match expected size %d", errRangeVersionChanged, response.TotalSize, meta.Size)
	}
	complete := cache.ByteRange{Start: 0, End: meta.Size}
	if response.Partial && response.Returned != requested {
		return 0, fmt.Errorf("%w: returned range [%d, %d), requested [%d, %d)", errMalformedRangeResponse, response.Returned.Start, response.Returned.End, requested.Start, requested.End)
	}
	if !response.Partial && response.Returned != complete {
		return 0, fmt.Errorf("%w: complete response returned [%d, %d), expected [%d, %d)", errMalformedRangeResponse, response.Returned.Start, response.Returned.End, complete.Start, complete.End)
	}
	meta, err := s.bindResponse(stream, response)
	if err != nil {
		return 0, err
	}

	var transferred int64
	offset := response.Returned.Start
	buffer := make([]byte, s.policy.chunkSize)
	for offset < response.Returned.End {
		if err := s.pauseWhileAhead(stream); err != nil {
			return transferred, err
		}
		want := min(int64(len(buffer)), response.Returned.End-offset)
		if response.Partial {
			s.mu.Lock()
			end := stream.end
			s.mu.Unlock()
			if offset >= end {
				// A newer stream took over the rest of this range.
				return transferred, nil
			}
			want = min(want, end-offset)
		}
		n, readErr := io.ReadFull(response.Body, buffer[:want])
		if n > 0 {
			transferred += int64(n)
			if readErr == nil && offset+int64(n) == response.Returned.End {
				// The bytes that complete the range are published only once
				// the body has been seen to end with them.
				if err := endOfResponse(response); err != nil {
					return transferred, err
				}
			}
			if stream.verifyCached {
				if err := s.matchCached(stream.file.path, meta, buffer[:n], offset); err != nil {
					return transferred, err
				}
			}
			if err := s.publish(stream, meta, buffer[:n], offset); err != nil {
				return transferred, err
			}
			offset += int64(n)
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) || errors.Is(readErr, io.ErrUnexpectedEOF) {
				return transferred, fmt.Errorf("%w: received %d of %d bytes", errShortRangeResponse, offset-response.Returned.Start, response.Returned.End-response.Returned.Start)
			}
			return transferred, readErr
		}
	}
	return transferred, nil
}

// endOfResponse reports an error when the body goes on past the bytes the
// response declared. A complete response that does is a larger file than the
// listing, so it is reported as a changed version.
func endOfResponse(response remoteRangeResponse) error {
	probe := []byte{0}
	n, err := response.Body.Read(probe)
	if n > 0 || errors.Is(err, errLongRangeResponse) {
		long := fmt.Errorf("%w: expected %d bytes", errLongRangeResponse, response.Returned.End-response.Returned.Start)
		if !response.Partial {
			return fmt.Errorf("%w: %w", errRangeVersionChanged, long)
		}
		return long
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

// matchCached checks data, which a response repeats at offset, against the
// bytes the cache already holds there. A difference means the response is
// another version of the file than the one those bytes, which a handle may
// have read, came from.
func (s *rangeStreams) matchCached(path string, meta cache.EntryMetadata, data []byte, offset int64) error {
	requested := cache.ByteRange{Start: offset, End: offset + int64(len(data))}
	gaps, err := s.store.MissingRanges(path, meta, requested)
	if err != nil {
		return err
	}
	cached := make([]byte, len(data))
	cursor := requested.Start
	match := func(end int64) error {
		if end <= cursor {
			return nil
		}
		want := cached[cursor-offset : end-offset]
		n, err := s.store.ReadRange(path, meta, want, cursor)
		if err != nil {
			return err
		}
		if n != len(want) || !bytes.Equal(want, data[cursor-offset:end-offset]) {
			return fmt.Errorf("%w: a response without a validator differs from the cached bytes at offset %d", errRangeVersionChanged, cursor)
		}
		return nil
	}
	for _, gap := range gaps {
		if err := match(gap.Start); err != nil {
			return err
		}
		cursor = gap.End
	}
	return match(requested.End)
}

// bindResponse checks that a response belongs to the version the file's
// cached bytes came from, and records its validator for later requests.
func (s *rangeStreams) bindResponse(stream *rangeStream, response remoteRangeResponse) (cache.EntryMetadata, error) {
	file := stream.file
	meta, uri, err := s.bindResponseLocked(stream, response)
	if err == nil && uri != "" && s.responseAccepted != nil {
		s.responseAccepted(file.path, uri)
	}
	return meta, err
}

func (s *rangeStreams) bindResponseLocked(stream *rangeStream, response remoteRangeResponse) (cache.EntryMetadata, string, error) {
	file := stream.file
	file.publishMu.Lock()
	defer file.publishMu.Unlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	meta := file.meta
	if s.files[file.path] != file || file.err != nil || stream.ctx.Err() != nil {
		return meta, "", context.Canceled
	}
	if modTime := response.File.ModTime(); !modTime.IsZero() && !modTime.Equal(meta.ModTime) {
		return meta, "", errRangeVersionChanged
	}
	if meta.ETag != "" && response.ETag != meta.ETag {
		return meta, "", errRangeVersionChanged
	}
	// Bytes cached from a response without a validator may already have been
	// read, and nothing proves another response holds the same file. A later
	// response must be complete, so it repeats them, and each byte it repeats
	// is compared with the cache; see matchCached.
	if response.Partial && (response.ETag == "" || file.unvalidated) {
		return meta, "", errRangeValidatorUnavailable
	}
	etag := response.ETag
	if file.unvalidated {
		stream.verifyCached = true
		etag = ""
	}
	if etag == "" {
		file.unvalidated = true
	}
	meta.ETag = etag
	file.meta = meta
	if !response.Partial && stream.kind != streamWhole {
		// The server answered a range request with the whole file. Keep it.
		s.becomeWholeLocked(stream)
	}
	return meta, response.File.DownloadUri, nil
}

// useWholeFile switches a file to one complete download after its server
// returned ranges that cannot be verified or stored. It reports false when the
// stream already was a complete download.
func (s *rangeStreams) useWholeFile(stream *rangeStream, cause error) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if stream.kind == streamWhole {
		return false
	}
	s.logger.Debug("RemoteFs: using a complete download for %v: %v", stream.file.path, cause)
	s.becomeWholeLocked(stream)
	return true
}

func (s *rangeStreams) becomeWholeLocked(stream *rangeStream) {
	file := stream.file
	stream.kind = streamWhole
	stream.start, stream.pos, stream.end = 0, 0, file.meta.Size
	file.wholeFile = true
	if len(s.wholeOnly) >= 4096 {
		clear(s.wholeOnly)
	}
	s.wholeOnly[file.path] = file.meta
	for _, other := range file.streams {
		if other != stream && other.active() {
			other.cancel()
		}
	}
	if session := file.session; session != nil {
		session.wholeFile = true
	}
	s.signalLocked(file)
}

func (s *rangeStreams) pauseWhileAhead(stream *rangeStream) error {
	s.mu.Lock()
	if !s.shouldPauseLocked(stream) {
		s.mu.Unlock()
		return nil
	}
	stream.paused = true
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		stream.paused = false
		s.mu.Unlock()
	}()

	timer := time.NewTimer(s.policy.pauseTimeout)
	defer timer.Stop()
	for {
		select {
		case <-stream.ctx.Done():
			return stream.ctx.Err()
		case <-timer.C:
			return errStreamIdle
		case <-stream.wake:
			s.mu.Lock()
			resume := s.canResumeLocked(stream)
			s.mu.Unlock()
			if resume {
				return nil
			}
		}
	}
}

// publish writes one chunk to the cache unless the stream or its file was
// stopped, then wakes the reads waiting on the file.
func (s *rangeStreams) publish(stream *rangeStream, meta cache.EntryMetadata, data []byte, offset int64) error {
	file := stream.file
	file.publishMu.Lock()
	defer file.publishMu.Unlock()

	s.mu.Lock()
	if s.files[file.path] != file || file.err != nil || stream.ctx.Err() != nil {
		err := file.err
		if err == nil {
			err = context.Canceled
		}
		s.mu.Unlock()
		return err
	}
	whole := stream.kind == streamWhole
	s.mu.Unlock()

	var n int
	var err error
	if whole {
		n, err = s.store.WriteCompleteRange(file.path, meta, data, offset)
	} else {
		n, err = s.store.WriteRange(file.path, meta, data, offset)
	}
	if err == nil && n != len(data) {
		err = io.ErrShortWrite
	}
	if err != nil {
		return fmt.Errorf("%w: %w", errRangeCacheWrite, err)
	}

	s.mu.Lock()
	now := time.Now()
	written := int64(len(data))
	stream.transferred += written
	if stream.firstByteAt.IsZero() {
		stream.firstByteAt = now
		stream.firstBytes = stream.transferred
	} else if elapsed := now.Sub(stream.firstByteAt).Seconds(); elapsed > 0 {
		stream.rate = float64(stream.transferred-stream.firstBytes) / elapsed
		if stream.transferred >= 4*int64(s.policy.chunkSize) {
			s.recentRate = 0.8*s.recentRate + 0.2*stream.rate
		}
	}
	stream.pos = max(stream.pos, offset+written)
	var reporter *transferReporter
	progress := written
	if session := file.session; session != nil {
		if whole {
			// A complete download that restarted from the first byte receives
			// bytes the transfer already counted. It reports how far into the
			// file it has come instead.
			progress = max(0, offset+written-session.transferred)
		}
		session.transferred += progress
		reporter = session.reporter
	}
	s.signalLocked(file)
	s.mu.Unlock()
	if reporter != nil && progress > 0 {
		reporter.Progress(progress)
	}
	return nil
}

func (s *rangeStreams) finishStream(stream *rangeStream, err error) {
	file := stream.file
	versionChanged := errors.Is(err, errRangeVersionChanged)
	if versionChanged {
		// Drop the file's cached bytes before any waiting read learns why.
		s.stopFile(file, err, stopInvalidate)
	} else {
		file.publishMu.Lock()
		s.mu.Lock()
		current := s.files[file.path] == file && file.err == nil
		s.mu.Unlock()
		if current {
			s.flush(file)
		}
		file.publishMu.Unlock()
	}

	s.mu.Lock()
	if err == nil && stream.ctx.Err() == nil && stream.pos >= stream.end && s.files[file.path] == file && file.err == nil {
		s.followProbeLocked(stream)
	}
	stream.done = true
	stream.err = err
	stream.cancel()
	if stream.orphan != nil {
		stream.orphan.Stop()
		stream.orphan = nil
	}
	delete(file.streams, stream.id)
	for _, reader := range file.readers {
		if reader.stream == stream {
			reader.stream = nil
		}
	}
	var completion streamCompletion
	if file.session != nil && err != nil && !versionChanged {
		file.session.err = err
	}
	if s.files[file.path] == file && len(file.streams) == 0 {
		completion = s.endSessionLocked(file, nil)
	}
	s.signalLocked(file)
	s.cleanupFileLocked(file)
	s.mu.Unlock()
	completion.publish()
}

// startSessionLocked opens a transfer row for the file's first download in a
// burst. It returns the row's reporter when it created one, so the caller can
// announce it after releasing s.mu.
func (s *rangeStreams) startSessionLocked(stream *rangeStream) *transferReporter {
	file := stream.file
	if file.session != nil {
		return nil
	}
	file.session = &streamSession{wholeFile: stream.kind == streamWhole}
	if s.newReporter == nil {
		return nil
	}
	reporter := s.newReporter(file.path, file.meta)
	file.session.reporter = reporter
	return reporter
}

// endSessionLocked closes the file's transfer row. err, or an error recorded
// by a failed stream, reports the row as failed.
func (s *rangeStreams) endSessionLocked(file *streamFile, err error) streamCompletion {
	session := file.session
	if session == nil {
		return streamCompletion{}
	}
	file.session = nil
	if err == nil {
		err = session.err
	}
	completion := streamCompletion{reporter: session.reporter, transferredBytes: session.transferred, err: err}
	if session.wholeFile && err == nil {
		completion.transferredBytes = file.meta.Size
	}
	return completion
}

// streamCompletion is a finished transfer row, published outside the lock.
type streamCompletion struct {
	reporter         *transferReporter
	transferredBytes int64
	err              error
}

func (completion streamCompletion) publish() {
	if completion.reporter == nil {
		return
	}
	err := completion.err
	if errors.Is(err, errRangeStreamsClosed) {
		// A cache clear or unmount stopped the download on purpose, and a
		// read that still needs the bytes starts a new one.
		err = fmt.Errorf("%w: %w", context.Canceled, err)
	}
	if err != nil {
		completion.reporter.Error(err, transferredBytesUnchanged)
		return
	}
	completion.reporter.Complete(completion.transferredBytes)
}

func (s *rangeStreams) position(stream *rangeStream) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return stream.pos
}

func (s *rangeStreams) hasReader(stream *rangeStream) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return stream.readers > 0 || stream.waiters > 0
}

func (s *rangeStreams) startRequestLocked(requested cache.ByteRange) uint64 {
	s.diagnostics.RequestsStarted++
	s.diagnostics.PlannedBytes += requested.End - requested.Start
	s.diagnostics.ActiveRequests++
	s.diagnostics.PeakActive = max(s.diagnostics.PeakActive, s.diagnostics.ActiveRequests)
	s.nextID++
	return s.nextID
}

func (s *rangeStreams) finishRequest(id uint64, requested cache.ByteRange, started time.Time, transferred int64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.diagnostics.ActiveRequests > 0 {
		s.diagnostics.ActiveRequests--
	}
	s.diagnostics.TransferredBytes += transferred
	status := "complete"
	switch {
	case errors.Is(err, context.Canceled), errors.Is(err, errStreamIdle):
		status = "canceled"
		s.diagnostics.RequestsCanceled++
	case err != nil:
		status = "failed"
		s.diagnostics.RequestsFailed++
	default:
		s.diagnostics.RequestsComplete++
	}
	s.diagnostics.RecentJobs = append(s.diagnostics.RecentJobs, rangeRequestDiagnostic{
		ID:                   id,
		Start:                requested.Start,
		End:                  requested.End,
		PlannedBytes:         requested.End - requested.Start,
		TransferredBytes:     transferred,
		DurationMilliseconds: float64(time.Since(started)) / float64(time.Millisecond),
		Status:               status,
	})
	if len(s.diagnostics.RecentJobs) > rangeRecentRequestLimit {
		s.diagnostics.RecentJobs = append(s.diagnostics.RecentJobs[:0], s.diagnostics.RecentJobs[len(s.diagnostics.RecentJobs)-rangeRecentRequestLimit:]...)
	}
}

const rangeRecentRequestLimit = 64

type rangeRequestDiagnostic struct {
	ID                   uint64  `json:"id"`
	Start                int64   `json:"start"`
	End                  int64   `json:"end"`
	PlannedBytes         int64   `json:"planned_bytes"`
	TransferredBytes     int64   `json:"transferred_bytes"`
	DurationMilliseconds float64 `json:"duration_ms"`
	Status               string  `json:"status"`
}

// rangeStreamDiagnostics is reported by the /debug/ranges endpoint.
type rangeStreamDiagnostics struct {
	SparseRangeReadsEnabled bool                     `json:"sparse_range_reads_enabled"`
	StartedAt               time.Time                `json:"started_at"`
	RequestsStarted         int64                    `json:"requests_started"`
	RequestsComplete        int64                    `json:"requests_complete"`
	RequestsCanceled        int64                    `json:"requests_canceled"`
	RequestsFailed          int64                    `json:"requests_failed"`
	PlannedBytes            int64                    `json:"planned_bytes"`
	TransferredBytes        int64                    `json:"transferred_bytes"`
	ActiveRequests          int                      `json:"active_requests"`
	PeakActive              int                      `json:"peak_active_requests"`
	RecentJobs              []rangeRequestDiagnostic `json:"recent_jobs"`
}

func newRangeStreamDiagnostics() rangeStreamDiagnostics {
	return rangeStreamDiagnostics{
		StartedAt:  time.Now().UTC(),
		RecentJobs: make([]rangeRequestDiagnostic, 0, rangeRecentRequestLimit),
	}
}

func (s *rangeStreams) diagnosticsSnapshot() rangeStreamDiagnostics {
	s.mu.Lock()
	defer s.mu.Unlock()
	snapshot := s.diagnostics
	snapshot.RecentJobs = append([]rangeRequestDiagnostic(nil), s.diagnostics.RecentJobs...)
	return snapshot
}

func (s *rangeStreams) resetDiagnostics() (rangeStreamDiagnostics, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.diagnostics.ActiveRequests != 0 {
		return s.diagnostics, fmt.Errorf("cannot reset range diagnostics with %d active requests", s.diagnostics.ActiveRequests)
	}
	s.diagnostics = newRangeStreamDiagnostics()
	return s.diagnostics, nil
}

// transientDownloadError reports whether a failed request is worth retrying:
// dropped connections, short bodies, timeouts, throttling and server errors.
func transientDownloadError(err error) bool {
	if err == nil || errors.Is(err, context.Canceled) || errors.Is(err, errRangeStreamsClosed) ||
		errors.Is(err, errRangeNotSatisfiable) || errors.Is(err, errRangeCacheWrite) {
		return false
	}
	if files_sdk.IsNotExist(err) || files_sdk.IsNotAuthenticated(err) {
		return false
	}
	var responseErr lib.ResponseError
	if errors.As(err, &responseErr) {
		return retryableStatus(responseErr.StatusCode)
	}
	var apiErr files_sdk.ResponseError
	if errors.As(err, &apiErr) && apiErr.HttpCode != 0 {
		return retryableStatus(apiErr.HttpCode)
	}
	return true
}

func retryableStatus(code int) bool {
	return code == http.StatusRequestTimeout || code == http.StatusTooManyRequests || code >= http.StatusInternalServerError
}
