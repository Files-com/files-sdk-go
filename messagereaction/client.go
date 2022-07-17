package message_reaction

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

func (i *Iter) MessageReaction() files_sdk.MessageReaction {
	return i.Current().(files_sdk.MessageReaction)
}

func (c *Client) List(ctx context.Context, params files_sdk.MessageReactionListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/message_reactions", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.MessageReactionCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.MessageReactionListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.MessageReactionFindParams) (messageReaction files_sdk.MessageReaction, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/message_reactions/{id}", Params: params, Entity: &messageReaction})
	return
}

func Find(ctx context.Context, params files_sdk.MessageReactionFindParams) (messageReaction files_sdk.MessageReaction, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.MessageReactionCreateParams) (messageReaction files_sdk.MessageReaction, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/message_reactions", Params: params, Entity: &messageReaction})
	return
}

func Create(ctx context.Context, params files_sdk.MessageReactionCreateParams) (messageReaction files_sdk.MessageReaction, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.MessageReactionDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/message_reactions/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.MessageReactionDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
