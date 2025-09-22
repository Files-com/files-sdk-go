// Package fsmount provides functionality to mount a Files.com file system using FUSE.
package fsmount

import (
	"fmt"
	"os"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/ignore"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/winfsp/cgofuse/fuse"

	gogitignore "github.com/sabhiram/go-gitignore"
)

const (
	// DefaultWriteConcurrency is the default number of concurrent file parts to upload.
	DefaultWriteConcurrency = 50

	// DefaultCacheTTL is the default cache TTL for the file system metadata.
	DefaultCacheTTL = 5 * time.Second

	// DefaultVolumeName is the default volume name for the file system.
	DefaultVolumeName = "Files.com"

	// DefaultDebugFuseLog is the default path to the fuse debug log. [Windows only]
	DefaultDebugFuseLog = "fuse.log"
)

// MountParams contains the parameters for mounting a Files.com file system using FUSE.
type MountParams struct {
	// Required. Files.com API configuration.
	Config *files_sdk.Config

	// Path to mount the file system.
	//
	// Optional on Windows. If provided, this is expected to be a drive letter
	// followed by a colon (e.g. "Z:"). If not specified, the letter closest to
	// the end of the Latin alphabet that is not already in use will be chosen.
	//
	// Required on MacOS and Linux, this is the path to the mount point (e.g. "/mnt/files").
	MountPoint string

	// Optional. Path to a temporary directory for storing files that don't belong on Files.com.
	// e.g. .DS_Store, Thumbs.db, etc... The full list of patterns is available in the ignore package
	// https://github.com/Files-com/files-sdk-go/tree/master/ignore/data
	//
	// Defaults to OS-specific temporary directory if not specified.
	TmpFsPath string

	// Optional. Volume name to display in Finder/Explorer. On Windows, this is also used as the
	// share name for the UNC path. Defaults to "Files.com".
	VolumeName string

	// Optional. Files.com path to mount as the root of the file system. Defaults to the site root.
	Root string

	// Optional. Number of concurrent file parts to upload. Defaults to 50.
	WriteConcurrency int

	// Optional. Cache TTL for the file system metadata. Defaults to 5 seconds.
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

var (
	// the registry of active mount hosts
	mntRegistry *mountRegistry

	// mutex to ensure only one mount operation is happening at a time
	mountMu sync.Mutex
)

// Mount initializes a Files.com file system and mounts it using FUSE.
func Mount(params MountParams) (*Host, error) {
	// only allow one mount operation at a time, to avoid multiple mounts
	// attempting to use the same mount point
	mountMu.Lock()
	defer mountMu.Unlock()

	logger := lib.NewLeveledLogger(params.Config.Logger)

	if mntRegistry == nil {
		mntRegistry = newRegistry(logger)

		// if the binary is built with the 'filescomfs_debug' tag
		//   start the debug server to allow for pprof profiling and other debug endpoints
		//   the listen address and port can be configured using the 'FILESCOMFS_DEBUG_PPROF_HOST'
		//   and 'FILESCOMFS_DEBUG_PPROF_PORT' environment variables, respectively. If not set, the
		//   default is 'localhost:6060'. If the debug server is already running, this is a no-op.
		// if the binary is not built with the 'filescomfs_debug' tag
		//   this is a no-op and the debug server will not be started
		mntRegistry.startPprof()
	}

	fs, err := newFs(params, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create file system: %w", err)
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

	// Create the file system host and mount it
	host := fuse.NewFileSystemHost(fs)
	host.SetCapReaddirPlus(true)
	go func() {
		mounted := host.Mount(fs.mountPoint, opts)
		if !mounted {
			fs.log.Error("Failed to mount file system at %s", fs.mountPoint)
			mntRegistry.remove(fs.mountPoint)
			return
		}
	}()

	return mntRegistry.add(fs.mountPoint, &Host{
		fuseHost: host,
		fs:       fs,
	})
}

func newFs(params MountParams, logger lib.LeveledLogger) (*Filescomfs, error) {
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

	if params.Root == "" {
		params.Root = "/"
	}

	// create the temporary directory if it doesn't exist
	if params.TmpFsPath == "" {
		tmpRoot, err := os.MkdirTemp("", "Files.com-v6-*")
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary local file system: %w", err)
		}
		params.TmpFsPath = tmpRoot
	}

	// The Filescomfs, RemoteFs, and LocalFs all share a single virtualfs instance to manage the
	// in-memory representation of the file system. This allows for consistent state management and
	// caching across the different file system implementations.
	vfs := newVirtualfs(params, logger)

	remotefs, err := newRemoteFs(params, vfs, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create remote file system: %w", err)
	}

	localfs := newLocalFs(params, vfs, logger)

	// create the Filescomfs instance
	fs := &Filescomfs{
		mountPoint:  mountPoint,
		remoteRoot:  params.Root,
		localFsRoot: params.TmpFsPath,
		log:         logger,
		remote:      remotefs,
		local:       localfs,
		vfs:         vfs,
	}

	ig, err := ignoreFromPatterns(params.IgnorePatterns)
	if err != nil {
		return nil, err
	}
	fs.ignore = ig

	logger.Info("Mounting Files.com file system at %s", params.MountPoint)
	return fs, nil
}

func ignoreFromPatterns(patterns []string) (*gogitignore.GitIgnore, error) {
	switch {
	case patterns == nil:
		// use OS-specific defaults + additional common patterns
		return ignore.NewWithDenyList(additionalIgnorePatterns()...)
	default:
		// use provided override patterns
		return ignore.New(patterns...)
	}
}

// Default mount options for all fuse implementations
func defaultMountOpts(params MountParams) []string {
	opts := []string{}
	opts = append(opts, "-o", "attr_timeout=1")
	opts = append(opts, "-o", "hard_remove") // avoids .fuse_hiddenXXXXXX files on delete
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
