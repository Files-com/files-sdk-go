// Package request provides the Files.com Request API client.
package request

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the Request API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.Request, error] {
	return files_sdk.IterAll[files_sdk.Request](i.Iter)
}

// Request returns the current resource. Call it only after Next returns true.
func (i *Iter) Request() files_sdk.Request {
	return i.Current().(files_sdk.Request)
}

// List returns an iterator for GET /requests.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) List(params files_sdk.RequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/requests", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.RequestCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.RequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// GetFolder returns an iterator for GET /requests/folders/{path}.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) GetFolder(params files_sdk.RequestGetFolderParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/requests/folders/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.RequestCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// GetFolder returns a listing iterator using the default configuration.
// See Client.GetFolder for the operation and paging behavior.
func GetFolder(params files_sdk.RequestGetFolderParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).GetFolder(params, opts...)
}

// Create calls POST /requests.
func (c *Client) Create(params files_sdk.RequestCreateParams, opts ...files_sdk.RequestResponseOption) (request files_sdk.Request, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/requests", Params: params, Entity: &request}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.RequestCreateParams, opts ...files_sdk.RequestResponseOption) (request files_sdk.Request, err error) {
	return (&Client{}).Create(params, opts...)
}

// Delete calls DELETE /requests/{id}.
func (c *Client) Delete(params files_sdk.RequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/requests/{id}", Params: params, Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(params files_sdk.RequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
