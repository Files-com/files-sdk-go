//go:build darwin

package file

import (
	"errors"
	"io/fs"
	"path/filepath"
)

// explicitTmpDownload is a caller-supplied temporary file path. On macOS the
// file lives inside its temporary download folder, which stays part of the
// confined name so it is removed with the file.
func explicitTmpDownload(path string) destinationPath {
	folder := filepath.Dir(path)
	return destinationPath{dir: filepath.Dir(folder), name: filepath.Join(filepath.Base(folder), filepath.Base(path))}
}

// tmpDownloadPathOnNotExist creates a temporary download folder, a common
// pattern on macOS, and returns the file, under its own name, inside it.
func tmpDownloadPathOnNotExist(final destinationPath, tmp destinationPath) (destinationPath, error) {
	if err := tmp.mkdirAll(); err != nil {
		return destinationPath{}, err
	}
	return tmp.withName(filepath.Join(tmp.name, final.base())), nil
}

// removeTmpDownloadFolder removes the temporary download folder once the file
// has been published out of it.
func removeTmpDownloadFolder(tmp destinationPath) error {
	return tmp.parent().remove()
}

func existingTmpDownloadFile(final destinationPath, tmp destinationPath) (destinationPath, bool) {
	file := tmp.withName(filepath.Join(tmp.name, final.base()))
	if _, err := file.stat(); err == nil {
		return file, true
	}
	return destinationPath{}, false
}

// removeTmpDownload removes a temporary download and its folder. The file may
// already be gone, for example after a replaced file was refused and removed,
// and the folder must still be removed then, or a later run would find it
// occupied and stage under an alternate name from then on.
func removeTmpDownload(tmp destinationPath) error {
	if err := tmp.remove(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return tmp.parent().remove()
}
