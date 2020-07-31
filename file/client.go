package file

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  files_sdk "github.com/Files-com/files-sdk-go"
)

type Client struct {
	files_sdk.Config
}


func (c *Client) Download (params files_sdk.FileDownloadParams) (files_sdk.File, error) {
  file := files_sdk.File{}
		path := "/files/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return file, err
	}
	if err := file.UnmarshalJSON(*data); err != nil {
	return file, err
	}

	return  file, nil
}

func Download (params files_sdk.FileDownloadParams) (files_sdk.File, error) {
  client := Client{}
  return client.Download (params)
}

func (c *Client) Create (params files_sdk.FileCreateParams) (files_sdk.File, error) {
  file := files_sdk.File{}
		path := "/files/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return file, err
	}
	if err := file.UnmarshalJSON(*data); err != nil {
	return file, err
	}

	return  file, nil
}

func Create (params files_sdk.FileCreateParams) (files_sdk.File, error) {
  client := Client{}
  return client.Create (params)
}

func (c *Client) Update (params files_sdk.FileUpdateParams) (files_sdk.File, error) {
  file := files_sdk.File{}
		path := "/files/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return file, err
	}
	if err := file.UnmarshalJSON(*data); err != nil {
	return file, err
	}

	return  file, nil
}

func Update (params files_sdk.FileUpdateParams) (files_sdk.File, error) {
  client := Client{}
  return client.Update (params)
}

func (c *Client) Delete (params files_sdk.FileDeleteParams) (files_sdk.File, error) {
  file := files_sdk.File{}
		path := "/files/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return file, err
	}
	if err := file.UnmarshalJSON(*data); err != nil {
	return file, err
	}

	return  file, nil
}

func Delete (params files_sdk.FileDeleteParams) (files_sdk.File, error) {
  client := Client{}
  return client.Delete (params)
}
