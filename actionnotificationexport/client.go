package action_notification_export

import (
	"context"
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(ctx context.Context, params files_sdk.ActionNotificationExportFindParams) (files_sdk.ActionNotificationExport, error) {
	actionNotificationExport := files_sdk.ActionNotificationExport{}
	if params.Id == 0 {
		return actionNotificationExport, lib.CreateError(params, "Id")
	}
	path := "/action_notification_exports/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return actionNotificationExport, err
	}
	if res.StatusCode == 204 {
		return actionNotificationExport, nil
	}
	if err := actionNotificationExport.UnmarshalJSON(*data); err != nil {
		return actionNotificationExport, err
	}

	return actionNotificationExport, nil
}

func Find(ctx context.Context, params files_sdk.ActionNotificationExportFindParams) (files_sdk.ActionNotificationExport, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.ActionNotificationExportCreateParams) (files_sdk.ActionNotificationExport, error) {
	actionNotificationExport := files_sdk.ActionNotificationExport{}
	path := "/action_notification_exports"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return actionNotificationExport, err
	}
	if res.StatusCode == 204 {
		return actionNotificationExport, nil
	}
	if err := actionNotificationExport.UnmarshalJSON(*data); err != nil {
		return actionNotificationExport, err
	}

	return actionNotificationExport, nil
}

func Create(ctx context.Context, params files_sdk.ActionNotificationExportCreateParams) (files_sdk.ActionNotificationExport, error) {
	return (&Client{}).Create(ctx, params)
}
