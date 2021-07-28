package file_action

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Copy(ctx context.Context, params files_sdk.FileActionCopyParams) (files_sdk.FileAction, error) {
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

func Copy(ctx context.Context, params files_sdk.FileActionCopyParams) (files_sdk.FileAction, error) {
	return (&Client{}).Copy(ctx, params)
}

func (c *Client) Move(ctx context.Context, params files_sdk.FileActionMoveParams) (files_sdk.FileAction, error) {
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

func Move(ctx context.Context, params files_sdk.FileActionMoveParams) (files_sdk.FileAction, error) {
	return (&Client{}).Move(ctx, params)
}

func (c *Client) BeginUpload(ctx context.Context, params files_sdk.FileActionBeginUploadParams) (files_sdk.FileUploadPartCollection, error) {
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

func BeginUpload(ctx context.Context, params files_sdk.FileActionBeginUploadParams) (files_sdk.FileUploadPartCollection, error) {
	return (&Client{}).BeginUpload(ctx, params)
}
