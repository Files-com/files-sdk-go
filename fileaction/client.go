package file_action

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  files_sdk "github.com/Files-com/files-sdk-go"
)

type Client struct {
	files_sdk.Config
}


func (c *Client) Copy (params files_sdk.FileActionCopyParams) (files_sdk.FileAction, error) {
  fileAction := files_sdk.FileAction{}
		path := "/file_actions/copy/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return fileAction, err
	}
	if err := fileAction.UnmarshalJSON(*data); err != nil {
	return fileAction, err
	}

	return  fileAction, nil
}

func Copy (params files_sdk.FileActionCopyParams) (files_sdk.FileAction, error) {
  client := Client{}
  return client.Copy (params)
}

func (c *Client) Move (params files_sdk.FileActionMoveParams) (files_sdk.FileAction, error) {
  fileAction := files_sdk.FileAction{}
		path := "/file_actions/move/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return fileAction, err
	}
	if err := fileAction.UnmarshalJSON(*data); err != nil {
	return fileAction, err
	}

	return  fileAction, nil
}

func Move (params files_sdk.FileActionMoveParams) (files_sdk.FileAction, error) {
  client := Client{}
  return client.Move (params)
}

func (c *Client) BeginUpload (params files_sdk.FileActionBeginUploadParams) (files_sdk.FileAction, error) {
  fileAction := files_sdk.FileAction{}
		path := "/file_actions/begin_upload/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return fileAction, err
	}
	if err := fileAction.UnmarshalJSON(*data); err != nil {
	return fileAction, err
	}

	return  fileAction, nil
}

func BeginUpload (params files_sdk.FileActionBeginUploadParams) (files_sdk.FileAction, error) {
  client := Client{}
  return client.BeginUpload (params)
}
