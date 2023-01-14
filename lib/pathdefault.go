//go:build !windows

package lib

func (p Path) NormalizePathSystemForAPI() Path {
	return p
}
