package fsmount

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	path_lib "path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/fsmount/events"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/mem"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/shell"
	fssync "github.com/Files-com/files-sdk-go/v3/fsmount/internal/sync"
	"github.com/Files-com/files-sdk-go/v3/ignore"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lock"
	"github.com/winfsp/cgofuse/fuse"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	api_key "github.com/Files-com/files-sdk-go/v3/apikey"
	file_migration "github.com/Files-com/files-sdk-go/v3/filemigration"
	ff "github.com/Files-com/files-sdk-go/v3/fsmount/internal/flags"
	lim "github.com/Files-com/files-sdk-go/v3/fsmount/internal/limit"
	gogitignore "github.com/sabhiram/go-gitignore"
)

const (
	// cacheWriteSize is the number of bytes read from an API response and written to the cache
	// before signalling any waiting cache readers that data is available
	cacheWriteSize = 128 * 1024 // 128KB

	// limits used in configuring the NewFuseOpLimiter for RemoteFs
	// S3 operations use HTTP/1.1 (6 connections is browser default for HTTP/1.1)
	downloadOpLimit = 6
	uploadOpLimit   = 6

	// Files.com API operations use HTTP/2 (multiplexed, single connection)
	otherOpLimit = 20

	// Global limit should accommodate peak concurrent operations
	globalOpLimit = 32
)

var (
	// compile time assertions that the cache implementations satisfy the fsCache interface
	_ cacheStore = (*disk.DiskCache)(nil)
	_ cacheStore = (*mem.MemoryCache)(nil)

	// webSyncInterval determines how frequently we ask Explorer to refresh any open folders.
	webSyncInterval = 15 * time.Second
)

const (
	accessMaskExecute = 1
	accessMaskWrite   = 2
	accessMaskRead    = 4
)

// RemoteFs is a file system that implements the logic for interacting with the Files.com API
// for a mounted file system. It handles all operations that are not handled by the LocalFs
// implementation, which is used for temporary files and files that should not be uploaded to
// Files.com. The Filescomfs implementation delegates operations to this implementation for
// all files who's source/destination is Files.com.
type RemoteFs struct {
	log              log.Logger
	vfs              *virtualfs
	config           *files_sdk.Config
	mountPoint       string
	localFsRoot      string
	root             string
	writeConcurrency int
	cacheTTL         time.Duration
	disableLocking   bool
	ignore           *gogitignore.GitIgnore

	providerBackend   ProviderBackend
	backend           remoteBackend
	currentUserId     int64
	uploadWorkingCopy func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error)
	lockMap           map[string]*lockInfo
	lockMapMutex      sync.Mutex
	loadDirMutexes    *fssync.PathMutex

	debugFuse bool

	initOnce sync.Once
	initTime time.Time

	cacheStore cacheStore

	gatesMu         sync.Mutex
	readyGates      map[string]*cache.ReadyGate
	gatePathMutexes *fssync.PathMutex

	events      events.EventPublisher
	transferSeq uint64

	ops *lim.FuseOpLimiter

	bufferPool *fssync.Pool[[]byte]

	webSyncTicker *time.Ticker
	stopCh        chan struct{}
	wg            sync.WaitGroup
}

// cacheStore defines the interface for the file system cache used by RemoteFs and allows for alternative
// implementations. e.g. an in-memory cache implementation vs a disk-based cache implementation.
type cacheStore interface {
	Read(path string, buff []byte, ofst int64) (n int, err error)
	ReadComplete(path string, meta cache.EntryMetadata, buff []byte, ofst int64) (n int, err error)
	ReadPartial(path string, buff []byte, ofst int64) (n int, err error)
	Write(path string, buff []byte, ofst int64) (n int, err error)
	WritePartial(path string, buff []byte, ofst int64) (n int, err error)
	Commit(path string, meta cache.EntryMetadata) error
	Delete(path string) bool
	DeletePartial(path string) bool
	StartMaintenance()
	StopMaintenance()
	// Pin increments the reference count for a file, preventing it from being evicted
	Pin(path string)
	PinPartial(path string)
	// Unpin decrements the reference count for a file
	Unpin(path string)
	UnpinPartial(path string)
}

// cacheReader wraps the cacheStore to provide an io.Reader interface for reading cached files.
// This is used to seed working copies with cached content for partial file updates.
type cacheReader struct {
	cache  cacheStore
	path   string
	size   int64
	offset int64
	logger log.Logger
}

func (cr *cacheReader) Read(p []byte) (n int, err error) {
	if cr.offset >= cr.size {
		return 0, io.EOF
	}

	n, err = cr.cache.Read(cr.path, p, cr.offset)
	if err != nil {
		cr.logger.Debug("cacheReader: error reading from cache at offset %d: %v", cr.offset, err)
		return 0, err
	}

	if n == 0 {
		// Cache doesn't have this data - return EOF to prevent blocking
		return 0, io.EOF
	}

	cr.offset += int64(n)
	return n, nil
}

// cacheWriterAdapter adapts an fsio.CacheWriter to io.Writer for use with io.TeeReader.
// This allows simultaneous writing to cache while reading from a source during uploads.
type cacheWriterAdapter struct {
	writer func([]byte, int64) (int, error)
	offset *int64
}

func (cw *cacheWriterAdapter) Write(p []byte) (n int, err error) {
	n, err = cw.writer(p, *cw.offset)
	if err == nil {
		*cw.offset += int64(n)
	}
	return n, err
}

type uploadWorkingCopyReader interface {
	io.Reader
	io.ReaderAt
	Stat() (os.FileInfo, error)
}

type uploadedFileMetadata struct {
	size    int64
	modTime time.Time
}

type uploadFileResult struct {
	metadata uploadedFileMetadata
	err      error
}

func (fs *RemoteFs) uploadedFileMetadata(path string, uploaded files_sdk.File, fallbackSize int64, fallbackModTime time.Time) uploadedFileMetadata {
	size := uploaded.Size
	if size == 0 && fallbackSize != 0 {
		fs.log.Debug("RemoteFs: upload response missing size for %v; using local size %d", path, fallbackSize)
		size = fallbackSize
	}
	modTime := uploaded.ModTime()
	if modTime.IsZero() {
		fs.log.Debug("RemoteFs: upload response missing mtime for %v; using local mtime %v", path, fallbackModTime)
		modTime = fallbackModTime
	}
	return uploadedFileMetadata{
		size:    size,
		modTime: modTime,
	}
}

func (fs *RemoteFs) providerUploadedSize(path string, uploaded ProviderEntry, fallbackSize int64) int64 {
	if uploaded.Size != 0 || fallbackSize == 0 {
		return uploaded.Size
	}
	fs.log.Warn("RemoteFs: provider upload response reported zero size for non-empty file %v; using local size %d", path, fallbackSize)
	return fallbackSize
}

type progressReader struct {
	reader   io.Reader
	progress func(int64)
}

func (r *progressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n > 0 && r.progress != nil {
		r.progress(int64(n))
	}
	return n, err
}

func cacheEntryMetadata(path string, size int64, modTime time.Time) cache.EntryMetadata {
	return cache.NewEntryMetadata(path, size, modTime)
}

type lockInfo struct {
	Fh   uint64
	Lock *files_sdk.Lock
}

func newRemoteFs(params MountParams, vfs *virtualfs, log log.Logger, cs cacheStore) (*RemoteFs, error) {
	if params.Root == "" {
		params.Root = "/"
	}

	limiter := lim.NewFuseOpLimiter(map[lim.FuseOpType]int64{
		lim.FuseOpDownload: downloadOpLimit,
		lim.FuseOpUpload:   uploadOpLimit,
		lim.FuseOpOther:    otherOpLimit,
	}, globalOpLimit)

	fs := &RemoteFs{
		log:              log,
		root:             params.Root,
		vfs:              vfs,
		mountPoint:       params.MountPoint,
		localFsRoot:      params.TmpFsPath,
		writeConcurrency: params.WriteConcurrency,
		cacheTTL:         params.CacheTTL,
		config:           params.Config,
		disableLocking:   params.DisableLocking,
		debugFuse:        params.DebugFuse,
		events:           params.EventPublisher,
		cacheStore:       cs,
		readyGates:       make(map[string]*cache.ReadyGate),
		gatePathMutexes:  fssync.NewPathMutex(),
		ops:              limiter,
		loadDirMutexes:   fssync.NewPathMutex(),
		bufferPool: fssync.NewPool(func() []byte {
			return make([]byte, cacheWriteSize)
		}),
	}
	if params.ProviderBackend != nil {
		fs.providerBackend = params.ProviderBackend
		fs.disableLocking = true
		fs.backend = &providerRemoteBackend{provider: params.ProviderBackend}
		fs.uploadWorkingCopy = func(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
			fileInfo, err := reader.Stat()
			if err != nil {
				return uploadedFileMetadata{}, err
			}
			remotePath := fs.remotePath(path)
			writeReader := io.Reader(reader)
			if node != nil {
				writeReader = &progressReader{
					reader:   reader,
					progress: fs.uploadProgressFunc(node),
				}
			}
			entry, err := params.ProviderBackend.Write(ctx, remotePath, writeReader, fileInfo.Size(), mtime)
			if err != nil {
				return uploadedFileMetadata{}, err
			}
			uploadedModTime := entry.ModTime
			if uploadedModTime.IsZero() {
				uploadedModTime = mtime
			}
			return uploadedFileMetadata{
				size:    fs.providerUploadedSize(remotePath, entry, fileInfo.Size()),
				modTime: uploadedModTime,
			}, nil
		}
	}

	// ensure write concurrency and cache TTL are positive
	if fs.writeConcurrency <= 0 {
		fs.writeConcurrency = DefaultWriteConcurrency
	}
	if fs.cacheTTL <= 0 {
		fs.cacheTTL = DefaultCacheTTL
	}
	if params.IgnorePatterns == nil || len(params.IgnorePatterns) > 0 {
		ignore, err := ignore.New(params.IgnorePatterns...)
		if err != nil {
			return nil, err
		}
		fs.ignore = ignore
	}
	if fs.events == nil {
		fs.events = &events.NoOpEventPublisher{}
	}
	return fs, nil
}

func (fs *RemoteFs) Init() {
	// Guard with a sync.Once because validation and cgofuse startup can both initialize the remote file system.
	fs.initOnce.Do(func() {
		if fs.backend == nil {
			fs.backend = &sdkRemoteBackend{
				fileClient:      &file.Client{Config: *fs.config},
				lockClient:      &lock.Client{Config: *fs.config},
				apiKeyClient:    &api_key.Client{Config: *fs.config},
				migrationClient: &file_migration.Client{Config: *fs.config},
			}
			fs.lockMap = make(map[string]*lockInfo)
		}

		// Skip the API key lookup on provider mounts: there is no Files.com API
		// key, locking is disabled, and currentUserId is unused on this path.
		if fs.providerBackend == nil {
			// no need to guard this with an operation limit since it's only called once during initialization
			key, err := fs.backend.findCurrent()
			if err != nil {
				fs.log.Error("Failed to find metadata for current API key, file exclusivity locks may not work as expected: %v", err)
				// set locking to false?
			}
			fs.currentUserId = key.UserId
		}
		// store the time the file system was initialized to use as the creation time for the root directory
		fs.initTime = time.Now()
		fs.log.Debug("RemoteFs: RemoteFs initialized successfully. Remote file system root: %s", fs.root)
	})
	// start the disk cache maintenance goroutine
	// this does not block and ensures only one goroutine is started
	fs.cacheStore.StartMaintenance()

	// start the web sync watcher on Windows, which periodically checks for changes made via the web interface
	// and notifies the OS to refresh the directory view
	if runtime.GOOS == "windows" {
		fs.startWebSyncWatcher(webSyncInterval)
	}
}

func (fs *RemoteFs) Destroy() {
	fs.log.Debug("RemoteFs: Destroy: removing all file locks")

	fs.lockMapMutex.Lock()
	defer fs.lockMapMutex.Unlock()
	for path, lockInfo := range fs.lockMap {
		fs.unlock(path, lockInfo.Fh)
	}

	fs.log.Debug("RemoteFs: Destroy: stopping cache maintenance")
	fs.cacheStore.StopMaintenance()
	fs.stopWebSync()
}

func (fs *RemoteFs) Validate() error {
	fs.Init()
	// Make sure the root directory can be listed.
	// no need to guard this with an operation limit since it's only called once during initialization
	it, err := fs.backend.listFor(files_sdk.FolderListForParams{
		Path: fs.remotePath("/"),
		ListParams: files_sdk.ListParams{
			PerPage: 1,
		},
	})
	if err == nil {
		it.Next() // Get 1 item. This is what actually triggers the API call.
		err = it.Err()
	}
	return err
}

func (fs *RemoteFs) denyIfKnownReadOnlyParent(path string, operation string) int {
	parentPath := path_lib.Dir(path)
	if parentPath == path {
		return 0
	}

	parent, ok := fs.vfs.fetch(parentPath)
	if !ok {
		return 0
	}

	if permissions := parent.getRemotePermissions(); permissions != "" && !parent.isWritable() {
		localPath, remotePath := fs.paths(path)
		fs.log.Debug(
			"mount_permission_denied source=cached_parent_permissions op=%s path=%q remote_path=%q local_path=%q parent=%q permissions=%q errc=%d errno=%s",
			operation,
			path,
			remotePath,
			localPath,
			parentPath,
			permissions,
			-fuse.EACCES,
			formatFuseErrno(-fuse.EACCES),
		)
		fs.log.Debug("RemoteFs: %s: parent directory is read-only: %v (%v) parent=%v permissions=%q", operation, remotePath, localPath, parentPath, permissions)
		return -fuse.EACCES
	}

	return 0
}

func (fs *RemoteFs) denyIfKnownReadOnlyPath(node *fsNode, operation string) int {
	if node == nil {
		return 0
	}

	if permissions := node.getRemotePermissions(); permissions != "" && !node.isWritable() {
		localPath, remotePath := fs.paths(node.path)
		fs.log.Debug(
			"mount_permission_denied source=cached_path_permissions op=%s path=%q remote_path=%q local_path=%q permissions=%q errc=%d errno=%s",
			operation,
			node.path,
			remotePath,
			localPath,
			permissions,
			-fuse.EACCES,
			formatFuseErrno(-fuse.EACCES),
		)
		fs.log.Debug("RemoteFs: %s: path is read-only: %v (%v) permissions=%q", operation, remotePath, localPath, permissions)
		return -fuse.EACCES
	}

	return 0
}

