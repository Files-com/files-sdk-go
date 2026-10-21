package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// As2OutgoingMessage is a Files.com API resource.
type As2OutgoingMessage struct {
	// Id of the AS2 Partner.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// ID of the Workspace associated with this AS2 Outgoing Message.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Id of the AS2 Partner associated with this message.
	As2PartnerId int64 `json:"as2_partner_id,omitempty" path:"as2_partner_id,omitempty" url:"as2_partner_id,omitempty"`
	// Id of the AS2 Station associated with this message.
	As2StationId int64 `json:"as2_station_id,omitempty" path:"as2_station_id,omitempty" url:"as2_station_id,omitempty"`
	// UUID assigned to this message.
	Uuid string `json:"uuid,omitempty" path:"uuid,omitempty" url:"uuid,omitempty"`
	// HTTP Headers sent with this message.
	HttpHeaders interface{} `json:"http_headers,omitempty" path:"http_headers,omitempty" url:"http_headers,omitempty"`
	// Result of processing.
	ProcessingResult string `json:"processing_result,omitempty" path:"processing_result,omitempty" url:"processing_result,omitempty"`
	// Result of processing description.
	ProcessingResultDescription string `json:"processing_result_description,omitempty" path:"processing_result_description,omitempty" url:"processing_result_description,omitempty"`
	// AS2 Message Integrity Check SHA1.
	Mic string `json:"mic,omitempty" path:"mic,omitempty" url:"mic,omitempty"`
	// AS2 Message Integrity Check SHA256.
	MicSha256 string `json:"mic_sha_256,omitempty" path:"mic_sha_256,omitempty" url:"mic_sha_256,omitempty"`
	// AS2 TO.
	As2To string `json:"as2_to,omitempty" path:"as2_to,omitempty" url:"as2_to,omitempty"`
	// AS2 FROM.
	As2From string `json:"as2_from,omitempty" path:"as2_from,omitempty" url:"as2_from,omitempty"`
	// AS2 Subject.
	As2Subject string `json:"as2_subject,omitempty" path:"as2_subject,omitempty" url:"as2_subject,omitempty"`
	// Date Header.
	Date string `json:"date,omitempty" path:"date,omitempty" url:"date,omitempty"`
	// AS2 Message Id.
	MessageId string `json:"message_id,omitempty" path:"message_id,omitempty" url:"message_id,omitempty"`
	// Encrypted Payload Body Size.
	BodySize string `json:"body_size,omitempty" path:"body_size,omitempty" url:"body_size,omitempty"`
	// Filename of the file being sent.
	AttachmentFilename string `json:"attachment_filename,omitempty" path:"attachment_filename,omitempty" url:"attachment_filename,omitempty"`
	// Message creation date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// HTTP Response Code received for this message.
	HttpResponseCode string `json:"http_response_code,omitempty" path:"http_response_code,omitempty" url:"http_response_code,omitempty"`
	// HTTP Headers received for this message.
	HttpResponseHeaders interface{} `json:"http_response_headers,omitempty" path:"http_response_headers,omitempty" url:"http_response_headers,omitempty"`
	// HTTP transmission duration in seceonds.
	HttpTransmissionDuration float64 `json:"http_transmission_duration,omitempty" path:"http_transmission_duration,omitempty" url:"http_transmission_duration,omitempty"`
	// Did the partner give a response body?
	MdnReceived *bool `json:"mdn_received,omitempty" path:"mdn_received,omitempty" url:"mdn_received,omitempty"`
	// Is the response in MDN format?
	MdnValid *bool `json:"mdn_valid,omitempty" path:"mdn_valid,omitempty" url:"mdn_valid,omitempty"`
	// MDN signature verified?
	MdnSignatureVerified *bool `json:"mdn_signature_verified,omitempty" path:"mdn_signature_verified,omitempty" url:"mdn_signature_verified,omitempty"`
	// MDN message id matched?
	MdnMessageIdMatched *bool `json:"mdn_message_id_matched,omitempty" path:"mdn_message_id_matched,omitempty" url:"mdn_message_id_matched,omitempty"`
	// MDN MIC matched?
	MdnMicMatched *bool `json:"mdn_mic_matched,omitempty" path:"mdn_mic_matched,omitempty" url:"mdn_mic_matched,omitempty"`
	// MDN disposition indicate a successful processing?
	MdnProcessingSuccess *bool `json:"mdn_processing_success,omitempty" path:"mdn_processing_success,omitempty" url:"mdn_processing_success,omitempty"`
	// URL to download the original file contents.
	RawUri string `json:"raw_uri,omitempty" path:"raw_uri,omitempty" url:"raw_uri,omitempty"`
	// URL to download the file contents encoded as smime.
	SmimeUri string `json:"smime_uri,omitempty" path:"smime_uri,omitempty" url:"smime_uri,omitempty"`
	// URL to download the file contents as smime with signature.
	SmimeSignedUri string `json:"smime_signed_uri,omitempty" path:"smime_signed_uri,omitempty" url:"smime_signed_uri,omitempty"`
	// URL to download the encrypted signed smime that is to sent as AS2 body.
	EncryptedUri string `json:"encrypted_uri,omitempty" path:"encrypted_uri,omitempty" url:"encrypted_uri,omitempty"`
	// URL to download the http response body.
	MdnResponseUri string `json:"mdn_response_uri,omitempty" path:"mdn_response_uri,omitempty" url:"mdn_response_uri,omitempty"`
}

// Identifier returns the resource ID.
func (a As2OutgoingMessage) Identifier() interface{} {
	return a.Id
}

// As2OutgoingMessageCollection is a list of As2OutgoingMessage resources.
type As2OutgoingMessageCollection []As2OutgoingMessage

// As2OutgoingMessageListParams contains the request parameters for GET /as2_outgoing_messages.
type As2OutgoingMessageListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `created_at` or `as2_partner_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `as2_station_id`, `workspace_id` or `as2_partner_id`. Valid field combinations are `[ as2_station_id, created_at ]`, `[ workspace_id, created_at ]`, `[ as2_partner_id, created_at ]`, `[ workspace_id, as2_station_id ]`, `[ workspace_id, as2_partner_id ]`, `[ workspace_id, as2_station_id, created_at ]` or `[ workspace_id, as2_partner_id, created_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *As2OutgoingMessage) UnmarshalJSON(data []byte) error {
	type as2OutgoingMessage As2OutgoingMessage
	var v as2OutgoingMessage
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = As2OutgoingMessage(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *As2OutgoingMessageCollection) UnmarshalJSON(data []byte) error {
	type as2OutgoingMessages As2OutgoingMessageCollection
	var v as2OutgoingMessages
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = As2OutgoingMessageCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *As2OutgoingMessageCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
