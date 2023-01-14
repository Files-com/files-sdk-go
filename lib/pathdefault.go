//go:build !windows

package lib

import "path/filepath"

func (p Path) NormalizePathSystemForAPI() Path {
	return Path{Path: filepath.Clean(p.Path)}.PruneStartingSlash()
}
