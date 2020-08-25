package payment

import (
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
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

func (c *Client) List(params files_sdk.PaymentListParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/payments"
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
	return i, nil
}

func List(params files_sdk.PaymentListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.PaymentFindParams) (files_sdk.AccountLineItem, error) {
	accountLineItem := files_sdk.AccountLineItem{}
	if params.Id == 0 {
		return accountLineItem, lib.CreateError(params, "Id")
	}
	path := "/payments/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return accountLineItem, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParms)
	if err != nil {
		return accountLineItem, err
	}
	if res.StatusCode == 204 {
		return accountLineItem, nil
	}
	if err := accountLineItem.UnmarshalJSON(*data); err != nil {
		return accountLineItem, err
	}

	return accountLineItem, nil
}

func Find(params files_sdk.PaymentFindParams) (files_sdk.AccountLineItem, error) {
	return (&Client{}).Find(params)
}
