package payment

import (
	"context"

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

func (i *Iter) Payment() files_sdk.Payment {
	return i.Current().(files_sdk.Payment)
}

func (c *Client) List(ctx context.Context, params files_sdk.PaymentListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/payments", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.PaymentCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.PaymentListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.PaymentFindParams) (accountLineItem files_sdk.AccountLineItem, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/payments/{id}", Params: params, Entity: &accountLineItem})
	return
}

func Find(ctx context.Context, params files_sdk.PaymentFindParams) (accountLineItem files_sdk.AccountLineItem, err error) {
	return (&Client{}).Find(ctx, params)
}
