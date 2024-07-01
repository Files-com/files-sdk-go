package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type MessageCommentReaction struct {
	Id     int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Emoji  string `json:"emoji,omitempty" path:"emoji,omitempty" url:"emoji,omitempty"`
	UserId int64  `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
}

func (m MessageCommentReaction) Identifier() interface{} {
	return m.Id
}

type MessageCommentReactionCollection []MessageCommentReaction

type MessageCommentReactionListParams struct {
	UserId           int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	Action           string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	MessageCommentId int64  `url:"message_comment_id,omitempty" required:"true" json:"message_comment_id,omitempty" path:"message_comment_id"`
	ListParams
}

type MessageCommentReactionFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type MessageCommentReactionCreateParams struct {
	UserId int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	Emoji  string `url:"emoji,omitempty" required:"true" json:"emoji,omitempty" path:"emoji"`
}

type MessageCommentReactionDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (m *MessageCommentReaction) UnmarshalJSON(data []byte) error {
	type messageCommentReaction MessageCommentReaction
	var v messageCommentReaction
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*m = MessageCommentReaction(v)
	return nil
}

func (m *MessageCommentReactionCollection) UnmarshalJSON(data []byte) error {
	type messageCommentReactions MessageCommentReactionCollection
	var v messageCommentReactions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
