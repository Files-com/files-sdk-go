// Package site provides the Files.com Site API client.
package site

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Client calls the Site API using its embedded Config.
type Client struct {
	files_sdk.Config
}

// Get calls GET /site.
func (c *Client) Get(opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/site", Entity: &site}, opts...)
	return
}

// Get calls Client.Get using the default configuration.
func Get(opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	return (&Client{}).Get(opts...)
}

// GetUsage calls GET /site/usage.
//
// API operation: Get the most recent usage snapshot (usage data for billing purposes) for a Site.
func (c *Client) GetUsage(opts ...files_sdk.RequestResponseOption) (usageSnapshot files_sdk.UsageSnapshot, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/site/usage", Entity: &usageSnapshot}, opts...)
	return
}

// GetUsage calls Client.GetUsage using the default configuration.
func GetUsage(opts ...files_sdk.RequestResponseOption) (usageSnapshot files_sdk.UsageSnapshot, err error) {
	return (&Client{}).GetUsage(opts...)
}

// Update calls PATCH /site.
func (c *Client) Update(params files_sdk.SiteUpdateParams, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/site", Params: params, Entity: &site}, opts...)
	return
}

// Update calls Client.Update using the default configuration.
func Update(params files_sdk.SiteUpdateParams, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	return (&Client{}).Update(params, opts...)
}

// UpdateWithMap calls PATCH /site using API parameter names as map keys.
// Include any path parameters in the map. Unlike optional struct fields, explicit
// zero values in the map are included in the request.
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/site", Params: params, Entity: &site}, opts...)
	return
}

// UpdateWithMap calls Client.UpdateWithMap using the default configuration.
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}
