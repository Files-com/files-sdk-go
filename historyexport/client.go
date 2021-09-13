package history_export

import (
	"context"
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(ctx context.Context, params files_sdk.HistoryExportFindParams) (files_sdk.HistoryExport, error) {
	historyExport := files_sdk.HistoryExport{}
	if params.Id == 0 {
		return historyExport, lib.CreateError(params, "Id")
	}
	path := "/history_exports/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return historyExport, err
	}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return historyExport, err
	}
	if res.StatusCode == 204 {
		return historyExport, nil
	}
	if err := historyExport.UnmarshalJSON(*data); err != nil {
		return historyExport, err
	}

	return historyExport, nil
}

func Find(ctx context.Context, params files_sdk.HistoryExportFindParams) (files_sdk.HistoryExport, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.HistoryExportCreateParams) (files_sdk.HistoryExport, error) {
	historyExport := files_sdk.HistoryExport{}
	path := "/history_exports"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return historyExport, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return historyExport, err
	}
	if res.StatusCode == 204 {
		return historyExport, nil
	}
	if err := historyExport.UnmarshalJSON(*data); err != nil {
		return historyExport, err
	}

	return historyExport, nil
}

func Create(ctx context.Context, params files_sdk.HistoryExportCreateParams) (files_sdk.HistoryExport, error) {
	return (&Client{}).Create(ctx, params)
}
