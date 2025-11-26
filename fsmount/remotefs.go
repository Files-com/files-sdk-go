package fsmount

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	path_lib "path"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/fsmount/events"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/mem"
	fsio "github.com/Files-com/files-sdk-go/v3/fsmount/internal/io"
	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
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

	// limits using in configuring the NewFuseOpLimiter for RemoteFs
	downloadOpLimit = 10
	uploadOpLimit   = 4
	otherOpLimit    = 4
	globalOpLimit   = 16
)

// compile time assertions that the cache implementations satisfy the fsCache interface
var _ cacheStore = (*disk.DiskCache)(nil)
var _ cacheStore = (*mem.MemoryCache)(nil)

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

	fileClient      *file.Client
	lockClient      *lock.Client
	apiKeyClient    *api_key.Client
	currentUserId   int64
	migrationClient *file_migration.Client
	lockMap         map[string]*lockInfo
	lockMapMutex    sync.Mutex
	loadDirMutex    sync.Mutex

	debugFuse bool

	initOnce sync.Once
	initTime time.Time

	cacheStore cacheStore

	gatesMu    sync.Mutex
	readyGates map[string]*cache.ReadyGate

	events events.EventPublisher

	ops *lim.FuseOpLimiter
}

// cacheStore defines the interface for the file system cache used by RemoteFs and allows for alternative
// implementations. e.g. an in-memory cache implementation vs a disk-based cache implementation.
type cacheStore interface {
	Read(path string, buff []byte, ofst int64) (n int, err error)
	Write(path string, buff []byte, ofst int64) (n int, err error)
	Delete(path string) bool
	StartMaintenance()
	StopMaintenance()
	// Pin increments the reference count for a file, preventing it from being evicted
	Pin(path string)
	// Unpin decrements the reference count for a file
	Unpin(path string)
}

// cacheReader wraps the cacheStore to provide an io.Reader interface for reading cached files.
// This is used to seed OrderedPipe with initial content for partial file updates.
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
		ops:              limiter,
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
	// Guard with a sync.Once because Init is called from fsmount.Mount, but cgofuse also calls Init
	// when it mounts the file system.
	fs.initOnce.Do(func() {
		if fs.fileClient == nil {
			fs.fileClient = &file.Client{Config: *fs.config}
			fs.lockClient = &lock.Client{Config: *fs.config}
			fs.apiKeyClient = &api_key.Client{Config: *fs.config}
			fs.migrationClient = &file_migration.Client{Config: *fs.config}
			fs.lockMap = make(map[string]*lockInfo)
		}

		// no need to guard this with an operation limit since it's only called once during initialization
		key, err := fs.apiKeyClient.FindCurrent()
		if err != nil {
			fs.log.Error("Failed to find metadata for current API key, file exclusivity locks may not work as expected: %v", err)
			// set locking to false?
		}
		fs.currentUserId = key.UserId
		// store the time the file system was initialized to use as the creation time for the root directory
		fs.initTime = time.Now()
		fs.log.Debug("RemoteFs: RemoteFs initialized successfully. Remote file system root: %s", fs.root)
	})
	// start the disk cache maintenance goroutine
	// this does not block and ensures only one goroutine is started
	fs.cacheStore.StartMaintenance()

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
}

