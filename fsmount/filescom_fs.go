package fsmount

import (
	"context"
	"errors"
	"os"
	"runtime"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/fsmount/events"
	"github.com/Files-com/files-sdk-go/v3/lib"
	gogitignore "github.com/sabhiram/go-gitignore"
	"github.com/winfsp/cgofuse/fuse"
)

const (
	folderNotEmpty = "processing-failure/folder-not-empty"
	blockSize      = 4096

	// The maximum time to wait for the Fsync operation to complete before timing out.
	// This is used to prevent the Fsync operation from hanging indefinitely if the remote API is unresponsive.
	fsyncTimeout = 30 * time.Second
)

var (
	// This is the length of time used when establishing a lock on a file.
	// This is used to prevent the lock from being held indefinitely if the lock is not released properly.
	fileLockSeconds = int64((1 * time.Hour).Seconds())
)

// Filescomfs is a file system that implements the fuse.FileSystem interface,
// allowing it to be mounted using FUSE. It provides a virtual files system
// interface to Files.com, allowing users to interact with their Files.com
// account as if it were a local file system.
type Filescomfs struct {
	remote         *RemoteFs
	local          *LocalFs
	vfs            *virtualfs
	log            lib.LeveledLogger
	mountPoint     string
	remoteRoot     string
	localFsRoot    string
	disableLocking bool
	ignore         *gogitignore.GitIgnore
	events         events.EventPublisher
	initOnce       sync.Once
}

// Init initializes the Filescomfs file system.
func (fs *Filescomfs) Init() {
	defer logPanics(fs.log)
	// Guard with a sync.Once because Init is called from fsmount.Mount, but cgofuse also calls Init
	// when it mounts the file system.
	fs.initOnce.Do(func() {
		fs.remote.Init()
		fs.local.Init()

		fs.log.Info("Files.com file system initialized successfully.")
		fs.log.Debug("Mount point: %s, Remote file system root: %s, Local file system root: %s", fs.mountPoint, fs.remote.root, fs.local.localFsRoot)
	})
}

func (fs *Filescomfs) Destroy() {
	defer logPanics(fs.log)
	fs.log.Info("Shutting down Files.com file system mounted at: %s", fs.mountPoint)
	fs.remote.Destroy()
	fs.local.Destroy()
	fs.vfs.destroy()
}

// Validate checks if the Filescomfs file system is valid by attempting to list the root directories of the RemoteFs and LocalFs.
func (fs *Filescomfs) Validate() error {
	defer logPanics(fs.log)
	fs.Init()
	lerr := fs.local.Validate()
	rerr := fs.remote.Validate()
	return errors.Join(lerr, rerr)
}

func (fs *Filescomfs) Statfs(path string, stat *fuse.Statfs_t) (errc int) {
	defer logPanics(fs.log)
	fs.log.Trace("Statfs: path=%v", path)

	totalBytes := remoteCapacityBytes()

	// TODO: get used bytes from the remote
	usedBytes := uint64(0)
	freeBytes := totalBytes - usedBytes

	stat.Bsize = blockSize
	stat.Frsize = blockSize
	stat.Blocks = totalBytes / blockSize
	stat.Bfree = freeBytes / blockSize
	stat.Bavail = freeBytes / blockSize

	return errc
}

func remoteCapacityBytes() uint64 {
	// the remote capacity is functionally unlimited, so return the largest
	// value that the OS will accept
	switch runtime.GOOS {
	case "darwin":
		// ~8TB - any larger and the drive shows up as zero capacity on macOS
		return uint64(1 << 43)
	default:
		// ~1PB
		return uint64(1 << 50)
	}
}

func (fs *Filescomfs) Mkdir(path string, mode uint32) (errc int) {
	defer logPanics(fs.log)
	// determine if the directory should be created locally or remotely based on the ignore patterns
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Mkdir(path, mode)
		fs.log.Trace("Filescomfs: Mkdir: creating directory remotely: path=%v, mode=%v, errc=%v", path, mode, errc)
		return errc
	}
	errc = fs.local.Mkdir(path, mode)
	fs.log.Trace("Filescomfs: Mkdir: creating directory locally: path=%v, mode=%v, errc=%v", path, mode, errc)
	return errc
}

