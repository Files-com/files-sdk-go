package request

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
	*files_sdk.Iter
	*Client
}

func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

func (i *Iter) Request() files_sdk.Request {
	return i.Current().(files_sdk.Request)
}

func (c *Client) List(ctx context.Context, params files_sdk.RequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/requests", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.RequestCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.RequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) GetFolder(ctx context.Context, params files_sdk.RequestGetFolderParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/requests/folders/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.RequestCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func GetFolder(ctx context.Context, params files_sdk.RequestGetFolderParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).GetFolder(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.RequestCreateParams, opts ...files_sdk.RequestResponseOption) (request files_sdk.Request, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/requests", Params: params, Entity: &request}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.RequestCreateParams, opts ...files_sdk.RequestResponseOption) (request files_sdk.Request, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.RequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/requests/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.RequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
