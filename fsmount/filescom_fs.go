package fsmount

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	path_lib "path"
	"path/filepath"
	"runtime"
	"slices"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	file_migration "github.com/Files-com/files-sdk-go/v3/filemigration"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lock"
	ignore "github.com/sabhiram/go-gitignore"
	"github.com/winfsp/cgofuse/fuse"
)

const (
	folderNotEmpty = "processing-failure/folder-not-empty"
	blockSize      = 4096
)

// Filescomfs is a filesystem that implements the fuse.FileSystem interface,
// allowing it to be mounted using FUSE. It provides a virtual filesystem
// interface to Files.com, allowing users to interact with their Files.com
// account as if it were a local filesystem.
type Filescomfs struct {
	fuse.FileSystemBase // implements fuse.FileSystem with no-op methods
	*virtualfs
	config           *files_sdk.Config
	mountPoint       string
	root             string
	writeConcurrency int
	cacheTTL         time.Duration
	disableLocking   bool
	ignore           *ignore.GitIgnore

	fileClient      *file.Client
	lockClient      *lock.Client
	migrationClient *file_migration.Client
	lockMap         map[string]*lockInfo
	lockMapMutex    sync.Mutex
	debugFuse       bool
}

type lockInfo struct {
	fh    uint64
	token string
}

// Init initializes the Filescomfs filesystem.
func (fs *Filescomfs) Init() {
	defer fs.logPanics()
	if fs.fileClient == nil {
		fs.fileClient = &file.Client{Config: *fs.config}
		fs.lockClient = &lock.Client{Config: *fs.config}
		fs.migrationClient = &file_migration.Client{Config: *fs.config}
		fs.lockMap = make(map[string]*lockInfo)
		fs.virtualfs = newVirtualfs(fs.config.Logger, fs.cacheTTL)
	}
}

func (fs *Filescomfs) Destroy() {
	fs.Debug("Destroying filesystem")

	for path, lockInfo := range fs.lockMap {
		fs.unlock(path, lockInfo.fh)
	}
}

// Validate checks if the Filescomfs filesystem is valid by attempting to list the root directory.
func (fs *Filescomfs) Validate() error {
	fs.Init()

	// Make sure we can list the root directory.
	it, err := fs.fileClient.ListFor(files_sdk.FolderListForParams{Path: fs.remotePath("/"), ListParams: files_sdk.ListParams{PerPage: 1}})
	if err == nil {
		it.Next() // Get 1 item. This is what actually triggers the API call.
		err = it.Err()
	}
	return err
}

func (fs *Filescomfs) Statfs(path string, stat *fuse.Statfs_t) (errc int) {
	defer fs.logPanics()
	fs.Trace("Statfs: path=%v", path)

	totalBytes := remoteCapacityBytes() // 1 PB?
	usedBytes := uint64(0)              // TODO: get used bytes from the remote
	freeBytes := totalBytes - usedBytes

	stat.Bsize = blockSize
	stat.Frsize = blockSize
	stat.Blocks = totalBytes / blockSize
	stat.Bfree = freeBytes / blockSize
	stat.Bavail = freeBytes / blockSize

	return errc
}

func remoteCapacityBytes() uint64 {
	// the remote capacity is functionally unlimited, but we need to return the largest
	// value that the OS will accept, otherwise it will show up as zero capacity.
	switch runtime.GOOS {
	case "darwin":
		return uint64(1 << 43) // ~8TB - any larger and the nfs drive shows up as zero capacity on macOS
	default:
		return uint64(1 << 50) // ~1PB
	}
}