func (fs *Filescomfs) Unlink(path string) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Unlink(path)
		fs.log.Trace("Filescomfs: Unlink: deleting file remotely: path=%v, errc=%v", path, errc)
		return errc
	}
	errc = fs.local.Unlink(path)
	fs.log.Trace("Filescomfs: Unlink: deleting file locally: path=%v, errc=%v", path, errc)
	return errc
}

func (fs *Filescomfs) Rmdir(path string) int {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc := fs.remote.Rmdir(path)
		fs.log.Trace("Filescomfs: Rmdir: removing directory remotely: path=%v, errc=%v", path, errc)
		return errc
	}
	errc := fs.local.Rmdir(path)
	fs.log.Trace("Filescomfs: Rmdir: removing directory locally: path=%v, errc=%v", path, errc)
	return errc
}

func (fs *Filescomfs) Rename(oldpath string, newpath string) (errc int) {
	defer logPanics(fs.log)
	fs.log.Trace("Filescomfs: Rename: renaming file: oldpath=%v, newpath=%v", oldpath, newpath)

	// for renames that stay within the same storage (local to local, remote to remote)
	// delegate to the appropriate file system's Rename method
	if fs.isStoredRemotely(oldpath) && fs.isStoredRemotely(newpath) {
		errc = fs.remote.Rename(oldpath, newpath)
		fs.log.Trace("Filescomfs: Rename: renaming file remotely: oldpath=%v, newpath=%v, errc=%v", oldpath, newpath, errc)
		return errc
	}
	if fs.isStoredLocally(oldpath) && fs.isStoredLocally(newpath) {
		errc = fs.local.Rename(oldpath, newpath)
		fs.log.Trace("Filescomfs: Rename: renaming file locally: oldpath=%v, newpath=%v, errc=%v", oldpath, newpath, errc)
		return errc
	}

	// for renames that cross storage boundaries (local to remote, remote to local)
	// perform the necessary upload/download and then delete the source file
	if fs.isStoredLocally(oldpath) && fs.isStoredRemotely(newpath) {
		fs.log.Trace("Filescomfs: Rename: renaming file from local to remote: oldpath=%v, newpath=%v", oldpath, newpath)
		// oldpath = /var/folders/xx/xx/T/filescomfs-xxxxxx/filename.txt
		// newpath = /remote/path/filename.txt
		oldFq := fs.local.fqPath(oldpath)
		if err := fs.remote.uploadFile(oldFq, newpath); err != nil {
			errc = -fuse.EIO
			fs.log.Error("Filescomfs: Rename: error uploading file from local to remote: oldpath=%v, newpath=%v, err=%v, errc=%v", oldpath, newpath, err, errc)
			return errc
		}
		fs.vfs.rename(oldpath, newpath)
		errc = fs.local.Unlink(oldpath)
		fs.log.Trace("Filescomfs: Rename: finished unlink after upload: oldpath=%v, newpath=%v, errc=%v", oldpath, newpath, errc)
		return errc
	}
	if fs.isStoredRemotely(oldpath) && fs.isStoredLocally(newpath) {
		fs.log.Trace("Filescomfs: Rename: renaming file from remote to local: oldpath=%v, newpath=%v", oldpath, newpath)
		// oldpath = /remote/path/filename.txt
		// newpath = /var/folders/xx/xx/T/filescomfs-xxxxxx/filename.txt
		// 1) If an upload is active, wait for it to finalize (bounded).
		if n, ok := fs.vfs.fetch(oldpath); ok {
			// Snapshot: (writer, owner, committed)
			if w, _, committed := n.writerSnapshot(); w != nil {
				// If bytes have been written, wait for finalize; if uncommitted, treat as busy.
				ctx, cancel := context.WithTimeout(context.Background(), fsyncTimeout)
				defer cancel()
				if committed {
					fs.log.Trace("Filescomfs: Rename: waiting for finalize of active upload before renaming: oldpath=%v, newpath=%v", oldpath, newpath)
					if err := n.waitForUploadIfFinalizing(ctx); err != nil {
						errc = -fuse.EAGAIN
						fs.log.Trace("Filescomfs: Rename: wait-for-finalize timed out: %v, errc=%v", err, errc)
						return errc
					}
				} else {
					// nothing has been written to the remote and the OS is trying to rename the
					// file. cancel the upload and proceed. The attempt to download it will 404,
					// and the fs can create an empty file locally to satisfy the rename.
					n.cancelUpload()
					fs.log.Trace("Filescomfs: Rename: canceled uncommitted active upload before renaming: oldpath=%v, newpath=%v", oldpath, newpath)
				}
			} else {
				fs.log.Trace("Filescomfs: Rename: no active upload to wait for: oldpath=%v, newpath=%v", oldpath, newpath)
			}
		} else {
			// TODO: maybe load parent directories to populate the vfs?
			errc = -fuse.ENOENT
			fs.log.Trace("Filescomfs: Rename: file to rename is not in the vfs: oldpath=%v, newpath=%v, errc=%v", oldpath, newpath, errc)
			return errc
		}

		// a best effort has been made to ensure that any active upload has finalized
		// before proceeding with the download
		fqNew := fs.local.fqPath(newpath)
		err := fs.remote.downloadFile(oldpath, fqNew)
		if err != nil {
			// if the file is not found on the remote, create an empty file locally to satisfy the rename
			if files_sdk.IsNotExist(err) {
				fs.log.Debug("Filescomfs: Rename: file not found on remote when downloading during rename: oldpath=%v, newpath=%v, err=%v", oldpath, newpath, err)
				errc = 0 // treat as success because the source would have been deleted by the rename anyway
				if _, err := os.Create(fqNew); err != nil {
					// create an empty file to satisfy the rename
					errc = -fuse.EIO
					fs.log.Error("Filescomfs: Rename: error creating empty file after source not found on remote during rename: oldpath=%v, newpath=%v, err=%v, errc=%v", oldpath, newpath, err, errc)
					return errc
				}
				return errc
			}
			errc = -fuse.EIO
			fs.log.Error("Filescomfs: Rename: error downloading file from remote to local: oldpath=%v, newpath=%v, err=%v, errc=%v", oldpath, newpath, err, errc)
			return errc
		}
		fs.vfs.rename(oldpath, newpath)

		// rename removes the node from the vfs, so calling Unlink will return ENOENT
		// so issue the delete directly
		errc = fs.remote.delete(oldpath)
		fs.log.Trace("Filescomfs: Rename: finished unlink after download: oldpath=%v, newpath=%v, errc=%v", oldpath, newpath, errc)
		return errc
	}
	errc = -fuse.ENOENT
	fs.log.Trace("Filescomfs: Rename: invalid rename operation: oldpath=%v, newpath=%v, errc=%v", oldpath, newpath, errc)
	return errc
}

