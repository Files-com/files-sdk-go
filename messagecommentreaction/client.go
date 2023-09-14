package message_comment_reaction

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
	listquery "github.com/Files-com/files-sdk-go/v3/listquery"
)

type Client struct {
	files_sdk.Config
}

type Iter struct {
	*files_sdk.Iter
	*Client
}

func (i *Iter) Reload(opts ...files_sdk.RequestResponseOption) files_sdk.IterI {
	return &Iter{Iter: i.Iter.Reload(opts...).(*files_sdk.Iter), Client: i.Client}
}

func (i *Iter) MessageCommentReaction() files_sdk.MessageCommentReaction {
	return i.Current().(files_sdk.MessageCommentReaction)
}

func (i *Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.MessageCommentReactionFindParams{}
	if id, ok := identifier.(int64); ok {
		params.Id = id
	}
	return i.Client.Find(params, opts...)
}

func (c *Client) List(params files_sdk.MessageCommentReactionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	i := &Iter{Iter: &files_sdk.Iter{}, Client: c}
	path, err := lib.BuildPath("/message_comment_reactions", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.MessageCommentReactionCollection{}
	i.Query = listquery.Build(c.Config, path, &list, opts...)
	return i, nil
}

func List(params files_sdk.MessageCommentReactionListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error) {
	return (&Client{}).List(params, opts...)
}

func (c *Client) Find(params files_sdk.MessageCommentReactionFindParams, opts ...files_sdk.RequestResponseOption) (messageCommentReaction files_sdk.MessageCommentReaction, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/message_comment_reactions/{id}", Params: params, Entity: &messageCommentReaction}, opts...)
	return
}

func Find(params files_sdk.MessageCommentReactionFindParams, opts ...files_sdk.RequestResponseOption) (messageCommentReaction files_sdk.MessageCommentReaction, err error) {
	return (&Client{}).Find(params, opts...)
}

func (c *Client) Create(params files_sdk.MessageCommentReactionCreateParams, opts ...files_sdk.RequestResponseOption) (messageCommentReaction files_sdk.MessageCommentReaction, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/message_comment_reactions", Params: params, Entity: &messageCommentReaction}, opts...)
	return
}

func Create(params files_sdk.MessageCommentReactionCreateParams, opts ...files_sdk.RequestResponseOption) (messageCommentReaction files_sdk.MessageCommentReaction, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) Delete(params files_sdk.MessageCommentReactionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/message_comment_reactions/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(params files_sdk.MessageCommentReactionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
