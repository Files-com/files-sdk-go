package file

import (
	"context"
	"io"
	goFs "io/fs"
	"net/http"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/folder"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Get(ctx context.Context, Path string) (files_sdk.File, error) {
	file := files_sdk.File{}
	path, err := lib.BuildPath("/files/{path}", map[string]string{"path": Path})
	if err != nil {
		return file, err
	}
	data, _, err := files_sdk.Call(ctx, "GET", c.Config, path, lib.Params{Params: lib.Interface()})
	if err != nil {
		return file, err
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Get(ctx context.Context, Path string) (files_sdk.File, error) {
	client := Client{}
	return client.Get(ctx, Path)
}

func (c *Client) Download(ctx context.Context, params files_sdk.FileDownloadParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path, err := lib.BuildPath("/files/{path}", params)
	if err != nil {
		return file, err
	}
	data, _, err := files_sdk.Call(ctx, "GET", c.Config, path, lib.Params{Params: lib.Interface()})
	if err != nil {
		return file, err
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}
	request, err := http.NewRequestWithContext(ctx, "GET", file.DownloadUri, nil)
	if err != nil {
		return file, err
	}
	resp, err := c.Config.GetHttpClient().Do(request)
	if err != nil {
		return file, err
	}
	if params.OnDownload != nil {
		params.OnDownload(resp)
	}
	_, err = io.Copy(params.Writer, resp.Body)
	if err != nil {
		return file, err
	}

	return file, nil
}

func Download(ctx context.Context, params files_sdk.FileDownloadParams) (files_sdk.File, error) {
	client := Client{}
	return client.Download(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.FileCreateParams) (file files_sdk.File, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/files/{path}", Params: params, Entity: &file})
	return
}

func Create(ctx context.Context, params files_sdk.FileCreateParams) (file files_sdk.File, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.FileUpdateParams) (file files_sdk.File, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/files/{path}", Params: params, Entity: &file})
	return
}

func Update(ctx context.Context, params files_sdk.FileUpdateParams) (file files_sdk.File, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.FileDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/files/{path}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.FileDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.FileFindParams) (file files_sdk.File, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/file_actions/metadata/{path}", Params: params, Entity: &file})
	return
}

func Find(ctx context.Context, params files_sdk.FileFindParams) (file files_sdk.File, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Copy(ctx context.Context, params files_sdk.FileCopyParams) (fileAction files_sdk.FileAction, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/file_actions/copy/{path}", Params: params, Entity: &fileAction})
	return
}

func Copy(ctx context.Context, params files_sdk.FileCopyParams) (fileAction files_sdk.FileAction, err error) {
	return (&Client{}).Copy(ctx, params)
}

func (c *Client) Move(ctx context.Context, params files_sdk.FileMoveParams) (fileAction files_sdk.FileAction, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/file_actions/move/{path}", Params: params, Entity: &fileAction})
	return
}

func Move(ctx context.Context, params files_sdk.FileMoveParams) (fileAction files_sdk.FileAction, err error) {
	return (&Client{}).Move(ctx, params)
}

func (c *Client) BeginUpload(ctx context.Context, params files_sdk.FileBeginUploadParams) (fileUploadPartCollection files_sdk.FileUploadPartCollection, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/file_actions/begin_upload/{path}", Params: params, Entity: &fileUploadPartCollection})
	return
}

func BeginUpload(ctx context.Context, params files_sdk.FileBeginUploadParams) (fileUploadPartCollection files_sdk.FileUploadPartCollection, err error) {
	return (&Client{}).BeginUpload(ctx, params)
}

func (c *Client) ListFor(ctx context.Context, params files_sdk.FolderListForParams) (*folder.Iter, error) {
	client := folder.Client{Config: c.Config}
	return client.ListFor(ctx, params)
}

func (c *Client) ListForRecursive(ctx context.Context, params files_sdk.FolderListForParams) (lib.IterI, error) {
	it := lib.IterChan{}.Init()

	go func(params files_sdk.FolderListForParams) {
		f := FS{}.Init(c.Config).WithContext(ctx)
		err := goFs.WalkDir(f, params.Path, func(path string, d goFs.DirEntry, err error) error {
			if path == "" && err == nil {
				return nil // Skip root directory
			}

			if err == nil {
				info, _ := d.Info()
				it.Send <- info.Sys()
			} else {
				it.SendError <- err
			}
			return err
		})
		if err != nil {
			it.Error.Store(err)
		}
		it.Stop <- true
	}(params)
	return it, nil
}
