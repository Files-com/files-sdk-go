package fsmount

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/fsmount/events"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/mem"
	lim "github.com/Files-com/files-sdk-go/v3/fsmount/internal/limit"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
	fssync "github.com/Files-com/files-sdk-go/v3/fsmount/internal/sync"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/winfsp/cgofuse/fuse"
)

type timeoutUploadError struct{}

func (timeoutUploadError) Error() string   { return "raw timeout detail" }
func (timeoutUploadError) Timeout() bool   { return true }
func (timeoutUploadError) Temporary() bool { return true }

type blockingDeleteCacheStore struct {
	cacheStore
	path    string
	started chan struct{}
	release chan struct{}
	once    sync.Once
}

func (s *blockingDeleteCacheStore) Delete(path string) bool {
	if path == s.path {
		s.once.Do(func() {
			close(s.started)
		})
		<-s.release
	}
	return s.cacheStore.Delete(path)
}

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
		log:             &log.NoOpLogger{},
		vfs:             vfs,
		cacheStore:      cacheStore,
		disableLocking:  true,
		lockMap:         make(map[string]*lockInfo),
		readyGates:      map[string]*cache.ReadyGate{},
		gatePathMutexes: fssync.NewPathMutex(),
		loadDirMutexes:  fssync.NewPathMutex(),
		backend:         &fakeRemoteBackend{},
		ops: lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
			lim.FuseOpDownload: downloadOpLimit,
			lim.FuseOpUpload:   uploadOpLimit,
			lim.FuseOpOther:    otherOpLimit,
		}, globalOpLimit),
		bufferPool: fssync.NewPool(func() []byte {
			return make([]byte, cacheWriteSize)
		}),
	}

	return fs, vfs, cacheStore
}

type captureEventPublisher struct {
	mu     sync.Mutex
	events []events.MountEvent
}

func (p *captureEventPublisher) Publish(event events.MountEvent) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, event)
}

func (p *captureEventPublisher) transferEvents() []events.TransferEvent {
	p.mu.Lock()
	defer p.mu.Unlock()

	transfers := make([]events.TransferEvent, 0, len(p.events))
	for _, event := range p.events {
		if transfer, ok := event.(events.TransferEvent); ok {
			transfers = append(transfers, transfer)
		}
	}
	return transfers
}

func (p *captureEventPublisher) waitForTransferEvents(t *testing.T, count int) []events.TransferEvent {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for {
		transfers := p.transferEvents()
		if len(transfers) >= count {
			return transfers
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %d transfer events, got %d", count, len(transfers))
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func assertStableTransferID(t *testing.T, transfers []events.TransferEvent) {
	t.Helper()
	if len(transfers) == 0 {
		t.Fatal("expected transfer events")
	}
	id := transfers[0].ID
	if id == "" {
		t.Fatal("expected non-empty transfer ID")
	}
	for i, transfer := range transfers {
		if transfer.ID != id {
			t.Fatalf("transfer event %d ID = %q, want %q", i, transfer.ID, id)
		}
	}
}

type fakeRemoteBackend struct {
	findFunc             func(params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)
	listForFunc          func(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (remoteFileIter, error)
	uploadFunc           func(opts ...file.UploadOption) error
	uploadWithResumeFunc func(opts ...file.UploadOption) (file.UploadResumable, error)
	moveFunc             func(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error)
	waitFunc             func(action files_sdk.FileAction, status func(files_sdk.FileMigration), opts ...files_sdk.RequestResponseOption) (files_sdk.FileMigration, error)
	downloadToFileFunc   func(params files_sdk.FileDownloadParams, filePath string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)
	downloadFunc         func(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)
	deleteFunc           func(params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) error
	createLockFunc       func(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.Lock, error)
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
	if b.listForFunc != nil {
		return b.listForFunc(params, opts...)
	}
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
	if b.uploadWithResumeFunc != nil {
		return b.uploadWithResumeFunc(opts...)
	}
	if b.uploadFunc != nil {
		return file.UploadResumable{}, b.uploadFunc(opts...)
	}
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
	if b.downloadFunc != nil {
		return b.downloadFunc(params, opts...)
	}
	return files_sdk.File{}, nil
}

func (b *fakeRemoteBackend) createLock(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.Lock, error) {
	if b.createLockFunc != nil {
		return b.createLockFunc(params, opts...)
	}
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

func fakeDownloadResponse(payload []byte, reportedSize int64) func(files_sdk.FileDownloadParams, ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return func(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
		resp := &http.Response{
			StatusCode:    http.StatusOK,
			Body:          io.NopCloser(bytes.NewReader(payload)),
			ContentLength: int64(len(payload)),
		}
		if _, err := files_sdk.BuildResponse(resp, opts...); err != nil {
			return files_sdk.File{}, err
		}
		return files_sdk.File{
			Path: params.File.Path,
			Type: "file",
			Size: reportedSize,
		}, nil
	}
}

type blockingDownloadReader struct {
	payload      []byte
	firstChunk   int
	offset       int
	firstWritten chan struct{}
	release      chan struct{}
	once         sync.Once
}

type cancelableBlockingDownloadReader struct {
	payload []byte
	ctx     context.Context
	started chan struct{}
	release chan struct{}
	once    sync.Once
	offset  int
}

func (r *cancelableBlockingDownloadReader) Read(p []byte) (int, error) {
	if r.offset >= len(r.payload) {
		return 0, io.EOF
	}
	r.once.Do(func() { close(r.started) })
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	case <-r.release:
		n := copy(p, r.payload[r.offset:])
		r.offset += n
		return n, nil
	}
}

func newBlockingDownloadReader(payload []byte, firstChunk int) *blockingDownloadReader {
	return &blockingDownloadReader{
		payload:      payload,
		firstChunk:   firstChunk,
		firstWritten: make(chan struct{}),
		release:      make(chan struct{}),
	}
}

func (r *blockingDownloadReader) Read(p []byte) (int, error) {
	if r.offset >= len(r.payload) {
		return 0, io.EOF
	}

	if r.offset == 0 {
		n := copy(p, r.payload[:min(r.firstChunk, len(r.payload))])
		r.offset += n
		r.once.Do(func() { close(r.firstWritten) })
		return n, nil
	}

	<-r.release
	n := copy(p, r.payload[r.offset:])
	r.offset += n
	return n, nil
}

func waitForPartialCacheBytes(t *testing.T, cacheStore cacheStore, path string, want int) []byte {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	buf := make([]byte, want)
	for {
		n, err := cacheStore.ReadPartial(path, buf, 0)
		if err != nil {
			t.Fatalf("partial cache Read failed: %v", err)
		}
		if n >= want {
			return buf[:n]
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %d partial cached bytes at %s, got %d", want, path, n)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

type notifyingCache struct {
	cacheStore
	writePartialPath string
	wrote            chan struct{}
	once             sync.Once
}

func (c *notifyingCache) WritePartial(path string, buff []byte, ofst int64) (int, error) {
	n, err := c.cacheStore.WritePartial(path, buff, ofst)
	if path == c.writePartialPath && n > 0 {
		c.once.Do(func() { close(c.wrote) })
	}
	return n, err
}

func newTestDiskCache(t *testing.T) *disk.DiskCache {
	t.Helper()

	cacheStore, err := disk.NewDiskCache(t.TempDir())
	if err != nil {
		t.Fatalf("NewDiskCache failed: %v", err)
	}
	return cacheStore
}

type captureMountLogger struct {
	mu           sync.Mutex
	visibleLines []string
}

func (l *captureMountLogger) Debug(format string, v ...any) { l.append("DEBUG", format, v...) }
func (l *captureMountLogger) Error(format string, v ...any) { l.append("ERROR", format, v...) }
func (l *captureMountLogger) Info(format string, v ...any)  { l.append("INFO", format, v...) }
func (l *captureMountLogger) Trace(format string, v ...any) { l.append("TRACE", format, v...) }
func (l *captureMountLogger) Warn(format string, v ...any)  { l.append("WARN", format, v...) }

func (l *captureMountLogger) append(level, format string, v ...any) {
	line := fmt.Sprintf("%s "+format, append([]any{level}, v...)...)
	l.mu.Lock()
	defer l.mu.Unlock()
	if level != "DEBUG" && level != "TRACE" {
		l.visibleLines = append(l.visibleLines, line)
	}
}

func (l *captureMountLogger) visibleJoined() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return strings.Join(l.visibleLines, "\n")
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
		log:             &log.NoOpLogger{},
		vfs:             vfs,
		cacheStore:      cacheStore,
		disableLocking:  true,
		lockMap:         make(map[string]*lockInfo),
		readyGates:      map[string]*cache.ReadyGate{},
		gatePathMutexes: fssync.NewPathMutex(),
		loadDirMutexes:  fssync.NewPathMutex(),
		backend:         &fakeRemoteBackend{},
		ops: lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
			lim.FuseOpDownload: downloadOpLimit,
			lim.FuseOpUpload:   uploadOpLimit,
			lim.FuseOpOther:    otherOpLimit,
		}, globalOpLimit),
		bufferPool: fssync.NewPool(func() []byte {
			return make([]byte, cacheWriteSize)
		}),
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

func TestTransferReporterThrottlesProgressEvents(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	publisher := &captureEventPublisher{}
	fs.events = publisher

	reporter := fs.newTransferReporter(events.TransferDirectionUpload, "/throttle.bin", transferProgressMinBytes*2)
	reporter.Queued()
	reporter.Progress(1)
	reporter.Progress(1)
	reporter.Progress(transferProgressMinBytes)
	reporter.Complete(transferProgressMinBytes + 2)

	transfers := publisher.transferEvents()
	if len(transfers) != 4 {
		t.Fatalf("transfer event count = %d, want 4", len(transfers))
	}
	assertStableTransferID(t, transfers)

	wantStatuses := []events.TransferStatus{
		events.TransferStatusQueued,
		events.TransferStatusTransferring,
		events.TransferStatusTransferring,
		events.TransferStatusComplete,
	}
	for i, want := range wantStatuses {
		if transfers[i].Status != want {
			t.Fatalf("transfer event %d status = %q, want %q", i, transfers[i].Status, want)
		}
	}
}

func TestTransferReporterRewindsNegativeProgress(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	publisher := &captureEventPublisher{}
	fs.events = publisher

	reporter := fs.newTransferReporter(events.TransferDirectionUpload, "/retry.bin", transferProgressMinBytes*2)
	reporter.Queued()
	reporter.Progress(transferProgressMinBytes)
	reporter.Progress(-transferProgressMinBytes * 2)
	reporter.Progress(transferProgressMinBytes)
	reporter.Complete(transferredBytesUnchanged)

	transfers := publisher.transferEvents()
	if len(transfers) != 4 {
		t.Fatalf("transfer event count = %d, want 4", len(transfers))
	}
	assertStableTransferID(t, transfers)

	if transfers[1].TransferredBytes != transferProgressMinBytes {
		t.Fatalf("first progress transferred bytes = %d, want %d", transfers[1].TransferredBytes, transferProgressMinBytes)
	}
	if transfers[2].TransferredBytes != transferProgressMinBytes {
		t.Fatalf("second progress transferred bytes = %d, want %d", transfers[2].TransferredBytes, transferProgressMinBytes)
	}
	last := transfers[len(transfers)-1]
	if last.Status != events.TransferStatusComplete {
		t.Fatalf("last status = %q, want complete", last.Status)
	}
	if last.TransferredBytes != transferProgressMinBytes {
		t.Fatalf("complete transferred bytes = %d, want %d", last.TransferredBytes, transferProgressMinBytes)
	}
}

func TestTransferReporterNoOpPublisherIsOptional(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	fs.events = &events.NoOpEventPublisher{}
	reporter := fs.newTransferReporter(events.TransferDirectionUpload, "/noop.bin", 10)
	reporter.Queued()
	reporter.Progress(5)
	reporter.Complete(10)
}

func TestRemoteFsWriteSessionUploadPublishesTransferEvents(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	publisher := &captureEventPublisher{}
	fs.events = publisher

	path := "/illustrator-transfer.ai"
	payload := bytes.Repeat([]byte("u"), int(transferProgressMinBytes)+1)
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		progress := fs.uploadProgressFunc(node)
		progress(int64(len(data) / 2))
		progress(int64(len(data) - len(data)/2))
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	errc, fh := fs.Create(path, fuse.O_RDWR|fuse.O_CREAT, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}
	defer fs.Release(path, fh)

	if n := fs.Write(path, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}
	if errc := fs.Flush(path, fh); errc != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errc)
	}

	transfers := publisher.waitForTransferEvents(t, 3)
	assertStableTransferID(t, transfers)
	if transfers[0].Status != events.TransferStatusQueued {
		t.Fatalf("first status = %q, want queued", transfers[0].Status)
	}
	if transfers[1].Status != events.TransferStatusTransferring {
		t.Fatalf("second status = %q, want transferring", transfers[1].Status)
	}
	last := transfers[len(transfers)-1]
	if last.Status != events.TransferStatusComplete {
		t.Fatalf("last status = %q, want complete", last.Status)
	}
	if last.Direction != events.TransferDirectionUpload {
		t.Fatalf("direction = %q, want upload", last.Direction)
	}
	if last.LocalPath != fs.localPath(path) {
		t.Fatalf("local path = %q, want %q", last.LocalPath, fs.localPath(path))
	}
	if last.RemotePath != fs.remotePath(path) {
		t.Fatalf("remote path = %q, want %q", last.RemotePath, fs.remotePath(path))
	}
	if last.Size != int64(len(payload)) || last.TransferredBytes != int64(len(payload)) {
		t.Fatalf("last size/transferred = %d/%d, want %d", last.Size, last.TransferredBytes, len(payload))
	}
	if last.EndedAt.IsZero() {
		t.Fatal("expected complete event to include EndedAt")
	}
}

func TestRemoteFsProviderWriteSessionUploadPublishesTransferProgress(t *testing.T) {
	_, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	publisher := &captureEventPublisher{}
	provider := &testProviderBackend{
		writeFunc: func(_ context.Context, path string, reader io.Reader, size int64, modTime time.Time) (ProviderEntry, error) {
			n, err := io.Copy(io.Discard, reader)
			if err != nil {
				return ProviderEntry{}, err
			}
			return ProviderEntry{
				Path:    path,
				Type:    ProviderTypeFile,
				Size:    n,
				ModTime: modTime,
			}, nil
		},
	}
	fs, err := newRemoteFs(MountParams{
		Config:          &files_sdk.Config{},
		ProviderBackend: provider,
		TmpFsPath:       t.TempDir(),
		EventPublisher:  publisher,
	}, vfs, &log.NoOpLogger{}, cacheStore)
	if err != nil {
		t.Fatalf("newRemoteFs failed: %v", err)
	}

	path := "/provider-write-session.bin"
	payload := bytes.Repeat([]byte("p"), int(transferProgressMinBytes)+1)
	errc, fh := fs.Create(path, fuse.O_RDWR|fuse.O_CREAT, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}
	defer fs.Release(path, fh)

	if n := fs.Write(path, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}
	if errc := fs.Flush(path, fh); errc != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errc)
	}

	transfers := publisher.waitForTransferEvents(t, 3)
	assertStableTransferID(t, transfers)
	if transfers[0].Status != events.TransferStatusQueued {
		t.Fatalf("first status = %q, want queued", transfers[0].Status)
	}
	if transfers[1].Status != events.TransferStatusTransferring {
		t.Fatalf("second status = %q, want transferring", transfers[1].Status)
	}
	last := transfers[len(transfers)-1]
	if last.Status != events.TransferStatusComplete {
		t.Fatalf("last status = %q, want complete", last.Status)
	}
	if last.TransferredBytes != int64(len(payload)) {
		t.Fatalf("complete transferred bytes = %d, want %d", last.TransferredBytes, len(payload))
	}
}

func TestRemoteFsProviderWriteSessionUploadUsesProviderReturnedModTime(t *testing.T) {
	_, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	providerModTime := time.Date(2026, 6, 8, 14, 30, 0, 0, time.UTC)
	provider := &testProviderBackend{
		writeFunc: func(_ context.Context, path string, reader io.Reader, size int64, modTime time.Time) (ProviderEntry, error) {
			n, err := io.Copy(io.Discard, reader)
			if err != nil {
				return ProviderEntry{}, err
			}
			return ProviderEntry{
				Path:    path,
				Type:    ProviderTypeFile,
				Size:    n,
				ModTime: providerModTime,
			}, nil
		},
	}
	fs, err := newRemoteFs(MountParams{
		Config:          &files_sdk.Config{},
		ProviderBackend: provider,
		TmpFsPath:       t.TempDir(),
	}, vfs, &log.NoOpLogger{}, cacheStore)
	if err != nil {
		t.Fatalf("newRemoteFs failed: %v", err)
	}

	path := "/provider-returned-mtime.bin"
	payload := []byte("provider returned mtime")
	errc, fh := fs.Create(path, fuse.O_RDWR|fuse.O_CREAT, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}
	defer fs.Release(path, fh)

	if n := fs.Write(path, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}
	if errc := fs.Flush(path, fh); errc != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errc)
	}

	node, ok := vfs.fetch(path)
	if !ok {
		t.Fatal("expected node to exist after Flush")
	}
	if !node.info.modTime.Equal(providerModTime) {
		t.Fatalf("node mtime = %v, want provider mtime %v", node.info.modTime, providerModTime)
	}

	buf := make([]byte, len(payload))
	n, err := cacheStore.ReadComplete(path, cacheEntryMetadata(path, int64(len(payload)), providerModTime), buf, 0)
	if err != nil {
		t.Fatalf("cache ReadComplete failed: %v", err)
	}
	if n != len(payload) || !bytes.Equal(buf[:n], payload) {
		t.Fatalf("cached bytes n=%d data=%q, want %q", n, string(buf[:n]), string(payload))
	}
}

