package message_comment_reaction

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

func (i *Iter) MessageCommentReaction() files_sdk.MessageCommentReaction {
	return i.Current().(files_sdk.MessageCommentReaction)
}

func (c *Client) List(params files_sdk.MessageCommentReactionListParams) *Iter {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/message_comment_reactions"

	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, i.ExportParams())
		defaultValue := make([]interface{}, 0)
        if err != nil {
          return &defaultValue, "", err
        }
		list := files_sdk.MessageCommentReactionCollection{}
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

func List(params files_sdk.MessageCommentReactionListParams) *Iter {
  client := Client{}
  return client.List (params)
}

func (c *Client) Find (params files_sdk.MessageCommentReactionFindParams) (files_sdk.MessageCommentReaction, error) {
  messageCommentReaction := files_sdk.MessageCommentReaction{}
  	path := "/message_comment_reactions/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return messageCommentReaction, err
	}
	if err := messageCommentReaction.UnmarshalJSON(*data); err != nil {
	return messageCommentReaction, err
	}

	return  messageCommentReaction, nil
}

func Find (params files_sdk.MessageCommentReactionFindParams) (files_sdk.MessageCommentReaction, error) {
  client := Client{}
  return client.Find (params)
}

func (c *Client) Create (params files_sdk.MessageCommentReactionCreateParams) (files_sdk.MessageCommentReaction, error) {
  messageCommentReaction := files_sdk.MessageCommentReaction{}
	  path := "/message_comment_reactions"
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return messageCommentReaction, err
	}
	if err := messageCommentReaction.UnmarshalJSON(*data); err != nil {
	return messageCommentReaction, err
	}

	return  messageCommentReaction, nil
}

func Create (params files_sdk.MessageCommentReactionCreateParams) (files_sdk.MessageCommentReaction, error) {
  client := Client{}
  return client.Create (params)
}

func (c *Client) Delete (params files_sdk.MessageCommentReactionDeleteParams) (files_sdk.MessageCommentReaction, error) {
  messageCommentReaction := files_sdk.MessageCommentReaction{}
  	path := "/message_comment_reactions/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return messageCommentReaction, err
	}
	if err := messageCommentReaction.UnmarshalJSON(*data); err != nil {
	return messageCommentReaction, err
	}

	return  messageCommentReaction, nil
}

func Delete (params files_sdk.MessageCommentReactionDeleteParams) (files_sdk.MessageCommentReaction, error) {
  client := Client{}
  return client.Delete (params)
}
