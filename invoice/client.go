package invoice

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

func (i *Iter) Invoice() files_sdk.Invoice {
	return i.Current().(files_sdk.Invoice)
}

func (c *Client) List(params files_sdk.InvoiceListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/invoices"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.InvoiceCollection{}
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

func List(params files_sdk.InvoiceListParams) *Iter {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.InvoiceFindParams) (files_sdk.AccountLineItem, error) {
	accountLineItem := files_sdk.AccountLineItem{}
	path := "/invoices/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
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

func Find(params files_sdk.InvoiceFindParams) (files_sdk.AccountLineItem, error) {
	return (&Client{}).Find(params)
}
