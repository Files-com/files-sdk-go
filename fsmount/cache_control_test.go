package fsmount

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
	fslog "github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
)

func TestDiskCacheSizeAndClearOnlyUseLiveMounts(t *testing.T) {
	basePath := t.TempDir()
	activePath := filepath.Join(basePath, "A", "cache")
	inactivePath := filepath.Join(basePath, "B", "cache")
	if err := os.MkdirAll(activePath, 0o755); err != nil {
		t.Fatal(err)
	}
	activeCache, err := disk.NewDiskCache(activePath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := activeCache.Write("/active", []byte("active"), 0); err != nil {
		t.Fatal(err)
	}
	inactiveFile := filepath.Join(inactivePath, "data", "inactive")
	if err := os.MkdirAll(filepath.Dir(inactiveFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(inactiveFile, []byte("inactive"), 0o644); err != nil {
		t.Fatal(err)
	}

	mountMu.Lock()
	previousRegistry := mntRegistry
	mntRegistry = &mountRegistry{
		hosts: map[string]*Host{
			"A": {fs: &Filescomfs{remote: &RemoteFs{cacheStore: activeCache}}},
		},
	}
	mountMu.Unlock()
	t.Cleanup(func() {
		mountMu.Lock()
		mntRegistry = previousRegistry
		mountMu.Unlock()
	})

	sizeBytes, err := DiskCacheSize()
	if err != nil {
		t.Fatal(err)
	}
	if want := int64(len("active")); sizeBytes != want {
		t.Fatalf("DiskCacheSize = %d, want %d", sizeBytes, want)
	}

	remainingBytes, err := ClearDiskCache()
	if err != nil {
		t.Fatal(err)
	}
	if remainingBytes != 0 {
		t.Fatalf("ClearDiskCache remaining bytes = %d, want 0", remainingBytes)
	}
	if n, err := activeCache.Read("/active", make([]byte, len("active")), 0); err != nil || n != 0 {
		t.Fatalf("active cache remained after clear: %d, %v", n, err)
	}
	contents, err := os.ReadFile(inactiveFile)
	if err != nil {
		t.Fatalf("cache not owned by this process was removed: %v", err)
	}
	if string(contents) != "inactive" {
		t.Fatalf("cache not owned by this process was changed: %q", contents)
	}
	if _, err := activeCache.Write("/after-clear", []byte("new"), 0); err != nil {
		t.Fatalf("live cache could not be reused after clear: %v", err)
	}
}

func TestClearDiskCacheRejectsDuplicateActiveRoot(t *testing.T) {
	cacheRoot := t.TempDir()
	firstCache, err := disk.NewDiskCache(cacheRoot)
	if err != nil {
		t.Fatal(err)
	}
	secondCache, err := disk.NewDiskCache(cacheRoot)
	if err != nil {
		t.Fatal(err)
	}
	payload := []byte("pinned")
	if _, err := firstCache.Write("/pinned", payload, 0); err != nil {
		t.Fatal(err)
	}
	firstCache.Pin("/pinned")
	defer firstCache.Unpin("/pinned")

	installCacheControlTestRegistry(t, map[string]*Host{
		"A": {fs: &Filescomfs{remote: &RemoteFs{cacheStore: firstCache}}},
		"B": {fs: &Filescomfs{remote: &RemoteFs{cacheStore: secondCache}}},
	})

	if _, err := ClearDiskCache(); !errors.Is(err, errDiskCachePathInUse) {
		t.Fatalf("ClearDiskCache error = %v, want %v", err, errDiskCachePathInUse)
	}
	buffer := make([]byte, len(payload))
	if n, err := firstCache.Read("/pinned", buffer, 0); err != nil || n != len(payload) || string(buffer) != string(payload) {
		t.Fatalf("pinned cache entry changed after rejected clear: n=%d data=%q err=%v", n, buffer, err)
	}
}

func TestNewCacheRejectsCanonicalDuplicateActiveRoot(t *testing.T) {
	cacheRoot := t.TempDir()
	existingCache, err := disk.NewDiskCache(cacheRoot)
	if err != nil {
		t.Fatal(err)
	}
	installCacheControlTestRegistry(t, map[string]*Host{
		"A": {fs: &Filescomfs{remote: &RemoteFs{cacheStore: existingCache}}},
	})

	aliasPath := filepath.Join(t.TempDir(), "cache-alias")
	if err := os.Symlink(cacheRoot, aliasPath); err != nil {
		t.Skipf("symlinks are unavailable: %v", err)
	}
	_, err = newCache(MountParams{
		DiskCacheEnabled: true,
		DiskCachePath:    aliasPath,
	}, &fslog.NoOpLogger{})
	if !errors.Is(err, errDiskCachePathInUse) {
		t.Fatalf("newCache error = %v, want %v", err, errDiskCachePathInUse)
	}
}

func TestNewCacheRejectsCaseVariantOfActiveRoot(t *testing.T) {
	parent := t.TempDir()
	cacheRoot := filepath.Join(parent, "CacheRoot")
	if err := os.Mkdir(cacheRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	caseVariant := filepath.Join(parent, "cACHErOOT")
	variantInfo, err := os.Stat(caseVariant)
	if err != nil {
		t.Skip("test filesystem is case-sensitive")
	}
	cacheInfo, err := os.Stat(cacheRoot)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(cacheInfo, variantInfo) {
		t.Skip("case variant does not identify the same directory")
	}

	existingCache, err := disk.NewDiskCache(cacheRoot)
	if err != nil {
		t.Fatal(err)
	}
	installCacheControlTestRegistry(t, map[string]*Host{
		"A": {fs: &Filescomfs{remote: &RemoteFs{cacheStore: existingCache}}},
	})

	_, err = newCache(MountParams{
		DiskCacheEnabled: true,
		DiskCachePath:    caseVariant,
	}, &fslog.NoOpLogger{})
	if !errors.Is(err, errDiskCachePathInUse) {
		t.Fatalf("newCache error = %v, want %v", err, errDiskCachePathInUse)
	}
}

func installCacheControlTestRegistry(t *testing.T, hosts map[string]*Host) {
	t.Helper()
	mountMu.Lock()
	previousRegistry := mntRegistry
	mntRegistry = &mountRegistry{hosts: hosts}
	mountMu.Unlock()
	t.Cleanup(func() {
		mountMu.Lock()
		mntRegistry = previousRegistry
		mountMu.Unlock()
	})
}
