package file_comment_reaction

import (
	"context"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(ctx context.Context, params files_sdk.FileCommentReactionCreateParams, opts ...files_sdk.RequestResponseOption) (fileCommentReaction files_sdk.FileCommentReaction, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/file_comment_reactions", Params: params, Entity: &fileCommentReaction}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.FileCommentReactionCreateParams, opts ...files_sdk.RequestResponseOption) (fileCommentReaction files_sdk.FileCommentReaction, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.FileCommentReactionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/file_comment_reactions/{id}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.FileCommentReactionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}
