//go:build windows
// +build windows

package fsmount

import (
	"os"
	"syscall"
	"time"

	"github.com/winfsp/cgofuse/fuse"
)

func (fs *LocalFs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	fq := fs.fqPath(path)
	info, err := os.Stat(fq)
	if err != nil {
		if !os.IsNotExist(err) {
			fs.log.Trace("LocalFs: Getattr: failed to stat file: path=%v, fh=%v, err=%v", fq, fh, err)
			return -fuse.EIO
		}
		fs.log.Trace("LocalFs: Getattr: file not found: path=%v, fh=%v", fq, fh)
		return -fuse.ENOENT
	}

	stat.Mode = uint32(info.Mode())
	stat.Size = int64(info.Size())
	stat.Mtim.Sec = info.ModTime().Unix()
	if sys := info.Sys(); sys != nil {
		if wstat, ok := sys.(*syscall.Win32FileAttributeData); ok {
			sec, nsec := ftToUnix(wstat.LastAccessTime)
			stat.Atim.Sec, stat.Atim.Nsec = sec, nsec
			sec, nsec = ftToUnix(wstat.LastWriteTime)
			stat.Ctim.Sec, stat.Ctim.Nsec = sec, nsec
			sec, nsec = ftToUnix(wstat.CreationTime)
			stat.Birthtim.Sec, stat.Birthtim.Nsec = sec, nsec
		}
	}
	return errc
}

const (
	ticksPerSecond = int64(1e7)
	epochDiff      = int64(11644473600) // seconds between 1601-01-01 and 1970-01-01 Seconds between 1601 and 1970
)

// Convert Windows FILETIME (100-ns ticks since 1601) -> Unix Sec/Nsec
func ftToUnix(ft syscall.Filetime) (sec int64, nsec int64) {
	t := (int64(ft.HighDateTime) << 32) + int64(ft.LowDateTime)
	if t == 0 {
		return 0, 0
	}
	sec = t/ticksPerSecond - epochDiff
	nsec = (t % ticksPerSecond) * 100
	return sec, nsec
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
	if sys := info.Sys(); sys != nil {
		if wstat, ok := sys.(*syscall.Win32FileAttributeData); ok {
			fs.log.Trace("LocalFs: createLocalNode: path=%v, stat=%v", path, wstat)
			sec, nsec := ftToUnix(wstat.CreationTime)
			creationTime = time.Unix(sec, nsec)
		}
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
