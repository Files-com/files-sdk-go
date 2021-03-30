package history_export_result

import (
	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
	listquery "github.com/Files-com/files-sdk-go/listquery"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*lib.Iter
}

func (i *Iter) HistoryExportResult() files_sdk.HistoryExportResult {
	return i.Current().(files_sdk.HistoryExportResult)
}

func (c *Client) List(params files_sdk.HistoryExportResultListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/history_export_results"
	i.ListParams = &params
	list := files_sdk.HistoryExportResultCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.HistoryExportResultListParams) (*Iter, error) {
	return (&Client{}).List(params)
}
