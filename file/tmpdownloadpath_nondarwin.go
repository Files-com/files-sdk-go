//go:build !darwin

package file

import (
	"errors"
	"io/fs"
	"path/filepath"
)

// explicitTmpDownload is a caller-supplied temporary file path, used as given.
func explicitTmpDownload(path string) destinationPath {
	folder := filepath.Dir(path)
	if filepath.Base(path) == boundPausedDownloadData && isBoundPausedDownloadElement(filepath.Base(folder)) {
		return destinationPath{dir: filepath.Dir(folder), name: filepath.Join(filepath.Base(folder), filepath.Base(path))}
	}
	return explicitDestination(path)
}

func preparePausedProofLegacy(entry destinationPath, _ string) (destinationPath, error) {
	return entry, refuseOccupied(entry)
}

func tmpDownloadPathOnNotExist(_ destinationPath, tmp destinationPath) (destinationPath, error) {
	return tmp, nil
}

// removeTmpDownloadFolder is a no-op: a temporary download is a file here, and
// it has just been published.
func removeTmpDownloadFolder(_ destinationPath) error {
	return nil
}

// tmpDownloadEntry is the reserved entry that carries a temporary download's
// name: the file itself.
func tmpDownloadEntry(tmp destinationPath) destinationPath {
	if container, ok := boundPausedDownloadContainer(tmp); ok {
		return container
	}
	return tmp
}

// tmpDownloadInEntry is the temporary download inside entry: the entry itself.
func tmpDownloadInEntry(entry destinationPath, _ string) destinationPath {
	if isBoundPausedDownloadElement(entry.base()) {
		return boundPausedDataIn(entry)
	}
	return entry
}

// publishPausedTmpDownload gives the claimed, verified file of tmp its paused
// name beside tmp and removes the claim's folder.
func publishPausedTmpDownload(claimed destinationPath, tmp destinationPath, element string) (destinationPath, error) {
	paused := tmp.withName(filepath.Join(filepath.Dir(tmp.name), element))
	if err := refuseOccupied(paused); err != nil {
		return destinationPath{}, err
	}
	if err := claimed.rename(paused); err != nil {
		return destinationPath{}, err
	}
	return paused, removeTmpDownloadClaimContainer(claimed, tmp)
}

func existingTmpDownloadFile(_ destinationPath, tmp destinationPath) (destinationPath, bool) {
	if isBoundPausedDownloadElement(tmp.base()) {
		return existingBoundPausedData(tmp)
	}
	if _, err := tmp.stat(); err == nil {
		return tmp, true
	}
	return destinationPath{}, false
}

// removeTmpDownload removes a temporary download. One that is already gone,
// for example after a replaced file was refused and removed, is not an error.
func removeTmpDownload(tmp destinationPath) error {
	if _, bound := boundPausedDownloadContainer(tmp); bound {
		return removeBoundTmpDownload(tmp)
	}
	if err := tmp.remove(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}
