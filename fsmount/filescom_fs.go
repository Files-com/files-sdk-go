//go:build windows

package fsmount

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	filepath "path"
	"slices"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	file_migration "github.com/Files-com/files-sdk-go/v3/filemigration"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/winfsp/cgofuse/fuse"
)

const (
	folderNotEmpty = "processing-failure/folder-not-empty"
	statCacheTime  = 1 * time.Second
)

type fsNode struct {
	fs          *Filescomfs
	path        string
	openCount   int
	stat        *fuse.Stat_t
	statExpires *time.Time
	mu          sync.Mutex
	writer      *io.PipeWriter
	reader      *io.PipeReader
	writeOffset int64
	partCache   map[int64][]byte
}

func (n *fsNode) updateStat(stat *fuse.Stat_t) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.stat = stat
	n.statExpires = lib.Time(time.Now().Add(statCacheTime))
}

func (n *fsNode) openWriter() {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.writer == nil {
		n.reader, n.writer = io.Pipe()
		n.partCache = make(map[int64][]byte)
		go func() {
			uploadOpts := []file.UploadOption{
				file.UploadWithReader(n.reader),
				file.UploadWithDestinationPath(n.path),
			}

			if err := n.fs.fileClient.Upload(uploadOpts...); err != nil {
				n.fs.error("Upload failed: %v", err)
			}

			n.reader.Close()
			n.reader = nil
		}()
	}
}

func (n *fsNode) write(buff []byte) (int, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	l, err := n.writer.Write(buff)
	if err != nil {
		return 0, err
	}

	// Remove the part from the cache. No-op if it's not in the cache.
	delete(n.partCache, n.writeOffset)

	n.writeOffset += int64(l)
	n.stat.Size = n.writeOffset
	n.statExpires = lib.Time(time.Now().Add(statCacheTime))

	return l, nil
}

func (n *fsNode) closeWriter() {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.writer != nil {
		n.writer.Close()
		n.writer = nil
		n.writeOffset = 0
		n.statExpires = nil
		n.partCache = nil

		// Wait for the reader to be closed so we know the upload is complete.
		for n.reader != nil {
			time.Sleep(1 * time.Second) // TODO: make this better
		}
	}
}

type Filescomfs struct {
	fuse.FileSystemBase
	root            string
	config          files_sdk.Config
	fileClient      *file.Client
	migrationClient *file_migration.Client
	openMap         map[string]*fsNode
	openMapMutex    sync.Mutex
}

func (self *Filescomfs) Init() {
	self.fileClient = &file.Client{Config: self.config}
	self.migrationClient = &file_migration.Client{Config: self.config}
	self.openMap = make(map[string]*fsNode)
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

	_, err := self.fileClient.Find(files_sdk.FileFindParams{Path: path})
	if err == nil {
		return -fuse.EEXIST
	}
	if !files_sdk.IsNotExist(err) {
		if errc = self.handleError(err); errc != 0 {
			return
		}
	}

	node := self.fetchNode(path)
	node.openWriter()

	return 0
}

func (self *Filescomfs) Mkdir(path string, mode uint32) (errc int) {
	path = self.absPath(path)
	self.trace("Mkdir: path=%v, mode=%v", path, mode)

	_, err := self.fileClient.CreateFolder(files_sdk.FolderCreateParams{Path: path})
	if files_sdk.IsExist(err) {
		return 0
	}
	return self.handleError(err)
}

func (self *Filescomfs) Unlink(path string) (errc int) {
	path = self.absPath(path)
	self.trace("Unlink: path=%v", path)

	node := self.fetchNode(path)
	node.closeWriter()

	file, err := self.fileClient.Find(files_sdk.FileFindParams{Path: path})
	if errc = self.handleError(err); errc != 0 {
		return
	}
	if file.IsDir() {
		return -fuse.EISDIR
	}

	err = self.fileClient.Delete(files_sdk.FileDeleteParams{Path: path})
	return self.handleError(err)
}

func (self *Filescomfs) Rmdir(path string) (errc int) {
	path = self.absPath(path)
	self.trace("Rmdir: path=%v", path)

	file, err := self.fileClient.Find(files_sdk.FileFindParams{Path: path})
	if files_sdk.IsNotExist(err) {
		return -fuse.ENOENT
	}
	if !file.IsDir() {
		return -fuse.ENOTDIR
	}

	params := files_sdk.FileDeleteParams{
		Path: path,
	}

	err = self.fileClient.Delete(params)
	return self.handleError(err)
}

func (self *Filescomfs) Rename(oldpath string, newpath string) (errc int) {
	oldpath = self.absPath(oldpath)
	newpath = self.absPath(newpath)
	self.trace("Rename: oldpath=%v, newpath=%v", oldpath, newpath)

	params := files_sdk.FileMoveParams{
		Path:        oldpath,
		Destination: newpath,
		Overwrite:   lib.Bool(true),
	}

	action, err := self.fileClient.Move(params)
	if errc = self.handleError(err); errc != 0 {
		return
	}

	err = self.waitForAction(action, "move")
	return self.handleError(err)
}

func (self *Filescomfs) Utimens(path string, tmsp []fuse.Timespec) (errc int) {
	path = self.absPath(path)
	self.trace("Utimens: path=%v, tmsp=%v", path, tmsp)

	node := self.fetchNode(path)
	node.closeWriter()

	params := files_sdk.FileUpdateParams{
		Path:          path,
		ProvidedMtime: lib.Time(tmsp[1].Time()),
	}

	_, err := self.fileClient.Update(params)
	return self.handleError(err)
}

