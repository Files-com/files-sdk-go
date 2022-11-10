package ip_address

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
	*lib.Iter
}

func (i *Iter) IpAddress() files_sdk.IpAddress {
	return i.Current().(files_sdk.IpAddress)
}

func (c *Client) List(ctx context.Context, params files_sdk.IpAddressListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/ip_addresses", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.IpAddressCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.IpAddressListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (i *Iter) PublicIpAddress() files_sdk.PublicIpAddress {
	return i.Current().(files_sdk.PublicIpAddress)
}

func (c *Client) GetExavaultReserved(ctx context.Context, params files_sdk.IpAddressGetExavaultReservedParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/ip_addresses/exavault-reserved", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PublicIpAddressCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func GetExavaultReserved(ctx context.Context, params files_sdk.IpAddressGetExavaultReservedParams) (*Iter, error) {
	return (&Client{}).GetExavaultReserved(ctx, params)
}

func (c *Client) GetReserved(ctx context.Context, params files_sdk.IpAddressGetReservedParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/ip_addresses/reserved", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PublicIpAddressCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func GetReserved(ctx context.Context, params files_sdk.IpAddressGetReservedParams) (*Iter, error) {
	return (&Client{}).GetReserved(ctx, params)
}
