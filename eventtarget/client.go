// Package event_target provides the Files.com EventTarget API client.
package event_target

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the EventTarget API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.EventTarget, error] {
	return files_sdk.IterAll[files_sdk.EventTarget](i.Iter)
}

// EventTarget returns the current resource. Call it only after Next returns true.
func (i *Iter) EventTarget() files_sdk.EventTarget {
	return i.Current().(files_sdk.EventTarget)
}

// LoadResource fetches a resource by its ID. identifier must have type
// int64.
func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.EventTargetFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

// List returns an iterator for GET /event_targets.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) List(params files_sdk.EventTargetListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/event_targets", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.EventTargetCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.EventTargetListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// Find calls GET /event_targets/{id}.
func (c *Client) Find(params files_sdk.EventTargetFindParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/event_targets/{id}", Params: params, Entity: &eventTarget}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.EventTargetFindParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	return (&Client{}).Find(params, opts...)
}

// Create calls POST /event_targets.
func (c *Client) Create(params files_sdk.EventTargetCreateParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/event_targets", Params: params, Entity: &eventTarget}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.EventTargetCreateParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	return (&Client{}).Create(params, opts...)
}

// Update calls PATCH /event_targets/{id}.
func (c *Client) Update(params files_sdk.EventTargetUpdateParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/event_targets/{id}", Params: params, Entity: &eventTarget}, opts...)
	return
}

// Update calls Client.Update using the default configuration.
func Update(params files_sdk.EventTargetUpdateParams, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	return (&Client{}).Update(params, opts...)
}

// UpdateWithMap calls PATCH /event_targets/{id} using API parameter names as map keys.
// Include any path parameters in the map. Unlike optional struct fields, explicit
// zero values in the map are included in the request.
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/event_targets/{id}", Params: params, Entity: &eventTarget}, opts...)
	return
}

// UpdateWithMap calls Client.UpdateWithMap using the default configuration.
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (eventTarget files_sdk.EventTarget, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

// Delete calls DELETE /event_targets/{id}.
func (c *Client) Delete(params files_sdk.EventTargetDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/event_targets/{id}", Params: params, Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(params files_sdk.EventTargetDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
