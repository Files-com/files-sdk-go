package folder

import (
	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
	listquery "github.com/Files-com/files-sdk-go/listquery"
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
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := lib.BuildPath("/folders/", params.Path)
	i.ListParams = &params
	list := files_sdk.FolderCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
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

func Create(params files_sdk.FolderCreateParams) (files_sdk.File, error) {
	return (&Client{}).Create(params)
}
