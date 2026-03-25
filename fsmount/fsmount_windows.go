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

// findAvailableDriveLetter searches for an available drive letter from Z: to D:
func findAvailableDriveLetter() (string, error) {
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
}

// mountPoint validates or selects a mount point (drive letter) for the FUSE file system.
// If mountPoint is empty, it finds the highest available drive letter by searching backward
// from Z: to D:.
// If mountPoint is provided, it validates that it is a single drive letter followed by a colon (e.g. "X:")
// and checks that the drive letter is not already in use.
// Returns the selected or validated mount point, or an error if no valid mount point is found.
func mountPoint(mountPoint string, useDefaultMountPoint bool) (string, error) {
	if err := validateMountPoint(mountPoint); err != nil {
		return "", err
	}

	if mountPoint == "" {
		// Find the highest available drive letter.
		return findAvailableDriveLetter()
	} else {
		_, err := os.Stat(mountPoint + string(os.PathSeparator))
		switch {
		case err == nil:
			if useDefaultMountPoint {
				// Mount point in use with useDefault=true: fall back to Z-D search
				return findAvailableDriveLetter()
			}
			return "", fmt.Errorf("mount point already in use")

		case os.IsNotExist(err):
			// ok — available
		default:
			if useDefaultMountPoint {
				// Mount point not available with useDefault=true: fall back to Z-D search
				return findAvailableDriveLetter()
			}
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

	// TODO: Decide if these options can be used. Certain applications like InDesign expect and actually
	// validate that the uid returned in the *fuse.Stat_t during Getattr matches what the program itself
	// sees when probing temp files on the file system.
	// opts = append(opts, "-o", "uid=-1")
	// opts = append(opts, "-o", "gid=-1")
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
	}
}
