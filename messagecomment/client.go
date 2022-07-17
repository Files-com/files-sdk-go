package message_comment

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

func (i *Iter) MessageComment() files_sdk.MessageComment {
	return i.Current().(files_sdk.MessageComment)
}

func (c *Client) List(ctx context.Context, params files_sdk.MessageCommentListParams) (*Iter, error) {
	i := &Iter{Iter: &lib.Iter{}}
	path, err := lib.BuildPath("/message_comments", params)
	if err != nil {
		return i, err
	}
	i.ListParams = &params
	list := files_sdk.MessageCommentCollection{}
	i.Query = listquery.Build(ctx, c.Config, path, &list)
	return i, nil
}

func List(ctx context.Context, params files_sdk.MessageCommentListParams) (*Iter, error) {
	return (&Client{}).List(ctx, params)
}

func (c *Client) Find(ctx context.Context, params files_sdk.MessageCommentFindParams) (messageComment files_sdk.MessageComment, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/message_comments/{id}", Params: params, Entity: &messageComment})
	return
}

func Find(ctx context.Context, params files_sdk.MessageCommentFindParams) (messageComment files_sdk.MessageComment, err error) {
	return (&Client{}).Find(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.MessageCommentCreateParams) (messageComment files_sdk.MessageComment, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/message_comments", Params: params, Entity: &messageComment})
	return
}

func Create(ctx context.Context, params files_sdk.MessageCommentCreateParams) (messageComment files_sdk.MessageComment, err error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Update(ctx context.Context, params files_sdk.MessageCommentUpdateParams) (messageComment files_sdk.MessageComment, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/message_comments/{id}", Params: params, Entity: &messageComment})
	return
}

func Update(ctx context.Context, params files_sdk.MessageCommentUpdateParams) (messageComment files_sdk.MessageComment, err error) {
	return (&Client{}).Update(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.MessageCommentDeleteParams) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/message_comments/{id}", Params: params, Entity: nil})
	return
}

func Delete(ctx context.Context, params files_sdk.MessageCommentDeleteParams) (err error) {
	return (&Client{}).Delete(ctx, params)
}
