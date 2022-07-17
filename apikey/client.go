package api_key

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

func (i *Iter) ApiKey() files_sdk.ApiKey {
	return i.Current().(files_sdk.ApiKey)
}

func (c *Client) List(ctx context.Context, params files_sdk.ApiKeyListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/api_keys", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ApiKeyCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.ApiKeyListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) FindCurrent(ctx context.Context) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/api_key", Params: lib.Interface(), Entity: &apiKey})
	return
}

func FindCurrent(ctx context.Context) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).FindCurrent(ctx)
}

func (c *Client) Find(ctx context.Context, params files_sdk.ApiKeyFindParams) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/api_keys/{id}", Params: params, Entity: &apiKey})
	return
}

func Find(ctx context.Context, params files_sdk.ApiKeyFindParams) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.ApiKeyCreateParams) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/api_keys", Params: params, Entity: &apiKey})
	return
}

func Create(ctx context.Context, params files_sdk.ApiKeyCreateParams) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) UpdateCurrent(ctx context.Context, params files_sdk.ApiKeyUpdateCurrentParams) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/api_key", Params: params, Entity: &apiKey})
	return
}

func UpdateCurrent(ctx context.Context, params files_sdk.ApiKeyUpdateCurrentParams) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).UpdateCurrent(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.ApiKeyUpdateParams) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/api_keys/{id}", Params: params, Entity: &apiKey})
	return
}

func Update(ctx context.Context, params files_sdk.ApiKeyUpdateParams) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) DeleteCurrent(ctx context.Context) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/api_key", Params: lib.Interface(), Entity: nil})
	return
}

func DeleteCurrent(ctx context.Context) (err error) {
	return (&Client{}).DeleteCurrent(ctx)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.ApiKeyDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/api_keys/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.ApiKeyDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
