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
	UserId    int64  `url:"user_id,omitempty"`
	Page      int    `url:"page,omitempty"`
	PerPage   int    `url:"per_page,omitempty"`
	Action    string `url:"action,omitempty"`
	MessageId int64  `url:"message_id,omitempty"`
	lib.ListParams
}

type MessageReactionFindParams struct {
	Id int64 `url:"-,omitempty"`
}

type MessageReactionCreateParams struct {
	UserId int64  `url:"user_id,omitempty"`
	Emoji  string `url:"emoji,omitempty"`
}

type MessageReactionDeleteParams struct {
	Id int64 `url:"-,omitempty"`
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
