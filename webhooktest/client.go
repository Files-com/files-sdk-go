package webhooktest

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(ctx context.Context, params files_sdk.WebhookTestCreateParams) (files_sdk.WebhookTest, error) {
	webhookTest := files_sdk.WebhookTest{}
	path := "/webhook_tests"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return webhookTest, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return webhookTest, err
	}
	if res.StatusCode == 204 {
		return webhookTest, nil
	}
	if err := webhookTest.UnmarshalJSON(*data); err != nil {
		return webhookTest, err
	}

	return webhookTest, nil
}

func Create(ctx context.Context, params files_sdk.WebhookTestCreateParams) (files_sdk.WebhookTest, error) {
	return (&Client{}).Create(ctx, params)
}
