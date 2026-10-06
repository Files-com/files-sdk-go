//go:build linux || windows

package fsmount

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/fsmount/events"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
)

const streamTestSize = 64 << 20

func TestRangeStreamsSmallFirstReadFetchesOnlyAProbe(t *testing.T) {
	h := newStreamHarness(t, newFakeStreamSource(streamTestSize), testStreamPolicy())
	h.mustRead(1, 0, 4096)
	h.streams.ReleaseReader(h.meta.Path, 1)
	h.waitIdle()

	if got := h.source.requestCount(); got != 1 {
		t.Fatalf("requests = %d, want 1: %v", got, h.source.requestLog())
	}
	if got := h.source.transferred.Load(); got != 1<<20 {
		t.Fatalf("transferred = %d, want one 1 MiB probe", got)
	}
}

func TestRangeStreamsSequentialReadUsesTwoRequests(t *testing.T) {
	h := newStreamHarness(t, newFakeStreamSource(streamTestSize), testStreamPolicy())
	for offset := int64(0); offset < streamTestSize; offset += 1 << 20 {
		h.mustRead(1, offset, 1<<20)
	}
	h.waitIdle()

	requests := h.source.requestLog()
	if len(requests) != 2 || requests[0].requested != (cache.ByteRange{Start: 0, End: 2 << 20}) || requests[1].requested != (cache.ByteRange{Start: 2 << 20, End: streamTestSize}) {
		t.Fatalf("requests = %v, want a probe then one stream to EOF", requests)
	}
}

func TestRangeStreamsNearReadJoinsStreamAndFarReadStartsNewRequest(t *testing.T) {
	source := newFakeStreamSource(streamTestSize)
	source.bytesPerSecond = 8 << 20
	h := newStreamHarness(t, source, testStreamPolicy())
	h.mustRead(1, 0, 1<<20)
	h.mustRead(1, 1<<20, 1<<20)
	h.mustRead(1, 4<<20, 1<<20)
	if got := source.requestCount(); got != 2 {
		t.Fatalf("requests after near read = %d, want 2: %v", got, source.requestLog())
	}

	started := time.Now()
	h.mustRead(1, 40<<20, 64<<10)
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("far read took %v; it waited for the existing stream instead of starting a request", elapsed)
	}
	requests := source.requestLog()
	if len(requests) != 3 || requests[2].requested.Start != 40<<20 {
		t.Fatalf("requests = %v, want a new request at the far offset", requests)
	}
}

func TestRangeStreamsResumeAfterTransientFailureWithoutReadError(t *testing.T) {
	source := newFakeStreamSource(streamTestSize)
	source.failOnce(2<<20, 1<<20)
	h := newStreamHarness(t, source, testStreamPolicy())
	for offset := int64(0); offset < 16<<20; offset += 1 << 20 {
		h.mustRead(1, offset, 1<<20)
	}

	var resumed bool
	for _, request := range source.requestLog() {
		if request.requested.Start == 3<<20 && request.ifMatch == source.etag {
			resumed = true
		}
	}
	if !resumed {
		t.Fatalf("requests = %v, want a resume at the failed offset that sends the ETag", source.requestLog())
	}
}

func TestRangeStreamsFailureReachesOnlyTheWaitingRead(t *testing.T) {
	source := newFakeStreamSource(streamTestSize)
	source.rejectNext.Store(1)
	h := newStreamHarness(t, source, testStreamPolicy())

	if err := h.read(1, 0, 4096); err == nil {
		t.Fatal("read succeeded, want the rejected request's error")
	}
	if err := h.read(1, 0, 4096); err != nil {
		t.Fatalf("read after the failure = %v, want a fresh request to succeed", err)
	}
}

func TestRangeStreamsPauseAheadOfReaderAndResume(t *testing.T) {
	h := newStreamHarness(t, newFakeStreamSource(streamTestSize), testStreamPolicy())
	h.mustRead(1, 0, 1<<20)
	h.mustRead(1, 1<<20, 1<<20)
	paused := h.source.waitStable(t)
	if paused > 12<<20 {
		t.Fatalf("transferred %d bytes for a reader at 2 MiB, want the stream to pause about 8 MiB ahead", paused)
	}

	for offset := int64(2 << 20); offset < 8<<20; offset += 1 << 20 {
		h.mustRead(1, offset, 1<<20)
	}
	if resumed := h.source.waitStable(t); resumed <= paused {
		t.Fatalf("transferred %d bytes after the reader caught up, want more than %d", resumed, paused)
	}
}

