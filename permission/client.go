package permission

import (
	"context"
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
	listquery "github.com/Files-com/files-sdk-go/v2/listquery"
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

func (c *Client) List(ctx context.Context, params files_sdk.PermissionListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/permissions"
	i.ListParams = &params
	list := files_sdk.PermissionCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.PermissionListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.PermissionCreateParams) (files_sdk.Permission, error) {
	permission := files_sdk.Permission{}
	path := "/permissions"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return permission, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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

func Create(ctx context.Context, params files_sdk.PermissionCreateParams) (files_sdk.Permission, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.PermissionDeleteParams) (files_sdk.Permission, error) {
	permission := files_sdk.Permission{}
	if params.Id == 0 {
		return permission, lib.CreateError(params, "Id")
	}
	path := "/permissions/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return permission, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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

func Delete(ctx context.Context, params files_sdk.PermissionDeleteParams) (files_sdk.Permission, error) {
	return (&Client{}).Delete(ctx, params)
}
