//go:build linux || windows

package disk_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	fscache "github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
)

func TestDiskCacheRangeWritesAndHoleReads(t *testing.T) {
	cacheStore := newRangeTestDiskCache(t, t.TempDir())
	path := "/ranges.bin"
	meta := rangeTestMetadata(path, 100)

	writes := []struct {
		offset int64
		data   []byte
	}{
		{offset: 80, data: bytes.Repeat([]byte("t"), 20)},
		{offset: 0, data: bytes.Repeat([]byte("b"), 10)},
		{offset: 40, data: bytes.Repeat([]byte("m"), 10)},
	}
	for _, write := range writes {
		n, err := cacheStore.WriteRange(path, meta, write.data, write.offset)
		if err != nil || n != len(write.data) {
			t.Fatalf("WriteRange(%d) = %d, %v; want %d, nil", write.offset, n, err, len(write.data))
		}
	}

	for _, write := range writes {
		buffer := make([]byte, len(write.data))
		n, err := cacheStore.ReadRange(path, meta, buffer, write.offset)
		if err != nil || n != len(write.data) || !bytes.Equal(buffer, write.data) {
			t.Fatalf("ReadRange(%d) = %d, %v, %q; want %q", write.offset, n, err, buffer, write.data)
		}
	}

	holeBuffer := bytes.Repeat([]byte{0xff}, 40)
	if n, err := cacheStore.ReadRange(path, meta, holeBuffer, 5); err != nil || n != 0 {
		t.Fatalf("ReadRange across hole = %d, %v; want cache miss", n, err)
	}

	missing, err := cacheStore.MissingRanges(path, meta, fscache.ByteRange{Start: 0, End: 100})
	if err != nil {
		t.Fatalf("MissingRanges failed: %v", err)
	}
	wantMissing := []fscache.ByteRange{{Start: 10, End: 40}, {Start: 50, End: 80}}
	if !reflect.DeepEqual(missing, wantMissing) {
		t.Errorf("MissingRanges = %#v, want %#v", missing, wantMissing)
	}
	complete, err := cacheStore.RangeEntryComplete(path, meta)
	if err != nil {
		t.Fatalf("RangeEntryComplete failed: %v", err)
	}
	if complete {
		t.Fatal("RangeEntryComplete = true with holes")
	}
	for _, gap := range wantMissing {
		if _, err := cacheStore.WriteRange(path, meta, bytes.Repeat([]byte("x"), int(gap.End-gap.Start)), gap.Start); err != nil {
			t.Fatalf("WriteRange missing gap %v failed: %v", gap, err)
		}
	}
	complete, err = cacheStore.RangeEntryComplete(path, meta)
	if err != nil {
		t.Fatalf("RangeEntryComplete after filling gaps failed: %v", err)
	}
	if !complete {
		t.Fatal("RangeEntryComplete = false after filling all gaps")
	}
}

func TestDiskCacheWriteCompleteRangePublishesCompleteEntry(t *testing.T) {
	cacheStore := newRangeTestDiskCache(t, t.TempDir())
	path := "/complete.bin"
	payload := []byte("complete response")
	meta := fscache.NewEntryMetadata(path, int64(len(payload)), time.Now().Round(0))

	split := len(payload) / 2
	if n, err := cacheStore.WriteCompleteRange(path, meta, payload[:split], 0); err != nil || n != split {
		t.Fatalf("first WriteCompleteRange = %d, %v; want %d, nil", n, err, split)
	}
	if n, err := cacheStore.WriteCompleteRange(path, meta, payload[split:], int64(split)); err != nil || n != len(payload)-split {
		t.Fatalf("second WriteCompleteRange = %d, %v; want %d, nil", n, err, len(payload)-split)
	}

	buffer := make([]byte, len(payload))
	if n, err := cacheStore.ReadRange(path, meta, buffer, 0); err != nil || n != len(payload) {
		t.Fatalf("ReadRange = %d, %v; want %d, nil", n, err, len(payload))
	}
	if !bytes.Equal(buffer, payload) {
		t.Fatalf("cached data = %q, want %q", buffer, payload)
	}
	if complete, err := cacheStore.RangeEntryComplete(path, meta); err != nil || !complete {
		t.Fatalf("RangeEntryComplete = %t, %v; want true, nil", complete, err)
	}
}