func (fs *Filescomfs) Mkdir(path string, mode uint32) (errc int) {
	defer fs.logPanics()
	localPath, remotePath := fs.paths(path)
	fs.Debug("Making dir: %v (%v) (mode=%v)", remotePath, localPath, mode)

	_, err := fs.fileClient.CreateFolder(files_sdk.FolderCreateParams{Path: remotePath})
	if files_sdk.IsExist(err) {
		return errc
	}

	// Windows File Explorer always tries to create the parent folder when writing a file, so don't
	// info-log until here in case the folder already exists.
	fs.Info("Creating folder: %v (%v)", remotePath, localPath)

	if errc = fs.handleError(path, err); errc != 0 {
		return errc
	}

	node := fs.getOrCreate(path, nodeTypeDir)
	node.updateSize(0)

	return errc
}

func (fs *Filescomfs) Unlink(path string) int {
	defer fs.logPanics()
	localPath, remotePath := fs.paths(path)

	if node, ok := fs.fetch(path); ok {
		// If the node is locked, we cannot delete it.
		if node.isLocked() || node.isWriterOpen() {
			fs.Info("Cannot delete locked file: %v (%v)", remotePath, localPath)
			return -fuse.ENOLCK
		}
	}

	// We may have been in the middle of writing the file, so don't log until here.
	fs.Info("Deleting file: %v (%v)", remotePath, localPath)

	return fs.delete(path)
}

func (fs *Filescomfs) Rmdir(path string) int {
	defer fs.logPanics()
	localPath, remotePath := fs.paths(path)
	fs.Info("Deleting folder: %v (%v)", remotePath, localPath)

	return fs.delete(path)
}

func (fs *Filescomfs) Rename(oldpath string, newpath string) (errc int) {
	defer fs.logPanics()
	oldLocalPath, oldRemotePath := fs.paths(oldpath)
	newLocalPath, newRemotePath := fs.paths(newpath)
	fs.Info("Renaming: %v to %v (%v to %v)", oldRemotePath, newRemotePath, oldLocalPath, newLocalPath)

	if node, ok := fs.fetch(oldpath); ok && node.isLocked() {
		fs.Info("Cannot rename locked file: %v (%v)", oldRemotePath, oldLocalPath)
		return -fuse.ENOLCK
	}

	params := files_sdk.FileMoveParams{
		Path:        oldRemotePath,
		Destination: newRemotePath,
		Overwrite:   lib.Ptr(true),
	}

	action, err := fs.fileClient.Move(params)
	if errc = fs.handleError(oldpath, err); errc != 0 {
		return errc
	}

	err = fs.waitForAction(action, "move")
	if errc = fs.handleError(oldpath, err); errc != 0 {
		return errc
	}

	fs.rename(oldpath, newpath)

	return errc
}

func (fs *Filescomfs) Utimens(path string, tmsp []fuse.Timespec) (errc int) {
	defer fs.logPanics()
	localPath, remotePath := fs.paths(path)
	modT := tmsp[1].Time()
	fs.Debug("Updating mtime for: %v (%v) (mtime=%v)", remotePath, localPath, modT)

	node, _ := fs.fetch(path)
	node.info.modTime = modT

	if node.isWriterOpen() {
		// If we're writing to the file, no need update the mtime. It will be updated when the write completes.
		return errc
	}

	params := files_sdk.FileUpdateParams{
		Path:          remotePath,
		ProvidedMtime: &node.info.modTime,
	}

	_, err := fs.fileClient.Update(params)
	return fs.handleError(path, err)
}

func (fs *Filescomfs) Create(path string, flags int, mode uint32) (errc int, fh uint64) {
	defer fs.logPanics()

	if fs.ignoreWrite(path) {
		return -fuse.ENOENT, fh
	}

	fh = rand.Uint64()
	localPath, remotePath := fs.paths(path)
	fs.Debug("Creating file: %v (%v) (flags=%v, mode=%v, fh=%v)", remotePath, localPath, flags, mode, fh)

	if errc = fs.loadParent(path); errc != 0 {
		return errc, fh
	}

	node, ok := fs.fetch(path)
	if ok && !node.infoExpired() {
		errc = -fuse.EEXIST
		return errc, fh
	}

	if !ok {
		node = fs.getOrCreate(path, nodeTypeFile)
	}

	node.updateSize(0)

	if errc = fs.lock(node, fh); errc != 0 {
		return errc, fh
	}

	if !node.isWriterOpen() {
		fs.Debug("Opening writer from Create: %v (%v)", remotePath, localPath)
		node.openWriter(fs, fh)
	}

	return errc, fh
}

