package share_group

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
	*files_sdk.Iter
	*Client
}

func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

func (i *Iter) ShareGroup() files_sdk.ShareGroup {
	return i.Current().(files_sdk.ShareGroup)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.ShareGroupFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(context.Background(), params, opts...)
}

func (c *Client) List(ctx context.Context, params files_sdk.ShareGroupListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/share_groups", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ShareGroupCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.ShareGroupListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.ShareGroupFindParams, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/share_groups/{id}", Params: params, Entity: &shareGroup}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.ShareGroupFindParams, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.ShareGroupCreateParams, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/share_groups", Params: params, Entity: &shareGroup}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.ShareGroupCreateParams, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.ShareGroupUpdateParams, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/share_groups/{id}", Params: params, Entity: &shareGroup}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.ShareGroupUpdateParams, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/share_groups/{id}", Params: params, Entity: &shareGroup}, opts...)
	return
}

func UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (shareGroup files_sdk.ShareGroup, err error) {
	return (&Client{}).UpdateWithMap(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.ShareGroupDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/share_groups/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.ShareGroupDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
