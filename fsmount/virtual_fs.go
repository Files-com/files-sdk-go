package fsmount

import (
	filepath "path"
	"sync"
	"time"
)

const (
	defaultCacheTTL = 1 * time.Second
)

type virtualfs struct {
	cacheTTL     time.Duration
	nodeMap      map[string]*fsNode
	nodeMapMutex sync.Mutex
}

func newVirtualfs(cacheTTL *time.Duration) *virtualfs {
	vfs := &virtualfs{
		nodeMap:  make(map[string]*fsNode),
		cacheTTL: defaultCacheTTL,
	}
	if cacheTTL != nil {
		vfs.cacheTTL = *cacheTTL
	}
	return vfs
}

func (self *virtualfs) fetch(path string) (*fsNode, bool) {
	self.nodeMapMutex.Lock()
	defer self.nodeMapMutex.Unlock()

	node, ok := self.nodeMap[path]
	return node, ok
}

func (self *virtualfs) getOrCreate(path string, dir bool) (node *fsNode) {
	self.nodeMapMutex.Lock()
	defer self.nodeMapMutex.Unlock()

	node, ok := self.nodeMap[path]
	if !ok {
		node = &fsNode{
			fs:   self,
			path: path,
		}
		node.updateInfo(fsNodeInfo{dir: dir})
		if dir {
			node.childPaths = make(map[string]struct{})
		}

		self.add(node)
	}

	return node
}

func (self *virtualfs) close(path string) (errc int) {
	if node, ok := self.fetch(path); ok {
		node.closeWriter()
	}

	// TODO: Remove nodes that haven't been accessed in a while.

	return
}

func (self *virtualfs) rename(oldPath string, newPath string) {
	node, ok := self.fetch(oldPath)
	if !ok {
		return
	}

	self.delete(oldPath)
	node.path = newPath

	self.nodeMapMutex.Lock()
	defer self.nodeMapMutex.Unlock()
	self.add(node)
}

func (self *virtualfs) add(node *fsNode) {
	self.nodeMap[node.path] = node

	parentPath := filepath.Dir(node.path)
	if parentPath != node.path {
		if parent, ok := self.nodeMap[parentPath]; ok {
			parent.childPaths[node.path] = struct{}{}
		}
	}
}

func (self *virtualfs) delete(path string) {
	self.nodeMapMutex.Lock()
	defer self.nodeMapMutex.Unlock()

	delete(self.nodeMap, path)

	parentPath := filepath.Dir(path)
	if parentPath != path {
		if parent, ok := self.nodeMap[parentPath]; ok {
			delete(parent.childPaths, path)
		}
	}
}
