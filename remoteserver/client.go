package remote_server

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
	listquery "github.com/Files-com/files-sdk-go/v2/listquery"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*files_sdk.Iter
	*Client
}

func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

func (i *Iter) RemoteServer() files_sdk.RemoteServer {
	return i.Current().(files_sdk.RemoteServer)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.RemoteServerFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(context.Background(), params, opts...)
}

func (c *Client) List(ctx context.Context, params files_sdk.RemoteServerListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/remote_servers", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.RemoteServerCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.RemoteServerListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.RemoteServerFindParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/remote_servers/{id}", Params: params, Entity: &remoteServer}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.RemoteServerFindParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) FindConfigurationFile(ctx context.Context, params files_sdk.RemoteServerFindConfigurationFileParams, opts ...files_sdk.RequestResponseOption) (remoteServerConfigurationFile files_sdk.RemoteServerConfigurationFile, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/remote_servers/{id}/configuration_file", Params: params, Entity: &remoteServerConfigurationFile}, opts...)
	return
}

func FindConfigurationFile(ctx context.Context, params files_sdk.RemoteServerFindConfigurationFileParams, opts ...files_sdk.RequestResponseOption) (remoteServerConfigurationFile files_sdk.RemoteServerConfigurationFile, err error) {
	return (&Client{}).FindConfigurationFile(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.RemoteServerCreateParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/remote_servers", Params: params, Entity: &remoteServer}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.RemoteServerCreateParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) ConfigurationFile(ctx context.Context, params files_sdk.RemoteServerConfigurationFileParams, opts ...files_sdk.RequestResponseOption) (remoteServerConfigurationFile files_sdk.RemoteServerConfigurationFile, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/remote_servers/{id}/configuration_file", Params: params, Entity: &remoteServerConfigurationFile}, opts...)
	return
}

func ConfigurationFile(ctx context.Context, params files_sdk.RemoteServerConfigurationFileParams, opts ...files_sdk.RequestResponseOption) (remoteServerConfigurationFile files_sdk.RemoteServerConfigurationFile, err error) {
	return (&Client{}).ConfigurationFile(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.RemoteServerUpdateParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/remote_servers/{id}", Params: params, Entity: &remoteServer}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.RemoteServerUpdateParams, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/remote_servers/{id}", Params: params, Entity: &remoteServer}, opts...)
	return
}

func UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (remoteServer files_sdk.RemoteServer, err error) {
	return (&Client{}).UpdateWithMap(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.RemoteServerDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/remote_servers/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.RemoteServerDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