func TestRangeStreamsStopAfterLastReaderLeaves(t *testing.T) {
	source := newFakeStreamSource(streamTestSize)
	source.bytesPerSecond = 32 << 20
	h := newStreamHarness(t, source, testStreamPolicy())
	h.mustRead(1, 0, 1<<20)
	h.mustRead(1, 1<<20, 1<<20)
	h.streams.ReleaseReader(h.meta.Path, 1)
	started := time.Now()
	for h.streams.hasActiveStreams() && time.Since(started) < 10*time.Second {
		time.Sleep(5 * time.Millisecond)
	}

	// The test policy's pause timeout is 5 s; only the reader grace stops it
	// sooner. The stop is timed when the download is canceled, because saving
	// what already arrived can take longer on a slow disk.
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("download kept running for %v after its last reader left", elapsed)
	}
	h.waitIdle()
}

func TestRangeStreamsWaitingReadTakesSlotFromReadAhead(t *testing.T) {
	policy := testStreamPolicy()
	policy.slots = 2
	policy.aheadSlots = 2
	source := newFakeStreamSource(streamTestSize)
	h := newStreamHarness(t, source, policy)
	other := h.addFile("/other")
	third := h.addFile("/third")
	waitFor := func(what string, done func() bool) {
		t.Helper()
		deadline := time.Now().Add(10 * time.Second)
		for !done() {
			if time.Now().After(deadline) {
				t.Fatalf("timed out waiting for %s", what)
			}
			time.Sleep(5 * time.Millisecond)
		}
	}

	// Each of the two files gets a range stream started by a waiting read, so
	// it takes a free slot and keeps it as read-ahead once that read is served.
	// The probe a first read starts has to finish first: a read past its middle
	// while it still held a slot would follow it with read-ahead, which is
	// dropped when no slot is free.
	//
	// Both range streams start, and so pause and start their pause timers,
	// after setupStarted. A stream paused for policy.pauseTimeout ends idle and
	// is counted as canceled too, so the cancellation below shows a preemption
	// only if it was counted within that long of setupStarted.
	setupStarted := time.Now()
	for _, meta := range []cache.EntryMetadata{h.meta, other} {
		h.mustReadFile(meta, 1, 0, 1<<20)
		waitFor(meta.Path+"'s probe to finish", func() bool { return h.streams.streamCountFor(meta.Path) == 0 })
		h.mustReadFile(meta, 1, 1<<20, 1<<20)
		h.mustReadFile(meta, 1, 2<<20, 1<<20)
	}
	waitFor("both files' read-ahead to pause holding a slot", func() bool {
		return h.streams.holdsPausedReadAhead(h.meta.Path) && h.streams.holdsPausedReadAhead(other.Path)
	})
	before := h.streams.diagnosticsSnapshot().RequestsCanceled

	started := time.Now()
	h.mustReadFile(third, 2, 0, 4096)
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("waiting read took %v with every slot held by read-ahead", elapsed)
	}
	after := h.streams.diagnosticsSnapshot()
	if elapsed := time.Since(setupStarted); elapsed >= policy.pauseTimeout {
		t.Fatalf("setup and the waiting read took %v, not under the %v pause timeout, so a read-ahead request may have ended idle instead of being stopped for the read", elapsed, policy.pauseTimeout)
	}
	if after.RequestsCanceled == before {
		t.Fatalf("no read-ahead request was stopped to make room: %v", after.RecentJobs)
	}
}

func TestRangeStreamsUnverifiableRangeUsesOneCompleteDownload(t *testing.T) {
	source := newFakeStreamSource(streamTestSize)
	source.etag = ""
	h := newStreamHarness(t, source, testStreamPolicy())
	h.mustRead(1, 20<<20, 64<<10)
	h.mustRead(1, 50<<20, 64<<10)
	h.streams.ReleaseReader(h.meta.Path, 1)
	h.waitIdle()

	requests := source.requestLog()
	if len(requests) != 2 || requests[0].complete || !requests[1].complete {
		t.Fatalf("requests = %v, want one discarded range then one complete download", requests)
	}
	h.mustBeComplete()
}

func TestRangeStreamsRangeAnsweredWithWholeFileKeepsIt(t *testing.T) {
	source := newFakeStreamSource(streamTestSize)
	source.ignoreRanges = true
	h := newStreamHarness(t, source, testStreamPolicy())
	h.mustRead(1, 20<<20, 64<<10)
	h.streams.ReleaseReader(h.meta.Path, 1)
	h.waitIdle()

	if got := source.requestCount(); got != 1 {
		t.Fatalf("requests = %d, want the whole-file answer to be kept: %v", got, source.requestLog())
	}
	h.mustBeComplete()
}

