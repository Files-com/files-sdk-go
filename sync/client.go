package sync

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
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

func (i *Iter) Sync() files_sdk.Sync {
	return i.Current().(files_sdk.Sync)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.SyncFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.SyncListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/syncs", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.SyncCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.SyncListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.SyncFindParams, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/syncs/{id}", Params: params, Entity: &sync}, opts...)
	return
}

func Find(params files_sdk.SyncFindParams, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.SyncCreateParams, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/syncs", Params: params, Entity: &sync}, opts...)
	return
}

func Create(params files_sdk.SyncCreateParams, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) CreateMigrateTo(opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/syncs/migrate_to_syncs", Entity: nil}, opts...)
	return
}

func CreateMigrateTo(opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).CreateMigrateTo(opts...)
}

func (c *Client) Update(params files_sdk.SyncUpdateParams, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/syncs/{id}", Params: params, Entity: &sync}, opts...)
	return
}

func Update(params files_sdk.SyncUpdateParams, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error) {
	return (&Client{}).Update(params, opts...)
}

func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/syncs/{id}", Params: params, Entity: &sync}, opts...)
	return
}

func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (sync files_sdk.Sync, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

func (c *Client) Delete(params files_sdk.SyncDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/syncs/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.SyncDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