func TestRemoteFsProviderRenameUploadPublishesTransferProgress(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	publisher := &captureEventPublisher{}
	fs.events = publisher

	payload := bytes.Repeat([]byte("p"), int(transferProgressMinBytes)+1)
	src := filepath.Join(t.TempDir(), "provider-rename-progress.bin")
	if err := os.WriteFile(src, payload, 0o600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	dst := "/provider-rename-progress.bin"

	fs.providerBackend = &testProviderBackend{
		writeFunc: func(_ context.Context, path string, reader io.Reader, size int64, modTime time.Time) (ProviderEntry, error) {
			n, err := io.Copy(io.Discard, reader)
			if err != nil {
				return ProviderEntry{}, err
			}
			return ProviderEntry{
				Path:    path,
				Type:    ProviderTypeFile,
				Size:    n,
				ModTime: modTime,
			}, nil
		},
	}

	if err := fs.uploadFile(src, dst); err != nil {
		t.Fatalf("uploadFile returned error: %v", err)
	}

	transfers := publisher.waitForTransferEvents(t, 3)
	assertStableTransferID(t, transfers)
	if transfers[0].Status != events.TransferStatusQueued {
		t.Fatalf("first status = %q, want queued", transfers[0].Status)
	}
	if transfers[1].Status != events.TransferStatusTransferring {
		t.Fatalf("second status = %q, want transferring", transfers[1].Status)
	}
	last := transfers[len(transfers)-1]
	if last.Status != events.TransferStatusComplete {
		t.Fatalf("last status = %q, want complete", last.Status)
	}
	if last.TransferredBytes != int64(len(payload)) {
		t.Fatalf("complete transferred bytes = %d, want %d", last.TransferredBytes, len(payload))
	}
}

func TestRemoteFsWriteSessionUploadPublishesErroredTransferEvent(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	publisher := &captureEventPublisher{}
	fs.events = publisher

	uploadErr := errors.New("backend exploded")
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		return uploadedFileMetadata{}, uploadErr
	}

	path := "/upload-error.ai"
	errc, fh := fs.Create(path, fuse.O_RDWR|fuse.O_CREAT, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}
	defer fs.Release(path, fh)

	if n := fs.Write(path, []byte("data"), 0, fh); n != 4 {
		t.Fatalf("Write returned %d, want 4", n)
	}
	if errc := fs.Flush(path, fh); errc != -fuse.EIO {
		t.Fatalf("Flush returned %d, want %d", errc, -fuse.EIO)
	}

	transfers := publisher.waitForTransferEvents(t, 2)
	assertStableTransferID(t, transfers)
	last := transfers[len(transfers)-1]
	if last.Status != events.TransferStatusErrored {
		t.Fatalf("last status = %q, want errored", last.Status)
	}
	if !strings.Contains(last.Error, uploadErr.Error()) {
		t.Fatalf("last error = %q, want to contain %q", last.Error, uploadErr.Error())
	}
}

func TestRemoteFsFillCachePublishesDownloadTransferEvents(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	publisher := &captureEventPublisher{}
	fs.events = publisher

	path := "/download-transfer.bin"
	payload := bytes.Repeat([]byte("d"), cacheWriteSize+3)
	modTime := time.Now().Add(-time.Minute).Round(0)
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(payload)),
		modTime:      modTime,
		creationTime: modTime,
	})
	fs.backend = &fakeRemoteBackend{
		downloadFunc: fakeDownloadResponse(payload, int64(len(payload))),
	}

	readyGate := cache.NewReadyGate()
	fs.fillCache(context.Background(), path, "https://example.invalid/download", cacheEntryMetadata(path, int64(len(payload)), modTime), readyGate, 0, false)

	transfers := publisher.waitForTransferEvents(t, 3)
	assertStableTransferID(t, transfers)
	if transfers[0].Status != events.TransferStatusQueued {
		t.Fatalf("first status = %q, want queued", transfers[0].Status)
	}
	if transfers[1].Status != events.TransferStatusTransferring {
		t.Fatalf("second status = %q, want transferring", transfers[1].Status)
	}
	last := transfers[len(transfers)-1]
	if last.Status != events.TransferStatusComplete {
		t.Fatalf("last status = %q, want complete", last.Status)
	}
	if last.Direction != events.TransferDirectionDownload {
		t.Fatalf("direction = %q, want download", last.Direction)
	}
	if last.LocalPath != fs.localPath(path) {
		t.Fatalf("local path = %q, want %q", last.LocalPath, fs.localPath(path))
	}
	if last.RemotePath != fs.remotePath(path) {
		t.Fatalf("remote path = %q, want %q", last.RemotePath, fs.remotePath(path))
	}
	if last.Size != int64(len(payload)) || last.TransferredBytes != int64(len(payload)) {
		t.Fatalf("last size/transferred = %d/%d, want %d", last.Size, last.TransferredBytes, len(payload))
	}
}

func TestRemoteFsDownloadFilePublishesTransferEventsWithStatFallback(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	publisher := &captureEventPublisher{}
	fs.events = publisher

	src := "/download-for-rename.bin"
	payload := []byte("downloaded through cross-boundary rename")
	dst := filepath.Join(t.TempDir(), "download-for-rename.bin")
	eventLocalPath := filepath.Join(t.TempDir(), "mount", "download-for-rename.bin")
	fs.backend = &fakeRemoteBackend{
		downloadToFileFunc: func(params files_sdk.FileDownloadParams, filePath string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
			if params.Path != src {
				t.Fatalf("download path = %q, want %q", params.Path, src)
			}
			if filePath != dst {
				t.Fatalf("download destination = %q, want %q", filePath, dst)
			}
			if err := os.WriteFile(filePath, payload, 0o600); err != nil {
				return files_sdk.File{}, err
			}
			return files_sdk.File{
				Path: params.Path,
				Type: "file",
				Size: 0,
			}, nil
		},
	}

	if err := fs.downloadFile(src, dst, eventLocalPath); err != nil {
		t.Fatalf("downloadFile returned error: %v", err)
	}

	transfers := publisher.transferEvents()
	if len(transfers) != 2 {
		t.Fatalf("transfer event count = %d, want 2", len(transfers))
	}
	assertStableTransferID(t, transfers)
	if transfers[0].Status != events.TransferStatusQueued {
		t.Fatalf("first status = %q, want queued", transfers[0].Status)
	}
	last := transfers[len(transfers)-1]
	if last.Status != events.TransferStatusComplete {
		t.Fatalf("last status = %q, want complete", last.Status)
	}
	if last.Direction != events.TransferDirectionDownload {
		t.Fatalf("direction = %q, want download", last.Direction)
	}
	if last.LocalPath != eventLocalPath {
		t.Fatalf("local path = %q, want %q", last.LocalPath, eventLocalPath)
	}
	if last.RemotePath != fs.remotePath(src) {
		t.Fatalf("remote path = %q, want %q", last.RemotePath, fs.remotePath(src))
	}
	if last.Size != int64(len(payload)) || last.TransferredBytes != int64(len(payload)) {
		t.Fatalf("last size/transferred = %d/%d, want %d", last.Size, last.TransferredBytes, len(payload))
	}
	if last.EndedAt.IsZero() {
		t.Fatal("expected complete event to include EndedAt")
	}
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

func TestRemoteFsReadIgnoresUncommittedDiskCacheAfterRestart(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	cacheStore := newTestDiskCache(t)
	fs.cacheStore = cacheStore

	path := "/large.bin"
	remotePayload := bytes.Repeat([]byte("r"), cacheWriteSize+17)
	partialPayload := remotePayload[:cacheWriteSize/2]
	if _, err := cacheStore.Write(path, partialPayload, 0); err != nil {
		t.Fatalf("cache Write failed: %v", err)
	}

	modTime := time.Now().Add(-time.Minute).Round(0)
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(remotePayload)),
		modTime:      modTime,
		creationTime: modTime,
	})
	node.setDownloadURI("https://example.invalid/download")

	downloadCalls := 0
	fs.backend = &fakeRemoteBackend{
		downloadFunc: func(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
			downloadCalls++
			return fakeDownloadResponse(remotePayload, int64(len(remotePayload)))(params, opts...)
		},
	}

	errc, fh := fs.Open(path, fuse.O_RDONLY)
	if errc != 0 {
		t.Fatalf("Open returned unexpected error: %d", errc)
	}
	defer fs.Release(path, fh)

	buf := make([]byte, len(remotePayload))
	n := fs.Read(path, buf, 0, fh)
	if n != len(remotePayload) {
		t.Fatalf("Read returned %d, want %d", n, len(remotePayload))
	}
	if !bytes.Equal(buf[:n], remotePayload) {
		t.Fatalf("Read returned stale cache payload prefix %q, want remote payload prefix %q", string(buf[:16]), string(remotePayload[:16]))
	}
	if downloadCalls != 1 {
		t.Fatalf("download calls = %d, want 1", downloadCalls)
	}
}

