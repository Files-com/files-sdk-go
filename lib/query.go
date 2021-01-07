package lib

import (
	"net/url"
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
	return filepath.Join(resourcePath, PathEscape(unescapedPath))
}