func (fs *Filescomfs) Utimens(path string, tmsp []fuse.Timespec) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Utimens(path, tmsp)
		fs.log.Trace("Filescomfs: Utimens: updating times remotely: path=%v, tmsp=%v, errc=%v", path, tmsp, errc)
		return errc
	}
	errc = fs.local.Utimens(path, tmsp)
	fs.log.Trace("Filescomfs: Utimens: updating times locally: path=%v, tmsp=%v, errc=%v", path, tmsp, errc)
	return errc
}

func (fs *Filescomfs) Create(path string, flags int, mode uint32) (errc int, fh uint64) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc, fh = fs.remote.Create(path, flags, mode)
		fs.log.Trace("Filescomfs: Create: creating file remotely: path=%v, flags=%v, mode=%v, errc=%v, fh=%v", path, flags, mode, errc, fh)
		return errc, fh
	}
	errc, fh = fs.local.Create(path, flags, mode)
	fs.log.Trace("Filescomfs: Create: creating file locally: path=%v, flags=%v, mode=%v, errc=%v, fh=%v", path, flags, mode, errc, fh)
	return errc, fh
}

func (fs *Filescomfs) Open(path string, flags int) (errc int, fh uint64) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc, fh = fs.remote.Open(path, flags)
		fs.log.Trace("Filescomfs: Open: opening file remotely: path=%v, flags=%v, errc=%v, fh=%v", path, flags, errc, fh)
		return errc, fh
	}
	errc, fh = fs.local.Open(path, flags)
	fs.log.Trace("Filescomfs: Open: opening file locally: path=%v, flags=%v, errc=%v, fh=%v", path, flags, errc, fh)
	return errc, fh
}

