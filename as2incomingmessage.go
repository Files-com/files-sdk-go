package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type As2IncomingMessage struct {
	Id                       int64           `json:"id,omitempty"`
	As2PartnerId             int64           `json:"as2_partner_id,omitempty"`
	As2StationId             int64           `json:"as2_station_id,omitempty"`
	Uuid                     string          `json:"uuid,omitempty"`
	ContentType              string          `json:"content_type,omitempty"`
	HttpHeaders              json.RawMessage `json:"http_headers,omitempty"`
	ActivityLog              string          `json:"activity_log,omitempty"`
	ProcessingResult         string          `json:"processing_result,omitempty"`
	Mic                      string          `json:"mic,omitempty"`
	MicAlgo                  string          `json:"mic_algo,omitempty"`
	As2To                    string          `json:"as2_to,omitempty"`
	As2From                  string          `json:"as2_from,omitempty"`
	MessageId                string          `json:"message_id,omitempty"`
	Subject                  string          `json:"subject,omitempty"`
	BodySize                 string          `json:"body_size,omitempty"`
	AttachmentFilename       string          `json:"attachment_filename,omitempty"`
	Ip                       string          `json:"ip,omitempty"`
	CreatedAt                time.Time       `json:"created_at,omitempty"`
	HttpResponseCode         string          `json:"http_response_code,omitempty"`
	HttpResponseHeaders      json.RawMessage `json:"http_response_headers,omitempty"`
	MessageReceived          *bool           `json:"message_received,omitempty"`
	MessageDecrypted         *bool           `json:"message_decrypted,omitempty"`
	MessageSignatureVerified *bool           `json:"message_signature_verified,omitempty"`
	MessageProcessingSuccess *bool           `json:"message_processing_success,omitempty"`
	MessageMdnReturned       *bool           `json:"message_mdn_returned,omitempty"`
	EncryptedUri             string          `json:"encrypted_uri,omitempty"`
	SmimeSignedUri           string          `json:"smime_signed_uri,omitempty"`
	SmimeUri                 string          `json:"smime_uri,omitempty"`
	RawUri                   string          `json:"raw_uri,omitempty"`
	MdnResponseUri           string          `json:"mdn_response_uri,omitempty"`
}

type As2IncomingMessageCollection []As2IncomingMessage

type As2IncomingMessageListParams struct {
	Cursor       string `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage      int64  `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	As2PartnerId int64  `url:"as2_partner_id,omitempty" required:"false" json:"as2_partner_id,omitempty"`
	lib.ListParams
}

func (a *As2IncomingMessage) UnmarshalJSON(data []byte) error {
	type as2IncomingMessage As2IncomingMessage
	var v as2IncomingMessage
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = As2IncomingMessage(v)
	return nil
}

func (a *As2IncomingMessageCollection) UnmarshalJSON(data []byte) error {
	type as2IncomingMessages []As2IncomingMessage
	var v as2IncomingMessages
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = As2IncomingMessageCollection(v)
	return nil
}

func (a *As2IncomingMessageCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}