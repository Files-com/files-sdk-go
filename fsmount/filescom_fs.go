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
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	file_migration "github.com/Files-com/files-sdk-go/v3/filemigration"
	"github.com/Files-com/files-sdk-go/v3/lib"
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
	fileClient       *file.Client
	migrationClient  *file_migration.Client
}

func (self *Filescomfs) Init() {
	defer self.logPanics()
	self.fileClient = &file.Client{Config: self.config}
	self.migrationClient = &file_migration.Client{Config: self.config}
	self.virtualfs = newVirtualfs(self.config.Logger, self.cacheTTL)
}

func (self *Filescomfs) Statfs(path string, stat *fuse.Statfs_t) (errc int) {
	defer self.logPanics()
	self.Trace("Statfs: path=%v", path)

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

func (self *Filescomfs) Mkdir(path string, mode uint32) (errc int) {
	defer self.logPanics()
	localPath, remotePath := self.paths(path)
	self.Debug("Making dir: %v (%v) (mode=%v)", remotePath, localPath, mode)

	_, err := self.fileClient.CreateFolder(files_sdk.FolderCreateParams{Path: remotePath})
	if files_sdk.IsExist(err) {
		return 0
	}

	// Windows File Explorer always tries to create the parent folder when writing a file, so don't
	// info-log until here in case the folder already exists.
	self.Info("Creating folder: %v (%v)", remotePath, localPath)

	if errc = self.handleError(path, err); errc != 0 {
		return
	}

	node := self.getOrCreate(path, true)
	node.updateSize(0)

	return
}

func (self *Filescomfs) Unlink(path string) (errc int) {
	defer self.logPanics()

	if node, ok := self.fetch(path); ok {
		// Close the file and wait for any writes to complete before deleting the file.
		node.closeWriter(true)
	}

	// We may have been in the middle of writing the file, so don't log until here.
	localPath, remotePath := self.paths(path)
	self.Info("Deleting file: %v (%v)", remotePath, localPath)

	return self.delete(path)
}

func (self *Filescomfs) Rmdir(path string) (errc int) {
	defer self.logPanics()
	localPath, remotePath := self.paths(path)
	self.Info("Deleting folder: %v (%v)", remotePath, localPath)

	return self.delete(path)
}

func (self *Filescomfs) Rename(oldpath string, newpath string) (errc int) {
	defer self.logPanics()
	oldLocalPath, oldRemotePath := self.paths(oldpath)
	newLocalPath, newRemotePath := self.paths(newpath)
	self.Info("Renaming: %v to %v (%v to %v)", oldRemotePath, newRemotePath, oldLocalPath, newLocalPath)

	params := files_sdk.FileMoveParams{
		Path:        oldRemotePath,
		Destination: newRemotePath,
		Overwrite:   lib.Ptr(true),
	}

	action, err := self.fileClient.Move(params)
	if errc = self.handleError(oldpath, err); errc != 0 {
		return
	}

	err = self.waitForAction(action, "move")
	if errc = self.handleError(oldpath, err); errc != 0 {
		return
	}

	self.rename(oldpath, newpath)

	return
}

func (self *Filescomfs) Utimens(path string, tmsp []fuse.Timespec) (errc int) {
	defer self.logPanics()
	localPath, remotePath := self.paths(path)
	self.Debug("Updating provided mtime: %v (%v) (mtime=%v)", remotePath, localPath, tmsp[1])

	node, _ := self.fetch(path)
	node.info.modTime = tmsp[1].Time()

	if node.isWriterOpen() {
		// If we're writing to the file, no need update the mtime. It will be updated when the write completes.
		// Don't wait for any writes to complete since those will slow down the process, especially when writing many small files.
		node.closeWriter(false)
		return 0
	}

	params := files_sdk.FileUpdateParams{
		Path:          remotePath,
		ProvidedMtime: &node.info.modTime,
	}

	_, err := self.fileClient.Update(params)
	return self.handleError(path, err)
}

func (self *Filescomfs) Create(path string, flags int, mode uint32) (errc int, fh uint64) {
	defer self.logPanics()
	fh = rand.Uint64()
	localPath, remotePath := self.paths(path)
	self.Debug("Creating file: %v (%v) (flags=%v, mode=%v, fh=%v)", remotePath, localPath, flags, mode, fh)

	if errc = self.loadParent(path); errc != 0 {
		return
	}

	node, ok := self.fetch(path)
	if ok && !node.infoExpired() {
		errc = -fuse.EEXIST
		return
	}

	if !ok {
		node = self.getOrCreate(path, false)
	}

	node.updateSize(0)

	if !node.isWriterOpen() {
		self.Info("Starting upload: %v (%v)", remotePath, localPath)
		node.openWriter(self, fh)
	}

	return
}

func (self *Filescomfs) Open(path string, flags int) (errc int, fh uint64) {
	defer self.logPanics()
	fh = rand.Uint64()
	self.Trace("Open: path=%v, flags=%v, fh=%v", path, flags, fh)

	self.getOrCreate(path, false)
	return
}

func (self *Filescomfs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	defer self.logPanics()
	self.Trace("Getattr: path=%v, fh=%v", path, fh)

	if node, ok := self.fetch(path); ok && !node.infoExpired() {
		self.Trace("Getattr: using cached stat, path=%v, size=%v, mtime=%v", path, node.info.size, node.info.modTime)
		getStat(node.info, stat)
		return
	}

	if errc = self.loadParent(path); errc != 0 {
		return
	}

	node, ok := self.fetch(path)
	if !ok || node.infoExpired() {
		localPath, remotePath := self.paths(path)
		self.Debug("File not found: %v (%v)", remotePath, localPath)
		return -fuse.ENOENT
	}

	self.Trace("Getattr: path=%v, size=%v, mtime=%v", path, node.info.size, node.info.modTime)
	getStat(node.info, stat)

	return
}

func (self *Filescomfs) Truncate(path string, size int64, fh uint64) (errc int) {
	defer self.logPanics()
	localPath, remotePath := self.paths(path)
	self.Debug("Truncating file: %v (%v) (size=%v, fh=%v)", remotePath, localPath, size, fh)

	node, _ := self.fetch(path)
	node.updateSize(size)

	if !node.isWriterOpen() {
		self.Info("Starting upload: %v (%v)", remotePath, localPath)
		node.openWriter(self, fh)
	}

	return
}

func (self *Filescomfs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	defer self.logPanics()
	buffLen := int64(len(buff))
	self.Trace("Read: path=%v, len=%v, ofst=%v, fh=%v", path, buffLen, ofst, fh)

	localPath, remotePath := self.paths(path)
	node, _ := self.fetch(path)

	if ofst > node.info.size {
		self.Trace("Read: offset %d is greater than file size %d, returning EOF", ofst, node.info.size)
		return 0
	}

	if node.isWriterOpen() {
		// We can't read while writing to the file, so close the writer and wait for it to finish.
		// Do this _before_ we log that we're starting the download.
		self.Debug("Finalizing upload before downloading: %v (%v) (fh=%v)", remotePath, localPath, fh)
		node.closeWriter(true)
	}

	if ofst == 0 && buffLen >= min(blockSize, node.info.size) {
		node.readerHandle = fh
		self.Info("Starting download: %v (%v)", remotePath, localPath)
	}

	headers := &http.Header{}
	headers.Set("Range", fmt.Sprintf("bytes=%v-%v", ofst, ofst+buffLen-1))
	file, err := self.fileClient.Download(
		files_sdk.FileDownloadParams{File: files_sdk.File{
			Path:        remotePath,
			DownloadUri: node.downloadUri,
		}},
		files_sdk.RequestHeadersOption(headers),
		files_sdk.ResponseBodyOption(func(reader io.ReadCloser) error {
			var err error
			n, err = io.ReadFull(reader, buff)
			if err == io.ErrUnexpectedEOF {
				return nil
			}
			return err
		}),
	)
	if errc := self.handleError(path, err); errc != 0 {
		return errc
	}

	node.downloadUri = file.DownloadUri

	self.Trace("Read: path=%v, ofst=%d, read %d bytes", path, ofst, n)

	return n
}

func (self *Filescomfs) Write(path string, buff []byte, ofst int64, fh uint64) (n int) {
	defer self.logPanics()
	self.Trace("Write: path=%v, len=%v, ofst=%v, fh=%v", path, len(buff), ofst, fh)

	node, _ := self.fetch(path)

	if !node.isWriterOpen() {
		localPath, remotePath := self.paths(path)
		self.Info("Starting upload: %v (%v)", remotePath, localPath)
		node.openWriter(self, fh)
	}

	n, err := node.writer.writeAt(buff, ofst)
	if errc := self.handleError(path, err); errc != 0 {
		return errc
	}

	return n
}

func (self *Filescomfs) Release(path string, fh uint64) (errc int) {
	defer self.logPanics()
	self.Trace("Release: path=%v, fh=%v", path, fh)

	if node, ok := self.fetch(path); ok && node.readerHandle == fh {
		localPath, remotePath := self.paths(path)
		self.Info("Download completed: %v (%v)", remotePath, localPath)
		node.readerHandle = 0
	}

	return self.close(path, fh)
}

func (self *Filescomfs) Opendir(path string) (errc int, fh uint64) {
	defer self.logPanics()
	fh = rand.Uint64()
	self.Trace("Opendir: path=%v, fh=%v", path, fh)

	self.getOrCreate(path, true)
	return
}

func (self *Filescomfs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	defer self.logPanics()
	localPath, remotePath := self.paths(path)
	self.Info("Listing folder: %v (%v)", remotePath, localPath)

	node, _ := self.fetch(path)
	if errc = self.loadDir(node); errc != 0 {
		return
	}

	fill(".", nil, 0)
	fill("..", nil, 0)

	for childPath := range node.childPaths {
		if childNode, ok := self.fetch(childPath); ok {
			fill(path_lib.Base(childPath), getStat(childNode.info, nil), 0)
		}
	}

	return
}

func (self *Filescomfs) Releasedir(path string, fh uint64) (errc int) {
	defer self.logPanics()
	self.Trace("Releasedir: path=%v, fh=%v", path, fh)

	return self.close(path, fh)
}

func (self *Filescomfs) writeFile(path string, reader io.Reader, mtime *time.Time) {
	localPath, remotePath := self.paths(path)
	uploadOpts := []file.UploadOption{
		file.UploadWithDestinationPath(remotePath),
		file.UploadWithReader(reader),
		file.UploadWithProvidedMtimePtr(mtime),
	}
	if self.writeConcurrency != nil {
		uploadOpts = append(uploadOpts, file.UploadWithManager(manager.ConcurrencyManager{}.New(*self.writeConcurrency)))
	}

	if err := self.fileClient.Upload(uploadOpts...); err != nil {
		self.Error("Upload failed: %v (%v): %v", remotePath, localPath, err)
		return
	}

	self.Info("Upload completed: %v (%v)", remotePath, localPath)
}

func (self *Filescomfs) paths(path string) (string, string) {
	return self.localPath(path), self.remotePath(path)
}

func (self *Filescomfs) localPath(path string) string {
	return filepath.Join(self.mountPoint, path)
}

func (self *Filescomfs) remotePath(path string) string {
	return path_lib.Join(self.root, path)
}

func (self *Filescomfs) handleError(path string, err error) int {
	if err != nil {
		localPath, remotePath := self.paths(path)
		self.Error("%v (%v): %v", remotePath, localPath, err)

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

func (self *Filescomfs) delete(path string) (errc int) {
	err := self.fileClient.Delete(files_sdk.FileDeleteParams{Path: self.remotePath(path)})
	if errc = self.handleError(path, err); errc != 0 {
		return
	}

	self.remove(path)
	return
}

func (self *Filescomfs) loadParent(path string) (errc int) {
	if path == "/" {
		// If we're at the root, we can't load the parent. Just make sure the root exists.
		_, errc = self.findDir(path)
		return
	}

	parentPath := path_lib.Dir(path)
	parent, ok := self.fetch(parentPath)

	// Make sure the parent is actually a directory that exists before attempting to load it.
	if !ok || parent.infoExpired() {
		parent, errc = self.findDir(parentPath)
		if errc != 0 {
			return
		}
	}

	if !parent.info.dir {
		// Don't log an error. Windows File Explorer sometimes treats shortcuts as parent directories.
		self.Trace("Parent of %s is not a directory %s", path, parentPath)
		return -fuse.ENOTDIR
	}

	return self.loadDir(parent)
}

func (self *Filescomfs) findDir(path string) (node *fsNode, errc int) {
	remotePath := self.remotePath(path)

	if remotePath == "/" {
		// Special case that we can't stat the root directory of a Files.com site.
		node = self.getOrCreate(path, true)
		node.updateInfo(fsNodeInfo{dir: true})
		return
	}

	item, err := self.fileClient.Find(files_sdk.FileFindParams{Path: remotePath})
	// Check for non-existence first so it doesn't get logged as an error, since this may be expected.
	if files_sdk.IsNotExist(err) {
		errc = -fuse.ENOENT
		return
	}
	if errc = self.handleError(path, err); errc != 0 {
		return
	}
	if !item.IsDir() {
		errc = -fuse.ENOTDIR
		return
	}

	node = self.createNode(path, item)

	return
}

func (self *Filescomfs) loadDir(node *fsNode) (errc int) {
	err := node.updateChildPaths(self.listDir)
	if errc = self.handleError(node.path, err); errc != 0 {
		return
	}

	return
}

func (self *Filescomfs) listDir(path string) (childPaths map[string]struct{}, err error) {
	it, err := self.fileClient.ListFor(files_sdk.FolderListForParams{Path: self.remotePath(path)})
	if err != nil {
		return
	}

	childPaths = make(map[string]struct{})

	for it.Next() {
		item := it.File()

		childPath := path_lib.Join(path, item.DisplayName)
		childPaths[childPath] = struct{}{}

		self.createNode(childPath, item)
	}
	err = it.Err()

	return
}

func (self *Filescomfs) createNode(path string, item files_sdk.File) (node *fsNode) {
	node = self.getOrCreate(path, item.IsDir())
	node.updateInfo(fsNodeInfo{
		dir:          item.IsDir(),
		size:         item.Size,
		modTime:      item.ModTime(),
		creationTime: item.CreatedAt,
	})

	return
}

func (self *Filescomfs) waitForAction(action files_sdk.FileAction, operation string) error {
	migration, err := self.migrationClient.Wait(action, func(migration files_sdk.FileMigration) {
		// noop
	})
	if err == nil && migration.Status != "completed" {
		return fmt.Errorf("%v did not complete successfully: %v", operation, migration.Status)
	}
	return err
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
	stat.Mtim = fuse.NewTimespec(info.modTime)
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