func TestDiskCacheCompleteRangeBatchesDurabilityWithoutDelayingReads(t *testing.T) {
	const (
		publishSize = rangeTestFlushBytes
		chunkSize   = 128 * 1024
	)
	path := "/batched-complete.bin"
	payload := bytes.Repeat([]byte("b"), publishSize+2*chunkSize)
	meta := rangeTestMetadata(path, int64(len(payload)))

	volatileRoot := t.TempDir()
	cacheStore := newRangeTestDiskCache(t, volatileRoot)
	if _, err := cacheStore.WriteCompleteRange(path, meta, payload[:chunkSize], 0); err != nil {
		t.Fatalf("first WriteCompleteRange failed: %v", err)
	}
	first := make([]byte, chunkSize)
	if n, err := cacheStore.ReadRange(path, meta, first, 0); err != nil || n != len(first) || !bytes.Equal(first, payload[:chunkSize]) {
		t.Fatalf("active ReadRange = %d, %v, correct=%t; want %d, nil, true", n, err, bytes.Equal(first, payload[:chunkSize]), len(first))
	}

	restarted := newRangeTestDiskCache(t, volatileRoot)
	missing, err := restarted.MissingRanges(path, meta, fscache.ByteRange{Start: 0, End: chunkSize})
	if err != nil || !reflect.DeepEqual(missing, []fscache.ByteRange{{Start: 0, End: chunkSize}}) {
		t.Fatalf("MissingRanges after restart = %#v, %v; want the unflushed range missing", missing, err)
	}
	if _, err := cacheStore.WriteCompleteRange(path, meta, payload[:chunkSize], 0); err != nil {
		t.Fatalf("WriteCompleteRange before Clear failed: %v", err)
	}
	if err := cacheStore.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}
	missing, err = cacheStore.MissingRanges(path, meta, fscache.ByteRange{Start: 0, End: chunkSize})
	if err != nil || !reflect.DeepEqual(missing, []fscache.ByteRange{{Start: 0, End: chunkSize}}) {
		t.Fatalf("cleared MissingRanges = %#v, %v; want the volatile range missing", missing, err)
	}

	cacheRoot := t.TempDir()
	cacheStore = newRangeTestDiskCache(t, cacheRoot)
	for offset := 0; offset < publishSize; offset += chunkSize {
		if _, err := cacheStore.WriteCompleteRange(path, meta, payload[offset:offset+chunkSize], int64(offset)); err != nil {
			t.Fatalf("WriteCompleteRange at %d failed: %v", offset, err)
		}
	}
	reopenedAfterBatch := newRangeTestDiskCache(t, cacheRoot)
	durable := make([]byte, publishSize)
	if n, err := reopenedAfterBatch.ReadRange(path, meta, durable, 0); err != nil || n != len(durable) || !bytes.Equal(durable, payload[:publishSize]) {
		t.Fatalf("durable ReadRange = %d, %v, correct=%t; want %d, nil, true", n, err, bytes.Equal(durable, payload[:publishSize]), len(durable))
	}

	if _, err := cacheStore.WriteCompleteRange(path, meta, payload[publishSize:publishSize+chunkSize], publishSize); err != nil {
		t.Fatalf("volatile WriteCompleteRange failed: %v", err)
	}
	restarted = newRangeTestDiskCache(t, cacheRoot)
	missing, err = restarted.MissingRanges(path, meta, fscache.ByteRange{Start: 0, End: publishSize + chunkSize})
	if err != nil || !reflect.DeepEqual(missing, []fscache.ByteRange{{Start: publishSize, End: publishSize + chunkSize}}) {
		t.Fatalf("MissingRanges after restart = %#v, %v; want only the unflushed suffix missing", missing, err)
	}
	if _, err := cacheStore.WriteCompleteRange(path, meta, payload[publishSize:], publishSize); err != nil {
		t.Fatalf("retried final WriteCompleteRange failed: %v", err)
	}
	complete := make([]byte, len(payload))
	if n, err := cacheStore.ReadRange(path, meta, complete, 0); err != nil || n != len(complete) || !bytes.Equal(complete, payload) {
		t.Fatalf("complete ReadRange = %d, %v, correct=%t; want %d, nil, true", n, err, bytes.Equal(complete, payload), len(complete))
	}
	reopenedComplete := newRangeTestDiskCache(t, cacheRoot)
	if n, err := reopenedComplete.ReadRange(path, meta, complete, 0); err != nil || n != len(complete) || !bytes.Equal(complete, payload) {
		t.Fatalf("reopened complete ReadRange = %d, %v, correct=%t; want %d, nil, true", n, err, bytes.Equal(complete, payload), len(complete))
	}
}

