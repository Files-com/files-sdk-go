package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type As2OutgoingMessage struct {
	Id                 int64           `json:"id,omitempty"`
	As2PartnerId       int64           `json:"as2_partner_id,omitempty"`
	Uuid               string          `json:"uuid,omitempty"`
	HttpHeaders        json.RawMessage `json:"http_headers,omitempty"`
	ActivityLog        string          `json:"activity_log,omitempty"`
	ProcessingResult   string          `json:"processing_result,omitempty"`
	Mic                string          `json:"mic,omitempty"`
	MessageId          string          `json:"message_id,omitempty"`
	BodySize           string          `json:"body_size,omitempty"`
	AttachmentFilename string          `json:"attachment_filename,omitempty"`
	CreatedAt          time.Time       `json:"created_at,omitempty"`
}

type As2OutgoingMessageCollection []As2OutgoingMessage

type As2OutgoingMessageListParams struct {
	Cursor       string `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage      int64  `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	As2PartnerId int64  `url:"as2_partner_id,omitempty" required:"false" json:"as2_partner_id,omitempty"`
	lib.ListParams
}

func (a *As2OutgoingMessage) UnmarshalJSON(data []byte) error {
	type as2OutgoingMessage As2OutgoingMessage
	var v as2OutgoingMessage
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = As2OutgoingMessage(v)
	return nil
}

func (a *As2OutgoingMessageCollection) UnmarshalJSON(data []byte) error {
	type as2OutgoingMessages []As2OutgoingMessage
	var v as2OutgoingMessages
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = As2OutgoingMessageCollection(v)
	return nil
}

func (a *As2OutgoingMessageCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
