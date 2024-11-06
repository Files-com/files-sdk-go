package export

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

func (i *Iter) Export() files_sdk.Export {
	return i.Current().(files_sdk.Export)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.ExportFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.ExportListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/exports", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ExportCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.ExportListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.ExportFindParams, opts ...files_sdk.RequestResponseOption) (export files_sdk.Export, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/exports/{id}", Params: params, Entity: &export}, opts...)
	return
}

func Find(params files_sdk.ExportFindParams, opts ...files_sdk.RequestResponseOption) (export files_sdk.Export, err error) {
	return (&Client{}).Find(params, opts...)
}
