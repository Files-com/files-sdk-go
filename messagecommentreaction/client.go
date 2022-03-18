package message_comment_reaction

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

func (i *Iter) MessageCommentReaction() files_sdk.MessageCommentReaction {
	return i.Current().(files_sdk.MessageCommentReaction)
}

func (c *Client) List(ctx context.Context, params files_sdk.MessageCommentReactionListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	params.ListParams.Set(params.Page, params.PerPage, params.Cursor, params.MaxPages)
	path := "/message_comment_reactions"
	i.ListParams = &params
	list := files_sdk.MessageCommentReactionCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.MessageCommentReactionListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.MessageCommentReactionFindParams) (files_sdk.MessageCommentReaction, error) {
	messageCommentReaction := files_sdk.MessageCommentReaction{}
	if params.Id == 0 {
		return messageCommentReaction, lib.CreateError(params, "Id")
	}
	path := "/message_comment_reactions/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "GET", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
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

func Find(ctx context.Context, params files_sdk.MessageCommentReactionFindParams) (files_sdk.MessageCommentReaction, error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.MessageCommentReactionCreateParams) (files_sdk.MessageCommentReaction, error) {
	messageCommentReaction := files_sdk.MessageCommentReaction{}
	path := "/message_comment_reactions"
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
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

func Create(ctx context.Context, params files_sdk.MessageCommentReactionCreateParams) (files_sdk.MessageCommentReaction, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.MessageCommentReactionDeleteParams) (files_sdk.MessageCommentReaction, error) {
	messageCommentReaction := files_sdk.MessageCommentReaction{}
	if params.Id == 0 {
		return messageCommentReaction, lib.CreateError(params, "Id")
	}
	path := "/message_comment_reactions/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams := lib.Params{Params: params}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil && res.Body != nil {
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

func Delete(ctx context.Context, params files_sdk.MessageCommentReactionDeleteParams) (files_sdk.MessageCommentReaction, error) {
	return (&Client{}).Delete(ctx, params)
}
