// Package remote_server provides the Files.com RemoteServer API client.
package remote_server

import (
	"iter"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

// Client calls the RemoteServer API using its embedded Config.
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
func (i *Iter) All() iter.Seq2[files_sdk.RemoteServer, error] {
	return files_sdk.IterAll[files_sdk.RemoteServer](i.Iter)
}

// RemoteServer returns the current resource. Call it only after Next returns true.
func (i *Iter) RemoteServer() files_sdk.RemoteServer {
	return i.Current().(files_sdk.RemoteServer)
}

// LoadResource fetches a resource by its ID. identifier must have type
// int64.
func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.RemoteServerFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

// List returns an iterator for GET /remote_servers.
// Pages are requested as iteration reaches them: range over Iter.All, or call
// Next and then check Err.
func (c *Client) List(params files_sdk.RemoteServerListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/remote_servers", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.RemoteServerCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

// List returns a listing iterator using the default configuration.
// See Client.List for the operation and paging behavior.
func List(params files_sdk.RemoteServerListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

// Find calls GET /remote_servers/{id}.
func (c *Client) Find(params files_sdk.RemoteServerFindParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/remote_servers/{id}", Params: params, Entity: &remoteServer}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.RemoteServerFindParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	return (&Client{}).Find(params, opts...)
}

// AgentNodes calls GET /remote_servers/{id}/agent_nodes.
//
// API operation: List Files.com Agent nodes.
func (c *Client) AgentNodes(params files_sdk.RemoteServerAgentNodesParams, opts ...files_sdk.RequestResponseOption) (agentNode files_sdk.AgentNode, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/remote_servers/{id}/agent_nodes", Params: params, Entity: &agentNode}, opts...)
	return
}

// AgentNodes calls Client.AgentNodes using the default configuration.
func AgentNodes(params files_sdk.RemoteServerAgentNodesParams, opts ...files_sdk.RequestResponseOption) (agentNode files_sdk.AgentNode, err error) {
	return (&Client{}).AgentNodes(params, opts...)
}

// FindConfigurationFile calls GET /remote_servers/{id}/configuration_file.
//
// API operation: Download configuration file (required for some Remote Server integrations, such as the Files.com Agent).
func (c *Client) FindConfigurationFile(params files_sdk.RemoteServerFindConfigurationFileParams, opts ...files_sdk.RequestResponseOption) (remoteServerConfigurationFile files_sdk.RemoteServerConfigurationFile, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/remote_servers/{id}/configuration_file", Params: params, Entity: &remoteServerConfigurationFile}, opts...)
	return
}

// FindConfigurationFile calls Client.FindConfigurationFile using the default configuration.
func FindConfigurationFile(params files_sdk.RemoteServerFindConfigurationFileParams, opts ...files_sdk.RequestResponseOption) (remoteServerConfigurationFile files_sdk.RemoteServerConfigurationFile, err error) {
	return (&Client{}).FindConfigurationFile(params, opts...)
}

// Create calls POST /remote_servers.
func (c *Client) Create(params files_sdk.RemoteServerCreateParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/remote_servers", Params: params, Entity: &remoteServer}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.RemoteServerCreateParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	return (&Client{}).Create(params, opts...)
}

// AgentPushUpdate calls POST /remote_servers/{id}/agent_push_update.
//
// API operation: Push update to Files Agent.
func (c *Client) AgentPushUpdate(params files_sdk.RemoteServerAgentPushUpdateParams, opts ...files_sdk.RequestResponseOption) (agentPushUpdate files_sdk.AgentPushUpdate, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/remote_servers/{id}/agent_push_update", Params: params, Entity: &agentPushUpdate}, opts...)
	return
}

// AgentPushUpdate calls Client.AgentPushUpdate using the default configuration.
func AgentPushUpdate(params files_sdk.RemoteServerAgentPushUpdateParams, opts ...files_sdk.RequestResponseOption) (agentPushUpdate files_sdk.AgentPushUpdate, err error) {
	return (&Client{}).AgentPushUpdate(params, opts...)
}

// AgentPushUpdateWithMap calls POST /remote_servers/{id}/agent_push_update using API parameter names as map keys.
// Include any path parameters in the map. Unlike optional struct fields, explicit
// zero values in the map are included in the request.
func (c *Client) AgentPushUpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (agentPushUpdate files_sdk.AgentPushUpdate, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/remote_servers/{id}/agent_push_update", Params: params, Entity: &agentPushUpdate}, opts...)
	return
}

// AgentPushUpdateWithMap calls Client.AgentPushUpdateWithMap using the default configuration.
func AgentPushUpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (agentPushUpdate files_sdk.AgentPushUpdate, err error) {
	return (&Client{}).AgentPushUpdateWithMap(params, opts...)
}

// Update calls PATCH /remote_servers/{id}.
func (c *Client) Update(params files_sdk.RemoteServerUpdateParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/remote_servers/{id}", Params: params, Entity: &remoteServer}, opts...)
	return
}

// Update calls Client.Update using the default configuration.
func Update(params files_sdk.RemoteServerUpdateParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	return (&Client{}).Update(params, opts...)
}

// UpdateWithMap calls PATCH /remote_servers/{id} using API parameter names as map keys.
// Include any path parameters in the map. Unlike optional struct fields, explicit
// zero values in the map are included in the request.
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/remote_servers/{id}", Params: params, Entity: &remoteServer}, opts...)
	return
}

// UpdateWithMap calls Client.UpdateWithMap using the default configuration.
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

// Delete calls DELETE /remote_servers/{id}.
func (c *Client) Delete(params files_sdk.RemoteServerDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/remote_servers/{id}", Params: params, Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(params files_sdk.RemoteServerDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
