package history

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

func (i *Iter) History() files_sdk.History {
	return i.Current().(files_sdk.History)
}

func (c *Client) ListForFile(params files_sdk.HistoryListForFileParams) (files_sdk.ActionCollection, error) {
	actionCollection := files_sdk.ActionCollection{}
	path := lib.BuildPath("/history/files/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return actionCollection, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
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
	path := lib.BuildPath("/history/folders/", params.Path)
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return actionCollection, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
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
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return actionCollection, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
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
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return actionCollection, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
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

func (c *Client) List(params files_sdk.HistoryListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/history"
	i.ListParams = &params
	list := files_sdk.HistoryCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.HistoryListParams) (*Iter, error) {
	return (&Client{}).List(params)
}
