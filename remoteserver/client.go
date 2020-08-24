package remote_server

import (
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
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

func (c *Client) List(params files_sdk.RemoteServerListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/remote_servers"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
		if err != nil {
			return &defaultValue, "", err
		}
		list := files_sdk.RemoteServerCollection{}
		if err := list.UnmarshalJSON(*data); err != nil {
			return &defaultValue, "", err
		}

		ret := make([]interface{}, len(list))
		for i, v := range list {
			ret[i] = v
		}
		cursor := res.Header.Get("X-Files-Cursor")
		return &ret, cursor, nil
	}
	i.ListParams = &params
	return i
}

func List(params files_sdk.RemoteServerListParams) *Iter {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.RemoteServerFindParams) (files_sdk.RemoteServer, error) {
	remoteServer := files_sdk.RemoteServer{}
	path := "/remote_servers/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
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

func Find(params files_sdk.RemoteServerFindParams) (files_sdk.RemoteServer, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.RemoteServerCreateParams) (files_sdk.RemoteServer, error) {
	remoteServer := files_sdk.RemoteServer{}
	path := "/remote_servers"
	data, res, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
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

func Create(params files_sdk.RemoteServerCreateParams) (files_sdk.RemoteServer, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.RemoteServerUpdateParams) (files_sdk.RemoteServer, error) {
	remoteServer := files_sdk.RemoteServer{}
	path := "/remote_servers/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
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

func Update(params files_sdk.RemoteServerUpdateParams) (files_sdk.RemoteServer, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.RemoteServerDeleteParams) (files_sdk.RemoteServer, error) {
	remoteServer := files_sdk.RemoteServer{}
	path := "/remote_servers/" + lib.QueryEscape(strconv.FormatInt(params.Id, 10)) + ""
	data, res, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
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

func Delete(params files_sdk.RemoteServerDeleteParams) (files_sdk.RemoteServer, error) {
	return (&Client{}).Delete(params)
}
