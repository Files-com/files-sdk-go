//go:build !darwin

package file

import (
	"errors"
	"io/fs"
)

// explicitTmpDownload is a caller-supplied temporary file path, used as given.
func explicitTmpDownload(path string) destinationPath {
	return explicitDestination(path)
}

func tmpDownloadPathOnNotExist(_ destinationPath, tmp destinationPath) (destinationPath, error) {
	return tmp, nil
}

// removeTmpDownloadFolder is a no-op: a temporary download is a file here, and
// it has just been published.
func removeTmpDownloadFolder(_ destinationPath) error {
	return nil
}

func existingTmpDownloadFile(_ destinationPath, tmp destinationPath) (destinationPath, bool) {
	if _, err := tmp.stat(); err == nil {
		return tmp, true
	}
	return destinationPath{}, false
}

// removeTmpDownload removes a temporary download. One that is already gone,
// for example after a replaced file was refused and removed, is not an error.
func removeTmpDownload(tmp destinationPath) error {
	if err := tmp.remove(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}
