// Package automation provides the Files.com Automation API client.
package automation

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the Automation API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.Automation, error] {
	return files_sdk.IterAll[files_sdk.Automation](i.Iter)
}

// Automation returns the current resource. Call it only after Next returns true.
func (i *Iter) Automation() files_sdk.Automation {
	return i.Current().(files_sdk.Automation)
}

// LoadResource fetches a resource by its ID. identifier must have type
// int64.
func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.AutomationFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

// List returns an iterator for GET /automations.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) List(params files_sdk.AutomationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/automations", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.AutomationCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.AutomationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// Find calls GET /automations/{id}.
func (c *Client) Find(params files_sdk.AutomationFindParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/automations/{id}", Params: params, Entity: &automation}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.AutomationFindParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error) {
	return (&Client{}).Find(params, opts...)
}

// GetAuthoringSchema calls GET /automations/authoring_schema.
//
// API operation: Show the Automation v2 authoring schema and active node catalog.
func (c *Client) GetAuthoringSchema(opts ...files_sdk.RequestResponseOption) (automationAuthoringSchema files_sdk.AutomationAuthoringSchema, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/automations/authoring_schema", Entity: &automationAuthoringSchema}, opts...)
	return
}

// GetAuthoringSchema calls Client.GetAuthoringSchema using the default configuration.
func GetAuthoringSchema(opts ...files_sdk.RequestResponseOption) (automationAuthoringSchema files_sdk.AutomationAuthoringSchema, err error) {
	return (&Client{}).GetAuthoringSchema(opts...)
}

// Create calls POST /automations.
func (c *Client) Create(params files_sdk.AutomationCreateParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/automations", Params: params, Entity: &automation}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.AutomationCreateParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error) {
	return (&Client{}).Create(params, opts...)
}

// Upgrade calls POST /automations/{id}/upgrade.
//
// API operation: Upgrade a legacy Automation to Automation v2.
func (c *Client) Upgrade(params files_sdk.AutomationUpgradeParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/automations/{id}/upgrade", Params: params, Entity: &automation}, opts...)
	return
}

// Upgrade calls Client.Upgrade using the default configuration.
func Upgrade(params files_sdk.AutomationUpgradeParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error) {
	return (&Client{}).Upgrade(params, opts...)
}

// ManualRun calls POST /automations/{id}/manual_run.
//
// API operation: Manually Run Automation.
func (c *Client) ManualRun(params files_sdk.AutomationManualRunParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/automations/{id}/manual_run", Params: params, Entity: nil}, opts...)
	return
}

// ManualRun calls Client.ManualRun using the default configuration.
func ManualRun(params files_sdk.AutomationManualRunParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).ManualRun(params, opts...)
}

// Update calls PATCH /automations/{id}.
func (c *Client) Update(params files_sdk.AutomationUpdateParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/automations/{id}", Params: params, Entity: &automation}, opts...)
	return
}

// Update calls Client.Update using the default configuration.
func Update(params files_sdk.AutomationUpdateParams, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error) {
	return (&Client{}).Update(params, opts...)
}

// UpdateWithMap calls PATCH /automations/{id} using API parameter names as map keys.
// Include any path parameters in the map. Unlike optional struct fields, explicit
// zero values in the map are included in the request.
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/automations/{id}", Params: params, Entity: &automation}, opts...)
	return
}

// UpdateWithMap calls Client.UpdateWithMap using the default configuration.
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (automation files_sdk.Automation, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

// Delete calls DELETE /automations/{id}.
func (c *Client) Delete(params files_sdk.AutomationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/automations/{id}", Params: params, Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(params files_sdk.AutomationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
