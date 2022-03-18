package remote_server

import (
	"context"
	"strconv"

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
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/remote_servers"
	i.ListParams = &params
	list := files_sdk.RemoteServerCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.RemoteServerListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.RemoteServerFindParams) (files_sdk.RemoteServer, error) {
	remoteServer := files_sdk.RemoteServer{}
	if params.Id == 0 {
		return remoteServer, lib.CreateError(params, "Id")
	}
	path := "/remote_servers/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return remoteServer, err
	}
	if res.StatusCode == 204 {
		return remoteServer, nil
	}
	if err := remoteServer.UnmarshalJSON(*data); err != nil {
		return remoteServer, err
	}

	return remoteServer, nil
}

func Find(ctx context.Context, params files_sdk.RemoteServerFindParams) (files_sdk.RemoteServer, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.RemoteServerCreateParams) (files_sdk.RemoteServer, error) {
	remoteServer := files_sdk.RemoteServer{}
	path := "/remote_servers"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return remoteServer, err
	}
	if res.StatusCode == 204 {
		return remoteServer, nil
	}
	if err := remoteServer.UnmarshalJSON(*data); err != nil {
		return remoteServer, err
	}

	return remoteServer, nil
}

func Create(ctx context.Context, params files_sdk.RemoteServerCreateParams) (files_sdk.RemoteServer, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.RemoteServerUpdateParams) (files_sdk.RemoteServer, error) {
	remoteServer := files_sdk.RemoteServer{}
	if params.Id == 0 {
		return remoteServer, lib.CreateError(params, "Id")
	}
	path := "/remote_servers/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "PATCH", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return remoteServer, err
	}
	if res.StatusCode == 204 {
		return remoteServer, nil
	}
	if err := remoteServer.UnmarshalJSON(*data); err != nil {
		return remoteServer, err
	}

	return remoteServer, nil
}

func Update(ctx context.Context, params files_sdk.RemoteServerUpdateParams) (files_sdk.RemoteServer, error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.RemoteServerDeleteParams) (files_sdk.RemoteServer, error) {
	remoteServer := files_sdk.RemoteServer{}
	if params.Id == 0 {
		return remoteServer, lib.CreateError(params, "Id")
	}
	path := "/remote_servers/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return remoteServer, err
	}
	if res.StatusCode == 204 {
		return remoteServer, nil
	}
	if err := remoteServer.UnmarshalJSON(*data); err != nil {
		return remoteServer, err
	}

	return remoteServer, nil
}

func Delete(ctx context.Context, params files_sdk.RemoteServerDeleteParams) (files_sdk.RemoteServer, error) {
	return (&Client{}).Delete(ctx, params)
}
