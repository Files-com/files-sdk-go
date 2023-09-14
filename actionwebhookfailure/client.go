package action_webhook_failure

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Retry(params files_sdk.ActionWebhookFailureRetryParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/action_webhook_failures/{id}/retry", Params: params, Entity: nil}, opts...)
	return
}

func Retry(params files_sdk.ActionWebhookFailureRetryParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Retry(params, opts...)
}
