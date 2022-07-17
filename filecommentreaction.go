package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type FileCommentReaction struct {
	Id            int64  `json:"id,omitempty" path:"id"`
	Emoji         string `json:"emoji,omitempty" path:"emoji"`
	UserId        int64  `json:"user_id,omitempty" path:"user_id"`
	FileCommentId int64  `json:"file_comment_id,omitempty" path:"file_comment_id"`
}

type FileCommentReactionCollection []FileCommentReaction

type FileCommentReactionCreateParams struct {
	UserId        int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	FileCommentId int64  `url:"file_comment_id,omitempty" required:"true" json:"file_comment_id,omitempty" path:"file_comment_id"`
	Emoji         string `url:"emoji,omitempty" required:"true" json:"emoji,omitempty" path:"emoji"`
}

type FileCommentReactionDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

func (f *FileCommentReaction) UnmarshalJSON(data []byte) error {
	type fileCommentReaction FileCommentReaction
	var v fileCommentReaction
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FileCommentReaction(v)
	return nil
}

func (f *FileCommentReactionCollection) UnmarshalJSON(data []byte) error {
	type fileCommentReactions FileCommentReactionCollection
	var v fileCommentReactions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FileCommentReactionCollection(v)
	return nil
}

func (f *FileCommentReactionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
