// Package external_event provides the Files.com ExternalEvent API client.
package external_event

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the ExternalEvent API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.ExternalEvent, error] {
	return files_sdk.IterAll[files_sdk.ExternalEvent](i.Iter)
}

// ExternalEvent returns the current resource. Call it only after Next returns true.
func (i *Iter) ExternalEvent() files_sdk.ExternalEvent {
	return i.Current().(files_sdk.ExternalEvent)
}

// LoadResource fetches a resource by its ID. identifier must have type
// int64.
func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.ExternalEventFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

// List returns an iterator for GET /external_events.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) List(params files_sdk.ExternalEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/external_events", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ExternalEventCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.ExternalEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// Find calls GET /external_events/{id}.
func (c *Client) Find(params files_sdk.ExternalEventFindParams, opts ...files_sdk.RequestResponseOption) (externalEvent files_sdk.ExternalEvent, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/external_events/{id}", Params: params, Entity: &externalEvent}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.ExternalEventFindParams, opts ...files_sdk.RequestResponseOption) (externalEvent files_sdk.ExternalEvent, err error) {
	return (&Client{}).Find(params, opts...)
}

// Create calls POST /external_events.
func (c *Client) Create(params files_sdk.ExternalEventCreateParams, opts ...files_sdk.RequestResponseOption) (externalEvent files_sdk.ExternalEvent, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/external_events", Params: params, Entity: &externalEvent}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.ExternalEventCreateParams, opts ...files_sdk.RequestResponseOption) (externalEvent files_sdk.ExternalEvent, err error) {
	return (&Client{}).Create(params, opts...)
}
