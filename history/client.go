package history

import (
	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
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

func (c *Client) ListForFile(params files_sdk.HistoryListForFileParams) (files_sdk.ActionCollection, error) {
	actionCollection := files_sdk.ActionCollection{}
	path := "/history/files/" + lib.QueryEscape(params.Path) + ""
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
		return actionCollection, err
	}
	if res.StatusCode == 204 {
		return actionCollection, nil
	}
	if err := actionCollection.UnmarshalJSON(*data); err != nil {
		return actionCollection, err
	}

	return actionCollection, nil
}

func ListForFile(params files_sdk.HistoryListForFileParams) (files_sdk.ActionCollection, error) {
	return (&Client{}).ListForFile(params)
}

func (c *Client) ListForFolder(params files_sdk.HistoryListForFolderParams) (files_sdk.ActionCollection, error) {
	actionCollection := files_sdk.ActionCollection{}
	path := "/history/folders/" + lib.QueryEscape(params.Path) + ""
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
		return actionCollection, err
	}
	if res.StatusCode == 204 {
		return actionCollection, nil
	}
	if err := actionCollection.UnmarshalJSON(*data); err != nil {
		return actionCollection, err
	}

	return actionCollection, nil
}

func ListForFolder(params files_sdk.HistoryListForFolderParams) (files_sdk.ActionCollection, error) {
	return (&Client{}).ListForFolder(params)
}

func (c *Client) ListForUser(params files_sdk.HistoryListForUserParams) (files_sdk.ActionCollection, error) {
	actionCollection := files_sdk.ActionCollection{}
	path := "/history/users/{user_id}"
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
		return actionCollection, err
	}
	if res.StatusCode == 204 {
		return actionCollection, nil
	}
	if err := actionCollection.UnmarshalJSON(*data); err != nil {
		return actionCollection, err
	}

	return actionCollection, nil
}

func ListForUser(params files_sdk.HistoryListForUserParams) (files_sdk.ActionCollection, error) {
	return (&Client{}).ListForUser(params)
}

func (c *Client) ListLogins(params files_sdk.HistoryListLoginsParams) (files_sdk.ActionCollection, error) {
	actionCollection := files_sdk.ActionCollection{}
	path := "/history/login"
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
		return actionCollection, err
	}
	if res.StatusCode == 204 {
		return actionCollection, nil
	}
	if err := actionCollection.UnmarshalJSON(*data); err != nil {
		return actionCollection, err
	}

	return actionCollection, nil
}

func ListLogins(params files_sdk.HistoryListLoginsParams) (files_sdk.ActionCollection, error) {
	return (&Client{}).ListLogins(params)
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
	return (&Client{}).List(params)
}
