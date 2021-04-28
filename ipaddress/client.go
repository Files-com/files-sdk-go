package ip_address

import (
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

func (c *Client) List(params files_sdk.IpAddressListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/ip_addresses"
	i.ListParams = &params
	list := files_sdk.IpAddressCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.IpAddressListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) GetReserved(params files_sdk.IpAddressGetReservedParams) (files_sdk.PublicIpAddressCollection, error) {
	publicIpAddressCollection := files_sdk.PublicIpAddressCollection{}
	path := "/ip_addresses/reserved"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return publicIpAddressCollection, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
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

func GetReserved(params files_sdk.IpAddressGetReservedParams) (files_sdk.PublicIpAddressCollection, error) {
	return (&Client{}).GetReserved(params)
}
