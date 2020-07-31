package api_key

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

func (i *Iter) ApiKey() files_sdk.ApiKey {
	return i.Current().(files_sdk.ApiKey)
}

func (c *Client) List(params files_sdk.ApiKeyListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/api_keys"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
        if err != nil {
          return &defaultValue, "", err
        }
		list := files_sdk.ApiKeyCollection{}
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

func List(params files_sdk.ApiKeyListParams) *Iter {
  client := Client{}
  return client.List (params)
}

func (c *Client) Find (params files_sdk.ApiKeyFindParams) (files_sdk.ApiKey, error) {
  apiKey := files_sdk.ApiKey{}
  	path := "/api_keys/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return apiKey, err
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
	return apiKey, err
	}

	return  apiKey, nil
}

func Find (params files_sdk.ApiKeyFindParams) (files_sdk.ApiKey, error) {
  client := Client{}
  return client.Find (params)
}

func (c *Client) FindCurrent () (files_sdk.ApiKey, error) {
  apiKey := files_sdk.ApiKey{}
	  path := "/api_key"
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(lib.Interface()))
	if err != nil {
	  return apiKey, err
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
	return apiKey, err
	}

	return  apiKey, nil
}

func FindCurrent () (files_sdk.ApiKey, error) {
  client := Client{}
  return client.FindCurrent ()
}

func (c *Client) Create (params files_sdk.ApiKeyCreateParams) (files_sdk.ApiKey, error) {
  apiKey := files_sdk.ApiKey{}
	  path := "/api_keys"
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return apiKey, err
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
	return apiKey, err
	}

	return  apiKey, nil
}

func Create (params files_sdk.ApiKeyCreateParams) (files_sdk.ApiKey, error) {
  client := Client{}
  return client.Create (params)
}

func (c *Client) Update (params files_sdk.ApiKeyUpdateParams) (files_sdk.ApiKey, error) {
  apiKey := files_sdk.ApiKey{}
  	path := "/api_keys/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return apiKey, err
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
	return apiKey, err
	}

	return  apiKey, nil
}

func Update (params files_sdk.ApiKeyUpdateParams) (files_sdk.ApiKey, error) {
  client := Client{}
  return client.Update (params)
}

func (c *Client) UpdateCurrent (params files_sdk.ApiKeyUpdateCurrentParams) (files_sdk.ApiKey, error) {
  apiKey := files_sdk.ApiKey{}
	  path := "/api_key"
	data, _, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return apiKey, err
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
	return apiKey, err
	}

	return  apiKey, nil
}

func UpdateCurrent (params files_sdk.ApiKeyUpdateCurrentParams) (files_sdk.ApiKey, error) {
  client := Client{}
  return client.UpdateCurrent (params)
}

func (c *Client) Delete (params files_sdk.ApiKeyDeleteParams) (files_sdk.ApiKey, error) {
  apiKey := files_sdk.ApiKey{}
  	path := "/api_keys/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return apiKey, err
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
	return apiKey, err
	}

	return  apiKey, nil
}

func Delete (params files_sdk.ApiKeyDeleteParams) (files_sdk.ApiKey, error) {
  client := Client{}
  return client.Delete (params)
}

func (c *Client) DeleteCurrent () (files_sdk.ApiKey, error) {
  apiKey := files_sdk.ApiKey{}
	  path := "/api_key"
	data, _, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(lib.Interface()))
	if err != nil {
	  return apiKey, err
	}
	if err := apiKey.UnmarshalJSON(*data); err != nil {
	return apiKey, err
	}

	return  apiKey, nil
}

func DeleteCurrent () (files_sdk.ApiKey, error) {
  client := Client{}
  return client.DeleteCurrent ()
}
