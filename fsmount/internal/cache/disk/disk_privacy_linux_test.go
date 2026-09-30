//go:build linux

package disk_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	fscache "github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
)

// The cache root belongs to the caller; a caller-selected 0755 root stays
// 0755. Everything the cache itself creates below it must be readable only by
// the owner, whatever the process umask allows.
func TestDiskCacheOwnedStorageIsPrivateWhileCallerRootIsKept(t *testing.T) {
	previous := syscall.Umask(0o022)
	t.Cleanup(func() { syscall.Umask(previous) })
	root := filepath.Join(t.TempDir(), "cache")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}

	cache, err := disk.NewDiskCache(root, disk.WithLruFlushInterval(time.Hour))
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}
	populate(t, cache)
	cache.StartMaintenance()
	cache.StopMaintenance() // persists LRU state

	if got := perm(t, root); got != 0o755 {
		t.Fatalf("caller root mode changed to %04o", got)
	}
	requireTreePrivate(t, root)
}

// A cache written by an earlier version has 0755 directories and 0644 files.
// Opening it must protect that state before it is used, keep its contents
// readable, and keep it private through Clear and later writes.
func TestNewDiskCacheProtectsLegacyPermissiveState(t *testing.T) {
	previous := syscall.Umask(0o022)
	t.Cleanup(func() { syscall.Umask(previous) })
	root := filepath.Join(t.TempDir(), "cache")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy, err := disk.NewDiskCache(root, disk.WithLruFlushInterval(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	populate(t, legacy)
	legacy.StartMaintenance()
	legacy.StopMaintenance()
	// Re-create the modes the earlier version left behind.
	loosen := func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == root {
			return err
		}
		if d.IsDir() {
			return os.Chmod(path, 0o755)
		}
		return os.Chmod(path, 0o644)
	}
	if err := filepath.WalkDir(root, loosen); err != nil {
		t.Fatal(err)
	}

	reopened, err := disk.NewDiskCache(root, disk.WithLruFlushInterval(time.Hour))
	if err != nil {
		t.Fatalf("NewDiskCache on legacy state failed: %v", err)
	}
	requireTreePrivate(t, root)
	buf := make([]byte, len(payload))
	n, err := reopened.ReadComplete(payloadPath, fscache.NewEntryMetadata(payloadPath, int64(len(payload)), payloadMtime), buf, 0)
	if err != nil || n != len(payload) || string(buf) != payload {
		t.Fatalf("legacy payload not readable after protection: n=%d err=%v data=%q", n, err, buf)
	}
	if err := reopened.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}
	populate(t, reopened)
	reopened.StartMaintenance()
	reopened.StopMaintenance()
	requireTreePrivate(t, root)
}

// A link planted inside cache-owned storage is refused rather than followed,
// and the target it points at is not modified.
func TestNewDiskCacheRefusesLinkInsideOwnedStorage(t *testing.T) {
	root := filepath.Join(privateTempDir(t), "cache")
	if err := os.MkdirAll(filepath.Join(root, "data"), 0o700); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(t.TempDir(), "victim")
	if err := os.WriteFile(target, []byte("unrelated"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(root, "data", "planted")); err != nil {
		t.Fatal(err)
	}
	if _, err := disk.NewDiskCache(root); err == nil {
		t.Fatal("expected NewDiskCache to refuse a link inside the data directory")
	}
	if got := perm(t, target); got != 0o644 {
		t.Fatalf("link target mode changed to %04o", got)
	}
}

// A legacy cache file that also has a name outside the cache shares its mode
// with that name. Opening the cache must refuse it before changing anything,
// so the caller's outside file keeps the access it had.
func TestNewDiskCacheRefusesLegacyFileWithOutsideName(t *testing.T) {
	previous := syscall.Umask(0o022)
	t.Cleanup(func() { syscall.Umask(previous) })
	root := filepath.Join(t.TempDir(), "cache")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	legacy, err := disk.NewDiskCache(root, disk.WithLruFlushInterval(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	populate(t, legacy)
	legacy.StartMaintenance()
	legacy.StopMaintenance()
	var payloadFile string
	err = filepath.WalkDir(filepath.Join(root, "data"), func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			payloadFile = p
		}
		return err
	})
	if err != nil || payloadFile == "" {
		t.Fatalf("no payload file found: %v", err)
	}
	if err := os.Chmod(payloadFile, 0o644); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "callers-copy.txt")
	if err := os.Link(payloadFile, outside); err != nil {
		t.Fatal(err)
	}

	if _, err := disk.NewDiskCache(root, disk.WithLruFlushInterval(time.Hour)); err == nil {
		t.Fatal("expected NewDiskCache to refuse a cache file that has a name outside the cache")
	}
	if got := perm(t, outside); got != 0o644 {
		t.Fatalf("the caller's outside name changed mode to %04o", got)
	}
	data, err := os.ReadFile(outside)
	if err != nil || string(data) != payload {
		t.Fatalf("the caller's outside file changed: %q, %v", data, err)
	}
}

// A root that every user may write to lets another user replace the cache's
// private directories, so it is refused before anything is stored there.
func TestNewDiskCacheRefusesWorldWritableRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "shared")
	if err := os.Mkdir(root, 0o777); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(root, 0o777); err != nil {
		t.Fatal(err)
	}
	if _, err := disk.NewDiskCache(root); err == nil {
		t.Fatal("expected NewDiskCache to refuse a world-writable root")
	}
}

const (
	payloadPath = "/confidential/payroll.txt"
	payload     = "cached bytes that only the owner may read"
	partialPath = "/confidential/in-progress.txt"
)

var payloadMtime = time.Unix(1, 0)

func populate(t *testing.T, cache *disk.DiskCache) {
	t.Helper()
	if _, err := cache.Write(payloadPath, []byte(payload), 0); err != nil {
		t.Fatal(err)
	}
	if err := cache.Commit(payloadPath, fscache.NewEntryMetadata(payloadPath, int64(len(payload)), payloadMtime)); err != nil {
		t.Fatal(err)
	}
	if _, err := cache.WritePartial(partialPath, []byte("partial bytes"), 0); err != nil {
		t.Fatal(err)
	}
}

func perm(t *testing.T, path string) fs.FileMode {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	return info.Mode().Perm()
}

// requireTreePrivate checks every entry below root (not root itself) and
// requires at least one of each kind of stored artifact, so a passing check is
// never vacuous.
func requireTreePrivate(t *testing.T, root string) {
	t.Helper()
	var files, dirs int
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == root {
			return err
		}
		mode := perm(t, path)
		if mode&^0o700 != 0 {
			t.Errorf("%s is readable beyond the owner: %04o", path, mode)
		}
		if d.IsDir() {
			dirs++
		} else {
			files++
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{filepath.Join(root, "state", "lru.json")} {
		if _, err := os.Stat(want); err != nil {
			t.Fatalf("expected %s to exist: %v", want, err)
		}
	}
	if files < 4 || dirs < 4 {
		t.Fatalf("expected payload, partial, metadata and LRU files below %s, found %d files in %d directories", root, files, dirs)
	}
}
