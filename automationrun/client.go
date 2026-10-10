// Package automation_run provides the Files.com AutomationRun API client.
package automation_run

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the AutomationRun API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.AutomationRun, error] {
	return files_sdk.IterAll[files_sdk.AutomationRun](i.Iter)
}

// AutomationRun returns the current resource. Call it only after Next returns true.
func (i *Iter) AutomationRun() files_sdk.AutomationRun {
	return i.Current().(files_sdk.AutomationRun)
}

// LoadResource fetches a resource by its ID. identifier must have type
// int64.
func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.AutomationRunFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

// List returns an iterator for GET /automation_runs.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) List(params files_sdk.AutomationRunListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/automation_runs", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.AutomationRunCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.AutomationRunListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// Find calls GET /automation_runs/{id}.
func (c *Client) Find(params files_sdk.AutomationRunFindParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/automation_runs/{id}", Params: params, Entity: &automationRun}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.AutomationRunFindParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error) {
	return (&Client{}).Find(params, opts...)
}

// FindNode calls GET /automation_runs/{id}/node.
func (c *Client) FindNode(params files_sdk.AutomationRunFindNodeParams, opts ...files_sdk.RequestResponseOption) (automationExecutionNode files_sdk.AutomationExecutionNode, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/automation_runs/{id}/node", Params: params, Entity: &automationExecutionNode}, opts...)
	return
}

// FindNode calls Client.FindNode using the default configuration.
func FindNode(params files_sdk.AutomationRunFindNodeParams, opts ...files_sdk.RequestResponseOption) (automationExecutionNode files_sdk.AutomationExecutionNode, err error) {
	return (&Client{}).FindNode(params, opts...)
}

// Cancel calls POST /automation_runs/{id}/cancel.
//
// API operation: Cancel Automation Run.
func (c *Client) Cancel(params files_sdk.AutomationRunCancelParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/automation_runs/{id}/cancel", Params: params, Entity: &automationRun}, opts...)
	return
}

// Cancel calls Client.Cancel using the default configuration.
func Cancel(params files_sdk.AutomationRunCancelParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error) {
	return (&Client{}).Cancel(params, opts...)
}

// Rerun calls POST /automation_runs/{id}/rerun.
//
// API operation: Re-run Automation from Node.
func (c *Client) Rerun(params files_sdk.AutomationRunRerunParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/automation_runs/{id}/rerun", Params: params, Entity: &automationRun}, opts...)
	return
}

// Rerun calls Client.Rerun using the default configuration.
func Rerun(params files_sdk.AutomationRunRerunParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error) {
	return (&Client{}).Rerun(params, opts...)
}
