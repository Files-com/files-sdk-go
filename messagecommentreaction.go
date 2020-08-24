package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type MessageCommentReaction struct {
	Id     int64  `json:"id,omitempty"`
	Emoji  string `json:"emoji,omitempty"`
	UserId int64  `json:"user_id,omitempty"`
}

type MessageCommentReactionCollection []MessageCommentReaction

type MessageCommentReactionListParams struct {
	UserId           int64  `url:"user_id,omitempty"`
	Page             int    `url:"page,omitempty"`
	PerPage          int    `url:"per_page,omitempty"`
	Action           string `url:"action,omitempty"`
	Cursor           string `url:"cursor,omitempty"`
	MessageCommentId int64  `url:"message_comment_id,omitempty"`
	lib.ListParams
}

type MessageCommentReactionFindParams struct {
	Id int64 `url:"-,omitempty"`
}

type MessageCommentReactionCreateParams struct {
	UserId int64  `url:"user_id,omitempty"`
	Emoji  string `url:"emoji,omitempty"`
}

type MessageCommentReactionDeleteParams struct {
	Id int64 `url:"-,omitempty"`
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
