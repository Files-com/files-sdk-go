package lib

import (
	"net/url"
	"strings"
)

var QueryEscape = url.QueryEscape

func PathEscape(path string) string {
	pathParts := strings.Split(path, "/")
	newParts := make([]string, len(pathParts))

	for _, part := range pathParts {
		newParts = append(newParts, url.PathEscape(part))
	}

	return strings.Join(newParts, "/")
}