func (fs *RemoteFs) Validate() error {
	fs.Init()
	// Make sure the root directory can be listed.
	// no need to guard this with an operation limit since it's only called once during initialization
	it, err := fs.fileClient.ListFor(files_sdk.FolderListForParams{
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

func (fs *RemoteFs) Mkdir(path string, mode uint32) (errc int) {
	localPath, remotePath := fs.paths(path)
	fs.log.Debug("RemoteFs: Mkdir: %v (%v) (mode=%o)", remotePath, localPath, mode)

	err := fs.ops.WithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		_, err := fs.fileClient.CreateFolder(files_sdk.FolderCreateParams{Path: remotePath})
		return err
	})
	if files_sdk.IsExist(err) {
		return errc
	}

	if errc = fs.handleError(path, err); errc != 0 {
		return errc
	}

	node := fs.vfs.getOrCreate(path, nodeTypeDir)
	node.updateSize(0)

	return errc
}

func (fs *RemoteFs) Unlink(path string) (errc int) {
	localPath, remotePath := fs.paths(path)

	node, exists := fs.vfs.fetch(path)
	if !exists {
		// If the node doesn't exist, it can not be deleted.
		fs.log.Debug("RemoteFs: Unlink: File not found: %v (%v)", remotePath, localPath)
		return errc
	}

	// If the node is locked, it can not be deleted.
	if node.isLocked() {
		fs.log.Info("Cannot delete locked file: %v (%v)", remotePath, localPath)
		return -fuse.ENOLCK
	}

	// If the node is being written to, cancel the upload and delete the file from the remote API.
	// This is necessary because the file may be in the middle of being written to, and the upload may not have completed yet.
	if node.isWriterOpen() {
		fs.log.Debug("RemoteFs: Unlink: Canceling upload for: %v (%v)", remotePath, localPath)
		node.cancelUpload()
	}

	// The fs may have been in the middle of writing the file, so don't log until here.
	fs.log.Info("Deleting file: %v (%v)", remotePath, localPath)
	return fs.delete(path)
}

func (fs *RemoteFs) Rmdir(path string) int {
	localPath, remotePath := fs.paths(path)
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

	defer node.expireInfo()

	// If there is no active upload for this node, proceed with the rename.
	if !node.isWriterOpen() {
		fs.log.Info("Renaming %v to %v (%v to %v)", oldRemotePath, newRemotePath, oldLocalPath, newLocalPath)

		params := files_sdk.FileMoveParams{
			Path:        oldRemotePath,
			Destination: newRemotePath,
			Overwrite:   lib.Ptr(true),
		}

		err := fs.ops.WithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
			action, err := fs.fileClient.Move(params)
			if err != nil {
				return err
			}
			return fs.waitForAction(ctx, action, "move")
		})
		if errc = fs.handleError(oldpath, err); errc != 0 {
			return errc
		}

		fs.rename(oldpath, newpath)
		_ = fs.cacheStore.Delete(oldpath)
		_ = fs.cacheStore.Delete(newpath)

		return errc
	}

	// There must be an active upload for this node. Update local VFS map immediately so listings/
	// lookups stay consistent. In order to support the pattern of writing a file to a temporary
	// name, then renaming once the upload is complete, the node, and upload must be updated to
	// reflect the new path. Before the upload is finalized, there is a callback that inspects the
	// file node to get the final path and upload ref to complete the upload.
	fs.rename(oldpath, newpath)
	node.clearDownloadURI()

	return errc
}

func (fs *RemoteFs) Utimens(path string, tmsp []fuse.Timespec) (errc int) {
	localPath, remotePath := fs.paths(path)
	modT := tmsp[1].Time()
	fs.log.Debug("RemoteFs: Utimens: Updating mtime for: %v (%v) (mtime=%v)", remotePath, localPath, modT)

	node := fs.vfs.getOrCreate(path, nodeTypeFile)
	node.info.modTime = modT

	if node.isWriterOpen() {
		// If the fs is writing to the file, no need update the mtime. It will be updated when the write completes.
		return errc
	}

	params := files_sdk.FileUpdateParams{
		Path:          remotePath,
		ProvidedMtime: &node.info.modTime,
	}

	err := fs.ops.WithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		_, err := fs.fileClient.Update(params)
		return err
	})

	return fs.handleError(path, err)
}

func (fs *RemoteFs) Create(path string, flags int, mode uint32) (errc int, fh uint64) {
	localPath, remotePath := fs.paths(path)
	fuseFlags := ff.NewFuseFlags(flags)
	fh, handle := fs.vfs.handles.Open(nil, fuseFlags)

	fs.log.Debug("RemoteFs: Create: Creating file: %v (%v) (flags=%v, mode=%o, fh=%v)", remotePath, localPath, fuseFlags, mode, fh)

	// Load the parent directory to populate the vfs nodes map
	if errc = fs.loadParent(path); errc != 0 {
		// Release handle on error since FUSE won't call Release() for failed creates
		fs.vfs.handles.Release(fh)
		return errc, fh
	}

	node, exists := fs.vfs.fetch(path)
	if exists && node.isWriterOpen() {
		// the node exists and has an open writer, if the writer's offset is
		// greater than zero, it means the file is actively being written to
		if node.writer.Offset() > 0 {
			fs.log.Info("Cannot create file while writing: %v (%v)", remotePath, localPath)
			// Release handle on error since FUSE won't call Release() for failed creates
			fs.vfs.handles.Release(fh)
			return -fuse.EEXIST, fh
		}
		// the node exists, and has an open writer, but the writer's offset is zero,
		// meaning the file was created but nothing has been written to it yet.
		// In this case, create a new file handle and return it. The writer will
		// only be closed when a file handle that has written data is released to
		// avoid creating multiple upload events for the same file.
		fs.log.Debug("RemoteFs: Create: File already exists, but no data has been written: %v (%v)", remotePath, localPath)
		handle.node = node
		return errc, fh
	}

	// TODO: decide if this makes sense. the node exists and the cache data is recent
	// so return an error for the Create call?
	if exists && !node.infoExpired() {
		fs.log.Debug("RemoteFs: Create: Node exists, cache data is recent, but no open writer: %v (%v)", remotePath, localPath)
		// Release handle on error since FUSE won't call Release() for failed creates
		fs.vfs.handles.Release(fh)
		return -fuse.EEXIST, fh
	}

	if !exists {
		node = fs.vfs.getOrCreate(path, nodeTypeFile)
	}

	node.updateSize(0)
	handle.node = node

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
		return errc, fh
	}

	// Mark write intent - writer will be created lazily on first Write
	node.markWriteIntent(fh)
	fs.log.Debug("RemoteFs: Create: marked write intent for %v (%v), fh=%v", remotePath, localPath, fh)

	return errc, fh
}

