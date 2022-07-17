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

func (i *Iter) History() files_sdk.History {
	return i.Current().(files_sdk.History)
}

func (c *Client) ListForFile(ctx context.Context, params files_sdk.HistoryListForFileParams) (actionCollection files_sdk.ActionCollection, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/history/files/{path}", Params: params, Entity: &actionCollection})
	return
}

func ListForFile(ctx context.Context, params files_sdk.HistoryListForFileParams) (actionCollection files_sdk.ActionCollection, err error) {
	return (&Client{}).ListForFile(ctx, params)
}

func (c *Client) ListForFolder(ctx context.Context, params files_sdk.HistoryListForFolderParams) (actionCollection files_sdk.ActionCollection, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/history/folders/{path}", Params: params, Entity: &actionCollection})
	return
}

func ListForFolder(ctx context.Context, params files_sdk.HistoryListForFolderParams) (actionCollection files_sdk.ActionCollection, err error) {
	return (&Client{}).ListForFolder(ctx, params)
}

func (c *Client) ListForUser(ctx context.Context, params files_sdk.HistoryListForUserParams) (actionCollection files_sdk.ActionCollection, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/history/users/{user_id}", Params: params, Entity: &actionCollection})
	return
}

func ListForUser(ctx context.Context, params files_sdk.HistoryListForUserParams) (actionCollection files_sdk.ActionCollection, err error) {
	return (&Client{}).ListForUser(ctx, params)
}

func (c *Client) ListLogins(ctx context.Context, params files_sdk.HistoryListLoginsParams) (actionCollection files_sdk.ActionCollection, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/history/login", Params: params, Entity: &actionCollection})
	return
}

func ListLogins(ctx context.Context, params files_sdk.HistoryListLoginsParams) (actionCollection files_sdk.ActionCollection, err error) {
	return (&Client{}).ListLogins(ctx, params)
}

func (c *Client) List(ctx context.Context, params files_sdk.HistoryListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/history", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.HistoryCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.HistoryListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}
