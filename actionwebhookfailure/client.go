package action_webhook_failure

import (
	"context"
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Retry(ctx context.Context, params files_sdk.ActionWebhookFailureRetryParams) (files_sdk.ActionWebhookFailure, error) {
	actionWebhookFailure := files_sdk.ActionWebhookFailure{}
	if params.Id == 0 {
		return actionWebhookFailure, lib.CreateError(params, "Id")
	}
	path := "/action_webhook_failures/" + strconv.FormatInt(params.Id, 10) + "/retry"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return actionWebhookFailure, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return actionWebhookFailure, err
	}
	if res.StatusCode == 204 {
		return actionWebhookFailure, nil
	}
	if err := actionWebhookFailure.UnmarshalJSON(*data); err != nil {
		return actionWebhookFailure, err
	}

	return actionWebhookFailure, nil
}

func Retry(ctx context.Context, params files_sdk.ActionWebhookFailureRetryParams) (files_sdk.ActionWebhookFailure, error) {
	return (&Client{}).Retry(ctx, params)
}
