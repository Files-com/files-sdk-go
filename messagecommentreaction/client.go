package message_comment_reaction

import (
	"context"

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

func (c *Client) List(ctx context.Context, params files_sdk.MessageCommentReactionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/message_comment_reactions", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.MessageCommentReactionCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list, opts...)
	return i, nil
}

func List(ctx context.Context, params files_sdk.MessageCommentReactionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.MessageCommentReactionFindParams, opts ...files_sdk.RequestResponseOption) (messageCommentReaction files_sdk.MessageCommentReaction, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/message_comment_reactions/{id}", Params: params, Entity: &messageCommentReaction}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.MessageCommentReactionFindParams, opts ...files_sdk.RequestResponseOption) (messageCommentReaction files_sdk.MessageCommentReaction, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Create(ctx context.Context, params files_sdk.MessageCommentReactionCreateParams, opts ...files_sdk.RequestResponseOption) (messageCommentReaction files_sdk.MessageCommentReaction, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/message_comment_reactions", Params: params, Entity: &messageCommentReaction}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.MessageCommentReactionCreateParams, opts ...files_sdk.RequestResponseOption) (messageCommentReaction files_sdk.MessageCommentReaction, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.MessageCommentReactionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/message_comment_reactions/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.MessageCommentReactionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
