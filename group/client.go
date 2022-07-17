package group

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

func (i *Iter) Group() files_sdk.Group {
	return i.Current().(files_sdk.Group)
}

func (c *Client) List(ctx context.Context, params files_sdk.GroupListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/groups", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.GroupCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.GroupListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.GroupFindParams) (group files_sdk.Group, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/groups/{id}", Params: params, Entity: &group})
	return
}

func Find(ctx context.Context, params files_sdk.GroupFindParams) (group files_sdk.Group, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.GroupCreateParams) (group files_sdk.Group, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/groups", Params: params, Entity: &group})
	return
}

func Create(ctx context.Context, params files_sdk.GroupCreateParams) (group files_sdk.Group, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.GroupUpdateParams) (group files_sdk.Group, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/groups/{id}", Params: params, Entity: &group})
	return
}

func Update(ctx context.Context, params files_sdk.GroupUpdateParams) (group files_sdk.Group, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.GroupDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/groups/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.GroupDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
