package fsmount

import (
	"context"
	"errors"
	"io"
	iofs "io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	ff "github.com/Files-com/files-sdk-go/v3/fsmount/internal/flags"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/winfsp/cgofuse/fuse"
)

// LocalFs is a used as a passthrough files system to the host operating system's file system for
// the local mount point. It is used for files that are considered temporary or that should not
// be uploaded to Files.com. In order to support programs that use the type of files that are not
// intended to be uploaded to Files.com, the Filescomfs can delegate those operations to this
// implementation.
type LocalFs struct {
	mountPoint      string
	localFsRoot     string
	log             lib.LeveledLogger
	vfs             *virtualfs
	initOnce        sync.Once
	initTime        time.Time
	deletable       chan string
	deleteRetries   map[string]int
	deleteRetriesMu sync.Mutex
	maintMu         sync.Mutex
	maintActive     bool
	maintCancel     context.CancelFunc
	wg              sync.WaitGroup
}

const maxDeleteAttempts = 10

func newLocalFs(params MountParams, vfs *virtualfs, ll lib.LeveledLogger) *LocalFs {
	return &LocalFs{
		mountPoint:    params.MountPoint,
		localFsRoot:   params.TmpFsPath,
		vfs:           vfs,
		log:           ll,
		deletable:     make(chan string, 100),
		deleteRetries: make(map[string]int),
	}
}

func (fs *LocalFs) Init() {
	// Guard with a sync.Once because Init is called from fsmount.Mount, but cgofuse also calls Init
	// when it mounts the file system.
	fs.initOnce.Do(func() {
		// store the time the file system was initialized to use as the creation time for the root directory
		fs.initTime = time.Now()
		fs.log.Debug("LocalFs initialized successfully. Local file system root: %s", fs.localFsRoot)
		fs.StartMaintenance()
	})
}

func (fs *LocalFs) Destroy() {
	root := filepath.Clean(fs.localFsRoot)
	tmp := filepath.Clean(os.TempDir())
	fs.log.Debug("LocalFs: Destroy: considering removal of local file system root: %v", root)

	fs.StopMaintenance()
	// only remove the local file system root if it is under the system temp directory
	if strings.HasPrefix(root+string(os.PathSeparator), tmp+string(os.PathSeparator)) {
		if err := os.RemoveAll(root); err != nil {
			fs.log.Debug("LocalFs: Destroy: failed to remove temporary local file system root: %v, err: %v", root, err)
			return
		}
		fs.log.Debug("LocalFs: Destroy: removed temporary local file system root: %v", root)
		return
	}
	fs.log.Debug("LocalFs: Destroy: refusing to remove non-temp TmpFsPath: %v (TempDir=%v)", root, tmp)
}

func (fs *LocalFs) Validate() error {
	fs.Init()
	return nil
}

func (fs *LocalFs) fqPath(path string) string {
	clean := filepath.Clean(path)
	// If absolute, prefer containment under localFsRoot; otherwise make it relative.
	if filepath.IsAbs(clean) {
		// If it's already under the temp root, allow it.
		if rel, err := filepath.Rel(fs.localFsRoot, clean); err == nil {
			rel = filepath.ToSlash(rel)
			if !strings.HasPrefix(rel, "..") {
				return clean
			}
		}
		// Otherwise, strip the leading separator to force a relative join.
		clean = strings.TrimPrefix(clean, string(os.PathSeparator))
	}
	fq := filepath.Join(fs.localFsRoot, clean)

	// Final containment check (defense-in-depth) to keep operations sandboxed.
	// Normalize both the sandbox root and the candidate path.
	rootClean := filepath.Clean(fs.localFsRoot)
	fqClean := filepath.Clean(fq)

	// Build "path with trailing separator" forms to ensure only whole
	// directory segments are matched (so /tmp/x doesn't match /tmp/xyz).
	fqCleanWithSep := fqClean + string(os.PathSeparator)
	rootCleanWithSep := rootClean + string(os.PathSeparator)

	// Check if the candidate is exactly the root…
	isSameAsRoot := fqClean == rootClean
	// …or a descendant of the root (using the trailing-sep forms to avoid false positives).
	isUnderRoot := strings.HasPrefix(fqCleanWithSep, rootCleanWithSep)

	// If it's neither the root nor under it, bail out to the sandbox root.
	if !(isSameAsRoot || isUnderRoot) {
		return rootClean
	}
	return fq
}

