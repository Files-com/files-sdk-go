//go:build linux || windows

package mem_test

import (
	"bytes"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	fscache "github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/mem"
)

func TestMemoryCacheRangeWritesDoNotTreatAllocatedPageHolesAsHits(t *testing.T) {
	cacheStore := newRangeTestMemoryCache(t)
	path := "/ranges.bin"
	meta := rangeTestMetadata(path, 40)

	writes := []struct {
		offset int64
		data   []byte
	}{
		{offset: 30, data: bytes.Repeat([]byte("t"), 10)},
		{offset: 0, data: bytes.Repeat([]byte("h"), 10)},
	}
	for _, write := range writes {
		n, err := cacheStore.WriteRange(path, meta, write.data, write.offset)
		if err != nil || n != len(write.data) {
			t.Fatalf("WriteRange(%d) = %d, %v; want %d, nil", write.offset, n, err, len(write.data))
		}
	}

	holeBuffer := bytes.Repeat([]byte{0xff}, 30)
	if n, err := cacheStore.ReadRange(path, meta, holeBuffer, 5); err != nil || n != 0 {
		t.Fatalf("ReadRange across same-page hole = %d, %v; want cache miss", n, err)
	}

	missing, err := cacheStore.MissingRanges(path, meta, fscache.ByteRange{Start: 0, End: 40})
	if err != nil {
		t.Fatalf("MissingRanges failed: %v", err)
	}
	wantMissing := []fscache.ByteRange{{Start: 10, End: 30}}
	if !reflect.DeepEqual(missing, wantMissing) {
		t.Fatalf("MissingRanges = %#v, want %#v", missing, wantMissing)
	}
	if complete, err := cacheStore.RangeEntryComplete(path, meta); err != nil || complete {
		t.Fatalf("RangeEntryComplete = %t, %v; want false, nil", complete, err)
	}

	if _, err := cacheStore.WriteCompleteRange(path, meta, bytes.Repeat([]byte("m"), 20), 10); err != nil {
		t.Fatalf("WriteCompleteRange gap failed: %v", err)
	}
	if complete, err := cacheStore.RangeEntryComplete(path, meta); err != nil || !complete {
		t.Fatalf("RangeEntryComplete = %t, %v; want true, nil", complete, err)
	}
	buffer := make([]byte, 40)
	n, err := cacheStore.ReadRange(path, meta, buffer, 0)
	if err != nil || n != len(buffer) {
		t.Fatalf("ReadRange complete = %d, %v; want %d, nil", n, err, len(buffer))
	}
	want := append(bytes.Repeat([]byte("h"), 10), bytes.Repeat([]byte("m"), 20)...)
	want = append(want, bytes.Repeat([]byte("t"), 10)...)
	if !bytes.Equal(buffer, want) {
		t.Fatalf("complete data = %q, want %q", buffer, want)
	}
}

func TestMemoryCacheRangeOverlapTracksPagesAndNewBytes(t *testing.T) {
	cacheStore := newRangeTestMemoryCache(t)
	path := "/overlap.bin"
	meta := rangeTestMetadata(path, 100)

	if _, err := cacheStore.WriteRange(path, meta, bytes.Repeat([]byte("a"), 10), 10); err != nil {
		t.Fatalf("first WriteRange failed: %v", err)
	}
	allocatedAfterFirst := cacheStore.Stats().SizeBytes.Load()
	if allocatedAfterFirst <= 10 {
		t.Fatalf("allocated bytes = %d, want page-based accounting", allocatedAfterFirst)
	}
	if _, err := cacheStore.WriteRange(path, meta, bytes.Repeat([]byte("b"), 10), 15); err != nil {
		t.Fatalf("overlapping WriteRange failed: %v", err)
	}
	if _, err := cacheStore.WriteRange(path, meta, bytes.Repeat([]byte("c"), 15), 10); err != nil {
		t.Fatalf("covered WriteRange failed: %v", err)
	}

	stats := cacheStore.Stats()
	if got := stats.SizeBytes.Load(); got != allocatedAfterFirst {
		t.Fatalf("allocated bytes = %d, want unchanged %d", got, allocatedAfterFirst)
	}
	if got := stats.WriteBytes.Load(); got != 15 {
		t.Fatalf("WriteBytes = %d, want 15 newly valid bytes", got)
	}

	buffer := make([]byte, 15)
	n, err := cacheStore.ReadRange(path, meta, buffer, 10)
	if err != nil || n != len(buffer) {
		t.Fatalf("ReadRange = %d, %v; want %d, nil", n, err, len(buffer))
	}
	want := append(bytes.Repeat([]byte("a"), 10), bytes.Repeat([]byte("b"), 5)...)
	if !bytes.Equal(buffer, want) {
		t.Fatalf("overlap data = %q, want %q", buffer, want)
	}
}

