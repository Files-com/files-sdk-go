package history

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
	listquery "github.com/Files-com/files-sdk-go/v2/listquery"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*files_sdk.Iter
	*Client
}

func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

func (i *Iter) Action() files_sdk.Action {
	return i.Current().(files_sdk.Action)
}

func (c *Client) ListForFile(params files_sdk.HistoryListForFileParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/history/files/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ActionCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func ListForFile(params files_sdk.HistoryListForFileParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).ListForFile(params, opts...)
}

func (c *Client) ListForFolder(params files_sdk.HistoryListForFolderParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/history/folders/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ActionCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func ListForFolder(params files_sdk.HistoryListForFolderParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).ListForFolder(params, opts...)
}

func (c *Client) ListForUser(params files_sdk.HistoryListForUserParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/history/users/{user_id}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ActionCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func ListForUser(params files_sdk.HistoryListForUserParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).ListForUser(params, opts...)
}

func (c *Client) ListLogins(params files_sdk.HistoryListLoginsParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/history/login", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ActionCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func ListLogins(params files_sdk.HistoryListLoginsParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).ListLogins(params, opts...)
}

func (c *Client) List(params files_sdk.HistoryListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/history", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ActionCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.HistoryListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}
