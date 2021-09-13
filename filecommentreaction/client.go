package file_comment_reaction

import (
	"context"
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(ctx context.Context, params files_sdk.FileCommentReactionCreateParams) (files_sdk.FileCommentReaction, error) {
	fileCommentReaction := files_sdk.FileCommentReaction{}
	path := "/file_comment_reactions"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return fileCommentReaction, err
	}
	data, res, err := files_sdk.Call(ctx, "POST", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return fileCommentReaction, err
	}
	if res.StatusCode == 204 {
		return fileCommentReaction, nil
	}
	if err := fileCommentReaction.UnmarshalJSON(*data); err != nil {
		return fileCommentReaction, err
	}

	return fileCommentReaction, nil
}

func Create(ctx context.Context, params files_sdk.FileCommentReactionCreateParams) (files_sdk.FileCommentReaction, error) {
	return (&Client{}).Create(ctx, params)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.FileCommentReactionDeleteParams) (files_sdk.FileCommentReaction, error) {
	fileCommentReaction := files_sdk.FileCommentReaction{}
	if params.Id == 0 {
		return fileCommentReaction, lib.CreateError(params, "Id")
	}
	path := "/file_comment_reactions/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return fileCommentReaction, err
	}
	data, res, err := files_sdk.Call(ctx, "DELETE", c.Config, path, exportedParams)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return fileCommentReaction, err
	}
	if res.StatusCode == 204 {
		return fileCommentReaction, nil
	}
	if err := fileCommentReaction.UnmarshalJSON(*data); err != nil {
		return fileCommentReaction, err
	}

	return fileCommentReaction, nil
}

func Delete(ctx context.Context, params files_sdk.FileCommentReactionDeleteParams) (files_sdk.FileCommentReaction, error) {
	return (&Client{}).Delete(ctx, params)
}