func (fs *Filescomfs) Open(path string, flags int) (errc int, fh uint64) {
	defer fs.logPanics()
	fh = rand.Uint64()
	fs.Debug("Open: path=%v, flags=%v, fh=%v", path, flags, fh)

	readOnly := flags == fuse.O_RDONLY
	node := fs.getOrCreate(path, nodeTypeFile)

	// If the requested op is read only, and the writer is not open,
	// return 0 and a file handle.  TODO: should the fh be ^uint64(0)
	if readOnly && !node.isWriterOpen() {
		return errc, fh
	}

	// If the requested op is read only, and the writer is already open,
	// return a busy status and a file handle. TODO: should the fh be ^uint64(0)
	if readOnly && node.isWriterOpen() {
		return -fuse.EBUSY, fh
	}
	// after this point, the requested op must be a write operation

	// return noent if the file is ignored
	if fs.ignoreWrite(path) {
		localPath, remotePath := fs.paths(path)
		fs.Debug("Ignoring file for upload: %v (%v)", remotePath, localPath)
		return -fuse.ENOENT, fh
	}

	errc = fs.lock(node, fh)
	if errc != 0 {
		return errc, fh
	}

	// open the writer and associate it with a file handle
	if !node.isWriterOpen() {
		node.openWriter(fs, fh)
	}

	return errc, fh
}

func (fs *Filescomfs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	defer fs.logPanics()
	fs.Trace("Getattr: path=%v, fh=%v", path, fh)

	if node, ok := fs.fetch(path); ok && !node.infoExpired() {
		fs.Trace("Getattr: using cached stat, path=%v, size=%v, mtime=%v", path, node.info.size, node.info.modTime)
		getStat(node.info, stat)
		return errc
	}

	if errc = fs.loadParent(path); errc != 0 {
		return errc
	}

	node, ok := fs.fetch(path)
	if !ok || node.infoExpired() {
		node = nil

		if fs.isLockFile(path) {
			if lockedNode, ok := fs.fetchLockTarget(path); ok && lockedNode.isLocked() {
				node = fs.getOrCreate(path, nodeTypeFile)
				node.updateInfo(fsNodeInfo{
					size:    int64(len(buildOwnerFile(lockedNode))),
					modTime: time.Now(),
				})
			}
		}

		if node == nil {
			localPath, remotePath := fs.paths(path)
			fs.Debug("File not found: %v (%v)", remotePath, localPath)
			return -fuse.ENOENT
		}
	}

	fs.Trace("Getattr: path=%v, size=%v, mtime=%v", path, node.info.size, node.info.modTime)
	getStat(node.info, stat)

	return errc
}

func (fs *Filescomfs) Truncate(path string, size int64, fh uint64) (errc int) {
	// The word truncate is overloaded here. The intention is to set the size of the
	// file to the size getting passed in, NOT to truncate the file to zero bytes.
	defer fs.logPanics()
	localPath, remotePath := fs.paths(path)
	fs.Debug("Truncating file: %v (%v) (size=%v, fh=%v)", remotePath, localPath, size, fh)

	node, _ := fs.fetch(path)
	node.updateSize(size)

	if !node.isWriterOpen() {
		fs.Debug("Opening writer from Truncate: %v (%v)", remotePath, localPath)
		node.openWriter(fs, fh)
	}

	return errc
}

