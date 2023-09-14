package webhooktest

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.WebhookTestCreateParams, opts ...files_sdk.RequestResponseOption) (webhookTest files_sdk.WebhookTest, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/webhook_tests", Params: params, Entity: &webhookTest}, opts...)
	return
}

func Create(params files_sdk.WebhookTestCreateParams, opts ...files_sdk.RequestResponseOption) (webhookTest files_sdk.WebhookTest, err error) {
	return (&Client{}).Create(params, opts...)
}
