package files_sdk

import (
  "encoding/json"
)

type FileCommentReaction struct {
  Id int `json:"id,omitempty"`
  Emoji string `json:"emoji,omitempty"`
  UserId int `json:"user_id,omitempty"`
  FileCommentId int `json:"file_comment_id,omitempty"`
}

type FileCommentReactionCollection []FileCommentReaction

type FileCommentReactionCreateParams struct {
  UserId int `url:"user_id,omitempty"`
  FileCommentId int `url:"file_comment_id,omitempty"`
  Emoji string `url:"emoji,omitempty"`
}

type FileCommentReactionDeleteParams struct {
  Id int `url:"-,omitempty"`
}


func (f *FileCommentReaction) UnmarshalJSON(data []byte) error {
	type fileCommentReaction FileCommentReaction
	var v fileCommentReaction
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FileCommentReaction(v)
	return nil
}

func (f *FileCommentReactionCollection) UnmarshalJSON(data []byte) error {
	type fileCommentReactions []FileCommentReaction
	var v fileCommentReactions
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FileCommentReactionCollection(v)
	return nil
}

