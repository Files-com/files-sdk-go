package sso_strategy

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

func (i *Iter) SsoStrategy() files_sdk.SsoStrategy {
	return i.Current().(files_sdk.SsoStrategy)
}

func (c *Client) List(ctx context.Context, params files_sdk.SsoStrategyListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/sso_strategies", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.SsoStrategyCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.SsoStrategyListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.SsoStrategyFindParams) (ssoStrategy files_sdk.SsoStrategy, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/sso_strategies/{id}", Params: params, Entity: &ssoStrategy})
	return
}

func Find(ctx context.Context, params files_sdk.SsoStrategyFindParams) (ssoStrategy files_sdk.SsoStrategy, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Sync(ctx context.Context, params files_sdk.SsoStrategySyncParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/sso_strategies/{id}/sync", Params: params, Entity: nil})
	return
}

func Sync(ctx context.Context, params files_sdk.SsoStrategySyncParams) (err error) {
	return (&Client{}).Sync(ctx, params)
}
