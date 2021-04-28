package message_comment_reaction

import (
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
	listquery "github.com/Files-com/files-sdk-go/listquery"
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
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/message_comment_reactions"
	i.ListParams = &params
	list := files_sdk.MessageCommentReactionCollection{}
	i.Query = listquery.Build(i, c.Config, path, &list)
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
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
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
