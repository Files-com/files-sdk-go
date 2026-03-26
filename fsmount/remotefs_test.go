package fsmount

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/mem"
	lim "github.com/Files-com/files-sdk-go/v3/fsmount/internal/limit"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
	fssync "github.com/Files-com/files-sdk-go/v3/fsmount/internal/sync"
	"github.com/winfsp/cgofuse/fuse"
)

func newTestRemoteFs(t *testing.T) (*RemoteFs, *virtualfs, cacheStore) {
	t.Helper()

	cacheStore, err := mem.NewMemoryCache()
	if err != nil {
		t.Fatalf("NewMemoryCache failed: %v", err)
	}

	vfs := &virtualfs{
		nodes:         make(map[string]*fsNode),
		handles:       &OpenHandles{entries: make(map[uint64]*fileHandle), log: &log.NoOpLogger{}},
		LeveledLogger: &log.NoOpLogger{},
		cacheTTL:      DefaultCacheTTL,
	}
	root := vfs.getOrCreate("/", nodeTypeDir)
	root.extendTtl()

	fs := &RemoteFs{
		log:            &log.NoOpLogger{},
		vfs:            vfs,
		cacheStore:     cacheStore,
		disableLocking: true,
		readyGates:     map[string]*cache.ReadyGate{},
		loadDirMutexes: fssync.NewPathMutex(),
		backend:        &fakeRemoteBackend{},
		ops: lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
			lim.FuseOpDownload: downloadOpLimit,
			lim.FuseOpUpload:   uploadOpLimit,
			lim.FuseOpOther:    otherOpLimit,
		}, globalOpLimit),
	}

	return fs, vfs, cacheStore
}

type fakeRemoteBackend struct {
	findFunc           func(params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)
	uploadFunc         func(opts ...file.UploadOption) error
	moveFunc           func(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error)
	waitFunc           func(action files_sdk.FileAction, status func(files_sdk.FileMigration), opts ...files_sdk.RequestResponseOption) (files_sdk.FileMigration, error)
	downloadToFileFunc func(params files_sdk.FileDownloadParams, filePath string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)
	deleteFunc         func(params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) error
}

func (b *fakeRemoteBackend) findCurrent(opts ...files_sdk.RequestResponseOption) (files_sdk.ApiKey, error) {
	return files_sdk.ApiKey{}, nil
}

func (b *fakeRemoteBackend) find(params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	if b.findFunc != nil {
		return b.findFunc(params, opts...)
	}
	return files_sdk.File{}, nil
}