func (fs *Filescomfs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		fs.log.Trace("Filescomfs: Getattr: getting attributes remotely: path=%v, fh=%v", path, fh)
		errc = fs.remote.Getattr(path, stat, fh)
		fs.log.Trace("Filescomfs: Getattr: got attributes remotely: path=%v, fh=%v, errc=%v", path, fh, errc)
		return errc
	}
	fs.log.Trace("Filescomfs: Getattr: getting attributes locally: path=%v, fh=%v", path, fh)
	errc = fs.local.Getattr(path, stat, fh)
	fs.log.Trace("Filescomfs: Getattr: got attributes locally: path=%v, fh=%v, errc=%v", path, fh, errc)
	return errc
}

func (fs *Filescomfs) Truncate(path string, size int64, fh uint64) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Truncate(path, size, fh)
		fs.log.Trace("Filescomfs: Truncate: truncating file remotely: path=%v, size=%v, fh=%v, errc=%v", path, size, fh, errc)
		return errc
	}
	errc = fs.local.Truncate(path, size, fh)
	fs.log.Trace("Filescomfs: Truncate: truncating file locally: path=%v, size=%v, fh=%v, errc=%v", path, size, fh, errc)
	return errc
}

func (fs *Filescomfs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		n = fs.remote.Read(path, buff, ofst, fh)
		fs.log.Trace("Filescomfs: Read: reading file remotely: path=%v, ofst=%v, fh=%v, len(buff)=%v, n=%d", path, ofst, fh, len(buff), n)
		return n
	}
	n = fs.local.Read(path, buff, ofst, fh)
	fs.log.Trace("Filescomfs: Read: reading file locally: path=%v, ofst=%v, fh=%v, len(buff)=%v, n=%v", path, ofst, fh, len(buff), n)
	return n
}

func (fs *Filescomfs) Write(path string, buff []byte, ofst int64, fh uint64) (n int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		n = fs.remote.Write(path, buff, ofst, fh)
		fs.log.Trace("Filescomfs: Write: writing file remotely: path=%v, ofst=%v, fh=%v, len(buff)=%v, n=%v", path, ofst, fh, len(buff), n)
		return n
	}
	n = fs.local.Write(path, buff, ofst, fh)
	fs.log.Trace("Filescomfs: Write: writing file locally: path=%v, ofst=%v, fh=%v, len(buff)=%v, n=%v", path, ofst, fh, len(buff), n)
	return n
}

func (fs *Filescomfs) Release(path string, fh uint64) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Release(path, fh)
		fs.log.Trace("Filescomfs: Release: releasing file remotely: path=%v, fh=%v, errc=%v", path, fh, errc)
		return errc
	}
	errc = fs.local.Release(path, fh)
	fs.log.Trace("Filescomfs: Release: releasing file locally: path=%v, fh=%v, errc=%v", path, fh, errc)
	return errc
}

func (fs *Filescomfs) Opendir(path string) (errc int, fh uint64) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc, fh = fs.remote.Opendir(path)
		fs.log.Trace("Filescomfs: Opendir: opening directory remotely: path=%v, errc=%v, fh=%v", path, errc, fh)
		return errc, fh
	}
	errc, fh = fs.local.Opendir(path)
	fs.log.Trace("Filescomfs: Opendir: opening directory locally: path=%v, errc=%v, fh=%v", path, errc, fh)
	return errc, fh
}

func (fs *Filescomfs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Readdir(path, fill, ofst, fh)
		fs.log.Trace("Filescomfs: Readdir: reading directory remotely: path=%v, ofst=%v, fh=%v, errc=%v", path, ofst, fh, errc)
		return errc
	}
	errc = fs.local.Readdir(path, fill, ofst, fh)
	fs.log.Trace("Filescomfs: Readdir: reading directory locally: path=%v, ofst=%v, fh=%v, errc=%v", path, ofst, fh, errc)
	return errc
}

func (fs *Filescomfs) Releasedir(path string, fh uint64) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Releasedir(path, fh)
		fs.log.Trace("Filescomfs: Releasedir: releasing directory remotely: path=%v, fh=%v, errc=%v", path, fh, errc)
		return errc
	}
	errc = fs.local.Releasedir(path, fh)
	fs.log.Trace("Filescomfs: Releasedir: releasing directory locally: path=%v, fh=%v, errc=%v", path, fh, errc)
	return errc
}

