package clickwrap

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

func (i *Iter) Clickwrap() files_sdk.Clickwrap {
	return i.Current().(files_sdk.Clickwrap)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.ClickwrapFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(context.Background(), params, opts...)
}

func (c *Client) List(ctx context.Context, params files_sdk.ClickwrapListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/clickwraps", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ClickwrapCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.ClickwrapListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.ClickwrapFindParams, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/clickwraps/{id}", Params: params, Entity: &clickwrap}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.ClickwrapFindParams, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.ClickwrapCreateParams, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/clickwraps", Params: params, Entity: &clickwrap}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.ClickwrapCreateParams, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.ClickwrapUpdateParams, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/clickwraps/{id}", Params: params, Entity: &clickwrap}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.ClickwrapUpdateParams, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/clickwraps/{id}", Params: params, Entity: &clickwrap}, opts...)
	return
}

func UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error) {
	return (&Client{}).UpdateWithMap(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.ClickwrapDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/clickwraps/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.ClickwrapDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
