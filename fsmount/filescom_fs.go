//go:build windows

package fsmount

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	path_lib "path"
	"path/filepath"
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

type Filescomfs struct {
	fuse.FileSystemBase
	*virtualfs
	mountPoint       string
	root             string
	writeConcurrency *int
	config           files_sdk.Config
	cacheTTL         *time.Duration
	disableLocking   bool
	ignore           *ignore.GitIgnore

	fileClient      *file.Client
	lockClient      *lock.Client
	migrationClient *file_migration.Client
	lockMap         map[string]*lockInfo
	lockMapMutex    sync.Mutex
}

type lockInfo struct {
	fh    uint64
	token string
}

func (fs *Filescomfs) Init() {
	defer fs.logPanics()
	if fs.fileClient == nil {
		fs.fileClient = &file.Client{Config: fs.config}
		fs.lockClient = &lock.Client{Config: fs.config}
		fs.migrationClient = &file_migration.Client{Config: fs.config}
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

func (fs *Filescomfs) Validate() (err error) {
	fs.Init()

	// Make sure we can list the root directory.
	it, err := fs.fileClient.ListFor(files_sdk.FolderListForParams{Path: fs.remotePath("/"), ListParams: files_sdk.ListParams{PerPage: 1}})
	if err == nil {
		it.Next() // Get 1 item. This is what actually triggers the API call.
		err = it.Err()
	}
	return
}

func (fs *Filescomfs) Statfs(path string, stat *fuse.Statfs_t) (errc int) {
	defer fs.logPanics()
	fs.Trace("Statfs: path=%v", path)

	totalBytes := uint64(1 << 50) // 1 PB?
	usedBytes := uint64(0)        // TODO: get used bytes
	freeBytes := totalBytes - usedBytes

	stat.Bsize = blockSize
	stat.Frsize = blockSize
	stat.Blocks = totalBytes / blockSize
	stat.Bfree = freeBytes / blockSize
	stat.Bavail = freeBytes / blockSize

	return 0
}

func (fs *Filescomfs) Mkdir(path string, mode uint32) (errc int) {
	defer fs.logPanics()
	localPath, remotePath := fs.paths(path)
	fs.Debug("Making dir: %v (%v) (mode=%v)", remotePath, localPath, mode)

	_, err := fs.fileClient.CreateFolder(files_sdk.FolderCreateParams{Path: remotePath})
	if files_sdk.IsExist(err) {
		return 0
	}

	// Windows File Explorer always tries to create the parent folder when writing a file, so don't
	// info-log until here in case the folder already exists.
	fs.Info("Creating folder: %v (%v)", remotePath, localPath)

	if errc = fs.handleError(path, err); errc != 0 {
		return
	}

	node := fs.getOrCreate(path, true)
	node.updateSize(0)

	return
}

func (fs *Filescomfs) Unlink(path string) (errc int) {
	defer fs.logPanics()
	localPath, remotePath := fs.paths(path)

	if node, ok := fs.fetch(path); ok {
		// Close the file and wait for any writes to complete before deleting the file.
		node.closeWriter(true)

		if node.isLocked() {
			fs.Info("Cannot delete locked file: %v (%v)", remotePath, localPath)
			return -fuse.ENOLCK
		}
	}

	// We may have been in the middle of writing the file, so don't log until here.
	fs.Info("Deleting file: %v (%v)", remotePath, localPath)

	return fs.delete(path)
}

func (fs *Filescomfs) Rmdir(path string) (errc int) {
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
		return
	}

	err = fs.waitForAction(action, "move")
	if errc = fs.handleError(oldpath, err); errc != 0 {
		return
	}

	fs.rename(oldpath, newpath)

	return
}

func (fs *Filescomfs) Utimens(path string, tmsp []fuse.Timespec) (errc int) {
	defer fs.logPanics()
	localPath, remotePath := fs.paths(path)
	fs.Debug("Updating provided mtime: %v (%v) (mtime=%v)", remotePath, localPath, tmsp[1])

	node, _ := fs.fetch(path)
	node.info.modTime = tmsp[1].Time()

	if node.isWriterOpen() {
		// If we're writing to the file, no need update the mtime. It will be updated when the write completes.
		return 0
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
		errc = -fuse.EEXIST
		return
	}

	fh = rand.Uint64()
	localPath, remotePath := fs.paths(path)
	fs.Debug("Creating file: %v (%v) (flags=%v, mode=%v, fh=%v)", remotePath, localPath, flags, mode, fh)

	if errc = fs.loadParent(path); errc != 0 {
		return
	}

	node, ok := fs.fetch(path)
	if ok && !node.infoExpired() {
		errc = -fuse.EEXIST
		return
	}

	if !ok {
		node = fs.getOrCreate(path, false)
	}

	node.updateSize(0)

	if errc = fs.lock(node, fh); errc != 0 {
		return
	}

	if !node.isWriterOpen() {
		fs.Info("Starting upload: %v (%v)", remotePath, localPath)
		node.openWriter(fs, fh)
	}

	return
}

func (fs *Filescomfs) Open(path string, flags int) (errc int, fh uint64) {
	defer fs.logPanics()
	isWrite := flags != fuse.O_RDONLY

	if isWrite && fs.ignoreWrite(path) {
		localPath, remotePath := fs.paths(path)
		fs.Debug("Ignoring file for upload: %v (%v)", remotePath, localPath)
		errc = -fuse.EACCES
		return
	}

	fh = rand.Uint64()
	fs.Trace("Open: path=%v, flags=%v, fh=%v", path, flags, fh)

	node := fs.getOrCreate(path, false)
	node.closeWriter(true)

	if isWrite {
		errc = fs.lock(node, fh)
	}

	return
}

func (fs *Filescomfs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	defer fs.logPanics()
	fs.Trace("Getattr: path=%v, fh=%v", path, fh)

	if node, ok := fs.fetch(path); ok && !node.infoExpired() {
		fs.Trace("Getattr: using cached stat, path=%v, size=%v, mtime=%v", path, node.info.size, node.info.modTime)
		getStat(node.info, stat)
		return
	}

	if errc = fs.loadParent(path); errc != 0 {
		return
	}

	node, ok := fs.fetch(path)
	if !ok || node.infoExpired() {
		node = nil

		if fs.isLockFile(path) {
			if lockedNode, ok := fs.fetchLockTarget(path); ok && lockedNode.isLocked() {
				node = fs.getOrCreate(path, false)
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

	return
}

func (fs *Filescomfs) Truncate(path string, size int64, fh uint64) (errc int) {
	defer fs.logPanics()
	localPath, remotePath := fs.paths(path)
	fs.Debug("Truncating file: %v (%v) (size=%v, fh=%v)", remotePath, localPath, size, fh)

	node, _ := fs.fetch(path)
	node.updateSize(size)

	if !node.isWriterOpen() {
		fs.Info("Starting upload: %v (%v)", remotePath, localPath)
		node.openWriter(fs, fh)
	}

	return
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
		// We can't read while writing to the file, so close the writer and wait for it to finish.
		// Do this _before_ we log that we're starting the download.
		fs.Debug("Finalizing upload before downloading: %v (%v) (fh=%v)", remotePath, localPath, fh)
		node.closeWriter(true)
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
	fs.Trace("Write: path=%v, len=%v, ofst=%v, fh=%v", path, len(buff), ofst, fh)

	node, _ := fs.fetch(path)

	if !node.isWriterOpen() {
		localPath, remotePath := fs.paths(path)
		fs.Info("Starting upload: %v (%v)", remotePath, localPath)
		node.openWriter(fs, fh)
	}

	n, err := node.writer.writeAt(buff, ofst)
	if errc := fs.handleError(path, err); errc != 0 {
		return errc
	}

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
		return
	}

	return fs.close(path, fh)
}

func (fs *Filescomfs) Opendir(path string) (errc int, fh uint64) {
	defer fs.logPanics()
	fh = rand.Uint64()
	fs.Trace("Opendir: path=%v, fh=%v", path, fh)

	fs.getOrCreate(path, true)
	return
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
		return
	}

	fill(".", nil, 0)
	fill("..", nil, 0)

	for childPath := range node.childPaths {
		if childNode, ok := fs.fetch(childPath); ok {
			fill(path_lib.Base(childPath), getStat(childNode.info, nil), 0)
		}
	}

	return
}

func (fs *Filescomfs) Releasedir(path string, fh uint64) (errc int) {
	defer fs.logPanics()
	fs.Trace("Releasedir: path=%v, fh=%v", path, fh)

	return fs.close(path, fh)
}

func (fs *Filescomfs) writeFile(path string, reader io.Reader, mtime *time.Time) {
	localPath, remotePath := fs.paths(path)
	uploadOpts := []file.UploadOption{
		file.UploadWithDestinationPath(remotePath),
		file.UploadWithReader(reader),
		file.UploadWithProvidedMtimePtr(mtime),
	}
	if fs.writeConcurrency != nil {
		uploadOpts = append(uploadOpts, file.UploadWithManager(manager.ConcurrencyManager{}.New(*fs.writeConcurrency)))
	}

	if err := fs.fileClient.Upload(uploadOpts...); err != nil {
		fs.Error("Upload failed: %v (%v): %v", remotePath, localPath, err)
		return
	}

	fs.Info("Upload completed: %v (%v)", remotePath, localPath)
}

func (fs *Filescomfs) lock(node *fsNode, fh uint64) (errc int) {
	if fs.disableLocking {
		return
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
		return
	}

	lock, err := fs.lockClient.Create(files_sdk.LockCreateParams{
		Path:                 remotePath,
		AllowAccessByAnyUser: lib.Ptr(true),
		Exclusive:            lib.Ptr(true),
		Recursive:            lib.Ptr(true),
		Timeout:              60 * 60, // 1 hour
	})
	if errc = fs.handleError(node.path, err); errc != 0 {
		return
	}

	// Update the local lock's path since it includes the full remote path. We build the full path ourselves.
	lock.Path = node.path

	fs.lockMap[node.path] = &lockInfo{fh: fh, token: lock.Token}
	return
}

func (fs *Filescomfs) unlock(path string, fh uint64) (errc int) {
	if fs.disableLocking {
		return
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
		return
	}

	localPath, remotePath := fs.paths(path)
	fs.Debug("Unlocking file: %v (%v)", remotePath, localPath)

	err := fs.lockClient.Delete(files_sdk.LockDeleteParams{
		Path:  remotePath,
		Token: lockInfo.token,
	})
	if errc = fs.handleError(path, err); errc != 0 {
		return
	}

	delete(fs.lockMap, path)
	return
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
		return
	}

	fs.remove(path)
	return
}

func (fs *Filescomfs) loadParent(path string) (errc int) {
	if path == "/" {
		// If we're at the root, we can't load the parent. Just make sure the root exists.
		_, errc = fs.findDir(path)
		return
	}

	parentPath := path_lib.Dir(path)
	parent, ok := fs.fetch(parentPath)

	// Make sure the parent is actually a directory that exists before attempting to load it.
	if !ok || parent.infoExpired() {
		parent, errc = fs.findDir(parentPath)
		if errc != 0 {
			return
		}
	}

	if !parent.info.dir {
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
		node = fs.getOrCreate(path, true)
		node.updateInfo(fsNodeInfo{dir: true})
		return
	}

	item, err := fs.fileClient.Find(files_sdk.FileFindParams{Path: remotePath})
	// Check for non-existence first so it doesn't get logged as an error, since this may be expected.
	if files_sdk.IsNotExist(err) {
		errc = -fuse.ENOENT
		return
	}
	if errc = fs.handleError(path, err); errc != 0 {
		return
	}
	if !item.IsDir() {
		errc = -fuse.ENOTDIR
		return
	}

	node = fs.createNode(path, item)

	return
}

func (fs *Filescomfs) loadDir(node *fsNode) (errc int) {
	err := node.updateChildPaths(fs.listDir)
	if errc = fs.handleError(node.path, err); errc != 0 {
		return
	}

	return
}

func (fs *Filescomfs) listDir(path string) (childPaths map[string]struct{}, err error) {
	it, err := fs.fileClient.ListFor(files_sdk.FolderListForParams{Path: fs.remotePath(path)})
	if err != nil {
		return
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
		return
	}

	if fs.disableLocking {
		return
	}

	locks, err := fs.lockClient.ListFor(files_sdk.LockListForParams{
		Path:            fs.remotePath(path),
		IncludeChildren: lib.Ptr(true),
	})
	if err != nil {
		return
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

	return
}

func (fs *Filescomfs) createNode(path string, item files_sdk.File) (node *fsNode) {
	node = fs.getOrCreate(path, item.IsDir())
	node.updateInfo(fsNodeInfo{
		dir:          item.IsDir(),
		size:         item.Size,
		modTime:      item.ModTime(),
		creationTime: item.CreatedAt,
	})

	return
}

func (fs *Filescomfs) waitForAction(action files_sdk.FileAction, operation string) error {
	migration, err := fs.migrationClient.Wait(action, func(migration files_sdk.FileMigration) {
		// noop
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

	if info.dir {
		stat.Mode = fuse.S_IFDIR | 0777
	} else {
		stat.Mode = fuse.S_IFREG | 0777
	}

	stat.Size = info.size
	stat.Mtim = fuse.NewTimespec(info.modTime.UTC().Truncate(time.Second))
	if info.creationTime != nil {
		stat.Birthtim = fuse.NewTimespec(*info.creationTime)
	}

	return stat
}

func isFolderNotEmpty(err error) bool {
	var re files_sdk.ResponseError
	ok := errors.As(err, &re)
	return ok && re.Type == folderNotEmpty
}
