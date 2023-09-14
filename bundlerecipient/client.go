package bundle_recipient

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

func (i *Iter) BundleRecipient() files_sdk.BundleRecipient {
	return i.Current().(files_sdk.BundleRecipient)
}

func (c *Client) List(params files_sdk.BundleRecipientListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/bundle_recipients", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.BundleRecipientCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.BundleRecipientListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Create(params files_sdk.BundleRecipientCreateParams, opts ...files_sdk.RequestResponseOption) (bundleRecipient files_sdk.BundleRecipient, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/bundle_recipients", Params: params, Entity: &bundleRecipient}, opts...)
	return
}

func Create(params files_sdk.BundleRecipientCreateParams, opts ...files_sdk.RequestResponseOption) (bundleRecipient files_sdk.BundleRecipient, err error) {
	return (&Client{}).Create(params, opts...)
}
