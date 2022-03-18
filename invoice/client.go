package invoice

import (
	"context"
	"strconv"

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

func (i *Iter) Invoice() files_sdk.Invoice {
	return i.Current().(files_sdk.Invoice)
}

func (c *Client) List(ctx context.Context, params files_sdk.InvoiceListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/invoices"
	i.ListParams = &params
	list := files_sdk.InvoiceCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.InvoiceListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.InvoiceFindParams) (files_sdk.AccountLineItem, error) {
	accountLineItem := files_sdk.AccountLineItem{}
	if params.Id == 0 {
		return accountLineItem, lib.CreateError(params, "Id")
	}
	path := "/invoices/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
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

func Find(ctx context.Context, params files_sdk.InvoiceFindParams) (files_sdk.AccountLineItem, error) {
	return (&Client{}).Find(ctx, params)
}
