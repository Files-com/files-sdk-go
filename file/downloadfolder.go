package file

import (
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/lib"
	"github.com/zenthangplus/goccm"
)

type Downloader interface {
	Download(files_sdk.FileDownloadParams) (files_sdk.File, error)
}

func downloadFolder(files []Entity, c Downloader, params DownloadFolderParams, rootDestination string, reporters ...Reporter) error {
	rootDestinationIsDir := false
	if rootDestination != "" && rootDestination[len(rootDestination)-1:] == string(os.PathSeparator) {
		rootDestinationIsDir = true
	} else {
		rootDestination, _ = filepath.Abs(rootDestination)
		fi, err := os.Stat(rootDestination)
		if err == nil && fi.Mode().IsDir() {
			rootDestinationIsDir = true
		}
	}
	goc := goccm.New(params.GetConcurrentDownloads())

	downloadCount := len(files)
	if len(files) > 1 {
		rootDestinationIsDir = true
	}
	signal := make(chan bool)

	sourceRootLen := len(strings.Split(params.Path, "/"))
	for _, entity := range files {
		if entity.error != nil {
			return entity.error
		}
		goc.Wait()
		go func(entity Entity) {
			defer func() {
				signal <- true
				goc.Done()
			}()
			file := files_sdk.File{Path: entity.file.Path, Size: entity.file.Size, Type: entity.file.Type}
			sep := strings.Split(file.Path, "/")
			r := int(math.Min(float64(len(sep)-1), float64(sourceRootLen)))
			filePathCompacted := strings.Join(sep[r:], string(os.PathSeparator))
			filePath, fileName := filepath.Split(filePathCompacted)
			var destinationPath string
			if rootDestinationIsDir {
				destinationPath = filepath.Join(rootDestination, filePath, fileName)
			} else {
				destinationPath = filepath.Join(rootDestination, filePath)
			}
			dir, _ := filepath.Split(destinationPath)
			_, err := os.Stat(dir)
			if os.IsNotExist(err) {
				os.MkdirAll(dir, 0755)
			}
			var out *os.File
			fileInfo, err := os.Stat(destinationPath)
			if err != nil && !os.IsNotExist(err) {
				if len(reporters) > 0 {
					reporters[0](0, file, destinationPath, err, "", downloadCount)
				}
			}

			if !os.IsNotExist(err) && params.Sync && !entity.file.Mtime.After(fileInfo.ModTime()) {
				//	Local version is the same or newer
				downloadCount -= 1
				return
			}
			out, err = os.Create(destinationPath + ".part")
			defer func() {
				err = os.Rename(destinationPath+".part", destinationPath)
				if err != nil {
					if len(reporters) > 0 {
						reporters[0](0, file, destinationPath, err, "", downloadCount)
					}
				}
			}()

			if err != nil {
				if len(reporters) > 0 {
					reporters[0](0, file, destinationPath, err, "", downloadCount)
				}
			}
			downloadParams := files_sdk.FileDownloadParams{Path: file.Path}
			writer := lib.ProgressWriter{Writer: out}
			writer.ProgressWatcher = func(incDownloadedBytes int64) {
				if len(reporters) > 0 {
					reporters[0](incDownloadedBytes, file, destinationPath, entity.error, "", downloadCount)
				}
			}
			downloadParams.Writer = writer
			downloadParams.OnDownload = func(response *http.Response) {
				file.Size = response.ContentLength
				if len(reporters) > 0 {
					reporters[0](0, file, destinationPath, entity.error, "", downloadCount)
				}
			}
			writer.ProgressWatcher(0)
			newFile, err := c.Download(downloadParams)
			if len(reporters) > 0 && err != nil {
				reporters[0](0, newFile, destinationPath, err, "", downloadCount)
			}
		}(entity)

	}
	for range files {
		<-signal
	}
	if downloadCount == 0 {
		if len(reporters) > 0 {
			reporters[0](0, files_sdk.File{}, params.Path, nil, "No files to download", 0)
		}
	}
	return nil
}
