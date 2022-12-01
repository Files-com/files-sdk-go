package action_webhook_failure

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Retry(ctx context.Context, params files_sdk.ActionWebhookFailureRetryParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/action_webhook_failures/{id}/retry", Params: params, Entity: nil}, opts...)
	return
}

func Retry(ctx context.Context, params files_sdk.ActionWebhookFailureRetryParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Retry(ctx, params, opts...)
}