func TestRemoteFsCommitReplacesLargerPinnedEntryWithSmallerContent(t *testing.T) {
	run := func(t *testing.T, cacheStore cacheStore) {
		t.Helper()

		fs, vfs, _ := newTestRemoteFs(t)
		defer vfs.destroy()
		fs.cacheStore = cacheStore

		path := "/shrink.bin"
		oldPayload := bytes.Repeat([]byte("o"), cacheWriteSize+1024)
		newPayload := bytes.Repeat([]byte("n"), cacheWriteSize/2)
		oldMtime := time.Now().Add(-time.Hour).Round(0)
		newMtime := oldMtime.Add(time.Second)
		if _, err := cacheStore.Write(path, oldPayload, 0); err != nil {
			t.Fatalf("old cache Write failed: %v", err)
		}
		if err := cacheStore.Commit(path, cacheEntryMetadata(path, int64(len(oldPayload)), oldMtime)); err != nil {
			t.Fatalf("old cache Commit failed: %v", err)
		}

		node := vfs.getOrCreate(path, nodeTypeFile)
		node.updateInfo(fsNodeInfo{
			nodeType:     nodeTypeFile,
			size:         int64(len(newPayload)),
			modTime:      newMtime,
			creationTime: newMtime,
		})
		node.setDownloadURI("https://example.invalid/download")
		fs.backend = &fakeRemoteBackend{
			downloadFunc: fakeDownloadResponse(newPayload, int64(len(newPayload))),
		}

		errc, fh := fs.Open(path, fuse.O_RDONLY)
		if errc != 0 {
			t.Fatalf("Open returned unexpected error: %d", errc)
		}
		defer fs.Release(path, fh)

		buf := make([]byte, len(newPayload))
		n := fs.Read(path, buf, 0, fh)
		if n != len(newPayload) {
			t.Fatalf("Read returned %d, want %d", n, len(newPayload))
		}
		if !bytes.Equal(buf[:n], newPayload) {
			t.Fatalf("Read returned %q, want %q", string(buf[:n]), string(newPayload))
		}

		committed := make([]byte, len(newPayload))
		n, err := cacheStore.ReadComplete(path, cacheEntryMetadata(path, int64(len(newPayload)), newMtime), committed, 0)
		if err != nil {
			t.Fatalf("ReadComplete failed: %v", err)
		}
		if n != len(newPayload) || !bytes.Equal(committed[:n], newPayload) {
			t.Fatalf("ReadComplete returned n=%d payload=%q, want %q", n, string(committed[:n]), string(newPayload))
		}
	}

	t.Run("disk", func(t *testing.T) {
		run(t, newTestDiskCache(t))
	})

	t.Run("memory", func(t *testing.T) {
		cacheStore, err := mem.NewMemoryCache()
		if err != nil {
			t.Fatalf("NewMemoryCache failed: %v", err)
		}
		run(t, cacheStore)
	})
}

func TestRemoteFsReadFailsAndDoesNotCommitShortDownload(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	cacheStore := newTestDiskCache(t)
	fs.cacheStore = cacheStore

	path := "/short-download.bin"
	shortPayload := bytes.Repeat([]byte("s"), cacheWriteSize/2)
	expectedSize := int64(cacheWriteSize + 1)
	modTime := time.Now().Add(-time.Minute).Round(0)
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         expectedSize,
		modTime:      modTime,
		creationTime: modTime,
	})
	node.setDownloadURI("https://example.invalid/download")

	fs.backend = &fakeRemoteBackend{
		downloadFunc: fakeDownloadResponse(shortPayload, expectedSize),
	}

	errc, fh := fs.Open(path, fuse.O_RDONLY)
	if errc != 0 {
		t.Fatalf("Open returned unexpected error: %d", errc)
	}
	defer fs.Release(path, fh)

	buf := make([]byte, expectedSize)
	n := fs.Read(path, buf, 0, fh)
	if n != -fuse.EIO {
		t.Fatalf("Read returned %d, want %d for short download", n, -fuse.EIO)
	}

	raw := make([]byte, len(shortPayload))
	if n, err := cacheStore.Read(path, raw, 0); err != nil || n != 0 {
		t.Fatalf("cache Read after short download returned n=%d err=%v, want empty cache", n, err)
	}
}

func TestRemoteFsReadDeletesPartialCacheAfterDownloadWaitersDrain(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	cacheStore := newTestDiskCache(t)
	fs.cacheStore = cacheStore

	path := "/cleanup-partial.bin"
	payload := bytes.Repeat([]byte("p"), cacheWriteSize+11)
	modTime := time.Now().Add(-time.Minute).Round(0)
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(payload)),
		modTime:      modTime,
		creationTime: modTime,
	})
	node.setDownloadURI("https://example.invalid/download")

	fs.backend = &fakeRemoteBackend{
		downloadFunc: fakeDownloadResponse(payload, int64(len(payload))),
	}

	errc, fh := fs.Open(path, fuse.O_RDONLY)
	if errc != 0 {
		t.Fatalf("Open returned unexpected error: %d", errc)
	}
	defer fs.Release(path, fh)

	buf := make([]byte, len(payload))
	if n := fs.Read(path, buf, 0, fh); n != len(payload) {
		t.Fatalf("Read returned %d, want %d", n, len(payload))
	}

	partial := make([]byte, len(payload))
	if n, err := cacheStore.ReadPartial(path, partial, 0); err != nil || n != 0 {
		t.Fatalf("partial cache Read after drained read returned n=%d err=%v, want empty partial cache", n, err)
	}

	committed := make([]byte, len(payload))
	if n, err := cacheStore.ReadComplete(path, cacheEntryMetadata(path, int64(len(payload)), modTime), committed, 0); err != nil || n != len(payload) {
		t.Fatalf("committed cache ReadComplete returned n=%d err=%v, want %d", n, err, len(payload))
	}
}

func TestRemoteFsPartialNamespaceDoesNotCollideWithSuffixPath(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	cacheStore := newTestDiskCache(t)
	fs.cacheStore = cacheStore

	path := "/namespace-source.bin"
	suffixPath := path + ".filescom-cache-partial"
	suffixPayload := []byte("real file cache entry")
	suffixModTime := time.Now().Add(-2 * time.Minute).Round(0)
	if _, err := cacheStore.Write(suffixPath, suffixPayload, 0); err != nil {
		t.Fatalf("suffix cache Write failed: %v", err)
	}
	if err := cacheStore.Commit(suffixPath, cacheEntryMetadata(suffixPath, int64(len(suffixPayload)), suffixModTime)); err != nil {
		t.Fatalf("suffix cache Commit failed: %v", err)
	}

	payload := bytes.Repeat([]byte("n"), cacheWriteSize+7)
	modTime := time.Now().Add(-time.Minute).Round(0)
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(payload)),
		modTime:      modTime,
		creationTime: modTime,
	})
	node.setDownloadURI("https://example.invalid/download")
	fs.backend = &fakeRemoteBackend{
		downloadFunc: fakeDownloadResponse(payload, int64(len(payload))),
	}

	errc, fh := fs.Open(path, fuse.O_RDONLY)
	if errc != 0 {
		t.Fatalf("Open returned unexpected error: %d", errc)
	}
	defer fs.Release(path, fh)

	buf := make([]byte, len(payload))
	if n := fs.Read(path, buf, 0, fh); n != len(payload) {
		t.Fatalf("Read returned %d, want %d", n, len(payload))
	}

	suffixBuf := make([]byte, len(suffixPayload))
	n, err := cacheStore.ReadComplete(suffixPath, cacheEntryMetadata(suffixPath, int64(len(suffixPayload)), suffixModTime), suffixBuf, 0)
	if err != nil {
		t.Fatalf("suffix ReadComplete failed: %v", err)
	}
	if n != len(suffixPayload) || string(suffixBuf[:n]) != string(suffixPayload) {
		t.Fatalf("suffix cache entry after partial cleanup = %q, want %q", string(suffixBuf[:n]), string(suffixPayload))
	}
}

func TestRemoteFsReadPinsPartialCacheDuringActiveDownload(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	cacheStore, err := mem.NewMemoryCache()
	if err != nil {
		t.Fatalf("NewMemoryCache failed: %v", err)
	}

	path := "/pinned-partial.bin"
	payload := bytes.Repeat([]byte("d"), cacheWriteSize+19)
	firstChunk := cacheWriteSize / 2
	observedCache := &notifyingCache{
		cacheStore:       cacheStore,
		writePartialPath: path,
		wrote:            make(chan struct{}),
	}
	fs.cacheStore = observedCache
	reader := newBlockingDownloadReader(payload, firstChunk)
	modTime := time.Now().Add(-time.Minute).Round(0)
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(payload)),
		modTime:      modTime,
		creationTime: modTime,
	})
	node.setDownloadURI("https://example.invalid/download")

	fs.backend = &fakeRemoteBackend{
		downloadFunc: func(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
			resp := &http.Response{
				StatusCode:    http.StatusOK,
				Body:          io.NopCloser(reader),
				ContentLength: int64(len(payload)),
			}
			if _, err := files_sdk.BuildResponse(resp, opts...); err != nil {
				return files_sdk.File{}, err
			}
			return files_sdk.File{Path: params.File.Path, Type: "file", Size: int64(len(payload))}, nil
		},
	}

	errc, fh := fs.Open(path, fuse.O_RDONLY)
	if errc != 0 {
		t.Fatalf("Open returned unexpected error: %d", errc)
	}
	defer fs.Release(path, fh)

	readDone := make(chan int, 1)
	go func() {
		buf := make([]byte, len(payload))
		readDone <- fs.Read(path, buf, 0, fh)
	}()

	<-reader.firstWritten
	select {
	case <-observedCache.wrote:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for partial cache write")
	}
	got := waitForPartialCacheBytes(t, cacheStore, path, firstChunk)
	if !bytes.Equal(got[:firstChunk], payload[:firstChunk]) {
		t.Fatalf("partial payload prefix = %q, want %q", string(got[:firstChunk]), string(payload[:firstChunk]))
	}

	_ = cacheStore.DeletePartial(path)
	got = waitForPartialCacheBytes(t, cacheStore, path, firstChunk)
	if !bytes.Equal(got[:firstChunk], payload[:firstChunk]) {
		t.Fatal("expected pinned active partial cache entry to survive Delete")
	}

	close(reader.release)
	select {
	case n := <-readDone:
		if n != len(payload) {
			t.Fatalf("Read returned %d, want %d", n, len(payload))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for read to finish")
	}

	partial := make([]byte, len(payload))
	if n, err := cacheStore.ReadPartial(path, partial, 0); err != nil || n != 0 {
		t.Fatalf("partial cache Read after read returned n=%d err=%v, want empty partial cache", n, err)
	}
}

func TestRemoteFsEnsureFullyCachedMissingNodeReturnsError(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/missing-hydration-node.bin"
	err := fs.ensureFullyCached(path, "https://example.invalid/download", 10, 0)
	if err == nil || !strings.Contains(err.Error(), "vfs node missing") {
		t.Fatalf("ensureFullyCached error = %v, want missing node error", err)
	}
	if _, ok := fs.peekGate(path); ok {
		t.Fatal("ensureFullyCached left a ready gate after missing node error")
	}
}

