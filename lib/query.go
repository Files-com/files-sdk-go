package lib

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var QueryEscape = url.QueryEscape

func PathEscape(path string) string {
	pathParts := strings.Split(path, "/")
	newParts := make([]string, len(pathParts))

	for i, part := range pathParts {
		newParts[i] = url.PathEscape(part)
	}

	return strings.Join(newParts, "/")
}

func BuildPath(resourcePath string, unescapedPath string) string {
	viaOS := filepath.Join(resourcePath, PathEscape(unescapedPath))
	return strings.Join(strings.Split(viaOS, string(os.PathSeparator)), "/")
}