func (fs *Filescomfs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	defer fs.logPanics()
	buffLen := int64(len(buff))
	fs.Trace("Read: path=%v, len=%v, ofst=%v, fh=%v", path, buffLen, ofst, fh)

	localPath, remotePath := fs.paths(path)
	node, _ := fs.fetch(path)

	if node.info.size == 0 {
		fs.Trace("Read: file is empty, returning EOF")
		return 0
	}

	if ofst > 0 && ofst >= node.info.size {
		fs.Trace("Read: offset %d is greater than file size %d, returning EOF", ofst, node.info.size)
		return 0
	}

	if node.isWriterOpen() {
		fs.Info("Cannot read from file while writing: %v (%v) (fh=%v)", remotePath, localPath, fh)
		return -fuse.EBUSY
	}

	if fs.isLockFile(path) {
		if lockedNode, ok := fs.fetchLockTarget(path); ok && lockedNode.isLocked() {
			ownerBuffer := buildOwnerFile(lockedNode)
			return copy(buff, ownerBuffer[ofst:])
		}
	}

	// Read up to the end of the file.
	buffLen = min(buffLen, node.info.size-ofst)

	if ofst == 0 && buffLen >= min(blockSize, node.info.size) {
		node.readerHandle = fh
		fs.Info("Starting download: %v (%v)", remotePath, localPath)
	}

	headers := &http.Header{}
	headers.Set("Range", fmt.Sprintf("bytes=%v-%v", ofst, ofst+buffLen-1))
	file, err := fs.fileClient.Download(
		files_sdk.FileDownloadParams{File: files_sdk.File{
			Path:        remotePath,
			DownloadUri: node.downloadUri,
		}},
		files_sdk.RequestHeadersOption(headers),
		files_sdk.ResponseOption(func(response *http.Response) error {
			var err error
			if err = lib.ResponseErrors(response, lib.IsStatus(http.StatusForbidden), lib.NotStatus(http.StatusPartialContent), files_sdk.APIError()); err != nil {
				return err
			}
			n, err = io.ReadAtLeast(response.Body, buff, int(buffLen))
			return err
		}),
	)
	if errc := fs.handleError(path, err); errc != 0 {
		return errc
	}

	node.downloadUri = file.DownloadUri

	fs.Trace("Read: path=%v, ofst=%d, read %d bytes", path, ofst, n)

	return n
}

func (fs *Filescomfs) Write(path string, buff []byte, ofst int64, fh uint64) (n int) {
	defer fs.logPanics()
	fs.Debug("Write: path=%v, len=%v, ofst=%v, fh=%v", path, len(buff), ofst, fh)

	node, _ := fs.fetch(path)

	if !node.isWriterOpen() {
		localPath, remotePath := fs.paths(path)
		fs.Debug("Opening writer from Write: %v (%v)", remotePath, localPath)
		node.openWriter(fs, fh)
	}

	n, err := node.writer.writeAt(buff, ofst)
	if errc := fs.handleError(path, err); errc != 0 {
		return errc
	}

	// If the write was successful, update the size of the node. This ensures
	// that Getattr calls return an up to date size that changes as the upload
	// progresses and allows the OS to more accurately calculate upload time
	// remaining. It has the additional benefit of keeping the infoExpires
	// field current, and keeps the node in the cache.
	node.updateSize(node.writer.offset)

	return n
}

func (fs *Filescomfs) Release(path string, fh uint64) (errc int) {
	defer fs.logPanics()
	fs.Trace("Release: path=%v, fh=%v", path, fh)
	defer fs.Trace("Released: path=%v, fh=%v", path, fh)

	if node, ok := fs.fetch(path); ok && node.readerHandle == fh {
		localPath, remotePath := fs.paths(path)
		fs.Info("Download completed: %v (%v)", remotePath, localPath)
		node.readerHandle = 0
	}

	if errc = fs.unlock(path, fh); errc != 0 {
		return errc
	}

	return fs.release(path, fh)
}