func (fs *RemoteFs) Mkdir(path string, mode uint32) (errc int) {
	localPath, remotePath := fs.paths(path)
	fs.log.Debug("RemoteFs: Mkdir: %v (%v) (mode=%o)", remotePath, localPath, mode)

	if node, ok := fs.vfs.fetch(path); ok && !node.infoExpired() {
		fs.log.Debug("RemoteFs: Mkdir: path already exists in VFS: %v (%v)", remotePath, localPath)
		return -fuse.EEXIST
	}

	if errc = fs.loadParent(path); errc != 0 {
		return errc
	}
	if errc = fs.denyIfKnownReadOnlyParent(path, "Mkdir"); errc != 0 {
		return errc
	}
	if node, ok := fs.vfs.fetch(path); ok && !node.infoExpired() {
		fs.log.Debug("RemoteFs: Mkdir: path discovered during parent refresh: %v (%v)", remotePath, localPath)
		return -fuse.EEXIST
	}

	err := fs.ops.TryWithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		_, err := fs.backend.createFolder(files_sdk.FolderCreateParams{Path: remotePath})
		return err
	})
	if errc = fs.handleError(path, err); errc != 0 {
		return errc
	}

	node := fs.vfs.getOrCreate(path, nodeTypeDir)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeDir,
		size:         0,
		modTime:      time.Now(),
		creationTime: time.Now(),
	})

	return errc
}

func (fs *RemoteFs) Unlink(path string) (errc int) {
	localPath, remotePath := fs.paths(path)
	node, errc := fs.fetchNodeWithParentRefresh(path)
	if errc != 0 {
		// If the node doesn't exist, it can not be deleted.
		fs.log.Debug("RemoteFs: Unlink: File not found: %v (%v)", remotePath, localPath)
		return errc
	}
	if node.info.nodeType == nodeTypeDir {
		fs.log.Debug("RemoteFs: Unlink: Path is a directory: %v (%v)", remotePath, localPath)
		return -fuse.EISDIR
	}
	if errc = fs.denyIfKnownReadOnlyParent(path, "Unlink"); errc != 0 {
		return errc
	}

	// If the node is locked, it can not be deleted.
	if node.isLocked() {
		fs.lockMapMutex.Lock()
		linfo, ok := fs.lockMap[path]
		fs.lockMapMutex.Unlock()
		if ok && fs.currentUserId == linfo.Lock.UserId {
			fs.log.Debug("RemoteFs: Unlink: allowing delete of same-user locked file: %v (%v)", remotePath, localPath)
		} else {
			fs.log.Info("Cannot delete locked file: %v (%v)", remotePath, localPath)
			return -fuse.ENOLCK
		}
	}

	// If the node is being written to, cancel the upload and delete the file from the remote API.
	// This is necessary because the file may be in the middle of being written to, and the upload may not have completed yet.
	if node.hasActiveWriteSession() {
		fs.log.Debug("RemoteFs: Unlink: Canceling active write session for: %v (%v)", remotePath, localPath)
		node.cancelUpload()
	}

	// The fs may have been in the middle of writing the file, so don't log until here.
	fs.log.Info("Deleting file: %v (%v)", remotePath, localPath)
	return fs.delete(path)
}

func (fs *RemoteFs) Rmdir(path string) (errc int) {
	localPath, remotePath := fs.paths(path)
	node, errc := fs.fetchNodeWithParentRefresh(path)
	if errc != 0 {
		// If the node doesn't exist, it can not be deleted.
		fs.log.Debug("RemoteFs: Rmdir: directory not found: %v (%v)", remotePath, localPath)
		return errc
	}

	if node.info.nodeType != nodeTypeDir {
		fs.log.Debug("RemoteFs: Rmdir: Path is not a directory: %v (%v)", remotePath, localPath)
		return -fuse.ENOTDIR
	}
	if errc = fs.denyIfKnownReadOnlyParent(path, "Rmdir"); errc != 0 {
		return errc
	}
	fs.log.Info("Deleting folder: %v (%v)", remotePath, localPath)

	return fs.delete(path)
}

func (fs *RemoteFs) Rename(oldpath string, newpath string) (errc int) {
	fs.log.Debug("RemoteFs: Rename: oldpath=%v, newpath=%v", oldpath, newpath)
	oldLocalPath, oldRemotePath := fs.paths(oldpath)
	newLocalPath, newRemotePath := fs.paths(newpath)

	node, ok := fs.vfs.fetch(oldpath)
	if !ok {
		return -fuse.ENOENT
	}
	if errc = fs.denyIfKnownReadOnlyParent(oldpath, "Rename"); errc != 0 {
		return errc
	}
	if errc = fs.denyIfKnownReadOnlyParent(newpath, "Rename"); errc != 0 {
		return errc
	}

	// Reserve a committed retained session before releasing its locks for the remote move. Writes and
	// truncates that were already preparing a local change finish first. No network operation runs
	// while either mutex is held.
	node.writeMu.Lock()
	session := node.writeSession
	waitForUpload := false
	if session == nil {
		if node.isUnmaterialized() {
			// Create can produce a local node without starting a write session or
			// creating a remote file. Rename that placeholder locally so callers
			// do not receive Not Found from a backend move of a nonexistent path.
			fs.log.Info("Renaming unmaterialized file %v to %v (%v to %v)", oldRemotePath, newRemotePath, oldLocalPath, newLocalPath)
			fs.rename(oldpath, newpath)
			_ = fs.cacheStore.Delete(oldpath)
			_ = fs.cacheStore.Delete(newpath)
			node.writeMu.Unlock()
			return 0
		}
		node.writeMu.Unlock()
		defer node.expireInfo()
	} else {
		session.mu.Lock()
		if session.renaming {
			session.mu.Unlock()
			node.writeMu.Unlock()
			return -fuse.EAGAIN
		}

		if !session.remoteCommitted && !session.uploadPathFixed {
			// The first upload has not captured its final destination yet. Update the local path so its
			// completion callback reads the new destination.
			fs.rename(oldpath, newpath)
			session.path = newpath
			node.expireInfo()
			session.mu.Unlock()
			node.writeMu.Unlock()
			return errc
		}

		// A later upload does not replace the already-committed remote object until it finishes.
		// Keep its destination at the old path, then move the completed object below so the old path
		// cannot be left behind as an orphan.
		waitForUpload = session.uploading || session.finalizing
		session.renaming = true
		node.writeMu.Unlock()
		for session.mutationCount > 0 {
			session.cond.Wait()
		}
		session.mu.Unlock()
	}
	if waitForUpload {
		if err := node.waitForUploadWithProgressTimeout(fsyncTimeout); err != nil {
			node.expireInfo()
			_ = node.finishWriteSessionRename(session, newpath, false)
			if errors.Is(err, context.DeadlineExceeded) {
				return -fuse.ETIMEDOUT
			}
			if errc := fs.handleUploadSessionError(oldpath, err, session); errc != 0 {
				return errc
			}
			return -fuse.EIO
		}
	}
	// With no write session, or with a retained session that has committed an upload, the remote
	// file already exists at oldRemotePath and must be moved explicitly.
	fs.log.Info("Renaming %v to %v (%v to %v)", oldRemotePath, newRemotePath, oldLocalPath, newLocalPath)

	err := fs.ops.TryWithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		action, err := fs.backend.move(files_sdk.FileMoveParams{
			Path:        oldRemotePath,
			Destination: newRemotePath,
			// FUSE Rename only provides old and new paths; it does not tell us
			// whether the caller already confirmed a replacement. Keep overwrite=true
			// so confirmed replacement flows complete instead of failing remotely.
			Overwrite: lib.Ptr(true),
		}, files_sdk.WithContext(ctx))
		if err != nil {
			return err
		}
		return fs.waitForAction(ctx, action, "move")
	})
	if errors.Is(err, lim.ErrNoSlotsAvailable) {
		if session != nil {
			node.expireInfo()
			_ = node.finishWriteSessionRename(session, newpath, false)
		}
		return -fuse.EAGAIN
	}
	if errc = fs.handleError(oldpath, err); errc != 0 {
		if session != nil {
			node.expireInfo()
			_ = node.finishWriteSessionRename(session, newpath, false)
		}
		return errc
	}

	// Stop an old source or destination download from republishing stale cache data after the move.
	// The per-path locks keep replacement downloads out until the VFS and cache reflect the new file.
	unlockGates := fs.lockGatePaths(oldpath, newpath)
	for _, readyGate := range fs.takeGatesForPaths(oldpath, newpath) {
		readyGate.CancelAndWait()
		readyGate.Cleanup()
	}
	fs.rename(oldpath, newpath)
	_ = fs.cacheStore.Delete(oldpath)
	_ = fs.cacheStore.Delete(newpath)
	unlockGates()
	if session != nil {
		node.expireInfo()
		if err := node.finishWriteSessionRename(session, newpath, true); err != nil {
			fs.log.Debug("RemoteFs: Rename: failed clearing retained write session for %v: %v", newpath, err)
		}
	}

	return errc
}

func (fs *RemoteFs) Utimens(path string, tmsp []fuse.Timespec) (errc int) {
	localPath, remotePath := fs.paths(path)
	modT := tmsp[1].Time()
	fs.log.Debug("RemoteFs: Utimens: Updating mtime for: %v (%v) (mtime=%v)", remotePath, localPath, modT)

	node, errc := fs.fetchNodeWithParentRefresh(path)
	if errc != 0 {
		return errc
	}

	node.info.modTime = modT

	if session := node.getWriteSession(); session != nil {
		session.mu.Lock()
		session.mtime = modT
		session.mtimeExplicit = true
		session.mu.Unlock()
		return errc
	}

	err := fs.ops.TryWithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		_, err := fs.backend.update(files_sdk.FileUpdateParams{
			Path:          remotePath,
			ProvidedMtime: &node.info.modTime,
		})
		return err
	})
	if errors.Is(err, lim.ErrNoSlotsAvailable) {
		return -fuse.EAGAIN
	}

	return fs.handleError(path, err)
}

func (fs *RemoteFs) Create(path string, flags int, mode uint32) (errc int, fh uint64) {
	localPath, remotePath := fs.paths(path)
	fuseFlags := ff.NewFuseFlags(flags)
	var handle *fileHandle

	fs.log.Debug("RemoteFs: Create: Creating file: %v (%v) (flags=%v, mode=%o)", remotePath, localPath, fuseFlags, mode)

	// Load the parent directory to populate the vfs nodes map
	if errc = fs.loadParent(path); errc != 0 {
		return errc, ^uint64(0)
	}
	if errc = fs.denyIfKnownReadOnlyParent(path, "Create"); errc != 0 {
		return errc, ^uint64(0)
	}

	node, exists := fs.vfs.fetch(path)
	if exists && node.hasActiveWriteSession() {
		fh, handle = fs.vfs.handles.Open(nil, fuseFlags)
		fs.log.Debug("RemoteFs: Create: joining existing write session: %v (%v)", remotePath, localPath)
		handle.node = node
		session := node.getWriteSession()
		session.addHandle(fh)
		fs.logWriteSessionMilestone(path, "create_join_existing_session", fh, session, "flags=%v mode=%o", fuseFlags, mode)
		return errc, fh
	}

	// TODO: decide if this makes sense. the node exists and the cache data is recent
	// so return an error for the Create call?
	if exists && !node.infoExpired() {
		fs.log.Debug("RemoteFs: Create: Node exists, cache data is recent, but no open writer: %v (%v)", remotePath, localPath)
		return -fuse.EEXIST, ^uint64(0)
	}

	if !exists {
		node = fs.vfs.getOrCreate(path, nodeTypeFile)
		node.markUnmaterialized()
	}

	fh, handle = fs.vfs.handles.Open(nil, fuseFlags)
	node.updateSize(0)
	handle.node = node
	fs.logWriteSessionMilestone(path, "create_handle_opened", fh, nil, "flags=%v mode=%o", fuseFlags, mode)

	// Pin the file in cache to prevent eviction while the handle is open
	// NOTE: If an error is returned below, Unpin() must be called manually since FUSE will
	// NOT call Release() for failed create operations.
	if fs.cacheStore != nil {
		fs.cacheStore.Pin(path)
	}

	if errc = fs.lock(node, fh); errc != 0 {
		// Lock failed - manually Unpin and Release handle since Release() won't be called for failed creates
		if fs.cacheStore != nil {
			fs.cacheStore.Unpin(path)
		}
		fs.vfs.handles.Release(fh)
		fs.logWriteSessionMilestone(path, "create_lock_failed", fh, nil, "errc=%d errno=%s", errc, formatFuseErrno(errc))
		return errc, fh
	}
	fs.logWriteSessionMilestone(path, "create_lock_acquired", fh, nil, "")

	fs.log.Debug("RemoteFs: Create: opened write-capable handle for %v (%v), fh=%v", remotePath, localPath, fh)
	fs.logWriteSessionMilestone(path, "create_write_handle_ready", fh, nil, "")

	return errc, fh
}