func TestDiskCacheEntryRewrittenAfterClearSurvivesUnpin(t *testing.T) {
	cacheStore := newRangeTestDiskCache(t, t.TempDir())
	path := "/cleared-while-open.bin"
	meta := rangeTestMetadata(path, 8)
	meta.ETag = `"stable"`
	cacheStore.Pin(path)
	if _, err := cacheStore.WriteRange(path, meta, []byte("old!"), 0); err != nil {
		t.Fatal(err)
	}
	if err := cacheStore.Clear(); err != nil {
		t.Fatal(err)
	}
	if _, err := cacheStore.WriteRange(path, meta, []byte("new!"), 0); err != nil {
		t.Fatal(err)
	}
	cacheStore.Unpin(path)

	buffer := make([]byte, 4)
	if n, err := cacheStore.ReadRange(path, meta, buffer, 0); err != nil || n != 4 || string(buffer) != "new!" {
		t.Fatalf("ReadRange after unpin = %d, %v, %q; want the bytes written after the clear", n, err, buffer)
	}
}

// A clear lets downloads resume once it has invalidated the cache, before it
// deletes the old files. What they write in the meantime is kept.
func TestDiskCacheClearKeepsEntriesWrittenBeforeItsFilesAreDeleted(t *testing.T) {
	cacheStore := newRangeTestDiskCache(t, t.TempDir())
	stale := "/stale.bin"
	staleMeta := rangeTestMetadata(stale, 4)
	redownloaded := "/redownloaded.bin"
	redownloadedMeta := rangeTestMetadata(redownloaded, 4)
	for _, entry := range []struct {
		path string
		meta fscache.EntryMetadata
	}{{stale, staleMeta}, {redownloaded, redownloadedMeta}} {
		if _, err := cacheStore.WriteCompleteRange(entry.path, entry.meta, []byte("old!"), 0); err != nil {
			t.Fatal(err)
		}
	}

	invalidated, err := cacheStore.Invalidate()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cacheStore.WriteCompleteRange(redownloaded, redownloadedMeta, []byte("new!"), 0); err != nil {
		t.Fatal(err)
	}
	open := "/open.bin"
	openMeta := rangeTestMetadata(open, 8)
	openMeta.ETag = `"stable"`
	cacheStore.Pin(open)
	if _, err := cacheStore.WriteRange(open, openMeta, []byte("head"), 0); err != nil {
		t.Fatal(err)
	}
	if err := invalidated.Remove(); err != nil {
		t.Fatal(err)
	}
	cacheStore.Unpin(open)

	buffer := make([]byte, 4)
	if n, err := cacheStore.ReadRange(stale, staleMeta, buffer, 0); err != nil || n != 0 {
		t.Fatalf("ReadRange of a cleared entry = %d, %v; want a miss", n, err)
	}
	if n, err := cacheStore.ReadRange(redownloaded, redownloadedMeta, buffer, 0); err != nil || n != 4 || string(buffer) != "new!" {
		t.Fatalf("ReadRange of an entry downloaded again during the clear = %d, %v, %q; want the new bytes", n, err, buffer)
	}
	if n, err := cacheStore.ReadRange(open, openMeta, buffer, 0); err != nil || n != 4 || string(buffer) != "head" {
		t.Fatalf("ReadRange of an open entry written during the clear = %d, %v, %q; want its bytes", n, err, buffer)
	}
}

