# InboxRecipient

`files_sdk.InboxRecipient` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/inboxrecipient`, named `inbox_recipient`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	inbox_recipient "github.com/Files-com/files-sdk-go/v3/inboxrecipient"
)
```

## Operations

### List

`GET /inbox_recipients`

List Inbox Recipients

```go
func (c *Client) List(params files_sdk.InboxRecipientListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.InboxRecipientListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.InboxRecipient` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.InboxRecipientListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `has_registrations`. |
| `InboxId` | `int64` | `inbox_id` | yes | List recipients for the inbox with this ID. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := inbox_recipient.List(files_sdk.InboxRecipientListParams{
	InboxId: 1,
})
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

### Create

`POST /inbox_recipients`

Create Inbox Recipient

```go
func (c *Client) Create(params files_sdk.InboxRecipientCreateParams, opts ...files_sdk.RequestResponseOption) (inboxRecipient files_sdk.InboxRecipient, err error)
func Create(params files_sdk.InboxRecipientCreateParams, opts ...files_sdk.RequestResponseOption) (inboxRecipient files_sdk.InboxRecipient, err error)
```

Parameters, `files_sdk.InboxRecipientCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `InboxId` | `int64` | `inbox_id` | yes | Inbox to share. |
| `Recipient` | `string` | `recipient` | yes | Email address to share this inbox with. |
| `Name` | `string` | `name` | no | Name of recipient. |
| `Company` | `string` | `company` | no | Company of recipient. |
| `Note` | `string` | `note` | no | Note to include in email. |
| `ShareAfterCreate` | `*bool` | `share_after_create` | no | Set to true to share the link with the recipient upon creation. |

Example:

```go
result, err := inbox_recipient.Create(files_sdk.InboxRecipientCreateParams{
	InboxId:   1,
	Recipient: "johndoe@gmail.com",
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.InboxRecipient` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Company` | `string` | `company` | The recipient's company. |
| `Name` | `string` | `name` | The recipient's name. |
| `Note` | `string` | `note` | A note sent to the recipient with the inbox. |
| `Recipient` | `string` | `recipient` | The recipient's email address. |
| `SentAt` | `*time.Time` | `sent_at` | When the Inbox was shared with this recipient. |
| `InboxId` | `int64` | `inbox_id` | Inbox to share. |
| `ShareAfterCreate` | `*bool` | `share_after_create` | Set to true to share the link with the recipient upon creation. |
