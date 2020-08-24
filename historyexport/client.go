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
	path := "/history_exports/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
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
	data, res, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
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
