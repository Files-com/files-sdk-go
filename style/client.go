// Package style provides the Files.com Style API client.
package style

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Client calls the Style API using its embedded Config.
type Client struct {
	files_sdk.Config
}

// Find calls GET /styles/{path}.
func (c *Client) Find(params files_sdk.StyleFindParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/styles/{path}", Params: params, Entity: &style}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.StyleFindParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	return (&Client{}).Find(params, opts...)
}

// Update calls PATCH /styles/{path}.
func (c *Client) Update(params files_sdk.StyleUpdateParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/styles/{path}", Params: params, Entity: &style}, opts...)
	return
}

// Update calls Client.Update using the default configuration.
func Update(params files_sdk.StyleUpdateParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	return (&Client{}).Update(params, opts...)
}

// UpdateWithMap calls PATCH /styles/{path} using API parameter names as map keys.
// Include any path parameters in the map. Unlike optional struct fields, explicit
// zero values in the map are included in the request.
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/styles/{path}", Params: params, Entity: &style}, opts...)
	return
}

// UpdateWithMap calls Client.UpdateWithMap using the default configuration.
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

// Delete calls DELETE /styles/{path}.
func (c *Client) Delete(params files_sdk.StyleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/styles/{path}", Params: params, Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(params files_sdk.StyleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
