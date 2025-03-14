//go:build windows

package fsmount

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	filepath "path"
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
)

type Filescomfs struct {
	fuse.FileSystemBase
	*virtualfs
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
	self.virtualfs = newVirtualfs(self.config, self.cacheTTL)
}

func (self *Filescomfs) Statfs(path string, stat *fuse.Statfs_t) (errc int) {
	defer self.logPanics()
	path = self.absPath(path)
	self.trace("Statfs: path=%v", path)

	blockSize := uint64(4096)
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
	path = self.absPath(path)
	self.trace("Mkdir: path=%v, mode=%v", path, mode)

	_, err := self.fileClient.CreateFolder(files_sdk.FolderCreateParams{Path: path})
	if files_sdk.IsExist(err) {
		return 0
	}
	if errc = self.handleError(err); errc != 0 {
		return
	}

	node := self.getOrCreate(path, true)
	node.updateSize(0)

	return
}

func (self *Filescomfs) Unlink(path string) (errc int) {
	defer self.logPanics()
	path = self.absPath(path)
	self.trace("Unlink: path=%v", path)

	if node, ok := self.fetch(path); ok {
		// Close the file and wait for any writes to complete before deleting the file.
		node.closeWriter(true)
	}

	return self.delete(path)
}

func (self *Filescomfs) Rmdir(path string) (errc int) {
	defer self.logPanics()
	path = self.absPath(path)
	self.trace("Rmdir: path=%v", path)

	return self.delete(path)
}

func (self *Filescomfs) Rename(oldpath string, newpath string) (errc int) {
	defer self.logPanics()
	oldpath = self.absPath(oldpath)
	newpath = self.absPath(newpath)
	self.trace("Rename: oldpath=%v, newpath=%v", oldpath, newpath)

	params := files_sdk.FileMoveParams{
		Path:        oldpath,
		Destination: newpath,
		Overwrite:   lib.Ptr(true),
	}

	action, err := self.fileClient.Move(params)
	if errc = self.handleError(err); errc != 0 {
		return
	}

	err = self.waitForAction(action, "move")
	if errc = self.handleError(err); errc != 0 {
		return
	}

	self.rename(oldpath, newpath)

	return
}

func (self *Filescomfs) Utimens(path string, tmsp []fuse.Timespec) (errc int) {
	defer self.logPanics()
	path = self.absPath(path)
	self.trace("Utimens: path=%v, tmsp=%v", path, tmsp)

	node, _ := self.fetch(path)
	node.info.modTime = tmsp[1].Time()

	if node.writer != nil {
		// If we're writing to the file, no need update the mtime. It will be updated when the write completes.
		// Don't wait for any writes to complete since those will slow down the process, especially when writing many small files.
		node.closeWriter(false)
		return 0
	}

	params := files_sdk.FileUpdateParams{
		Path:          path,
		ProvidedMtime: &node.info.modTime,
	}

	_, err := self.fileClient.Update(params)
	return self.handleError(err)
}

func (self *Filescomfs) Create(path string, flags int, mode uint32) (errc int, fh uint64) {
	defer self.logPanics()
	path = self.absPath(path)
	fh = rand.Uint64()
	self.trace("Create: path=%v, flags=%v, mode=%v, fh=%v", path, flags, mode, fh)

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
	node.openWriter(self, fh)

	return
}

func (self *Filescomfs) Open(path string, flags int) (errc int, fh uint64) {
	defer self.logPanics()
	path = self.absPath(path)
	fh = rand.Uint64()
	self.trace("Open: path=%v, flags=%v, fh=%v", path, flags, fh)

	self.getOrCreate(path, false)
	return
}

func (self *Filescomfs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	defer self.logPanics()
	path = self.absPath(path)
	self.trace("Getattr: path=%v", path)

	if path == "/" {
		// Special case that we can't stat the root directory of a Files.com site.
		stat.Mode = fuse.S_IFDIR | 0777
		return
	}

	if node, ok := self.fetch(path); ok && !node.infoExpired() {
		self.trace("Getattr: using cached stat, path=%v, size=%v, mtime=%v", path, node.info.size, node.info.modTime)
		getStat(node.info, stat)
		return
	}

	if errc = self.loadParent(path); errc != 0 {
		return
	}

	node, ok := self.fetch(path)
	if !ok || node.infoExpired() {
		self.trace("Getattr: path=%v, not found", path)
		return -fuse.ENOENT
	}

	self.trace("Getattr: path=%v, size=%v, mtime=%v", path, node.info.size, node.info.modTime)
	getStat(node.info, stat)

	return
}

func (self *Filescomfs) Truncate(path string, size int64, fh uint64) (errc int) {
	defer self.logPanics()
	path = self.absPath(path)
	self.trace("Truncate: path=%v, size=%v, fh=%v", path, size, fh)

	node, _ := self.fetch(path)
	node.updateSize(size)
	node.openWriter(self, fh)

	return
}

