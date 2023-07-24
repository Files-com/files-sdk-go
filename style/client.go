package style

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(params files_sdk.StyleFindParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/styles/{path}", Params: params, Entity: &style}, opts...)
	return
}

func Find(params files_sdk.StyleFindParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Update(params files_sdk.StyleUpdateParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/styles/{path}", Params: params, Entity: &style}, opts...)
	return
}

func Update(params files_sdk.StyleUpdateParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/styles/{path}", Params: params, Entity: &style}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.StyleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/styles/{path}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.StyleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
