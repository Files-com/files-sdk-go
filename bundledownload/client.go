package bundle_download

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

func (i *Iter) BundleDownload() files_sdk.BundleDownload {
	return i.Current().(files_sdk.BundleDownload)
}

func (c *Client) List(ctx context.Context, params files_sdk.BundleDownloadListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/bundle_downloads"
	i.ListParams = &params
	list := files_sdk.BundleDownloadCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.BundleDownloadListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}
