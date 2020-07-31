package lock

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

func (i *Iter) Lock() files_sdk.Lock {
	return i.Current().(files_sdk.Lock)
}

func (c *Client) ListFor(params files_sdk.LockListForParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/locks/" + lib.QueryEscape(params.Path) + ""

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
        if err != nil {
          return &defaultValue, "", err
        }
		list := files_sdk.LockCollection{}
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

func ListFor(params files_sdk.LockListForParams) *Iter {
  client := Client{}
  return client.ListFor (params)
}

func (c *Client) Create (params files_sdk.LockCreateParams) (files_sdk.Lock, error) {
  lock := files_sdk.Lock{}
		path := "/locks/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return lock, err
	}
	if err := lock.UnmarshalJSON(*data); err != nil {
	return lock, err
	}

	return  lock, nil
}

func Create (params files_sdk.LockCreateParams) (files_sdk.Lock, error) {
  client := Client{}
  return client.Create (params)
}

func (c *Client) Delete (params files_sdk.LockDeleteParams) (files_sdk.Lock, error) {
  lock := files_sdk.Lock{}
		path := "/locks/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return lock, err
	}
	if err := lock.UnmarshalJSON(*data); err != nil {
	return lock, err
	}

	return  lock, nil
}

func Delete (params files_sdk.LockDeleteParams) (files_sdk.Lock, error) {
  client := Client{}
  return client.Delete (params)
}
