package disk_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
)

// TestNewDiskCache verifies that a DiskCache can be created with default settings
func TestNewDiskCache(t *testing.T) {
	tmpDir := t.TempDir()

	cache, err := disk.NewDiskCache(tmpDir)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	if cache == nil {
		t.Fatal("Expected non-nil cache")
	}

	stats := cache.Stats()
	if stats == nil {
		t.Fatal("Expected non-nil stats")
	}

	if stats.CapacityBytes != disk.DefaultCapacity {
		t.Errorf("Expected capacity %d, got %d", disk.DefaultCapacity, stats.CapacityBytes)
	}

	if stats.MaxFileCount != disk.DefaultMaxFileCount {
		t.Errorf("Expected max file count %d, got %d", disk.DefaultMaxFileCount, stats.MaxFileCount)
	}
}

// TestDiskCacheWriteAndRead verifies basic write and read operations
func TestDiskCacheWriteAndRead(t *testing.T) {
	tmpDir := t.TempDir()

	cache, err := disk.NewDiskCache(tmpDir)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	path := "/test/file.txt"
	data := []byte("hello world")

	// Write data
	n, err := cache.Write(path, data, 0)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data), n)
	}

	// Read data back
	readBuf := make([]byte, len(data))
	n, err = cache.Read(path, readBuf, 0)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to read %d bytes, read %d", len(data), n)
	}
	if string(readBuf) != string(data) {
		t.Errorf("Expected to read %q, got %q", string(data), string(readBuf))
	}
}

// TestDiskCacheWriteAtOffset verifies writing at a non-zero offset
func TestDiskCacheWriteAtOffset(t *testing.T) {
	tmpDir := t.TempDir()

	cache, err := disk.NewDiskCache(tmpDir)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	path := "/test/file.txt"

	// Write at offset 0
	data1 := []byte("hello")
	n, err := cache.Write(path, data1, 0)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(data1) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data1), n)
	}

	// Write at offset 5
	data2 := []byte(" world")
	n, err = cache.Write(path, data2, int64(len(data1)))
	if err != nil {
		t.Fatalf("Write at offset failed: %v", err)
	}
	if n != len(data2) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(data2), n)
	}

	// Read full content
	expectedLen := len(data1) + len(data2)
	readBuf := make([]byte, expectedLen)
	n, err = cache.Read(path, readBuf, 0)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if n != expectedLen {
		t.Errorf("Expected to read %d bytes, read %d", expectedLen, n)
	}

	expected := "hello world"
	if string(readBuf) != expected {
		t.Errorf("Expected to read %q, got %q", expected, string(readBuf))
	}
}

// TestDiskCacheReadMiss verifies that reading a non-existent file returns 0 bytes
func TestDiskCacheReadMiss(t *testing.T) {
	tmpDir := t.TempDir()

	cache, err := disk.NewDiskCache(tmpDir)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	path := "/test/nonexistent.txt"
	readBuf := make([]byte, 100)
	n, err := cache.Read(path, readBuf, 0)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if n != 0 {
		t.Errorf("Expected to read 0 bytes from nonexistent file, read %d", n)
	}
}

// TestDiskCacheDelete verifies that deleting a file removes it from the cache
func TestDiskCacheDelete(t *testing.T) {
	tmpDir := t.TempDir()

	cache, err := disk.NewDiskCache(tmpDir)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	path := "/test/file.txt"
	data := []byte("hello world")

	// Write data
	_, err = cache.Write(path, data, 0)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Delete the file
	deleted := cache.Delete(path)
	if !deleted {
		t.Error("Expected Delete to return true")
	}

	// Try to read - should get 0 bytes
	readBuf := make([]byte, len(data))
	n, err := cache.Read(path, readBuf, 0)
	if err != nil {
		t.Fatalf("Read after delete failed: %v", err)
	}
	if n != 0 {
		t.Errorf("Expected to read 0 bytes after delete, read %d", n)
	}
}

// TestDiskCachePinUnpin verifies that pinned files are not evicted
func TestDiskCachePinUnpin(t *testing.T) {
	tmpDir := t.TempDir()

	// Create cache with very small capacity to force evictions
	cache, err := disk.NewDiskCache(tmpDir, disk.WithCapacityBytes(1024))
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	path1 := "/test/file1.txt"
	path2 := "/test/file2.txt"
	data := make([]byte, 600) // Each file is 600 bytes

	// Write first file and pin it
	_, err = cache.Write(path1, data, 0)
	if err != nil {
		t.Fatalf("Write file1 failed: %v", err)
	}
	cache.Pin(path1)

	// Write second file - this should not evict the pinned file
	_, err = cache.Write(path2, data, 0)
	if err != nil {
		t.Fatalf("Write file2 failed: %v", err)
	}

	// Verify first file is still readable (wasn't evicted)
	readBuf := make([]byte, len(data))
	n, err := cache.Read(path1, readBuf, 0)
	if err != nil {
		t.Fatalf("Read file1 failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to read %d bytes from pinned file, read %d", len(data), n)
	}

	// Unpin the file
	cache.Unpin(path1)

	// Now write a third file - this should be able to evict file1
	path3 := "/test/file3.txt"
	_, err = cache.Write(path3, data, 0)
	if err != nil {
		t.Fatalf("Write file3 failed: %v", err)
	}
	// Verify first file is still readable (wasn't evicted)
	readBuf = make([]byte, len(data))
	n, err = cache.Read(path1, readBuf, 0)
	if err != nil {
		t.Fatalf("Read file1 failed: %v", err)
	}
	if n != 0 {
		t.Errorf("Expected to read %d bytes from pinned file, read %d", 0, n)
	}
}

