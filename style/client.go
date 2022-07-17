package style

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(ctx context.Context, params files_sdk.StyleFindParams) (style files_sdk.Style, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/styles/{path}", Params: params, Entity: &style})
	return
}

func Find(ctx context.Context, params files_sdk.StyleFindParams) (style files_sdk.Style, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.StyleUpdateParams) (style files_sdk.Style, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/styles/{path}", Params: params, Entity: &style})
	return
}

func Update(ctx context.Context, params files_sdk.StyleUpdateParams) (style files_sdk.Style, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.StyleDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/styles/{path}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.StyleDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
