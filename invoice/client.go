package invoice

import (
	"strconv"

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

func (i *Iter) Invoice() files_sdk.Invoice {
	return i.Current().(files_sdk.Invoice)
}

func (c *Client) List(params files_sdk.InvoiceListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/invoices"
	i.ListParams = &params
	list := files_sdk.InvoiceCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.InvoiceListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.InvoiceFindParams) (files_sdk.AccountLineItem, error) {
	accountLineItem := files_sdk.AccountLineItem{}
	if params.Id == 0 {
		return accountLineItem, lib.CreateError(params, "Id")
	}
	path := "/invoices/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return accountLineItem, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
