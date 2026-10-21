# As2OutgoingMessage

`files_sdk.As2OutgoingMessage` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/as2outgoingmessage`, named `as2_outgoing_message`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	as2_outgoing_message "github.com/Files-com/files-sdk-go/v3/as2outgoingmessage"
)
```

## Operations

### List

`GET /as2_outgoing_messages`

List AS2 Outgoing Messages

```go
func (c *Client) List(params files_sdk.As2OutgoingMessageListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.As2OutgoingMessageListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.As2OutgoingMessage` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.As2OutgoingMessageListParams`:

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
it, err := as2_outgoing_message.List(files_sdk.As2OutgoingMessageListParams{})
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

## Model

`files_sdk.As2OutgoingMessage` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Id of the AS2 Partner. |
| `WorkspaceId` | `int64` | `workspace_id` | ID of the Workspace associated with this AS2 Outgoing Message. |
| `As2PartnerId` | `int64` | `as2_partner_id` | Id of the AS2 Partner associated with this message. |
| `As2StationId` | `int64` | `as2_station_id` | Id of the AS2 Station associated with this message. |
| `Uuid` | `string` | `uuid` | UUID assigned to this message. |
| `HttpHeaders` | `interface{}` | `http_headers` | HTTP Headers sent with this message. |
| `ProcessingResult` | `string` | `processing_result` | Result of processing. |
| `ProcessingResultDescription` | `string` | `processing_result_description` | Result of processing description. |
| `Mic` | `string` | `mic` | AS2 Message Integrity Check SHA1 |
| `MicSha256` | `string` | `mic_sha_256` | AS2 Message Integrity Check SHA256 |
| `As2To` | `string` | `as2_to` | AS2 TO |
| `As2From` | `string` | `as2_from` | AS2 FROM |
| `As2Subject` | `string` | `as2_subject` | AS2 Subject |
| `Date` | `string` | `date` | Date Header |
| `MessageId` | `string` | `message_id` | AS2 Message Id |
| `BodySize` | `string` | `body_size` | Encrypted Payload Body Size |
| `AttachmentFilename` | `string` | `attachment_filename` | Filename of the file being sent. |
| `CreatedAt` | `*time.Time` | `created_at` | Message creation date/time |
| `HttpResponseCode` | `string` | `http_response_code` | HTTP Response Code received for this message |
| `HttpResponseHeaders` | `interface{}` | `http_response_headers` | HTTP Headers received for this message. |
| `HttpTransmissionDuration` | `float64` | `http_transmission_duration` | HTTP transmission duration in seceonds |
| `MdnReceived` | `*bool` | `mdn_received` | Did the partner give a response body? |
| `MdnValid` | `*bool` | `mdn_valid` | Is the response in MDN format? |
| `MdnSignatureVerified` | `*bool` | `mdn_signature_verified` | MDN signature verified? |
| `MdnMessageIdMatched` | `*bool` | `mdn_message_id_matched` | MDN message id matched? |
| `MdnMicMatched` | `*bool` | `mdn_mic_matched` | MDN MIC matched? |
| `MdnProcessingSuccess` | `*bool` | `mdn_processing_success` | MDN disposition indicate a successful processing? |
| `RawUri` | `string` | `raw_uri` | URL to download the original file contents |
| `SmimeUri` | `string` | `smime_uri` | URL to download the file contents encoded as smime |
| `SmimeSignedUri` | `string` | `smime_signed_uri` | URL to download the file contents as smime with signature |
| `EncryptedUri` | `string` | `encrypted_uri` | URL to download the encrypted signed smime that is to sent as AS2 body |
| `MdnResponseUri` | `string` | `mdn_response_uri` | URL to download the http response body |
