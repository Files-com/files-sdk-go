package history

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

func (i *Iter) History() files_sdk.History {
	return i.Current().(files_sdk.History)
}

func (c *Client) ListForFile (params files_sdk.HistoryListForFileParams) (files_sdk.History, error) {
  history := files_sdk.History{}
		path := "/history/files/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return history, err
	}
	if err := history.UnmarshalJSON(*data); err != nil {
	return history, err
	}

	return  history, nil
}

func ListForFile (params files_sdk.HistoryListForFileParams) (files_sdk.History, error) {
  client := Client{}
  return client.ListForFile (params)
}

func (c *Client) ListForFolder (params files_sdk.HistoryListForFolderParams) (files_sdk.History, error) {
  history := files_sdk.History{}
		path := "/history/folders/" + lib.QueryEscape(params.Path) + ""
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return history, err
	}
	if err := history.UnmarshalJSON(*data); err != nil {
	return history, err
	}

	return  history, nil
}

func ListForFolder (params files_sdk.HistoryListForFolderParams) (files_sdk.History, error) {
  client := Client{}
  return client.ListForFolder (params)
}

func (c *Client) ListForUser (params files_sdk.HistoryListForUserParams) (files_sdk.History, error) {
  history := files_sdk.History{}
	  path := "/history/users/{user_id}"
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return history, err
	}
	if err := history.UnmarshalJSON(*data); err != nil {
	return history, err
	}

	return  history, nil
}

func ListForUser (params files_sdk.HistoryListForUserParams) (files_sdk.History, error) {
  client := Client{}
  return client.ListForUser (params)
}

func (c *Client) ListLogins (params files_sdk.HistoryListLoginsParams) (files_sdk.History, error) {
  history := files_sdk.History{}
	  path := "/history/login"
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return history, err
	}
	if err := history.UnmarshalJSON(*data); err != nil {
	return history, err
	}

	return  history, nil
}

func ListLogins (params files_sdk.HistoryListLoginsParams) (files_sdk.History, error) {
  client := Client{}
  return client.ListLogins (params)
}

func (c *Client) List(params files_sdk.HistoryListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/history"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
        if err != nil {
          return &defaultValue, "", err
        }
		list := files_sdk.HistoryCollection{}
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

func List(params files_sdk.HistoryListParams) *Iter {
  client := Client{}
  return client.List (params)
}
