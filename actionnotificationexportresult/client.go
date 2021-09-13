package action_notification_export_result

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

func (i *Iter) ActionNotificationExportResult() files_sdk.ActionNotificationExportResult {
	return i.Current().(files_sdk.ActionNotificationExportResult)
}

func (c *Client) List(ctx context.Context, params files_sdk.ActionNotificationExportResultListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/action_notification_export_results"
	i.ListParams = &params
	list := files_sdk.ActionNotificationExportResultCollection{}
	i.Query = listquery.Build(ctx, i, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.ActionNotificationExportResultListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}
