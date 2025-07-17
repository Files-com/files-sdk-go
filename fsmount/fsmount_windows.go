//go:build windows

package fsmount

import (
	"fmt"
	"os"
)

func mountPoint(mountPoint string) (string, error) {
	if mountPoint == "" {
		// Find the highest available drive letter.
		for l := 'Z'; l >= 'D'; l-- {
			drive := string(l) + ":"
			if _, err := os.Stat(drive + string(os.PathSeparator)); os.IsNotExist(err) {
				return drive, nil
			}
		}

		return "", fmt.Errorf("no available drive letters")
	} else {
		if len(mountPoint) != 2 || mountPoint[1] != ':' {
			return "", fmt.Errorf("invalid mount point")
		}

		_, err := os.Stat(mountPoint + string(os.PathSeparator))
		if err == nil || !os.IsNotExist(err) {
			return "", fmt.Errorf("mount point already in use")
		}
	}

	return mountPoint, nil
}

func mountOpts(params MountParams) []string {
	opts := defaultMountOpts(params)

	if params.DebugFuse {
		logfile := "fuse.log"
		if params.DebugFuseLog != "" {
			logfile = params.DebugFuseLog
		}
		opts = append(opts, "-o", "DebugLog="+logfile)
	}

	opts = append(opts, "-o", "uid=-1")
	opts = append(opts, "-o", "gid=-1")
	if params.VolumeName != "" {
		opts = append(opts, "--VolumePrefix=\\Files.com\\"+params.VolumeName)
	}

	opts = append(opts, "-o", "FileSystemName=Files.com")
	return opts
}
