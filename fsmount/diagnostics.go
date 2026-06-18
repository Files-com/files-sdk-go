package fsmount

import (
	"fmt"
	"time"

	"github.com/winfsp/cgofuse/fuse"
)

const mountDiagnosticsSlowThreshold = time.Second

func mountCallbackErrc(result int) int {
	if result < 0 {
		return result
	}
	return 0
}

func formatFuseErrno(errc int) string {
	if errc == 0 {
		return "OK"
	}

	code := errc
	if code < 0 {
		code = -code
	}

	switch code {
	case fuse.EACCES:
		return "EACCES"
	case fuse.EAGAIN:
		return "EAGAIN"
	case fuse.EBADF:
		return "EBADF"
	case fuse.EEXIST:
		return "EEXIST"
	case fuse.EIO:
		return "EIO"
	case fuse.EISDIR:
		return "EISDIR"
	case fuse.ENOENT:
		return "ENOENT"
	case fuse.ENOLCK:
		return "ENOLCK"
	case fuse.ENOTDIR:
		return "ENOTDIR"
	case fuse.ENOTEMPTY:
		return "ENOTEMPTY"
	case fuse.ENOSYS:
		return "ENOSYS"
	case fuse.EPERM:
		return "EPERM"
	case fuse.ETIMEDOUT:
		return "ETIMEDOUT"
	default:
		return fmt.Sprintf("errno_%d", code)
	}
}

func (fs *Filescomfs) storageForPath(path string) string {
	if fs.isStoredRemotely(path) {
		return "remote"
	}
	return "local"
}

func (fs *Filescomfs) logMountCallback(op, path, storage string, start time.Time, result int, detailFmt string, detailArgs ...any) {
	duration := time.Since(start)
	errc := mountCallbackErrc(result)

	format := "mount_callback op=%s path=%q storage=%s result=%d errc=%d errno=%s duration=%s"
	args := []any{
		op,
		path,
		storage,
		result,
		errc,
		formatFuseErrno(errc),
		duration,
	}
	if detailFmt != "" {
		format += " " + detailFmt
		args = append(args, detailArgs...)
	}

	fs.log.Trace(format, args...)

	if !shouldLogMountCallback(op, errc, duration) {
		return
	}

	fs.log.Debug(format, args...)
}

func shouldLogMountCallback(op string, errc int, duration time.Duration) bool {
	if duration >= mountDiagnosticsSlowThreshold {
		return true
	}
	if errc == 0 {
		return false
	}

	// macOS Finder and Windows Explorer both probe for optional sidecar and
	// metadata paths during normal browsing. Fast lookup misses are expected
	// noise; slow misses still log because they can explain visible stalls.
	if op == "Getattr" && errc == -fuse.ENOENT {
		return false
	}

	return true
}

func (fs *RemoteFs) logWriteSessionMilestone(path, stage string, fh uint64, session *writeSession, detailFmt string, detailArgs ...any) {
	localPath, remotePath := fs.paths(path)
	details := ""
	if detailFmt != "" {
		details = " " + fmt.Sprintf(detailFmt, detailArgs...)
	}

	if session == nil {
		fs.log.Debug(
			"mount_write_session stage=%s path=%q remote_path=%q local_path=%q fh=%d session=false%s",
			stage,
			path,
			remotePath,
			localPath,
			fh,
			details,
		)
		return
	}

	snap := session.snapshot()
	lastUploadErr := ""
	if snap.lastUploadErr != nil {
		lastUploadErr = snap.lastUploadErr.Error()
	}

	fs.log.Debug(
		"mount_write_session stage=%s path=%q remote_path=%q local_path=%q fh=%d session=true hydrated=%t dirty=%t uploading=%t finalizing=%t current_size=%d handle_count=%d last_upload_err=%q%s",
		stage,
		path,
		remotePath,
		localPath,
		fh,
		snap.hydrated,
		snap.dirty,
		snap.uploading,
		snap.finalizing,
		snap.currentSize,
		snap.handleCount,
		lastUploadErr,
		details,
	)
}
