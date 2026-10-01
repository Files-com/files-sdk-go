//go:build windows

package disk

import (
	"os"
	"testing"
	"time"
	"unsafe"

	fscache "github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"golang.org/x/sys/windows"
)

type fileStandardInfo struct {
	AllocationSize int64
	EndOfFile      int64
	NumberOfLinks  uint32
	DeletePending  byte
	Directory      byte
	_              [2]byte
}

func TestDiskCacheTailRangeUsesSparseAllocation(t *testing.T) {
	cacheStore, err := NewDiskCache(t.TempDir())
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}
	t.Cleanup(cacheStore.StopMaintenance)
	const fileSize = 128 << 20
	payload := make([]byte, 4096)
	for i := range payload {
		payload[i] = 1
	}
	path := "/tail-range.bin"
	meta := fscache.NewEntryMetadata(path, fileSize, time.Now())
	offset := int64(fileSize - len(payload))
	if _, err := cacheStore.WriteRange(path, meta, payload, offset); err != nil {
		t.Fatalf("WriteRange failed: %v", err)
	}

	file, err := os.Open(cacheStore.entryPath(path))
	if err != nil {
		t.Fatalf("Open cache entry failed: %v", err)
	}
	defer file.Close()
	info, err := getFileStandardInfo(file)
	if err != nil {
		t.Fatalf("getFileStandardInfo failed: %v", err)
	}
	if info.EndOfFile != fileSize {
		t.Fatalf("logical size = %d, want %d", info.EndOfFile, fileSize)
	}
	if info.AllocationSize >= info.EndOfFile {
		t.Fatalf("allocated size = %d, want less than logical size %d", info.AllocationSize, info.EndOfFile)
	}
	if got := cacheStore.SizeBytes(); got != int64(len(payload)) {
		t.Fatalf("cache accounting = %d, want %d", got, len(payload))
	}
}

func TestMarkSparseFileAvoidsAllocatingHoles(t *testing.T) {
	path := t.TempDir() + `\sparse-cache-entry`
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		t.Fatalf("OpenFile failed: %v", err)
	}
	defer file.Close()

	if err := markSparseFile(file); err != nil {
		t.Fatalf("markSparseFile failed: %v", err)
	}

	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		t.Fatalf("UTF16PtrFromString failed: %v", err)
	}
	attributes, err := windows.GetFileAttributes(pathPtr)
	if err != nil {
		t.Fatalf("GetFileAttributes failed: %v", err)
	}
	if attributes&windows.FILE_ATTRIBUTE_SPARSE_FILE == 0 {
		t.Fatalf("file attributes %#x do not include FILE_ATTRIBUTE_SPARSE_FILE", attributes)
	}

	const writeOffset = 64 << 20
	payload := make([]byte, 4096)
	for i := range payload {
		payload[i] = 1
	}
	if _, err := file.WriteAt(payload, writeOffset); err != nil {
		t.Fatalf("WriteAt failed: %v", err)
	}
	if err := file.Sync(); err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	info, err := getFileStandardInfo(file)
	if err != nil {
		t.Fatalf("getFileStandardInfo failed: %v", err)
	}
	if info.EndOfFile != writeOffset+int64(len(payload)) {
		t.Fatalf("logical size = %d, want %d", info.EndOfFile, writeOffset+int64(len(payload)))
	}
	if info.AllocationSize >= info.EndOfFile {
		t.Fatalf("allocated size = %d, want less than logical size %d", info.AllocationSize, info.EndOfFile)
	}
}

func getFileStandardInfo(file *os.File) (fileStandardInfo, error) {
	var info fileStandardInfo
	err := windows.GetFileInformationByHandleEx(
		windows.Handle(file.Fd()),
		windows.FileStandardInfo,
		(*byte)(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info)),
	)
	return info, err
}
