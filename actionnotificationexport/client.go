// Package action_notification_export provides the Files.com ActionNotificationExport API client.
package action_notification_export

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Client calls the ActionNotificationExport API using its embedded Config.
type Client struct {
	files_sdk.Config
}

// Find calls GET /action_notification_exports/{id}.
func (c *Client) Find(params files_sdk.ActionNotificationExportFindParams, opts ...files_sdk.RequestResponseOption) (actionNotificationExport files_sdk.ActionNotificationExport, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/action_notification_exports/{id}", Params: params, Entity: &actionNotificationExport}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.ActionNotificationExportFindParams, opts ...files_sdk.RequestResponseOption) (actionNotificationExport files_sdk.ActionNotificationExport, err error) {
	return (&Client{}).Find(params, opts...)
}

// Create calls POST /action_notification_exports.
func (c *Client) Create(params files_sdk.ActionNotificationExportCreateParams, opts ...files_sdk.RequestResponseOption) (actionNotificationExport files_sdk.ActionNotificationExport, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/action_notification_exports", Params: params, Entity: &actionNotificationExport}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.ActionNotificationExportCreateParams, opts ...files_sdk.RequestResponseOption) (actionNotificationExport files_sdk.ActionNotificationExport, err error) {
	return (&Client{}).Create(params, opts...)
}
