package history_export

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(params files_sdk.HistoryExportFindParams, opts ...files_sdk.RequestResponseOption) (historyExport files_sdk.HistoryExport, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/history_exports/{id}", Params: params, Entity: &historyExport}, opts...)
	return
}

func Find(params files_sdk.HistoryExportFindParams, opts ...files_sdk.RequestResponseOption) (historyExport files_sdk.HistoryExport, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.HistoryExportCreateParams, opts ...files_sdk.RequestResponseOption) (historyExport files_sdk.HistoryExport, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/history_exports", Params: params, Entity: &historyExport}, opts...)
	return
}

func Create(params files_sdk.HistoryExportCreateParams, opts ...files_sdk.RequestResponseOption) (historyExport files_sdk.HistoryExport, err error) {
	return (&Client{}).Create(params, opts...)
}