func TestRangeStreamsVersionChangeFailsReadAndDropsCachedRanges(t *testing.T) {
	source := newFakeStreamSource(streamTestSize)
	h := newStreamHarness(t, source, testStreamPolicy())
	var changed atomic.Int32
	h.streams.versionChanged = func(string) { changed.Add(1) }
	h.mustRead(1, 0, 4096)
	h.waitIdle()
	source.setETag(`"replacement"`)

	if err := h.read(1, 40<<20, 4096); !errors.Is(err, errRangeVersionChanged) {
		t.Fatalf("read after replacement = %v, want %v", err, errRangeVersionChanged)
	}
	missing, err := h.store.MissingRanges(h.meta.Path, h.meta, cache.ByteRange{Start: 0, End: 4096})
	if err != nil || len(missing) == 0 {
		t.Fatalf("old version's bytes are still cached: missing=%v err=%v", missing, err)
	}
	if changed.Load() == 0 {
		t.Fatal("the version change was not reported")
	}
}

func TestRangeStreamsDisabledSparseReadsUseOneCompleteDownload(t *testing.T) {
	policy := testStreamPolicy()
	policy.sparse = false
	source := newFakeStreamSource(streamTestSize)
	h := newStreamHarness(t, source, policy)
	h.mustRead(1, 30<<20, 64<<10)
	h.mustRead(1, 0, 64<<10)
	h.streams.ReleaseReader(h.meta.Path, 1)
	h.waitIdle()

	requests := source.requestLog()
	if len(requests) != 1 || !requests[0].complete {
		t.Fatalf("requests = %v, want one complete download", requests)
	}
	h.mustBeComplete()
}

// A complete download that drops starts again from the first byte. Its
// transfer row keeps counting how far into the file it has come, so it never
// reports more of the file than is cached while the prefix downloads again.
func TestRangeStreamsRestartedCompleteDownloadReportsOnlyCachedProgress(t *testing.T) {
	policy := testStreamPolicy()
	policy.sparse = false
	source := newFakeStreamSource(streamTestSize)
	source.failOnce(0, streamTestSize/2)
	h := newStreamHarness(t, source, policy)
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()
	publisher := &captureEventPublisher{}
	var mu sync.Mutex
	var overReported []string
	// Progress is published from the download as each chunk is cached.
	fs.events = publisherFunc(func(event events.MountEvent) {
		publisher.Publish(event)
		transfer, ok := event.(events.TransferEvent)
		if !ok || transfer.Status != events.TransferStatusTransferring {
			return
		}
		missing, err := h.store.MissingRanges(h.meta.Path, h.meta, cache.ByteRange{End: h.meta.Size})
		cached := h.meta.Size - byteRangeCount(missing)
		if err != nil || transfer.TransferredBytes > cached {
			mu.Lock()
			overReported = append(overReported, fmt.Sprintf("%d reported with %d cached (%v)", transfer.TransferredBytes, cached, err))
			mu.Unlock()
		}
	})
	h.streams.newReporter = func(path string, meta cache.EntryMetadata) *transferReporter {
		return fs.newTransferReporter(events.TransferDirectionDownload, path, meta.Size)
	}
	h.mustRead(1, streamTestSize-(64<<10), 64<<10)
	h.streams.ReleaseReader(h.meta.Path, 1)
	h.waitIdle()

	if requests := source.requestLog(); len(requests) != 2 {
		t.Fatalf("requests = %v, want the dropped download and its restart", requests)
	}
	transfers := publisher.waitForTerminalTransfer(t)
	mu.Lock()
	defer mu.Unlock()
	if len(overReported) > 0 {
		t.Fatalf("progress ran ahead of the cached bytes: %v", overReported[0])
	}
	if last := transfers[len(transfers)-1]; last.Status != events.TransferStatusComplete || last.TransferredBytes != streamTestSize {
		t.Fatalf("last transfer event = %s with %d bytes, want complete with %d", last.Status, last.TransferredBytes, streamTestSize)
	}
}

// Transfer events go to the Desktop app, which can stop reading them. A
// download waiting to report that it is queued holds up only its own file.
func TestRangeStreamsBlockedTransferEventDoesNotStallOtherFiles(t *testing.T) {
	h := newStreamHarness(t, newFakeStreamSource(streamTestSize), testStreamPolicy())
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()
	stalled := h.addFile("/stalled-events.bin")
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseEvents := func() { releaseOnce.Do(func() { close(release) }) }
	defer releaseEvents()
	queuedBlocked := make(chan struct{}, 1)
	fs.events = publisherFunc(func(event events.MountEvent) {
		transfer, ok := event.(events.TransferEvent)
		if ok && transfer.Status == events.TransferStatusQueued && strings.Contains(transfer.RemotePath, "stalled-events") {
			select {
			case queuedBlocked <- struct{}{}:
			default:
			}
			<-release
		}
	})
	h.streams.newReporter = func(path string, meta cache.EntryMetadata) *transferReporter {
		return fs.newTransferReporter(events.TransferDirectionDownload, path, meta.Size)
	}
	stalledRead := make(chan error, 1)
	go func() { stalledRead <- h.readFile(stalled, 1, 0, 4096) }()
	<-queuedBlocked

	otherRead := make(chan error, 1)
	go func() { otherRead <- h.readFile(h.meta, 2, 0, 4096) }()
	select {
	case err := <-otherRead:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		releaseEvents()
		t.Fatal("a read of another file waited for a blocked transfer event")
	}
	releaseEvents()
	if err := <-stalledRead; err != nil {
		t.Fatal(err)
	}
}

