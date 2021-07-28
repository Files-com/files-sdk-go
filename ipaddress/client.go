package ip_address

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
	listquery "github.com/Files-com/files-sdk-go/listquery"
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
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/ip_addresses"
	i.ListParams = &params
	list := files_sdk.IpAddressCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.IpAddressListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) GetReserved(ctx context.Context, params files_sdk.IpAddressGetReservedParams) (files_sdk.PublicIpAddressCollection, error) {
	publicIpAddressCollection := files_sdk.PublicIpAddressCollection{}
	path := "/ip_addresses/reserved"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return publicIpAddressCollection, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return publicIpAddressCollection, err
	}
	if res.StatusCode == 204 {
		return publicIpAddressCollection, nil
	}
	if err := publicIpAddressCollection.UnmarshalJSON(*data); err != nil {
		return publicIpAddressCollection, err
	}

	return publicIpAddressCollection, nil
}

func GetReserved(ctx context.Context, params files_sdk.IpAddressGetReservedParams) (files_sdk.PublicIpAddressCollection, error) {
	return (&Client{}).GetReserved(ctx, params)
}
