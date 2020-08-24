package ip_address

import (
	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
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

func (c *Client) List(params files_sdk.IpAddressListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/ip_addresses"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.IpAddressCollection{}
		if err := list.UnmarshalJSON(*data); err != nil {
			return &defaultValue, "", err
		}

		ret := make([]interface{}, len(list))
		for i, v := range list {
			ret[i] = v
		}
		cursor := res.Header.Get("X-Files-Cursor")
		return &ret, cursor, nil
	}
	i.ListParams = &params
	return i
}

func List(params files_sdk.IpAddressListParams) *Iter {
	return (&Client{}).List(params)
}

func (c *Client) GetReserved(params files_sdk.IpAddressGetReservedParams) (files_sdk.PublicIpAddressCollection, error) {
	publicIpAddressCollection := files_sdk.PublicIpAddressCollection{}
	path := "/ip_addresses/reserved"
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
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