func (fs *LocalFs) Mkdir(path string, mode uint32) (errc int) {
	path = fs.fqPath(path)
	if err := os.MkdirAll(path, os.FileMode(mode)); err != nil {
		fs.log.Debug("LocalFs: Mkdir: failed to create directory: path=%v, mode=%o, err=%v", path, mode, err)
		return -fuse.EIO
	}
	return errc
}

func (fs *LocalFs) Unlink(path string) (errc int) {
	// defer deletion to a background goroutine to avoid blocking the FUSE operation
	// this entire file system will be cleaned up on unmount, so best effort deletion is sufficient
	path = fs.fqPath(path)
	fs.vfs.remove(path)
	fs.deletable <- fs.fqPath(path)
	return errc
}

// StartMaintenance starts the background deletion goroutine if it is not already running.
func (fs *LocalFs) StartMaintenance() {
	fs.maintMu.Lock()
	if fs.maintActive {
		fs.maintMu.Unlock()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	fs.maintCancel = cancel
	fs.maintActive = true
	fs.maintMu.Unlock()
	fs.wg.Add(1)
	go func() {
		defer fs.wg.Done()
		fs.maintenanceLoop(ctx)
	}()
}

// StopMaintenance stops the background deletion goroutine if it is running.
func (fs *LocalFs) StopMaintenance() {
	fs.maintMu.Lock()
	if !fs.maintActive {
		fs.maintMu.Unlock()
		return
	}
	cancel := fs.maintCancel
	fs.maintMu.Unlock()
	// Trigger shutdown and wait for the goroutine to exit cleanly.
	if cancel != nil {
		cancel()
	}
	fs.wg.Wait()

	fs.maintMu.Lock()
	// Clear struct state under the lock so a concurrent Start can proceed.
	fs.maintCancel = nil
	fs.maintActive = false
	fs.maintMu.Unlock()
}

// maintenanceLoop runs periodic maintenance tasks until the context is cancelled.
func (fs *LocalFs) maintenanceLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case path := <-fs.deletable:
			fs.handleDeletion(ctx, path)
		case <-ctx.Done():
			return
		case <-ticker.C:
			// periodic wake-up to check context
		}
	}
}

// handleDeletion attempts to delete a path and schedules a retry if needed.
func (fs *LocalFs) handleDeletion(ctx context.Context, path string) {
	fs.deleteRetriesMu.Lock()
	defer fs.deleteRetriesMu.Unlock()

	err := os.RemoveAll(path)
	if err != nil && !errors.Is(err, iofs.ErrNotExist) {
		retries := fs.deleteRetries[path]
		retries++
		if retries < maxDeleteAttempts {
			fs.deleteRetries[path] = retries
			fs.log.Debug("LocalFs: handleDeletion: failed to remove path (attempt %d/%d): path=%v, err=%v", retries, maxDeleteAttempts, path, err)
			fs.scheduleRetry(ctx, path)
		} else {
			delete(fs.deleteRetries, path)
			fs.log.Debug("LocalFs: handleDeletion: giving up on path after %d attempts: path=%v, err=%v", retries, path, err)
		}
	} else {
		delete(fs.deleteRetries, path)
		fs.log.Trace("LocalFs: handleDeletion: removed path: path=%v", path)
	}
}

// scheduleRetry re-queues a path for deletion after a delay, respecting context cancellation.
func (fs *LocalFs) scheduleRetry(ctx context.Context, path string) {
	go func() {
		select {
		case <-time.After(time.Second):
			select {
			case fs.deletable <- path:
			case <-ctx.Done():
				// Context cancelled while trying to send
			}
		case <-ctx.Done():
			// Context cancelled during delay
		}
	}()
}