func (fs *RemoteFs) Open(path string, flags int) (errc int, fh uint64) {
	fuseFlags := ff.NewFuseFlags(flags)
	node, errc := fs.fetchNodeWithParentRefresh(path)
	if errc != 0 {
		return errc, ^uint64(0)
	}

	if !fuseFlags.IsReadOnly() {
		if errc := fs.denyIfKnownReadOnlyPath(node, "Open"); errc != 0 {
			return errc, ^uint64(0)
		}
	}

	fh, handle := fs.vfs.handles.Open(node, fuseFlags)
	fs.log.Trace("RemoteFs: Open: path=%v, flags=%v, fh=%v", path, fuseFlags, fh)

	// Pin the file in cache to prevent eviction while the handle is open.
	// This must happen early, before any error returns, to ensure the file remains
	// cached for the lifetime of the handle (both read and write operations).
	// NOTE: If an error is returned below, Unpin() must be called manually since FUSE will
	// NOT call Release() for failed open operations.
	if fs.cacheStore != nil {
		fs.cacheStore.Pin(path)
	}

	// TODO: this can succeed even if the file doesn't exist. The file may be created
	// later when the file is written to, or it may never be created if the file
	// is never written to. Decide if this is the desired behavior.
	if handle.IsReadOnly() {
		return errc, fh
	}
	// after this point, the requested op must be a write operation

	if errc = fs.lock(node, fh); errc != 0 {
		// Lock failed - manually Unpin and Release handle since Release() won't be called for failed opens
		if fs.cacheStore != nil {
			fs.cacheStore.Unpin(path)
		}
		fs.vfs.handles.Release(fh)
		fs.logWriteSessionMilestone(path, "open_lock_failed", fh, nil, "flags=%v errc=%d errno=%s", fuseFlags, errc, formatFuseErrno(errc))
		return errc, fh
	}
	fs.logWriteSessionMilestone(path, "open_lock_acquired", fh, nil, "flags=%v", fuseFlags)

	// A node with an active write session is already in write-owned state.
	// Additional write-capable handles join the session instead of being rejected.
	if session := node.getWriteSession(); session != nil {
		session.addHandle(fh)
		fs.log.Debug("RemoteFs: Open: joined active write session for path=%v, fh=%v", path, fh)
		fs.logWriteSessionMilestone(path, "open_join_existing_session", fh, session, "flags=%v", fuseFlags)
		return errc, fh
	}
	fs.log.Debug("RemoteFs: Open: opened write-capable handle for path=%v, fh=%v", path, fh)
	fs.logWriteSessionMilestone(path, "open_write_handle_ready", fh, nil, "flags=%v", fuseFlags)

	return errc, fh
}

func (fs *RemoteFs) fetchNodeWithParentRefresh(path string) (node *fsNode, errc int) {
	if node, ok := fs.vfs.fetch(path); ok && !node.infoExpired() {
		return node, 0
	}

	if errc = fs.loadParent(path); errc != 0 {
		return nil, errc
	}

	node, ok := fs.vfs.fetch(path)
	if !ok {
		return nil, -fuse.ENOENT
	}

	return node, 0
}

func (fs *RemoteFs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	fs.vfs.ensureContextOwner()
	if stat == nil {
		stat = &fuse.Stat_t{}
	}
	// If the file handle is open, extend the TTL of the open handle. The info may have expired,
	// but the handle is still open, meaning the OS is still using the file. This can happen if there
	// are multiple simultaneous uploads, but they haven't all received a write request in the last
	// cacheTTL duration. If the Getattr call returns an error, the OS will remove the file from the
	// Explorer/Finder window until the upload completes, and a subsequent Getattr call succeeds, which
	// is a bad user experience.
	fs.vfs.handles.ExtendOpenHandleTtls()
	if node, exists := fs.vfs.fetch(path); exists && !node.infoExpired() {
		if session := node.getWriteSession(); session != nil {
			snap := session.snapshot()
			info := node.info
			info.size = snap.currentSize
			info.modTime = snap.mtime
			getStat(info, stat, fs.vfs.uid, fs.vfs.gid)
			fs.log.Trace("RemoteFs: Getattr: returning cached write-session stat: path=%v, fh=%v, stat=%s", path, fh, formatFuseStat(stat))
			return errc
		}
		getStat(node.info, stat, fs.vfs.uid, fs.vfs.gid)
		fs.log.Trace("RemoteFs: Getattr: returning cached node stat: path=%v, fh=%v, stat=%s", path, fh, formatFuseStat(stat))
		return errc
	}

	if errc = fs.loadParent(path); errc != 0 {
		fs.log.Trace("RemoteFs: Getattr: loadParent failed: path=%v, fh=%v, errc=%v", path, fh, errc)
		return errc
	}

	node, exists := fs.vfs.fetch(path)
	if !exists {
		node = nil

		if fs.isLockFile(path) {
			if lockedNode, exists := fs.vfs.fetchLockTarget(path); exists && lockedNode.isLocked() {
				node = fs.vfs.getOrCreate(path, nodeTypeFile)
				node.updateInfo(fsNodeInfo{
					size:    int64(len(buildOwnerFile(lockedNode))),
					modTime: time.Now(),
				})
			}
		}

		if node == nil {
			localPath, remotePath := fs.paths(path)
			fs.log.Trace("RemoteFs: Getattr: file not found: %v (%v), fh=%v", remotePath, localPath, fh)
			return -fuse.ENOENT
		}
	}
	if session := node.getWriteSession(); session != nil {
		snap := session.snapshot()
		info := node.info
		info.size = snap.currentSize
		info.modTime = snap.mtime
		getStat(info, stat, fs.vfs.uid, fs.vfs.gid)
		return errc
	}
	getStat(node.info, stat, fs.vfs.uid, fs.vfs.gid)

	return errc
}

func (fs *RemoteFs) Truncate(path string, size int64, fh uint64) (errc int) {
	// The word truncate is overloaded here. The intention is to set the size of the
	// file to the size getting passed in, NOT to truncate the file to zero bytes.

	localPath, remotePath := fs.paths(path)
	fs.log.Debug("RemoteFs: Truncate: %v (%v) (size=%v, fh=%v)", remotePath, localPath, size, fh)

	var node *fsNode
	if fh != ^uint64(0) {
		if handle, existingNode, ok := fs.vfs.handles.Lookup(fh); ok && handle != nil && existingNode != nil {
			node = existingNode
		}
	}
	if node == nil {
		node, errc = fs.fetchNodeWithParentRefresh(path)
		if errc != 0 {
			return errc
		}
	}

	// Invalidate any cached content. The size has changed, so cached data is stale.
	// Without this, a subsequent write could load stale cached content as the
	// working copy baseline and preserve data from the wrong version of the file.
	fs.cacheStore.Delete(path)

	session, _, err := node.beginWriteSessionMutation(path)
	if err != nil {
		fs.log.Error("RemoteFs: Truncate: failed to create write session for %v: %v", path, err)
		return -fuse.EIO
	}
	defer session.endMutation()
	path = session.snapshot().path
	session.addHandle(fh)
	if err := fs.ensureWriteSessionBaseline(path, node, session, size == 0, fh); err != nil {
		if errc := fs.handleUploadSessionError(path, err, session); errc != 0 {
			return errc
		}
		return -fuse.EIO
	}
	if err := fs.truncateWorkingCopy(session, size); err != nil {
		fs.log.Error("RemoteFs: Truncate: working copy truncate failed for %v: %s", path, uploadLogMessage(err))
		return -fuse.EIO
	}
	node.extendTtl()
	fs.log.Debug("RemoteFs: Truncate: updated working copy for %v (%v)", remotePath, localPath)

	return errc
}

func (fs *RemoteFs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	buffLen := int64(len(buff))
	// unlikely, but guard against a zero-length read
	if buffLen == 0 {
		return 0
	}

	fs.log.Trace("RemoteFs: Read: path=%v, len=%v, ofst=%v, fh=%v", path, buffLen, ofst, fh)

	handle, node, ok := fs.vfs.handles.Lookup(fh)
	if !ok {
		fs.log.Debug("RemoteFs: Read: file handle %v not found for path %v", fh, path)
		return -fuse.EBADF
	}

	if n, sessionOwned, err := node.readFromWriteSession(buff, ofst); sessionOwned {
		if err != nil {
			fs.log.Debug("RemoteFs: Read: write session read failed for path=%v: %v", path, err)
			return -fuse.EIO
		}
		if n > 0 {
			handle.incrementRead(int64(n))
		}
		return n
	}

	// the following operations all benefit from knowing the file size, so attempt to get the most
	// up-to-date size possible before proceeding.
	size := node.info.size
	if node.infoExpired() || size <= 0 {
		// in order to get the most up-to-date size the parent directory's info must be expired as well,
		// otherwise the getattr call will return a cached stat that may not reflect the current size of the file.
		fs.vfs.expireNodeInfo(path)
		var st fuse.Stat_t
		// ignore errno and fall back to range-from-ofst if still unknown
		_ = fs.Getattr(path, &st, fh)
		size = node.info.size
	}

	// make sure offset is not negative, and if size is known, that offset is not past EOF
	if ofst < 0 {
		ofst = 0
	}

	if size == 0 {
		fs.log.Trace("RemoteFs: Read: file is empty, returning EOF")
		return 0
	}

	// attempting to read past EOF
	if size > 0 && ofst >= size {
		fs.log.Trace("RemoteFs: Read: offset %d is greater than file size %d, returning EOF", ofst, size)
		return 0
	}

	// TODO: determine if this is still needed, or if it needs to move to the localfs, these files are written to the local fs and not
	// stored remotely anymore
	if fs.isLockFile(path) {
		if lockedNode, ok := fs.vfs.fetchLockTarget(path); ok && lockedNode.isLocked() {
			ownerBuffer := buildOwnerFile(lockedNode)
			return copy(buff, ownerBuffer[ofst:])
		}
	}

	// Attempt to read from the cache only when no download is in progress for this path.
	// If a download is active, the cache entry is partially written and a Read call may return
	// fewer bytes than len(buff) (a short read). By skipping the early cache check when a gate is
	// active, all reads during an active download go through WaitFor, which blocks until the
	// full requested range is available, guaranteeing a non-short read.
	if _, isDownloading := fs.peekGate(path); !isDownloading {
		n, err := fs.cacheStore.ReadComplete(path, cacheEntryMetadata(path, size, node.info.modTime), buff, ofst)
		if err != nil {
			fs.log.Debug("RemoteFs: Read: cache.Read error: %v", err)
		}
		if n > 0 {
			handle.incrementRead(int64(n))
			fs.log.Trace("RemoteFs: Read: readAt: path=%v, ofst=%d, read %d bytes from cache", path, ofst, n)
			return n
		}
	}

	// At this point, the read request could not be satisfied from the working copy
	// or the disk cache, so read from the remote API.
	endOffset := ofst + int64(len(buff))
	readyGate, exists := fs.findOrCreateGate(path)
	readyGate.Add()
	defer fs.releaseGateWaiter(path, readyGate)
	if !exists {
		// start the download in a new goroutine, which will populate the disk cache
		go fs.fillCache(context.Background(), path, node.downloadUri, cacheEntryMetadata(path, size, node.info.modTime), readyGate, fh, true)
	}

	// wait for the requested range to be available in the cache
	if err := readyGate.WaitFor(endOffset); err != nil {
		// adjust on EOF: serve whatever is available
		if err != io.EOF {
			if errc := fs.handleError(path, err); errc != 0 {
				return errc
			}
		}
		// err is EOF: serve whatever is available
		avail := readyGate.Available()
		if avail <= ofst {
			return 0
		} // nothing available
		endOffset = avail
	}
	// now read from cache
	want := min(endOffset-ofst, int64(len(buff)))
	n, err := fs.cacheStore.ReadPartial(path, buff[:want], ofst)
	if err != nil {
		fs.log.Debug("RemoteFs: Read: diskCache.Read error after WaitFor: %v", err)
		return -fuse.EAGAIN
	}
	fs.log.Trace("RemoteFs: Read: ok path=%v ofst=%d read=%d", path, ofst, n)
	handle.incrementRead(int64(n))
	return n
}

func (fs *RemoteFs) Write(path string, buff []byte, ofst int64, fh uint64) (n int) {
	_, node, ok := fs.vfs.handles.Lookup(fh)
	if !ok {
		fs.log.Debug("RemoteFs: Write: file handle %v not found for path %v", fh, path)
		return -fuse.EBADF
	}

	if err := node.poisonedWriteSessionErr(); err != nil {
		fs.log.Debug("RemoteFs: Write: poisoned write session for path=%v: %s", path, uploadLogMessage(err))
		return -fuse.EIO
	}

	session, created, err := node.beginWriteSessionMutation(path)
	if err != nil {
		fs.log.Error("RemoteFs: Write: failed to create write session for %v: %v", path, err)
		fs.logWriteSessionMilestone(path, "write_session_create_failed", fh, nil, "err=%q", err.Error())
		return -fuse.EIO
	}
	defer session.endMutation()
	path = session.snapshot().path
	if created {
		node.setPendingVisible()
		fs.logWriteSessionMilestone(path, "working_copy_created", fh, session, "offset=%d bytes=%d", ofst, len(buff))
	}
	session.addHandle(fh)

	if err := fs.ensureWriteSessionBaseline(path, node, session, false, fh); err != nil {
		fs.logWriteSessionMilestone(path, "baseline_hydration_failed", fh, session, "err=%q", err.Error())
		if errc := fs.handleUploadSessionError(path, err, session); errc != 0 {
			return errc
		}
		return -fuse.EIO
	}

	written, err := fs.writeToWorkingCopy(session, buff, ofst)
	if err != nil {
		fs.log.Error("RemoteFs: Write: working copy write failed for %v: %v", path, err)
		fs.logWriteSessionMilestone(path, "working_copy_write_failed", fh, session, "offset=%d bytes=%d err=%q", ofst, len(buff), err.Error())
		return -fuse.EIO
	}

	if written > 0 {
		node.extendTtl()
	}

	return written
}

func (fs *RemoteFs) ensureWriteSessionBaseline(path string, node *fsNode, session *writeSession, truncateToZero bool, fh uint64) error {
	return node.ensureWriteSessionHydrated(func(session *writeSession) error {
		session.mu.Lock()
		file := session.workingCopy
		session.mu.Unlock()
		if file == nil {
			return fmt.Errorf("working copy missing for %s", path)
		}

		if err := file.Truncate(0); err != nil {
			return err
		}
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return err
		}

		if truncateToZero || node.info.size <= 0 {
			session.mu.Lock()
			session.baselineSize = node.info.size
			session.currentSize = 0
			if session.mtime.IsZero() {
				session.mtime = time.Now()
			}
			session.mu.Unlock()
			return nil
		}

		if err := fs.populateWorkingCopyFromRemoteOrCache(path, node, session, fh); err != nil {
			return err
		}
		return nil
	})
}