func (b *fakeRemoteBackend) listFor(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (remoteFileIter, error) {
	return nil, nil
}

func (b *fakeRemoteBackend) createFolder(params files_sdk.FolderCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return files_sdk.File{}, nil
}

func (b *fakeRemoteBackend) move(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
	if b.moveFunc != nil {
		return b.moveFunc(params, opts...)
	}
	return files_sdk.FileAction{}, nil
}

func (b *fakeRemoteBackend) update(params files_sdk.FileUpdateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return files_sdk.File{}, nil
}

func (b *fakeRemoteBackend) uploadWithResume(opts ...file.UploadOption) (file.UploadResumable, error) {
	return file.UploadResumable{}, nil
}

func (b *fakeRemoteBackend) upload(opts ...file.UploadOption) error {
	if b.uploadFunc != nil {
		return b.uploadFunc(opts...)
	}
	return nil
}

func (b *fakeRemoteBackend) downloadToFile(params files_sdk.FileDownloadParams, filePath string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	if b.downloadToFileFunc != nil {
		return b.downloadToFileFunc(params, filePath, opts...)
	}
	return files_sdk.File{}, nil
}

func (b *fakeRemoteBackend) download(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return files_sdk.File{}, nil
}

func (b *fakeRemoteBackend) createLock(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.Lock, error) {
	return files_sdk.Lock{}, nil
}

func (b *fakeRemoteBackend) deleteLock(params files_sdk.LockDeleteParams, opts ...files_sdk.RequestResponseOption) error {
	return nil
}

func (b *fakeRemoteBackend) listLocksFor(params files_sdk.LockListForParams, opts ...files_sdk.RequestResponseOption) (remoteLockIter, error) {
	return nil, nil
}

func (b *fakeRemoteBackend) delete(params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) error {
	if b.deleteFunc != nil {
		return b.deleteFunc(params, opts...)
	}
	return nil
}

func (b *fakeRemoteBackend) wait(action files_sdk.FileAction, status func(files_sdk.FileMigration), opts ...files_sdk.RequestResponseOption) (files_sdk.FileMigration, error) {
	if b.waitFunc != nil {
		return b.waitFunc(action, status, opts...)
	}
	return files_sdk.FileMigration{Status: "completed"}, nil
}

func newTestFilescomfs(t *testing.T) (*Filescomfs, *RemoteFs, *LocalFs, *virtualfs, cacheStore) {
	t.Helper()

	cacheStore, err := mem.NewMemoryCache()
	if err != nil {
		t.Fatalf("NewMemoryCache failed: %v", err)
	}

	vfs := &virtualfs{
		nodes:         make(map[string]*fsNode),
		handles:       &OpenHandles{entries: make(map[uint64]*fileHandle), log: &log.NoOpLogger{}},
		LeveledLogger: &log.NoOpLogger{},
		cacheTTL:      DefaultCacheTTL,
	}
	root := vfs.getOrCreate("/", nodeTypeDir)
	root.extendTtl()

	params := MountParams{
		TmpFsPath: t.TempDir(),
		CacheTTL:  DefaultCacheTTL,
	}
	ig, err := ignoreFromPatterns(nil)
	if err != nil {
		t.Fatalf("ignoreFromPatterns failed: %v", err)
	}

	remote := &RemoteFs{
		log:            &log.NoOpLogger{},
		vfs:            vfs,
		cacheStore:     cacheStore,
		disableLocking: true,
		readyGates:     map[string]*cache.ReadyGate{},
		loadDirMutexes: fssync.NewPathMutex(),
		backend:        &fakeRemoteBackend{},
		ops: lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
			lim.FuseOpDownload: downloadOpLimit,
			lim.FuseOpUpload:   uploadOpLimit,
			lim.FuseOpOther:    otherOpLimit,
		}, globalOpLimit),
	}
	local := newLocalFs(params, vfs, &log.NoOpLogger{})
	local.Init()

	fs := &Filescomfs{
		remote:      remote,
		local:       local,
		vfs:         vfs,
		log:         &log.NoOpLogger{},
		localFsRoot: params.TmpFsPath,
		ignore:      ig,
	}

	return fs, remote, local, vfs, cacheStore
}

