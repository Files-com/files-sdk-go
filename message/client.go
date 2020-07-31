package message

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

func (i *Iter) Message() files_sdk.Message {
	return i.Current().(files_sdk.Message)
}

func (c *Client) List(params files_sdk.MessageListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/messages"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
        if err != nil {
          return &defaultValue, "", err
        }
		list := files_sdk.MessageCollection{}
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

func List(params files_sdk.MessageListParams) *Iter {
  client := Client{}
  return client.List (params)
}

func (c *Client) Find (params files_sdk.MessageFindParams) (files_sdk.Message, error) {
  message := files_sdk.Message{}
  	path := "/messages/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return message, err
	}
	if err := message.UnmarshalJSON(*data); err != nil {
	return message, err
	}

	return  message, nil
}

func Find (params files_sdk.MessageFindParams) (files_sdk.Message, error) {
  client := Client{}
  return client.Find (params)
}

func (c *Client) Create (params files_sdk.MessageCreateParams) (files_sdk.Message, error) {
  message := files_sdk.Message{}
	  path := "/messages"
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return message, err
	}
	if err := message.UnmarshalJSON(*data); err != nil {
	return message, err
	}

	return  message, nil
}

func Create (params files_sdk.MessageCreateParams) (files_sdk.Message, error) {
  client := Client{}
  return client.Create (params)
}

func (c *Client) Update (params files_sdk.MessageUpdateParams) (files_sdk.Message, error) {
  message := files_sdk.Message{}
  	path := "/messages/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return message, err
	}
	if err := message.UnmarshalJSON(*data); err != nil {
	return message, err
	}

	return  message, nil
}

func Update (params files_sdk.MessageUpdateParams) (files_sdk.Message, error) {
  client := Client{}
  return client.Update (params)
}

func (c *Client) Delete (params files_sdk.MessageDeleteParams) (files_sdk.Message, error) {
  message := files_sdk.Message{}
  	path := "/messages/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return message, err
	}
	if err := message.UnmarshalJSON(*data); err != nil {
	return message, err
	}

	return  message, nil
}

func Delete (params files_sdk.MessageDeleteParams) (files_sdk.Message, error) {
  client := Client{}
  return client.Delete (params)
}
