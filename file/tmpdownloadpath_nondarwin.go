//go:build !darwin

package file

import "context"

// explicitTmpDownload is a caller-supplied temporary file path, used as given.
func explicitTmpDownload(path string) destinationPath {
	return explicitDestination(path)
}

func tmpDownloadPathOnNotExist(_ destinationPath, tmp destinationPath) (destinationPath, error) {
	return tmp, nil
}

func finalizeTmpDownload(ctx context.Context, tmp destinationPath, final destinationPath) error {
	return tmp.moveTo(ctx, final)
}

func existingTmpDownloadFile(_ destinationPath, tmp destinationPath) (destinationPath, bool) {
	if _, err := tmp.stat(); err == nil {
		return tmp, true
	}
	return destinationPath{}, false
}

func removeTmpDownload(tmp destinationPath) error {
	return tmp.remove()
}
