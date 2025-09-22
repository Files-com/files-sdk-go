package fsmount

import (
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/winfsp/cgofuse/fuse"
)

var (
	errMountPointInUse = fmt.Errorf("mount point in use")
	errNilHost         = fmt.Errorf("nil host provided")
)

// mountRegistry manages a collection of active mount hosts.
type mountRegistry struct {
	hosts   map[string]*Host
	hostsMu sync.Mutex

	log    lib.LeveledLogger
	dbgSrv *http.Server
}

func newRegistry(logger lib.LeveledLogger) *mountRegistry {
	return &mountRegistry{
		hosts: make(map[string]*Host),
		log:   logger,
	}
}

func (reg *mountRegistry) add(mountPoint string, host *Host) (*Host, error) {
	reg.hostsMu.Lock()
	defer reg.hostsMu.Unlock()
	if _, exists := reg.hosts[mountPoint]; exists {
		return nil, errMountPointInUse
	}
	if host == nil {
		return nil, errNilHost
	}
	reg.hosts[mountPoint] = host
	return host, nil
}

func (reg *mountRegistry) get(mountPoint string) (*Host, bool) {
	reg.hostsMu.Lock()
	defer reg.hostsMu.Unlock()
	host, ok := reg.hosts[mountPoint]
	return host, ok
}

func (reg *mountRegistry) remove(mountPoint string) {
	reg.hostsMu.Lock()
	defer reg.hostsMu.Unlock()
	delete(reg.hosts, mountPoint)
}

func (reg *mountRegistry) list() []string {
	reg.hostsMu.Lock()
	defer reg.hostsMu.Unlock()
	mounts := make([]string, 0, len(reg.hosts))
	for mnt := range reg.hosts {
		mounts = append(mounts, mnt)
	}
	// sort the mount points for consistent ordering
	sort.Strings(mounts)
	return mounts
}

// Host acts as a wrapper around a fuse.FileSystemHost and an fsmount.Filescomfs to allow
// interception unmount, and notify calls. This is primarily to facilitate calling Unmount on macOS
// because the unmount action on macOS does not reliably propagate to the underlying file system
// implementations, which means they don't reliably have the opportunity to clean up resources.
type Host struct {
	fuseHost *fuse.FileSystemHost
	fs       *Filescomfs
}

// Unmount unmounts the file system and cleans up resources.
func (h *Host) Unmount() bool {
	if h.fuseHost == nil || h.fs == nil {
		return false
	}
	// unmount the fuse host first to stop any further file system operations
	unmounted := h.fuseHost.Unmount()
	if !unmounted {
		return false
	}

	// call to destroy the underlying Filescomfs instance
	if h.fs != nil {
		h.fs.Destroy()
	}

	// remove from the registry
	mntRegistry.remove(h.fs.mountPoint)
	return unmounted
}

// Notify sends a notification to the FUSE host about changes to a specific path.
// This can be used to inform the FUSE layer that a file or directory has changed,
// prompting it to refresh its cache or take other appropriate actions.
func (h *Host) Notify(path string, action uint32) bool {
	if h.fuseHost == nil {
		return false
	}
	return h.fuseHost.Notify(path, action)
}
