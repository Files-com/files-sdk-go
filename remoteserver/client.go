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
	*lib.Iter
}

func (i *Iter) RemoteServer() files_sdk.RemoteServer {
	return i.Current().(files_sdk.RemoteServer)
}

func (c *Client) List(ctx context.Context, params files_sdk.RemoteServerListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/remote_servers", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.RemoteServerCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.RemoteServerListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.RemoteServerFindParams) (remoteServer files_sdk.RemoteServer, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/remote_servers/{id}", Params: params, Entity: &remoteServer})
	return
}

func Find(ctx context.Context, params files_sdk.RemoteServerFindParams) (remoteServer files_sdk.RemoteServer, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) FindConfigurationFile(ctx context.Context, params files_sdk.RemoteServerFindConfigurationFileParams) (remoteServerConfigurationFile files_sdk.RemoteServerConfigurationFile, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/remote_servers/{id}/configuration_file", Params: params, Entity: &remoteServerConfigurationFile})
	return
}

func FindConfigurationFile(ctx context.Context, params files_sdk.RemoteServerFindConfigurationFileParams) (remoteServerConfigurationFile files_sdk.RemoteServerConfigurationFile, err error) {
	return (&Client{}).FindConfigurationFile(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.RemoteServerCreateParams) (remoteServer files_sdk.RemoteServer, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/remote_servers", Params: params, Entity: &remoteServer})
	return
}

func Create(ctx context.Context, params files_sdk.RemoteServerCreateParams) (remoteServer files_sdk.RemoteServer, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) ConfigurationFile(ctx context.Context, params files_sdk.RemoteServerConfigurationFileParams) (remoteServerConfigurationFile files_sdk.RemoteServerConfigurationFile, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/remote_servers/{id}/configuration_file", Params: params, Entity: &remoteServerConfigurationFile})
	return
}

func ConfigurationFile(ctx context.Context, params files_sdk.RemoteServerConfigurationFileParams) (remoteServerConfigurationFile files_sdk.RemoteServerConfigurationFile, err error) {
	return (&Client{}).ConfigurationFile(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.RemoteServerUpdateParams) (remoteServer files_sdk.RemoteServer, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/remote_servers/{id}", Params: params, Entity: &remoteServer})
	return
}

func Update(ctx context.Context, params files_sdk.RemoteServerUpdateParams) (remoteServer files_sdk.RemoteServer, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.RemoteServerDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/remote_servers/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.RemoteServerDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
