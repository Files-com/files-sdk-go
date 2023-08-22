package bundle_registration

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

func (i *Iter) BundleRegistration() files_sdk.BundleRegistration {
	return i.Current().(files_sdk.BundleRegistration)
}

func (c *Client) List(params files_sdk.BundleRegistrationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/bundle_registrations", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.BundleRegistrationCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.BundleRegistrationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Create(params files_sdk.BundleRegistrationCreateParams, opts ...files_sdk.RequestResponseOption) (bundleRegistration files_sdk.BundleRegistration, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/bundle_registrations", Params: params, Entity: &bundleRegistration}, opts...)
	return
}

func Create(params files_sdk.BundleRegistrationCreateParams, opts ...files_sdk.RequestResponseOption) (bundleRegistration files_sdk.BundleRegistration, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) LastActivity(params files_sdk.BundleRegistrationLastActivityParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/bundle_registrations/last_activity", Params: params, Entity: nil}, opts...)
	return
}

func LastActivity(params files_sdk.BundleRegistrationLastActivityParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).LastActivity(params, opts...)
}
