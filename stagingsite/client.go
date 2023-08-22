package staging_site

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

func (i *Iter) Site() files_sdk.Site {
	return i.Current().(files_sdk.Site)
}

func (c *Client) List(params files_sdk.StagingSiteListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/staging_sites", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.SiteCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.StagingSiteListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Create(params files_sdk.StagingSiteCreateParams, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/staging_sites", Params: params, Entity: &site}, opts...)
	return
}

func Create(params files_sdk.StagingSiteCreateParams, opts ...files_sdk.RequestResponseOption) (site files_sdk.Site, err error) {
	return (&Client{}).Create(params, opts...)
}
