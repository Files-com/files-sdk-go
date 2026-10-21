package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// As2IncomingMessage is a Files.com API resource.
type As2IncomingMessage struct {
	// Id of the AS2 Partner.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// ID of the Workspace associated with this AS2 Incoming Message.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Id of the AS2 Partner associated with this message.
	As2PartnerId int64 `json:"as2_partner_id,omitempty" path:"as2_partner_id,omitempty" url:"as2_partner_id,omitempty"`
	// Id of the AS2 Station associated with this message.
	As2StationId int64 `json:"as2_station_id,omitempty" path:"as2_station_id,omitempty" url:"as2_station_id,omitempty"`
	// UUID assigned to this message.
	Uuid string `json:"uuid,omitempty" path:"uuid,omitempty" url:"uuid,omitempty"`
	// Content Type header of the incoming message.
	ContentType string `json:"content_type,omitempty" path:"content_type,omitempty" url:"content_type,omitempty"`
	// HTTP Headers sent with this message.
	HttpHeaders interface{} `json:"http_headers,omitempty" path:"http_headers,omitempty" url:"http_headers,omitempty"`
	// Result of processing.
	ProcessingResult string `json:"processing_result,omitempty" path:"processing_result,omitempty" url:"processing_result,omitempty"`
	// Result of processing description.
	ProcessingResultDescription string `json:"processing_result_description,omitempty" path:"processing_result_description,omitempty" url:"processing_result_description,omitempty"`
	// AS2 Message Integrity Check.
	Mic string `json:"mic,omitempty" path:"mic,omitempty" url:"mic,omitempty"`
	// AS2 Message Integrity Check Algorithm Used.
	MicAlgo string `json:"mic_algo,omitempty" path:"mic_algo,omitempty" url:"mic_algo,omitempty"`
	// AS2 TO header of message.
	As2To string `json:"as2_to,omitempty" path:"as2_to,omitempty" url:"as2_to,omitempty"`
	// AS2 FROM header of message.
	As2From string `json:"as2_from,omitempty" path:"as2_from,omitempty" url:"as2_from,omitempty"`
	// AS2 Message Id.
	MessageId string `json:"message_id,omitempty" path:"message_id,omitempty" url:"message_id,omitempty"`
	// AS2 Subject Header.
	Subject string `json:"subject,omitempty" path:"subject,omitempty" url:"subject,omitempty"`
	// Date Header.
	Date string `json:"date,omitempty" path:"date,omitempty" url:"date,omitempty"`
	// Encrypted Payload Body Size.
	BodySize string `json:"body_size,omitempty" path:"body_size,omitempty" url:"body_size,omitempty"`
	// Filename of the file being received.
	AttachmentFilename string `json:"attachment_filename,omitempty" path:"attachment_filename,omitempty" url:"attachment_filename,omitempty"`
	// IP Address of the Sender.
	Ip string `json:"ip,omitempty" path:"ip,omitempty" url:"ip,omitempty"`
	// Message creation date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// HTTP Response Code sent for this message.
	HttpResponseCode string `json:"http_response_code,omitempty" path:"http_response_code,omitempty" url:"http_response_code,omitempty"`
	// HTTP Headers sent for this message.
	HttpResponseHeaders interface{} `json:"http_response_headers,omitempty" path:"http_response_headers,omitempty" url:"http_response_headers,omitempty"`
	// Incoming Message Recipient(the Client Cert used to encrypt this message)'s serial.
	RecipientSerial string `json:"recipient_serial,omitempty" path:"recipient_serial,omitempty" url:"recipient_serial,omitempty"`
	// Incoming Message Recipient(the Client Cert used to encrypt this message)'s serial in hex format.
	HexRecipientSerial string `json:"hex_recipient_serial,omitempty" path:"hex_recipient_serial,omitempty" url:"hex_recipient_serial,omitempty"`
	// Incoming Message Recipient(the Client Cert used to encrypt this message)'s issuer.
	RecipientIssuer string `json:"recipient_issuer,omitempty" path:"recipient_issuer,omitempty" url:"recipient_issuer,omitempty"`
	// Message body received?
	MessageReceived *bool `json:"message_received,omitempty" path:"message_received,omitempty" url:"message_received,omitempty"`
	// Message decrypted successfully?
	MessageDecrypted *bool `json:"message_decrypted,omitempty" path:"message_decrypted,omitempty" url:"message_decrypted,omitempty"`
	// Message signature verified?
	MessageSignatureVerified *bool `json:"message_signature_verified,omitempty" path:"message_signature_verified,omitempty" url:"message_signature_verified,omitempty"`
	// Message processed successfully?
	MessageProcessingSuccess *bool `json:"message_processing_success,omitempty" path:"message_processing_success,omitempty" url:"message_processing_success,omitempty"`
	// MDN returned?
	MessageMdnReturned *bool `json:"message_mdn_returned,omitempty" path:"message_mdn_returned,omitempty" url:"message_mdn_returned,omitempty"`
	// URL to download the encrypted signed smime that is to sent as AS2 body.
	EncryptedUri string `json:"encrypted_uri,omitempty" path:"encrypted_uri,omitempty" url:"encrypted_uri,omitempty"`
	// URL to download the file contents as smime with signature.
	SmimeSignedUri string `json:"smime_signed_uri,omitempty" path:"smime_signed_uri,omitempty" url:"smime_signed_uri,omitempty"`
	// URL to download the file contents encoded as smime.
	SmimeUri string `json:"smime_uri,omitempty" path:"smime_uri,omitempty" url:"smime_uri,omitempty"`
	// URL to download the original file contents.
	RawUri string `json:"raw_uri,omitempty" path:"raw_uri,omitempty" url:"raw_uri,omitempty"`
	// URL to download the http response body.
	MdnResponseUri string `json:"mdn_response_uri,omitempty" path:"mdn_response_uri,omitempty" url:"mdn_response_uri,omitempty"`
}

// Identifier returns the resource ID.
func (a As2IncomingMessage) Identifier() interface{} {
	return a.Id
}

// As2IncomingMessageCollection is a list of As2IncomingMessage resources.
type As2IncomingMessageCollection []As2IncomingMessage

// As2IncomingMessageListParams contains the request parameters for GET /as2_incoming_messages.
type As2IncomingMessageListParams struct {
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
func (a *As2IncomingMessage) UnmarshalJSON(data []byte) error {
	type as2IncomingMessage As2IncomingMessage
	var v as2IncomingMessage
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = As2IncomingMessage(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *As2IncomingMessageCollection) UnmarshalJSON(data []byte) error {
	type as2IncomingMessages As2IncomingMessageCollection
	var v as2IncomingMessages
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = As2IncomingMessageCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *As2IncomingMessageCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
