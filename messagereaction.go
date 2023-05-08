package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type MessageReaction struct {
	Id     int64  `json:"id,omitempty" path:"id"`
	Emoji  string `json:"emoji,omitempty" path:"emoji"`
	UserId int64  `json:"user_id,omitempty" path:"user_id"`
}

func (m MessageReaction) Identifier() interface{} {
	return m.Id
}

type MessageReactionCollection []MessageReaction

type MessageReactionListParams struct {
	UserId    int64 `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	MessageId int64 `url:"message_id,omitempty" required:"true" json:"message_id,omitempty" path:"message_id"`
	ListParams
}

type MessageReactionFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type MessageReactionCreateParams struct {
	UserId int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	Emoji  string `url:"emoji,omitempty" required:"true" json:"emoji,omitempty" path:"emoji"`
}

type MessageReactionDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (m *MessageReaction) UnmarshalJSON(data []byte) error {
	type messageReaction MessageReaction
	var v messageReaction
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*m = MessageReaction(v)
	return nil
}

func (m *MessageReactionCollection) UnmarshalJSON(data []byte) error {
	type messageReactions MessageReactionCollection
	var v messageReactions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*m = MessageReactionCollection(v)
	return nil
}

func (m *MessageReactionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*m))
	for i, v := range *m {
		ret[i] = v
	}

	return &ret
}