func TestDiskCacheRangeOverlapAccounting(t *testing.T) {
	cacheStore := newRangeTestDiskCache(t, t.TempDir())
	path := "/overlap.bin"
	meta := rangeTestMetadata(path, 100)

	if _, err := cacheStore.WriteRange(path, meta, bytes.Repeat([]byte("a"), 10), 10); err != nil {
		t.Fatalf("first WriteRange failed: %v", err)
	}
	if _, err := cacheStore.WriteRange(path, meta, bytes.Repeat([]byte("b"), 10), 15); err != nil {
		t.Fatalf("overlapping WriteRange failed: %v", err)
	}
	if _, err := cacheStore.WriteRange(path, meta, bytes.Repeat([]byte("c"), 15), 10); err != nil {
		t.Fatalf("covered WriteRange failed: %v", err)
	}

	stats := cacheStore.Stats()
	if got := stats.SizeBytes.Load(); got != 15 {
		t.Errorf("SizeBytes = %d, want 15", got)
	}
	if got := stats.FileCount.Load(); got != 1 {
		t.Errorf("FileCount = %d, want 1", got)
	}
	if got := stats.WriteBytes.Load(); got != 15 {
		t.Errorf("WriteBytes = %d, want 15 newly covered bytes", got)
	}

	buffer := make([]byte, 15)
	n, err := cacheStore.ReadRange(path, meta, buffer, 10)
	if err != nil || n != len(buffer) {
		t.Fatalf("ReadRange = %d, %v; want %d, nil", n, err, len(buffer))
	}
	want := append(bytes.Repeat([]byte("a"), 10), bytes.Repeat([]byte("b"), 5)...)
	if !bytes.Equal(buffer, want) {
		t.Errorf("overlap data = %q, want %q", buffer, want)
	}
}

func TestDiskCacheRangeEntrySurvivesReconstruction(t *testing.T) {
	root := t.TempDir()
	path := "/restart.bin"
	meta := rangeTestMetadata(path, 1<<20)
	offset := meta.Size - 16
	payload := []byte("tail survives!!!")

	first := newRangeTestDiskCache(t, root)
	if _, err := first.WriteRange(path, meta, payload, offset); err != nil {
		t.Fatalf("WriteRange failed: %v", err)
	}

	if err := first.FlushRanges(path, meta); err != nil {
		t.Fatal(err)
	}
	second := newRangeTestDiskCache(t, root)
	if got := second.SizeBytes(); got != int64(len(payload)) {
		t.Fatalf("reconstructed SizeBytes = %d, want %d", got, len(payload))
	}
	if got := second.Stats().FileCount.Load(); got != 1 {
		t.Fatalf("reconstructed FileCount = %d, want 1", got)
	}
	buffer := make([]byte, len(payload))
	n, err := second.ReadRange(path, meta, buffer, offset)
	if err != nil || n != len(payload) || !bytes.Equal(buffer, payload) {
		t.Fatalf("reconstructed ReadRange = %d, %v, %q; want %q", n, err, buffer, payload)
	}
	if n, err := second.ReadRange(path, meta, make([]byte, 1), 0); err != nil || n != 0 {
		t.Fatalf("uncached prefix ReadRange = %d, %v; want cache miss", n, err)
	}
}

