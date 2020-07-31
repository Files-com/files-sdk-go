package files_sdk

import (
  lib "github.com/Files-com/files-sdk-go/lib"
  "encoding/json"
)

type Message struct {
  Id int `json:"id,omitempty"`
  Subject string `json:"subject,omitempty"`
  Body string `json:"body,omitempty"`
  Comments []string `json:"comments,omitempty"`
  UserId int `json:"user_id,omitempty"`
  ProjectId int `json:"project_id,omitempty"`
}

type MessageCollection []Message

type MessageListParams struct {
  UserId int `url:"user_id,omitempty"`
  Page int `url:"page,omitempty"`
  PerPage int `url:"per_page,omitempty"`
  Action string `url:"action,omitempty"`
  ProjectId int `url:"project_id,omitempty"`
  lib.ListParams
}

type MessageFindParams struct {
  Id int `url:"-,omitempty"`
}

type MessageCreateParams struct {
  UserId int `url:"user_id,omitempty"`
  ProjectId int `url:"project_id,omitempty"`
  Subject string `url:"subject,omitempty"`
  Body string `url:"body,omitempty"`
}

type MessageUpdateParams struct {
  Id int `url:"-,omitempty"`
  ProjectId int `url:"project_id,omitempty"`
  Subject string `url:"subject,omitempty"`
  Body string `url:"body,omitempty"`
}

type MessageDeleteParams struct {
  Id int `url:"-,omitempty"`
}


func (m *Message) UnmarshalJSON(data []byte) error {
	type message Message
	var v message
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = Message(v)
	return nil
}

func (m *MessageCollection) UnmarshalJSON(data []byte) error {
	type messages []Message
	var v messages
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*m = MessageCollection(v)
	return nil
}