func (fs *RemoteFs) populateWorkingCopyFromRemoteOrCache(path string, node *fsNode, session *writeSession, fh uint64) error {
	if err := fs.ensureFullyCached(path, node.downloadUri, node.info.size, fh); err != nil {
		if files_sdk.IsNotExist(err) {
			node.markDeleted()
			session.mu.Lock()
			session.baselineSize = 0
			session.currentSize = 0
			session.mu.Unlock()
			return nil
		}
		return err
	}

	buf := make([]byte, cacheWriteSize)
	var ofst int64
	for {
		n, err := fs.cacheStore.Read(path, buf, ofst)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := session.workingCopy.WriteAt(buf[:n], ofst); err != nil {
			return err
		}
		ofst += int64(n)
		if ofst >= node.info.size {
			break
		}
	}

	session.mu.Lock()
	session.baselineSize = node.info.size
	session.currentSize = node.info.size
	if session.mtime.IsZero() {
		session.mtime = node.info.modTime
	}
	session.mu.Unlock()
	return nil
}

func (fs *RemoteFs) writeToWorkingCopy(session *writeSession, buff []byte, ofst int64) (int, error) {
	session.mu.Lock()
	defer session.mu.Unlock()

	if session.lastUploadErr != nil {
		return 0, session.lastUploadErr
	}
	if session.workingCopy == nil {
		return 0, io.ErrClosedPipe
	}

	n, err := session.workingCopy.WriteAt(buff, ofst)
	if err != nil {
		return n, err
	}

	end := ofst + int64(n)
	if end > session.currentSize {
		session.currentSize = end
	}
	session.dirty = true
	return n, nil
}

func (fs *RemoteFs) truncateWorkingCopy(session *writeSession, size int64) error {
	session.mu.Lock()
	defer session.mu.Unlock()

	if session.lastUploadErr != nil {
		return session.lastUploadErr
	}
	if session.workingCopy == nil {
		return io.ErrClosedPipe
	}
	if err := session.workingCopy.Truncate(size); err != nil {
		return err
	}
	session.currentSize = size
	session.dirty = true
	return nil
}

func (fs *RemoteFs) refreshReadCacheFromWorkingCopy(path string, session *writeSession, size int64, mtime time.Time) error {
	if fs.cacheStore == nil {
		return nil
	}
	if session.workingCopyPath == "" {
		return nil
	}

	cacheFile, err := os.Open(session.workingCopyPath)
	if err != nil {
		return err
	}
	defer cacheFile.Close()

	return fs.refreshReadCacheFromReader(path, cacheFile, size, mtime)
}

func (fs *RemoteFs) refreshReadCacheFromFile(path string, filePath string, size int64, mtime time.Time) error {
	if fs.cacheStore == nil {
		return nil
	}

	cacheFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer cacheFile.Close()

	return fs.refreshReadCacheFromReader(path, cacheFile, size, mtime)
}

func (fs *RemoteFs) refreshReadCacheFromReader(path string, reader io.Reader, size int64, mtime time.Time) error {
	_ = fs.cacheStore.DeletePartial(path)
	fs.cacheStore.PinPartial(path)
	defer func() {
		fs.cacheStore.UnpinPartial(path)
		_ = fs.cacheStore.DeletePartial(path)
	}()

	if _, err := fs.writePartialCacheFromReader(path, reader, size); err != nil {
		return err
	}
	return fs.commitCacheEntryFromPartial(path, path, cacheEntryMetadata(path, size, mtime), false)
}

func (fs *RemoteFs) writePartialCacheFromReader(path string, reader io.Reader, expectedSize int64) (int64, error) {
	buf := make([]byte, cacheWriteSize)
	var ofst int64
	for {
		n, readErr := reader.Read(buf)
		if n > 0 {
			if _, err := fs.cacheStore.WritePartial(path, buf[:n], ofst); err != nil {
				return ofst, err
			}
			ofst += int64(n)
		}
		if readErr == nil {
			continue
		}
		if errors.Is(readErr, io.EOF) {
			if expectedSize >= 0 && ofst != expectedSize {
				return ofst, io.ErrUnexpectedEOF
			}
			return ofst, nil
		}
		return ofst, readErr
	}
}

func (fs *RemoteFs) flushWriteSession(path string, node *fsNode, fh uint64) (errc int) {
	session := node.getWriteSession()
	if session == nil {
		fs.logWriteSessionMilestone(path, "flush_no_session", fh, nil, "")
		return 0
	}

	session.mu.Lock()
	for session.renaming && session.dirty {
		session.cond.Wait()
		path = session.path
	}
	if session.lastUploadErr != nil {
		err := session.lastUploadErr
		session.mu.Unlock()
		fs.log.Debug("RemoteFs: flushWriteSession: poisoned session for %v: %s", path, uploadLogMessage(err))
		fs.logWriteSessionMilestone(path, "flush_poisoned", fh, session, "err=%q", uploadLogMessage(err))
		return -fuse.EIO
	}
	if !session.dirty && !session.uploading && !session.finalizing {
		session.mu.Unlock()
		fs.logWriteSessionMilestone(path, "flush_clean", fh, session, "")
		return 0
	}
	if session.uploading || session.finalizing {
		session.mu.Unlock()
		fs.logWriteSessionMilestone(path, "flush_wait_existing_upload", fh, session, "")
		if err := node.waitForUploadWithProgressTimeout(fsyncTimeout); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				fs.logWriteSessionMilestone(path, "flush_wait_timeout", fh, session, "errc=%d errno=%s", -fuse.ETIMEDOUT, formatFuseErrno(-fuse.ETIMEDOUT))
				return -fuse.ETIMEDOUT
			}
			fs.logWriteSessionMilestone(path, "flush_wait_failed", fh, session, "err=%q errc=%d errno=%s", err.Error(), -fuse.EIO, formatFuseErrno(-fuse.EIO))
			return -fuse.EIO
		}
		if err := node.poisonedWriteSessionErr(); err != nil {
			fs.logWriteSessionMilestone(path, "flush_wait_poisoned", fh, session, "err=%q errc=%d errno=%s", uploadLogMessage(err), -fuse.EIO, formatFuseErrno(-fuse.EIO))
			return -fuse.EIO
		}
		fs.logWriteSessionMilestone(path, "flush_wait_completed", fh, session, "")
		return 0
	}
	session.finalizing = true
	session.uploading = true
	session.mu.Unlock()
	fs.logWriteSessionMilestone(path, "flush_start_finalize", fh, session, "")

	go fs.finalizeUploadFromWorkingCopy(path, node, session, fh)

	if err := node.waitForUploadWithProgressTimeout(fsyncTimeout); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fs.logWriteSessionMilestone(path, "flush_finalize_timeout", fh, session, "errc=%d errno=%s", -fuse.ETIMEDOUT, formatFuseErrno(-fuse.ETIMEDOUT))
			return -fuse.ETIMEDOUT
		}
		fs.logWriteSessionMilestone(path, "flush_finalize_wait_failed", fh, session, "err=%q errc=%d errno=%s", err.Error(), -fuse.EIO, formatFuseErrno(-fuse.EIO))
		return -fuse.EIO
	}
	if err := node.poisonedWriteSessionErr(); err != nil {
		fs.logWriteSessionMilestone(path, "flush_finalize_poisoned", fh, session, "err=%q errc=%d errno=%s", uploadLogMessage(err), -fuse.EIO, formatFuseErrno(-fuse.EIO))
		return -fuse.EIO
	}
	fs.logWriteSessionMilestone(path, "flush_completed", fh, session, "")
	return 0
}

func (fs *RemoteFs) finalizeUploadFromWorkingCopy(path string, node *fsNode, session *writeSession, fh uint64) {
	localPath, remotePath := fs.paths(path)
	fs.log.Info("Starting upload from working copy: %v (%v)", remotePath, localPath)
	fs.logWriteSessionMilestone(path, "upload_finalize_started", fh, session, "")
	sessionSnapshot := session.snapshot()
	transfer := fs.newTransferReporter(events.TransferDirectionUpload, path, sessionSnapshot.currentSize)
	transfer.Queued()

	reader, err := os.Open(session.workingCopyPath)
	if err != nil {
		fs.logWriteSessionMilestone(path, "upload_open_working_copy_failed", fh, session, "err=%q", err.Error())
		transfer.Error(err, transferredBytesUnchanged)
		node.clearPendingVisible()
		node.writeSessionFinishUpload(session.snapshot().currentSize, err)
		return
	}
	defer reader.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if _, ok := node.writeSessionStartUpload(cancel, transfer); !ok {
		err := fmt.Errorf("write session missing for %s", path)
		fs.logWriteSessionMilestone(path, "upload_start_missing_session", fh, session, "err=%q", err.Error())
		transfer.Error(err, transferredBytesUnchanged)
		node.clearPendingVisible()
		return
	}
	fs.logWriteSessionMilestone(path, "upload_started", fh, session, "size=%d", sessionSnapshot.currentSize)

	session.mu.Lock()
	mtime := session.mtime
	if session.dirty && !session.mtimeExplicit {
		mtime = time.Now()
	}
	session.mu.Unlock()

	uploaded, err := fs.uploadWorkingCopyWithSDK(ctx, node, path, reader, mtime, fh)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			fs.log.Info("Upload canceled from working copy: %v (%v)", remotePath, localPath)
		} else if !files_sdk.IsNotExist(err) {
			fs.logUploadFailure("Upload failed from working copy", remotePath, localPath, err)
		}
		class, mappedErrc := classifyMountError(err)
		fs.log.Debug(
			"mount_upload_error source=backend_upload path=%q remote_path=%q local_path=%q class=%s mapped_errc=%d mapped_errno=%s callback_errc=%d callback_errno=%s message=%q",
			path,
			remotePath,
			localPath,
			class,
			mappedErrc,
			formatFuseErrno(mappedErrc),
			-fuse.EIO,
			formatFuseErrno(-fuse.EIO),
			uploadLogMessage(err),
		)
		fs.logWriteSessionMilestone(path, "upload_failed", fh, session, "err=%q class=%s mapped_errc=%d mapped_errno=%s message=%q", err.Error(), class, mappedErrc, formatFuseErrno(mappedErrc), uploadLogMessage(err))
		transfer.Error(err, transferredBytesUnchanged)
		node.clearPendingVisible()
		_ = node.writeSessionFinishUpload(session.snapshot().currentSize, err)
		return
	}
	fs.logWriteSessionMilestone(path, "upload_succeeded", fh, session, "size=%d", uploaded.size)
	node.markMaterialized()

	if err := fs.refreshReadCacheFromWorkingCopy(path, session, uploaded.size, uploaded.modTime); err != nil {
		// The remote upload has already succeeded, so a cache refresh failure should only
		// invalidate the local cache. Do not poison or finish the write session with this
		// error, or the caller would see a successful save as a failed upload.
		fs.log.Error("Error refreshing cache from working copy: %v (%v): %v", remotePath, localPath, err)
		fs.logWriteSessionMilestone(path, "upload_cache_refresh_failed", fh, session, "err=%q", err.Error())
		_ = fs.cacheStore.Delete(path)
	}

	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         uploaded.size,
		modTime:      uploaded.modTime,
		creationTime: node.info.creationTime,
		uid:          node.info.uid,
		gid:          node.info.gid,
	})
	session.mu.Lock()
	session.mtime = uploaded.modTime
	session.mtimeExplicit = false
	session.mu.Unlock()
	fs.log.Info("Upload completed from working copy: %v (%v)", remotePath, localPath)
	transfer.Complete(sessionSnapshot.currentSize)
	_ = node.writeSessionFinishUpload(uploaded.size, nil)
	fs.logWriteSessionMilestone(path, "upload_finalize_completed", fh, session, "size=%d", uploaded.size)
}

func (fs *RemoteFs) uploadWorkingCopyWithSDK(ctx context.Context, node *fsNode, path string, reader uploadWorkingCopyReader, mtime time.Time, fh uint64) (uploadedFileMetadata, error) {
	if fs.uploadWorkingCopy != nil {
		return fs.uploadWorkingCopy(ctx, node, path, reader, mtime, fh)
	}

	localPath, remotePath := fs.paths(path)
	fileInfo, err := reader.Stat()
	if err != nil {
		return uploadedFileMetadata{}, err
	}
	opts := []file.UploadOption{
		file.UploadWithContext(ctx),
		file.UploadWithDestinationPath(remotePath),
		file.UploadWithReaderAt(reader),
		file.UploadWithSize(fileInfo.Size()),
		file.UploadWithProvidedMtime(mtime),
		file.UploadWithProgress(fs.uploadProgressFunc(node)),
		file.WithUploadStartedCallback(func(part files_sdk.FileUploadPart) {
			fs.log.Debug("RemoteFs: Uploading part number %d, of: %v, ref: '%v'", part.PartNumber, remotePath, part.Ref)
			node.captureRef(part.Ref)
		}),
		file.WithUploadRenamedCallback(func() (string, string) {
			finalRemotePath, ref := fs.finalizeUploadPathAndRef(node)
			if remotePath != finalRemotePath {
				fs.log.Debug("RemoteFs: finalizeUploadFromWorkingCopy: in progress upload renamed from: %v to %v", remotePath, finalRemotePath)
			}
			return finalRemotePath, ref
		}),
	}
	if fs.writeConcurrency != 0 {
		opts = append(opts, file.UploadWithManager(manager.ConcurrencyManager{}.New(fs.writeConcurrency)))
	}
	u, err := fs.backend.uploadWithResume(opts...)
	if err != nil {
		fs.log.Debug("RemoteFs: uploadWorkingCopyWithSDK failed for %v (%v): %s", remotePath, localPath, uploadLogMessage(err))
		return uploadedFileMetadata{}, err
	}
	return fs.uploadedFileMetadata(path, u.File, fileInfo.Size(), mtime), nil
}

func (fs *RemoteFs) finalizeUploadPathAndRef(node *fsNode) (string, string) {
	path, ref := node.captureUploadPathAndRef()
	return fs.remotePath(path), ref
}