func TestRemoteFsPublicWriteReadGetattrUsesWorkingCopy(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/design.ai"
	errc, fh := fs.Create(path, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	payload := []byte("session-backed-data")
	if n := fs.Write(path, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}

	buf := make([]byte, len(payload))
	if n := fs.Read(path, buf, 0, fh); n != len(payload) {
		t.Fatalf("Read returned %d, want %d", n, len(payload))
	}
	if string(buf) != string(payload) {
		t.Fatalf("Read returned %q, want %q", string(buf), string(payload))
	}

	var stat fuse.Stat_t
	if errc := fs.Getattr(path, &stat, fh); errc != 0 {
		t.Fatalf("Getattr returned unexpected error: %d", errc)
	}
	if stat.Size != int64(len(payload)) {
		t.Fatalf("Getattr size = %d, want %d", stat.Size, len(payload))
	}
	if stat.Nlink != 1 {
		t.Fatalf("Getattr nlink = %d, want 1", stat.Nlink)
	}
}

func TestRemoteFsOpenMissingFileReturnsENOENTAndDoesNotCreateNode(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/missing.ai"
	errc, fh := fs.Open(path, fuse.O_RDWR)
	if errc != -fuse.ENOENT {
		t.Fatalf("Open returned %d, want %d", errc, -fuse.ENOENT)
	}
	if fh != ^uint64(0) {
		t.Fatalf("Open fh = %d, want invalid handle", fh)
	}
	if _, ok := vfs.fetch(path); ok {
		t.Fatal("expected missing path to remain absent from VFS after failed open")
	}
}

func TestRemoteFsTruncateMissingFileReturnsENOENTAndDoesNotCreateNode(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/missing.ai"
	errc := fs.Truncate(path, 0, ^uint64(0))
	if errc != -fuse.ENOENT {
		t.Fatalf("Truncate returned %d, want %d", errc, -fuse.ENOENT)
	}
	if _, ok := vfs.fetch(path); ok {
		t.Fatal("expected missing path to remain absent from VFS after failed truncate")
	}
}

func TestRemoteFsUnlinkAllowsSameUserLockedDirtyWriteSession(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	deleteCalls := 0
	fs.backend = &fakeRemoteBackend{
		deleteFunc: func(params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) error {
			deleteCalls++
			return nil
		},
	}

	path := "/dirty.ai"
	errc, fh := fs.Create(path, fuse.O_RDWR|fuse.O_CREAT, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	if n := fs.Write(path, []byte("data"), 0, fh); n != 4 {
		t.Fatalf("Write returned %d, want 4", n)
	}

	node, ok := vfs.fetch(path)
	if !ok {
		t.Fatal("expected created path to exist in VFS")
	}
	node.setLockOwner("current-user")
	fs.currentUserId = 123
	fs.lockMap = map[string]*lockInfo{
		path: {
			Fh: fh,
			Lock: &files_sdk.Lock{
				UserId:   123,
				Username: "current-user",
			},
		},
	}

	if errc := fs.Unlink(path); errc != 0 {
		t.Fatalf("Unlink returned %d, want 0", errc)
	}
	if deleteCalls != 1 {
		t.Fatalf("delete called %d times, want 1", deleteCalls)
	}
	if _, ok := vfs.fetch(path); ok {
		t.Fatal("expected path to be removed from VFS after unlink")
	}
}

func TestRemoteFsMkdirOnExistingFileReturnsEEXISTWithoutChangingNodeType(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/existing-file"
	errc, fh := fs.Create(path, fuse.O_RDWR|fuse.O_CREAT, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}
	defer fs.Release(path, fh)

	node, ok := vfs.fetch(path)
	if !ok {
		t.Fatal("expected created file to exist in VFS")
	}
	if node.info.nodeType != nodeTypeFile {
		t.Fatalf("initial node type = %v, want file", node.info.nodeType)
	}

	if errc := fs.Mkdir(path, 0o755); errc != -fuse.EEXIST {
		t.Fatalf("Mkdir returned %d, want %d", errc, -fuse.EEXIST)
	}

	node, ok = vfs.fetch(path)
	if !ok {
		t.Fatal("expected path to remain in VFS after failed Mkdir")
	}
	if node.info.nodeType != nodeTypeFile {
		t.Fatalf("node type after failed Mkdir = %v, want file", node.info.nodeType)
	}

	if errc := fs.Rmdir(path); errc != -fuse.ENOTDIR {
		t.Fatalf("Rmdir returned %d, want %d", errc, -fuse.ENOTDIR)
	}
	if errc := fs.Unlink(path); errc != 0 {
		t.Fatalf("Unlink returned %d, want 0", errc)
	}
}

func TestRemoteFsCreateExOpenExistingOnMissingFileReturnsENOENT(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/missing.ai"
	fi := &fuse.FileInfo_t{Flags: fuse.O_RDWR}

	errc := fs.CreateEx(path, 0o644, fi)
	if errc != -fuse.ENOENT {
		t.Fatalf("CreateEx returned %d, want %d", errc, -fuse.ENOENT)
	}
	if fi.Fh != ^uint64(0) {
		t.Fatalf("CreateEx fh = %d, want invalid handle", fi.Fh)
	}
	if _, ok := vfs.fetch(path); ok {
		t.Fatal("expected missing path to remain absent from VFS after failed CreateEx open")
	}
}

func TestRemoteFsCreateExTruncateExistingOnMissingFileReturnsENOENT(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/missing.ai"
	fi := &fuse.FileInfo_t{Flags: fuse.O_RDWR | fuse.O_TRUNC}

	errc := fs.CreateEx(path, 0o644, fi)
	if errc != -fuse.ENOENT {
		t.Fatalf("CreateEx returned %d, want %d", errc, -fuse.ENOENT)
	}
	if fi.Fh != ^uint64(0) {
		t.Fatalf("CreateEx fh = %d, want invalid handle", fi.Fh)
	}
	if _, ok := vfs.fetch(path); ok {
		t.Fatal("expected missing path to remain absent from VFS after failed CreateEx truncate")
	}
}

func TestFilescomfsCreateExOpenExistingDoesNotCreatePhantomRemoteNode(t *testing.T) {
	fs, _, _, vfs, _ := newTestFilescomfs(t)
	defer vfs.destroy()

	path := "/missing.ai"
	fi := &fuse.FileInfo_t{Flags: fuse.O_RDWR}

	errc := fs.CreateEx(path, 0o644, fi)
	if errc != -fuse.ENOENT {
		t.Fatalf("CreateEx returned %d, want %d", errc, -fuse.ENOENT)
	}
	if _, ok := vfs.fetch(path); ok {
		t.Fatal("expected missing path to remain absent from VFS after failed Filescomfs CreateEx open")
	}
}

func TestFilescomfsCreateExTruncateExistingDoesNotCreatePhantomRemoteNode(t *testing.T) {
	fs, _, _, vfs, _ := newTestFilescomfs(t)
	defer vfs.destroy()

	path := "/missing.ai"
	fi := &fuse.FileInfo_t{Flags: fuse.O_RDWR | fuse.O_TRUNC}

	errc := fs.CreateEx(path, 0o644, fi)
	if errc != -fuse.ENOENT {
		t.Fatalf("CreateEx returned %d, want %d", errc, -fuse.ENOENT)
	}
	if _, ok := vfs.fetch(path); ok {
		t.Fatal("expected missing path to remain absent from VFS after failed Filescomfs CreateEx truncate")
	}
}

func TestRemoteFsFailedOpenAndTruncateDoNotPoisonParentPathLookup(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	fs.backend = &fakeRemoteBackend{
		findFunc: func(params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
			return files_sdk.File{}, files_sdk.ResponseError{Type: string(files_sdk.ErrFileNotFound)}
		},
	}

	missingFile := "/ghost.ai"
	childPath := "/ghost.ai/bar"

	if errc, _ := fs.Open(missingFile, fuse.O_RDWR); errc != -fuse.ENOENT {
		t.Fatalf("Open returned %d, want %d", errc, -fuse.ENOENT)
	}
	if errc := fs.Truncate(missingFile, 0, ^uint64(0)); errc != -fuse.ENOENT {
		t.Fatalf("Truncate returned %d, want %d", errc, -fuse.ENOENT)
	}

	errc, _ := fs.Create(childPath, fuse.O_RDWR|fuse.O_CREAT|fuse.O_EXCL, 0o644)
	if errc != -fuse.ENOENT {
		t.Fatalf("Create child returned %d, want %d", errc, -fuse.ENOENT)
	}
	if _, ok := vfs.fetch(missingFile); ok {
		t.Fatal("expected failed operations not to leave phantom parent file node behind")
	}
}

func TestRemoteFsPublicRenameMovesActiveWriteSession(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	oldPath := "/draft.ai"
	newPath := "/final.ai"

	errc, fh := fs.Create(oldPath, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	payload := []byte("rename-me")
	if n := fs.Write(oldPath, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}

	if errc := fs.Rename(oldPath, newPath); errc != 0 {
		t.Fatalf("Rename returned unexpected error: %d", errc)
	}

	buf := make([]byte, len(payload))
	if n := fs.Read(newPath, buf, 0, fh); n != len(payload) {
		t.Fatalf("Read after rename returned %d, want %d", n, len(payload))
	}
	if string(buf) != string(payload) {
		t.Fatalf("Read after rename returned %q, want %q", string(buf), string(payload))
	}

	if _, ok := vfs.fetch(oldPath); ok {
		t.Fatal("expected old path to be absent after rename")
	}
	if _, ok := vfs.fetch(newPath); !ok {
		t.Fatal("expected new path to exist after rename")
	}
}

func TestRemoteFsFlushAfterActiveRenameUploadsToNewPath(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	oldPath := "/draft.ai"
	newPath := "/final.ai"

	var uploadedPath string
	var uploaded []byte
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (int64, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return 0, err
		}
		uploadedPath = path
		uploaded = append([]byte(nil), data...)
		return int64(len(data)), nil
	}

	errc, fh := fs.Create(oldPath, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	payload := []byte("rename-then-flush")
	if n := fs.Write(oldPath, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}

	if errc := fs.Rename(oldPath, newPath); errc != 0 {
		t.Fatalf("Rename returned unexpected error: %d", errc)
	}
	if errc := fs.Flush(newPath, fh); errc != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errc)
	}

	if uploadedPath != newPath {
		t.Fatalf("uploaded path = %q, want %q", uploadedPath, newPath)
	}
	if !bytes.Equal(uploaded, payload) {
		t.Fatalf("uploaded payload = %q, want %q", string(uploaded), string(payload))
	}
	if _, ok := vfs.fetch(oldPath); ok {
		t.Fatal("expected old path to be absent after active rename")
	}
	if _, ok := vfs.fetch(newPath); !ok {
		t.Fatal("expected new path to exist after active rename")
	}
}

func TestRemoteFsPublicTruncateZeroSkipsHydrationAndResetsWorkingCopy(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/existing.ai"
	original := []byte("old remote contents")
	if _, err := cacheStore.Write(path, original, 0); err != nil {
		t.Fatalf("cache Write failed: %v", err)
	}

	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(original)),
		modTime:      time.Now(),
		creationTime: time.Now(),
	})
	node.setDownloadURI("https://example.invalid/download")

	errc, fh := fs.Open(path, fuse.O_RDWR)
	if errc != 0 {
		t.Fatalf("Open returned unexpected error: %d", errc)
	}

	if errc := fs.Truncate(path, 0, fh); errc != 0 {
		t.Fatalf("Truncate returned unexpected error: %d", errc)
	}

	replacement := []byte("new")
	if n := fs.Write(path, replacement, 0, fh); n != len(replacement) {
		t.Fatalf("Write returned %d, want %d", n, len(replacement))
	}

	buf := make([]byte, len(replacement))
	if n := fs.Read(path, buf, 0, fh); n != len(replacement) {
		t.Fatalf("Read returned %d, want %d", n, len(replacement))
	}
	if string(buf) != string(replacement) {
		t.Fatalf("Read returned %q, want %q", string(buf), string(replacement))
	}

	var stat fuse.Stat_t
	if errc := fs.Getattr(path, &stat, fh); errc != 0 {
		t.Fatalf("Getattr returned unexpected error: %d", errc)
	}
	if stat.Size != int64(len(replacement)) {
		t.Fatalf("Getattr size = %d, want %d", stat.Size, len(replacement))
	}
}

