package action_notification_export

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(ctx context.Context, params files_sdk.ActionNotificationExportFindParams) (actionNotificationExport files_sdk.ActionNotificationExport, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/action_notification_exports/{id}", Params: params, Entity: &actionNotificationExport})
	return
}

func Find(ctx context.Context, params files_sdk.ActionNotificationExportFindParams) (actionNotificationExport files_sdk.ActionNotificationExport, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.ActionNotificationExportCreateParams) (actionNotificationExport files_sdk.ActionNotificationExport, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/action_notification_exports", Params: params, Entity: &actionNotificationExport})
	return
}

func Create(ctx context.Context, params files_sdk.ActionNotificationExportCreateParams) (actionNotificationExport files_sdk.ActionNotificationExport, err error) {
	return (&Client{}).Create(ctx, params)
}
