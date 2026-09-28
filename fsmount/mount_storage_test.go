//go:build linux || windows

package fsmount

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
	"github.com/winfsp/cgofuse/fuse"
)

// The mount's private storage is a subdirectory the mount creates inside the
// caller-selected TmpFsPath. Unmounting removes exactly that subdirectory: the
// caller's directory and its other contents stay.
func TestFilescomfsDestroyRemovesOnlyTheMountsOwnStorage(t *testing.T) {
	parent := t.TempDir()
	keep := filepath.Join(parent, "callers-own-file")
	if err := os.WriteFile(keep, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	storage, err := newMountStorage(parent)
	if err != nil {
		t.Fatalf("newMountStorage failed: %v", err)
	}
	if filepath.Dir(storage.root) != parent {
		t.Fatalf("storage root %s is not directly inside the selected directory %s", storage.root, parent)
	}
	if filepath.Dir(storage.local) != storage.root || filepath.Dir(storage.write) != storage.root {
		t.Fatalf("local %s and write %s must live inside %s", storage.local, storage.write, storage.root)
	}

	fs, _, local, _, _ := newTestFilescomfs(t)
	local.localFsRoot = storage.local
	fs.storage = storage
	errc, fh := local.Create("/report.tmp", fuse.O_CREAT|fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned %d", errc)
	}
	if n := local.Write("/report.tmp", []byte("scratch"), 0, fh); n != len("scratch") {
		t.Fatalf("Write returned %d", n)
	}
	if errc := local.Release("/report.tmp", fh); errc != 0 {
		t.Fatalf("Release returned %d", errc)
	}
	if _, err := os.Stat(filepath.Join(storage.local, "report.tmp")); err != nil {
		t.Fatalf("scratch file was not written inside the mount's storage: %v", err)
	}
	fs.Destroy()

	if _, err := os.Stat(storage.root); !os.IsNotExist(err) {
		t.Fatalf("mount storage %s should be removed on Destroy, stat err=%v", storage.root, err)
	}
	if data, err := os.ReadFile(keep); err != nil || string(data) != "keep" {
		t.Fatalf("the caller's own directory contents must survive: %v", err)
	}
	if _, err := os.Stat(parent); err != nil {
		t.Fatalf("the caller's directory must survive: %v", err)
	}
}

// Working copies of files being uploaded hold their content; they must be
// created inside the mount's private storage, and a session created without a
// mount still gets a private directory of its own rather than the shared
// temporary directory.
func TestNewWriteSessionKeepsWorkingCopyInPrivateDirectory(t *testing.T) {
	parent := t.TempDir()
	storage, err := newMountStorage(parent)
	if err != nil {
		t.Fatal(err)
	}
	defer storage.remove()

	session, err := newWriteSession(storage.write, "/doc.txt", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(session.workingCopyPath) != storage.write {
		t.Fatalf("working copy %s is not inside %s", session.workingCopyPath, storage.write)
	}
	if err := session.closeAndRemoveWorkingCopy(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(storage.write); err != nil {
		t.Fatalf("the mount's write directory must remain for later sessions: %v", err)
	}

	standalone, err := newWriteSession("", "/doc.txt", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Dir(standalone.workingCopyPath)
	if filepath.Clean(dir) == filepath.Clean(os.TempDir()) {
		t.Fatalf("a standalone working copy must not be created directly in %s", os.TempDir())
	}
	requireOwnerOnlyDir(t, dir)
	if err := standalone.closeAndRemoveWorkingCopy(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("the standalone session's private directory should be removed with it: %v", err)
	}
}

// An application cannot widen the mount's storage root through the mount.
func TestLocalFsChmodOnRootKeepsStoragePrivate(t *testing.T) {
	parent := t.TempDir()
	storage, err := newMountStorage(parent)
	if err != nil {
		t.Fatal(err)
	}
	defer storage.remove()
	local := newLocalFs(MountParams{TmpFsPath: storage.local}, &virtualfs{nodes: map[string]*fsNode{}, handles: &OpenHandles{entries: map[uint64]*fileHandle{}, log: &log.NoOpLogger{}}, LeveledLogger: &log.NoOpLogger{}}, &log.NoOpLogger{})
	if errc := local.Chmod("/", 0o777); errc != 0 {
		t.Fatalf("Chmod returned %d", errc)
	}
	requireOwnerOnlyDir(t, storage.local)
}
