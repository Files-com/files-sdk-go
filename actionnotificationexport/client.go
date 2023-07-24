package action_notification_export

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(params files_sdk.ActionNotificationExportFindParams, opts ...files_sdk.RequestResponseOption) (actionNotificationExport files_sdk.ActionNotificationExport, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/action_notification_exports/{id}", Params: params, Entity: &actionNotificationExport}, opts...)
	return
}

func Find(params files_sdk.ActionNotificationExportFindParams, opts ...files_sdk.RequestResponseOption) (actionNotificationExport files_sdk.ActionNotificationExport, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.ActionNotificationExportCreateParams, opts ...files_sdk.RequestResponseOption) (actionNotificationExport files_sdk.ActionNotificationExport, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/action_notification_exports", Params: params, Entity: &actionNotificationExport}, opts...)
	return
}

func Create(params files_sdk.ActionNotificationExportCreateParams, opts ...files_sdk.RequestResponseOption) (actionNotificationExport files_sdk.ActionNotificationExport, err error) {
	return (&Client{}).Create(params, opts...)
}