func TestDiskCacheReopenAccountsEntriesWrittenDuringStartupScan(t *testing.T) {
	root := t.TempDir()
	path := "/resumed.bin"
	meta := rangeTestMetadata(path, 1<<20)
	meta.ETag = `"stable"`
	first := newRangeTestDiskCache(t, root)
	if _, err := first.WriteRange(path, meta, bytes.Repeat([]byte("a"), 4096), 0); err != nil {
		t.Fatal(err)
	}
	if err := first.FlushRanges(path, meta); err != nil {
		t.Fatal(err)
	}
	first.StopMaintenance()
	abandoned := filepath.Join(root, "data", "00", "data-without-metadata")
	if err := os.MkdirAll(filepath.Dir(abandoned), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abandoned, []byte("orphan"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Whichever of the startup scan and this write reaches the entry first,
	// the entry is counted once, including the bytes it already had.
	second := newRangeTestDiskCache(t, root)
	if _, err := second.WriteRange(path, meta, bytes.Repeat([]byte("b"), 4096), 4096); err != nil {
		t.Fatal(err)
	}
	if got := second.SizeBytes(); got != 8192 {
		t.Fatalf("SizeBytes after reopening and writing = %d, want 8192", got)
	}
	if got := second.Stats().FileCount.Load(); got != 1 {
		t.Fatalf("FileCount = %d, want 1", got)
	}
	if _, err := os.Stat(abandoned); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("data without metadata survived the startup scan: %v", err)
	}
}

func TestDiskCacheRangeEvictionRemovesDataAndMetadata(t *testing.T) {
	root := t.TempDir()
	cacheStore, err := disk.NewDiskCache(root, disk.WithCapacityBytes(10))
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}
	firstPath := "/first.bin"
	secondPath := "/second.bin"
	firstMeta := rangeTestMetadata(firstPath, 10)
	secondMeta := rangeTestMetadata(secondPath, 10)
	if _, err := cacheStore.WriteRange(firstPath, firstMeta, bytes.Repeat([]byte("a"), 10), 0); err != nil {
		t.Fatalf("first WriteRange failed: %v", err)
	}
	if _, err := cacheStore.WriteRange(secondPath, secondMeta, bytes.Repeat([]byte("b"), 10), 0); err != nil {
		t.Fatalf("second WriteRange failed: %v", err)
	}

	if n, err := cacheStore.ReadRange(firstPath, firstMeta, make([]byte, 10), 0); err != nil || n != 0 {
		t.Fatalf("evicted ReadRange = %d, %v; want cache miss", n, err)
	}
	if n, err := cacheStore.ReadRange(secondPath, secondMeta, make([]byte, 10), 0); err != nil || n != 10 {
		t.Fatalf("retained ReadRange = %d, %v; want 10, nil", n, err)
	}
	metadataFiles, err := filepath.Glob(filepath.Join(root, "state", "metadata", "*", "*.json"))
	if err != nil {
		t.Fatalf("Glob metadata failed: %v", err)
	}
	if len(metadataFiles) != 1 {
		t.Fatalf("metadata files = %v, want one retained entry", metadataFiles)
	}
	if got := cacheStore.SizeBytes(); got != 10 {
		t.Fatalf("SizeBytes = %d, want 10", got)
	}
}

func TestDiskCacheClearDefersPinnedRangeRemoval(t *testing.T) {
	cacheStore := newRangeTestDiskCache(t, t.TempDir())
	path := "/pinned-range.bin"
	meta := rangeTestMetadata(path, 100)
	payload := []byte("tail")
	if _, err := cacheStore.WriteRange(path, meta, payload, 96); err != nil {
		t.Fatalf("WriteRange failed: %v", err)
	}
	cacheStore.Pin(path)

	if err := cacheStore.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}
	if n, err := cacheStore.ReadRange(path, meta, make([]byte, len(payload)), 96); err != nil || n != 0 {
		t.Fatalf("ReadRange after Clear = %d, %v; want cache miss", n, err)
	}
	if got := cacheStore.SizeBytes(); got != int64(len(payload)) {
		t.Fatalf("SizeBytes while pinned = %d, want %d", got, len(payload))
	}

	cacheStore.Unpin(path)
	if got := cacheStore.SizeBytes(); got != 0 {
		t.Fatalf("SizeBytes after Unpin = %d, want 0", got)
	}
}

