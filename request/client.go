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
	*lib.Iter
}

func (i *Iter) Request() files_sdk.Request {
	return i.Current().(files_sdk.Request)
}

func (c *Client) List(ctx context.Context, params files_sdk.RequestListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/requests", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.RequestCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.RequestListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) GetFolder(ctx context.Context, params files_sdk.RequestGetFolderParams) (requestCollection files_sdk.RequestCollection, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/requests/folders/{path}", Params: params, Entity: &requestCollection})
	return
}

func GetFolder(ctx context.Context, params files_sdk.RequestGetFolderParams) (requestCollection files_sdk.RequestCollection, err error) {
	return (&Client{}).GetFolder(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.RequestCreateParams) (request files_sdk.Request, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/requests", Params: params, Entity: &request})
	return
}

func Create(ctx context.Context, params files_sdk.RequestCreateParams) (request files_sdk.Request, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.RequestDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/requests/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.RequestDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