func (fs *RemoteFs) logUploadFailure(prefix string, remotePath string, localPath string, err error) {
	fs.log.Error("%s: %v (%v): %s", prefix, remotePath, localPath, uploadLogMessage(err))
}

func uploadLogMessage(err error) string {
	if classified, ok := lib.ClassifyS3Error(err); ok {
		return classified.Message
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return "Network timeout"
	}
	if os.IsTimeout(err) {
		return "Network timeout"
	}
	if errors.Is(err, syscall.ECONNREFUSED) {
		return "Connection refused"
	}
	if errors.Is(err, syscall.ECONNRESET) {
		return "Connection reset by peer"
	}
	if errors.Is(err, syscall.EPIPE) {
		return "Client disconnected during transfer"
	}

	return "Error returned by the remote service"
}

func (fs *RemoteFs) Release(path string, fh uint64) (errc int) {
	fs.log.Trace("RemoteFs: Release: path=%v, fh=%v", path, fh)
	handle, node, ok := fs.vfs.handles.Lookup(fh)

	// Remove the handle from the set of open handles in all cases,
	// the host FS is finished with the handle, so it will not be used again.
	defer fs.vfs.handles.Release(fh)

	if !ok {
		// This is an unexpected condition. Why is the OS calling to release
		// a file handle that was never opened? There's no file handle to release,
		// so log the error and try to unlock the path if it was locked.
		fs.log.Debug("RemoteFs: Release: file handle not found for path: %v, fh: %v", path, fh)

		// unlock is a no-op if the path/handle combo doesn't match an existing lock
		if errc = fs.unlock(path, fh); errc != 0 {
			return errc
		}
		return errc
	}

	// Unpin the file from cache when the handle is closed, allowing it to be evicted if needed.
	// This must be done for ALL handle types (read-only and write) since Pin() is called
	// in both Open() and Create() for all handle types.
	defer func() {
		if fs.cacheStore != nil {
			fs.cacheStore.Unpin(path)
		}
	}()

	// If the handle is a read only operation, there's nothing left to do.
	if handle.IsReadOnly() {
		fs.log.Trace("RemoteFs: Release: closed handle for path=%v, fh=%v", path, fh)
		return errc
	}

	if session := node.getWriteSession(); session != nil {
		iOwn := session.hasHandle(fh)
		defer node.expireInfo()
		if iOwn {
			remaining := session.removeHandle(fh)
			if errc := fs.flushWriteSession(path, node, fh); errc != 0 {
				return errc
			}
			if remaining == 0 {
				if err := node.clearWriteSession(); err != nil {
					fs.log.Debug("RemoteFs: Release: failed clearing write session for %v: %v", path, err)
				}
			}
		}
		return fs.unlock(path, fh)
	}

	return fs.unlock(path, fh)
}

func (fs *RemoteFs) Opendir(path string) (errc int, fh uint64) {
	node := fs.vfs.getOrCreate(path, nodeTypeDir)
	fh, _ = fs.vfs.handles.Open(node, ff.NewFuseFlags(fuse.O_RDONLY))
	fs.log.Trace("RemoteFs: Opendir: path=%v, fh=%v", path, fh)
	return errc, fh
}

func (fs *RemoteFs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {

	readdirStart := time.Now()
	defer func() {
		if elapsed := time.Since(readdirStart); elapsed > 2*time.Second {
			localPath, remotePath := fs.paths(path)
			fs.log.Debug("RemoteFs: Readdir: slow response (%v): %v (%v)", elapsed, remotePath, localPath)
		}
	}()

	localPath, remotePath := fs.paths(path)

	// This happens a lot, so log at trace level.
	fs.log.Trace("RemoteFs: Readdir: Listing folder: %v (%v)", remotePath, localPath)

	fillNode, _ := fs.vfs.fetch(path)

	// Force a load of the directory entries from the remote to make sure
	// the local vfs representation is up to date.
	start := time.Now()
	if errc = fs.loadDir(fillNode); errc != 0 {
		return errc
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		fs.log.Debug("RemoteFs: loadDir: slow response (%v): %v (%v)", elapsed, remotePath, localPath)
	}

	fill(".", nil, 0)
	fill("..", nil, 0)

	// construct a list of child entries for the current directory
	entries := make([]string, 0, len(fillNode.childPaths))
	for childPath := range fillNode.childPaths {
		entries = append(entries, childPath)
	}

	// make sure to append any open paths that are not already in the entries list
	// this ensures that uploads in progress are visible in the directory listing
	handles := fs.vfs.handles.OpenHandles()
	for _, handle := range handles {
		openNode := handle.node
		// skip the remote root node, the node that represents the current directory being read, and nodes
		// that are rooted in the local fs (e.g. temporary files created by the OS)
		if openNode.path == fs.root || openNode.path == path || strings.HasPrefix(openNode.path, fs.localFsRoot) {
			continue
		}
		if !slices.Contains(entries, openNode.path) && path == path_lib.Dir(openNode.path) {
			fs.log.Debug("RemoteFs: Readdir: Child entries %v: for path %s, does not include open handle: %v, adding %v", entries, path, handle, openNode.path)
			entries = append(entries, openNode.path)
		}
	}

	// include pending-visible nodes so that files being uploaded remain visible
	// even after the handle is released but before the remote listing catches up
	for p := range fs.vfs.pendingVisibleChildPaths(path) {
		if !slices.Contains(entries, p) {
			entries = append(entries, p)
		}
	}

	// sort the entries in order to provide a consistent sort order when calling fill
	slices.Sort(entries)
	for _, entryPath := range entries {
		if entryNode, ok := fs.vfs.fetch(entryPath); ok {
			fs.log.Trace("RemoteFs: Readdir: Calling fill for entry: %v (%v)", entryPath, entryPath)
			fill(path_lib.Base(entryPath), getStat(entryNode.info, nil, fs.vfs.uid, fs.vfs.gid), 0)
		} else {
			// This can happen if the OS has opened multiple handles for a single node, and Unlink
			// is called on a path before all the handles are Released. In this case, the node will
			// be removed from the vfs, but the handle will still exist in the handles map, and the
			// Readdir call will attempt to fill the entry for the node that no longer exists in the vfs.
			fs.log.Debug("RemoteFs: Readdir: entry node not found: %v (%v)", path_lib.Base(entryPath), entryPath)
		}
	}

	return errc
}

func (fs *RemoteFs) Releasedir(path string, fh uint64) (errc int) {
	fs.log.Trace("RemoteFs: Releasedir: path=%v, fh=%v", path, fh)
	fs.vfs.handles.Release(fh)
	return errc
}

// Chmod changes the permission bits of a file.
// Files.com does not support POSIX permissions, but certain operations may fail
// if calling Chmod returns an error, so this implementation is a no-op that returns success.
func (fs *RemoteFs) Chmod(path string, mode uint32) int {
	fs.log.Debug("RemoteFs: Chmod: path=%v, mode=%o", path, mode)
	return 0
}

// Fsync attempts to synchronize file contents.
// If an upload is already finalizing in the background, wait for it to complete.
// If the writer is still open, do not finalize it here; some applications issue
// repeated Fsync calls while continuing to write through the same handle.
func (fs *RemoteFs) Fsync(path string, datasync bool, fh uint64) (errc int) {
	fs.log.Debug("RemoteFs: Fsync: path=%v, datasync=%v, fh=%v", path, datasync, fh)

	_, node, ok := fs.vfs.handles.Lookup(fh)
	if !ok {
		fs.log.Debug("RemoteFs: Fsync: file handle not found for path: %v, fh: %v", path, fh)
		return -fuse.EBADF
	}

	if session := node.getWriteSession(); session != nil {
		if !session.hasHandle(fh) {
			return 0
		}
		if err := node.poisonedWriteSessionErr(); err != nil {
			fs.log.Debug("RemoteFs: Fsync: poisoned write session for path=%v: %s", path, uploadLogMessage(err))
			return -fuse.EIO
		}
		return fs.flushWriteSession(path, node, fh)
	}

	if fh != ^uint64(0) {
		fs.log.Debug("RemoteFs: Fsync: handle does not own writer for path=%v, fh=%v", path, fh)
		return 0
	}
	return 0
}

// copyWriterToCache copies the writer's temp file content to the cache.
// This is called before closing the writer to preserve the uploaded content for subsequent writes.
func (fs *RemoteFs) uploadProgressFunc(node *fsNode) func(int64) {
	return fs.uploadProgressFuncWithTransfer(node, nil)
}

func (fs *RemoteFs) uploadProgressFuncWithTransfer(node *fsNode, transfer *transferReporter) func(int64) {
	return func(delta int64) {
		// Extend the node's TTL and keep track of bytes written for logging/sweeping.
		node.extendTtl()
		if sessionTransfer := node.writeSessionRecordProgress(delta); sessionTransfer != nil {
			sessionTransfer.Progress(delta)
			return
		}
		if transfer != nil {
			transfer.Progress(delta)
		}
	}
}

// this is a convenience method for uploading a file from the local file system to the remote backend
// for use by the Rename operation when moving a file from the LocalFs to the RemoteFs.
func (fs *RemoteFs) uploadFile(src, dst string) error {
	fs.log.Debug("Uploading file: %v to %v", src, dst)

	uploadFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer uploadFile.Close()

	// Get file info for size and mtime
	fileInfo, err := uploadFile.Stat()
	if err != nil {
		return err
	}

	localPath, remotePath := fs.paths(dst)
	transfer := fs.newTransferReporter(events.TransferDirectionUpload, dst, fileInfo.Size())
	transfer.Queued()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node := fs.vfs.getOrCreate(dst, nodeTypeFile)
	if fs.providerBackend != nil {
		var uploaded ProviderEntry
		err := fs.ops.WithLimit(ctx, lim.FuseOpUpload, func(ctx context.Context) error {
			var writeErr error
			uploadReader := &progressReader{
				reader:   uploadFile,
				progress: fs.uploadProgressFuncWithTransfer(node, transfer),
			}
			uploaded, writeErr = fs.providerBackend.Write(ctx, remotePath, uploadReader, fileInfo.Size(), fileInfo.ModTime())
			return writeErr
		})
		if err != nil {
			if !errors.Is(err, context.Canceled) && !files_sdk.IsNotExist(err) {
				fs.logUploadFailure("Error uploading file during rename", remotePath, localPath, err)
			}
			transfer.Error(err, transferredBytesUnchanged)
			return err
		}

		modTime := uploaded.ModTime
		if modTime.IsZero() {
			modTime = fileInfo.ModTime()
		}
		// Provider writes consume the source reader directly. Cache population is
		// intentionally serialized after Write for now; SDK uploads can do this in
		// parallel because they use the SDK upload manager and a separate cache reader.
		if err := fs.refreshReadCacheFromFile(dst, src, fileInfo.Size(), modTime); err != nil {
			_ = fs.cacheStore.Delete(dst)
			fs.log.Error("Error refreshing cache during provider rename upload; invalidating cache entry: %v (%v): %v", remotePath, localPath, err)
		}

		size := fs.providerUploadedSize(remotePath, uploaded, fileInfo.Size())
		fs.log.Info("Upload completed during rename: %v (%v).", remotePath, localPath)
		transfer.Complete(size)
		node.updateInfo(fsNodeInfo{
			nodeType:     nodeTypeFile,
			size:         size,
			modTime:      modTime,
			creationTime: node.info.creationTime,
		})
		return nil
	}

	// Open an independent file handle so cache population can read sequentially
	// in parallel while the SDK upload uses the seekable reader above.
	cacheFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer cacheFile.Close()

	uploadResultCh := make(chan uploadFileResult, 1)
	cacheErrCh := make(chan error, 1)
	_ = fs.cacheStore.DeletePartial(dst)
	fs.cacheStore.PinPartial(dst)
	defer func() {
		fs.cacheStore.UnpinPartial(dst)
		_ = fs.cacheStore.DeletePartial(dst)
	}()

	go func() {
		result := uploadFileResult{
			metadata: uploadedFileMetadata{
				size:    fileInfo.Size(),
				modTime: fileInfo.ModTime(),
			},
		}
		result.err = fs.ops.WithLimit(ctx, lim.FuseOpUpload, func(ctx context.Context) error {
			uploadOpts := []file.UploadOption{
				file.UploadWithContext(ctx),
				file.UploadWithReaderAt(uploadFile),
				file.UploadWithSize(fileInfo.Size()),
				file.UploadWithDestinationPath(remotePath),
				file.UploadWithProvidedMtime(fileInfo.ModTime()),
				file.UploadWithProgress(fs.uploadProgressFuncWithTransfer(node, transfer)),
				file.WithUploadStartedCallback(func(part files_sdk.FileUploadPart) {
					fs.log.Debug("RemoteFs: uploadFile: uploading part number %d, of: %v, ref: '%v'", part.PartNumber, remotePath, part.Ref)
				}),
			}
			if fs.writeConcurrency != 0 {
				uploadOpts = append(uploadOpts, file.UploadWithManager(manager.ConcurrencyManager{}.New(fs.writeConcurrency)))
			}
			u, err := fs.backend.uploadWithResume(uploadOpts...)
			if err != nil {
				return err
			}
			result.metadata = fs.uploadedFileMetadata(dst, u.File, fileInfo.Size(), fileInfo.ModTime())
			return nil
		})
		uploadResultCh <- result
	}()

	go func() {
		_, err := fs.writePartialCacheFromReader(dst, cacheFile, fileInfo.Size())
		cacheErrCh <- err
	}()

	var uploadResult uploadFileResult
	var cacheErr error
	for i := 0; i < 2; i++ {
		select {
		case uploadResult = <-uploadResultCh:
			if uploadResult.err != nil {
				cancel()
			}
		case cacheErr = <-cacheErrCh:
			// The cache goroutine writes to the partial namespace, which is cleaned up by the
			// function defer. A prior committed dst entry remains valid unless the
			// upload succeeds and updates the remote version.
		}
	}

	if cacheErr != nil {
		fs.log.Error("Error populating cache during rename upload; invalidating cache entry: %v (%v): %v", remotePath, localPath, cacheErr)
	}

	if uploadResult.err != nil {
		if !errors.Is(uploadResult.err, context.Canceled) && !files_sdk.IsNotExist(uploadResult.err) {
			fs.logUploadFailure("Error uploading file during rename", remotePath, localPath, uploadResult.err)
		}
		transfer.Error(uploadResult.err, transferredBytesUnchanged)
		return uploadResult.err
	}
	if cacheErr == nil {
		if err := fs.commitCacheEntryFromPartial(dst, dst, cacheEntryMetadata(dst, uploadResult.metadata.size, uploadResult.metadata.modTime), false); err != nil {
			_ = fs.cacheStore.Delete(dst)
			fs.log.Error("Error committing cache during rename upload; invalidating cache entry: %v (%v): %v", remotePath, localPath, err)
		}
	}

	// Update the node's info after the upload succeeds. Cache commit failures are logged above
	// and treated as cache misses because the remote write has already completed.
	fs.log.Info("Upload completed during rename: %v (%v).", remotePath, localPath)
	transfer.Complete(fileInfo.Size())
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         uploadResult.metadata.size,
		modTime:      uploadResult.metadata.modTime,
		creationTime: node.info.creationTime,
	})

	return nil
}

