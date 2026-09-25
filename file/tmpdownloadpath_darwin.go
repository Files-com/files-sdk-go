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

// tmpDownloadEntry is the reserved entry that carries a temporary download's
// name: on macOS the folder holding the file.
func tmpDownloadEntry(tmp destinationPath) destinationPath {
	return tmp.parent()
}

// tmpDownloadInEntry is the temporary download named base inside entry.
func tmpDownloadInEntry(entry destinationPath, base string) destinationPath {
	return entry.withName(filepath.Join(entry.name, base))
}

// publishPausedTmpDownload gives the claimed, verified file of tmp its paused
// name. On macOS the claim lives inside the stage folder, so the folder is
// renamed to the paused element and the file inside it back to its own name;
// until that second step the paused folder holds nothing under the name a
// lookup uses.
func publishPausedTmpDownload(claimed destinationPath, tmp destinationPath, element string) (destinationPath, error) {
	folder := tmp.parent()
	target := folder.withName(filepath.Join(filepath.Dir(folder.name), element))
	if err := refuseOccupied(target); err != nil {
		return destinationPath{}, err
	}
	if err := folder.rename(target); err != nil {
		return destinationPath{}, err
	}
	moved := target.withName(filepath.Join(target.name, claimed.base()))
	paused := target.withName(filepath.Join(target.name, tmp.base()))
	if err := moved.rename(paused); err != nil {
		return destinationPath{}, errors.Join(err, target.rename(folder))
	}
	return paused, nil
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
