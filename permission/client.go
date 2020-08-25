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

func (c *Client) List(params files_sdk.PermissionListParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/permissions"
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
	return i, nil
}

func List(params files_sdk.PermissionListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Create(params files_sdk.PermissionCreateParams) (files_sdk.Permission, error) {
	permission := files_sdk.Permission{}
	path := "/permissions"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return permission, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
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
	if params.Id == 0 {
		return permission, lib.CreateError(params, "Id")
	}
	path := "/permissions/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return permission, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParms)
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
