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

func (c *Client) GetReserved(ctx context.Context, params files_sdk.IpAddressGetReservedParams) (publicIpAddressCollection files_sdk.PublicIpAddressCollection, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/ip_addresses/reserved", Params: params, Entity: &publicIpAddressCollection})
	return
}

func GetReserved(ctx context.Context, params files_sdk.IpAddressGetReservedParams) (publicIpAddressCollection files_sdk.PublicIpAddressCollection, err error) {
	return (&Client{}).GetReserved(ctx, params)
}