// this is a convenience method for downloading a file from the remote API to the local file system
// for use by the Rename operation when moving a file from the RemoteFs to the LocalFs.
func (fs *RemoteFs) downloadFile(src, dst string, eventLocalPath string) error {
	fs.log.Debug("RemoteFs: Downloading file: %v to %v", src, dst)
	transfer := fs.newTransferReporterForPaths(events.TransferDirectionDownload, eventLocalPath, fs.remotePath(src), 0)
	transfer.Queued()
	var downloaded files_sdk.File
	err := fs.ops.WithLimit(context.Background(), lim.FuseOpDownload, func(ctx context.Context) error {
		var err error
		downloaded, err = fs.backend.downloadToFile(files_sdk.FileDownloadParams{Path: src}, dst)
		return err
	})
	if err != nil {
		transfer.Error(err, transferredBytesUnchanged)
		return err
	}

	size := downloaded.Size
	if size == 0 {
		if info, statErr := os.Stat(dst); statErr == nil {
			size = info.Size()
		}
	}
	transfer.Complete(size)
	return nil
}

func (fs *RemoteFs) findOrCreateGate(path string) (*cache.ReadyGate, bool) {
	fs.gatePathMutexes.Lock(path)
	defer fs.gatePathMutexes.Unlock(path)
	fs.gatesMu.Lock()
	defer fs.gatesMu.Unlock()
	if fs.readyGates == nil {
		fs.readyGates = map[string]*cache.ReadyGate{}
	}
	if s, ok := fs.readyGates[path]; ok {
		return s, true
	}
	s := cache.NewReadyGate()
	fs.readyGates[path] = s
	return s, false
}

func (fs *RemoteFs) lockGatePaths(paths ...string) func() {
	slices.Sort(paths)
	paths = slices.Compact(paths)
	for _, path := range paths {
		fs.gatePathMutexes.Lock(path)
	}
	return func() {
		for i := len(paths) - 1; i >= 0; i-- {
			fs.gatePathMutexes.Unlock(paths[i])
		}
	}
}

func (fs *RemoteFs) takeGatesForPaths(paths ...string) []*cache.ReadyGate {
	fs.gatesMu.Lock()
	defer fs.gatesMu.Unlock()
	var gates []*cache.ReadyGate
	for _, path := range paths {
		if readyGate := fs.readyGates[path]; readyGate != nil {
			gates = append(gates, readyGate)
			delete(fs.readyGates, path)
		}
	}
	return gates
}

func (fs *RemoteFs) removeGate(path string, s *cache.ReadyGate) {
	fs.gatesMu.Lock()
	if cur, ok := fs.readyGates[path]; ok && cur == s {
		delete(fs.readyGates, path)
	}
	fs.gatesMu.Unlock()
}

func (fs *RemoteFs) releaseGateWaiter(path string, readyGate *cache.ReadyGate) {
	if readyGate.Done() {
		fs.removeGate(path, readyGate)
		readyGate.Cleanup()
	}
}

func (fs *RemoteFs) commitCacheEntryFromPartial(srcPath string, dstPath string, meta cache.EntryMetadata, deleteSource bool) error {
	_ = fs.cacheStore.Delete(dstPath)
	fs.cacheStore.Pin(dstPath)
	committed := false
	defer func() {
		fs.cacheStore.Unpin(dstPath)
		if !committed {
			_ = fs.cacheStore.Delete(dstPath)
		}
	}()

	buf := make([]byte, cacheWriteSize)
	var copied int64
	for copied < meta.Size {
		want := min(int64(len(buf)), meta.Size-copied)
		n, err := fs.cacheStore.ReadPartial(srcPath, buf[:want], copied)
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrUnexpectedEOF
		}
		written, err := fs.cacheStore.Write(dstPath, buf[:n], copied)
		if err != nil {
			return err
		}
		if written != n {
			return io.ErrShortWrite
		}
		copied += int64(n)
	}

	if meta.Size == 0 {
		if _, err := fs.cacheStore.Write(dstPath, nil, 0); err != nil {
			return err
		}
	}

	if err := fs.cacheStore.Commit(dstPath, meta); err != nil {
		return err
	}
	committed = true
	if deleteSource {
		_ = fs.cacheStore.DeletePartial(srcPath)
	}
	return nil
}

// ensureFullyCached ensures the remote file at path is fully downloaded to the cache.
// If a download is already in progress it joins that download; otherwise it starts one.
// It blocks until all size bytes are available in the cache.
//
// A fast-path probe checks whether the last byte of a complete, metadata-matching cache
// entry is already readable, which means no network round-trip is needed.
func (fs *RemoteFs) ensureFullyCached(path, uri string, size int64, fh uint64) error {
	if size <= 0 {
		return nil
	}
	if node, ok := fs.vfs.fetch(path); ok {
		meta := cacheEntryMetadata(path, size, node.info.modTime)
		probe := [1]byte{}
		if n, _ := fs.cacheStore.ReadComplete(path, meta, probe[:], size-1); n == 1 {
			return nil
		}
	}
	// Slow path: join or start a download and wait for the full file.
	readyGate, exists := fs.findOrCreateGate(path)
	readyGate.Add()
	defer fs.releaseGateWaiter(path, readyGate)
	if !exists {
		node, ok := fs.vfs.fetch(path)
		if !ok {
			err := fmt.Errorf("ensureFullyCached: vfs node missing for %s", path)
			readyGate.Finish(err, 0)
			fs.removeGate(path, readyGate)
			return err
		}
		go fs.fillCache(context.Background(), path, uri, cacheEntryMetadata(path, size, node.info.modTime), readyGate, fh, false)
	}
	if err := readyGate.WaitFor(size); err != nil {
		return err
	}
	return nil
}

// peekGate returns the existing ready gate for path if one is present, without creating one.
// It is used to check whether a download is currently in progress for a given path.
func (fs *RemoteFs) peekGate(path string) (*cache.ReadyGate, bool) {
	fs.gatesMu.Lock()
	defer fs.gatesMu.Unlock()
	if fs.readyGates == nil {
		return nil, false
	}
	s, ok := fs.readyGates[path]
	return s, ok
}

func (fs *RemoteFs) fillCache(ctx context.Context, path string, uri string, meta cache.EntryMetadata, readyGate *cache.ReadyGate, fh uint64, cancelOnActiveWrite bool) {
	ctx, cancel := context.WithCancel(ctx)
	readyGate.SetTotal(meta.Size)
	readyGate.SetCancel(cancel)
	readyGate.SetCleanup(func() {
		fs.cacheStore.UnpinPartial(path)
		_ = fs.cacheStore.DeletePartial(path)
	})
	defer func() {
		cancel()
		fs.removeGate(path, readyGate)
		if readyGate.Drained() {
			readyGate.Cleanup()
		}
	}()
	_ = fs.cacheStore.DeletePartial(path)
	fs.cacheStore.PinPartial(path)
	transfer := fs.newTransferReporter(events.TransferDirectionDownload, path, meta.Size)
	transfer.Queued()

	var f files_sdk.File
	var err error
	err = fs.ops.WithLimit(ctx, lim.FuseOpDownload, func(ctx context.Context) error {
		f, err = fs.backend.download(
			files_sdk.FileDownloadParams{File: files_sdk.File{Path: fs.remotePath(path), DownloadUri: uri}},
			files_sdk.WithContext(ctx),
			files_sdk.ResponseOption(func(resp *http.Response) error {
				if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
					return files_sdk.APIError()(resp)
				}
				defer resp.Body.Close()

				// Get buffer from pool
				buf := fs.bufferPool.Get()
				defer fs.bufferPool.Put(buf)

				// while downloading a file from the remote API, write data to the disk cache in chunks and update
				// the ready gate every cacheWriteSize bytes to signal that data is available for reading.
				var off int64 = 0
				for {
					if cancelOnActiveWrite {
						if node, ok := fs.vfs.fetch(path); ok && node.hasHydratedWriteSession() {
							// A public read should not keep serving old remote bytes after
							// a local write has its baseline. Before hydration completes,
							// this same download may be supplying that baseline.
							readyGate.Finish(context.Canceled, off)
							transfer.Error(context.Canceled, off)
							return nil
						}
					}

					nr, er := resp.Body.Read(buf)
					if nr > 0 {
						// TODO: consider altering Write to keep data in memory and periodically flush to disk
						// to reduce the number of disk writes. This would require more memory usage, but would
						// improve read and write performance by avoiding constantly opening/closing the file.
						written, err := fs.cacheStore.WritePartial(path, buf[:nr], off)
						if err != nil || written != nr {
							// there was an error writing to the disk cache, or not all bytes that were read from the
							// remote API were written to the disk cache.
							cacheErr := fmt.Errorf("error writing to disk cache for %v: %v", path, err)
							readyGate.Finish(cacheErr, off)
							transfer.Error(cacheErr, off)
							return cacheErr
						}
						off += int64(written)
						readyGate.SetAvailable(off)
						transfer.Progress(int64(written))
					}
					if er != nil {
						if er == io.EOF {
							if off != meta.Size {
								err := io.ErrUnexpectedEOF
								readyGate.Finish(err, off)
								transfer.Error(err, off)
								return err
							}
							if err := fs.commitCacheEntryFromPartial(path, path, meta, false); err != nil {
								readyGate.Finish(err, off)
								transfer.Error(err, off)
								return err
							}
							readyGate.Finish(nil, off)
							transfer.Complete(off)
							return nil
						}
						readyGate.Finish(er, off)
						transfer.Error(er, off)
						return er
					}
					// TODO: consider canceling the download if there are no active readers/waiters
					// after a certain period of time
				}
			}),
		)
		return err
	})

	if err != nil {
		readyGate.Finish(err, -1)
		transfer.Error(err, transferredBytesUnchanged)
		return
	}
	if f.Size > 0 {
		localPath, remotePath := fs.paths(path)
		fs.log.Info("Download complete: %v (%v), size=%v fh=%v", remotePath, localPath, f.Size, fh)
	}
}

// rename updates local bookkeeping for a path change.
// It renames the node in the VFS, migrates any lock entry, and clears stale URLs.
func (fs *RemoteFs) rename(oldpath, newpath string) {
	// 1) Update the in-memory node map (and parent childPaths) first.
	//    vfs.rename should handle moving the node and fixing parent listings.
	node := fs.vfs.rename(oldpath, newpath)

	// 2) Move any lock entry old -> new.
	fs.lockMapMutex.Lock()
	if li, ok := fs.lockMap[oldpath]; ok {
		delete(fs.lockMap, oldpath)
		fs.lockMap[newpath] = li
	}
	fs.lockMapMutex.Unlock()

	// 3) Clear any cached presigned URL for this node (path changed).
	if node != nil {
		node.clearDownloadURI()
	}
}