func TestRemoteFsPublicFlushUploadsWorkingCopyAndRefreshesCache(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	var uploaded []byte
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (int64, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return 0, err
		}
		uploaded = append([]byte(nil), data...)
		return int64(len(data)), nil
	}

	path := "/flush.ai"
	errc, fh := fs.Create(path, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	payload := []byte("flush-payload")
	if n := fs.Write(path, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}

	if errc := fs.Flush(path, fh); errc != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errc)
	}

	if !bytes.Equal(uploaded, payload) {
		t.Fatalf("uploaded payload = %q, want %q", string(uploaded), string(payload))
	}

	buf := make([]byte, len(payload))
	n, err := cacheStore.Read(path, buf, 0)
	if err != nil {
		t.Fatalf("cache Read failed: %v", err)
	}
	if n != len(payload) {
		t.Fatalf("cache Read returned %d, want %d", n, len(payload))
	}
	if !bytes.Equal(buf[:n], payload) {
		t.Fatalf("cache payload = %q, want %q", string(buf[:n]), string(payload))
	}

	if errc := fs.Release(path, fh); errc != 0 {
		t.Fatalf("Release returned unexpected error: %d", errc)
	}
	if node, ok := vfs.fetch(path); ok && node.getWriteSession() != nil {
		t.Fatal("expected write session to be cleared after successful release")
	}
}