type publisherFunc func(events.MountEvent)

func (f publisherFunc) Publish(event events.MountEvent) { f(event) }

func TestRangeStreamsUnvalidatedPartialEntryStartsOverWithoutError(t *testing.T) {
	source := newFakeStreamSource(streamTestSize)
	h := newStreamHarness(t, source, testStreamPolicy())
	unvalidated := h.meta
	if _, err := h.store.WriteCompleteRange(unvalidated.Path, unvalidated, source.payload[:1<<20], 0); err != nil {
		t.Fatal(err)
	}
	if err := h.store.FlushRanges(unvalidated.Path, unvalidated); err != nil {
		t.Fatal(err)
	}

	if err := h.read(1, 8<<20, 4096); err != nil {
		t.Fatalf("read with an earlier unvalidated partial entry = %v, want success", err)
	}
}

func TestRangeStreamsEnsureCompleteDownloadsOnlyTheGaps(t *testing.T) {
	source := newFakeStreamSource(streamTestSize)
	h := newStreamHarness(t, source, testStreamPolicy())
	h.mustRead(1, 10<<20, 4096)
	h.mustRead(2, 40<<20, 4096)
	h.waitIdle()
	cached := source.transferred.Load()

	if err := h.streams.EnsureComplete(context.Background(), h.meta.Path, h.meta); err != nil {
		t.Fatal(err)
	}
	h.mustBeComplete()
	if got := source.transferred.Load(); got != streamTestSize {
		t.Fatalf("transferred %d bytes in total with %d already cached, want each byte once", got, cached)
	}
}

func TestRangeStreamsCancelPathStopsPublishing(t *testing.T) {
	source := newFakeStreamSource(streamTestSize)
	source.bytesPerSecond = 16 << 20
	h := newStreamHarness(t, source, testStreamPolicy())
	h.mustRead(1, 0, 1<<20)
	h.mustRead(1, 1<<20, 1<<20)

	h.streams.CancelPath(h.meta.Path)
	before := h.cachedBytes()
	time.Sleep(100 * time.Millisecond)
	if after := h.cachedBytes(); after != before {
		t.Fatalf("cached bytes grew from %d to %d after CancelPath returned", before, after)
	}
}

func waitFor(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for condition")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestRangeStreamsWaitingReadGetsSlotOnceOtherReadsAreSatisfied(t *testing.T) {
	policy := testStreamPolicy()
	policy.slots = 2
	policy.aheadSlots = 2
	source := &gatedStreamSource{payload: make([]byte, streamTestSize)}
	store, err := disk.NewDiskCache(t.TempDir(), disk.WithCapacityBytes(8<<30))
	if err != nil {
		t.Fatal(err)
	}
	streams := newRangeStreams(store, source, policy)
	t.Cleanup(streams.Close)
	modTime := time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)
	read := cache.ByteRange{Start: 0, End: 4096}
	wait := func(reader uint64, path string, timeout time.Duration) <-chan error {
		result := make(chan error, 1)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			meta := cache.NewEntryMetadata(path, streamTestSize, modTime)
			result <- streams.Wait(ctx, reader, path, meta, read)
		}()
		return result
	}

	first, second := wait(1, "/a", 5*time.Second), wait(2, "/b", 5*time.Second)
	waitFor(t, func() bool { return source.started.Load() == 2 })
	third := wait(3, "/c", 3*time.Second)
	waitFor(t, func() bool {
		streams.mu.Lock()
		defer streams.mu.Unlock()
		return streams.demandQueued == 1
	})

	// Both first reads get their bytes; their streams keep their slots.
	source.allow.Store(256 << 10)
	for _, result := range []<-chan error{first, second} {
		if err := <-result; err != nil {
			t.Fatalf("first reads: %v", err)
		}
	}
	if err := <-third; err != nil {
		t.Fatalf("queued read never got a slot from the satisfied streams: %v", err)
	}
}

// A download without a validator that drops is retried from the first byte.
// When the retry returns other bytes, a handle that read the first response
// must not go on to read the second one: the read fails as a changed version.
func TestRangeStreamsRetriedUnvalidatedDownloadDoesNotMixResponses(t *testing.T) {
	const size = 2 << 20
	source := &changingWholeSource{
		first:  bytes.Repeat([]byte{1}, size),
		later:  bytes.Repeat([]byte{2}, size),
		failAt: 512 << 10,
	}
	store, streams := newUnvalidatedTestStreams(t, source)
	meta := cache.NewEntryMetadata("/small", size, time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC))

	head := readThrough(t, store, streams, meta, 0, 4096)
	if head[0] != 1 {
		t.Fatalf("head byte = %d, want the first response", head[0])
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := streams.Wait(ctx, 1, meta.Path, meta, cache.ByteRange{Start: size - 4096, End: size})
	if !errors.Is(err, errRangeVersionChanged) {
		tail := make([]byte, 4096)
		n, _ := store.ReadRange(meta.Path, meta, tail, size-4096)
		t.Fatalf("tail read after the retry = %v with %d bytes (first byte %d), want a changed version", err, n, tail[0])
	}
}

