package webhooktest

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(ctx context.Context, params files_sdk.WebhookTestCreateParams) (webhookTest files_sdk.WebhookTest, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/webhook_tests", Params: params, Entity: &webhookTest})
	return
}

func Create(ctx context.Context, params files_sdk.WebhookTestCreateParams) (webhookTest files_sdk.WebhookTest, err error) {
	return (&Client{}).Create(ctx, params)
}
