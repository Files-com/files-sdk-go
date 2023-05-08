package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type As2OutgoingMessage struct {
	Id                          int64           `json:"id,omitempty" path:"id"`
	As2PartnerId                int64           `json:"as2_partner_id,omitempty" path:"as2_partner_id"`
	As2StationId                int64           `json:"as2_station_id,omitempty" path:"as2_station_id"`
	Uuid                        string          `json:"uuid,omitempty" path:"uuid"`
	HttpHeaders                 json.RawMessage `json:"http_headers,omitempty" path:"http_headers"`
	ActivityLog                 string          `json:"activity_log,omitempty" path:"activity_log"`
	ProcessingResult            string          `json:"processing_result,omitempty" path:"processing_result"`
	ProcessingResultDescription string          `json:"processing_result_description,omitempty" path:"processing_result_description"`
	Mic                         string          `json:"mic,omitempty" path:"mic"`
	MicSha256                   string          `json:"mic_sha_256,omitempty" path:"mic_sha_256"`
	As2To                       string          `json:"as2_to,omitempty" path:"as2_to"`
	As2From                     string          `json:"as2_from,omitempty" path:"as2_from"`
	Date                        string          `json:"date,omitempty" path:"date"`
	MessageId                   string          `json:"message_id,omitempty" path:"message_id"`
	BodySize                    string          `json:"body_size,omitempty" path:"body_size"`
	AttachmentFilename          string          `json:"attachment_filename,omitempty" path:"attachment_filename"`
	CreatedAt                   *time.Time      `json:"created_at,omitempty" path:"created_at"`
	HttpResponseCode            string          `json:"http_response_code,omitempty" path:"http_response_code"`
	HttpResponseHeaders         json.RawMessage `json:"http_response_headers,omitempty" path:"http_response_headers"`
	HttpTransmissionDuration    string          `json:"http_transmission_duration,omitempty" path:"http_transmission_duration"`
	MdnReceived                 *bool           `json:"mdn_received,omitempty" path:"mdn_received"`
	MdnValid                    *bool           `json:"mdn_valid,omitempty" path:"mdn_valid"`
	MdnSignatureVerified        *bool           `json:"mdn_signature_verified,omitempty" path:"mdn_signature_verified"`
	MdnMessageIdMatched         *bool           `json:"mdn_message_id_matched,omitempty" path:"mdn_message_id_matched"`
	MdnMicMatched               *bool           `json:"mdn_mic_matched,omitempty" path:"mdn_mic_matched"`
	MdnProcessingSuccess        *bool           `json:"mdn_processing_success,omitempty" path:"mdn_processing_success"`
	RawUri                      string          `json:"raw_uri,omitempty" path:"raw_uri"`
	SmimeUri                    string          `json:"smime_uri,omitempty" path:"smime_uri"`
	SmimeSignedUri              string          `json:"smime_signed_uri,omitempty" path:"smime_signed_uri"`
	EncryptedUri                string          `json:"encrypted_uri,omitempty" path:"encrypted_uri"`
	MdnResponseUri              string          `json:"mdn_response_uri,omitempty" path:"mdn_response_uri"`
}

func (a As2OutgoingMessage) Identifier() interface{} {
	return a.Id
}

type As2OutgoingMessageCollection []As2OutgoingMessage

type As2OutgoingMessageListParams struct {
	SortBy       json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter       json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt     json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq   json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt     json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq   json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	As2PartnerId int64           `url:"as2_partner_id,omitempty" required:"false" json:"as2_partner_id,omitempty" path:"as2_partner_id"`
	ListParams
}

func (a *As2OutgoingMessage) UnmarshalJSON(data []byte) error {
	type as2OutgoingMessage As2OutgoingMessage
	var v as2OutgoingMessage
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = As2OutgoingMessage(v)
	return nil
}

func (a *As2OutgoingMessageCollection) UnmarshalJSON(data []byte) error {
	type as2OutgoingMessages As2OutgoingMessageCollection
	var v as2OutgoingMessages
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
