package sso_strategy

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  files_sdk "github.com/Files-com/files-sdk-go"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*lib.Iter
}

func (i *Iter) SsoStrategy() files_sdk.SsoStrategy {
	return i.Current().(files_sdk.SsoStrategy)
}

func (c *Client) List(params files_sdk.SsoStrategyListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/sso_strategies"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
        if err != nil {
          return &defaultValue, "", err
        }
		list := files_sdk.SsoStrategyCollection{}
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

func List(params files_sdk.SsoStrategyListParams) *Iter {
  client := Client{}
  return client.List (params)
}

func (c *Client) Find (params files_sdk.SsoStrategyFindParams) (files_sdk.SsoStrategy, error) {
  ssoStrategy := files_sdk.SsoStrategy{}
  	path := "/sso_strategies/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return ssoStrategy, err
	}
	if err := ssoStrategy.UnmarshalJSON(*data); err != nil {
	return ssoStrategy, err
	}

	return  ssoStrategy, nil
}

func Find (params files_sdk.SsoStrategyFindParams) (files_sdk.SsoStrategy, error) {
  client := Client{}
  return client.Find (params)
}
