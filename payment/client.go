package payment

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

func (i *Iter) Payment() files_sdk.Payment {
	return i.Current().(files_sdk.Payment)
}

func (c *Client) List(params files_sdk.PaymentListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/payments"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
        if err != nil {
          return &defaultValue, "", err
        }
		list := files_sdk.PaymentCollection{}
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

func List(params files_sdk.PaymentListParams) *Iter {
  client := Client{}
  return client.List (params)
}

func (c *Client) Find (params files_sdk.PaymentFindParams) (files_sdk.Payment, error) {
  payment := files_sdk.Payment{}
  	path := "/payments/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return payment, err
	}
	if err := payment.UnmarshalJSON(*data); err != nil {
	return payment, err
	}

	return  payment, nil
}

func Find (params files_sdk.PaymentFindParams) (files_sdk.Payment, error) {
  client := Client{}
  return client.Find (params)
}
