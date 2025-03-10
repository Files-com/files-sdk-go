//go:build windows

package fsmount

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	filepath "path"
	"slices"
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
	vfs              virtualfs
	root             string
	writeConcurrency *int
	config           files_sdk.Config
	fileClient       *file.Client
	migrationClient  *file_migration.Client
}

func (self *Filescomfs) Init() {
	self.fileClient = &file.Client{Config: self.config}
	self.migrationClient = &file_migration.Client{Config: self.config}
	self.vfs.init()
}

func (self *Filescomfs) Statfs(path string, stat *fuse.Statfs_t) (errc int) {
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

func (self *Filescomfs) Mknod(path string, mode uint32, dev uint64) (errc int) {
	path = self.absPath(path)
	self.trace("Mknod: path=%v, mode=%v, dev=%v", path, mode, dev)

	if errc = self.loadParent(path); errc != 0 {
		return
	}

	node, ok := self.vfs.fetch(path)
	if ok && !node.infoExpired() {
		return -fuse.EEXIST
	}

	if !ok {
		node = self.vfs.getOrCreate(path, false)
	}

	node.openWriter(self)

	return 0
}

func (self *Filescomfs) Mkdir(path string, mode uint32) (errc int) {
	path = self.absPath(path)
	self.trace("Mkdir: path=%v, mode=%v", path, mode)

	_, err := self.fileClient.CreateFolder(files_sdk.FolderCreateParams{Path: path})
	if files_sdk.IsExist(err) {
		return 0
	}
	if errc = self.handleError(err); errc != 0 {
		return
	}

	self.vfs.getOrCreate(path, true)

	return
}

func (self *Filescomfs) Unlink(path string) (errc int) {
	path = self.absPath(path)
	self.trace("Unlink: path=%v", path)

	if node, ok := self.vfs.fetch(path); ok {
		node.closeWriter()
		// Wait for the upload to complete before deleting the file.
		node.waitForUploadCompletion()
	}

	return self.delete(path)
}

func (self *Filescomfs) Rmdir(path string) (errc int) {
	path = self.absPath(path)
	self.trace("Rmdir: path=%v", path)

	return self.delete(path)
}

func (self *Filescomfs) Rename(oldpath string, newpath string) (errc int) {
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

	self.vfs.rename(oldpath, newpath)

	return
}

func (self *Filescomfs) Utimens(path string, tmsp []fuse.Timespec) (errc int) {
	path = self.absPath(path)
	self.trace("Utimens: path=%v, tmsp=%v", path, tmsp)

	node, _ := self.vfs.fetch(path)
	node.info.modTime = tmsp[1].Time()

	if node.writer != nil {
		// If we're writing to the file, no need update the mtime. It will be updated when the write completes.
		node.closeWriter()
		return 0
	}

	params := files_sdk.FileUpdateParams{
		Path:          path,
		ProvidedMtime: &node.info.modTime,
	}

	_, err := self.fileClient.Update(params)
	return self.handleError(err)
}

func (self *Filescomfs) Open(path string, flags int) (errc int, fh uint64) {
	path = self.absPath(path)
	self.trace("Open: path=%v, flags=%v", path, flags)

	self.vfs.getOrCreate(path, false)
	return
}

func (self *Filescomfs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	path = self.absPath(path)
	self.trace("Getattr: path=%v", path)

	if path == "/" {
		// Special case that we can't stat the root directory of a Files.com site.
		stat.Mode = fuse.S_IFDIR | 0777
		return
	}

	if node, ok := self.vfs.fetch(path); ok && !node.infoExpired() {
		self.trace("Getattr: using cached stat, path=%v, size=%v, mtime=%v", path, node.info.size, node.info.modTime)
		getStat(node.info, stat)
		return
	}

	if errc = self.loadParent(path); errc != 0 {
		return
	}

	node, ok := self.vfs.fetch(path)
	if !ok || node.infoExpired() {
		self.trace("Getattr: path=%v not found", path)
		return -fuse.ENOENT
	}

	getStat(node.info, stat)

	return
}

// This is needed in order to support file overwrites, but we don't actually need to do the truncate.
func (self *Filescomfs) Truncate(path string, size int64, fh uint64) (errc int) {
	path = self.absPath(path)
	self.trace("Truncate: path=%v, size=%v", path, size)

	return
}

func (self *Filescomfs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	path = self.absPath(path)
	self.trace("Read: path=%v, len=%v, ofst=%v", path, len(buff), ofst)

	node, _ := self.vfs.fetch(path)

	if ofst >= node.info.size {
		self.trace("Read: offset %d is greater than file size %d, returning EOF", ofst, node.info.size)
		return 0
	}

	headers := &http.Header{}
	headers.Set("Range", fmt.Sprintf("bytes=%v-%v", ofst, ofst+int64(len(buff))-1))
	_, err := self.fileClient.Download(
		files_sdk.FileDownloadParams{Path: path},
		files_sdk.RequestHeadersOption(headers),
		files_sdk.ResponseBodyOption(func(closer io.ReadCloser) error {
			var err error
			n, err = io.ReadFull(closer, buff)
			if err == io.ErrUnexpectedEOF {
				return nil
			}
			return err
		}),
	)
	if errc := self.handleError(err); errc != 0 {
		return errc
	}

	return n
}

func (self *Filescomfs) Write(path string, buff []byte, ofst int64, fh uint64) (n int) {
	path = self.absPath(path)
	self.trace("Write: path=%v, len=%v, ofst=%v", path, len(buff), ofst)

	node, _ := self.vfs.fetch(path)

	if ofst < node.writeOffset {
		// This happens on Windows when a write operation is paused. It writes a 56 byte buffer at
		// offset 0. It's unclear how to handle this to properly resume the write.
		self.trace("Write: path=%v, offset %d is less than write offset %d, closing writer", path, ofst, node.writeOffset)
		node.closeWriter()
		return len(buff)
	}

	node.openWriter(self)

	if ofst > node.writeOffset {
		// Sometimes parts come in out of order. We need to cache them until it's time to write them.
		self.trace("Write: path=%v, offset %d is greater than write offset %d, caching", path, ofst, node.writeOffset)
		// TODO: Allow for configuring the cache size.
		node.partCache[ofst] = slices.Clone(buff)
		// Return that we wrote the full buffer, otherwise fuse will eventually fail the write.
		return len(buff)
	}

	n, err := node.write(buff)
	if errc := self.handleError(err); errc != 0 {
		return errc
	}

	self.trace("Write: path=%v, wrote %d bytes, new write offset is %d", path, n, node.writeOffset)

	for {
		part, ok := node.partCache[node.writeOffset]
		if !ok {
			break
		}

		l, err := node.write(part)
		if errc := self.handleError(err); errc != 0 {
			return errc
		}

		self.trace("Write: path=%v, wrote %d bytes, new write offset is %d", path, l, node.writeOffset)
	}

	return n
}

func (self *Filescomfs) Release(path string, fh uint64) (errc int) {
	path = self.absPath(path)
	self.trace("Release: path=%v", path)

	return self.vfs.close(path)
}

func (self *Filescomfs) Opendir(path string) (errc int, fh uint64) {
	path = self.absPath(path)
	self.trace("Opendir: path=%v", path)

	self.vfs.getOrCreate(path, true)
	return
}

func (self *Filescomfs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	path = self.absPath(path)
	self.trace("Readdir: path=%v", path)

	node, _ := self.vfs.fetch(path)
	if errc = self.loadDir(node); errc != 0 {
		return
	}

	fill(".", nil, 0)
	fill("..", nil, 0)

	for childPath := range node.childPaths {
		if childNode, ok := self.vfs.fetch(childPath); ok {
			fill(filepath.Base(childPath), getStat(childNode.info, nil), 0)
		}
	}

	return
}

func (self *Filescomfs) Releasedir(path string, fh uint64) (errc int) {
	path = self.absPath(path)
	self.trace("Releasedir: path=%v", path)

	return self.vfs.close(path)
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

	self.vfs.delete(path)
	return
}

func (self *Filescomfs) loadParent(path string) (errc int) {
	parentPath := filepath.Dir(path)
	parent := self.vfs.getOrCreate(parentPath, true)
	return self.loadDir(parent)
}

func (self *Filescomfs) loadDir(node *fsNode) (errc int) {
	if !node.childPathsExpired() {
		return
	}

	it, err := self.fileClient.ListFor(files_sdk.FolderListForParams{Path: node.path})
	if errc = self.handleError(err); errc != 0 {
		return
	}

	childPaths := make(map[string]struct{})

	for it.Next() {
		item := it.File()

		childPath := filepath.Join(node.path, item.DisplayName)
		childPaths[childPath] = struct{}{}

		childNode := self.vfs.getOrCreate(childPath, item.IsDir())
		childNode.updateInfo(fsNodeInfo{
			dir:          item.IsDir(),
			size:         item.Size,
			modTime:      item.ModTime(),
			creationTime: item.CreatedAt,
		})
	}
	if errc = self.handleError(it.Err()); errc != 0 {
		return
	}

	node.updateChildPaths(childPaths)

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

func (self *Filescomfs) error(format string, args ...any) {
	self.log("ERROR", format, args...)
}

func (self *Filescomfs) trace(format string, args ...any) {
	self.log("TRACE", format, args...)
}

func (self *Filescomfs) log(level string, format string, args ...any) {
	format = fmt.Sprintf("[%v] %v", level, format)
	self.config.Printf(format, args...)
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
