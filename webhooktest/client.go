package webhooktest

import (
	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.WebhookTestCreateParams) (files_sdk.WebhookTest, error) {
	webhookTest := files_sdk.WebhookTest{}
	path := "/webhook_tests"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return webhookTest, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
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

func Create(params files_sdk.WebhookTestCreateParams) (files_sdk.WebhookTest, error) {
	return (&Client{}).Create(params)
}
