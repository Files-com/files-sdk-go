package site

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Get(opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/site", Params: lib.Interface(), Entity: &site}, opts...)
	return
}

func Get(opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	return (&Client{}).Get(opts...)
}

func (c *Client) GetUsage(opts ...files_sdk.RequestResponseOption) (usageSnapshot files_sdk.UsageSnapshot, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/site/usage", Params: lib.Interface(), Entity: &usageSnapshot}, opts...)
	return
}

func GetUsage(opts ...files_sdk.RequestResponseOption) (usageSnapshot files_sdk.UsageSnapshot, err error) {
	return (&Client{}).GetUsage(opts...)
}

func (c *Client) Update(params files_sdk.SiteUpdateParams, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/site", Params: params, Entity: &site}, opts...)
	return
}

func Update(params files_sdk.SiteUpdateParams, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/site", Params: params, Entity: &site}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}
