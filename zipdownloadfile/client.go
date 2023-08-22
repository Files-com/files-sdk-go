package zip_download_file

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
	listquery "github.com/Files-com/files-sdk-go/v2/listquery"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*files_sdk.Iter
	*Client
}

func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

func (i *Iter) ZipDownloadFiles() files_sdk.ZipDownloadFiles {
	return i.Current().(files_sdk.ZipDownloadFiles)
}

func (c *Client) Create(params files_sdk.ZipDownloadFileCreateParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/zip_download_files", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ZipDownloadFilesCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func Create(params files_sdk.ZipDownloadFileCreateParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).Create(params, opts...)
}