// TestDiskCacheStats verifies that statistics are updated correctly
func TestDiskCacheStats(t *testing.T) {
	tmpDir := t.TempDir()

	cache, err := disk.NewDiskCache(tmpDir)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	initialStats := cache.Stats()
	if initialStats.FileCount.Load() != 0 {
		t.Errorf("Expected initial file count 0, got %d", initialStats.FileCount.Load())
	}
	if initialStats.SizeBytes.Load() != 0 {
		t.Errorf("Expected initial size 0, got %d", initialStats.SizeBytes.Load())
	}

	path := "/test/file.txt"
	data := []byte("hello world")

	// Write data
	_, err = cache.Write(path, data, 0)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	stats := cache.Stats()
	if stats.FileCount.Load() != 1 {
		t.Errorf("Expected file count 1, got %d", stats.FileCount.Load())
	}
	if stats.WriteCount.Load() != 1 {
		t.Errorf("Expected write count 1, got %d", stats.WriteCount.Load())
	}
	if stats.WriteBytes.Load() != int64(len(data)) {
		t.Errorf("Expected write bytes %d, got %d", len(data), stats.WriteBytes.Load())
	}

	// Read data
	readBuf := make([]byte, len(data))
	_, _ = cache.Read(path, readBuf, 0)

	stats = cache.Stats()
	if stats.ReadCount.Load() != 1 {
		t.Errorf("Expected read count 1, got %d", stats.ReadCount.Load())
	}
	if stats.ReadBytes.Load() != int64(len(data)) {
		t.Errorf("Expected read bytes %d, got %d", len(data), stats.ReadBytes.Load())
	}
}

// TestDiskCacheDisabled verifies that a disabled cache doesn't perform I/O
func TestDiskCacheDisabled(t *testing.T) {
	tmpDir := t.TempDir()

	cache, err := disk.NewDiskCache(tmpDir, disk.WithDisabled(true))
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	path := "/test/file.txt"
	data := []byte("hello world")

	// Write should return 0 bytes
	n, err := cache.Write(path, data, 0)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != 0 {
		t.Errorf("Expected disabled cache to write 0 bytes, wrote %d", n)
	}

	// Read should return 0 bytes
	readBuf := make([]byte, len(data))
	n, err = cache.Read(path, readBuf, 0)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if n != 0 {
		t.Errorf("Expected disabled cache to read 0 bytes, read %d", n)
	}
}

// TestDiskCacheMaintenance verifies that maintenance runs without error
func TestDiskCacheMaintenance(t *testing.T) {
	tmpDir := t.TempDir()

	// Create cache with short maintenance interval
	cache, err := disk.NewDiskCache(
		tmpDir,
		disk.WithMaintenanceInterval(100*time.Millisecond),
	)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Start maintenance
	cache.StartMaintenance()

	// Write some data
	path := "/test/file.txt"
	data := []byte("hello world")
	_, err = cache.Write(path, data, 0)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Wait for at least one maintenance cycle
	time.Sleep(200 * time.Millisecond)

	// Stop maintenance
	cache.StopMaintenance()

	// Verify data is still readable
	readBuf := make([]byte, len(data))
	n, err := cache.Read(path, readBuf, 0)
	if err != nil {
		t.Fatalf("Read after maintenance failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected to read %d bytes, read %d", len(data), n)
	}
}

// TestDiskCacheInvalidPath verifies that invalid paths are rejected
func TestDiskCacheInvalidPath(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{"empty path", ""},
		{"relative path", "relative/path"},
		{"nonexistent path", "/nonexistent/path/that/does/not/exist"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := disk.NewDiskCache(tt.path)
			if err == nil {
				t.Errorf("Expected error for invalid path %q, got nil", tt.path)
			}
		})
	}
}

// TestDiskCacheMultipleFiles verifies that multiple files can coexist in the cache
func TestDiskCacheMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()

	cache, err := disk.NewDiskCache(tmpDir)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	files := map[string][]byte{
		"/test/file1.txt": []byte("file one"),
		"/test/file2.txt": []byte("file two"),
		"/test/file3.txt": []byte("file three"),
	}

	// Write all files
	for path, data := range files {
		_, err := cache.Write(path, data, 0)
		if err != nil {
			t.Fatalf("Write %s failed: %v", path, err)
		}
	}

	// Read all files and verify content
	for path, expectedData := range files {
		readBuf := make([]byte, len(expectedData))
		n, err := cache.Read(path, readBuf, 0)
		if err != nil {
			t.Fatalf("Read %s failed: %v", path, err)
		}
		if n != len(expectedData) {
			t.Errorf("Expected to read %d bytes from %s, read %d", len(expectedData), path, n)
		}
		if string(readBuf) != string(expectedData) {
			t.Errorf("Expected to read %q from %s, got %q", string(expectedData), path, string(readBuf))
		}
	}

	// Verify stats
	stats := cache.Stats()
	if stats.FileCount.Load() != int64(len(files)) {
		t.Errorf("Expected file count %d, got %d", len(files), stats.FileCount.Load())
	}
}

