// Package exavault_api_request_log provides the Files.com ExavaultApiRequestLog API client.
package exavault_api_request_log

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the ExavaultApiRequestLog API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.ExavaultApiRequestLog, error] {
	return files_sdk.IterAll[files_sdk.ExavaultApiRequestLog](i.Iter)
}

// ExavaultApiRequestLog returns the current resource. Call it only after Next returns true.
func (i *Iter) ExavaultApiRequestLog() files_sdk.ExavaultApiRequestLog {
	return i.Current().(files_sdk.ExavaultApiRequestLog)
}

// List returns an iterator for GET /exavault_api_request_logs.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) List(params files_sdk.ExavaultApiRequestLogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/exavault_api_request_logs", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ExavaultApiRequestLogCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.ExavaultApiRequestLogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}
