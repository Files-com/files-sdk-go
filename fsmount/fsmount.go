//go:build windows

package fsmount

import (
	"fmt"
	"os"
	"runtime"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/ignore"
	"github.com/winfsp/cgofuse/fuse"
)

type MountParams struct {
	// Optional. Path to mount the filesystem. On Windows, this is expected to be a drive letter
	// followed by a colon (e.g. "Z:"). If not specified, the highest available drive letter will
	// be used.
	MountPoint string

	// Optional. Volume name to display in Finder/Explorer. On Windows, this is also used as the
	// share name for the UNC path. The server name will be "Files.com".
	VolumeName string

	// Optional. Files.com path to mount as the root of the filesystem. Defaults to the site root.
	Root string

	// Optional. Number of concurrent file parts to upload. Defaults to 50.
	WriteConcurrency *int

	// Optional. Cache TTL for the filesystem metadata. Defaults to 1 second.
	CacheTTL *time.Duration

	// Optional. Disable use of Files.com locks when writing files. Defaults to false.
	DisableLocking bool

	// Optional. List of patterns to ignore when creating files and directories. Defaults to
	// OS-specific defaults. To ignore no patterns, pass an empty slice.
	IgnorePatterns []string

	Config files_sdk.Config

	// Optional. If set to true, will initialize fuse configured to provide extra debug information.
	// Defaults to false.
	DebugFuse bool
	// Optional. The path to the fuse debug log. Only used if DebugFuse is set to true.
	// Defaults to fuse.log
	DebugFuseLog string
}

type MountHost interface {
	Unmount() bool
}

func Mount(params MountParams) (MountHost, error) {
	mountPoint, err := getMountPoint(params.MountPoint)
	if err != nil {
		return nil, err
	}

	fs := &Filescomfs{
		mountPoint:       mountPoint,
		root:             params.Root,
		writeConcurrency: params.WriteConcurrency,
		cacheTTL:         params.CacheTTL,
		config:           params.Config,
		disableLocking:   params.DisableLocking,
		debugFuse:        params.DebugFuse,
	}

	if params.IgnorePatterns == nil || len(params.IgnorePatterns) > 0 {
		fs.ignore, err = ignore.New(params.IgnorePatterns...)
		if err != nil {
			return nil, err
		}
	}

	if err := fs.Validate(); err != nil {
		return nil, err
	}

	host := fuse.NewFileSystemHost(fs)
	host.SetCapReaddirPlus(true)

	options := []string{"-o", "attr_timeout=1"}

	if fs.debugFuse {
		logfile := "fuse.log"
		if params.DebugFuseLog != "" {
			logfile = params.DebugFuseLog
		}
		options = append(options, "-o", "debug")
		options = append(options, "-o", "DebugLog="+logfile)
	}

	if runtime.GOOS == "windows" {
		options = append(options, "-o", "uid=-1")
		options = append(options, "-o", "gid=-1")
		if params.VolumeName != "" {
			options = append(options, "--VolumePrefix=\\Files.com\\"+params.VolumeName)
		}
	} else {
		if params.VolumeName != "" {
			options = append(options, "-o", "volname="+params.VolumeName)
		}
	}
	options = append(options, "-o", "FileSystemName=Files.com")

	if err := initFuse(); err != nil {
		return nil, err
	}

	go func() {
		host.Mount(mountPoint, options)
	}()

	return host, nil
}

func getMountPoint(mountPoint string) (string, error) {
	if runtime.GOOS == "windows" {
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
	}

	return mountPoint, nil
}

// Test if the fuse library can be loaded, and gracefully handle any error.
func initFuse() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// This will panic if the fuse library cannot be loaded.
	fuse.OptParse([]string{}, "")
	return
}
