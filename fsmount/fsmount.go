// Package fsmount provides functionality to mount a Files.com file system using FUSE.
package fsmount

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	api_key "github.com/Files-com/files-sdk-go/v3/apikey"
	dc "github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/disk"
	mc "github.com/Files-com/files-sdk-go/v3/fsmount/internal/cache/mem"

	"github.com/Files-com/files-sdk-go/v3/ignore"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/winfsp/cgofuse/fuse"

	"github.com/Files-com/files-sdk-go/v3/fsmount/internal/log"
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
	// Required on MacOS and Linux, this is the path to the directory where mount points
	// will be located (e.g. "/mnt/files").
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

	// DiskCacheEnabled determines whether to use a disk-based cache for file data.
	// If false, an in-memory cache will be used instead.
	DiskCacheEnabled bool

	// DiskCachePath specifies the file system directory where disk cache data is stored.
	// The directory must be writable and have sufficient space for the cache.
	//
	// Ignored if DiskCacheEnabled is set to false
	DiskCachePath string

	// DiskCacheParams contains the configuration parameters for the disk cache.
	//
	// Ignored if DiskCacheEnabled is set to false
	DiskCacheParams CacheParams

	// MemoryCacheParams contains the configuration parameters for the in-memory cache.
	//
	// Ignored if DiskCacheEnabled is set to true
	MemoryCacheParams CacheParams
}

// CacheParams defines the configuration parameters for the file system cache.
// It controls cache behavior including capacity limits, timing settings, and storage location.
type CacheParams struct {
	// CapacityBytes specifies the maximum size of the cache in bytes.
	// When this limit is reached, older entries will be evicted to make room for new ones.
	CapacityBytes int64

	// MaintenanceInterval defines how frequently cache maintenance tasks run,
	// such as cleaning up expired entries and enforcing capacity limits.
	MaintenanceInterval time.Duration

	// MaxAge specifies the maximum duration a cache entry can exist
	// before it is considered stale and eligible for eviction.
	MaxAge time.Duration

	// MaxFileCount limits the total number of files that can be stored in the cache.
	// When this limit is reached, older files will be removed to accommodate new ones.
	MaxFileCount int64
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

	// make sure the API key is valid before attempting to mount
	apiKeyClient := &api_key.Client{Config: *params.Config}
	_, err := apiKeyClient.FindCurrent()
	if err != nil {
		return nil, err
	}

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
	// the mountPoint function is platform specific
	mountPoint, err := mountPoint(params.MountPoint)
	if err != nil {
		return nil, err
	}
	params.MountPoint = mountPoint

	if params.Root == "" {
		params.Root = "/"
	}

	tmpRoot, err := tmpFsPath(params.TmpFsPath)
	if err != nil {
		return nil, err
	}
	params.TmpFsPath = tmpRoot

	// newCache creates either a disk or memory cache based on if the mount parameters
	// have disk caching enabled or disabled.
	cache, err := newCache(params, logger)
	if err != nil {
		return nil, err
	}

	// The Filescomfs, RemoteFs, and LocalFs all share a single virtualfs instance to manage the
	// in-memory representation of the file system. This allows for consistent state management and
	// caching across the different file system implementations.
	vfs := newVirtualfs(params, logger)

	remotefs, err := newRemoteFs(params, vfs, logger, cache)
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

func newCache(params MountParams, log log.Logger) (cacheStore, error) {
	if params.DiskCacheEnabled {
		cachePath, err := diskCachePath(params)
		if err != nil {
			return nil, err
		}
		opts := []dc.Option{
			dc.WithLogger(log),
		}
		if params.DiskCacheParams.CapacityBytes > 0 {
			opts = append(opts, dc.WithCapacityBytes(params.DiskCacheParams.CapacityBytes))
		}
		if params.DiskCacheParams.MaintenanceInterval > 0 {
			opts = append(opts, dc.WithMaintenanceInterval(params.DiskCacheParams.MaintenanceInterval))
		}
		if params.DiskCacheParams.MaxAge > 0 {
			opts = append(opts, dc.WithMaxAge(params.DiskCacheParams.MaxAge))
		}
		if params.DiskCacheParams.MaxFileCount > 0 {
			opts = append(opts, dc.WithMaxFileCount(params.DiskCacheParams.MaxFileCount))
		}
		cache, err := dc.NewDiskCache(cachePath, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create disk cache: %w", err)
		}
		return cache, nil
	}

	// if the disk cache has been disabled, use an in memory cache
	opts := []mc.Option{
		mc.WithLogger(log),
	}
	if params.MemoryCacheParams.CapacityBytes > 0 {
		opts = append(opts, mc.WithCapacityBytes(params.MemoryCacheParams.CapacityBytes))
	}
	if params.MemoryCacheParams.MaxAge > 0 {
		opts = append(opts, mc.WithMaxAge(params.MemoryCacheParams.MaxAge))
	}
	if params.MemoryCacheParams.MaxFileCount > 0 {
		opts = append(opts, mc.WithMaxFileCount(params.MemoryCacheParams.MaxFileCount))
	}
	if params.MemoryCacheParams.MaintenanceInterval > 0 {
		opts = append(opts, mc.WithMaintenanceInterval(params.MemoryCacheParams.MaintenanceInterval))
	}
	cache, err := mc.NewMemoryCache(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory cache: %w", err)
	}
	return cache, nil
}

// tmpFsPath returns the path to the temporary file system directory.
// If the provided path is empty, it defaults to an application specific path appended to
// the OS-specific temporary directory.
// If the provided path does not exist, it is created.
// If the provided path exists but is not a directory, an error is returned.
func tmpFsPath(path string) (string, error) {
	if path != "" {
		st, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				if err := os.MkdirAll(path, 0o700); err != nil {
					return "", fmt.Errorf("failed to create temporary file system path: %w", err)
				}
				return path, nil
			}
			return "", fmt.Errorf("failed to access temporary file system path: %w", err)
		}
		if !st.IsDir() {
			return "", fmt.Errorf("temporary file system path is not a directory")
		}
		return path, nil
	}
	path, err := os.MkdirTemp("", "Files.com-v6-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary local file system: %w", err)
	}
	return path, nil
}