func TestRemoteFsCommitCacheEntryFromPartialCleansDestinationOnError(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	src := "/commit-source.bin"
	dst := "/commit-destination.bin"
	payload := []byte("short")
	if _, err := cacheStore.WritePartial(src, payload, 0); err != nil {
		t.Fatalf("source partial cache Write failed: %v", err)
	}

	err := fs.commitCacheEntryFromPartial(src, dst, cacheEntryMetadata(dst, int64(len(payload)+1), time.Now()), false)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("commitCacheEntryFromPartial error = %v, want %v", err, io.ErrUnexpectedEOF)
	}

	buf := make([]byte, len(payload))
	if n, err := cacheStore.Read(dst, buf, 0); err != nil || n != 0 {
		t.Fatalf("destination cache after failed commit returned n=%d err=%v, want empty", n, err)
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
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		uploadedPath = path
		uploaded = append([]byte(nil), data...)
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}
	moveCalls := 0
	fs.backend = &fakeRemoteBackend{
		moveFunc: func(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			moveCalls++
			return files_sdk.FileAction{}, nil
		},
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
	if moveCalls != 0 {
		t.Fatalf("backend move calls before upload = %d, want 0", moveCalls)
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

func TestRemoteFsRenameInProgressUploadUsesSessionDestination(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	temporaryPath := "/document.bin.~tmp"
	finalPath := "/document.bin"
	uploadStarted := make(chan string, 1)
	finishUpload := make(chan struct{})
	finalizedPath := make(chan string, 1)
	var finishOnce sync.Once
	releaseUpload := func() {
		finishOnce.Do(func() {
			close(finishUpload)
		})
	}
	defer releaseUpload()

	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		uploadStarted <- path
		<-finishUpload
		path, _ = fs.finalizeUploadPathAndRef(node)
		finalizedPath <- path
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	moveCalls := 0
	fs.backend = &fakeRemoteBackend{
		moveFunc: func(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			moveCalls++
			return files_sdk.FileAction{}, nil
		},
	}

	errno, fh := fs.Create(temporaryPath, fuse.O_RDWR, 0o644)
	if errno != 0 {
		t.Fatalf("Create returned unexpected error: %d", errno)
	}
	payload := []byte("pending-upload")
	if n := fs.Write(temporaryPath, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}

	flushResult := make(chan int, 1)
	go func() {
		flushResult <- fs.Flush(temporaryPath, fh)
	}()

	select {
	case path := <-uploadStarted:
		if path != temporaryPath {
			t.Fatalf("upload started at %q, want %q", path, temporaryPath)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for upload to start")
	}

	if errno := fs.Rename(temporaryPath, finalPath); errno != 0 {
		t.Fatalf("Rename returned unexpected error: %d", errno)
	}
	if moveCalls != 0 {
		t.Fatalf("backend move calls during upload = %d, want 0", moveCalls)
	}

	releaseUpload()
	select {
	case path := <-finalizedPath:
		if path != finalPath {
			t.Fatalf("upload finalized at %q, want %q", path, finalPath)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for upload final path")
	}
	select {
	case errno := <-flushResult:
		if errno != 0 {
			t.Fatalf("Flush returned unexpected error: %d", errno)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Flush")
	}
}

func TestRemoteFsRenameAfterFirstUploadDestinationCapturedMovesRemoteFile(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	temporaryPath := "/captured-destination.bin.~tmp"
	finalPath := "/captured-destination.bin"
	destinationCaptured := make(chan string, 1)
	finishUpload := make(chan struct{})
	var finishOnce sync.Once
	releaseUpload := func() {
		finishOnce.Do(func() {
			close(finishUpload)
		})
	}
	defer releaseUpload()

	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		path, _ = fs.finalizeUploadPathAndRef(node)
		destinationCaptured <- path
		<-finishUpload
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	moveCalled := make(chan files_sdk.FileMoveParams, 1)
	fs.backend = &fakeRemoteBackend{
		moveFunc: func(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			moveCalled <- params
			return files_sdk.FileAction{}, nil
		},
	}

	errno, fh := fs.Create(temporaryPath, fuse.O_RDWR, 0o644)
	if errno != 0 {
		t.Fatalf("Create returned unexpected error: %d", errno)
	}
	payload := []byte("first-upload")
	if n := fs.Write(temporaryPath, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}

	flushResult := make(chan int, 1)
	go func() {
		flushResult <- fs.Flush(temporaryPath, fh)
	}()
	select {
	case path := <-destinationCaptured:
		if path != temporaryPath {
			t.Fatalf("upload captured destination %q, want %q", path, temporaryPath)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for upload destination capture")
	}

	renameResult := make(chan int, 1)
	go func() {
		renameResult <- fs.Rename(temporaryPath, finalPath)
	}()
	select {
	case errno := <-renameResult:
		t.Fatalf("Rename returned %d before the first upload completed", errno)
	case <-time.After(100 * time.Millisecond):
	}
	select {
	case params := <-moveCalled:
		t.Fatalf("backend move started before the first upload completed: %+v", params)
	default:
	}

	releaseUpload()
	select {
	case params := <-moveCalled:
		if params.Path != temporaryPath || params.Destination != finalPath {
			t.Fatalf("backend move = %q to %q, want %q to %q", params.Path, params.Destination, temporaryPath, finalPath)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for backend move")
	}
	select {
	case errno := <-flushResult:
		if errno != 0 {
			t.Fatalf("Flush returned unexpected error: %d", errno)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Flush")
	}
	select {
	case errno := <-renameResult:
		if errno != 0 {
			t.Fatalf("Rename returned unexpected error: %d", errno)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Rename")
	}

	node, ok := vfs.fetch(finalPath)
	if !ok {
		t.Fatal("expected destination path after Rename")
	}
	session := node.getWriteSession()
	if session == nil {
		t.Fatal("expected retained session after Rename")
	}
	if path := session.snapshot().path; path != finalPath {
		t.Fatalf("retained session path = %q, want %q", path, finalPath)
	}
}

func TestRemoteFsRenameInProgressReuploadMovesCommittedRemoteFile(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	temporaryPath := "/reupload.bin.~tmp"
	finalPath := "/reupload.bin"
	secondUploadStarted := make(chan string, 1)
	finishSecondUpload := make(chan struct{})
	secondUploadFinalized := make(chan string, 1)
	var finishOnce sync.Once
	releaseSecondUpload := func() {
		finishOnce.Do(func() {
			close(finishSecondUpload)
		})
	}
	defer releaseSecondUpload()

	uploadCount := 0
	var uploadMu sync.Mutex
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		uploadMu.Lock()
		uploadCount++
		currentUpload := uploadCount
		uploadMu.Unlock()
		if currentUpload == 2 {
			secondUploadStarted <- path
			<-finishSecondUpload
			path, _ = fs.finalizeUploadPathAndRef(node)
			secondUploadFinalized <- path
		}
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	moveCalled := make(chan files_sdk.FileMoveParams, 1)
	fs.backend = &fakeRemoteBackend{
		moveFunc: func(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			moveCalled <- params
			return files_sdk.FileAction{}, nil
		},
	}

	errno, fh := fs.Create(temporaryPath, fuse.O_RDWR, 0o644)
	if errno != 0 {
		t.Fatalf("Create returned unexpected error: %d", errno)
	}
	initialPayload := []byte("initial")
	if n := fs.Write(temporaryPath, initialPayload, 0, fh); n != len(initialPayload) {
		t.Fatalf("initial Write returned %d, want %d", n, len(initialPayload))
	}
	if errno := fs.Flush(temporaryPath, fh); errno != 0 {
		t.Fatalf("initial Flush returned unexpected error: %d", errno)
	}

	updatedPayload := []byte("updated")
	if n := fs.Write(temporaryPath, updatedPayload, 0, fh); n != len(updatedPayload) {
		t.Fatalf("updated Write returned %d, want %d", n, len(updatedPayload))
	}
	flushResult := make(chan int, 1)
	go func() {
		flushResult <- fs.Flush(temporaryPath, fh)
	}()
	select {
	case path := <-secondUploadStarted:
		if path != temporaryPath {
			t.Fatalf("second upload started at %q, want %q", path, temporaryPath)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for second upload to start")
	}

	renameResult := make(chan int, 1)
	go func() {
		renameResult <- fs.Rename(temporaryPath, finalPath)
	}()
	select {
	case errno := <-renameResult:
		t.Fatalf("Rename returned %d before second upload completed", errno)
	case <-time.After(100 * time.Millisecond):
	}
	select {
	case params := <-moveCalled:
		t.Fatalf("backend move started before second upload completed: %+v", params)
	default:
	}

	releaseSecondUpload()
	select {
	case path := <-secondUploadFinalized:
		if path != temporaryPath {
			t.Fatalf("second upload finalized at %q, want %q", path, temporaryPath)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for second upload final path")
	}
	select {
	case errno := <-flushResult:
		if errno != 0 {
			t.Fatalf("second Flush returned unexpected error: %d", errno)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for second Flush")
	}
	select {
	case params := <-moveCalled:
		if params.Path != temporaryPath || params.Destination != finalPath {
			t.Fatalf("backend move = %q to %q, want %q to %q", params.Path, params.Destination, temporaryPath, finalPath)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for backend move")
	}
	select {
	case errno := <-renameResult:
		if errno != 0 {
			t.Fatalf("Rename returned unexpected error: %d", errno)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Rename")
	}

	node, ok := vfs.fetch(finalPath)
	if !ok {
		t.Fatal("expected destination path after Rename")
	}
	session := node.getWriteSession()
	if session == nil {
		t.Fatal("expected retained session after Rename")
	}
	if path := session.snapshot().path; path != finalPath {
		t.Fatalf("retained session path = %q, want %q", path, finalPath)
	}
}

func TestRemoteFsRenameCompletedUploadWithRetainedHandleMovesRemoteFile(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	temporaryPath := "/image.png.~tmp"
	finalPath := "/image.png"

	var uploadedPaths []string
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		uploadedPaths = append(uploadedPaths, path)
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	moveCalls := 0
	fs.backend = &fakeRemoteBackend{
		moveFunc: func(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			moveCalls++
			if params.Path != temporaryPath {
				t.Fatalf("move source = %q, want %q", params.Path, temporaryPath)
			}
			if params.Destination != finalPath {
				t.Fatalf("move destination = %q, want %q", params.Destination, finalPath)
			}
			return files_sdk.FileAction{}, nil
		},
	}

	errno, fh := fs.Create(temporaryPath, fuse.O_RDWR, 0o644)
	if errno != 0 {
		t.Fatalf("Create returned unexpected error: %d", errno)
	}

	payload := []byte("rotated-image")
	if n := fs.Write(temporaryPath, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}
	if errno := fs.Flush(temporaryPath, fh); errno != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errno)
	}

	if len(uploadedPaths) != 1 || uploadedPaths[0] != temporaryPath {
		t.Fatalf("uploaded paths = %q, want [%q]", uploadedPaths, temporaryPath)
	}
	node, ok := vfs.fetch(temporaryPath)
	if !ok {
		t.Fatal("expected temporary path to exist after upload")
	}
	session := node.getWriteSession()
	if session == nil {
		t.Fatal("expected completed write session to remain attached while handle is open")
	}
	snapshot := session.snapshot()
	if snapshot.dirty || snapshot.uploading || snapshot.finalizing {
		t.Fatalf(
			"write session state after upload = dirty:%t uploading:%t finalizing:%t, want clean and completed",
			snapshot.dirty,
			snapshot.uploading,
			snapshot.finalizing,
		)
	}
	if snapshot.handleCount != 1 {
		t.Fatalf("write session handle count = %d, want 1", snapshot.handleCount)
	}

	if errno := fs.Rename(temporaryPath, finalPath); errno != 0 {
		t.Fatalf("Rename returned unexpected error: %d", errno)
	}
	if moveCalls != 1 {
		t.Fatalf("backend move calls = %d, want 1", moveCalls)
	}

	updatedPayload := []byte("updated-image")
	if n := fs.Write(finalPath, updatedPayload, 0, fh); n != len(updatedPayload) {
		t.Fatalf("Write after rename returned %d, want %d", n, len(updatedPayload))
	}
	if errno := fs.Flush(finalPath, fh); errno != 0 {
		t.Fatalf("Flush after rename returned unexpected error: %d", errno)
	}
	if len(uploadedPaths) != 2 || uploadedPaths[1] != finalPath {
		t.Fatalf("uploaded paths after retained-handle write = %q, want second path %q", uploadedPaths, finalPath)
	}
}

func TestRemoteFsRenameSerializesStartingMutationAndNextUpload(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	temporaryPath := "/starting-mutation.bin.~tmp"
	finalPath := "/starting-mutation.bin"
	var uploadedPaths []string
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		uploadedPaths = append(uploadedPaths, path)
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	moveStarted := make(chan struct{})
	finishMove := make(chan struct{})
	var finishOnce sync.Once
	releaseMove := func() {
		finishOnce.Do(func() {
			close(finishMove)
		})
	}
	defer releaseMove()
	fs.backend = &fakeRemoteBackend{
		moveFunc: func(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			close(moveStarted)
			<-finishMove
			return files_sdk.FileAction{}, nil
		},
	}

	errno, fh := fs.Create(temporaryPath, fuse.O_RDWR, 0o644)
	if errno != 0 {
		t.Fatalf("Create returned unexpected error: %d", errno)
	}
	initialPayload := []byte("initial")
	if n := fs.Write(temporaryPath, initialPayload, 0, fh); n != len(initialPayload) {
		t.Fatalf("initial Write returned %d, want %d", n, len(initialPayload))
	}
	if errno := fs.Flush(temporaryPath, fh); errno != 0 {
		t.Fatalf("initial Flush returned unexpected error: %d", errno)
	}

	node, ok := vfs.fetch(temporaryPath)
	if !ok {
		t.Fatal("expected temporary path after upload")
	}
	session, created, err := node.beginWriteSessionMutation(temporaryPath)
	if err != nil {
		t.Fatalf("beginWriteSessionMutation returned unexpected error: %v", err)
	}
	if created {
		t.Fatal("expected to reuse retained write session")
	}
	mutationActive := true
	defer func() {
		if mutationActive {
			session.endMutation()
		}
	}()

	renameResult := make(chan int, 1)
	go func() {
		renameResult <- fs.Rename(temporaryPath, finalPath)
	}()
	select {
	case <-moveStarted:
		t.Fatal("backend move started before the existing mutation completed")
	case <-time.After(100 * time.Millisecond):
	}

	updatedPayload := []byte("updated")
	if _, err := fs.writeToWorkingCopy(session, updatedPayload, 0); err != nil {
		t.Fatalf("writeToWorkingCopy returned unexpected error: %v", err)
	}
	session.endMutation()
	mutationActive = false

	select {
	case <-moveStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for backend move after mutation completed")
	}

	flushResult := make(chan int, 1)
	go func() {
		flushResult <- fs.Flush(temporaryPath, fh)
	}()
	select {
	case errno := <-flushResult:
		t.Fatalf("Flush returned %d before backend move completed", errno)
	case <-time.After(100 * time.Millisecond):
	}

	releaseMove()
	select {
	case errno := <-renameResult:
		if errno != 0 {
			t.Fatalf("Rename returned unexpected error: %d", errno)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Rename")
	}
	select {
	case errno := <-flushResult:
		if errno != 0 {
			t.Fatalf("Flush returned unexpected error: %d", errno)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Flush")
	}

	if len(uploadedPaths) != 2 || uploadedPaths[0] != temporaryPath || uploadedPaths[1] != finalPath {
		t.Fatalf("uploaded paths = %q, want [%q %q]", uploadedPaths, temporaryPath, finalPath)
	}
}

func TestRemoteFsRenameWaiterRetriesAfterRetainedSessionCleared(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	temporaryPath := "/detached-session.bin.~tmp"
	finalPath := "/detached-session.bin"
	initialPayload := []byte("initial-data")
	updatedPayload := []byte("new")
	expectedPayload := append([]byte(nil), initialPayload...)
	copy(expectedPayload, updatedPayload)
	cacheDeleteStarted := make(chan struct{})
	releaseCacheDelete := make(chan struct{})
	var releaseCacheOnce sync.Once
	releaseCache := func() {
		releaseCacheOnce.Do(func() {
			close(releaseCacheDelete)
		})
	}
	defer releaseCache()
	fs.cacheStore = &blockingDeleteCacheStore{
		cacheStore: cacheStore,
		path:       finalPath,
		started:    cacheDeleteStarted,
		release:    releaseCacheDelete,
	}
	var uploadedPaths []string
	var uploadedPayloads [][]byte
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		uploadedPaths = append(uploadedPaths, path)
		uploadedPayloads = append(uploadedPayloads, append([]byte(nil), data...))
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	moveStarted := make(chan struct{})
	finishMove := make(chan struct{})
	var finishOnce sync.Once
	releaseMove := func() {
		finishOnce.Do(func() {
			close(finishMove)
		})
	}
	defer releaseMove()
	var downloadedPaths []string
	fs.backend = &fakeRemoteBackend{
		moveFunc: func(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			close(moveStarted)
			<-finishMove
			return files_sdk.FileAction{}, nil
		},
		downloadFunc: func(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
			downloadedPaths = append(downloadedPaths, params.File.Path)
			return fakeDownloadResponse(initialPayload, int64(len(initialPayload)))(params, opts...)
		},
	}

	errno, waitingFH := fs.Create(temporaryPath, fuse.O_RDWR, 0o644)
	if errno != 0 {
		t.Fatalf("Create returned unexpected error: %d", errno)
	}
	errno, writerFH := fs.Open(temporaryPath, fuse.O_RDWR)
	if errno != 0 {
		t.Fatalf("Open returned unexpected error: %d", errno)
	}
	if n := fs.Write(temporaryPath, initialPayload, 0, writerFH); n != len(initialPayload) {
		t.Fatalf("initial Write returned %d, want %d", n, len(initialPayload))
	}
	if errno := fs.Flush(temporaryPath, writerFH); errno != 0 {
		t.Fatalf("initial Flush returned unexpected error: %d", errno)
	}

	node, ok := vfs.fetch(temporaryPath)
	if !ok {
		t.Fatal("expected temporary path after initial upload")
	}
	retainedSession := node.getWriteSession()
	if retainedSession == nil {
		t.Fatal("expected retained write session")
	}
	if handles := retainedSession.snapshot().handleCount; handles != 1 {
		t.Fatalf("retained session handle count = %d, want 1", handles)
	}
	stalePayload := bytes.Repeat([]byte("x"), len(initialPayload))
	if _, err := cacheStore.Write(finalPath, stalePayload, 0); err != nil {
		t.Fatalf("stale destination cache Write failed: %v", err)
	}
	if err := cacheStore.Commit(finalPath, cacheEntryMetadata(finalPath, int64(len(stalePayload)), node.info.modTime)); err != nil {
		t.Fatalf("stale destination cache Commit failed: %v", err)
	}
	cacheStore.Pin(finalPath)
	defer cacheStore.Unpin(finalPath)

	renameResult := make(chan int, 1)
	go func() {
		renameResult <- fs.Rename(temporaryPath, finalPath)
	}()
	select {
	case <-moveStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for backend move")
	}

	writeResult := make(chan int, 1)
	go func() {
		writeResult <- fs.Write(temporaryPath, updatedPayload, 0, waitingFH)
	}()
	select {
	case n := <-writeResult:
		t.Fatalf("waiting Write returned %d before backend move completed", n)
	case <-time.After(100 * time.Millisecond):
	}

	if errno := fs.Release(temporaryPath, writerFH); errno != 0 {
		t.Fatalf("Release returned unexpected error: %d", errno)
	}
	releaseMove()
	select {
	case <-cacheDeleteStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for destination cache invalidation")
	}
	select {
	case n := <-writeResult:
		t.Fatalf("waiting Write returned %d before destination cache invalidation completed", n)
	case <-time.After(100 * time.Millisecond):
	}
	releaseCache()

	select {
	case errno := <-renameResult:
		if errno != 0 {
			t.Fatalf("Rename returned unexpected error: %d", errno)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Rename")
	}
	select {
	case n := <-writeResult:
		if n != len(updatedPayload) {
			t.Fatalf("waiting Write returned %d, want %d", n, len(updatedPayload))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Write")
	}

	node, ok = vfs.fetch(finalPath)
	if !ok {
		t.Fatal("expected destination path after Rename")
	}
	newSession := node.getWriteSession()
	if newSession == nil {
		t.Fatal("expected waiting Write to attach a replacement session")
	}
	if newSession == retainedSession {
		t.Fatal("waiting Write reused detached retained session")
	}
	if path := newSession.snapshot().path; path != finalPath {
		t.Fatalf("replacement session path = %q, want %q", path, finalPath)
	}
	if !newSession.hasHandle(waitingFH) {
		t.Fatal("replacement session does not contain waiting handle")
	}
	if len(downloadedPaths) != 1 || downloadedPaths[0] != finalPath {
		t.Fatalf("downloaded paths = %q, want [%q]", downloadedPaths, finalPath)
	}

	if errno := fs.Flush(finalPath, waitingFH); errno != 0 {
		t.Fatalf("replacement session Flush returned unexpected error: %d", errno)
	}
	if len(uploadedPaths) != 2 || uploadedPaths[0] != temporaryPath || uploadedPaths[1] != finalPath {
		t.Fatalf("uploaded paths = %q, want [%q %q]", uploadedPaths, temporaryPath, finalPath)
	}
	if len(uploadedPayloads) != 2 {
		t.Fatalf("uploaded payload count = %d, want 2", len(uploadedPayloads))
	}
	if !bytes.Equal(uploadedPayloads[1], expectedPayload) {
		t.Fatalf("second uploaded payload = %q, want %q", uploadedPayloads[1], expectedPayload)
	}
	if errno := fs.Release(finalPath, waitingFH); errno != 0 {
		t.Fatalf("waiting handle Release returned unexpected error: %d", errno)
	}
}

func TestRemoteFsRenameCancelsInFlightDestinationDownload(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	temporaryPath := "/in-flight-cache.bin.~tmp"
	finalPath := "/in-flight-cache.bin"
	newPayload := []byte("new-destination")
	stalePayload := []byte("old-destination")
	modTime := time.Now().Add(-time.Minute).Round(0)
	sourceNode := vfs.getOrCreate(temporaryPath, nodeTypeFile)
	sourceNode.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(newPayload)),
		modTime:      modTime,
		creationTime: modTime,
	})
	destinationNode := vfs.getOrCreate(finalPath, nodeTypeFile)
	destinationNode.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(stalePayload)),
		modTime:      modTime,
		creationTime: modTime,
	})
	destinationNode.setDownloadURI("https://example.invalid/old-destination")

	downloadStarted := make(chan struct{})
	releaseDownload := make(chan struct{})
	var releaseOnce sync.Once
	release := func() {
		releaseOnce.Do(func() { close(releaseDownload) })
	}
	defer release()
	fs.backend = &fakeRemoteBackend{
		downloadFunc: func(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
			reader := &cancelableBlockingDownloadReader{
				payload: stalePayload,
				ctx:     files_sdk.ContextOption(opts),
				started: downloadStarted,
				release: releaseDownload,
			}
			resp := &http.Response{
				StatusCode:    http.StatusOK,
				Body:          io.NopCloser(reader),
				ContentLength: int64(len(stalePayload)),
			}
			if _, err := files_sdk.BuildResponse(resp, opts...); err != nil {
				return files_sdk.File{}, err
			}
			return files_sdk.File{Path: params.File.Path, Type: "file", Size: int64(len(stalePayload))}, nil
		},
	}

	downloadResult := make(chan error, 1)
	go func() {
		downloadResult <- fs.ensureFullyCached(finalPath, destinationNode.downloadUri, int64(len(stalePayload)), 0)
	}()
	select {
	case <-downloadStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for destination download")
	}

	renameResult := make(chan int, 1)
	go func() {
		renameResult <- fs.Rename(temporaryPath, finalPath)
	}()
	select {
	case errno := <-renameResult:
		if errno != 0 {
			t.Fatalf("Rename returned unexpected error: %d", errno)
		}
	case <-time.After(2 * time.Second):
		release()
		t.Fatal("timed out waiting for Rename to cancel the destination download")
	}
	release()
	select {
	case err := <-downloadResult:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("destination download error = %v, want context.Canceled", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for destination download cancellation")
	}

	buf := make([]byte, len(stalePayload))
	if n, err := cacheStore.Read(finalPath, buf, 0); err != nil || n != 0 {
		t.Fatalf("destination cache after Rename = %d bytes, %v; want empty", n, err)
	}
}

func TestRemoteFsRenameCompletedUploadFailureKeepsOriginalPaths(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	temporaryPath := "/report.dat.~tmp"
	finalPath := "/report.dat"
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	moveCalls := 0
	fs.backend = &fakeRemoteBackend{
		moveFunc: func(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			moveCalls++
			return files_sdk.FileAction{}, errors.New("move failed")
		},
	}

	errno, fh := fs.Create(temporaryPath, fuse.O_RDWR, 0o644)
	if errno != 0 {
		t.Fatalf("Create returned unexpected error: %d", errno)
	}
	payload := []byte("completed-upload")
	if n := fs.Write(temporaryPath, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}
	if errno := fs.Flush(temporaryPath, fh); errno != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errno)
	}

	if errno := fs.Rename(temporaryPath, finalPath); errno == 0 {
		t.Fatal("expected Rename to fail when backend move fails")
	}
	if moveCalls != 1 {
		t.Fatalf("backend move calls = %d, want 1", moveCalls)
	}
	node, ok := vfs.fetch(temporaryPath)
	if !ok {
		t.Fatal("expected original VFS path to remain after failed move")
	}
	if _, ok := vfs.fetch(finalPath); ok {
		t.Fatal("expected destination VFS path to remain absent after failed move")
	}
	session := node.getWriteSession()
	if session == nil {
		t.Fatal("expected retained write session after failed move")
	}
	if path := session.snapshot().path; path != temporaryPath {
		t.Fatalf("write session path after failed move = %q, want %q", path, temporaryPath)
	}
}

func TestRemoteFsRenameCompletedUploadSerializesConcurrentWrite(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	temporaryPath := "/serialized.bin.~tmp"
	finalPath := "/serialized.bin"
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	moveStarted := make(chan struct{})
	finishMove := make(chan struct{})
	var finishOnce sync.Once
	releaseMove := func() {
		finishOnce.Do(func() {
			close(finishMove)
		})
	}
	defer releaseMove()
	fs.backend = &fakeRemoteBackend{
		moveFunc: func(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			close(moveStarted)
			<-finishMove
			return files_sdk.FileAction{}, nil
		},
	}

	errno, fh := fs.Create(temporaryPath, fuse.O_RDWR, 0o644)
	if errno != 0 {
		t.Fatalf("Create returned unexpected error: %d", errno)
	}
	initialPayload := []byte("initial")
	if n := fs.Write(temporaryPath, initialPayload, 0, fh); n != len(initialPayload) {
		t.Fatalf("Write returned %d, want %d", n, len(initialPayload))
	}
	if errno := fs.Flush(temporaryPath, fh); errno != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errno)
	}

	renameResult := make(chan int, 1)
	go func() {
		renameResult <- fs.Rename(temporaryPath, finalPath)
	}()
	select {
	case <-moveStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for backend move")
	}

	node, ok := vfs.fetch(temporaryPath)
	if !ok {
		t.Fatal("expected source VFS path while backend move is pending")
	}
	hasSessionResult := make(chan bool, 1)
	go func() {
		hasSessionResult <- node.hasActiveWriteSession()
	}()
	select {
	case hasSession := <-hasSessionResult:
		if !hasSession {
			t.Fatal("expected retained write session while backend move is pending")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("write-session inspection blocked while backend move was pending")
	}

	readResult := make(chan int, 1)
	readBuffer := make([]byte, len(initialPayload))
	go func() {
		readResult <- fs.Read(temporaryPath, readBuffer, 0, fh)
	}()
	select {
	case n := <-readResult:
		if n != len(initialPayload) || !bytes.Equal(readBuffer, initialPayload) {
			t.Fatalf("Read during backend move returned %d bytes %q, want %d bytes %q", n, readBuffer, len(initialPayload), initialPayload)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Read blocked while backend move was pending")
	}

	updatedPayload := []byte("updated")
	writeResult := make(chan int, 1)
	go func() {
		writeResult <- fs.Write(finalPath, updatedPayload, 0, fh)
	}()
	select {
	case n := <-writeResult:
		t.Fatalf("concurrent Write returned %d before backend move completed", n)
	case <-time.After(100 * time.Millisecond):
	}

	releaseMove()
	select {
	case errno := <-renameResult:
		if errno != 0 {
			t.Fatalf("Rename returned unexpected error: %d", errno)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Rename")
	}
	select {
	case n := <-writeResult:
		if n != len(updatedPayload) {
			t.Fatalf("Write after backend move returned %d, want %d", n, len(updatedPayload))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for concurrent Write")
	}

	node, ok = vfs.fetch(finalPath)
	if !ok {
		t.Fatal("expected destination VFS path after successful move")
	}
	session := node.getWriteSession()
	if session == nil {
		t.Fatal("expected retained write session after successful move")
	}
	if path := session.snapshot().path; path != finalPath {
		t.Fatalf("write session path after serialized move = %q, want %q", path, finalPath)
	}
}

func TestRemoteFsRenameCompletedUploadDoesNotBlockRelease(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	temporaryPath := "/release.bin.~tmp"
	finalPath := "/release.bin"
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	moveStarted := make(chan struct{})
	finishMove := make(chan struct{})
	var finishOnce sync.Once
	releaseMove := func() {
		finishOnce.Do(func() {
			close(finishMove)
		})
	}
	defer releaseMove()
	fs.backend = &fakeRemoteBackend{
		moveFunc: func(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
			close(moveStarted)
			<-finishMove
			return files_sdk.FileAction{}, nil
		},
	}

	errno, fh := fs.Create(temporaryPath, fuse.O_RDWR, 0o644)
	if errno != 0 {
		t.Fatalf("Create returned unexpected error: %d", errno)
	}
	payload := []byte("release-during-rename")
	if n := fs.Write(temporaryPath, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}
	if errno := fs.Flush(temporaryPath, fh); errno != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errno)
	}

	renameResult := make(chan int, 1)
	go func() {
		renameResult <- fs.Rename(temporaryPath, finalPath)
	}()
	select {
	case <-moveStarted:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for backend move")
	}

	releaseResult := make(chan int, 1)
	go func() {
		releaseResult <- fs.Release(temporaryPath, fh)
	}()
	select {
	case errno := <-releaseResult:
		if errno != 0 {
			t.Fatalf("Release returned unexpected error: %d", errno)
		}
	case <-time.After(2 * time.Second):
		releaseMove()
		t.Fatal("Release blocked while backend move was pending")
	}

	releaseMove()
	select {
	case errno := <-renameResult:
		if errno != 0 {
			t.Fatalf("Rename returned unexpected error: %d", errno)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Rename")
	}

	node, ok := vfs.fetch(finalPath)
	if !ok {
		t.Fatal("expected destination VFS path after successful move")
	}
	if session := node.getWriteSession(); session != nil {
		t.Fatal("expected retained write session to be cleared after pending Release")
	}
}

func TestRemoteFsFinalizeUploadPathAndRefUsesMountedRootAfterActiveRename(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	fs.root = "/o_test"

	oldPath := "/draft.ai"
	newPath := "/final.ai"

	errc, fh := fs.Create(oldPath, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	payload := []byte("rename-then-finalize")
	if n := fs.Write(oldPath, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}

	if errc := fs.Rename(oldPath, newPath); errc != 0 {
		t.Fatalf("Rename returned unexpected error: %d", errc)
	}

	node, ok := vfs.fetch(newPath)
	if !ok {
		t.Fatal("expected renamed path to exist in VFS")
	}

	finalPath, _ := fs.finalizeUploadPathAndRef(node)
	if finalPath != "/o_test/final.ai" {
		t.Fatalf("finalize upload path = %q, want %q", finalPath, "/o_test/final.ai")
	}
}

func TestRemoteFsRenameUploadFailureLogsSanitizedStorageError(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	logger := &captureMountLogger{}
	fs.log = logger

	rawMessage := "Your socket connection to the server was not read from or written to within the timeout period.\nIdle connections will be closed."
	fs.backend = &fakeRemoteBackend{
		uploadFunc: func(opts ...file.UploadOption) error {
			return lib.S3Error{Code: "RequestTimeout", Message: rawMessage}
		},
	}

	src := filepath.Join(t.TempDir(), "rename-timeout.pdf")
	if err := os.WriteFile(src, []byte("rename-timeout"), 0o600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	err := fs.uploadFile(src, "/rename-timeout.pdf")
	if err == nil {
		t.Fatal("expected rename upload to fail")
	}

	logs := logger.visibleJoined()
	if count := strings.Count(logs, "Error uploading file during rename:"); count != 1 {
		t.Fatalf("rename failure log count = %d, want 1. logs:\n%s", count, logs)
	}
	if !strings.Contains(logs, "Upload to backend storage timed out") {
		t.Fatalf("expected safe timeout message, got:\n%s", logs)
	}
	for _, raw := range []string{"RequestTimeout", "Your socket connection", "Idle connections"} {
		if strings.Contains(logs, raw) {
			t.Fatalf("expected user-visible logs to omit raw provider text %q, got:\n%s", raw, logs)
		}
	}
}

func TestRemoteFsRenameUploadFailureDoesNotCommitDiskCache(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	cacheStore := newTestDiskCache(t)
	fs.cacheStore = cacheStore

	uploadErr := errors.New("upload failed")
	fs.backend = &fakeRemoteBackend{
		uploadFunc: func(opts ...file.UploadOption) error {
			return uploadErr
		},
	}

	src := filepath.Join(t.TempDir(), "rename-failure.pdf")
	payload := bytes.Repeat([]byte("u"), cacheWriteSize+3)
	if err := os.WriteFile(src, payload, 0o600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	dst := "/rename-failure.pdf"
	if err := fs.uploadFile(src, dst); !errors.Is(err, uploadErr) {
		t.Fatalf("uploadFile error = %v, want %v", err, uploadErr)
	}

	buf := make([]byte, len(payload))
	if n, err := cacheStore.Read(dst, buf, 0); err != nil || n != 0 {
		t.Fatalf("cache Read after failed upload returned n=%d err=%v, want empty cache", n, err)
	}
}

func TestRemoteFsUploadFileUsesUploadResponseMetadataForPostUploadReads(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	src := filepath.Join(t.TempDir(), "rename-cache.bin")
	dst := "/rename-cache.bin"
	payload := []byte("rename-upload-cache")
	if err := os.WriteFile(src, payload, 0o600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	localMtime := time.Date(2026, 5, 21, 11, 30, 45, 0, time.UTC)
	if err := os.Chtimes(src, localMtime, localMtime); err != nil {
		t.Fatalf("Chtimes failed: %v", err)
	}

	remoteMtime := localMtime.Add(30 * time.Second)
	remoteFile := files_sdk.File{
		DisplayName: "rename-cache.bin",
		Path:        dst,
		Type:        "file",
		Size:        int64(len(payload)),
		Mtime:       &remoteMtime,
		DownloadUri: "https://example.invalid/download",
	}

	backend := fs.backend.(*fakeRemoteBackend)
	backend.uploadWithResumeFunc = func(opts ...file.UploadOption) (file.UploadResumable, error) {
		return file.UploadResumable{File: remoteFile}, nil
	}
	downloadCalls := 0
	backend.downloadFunc = func(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
		downloadCalls++
		return fakeDownloadResponse(payload, int64(len(payload)))(params, opts...)
	}

	if err := fs.uploadFile(src, dst); err != nil {
		t.Fatalf("uploadFile returned error: %v", err)
	}

	node, ok := vfs.fetch(dst)
	if !ok {
		t.Fatal("expected node to exist after uploadFile")
	}
	if !node.info.modTime.Equal(remoteMtime) {
		t.Fatalf("node mtime = %v, want upload response mtime %v", node.info.modTime, remoteMtime)
	}

	buf := make([]byte, len(payload))
	n, err := cacheStore.ReadComplete(dst, cacheEntryMetadata(dst, int64(len(payload)), remoteMtime), buf, 0)
	if err != nil {
		t.Fatalf("cache ReadComplete failed: %v", err)
	}
	if n != len(payload) || !bytes.Equal(buf[:n], payload) {
		t.Fatalf("cache ReadComplete returned n=%d payload=%q, want %q", n, string(buf[:n]), string(payload))
	}

	fs.createNode(dst, remoteFile)

	errc, fh := fs.Open(dst, fuse.O_RDONLY)
	if errc != 0 {
		t.Fatalf("Open returned unexpected error: %d", errc)
	}
	defer fs.Release(dst, fh)

	readBuf := make([]byte, len(payload))
	if n := fs.Read(dst, readBuf, 0, fh); n != len(payload) {
		t.Fatalf("Read returned %d, want %d", n, len(payload))
	}
	if !bytes.Equal(readBuf, payload) {
		t.Fatalf("Read returned %q, want %q", string(readBuf), string(payload))
	}
	if downloadCalls != 0 {
		t.Fatalf("download calls = %d, want 0", downloadCalls)
	}
}

func TestRemoteFsPublicTruncateZeroSkipsHydrationAndResetsWorkingCopy(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/existing.ai"
	original := []byte("old remote contents")
	modTime := time.Now().Round(0)
	if _, err := cacheStore.Write(path, original, 0); err != nil {
		t.Fatalf("cache Write failed: %v", err)
	}
	if err := cacheStore.Commit(path, cacheEntryMetadata(path, int64(len(original)), modTime)); err != nil {
		t.Fatalf("cache Commit failed: %v", err)
	}

	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(original)),
		modTime:      modTime,
		creationTime: modTime,
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
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		uploaded = append([]byte(nil), data...)
		return testUploadedMetadata(int64(len(data)), mtime), nil
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

	node, ok := vfs.fetch(path)
	if !ok {
		t.Fatal("expected node to exist after Flush")
	}

	buf := make([]byte, len(payload))
	n, err := cacheStore.ReadComplete(path, cacheEntryMetadata(path, int64(len(payload)), node.info.modTime), buf, 0)
	if err != nil {
		t.Fatalf("cache ReadComplete failed: %v", err)
	}
	if n != len(payload) {
		t.Fatalf("cache ReadComplete returned %d, want %d", n, len(payload))
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

func TestRemoteFsFlushUsesUploadResponseMetadataForPostUploadReads(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/post-upload.bin"
	payload := []byte("post-upload-cache")
	remoteMtime := time.Date(2026, 5, 21, 12, 30, 45, 0, time.UTC)
	remoteFile := files_sdk.File{
		DisplayName: "post-upload.bin",
		Path:        path,
		Type:        "file",
		Size:        int64(len(payload)),
		Mtime:       &remoteMtime,
		DownloadUri: "https://example.invalid/download",
	}

	backend := fs.backend.(*fakeRemoteBackend)
	backend.uploadWithResumeFunc = func(opts ...file.UploadOption) (file.UploadResumable, error) {
		return file.UploadResumable{File: remoteFile}, nil
	}
	downloadCalls := 0
	backend.downloadFunc = func(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
		downloadCalls++
		return fakeDownloadResponse(payload, int64(len(payload)))(params, opts...)
	}

	errc, fh := fs.Create(path, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	if n := fs.Write(path, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}

	if errc := fs.Flush(path, fh); errc != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errc)
	}

	node, ok := vfs.fetch(path)
	if !ok {
		t.Fatal("expected node to exist after Flush")
	}
	if !node.info.modTime.Equal(remoteMtime) {
		t.Fatalf("node mtime = %v, want upload response mtime %v", node.info.modTime, remoteMtime)
	}

	buf := make([]byte, len(payload))
	n, err := cacheStore.ReadComplete(path, cacheEntryMetadata(path, int64(len(payload)), remoteMtime), buf, 0)
	if err != nil {
		t.Fatalf("cache ReadComplete failed: %v", err)
	}
	if n != len(payload) || !bytes.Equal(buf[:n], payload) {
		t.Fatalf("cache ReadComplete returned n=%d payload=%q, want %q", n, string(buf[:n]), string(payload))
	}

	if errc := fs.Release(path, fh); errc != 0 {
		t.Fatalf("Release returned unexpected error: %d", errc)
	}

	fs.createNode(path, remoteFile)

	errc, readFh := fs.Open(path, fuse.O_RDONLY)
	if errc != 0 {
		t.Fatalf("Open returned unexpected error: %d", errc)
	}
	defer fs.Release(path, readFh)

	readBuf := make([]byte, len(payload))
	if n := fs.Read(path, readBuf, 0, readFh); n != len(payload) {
		t.Fatalf("Read returned %d, want %d", n, len(payload))
	}
	if !bytes.Equal(readBuf, payload) {
		t.Fatalf("Read returned %q, want %q", string(readBuf), string(payload))
	}
	if downloadCalls != 0 {
		t.Fatalf("download calls = %d, want 0", downloadCalls)
	}
}

func TestRemoteFsWorkingCopyUploadSuccessLogsLifecycle(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	logger := &captureMountLogger{}
	fs.log = logger

	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	path := "/lifecycle-success.ai"
	errc, fh := fs.Create(path, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	payload := []byte("lifecycle-success")
	if n := fs.Write(path, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}
	if errc := fs.Flush(path, fh); errc != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errc)
	}

	logs := logger.visibleJoined()
	if count := strings.Count(logs, "Starting upload from working copy:"); count != 1 {
		t.Fatalf("start log count = %d, want 1. logs:\n%s", count, logs)
	}
	if count := strings.Count(logs, "Upload completed from working copy:"); count != 1 {
		t.Fatalf("completion log count = %d, want 1. logs:\n%s", count, logs)
	}
}

func TestRemoteFsPublicFlushPoisonsSessionAfterUploadFailure(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	uploadErr := errors.New("upload failed")
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		_, _ = io.ReadAll(reader)
		return uploadedFileMetadata{}, uploadErr
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

func TestRemoteFsWorkingCopyUploadFailureLogsSanitizedStorageError(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	logger := &captureMountLogger{}
	fs.log = logger

	rawMessage := "Your socket connection to the server was not read from or written to within the timeout period.\nIdle connections will be closed."
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		_, _ = io.ReadAll(reader)
		return uploadedFileMetadata{}, lib.S3Error{Code: "RequestTimeout", Message: rawMessage}
	}

	path := "/timeout.pdf"
	errc, fh := fs.Create(path, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	payload := []byte("upload-timeout")
	if n := fs.Write(path, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}
	if errc := fs.Flush(path, fh); errc == 0 {
		t.Fatal("expected Flush to fail after upload timeout")
	}
	if errc := fs.Truncate(path, int64(len(payload)+1), fh); errc == 0 {
		t.Fatal("expected Truncate to fail after failed upload")
	}

	logs := logger.visibleJoined()
	if count := strings.Count(logs, "Starting upload from working copy:"); count != 1 {
		t.Fatalf("start log count = %d, want 1. logs:\n%s", count, logs)
	}
	if count := strings.Count(logs, "Upload failed from working copy:"); count != 1 {
		t.Fatalf("failure log count = %d, want 1. logs:\n%s", count, logs)
	}
	if !strings.Contains(logs, "Upload to backend storage timed out") {
		t.Fatalf("expected safe Files.com-facing timeout message, got:\n%s", logs)
	}
	if count := strings.Count(logs, "Upload to backend storage timed out"); count != 2 {
		t.Fatalf("safe timeout message count = %d, want 2 after initial failure and follow-up operation. logs:\n%s", count, logs)
	}
	if strings.Contains(logs, "error_class=") {
		t.Fatalf("expected user-visible logs to omit internal error class, got:\n%s", logs)
	}
	for _, raw := range []string{"RequestTimeout", "Your socket connection", "Idle connections"} {
		if strings.Contains(logs, raw) {
			t.Fatalf("expected user-visible logs to omit raw provider text %q, got:\n%s", raw, logs)
		}
	}
	for _, line := range logger.visibleLines {
		if strings.Contains(line, "\n") || strings.Contains(line, "\r") {
			t.Fatalf("expected each log entry to stay on one line, got %q", line)
		}
	}
}

func TestRemoteFsWorkingCopyUploadFailureLogsSafeFallbackForNonS3Error(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	logger := &captureMountLogger{}
	fs.log = logger

	rawMessage := "dial tcp 10.0.0.1:443: connect: raw provider network failure"
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		_, _ = io.ReadAll(reader)
		return uploadedFileMetadata{}, errors.New(rawMessage)
	}

	path := "/network-error.pdf"
	errc, fh := fs.Create(path, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	payload := []byte("upload-network-error")
	if n := fs.Write(path, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}
	if errc := fs.Flush(path, fh); errc == 0 {
		t.Fatal("expected Flush to fail after upload error")
	}
	if errc := fs.Truncate(path, int64(len(payload)+1), fh); errc == 0 {
		t.Fatal("expected Truncate to fail after failed upload")
	}

	logs := logger.visibleJoined()
	if count := strings.Count(logs, "Upload failed from working copy:"); count != 1 {
		t.Fatalf("failure log count = %d, want 1. logs:\n%s", count, logs)
	}
	if !strings.Contains(logs, "Error returned by the remote service") {
		t.Fatalf("expected safe fallback message, got:\n%s", logs)
	}
	if count := strings.Count(logs, "Error returned by the remote service"); count != 2 {
		t.Fatalf("safe fallback message count = %d, want 2 after initial failure and follow-up operation. logs:\n%s", count, logs)
	}
	if strings.Contains(logs, rawMessage) {
		t.Fatalf("expected user-visible logs to omit raw error %q, got:\n%s", rawMessage, logs)
	}
}

func TestUploadLogMessageTransportFallbacks(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		message string
	}{
		{
			name:    "net timeout",
			err:     &net.OpError{Op: "dial", Net: "tcp", Err: timeoutUploadError{}},
			message: "Network timeout",
		},
		{
			name:    "os timeout",
			err:     os.ErrDeadlineExceeded,
			message: "Network timeout",
		},
		{
			name:    "connection refused",
			err:     fmt.Errorf("wrapped: %w", syscall.ECONNREFUSED),
			message: "Connection refused",
		},
		{
			name:    "connection reset",
			err:     fmt.Errorf("wrapped: %w", syscall.ECONNRESET),
			message: "Connection reset by peer",
		},
		{
			name:    "broken pipe",
			err:     fmt.Errorf("wrapped: %w", syscall.EPIPE),
			message: "Client disconnected during transfer",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := uploadLogMessage(test.err); got != test.message {
				t.Fatalf("uploadLogMessage() = %q, want %q", got, test.message)
			}
		})
	}
}

func TestFormatFuseErrno(t *testing.T) {
	tests := []struct {
		name string
		errc int
		want string
	}{
		{name: "success", errc: 0, want: "OK"},
		{name: "negative errno", errc: -fuse.EACCES, want: "EACCES"},
		{name: "positive errno", errc: fuse.ENOENT, want: "ENOENT"},
		{name: "unknown errno", errc: -9999, want: "errno_9999"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := formatFuseErrno(test.errc); got != test.want {
				t.Fatalf("formatFuseErrno(%d) = %q, want %q", test.errc, got, test.want)
			}
		})
	}
}

func TestClassifyMountError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantClass string
		wantErrc  int
	}{
		{name: "nil", err: nil, wantClass: "none", wantErrc: 0},
		{name: "not authenticated", err: files_sdk.ResponseError{Type: string(files_sdk.ErrInvalidCredentials)}, wantClass: "not_authenticated", wantErrc: -fuse.EPERM},
		{name: "not exist", err: files_sdk.ResponseError{Type: string(files_sdk.ErrFileNotFound)}, wantClass: "not_exist", wantErrc: -fuse.ENOENT},
		{name: "exists", err: files_sdk.ResponseError{Type: string(files_sdk.ErrDestinationExists)}, wantClass: "exist", wantErrc: -fuse.EEXIST},
		{name: "no slots available", err: lim.ErrNoSlotsAvailable, wantClass: "no_slots_available", wantErrc: -fuse.EAGAIN},
		{name: "folder not empty", err: files_sdk.ResponseError{Type: string(files_sdk.ErrFolderNotEmpty)}, wantClass: "folder_not_empty", wantErrc: -fuse.ENOTEMPTY},
		{name: "resource locked message", err: errors.New("resource locked by another operation"), wantClass: "resource_locked", wantErrc: -fuse.EAGAIN},
		{name: "resource locked type without message", err: files_sdk.ResponseError{Type: string(files_sdk.ErrResourceLocked)}, wantClass: "unknown", wantErrc: -fuse.EIO},
		{name: "unknown", err: errors.New("backend exploded"), wantClass: "unknown", wantErrc: -fuse.EIO},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotClass, gotErrc := classifyMountError(test.err)
			if gotClass != test.wantClass {
				t.Fatalf("classifyMountError() class = %q, want %q", gotClass, test.wantClass)
			}
			if gotErrc != test.wantErrc {
				t.Fatalf("classifyMountError() errc = %d (%s), want %d (%s)", gotErrc, formatFuseErrno(gotErrc), test.wantErrc, formatFuseErrno(test.wantErrc))
			}
		})
	}
}

func TestRemoteFsWorkingCopyUploadCancellationLogsCanceled(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	logger := &captureMountLogger{}
	fs.log = logger

	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		_, _ = io.ReadAll(reader)
		return uploadedFileMetadata{}, context.Canceled
	}

	path := "/cancel.pdf"
	errc, fh := fs.Create(path, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	payload := []byte("upload-canceled")
	if n := fs.Write(path, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}
	_ = fs.Flush(path, fh)

	logs := logger.visibleJoined()
	if count := strings.Count(logs, "Starting upload from working copy:"); count != 1 {
		t.Fatalf("start log count = %d, want 1. logs:\n%s", count, logs)
	}
	if count := strings.Count(logs, "Upload canceled from working copy:"); count != 1 {
		t.Fatalf("cancellation log count = %d, want 1. logs:\n%s", count, logs)
	}
	if strings.Contains(logs, "Upload failed from working copy:") {
		t.Fatalf("expected cancellation not to log failure, got:\n%s", logs)
	}
}

func TestRemoteFsInPlaceWritesAndFlushDoNotChangeSizeUntilTruncate(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/dense-save.indd"
	initialSize := 2166784
	finalSize := 2498560

	initial := bytes.Repeat([]byte("a"), initialSize)
	initialMtime := time.Now().Round(0)
	if _, err := cacheStore.Write(path, initial, 0); err != nil {
		t.Fatalf("cache Write failed: %v", err)
	}
	if err := cacheStore.Commit(path, cacheEntryMetadata(path, int64(initialSize), initialMtime)); err != nil {
		t.Fatalf("cache Commit failed: %v", err)
	}

	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(initialSize),
		modTime:      initialMtime,
		creationTime: initialMtime,
	})
	node.setDownloadURI("https://example.invalid/download")

	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		return testUploadedMetadata(int64(len(data)), mtime), nil
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

func TestRemoteFsSparseInPlaceWriteHydratesWhenCacheMtimeStale(t *testing.T) {
	fs, vfs, cacheStore := newTestRemoteFs(t)
	defer vfs.destroy()

	const lineSize = 16
	const lineCount = 1024

	path := "/sparse-stale-mtime.txt"
	original := make([]byte, 0, lineCount*lineSize)
	for i := range lineCount {
		original = append(original, fmt.Sprintf("ORIGINAL-%06d\n", i)...)
	}

	cacheMtime := time.Now().Add(-time.Hour).Round(0)
	nodeMtime := cacheMtime.Add(time.Second)
	if _, err := cacheStore.Write(path, original, 0); err != nil {
		t.Fatalf("cache Write failed: %v", err)
	}
	if err := cacheStore.Commit(path, cacheEntryMetadata(path, int64(len(original)), cacheMtime)); err != nil {
		t.Fatalf("cache Commit failed: %v", err)
	}

	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(original)),
		modTime:      nodeMtime,
		creationTime: nodeMtime,
	})
	node.setDownloadURI("https://example.invalid/download")

	fs.backend = &fakeRemoteBackend{
		downloadFunc: fakeDownloadResponse(original, int64(len(original))),
	}

	var uploaded []byte
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		uploaded = append(uploaded[:0], data...)
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	errc, fh := fs.Open(path, fuse.O_RDWR)
	if errc != 0 {
		t.Fatalf("Open returned unexpected error: %d", errc)
	}
	defer fs.Release(path, fh)

	writes := []struct {
		offset  int64
		content []byte
	}{
		{0, []byte("CHANGE01-000000\n")},
		{int64(lineCount-1) * lineSize, []byte("CHANGE04-FINALE\n")},
	}
	for _, w := range writes {
		if n := fs.Write(path, w.content, w.offset, fh); n != len(w.content) {
			t.Fatalf("Write at offset %d returned %d, want %d", w.offset, n, len(w.content))
		}
	}
	if errc := fs.Flush(path, fh); errc != 0 {
		t.Fatalf("Flush returned unexpected error: %d", errc)
	}
	if len(uploaded) != len(original) {
		t.Fatalf("uploaded size = %d, want %d", len(uploaded), len(original))
	}

	unchangedLine := 100
	offset := unchangedLine * lineSize
	got := string(uploaded[offset : offset+lineSize])
	want := fmt.Sprintf("ORIGINAL-%06d\n", unchangedLine)
	if got != want {
		t.Fatalf("line %d after sparse write = %q, want %q", unchangedLine, got, want)
	}
}

func TestRemoteFsHydrationJoinedToPublicReadDownloadDoesNotCancel(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	cacheStore := newTestDiskCache(t)
	fs.cacheStore = cacheStore

	path := "/shared-download-baseline.txt"
	payload := bytes.Repeat([]byte("r"), cacheWriteSize+64)
	reader := newBlockingDownloadReader(payload, cacheWriteSize/2)
	modTime := time.Now().Add(-time.Minute).Round(0)
	node := vfs.getOrCreate(path, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         int64(len(payload)),
		modTime:      modTime,
		creationTime: modTime,
	})
	node.setDownloadURI("https://example.invalid/download")

	fs.backend = &fakeRemoteBackend{
		downloadFunc: func(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
			resp := &http.Response{
				StatusCode:    http.StatusOK,
				Body:          io.NopCloser(reader),
				ContentLength: int64(len(payload)),
			}
			if _, err := files_sdk.BuildResponse(resp, opts...); err != nil {
				return files_sdk.File{}, err
			}
			return files_sdk.File{Path: params.File.Path, Type: "file", Size: int64(len(payload))}, nil
		},
	}
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		data, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		return testUploadedMetadata(int64(len(data)), mtime), nil
	}

	errc, readFh := fs.Open(path, fuse.O_RDONLY)
	if errc != 0 {
		t.Fatalf("read Open returned unexpected error: %d", errc)
	}
	defer fs.Release(path, readFh)

	readDone := make(chan int, 1)
	go func() {
		buf := make([]byte, len(payload))
		readDone <- fs.Read(path, buf, 0, readFh)
	}()
	<-reader.firstWritten

	errc, writeFh := fs.Open(path, fuse.O_RDWR)
	if errc != 0 {
		t.Fatalf("write Open returned unexpected error: %d", errc)
	}
	defer fs.Release(path, writeFh)

	writeDone := make(chan int, 1)
	go func() {
		writeDone <- fs.Write(path, []byte("W"), 0, writeFh)
	}()

	deadline := time.Now().Add(2 * time.Second)
	for !node.hasActiveWriteSession() {
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for active write session")
		}
		time.Sleep(10 * time.Millisecond)
	}

	close(reader.release)

	select {
	case n := <-writeDone:
		if n != 1 {
			t.Fatalf("Write returned %d, want 1", n)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for write to finish")
	}

	select {
	case n := <-readDone:
		if n != len(payload) {
			t.Fatalf("Read returned %d, want %d", n, len(payload))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for read to finish")
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
	if err := cacheStore.Commit(path, cacheEntryMetadata(path, int64(len(initial)), initialMtime)); err != nil {
		t.Fatalf("cache Commit failed: %v", err)
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
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		_, err := io.ReadAll(reader)
		if err != nil {
			return uploadedFileMetadata{}, err
		}
		uploadedMtime = mtime
		return testUploadedMetadata(int64(len(initial)), mtime), nil
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

// fakeFileIter is a simple iterator over a slice of files_sdk.File for testing.
type fakeFileIter struct {
	files []files_sdk.File
	idx   int
}

func (it *fakeFileIter) Next() bool {
	if it.idx < len(it.files) {
		it.idx++
		return true
	}
	return false
}

func (it *fakeFileIter) File() files_sdk.File {
	return it.files[it.idx-1]
}

func (it *fakeFileIter) Err() error {
	return nil
}

func TestPendingVisibleChildPathsReturnsPendingNodes(t *testing.T) {
	_, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	dir := vfs.getOrCreate("/uploads", nodeTypeDir)
	dir.extendTtl()

	a := vfs.getOrCreate("/uploads/a.txt", nodeTypeFile)
	a.setPendingVisible()

	b := vfs.getOrCreate("/uploads/b.txt", nodeTypeFile)
	b.setPendingVisible()

	// c exists but is not pending — already confirmed remote
	vfs.getOrCreate("/uploads/c.txt", nodeTypeFile)

	// d is pending but in a different directory
	d := vfs.getOrCreate("/other/d.txt", nodeTypeFile)
	d.setPendingVisible()

	pending := vfs.pendingVisibleChildPaths("/uploads")
	if len(pending) != 2 {
		t.Fatalf("expected 2 pending paths, got %d: %v", len(pending), pending)
	}
	if _, ok := pending["/uploads/a.txt"]; !ok {
		t.Fatal("expected /uploads/a.txt in pending set")
	}
	if _, ok := pending["/uploads/b.txt"]; !ok {
		t.Fatal("expected /uploads/b.txt in pending set")
	}
}

func TestReaddirShowsPendingVisibleFilesAfterRemoteRefresh(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	dir := vfs.getOrCreate("/docs", nodeTypeDir)
	dir.extendTtl()

	type createdFile struct {
		path string
		fh   uint64
	}
	var created []createdFile

	// Create 3 files through the real Create→Write path
	for _, name := range []string{"one.pdf", "two.pdf", "three.pdf"} {
		path := fmt.Sprintf("/docs/%s", name)
		errc, fh := fs.Create(path, fuse.O_RDWR, 0o644)
		if errc != 0 {
			t.Fatalf("Create %s returned unexpected error: %d", name, errc)
		}
		payload := []byte("data")
		if n := fs.Write(path, payload, 0, fh); n != len(payload) {
			t.Fatalf("Write %s returned %d, want %d", name, n, len(payload))
		}
		created = append(created, createdFile{path: path, fh: fh})
	}

	// Release the handles before the refresh so the listing must rely on
	// pendingVisible rather than the open-handle merge.
	for _, file := range created {
		if errc := fs.Release(file.path, file.fh); errc != 0 {
			t.Fatalf("Release %s returned unexpected error: %d", file.path, errc)
		}
	}

	// Remote listing only returns one.pdf (the other two haven't propagated yet)
	backend := fs.backend.(*fakeRemoteBackend)
	backend.listForFunc = func(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (remoteFileIter, error) {
		return &fakeFileIter{files: []files_sdk.File{
			{DisplayName: "one.pdf", Size: 4, Type: "file"},
		}}, nil
	}

	// Expire directory info so loadDir refreshes from remote
	dir.expireInfo()

	// Readdir should show all 3 files: 1 from remote + 2 pending-visible
	var names []string
	fs.Readdir("/docs", func(name string, stat *fuse.Stat_t, ofst int64) bool {
		if name != "." && name != ".." {
			names = append(names, name)
		}
		return true
	}, 0, 0)

	slices.Sort(names)
	expected := []string{"one.pdf", "three.pdf", "two.pdf"}
	if !slices.Equal(names, expected) {
		t.Fatalf("Readdir names = %v, want %v", names, expected)
	}
}

func TestCreateNodeClearsPendingVisible(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	vfs.getOrCreate("/", nodeTypeDir)

	node := vfs.getOrCreate("/report.xlsx", nodeTypeFile)
	node.setPendingVisible()

	if !node.isPendingVisible() {
		t.Fatal("expected pendingVisible to be true after setPendingVisible")
	}

	// Simulate remote confirmation via createNode (called during listDir)
	confirmed := fs.createNode("/report.xlsx", files_sdk.File{
		DisplayName: "report.xlsx",
		Size:        500,
		Type:        "file",
	})

	if confirmed.isPendingVisible() {
		t.Fatal("expected pendingVisible to be false after createNode confirms remote existence")
	}
}

func TestCreateNodeStoresRemotePermissions(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	node := fs.createNode("/report.xlsx", files_sdk.File{
		DisplayName: "report.xlsx",
		Size:        500,
		Type:        "file",
		Permissions: "lrwd",
	})

	if got := node.getRemotePermissions(); got != "lrwd" {
		t.Fatalf("remote permissions = %q, want %q", got, "lrwd")
	}
	if !node.isReadable() {
		t.Fatal("expected node to be readable")
	}
	if !node.isWritable() {
		t.Fatal("expected node to be writable")
	}
}

func TestRemoteFsLockSkipsCreateWhenParentDirectoryIsReadOnly(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	fs.disableLocking = false
	fs.backend = &fakeRemoteBackend{
		createLockFunc: func(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.Lock, error) {
			t.Fatalf("expected remote lock creation to be skipped, got params %+v", params)
			return files_sdk.Lock{}, nil
		},
	}

	fs.createNode("/readonly", files_sdk.File{
		Path:        "/readonly",
		Type:        "directory",
		Permissions: "lr",
	})
	node := vfs.getOrCreate("/readonly/report.xlsx", nodeTypeFile)

	if errc := fs.lock(node, 123); errc != 0 {
		t.Fatalf("lock returned %d, want 0", errc)
	}
}

func TestRemoteFsLockCreatesRemoteLockWhenParentPermissionsUnknown(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	fs.disableLocking = false
	createLockCalled := false
	fs.backend = &fakeRemoteBackend{
		createLockFunc: func(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.Lock, error) {
			createLockCalled = true
			if params.Recursive == nil || *params.Recursive {
				t.Fatalf("lock recursive = %v, want false", params.Recursive)
			}
			return files_sdk.Lock{}, nil
		},
	}

	vfs.getOrCreate("/unknown", nodeTypeDir)
	node := vfs.getOrCreate("/unknown/report.xlsx", nodeTypeFile)

	if errc := fs.lock(node, 123); errc != 0 {
		t.Fatalf("lock returned %d, want 0", errc)
	}
	if !createLockCalled {
		t.Fatal("expected remote lock creation when parent permissions are unknown")
	}
}

func TestRemoteFsCreateReturnsEACCESWhenParentDirectoryIsReadOnly(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	parent := vfs.getOrCreate("/readonly", nodeTypeDir)
	parent.setRemotePermissions("lr")
	parent.extendTtl()

	errc, fh := fs.Create("/readonly/report.xlsx", fuse.O_RDWR|fuse.O_CREAT, 0o644)
	if errc != -fuse.EACCES {
		t.Fatalf("Create returned %d, want %d", errc, -fuse.EACCES)
	}
	if fh != ^uint64(0) {
		t.Fatalf("Create fh = %d, want invalid handle", fh)
	}
	if _, ok := vfs.fetch("/readonly/report.xlsx"); ok {
		t.Fatal("expected denied create not to add a child node")
	}
}

func TestRemoteFsOpenWriteReturnsEACCESWhenFileIsReadOnly(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	fs.createNode("/readonly.txt", files_sdk.File{
		Path:        "/readonly.txt",
		Type:        "file",
		Permissions: "lr",
	})

	errc, fh := fs.Open("/readonly.txt", fuse.O_RDWR)
	if errc != -fuse.EACCES {
		t.Fatalf("Open returned %d, want %d", errc, -fuse.EACCES)
	}
	if fh != ^uint64(0) {
		t.Fatalf("Open fh = %d, want invalid handle", fh)
	}
}

func TestRemoteFsAccessUsesKnownRemotePermissions(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	fs.createNode("/readonly.txt", files_sdk.File{
		Path:        "/readonly.txt",
		Type:        "file",
		Permissions: "lr",
	})
	fs.createNode("/writable.txt", files_sdk.File{
		Path:        "/writable.txt",
		Type:        "file",
		Permissions: "lrwd",
	})
	vfs.getOrCreate("/unknown.txt", nodeTypeFile)

	if errc := fs.Access("/readonly.txt", accessMaskRead); errc != 0 {
		t.Fatalf("Access(readonly, R_OK) returned %d, want 0", errc)
	}
	if errc := fs.Access("/readonly.txt", accessMaskWrite); errc != -fuse.EACCES {
		t.Fatalf("Access(readonly, W_OK) returned %d, want %d", errc, -fuse.EACCES)
	}
	if errc := fs.Access("/writable.txt", accessMaskWrite); errc != 0 {
		t.Fatalf("Access(writable, W_OK) returned %d, want 0", errc)
	}
	if errc := fs.Access("/unknown.txt", accessMaskWrite); errc != 0 {
		t.Fatalf("Access(unknown, W_OK) returned %d, want 0", errc)
	}
}

func TestGetStatUsesReadOnlyModesForKnownRemotePermissions(t *testing.T) {
	fileInfo := fsNodeInfo{nodeType: nodeTypeFile, remotePermissions: "lr"}
	fileStat := getStat(fileInfo, nil, 0, 0)
	wantFileMode := platformFileMode(0o444) & 0o777
	if got := fileStat.Mode & 0o777; got != wantFileMode {
		t.Fatalf("file mode = %o, want %o", got, wantFileMode)
	}

	dirInfo := fsNodeInfo{nodeType: nodeTypeDir, remotePermissions: "lr"}
	dirStat := getStat(dirInfo, nil, 0, 0)
	wantDirMode := platformFileMode(0o555) & 0o777
	if got := dirStat.Mode & 0o777; got != wantDirMode {
		t.Fatalf("dir mode = %o, want %o", got, wantDirMode)
	}

	writableFileInfo := fsNodeInfo{nodeType: nodeTypeFile, remotePermissions: "lrwd"}
	writableFileStat := getStat(writableFileInfo, nil, 0, 0)
	wantWritableFileMode := platformFileMode(0o644) & 0o777
	if got := writableFileStat.Mode & 0o777; got != wantWritableFileMode {
		t.Fatalf("writable file mode = %o, want %o", got, wantWritableFileMode)
	}
}

func TestRemoteFsMkdirReturnsEACCESWhenParentDirectoryIsReadOnly(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	parent := vfs.getOrCreate("/readonly", nodeTypeDir)
	parent.setRemotePermissions("lr")
	parent.extendTtl()

	if errc := fs.Mkdir("/readonly/child", 0o755); errc != -fuse.EACCES {
		t.Fatalf("Mkdir returned %d, want %d", errc, -fuse.EACCES)
	}
	if _, ok := vfs.fetch("/readonly/child"); ok {
		t.Fatal("expected denied mkdir not to add a child node")
	}
}

func TestRemoteFsUnlinkReturnsEACCESWhenParentDirectoryIsReadOnly(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	parent := vfs.getOrCreate("/readonly", nodeTypeDir)
	parent.setRemotePermissions("lr")
	parent.extendTtl()
	node := vfs.getOrCreate("/readonly/report.xlsx", nodeTypeFile)
	node.updateInfo(fsNodeInfo{nodeType: nodeTypeFile})
	node.extendTtl()

	if errc := fs.Unlink("/readonly/report.xlsx"); errc != -fuse.EACCES {
		t.Fatalf("Unlink returned %d, want %d", errc, -fuse.EACCES)
	}
}

func TestRemoteFsRmdirReturnsEACCESWhenParentDirectoryIsReadOnly(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	parent := vfs.getOrCreate("/readonly", nodeTypeDir)
	parent.setRemotePermissions("lr")
	parent.extendTtl()
	child := vfs.getOrCreate("/readonly/child", nodeTypeDir)
	child.updateInfo(fsNodeInfo{nodeType: nodeTypeDir})
	child.extendTtl()

	if errc := fs.Rmdir("/readonly/child"); errc != -fuse.EACCES {
		t.Fatalf("Rmdir returned %d, want %d", errc, -fuse.EACCES)
	}
}

func TestRemoteFsRenameReturnsEACCESWhenDestinationParentDirectoryIsReadOnly(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	srcParent := vfs.getOrCreate("/writable", nodeTypeDir)
	srcParent.setRemotePermissions("lrwd")
	srcParent.extendTtl()
	dstParent := vfs.getOrCreate("/readonly", nodeTypeDir)
	dstParent.setRemotePermissions("lr")
	dstParent.extendTtl()
	node := vfs.getOrCreate("/writable/report.xlsx", nodeTypeFile)
	node.updateInfo(fsNodeInfo{nodeType: nodeTypeFile})
	node.extendTtl()

	if errc := fs.Rename("/writable/report.xlsx", "/readonly/report.xlsx"); errc != -fuse.EACCES {
		t.Fatalf("Rename returned %d, want %d", errc, -fuse.EACCES)
	}
}

func TestUploadFailureClearsPendingVisible(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/fail.bin"

	// Force the upload to fail
	fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, p string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
		return uploadedFileMetadata{}, fmt.Errorf("simulated upload failure")
	}

	errc, fh := fs.Create(path, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	node, ok := vfs.fetch(path)
	if !ok {
		t.Fatal("expected node to exist after Create")
	}

	payload := []byte("will fail")
	if n := fs.Write(path, payload, 0, fh); n != len(payload) {
		t.Fatalf("Write returned %d, want %d", n, len(payload))
	}

	if !node.isPendingVisible() {
		t.Fatal("expected pendingVisible to be true after Write creates session")
	}

	// Flush triggers the upload which will fail
	fs.Flush(path, fh)

	if node.isPendingVisible() {
		t.Fatal("expected pendingVisible to be false after upload failure")
	}
}

func TestCreateWithoutWriteDoesNotSetPendingVisible(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/empty.txt"
	errc, fh := fs.Create(path, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	node, ok := vfs.fetch(path)
	if !ok {
		t.Fatal("expected node to exist after Create")
	}
	if node.isPendingVisible() {
		t.Fatal("expected pendingVisible to be false after Create with no Write")
	}

	fs.Release(path, fh)
}

func TestTruncateDoesNotSetPendingVisible(t *testing.T) {
	fs, vfs, _ := newTestRemoteFs(t)
	defer vfs.destroy()

	path := "/trunc.bin"
	errc, fh := fs.Create(path, fuse.O_RDWR, 0o644)
	if errc != 0 {
		t.Fatalf("Create returned unexpected error: %d", errc)
	}

	node, ok := vfs.fetch(path)
	if !ok {
		t.Fatal("expected node to exist after Create")
	}

	if errc := fs.Truncate(path, 0, fh); errc != 0 {
		t.Fatalf("Truncate returned unexpected error: %d", errc)
	}

	if node.isPendingVisible() {
		t.Fatal("expected pendingVisible to be false after Truncate without Write")
	}
}
