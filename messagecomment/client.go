package message_comment

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

func (i *Iter) MessageComment() files_sdk.MessageComment {
	return i.Current().(files_sdk.MessageComment)
}

func (c *Client) List(params files_sdk.MessageCommentListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/message_comments"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
        if err != nil {
          return &defaultValue, "", err
        }
		list := files_sdk.MessageCommentCollection{}
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

func List(params files_sdk.MessageCommentListParams) *Iter {
  client := Client{}
  return client.List (params)
}

func (c *Client) Find (params files_sdk.MessageCommentFindParams) (files_sdk.MessageComment, error) {
  messageComment := files_sdk.MessageComment{}
  	path := "/message_comments/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return messageComment, err
	}
	if err := messageComment.UnmarshalJSON(*data); err != nil {
	return messageComment, err
	}

	return  messageComment, nil
}

func Find (params files_sdk.MessageCommentFindParams) (files_sdk.MessageComment, error) {
  client := Client{}
  return client.Find (params)
}

func (c *Client) Create (params files_sdk.MessageCommentCreateParams) (files_sdk.MessageComment, error) {
  messageComment := files_sdk.MessageComment{}
	  path := "/message_comments"
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return messageComment, err
	}
	if err := messageComment.UnmarshalJSON(*data); err != nil {
	return messageComment, err
	}

	return  messageComment, nil
}

func Create (params files_sdk.MessageCommentCreateParams) (files_sdk.MessageComment, error) {
  client := Client{}
  return client.Create (params)
}

func (c *Client) Update (params files_sdk.MessageCommentUpdateParams) (files_sdk.MessageComment, error) {
  messageComment := files_sdk.MessageComment{}
  	path := "/message_comments/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("PATCH", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return messageComment, err
	}
	if err := messageComment.UnmarshalJSON(*data); err != nil {
	return messageComment, err
	}

	return  messageComment, nil
}

func Update (params files_sdk.MessageCommentUpdateParams) (files_sdk.MessageComment, error) {
  client := Client{}
  return client.Update (params)
}

func (c *Client) Delete (params files_sdk.MessageCommentDeleteParams) (files_sdk.MessageComment, error) {
  messageComment := files_sdk.MessageComment{}
  	path := "/message_comments/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return messageComment, err
	}
	if err := messageComment.UnmarshalJSON(*data); err != nil {
	return messageComment, err
	}

	return  messageComment, nil
}

func Delete (params files_sdk.MessageCommentDeleteParams) (files_sdk.MessageComment, error) {
  client := Client{}
  return client.Delete (params)
}