// The same retry with the same bytes carries on, so a proxied download that
// drops still finishes without an error.
func TestRangeStreamsRetriedUnvalidatedDownloadResumesWhenBytesMatch(t *testing.T) {
	const size = 2 << 20
	payload := deterministicRangeTestPayload(size)
	source := &changingWholeSource{first: payload, later: payload, failAt: 512 << 10}
	store, streams := newUnvalidatedTestStreams(t, source)
	meta := cache.NewEntryMetadata("/small", size, time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC))

	readThrough(t, store, streams, meta, 0, 4096)
	if tail := readThrough(t, store, streams, meta, size-4096, 4096); !bytes.Equal(tail, payload[size-4096:]) {
		t.Fatal("tail read after the retry returned the wrong bytes")
	}
	if got := source.requests.Load(); got != 2 {
		t.Fatalf("requests = %d, want the dropped download and its retry", got)
	}
}

// A complete response that goes on past the listed size is a larger file.
// Its bytes must not complete the entry.
func TestRangeStreamsLongerCompleteResponseDoesNotCompleteEntry(t *testing.T) {
	const size = 256 << 10
	source := &longCompleteSource{payload: bytes.Repeat([]byte{7}, size+1)}
	store, streams := newUnvalidatedTestStreams(t, source)
	meta := cache.NewEntryMetadata("/grown", size, time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := streams.Wait(ctx, 1, meta.Path, meta, cache.ByteRange{End: size}); !errors.Is(err, errRangeVersionChanged) {
		t.Fatalf("read of a response longer than the file = %v, want a changed version", err)
	}
	if complete, _ := store.RangeEntryComplete(meta.Path, meta); complete {
		t.Fatal("a response longer than the file completed the cache entry")
	}
}

func newUnvalidatedTestStreams(t *testing.T, source rangeSource) (*disk.DiskCache, *rangeStreams) {
	t.Helper()
	store, err := disk.NewDiskCache(t.TempDir(), disk.WithCapacityBytes(8<<30))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(store.StopMaintenance)
	streams := newRangeStreams(store, source, testStreamPolicy())
	t.Cleanup(streams.Close)
	return store, streams
}

func readThrough(t *testing.T, store *disk.DiskCache, streams *rangeStreams, meta cache.EntryMetadata, offset, length int64) []byte {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := streams.Wait(ctx, 1, meta.Path, meta, cache.ByteRange{Start: offset, End: offset + length}); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, length)
	if n, err := store.ReadRange(meta.Path, meta, buffer, offset); err != nil || int64(n) != length {
		t.Fatalf("cache read = %d, %v", n, err)
	}
	return buffer
}

func TestRangeStreamsReadBehindAStreamStartsItsOwnRequest(t *testing.T) {
	source := newFakeStreamSource(streamTestSize)
	source.bytesPerSecond = 16 << 20
	h := newStreamHarness(t, source, testStreamPolicy())
	h.mustRead(1, 0, 1<<20)
	h.mustRead(1, 1<<20, 1<<20)
	waitFor(t, func() bool { return h.cachedBytes() >= 5<<20 })
	// Bytes behind the stream disappear, as when a listing drops the entry.
	h.store.Delete(h.meta.Path)

	started := time.Now()
	h.mustRead(1, 3<<20, 4096)
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("read took %v waiting on a stream that had already passed it", elapsed)
	}
}

func TestRangeStreamsUnsatisfiableRangeReportsVersionChange(t *testing.T) {
	h := newStreamHarness(t, newFakeStreamSource(streamTestSize), testStreamPolicy())
	var changed atomic.Int32
	h.streams.versionChanged = func(string) { changed.Add(1) }
	h.streams.source = rangeSourceFunc(func(context.Context, string, cache.EntryMetadata, cache.ByteRange) (remoteRangeResponse, error) {
		return remoteRangeResponse{}, fmt.Errorf("%w: bytes */1024", errRangeNotSatisfiable)
	})

	if err := h.read(1, 8<<20, 4096); !errors.Is(err, errRangeVersionChanged) {
		t.Fatalf("read of a range the server cannot satisfy = %v, want %v", err, errRangeVersionChanged)
	}
	if changed.Load() == 0 {
		t.Fatal("the shrunken file was not reported, so its listing would stay stale")
	}
}