func TestDiskCacheMaintenanceExpiresIncompleteAndCompleteRanges(t *testing.T) {
	root := t.TempDir()
	cacheStore, err := disk.NewDiskCache(
		root,
		disk.WithMaxAge(10*time.Millisecond),
		disk.WithMaintenanceInterval(5*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}
	incompletePath := "/incomplete-expired.bin"
	completePath := "/complete-expired.bin"
	incompleteMeta := rangeTestMetadata(incompletePath, 10)
	completeMeta := rangeTestMetadata(completePath, 10)
	if _, err := cacheStore.WriteRange(incompletePath, incompleteMeta, []byte("tail"), 6); err != nil {
		t.Fatalf("incomplete WriteRange failed: %v", err)
	}
	if _, err := cacheStore.WriteRange(completePath, completeMeta, []byte("0123456789"), 0); err != nil {
		t.Fatalf("complete WriteRange failed: %v", err)
	}

	dataFiles, err := filepath.Glob(filepath.Join(root, "data", "*", "*"))
	if err != nil {
		t.Fatalf("Glob data failed: %v", err)
	}
	oldTime := time.Now().Add(-time.Hour)
	for _, dataFile := range dataFiles {
		if err := os.Chtimes(dataFile, oldTime, oldTime); err != nil {
			t.Fatalf("Chtimes(%s) failed: %v", dataFile, err)
		}
	}

	cacheStore.StartMaintenance()
	deadline := time.Now().Add(5 * time.Second)
	for cacheStore.SizeBytes() != 0 {
		if time.Now().After(deadline) {
			cacheStore.StopMaintenance()
			t.Fatalf("expired range entries remain: size=%d", cacheStore.SizeBytes())
		}
		time.Sleep(10 * time.Millisecond)
	}
	cacheStore.StopMaintenance()
	if got := cacheStore.Stats().FileCount.Load(); got != 0 {
		t.Fatalf("FileCount after expiration = %d, want 0", got)
	}
	if got := cacheStore.SizeBytes(); got != 0 {
		t.Fatalf("SizeBytes after expiration = %d, want 0", got)
	}
	metadataFiles, err := filepath.Glob(filepath.Join(root, "state", "metadata", "*", "*.json"))
	if err != nil {
		t.Fatalf("Glob metadata failed: %v", err)
	}
	if len(metadataFiles) != 0 {
		t.Fatalf("expired metadata files remain: %v", metadataFiles)
	}
}

func TestDiskCacheRangeVersionChangeInvalidatesEntry(t *testing.T) {
	cacheStore := newRangeTestDiskCache(t, t.TempDir())
	path := "/versioned.bin"
	oldMeta := rangeTestMetadata(path, 10)
	newMeta := fscache.NewEntryMetadata(path, 10, oldMeta.ModTime.Add(time.Second))
	if _, err := cacheStore.WriteRange(path, oldMeta, bytes.Repeat([]byte("a"), 10), 0); err != nil {
		t.Fatalf("WriteRange failed: %v", err)
	}

	if n, err := cacheStore.ReadRange(path, newMeta, make([]byte, 10), 0); err != nil || n != 0 {
		t.Fatalf("stale ReadRange = %d, %v; want cache miss", n, err)
	}
	if got := cacheStore.SizeBytes(); got != 0 {
		t.Fatalf("SizeBytes after invalidation = %d, want 0", got)
	}
}

func TestDiskCacheRangeVersionChangeReplacesPinnedEntry(t *testing.T) {
	cacheStore := newRangeTestDiskCache(t, t.TempDir())
	path := "/pinned-version.bin"
	oldMeta := rangeTestMetadata(path, 4)
	newMeta := fscache.NewEntryMetadata(path, 4, oldMeta.ModTime.Add(time.Second))
	if _, err := cacheStore.WriteRange(path, oldMeta, []byte("old!"), 0); err != nil {
		t.Fatalf("old WriteRange failed: %v", err)
	}
	cacheStore.Pin(path)

	if _, err := cacheStore.WriteRange(path, newMeta, []byte("new!"), 0); err != nil {
		t.Fatalf("pinned WriteRange failed: %v", err)
	}
	if n, err := cacheStore.ReadRange(path, oldMeta, make([]byte, 4), 0); err != nil || n != 0 {
		t.Fatalf("invalidated old ReadRange = %d, %v; want cache miss", n, err)
	}
	buffer := make([]byte, 4)
	if n, err := cacheStore.ReadRange(path, newMeta, buffer, 0); err != nil || n != 4 || string(buffer) != "new!" {
		t.Fatalf("new ReadRange = %d, %v, %q; want 4, nil, new!", n, err, buffer)
	}
	cacheStore.Unpin(path)
}

func TestDiskCacheFailedFlushDoesNotSurviveRestart(t *testing.T) {
	root := t.TempDir()
	cacheStore := newRangeTestDiskCache(t, root)
	path := "/interrupted.bin"
	meta := rangeTestMetadata(path, 100)

	if _, err := cacheStore.WriteRange(path, meta, []byte("data"), 50); err != nil {
		t.Fatal(err)
	}
	metadataDir := filepath.Join(root, "state", "metadata")
	if err := os.RemoveAll(metadataDir); err != nil {
		t.Fatalf("RemoveAll metadata directory failed: %v", err)
	}
	if err := os.WriteFile(metadataDir, []byte("block metadata directory"), 0o644); err != nil {
		t.Fatalf("WriteFile metadata blocker failed: %v", err)
	}
	if err := cacheStore.FlushRanges(path, meta); err == nil {
		t.Fatal("FlushRanges succeeded with unavailable metadata directory")
	}

	if err := os.Remove(metadataDir); err != nil {
		t.Fatalf("Remove metadata blocker failed: %v", err)
	}
	if err := os.MkdirAll(metadataDir, 0o755); err != nil {
		t.Fatalf("MkdirAll metadata directory failed: %v", err)
	}

	second := newRangeTestDiskCache(t, root)
	if got := second.Stats().FileCount.Load(); got != 0 {
		t.Fatalf("FileCount after reconstruction = %d, want 0", got)
	}
	dataFiles, err := filepath.Glob(filepath.Join(root, "data", "*", "*"))
	if err != nil {
		t.Fatalf("Glob data failed: %v", err)
	}
	if len(dataFiles) != 0 {
		t.Fatalf("orphaned data files remain after reconstruction: %v", dataFiles)
	}
}

func TestDiskCacheRangeRejectsOutOfBoundsWrite(t *testing.T) {
	cacheStore := newRangeTestDiskCache(t, t.TempDir())
	meta := rangeTestMetadata("/bounds.bin", 10)
	_, err := cacheStore.WriteRange(meta.Path, meta, make([]byte, 2), 9)
	if !errors.Is(err, fscache.ErrInvalidByteRange) {
		t.Fatalf("WriteRange error = %v, want ErrInvalidByteRange", err)
	}
}

func TestDiskCacheDisabledMissingRangesStillValidatesRequest(t *testing.T) {
	cacheStore, err := disk.NewDiskCache(t.TempDir(), disk.WithDisabled(true))
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}
	meta := rangeTestMetadata("/disabled.bin", 10)
	_, err = cacheStore.MissingRanges(meta.Path, meta, fscache.ByteRange{Start: 9, End: 11})
	if !errors.Is(err, fscache.ErrInvalidByteRange) {
		t.Fatalf("MissingRanges error = %v, want ErrInvalidByteRange", err)
	}
}

