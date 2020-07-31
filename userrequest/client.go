package user_request

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

func (i *Iter) UserRequest() files_sdk.UserRequest {
	return i.Current().(files_sdk.UserRequest)
}

func (c *Client) List(params files_sdk.UserRequestListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/user_requests"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
        if err != nil {
          return &defaultValue, "", err
        }
		list := files_sdk.UserRequestCollection{}
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

func List(params files_sdk.UserRequestListParams) *Iter {
  client := Client{}
  return client.List (params)
}

func (c *Client) Find (params files_sdk.UserRequestFindParams) (files_sdk.UserRequest, error) {
  userRequest := files_sdk.UserRequest{}
  	path := "/user_requests/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return userRequest, err
	}
	if err := userRequest.UnmarshalJSON(*data); err != nil {
	return userRequest, err
	}

	return  userRequest, nil
}

func Find (params files_sdk.UserRequestFindParams) (files_sdk.UserRequest, error) {
  client := Client{}
  return client.Find (params)
}

func (c *Client) Create (params files_sdk.UserRequestCreateParams) (files_sdk.UserRequest, error) {
  userRequest := files_sdk.UserRequest{}
	  path := "/user_requests"
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return userRequest, err
	}
	if err := userRequest.UnmarshalJSON(*data); err != nil {
	return userRequest, err
	}

	return  userRequest, nil
}

func Create (params files_sdk.UserRequestCreateParams) (files_sdk.UserRequest, error) {
  client := Client{}
  return client.Create (params)
}

func (c *Client) Delete (params files_sdk.UserRequestDeleteParams) (files_sdk.UserRequest, error) {
  userRequest := files_sdk.UserRequest{}
  	path := "/user_requests/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return userRequest, err
	}
	if err := userRequest.UnmarshalJSON(*data); err != nil {
	return userRequest, err
	}

	return  userRequest, nil
}

func Delete (params files_sdk.UserRequestDeleteParams) (files_sdk.UserRequest, error) {
  client := Client{}
  return client.Delete (params)
}
