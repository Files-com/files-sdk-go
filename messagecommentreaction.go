package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type MessageCommentReaction struct {
	Id     int64  `json:"id,omitempty"`
	Emoji  string `json:"emoji,omitempty"`
	UserId int64  `json:"user_id,omitempty"`
}

type MessageCommentReactionCollection []MessageCommentReaction

type MessageCommentReactionListParams struct {
	UserId           int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	Cursor           string `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage          int64  `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	MessageCommentId int64  `url:"message_comment_id,omitempty" required:"true" json:"message_comment_id,omitempty"`
	lib.ListParams
}

type MessageCommentReactionFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

type MessageCommentReactionCreateParams struct {
	UserId int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty"`
	Emoji  string `url:"emoji,omitempty" required:"true" json:"emoji,omitempty"`
}

type MessageCommentReactionDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

func (m *MessageCommentReaction) UnmarshalJSON(data []byte) error {
	type messageCommentReaction MessageCommentReaction
	var v messageCommentReaction
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = MessageCommentReaction(v)
	return nil
}

func (m *MessageCommentReactionCollection) UnmarshalJSON(data []byte) error {
	type messageCommentReactions []MessageCommentReaction
	var v messageCommentReactions
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = MessageCommentReactionCollection(v)
	return nil
}

func (m *MessageCommentReactionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*m))
	for i, v := range *m {
		ret[i] = v
	}

	return &ret
}
