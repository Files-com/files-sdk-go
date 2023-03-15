package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type As2IncomingMessage struct {
	Id                          int64           `json:"id,omitempty" path:"id"`
	As2PartnerId                int64           `json:"as2_partner_id,omitempty" path:"as2_partner_id"`
	As2StationId                int64           `json:"as2_station_id,omitempty" path:"as2_station_id"`
	Uuid                        string          `json:"uuid,omitempty" path:"uuid"`
	ContentType                 string          `json:"content_type,omitempty" path:"content_type"`
	HttpHeaders                 json.RawMessage `json:"http_headers,omitempty" path:"http_headers"`
	ActivityLog                 string          `json:"activity_log,omitempty" path:"activity_log"`
	ProcessingResult            string          `json:"processing_result,omitempty" path:"processing_result"`
	ProcessingResultDescription string          `json:"processing_result_description,omitempty" path:"processing_result_description"`
	Mic                         string          `json:"mic,omitempty" path:"mic"`
	MicAlgo                     string          `json:"mic_algo,omitempty" path:"mic_algo"`
	As2To                       string          `json:"as2_to,omitempty" path:"as2_to"`
	As2From                     string          `json:"as2_from,omitempty" path:"as2_from"`
	MessageId                   string          `json:"message_id,omitempty" path:"message_id"`
	Subject                     string          `json:"subject,omitempty" path:"subject"`
	Date                        string          `json:"date,omitempty" path:"date"`
	BodySize                    string          `json:"body_size,omitempty" path:"body_size"`
	AttachmentFilename          string          `json:"attachment_filename,omitempty" path:"attachment_filename"`
	Ip                          string          `json:"ip,omitempty" path:"ip"`
	CreatedAt                   *time.Time      `json:"created_at,omitempty" path:"created_at"`
	HttpResponseCode            string          `json:"http_response_code,omitempty" path:"http_response_code"`
	HttpResponseHeaders         json.RawMessage `json:"http_response_headers,omitempty" path:"http_response_headers"`
	RecipientSerial             string          `json:"recipient_serial,omitempty" path:"recipient_serial"`
	HexRecipientSerial          string          `json:"hex_recipient_serial,omitempty" path:"hex_recipient_serial"`
	RecipientIssuer             string          `json:"recipient_issuer,omitempty" path:"recipient_issuer"`
	MessageReceived             *bool           `json:"message_received,omitempty" path:"message_received"`
	MessageDecrypted            *bool           `json:"message_decrypted,omitempty" path:"message_decrypted"`
	MessageSignatureVerified    *bool           `json:"message_signature_verified,omitempty" path:"message_signature_verified"`
	MessageProcessingSuccess    *bool           `json:"message_processing_success,omitempty" path:"message_processing_success"`
	MessageMdnReturned          *bool           `json:"message_mdn_returned,omitempty" path:"message_mdn_returned"`
	EncryptedUri                string          `json:"encrypted_uri,omitempty" path:"encrypted_uri"`
	SmimeSignedUri              string          `json:"smime_signed_uri,omitempty" path:"smime_signed_uri"`
	SmimeUri                    string          `json:"smime_uri,omitempty" path:"smime_uri"`
	RawUri                      string          `json:"raw_uri,omitempty" path:"raw_uri"`
	MdnResponseUri              string          `json:"mdn_response_uri,omitempty" path:"mdn_response_uri"`
}

type As2IncomingMessageCollection []As2IncomingMessage

type As2IncomingMessageListParams struct {
	SortBy       json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter       json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt     json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq   json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt     json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq   json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	As2PartnerId int64           `url:"as2_partner_id,omitempty" required:"false" json:"as2_partner_id,omitempty" path:"as2_partner_id"`
	lib.ListParams
}

func (a *As2IncomingMessage) UnmarshalJSON(data []byte) error {
	type as2IncomingMessage As2IncomingMessage
	var v as2IncomingMessage
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = As2IncomingMessage(v)
	return nil
}

func (a *As2IncomingMessageCollection) UnmarshalJSON(data []byte) error {
	type as2IncomingMessages As2IncomingMessageCollection
	var v as2IncomingMessages
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
