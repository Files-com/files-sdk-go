package file

import (
	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/folder"
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

type DownloadFolderParams struct {
	files_sdk.FolderListForParams
	Sync                bool
	ConcurrentDownloads int
}

func (d DownloadFolderParams) GetConcurrentDownloads() int {
	if d.ConcurrentDownloads == 0 {
		return 10
	}

	return d.ConcurrentDownloads
}

type Reporter func(incDownloadedBytes int64, file files_sdk.File, destination string, err error, onlyMessage string, totalFiles int)

func (c *Client) DownloadFolder(params DownloadFolderParams, rootDestination string, reporters ...Reporter) error {
	return downloadFolder(c.index(params.FolderListForParams), c, params, rootDestination, reporters...)
}

type Entity struct {
	file files_sdk.Folder
	error
}

func (c *Client) index(params files_sdk.FolderListForParams) []Entity {
	var files []Entity
	folderClient := folder.Client{Config: c.Config}
	it, err := folderClient.ListFor(params)

	if err != nil {
		files = append(files, Entity{file: files_sdk.Folder{Path: params.Path, Type: "error"}, error: err})
	}

	for it.Next() {
		entry := it.Folder()
		switch entry.Type {
		case "directory":
			files = append(files, c.index(files_sdk.FolderListForParams{Path: entry.Path})...)
		case "file":
			files = append(files, Entity{file: entry})
		default:
			files = append(files, Entity{file: entry, error: fmt.Errorf("unknown file type %v", entry.Type)})
		}
	}

	if it.Err() != nil {
		files = append(files, Entity{file: files_sdk.Folder{Path: params.Path}, error: it.Err()})
	}
	return files
}
