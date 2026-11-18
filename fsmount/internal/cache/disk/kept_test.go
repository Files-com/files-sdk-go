//go:build linux || windows

package disk_test

import (
	"bytes"
	"testing"
	"time"

	fscache "github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
)

func writeCommitted(t *testing.T, dc *disk.DiskCache, path string, data []byte) {
	t.Helper()
	if _, err := dc.Write(path, data, 0); err != nil {
		t.Fatalf("Write %s failed: %v", path, err)
	}
	if err := dc.Commit(path, fscache.NewEntryMetadata(path, int64(len(data)), time.Now())); err != nil {
		t.Fatalf("Commit %s failed: %v", path, err)
	}
}

func readable(dc *disk.DiskCache, path string, want []byte) bool {
	buf := make([]byte, len(want))
	n, err := dc.Read(path, buf, 0)
	return err == nil && n == len(want) && bytes.Equal(buf, want)
}

func TestDiskCacheKeptEntryIsNotEvictedByTheFileCountLimit(t *testing.T) {
	dc, err := disk.NewDiskCache(privateTempDir(t), disk.WithMaxFileCount(2))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dc.StopMaintenance)

	kept := []byte("kept for offline use")
	writeCommitted(t, dc, "/docs/kept.txt", kept)
	if err := dc.SetKept([]string{"/docs/kept.txt"}, true); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/other/1.txt", "/other/2.txt", "/other/3.txt", "/other/4.txt"} {
		writeCommitted(t, dc, path, []byte(path))
	}

	if !dc.IsKept("/docs/kept.txt") {
		t.Fatal("kept entry is no longer kept")
	}
	if !readable(dc, "/docs/kept.txt", kept) {
		t.Fatal("kept entry was evicted by the file count limit")
	}
	if readable(dc, "/other/1.txt", []byte("/other/1.txt")) {
		t.Fatal("the oldest entry that isn't kept was not evicted")
	}
}

func TestDiskCacheKeptEntriesAreRememberedAfterARestart(t *testing.T) {
	dir := privateTempDir(t)
	dc, err := disk.NewDiskCache(dir)
	if err != nil {
		t.Fatal(err)
	}
	writeCommitted(t, dc, "/docs/a.txt", []byte("alpha"))
	if err := dc.SetKept([]string{"/docs/a.txt"}, true); err != nil {
		t.Fatal(err)
	}
	dc.StopMaintenance()

	reopened, err := disk.NewDiskCache(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(reopened.StopMaintenance)
	if !reopened.IsKept("/docs/a.txt") {
		t.Fatal("kept entry was forgotten after a restart")
	}
	if !readable(reopened, "/docs/a.txt", []byte("alpha")) {
		t.Fatal("kept entry can't be read after a restart")
	}
}

func TestDiskCacheReleasedEntryIsEvictedAgain(t *testing.T) {
	dc, err := disk.NewDiskCache(privateTempDir(t), disk.WithMaxFileCount(2))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dc.StopMaintenance)

	writeCommitted(t, dc, "/docs/a.txt", []byte("alpha"))
	if err := dc.SetKept([]string{"/docs/a.txt"}, true); err != nil {
		t.Fatal(err)
	}
	if err := dc.SetKept([]string{"/docs/a.txt"}, false); err != nil {
		t.Fatal(err)
	}
	if dc.IsKept("/docs/a.txt") {
		t.Fatal("released entry is still kept")
	}
	for _, path := range []string{"/other/1.txt", "/other/2.txt", "/other/3.txt"} {
		writeCommitted(t, dc, path, []byte(path))
	}
	if readable(dc, "/docs/a.txt", []byte("alpha")) {
		t.Fatal("released entry was not evicted by the file count limit")
	}
}

func TestDiskCacheDeleteRemovesAKeptEntry(t *testing.T) {
	dc, err := disk.NewDiskCache(privateTempDir(t))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dc.StopMaintenance)

	writeCommitted(t, dc, "/docs/a.txt", []byte("alpha"))
	if err := dc.SetKept([]string{"/docs/a.txt"}, true); err != nil {
		t.Fatal(err)
	}
	// The file changed or was removed on the storage.
	if !dc.Delete("/docs/a.txt") {
		t.Fatal("Delete of a kept entry failed")
	}
	if readable(dc, "/docs/a.txt", []byte("alpha")) {
		t.Fatal("deleted kept entry can still be read")
	}
}

func TestDiskCacheDeleteWhenClosedWaitsForTheLastHandle(t *testing.T) {
	dc, err := disk.NewDiskCache(privateTempDir(t))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dc.StopMaintenance)

	writeCommitted(t, dc, "/docs/closed.txt", []byte("closed"))
	if dc.DeleteWhenClosed("/docs/closed.txt") {
		t.Fatal("DeleteWhenClosed of a file that isn't open reported it open")
	}

	writeCommitted(t, dc, "/docs/open.txt", []byte("open"))
	dc.Pin("/docs/open.txt")
	if !dc.DeleteWhenClosed("/docs/open.txt") {
		t.Fatal("DeleteWhenClosed of an open file reported it closed")
	}
	if !readable(dc, "/docs/open.txt", []byte("open")) {
		t.Fatal("an open file lost its data before it was closed")
	}
	dc.Unpin("/docs/open.txt")
	if readable(dc, "/docs/open.txt", []byte("open")) {
		t.Fatal("the file is still cached after its last handle closed")
	}
}

func TestDiskCacheKeepingAgainCancelsDeleteWhenClosed(t *testing.T) {
	dc, err := disk.NewDiskCache(privateTempDir(t))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dc.StopMaintenance)

	writeCommitted(t, dc, "/docs/a.txt", []byte("alpha"))
	dc.Pin("/docs/a.txt")
	dc.DeleteWhenClosed("/docs/a.txt")
	if err := dc.SetKept([]string{"/docs/a.txt"}, true); err != nil {
		t.Fatal(err)
	}
	dc.Unpin("/docs/a.txt")
	if !readable(dc, "/docs/a.txt", []byte("alpha")) {
		t.Fatal("an entry kept again before it was closed was deleted on close")
	}
}