func TestRemoteFsPublicFlushPoisonsSessionAfterUploadFailure(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	uploadErr := errors.New("upload failed")
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (int64, error) {
		_, _ = io.ReadAll(reader)
		return 0, uploadErr
	}

	path := "/poison.ai"
	errc, fh := fs.Create(path, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	payload := []byte("poison-me")
	if n := fs.Write(path, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}

	if errc := fs.Flush(path, fh); errc == 0 {
		t.Fatal("expected Flush to fail after upload error")
	}

	buf := make([]byte, len(payload))
	if n := fs.Read(path, buf, 0, fh); n != len(payload) {
		t.Fatalf("Read returned %d, want %d", n, len(payload))
	}
	if !bytes.Equal(buf, payload) {
		t.Fatalf("Read returned %q, want %q", string(buf), string(payload))
	}

	if n := fs.Write(path, []byte("x"), 0, fh); n >= 0 {
		t.Fatalf("expected poisoned session write to fail, got %d", n)
	}

	if errc := fs.Release(path, fh); errc == 0 {
		t.Fatal("expected Release to surface poisoned session error")
	}
	if node, ok := vfs.fetch(path); !ok || node.getWriteSession() == nil {
		t.Fatal("expected poisoned write session to remain after failed release")
	}
}

func TestRemoteFsInPlaceWritesAndFlushDoNotChangeSizeUntilTruncate(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/dense-save.indd"
	initialSize := 2166784
	finalSize := 2498560

	initial := bytes.Repeat([]byte("a"), initialSize)
	if _, err := cacheStore.Write(path, initial, 0); err != nil {
		t.Fatalf("cache Write failed: %v", err)
	}

	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(initialSize),
		modTime:      time.Now(),
		creationTime: time.Now(),
	})
	node.setDownloadURI("https://example.invalid/download")

	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (int64, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return 0, err
		}
		return int64(len(data)), nil
	}

	errc, fh := fs.Open(path, fuse.O_RDWR)
	if errc != 0 {
		t.Fatalf("Open returned unexpected error: %d", errc)
	}

	chunk := bytes.Repeat([]byte("b"), 4096)
	for offset := int64(8192); offset < 8192+(32*4096); offset += 4096 {
		if n := fs.Write(path, chunk, offset, fh); n != len(chunk) {
			t.Fatalf("Write at offset %d returned %d, want %d", offset, n, len(chunk))
		}
		if errc := fs.Flush(path, fh); errc != 0 {
			t.Fatalf("Flush at offset %d returned unexpected error: %d", offset, errc)
		}
	}

	var stat fuse.Stat_t
	if errc := fs.Getattr(path, &stat, fh); errc != 0 {
		t.Fatalf("Getattr after in-place writes returned unexpected error: %d", errc)
	}
	if stat.Size != int64(initialSize) {
		t.Fatalf("Getattr size after in-place writes = %d, want %d", stat.Size, initialSize)
	}

	if errc := fs.Truncate(path, int64(finalSize), fh); errc != 0 {
		t.Fatalf("Truncate returned unexpected error: %d", errc)
	}
	if errc := fs.Flush(path, fh); errc != 0 {
		t.Fatalf("Flush after truncate returned unexpected error: %d", errc)
	}

	if errc := fs.Getattr(path, &stat, fh); errc != 0 {
		t.Fatalf("Getattr after truncate returned unexpected error: %d", errc)
	}
	if stat.Size != int64(finalSize) {
		t.Fatalf("Getattr size after truncate = %d, want %d", stat.Size, finalSize)
	}
}

