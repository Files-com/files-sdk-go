package settings_change

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

func (i *Iter) SettingsChange() files_sdk.SettingsChange {
	return i.Current().(files_sdk.SettingsChange)
}

func (c *Client) List(ctx context.Context, params files_sdk.SettingsChangeListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/settings_changes"
	i.ListParams = &params
	list := files_sdk.SettingsChangeCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.SettingsChangeListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}
