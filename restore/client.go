package restore

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
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

func (i *Iter) Restore() files_sdk.Restore {
	return i.Current().(files_sdk.Restore)
}

func (c *Client) List(params files_sdk.RestoreListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/restores", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.RestoreCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.RestoreListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Create(params files_sdk.RestoreCreateParams, opts ...files_sdk.RequestResponseOption) (restore files_sdk.Restore, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/restores", Params: params, Entity: &restore}, opts...)
	return
}

func Create(params files_sdk.RestoreCreateParams, opts ...files_sdk.RequestResponseOption) (restore files_sdk.Restore, err error) {
	return (&Client{}).Create(params, opts...)
}
