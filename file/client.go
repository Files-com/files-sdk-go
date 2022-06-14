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
	path := lib.BuildPath("/files/", Path)
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
	path := lib.BuildPath("/files/", params.Path)
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

func (c *Client) Create(ctx context.Context, params files_sdk.FileCreateParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/files/", params.Path)
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}

	return file, file.UnmarshalJSON(*data)
}

func Create(ctx context.Context, params files_sdk.FileCreateParams) (files_sdk.File, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.FileUpdateParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/files/", params.Path)
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}

	return file, file.UnmarshalJSON(*data)
}

func Update(ctx context.Context, params files_sdk.FileUpdateParams) (files_sdk.File, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.FileDeleteParams) error {
	file := files_sdk.File{}
	path := lib.BuildPath("/files/", params.Path)
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	if res.StatusCode == 204 {
		return nil
	}

	return file.UnmarshalJSON(*data)
}

func Delete(ctx context.Context, params files_sdk.FileDeleteParams) error {
	return (&Client{}).Delete(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.FileFindParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/file_actions/metadata/", params.Path)
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}

	return file, file.UnmarshalJSON(*data)
}

func Find(ctx context.Context, params files_sdk.FileFindParams) (files_sdk.File, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Copy(ctx context.Context, params files_sdk.FileCopyParams) (files_sdk.FileAction, error) {
	fileAction := files_sdk.FileAction{}
	path := lib.BuildPath("/file_actions/copy/", params.Path)
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return fileAction, err
	}
	if res.StatusCode == 204 {
		return fileAction, nil
	}

	return fileAction, fileAction.UnmarshalJSON(*data)
}

func Copy(ctx context.Context, params files_sdk.FileCopyParams) (files_sdk.FileAction, error) {
	return (&Client{}).Copy(ctx, params)
}

func (c *Client) Move(ctx context.Context, params files_sdk.FileMoveParams) (files_sdk.FileAction, error) {
	fileAction := files_sdk.FileAction{}
	path := lib.BuildPath("/file_actions/move/", params.Path)
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return fileAction, err
	}
	if res.StatusCode == 204 {
		return fileAction, nil
	}

	return fileAction, fileAction.UnmarshalJSON(*data)
}

func Move(ctx context.Context, params files_sdk.FileMoveParams) (files_sdk.FileAction, error) {
	return (&Client{}).Move(ctx, params)
}

func (c *Client) BeginUpload(ctx context.Context, params files_sdk.FileBeginUploadParams) (files_sdk.FileUploadPartCollection, error) {
	fileUploadPartCollection := files_sdk.FileUploadPartCollection{}
	path := lib.BuildPath("/file_actions/begin_upload/", params.Path)
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return fileUploadPartCollection, err
	}
	if res.StatusCode == 204 {
		return fileUploadPartCollection, nil
	}

	return fileUploadPartCollection, fileUploadPartCollection.UnmarshalJSON(*data)
}

func BeginUpload(ctx context.Context, params files_sdk.FileBeginUploadParams) (files_sdk.FileUploadPartCollection, error) {
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
