# As2IncomingMessage

`files_sdk.As2IncomingMessage` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/as2incomingmessage`, named `as2_incoming_message`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	as2_incoming_message "github.com/Files-com/files-sdk-go/v3/as2incomingmessage"
)
```

## Operations

### List

`GET /as2_incoming_messages`

List AS2 Incoming Messages

```go
func (c *Client) List(params files_sdk.As2IncomingMessageListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.As2IncomingMessageListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.As2IncomingMessage` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.As2IncomingMessageListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `created_at` or `as2_partner_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `as2_station_id`, `workspace_id` or `as2_partner_id`. Valid field combinations are `[ as2_station_id, created_at ]`, `[ workspace_id, created_at ]`, `[ as2_partner_id, created_at ]`, `[ workspace_id, as2_station_id ]`, `[ workspace_id, as2_partner_id ]`, `[ workspace_id, as2_station_id, created_at ]` or `[ workspace_id, as2_partner_id, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := as2_incoming_message.List(files_sdk.As2IncomingMessageListParams{})
if err != nil {
	return err
}
for item, err := range it.All() {
	if err != nil {
		return err
	}
	fmt.Println(item)
}
```

### Replay

`POST /as2_incoming_messages/{id}/replay`

Replay AS2 Incoming Message

```go
func (c *Client) Replay(params files_sdk.As2IncomingMessageReplayParams, opts ...files_sdk.RequestResponseOption) (as2IncomingMessage files_sdk.As2IncomingMessage, err error)
func Replay(params files_sdk.As2IncomingMessageReplayParams, opts ...files_sdk.RequestResponseOption) (as2IncomingMessage files_sdk.As2IncomingMessage, err error)
```

Parameters, `files_sdk.As2IncomingMessageReplayParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | As2 Incoming Message ID. |

Example:

```go
result, err := as2_incoming_message.Replay(files_sdk.As2IncomingMessageReplayParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.As2IncomingMessage` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | ID of this AS2 incoming message record. |
| `ReplayedFromMessageId` | `int64` | `replayed_from_message_id` | ID of the incoming message record this message was replayed from. Empty for messages received directly from a trading partner. |
| `WorkspaceId` | `int64` | `workspace_id` | ID of the Workspace associated with this AS2 Incoming Message. |
| `As2PartnerId` | `int64` | `as2_partner_id` | Id of the AS2 Partner associated with this message. |
| `As2StationId` | `int64` | `as2_station_id` | Id of the AS2 Station associated with this message. |
| `Uuid` | `string` | `uuid` | UUID assigned to this message. |
| `ContentType` | `string` | `content_type` | Content Type header of the incoming message. |
| `HttpHeaders` | `interface{}` | `http_headers` | HTTP Headers sent with this message. |
| `ProcessingResult` | `string` | `processing_result` | Result of processing. |
| `ProcessingResultDescription` | `string` | `processing_result_description` | Result of processing description. |
| `Mic` | `string` | `mic` | AS2 Message Integrity Check |
| `MicAlgo` | `string` | `mic_algo` | AS2 Message Integrity Check Algorithm Used |
| `As2To` | `string` | `as2_to` | AS2 TO header of message |
| `As2From` | `string` | `as2_from` | AS2 FROM header of message |
| `MessageId` | `string` | `message_id` | AS2 Message Id |
| `Subject` | `string` | `subject` | AS2 Subject Header |
| `Date` | `string` | `date` | Date Header |
| `BodySize` | `string` | `body_size` | Encrypted Payload Body Size |
| `AttachmentFilename` | `string` | `attachment_filename` | Filename of the file being received. |
| `Ip` | `string` | `ip` | IP Address of the Sender |
| `CreatedAt` | `*time.Time` | `created_at` | Message creation date/time |
| `HttpResponseCode` | `string` | `http_response_code` | HTTP Response Code sent for this message |
| `HttpResponseHeaders` | `interface{}` | `http_response_headers` | HTTP Headers sent for this message. |
| `RecipientSerial` | `string` | `recipient_serial` | Incoming Message Recipient(the Client Cert used to encrypt this message)'s serial |
| `HexRecipientSerial` | `string` | `hex_recipient_serial` | Incoming Message Recipient(the Client Cert used to encrypt this message)'s serial in hex format. |
| `RecipientIssuer` | `string` | `recipient_issuer` | Incoming Message Recipient(the Client Cert used to encrypt this message)'s issuer |
| `MessageReceived` | `*bool` | `message_received` | Message body received? |
| `MessageDecrypted` | `*bool` | `message_decrypted` | Message decrypted successfully? |
| `MessageSignatureVerified` | `*bool` | `message_signature_verified` | Message signature verified? |
| `MessageProcessingSuccess` | `*bool` | `message_processing_success` | Message processed successfully? |
| `MessageMdnReturned` | `*bool` | `message_mdn_returned` | MDN returned? |
| `EncryptedUri` | `string` | `encrypted_uri` | URL to download the encrypted signed smime that is to sent as AS2 body |
| `SmimeSignedUri` | `string` | `smime_signed_uri` | URL to download the file contents as smime with signature |
| `SmimeUri` | `string` | `smime_uri` | URL to download the file contents encoded as smime |
| `RawUri` | `string` | `raw_uri` | URL to download the original file contents |
| `MdnResponseUri` | `string` | `mdn_response_uri` | URL to download the http response body |