func (fs *Filescomfs) Opendir(path string) (errc int, fh uint64) {
	defer fs.logPanics()
	fh = rand.Uint64()
	fs.Trace("Opendir: path=%v, fh=%v", path, fh)

	fs.getOrCreate(path, nodeTypeDir)
	return errc, fh
}

func (fs *Filescomfs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	defer fs.logPanics()
	localPath, remotePath := fs.paths(path)
	fs.Info("Listing folder: %v (%v)", remotePath, localPath)

	node, _ := fs.fetch(path)
	if errc = fs.loadDir(node); errc != 0 {
		return errc
	}

	fill(".", nil, 0)
	fill("..", nil, 0)

	// provide a consistent sort order when calling fill
	keys := make([]string, len(node.childPaths))
	pos := 0
	for childPath := range node.childPaths {
		keys[pos] = childPath
		pos++
	}
	slices.Sort(keys)
	for _, childPath := range keys {
		if childNode, ok := fs.fetch(childPath); ok {
			fill(path_lib.Base(childPath), getStat(childNode.info, nil), 0)
		}
	}

	return errc
}

func (fs *Filescomfs) Releasedir(path string, fh uint64) int {
	defer fs.logPanics()
	fs.Trace("Releasedir: path=%v, fh=%v", path, fh)

	return fs.release(path, fh)
}

func (fs *Filescomfs) writeFile(path string, reader io.Reader, mtime time.Time) {
	localPath, remotePath := fs.paths(path)
	fs.Info("Starting upload: %v (%v)", remotePath, localPath)
	uploadOpts := []file.UploadOption{
		file.UploadWithDestinationPath(remotePath),
		file.UploadWithReader(reader),
		file.UploadWithProvidedMtime(mtime),
	}
	if fs.writeConcurrency != 0 {
		uploadOpts = append(uploadOpts, file.UploadWithManager(manager.ConcurrencyManager{}.New(fs.writeConcurrency)))
	}

	start := time.Now()
	u, err := fs.fileClient.UploadWithResume(uploadOpts...)
	if err != nil {
		fs.Error("Upload failed: %v (%v): %v", remotePath, localPath, err)
		return
	}

	fs.Info("Upload completed: %v (%v).", remotePath, localPath)
	fs.Debug("Bytes: %v, Duration: %v", u.Size, time.Since(start))

}

func (fs *Filescomfs) lock(node *fsNode, fh uint64) (errc int) {
	if fs.disableLocking {
		return errc
	}

	node.lockMutex.Lock()
	defer node.lockMutex.Unlock()

	fs.lockMapMutex.Lock()
	defer fs.lockMapMutex.Unlock()

	localPath, remotePath := fs.paths(node.path)
	fs.Debug("Locking file: %v (%v)", remotePath, localPath)

	if node.isLocked() {
		fs.Debug("File is already locked by %v: %v (%v)", node.info.lockOwner, remotePath, localPath)
		errc = -fuse.ENOLCK
		return errc
	}

	lock, err := fs.lockClient.Create(files_sdk.LockCreateParams{
		Path:                 remotePath,
		AllowAccessByAnyUser: lib.Ptr(true),
		Exclusive:            lib.Ptr(true),
		Recursive:            lib.Ptr(true),
		Timeout:              60 * 60, // 1 hour
	})
	if errc = fs.handleError(node.path, err); errc != 0 {
		return errc
	}

	// Update the local lock's path since it includes the full remote path. We build the full path ourselves.
	lock.Path = node.path

	fs.lockMap[node.path] = &lockInfo{fh: fh, token: lock.Token}
	return errc
}

