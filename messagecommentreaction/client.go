package message_comment_reaction

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

func (i *Iter) MessageCommentReaction() files_sdk.MessageCommentReaction {
	return i.Current().(files_sdk.MessageCommentReaction)
}

func (c *Client) List(params files_sdk.MessageCommentReactionListParams) (*Iter, error) {
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	i := &Iter{Iter: &lib.Iter{}}
	path := "/message_comment_reactions"
	i.ListParams = &params
	exportParams, err := i.ExportParams()
	if err != nil {
		return i, err
	}
	i.Query = func() (*[]interface{}, string, error) {
		data, res, err := files_sdk.Call("GET", c.Config, path, exportParams)
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
	return i, nil
}

func List(params files_sdk.MessageCommentReactionListParams) (*Iter, error) {
	return (&Client{}).List(params)
}

func (c *Client) Find(params files_sdk.MessageCommentReactionFindParams) (files_sdk.MessageCommentReaction, error) {
	messageCommentReaction := files_sdk.MessageCommentReaction{}
	if params.Id == 0 {
		return messageCommentReaction, lib.CreateError(params, "Id")
	}
	path := "/message_comment_reactions/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return messageCommentReaction, err
	}
	data, res, err := files_sdk.Call("GET", c.Config, path, exportedParams)
	if err != nil {
		return messageCommentReaction, err
	}
	if res.StatusCode == 204 {
		return messageCommentReaction, nil
	}
	if err := messageCommentReaction.UnmarshalJSON(*data); err != nil {
		return messageCommentReaction, err
	}

	return messageCommentReaction, nil
}

func Find(params files_sdk.MessageCommentReactionFindParams) (files_sdk.MessageCommentReaction, error) {
	return (&Client{}).Find(params)
}

func (c *Client) Create(params files_sdk.MessageCommentReactionCreateParams) (files_sdk.MessageCommentReaction, error) {
	messageCommentReaction := files_sdk.MessageCommentReaction{}
	path := "/message_comment_reactions"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return messageCommentReaction, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
	if err != nil {
		return messageCommentReaction, err
	}
	if res.StatusCode == 204 {
		return messageCommentReaction, nil
	}
	if err := messageCommentReaction.UnmarshalJSON(*data); err != nil {
		return messageCommentReaction, err
	}

	return messageCommentReaction, nil
}

func Create(params files_sdk.MessageCommentReactionCreateParams) (files_sdk.MessageCommentReaction, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Delete(params files_sdk.MessageCommentReactionDeleteParams) (files_sdk.MessageCommentReaction, error) {
	messageCommentReaction := files_sdk.MessageCommentReaction{}
	if params.Id == 0 {
		return messageCommentReaction, lib.CreateError(params, "Id")
	}
	path := "/message_comment_reactions/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return messageCommentReaction, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
	if err != nil {
		return messageCommentReaction, err
	}
	if res.StatusCode == 204 {
		return messageCommentReaction, nil
	}
	if err := messageCommentReaction.UnmarshalJSON(*data); err != nil {
		return messageCommentReaction, err
	}

	return messageCommentReaction, nil
}

func Delete(params files_sdk.MessageCommentReactionDeleteParams) (files_sdk.MessageCommentReaction, error) {
	return (&Client{}).Delete(params)
}
