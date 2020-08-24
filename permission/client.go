package permission

import (
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*lib.Iter
}

func (i *Iter) Permission() files_sdk.Permission {
	return i.Current().(files_sdk.Permission)
}

func (c *Client) List(params files_sdk.PermissionListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/permissions"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.PermissionCollection{}
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

func List(params files_sdk.PermissionListParams) *Iter {
	return (&Client{}).List(params)
}

func (c *Client) Create(params files_sdk.PermissionCreateParams) (files_sdk.Permission, error) {
	permission := files_sdk.Permission{}
	path := "/permissions"
	data, res, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
		return permission, err
	}
	if res.StatusCode == 204 {
		return permission, nil
	}
	if err := permission.UnmarshalJSON(*data); err != nil {
		return permission, err
	}

	return permission, nil
}

func Create(params files_sdk.PermissionCreateParams) (files_sdk.Permission, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Delete(params files_sdk.PermissionDeleteParams) (files_sdk.Permission, error) {
	permission := files_sdk.Permission{}
	path := "/permissions/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
	if err != nil {
		return permission, err
	}
	if res.StatusCode == 204 {
		return permission, nil
	}
	if err := permission.UnmarshalJSON(*data); err != nil {
		return permission, err
	}

	return permission, nil
}

func Delete(params files_sdk.PermissionDeleteParams) (files_sdk.Permission, error) {
	return (&Client{}).Delete(params)
}