func (fs *Filescomfs) Chmod(path string, mode uint32) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Chmod(path, mode)
		fs.log.Trace("Filescomfs: Chmod: changing mode remotely: path=%v, mode=%v, errc=%v", path, mode, errc)
		return errc
	}
	errc = fs.local.Chmod(path, mode)
	fs.log.Trace("Filescomfs: Chmod: changing mode locally: path=%v, mode=%v, errc=%v", path, mode, errc)
	return errc
}

func (fs *Filescomfs) Fsync(path string, datasync bool, fh uint64) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Fsync(path, datasync, fh)
		fs.log.Trace("Filescomfs: Fsync: syncing file remotely: path=%v, datasync=%v, fh=%v, errc=%v", path, datasync, fh, errc)
		return errc
	}
	errc = fs.local.Fsync(path, datasync, fh)
	fs.log.Trace("Filescomfs: Fsync: syncing file locally: path=%v, datasync=%v, fh=%v, errc=%v", path, datasync, fh, errc)
	return errc
}

func (fs *Filescomfs) isStoredRemotely(path string) bool {
	switch {
	case path == "/":
		return true
	case fs.isIgnoreFile(path):
		return false
	case fs.isLockFile(path):
		return false
	default:
		return true
	}
}

func (fs *Filescomfs) isStoredLocally(path string) bool {
	return !fs.isStoredRemotely(path)
}

func (fs *Filescomfs) isIgnoreFile(path string) bool {
	if fs.ignore == nil {
		return false
	}
	return fs.ignore.MatchesPath(path)
}

func (fs *Filescomfs) isLockFile(path string) bool {
	return isMsOfficeOwnerFile(path) && !fs.disableLocking
}

func getStat(info fsNodeInfo, stat *fuse.Stat_t) *fuse.Stat_t {
	if stat == nil {
		stat = &fuse.Stat_t{}
	}

	if info.nodeType == nodeTypeDir {
		stat.Mode = fuse.S_IFDIR | 0777
	} else {
		stat.Mode = fuse.S_IFREG | 0777
	}

	stat.Size = info.size
	stat.Mtim = fuse.NewTimespec(info.modTime.UTC().Truncate(time.Second))
	if !info.creationTime.IsZero() {
		btime := fuse.NewTimespec(info.creationTime.UTC().Truncate(time.Second))
		stat.Birthtim = btime
	}

	return stat
}

func isFolderNotEmpty(err error) bool {
	var re files_sdk.ResponseError
	ok := errors.As(err, &re)
	return ok && re.Type == folderNotEmpty
}

// Mknod creates a file node.
func (fs *Filescomfs) Mknod(path string, mode uint32, dev uint64) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Mknod(path, mode, dev)
		fs.log.Trace("Filescomfs: Mknod: creating node remotely: path=%v, mode=%v, dev=%v, errc=%v", path, mode, dev, errc)
		return errc
	}
	errc = fs.local.Mknod(path, mode, dev)
	fs.log.Trace("Filescomfs: Mknod: creating node locally: path=%v, mode=%v, dev=%v, errc=%v", path, mode, dev, errc)
	return errc
}

// Link creates a hard link to a file.
func (fs *Filescomfs) Link(oldpath string, newpath string) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(newpath) {
		errc = fs.remote.Link(oldpath, newpath)
		fs.log.Trace("Filescomfs: Link: creating hard link remotely: oldpath=%v, newpath=%v, errc=%v", oldpath, newpath, errc)
		return errc
	}
	errc = fs.local.Link(oldpath, newpath)
	fs.log.Trace("Filescomfs: Link: creating hard link locally: oldpath=%v, newpath=%v, errc=%v", oldpath, newpath, errc)
	return errc
}

// Symlink creates a symbolic link.
func (fs *Filescomfs) Symlink(target string, newpath string) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(newpath) {
		errc = fs.remote.Symlink(target, newpath)
		fs.log.Trace("Filescomfs: Symlink: creating symlink remotely: target=%v, newpath=%v, errc=%v", target, newpath, errc)
		return errc
	}
	errc = fs.local.Symlink(target, newpath)
	fs.log.Trace("Filescomfs: Symlink: creating symlink locally: target=%v, newpath=%v, errc=%v", target, newpath, errc)
	return errc
}

