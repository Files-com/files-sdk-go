package lib

import (
	"os"
	"path/filepath"
	"strings"
)

type Path struct {
	Path string
}

func (p Path) Pop() string {
	_, last := filepath.Split(strings.TrimSuffix(p.Path, string(os.PathSeparator)))
	return last
}

func (p Path) EndingSlash() bool {
	if p.Path == "" {
		return false
	}
	return p.Path[len(p.Path)-1:] == string(os.PathSeparator)
}

func (p Path) PruneStartingSlash() Path {
	if p.Path == "" {
		return p
	}

	if p.Path[0:1] == string(os.PathSeparator) {
		return Path{Path: p.Path[1:]}
	}
	return p
}

func (p Path) PruneEndingSlash() Path {
	if !p.EndingSlash() {
		return Path{Path: p.Path}
	}

	return Path{Path: p.Path[0 : len(p.Path)-1]}
}

func (p Path) ConvertEmptyToRoot() Path {
	if p.Path == "" {
		return Path{Path: "."}
	}

	return p
}

func (p Path) Join(str ...string) Path {
	return Path{Path: filepath.Join(p.String(), filepath.Join(str...))}
}

func (p Path) String() string {
	return p.Path
}

func (p Path) NormalizePathSystemForAPI() Path {
	return Path{
		Path: UrlJoinNoEscape(strings.Split(filepath.Clean(p.Path), string(os.PathSeparator))...),
	}.PruneStartingSlash()
}