func (fs *LocalFs) Rmdir(path string) (errc int) {
	path = fs.fqPath(path)
	if err := os.RemoveAll(path); err != nil {
		fs.log.Debug("LocalFs: Rmdir: failed to remove directory: path=%v, err=%v", path, err)
		return -fuse.EIO
	}
	fs.vfs.remove(path)
	return errc
}

func (fs *LocalFs) Rename(oldpath string, newpath string) (errc int) {
	oldpath = fs.fqPath(oldpath)
	newpath = fs.fqPath(newpath)
	if err := os.Rename(oldpath, newpath); err != nil {
		fs.log.Debug("LocalFs: Rename: failed to rename file: oldpath=%v, newpath=%v, err=%v", oldpath, newpath, err)
		return -fuse.EIO
	}
	fs.rename(oldpath, newpath)
	return errc
}

func (fs *LocalFs) rename(oldpath, newpath string) {
	// handles moving the node and fixing parent listings.
	node := fs.vfs.rename(oldpath, newpath)

	// clear any cached presigned URL for this node (path changed).
	if node != nil {
		node.clearDownloadURI()
	}
}

func (fs *LocalFs) Utimens(path string, tmsp []fuse.Timespec) (errc int) {
	path = fs.fqPath(path)
	if len(tmsp) != 2 {
		fs.log.Debug("LocalFs: Utimens: invalid number of timespecs provided: path=%v, tmsp=%v", path, tmsp)
		return -fuse.EINVAL
	}
	if err := os.Chtimes(path, tmsp[0].Time(), tmsp[1].Time()); err != nil {
		fs.log.Debug("LocalFs: Utimens: failed to change times for path: path=%v, tmsp=%v, err=%v", path, tmsp, err)
		return -fuse.EIO
	}
	return errc
}

func (fs *LocalFs) Create(path string, flags int, mode uint32) (errc int, fh uint64) {
	path = fs.fqPath(path)
	fuseFlags := ff.NewFuseFlags(flags)
	fs.log.Trace("LocalFs: Create: path=%v, flags=%v, mode=%o", path, fuseFlags, mode)
	return fs.open(path, flags, mode)
}

func (fs *LocalFs) Open(path string, flags int) (errc int, fh uint64) {
	path = fs.fqPath(path)
	fuseFlags := ff.NewFuseFlags(flags)
	fs.log.Trace("LocalFs: Open: path=%v, flags=%v", path, fuseFlags)
	return fs.open(path, flags, 0)
}

func (fs *LocalFs) open(path string, flags int, mode uint32) (errc int, fh uint64) {
	dpath := filepath.Dir(path)
	if err := os.MkdirAll(dpath, 0o755); err != nil {
		fs.log.Debug("LocalFs: open: failed to create parent directories: path=%v, flags=%v, mode=%o, err=%v", path, flags, mode, err)
		return -fuse.EIO, ^uint64(0)
	}
	fuseFlags := ff.NewFuseFlags(flags)
	f, err := os.OpenFile(path, fuseFlags.AsOsFlags(), os.FileMode(mode))
	if err != nil {
		// this is expected in some cases, like .DS_Store files on macOS, so log at Debug level
		errc = toErrno(err)
		fs.log.Debug("LocalFs: open: failed to open file: path=%v, flags=%v, mode=%o, err=%v, errno=%v", path, fuseFlags, mode, err, errc)
		return errc, ^uint64(0)
	}
	// create a new fsNode for the file and open a file handle for it in the virtual file system
	node := fs.vfs.getOrCreate(path, nodeTypeFile)
	fh, _ = fs.vfs.handles.OpenWithFile(node, fuseFlags, f)
	fs.log.Trace("LocalFs: open: succeeded path=%v, flags=%v, mode=%o fh=%v", path, fuseFlags, mode, fh)
	return errc, fh
}

