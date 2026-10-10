// Package expectation_incident provides the Files.com ExpectationIncident API client.
package expectation_incident

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the ExpectationIncident API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.ExpectationIncident, error] {
	return files_sdk.IterAll[files_sdk.ExpectationIncident](i.Iter)
}

// ExpectationIncident returns the current resource. Call it only after Next returns true.
func (i *Iter) ExpectationIncident() files_sdk.ExpectationIncident {
	return i.Current().(files_sdk.ExpectationIncident)
}

// LoadResource fetches a resource by its ID. identifier must have type
// int64.
func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.ExpectationIncidentFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

// List returns an iterator for GET /expectation_incidents.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) List(params files_sdk.ExpectationIncidentListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/expectation_incidents", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.ExpectationIncidentCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.ExpectationIncidentListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// Find calls GET /expectation_incidents/{id}.
func (c *Client) Find(params files_sdk.ExpectationIncidentFindParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/expectation_incidents/{id}", Params: params, Entity: &expectationIncident}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.ExpectationIncidentFindParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	return (&Client{}).Find(params, opts...)
}

// Resolve calls POST /expectation_incidents/{id}/resolve.
//
// API operation: Resolve an expectation incident.
func (c *Client) Resolve(params files_sdk.ExpectationIncidentResolveParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/expectation_incidents/{id}/resolve", Params: params, Entity: &expectationIncident}, opts...)
	return
}

// Resolve calls Client.Resolve using the default configuration.
func Resolve(params files_sdk.ExpectationIncidentResolveParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	return (&Client{}).Resolve(params, opts...)
}

// Snooze calls POST /expectation_incidents/{id}/snooze.
//
// API operation: Snooze an expectation incident until a specified time.
func (c *Client) Snooze(params files_sdk.ExpectationIncidentSnoozeParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/expectation_incidents/{id}/snooze", Params: params, Entity: &expectationIncident}, opts...)
	return
}

// Snooze calls Client.Snooze using the default configuration.
func Snooze(params files_sdk.ExpectationIncidentSnoozeParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	return (&Client{}).Snooze(params, opts...)
}

// Acknowledge calls POST /expectation_incidents/{id}/acknowledge.
//
// API operation: Acknowledge an expectation incident.
func (c *Client) Acknowledge(params files_sdk.ExpectationIncidentAcknowledgeParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/expectation_incidents/{id}/acknowledge", Params: params, Entity: &expectationIncident}, opts...)
	return
}

// Acknowledge calls Client.Acknowledge using the default configuration.
func Acknowledge(params files_sdk.ExpectationIncidentAcknowledgeParams, opts ...files_sdk.RequestResponseOption) (expectationIncident files_sdk.ExpectationIncident, err error) {
	return (&Client{}).Acknowledge(params, opts...)
}
