package file

import (
	"context"
	"io"
	"net/http"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(ctx context.Context, Path string) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/files/", Path)
	exportParams, err := lib.ExportParams(lib.Interface())
	if err != nil {
		return file, err
	}
	data, _, err := files_sdk.Call(ctx, "GET", c.Config, path, exportParams)
	if err != nil {
		return file, err
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Find(ctx context.Context, Path string) (files_sdk.File, error) {
	client := Client{}
	return client.Find(ctx, Path)
}

func (c *Client) Download(ctx context.Context, params files_sdk.FileDownloadParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/files/", params.Path)
	exportParams, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, _, err := files_sdk.Call(ctx, "GET", c.Config, path, exportParams)
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
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Create(ctx context.Context, params files_sdk.FileCreateParams) (files_sdk.File, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.FileUpdateParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/files/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Update(ctx context.Context, params files_sdk.FileUpdateParams) (files_sdk.File, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.FileDeleteParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/files/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Delete(ctx context.Context, params files_sdk.FileDeleteParams) (files_sdk.File, error) {
	return (&Client{}).Delete(ctx, params)
}

func (c *Client) Metadata(ctx context.Context, params files_sdk.FileMetadataParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/file_actions/metadata/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Metadata(ctx context.Context, params files_sdk.FileMetadataParams) (files_sdk.File, error) {
	return (&Client{}).Metadata(ctx, params)
}

func (c *Client) Copy(ctx context.Context, params files_sdk.FileCopyParams) (files_sdk.FileAction, error) {
	fileAction := files_sdk.FileAction{}
	path := lib.BuildPath("/file_actions/copy/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return fileAction, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return fileAction, err
	}
	if res.StatusCode == 204 {
		return fileAction, nil
	}
	if err := fileAction.UnmarshalJSON(*data); err != nil {
		return fileAction, err
	}

	return fileAction, nil
}

func Copy(ctx context.Context, params files_sdk.FileCopyParams) (files_sdk.FileAction, error) {
	return (&Client{}).Copy(ctx, params)
}

func (c *Client) Move(ctx context.Context, params files_sdk.FileMoveParams) (files_sdk.FileAction, error) {
	fileAction := files_sdk.FileAction{}
	path := lib.BuildPath("/file_actions/move/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return fileAction, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return fileAction, err
	}
	if res.StatusCode == 204 {
		return fileAction, nil
	}
	if err := fileAction.UnmarshalJSON(*data); err != nil {
		return fileAction, err
	}

	return fileAction, nil
}

func Move(ctx context.Context, params files_sdk.FileMoveParams) (files_sdk.FileAction, error) {
	return (&Client{}).Move(ctx, params)
}

func (c *Client) BeginUpload(ctx context.Context, params files_sdk.FileBeginUploadParams) (files_sdk.FileUploadPartCollection, error) {
	fileUploadPartCollection := files_sdk.FileUploadPartCollection{}
	path := lib.BuildPath("/file_actions/begin_upload/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return fileUploadPartCollection, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return fileUploadPartCollection, err
	}
	if res.StatusCode == 204 {
		return fileUploadPartCollection, nil
	}
	if err := fileUploadPartCollection.UnmarshalJSON(*data); err != nil {
		return fileUploadPartCollection, err
	}

	return fileUploadPartCollection, nil
}

func BeginUpload(ctx context.Context, params files_sdk.FileBeginUploadParams) (files_sdk.FileUploadPartCollection, error) {
	return (&Client{}).BeginUpload(ctx, params)
}