func toErrno(err error) int {
	if err == nil {
		return 0
	}
	switch {
	case errors.Is(err, iofs.ErrInvalid):
		return -fuse.EINVAL
	case errors.Is(err, iofs.ErrPermission):
		return -fuse.EPERM
	case errors.Is(err, iofs.ErrExist):
		return -fuse.EEXIST
	case errors.Is(err, iofs.ErrNotExist):
		return -fuse.ENOENT
	case errors.Is(err, iofs.ErrClosed):
		return -fuse.EBADF
	default:
		return -fuse.EIO
	}
}

func (fs *LocalFs) Truncate(path string, size int64, fh uint64) (errc int) {
	path = fs.fqPath(path)
	if err := os.Truncate(path, size); err != nil {
		fs.log.Debug("LocalFs: Truncate: failed to truncate file: path=%v, size=%v, fh=%v, err=%v", path, size, fh, err)
		return -fuse.EIO
	}
	return errc
}

func (fs *LocalFs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	path = fs.fqPath(path)
	handle, _, ok := fs.vfs.handles.Lookup(fh)
	if !ok {
		fs.log.Debug("LocalFs: Read: invalid file handle: path=%v, ofst=%v, fh=%v", path, ofst, fh)
		return -fuse.EBADF
	}
	n, err := handle.localFile.ReadAt(buff, ofst)
	if err != nil && !errors.Is(err, io.EOF) {
		fs.log.Debug("LocalFs: Read: failed to read file: path=%v, ofst=%v, fh=%v, err=%v", path, ofst, fh, err)
		return -fuse.EIO
	}
	handle.incrementRead(int64(n))
	return n
}

func (fs *LocalFs) Write(path string, buff []byte, ofst int64, fh uint64) (n int) {
	path = fs.fqPath(path)
	handle, _, ok := fs.vfs.handles.Lookup(fh)
	if !ok {
		fs.log.Debug("LocalFs: Write: invalid file handle: path=%v, ofst=%v, fh=%v", path, ofst, fh)
		return -fuse.EBADF
	}
	n, err := handle.localFile.WriteAt(buff, ofst)
	if err != nil {
		fs.log.Debug("LocalFs: Write: failed to write file: path=%v, ofst=%v, fh=%v, err=%v", path, ofst, fh, err)
		return -fuse.EIO
	}
	fs.log.Trace("LocalFs: Write: path=%v, len=%v, ofst=%v, fh=%v", path, len(buff), ofst, fh)
	return n
}

func (fs *LocalFs) Release(path string, fh uint64) (errc int) {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Release: path=%v, fh=%v", path, fh)
	handle, _, ok := fs.vfs.handles.Lookup(fh)
	if !ok {
		fs.log.Debug("LocalFs: Release: invalid file handle: path=%v, fh=%v", path, fh)
		return -fuse.EBADF
	}
	defer fs.vfs.handles.Release(fh)

	// This should always be non-nil since it's set when opening the handle,
	// but check for nil before attempting to close it.
	if handle.localFile != nil {
		if err := handle.localFile.Close(); err != nil {
			fs.log.Debug("LocalFs: Release: failed to close file: path=%v, fh=%v, err=%v", path, fh, err)
			return -fuse.EIO
		}
	}
	fs.log.Trace("LocalFs: Release: succeeded path=%v, fh=%v", path, fh)
	return errc
}

func (fs *LocalFs) Opendir(path string) (errc int, fh uint64) {
	// This is largely a no-op since directory handles aren't needed on the local file system,
	// but allocating a file handle allows tracking of open directories in the virtual file system
	// and aligns with releasedir which will close the handle.
	path = fs.fqPath(path)
	node := fs.vfs.getOrCreate(path, nodeTypeDir)
	fh, _ = fs.vfs.handles.Open(node, ff.NewFuseFlags(fuse.O_RDONLY))
	return 0, fh
}

