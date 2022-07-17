package history_export

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(ctx context.Context, params files_sdk.HistoryExportFindParams) (historyExport files_sdk.HistoryExport, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/history_exports/{id}", Params: params, Entity: &historyExport})
	return
}

func Find(ctx context.Context, params files_sdk.HistoryExportFindParams) (historyExport files_sdk.HistoryExport, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.HistoryExportCreateParams) (historyExport files_sdk.HistoryExport, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/history_exports", Params: params, Entity: &historyExport})
	return
}

func Create(ctx context.Context, params files_sdk.HistoryExportCreateParams) (historyExport files_sdk.HistoryExport, err error) {
	return (&Client{}).Create(ctx, params)
}
