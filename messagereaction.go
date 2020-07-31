package files_sdk

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  "encoding/json"
)

type MessageReaction struct {
  Id int `json:"id,omitempty"`
  Emoji string `json:"emoji,omitempty"`
  UserId int `json:"user_id,omitempty"`
}

type MessageReactionCollection []MessageReaction

type MessageReactionListParams struct {
  UserId int `url:"user_id,omitempty"`
  Page int `url:"page,omitempty"`
  PerPage int `url:"per_page,omitempty"`
  Action string `url:"action,omitempty"`
  MessageId int `url:"message_id,omitempty"`
  lib.ListParams
}

type MessageReactionFindParams struct {
  Id int `url:"-,omitempty"`
}

type MessageReactionCreateParams struct {
  UserId int `url:"user_id,omitempty"`
  Emoji string `url:"emoji,omitempty"`
}

type MessageReactionDeleteParams struct {
  Id int `url:"-,omitempty"`
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

