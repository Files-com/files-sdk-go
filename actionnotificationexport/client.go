package action_notification_export

import (
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(params files_sdk.ActionNotificationExportFindParams) (files_sdk.ActionNotificationExport, error) {
	actionNotificationExport := files_sdk.ActionNotificationExport{}
	if params.Id == 0 {
		return actionNotificationExport, lib.CreateError(params, "Id")
	}
	path := "/action_notification_exports/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return actionNotificationExport, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
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

func Find(params files_sdk.ActionNotificationExportFindParams) (files_sdk.ActionNotificationExport, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.ActionNotificationExportCreateParams) (files_sdk.ActionNotificationExport, error) {
	actionNotificationExport := files_sdk.ActionNotificationExport{}
	path := "/action_notification_exports"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return actionNotificationExport, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
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

func Create(params files_sdk.ActionNotificationExportCreateParams) (files_sdk.ActionNotificationExport, error) {
	return (&Client{}).Create(params)
}
