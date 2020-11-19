package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Files-com/files-sdk-go/folder"
	"github.com/zenthangplus/goccm"

	files_sdk "github.com/Files-com/files-sdk-go"
)

func (c *Client) DownloadToFile(params files_sdk.FileDownloadParams, filePath string) (files_sdk.File, error) {
	out, err := os.Create(filePath)
	if err != nil {
		return files_sdk.File{}, err
	}
	params.Writer = out
	return c.Download(params)
}

func DownloadToFile(params files_sdk.FileDownloadParams, filePath string) (files_sdk.File, error) {
	return (&Client{}).DownloadToFile(params, filePath)
}

type fileDownload struct {
	files_sdk.File
	destination string
	error
}

func (c *Client) DownloadFolder(params files_sdk.FolderListForParams, rootDestination string, reporters ...func(file files_sdk.File, destination string, err error)) error {
	queuedDownloads := 0
	goc := goccm.New(c.Config.MaxConcurrentConnections())
	fileChannel := make(chan fileDownload)

	err := c.downloadNode(params, rootDestination, &queuedDownloads, goc, &fileChannel)
	if err != nil {
		return err
	}

	for i := 0; i < queuedDownloads; i++ {
		download := <-fileChannel
		if len(reporters) > 0 {
			reporters[0](download.File, download.destination, download.error)
		}
	}
	return nil
}

func (c *Client) downloadNode(params files_sdk.FolderListForParams, rootDestination string, queuedDownloads *int, goc goccm.ConcurrencyManager, fileChannel *chan fileDownload) error {
	folderClient := folder.Client{Config: c.Config}
	it, err := folderClient.ListFor(params)

	if err != nil {
		return err
	}

	action := func(entry files_sdk.Folder) {
		goc.Wait()
		destinationPath := filepath.Join(rootDestination, entry.Path)
		download := fileDownload{destination: destinationPath, File: files_sdk.File{Path: entry.Path, Type: entry.Type}}
		switch entry.Type {
		case "directory":
			_, err := os.Stat(destinationPath)
			if os.IsNotExist(err) {
				os.MkdirAll(destinationPath, 0755)
			}
			err = c.downloadNode(files_sdk.FolderListForParams{Path: entry.Path}, rootDestination, queuedDownloads, goc, fileChannel)
			if err != nil {
				download.error = err
			}
		case "file":
			dir, _ := filepath.Split(destinationPath)
			_, err := os.Stat(dir)
			if os.IsNotExist(err) {
				os.MkdirAll(dir, 0755)
			}
			file, err := c.DownloadToFile(files_sdk.FileDownloadParams{Path: entry.Path}, destinationPath)
			download.error = err
			download.File = file
		default:
			if it.Err() != nil {
				download.error = it.Err()
			} else {
				download.error = fmt.Errorf("unknown file type %v", entry.Type)
			}
		}
		*fileChannel <- download
		goc.Done()
	}
	for it.Next() {
		*queuedDownloads++
		entry := it.Folder()
		go action(entry)
	}
	return nil
}
