//go:build windows

package fsmount

import (
	"fmt"
	"runtime"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/winfsp/cgofuse/fuse"
)

type MountParams struct {
	MountPoint       string
	VolumeName       string
	Root             string
	WriteConcurrency *int
	CacheTTL         *time.Duration
	Config           files_sdk.Config
}

type MountHost interface {
	Unmount() bool
}

func Mount(params MountParams) (MountHost, error) {
	fs := &Filescomfs{
		mountPoint:       params.MountPoint,
		root:             params.Root,
		writeConcurrency: params.WriteConcurrency,
		cacheTTL:         params.CacheTTL,
		config:           params.Config,
	}
	host := fuse.NewFileSystemHost(fs)
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

	if err := initFuse(); err != nil {
		return nil, err
	}

	go func() {
		host.Mount(params.MountPoint, options)
	}()

	return host, nil
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