// Readlink reads the target of a symbolic link.
func (fs *Filescomfs) Readlink(path string) (errc int, target string) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc, target = fs.remote.Readlink(path)
		fs.log.Trace("Filescomfs: Readlink: reading symlink remotely: path=%v, errc=%v, target=%v", path, errc, target)
		return errc, target
	}
	errc, target = fs.local.Readlink(path)
	fs.log.Trace("Filescomfs: Readlink: reading symlink locally: path=%v, errc=%v, target=%v", path, errc, target)
	return errc, target
}

// Chown changes the owner and group of a file.
func (fs *Filescomfs) Chown(path string, uid uint32, gid uint32) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Chown(path, uid, gid)
		fs.log.Trace("Filescomfs: Chown: changing owner/group remotely: path=%v, uid=%v, gid=%v, errc=%v", path, uid, gid, errc)
		return errc
	}
	errc = fs.local.Chown(path, uid, gid)
	fs.log.Trace("Filescomfs: Chown: changing owner/group locally: path=%v, uid=%v, gid=%v, errc=%v", path, uid, gid, errc)
	return errc
}

// Access checks file access permissions.
func (fs *Filescomfs) Access(path string, mask uint32) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Access(path, mask)
		fs.log.Trace("Filescomfs: Access: checking access remotely: path=%v, mask=%v, errc=%v", path, mask, errc)
		return errc
	}
	errc = fs.local.Access(path, mask)
	fs.log.Trace("Filescomfs: Access: checking access locally: path=%v, mask=%v, errc=%v", path, mask, errc)
	return errc
}

// Flush flushes cached file data.
func (fs *Filescomfs) Flush(path string, fh uint64) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Flush(path, fh)
		fs.log.Trace("Filescomfs: Flush: flushing file remotely: path=%v, fh=%v, errc=%v", path, fh, errc)
		return errc
	}
	errc = fs.local.Flush(path, fh)
	fs.log.Trace("Filescomfs: Flush: flushing file locally: path=%v, fh=%v, errc=%v", path, fh, errc)
	return errc
}

// Fsyncdir synchronizes directory contents.
func (fs *Filescomfs) Fsyncdir(path string, datasync bool, fh uint64) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Fsyncdir(path, datasync, fh)
		fs.log.Trace("Filescomfs: Fsyncdir: syncing directory remotely: path=%v, datasync=%v, fh=%v, errc=%v", path, datasync, fh, errc)
		return errc
	}
	errc = fs.local.Fsyncdir(path, datasync, fh)
	fs.log.Trace("Filescomfs: Fsyncdir: syncing directory locally: path=%v, datasync=%v, fh=%v, errc=%v", path, datasync, fh, errc)
	return errc
}

// Getxattr gets extended attributes.
func (fs *Filescomfs) Getxattr(path string, name string) (errc int, value []byte) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc, value = fs.remote.Getxattr(path, name)
		fs.log.Trace("Filescomfs: Getxattr: getting xattr remotely: path=%v, name=%v, errc=%v, valueLen=%v", path, name, errc, len(value))
		return errc, value
	}
	errc, value = fs.local.Getxattr(path, name)
	fs.log.Trace("Filescomfs: Getxattr: getting xattr locally: path=%v, name=%v, errc=%v, valueLen=%v", path, name, errc, len(value))
	return errc, value
}

// Setxattr sets extended attributes.
func (fs *Filescomfs) Setxattr(path string, name string, value []byte, flags int) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Setxattr(path, name, value, flags)
		fs.log.Trace("Filescomfs: Setxattr: setting xattr remotely: path=%v, name=%v, flags=%v, valueLen=%v, errc=%v", path, name, flags, len(value), errc)
		return errc
	}
	errc = fs.local.Setxattr(path, name, value, flags)
	fs.log.Trace("Filescomfs: Setxattr: setting xattr locally: path=%v, name=%v, flags=%v, valueLen=%v, errc=%v", path, name, flags, len(value), errc)
	return errc
}

// Removexattr removes extended attributes.
func (fs *Filescomfs) Removexattr(path string, name string) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Removexattr(path, name)
		fs.log.Trace("Filescomfs: Removexattr: removing xattr remotely: path=%v, name=%v, errc=%v", path, name, errc)
		return errc
	}
	errc = fs.local.Removexattr(path, name)
	fs.log.Trace("Filescomfs: Removexattr: removing xattr locally: path=%v, name=%v, errc=%v", path, name, errc)
	return errc
}

