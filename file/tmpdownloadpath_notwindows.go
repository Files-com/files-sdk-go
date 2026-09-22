//go:build !windows

package file

// canonicalLocalPath returns path unchanged. Only NTFS keeps a second, 8.3 name
// for a file, so on every other filesystem the name a path is written with is
// the name the file has. Names that differ only in case or in how a character
// is spelled are handled by the folded comparison in isReservedTempDownloadElement.
func canonicalLocalPath(path string) (string, error) {
	return path, nil
}