func (fs *RemoteFs) lock(node *fsNode, fh uint64) (errc int) {
	if fs.disableLocking {
		fs.logWriteSessionMilestone(node.path, "lock_skipped_disabled", fh, nil, "")
		return errc
	}

	node.lockMutex.Lock()
	defer node.lockMutex.Unlock()

	localPath, remotePath := fs.paths(node.path)
	fs.log.Debug("RemoteFs: lock: file %v (%v) fh=%v", remotePath, localPath, fh)

	fs.lockMapMutex.Lock()
	linfo, ok := fs.lockMap[node.path]
	fs.lockMapMutex.Unlock()
	if ok {
		if fs.currentUserId == linfo.Lock.UserId {
			node.setLockOwner(linfo.Lock.Username)
			fs.log.Debug("RemoteFs: lock: reusing existing same-user lock for %v (%v) fh=%v", remotePath, localPath, fh)
			fs.logWriteSessionMilestone(node.path, "lock_reused_same_user", fh, nil, "owner=%q", linfo.Lock.Username)
			return 0
		}
		node.setLockOwner(linfo.Lock.Username)
		fs.log.Error("File '%v' is already locked by %v:", remotePath, linfo.Lock.Username)
		fs.logWriteSessionMilestone(node.path, "lock_conflict", fh, nil, "owner=%q errc=%d errno=%s", linfo.Lock.Username, -fuse.ENOLCK, formatFuseErrno(-fuse.ENOLCK))
		return -fuse.ENOLCK
	}

	if node.isLocked() {
		fs.log.Error("File is already locked by %v: %v (%v) fh=%v", node.info.lockOwner, remotePath, localPath, fh)
		errc = -fuse.ENOLCK
		fs.logWriteSessionMilestone(node.path, "lock_node_already_locked", fh, nil, "owner=%q errc=%d errno=%s", node.info.lockOwner, errc, formatFuseErrno(errc))
		return errc
	}

	parentPath := path_lib.Dir(node.path)
	if parentPath != node.path {
		if parent, ok := fs.vfs.fetch(parentPath); ok {
			if permissions := parent.getRemotePermissions(); permissions != "" && !parent.isWritable() {
				fs.log.Debug("RemoteFs: lock: skipping lock for %v (%v) because parent directory %v has non-writable remote permissions %q", remotePath, localPath, parentPath, permissions)
				fs.logWriteSessionMilestone(node.path, "lock_skipped_readonly_parent", fh, nil, "parent=%q permissions=%q", parentPath, permissions)
				return 0
			}
		}
	}

	// Make API call without holding lockMapMutex
	var lock files_sdk.Lock
	var err error
	fs.logWriteSessionMilestone(node.path, "lock_create_started", fh, nil, "")
	err = fs.ops.TryWithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		lock, err = fs.backend.createLock(files_sdk.LockCreateParams{
			Path:                 remotePath,
			AllowAccessByAnyUser: lib.Ptr(true),
			Exclusive:            lib.Ptr(true),
			Recursive:            lib.Ptr(false),
			Timeout:              fileLockSeconds,
		})
		return err
	})
	if errors.Is(err, lim.ErrNoSlotsAvailable) {
		fs.logWriteSessionMilestone(node.path, "lock_create_no_slots", fh, nil, "errc=%d errno=%s", -fuse.EAGAIN, formatFuseErrno(-fuse.EAGAIN))
		return -fuse.EAGAIN
	}

	if files_sdk.IsExist(err) {
		// the file is already locked, if it's in the lock map and not owned by this user, return ENOLCK
		fs.lockMapMutex.Lock()
		linfo, ok := fs.lockMap[node.path]
		fs.lockMapMutex.Unlock()

		if ok && fs.currentUserId != linfo.Lock.UserId {
			node.setLockOwner(linfo.Lock.Username)
			fs.log.Error("File '%v' is already locked by %v:", remotePath, linfo.Lock.Username)
			fs.logWriteSessionMilestone(node.path, "lock_backend_conflict", fh, nil, "owner=%q errc=%d errno=%s", linfo.Lock.Username, -fuse.ENOLCK, formatFuseErrno(-fuse.ENOLCK))
			return -fuse.ENOLCK
		}
		if ok && fs.currentUserId == linfo.Lock.UserId {
			// If the lock is already held by the current user, treat it as a success.
			node.setLockOwner(linfo.Lock.Username)
			fs.log.Debug("RemoteFs: lock: File is already locked by current user %v: %v (%v) fh=%v", fs.currentUserId, remotePath, localPath, fh)
			fs.logWriteSessionMilestone(node.path, "lock_backend_same_user", fh, nil, "owner=%q", linfo.Lock.Username)
			return 0
		}
		if node.uploadActive() {
			node.setLockOwner("current-user")
			fs.log.Debug("RemoteFs: lock: treating backend lock conflict as in-flight upload reuse for %v (%v) fh=%v", remotePath, localPath, fh)
			fs.logWriteSessionMilestone(node.path, "lock_backend_conflict_upload_active", fh, nil, "")
			return 0
		}
	}

	if errc = fs.handleError(node.path, err); errc != 0 {
		fs.logWriteSessionMilestone(node.path, "lock_create_failed", fh, nil, "errc=%d errno=%s", errc, formatFuseErrno(errc))
		return errc
	}

	// Store the lock in the map - only hold mutex for this operation
	fs.lockMapMutex.Lock()
	fs.lockMap[node.path] = &lockInfo{Fh: fh, Lock: &lock}
	fs.lockMapMutex.Unlock()
	node.setLockOwner(lock.Username)

	fs.log.Debug("RemoteFs: lock: created owner=%v, path=%v, fh=%v", lock.Username, remotePath, fh)
	fs.logWriteSessionMilestone(node.path, "lock_created", fh, nil, "owner=%q", lock.Username)
	return errc
}

func (fs *RemoteFs) unlock(path string, fh uint64) (errc int) {
	if fs.disableLocking {
		fs.logWriteSessionMilestone(path, "unlock_skipped_disabled", fh, nil, "")
		return errc
	}

	// If the node exists, prevent locking while unlocking.
	// If the node was renamed/moved, it may still need to be unlocked.
	if node, ok := fs.vfs.fetch(path); ok {
		node.lockMutex.Lock()
		defer node.lockMutex.Unlock()
	}

	// Check if lock exists in map - only hold mutex briefly
	fs.lockMapMutex.Lock()
	lockInfo, ok := fs.lockMap[path]
	fs.lockMapMutex.Unlock()

	if !ok {
		// If the lock map doesn't have an entry for this path, it means the file
		// was never locked, or it was locked by a different file handle.
		fs.log.Debug("RemoteFs: unlock: File not locked: %v fh=%v", path, fh)
		fs.logWriteSessionMilestone(path, "unlock_no_lock", fh, nil, "")
		return errc
	}
	if lockInfo.Fh != fh {
		// This is fine. It just means the file either wasn't locked or it was locked by a different file handle.
		fs.log.Debug("RemoteFs: unlock: File not locked by this handle: %v fh=%v", path, fh)
		fs.logWriteSessionMilestone(path, "unlock_different_handle", fh, nil, "lock_fh=%d", lockInfo.Fh)
		return errc
	}

	localPath, remotePath := fs.paths(path)
	fs.log.Debug("RemoteFs: unlock: file %v (%v)", remotePath, localPath)

	// Make API call without holding lockMapMutex
	fs.logWriteSessionMilestone(path, "unlock_delete_started", fh, nil, "")
	err := fs.ops.TryWithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		return fs.backend.deleteLock(files_sdk.LockDeleteParams{
			Path:  remotePath,
			Token: lockInfo.Lock.Token,
		})
	})
	if errors.Is(err, lim.ErrNoSlotsAvailable) {
		fs.logWriteSessionMilestone(path, "unlock_delete_no_slots", fh, nil, "errc=%d errno=%s", -fuse.EAGAIN, formatFuseErrno(-fuse.EAGAIN))
		return -fuse.EAGAIN
	}

	if files_sdk.IsNotExist(err) {
		// If the lock was already deleted, consider it a success.
		fs.log.Debug("RemoteFs: unlock: %v (%v) err=%v", remotePath, localPath, err)
		fs.lockMapMutex.Lock()
		delete(fs.lockMap, path)
		fs.lockMapMutex.Unlock()
		if node, ok := fs.vfs.fetch(path); ok {
			node.setLockOwner("")
		}
		fs.logWriteSessionMilestone(path, "unlock_already_deleted", fh, nil, "")
		return errc
	}
	// for any other error, handle it normally
	if errc = fs.handleError(path, err); errc != 0 {
		fs.logWriteSessionMilestone(path, "unlock_delete_failed", fh, nil, "errc=%d errno=%s", errc, formatFuseErrno(errc))
		return errc
	}

	// Delete from map - only hold mutex for this operation
	fs.lockMapMutex.Lock()
	delete(fs.lockMap, path)
	fs.lockMapMutex.Unlock()
	if node, ok := fs.vfs.fetch(path); ok {
		node.setLockOwner("")
	}
	fs.logWriteSessionMilestone(path, "unlock_deleted", fh, nil, "")

	return errc
}

func (fs *RemoteFs) paths(path string) (string, string) {
	return fs.localPath(path), fs.remotePath(path)
}

func (fs *RemoteFs) localPath(path string) string {
	return filepath.Join(fs.mountPoint, path)
}

func (fs *RemoteFs) remotePath(path string) string {
	return path_lib.Join(fs.root, path)
}

func (fs *RemoteFs) handleError(path string, err error) int {
	if err != nil {
		return fs.handleErrorMessage(path, err, err.Error())
	}
	return 0
}

func (fs *RemoteFs) handleUploadSessionError(path string, err error, session *writeSession) int {
	if writeSessionLastUploadErr(session, err) {
		return fs.handleErrorMessage(path, err, uploadLogMessage(err))
	}
	return fs.handleError(path, err)
}

func (fs *RemoteFs) handleErrorMessage(path string, err error, message string) int {
	if err != nil {
		localPath, remotePath := fs.paths(path)
		fs.log.Error("%v (%v): %s", remotePath, localPath, message)

		class, errc := classifyMountError(err)
		if class == "not_authenticated" {
			fs.events.Publish(events.AuthenticationFailedEvent{
				Reason: err.Error(),
			})
		}
		fs.log.Debug(
			"mount_error_mapped source=api path=%q remote_path=%q local_path=%q class=%s errc=%d errno=%s message=%q",
			path,
			remotePath,
			localPath,
			class,
			errc,
			formatFuseErrno(errc),
			message,
		)
		return errc
	}
	return 0
}

func classifyMountError(err error) (string, int) {
	if err == nil {
		return "none", 0
	}
	if files_sdk.IsNotAuthenticated(err) {
		return "not_authenticated", -fuse.EPERM
	}
	if files_sdk.IsNotExist(err) {
		return "not_exist", -fuse.ENOENT
	}
	if files_sdk.IsExist(err) {
		return "exist", -fuse.EEXIST
	}
	if errors.Is(err, lim.ErrNoSlotsAvailable) {
		return "no_slots_available", -fuse.EAGAIN
	}
	if isFolderNotEmpty(err) {
		return "folder_not_empty", -fuse.ENOTEMPTY
	}
	if isResourceLocked(err) {
		return "resource_locked", -fuse.EAGAIN
	}
	return "unknown", -fuse.EIO
}

func writeSessionLastUploadErr(session *writeSession, err error) bool {
	if session == nil || err == nil {
		return false
	}
	session.mu.Lock()
	lastUploadErr := session.lastUploadErr
	session.mu.Unlock()
	return lastUploadErr != nil && errors.Is(err, lastUploadErr)
}

func isResourceLocked(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "resource locked") || strings.Contains(msg, "exclusive lock")
}

func (fs *RemoteFs) delete(path string) (errc int) {
	err := fs.ops.TryWithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		return fs.backend.delete(files_sdk.FileDeleteParams{Path: fs.remotePath(path)})
	})
	if errors.Is(err, lim.ErrNoSlotsAvailable) {
		return -fuse.EAGAIN
	}

	// if there's an error, and it's a not-found, consider it a success.
	if files_sdk.IsNotExist(err) {
		fs.finalizeDelete(path)
		return errc
	}
	// for any other error, handle it normally
	if errc = fs.handleError(path, err); errc != 0 {
		return errc
	}
	fs.finalizeDelete(path)
	return errc
}

func (fs *RemoteFs) finalizeDelete(path string) {
	if node, ok := fs.vfs.fetch(path); ok {
		node.clearPendingVisible()
		node.markDeleted()
	}
	fs.vfs.remove(path)
	_ = fs.cacheStore.Delete(path)
}

func (fs *RemoteFs) loadParent(path string) (errc int) {
	start := time.Now()
	parentPath := path_lib.Dir(path)
	defer func() {
		elapsed := time.Since(start)
		if errc == 0 && elapsed < mountDiagnosticsSlowThreshold {
			return
		}
		fs.log.Debug(
			"mount_metadata_refresh op=loadParent requested_path=%q parent=%q errc=%d errno=%s duration=%s",
			path,
			parentPath,
			errc,
			formatFuseErrno(errc),
			elapsed,
		)
	}()

	if path == "/" {
		// If loading at the root, the parent can't be loaded. Just make sure the root exists.
		_, errc = fs.findDir(path)
		return errc
	}

	parent, ok := fs.vfs.fetch(parentPath)

	// Make sure the parent is actually a directory that exists before attempting to load it.
	if !ok || parent.infoExpired() {
		parent, errc = fs.findDir(parentPath)
		if errc != 0 {
			return errc
		}
	}

	if parent.info.nodeType != nodeTypeDir {
		// Don't log an error. Windows File Explorer sometimes treats shortcuts as parent directories.
		fs.log.Trace("RemoteFs: loadParent: Parent of %s is not a directory %s", path, parentPath)
		return -fuse.ENOTDIR
	}

	return fs.loadDir(parent)
}

func (fs *RemoteFs) findDir(path string) (node *fsNode, errc int) {
	remotePath := fs.remotePath(path)

	if remotePath == "/" {
		// Special case that the root directory of a Files.com site can't be stat'd.
		node = fs.vfs.getOrCreate(path, nodeTypeDir)
		node.updateInfo(fsNodeInfo{
			nodeType:     nodeTypeDir,
			creationTime: fs.initTime,
			modTime:      time.Now(),
		})
		return node, errc
	}

	var item files_sdk.File
	var err error
	err = fs.ops.TryWithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		item, err = fs.backend.find(files_sdk.FileFindParams{Path: remotePath})
		return err
	})
	if errors.Is(err, lim.ErrNoSlotsAvailable) {
		return nil, -fuse.EAGAIN
	}

	// Check for non-existence first so it doesn't get logged as an error, since this may be expected.
	if files_sdk.IsNotExist(err) {
		errc = -fuse.ENOENT
		return node, errc
	}
	if errc = fs.handleError(path, err); errc != 0 {
		return nil, errc
	}
	if !item.IsDir() {
		errc = -fuse.ENOTDIR
		return node, errc
	}

	node = fs.createNode(path, item)

	return node, errc
}

func (fs *RemoteFs) loadDir(node *fsNode) (errc int) {
	start := time.Now()
	refreshed := false
	defer func() {
		elapsed := time.Since(start)
		if errc == 0 && elapsed < mountDiagnosticsSlowThreshold {
			return
		}
		localPath, remotePath := fs.paths(node.path)
		fs.log.Debug(
			"mount_metadata_refresh op=loadDir path=%q remote_path=%q local_path=%q refreshed=%t errc=%d errno=%s duration=%s",
			node.path,
			remotePath,
			localPath,
			refreshed,
			errc,
			formatFuseErrno(errc),
			elapsed,
		)
	}()

	fs.loadDirMutexes.Lock(node.path)
	defer fs.loadDirMutexes.Unlock(node.path)
	if node.infoExpired() {
		refreshed = true
		fs.log.Debug("RemoteFs: loadDir: Refreshing directory listing: %v", node.path)
		err := node.updateChildPaths(fs.listDir)
		if errors.Is(err, lim.ErrNoSlotsAvailable) {
			return -fuse.EAGAIN
		}
		if errc = fs.handleError(node.path, err); errc != 0 {
			return errc
		}
	} else {
		fs.log.Trace("RemoteFs: loadDir: Skipping load of directory, info not expired: %v", node.path)
	}
	return errc
}

