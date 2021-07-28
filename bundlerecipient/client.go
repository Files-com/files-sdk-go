package bundle_recipient

import (
	"context"

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

func (i *Iter) BundleRecipient() files_sdk.BundleRecipient {
	return i.Current().(files_sdk.BundleRecipient)
}

func (c *Client) List(ctx context.Context, params files_sdk.BundleRecipientListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/bundle_recipients"
	i.ListParams = &params
	list := files_sdk.BundleRecipientCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.BundleRecipientListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.BundleRecipientCreateParams) (files_sdk.BundleRecipient, error) {
	bundleRecipient := files_sdk.BundleRecipient{}
	path := "/bundle_recipients"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return bundleRecipient, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return bundleRecipient, err
	}
	if res.StatusCode == 204 {
		return bundleRecipient, nil
	}
	if err := bundleRecipient.UnmarshalJSON(*data); err != nil {
		return bundleRecipient, err
	}

	return bundleRecipient, nil
}

func Create(ctx context.Context, params files_sdk.BundleRecipientCreateParams) (files_sdk.BundleRecipient, error) {
	return (&Client{}).Create(ctx, params)
}
