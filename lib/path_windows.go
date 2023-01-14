package lib

import (
	"os"
	"path/filepath"
	"strings"
)

func (p Path) NormalizePathSystemForAPI() Path {
	parts := strings.Split(filepath.Clean(p.Path), string(os.PathSeparator))
	return Path{Path: strings.Join(parts, "/")}.PruneStartingSlash()
}
