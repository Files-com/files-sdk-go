package app

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

func (i *Iter) App() files_sdk.App {
	return i.Current().(files_sdk.App)
}

func (c *Client) List(ctx context.Context, params files_sdk.AppListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/apps", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.AppCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.AppListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}
