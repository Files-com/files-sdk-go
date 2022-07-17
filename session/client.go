package session

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(ctx context.Context, params files_sdk.SessionCreateParams) (session files_sdk.Session, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/sessions", Params: params, Entity: &session})
	return
}

func Create(ctx context.Context, params files_sdk.SessionCreateParams) (session files_sdk.Session, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/sessions", Params: lib.Interface(), Entity: nil})
	return
}

func Delete(ctx context.Context) (err error) {
	return (&Client{}).Delete(ctx)
}
