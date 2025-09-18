package fsmount

import (
	"fmt"
	"os"
	path_lib "path"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

const (
	officeOwnerFilePrefix = "~$"
	officeOwnerNameLength = 54 // Excel uses 54. Word uses 53, but it accepts 54
)

// nodeType represents the type of a file system node, either a file or a directory.
type nodeType int

const (
	nodeTypeFile nodeType = iota
	nodeTypeDir
)

func (nt nodeType) String() string {
	switch nt {
	case nodeTypeFile:
		return "file"
	case nodeTypeDir:
		return "directory"
	default:
		return "unknown"
	}
}

type fileHandle struct {
	// ID of the file handle, unique within the OpenHandles instance.
	id uint64

	// The node this handle is associated with
	node *fsNode

	// The number of bytes read from the file if opened for reading
	bytesRead atomic.Int64

	// readAt is the time when the file was last read
	readAt time.Time

	// The flags used when opening the file
	FuseFlags

	// If the file is stored locally, this is the *os.File represented by this handle
	// If the file is stored remotely, this will be nil.
	localFile *os.File
}

func (fh *fileHandle) String() string {
	return fmt.Sprintf("fileHandle{id: %v, node: %v, flags: %v}", fh.id, fh.node, fh.FuseFlags)
}

// isWriteOp checks if the file handle was opened as a write operation.
func (fh *fileHandle) isWriteOp() bool {
	return !fh.IsReadOnly()
}

// incrementRead increments the number of bytes read from the file
func (fh *fileHandle) incrementRead(n int64) {
	fh.bytesRead.Add(n)
	fh.readAt = time.Now()
}

type virtualfs struct {
	cacheTTL time.Duration

	// map from path to fsNode
	nodes map[string]*fsNode

	// map from handle ID to fileHandle
	handles *OpenHandles

	// mutex to protect access to nodes
	nodesMu sync.Mutex

	lib.LeveledLogger
}

func newVirtualfs(params MountParams, ll lib.LeveledLogger) *virtualfs {
	vfs := &virtualfs{
		nodes:         make(map[string]*fsNode),
		handles:       NewOpenHandles(ll),
		LeveledLogger: ll,
		cacheTTL:      DefaultCacheTTL,
	}
	if params.CacheTTL > 0 {
		vfs.cacheTTL = params.CacheTTL
	}
	return vfs
}

func (vfs *virtualfs) destroy() {
	vfs.handles.Close()
}

func (vfs *virtualfs) fetch(path string) (*fsNode, bool) {
	vfs.nodesMu.Lock()
	defer vfs.nodesMu.Unlock()

	node, ok := vfs.nodes[path]
	return node, ok
}

func (vfs *virtualfs) getOrCreate(path string, nt nodeType) (node *fsNode) {
	vfs.nodesMu.Lock()
	defer vfs.nodesMu.Unlock()

	node, ok := vfs.nodes[path]
	if !ok {
		node = &fsNode{
			path:     path,
			cacheTTL: vfs.cacheTTL,
			logger:   vfs.LeveledLogger,
		}
		node.updateInfo(fsNodeInfo{
			nodeType:     nt,
			creationTime: time.Now(),
			modTime:      time.Now(),
		})
		if nt == nodeTypeDir {
			node.childPaths = make(map[string]struct{})
		}

		vfs.add(node)
	}

	return node
}

func (vfs *virtualfs) rename(oldPath string, newPath string) *fsNode {
	node, ok := vfs.fetch(oldPath)
	if !ok {
		return nil
	}

	vfs.remove(oldPath)
	node.path = newPath

	vfs.nodesMu.Lock()
	vfs.add(node)
	vfs.nodesMu.Unlock()

	return node
}

func (vfs *virtualfs) add(node *fsNode) {
	vfs.nodes[node.path] = node

	parentPath := path_lib.Dir(node.path)
	if parentPath != node.path {
		if parent, ok := vfs.nodes[parentPath]; ok {
			parent.childPaths[node.path] = struct{}{}
		}
	}
}

func (vfs *virtualfs) remove(path string) {
	vfs.nodesMu.Lock()
	defer vfs.nodesMu.Unlock()

	delete(vfs.nodes, path)

	parentPath := path_lib.Dir(path)
	if parentPath != path {
		if parent, ok := vfs.nodes[parentPath]; ok {
			delete(parent.childPaths, path)
		}
	}
}

func (vfs *virtualfs) expireNodeInfo(path string) {
	node := vfs.getOrCreate(path, nodeTypeDir)
	node.expireInfo()
	parentPath := path_lib.Dir(path)
	parent := vfs.getOrCreate(parentPath, nodeTypeDir)
	parent.expireInfo()
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

func logPanics(log lib.LeveledLogger) {
	if r := recover(); r != nil {
		log.Error("Panic: %v\nStack trace:\n%s", r, debug.Stack())
		panic(r)
	}
}
