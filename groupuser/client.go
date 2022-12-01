package group_user

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

func (i *Iter) GroupUser() files_sdk.GroupUser {
	return i.Current().(files_sdk.GroupUser)
}

func (c *Client) List(ctx context.Context, params files_sdk.GroupUserListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/group_users", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.GroupUserCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.GroupUserListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.GroupUserCreateParams, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/group_users", Params: params, Entity: &groupUser}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.GroupUserCreateParams, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.GroupUserUpdateParams, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/group_users/{id}", Params: params, Entity: &groupUser}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.GroupUserUpdateParams, opts ...files_sdk.RequestResponseOption) (groupUser files_sdk.GroupUser, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.GroupUserDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/group_users/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.GroupUserDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
