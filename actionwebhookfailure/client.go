package action_webhook_failure

import (
	"context"
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Retry(ctx context.Context, params files_sdk.ActionWebhookFailureRetryParams) error {
	actionWebhookFailure := files_sdk.ActionWebhookFailure{}
	if params.Id == 0 {
		return lib.CreateError(params, "Id")
	}
	path := "/action_webhook_failures/" + strconv.FormatInt(params.Id, 10) + "/retry"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	if res.StatusCode == 204 {
		return nil
	}

	return actionWebhookFailure.UnmarshalJSON(*data)
}

func Retry(ctx context.Context, params files_sdk.ActionWebhookFailureRetryParams) error {
	return (&Client{}).Retry(ctx, params)
}
