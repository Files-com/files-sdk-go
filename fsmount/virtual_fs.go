package fsmount

import (
	"fmt"
	path_lib "path"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

const (
	defaultCacheTTL       = 5 * time.Second
	officeOwnerFilePrefix = "~$"
	officeOwnerNameLength = 54 // Excel uses 54. Word uses 53, but it accepts 54, so we'll use that.
)

type virtualfs struct {
	cacheTTL     time.Duration
	nodeMap      map[string]*fsNode
	nodeMapMutex sync.Mutex
	lib.LeveledLogger
}

func newVirtualfs(logger lib.Logger, cacheTTL *time.Duration) *virtualfs {
	vfs := &virtualfs{
		nodeMap:       make(map[string]*fsNode),
		LeveledLogger: lib.NewLeveledLogger(logger),
		cacheTTL:      defaultCacheTTL,
	}
	if cacheTTL != nil {
		vfs.cacheTTL = *cacheTTL
	}
	return vfs
}

func (vfs *virtualfs) fetch(path string) (*fsNode, bool) {
	vfs.nodeMapMutex.Lock()
	defer vfs.nodeMapMutex.Unlock()

	node, ok := vfs.nodeMap[path]
	return node, ok
}

func (vfs *virtualfs) getOrCreate(path string, dir bool) (node *fsNode) {
	vfs.nodeMapMutex.Lock()
	defer vfs.nodeMapMutex.Unlock()

	node, ok := vfs.nodeMap[path]
	if !ok {
		node = &fsNode{
			fs:   vfs,
			path: path,
		}
		node.updateInfo(fsNodeInfo{dir: dir})
		if dir {
			node.childPaths = make(map[string]struct{})
		}

		vfs.add(node)
	}

	return node
}

func (vfs *virtualfs) close(path string, handle uint64) (errc int) {
	if node, ok := vfs.fetch(path); ok {
		node.closeWriterByHandle(handle)
	}

	// TODO: Remove nodes that haven't been accessed in a while.

	return
}

func (vfs *virtualfs) rename(oldPath string, newPath string) {
	node, ok := vfs.fetch(oldPath)
	if !ok {
		return
	}

	vfs.remove(oldPath)
	node.path = newPath

	vfs.nodeMapMutex.Lock()
	defer vfs.nodeMapMutex.Unlock()
	vfs.add(node)
}

func (vfs *virtualfs) add(node *fsNode) {
	vfs.nodeMap[node.path] = node

	parentPath := path_lib.Dir(node.path)
	if parentPath != node.path {
		if parent, ok := vfs.nodeMap[parentPath]; ok {
			parent.childPaths[node.path] = struct{}{}
		}
	}
}

func (vfs *virtualfs) remove(path string) {
	vfs.nodeMapMutex.Lock()
	defer vfs.nodeMapMutex.Unlock()

	delete(vfs.nodeMap, path)

	parentPath := path_lib.Dir(path)
	if parentPath != path {
		if parent, ok := vfs.nodeMap[parentPath]; ok {
			delete(parent.childPaths, path)
		}
	}
}

func (vfs *virtualfs) fetchLockTarget(path string) (*fsNode, bool) {
	if !isMsOfficeOwnerFile(path) {
		return nil, false
	}

	lockSuffix := path_lib.Base(path)[len(officeOwnerFilePrefix):]

	if parent, ok := vfs.fetch(path_lib.Dir(path)); ok {
		for childPath := range parent.childPaths {
			if strings.HasSuffix(childPath, lockSuffix) && !isMsOfficeOwnerFile(childPath) {
				return vfs.fetch(childPath)
			}
		}
	}

	return nil, false
}

func isMsOfficeOwnerFile(path string) bool {
	filename := path_lib.Base(path)
	return strings.HasPrefix(filename, officeOwnerFilePrefix)
}

func buildOwnerFile(node *fsNode) []byte {
	owner := node.info.lockOwner
	length := officeOwnerNameLength

	// Truncate the owner name if it's too long.
	if len(owner) > length {
		owner = owner[:length]
	}

	// Prefix the owner name with a byte indicating its length. Do this _after_ truncating the name.
	owner = fmt.Sprintf("%c%s", byte(len(owner)), owner)
	length++

	// Create a buffer and write the owner name in both single-byte and double-byte formats.
	ownerBuffer := make([]byte, length*3)
	for i, b := range []byte(owner) {
		ownerBuffer[i] = b
		ownerBuffer[length+(i*2)] = b
	}
	return ownerBuffer
}

func (vfs *virtualfs) logPanics() {
	if r := recover(); r != nil {
		vfs.Error("Panic: %v\nStack trace:\n%s", r, debug.Stack())
		panic(r)
	}
}
