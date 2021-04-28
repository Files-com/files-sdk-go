package file_comment_reaction

import (
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.FileCommentReactionCreateParams) (files_sdk.FileCommentReaction, error) {
	fileCommentReaction := files_sdk.FileCommentReaction{}
	path := "/file_comment_reactions"
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return fileCommentReaction, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParams)
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

func Create(params files_sdk.FileCommentReactionCreateParams) (files_sdk.FileCommentReaction, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Delete(params files_sdk.FileCommentReactionDeleteParams) (files_sdk.FileCommentReaction, error) {
	fileCommentReaction := files_sdk.FileCommentReaction{}
	if params.Id == 0 {
		return fileCommentReaction, lib.CreateError(params, "Id")
	}
	path := "/file_comment_reactions/" + strconv.FormatInt(params.Id, 10) + ""
	exportedParams, err := lib.ExportParams(params)
	if err != nil {
		return fileCommentReaction, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParams)
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

func Delete(params files_sdk.FileCommentReactionDeleteParams) (files_sdk.FileCommentReaction, error) {
	return (&Client{}).Delete(params)
}