// diskCachePath returns the path to the directory for the disk cache.
// If the provided path is empty, it defaults to an application specific path appended to
// the OS-specific cache directory.
// If the provided path does not exist, it is created.
// If the provided path exists but is not a directory, an error is returned.
func diskCachePath(params MountParams) (string, error) {
	path := params.DiskCachePath
	if path != "" {
		st, err := os.Stat(path)
		if err != nil {
			return "", fmt.Errorf("failed to stat local cache path: %w", err)
		}
		if !st.IsDir() {
			return "", fmt.Errorf("local cache path is not a directory")
		}
		return path, nil
	}

	// the passed in path is empty, so use the OS-specific cache directory
	// e.g. /Users/username/Library/Caches on MacOS
	//      /home/username/.cache on Linux
	//      C:\Users\username\AppData\Local on Windows
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to locate user cache directory: %w", err)
	}

	// use a cache directory specific to this mount point
	// to avoid conflicts with multiple mounts
	// e.g. /Users/username/Library/Caches/Files.com/v6/A/cache
	var mountBase string
	switch runtime.GOOS {
	case "windows":
		// On Windows, the mount point is a drive letter like "X:",
		// so use the drive letter as the directory name
		mountBase = params.MountPoint[:1]
	default:
		// On MacOS and Linux, the mount point is a path like "/mnt/files/A",
		// so use the base name of the path as the directory name
		mountBase = filepath.Base(params.MountPoint)
	}

	if mountBase == string(os.PathSeparator) || mountBase == "" {
		return "", fmt.Errorf("failed to locate cache directory: expected path or drive letter: got '%s'", mountBase)
	}
	path = filepath.Join(cacheDir, "Files.com", "v6", mountBase, "cache")
	if err := os.MkdirAll(path, 0o700); err != nil {
		return "", fmt.Errorf("failed to create mount specific cache directory: %w", err)
	}
	return path, nil
}

func ignoreFromPatterns(patterns []string) (*gogitignore.GitIgnore, error) {
	switch patterns {
	case nil:
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