func (fs *Filescomfs) unlock(path string, fh uint64) (errc int) {
	if fs.disableLocking {
		return errc
	}

	// If we have a node, prevent locking while we're unlocking.
	// If the node was renamed/moved, we won't have a node, but we may still need to unlock it.
	if node, ok := fs.fetch(path); ok {
		node.lockMutex.Lock()
		defer node.lockMutex.Unlock()
	}

	fs.lockMapMutex.Lock()
	defer fs.lockMapMutex.Unlock()

	lockInfo, ok := fs.lockMap[path]
	if !ok || lockInfo.fh != fh {
		// This is fine. It just means the file either wasn't locked or it was locked by a different file handle.
		return errc
	}

	localPath, remotePath := fs.paths(path)
	fs.Debug("Unlocking file: %v (%v)", remotePath, localPath)

	err := fs.lockClient.Delete(files_sdk.LockDeleteParams{
		Path:  remotePath,
		Token: lockInfo.token,
	})
	if errc = fs.handleError(path, err); errc != 0 {
		return errc
	}

	delete(fs.lockMap, path)
	return errc
}

func (fs *Filescomfs) paths(path string) (string, string) {
	return fs.localPath(path), fs.remotePath(path)
}

func (fs *Filescomfs) localPath(path string) string {
	return filepath.Join(fs.mountPoint, path)
}

func (fs *Filescomfs) remotePath(path string) string {
	return path_lib.Join(fs.root, path)
}

func (fs *Filescomfs) handleError(path string, err error) int {
	if err != nil {
		localPath, remotePath := fs.paths(path)
		fs.Error("%v (%v): %v", remotePath, localPath, err)

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

func (fs *Filescomfs) delete(path string) (errc int) {
	err := fs.fileClient.Delete(files_sdk.FileDeleteParams{Path: fs.remotePath(path)})
	if errc = fs.handleError(path, err); errc != 0 {
		return errc
	}

	fs.remove(path)
	return errc
}

func (fs *Filescomfs) loadParent(path string) (errc int) {
	if path == "/" {
		// If we're at the root, we can't load the parent. Just make sure the root exists.
		_, errc = fs.findDir(path)
		return errc
	}

	parentPath := path_lib.Dir(path)
	parent, ok := fs.fetch(parentPath)

	// Make sure the parent is actually a directory that exists before attempting to load it.
	if !ok || parent.infoExpired() {
		parent, errc = fs.findDir(parentPath)
		if errc != 0 {
			return errc
		}
	}

	if parent.info.nodeType != nodeTypeDir {
		// Don't log an error. Windows File Explorer sometimes treats shortcuts as parent directories.
		fs.Trace("Parent of %s is not a directory %s", path, parentPath)
		return -fuse.ENOTDIR
	}

	return fs.loadDir(parent)
}

func (fs *Filescomfs) findDir(path string) (node *fsNode, errc int) {
	remotePath := fs.remotePath(path)

	if remotePath == "/" {
		// Special case that we can't stat the root directory of a Files.com site.
		node = fs.getOrCreate(path, nodeTypeDir)
		node.updateInfo(fsNodeInfo{nodeType: nodeTypeDir})
		return node, errc
	}

	item, err := fs.fileClient.Find(files_sdk.FileFindParams{Path: remotePath})
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

func (fs *Filescomfs) loadDir(node *fsNode) (errc int) {
	err := node.updateChildPaths(fs.listDir)
	if errc = fs.handleError(node.path, err); errc != 0 {
		return errc
	}

	return errc
}

func (fs *Filescomfs) listDir(path string) (childPaths map[string]struct{}, err error) {
	it, err := fs.fileClient.ListFor(files_sdk.FolderListForParams{Path: fs.remotePath(path)})
	if err != nil {
		return nil, err
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
		return childPaths, err
	}

	if fs.disableLocking {
		return childPaths, err
	}

	locks, err := fs.lockClient.ListFor(files_sdk.LockListForParams{
		Path:            fs.remotePath(path),
		IncludeChildren: lib.Ptr(true),
	})
	if err != nil {
		return childPaths, err
	}

	for locks.Next() {
		lock := locks.Lock()
		childPath := path_lib.Join(path, path_lib.Base(lock.Path))

		// Ignore paths where *we* hold the lock.
		if _, ok := fs.lockMap[childPath]; ok {
			continue
		}

		if child, ok := fs.fetch(childPath); ok {
			child.info.lockOwner = lock.Username
		}
	}
	err = locks.Err()

	return childPaths, err
}

func (fs *Filescomfs) createNode(path string, item files_sdk.File) *fsNode {
	var nt nodeType
	if item.IsDir() {
		nt = nodeTypeDir
	} else {
		nt = nodeTypeFile
	}
	node := fs.getOrCreate(path, nt)
	node.updateInfo(fsNodeInfo{
		nodeType:     nt,
		size:         item.Size,
		modTime:      item.ModTime(),
		creationTime: item.CreationTime(),
	})

	return node
}

func (fs *Filescomfs) waitForAction(action files_sdk.FileAction, operation string) error {
	migration, err := fs.migrationClient.Wait(action, func(migration files_sdk.FileMigration) {
		fs.Trace("waiting for migration")
	})
	if err == nil && migration.Status != "completed" {
		return fmt.Errorf("%v did not complete successfully: %v", operation, migration.Status)
	}
	return err
}

func (fs *Filescomfs) ignoreWrite(path string) bool {
	return fs.isIgnoreFile(path) || fs.isLockFile(path)
}

func (fs *Filescomfs) isIgnoreFile(path string) bool {
	return fs.ignore != nil && fs.ignore.MatchesPath(path)
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
		stat.Birthtim = fuse.NewTimespec(info.creationTime)
	}

	return stat
}

func isFolderNotEmpty(err error) bool {
	var re files_sdk.ResponseError
	ok := errors.As(err, &re)
	return ok && re.Type == folderNotEmpty
}

// Methods below are part of the fuse.FileSystemInterface, but not supported by
// this implementation. They exist here to support logging for visibility of how
// the underlying fuse layer calls into this implementation.

// Mknod creates a file node.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *Filescomfs) Mknod(path string, mode uint32, dev uint64) int {
	fs.Trace("Mknod: path=%v, mode=%v, dev=%v", path, mode, dev)
	return -fuse.ENOSYS
}

// Link creates a hard link to a file.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *Filescomfs) Link(oldpath string, newpath string) int {
	fs.Trace("Link: old=%v, new=%v", oldpath, newpath)
	return -fuse.ENOSYS
}