func (fs *RemoteFs) Open(path string, flags int) (errc int, fh uint64) {
	fuseFlags := ff.NewFuseFlags(flags)
	node := fs.vfs.getOrCreate(path, nodeTypeFile)
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
		return errc, fh
	}

	// Single-writer enforcement: block if writer exists and has writes
	if node.writerIsOpen() {
		_, owner, hasWrites := node.writerSnapshot()
		if hasWrites && owner != fh {
			fs.log.Debug("RemoteFs: Open: writer already active with writes for path=%v, owner=%v", path, owner)
			// EBUSY error - manually Unpin and Release handle since Release() won't be called for failed opens
			if fs.cacheStore != nil {
				fs.cacheStore.Unpin(path)
			}
			fs.vfs.handles.Release(fh)
			return -fuse.EBUSY, fh
		}
	}

	// Mark write intent - writer will be created lazily on first Write
	node.markWriteIntent(fh)
	fs.log.Debug("RemoteFs: Open: marked write intent for path=%v, fh=%v", path, fh)

	return errc, fh
}

func (fs *RemoteFs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	fs.log.Trace("RemoteFs: Getattr: path=%v, fh=%v", path, fh)
	// If the file handle is open, extend the TTL of the open handle. The info may have expired,
	// but the handle is still open, meaning the OS is still using the file. This can happen if there
	// are multiple simultaneous uploads, but they haven't all received a write request in the last
	// cacheTTL duration. If the Getattr call returns an error, the OS will remove the file from the
	// Explorer/Finder window until the upload completes, and a subsequent Getattr call succeeds, which
	// is a bad user experience.
	fs.vfs.handles.ExtendOpenHandleTtls()
	if node, exists := fs.vfs.fetch(path); exists && !node.infoExpired() {
		fs.log.Trace("RemoteFs: Getattr: using cached stat, path=%v, size=%v, mtime=%v", path, node.info.size, node.info.modTime)
		getStat(node.info, stat)
		return errc
	}

	if errc = fs.loadParent(path); errc != 0 {
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
			fs.log.Trace("RemoteFs: Getattr: File not found: %v (%v)", remotePath, localPath)
			return -fuse.ENOENT
		}
	}

	fs.log.Trace("RemoteFs: Getattr: path=%v, size=%v, mtime=%v", path, node.info.size, node.info.modTime)
	getStat(node.info, stat)

	return errc
}

