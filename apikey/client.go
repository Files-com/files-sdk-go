// Package api_key provides the Files.com ApiKey API client.
package api_key

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the ApiKey API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.ApiKey, error] {
	return files_sdk.IterAll[files_sdk.ApiKey](i.Iter)
}

// ApiKey returns the current resource. Call it only after Next returns true.
func (i *Iter) ApiKey() files_sdk.ApiKey {
	return i.Current().(files_sdk.ApiKey)
}

// LoadResource fetches a resource by its ID. identifier must have type
// int64.
func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.ApiKeyFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

// List returns an iterator for GET /api_keys.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) List(params files_sdk.ApiKeyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/api_keys", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ApiKeyCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.ApiKeyListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// FindCurrent calls GET /api_key.
//
// API operation: Show information about current API key.  (Requires current API connection to be using an API key.).
func (c *Client) FindCurrent(opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/api_key", Entity: &apiKey}, opts...)
	return
}

// FindCurrent calls Client.FindCurrent using the default configuration.
func FindCurrent(opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).FindCurrent(opts...)
}

// Find calls GET /api_keys/{id}.
func (c *Client) Find(params files_sdk.ApiKeyFindParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/api_keys/{id}", Params: params, Entity: &apiKey}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.ApiKeyFindParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).Find(params, opts...)
}

// Create calls POST /api_keys.
func (c *Client) Create(params files_sdk.ApiKeyCreateParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/api_keys", Params: params, Entity: &apiKey}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.ApiKeyCreateParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).Create(params, opts...)
}

// UpdateCurrent calls PATCH /api_key.
//
// API operation: Update current API key.  (Requires current API connection to be using an API key.).
func (c *Client) UpdateCurrent(params files_sdk.ApiKeyUpdateCurrentParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/api_key", Params: params, Entity: &apiKey}, opts...)
	return
}

// UpdateCurrent calls Client.UpdateCurrent using the default configuration.
func UpdateCurrent(params files_sdk.ApiKeyUpdateCurrentParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).UpdateCurrent(params, opts...)
}

// UpdateCurrentWithMap calls PATCH /api_key using API parameter names as map keys.
// Include any path parameters in the map. Unlike optional struct fields, explicit
// zero values in the map are included in the request.
func (c *Client) UpdateCurrentWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/api_key", Params: params, Entity: &apiKey}, opts...)
	return
}

// UpdateCurrentWithMap calls Client.UpdateCurrentWithMap using the default configuration.
func UpdateCurrentWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).UpdateCurrentWithMap(params, opts...)
}

// Update calls PATCH /api_keys/{id}.
func (c *Client) Update(params files_sdk.ApiKeyUpdateParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/api_keys/{id}", Params: params, Entity: &apiKey}, opts...)
	return
}

// Update calls Client.Update using the default configuration.
func Update(params files_sdk.ApiKeyUpdateParams, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).Update(params, opts...)
}

// UpdateWithMap calls PATCH /api_keys/{id} using API parameter names as map keys.
// Include any path parameters in the map. Unlike optional struct fields, explicit
// zero values in the map are included in the request.
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/api_keys/{id}", Params: params, Entity: &apiKey}, opts...)
	return
}

// UpdateWithMap calls Client.UpdateWithMap using the default configuration.
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (apiKey files_sdk.ApiKey, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

// DeleteCurrent calls DELETE /api_key.
//
// API operation: Delete current API key.  (Requires current API connection to be using an API key.).
func (c *Client) DeleteCurrent(opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/api_key", Entity: nil}, opts...)
	return
}

// DeleteCurrent calls Client.DeleteCurrent using the default configuration.
func DeleteCurrent(opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).DeleteCurrent(opts...)
}

// Delete calls DELETE /api_keys/{id}.
func (c *Client) Delete(params files_sdk.ApiKeyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/api_keys/{id}", Params: params, Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(params files_sdk.ApiKeyDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
