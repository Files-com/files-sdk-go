//go:build windows

package fsmount

import (
	"fmt"
	"os"
	"regexp"
)

var (
	validMountPointRe = regexp.MustCompile(`^[A-Za-z]:$`)
)

// mountPoint validates or selects a mount point (drive letter) for the FUSE file system.
// If mountPoint is empty, it finds the highest available drive letter by searching backward
// from Z: to D:.
// If mountPoint is provided, it validates that it is a single drive letter followed by a colon (e.g. "X:")
// and checks that the drive letter is not already in use.
// Returns the selected or validated mount point, or an error if no valid mount point is found.
func mountPoint(mountPoint string) (string, error) {
	if err := validateMountPoint(mountPoint); err != nil {
		return "", err
	}
	if mountPoint == "" {
		// Find the highest available drive letter.
		for l := 'Z'; l >= 'D'; l-- {
			drive := string(l) + ":"
			// check if the candidate is already reserved
			if _, ok := mntRegistry.get(drive); ok {
				continue
			}
			if _, err := os.Stat(drive + string(os.PathSeparator)); os.IsNotExist(err) {
				return drive, nil
			}
		}

		return "", fmt.Errorf("no available drive letters")
	} else {
		_, err := os.Stat(mountPoint + string(os.PathSeparator))
		switch {
		case err == nil:
			return "", fmt.Errorf("mount point already in use")
		case os.IsNotExist(err):
			// ok — available
		default:
			return "", fmt.Errorf("requested mount point not available: %w", err)
		}
	}

	return mountPoint, nil
}

// expect a drive letter like "X:" or ""
func validateMountPoint(mountPoint string) error {
	if mountPoint == "" {
		return nil
	}
	if !validMountPointRe.MatchString(mountPoint) {
		return fmt.Errorf("mount point must be a drive letter followed by a colon (e.g. 'X:')")
	}
	return nil
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
	opts = append(opts, "-o", "DirInfoTimeout=1")
	return opts
}

func additionalIgnorePatterns() []string {
	return []string{
		// Microsoft Office lock/owner files (sidecar next to the doc)
		// ~$*
		"~$*",
		// Office scratch temp files
		// ~WR*.tmp
		// ~DF*.tmp
		// AD70B1.tmp
		// AD70B13.tmp
		// AD70B13E.tmp
		"~WR*.tmp",
		"~DF*.tmp",
		"[0-9A-F][0-9A-F][0-9A-F][0-9A-F][0-9A-F][0-9A-F].tmp",
		"[0-9A-F][0-9A-F][0-9A-F][0-9A-F][0-9A-F][0-9A-F][0-9A-F].tmp",
		"[0-9A-F][0-9A-F][0-9A-F][0-9A-F][0-9A-F][0-9A-F][0-9A-F][0-9A-F].tmp",
		// Adobe InDesign temp files
		// test4c4a9d1c-5b46.TMP
		"test[0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f]-[0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f].TMP",
	}
}