func (fs *RemoteFs) Truncate(path string, size int64, fh uint64) (errc int) {
	// The word truncate is overloaded here. The intention is to set the size of the
	// file to the size getting passed in, NOT to truncate the file to zero bytes.

	localPath, remotePath := fs.paths(path)
	fs.log.Debug("RemoteFs: Truncate: %v (%v) (size=%v, fh=%v)", remotePath, localPath, size, fh)

	node := fs.vfs.getOrCreate(path, nodeTypeFile)
	node.updateSize(size)

	// Mark write intent - actual truncation will happen if data is written
	// Per requirements: O_TRUNC creates writer but waits for actual data before uploading
	node.markWriteIntent(fh)
	fs.log.Debug("RemoteFs: Truncate: marked write intent for %v (%v)", remotePath, localPath)

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

	// Attempt to read from the temporary file backing the writer. If the read can't be satisfied
	// from the temporary file, it will return zero bytes, and the logic will fall through to
	// reading from the cache or remote API.
	if node.writerIsOpen() {
		n = node.readFromWriter(buff, ofst)
		if n > 0 {
			handle.incrementRead(int64(n))
			fs.log.Trace("RemoteFs: Read: readAt: path=%v, ofst=%d, read %d bytes from writer pipe", path, ofst, n)
			return n
		}
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

	if !node.writerIsOpen() && size == 0 {
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

	// Attempt to read from the disk cache. If the read can't be satisfied from the disk cache, it
	// will return zero bytes, and the logic will fall through to reading from the remote API.
	n, err := fs.cacheStore.Read(path, buff, ofst)
	if err != nil {
		fs.log.Debug("RemoteFs: Read: diskCache.Read error: %v", err)
	}

	if n > 0 {
		handle.incrementRead(int64(n))
		fs.log.Trace("RemoteFs: Read: readAt: path=%v, ofst=%d, read %d bytes from disk cache", path, ofst, n)
		return n
	}

	// At this point, the read request could not be satisfied from the temporary file backing the
	// writer or the disk cache, so read from the remote API.
	endOffset := ofst + int64(len(buff))
	readyGate, exists := fs.findOrCreateGate(path)
	if !exists {
		// start the download in a new goroutine, which will populate the disk cache
		go fs.fillCache(context.Background(), path, node.downloadUri, readyGate, fh)
	}
	readyGate.Add()
	defer readyGate.Done()

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
	n, err = fs.cacheStore.Read(path, buff[:want], ofst)
	if err != nil {
		fs.log.Debug("RemoteFs: Read: diskCache.Read error after WaitFor: %v", err)
		return -fuse.EAGAIN
	}
	fs.log.Trace("RemoteFs: Read: ok path=%v ofst=%d read=%d", path, ofst, n)
	handle.incrementRead(int64(n))
	return n
}

func (fs *RemoteFs) Write(path string, buff []byte, ofst int64, fh uint64) (n int) {
	fs.log.Debug("RemoteFs: Write: path=%v, len=%v, ofst=%v, fh=%v", path, len(buff), ofst, fh)

	_, node, ok := fs.vfs.handles.Lookup(fh)

	if !ok {
		fs.log.Debug("RemoteFs: Write: file handle %v not found for path %v", fh, path)
		return -fuse.EBADF
	}

	// Lazy-create writer on first write with initial content and cache writer support
	writer, created, err := node.ensureWriter(fs, fh, func() io.Reader {
		// Get initial content from cache for partial file updates
		if node.info.size > 0 {
			// Try to get content from cache
			cacheReader := &cacheReader{
				cache:  fs.cacheStore,
				path:   path,
				size:   node.info.size,
				logger: fs.log,
			}
			return cacheReader
		}
		return nil
	}, func() fsio.CacheWriter {
		// Return a cache writer that updates cache in real-time as writes occur
		return func(data []byte, offset int64) (int, error) {
			return fs.cacheStore.Write(path, data, offset)
		}
	})

	if err != nil {
		fs.log.Error("RemoteFs: Write: failed to ensure writer for %v: %v", path, err)
		return -fuse.EIO
	}

	if created {
		fs.log.Debug("RemoteFs: Write: created new writer for path=%v, fh=%v", path, fh)
	}

	// Write the buffer to the ordered pipe
	n, writeErr := writer.WriteAt(buff, ofst)
	if errc := fs.handleError(path, writeErr); errc != 0 {
		return errc
	}
	return n
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

	// Check if this handle owns a writer with uncommitted writes.
	_, owner, hasWrites := node.writerSnapshot()
	iOwn := owner == fh

	// set the node's expire time to expired so that the next Getattr call will trigger a reload of the node's info from the remote API.
	defer node.expireInfo()

	// If this handle owns the writer and has uncommitted writes, finalize the upload.
	if iOwn && hasWrites {
		fs.log.Debug("RemoteFs: Release: finalizing upload for path=%v, fh=%v", path, fh)
		if errc := fs.fsyncNode(path, node, fh); errc != 0 {
			return errc
		}
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

	localPath, remotePath := fs.paths(path)

	// This happens a lot, so log at trace level.
	fs.log.Trace("RemoteFs: Readdir: Listing folder: %v (%v)", remotePath, localPath)

	fillNode, _ := fs.vfs.fetch(path)

	// Force a load of the directory entries from the remote to make sure
	// the local vfs representation is up to date.
	if errc = fs.loadDir(fillNode); errc != 0 {
		return errc
	}

	fill(".", nil, 0)
	fill("..", nil, 0)

	// construct a list of child entries for the current directory
	entries := make([]string, len(fillNode.childPaths))
	pos := 0
	for childPath := range fillNode.childPaths {
		entries[pos] = childPath
		pos++
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
		if !slices.Contains(entries, openNode.path) && path == filepath.Dir(openNode.path) {
			fs.log.Debug("RemoteFs: Readdir: Child entries %v: for path %s, does not include open handle: %v, adding %v", entries, path, handle, openNode.path)
			entries = append(entries, openNode.path)
		}
	}

	// sort the entries in order to provide a consistent sort order when calling fill
	slices.Sort(entries)
	for _, entryPath := range entries {
		if entryNode, ok := fs.vfs.fetch(entryPath); ok {
			fs.log.Trace("RemoteFs: Readdir: Calling fill for entry: %v (%v)", entryPath, entryPath)
			fill(path_lib.Base(entryPath), getStat(entryNode.info, nil), 0)
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
// If an upload is active but the writer is already closed (finalizing in background),
// wait for completion. Otherwise, fall back to ENOSYS.
func (fs *RemoteFs) Fsync(path string, datasync bool, fh uint64) (errc int) {
	fs.log.Debug("RemoteFs: Fsync: path=%v, datasync=%v, fh=%v", path, datasync, fh)

	_, node, ok := fs.vfs.handles.Lookup(fh)
	if !ok {
		fs.log.Debug("RemoteFs: Fsync: file handle not found for path: %v, fh: %v", path, fh)
		return -fuse.EBADF
	}

	// Check if this handle owns a writer with uncommitted writes.
	_, owner, hasWrites := node.writerSnapshot()
	if owner != fh {
		// This handle doesn't own the writer, nothing to fsync.
		fs.log.Debug("RemoteFs: Fsync: handle does not own writer for path=%v, fh=%v", path, fh)
		return 0
	}

	if !hasWrites {
		// No writes to sync.
		fs.log.Debug("RemoteFs: Fsync: no writes to sync for path=%v, fh=%v", path, fh)
		return 0
	}

	// Finalize the upload.
	return fs.fsyncNode(path, node, fh)
}

// fsyncNode finalizes an upload for a node with uncommitted writes.
// It closes the writer, waits for the upload to complete, then discards the writer.
// Note: Cache is already updated in real-time via the CacheWriter callback set during ensureWriter.
func (fs *RemoteFs) fsyncNode(path string, node *fsNode, fh uint64) (errc int) {
	fs.log.Debug("RemoteFs: fsyncNode: closing writer for path=%v, fh=%v", path, fh)

	// Close the writer, which signals that no more data will be written.
	if err := node.closeWriter(); err != nil {
		fs.log.Debug("RemoteFs: fsyncNode: error closing writer for %v: %v", path, err)
		return -fuse.EIO
	}

	// Wait for the upload to complete with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), fsyncTimeout)
	defer cancel()

	if err := node.waitForUploadIfFinalizing(ctx); err != nil {
		fs.log.Debug("RemoteFs: fsyncNode: error waiting for upload to complete for %v: %v", path, err)
		return -fuse.ETIMEDOUT
	}

	// Discard the writer after successful upload.
	node.discardWriter()
	fs.log.Debug("RemoteFs: fsyncNode: upload complete and writer discarded for path=%v, fh=%v", path, fh)

	return 0
}

// copyWriterToCache copies the writer's temp file content to the cache.
// This is called before closing the writer to preserve the uploaded content for subsequent writes.
func (fs *RemoteFs) uploadProgressFunc(node *fsNode) func(int64) {
	return func(delta int64) {
		// If the write was successful, extend the node's ttl and keep track of the number
		// of bytes written for logging purposes.
		// Extend the node's TTL and keep track of bytes written for logging/sweeping.
		node.extendTtl()
		node.recordProgress(delta)
	}
}

func (fs *RemoteFs) writeFile(path string, reader io.Reader, mtime time.Time, fh uint64) {
	localPath, remotePath := fs.paths(path)
	fs.log.Info("Starting upload: %v (%v)", remotePath, localPath)
	_, node, _ := fs.vfs.handles.Lookup(fh)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	node.startUpload(path, cancel)
	uploadOpts := []file.UploadOption{
		file.UploadWithContext(ctx),
		file.UploadWithDestinationPath(remotePath),
		file.UploadWithReader(reader),
		file.UploadWithProvidedMtime(mtime),
		file.UploadWithProgress(fs.uploadProgressFunc(node)),

		// Using the WithUploadStartedCallback option allows capturing
		// the upload reference as soon as the upload starts, which is needed in
		// order to support renaming the file during an active upload.
		file.WithUploadStartedCallback(func(part files_sdk.FileUploadPart) {
			fs.log.Debug("RemoteFs: Uploading part number %d, of: %v, ref: '%v'", part.PartNumber, remotePath, part.Ref)
			node.captureRef(part.Ref)
		}),

		// Using the WithUploadRenamedCallback option allows renaming the upload
		// while it is in progress. This is needed in order to support the pattern
		// of writing a file to a temporary name, then renaming it to the final
		// name once the upload is complete.
		file.WithUploadRenamedCallback(func() (string, string) {
			if remotePath != node.path {
				fs.log.Debug("RemoteFs: writeFile: in progress upload renamed from: %v to %v", remotePath, node.path)
			}
			return node.pathAndRef()
		}),
	}
	if fs.writeConcurrency != 0 {
		uploadOpts = append(uploadOpts, file.UploadWithManager(manager.ConcurrencyManager{}.New(fs.writeConcurrency)))
	}

	start := time.Now()

	var u file.UploadResumable
	var err error
	err = fs.ops.WithLimit(ctx, lim.FuseOpUpload, func(ctx context.Context) error {
		u, err = fs.fileClient.UploadWithResume(uploadOpts...)
		return err
	})

	if err != nil {
		// this is only an error if the upload was not cancelled. If the upload was cancelled, it should not be logged as an error.
		if !errors.Is(err, context.Canceled) && !files_sdk.IsNotExist(err) {
			fs.log.Error("Error uploading file: %v (%v): %v", remotePath, localPath, err)
		}
		return
	}
	node.closeUpload(u.Size)

	// Note: Cache was already updated in fsyncNode() by copying from the writer's temp file
	// before it was closed. No need to re-download from remote.

	fs.log.Info("Upload completed: %v (%v).", remotePath, localPath)
	fs.log.Debug("RemoteFs: Bytes: %v, Duration: %v, fh: %v", u.Size, time.Since(start), fh)
}

// this is a convenience method for uploading a file from the local file system to the remote API
// for use by the Rename operation when moving a file from the LocalFs to the RemoteFs.
func (fs *RemoteFs) uploadFile(src, dst string) error {
	fs.log.Debug("Uploading file: %v to %v", src, dst)

	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get file info for size and mtime
	fileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create a cache writer adapter with TeeReader
	var offset int64
	cacheWriter := &cacheWriterAdapter{
		writer: func(data []byte, off int64) (int, error) {
			return fs.cacheStore.Write(dst, data, off)
		},
		offset: &offset,
	}
	teeReader := io.TeeReader(srcFile, cacheWriter)

	// Upload with the TeeReader (writes to cache as it reads)
	err = fs.ops.WithLimit(context.Background(), lim.FuseOpUpload, func(ctx context.Context) error {
		return fs.fileClient.Upload(
			file.UploadWithReader(teeReader),
			file.UploadWithDestinationPath(dst),
		)
	})
	if err != nil {
		return err
	}

	// Update the node's info with the uploaded file details
	node := fs.vfs.getOrCreate(dst, nodeTypeFile)
	node.updateInfo(fsNodeInfo{
		nodeType:     nodeTypeFile,
		size:         fileInfo.Size(),
		modTime:      fileInfo.ModTime(),
		creationTime: node.info.creationTime,
	})

	return nil
}

// this is a convenience method for downloading a file from the remote API to the local file system
// for use by the Rename operation when moving a file from the RemoteFs to the LocalFs.
func (fs *RemoteFs) downloadFile(src, dst string) error {
	fs.log.Debug("RemoteFs: Downloading file: %v to %v", src, dst)
	err := fs.ops.WithLimit(context.Background(), lim.FuseOpDownload, func(ctx context.Context) error {
		_, err := fs.fileClient.DownloadToFile(files_sdk.FileDownloadParams{Path: src}, dst)
		return err
	})

	return err
}

func (fs *RemoteFs) findOrCreateGate(path string) (*cache.ReadyGate, bool) {
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

func (fs *RemoteFs) removeGate(path string, s *cache.ReadyGate) {
	fs.gatesMu.Lock()
	if cur, ok := fs.readyGates[path]; ok && cur == s {
		delete(fs.readyGates, path)
	}
	fs.gatesMu.Unlock()
}

func (fs *RemoteFs) fillCache(ctx context.Context, path string, uri string, readyGate *cache.ReadyGate, fh uint64) {
	_, cancel := context.WithCancel(ctx)
	readyGate.SetCancel(cancel)
	defer func() {
		cancel()
		fs.removeGate(path, readyGate)
	}()

	var f files_sdk.File
	var err error
	err = fs.ops.WithLimit(ctx, lim.FuseOpDownload, func(ctx context.Context) error {
		f, err = fs.fileClient.Download(
			files_sdk.FileDownloadParams{File: files_sdk.File{Path: fs.remotePath(path), DownloadUri: uri}},
			files_sdk.WithContext(ctx),
			files_sdk.ResponseOption(func(resp *http.Response) error {
				if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
					return files_sdk.APIError()(resp)
				}
				defer resp.Body.Close()

				// while downloading a file from the remote API, write data to the disk cache in chunks and update
				// the ready gate every cacheWriteSize bytes to signal that data is available for reading.
				buf := make([]byte, cacheWriteSize)
				var off int64 = 0
				for {
					nr, er := resp.Body.Read(buf)
					if nr > 0 {
						// TODO: consider altering Write to keep data in memory and periodically flush to disk
						// to reduce the number of disk writes. This would require more memory usage, but would
						// improve read and write performance by avoiding constantly opening/closing the file.
						written, err := fs.cacheStore.Write(path, buf[:nr], off)
						if err != nil || written != nr {
							// there was an error writing to the disk cache, or not all bytes that were read from the
							// remote API were written to the disk cache.
							readyGate.Finish(fmt.Errorf("error writing to disk cache for %v: %v", path, err), off)
							return err
						}
						off += int64(written)
						readyGate.SetAvailable(off)
					}
					if er != nil {
						if er == io.EOF {
							readyGate.Finish(nil, off)
							return nil
						}
						readyGate.Finish(er, off)
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
		return errc
	}

	node.lockMutex.Lock()
	defer node.lockMutex.Unlock()

	fs.lockMapMutex.Lock()
	defer fs.lockMapMutex.Unlock()

	localPath, remotePath := fs.paths(node.path)
	fs.log.Debug("RemoteFs: lock: file %v (%v) fh=%v", remotePath, localPath, fh)

	if node.isLocked() {
		fs.log.Error("File is already locked by %v: %v (%v) fh=%v", node.info.lockOwner, remotePath, localPath, fh)
		errc = -fuse.ENOLCK
		return errc
	}

	var lock files_sdk.Lock
	var err error
	err = fs.ops.WithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		lock, err = fs.lockClient.Create(files_sdk.LockCreateParams{
			Path:                 remotePath,
			AllowAccessByAnyUser: lib.Ptr(true),
			Exclusive:            lib.Ptr(true),
			Recursive:            lib.Ptr(true),
			Timeout:              fileLockSeconds,
		})
		return err
	})

	if files_sdk.IsExist(err) {
		// the file is already locked, if it's in the lock map and not owned by this user, return ENOLCK
		linfo, ok := fs.lockMap[node.path]
		if ok && fs.currentUserId != linfo.Lock.UserId {
			fs.log.Error("File '%v' is already locked by %v:", remotePath, linfo.Lock.Username)
			return -fuse.ENOLCK
		}
		if ok && fs.currentUserId == linfo.Lock.UserId {
			// If the lock is already held by the current user, treat it as a success.
			fs.log.Debug("RemoteFs: lock: File is already locked by current user %v: %v (%v) fh=%v", fs.currentUserId, remotePath, localPath, fh)
			return 0
		}
	}

	if errc = fs.handleError(node.path, err); errc != 0 {
		return errc
	}
	fs.log.Debug("RemoteFs: lock: created owner=%v, path=%v, fh=%v", lock.Username, remotePath, fh)
	fs.lockMap[node.path] = &lockInfo{Fh: fh, Lock: &lock}
	return errc
}

func (fs *RemoteFs) unlock(path string, fh uint64) (errc int) {
	if fs.disableLocking {
		return errc
	}

	// If the node exists, prevent locking while unlocking.
	// If the node was renamed/moved, it may still need to be unlocked.
	if node, ok := fs.vfs.fetch(path); ok {
		node.lockMutex.Lock()
		defer node.lockMutex.Unlock()
	}

	fs.lockMapMutex.Lock()
	defer fs.lockMapMutex.Unlock()

	lockInfo, ok := fs.lockMap[path]
	if !ok {
		// If the lock map doesn't have an entry for this path, it means the file
		// was never locked, or it was locked by a different file handle.
		fs.log.Debug("RemoteFs: unlock: File not locked: %v fh=%v", path, fh)
		return errc
	}
	if lockInfo.Fh != fh {
		// This is fine. It just means the file either wasn't locked or it was locked by a different file handle.
		fs.log.Debug("RemoteFs: unlock: File not locked by this handle: %v fh=%v", path, fh)
		return errc
	}

	localPath, remotePath := fs.paths(path)
	fs.log.Debug("RemoteFs: unlock: file %v (%v)", remotePath, localPath)

	err := fs.ops.WithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		err := fs.lockClient.Delete(files_sdk.LockDeleteParams{
			Path:  remotePath,
			Token: lockInfo.Lock.Token,
		})
		return err
	})

	if files_sdk.IsNotExist(err) {
		// If the lock was already deleted, consider it a success.
		fs.log.Debug("RemoteFs: unlock: %v (%v) err=%v", remotePath, localPath, err)
		delete(fs.lockMap, path)
		return errc
	}
	// for any other error, handle it normally
	if errc = fs.handleError(path, err); errc != 0 {
		return errc
	}

	delete(fs.lockMap, path)
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
		localPath, remotePath := fs.paths(path)
		fs.log.Error("%v (%v): %v", remotePath, localPath, err)

		if files_sdk.IsNotAuthenticated(err) {
			fs.events.Publish(events.AuthenticationFailedEvent{
				Reason: err.Error(),
			})
			return -fuse.EPERM
		}
		if files_sdk.IsNotExist(err) {
			return -fuse.ENOENT
		}
		if files_sdk.IsExist(err) {
			return -fuse.EEXIST
		}
		if isFolderNotEmpty(err) {
			return -fuse.ENOTEMPTY
		}
		return -fuse.EIO
	}
	return 0
}

func (fs *RemoteFs) delete(path string) (errc int) {
	err := fs.ops.WithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		err := fs.fileClient.Delete(files_sdk.FileDeleteParams{Path: fs.remotePath(path)})
		return err
	})

	// if there's an error, and it's a not-found, consider it a success.
	if files_sdk.IsNotExist(err) {
		// if the delete was successful, remove the node from the vfs
		fs.vfs.remove(path)
		_ = fs.cacheStore.Delete(path)
		return errc
	}
	// for any other error, handle it normally
	if errc = fs.handleError(path, err); errc != 0 {
		return errc
	}
	// if the delete was successful, remove the node from the vfs
	fs.vfs.remove(path)
	_ = fs.cacheStore.Delete(path)
	return errc
}

func (fs *RemoteFs) loadParent(path string) (errc int) {
	if path == "/" {
		// If loading at the root, the parent can't be loaded. Just make sure the root exists.
		_, errc = fs.findDir(path)
		return errc
	}

	parentPath := path_lib.Dir(path)
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
	err = fs.ops.WithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		item, err = fs.fileClient.Find(files_sdk.FileFindParams{Path: remotePath})
		return err
	})

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
	fs.loadDirMutex.Lock()
	defer fs.loadDirMutex.Unlock()
	if node.infoExpired() {
		fs.log.Debug("RemoteFs: loadDir: Refreshing directory listing: %v", node.path)
		err := node.updateChildPaths(fs.listDir)
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

	opErr = fs.ops.WithLimit(context.Background(), lim.FuseOpOther, func(ctx context.Context) error {
		it, err := fs.fileClient.ListFor(files_sdk.FolderListForParams{Path: fs.remotePath(path)})
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

		locks, err := fs.lockClient.ListFor(files_sdk.LockListForParams{
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

func (fs *RemoteFs) createNode(path string, item files_sdk.File) *fsNode {
	var nt nodeType
	if item.IsDir() {
		nt = nodeTypeDir
	} else {
		nt = nodeTypeFile
	}
	// best-effort invalidate stale data
	if prev, ok := fs.vfs.fetch(path); ok && prev.info.nodeType == nodeTypeFile {
		if prev.info.size != item.Size || !prev.info.modTime.Equal(item.ModTime()) {
			_ = fs.cacheStore.Delete(path)
		}
	}

	node := fs.vfs.getOrCreate(path, nt)
	node.updateInfo(fsNodeInfo{
		nodeType:     nt,
		size:         item.Size,
		modTime:      item.ModTime(),
		creationTime: item.CreationTime(),
	})

	return node
}

func (fs *RemoteFs) waitForAction(ctx context.Context, action files_sdk.FileAction, operation string) error {
	var migration files_sdk.FileMigration
	var err error
	err = fs.ops.WithLimit(ctx, lim.FuseOpOther, func(ctx context.Context) error {
		migration, err = fs.migrationClient.Wait(action, func(migration files_sdk.FileMigration) {
			fs.log.Trace("RemoteFs: watchForAction: waiting for migration")
		})
		return err
	})

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
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *RemoteFs) Chown(path string, uid uint32, gid uint32) int {
	fs.log.Trace("RemoteFs: Chown: path=%v, uid=%v, gid=%v", path, uid, gid)
	return -fuse.ENOSYS
}

// Access checks file access permissions.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *RemoteFs) Access(path string, mask uint32) int {
	fs.log.Trace("RemoteFs: Access: path=%v, mask=%v", path, mask)
	return -fuse.ENOSYS
}

// Flush flushes cached file data.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *RemoteFs) Flush(path string, fh uint64) int {
	fs.log.Trace("RemoteFs: Flush: path=%v, fh=%v", path, fh)
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
	errc, fh := fs.Create(path, fi.Flags, mode)
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
