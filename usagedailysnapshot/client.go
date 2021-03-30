package usage_daily_snapshot

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

func (i *Iter) UsageDailySnapshot() files_sdk.UsageDailySnapshot {
	return i.Current().(files_sdk.UsageDailySnapshot)
}

func (c *Client) List(params files_sdk.UsageDailySnapshotListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/usage_daily_snapshots"
	i.ListParams = &params
	list := files_sdk.UsageDailySnapshotCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.UsageDailySnapshotListParams) (*Iter, error) {
	return (&Client{}).List(params)
}
