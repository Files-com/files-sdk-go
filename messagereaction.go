package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/lib"
)

type MessageReaction struct {
	Id     int64  `json:"id,omitempty"`
	Emoji  string `json:"emoji,omitempty"`
	UserId int64  `json:"user_id,omitempty"`
}

type MessageReactionCollection []MessageReaction

type MessageReactionListParams struct {
	UserId    int64  `url:"user_id,omitempty" required:"false"`
	Page      int    `url:"page,omitempty" required:"false"`
	PerPage   int    `url:"per_page,omitempty" required:"false"`
	Action    string `url:"action,omitempty" required:"false"`
	Cursor    string `url:"cursor,omitempty" required:"false"`
	MessageId int64  `url:"message_id,omitempty" required:"true"`
	lib.ListParams
}

type MessageReactionFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type MessageReactionCreateParams struct {
	UserId int64  `url:"user_id,omitempty" required:"false"`
	Emoji  string `url:"emoji,omitempty" required:"true"`
}

type MessageReactionDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (m *MessageReaction) UnmarshalJSON(data []byte) error {
	type messageReaction MessageReaction
	var v messageReaction
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = MessageReaction(v)
	return nil
}

func (m *MessageReactionCollection) UnmarshalJSON(data []byte) error {
	type messageReactions []MessageReaction
	var v messageReactions
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = MessageReactionCollection(v)
	return nil
}