func TestMemoryCacheRangeVersionChangeInvalidatesEntry(t *testing.T) {
	cacheStore := newRangeTestMemoryCache(t)
	path := "/versioned.bin"
	oldMeta := rangeTestMetadata(path, 10)
	newMeta := fscache.NewEntryMetadata(path, 10, oldMeta.ModTime.Add(time.Second))
	if _, err := cacheStore.WriteRange(path, oldMeta, bytes.Repeat([]byte("a"), 10), 0); err != nil {
		t.Fatalf("WriteRange failed: %v", err)
	}

	if n, err := cacheStore.ReadRange(path, newMeta, make([]byte, 10), 0); err != nil || n != 0 {
		t.Fatalf("stale ReadRange = %d, %v; want cache miss", n, err)
	}
	if got := cacheStore.Stats().FileCount.Load(); got != 0 {
		t.Fatalf("FileCount after invalidation = %d, want 0", got)
	}
	missing, err := cacheStore.MissingRanges(path, newMeta, fscache.ByteRange{Start: 0, End: 10})
	if err != nil || !reflect.DeepEqual(missing, []fscache.ByteRange{{Start: 0, End: 10}}) {
		t.Fatalf("MissingRanges after invalidation = %#v, %v", missing, err)
	}
}

func TestMemoryCacheRangeEvictionRemovesWholeEntry(t *testing.T) {
	cacheStore, err := mem.NewMemoryCache(mem.WithMaxFileCount(1))
	if err != nil {
		t.Fatalf("NewMemoryCache failed: %v", err)
	}
	firstPath := "/first.bin"
	secondPath := "/second.bin"
	if _, err := cacheStore.WriteRange(firstPath, rangeTestMetadata(firstPath, 10), bytes.Repeat([]byte("a"), 10), 0); err != nil {
		t.Fatalf("first WriteRange failed: %v", err)
	}
	if _, err := cacheStore.WriteRange(secondPath, rangeTestMetadata(secondPath, 10), bytes.Repeat([]byte("b"), 10), 0); err != nil {
		t.Fatalf("second WriteRange failed: %v", err)
	}

	if n, err := cacheStore.ReadRange(firstPath, rangeTestMetadata(firstPath, 10), make([]byte, 10), 0); err != nil || n != 0 {
		t.Fatalf("evicted ReadRange = %d, %v; want cache miss", n, err)
	}
	if n, err := cacheStore.ReadRange(secondPath, rangeTestMetadata(secondPath, 10), make([]byte, 10), 0); err != nil || n != 10 {
		t.Fatalf("retained ReadRange = %d, %v; want 10, nil", n, err)
	}
	if got := cacheStore.Stats().FileCount.Load(); got != 1 {
		t.Fatalf("FileCount = %d, want 1", got)
	}
}

func TestMemoryCacheRangeVersionChangeReplacesPinnedEntry(t *testing.T) {
	cacheStore := newRangeTestMemoryCache(t)
	path := "/pinned.bin"
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

func TestMemoryCacheRangeWriteAndDeleteCannotPublishRemovedBytes(t *testing.T) {
	for i := 0; i < 200; i++ {
		cacheStore := newRangeTestMemoryCache(t)
		path := "/delete-race.bin"
		meta := rangeTestMetadata(path, 20)
		if _, err := cacheStore.WriteRange(path, meta, bytes.Repeat([]byte("a"), 10), 0); err != nil {
			t.Fatalf("initial WriteRange failed: %v", err)
		}

		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(2)
		var writeErr error
		go func() {
			defer wg.Done()
			<-start
			_, writeErr = cacheStore.WriteRange(path, meta, bytes.Repeat([]byte("b"), 10), 10)
		}()
		go func() {
			defer wg.Done()
			<-start
			cacheStore.Delete(path)
		}()
		close(start)
		wg.Wait()
		if writeErr != nil {
			t.Fatalf("concurrent WriteRange failed: %v", writeErr)
		}

		if n, err := cacheStore.ReadRange(path, meta, make([]byte, 20), 0); err != nil || n != 0 {
			t.Fatalf("iteration %d ReadRange after concurrent delete = %d, %v; want cache miss", i, n, err)
		}
	}
}

func TestMemoryCacheRangeRejectsOutOfBoundsWrite(t *testing.T) {
	cacheStore := newRangeTestMemoryCache(t)
	meta := rangeTestMetadata("/bounds.bin", 10)
	_, err := cacheStore.WriteRange(meta.Path, meta, make([]byte, 2), 9)
	if !errors.Is(err, fscache.ErrInvalidByteRange) {
		t.Fatalf("WriteRange error = %v, want ErrInvalidByteRange", err)
	}
}

func newRangeTestMemoryCache(t *testing.T) *mem.MemoryCache {
	t.Helper()
	cacheStore, err := mem.NewMemoryCache()
	if err != nil {
		t.Fatalf("NewMemoryCache failed: %v", err)
	}
	return cacheStore
}

func rangeTestMetadata(path string, size int64) fscache.EntryMetadata {
	return fscache.NewEntryMetadata(path, size, time.Date(2026, time.August, 13, 12, 0, 0, 0, time.UTC))
}
