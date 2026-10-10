// Package file_comment_reaction provides the Files.com FileCommentReaction API client.
package file_comment_reaction

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Client calls the FileCommentReaction API using its embedded Config.
type Client struct {
	files_sdk.Config
}

// Create calls POST /file_comment_reactions.
func (c *Client) Create(params files_sdk.FileCommentReactionCreateParams, opts ...files_sdk.RequestResponseOption) (fileCommentReaction files_sdk.FileCommentReaction, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/file_comment_reactions", Params: params, Entity: &fileCommentReaction}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.FileCommentReactionCreateParams, opts ...files_sdk.RequestResponseOption) (fileCommentReaction files_sdk.FileCommentReaction, err error) {
	return (&Client{}).Create(params, opts...)
}

// Delete calls DELETE /file_comment_reactions/{id}.
func (c *Client) Delete(params files_sdk.FileCommentReactionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/file_comment_reactions/{id}", Params: params, Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(params files_sdk.FileCommentReactionDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}