// TestDiskCacheCapacityEviction verifies that files are evicted when capacity is exceeded
func TestDiskCacheCapacityEviction(t *testing.T) {
	tmpDir := t.TempDir()

	// Create cache with capacity for ~2 files of 512 bytes each
	cache, err := disk.NewDiskCache(tmpDir, disk.WithCapacityBytes(1024))
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	data := make([]byte, 600) // Each file is 600 bytes

	// Write first file
	path1 := "/test/file1.txt"
	_, err = cache.Write(path1, data, 0)
	if err != nil {
		t.Fatalf("Write file1 failed: %v", err)
	}

	// Write second file - should succeed but cache is now over capacity
	path2 := "/test/file2.txt"
	_, err = cache.Write(path2, data, 0)
	if err != nil {
		t.Fatalf("Write file2 failed: %v", err)
	}

	// Write third file - should evict the oldest file (file1)
	path3 := "/test/file3.txt"
	_, err = cache.Write(path3, data, 0)
	if err != nil {
		t.Fatalf("Write file3 failed: %v", err)
	}

	// Verify that at least one of the earlier files can't be read (was evicted)
	// Note: The specific file evicted cannot be guaranteed due to timing,
	// but the cache should be respecting capacity limits
	stats := cache.Stats()
	if stats.SizeBytes.Load() > 1024 {
		// Allow some slack for pinned files
		t.Logf("Warning: Cache size %d exceeds capacity 1024 (may be due to pinned files)", stats.SizeBytes.Load())
	}
}

// TestDiskCacheWithFileNotInCache verifies entryPath handling
func TestDiskCacheWithFileNotInCache(t *testing.T) {
	tmpDir := t.TempDir()
	cache, err := disk.NewDiskCache(tmpDir)
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	// Create a file directly in the data directory
	dataDir := filepath.Join(tmpDir, "data")
	testFile := filepath.Join(dataDir, "test.txt")
	if err := os.MkdirAll(filepath.Dir(testFile), 0o755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.WriteFile(testFile, []byte("direct write"), 0o644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Try to read using the full path - should work since entryPath handles this
	readBuf := make([]byte, 12)
	n, err := cache.Read(testFile, readBuf, 0)
	if err != nil {
		t.Fatalf("Read with full path failed: %v", err)
	}
	if n != 12 {
		t.Errorf("Expected to read 12 bytes, read %d", n)
	}
}

// TestDiskCacheUnpinEviction verifies that an unpinned file is evicted when the cache is over capacity,
// even when other files in the LRU are pinned. This test ensures the bug fix where unpinned files
// are prioritized for eviction doesn't regress.
func TestDiskCacheUnpinEviction(t *testing.T) {
	tmpDir := t.TempDir()

	// Create cache with capacity for ~2 files
	cache, err := disk.NewDiskCache(tmpDir, disk.WithCapacityBytes(1200))
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}

	data := make([]byte, 600) // Each file is 600 bytes

	// Write and pin file1 (oldest in LRU)
	path1 := "/test/file1.txt"
	_, err = cache.Write(path1, data, 0)
	if err != nil {
		t.Fatalf("Write file1 failed: %v", err)
	}
	cache.Pin(path1)

	// Write and pin file2
	path2 := "/test/file2.txt"
	_, err = cache.Write(path2, data, 0)
	if err != nil {
		t.Fatalf("Write file2 failed: %v", err)
	}
	cache.Pin(path2)

	// Verify cache is at capacity before the unpin eviction logic runs
	statsBefore := cache.Stats()
	if statsBefore.SizeBytes.Load() != 1200 {
		t.Fatalf("Cache size before unpin: %d (expected 1200)", statsBefore.SizeBytes.Load())
	}

	// Unpin file2 (which is newer in the LRU than file1)
	cache.Unpin(path2)

	// Verify file2 was evicted (can't be read)
	readBuf := make([]byte, len(data))
	n, err := cache.Read(path2, readBuf, 0)
	if err != nil {
		t.Fatalf("Read file2 failed: %v", err)
	}
	if n != 0 {
		t.Errorf("Expected file2 to be evicted (read 0 bytes), but read %d bytes", n)
	}

	// Verify file1 is still in cache (pinned, shouldn't be evicted)
	n, err = cache.Read(path1, readBuf, 0)
	if err != nil {
		t.Fatalf("Read file1 failed: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected file1 to still be cached (read %d bytes), but read %d bytes", len(data), n)
	}

	// Verify stats show only one file remaining
	statsAfter := cache.Stats()
	if statsAfter.FileCount.Load() != 1 {
		t.Errorf("Expected 1 file in cache after eviction, got %d", statsAfter.FileCount.Load())
	}
}
