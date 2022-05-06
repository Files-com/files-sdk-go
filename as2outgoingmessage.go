package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type As2OutgoingMessage struct {
	Id                          int64           `json:"id,omitempty"`
	As2PartnerId                int64           `json:"as2_partner_id,omitempty"`
	As2StationId                int64           `json:"as2_station_id,omitempty"`
	Uuid                        string          `json:"uuid,omitempty"`
	HttpHeaders                 json.RawMessage `json:"http_headers,omitempty"`
	ActivityLog                 string          `json:"activity_log,omitempty"`
	ProcessingResult            string          `json:"processing_result,omitempty"`
	ProcessingResultDescription string          `json:"processing_result_description,omitempty"`
	Mic                         string          `json:"mic,omitempty"`
	MicSha256                   string          `json:"mic_sha_256,omitempty"`
	As2To                       string          `json:"as2_to,omitempty"`
	As2From                     string          `json:"as2_from,omitempty"`
	Date                        string          `json:"date,omitempty"`
	MessageId                   string          `json:"message_id,omitempty"`
	BodySize                    string          `json:"body_size,omitempty"`
	AttachmentFilename          string          `json:"attachment_filename,omitempty"`
	CreatedAt                   time.Time       `json:"created_at,omitempty"`
	HttpResponseCode            string          `json:"http_response_code,omitempty"`
	HttpResponseHeaders         json.RawMessage `json:"http_response_headers,omitempty"`
	HttpTransmissionDuration    float32         `json:"http_transmission_duration,omitempty"`
	MdnReceived                 *bool           `json:"mdn_received,omitempty"`
	MdnValid                    *bool           `json:"mdn_valid,omitempty"`
	MdnSignatureVerified        *bool           `json:"mdn_signature_verified,omitempty"`
	MdnMessageIdMatched         *bool           `json:"mdn_message_id_matched,omitempty"`
	MdnMicMatched               *bool           `json:"mdn_mic_matched,omitempty"`
	MdnProcessingSuccess        *bool           `json:"mdn_processing_success,omitempty"`
	RawUri                      string          `json:"raw_uri,omitempty"`
	SmimeUri                    string          `json:"smime_uri,omitempty"`
	SmimeSignedUri              string          `json:"smime_signed_uri,omitempty"`
	EncryptedUri                string          `json:"encrypted_uri,omitempty"`
	MdnResponseUri              string          `json:"mdn_response_uri,omitempty"`
}

type As2OutgoingMessageCollection []As2OutgoingMessage

type As2OutgoingMessageListParams struct {
	Cursor       string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage      int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy       json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter       json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt     json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq   json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike   json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt     json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq   json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	As2PartnerId int64           `url:"as2_partner_id,omitempty" required:"false" json:"as2_partner_id,omitempty"`
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
