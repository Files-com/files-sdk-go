//go:build darwin

package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func tmpDownloadPath(path string) string {
	return _tmpDownloadPath(path, 0)
}

func _tmpDownloadPath(path string, index int) string {
	var name string

	if index == 0 {
		name = fmt.Sprintf("%v.download", path)
	} else {
		name = fmt.Sprintf("%v (%v).download", path, index)
	}
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		os.MkdirAll(name, 0755)
		_, fileName := filepath.Split(path)
		return filepath.Join(name, fileName)
	}
	return _tmpDownloadPath(path, index+1)
}

func finalizeTmpDownload(tmpName string, finalPath string) error {
	err := os.Rename(tmpName, finalPath)
	if err != nil {
		return err
	}
	downloadPackage, _ := filepath.Split(tmpName)
	return os.Remove(downloadPackage)
}
