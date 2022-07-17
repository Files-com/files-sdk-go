package permission

import (
	"context"

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
	path, err := lib.BuildPath("/permissions", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PermissionCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.PermissionListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.PermissionCreateParams) (permission files_sdk.Permission, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/permissions", Params: params, Entity: &permission})
	return
}

func Create(ctx context.Context, params files_sdk.PermissionCreateParams) (permission files_sdk.Permission, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.PermissionDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/permissions/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.PermissionDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
