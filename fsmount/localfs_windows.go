//go:build windows
// +build windows

package fsmount

import (
	"os"
	"syscall"
	"time"

	"github.com/winfsp/cgofuse/fuse"
	"golang.org/x/sys/windows"
)

func (fs *LocalFs) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	fs.vfs.ensureContextOwner()
	if stat == nil {
		stat = &fuse.Stat_t{}
	}
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
	node, ok := fs.vfs.fetch(path)
	if ok && node.info.uid != 0 {
		stat.Uid = node.info.uid
	} else {
		stat.Uid = fs.vfs.uid
	}
	if ok && node.info.gid != 0 {
		stat.Gid = node.info.gid
	} else {
		stat.Gid = fs.vfs.gid
	}
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

func openLocalFile(path string, flags int, mode os.FileMode) (*os.File, error) {
	pathPtr, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}

	access := uint32(0)
	switch flags & (os.O_RDONLY | os.O_WRONLY | os.O_RDWR) {
	case os.O_WRONLY:
		access = windows.GENERIC_WRITE
	case os.O_RDWR:
		access = windows.GENERIC_READ | windows.GENERIC_WRITE
	default:
		access = windows.GENERIC_READ
	}

	if flags&os.O_APPEND != 0 {
		access |= windows.FILE_APPEND_DATA
	}

	shareMode := uint32(windows.FILE_SHARE_READ | windows.FILE_SHARE_WRITE | windows.FILE_SHARE_DELETE)

	var creationDisposition uint32
	switch {
	case flags&(os.O_CREATE|os.O_EXCL) == (os.O_CREATE | os.O_EXCL):
		creationDisposition = windows.CREATE_NEW
	case flags&(os.O_CREATE|os.O_TRUNC) == (os.O_CREATE | os.O_TRUNC):
		creationDisposition = windows.CREATE_ALWAYS
	case flags&os.O_CREATE != 0:
		creationDisposition = windows.OPEN_ALWAYS
	case flags&os.O_TRUNC != 0:
		creationDisposition = windows.TRUNCATE_EXISTING
	default:
		creationDisposition = windows.OPEN_EXISTING
	}

	handle, err := windows.CreateFile(
		pathPtr,
		access,
		shareMode,
		nil,
		creationDisposition,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return nil, err
	}

	return os.NewFile(uintptr(handle), path), nil
}
