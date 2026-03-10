package fsmount

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/mem"
	ff "github.com/Files-com/files-sdk-go/v3/fsmount/internal/flags"
	fsio "github.com/Files-com/files-sdk-go/v3/fsmount/internal/io"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
	fssync "github.com/Files-com/files-sdk-go/v3/fsmount/internal/sync"
	"github.com/winfsp/cgofuse/fuse"
)

type noopFSWriter struct{}

func (noopFSWriter) writeFile(path string, reader io.Reader, mtime time.Time, fh uint64) {}

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
		readyGates:     map[string]*cache.ReadyGate{},
		loadDirMutexes: fssync.NewPathMutex(),
	}

	return fs, vfs, cacheStore
}

func TestInitialContentForWriteDowngradesNotFound(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	node := &fsNode{
		path:     "/design.ai",
		cacheTTL: DefaultCacheTTL,
		logger:   &log.NoOpLogger{},
	}
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         128,
		modTime:      time.Now(),
		creationTime: time.Now(),
	})
	node.setDownloadURI("https://example.invalid/download")

	gate := cache.NewReadyGate()
	gate.Finish(files_sdk.ResponseError{
		Type:         "not-found/file",
		Title:        "Not Found",
		ErrorMessage: "Not Found. This may be related to your permissions.",
	}, -1)
	fs.readyGates[node.path] = gate

	reader, err := fs.initialContentForWrite(node.path, node, 42)
	if err != nil {
		t.Fatalf("initialContentForWrite returned unexpected error: %v", err)
	}
	if reader != nil {
		t.Fatal("expected no initial content reader when old version is missing")
	}
	if node.info.size != 0 {
		t.Fatalf("expected node size to be reset after missing initial content, got %d", node.info.size)
	}
	if node.downloadUri != "" {
		t.Fatalf("expected download URI to be cleared, got %q", node.downloadUri)
	}
	if !node.infoExpired() {
		t.Fatal("expected node info to be expired after missing initial content")
	}
}

func TestInitialContentForWritePropagatesNonNotFound(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	node := &fsNode{
		path:     "/design.ai",
		cacheTTL: DefaultCacheTTL,
		logger:   &log.NoOpLogger{},
	}
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         128,
		modTime:      time.Now(),
		creationTime: time.Now(),
	})

	gate := cache.NewReadyGate()
	gate.Finish(errors.New("cache preload failed"), -1)
	fs.readyGates[node.path] = gate

	reader, err := fs.initialContentForWrite(node.path, node, 42)
	if err == nil {
		t.Fatal("expected preload error to be returned")
	}
	if reader != nil {
		t.Fatal("expected no reader when preload fails")
	}
	if node.info.size != 128 {
		t.Fatalf("expected node size to remain unchanged on non-not-found error, got %d", node.info.size)
	}
}

func TestInitialContentForWriteReturnsCachedReader(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/design.ai"
	data := []byte("prior file contents")
	if _, err := cacheStore.Write(path, data, 0); err != nil {
		t.Fatalf("cache Write failed: %v", err)
	}

	node := &fsNode{
		path:     path,
		cacheTTL: DefaultCacheTTL,
		logger:   &log.NoOpLogger{},
	}
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(data)),
		modTime:      time.Now(),
		creationTime: time.Now(),
	})

	reader, err := fs.initialContentForWrite(path, node, 42)
	if err != nil {
		t.Fatalf("initialContentForWrite returned unexpected error: %v", err)
	}
	if reader == nil {
		t.Fatal("expected cache reader for cached initial content")
	}

	got, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	if string(got) != string(data) {
		t.Fatalf("unexpected cached content: got %q want %q", string(got), string(data))
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

func TestFinalizeDeleteAfterCanceledUploadClearsWriterAndUpload(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/design.ai"
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         64,
		modTime:      time.Now(),
		creationTime: time.Now(),
	})

	fh, _ := vfs.handles.Open(node, ff.NewFuseFlags(fuse.O_RDWR))
	node.markWriteIntent(fh)

	writer, created, err := node.ensureWriter(noopFSWriter{}, fh, func() (io.Reader, error) {
		return nil, nil
	}, func() fsio.CacheWriter {
		return func(data []byte, offset int64) (int, error) {
			return cacheStore.Write(path, data, offset)
		}
	})
	if err != nil {
		t.Fatalf("ensureWriter returned unexpected error: %v", err)
	}
	if !created {
		t.Fatal("expected writer to be created")
	}
	if _, err := writer.WriteAt([]byte("new"), 0); err != nil {
		t.Fatalf("WriteAt failed: %v", err)
	}

	_, cancel := context.WithCancel(context.Background())
	node.startUpload(path, cancel)

	node.cancelUpload()
	fs.finalizeDelete(path)

	if node.writerIsOpen() {
		t.Fatal("expected writer to be closed after canceled upload and finalizeDelete")
	}
	if node.upload != nil {
		t.Fatal("expected upload to be cleared after cancelUpload")
	}
	if node.info.size != 0 {
		t.Fatalf("expected node size to be reset after finalizeDelete, got %d", node.info.size)
	}
}

func TestSurvivingHandleCanCreateReplacementWriterAfterFinalizeDelete(t *testing.T) {
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

	fh, _ := vfs.handles.Open(node, ff.NewFuseFlags(fuse.O_RDWR))
	node.markWriteIntent(fh)

	fs.finalizeDelete(path)

	reader, err := fs.initialContentForWrite(path, node, fh)
	if err != nil {
		t.Fatalf("initialContentForWrite returned unexpected error: %v", err)
	}
	if reader != nil {
		t.Fatal("expected no initial content after finalizeDelete")
	}

	writer, created, err := node.ensureWriter(noopFSWriter{}, fh, func() (io.Reader, error) {
		return fs.initialContentForWrite(path, node, fh)
	}, func() fsio.CacheWriter {
		return func(data []byte, offset int64) (int, error) {
			return cacheStore.Write(path, data, offset)
		}
	})
	if err != nil {
		t.Fatalf("ensureWriter returned unexpected error: %v", err)
	}
	if !created {
		t.Fatal("expected writer to be created for replacement write")
	}

	replacement := []byte("new illustrator contents")
	n, err := writer.WriteAt(replacement, 0)
	if err != nil {
		t.Fatalf("WriteAt failed: %v", err)
	}
	if n != len(replacement) {
		t.Fatalf("unexpected write length: got %d want %d", n, len(replacement))
	}
}
