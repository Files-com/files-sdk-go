package ip_address

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  files_sdk "github.com/Files-com/files-sdk-go"
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
  client := Client{}
  return client.List (params)
}

func (c *Client) GetReserved (params files_sdk.IpAddressGetReservedParams) (files_sdk.IpAddress, error) {
  ipAddress := files_sdk.IpAddress{}
	  path := "/ip_addresses/reserved"
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return ipAddress, err
	}
	if err := ipAddress.UnmarshalJSON(*data); err != nil {
	return ipAddress, err
	}

	return  ipAddress, nil
}

func GetReserved (params files_sdk.IpAddressGetReservedParams) (files_sdk.IpAddress, error) {
  client := Client{}
  return client.GetReserved (params)
}
