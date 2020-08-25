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

func (c *Client) List(params files_sdk.IpAddressListParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/ip_addresses"
	i.ListParams = &params
	exportParams, err := i.ExportParams()
	if err != nil {
		return i, err
	}
	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, exportParams)
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
	return i, nil
}

func List(params files_sdk.IpAddressListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) GetReserved(params files_sdk.IpAddressGetReservedParams) (files_sdk.PublicIpAddressCollection, error) {
	publicIpAddressCollection := files_sdk.PublicIpAddressCollection{}
	path := "/ip_addresses/reserved"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return publicIpAddressCollection, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParms)
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
