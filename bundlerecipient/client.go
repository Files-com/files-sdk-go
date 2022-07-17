package bundle_recipient

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

func (i *Iter) BundleRecipient() files_sdk.BundleRecipient {
	return i.Current().(files_sdk.BundleRecipient)
}

func (c *Client) List(ctx context.Context, params files_sdk.BundleRecipientListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/bundle_recipients", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.BundleRecipientCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.BundleRecipientListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.BundleRecipientCreateParams) (bundleRecipient files_sdk.BundleRecipient, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/bundle_recipients", Params: params, Entity: &bundleRecipient})
	return
}

func Create(ctx context.Context, params files_sdk.BundleRecipientCreateParams) (bundleRecipient files_sdk.BundleRecipient, err error) {
	return (&Client{}).Create(ctx, params)
}
