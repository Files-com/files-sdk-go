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

func (c *Client) List(ctx context.Context, params files_sdk.GroupUserListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/group_users", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.GroupUserCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.GroupUserListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.GroupUserCreateParams) (groupUser files_sdk.GroupUser, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/group_users", Params: params, Entity: &groupUser})
	return
}

func Create(ctx context.Context, params files_sdk.GroupUserCreateParams) (groupUser files_sdk.GroupUser, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.GroupUserUpdateParams) (groupUser files_sdk.GroupUser, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/group_users/{id}", Params: params, Entity: &groupUser})
	return
}

func Update(ctx context.Context, params files_sdk.GroupUserUpdateParams) (groupUser files_sdk.GroupUser, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.GroupUserDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/group_users/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.GroupUserDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