// rangeTestFlushBytes keeps durability batches small enough to cross in tests.
const rangeTestFlushBytes = 4 * 1024 * 1024

func newRangeTestDiskCache(t *testing.T, root string) *disk.DiskCache {
	t.Helper()
	cacheStore, err := disk.NewDiskCache(root, disk.WithRangeFlushBytes(rangeTestFlushBytes))
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}
	// Windows cannot remove the test directory while a data file is open.
	t.Cleanup(cacheStore.StopMaintenance)
	return cacheStore
}

func rangeTestMetadata(path string, size int64) fscache.EntryMetadata {
	return fscache.NewEntryMetadata(path, size, time.Date(2026, time.August, 13, 12, 0, 0, 0, time.UTC))
}

// Durability follows received byte count, even when the first batch is in the
// middle of the file. An interrupted next batch must not advertise its holes.
func TestDiskCacheOutOfOrderBatchesSurviveRestartWithoutVolatileBytes(t *testing.T) {
	const mib = 1024 * 1024
	root := t.TempDir()
	store := newRangeTestDiskCache(t, root)
	meta := rangeTestMetadata("/out-of-order", 12*mib)
	meta.ETag = `"stable"`
	block := bytes.Repeat([]byte("x"), mib)
	for _, offset := range []int64{10 * mib, 4 * mib, 8 * mib, 2 * mib} {
		if _, err := store.WriteRange(meta.Path, meta, block, offset); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := store.WriteRange(meta.Path, meta, block[:128*1024], 0); err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 128*1024)
	if n, err := store.ReadRange(meta.Path, meta, data, 0); n != len(data) || err != nil || !bytes.Equal(data, block[:len(data)]) {
		t.Fatalf("live prefix = %d, %v", n, err)
	}
	store = newRangeTestDiskCache(t, root)
	if n, err := store.ReadRange(meta.Path, meta, data, 0); n != 0 || err != nil {
		t.Fatalf("unpublished prefix survived: %d, %v", n, err)
	}
	for _, offset := range []int64{10 * mib, 4 * mib, 8 * mib, 2 * mib} {
		data := make([]byte, mib)
		if n, err := store.ReadRange(meta.Path, meta, data, offset); n != mib || err != nil || !bytes.Equal(data, block) {
			t.Fatalf("durable range at %d = %d, %v", offset, n, err)
		}
	}
	if got := store.SizeBytes(); got != 4*mib {
		t.Fatalf("reopened coverage = %d", got)
	}
	if _, err := store.WriteRange(meta.Path, meta, []byte("last"), 0); err != nil {
		t.Fatal(err)
	}
	if err := store.FlushRanges(meta.Path, meta); err != nil {
		t.Fatal(err)
	}
	store = newRangeTestDiskCache(t, root)
	data = make([]byte, 4)
	if n, err := store.ReadRange(meta.Path, meta, data, 0); n != 4 || err != nil || string(data) != "last" {
		t.Fatalf("final partial batch = %q, %d, %v", data, n, err)
	}
	stored, found, err := store.RangeMetadata(meta.Path, meta)
	if err != nil || !found || stored.ETag != meta.ETag {
		t.Fatalf("lost validator: %+v, %t, %v", stored, found, err)
	}
}