func TestRemoteFsGetattrKeepsStableMtimeDuringWriteSessionAndPublishesOnFlush(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/illustrator.ai"
	initial := []byte("existing remote contents")
	initialMtime := time.Now().Add(-time.Hour).Round(0)

	if _, err := cacheStore.Write(path, initial, 0); err != nil {
		t.Fatalf("cache Write failed: %v", err)
	}

	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(initial)),
		modTime:      initialMtime,
		creationTime: initialMtime,
	})
	node.setDownloadURI("https://example.invalid/download")

	var uploadedMtime time.Time
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (int64, error) {
		_, err := io.ReadAll(reader)
		if err != nil {
			return 0, err
		}
		uploadedMtime = mtime
		return int64(len(initial)), nil
	}

	errc, fh := fs.Open(path, fuse.O_RDWR)
	if errc != 0 {
		t.Fatalf("Open returned unexpected error: %d", errc)
	}

	if n := fs.Write(path, []byte("X"), 0, fh); n != 1 {
		t.Fatalf("Write returned %d, want 1", n)
	}

	var stat fuse.Stat_t
	if errc := fs.Getattr(path, &stat, fh); errc != 0 {
		t.Fatalf("Getattr during write session returned unexpected error: %d", errc)
	}

	gotSessionMtime := stat.Mtim.Time().Round(0)
	if !gotSessionMtime.Equal(initialMtime) {
		t.Fatalf("Getattr mtime during write session = %v, want stable %v", gotSessionMtime, initialMtime)
	}

	if errc := fs.Flush(path, fh); errc != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errc)
	}

	if uploadedMtime.IsZero() {
		t.Fatal("expected upload to receive a published mtime")
	}
	if !uploadedMtime.After(initialMtime) {
		t.Fatalf("uploaded mtime = %v, want after %v", uploadedMtime, initialMtime)
	}

	if errc := fs.Getattr(path, &stat, fh); errc != 0 {
		t.Fatalf("Getattr after flush returned unexpected error: %d", errc)
	}

	gotFlushedMtime := stat.Mtim.Time()
	if !gotFlushedMtime.Equal(uploadedMtime) {
		t.Fatalf("Getattr mtime after flush = %v, want %v", gotFlushedMtime, uploadedMtime)
	}
}