// Listxattr lists extended attributes.
func (fs *Filescomfs) Listxattr(path string, fill func(name string) bool) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Listxattr(path, fill)
		fs.log.Trace("Filescomfs: Listxattr: listing xattr remotely: path=%v, errc=%v", path, errc)
		return errc
	}
	errc = fs.local.Listxattr(path, fill)
	fs.log.Trace("Filescomfs: Listxattr: listing xattr locally: path=%v, errc=%v", path, errc)
	return errc
}

// CreateEx is similar to Create except that it allows direct manipulation of the FileInfo_t struct.
func (fs *Filescomfs) CreateEx(path string, mode uint32, fi *fuse.FileInfo_t) (errc int) {
	defer logPanics(fs.log)
	fs.log.Trace("Filescomfs: CreateEx: path=%v, mode=%v, fi=%v", path, mode, fi)
	errc, fh := fs.Create(path, fi.Flags, mode)
	fs.log.Trace("Filescomfs: CreateEx: created file: path=%v, mode=%v, flags=%v, errc=%v, fh=%v", path, mode, fi.Flags, errc, fh)
	fi.Fh = uint64(fh)
	return errc
}

// OpenEx is similar to Open except that it allows direct manipulation of the FileInfo_t struct.
func (fs *Filescomfs) OpenEx(path string, fi *fuse.FileInfo_t) (errc int) {
	defer logPanics(fs.log)
	fs.log.Trace("Filescomfs: OpenEx: path=%v, fi=%v", path, fi)
	errc, fh := fs.Open(path, fi.Flags)
	fs.log.Trace("Filescomfs: OpenEx: opened file: path=%v, flags=%v, errc=%v, fh=%v", path, fi.Flags, errc, fh)
	fi.Fh = uint64(fh)
	return errc
}

// Getpath allows a case-insensitive file system to report the correct case of a file path.
func (fs *Filescomfs) Getpath(path string, fh uint64) (errc int, result string) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc, result = fs.remote.Getpath(path, fh)
		fs.log.Trace("Filescomfs: Getpath: getting path remotely: path=%v, fh=%v, errc=%v, result=%v", path, fh, errc, result)
		return errc, result
	}
	errc, result = fs.local.Getpath(path, fh)
	fs.log.Trace("Filescomfs: Getpath: getting path locally: path=%v, fh=%v, errc=%v, result=%v", path, fh, errc, result)
	return errc, result
}

// Chflags changes the BSD file flags (Windows file attributes).
func (fs *Filescomfs) Chflags(path string, flags uint32) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Chflags(path, flags)
		fs.log.Trace("Filescomfs: Chflags: changing flags remotely: path=%v, flags=%v, errc=%v", path, flags, errc)
		return errc
	}
	errc = fs.local.Chflags(path, flags)
	fs.log.Trace("Filescomfs: Chflags: changing flags locally: path=%v, flags=%v, errc=%v", path, flags, errc)
	return errc
}

// Setcrtime changes the file creation (birth) time.
func (fs *Filescomfs) Setcrtime(path string, tmsp fuse.Timespec) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Setcrtime(path, tmsp)
		fs.log.Trace("Filescomfs: Setcrtime: setting creation time remotely: path=%v, tmsp=%v, errc=%v", path, tmsp, errc)
		return errc
	}
	errc = fs.local.Setcrtime(path, tmsp)
	fs.log.Trace("Filescomfs: Setcrtime: setting creation time locally: path=%v, tmsp=%v, errc=%v", path, tmsp, errc)
	return errc
}

// Setchgtime changes the file change (ctime) time.
func (fs *Filescomfs) Setchgtime(path string, tmsp fuse.Timespec) (errc int) {
	defer logPanics(fs.log)
	if fs.isStoredRemotely(path) {
		errc = fs.remote.Setchgtime(path, tmsp)
		fs.log.Trace("Filescomfs: Setchgtime: setting change time remotely: path=%v, tmsp=%v, errc=%v", path, tmsp, errc)
		return errc
	}
	errc = fs.local.Setchgtime(path, tmsp)
	fs.log.Trace("Filescomfs: Setchgtime: setting change time locally: path=%v, tmsp=%v, errc=%v", path, tmsp, errc)
	return errc
}
