// Package webhooktest provides the Files.com WebhookTest API client.
package webhooktest

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Client calls the WebhookTest API using its embedded Config.
type Client struct {
	files_sdk.Config
}

// Create calls POST /webhook_tests.
func (c *Client) Create(params files_sdk.WebhookTestCreateParams, opts ...files_sdk.RequestResponseOption) (webhookTest files_sdk.WebhookTest, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/webhook_tests", Params: params, Entity: &webhookTest}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.WebhookTestCreateParams, opts ...files_sdk.RequestResponseOption) (webhookTest files_sdk.WebhookTest, err error) {
	return (&Client{}).Create(params, opts...)
}
