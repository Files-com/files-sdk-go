package lock

import (
	"context"

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

func (i *Iter) Lock() files_sdk.Lock {
	return i.Current().(files_sdk.Lock)
}

func (c *Client) ListFor(ctx context.Context, params files_sdk.LockListForParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/locks/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.LockCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func ListFor(ctx context.Context, params files_sdk.LockListForParams) (*Iter, error) {
	return (&Client{}).ListFor(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.LockCreateParams) (lock files_sdk.Lock, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/locks/{path}", Params: params, Entity: &lock})
	return
}

func Create(ctx context.Context, params files_sdk.LockCreateParams) (lock files_sdk.Lock, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.LockDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/locks/{path}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.LockDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