func (self *Filescomfs) Open(path string, flags int) (errc int, fh uint64) {
	path = self.absPath(path)
	self.trace("Open: path=%v, flags=%v", path, flags)

	return self.openNode(path, false)
}

func (self *Filescomfs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	path = self.absPath(path)
	self.trace("Getattr: path=%v", path)

	node := self.fetchNode(path)

	if fuse.S_IFDIR == node.stat.Mode&fuse.S_IFMT || (node.statExpires != nil && node.statExpires.After(time.Now())) || node.writer != nil {
		self.trace("Getattr: using cached stat, path=%v, size=%v, mtime=%v", path, node.stat.Size, node.stat.Mtim)
		*stat = *node.stat
		return 0
	}

	item, err := self.fileClient.Find(files_sdk.FileFindParams{Path: path})
	if errc = self.handleError(err); errc != 0 {
		return
	}

	getItemMetadata(item, stat)
	node.updateStat(stat)

	return 0
}

// TODO: this is needed in order to support file overwrites, but do we need to actually truncate to the given size?
func (self *Filescomfs) Truncate(path string, size int64, fh uint64) (errc int) {
	path = self.absPath(path)
	self.trace("Truncate: path=%v, size=%v", path, size)

	return 0
}

func (self *Filescomfs) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	path = self.absPath(path)
	self.trace("Read: path=%v, len=%v, ofst=%v", path, len(buff), ofst)

	node := self.fetchNode(path)

	if ofst >= node.stat.Size {
		self.trace("Read: offset %d is greater than file size %d, returning EOF", ofst, node.stat.Size)
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

	node := self.fetchNode(path)

	if ofst < node.writeOffset {
		// This happens on Windows when a write operation is paused. It writes a 56 byte buffer at
		// offset 0. It's unclear how to handle this to properly resume the write.
		self.trace("Write: path=%v, offset %d is less than write offset %d, closing writer", path, ofst, node.writeOffset)
		node.closeWriter()
		return len(buff)
	}

	node.openWriter()

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

	return self.closeNode(path)
}

func (self *Filescomfs) Opendir(path string) (errc int, fh uint64) {
	path = self.absPath(path)
	self.trace("Opendir: path=%v", path)

	return self.openNode(path, true)
}

func (self *Filescomfs) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	path = self.absPath(path)
	self.trace("Readdir: path=%v", path)

	it, err := self.fileClient.ListFor(files_sdk.FolderListForParams{Path: path})
	if errc = self.handleError(err); errc != 0 {
		return
	}

	fill(".", nil, 0)
	fill("..", nil, 0)

	for it.Next() {
		item := it.File()
		if !fill(item.DisplayName, getItemMetadata(item, nil), 0) {
			break
		}
	}

	return self.handleError(it.Err())
}

func (self *Filescomfs) Releasedir(path string, fh uint64) (errc int) {
	path = self.absPath(path)
	self.trace("Releasedir: path=%v", path)

	return self.closeNode(path)
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

func (self *Filescomfs) openNode(path string, dir bool) (errc int, fh uint64) {
	node := self.fetchNode(path)
	if node.openCount == 0 {
		if dir {
			node.stat.Mode = fuse.S_IFDIR | 0777
		} else {
			node.stat.Mode = fuse.S_IFREG | 0777
		}
	}

	if dir && fuse.S_IFDIR != node.stat.Mode&fuse.S_IFMT {
		self.error("%v is not a directory", path)
		return -fuse.ENOTDIR, ^uint64(0)
	} else if !dir && fuse.S_IFDIR == node.stat.Mode&fuse.S_IFMT {
		self.error("%v is a directory", path)
		return -fuse.EISDIR, ^uint64(0)
	}

	node.openCount++
	return
}

func (self *Filescomfs) fetchNode(path string) (node *fsNode) {
	self.openMapMutex.Lock()
	defer self.openMapMutex.Unlock()

	node, ok := self.openMap[path]
	if !ok {
		stat := &fuse.Stat_t{}
		node = &fsNode{
			fs:   self,
			path: path,
			stat: stat,
		}
		self.openMap[path] = node

		if path == "/" {
			stat.Mode = fuse.S_IFDIR | 0777
		}
	}

	return node
}

func (self *Filescomfs) closeNode(path string) (errc int) {
	self.openMapMutex.Lock()
	defer self.openMapMutex.Unlock()

	node, ok := self.openMap[path]
	if !ok {
		self.error("file not found in open map: %v", path)
		return -fuse.EBADF
	}

	node.closeWriter()

	node.openCount--
	if node.openCount < 0 {
		self.error("openCount is negative: %v", path)
		return -fuse.EBADF
	}

	// TODO: Remove closed nodes that haven't been accessed in a while.

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

func getItemMetadata(item files_sdk.File, stat *fuse.Stat_t) *fuse.Stat_t {
	if stat == nil {
		stat = &fuse.Stat_t{}
	}

	if item.IsDir() {
		stat.Mode = fuse.S_IFDIR | 0777
	} else {
		stat.Mode = fuse.S_IFREG | 0777
	}

	stat.Size = item.Size
	stat.Mtim = fuse.NewTimespec(item.ModTime())
	if item.CreatedAt != nil {
		stat.Birthtim = fuse.NewTimespec(*item.CreatedAt)
	}

	return stat
}

func isFolderNotEmpty(err error) bool {
	var re files_sdk.ResponseError
	ok := errors.As(err, &re)
	return ok && re.Type == folderNotEmpty
}
