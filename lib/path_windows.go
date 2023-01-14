package lib

import (
	"os"
	"strings"
)

func (p Path) NormalizePathSystemForAPI() Path {
	parts := strings.Split(p.Path, string(os.PathSeparator))
	return Path{Path: strings.Join(parts, "/")}
}
