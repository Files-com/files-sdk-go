// Package file_comment provides the Files.com FileComment API client.
package file_comment

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the FileComment API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.FileComment, error] {
	return files_sdk.IterAll[files_sdk.FileComment](i.Iter)
}

// FileComment returns the current resource. Call it only after Next returns true.
func (i *Iter) FileComment() files_sdk.FileComment {
	return i.Current().(files_sdk.FileComment)
}

// Iterate returns an iterator for the path in identifier, which must be a string.
func (i *Iter) Iterate(identifier interface{}, opts ...files_sdk.RequestResponseOption) (files_sdk.IterI, error) {
	params := files_sdk.FileCommentListForParams{}
	if path, ok := identifier.(string); ok {
		params.Path = path
	}
	return i.Client.ListFor(params, opts...)
}

// ListFor returns an iterator for GET /file_comments/files/{path}.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
//
// API operation: List File Comments by Path.
func (c *Client) ListFor(params files_sdk.FileCommentListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/file_comments/files/{path}", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.FileCommentCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// ListFor returns a listing iterator using the default configuration.
// See Client.ListFor for the operation and paging behavior.
func ListFor(params files_sdk.FileCommentListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).ListFor(params, opts...)
}

// Create calls POST /file_comments.
func (c *Client) Create(params files_sdk.FileCommentCreateParams, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/file_comments", Params: params, Entity: &fileComment}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.FileCommentCreateParams, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error) {
	return (&Client{}).Create(params, opts...)
}

// Update calls PATCH /file_comments/{id}.
func (c *Client) Update(params files_sdk.FileCommentUpdateParams, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/file_comments/{id}", Params: params, Entity: &fileComment}, opts...)
	return
}

// Update calls Client.Update using the default configuration.
func Update(params files_sdk.FileCommentUpdateParams, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error) {
	return (&Client{}).Update(params, opts...)
}

// UpdateWithMap calls PATCH /file_comments/{id} using API parameter names as map keys.
// Include any path parameters in the map. Unlike optional struct fields, explicit
// zero values in the map are included in the request.
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/file_comments/{id}", Params: params, Entity: &fileComment}, opts...)
	return
}

// UpdateWithMap calls Client.UpdateWithMap using the default configuration.
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (fileComment files_sdk.FileComment, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

// Delete calls DELETE /file_comments/{id}.
func (c *Client) Delete(params files_sdk.FileCommentDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/file_comments/{id}", Params: params, Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(params files_sdk.FileCommentDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
