package sso_strategy

import (
	"context"
	"strconv"

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
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/sso_strategies"
	i.ListParams = &params
	list := files_sdk.SsoStrategyCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.SsoStrategyListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.SsoStrategyFindParams) (files_sdk.SsoStrategy, error) {
	ssoStrategy := files_sdk.SsoStrategy{}
	if params.Id == 0 {
		return ssoStrategy, lib.CreateError(params, "Id")
	}
	path := "/sso_strategies/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return ssoStrategy, err
	}
	if res.StatusCode == 204 {
		return ssoStrategy, nil
	}
	if err := ssoStrategy.UnmarshalJSON(*data); err != nil {
		return ssoStrategy, err
	}

	return ssoStrategy, nil
}

func Find(ctx context.Context, params files_sdk.SsoStrategyFindParams) (files_sdk.SsoStrategy, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Sync(ctx context.Context, params files_sdk.SsoStrategySyncParams) (files_sdk.SsoStrategy, error) {
	ssoStrategy := files_sdk.SsoStrategy{}
	if params.Id == 0 {
		return ssoStrategy, lib.CreateError(params, "Id")
	}
	path := "/sso_strategies/" + strconv.FormatInt(params.Id, 10) + "/sync"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return ssoStrategy, err
	}
	if res.StatusCode == 204 {
		return ssoStrategy, nil
	}
	if err := ssoStrategy.UnmarshalJSON(*data); err != nil {
		return ssoStrategy, err
	}

	return ssoStrategy, nil
}

func Sync(ctx context.Context, params files_sdk.SsoStrategySyncParams) (files_sdk.SsoStrategy, error) {
	return (&Client{}).Sync(ctx, params)
}
