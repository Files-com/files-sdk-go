package ip_address

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

func (i *Iter) IpAddress() files_sdk.IpAddress {
	return i.Current().(files_sdk.IpAddress)
}

func (c *Client) List(params files_sdk.IpAddressListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/ip_addresses", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.IpAddressCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.IpAddressListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (i *Iter) PublicIpAddress() files_sdk.PublicIpAddress {
	return i.Current().(files_sdk.PublicIpAddress)
}

func (c *Client) GetActive(params files_sdk.IpAddressGetActiveParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/ip_addresses/active", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PublicIpAddressCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func GetActive(params files_sdk.IpAddressGetActiveParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).GetActive(params, opts...)
}

func (c *Client) GetExavaultReserved(params files_sdk.IpAddressGetExavaultReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/ip_addresses/exavault-reserved", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PublicIpAddressCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func GetExavaultReserved(params files_sdk.IpAddressGetExavaultReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).GetExavaultReserved(params, opts...)
}

func (c *Client) GetReserved(params files_sdk.IpAddressGetReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/ip_addresses/reserved", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PublicIpAddressCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func GetReserved(params files_sdk.IpAddressGetReservedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).GetReserved(params, opts...)
}
