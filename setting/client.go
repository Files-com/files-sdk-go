package setting

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

func (c *Client) Languages(opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/settings/languages", Params: lib.Interface(), Entity: nil}, opts...)
	return
}

func Languages(opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Languages(opts...)
}

func (i *Iter) Settings() files_sdk.Settings {
	return i.Current().(files_sdk.Settings)
}

func (c *Client) List(params files_sdk.SettingListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/settings", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.SettingsCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.SettingListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) GetDomain(params files_sdk.SettingGetDomainParams, opts ...files_sdk.RequestResponseOption) (settings files_sdk.Settings, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/settings/domain", Params: params, Entity: &settings}, opts...)
	return
}

func GetDomain(params files_sdk.SettingGetDomainParams, opts ...files_sdk.RequestResponseOption) (settings files_sdk.Settings, err error) {
	return (&Client{}).GetDomain(params, opts...)
}
