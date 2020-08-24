package site

import (
	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Get(params files_sdk.SiteGetParams) (files_sdk.Site, error) {
	site := files_sdk.Site{}
	path := "/site"
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
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

func Get(params files_sdk.SiteGetParams) (files_sdk.Site, error) {
	return (&Client{}).Get(params)
}

func (c *Client) GetUsage(params files_sdk.SiteGetUsageParams) (files_sdk.UsageSnapshot, error) {
	usageSnapshot := files_sdk.UsageSnapshot{}
	path := "/site/usage"
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
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

func GetUsage(params files_sdk.SiteGetUsageParams) (files_sdk.UsageSnapshot, error) {
	return (&Client{}).GetUsage(params)
}

func (c *Client) Update(params files_sdk.SiteUpdateParams) (files_sdk.Site, error) {
	site := files_sdk.Site{}
	path := "/site"
	data, res, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
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

func Update(params files_sdk.SiteUpdateParams) (files_sdk.Site, error) {
	return (&Client{}).Update(params)
}
