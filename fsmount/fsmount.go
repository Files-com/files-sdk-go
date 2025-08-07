// Package fsmount provides functionality to mount a Files.com filesystem using FUSE.
package fsmount

import (
	"fmt"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/ignore"
	"github.com/winfsp/cgofuse/fuse"
)

const (
	// DefaultWriteConcurrency is the default number of concurrent file parts to upload.
	DefaultWriteConcurrency = 50

	// DefaultCacheTTL is the default cache TTL for the filesystem metadata.
	DefaultCacheTTL = 5 * time.Second

	// DefaultVolumeName is the default volume name for the filesystem.
	DefaultVolumeName = "Files.com"

	// DefaultDebugFuseLog is the default path to the fuse debug log. [Windows only]
	DefaultDebugFuseLog = "fuse.log"
)

// MountParams contains the parameters for mounting a Files.com filesystem using FUSE.
type MountParams struct {
	// Required. Files.com API configuration.
	Config *files_sdk.Config

	// Path to mount the filesystem.
	//
	// Optional on Windows. If provided, this is expected to be a drive letter
	// followed by a colon (e.g. "Z:"). If not specified, the letter closest to
	// the end of the Latin alphabet that is not already in use will be chosen.
	//
	// Required on MacOS and Linux, this is the path to the mount point (e.g. "/mnt/files").
	MountPoint string

	// Optional. Volume name to display in Finder/Explorer. On Windows, this is also used as the
	// share name for the UNC path. Defaults to "Files.com".
	VolumeName string

	// Optional. Files.com path to mount as the root of the filesystem. Defaults to the site root.
	Root string

	// Optional. Number of concurrent file parts to upload. Defaults to 50.
	WriteConcurrency int

	// Optional. Cache TTL for the filesystem metadata. Defaults to 5 seconds.
	CacheTTL time.Duration

	// Optional. Disable use of Files.com locks when writing files. Defaults to false.
	DisableLocking bool

	// Optional. List of patterns to ignore when creating files and directories. Defaults to
	// OS-specific defaults. To ignore no patterns, pass an empty slice.
	IgnorePatterns []string

	// Optional. If set to true, will initialize fuse configured to provide extra debug information.
	// Defaults to false.
	DebugFuse bool

	// Optional. The path to the fuse debug log. Only used if DebugFuse is set to true.
	// Defaults to fuse.log [Windows only]
	DebugFuseLog string

	// Optional. The path to the icon to display in Finder. If not specified, the default icon
	// for a network drive is used. [MacOS only]
	IconPath string
}

// MountHost defines the interface for a mounted Files.com filesystem.
type MountHost interface {
	Unmount() bool
}

// Mount initializes a Files.com filesystem and mounts it using FUSE.
func Mount(params MountParams) (MountHost, error) {
	fs, err := newFs(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create filesystem: %w", err)
	}

	// test that the fs can list the root
	if err := fs.Validate(); err != nil {
		return nil, err
	}

	// get the platform specific mount options
	opts := mountOpts(params)

	// test that the fuse library can be loaded
	if err := loadFuse(); err != nil {
		return nil, err
	}

	// Create the filesystem host and mount it
	host := fuse.NewFileSystemHost(fs)
	host.SetCapReaddirPlus(true)
	go func() {
		host.Mount(fs.mountPoint, opts)
	}()

	return host, nil
}

func newFs(params MountParams) (*Filescomfs, error) {
	// return early if config is nil or the mount point can't
	// be determined
	if params.Config == nil {
		return nil, fmt.Errorf("config is required")
	}
	// get the platform specific mount point
	mountPoint, err := mountPoint(params.MountPoint)
	if err != nil {
		return nil, err
	}
	params.MountPoint = mountPoint

	fs := &Filescomfs{
		mountPoint:       params.MountPoint,
		root:             params.Root,
		writeConcurrency: params.WriteConcurrency,
		cacheTTL:         params.CacheTTL,
		config:           params.Config,
		disableLocking:   params.DisableLocking,
		debugFuse:        params.DebugFuse,
	}

	// ensure write concurrency and cache TTL are positive
	if fs.writeConcurrency <= 0 {
		fs.writeConcurrency = DefaultWriteConcurrency
	}
	if fs.cacheTTL <= 0 {
		fs.cacheTTL = DefaultCacheTTL
	}
	if params.IgnorePatterns == nil || len(params.IgnorePatterns) > 0 {
		fs.ignore, _ = ignore.New(params.IgnorePatterns...)
	}
	return fs, nil
}

// Default mount options for all fuse implementations
func defaultMountOpts(params MountParams) []string {
	opts := []string{}
	opts = append(opts, "-o", "attr_timeout=1")
	if params.DebugFuse {
		// enables debug logging in the underlying fuse implementation
		opts = append(opts, "-o", "debug")
	}
	volname := DefaultVolumeName
	if params.VolumeName != "" {
		volname = params.VolumeName
	}
	// sets the volume name that is passed to the fuse implementation
	opts = append(opts, "-o", "volname="+volname)
	return opts
}

// Test if the fuse library can be loaded, and gracefully handle any error.
func loadFuse() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// This will panic if the fuse library cannot be loaded.
	_, _ = fuse.OptParse([]string{}, "")
	return
}
