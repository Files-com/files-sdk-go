package webhooktest

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(ctx context.Context, params files_sdk.WebhookTestCreateParams) (files_sdk.WebhookTest, error) {
	webhookTest := files_sdk.WebhookTest{}
	path := "/webhook_tests"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return webhookTest, err
	}
	if res.StatusCode == 204 {
		return webhookTest, nil
	}

	return webhookTest, webhookTest.UnmarshalJSON(*data)
}

func Create(ctx context.Context, params files_sdk.WebhookTestCreateParams) (files_sdk.WebhookTest, error) {
	return (&Client{}).Create(ctx, params)
}
