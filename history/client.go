package history

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
	listquery "github.com/Files-com/files-sdk-go/v2/listquery"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*lib.Iter
}

func (i *Iter) Action() files_sdk.Action {
	return i.Current().(files_sdk.Action)
}

func (c *Client) ListForFile(ctx context.Context, params files_sdk.HistoryListForFileParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/history/files/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ActionCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func ListForFile(ctx context.Context, params files_sdk.HistoryListForFileParams) (*Iter, error) {
	return (&Client{}).ListForFile(ctx, params)
}

func (c *Client) ListForFolder(ctx context.Context, params files_sdk.HistoryListForFolderParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/history/folders/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ActionCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func ListForFolder(ctx context.Context, params files_sdk.HistoryListForFolderParams) (*Iter, error) {
	return (&Client{}).ListForFolder(ctx, params)
}

func (c *Client) ListForUser(ctx context.Context, params files_sdk.HistoryListForUserParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/history/users/{user_id}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ActionCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func ListForUser(ctx context.Context, params files_sdk.HistoryListForUserParams) (*Iter, error) {
	return (&Client{}).ListForUser(ctx, params)
}

func (c *Client) ListLogins(ctx context.Context, params files_sdk.HistoryListLoginsParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/history/login", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ActionCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func ListLogins(ctx context.Context, params files_sdk.HistoryListLoginsParams) (*Iter, error) {
	return (&Client{}).ListLogins(ctx, params)
}

func (c *Client) List(ctx context.Context, params files_sdk.HistoryListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/history", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ActionCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.HistoryListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}
