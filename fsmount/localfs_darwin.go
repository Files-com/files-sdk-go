//go:build darwin
// +build darwin

package fsmount

import (
	"os"
	"syscall"
	"time"

	"github.com/winfsp/cgofuse/fuse"
)

func (fs *LocalFs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	fq := fs.fqPath(path)
	stgo := syscall.Stat_t{}
	if err := syscall.Lstat(fq, &stgo); err != nil {
		if !os.IsNotExist(err) {
			fs.log.Trace("LocalFs: Getattr: failed to lstat file: path=%v, fh=%v, err=%v", fq, fh, err)
			return -fuse.EIO
		}
		fs.log.Trace("LocalFs: Getattr: file not found: path=%v, fh=%v", fq, fh)
		return -fuse.ENOENT
	}
	stat.Dev = uint64(stgo.Dev)
	stat.Ino = uint64(stgo.Ino)
	stat.Mode = uint32(stgo.Mode)
	stat.Nlink = uint32(stgo.Nlink)
	stat.Uid = uint32(stgo.Uid)
	stat.Gid = uint32(stgo.Gid)
	stat.Rdev = uint64(stgo.Rdev)
	stat.Size = int64(stgo.Size)
	stat.Atim.Sec, stat.Atim.Nsec = stgo.Atimespec.Sec, stgo.Atimespec.Nsec
	stat.Mtim.Sec, stat.Mtim.Nsec = stgo.Mtimespec.Sec, stgo.Mtimespec.Nsec
	stat.Ctim.Sec, stat.Ctim.Nsec = stgo.Ctimespec.Sec, stgo.Ctimespec.Nsec
	stat.Blksize = int64(stgo.Blksize)
	stat.Blocks = int64(stgo.Blocks)
	stat.Birthtim.Sec, stat.Birthtim.Nsec = stgo.Birthtimespec.Sec, stgo.Birthtimespec.Nsec

	fs.log.Trace("LocalFs: Getattr: path=%v, stat=%v, fh=%v", fq, stat, fh)
	return errc
}

func (fs *LocalFs) createLocalNode(path string, entry os.DirEntry) (*fsNode, error) {
	var nt nodeType
	if entry.IsDir() {
		nt = nodeTypeDir
	} else {
		nt = nodeTypeFile
	}
	info, err := entry.Info()
	if err != nil {
		return nil, err
	}
	var creationTime time.Time
	// On unix, os.FileInfo does not provide creation time, so we need to use syscall.Stat_t
	// to get the creation time if available.
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		creationTime = time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)
	}
	node := fs.vfs.getOrCreate(path, nt)
	node.updateInfo(fsNodeInfo{
		nodeType:     nt,
		size:         info.Size(),
		modTime:      info.ModTime(),
		creationTime: creationTime,
	})

	return node, nil
}