// streamHarness drives rangeStreams the way RemoteFs.Read does, against a
// real disk cache.
type streamHarness struct {
	t       *testing.T
	store   *disk.DiskCache
	source  *fakeStreamSource
	streams *rangeStreams
	meta    cache.EntryMetadata
}

func testStreamPolicy() streamPolicy {
	policy := defaultStreamPolicy(true, 0)
	policy.readerGrace = 50 * time.Millisecond
	policy.pauseTimeout = 5 * time.Second
	policy.retryInitial = 5 * time.Millisecond
	policy.retryMax = 20 * time.Millisecond
	policy.retryGiveUp = time.Second
	return policy
}

func newStreamHarness(t *testing.T, source *fakeStreamSource, policy streamPolicy) *streamHarness {
	t.Helper()
	store, err := disk.NewDiskCache(t.TempDir(), disk.WithCapacityBytes(8<<30))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(store.StopMaintenance)
	streams := newRangeStreams(store, source, policy)
	t.Cleanup(streams.Close)
	return &streamHarness{
		t:       t,
		store:   store,
		source:  source,
		streams: streams,
		meta:    cache.NewEntryMetadata("/file", source.size, time.Date(2026, time.September, 1, 0, 0, 0, 0, time.UTC)),
	}
}

func (h *streamHarness) addFile(path string) cache.EntryMetadata {
	return cache.NewEntryMetadata(path, h.meta.Size, h.meta.ModTime)
}

func (h *streamHarness) read(reader uint64, offset, length int64) error {
	return h.readFile(h.meta, reader, offset, length)
}

