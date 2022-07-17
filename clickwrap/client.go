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
	*lib.Iter
}

func (i *Iter) Clickwrap() files_sdk.Clickwrap {
	return i.Current().(files_sdk.Clickwrap)
}

func (c *Client) List(ctx context.Context, params files_sdk.ClickwrapListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/clickwraps", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ClickwrapCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.ClickwrapListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.ClickwrapFindParams) (clickwrap files_sdk.Clickwrap, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/clickwraps/{id}", Params: params, Entity: &clickwrap})
	return
}

func Find(ctx context.Context, params files_sdk.ClickwrapFindParams) (clickwrap files_sdk.Clickwrap, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.ClickwrapCreateParams) (clickwrap files_sdk.Clickwrap, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/clickwraps", Params: params, Entity: &clickwrap})
	return
}

func Create(ctx context.Context, params files_sdk.ClickwrapCreateParams) (clickwrap files_sdk.Clickwrap, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.ClickwrapUpdateParams) (clickwrap files_sdk.Clickwrap, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/clickwraps/{id}", Params: params, Entity: &clickwrap})
	return
}

func Update(ctx context.Context, params files_sdk.ClickwrapUpdateParams) (clickwrap files_sdk.Clickwrap, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.ClickwrapDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/clickwraps/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.ClickwrapDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
