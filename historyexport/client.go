package history_export

import (
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Find(params files_sdk.HistoryExportFindParams) (files_sdk.HistoryExport, error) {
	historyExport := files_sdk.HistoryExport{}
	if params.Id == 0 {
		return historyExport, lib.CreateError(params, "Id")
	}
	path := "/history_exports/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return historyExport, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParms)
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

func Find(params files_sdk.HistoryExportFindParams) (files_sdk.HistoryExport, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.HistoryExportCreateParams) (files_sdk.HistoryExport, error) {
	historyExport := files_sdk.HistoryExport{}
	path := "/history_exports"
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return historyExport, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
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

func Create(params files_sdk.HistoryExportCreateParams) (files_sdk.HistoryExport, error) {
	return (&Client{}).Create(params)
}