func TestFinalizeDeleteResetsSurvivingHandleState(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/design.ai"
	original := []byte("old illustrator contents")
	if _, err := cacheStore.Write(path, original, 0); err != nil {
		t.Fatalf("cache Write failed: %v", err)
	}

	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(original)),
		modTime:      time.Now(),
		creationTime: time.Now(),
	})
	node.setDownloadURI("https://example.invalid/download")
	node.extendTtl()

	// Simulate the effective state after Unlink succeeds: the path is removed
	// from the VFS/cache, but the open handle still points at the same node.
	fs.finalizeDelete(path)

	if _, ok := vfs.fetch(path); ok {
		t.Fatal("expected path to be removed from VFS after finalizeDelete")
	}
	if node.info.size != 0 {
		t.Fatalf("expected surviving handle node size to be reset, got %d", node.info.size)
	}
	if node.downloadUri != "" {
		t.Fatalf("expected surviving handle download URI to be cleared, got %q", node.downloadUri)
	}
	if !node.infoExpired() {
		t.Fatal("expected surviving handle node info to be expired after finalizeDelete")
	}

	buf := make([]byte, len(original))
	if n, _ := cacheStore.Read(path, buf, 0); n != 0 {
		t.Fatalf("expected cached content to be cleared after finalizeDelete, read %d bytes", n)
	}
}

func TestFinalizeDeleteIsIdempotent(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/design.ai"
	data := []byte("old illustrator contents")
	if _, err := cacheStore.Write(path, data, 0); err != nil {
		t.Fatalf("cache Write failed: %v", err)
	}

	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(data)),
		modTime:      time.Now(),
		creationTime: time.Now(),
	})
	node.setDownloadURI("https://example.invalid/download")

	fs.finalizeDelete(path)
	fs.finalizeDelete(path)

	if _, ok := vfs.fetch(path); ok {
		t.Fatal("expected path to remain absent after repeated finalizeDelete")
	}
	buf := make([]byte, len(data))
	if n, _ := cacheStore.Read(path, buf, 0); n != 0 {
		t.Fatalf("expected cached content to stay cleared after repeated finalizeDelete, read %d bytes", n)
	}
}

