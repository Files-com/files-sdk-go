package snapshot

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
	*files_sdk.Iter
	*Client
}

func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

func (i *Iter) Snapshot() files_sdk.Snapshot {
	return i.Current().(files_sdk.Snapshot)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.SnapshotFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(context.Background(), params, opts...)
}

func (c *Client) List(ctx context.Context, params files_sdk.SnapshotListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/snapshots", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.SnapshotCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.SnapshotListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.SnapshotFindParams, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/snapshots/{id}", Params: params, Entity: &snapshot}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.SnapshotFindParams, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/snapshots", Params: lib.Interface(), Entity: &snapshot}, opts...)
	return
}

func Create(ctx context.Context, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error) {
	return (&Client{}).Create(ctx, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.SnapshotUpdateParams, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/snapshots/{id}", Params: params, Entity: &snapshot}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.SnapshotUpdateParams, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/snapshots/{id}", Params: params, Entity: &snapshot}, opts...)
	return
}

func UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error) {
	return (&Client{}).UpdateWithMap(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.SnapshotDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/snapshots/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.SnapshotDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
