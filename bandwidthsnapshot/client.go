package bandwidth_snapshot

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

func (i *Iter) BandwidthSnapshot() files_sdk.BandwidthSnapshot {
	return i.Current().(files_sdk.BandwidthSnapshot)
}

func (c *Client) List(params files_sdk.BandwidthSnapshotListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/bandwidth_snapshots"
	i.ListParams = &params
	list := files_sdk.BandwidthSnapshotCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
	return i, nil
}

func List(params files_sdk.BandwidthSnapshotListParams) (*Iter, error) {
	return (&Client{}).List(params)
}