// Symlink creates a symbolic link.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *Filescomfs) Symlink(target string, newpath string) int {
	fs.Trace("Symlink: target=%v, newpath=%v", target, newpath)
	return -fuse.ENOSYS
}

// Readlink reads the target of a symbolic link.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *Filescomfs) Readlink(path string) (int, string) {
	fs.Trace("Readlink: path=%v", path)
	return -fuse.ENOSYS, ""
}

// Chmod changes the permission bits of a file.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *Filescomfs) Chmod(path string, mode uint32) int {
	fs.Trace("Chmod: path=%v, mode=%v", path, mode)
	return -fuse.ENOSYS
}

// Chown changes the owner and group of a file.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *Filescomfs) Chown(path string, uid uint32, gid uint32) int {
	fs.Trace("Chown: path=%v, uid=%v, gid=%v", path, uid, gid)
	return -fuse.ENOSYS
}

// Access checks file access permissions.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *Filescomfs) Access(path string, mask uint32) int {
	fs.Trace("Access: path=%v, mask=%v", path, mask)
	return -fuse.ENOSYS
}

// Flush flushes cached file data.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *Filescomfs) Flush(path string, fh uint64) int {
	fs.Trace("Flush: path=%v, fh=%v", path, fh)
	return -fuse.ENOSYS
}

// Fsync synchronizes file contents.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *Filescomfs) Fsync(path string, datasync bool, fh uint64) int {
	fs.Trace("Fsync: path=%v, datasync=%v, fh=%v", path, datasync, fh)
	return -fuse.ENOSYS
}

