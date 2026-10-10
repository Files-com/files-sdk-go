// Package sync_run provides the Files.com SyncRun API client.
package sync_run

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the SyncRun API using its embedded Config.
type Client struct {
	files_sdk.Config
}

// Iter traverses a paginated API response. Range over All to get each resource
// with a nil error, followed by the error that stopped the listing, if any.
// Alternatively, call Next before reading the current resource, then check Err
// after Next returns false.
type Iter struct {
	*files_sdk.Iter
	*Client
}

// Reload returns a new iterator for the same listing, starting at the first page.
// See files_sdk.Iter.Reload for parameter and request-option handling.
func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

// All returns an iterator over the listing's remaining resources, each with a
// nil error, followed by the error that stops the listing, if any. Pages are
// requested as the loop reaches them. See files_sdk.IterAll.
func (i *Iter) All() iter.Seq2[files_sdk.SyncRun, error] {
	return files_sdk.IterAll[files_sdk.SyncRun](i.Iter)
}

// SyncRun returns the current resource. Call it only after Next returns true.
func (i *Iter) SyncRun() files_sdk.SyncRun {
	return i.Current().(files_sdk.SyncRun)
}

// LoadResource fetches a resource by its ID. identifier must have type
// int64.
func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.SyncRunFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

// List returns an iterator for GET /sync_runs.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) List(params files_sdk.SyncRunListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/sync_runs", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.SyncRunCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.SyncRunListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// Find calls GET /sync_runs/{id}.
func (c *Client) Find(params files_sdk.SyncRunFindParams, opts ...files_sdk.RequestResponseOption) (syncRun files_sdk.SyncRun, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/sync_runs/{id}", Params: params, Entity: &syncRun}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.SyncRunFindParams, opts ...files_sdk.RequestResponseOption) (syncRun files_sdk.SyncRun, err error) {
	return (&Client{}).Find(params, opts...)
}