func (h *streamHarness) readFile(meta cache.EntryMetadata, reader uint64, offset, length int64) error {
	requested := cache.ByteRange{Start: offset, End: min(offset+length, meta.Size)}
	buffer := make([]byte, requested.End-requested.Start)
	if n, _ := h.store.ReadRange(meta.Path, meta, buffer, offset); n == len(buffer) {
		h.streams.NoteRead(reader, meta.Path, meta, requested)
		return h.source.verify(buffer, offset)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.streams.Wait(ctx, reader, meta.Path, meta, requested); err != nil {
		return err
	}
	n, err := h.store.ReadRange(meta.Path, meta, buffer, offset)
	if err != nil || n != len(buffer) {
		return fmt.Errorf("cache read after wait = %d, %v; want %d bytes", n, err, len(buffer))
	}
	return h.source.verify(buffer, offset)
}

func (h *streamHarness) mustRead(reader uint64, offset, length int64) {
	h.t.Helper()
	h.mustReadFile(h.meta, reader, offset, length)
}

func (h *streamHarness) mustReadFile(meta cache.EntryMetadata, reader uint64, offset, length int64) {
	h.t.Helper()
	if err := h.readFile(meta, reader, offset, length); err != nil {
		h.t.Fatalf("read %s [%d, %d): %v", meta.Path, offset, offset+length, err)
	}
}

func (h *streamHarness) waitIdle() {
	h.t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for h.streams.diagnosticsSnapshot().ActiveRequests != 0 || h.streams.hasStreams() {
		if time.Now().After(deadline) {
			h.t.Fatalf("downloads still active: %v", h.source.requestLog())
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func (h *streamHarness) mustBeComplete() {
	h.t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for {
		complete, err := h.store.RangeEntryComplete(h.meta.Path, h.meta)
		if err != nil {
			h.t.Fatal(err)
		}
		if complete {
			break
		}
		if time.Now().After(deadline) {
			h.t.Fatal("file never became complete")
		}
		time.Sleep(5 * time.Millisecond)
	}
	buffer := make([]byte, h.meta.Size)
	if n, err := h.store.ReadRange(h.meta.Path, h.meta, buffer, 0); err != nil || int64(n) != h.meta.Size {
		h.t.Fatalf("complete read = %d, %v", n, err)
	}
	if err := h.source.verify(buffer, 0); err != nil {
		h.t.Fatal(err)
	}
}

func (h *streamHarness) cachedBytes() int64 {
	missing, err := h.store.MissingRanges(h.meta.Path, h.meta, cache.ByteRange{Start: 0, End: h.meta.Size})
	if err != nil {
		h.t.Fatal(err)
	}
	return h.meta.Size - byteRangeCount(missing)
}

// hasActiveStreams reports whether any stream is still downloading or about to.
func (s *rangeStreams) hasActiveStreams() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, file := range s.files {
		for _, stream := range file.streams {
			if stream.active() {
				return true
			}
		}
	}
	return false
}

// streamCountFor reports how many streams path has.
func (s *rangeStreams) streamCountFor(path string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if file := s.files[path]; file != nil {
		return len(file.streams)
	}
	return 0
}

// holdsPausedReadAhead reports whether path has a range stream holding a
// download slot that no read is waiting on, paused ahead of its reader with
// its request still open.
func (s *rangeStreams) holdsPausedReadAhead(path string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if file := s.files[path]; file != nil {
		for _, stream := range file.streams {
			if stream.kind == streamRange && stream.active() && stream.hasSlot && stream.waiters == 0 && stream.paused {
				return true
			}
		}
	}
	return false
}

func (s *rangeStreams) hasStreams() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, file := range s.files {
		if len(file.streams) > 0 {
			return true
		}
	}
	return false
}

// fakeStreamSource serves one payload the way the SDK transport returns
// validated responses. Bodies only produce bytes as they are read.
type fakeStreamSource struct {
	size           int64
	payload        []byte
	bytesPerSecond int64
	ignoreRanges   bool
	rejectNext     atomic.Int32
	transferred    atomic.Int64

	mu       sync.Mutex
	etag     string
	requests []fakeStreamRequest
	failAt   map[int64]int64
}

type fakeStreamRequest struct {
	requested cache.ByteRange
	complete  bool
	ifMatch   string
}

func (r fakeStreamRequest) String() string {
	return fmt.Sprintf("{[%d, %d) complete=%t if-match=%q}", r.requested.Start, r.requested.End, r.complete, r.ifMatch)
}

func newFakeStreamSource(size int64) *fakeStreamSource {
	payload := make([]byte, size)
	for i := range payload {
		payload[i] = byte(i*31 + i/65536)
	}
	return &fakeStreamSource{size: size, payload: payload, etag: `"v1"`, failAt: make(map[int64]int64)}
}

func (f *fakeStreamSource) DownloadRange(ctx context.Context, _ string, meta cache.EntryMetadata, requested cache.ByteRange) (remoteRangeResponse, error) {
	if f.ignoreRanges {
		return f.respond(ctx, meta, cache.ByteRange{Start: 0, End: f.size}, false)
	}
	return f.respond(ctx, meta, requested, true)
}

func (f *fakeStreamSource) DownloadComplete(ctx context.Context, _ string, meta cache.EntryMetadata, _ cache.ByteRange) (remoteRangeResponse, error) {
	return f.respond(ctx, meta, cache.ByteRange{Start: 0, End: f.size}, false)
}

func (f *fakeStreamSource) respond(ctx context.Context, meta cache.EntryMetadata, returned cache.ByteRange, partial bool) (remoteRangeResponse, error) {
	f.mu.Lock()
	f.requests = append(f.requests, fakeStreamRequest{requested: returned, complete: !partial, ifMatch: meta.ETag})
	etag := f.etag
	failAfter, fail := f.failAt[returned.Start]
	delete(f.failAt, returned.Start)
	f.mu.Unlock()

	if f.rejectNext.Add(-1) >= 0 {
		return remoteRangeResponse{}, files_sdk.ResponseError{HttpCode: 400, Title: "rejected"}
	}
	if meta.ETag != "" && meta.ETag != etag {
		return remoteRangeResponse{}, errRangeVersionChanged
	}
	body := &fakeStreamBody{ctx: ctx, source: f, pos: returned.Start, end: returned.End, failAt: -1}
	if fail {
		body.failAt = returned.Start + failAfter
	}
	return remoteRangeResponse{ETag: etag, Returned: returned, TotalSize: f.size, Partial: partial, Body: body}, nil
}

func (f *fakeStreamSource) failOnce(start, after int64) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.failAt[start] = after
}

func (f *fakeStreamSource) setETag(etag string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.etag = etag
}

func (f *fakeStreamSource) requestLog() []fakeStreamRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]fakeStreamRequest(nil), f.requests...)
}

func (f *fakeStreamSource) requestCount() int {
	return len(f.requestLog())
}

// waitStable returns the transferred byte count once it stops changing.
func (f *fakeStreamSource) waitStable(t *testing.T) int64 {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	last := f.transferred.Load()
	stableSince := time.Now()
	for time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
		if current := f.transferred.Load(); current != last {
			last, stableSince = current, time.Now()
			continue
		}
		if time.Since(stableSince) >= 150*time.Millisecond {
			return last
		}
	}
	t.Fatalf("transfer never settled; %d bytes so far", last)
	return last
}

func (f *fakeStreamSource) verify(data []byte, offset int64) error {
	if !bytes.Equal(data, f.payload[offset:offset+int64(len(data))]) {
		return fmt.Errorf("bytes at [%d, %d) do not match the source", offset, offset+int64(len(data)))
	}
	return nil
}

type fakeStreamBody struct {
	ctx    context.Context
	source *fakeStreamSource
	pos    int64
	end    int64
	failAt int64
}

