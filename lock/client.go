package lock

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

func (i *Iter) Lock() files_sdk.Lock {
	return i.Current().(files_sdk.Lock)
}

func (i *Iter) Iterate(identifier interface{}, opts ...files_sdk.RequestResponseOption) (files_sdk.IterI, error) {
	params := files_sdk.LockListForParams{}
	if path, ok := identifier.(string); ok {
		params.Path = path
	}
	return i.Client.ListFor(params, opts...)
}

func (c *Client) ListFor(params files_sdk.LockListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/locks/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.LockCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func ListFor(params files_sdk.LockListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).ListFor(params, opts...)
}

func (c *Client) Create(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (lock files_sdk.Lock, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/locks/{path}", Params: params, Entity: &lock}, opts...)
	return
}

func Create(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (lock files_sdk.Lock, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Delete(params files_sdk.LockDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/locks/{path}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.LockDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
