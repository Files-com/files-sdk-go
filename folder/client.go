package folder

import (
	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
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

func (c *Client) ListFor(params files_sdk.FolderListForParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := lib.BuildPath("/folders/", params.Path)
	i.ListParams = &params
	exportParams, err := i.ExportParams()
	if err != nil {
		return i, err
	}
	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, exportParams)
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
	return i, nil
}

func ListFor(params files_sdk.FolderListForParams) (*Iter, error) {
	return (&Client{}).ListFor(params)
}

func (c *Client) Create(params files_sdk.FolderCreateParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := lib.BuildPath("/folders/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
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

func Create(params files_sdk.FolderCreateParams) (files_sdk.File, error) {
	return (&Client{}).Create(params)
}
