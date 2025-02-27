//go:build windows

package fsmount

import (
	"runtime"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/winfsp/cgofuse/fuse"
)

type MountParams struct {
	MountPoint string
	VolumeName string
	Root       string
	Config     files_sdk.Config
}

type MountHost interface {
	Unmount() bool
}

func Mount(params MountParams) MountHost {
	fs := &Filescomfs{
		root:   params.Root,
		config: params.Config,
	}
	host := fuse.NewFileSystemHost(fs)
	host.SetDirectIO(true)
	host.SetCapReaddirPlus(true)

	options := []string{"-o", "attr_timeout=1"}
	if runtime.GOOS == "windows" {
		options = append(options, "-o", "uid=-1")
		options = append(options, "-o", "gid=-1")
	}
	if params.VolumeName != "" {
		options = append(options, "-o", "volname="+params.VolumeName)
	}
	options = append(options, "-o", "FileSystemName=Files.com")

	go func() {
		host.Mount(params.MountPoint, options)
	}()
	return host
}
