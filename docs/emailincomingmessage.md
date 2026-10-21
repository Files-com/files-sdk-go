# EmailIncomingMessage

`files_sdk.EmailIncomingMessage` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/emailincomingmessage`, named `email_incoming_message`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	email_incoming_message "github.com/Files-com/files-sdk-go/v3/emailincomingmessage"
)
```

## Operations

### List

`GET /email_incoming_messages`

List Email Incoming Messages

```go
func (c *Client) List(params files_sdk.EmailIncomingMessageListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.EmailIncomingMessageListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.EmailIncomingMessage` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.EmailIncomingMessageListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `sender`, `status` or `inbox_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `inbox_id`, `sender` or `status`. Valid field combinations are `[ inbox_id, created_at ]`, `[ sender, created_at ]`, `[ status, created_at ]`, `[ inbox_id, status ]`, `[ status, sender ]`, `[ inbox_id, status, created_at ]`, `[ status, sender, created_at ]`, `[ inbox_id, status, sender ]` or `[ inbox_id, status, sender, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `sender`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := email_incoming_message.List(files_sdk.EmailIncomingMessageListParams{})
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

`files_sdk.EmailIncomingMessage` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Id of the Email Incoming Message |
| `InboxId` | `int64` | `inbox_id` | Id of the Inbox associated with this message |
| `Sender` | `string` | `sender` | Sender of the email |
| `SenderName` | `string` | `sender_name` | Sender name |
| `Status` | `string` | `status` | Status of the message |
| `Body` | `string` | `body` | Body of the email |
| `Message` | `string` | `message` | Message describing the failure |
| `CreatedAt` | `*time.Time` | `created_at` | Message creation date/time |
| `InboxTitle` | `string` | `inbox_title` | Title of the Inbox associated with this message |