func TestRemoteFsCreateNodePreservesExistingCreationTimeForFiles(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/illustrator.ai"
	originalCreationTime := time.Now().Add(-2 * time.Hour).Round(0)
	updatedModTime := time.Now().Add(-time.Minute).Round(0)

	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         10,
		modTime:      originalCreationTime,
		creationTime: originalCreationTime,
	})

	fs.createNode(path, files_sdk.File{
		Path:      path,
		Type:      "file",
		Size:      20,
		Mtime:     &updatedModTime,
		CreatedAt: ptrToTime(time.Now().Round(0)),
	})

	refreshed, ok := vfs.fetch(path)
	if !ok {
		t.Fatal("expected node to exist after createNode refresh")
	}
	if !refreshed.info.creationTime.Equal(originalCreationTime) {
		t.Fatalf("creation time = %v, want preserved %v", refreshed.info.creationTime, originalCreationTime)
	}
	if !refreshed.info.modTime.Equal(updatedModTime) {
		t.Fatalf("mod time = %v, want %v", refreshed.info.modTime, updatedModTime)
	}
}

func ptrToTime(t time.Time) *time.Time {
	return &t
}

func TestFilescomfsLocalToRemoteRenameUploadsContentAndPreservesCreationTime(t *testing.T) {
	fs, remote, local, vfs, cacheStore := newTestFilescomfs(t)
	defer vfs.destroy()
	defer local.Destroy()

	remotePath := "/Shepard/1.ai"
	localTmpPath := "/Shepard/~ai-ee7c8b6f-000b-4477-94b5-b0b4935a5a94_.tmp"
	originalCreationTime := time.Now().Add(-3 * time.Hour).Round(0)

	existing := vfs.getOrCreate(remotePath, nodeTypeFile)
	existing.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         123,
		modTime:      originalCreationTime.Add(time.Hour),
		creationTime: originalCreationTime,
	})

	errc, fh := local.Create(localTmpPath, fuse.O_RDWR|fuse.O_CREAT, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}
	payload := []byte("illustrator-save-as")
	if n := local.Write(localTmpPath, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}
	if errc := local.Release(localTmpPath, fh); errc != 0 {
		t.Fatalf("Release returned unexpected error: %d", errc)
	}

	remote.backend = &fakeRemoteBackend{
		uploadFunc: func(opts ...file.UploadOption) error {
			return nil
		},
	}

	if errc := fs.Rename(localTmpPath, remotePath); errc != 0 {
		t.Fatalf("Rename returned unexpected error: %d", errc)
	}

	var stat fuse.Stat_t
	if errc := fs.Getattr(remotePath, &stat, ^uint64(0)); errc != 0 {
		t.Fatalf("Getattr returned unexpected error: %d", errc)
	}
	if stat.Size != int64(len(payload)) {
		t.Fatalf("Getattr size = %d, want %d", stat.Size, len(payload))
	}
	if stat.Nlink != 1 {
		t.Fatalf("Getattr nlink = %d, want 1", stat.Nlink)
	}
	if got := stat.Birthtim.Time().Round(0); !got.Equal(originalCreationTime) {
		t.Fatalf("Getattr birth time = %v, want preserved %v", got, originalCreationTime)
	}

	buf := make([]byte, len(payload))
	n, err := cacheStore.Read(remotePath, buf, 0)
	if err != nil {
		t.Fatalf("cache Read failed: %v", err)
	}
	if n != len(payload) {
		t.Fatalf("cache Read returned %d, want %d", n, len(payload))
	}
	if !bytes.Equal(buf[:n], payload) {
		t.Fatalf("cache payload = %q, want %q", string(buf[:n]), string(payload))
	}

	oldFq := local.fqPath(localTmpPath)
	deadline := time.Now().Add(2 * time.Second)
	for {
		_, err := os.Stat(oldFq)
		if errors.Is(err, os.ErrNotExist) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("expected local temp file to be removed: %s", oldFq)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