func (fs *LocalFs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {

	path = fs.fqPath(path)
	entries, err := os.ReadDir(path)
	if err != nil {
		fs.log.Debug("LocalFs: Readdir: failed to open directory: path=%v, fh=%v, err=%v", path, fh, err)
		return -fuse.EIO
	}

	// TODO: see if these to lookups agree
	// handle, node, ok := fs.vfs.handles.Lookup(fh)
	// fillNode, fok := fs.vfs.fetch(path)
	fill(".", nil, 0)
	fill("..", nil, 0)

	// no need to sort the entries since os.ReadDir returns them in sorted order
	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		node, err := fs.createLocalNode(entryPath, entry)
		if err != nil {
			fs.log.Debug("LocalFs: Readdir: failed to create local node for entry: path=%v, entry=%v, err=%v", path, entryPath, err)
			continue
		}
		fs.log.Trace("LocalFs: Readdir: Calling fill for entry: %v", entryPath)
		fill(entry.Name(), getStat(node.info, nil), 0)
	}
	return 0
}

func (fs *LocalFs) Releasedir(path string, fh uint64) (errc int) {
	_, _ = fs.vfs.handles.Release(fh)
	return errc
}

// Chmod changes the permission bits of a file.
func (fs *LocalFs) Chmod(path string, mode uint32) int {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Chmod: path=%v, mode=%o", path, mode)
	if err := os.Chmod(path, os.FileMode(mode)); err != nil {
		fs.log.Debug("LocalFs: Chmod: failed to change mode for path: path=%v, mode=%o, err=%v", path, mode, err)
		return -fuse.EIO
	}
	return 0
}

// Fsync attempts to synchronize file contents.
func (fs *LocalFs) Fsync(path string, datasync bool, fh uint64) (errc int) {
	path = fs.fqPath(path)
	handle, _, ok := fs.vfs.handles.Lookup(fh)
	if !ok {
		fs.log.Debug("LocalFs: Fsync: invalid file handle: path=%v, datasync=%v, fh=%v", path, datasync, fh)
		return -fuse.EBADF
	}
	if handle.localFile != nil {
		if err := handle.localFile.Sync(); err != nil {
			fs.log.Debug("LocalFs: Fsync: failed to sync file: path=%v, datasync=%v, fh=%v, err=%v", path, datasync, fh, err)
			return -fuse.EIO
		}
	} else {
		fs.log.Debug("LocalFs: Fsync: no local file to sync: path=%v, datasync=%v, fh=%v", path, datasync, fh)
		return -fuse.EIO
	}

	fs.log.Trace("LocalFs: Fsync: path=%v, datasync=%v, fh=%v", path, datasync, fh)
	return errc
}

// Methods below are part of the fuse.FileSystemInterface, but not supported by
// this implementation. They exist here to support logging for visibility of how
// the underlying fuse layer calls into this implementation.

// Mknod creates a file node.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *LocalFs) Mknod(path string, mode uint32, dev uint64) int {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Mknod: path=%v, mode=%o, dev=%v", path, mode, dev)
	return -fuse.ENOSYS
}

// Link creates a hard link to a file.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *LocalFs) Link(oldpath string, newpath string) int {
	fs.log.Trace("LocalFs: Link: old=%v, new=%v", oldpath, newpath)
	return -fuse.ENOSYS
}

// Symlink creates a symbolic link.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *LocalFs) Symlink(target string, newpath string) int {
	fs.log.Trace("LocalFs: Symlink: target=%v, newpath=%v", target, newpath)
	return -fuse.ENOSYS
}

// Readlink reads the target of a symbolic link.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *LocalFs) Readlink(path string) (int, string) {
	fs.log.Trace("LocalFs: Readlink: path=%v", path)
	return -fuse.ENOSYS, ""
}

// Chown changes the owner and group of a file.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *LocalFs) Chown(path string, uid uint32, gid uint32) int {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Chown: path=%v, uid=%v, gid=%v", path, uid, gid)
	return -fuse.ENOSYS
}

// Access checks file access permissions.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *LocalFs) Access(path string, mask uint32) int {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Access: path=%v, mask=%v", path, mask)
	return -fuse.ENOSYS
}

// Flush flushes cached file data.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *LocalFs) Flush(path string, fh uint64) int {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Flush: path=%v, fh=%v", path, fh)
	return 0
}

