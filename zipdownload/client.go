package zip_download

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.ZipDownloadCreateParams, opts ...files_sdk.RequestResponseOption) (zipDownload files_sdk.ZipDownload, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/zip_downloads", Params: params, Entity: &zipDownload}, opts...)
	return
}

func Create(params files_sdk.ZipDownloadCreateParams, opts ...files_sdk.RequestResponseOption) (zipDownload files_sdk.ZipDownload, err error) {
	return (&Client{}).Create(params, opts...)
}
