package style

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  files_sdk "github.com/Files-com/files-sdk-go"
)

type Client struct {
	files_sdk.Config
}


func (c *Client) Find (params files_sdk.StyleFindParams) (files_sdk.Style, error) {
  style := files_sdk.Style{}
		path := "/styles/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return style, err
	}
	if err := style.UnmarshalJSON(*data); err != nil {
	return style, err
	}

	return  style, nil
}

func Find (params files_sdk.StyleFindParams) (files_sdk.Style, error) {
  client := Client{}
  return client.Find (params)
}

func (c *Client) Update (params files_sdk.StyleUpdateParams) (files_sdk.Style, error) {
  style := files_sdk.Style{}
		path := "/styles/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return style, err
	}
	if err := style.UnmarshalJSON(*data); err != nil {
	return style, err
	}

	return  style, nil
}

func Update (params files_sdk.StyleUpdateParams) (files_sdk.Style, error) {
  client := Client{}
  return client.Update (params)
}

func (c *Client) Delete (params files_sdk.StyleDeleteParams) (files_sdk.Style, error) {
  style := files_sdk.Style{}
		path := "/styles/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return style, err
	}
	if err := style.UnmarshalJSON(*data); err != nil {
	return style, err
	}

	return  style, nil
}

func Delete (params files_sdk.StyleDeleteParams) (files_sdk.Style, error) {
  client := Client{}
  return client.Delete (params)
}
