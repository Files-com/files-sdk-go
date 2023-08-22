package preview

import (
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

func (i *Iter) Preview() files_sdk.Preview {
	return i.Current().(files_sdk.Preview)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.PreviewFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.PreviewListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/previews", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PreviewCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.PreviewListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.PreviewFindParams, opts ...files_sdk.RequestResponseOption) (preview files_sdk.Preview, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/previews/{id}", Params: params, Entity: &preview}, opts...)
	return
}

func Find(params files_sdk.PreviewFindParams, opts ...files_sdk.RequestResponseOption) (preview files_sdk.Preview, err error) {
	return (&Client{}).Find(params, opts...)
}
