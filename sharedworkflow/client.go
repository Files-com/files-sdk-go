// Package shared_workflow is the Files.com API client package for the files_sdk.SharedWorkflow model.
package shared_workflow

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the SharedWorkflow API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.SharedWorkflow, error] {
	return files_sdk.IterAll[files_sdk.SharedWorkflow](i.Iter)
}

// SharedWorkflow returns the current resource. Call it only after Next returns true.
func (i *Iter) SharedWorkflow() files_sdk.SharedWorkflow {
	return i.Current().(files_sdk.SharedWorkflow)
}

// LoadResource fetches a resource by its ID. identifier must have type
// int64.
func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.SharedWorkflowFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

// List returns an iterator for GET /shared_workflows.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) List(params files_sdk.SharedWorkflowListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/shared_workflows", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.SharedWorkflowCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.SharedWorkflowListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// Find calls GET /shared_workflows/{id}.
func (c *Client) Find(params files_sdk.SharedWorkflowFindParams, opts ...files_sdk.RequestResponseOption) (sharedWorkflow files_sdk.SharedWorkflow, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/shared_workflows/{id}", Params: params, Entity: &sharedWorkflow}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.SharedWorkflowFindParams, opts ...files_sdk.RequestResponseOption) (sharedWorkflow files_sdk.SharedWorkflow, err error) {
	return (&Client{}).Find(params, opts...)
}

// FindRun calls GET /shared_workflows/{id}/runs/{run_id}.
//
// API operation: Show a Shared Workflow run started by the recipient.
func (c *Client) FindRun(params files_sdk.SharedWorkflowFindRunParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/shared_workflows/{id}/runs/{run_id}", Params: params, Entity: &automationRun}, opts...)
	return
}

// FindRun calls Client.FindRun using the default configuration.
func FindRun(params files_sdk.SharedWorkflowFindRunParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error) {
	return (&Client{}).FindRun(params, opts...)
}

// RunWorkflow calls POST /shared_workflows/{id}/run_workflow.
//
// API operation: Run a Shared Workflow.
func (c *Client) RunWorkflow(params files_sdk.SharedWorkflowRunWorkflowParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/shared_workflows/{id}/run_workflow", Params: params, Entity: &automationRun}, opts...)
	return
}

// RunWorkflow calls Client.RunWorkflow using the default configuration.
func RunWorkflow(params files_sdk.SharedWorkflowRunWorkflowParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error) {
	return (&Client{}).RunWorkflow(params, opts...)
}
