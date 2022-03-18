package bundle_registration

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

func (i *Iter) BundleRegistration() files_sdk.BundleRegistration {
	return i.Current().(files_sdk.BundleRegistration)
}

func (c *Client) List(ctx context.Context, params files_sdk.BundleRegistrationListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/bundle_registrations"
	i.ListParams = &params
	list := files_sdk.BundleRegistrationCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.BundleRegistrationListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}