func (fs *RemoteFs) listDir(path string) (childPaths map[string]struct{}, opErr error) {
	fs.log.Trace("RemoteFs: listDir: Listing directory: %v", path)

	opErr = fs.ops.TryWithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		it, err := fs.backend.listFor(files_sdk.FolderListForParams{Path: fs.remotePath(path)})
		if err != nil {
			return err
		}

		childPaths = make(map[string]struct{})

		for it.Next() {
			item := it.File()

			childPath := path_lib.Join(path, item.DisplayName)
			childPaths[childPath] = struct{}{}

			fs.createNode(childPath, item)
		}
		err = it.Err()
		if err != nil {
			return err
		}

		if fs.disableLocking {
			return err
		}

		locks, err := fs.backend.listLocksFor(files_sdk.LockListForParams{
			Path:            fs.remotePath(path),
			IncludeChildren: lib.Ptr(true),
		})

		if err != nil {
			return err
		}

		for locks.Next() {
			lock := locks.Lock()
			childPath := path_lib.Join(path, path_lib.Base(lock.Path))

			// Ignore paths where the lock is held by this file system.
			if _, ok := fs.lockMap[childPath]; ok {
				continue
			}

			if child, ok := fs.vfs.fetch(childPath); ok {
				fs.log.Trace("RemoteFs: listDir: Found lock for child path %v, setting lock owner to %v", childPath, lock.Username)
				child.info.lockOwner = lock.Username
			}
		}
		return locks.Err()
	})

	return childPaths, opErr
}

// startWebSyncWatcher starts a periodic watcher that notifies all cached directories
// to catch changes made via the web interface. This triggers Explorer/Finder to refresh.
func (fs *RemoteFs) startWebSyncWatcher(interval time.Duration) {
	if runtime.GOOS != "windows" {
		return
	}
	if fs.webSyncTicker != nil {
		return
	}

	fs.webSyncTicker = time.NewTicker(interval)
	if fs.stopCh == nil {
		fs.stopCh = make(chan struct{})
	}
	fs.wg.Add(1)

	go func() {
		defer fs.wg.Done()
		for {
			select {
			case <-fs.webSyncTicker.C:
				if fs.vfs != nil {
					var dirs []string
					if fs.vfs.handles != nil {
						dirs = fs.vfs.handles.OpenDirectoryPaths()
					}
					// If nothing is open, skip this tick
					if len(dirs) == 0 {
						continue
					}
					for _, dir := range dirs {
						path := fs.localPath(dir)
						if err := shell.NotifyUpdatedDir(path); err != nil {
							fs.log.Error("shell notify failed for %s: %v", path, err)
						}
					}
				}
			case <-fs.stopCh:
				return
			}
		}
	}()
}

// stopWebSync stops the web sync watcher and related goroutines.
func (fs *RemoteFs) stopWebSync() {
	if fs.stopCh != nil {
		close(fs.stopCh)
		fs.stopCh = nil
	}

	if fs.webSyncTicker != nil {
		fs.webSyncTicker.Stop()
		fs.webSyncTicker = nil
	}

	fs.wg.Wait()
}

func (fs *RemoteFs) createNode(path string, item files_sdk.File) *fsNode {
	var nt nodeType
	if item.IsDir() {
		nt = nodeTypeDir
	} else {
		nt = nodeTypeFile
	}
	var existingCreationTime time.Time
	// best-effort invalidate stale data
	if prev, ok := fs.vfs.fetch(path); ok && prev.info.nodeType == nodeTypeFile {
		if prev.info.size != item.Size || !prev.info.modTime.Equal(item.ModTime()) {
			_ = fs.cacheStore.Delete(path)
		}
		existingCreationTime = prev.info.creationTime
	}

	node := fs.vfs.getOrCreate(path, nt)
	node.markMaterialized()
	if nt == nodeTypeFile && node.hasActiveWriteSession() {
		node.setRemotePermissions(item.Permissions)
		return node
	}
	node.clearPendingVisible()
	creationTime := item.CreationTime()
	if nt == nodeTypeFile && !existingCreationTime.IsZero() {
		creationTime = existingCreationTime
	}
	node.updateInfo(fsNodeInfo{
		nodeType:     nt,
		size:         item.Size,
		modTime:      item.ModTime(),
		creationTime: creationTime,
	})
	node.setRemotePermissions(item.Permissions)
	if item.DownloadUri != "" {
		node.setDownloadURI(item.DownloadUri)
	}

	return node
}

func (fs *RemoteFs) waitForAction(ctx context.Context, action files_sdk.FileAction, operation string) error {
	var migration files_sdk.FileMigration
	var err error
	err = fs.ops.TryWithLimit(ctx, lim.FuseOpOther, func(ctx context.Context) error {
		migration, err = fs.backend.wait(action, func(migration files_sdk.FileMigration) {
			fs.log.Trace("RemoteFs: watchForAction: waiting for migration")
		}, files_sdk.WithContext(ctx))
		return err
	})
	if errors.Is(err, lim.ErrNoSlotsAvailable) {
		return lim.ErrNoSlotsAvailable
	}

	if err == nil && migration.Status != "completed" {
		return fmt.Errorf("%v did not complete successfully: %v", operation, migration.Status)
	}
	return err
}

func (fs *RemoteFs) isLockFile(path string) bool {
	return isMsOfficeOwnerFile(path) && !fs.disableLocking
}

// Methods below are part of the fuse.FileSystemInterface, but not supported by
// this implementation. They exist here to support logging for visibility of how
// the underlying fuse layer calls into this implementation.

// Mknod creates a file node.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *RemoteFs) Mknod(path string, mode uint32, dev uint64) int {
	fs.log.Trace("RemoteFs: Mknod: path=%v, mode=%o, dev=%v", path, mode, dev)
	return -fuse.ENOSYS
}

// Link creates a hard link to a file.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *RemoteFs) Link(oldpath string, newpath string) int {
	fs.log.Trace("RemoteFs: Link: old=%v, new=%v", oldpath, newpath)
	return -fuse.ENOSYS
}

// Symlink creates a symbolic link.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *RemoteFs) Symlink(target string, newpath string) int {
	fs.log.Trace("RemoteFs: Symlink: target=%v, newpath=%v", target, newpath)
	return -fuse.ENOSYS
}

// Readlink reads the target of a symbolic link.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *RemoteFs) Readlink(path string) (int, string) {
	fs.log.Trace("RemoteFs: Readlink: path=%v", path)
	return -fuse.ENOSYS, ""
}

// Chown changes the owner and group of a file.
// On Windows this is treated as a no-op success for compatibility.
func (fs *RemoteFs) Chown(path string, uid uint32, gid uint32) int {
	node := fs.vfs.getOrCreate(path, nodeTypeFile)
	node.setOwner(uid, gid)
	fs.log.Debug("RemoteFs: Chown: path=%v, uid=%v, gid=%v -> errc=0", path, uid, gid)
	return 0
}

// Access checks file access permissions.
func (fs *RemoteFs) Access(path string, mask uint32) int {
	fs.log.Trace("RemoteFs: Access: path=%v, mask=%v", path, mask)
	node, errc := fs.fetchNodeWithParentRefresh(path)
	if errc != 0 {
		return errc
	}

	permissions := node.getRemotePermissions()
	if permissions == "" {
		return 0
	}

	if mask&accessMaskRead != 0 && !node.isReadable() {
		localPath, remotePath := fs.paths(path)
		fs.log.Debug("mount_permission_denied source=access_permissions op=Access path=%q remote_path=%q local_path=%q mask=%d permissions=%q errc=%d errno=%s", path, remotePath, localPath, mask, permissions, -fuse.EACCES, formatFuseErrno(-fuse.EACCES))
		return -fuse.EACCES
	}
	if mask&accessMaskWrite != 0 && !node.isWritable() {
		localPath, remotePath := fs.paths(path)
		fs.log.Debug("mount_permission_denied source=access_permissions op=Access path=%q remote_path=%q local_path=%q mask=%d permissions=%q errc=%d errno=%s", path, remotePath, localPath, mask, permissions, -fuse.EACCES, formatFuseErrno(-fuse.EACCES))
		return -fuse.EACCES
	}
	if mask&accessMaskExecute != 0 {
		if node.info.nodeType == nodeTypeDir {
			if !node.isListable() {
				localPath, remotePath := fs.paths(path)
				fs.log.Debug("mount_permission_denied source=access_permissions op=Access path=%q remote_path=%q local_path=%q mask=%d permissions=%q errc=%d errno=%s", path, remotePath, localPath, mask, permissions, -fuse.EACCES, formatFuseErrno(-fuse.EACCES))
				return -fuse.EACCES
			}
		} else {
			localPath, remotePath := fs.paths(path)
			fs.log.Debug("mount_permission_denied source=access_permissions op=Access path=%q remote_path=%q local_path=%q mask=%d permissions=%q errc=%d errno=%s", path, remotePath, localPath, mask, permissions, -fuse.EACCES, formatFuseErrno(-fuse.EACCES))
			return -fuse.EACCES
		}
	}

	return 0
}

// Flush flushes cached file data.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *RemoteFs) Flush(path string, fh uint64) int {
	fs.log.Trace("RemoteFs: Flush: path=%v, fh=%v", path, fh)

	handle, node, ok := fs.vfs.handles.Lookup(fh)
	if !ok || handle.IsReadOnly() {
		return 0
	}

	// On macOS, close(2) returns after Flush completes — not after Release.
	// Release runs asynchronously, so finalizing the upload there causes
	// fuse_do_release to block, which triggers a libfuse-t assertion failure
	// (open_count > 0) when a second release arrives before the first completes.
	// Finalizing here ensures Release returns quickly and the upload is committed
	// before the OS reports the copy as done.
	if session := node.getWriteSession(); session != nil {
		if !session.hasHandle(fh) {
			return 0
		}
		if err := node.poisonedWriteSessionErr(); err != nil {
			return -fuse.EIO
		}
		return fs.flushWriteSession(path, node, fh)
	}

	return 0
}

// Fsyncdir synchronizes directory contents.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *RemoteFs) Fsyncdir(path string, datasync bool, fh uint64) int {
	fs.log.Trace("RemoteFs: Fsyncdir: path=%v, datasync=%v, fh=%v", path, datasync, fh)
	return 0
}

// The [Foo]xattr implementations below explicitly return 0 to indicate that
// extended attributes are "supported" in order to ensure that the other xattr
// methods are called for debugging visibility, but are all no-op implementations.

// Getxattr gets extended attributes.
// Any return value other than -fuse.ENOSYS indicates support for extended
// attributes, but also expects Setxattr, Listxattr, and Removexattr to exist
// for extended attribute support.
func (fs *RemoteFs) Getxattr(path string, name string) (int, []byte) {
	fs.log.Debug("RemoteFs: Getxattr: path=%v, name=%v", path, name)
	return 0, []byte{}
}

// Setxattr sets extended attributes.
func (fs *RemoteFs) Setxattr(path string, name string, value []byte, flags int) int {
	fuseFlags := ff.NewFuseFlags(flags)
	fs.log.Debug("RemoteFs: Setxattr: path=%v, name=%v, value=%v flags=%v", path, name, value, fuseFlags)
	return 0
}

// Removexattr removes extended attributes.
func (fs *RemoteFs) Removexattr(path string, name string) int {
	fs.log.Debug("RemoteFs: Removexattr: path=%v, name=%v", path, name)
	return 0
}

// Listxattr lists extended attributes.
func (fs *RemoteFs) Listxattr(path string, fill func(name string) bool) int {
	fs.log.Debug("RemoteFs: Listxattr: path=%v", path)
	return 0
}

// FileSystemOpenEx is the interface that wraps the OpenEx and CreateEx methods.

// OpenEx and CreateEx are similar to Open and Create except that they allow
// direct manipulation of the FileInfo_t struct (which is analogous to the
// FUSE struct fuse_file_info). If implemented, they are preferred over
// Open and Create.
func (fs *RemoteFs) CreateEx(path string, mode uint32, fi *fuse.FileInfo_t) int {
	fs.log.Debug("RemoteFs: CreateEx: path=%v, mode=%o, fi=%v", path, mode, fi)
	fuseFlags := ff.NewFuseFlags(fi.Flags)
	var errc int
	var fh uint64
	if fuseFlags.IsCreate() {
		errc, fh = fs.Create(path, fi.Flags, mode)
	} else {
		errc, fh = fs.Open(path, fi.Flags)
	}
	fi.Fh = fh
	return errc
}

func (fs *RemoteFs) OpenEx(path string, fi *fuse.FileInfo_t) int {
	fs.log.Debug("RemoteFs: OpenEx: path=%v, fi=%v", path, fi)
	errc, fh := fs.Open(path, fi.Flags)
	fi.Fh = fh
	return errc
}

// Getpath is part of the FileSystemGetpath interface and
// allows a case-insensitive file system to report the correct case of a file path.
func (fs *RemoteFs) Getpath(path string, fh uint64) (int, string) {
	fs.log.Trace("RemoteFs: Getpath: path=%v, fh=%v", path, fh)
	return -fuse.ENOSYS, path
}

// Chflags is part of the FileSystemChflags interface and
// changes the BSD file flags (Windows file attributes).
func (fs *RemoteFs) Chflags(path string, flags uint32) int {
	fs.log.Trace("RemoteFs: Chflags: path=%v, flags=%v", path, flags)
	return -fuse.ENOSYS
}

// Setcrtime is part of the FileSystemSetcrtime interface and
// changes the file creation (birth) time.
func (fs *RemoteFs) Setcrtime(path string, tmsp fuse.Timespec) int {
	fs.log.Trace("RemoteFs: Setcrtime: path=%v, tmsp=%v", path, tmsp)
	return -fuse.ENOSYS
}

// Setchgtime is part of the FileSystemSetchgtime interface and
// changes the file change (ctime) time.
func (fs *RemoteFs) Setchgtime(path string, tmsp fuse.Timespec) int {
	fs.log.Trace("RemoteFs: Setchgtime: path=%v, tmsp=%v", path, tmsp)
	return -fuse.ENOSYS
}
