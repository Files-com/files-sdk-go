// Package folder provides the Files.com Folder API client.
package folder

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the Folder API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.File, error] {
	return files_sdk.IterAll[files_sdk.File](i.Iter)
}

// File returns the current resource. Call it only after Next returns true.
func (i *Iter) File() files_sdk.File {
	return i.Current().(files_sdk.File)
}

// Iterate returns an iterator for the path in identifier, which must be a string.
func (i *Iter) Iterate(identifier interface{}, opts ...files_sdk.RequestResponseOption) (files_sdk.IterI, error) {
	params := files_sdk.FolderListForParams{}
	if path, ok := identifier.(string); ok {
		params.Path = path
	}
	return i.Client.ListFor(params, opts...)
}

// ListFor returns an iterator for GET /folders/{path}.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
//
// API operation: List Folders by Path.
func (c *Client) ListFor(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/folders/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.FileCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// ListFor returns a listing iterator using the default configuration.
// See Client.ListFor for the operation and paging behavior.
func ListFor(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).ListFor(params, opts...)
}

// Create calls POST /folders/{path}.
func (c *Client) Create(params files_sdk.FolderCreateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/folders/{path}", Params: params, Entity: &file}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.FolderCreateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	return (&Client{}).Create(params, opts...)
}
