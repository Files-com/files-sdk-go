//go:build linux
// +build linux

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
	if err := syscall.Lstat(path, &stgo); err != nil {
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
	stat.Atim.Sec, stat.Atim.Nsec = stgo.Atim.Sec, stgo.Atim.Nsec
	stat.Mtim.Sec, stat.Mtim.Nsec = stgo.Mtim.Sec, stgo.Mtim.Nsec
	stat.Ctim.Sec, stat.Ctim.Nsec = stgo.Ctim.Sec, stgo.Ctim.Nsec
	stat.Blksize = int64(stgo.Blksize)
	stat.Blocks = int64(stgo.Blocks)

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
	// TODO: double check this on linux
	// On linux, os.FileInfo does not provide creation time
	node := fs.vfs.getOrCreate(path, nt)
	node.updateInfo(fsNodeInfo{
		nodeType:     nt,
		size:         info.Size(),
		modTime:      info.ModTime(),
		creationTime: creationTime,
	})

	return node, nil
}