// Fsyncdir synchronizes directory contents.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *LocalFs) Fsyncdir(path string, datasync bool, fh uint64) int {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Fsyncdir: path=%v, datasync=%v, fh=%v", path, datasync, fh)
	return -fuse.ENOSYS
}

// The [Foo]xattr implementations below explicitly return 0 to indicate that
// extended attributes are "supported" in order to ensure that the other xattr
// methods are called for debugging visibility, but are all no-op implementations.

// Getxattr gets extended attributes.
// Any return value other than -fuse.ENOSYS indicates support for extended
// attributes, but also expects Setxattr, Listxattr, and Removexattr to exist
// for extended attribute support.
func (fs *LocalFs) Getxattr(path string, name string) (int, []byte) {
	fs.log.Trace("LocalFs: Getxattr: path=%v, name=%v", path, name)
	return 0, []byte{}
}

// Setxattr sets extended attributes.
func (fs *LocalFs) Setxattr(path string, name string, value []byte, flags int) int {
	path = fs.fqPath(path)
	fuseFlags := ff.NewFuseFlags(flags)
	fs.log.Trace("LocalFs: Setxattr: path=%v, name=%v, value=%v flags=%v", path, name, value, fuseFlags)
	return 0
}

// Removexattr removes extended attributes.
func (fs *LocalFs) Removexattr(path string, name string) int {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Removexattr: path=%v, name=%v", path, name)
	return 0
}

// Listxattr lists extended attributes.
func (fs *LocalFs) Listxattr(path string, fill func(name string) bool) int {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Listxattr: path=%v", path)
	return 0
}

// FileSystemOpenEx is the interface that wraps the OpenEx and CreateEx methods.

// OpenEx and CreateEx are similar to Open and Create except that they allow
// direct manipulation of the FileInfo_t struct (which is analogous to the
// FUSE struct fuse_file_info). If implemented, they are preferred over
// Open and Create.
func (fs *LocalFs) CreateEx(path string, mode uint32, fi *fuse.FileInfo_t) int {
	fs.log.Trace("LocalFs: CreateEx: path=%v, mode=%o, fi=%v", path, mode, fi)
	errc, fh := fs.Create(path, fi.Flags, mode)
	fi.Fh = fh
	return errc
}

func (fs *LocalFs) OpenEx(path string, fi *fuse.FileInfo_t) int {
	fs.log.Trace("LocalFs: OpenEx: path=%v, fi=%v", path, fi)
	errc, fh := fs.Open(path, fi.Flags)
	fi.Fh = fh
	return errc
}

// Getpath is part of the FileSystemGetpath interface and
// allows a case-insensitive file system to report the correct case of a file path.
func (fs *LocalFs) Getpath(path string, fh uint64) (int, string) {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Getpath: path=%v, fh=%v", path, fh)
	return -fuse.ENOSYS, path
}

// Chflags is part of the FileSystemChflags interface and
// changes the BSD file flags (Windows file attributes).
func (fs *LocalFs) Chflags(path string, flags uint32) int {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Chflags: path=%v, flags=%v", path, flags)
	return -fuse.ENOSYS
}

// Setcrtime is part of the FileSystemSetcrtime interface and
// changes the file creation (birth) time.
func (fs *LocalFs) Setcrtime(path string, tmsp fuse.Timespec) int {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Setcrtime: path=%v, tmsp=%v", path, tmsp)
	return -fuse.ENOSYS
}

// Setchgtime is part of the FileSystemSetchgtime interface and
// changes the file change (ctime) time.
func (fs *LocalFs) Setchgtime(path string, tmsp fuse.Timespec) int {
	path = fs.fqPath(path)
	fs.log.Trace("LocalFs: Setchgtime: path=%v, tmsp=%v", path, tmsp)
	if err := os.Chtimes(path, time.Now(), tmsp.Time()); err != nil {
		fs.log.Debug("LocalFs: Setchgtime: failed to change change time for path: path=%v, tmsp=%v, err=%v", path, tmsp, err)
		return -fuse.EIO
	}
	return 0
}