func TestDiskCacheConcurrentRangeWritersPreserveEveryByte(t *testing.T) {
	const chunkSize = 128 * 1024
	root := t.TempDir()
	store := newRangeTestDiskCache(t, root)
	meta := rangeTestMetadata("/concurrent", 64*chunkSize)
	results := make(chan error, 64)
	for index := 0; index < 64; index++ {
		go func() {
			_, err := store.WriteRange(meta.Path, meta, bytes.Repeat([]byte{byte(index)}, chunkSize), int64(index*chunkSize))
			results <- err
		}()
	}
	for index := 0; index < 64; index++ {
		if err := <-results; err != nil {
			t.Fatal(err)
		}
	}
	if err := store.FlushRanges(meta.Path, meta); err != nil {
		t.Fatal(err)
	}
	store = newRangeTestDiskCache(t, root)
	for index := 0; index < 64; index++ {
		data := make([]byte, chunkSize)
		if n, err := store.ReadRange(meta.Path, meta, data, int64(index*chunkSize)); n != chunkSize || err != nil || !bytes.Equal(data, bytes.Repeat([]byte{byte(index)}, chunkSize)) {
			t.Fatalf("chunk %d = %d, %v", index, n, err)
		}
	}
	if got := store.SizeBytes(); got != meta.Size {
		t.Fatalf("coverage = %d, want %d", got, meta.Size)
	}
}
