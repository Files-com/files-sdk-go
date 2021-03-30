package usage_snapshot

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

func (i *Iter) UsageSnapshot() files_sdk.UsageSnapshot {
	return i.Current().(files_sdk.UsageSnapshot)
}

func (c *Client) List(params files_sdk.UsageSnapshotListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/usage_snapshots"
	i.ListParams = &params
	list := files_sdk.UsageSnapshotCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.UsageSnapshotListParams) (*Iter, error) {
	return (&Client{}).List(params)
}
