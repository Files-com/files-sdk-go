package history_export

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  files_sdk "github.com/Files-com/files-sdk-go"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*lib.Iter
}

func (i *Iter) HistoryExport() files_sdk.HistoryExport {
	return i.Current().(files_sdk.HistoryExport)
}

func (c *Client) List(params files_sdk.HistoryExportListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/history_exports"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
        if err != nil {
          return &defaultValue, "", err
        }
		list := files_sdk.HistoryExportCollection{}
		if err := list.UnmarshalJSON(*data); err != nil {
          return &defaultValue, "", err
        }

		ret := make([]interface{}, len(list))
		for i, v := range list {
			ret[i] = v
		}
		cursor := res.Header.Get("X-Files-Cursor")
		return &ret, cursor, nil
	}
	i.ListParams = &params
	return i
}

func List(params files_sdk.HistoryExportListParams) *Iter {
  client := Client{}
  return client.List (params)
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

func (c *Client) Delete (params files_sdk.HistoryExportDeleteParams) (files_sdk.HistoryExport, error) {
  historyExport := files_sdk.HistoryExport{}
  	path := "/history_exports/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return historyExport, err
	}
	if err := historyExport.UnmarshalJSON(*data); err != nil {
	return historyExport, err
	}

	return  historyExport, nil
}

func Delete (params files_sdk.HistoryExportDeleteParams) (files_sdk.HistoryExport, error) {
  client := Client{}
  return client.Delete (params)
}
