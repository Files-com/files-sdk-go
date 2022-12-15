package style

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(ctx context.Context, params files_sdk.StyleFindParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/styles/{path}", Params: params, Entity: &style}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.StyleFindParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.StyleUpdateParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/styles/{path}", Params: params, Entity: &style}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.StyleUpdateParams, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/styles/{path}", Params: params, Entity: &style}, opts...)
	return
}

func UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (style files_sdk.Style, err error) {
	return (&Client{}).UpdateWithMap(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.StyleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/styles/{path}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.StyleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
