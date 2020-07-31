package file_comment_reaction

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  files_sdk "github.com/Files-com/files-sdk-go"
)

type Client struct {
	files_sdk.Config
}


func (c *Client) Create (params files_sdk.FileCommentReactionCreateParams) (files_sdk.FileCommentReaction, error) {
  fileCommentReaction := files_sdk.FileCommentReaction{}
	  path := "/file_comment_reactions"
	data, _, err := files_sdk.Call("POST", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return fileCommentReaction, err
	}
	if err := fileCommentReaction.UnmarshalJSON(*data); err != nil {
	return fileCommentReaction, err
	}

	return  fileCommentReaction, nil
}

func Create (params files_sdk.FileCommentReactionCreateParams) (files_sdk.FileCommentReaction, error) {
  client := Client{}
  return client.Create (params)
}

func (c *Client) Delete (params files_sdk.FileCommentReactionDeleteParams) (files_sdk.FileCommentReaction, error) {
  fileCommentReaction := files_sdk.FileCommentReaction{}
  	path := "/file_comment_reactions/" + lib.QueryEscape(string(params.Id)) + ""
	data, _, err := files_sdk.Call("DELETE", c.Config, path, lib.ExportParams(params))
	if err != nil {
	  return fileCommentReaction, err
	}
	if err := fileCommentReaction.UnmarshalJSON(*data); err != nil {
	return fileCommentReaction, err
	}

	return  fileCommentReaction, nil
}

func Delete (params files_sdk.FileCommentReactionDeleteParams) (files_sdk.FileCommentReaction, error) {
  client := Client{}
  return client.Delete (params)
}