func (self *Filescomfs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	defer self.logPanics()
	path = self.absPath(path)
	self.trace("Read: path=%v, len=%v, ofst=%v, fh=%v", path, len(buff), ofst, fh)

	node, _ := self.fetch(path)

	if ofst >= node.info.size {
		self.trace("Read: offset %d is greater than file size %d, returning EOF", ofst, node.info.size)
		return 0
	}

	headers := &http.Header{}
	headers.Set("Range", fmt.Sprintf("bytes=%v-%v", ofst, ofst+int64(len(buff))-1))
	file, err := self.fileClient.Download(
		files_sdk.FileDownloadParams{File: files_sdk.File{
			Path:        node.path,
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
	if errc := self.handleError(err); errc != 0 {
		return errc
	}

	node.downloadUri = file.DownloadUri

	self.trace("Read: path=%v, ofst=%d, read %d bytes", path, ofst, n)

	return n
}

func (self *Filescomfs) Write(path string, buff []byte, ofst int64, fh uint64) (n int) {
	defer self.logPanics()
	path = self.absPath(path)
	self.trace("Write: path=%v, len=%v, ofst=%v, fh=%v", path, len(buff), ofst, fh)

	node, _ := self.fetch(path)

	if !node.openWriter(self, fh) {
		self.error("Write: path=%v, fh=%v, writer is already open for a different handle", path)
		return -fuse.EIO
	}

	n, err := node.writer.writeAt(buff, ofst)
	if errc := self.handleError(err); errc != 0 {
		return errc
	}

	return n
}

func (self *Filescomfs) Release(path string, fh uint64) (errc int) {
	defer self.logPanics()
	path = self.absPath(path)
	self.trace("Release: path=%v, fh=%v", path, fh)

	return self.close(path, fh)
}

func (self *Filescomfs) Opendir(path string) (errc int, fh uint64) {
	defer self.logPanics()
	path = self.absPath(path)
	fh = rand.Uint64()
	self.trace("Opendir: path=%v, fh=%v", path, fh)

	self.getOrCreate(path, true)
	return
}

func (self *Filescomfs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	defer self.logPanics()
	path = self.absPath(path)
	self.trace("Readdir: path=%v, ofst=%v, fh=%v", path, ofst, fh)

	node, _ := self.fetch(path)
	if errc = self.loadDir(node); errc != 0 {
		return
	}

	fill(".", nil, 0)
	fill("..", nil, 0)

	for childPath := range node.childPaths {
		if childNode, ok := self.fetch(childPath); ok {
			fill(filepath.Base(childPath), getStat(childNode.info, nil), 0)
		}
	}

	return
}

func (self *Filescomfs) Releasedir(path string, fh uint64) (errc int) {
	defer self.logPanics()
	path = self.absPath(path)
	self.trace("Releasedir: path=%v, fh=%v", path, fh)

	return self.close(path, fh)
}

func (self *Filescomfs) writeFile(path string, reader io.Reader, mtime *time.Time) {
	uploadOpts := []file.UploadOption{
		file.UploadWithDestinationPath(path),
		file.UploadWithReader(reader),
		file.UploadWithProvidedMtimePtr(mtime),
	}
	if self.writeConcurrency != nil {
		uploadOpts = append(uploadOpts, file.UploadWithManager(manager.ConcurrencyManager{}.New(*self.writeConcurrency)))
	}

	if err := self.fileClient.Upload(uploadOpts...); err != nil {
		self.error("Upload failed: %v", err)
	}
}

func (self *Filescomfs) absPath(path string) string {
	return filepath.Join(self.root, path)
}

func (self *Filescomfs) handleError(err error) int {
	if err != nil {
		self.error(err.Error())

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
	err := self.fileClient.Delete(files_sdk.FileDeleteParams{Path: path})
	if errc = self.handleError(err); errc != 0 {
		return
	}

	self.remove(path)
	return
}

func (self *Filescomfs) loadParent(path string) (errc int) {
	parentPath := filepath.Dir(path)
	parent := self.getOrCreate(parentPath, true)
	return self.loadDir(parent)
}

func (self *Filescomfs) loadDir(node *fsNode) (errc int) {
	err := node.updateChildPaths(self.listDir)
	if errc = self.handleError(err); errc != 0 {
		return
	}

	return
}

func (self *Filescomfs) listDir(path string) (childPaths map[string]struct{}, err error) {
	it, err := self.fileClient.ListFor(files_sdk.FolderListForParams{Path: path})
	if err != nil {
		return
	}

	childPaths = make(map[string]struct{})

	for it.Next() {
		item := it.File()

		childPath := filepath.Join(path, item.DisplayName)
		childPaths[childPath] = struct{}{}

		childNode := self.getOrCreate(childPath, item.IsDir())
		childNode.updateInfo(fsNodeInfo{
			dir:          item.IsDir(),
			size:         item.Size,
			modTime:      item.ModTime(),
			creationTime: item.CreatedAt,
		})
	}
	err = it.Err()

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
