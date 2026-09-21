//go:build darwin

package file

import (
	"context"
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

func finalizeTmpDownload(ctx context.Context, tmp destinationPath, final destinationPath) error {
	if err := tmp.moveTo(ctx, final); err != nil {
		return err
	}
	return tmp.parent().remove()
}

func existingTmpDownloadFile(final destinationPath, tmp destinationPath) (destinationPath, bool) {
	file := tmp.withName(filepath.Join(tmp.name, final.base()))
	if _, err := file.stat(); err == nil {
		return file, true
	}
	return destinationPath{}, false
}

func removeTmpDownload(tmp destinationPath) error {
	if err := tmp.remove(); err != nil {
		return err
	}
	return tmp.parent().remove()
}