// Fsyncdir synchronizes directory contents.
// The return value of -fuse.ENOSYS indicates the method is not supported.
func (fs *Filescomfs) Fsyncdir(path string, datasync bool, fh uint64) int {
	fs.Trace("Fsyncdir: path=%v, datasync=%v, fh=%v", path, datasync, fh)
	return -fuse.ENOSYS
}

// The [Foo]xattr implementations below explicitly return 0 to indicate that
// extended attributes are "supported" in order to ensure that the other xattr
// methods are called for debugging visibility, but are all no-op implementations.

// Getxattr gets extended attributes.
// Any return value other than -fuse.ENOSYS indicates support for extended
// attributes, but also expects Setxattr, Listxattr, and Removexattr to exist
// for extended attribute support.
func (fs *Filescomfs) Getxattr(path string, name string) (int, []byte) {
	fs.Debug("Getxattr: path=%v, name=%v", path, name)
	return 0, []byte{}
}

// Setxattr sets extended attributes.
func (fs *Filescomfs) Setxattr(path string, name string, value []byte, flags int) int {
	fs.Debug("Setxattr: path=%v, name=%v, value=%v flags=%v", path, name, value, flags)
	return 0
}

// Removexattr removes extended attributes.
func (fs *Filescomfs) Removexattr(path string, name string) int {
	fs.Debug("Removexattr: path=%v, name=%v", path, name)
	return 0
}

// Listxattr lists extended attributes.
func (fs *Filescomfs) Listxattr(path string, fill func(name string) bool) int {
	fs.Debug("Listxattr: path=%v", path)
	return 0
}

// FileSystemOpenEx is the interface that wraps the OpenEx and CreateEx methods.

// OpenEx and CreateEx are similar to Open and Create except that they allow
// direct manipulation of the FileInfo_t struct (which is analogous to the
// FUSE struct fuse_file_info). If implemented, they are preferred over
// Open and Create.
func (fs *Filescomfs) CreateEx(path string, mode uint32, fi *fuse.FileInfo_t) int {
	fs.Debug("CreateEx: path=%v, mode=%v, fi=%v", path, mode, fi)
	errc, fh := fs.Create(path, fi.Flags, mode)
	fi.Fh = fh
	return errc
}

func (fs *Filescomfs) OpenEx(path string, fi *fuse.FileInfo_t) int {
	fs.Debug("OpenEx: path=%v, fi=%v", path, fi)
	errc, fh := fs.Open(path, fi.Flags)
	fi.Fh = fh
	return errc
}

// Getpath is part of the FileSystemGetpath interface and
// allows a case-insensitive file system to report the correct case of a file path.
func (fs *Filescomfs) Getpath(path string, fh uint64) (int, string) {
	fs.Trace("Getpath: path=%v, fh=%v", path, fh)
	return -fuse.ENOSYS, path
}

// Chflags is part of the FileSystemChflags interface and
// changes the BSD file flags (Windows file attributes).
func (fs *Filescomfs) Chflags(path string, flags uint32) int {
	fs.Trace("Chflags: path=%v, flags=%v", path, flags)
	return -fuse.ENOSYS
}

// Setcrtime is part of the FileSystemSetcrtime interface and
// changes the file creation (birth) time.
func (fs *Filescomfs) Setcrtime(path string, tmsp fuse.Timespec) int {
	fs.Trace("Setcrtime: path=%v, tmsp=%v", path, tmsp)
	return -fuse.ENOSYS
}

// Setchgtime is part of the FileSystemSetchgtime interface and
// changes the file change (ctime) time.
func (fs *Filescomfs) Setchgtime(path string, tmsp fuse.Timespec) int {
	fs.Trace("Setchgtime: path=%v, tmsp=%v", path, tmsp)
	return -fuse.ENOSYS
}
