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

func (c *Client) List(ctx context.Context, params files_sdk.GroupListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/groups", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.GroupCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.GroupListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.GroupFindParams, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/groups/{id}", Params: params, Entity: &group}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.GroupFindParams, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.GroupCreateParams, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/groups", Params: params, Entity: &group}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.GroupCreateParams, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.GroupUpdateParams, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/groups/{id}", Params: params, Entity: &group}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.GroupUpdateParams, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/groups/{id}", Params: params, Entity: &group}, opts...)
	return
}

func UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (group files_sdk.Group, err error) {
	return (&Client{}).UpdateWithMap(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.GroupDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/groups/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.GroupDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
