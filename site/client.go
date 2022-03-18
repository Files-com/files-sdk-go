package site

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Get(ctx context.Context) (files_sdk.Site, error) {
	site := files_sdk.Site{}
	path := "/site"
	exportedParams := lib.Params{Params: lib.Interface()}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return site, err
	}
	if res.StatusCode == 204 {
		return site, nil
	}
	if err := site.UnmarshalJSON(*data); err != nil {
		return site, err
	}

	return site, nil
}

func Get(ctx context.Context) (files_sdk.Site, error) {
	return (&Client{}).Get(ctx)
}

func (c *Client) GetUsage(ctx context.Context) (files_sdk.UsageSnapshot, error) {
	usageSnapshot := files_sdk.UsageSnapshot{}
	path := "/site/usage"
	exportedParams := lib.Params{Params: lib.Interface()}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return usageSnapshot, err
	}
	if res.StatusCode == 204 {
		return usageSnapshot, nil
	}
	if err := usageSnapshot.UnmarshalJSON(*data); err != nil {
		return usageSnapshot, err
	}

	return usageSnapshot, nil
}

func GetUsage(ctx context.Context) (files_sdk.UsageSnapshot, error) {
	return (&Client{}).GetUsage(ctx)
}

func (c *Client) Update(ctx context.Context, params files_sdk.SiteUpdateParams) (files_sdk.Site, error) {
	site := files_sdk.Site{}
	path := "/site"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return site, err
	}
	if res.StatusCode == 204 {
		return site, nil
	}
	if err := site.UnmarshalJSON(*data); err != nil {
		return site, err
	}

	return site, nil
}

func Update(ctx context.Context, params files_sdk.SiteUpdateParams) (files_sdk.Site, error) {
	return (&Client{}).Update(ctx, params)
}
