package folder

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  files_sdk "github.com/Files-com/files-sdk-go"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*lib.Iter
}

func (i *Iter) Folder() files_sdk.Folder {
	return i.Current().(files_sdk.Folder)
}

func (c *Client) ListFor(params files_sdk.FolderListForParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/folders/" + lib.QueryEscape(params.Path) + ""

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
        if err != nil {
          return &defaultValue, "", err
        }
		list := files_sdk.FolderCollection{}
		if err := list.UnmarshalJSON(*data); err != nil {
          return &defaultValue, "", err
        }

		ret := make([]interface{}, len(list))
		for i, v := range list {
			ret[i] = v
		}
		cursor := res.Header.Get("X-Files-Cursor")
		return &ret, cursor, nil
	}
	i.ListParams = &params
	return i
}

func ListFor(params files_sdk.FolderListForParams) *Iter {
  client := Client{}
  return client.ListFor (params)
}

func (c *Client) Create (params files_sdk.FolderCreateParams) (files_sdk.Folder, error) {
  folder := files_sdk.Folder{}
		path := "/folders/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return folder, err
	}
	if err := folder.UnmarshalJSON(*data); err != nil {
	return folder, err
	}

	return  folder, nil
}

func Create (params files_sdk.FolderCreateParams) (files_sdk.Folder, error) {
  client := Client{}
  return client.Create (params)
}
