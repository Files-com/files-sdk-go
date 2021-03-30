package lock

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

func (i *Iter) Lock() files_sdk.Lock {
	return i.Current().(files_sdk.Lock)
}

func (c *Client) ListFor(params files_sdk.LockListForParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := lib.BuildPath("/locks/", params.Path)
	i.ListParams = &params
	list := files_sdk.LockCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func ListFor(params files_sdk.LockListForParams) (*Iter, error) {
	return (&Client{}).ListFor(params)
}

func (c *Client) Create(params files_sdk.LockCreateParams) (files_sdk.Lock, error) {
	lock := files_sdk.Lock{}
	path := lib.BuildPath("/locks/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return lock, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	if err != nil {
		return lock, err
	}
	if res.StatusCode == 204 {
		return lock, nil
	}
	if err := lock.UnmarshalJSON(*data); err != nil {
		return lock, err
	}

	return lock, nil
}

func Create(params files_sdk.LockCreateParams) (files_sdk.Lock, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Delete(params files_sdk.LockDeleteParams) (files_sdk.Lock, error) {
	lock := files_sdk.Lock{}
	path := lib.BuildPath("/locks/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return lock, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	if err != nil {
		return lock, err
	}
	if res.StatusCode == 204 {
		return lock, nil
	}
	if err := lock.UnmarshalJSON(*data); err != nil {
		return lock, err
	}

	return lock, nil
}

func Delete(params files_sdk.LockDeleteParams) (files_sdk.Lock, error) {
	return (&Client{}).Delete(params)
}