func (b *fakeStreamBody) Read(buffer []byte) (int, error) {
	if err := b.ctx.Err(); err != nil {
		return 0, err
	}
	if b.pos >= b.end {
		return 0, io.EOF
	}
	n := min(int64(len(buffer)), b.end-b.pos, 64<<10)
	if b.failAt >= 0 {
		if b.pos >= b.failAt {
			return 0, errors.New("connection reset by peer")
		}
		n = min(n, b.failAt-b.pos)
	}
	if rate := b.source.bytesPerSecond; rate > 0 {
		select {
		case <-b.ctx.Done():
			return 0, b.ctx.Err()
		case <-time.After(time.Duration(float64(n) / float64(rate) * float64(time.Second))):
		}
	}
	copy(buffer, b.source.payload[b.pos:b.pos+n])
	b.pos += n
	b.source.transferred.Add(n)
	return int(n), nil
}

func (b *fakeStreamBody) Close() error { return nil }

// gatedStreamSource releases at most allow bytes of every response.
type gatedStreamSource struct {
	payload []byte
	allow   atomic.Int64
	started atomic.Int32
}

func (g *gatedStreamSource) DownloadRange(ctx context.Context, _ string, _ cache.EntryMetadata, requested cache.ByteRange) (remoteRangeResponse, error) {
	g.started.Add(1)
	body := &gatedStreamBody{ctx: ctx, source: g, start: requested.Start, pos: requested.Start, end: requested.End}
	return remoteRangeResponse{ETag: `"v1"`, Returned: requested, TotalSize: int64(len(g.payload)), Partial: true, Body: body}, nil
}

type gatedStreamBody struct {
	ctx             context.Context
	source          *gatedStreamSource
	start, pos, end int64
}

func (b *gatedStreamBody) Read(buffer []byte) (int, error) {
	for {
		if err := b.ctx.Err(); err != nil {
			return 0, err
		}
		if b.pos >= b.end {
			return 0, io.EOF
		}
		if allowed := b.start + b.source.allow.Load(); b.pos < allowed {
			n := min(int64(len(buffer)), b.end-b.pos, allowed-b.pos, 64<<10)
			copy(buffer, b.source.payload[b.pos:b.pos+n])
			b.pos += n
			return int(n), nil
		}
		time.Sleep(time.Millisecond)
	}
}

func (b *gatedStreamBody) Close() error { return nil }

// longCompleteSource answers with the whole file and no length, as a chunked
// response does, so the listed size is trusted, and sends every byte of
// payload.
type longCompleteSource struct{ payload []byte }

func (l *longCompleteSource) DownloadRange(ctx context.Context, path string, meta cache.EntryMetadata, requested cache.ByteRange) (remoteRangeResponse, error) {
	return l.DownloadComplete(ctx, path, meta, requested)
}

func (l *longCompleteSource) DownloadComplete(_ context.Context, _ string, meta cache.EntryMetadata, _ cache.ByteRange) (remoteRangeResponse, error) {
	body := newExactLengthReadCloser(io.NopCloser(bytes.NewReader(l.payload)), meta.Size)
	return remoteRangeResponse{Returned: cache.ByteRange{End: meta.Size}, TotalSize: meta.Size, Body: body}, nil
}

// changingWholeSource answers with the whole file and no validator. The first
// response drops after failAt bytes; later ones carry different content of
// the same size, as after an unnoticed remote edit.
type changingWholeSource struct {
	first, later []byte
	failAt       int64
	requests     atomic.Int32
}

func (c *changingWholeSource) DownloadRange(ctx context.Context, path string, meta cache.EntryMetadata, requested cache.ByteRange) (remoteRangeResponse, error) {
	return c.DownloadComplete(ctx, path, meta, requested)
}

func (c *changingWholeSource) DownloadComplete(ctx context.Context, _ string, _ cache.EntryMetadata, _ cache.ByteRange) (remoteRangeResponse, error) {
	payload, failAt := c.later, int64(-1)
	if c.requests.Add(1) == 1 {
		payload, failAt = c.first, c.failAt
	}
	size := int64(len(payload))
	return remoteRangeResponse{Returned: cache.ByteRange{Start: 0, End: size}, TotalSize: size, Body: &changingWholeBody{ctx: ctx, payload: payload, failAt: failAt}}, nil
}

type changingWholeBody struct {
	ctx     context.Context
	payload []byte
	pos     int64
	failAt  int64
}

func (b *changingWholeBody) Read(buffer []byte) (int, error) {
	if err := b.ctx.Err(); err != nil {
		return 0, err
	}
	if b.failAt >= 0 && b.pos >= b.failAt {
		return 0, errors.New("connection reset by peer")
	}
	if b.pos >= int64(len(b.payload)) {
		return 0, io.EOF
	}
	end := min(int64(len(b.payload)), b.pos+int64(len(buffer)), b.pos+64<<10)
	if b.failAt >= 0 {
		end = min(end, b.failAt)
	}
	n := copy(buffer, b.payload[b.pos:end])
	b.pos += int64(n)
	return n, nil
}

func (b *changingWholeBody) Close() error { return nil }
