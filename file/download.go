package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Files-com/files-sdk-go/file/manager"

	"github.com/Files-com/files-sdk-go/file/status"

	"encoding/json"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/folder"
)

type DownloadRetryParams struct {
	status.File
	*manager.Manager
	status.Reporter
}

func (c *Client) DownloadRetry(ctx context.Context, params DownloadRetryParams) status.Job {
	rootDestination, _ := filepath.Split(params.LocalPath)
	return c.DownloadFolder(ctx,
		DownloadFolderParams{
			FolderListForParams: files_sdk.FolderListForParams{Path: params.File.RemotePath},
			Sync:                params.Sync,
			Manager:             params.Manager,
			Reporter:            params.Reporter,
			RootDestination:     rootDestination,
			JobId:               params.Job.Id,
		})
}

func (c *Client) DownloadToFile(ctx context.Context, params files_sdk.FileDownloadParams, filePath string) (files_sdk.File, error) {
	out, err := os.Create(filePath)
	if err != nil {
		return files_sdk.File{}, err
	}
	params.Writer = out
	return c.Download(ctx, params)
}

func DownloadToFile(ctx context.Context, params files_sdk.FileDownloadParams, filePath string) (files_sdk.File, error) {
	return (&Client{}).DownloadToFile(ctx, params, filePath)
}

type DownloadFolderParams struct {
	files_sdk.FolderListForParams
	Sync bool
	*manager.Manager
	status.Reporter
	RootDestination string
	JobId           string
}

func (c *Client) DownloadFolder(ctx context.Context, params DownloadFolderParams) status.Job {
	return downloadFolder(ctx, c.index(ctx, params.FolderListForParams), c, params)
}

type Entity struct {
	file files_sdk.File
	error
}

func (c *Client) index(ctx context.Context, params files_sdk.FolderListForParams) []Entity {
	var files []Entity
	folderClient := folder.Client{Config: c.Config}
	it, err := folderClient.ListFor(ctx, params)

	if err != nil {
		files = append(files, Entity{file: files_sdk.File{Path: params.Path, Type: "error"}, error: err})
	}

	for it.Next() {
		b, err := json.Marshal(it.Folder())
		if err != nil {
			files = append(files, Entity{file: files_sdk.File{Path: params.Path, Type: "error"}, error: err})
			continue
		}
		entry := files_sdk.File{}
		err = entry.UnmarshalJSON(b)
		if err != nil {
			files = append(files, Entity{file: files_sdk.File{Path: params.Path, Type: "error"}, error: err})
			continue
		}
		switch entry.Type {
		case "directory":
			files = append(files, c.index(ctx, files_sdk.FolderListForParams{Path: entry.Path})...)
		case "file":
			files = append(files, Entity{file: entry})
		default:
			files = append(files, Entity{file: entry, error: fmt.Errorf("unknown file type %v", entry.Type)})
		}
	}

	if it.Err() != nil {
		files = append(files, Entity{file: files_sdk.File{Path: params.Path}, error: it.Err()})
	}
	return files
}
