package bundle_download

import (
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

func (i *Iter) BundleDownload() files_sdk.BundleDownload {
	return i.Current().(files_sdk.BundleDownload)
}

func (c *Client) List(params files_sdk.BundleDownloadListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/bundle_downloads"
	i.ListParams = &params
	list := files_sdk.BundleDownloadCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.BundleDownloadListParams) (*Iter, error) {
	return (&Client{}).List(params)
}
