package history_export

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  files_sdk "github.com/Files-com/files-sdk-go"
)

type Client struct {
	files_sdk.Config
}


func (c *Client) Find (params files_sdk.HistoryExportFindParams) (files_sdk.HistoryExport, error) {
  historyExport := files_sdk.HistoryExport{}
  	path := "/history_exports/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return historyExport, err
	}
	if err := historyExport.UnmarshalJSON(*data); err != nil {
	return historyExport, err
	}

	return  historyExport, nil
}

func Find (params files_sdk.HistoryExportFindParams) (files_sdk.HistoryExport, error) {
  client := Client{}
  return client.Find (params)
}

func (c *Client) Create (params files_sdk.HistoryExportCreateParams) (files_sdk.HistoryExport, error) {
  historyExport := files_sdk.HistoryExport{}
	  path := "/history_exports"
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return historyExport, err
	}
	if err := historyExport.UnmarshalJSON(*data); err != nil {
	return historyExport, err
	}

	return  historyExport, nil
}

func Create (params files_sdk.HistoryExportCreateParams) (files_sdk.HistoryExport, error) {
  client := Client{}
  return client.Create (params)
}
