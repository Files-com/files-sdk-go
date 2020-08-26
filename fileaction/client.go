package file_action

import (
	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Copy(params files_sdk.FileActionCopyParams) (files_sdk.FileAction, error) {
	fileAction := files_sdk.FileAction{}
	path := "/file_actions/copy/" + lib.QueryEscape(params.Path) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return fileAction, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
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

func Copy(params files_sdk.FileActionCopyParams) (files_sdk.FileAction, error) {
	return (&Client{}).Copy(params)
}

func (c *Client) Move(params files_sdk.FileActionMoveParams) (files_sdk.FileAction, error) {
	fileAction := files_sdk.FileAction{}
	path := "/file_actions/move/" + lib.QueryEscape(params.Path) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return fileAction, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
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

func Move(params files_sdk.FileActionMoveParams) (files_sdk.FileAction, error) {
	return (&Client{}).Move(params)
}

func (c *Client) BeginUpload(params files_sdk.FileActionBeginUploadParams) (files_sdk.FileUploadPartCollection, error) {
	fileUploadPartCollection := files_sdk.FileUploadPartCollection{}
	path := "/file_actions/begin_upload/" + lib.QueryEscape(params.Path) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return fileUploadPartCollection, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
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

func BeginUpload(params files_sdk.FileActionBeginUploadParams) (files_sdk.FileUploadPartCollection, error) {
	return (&Client{}).BeginUpload(params)
}
