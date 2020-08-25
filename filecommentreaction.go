package files_sdk

import (
	"encoding/json"
)

type FileCommentReaction struct {
	Id            int64  `json:"id,omitempty"`
	Emoji         string `json:"emoji,omitempty"`
	UserId        int64  `json:"user_id,omitempty"`
	FileCommentId int64  `json:"file_comment_id,omitempty"`
}

type FileCommentReactionCollection []FileCommentReaction

type FileCommentReactionCreateParams struct {
	UserId        int64  `url:"user_id,omitempty" required:"false"`
	FileCommentId int64  `url:"file_comment_id,omitempty" required:"true"`
	Emoji         string `url:"emoji,omitempty" required:"true"`
}

type FileCommentReactionDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
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
