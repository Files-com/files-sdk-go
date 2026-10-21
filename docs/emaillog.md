# EmailLog

`files_sdk.EmailLog` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/emaillog`, named `email_log`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	email_log "github.com/Files-com/files-sdk-go/v3/emaillog"
)
```

## Operations

### List

`GET /email_logs`

List Email Logs

```go
func (c *Client) List(params files_sdk.EmailLogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.EmailLogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.EmailLog` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.EmailLogListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `status` and `created_at`. Valid field combinations are `[ status ]`, `[ created_at ]` or `[ status, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ status ]`, `[ created_at ]` or `[ status, created_at ]`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ status ]`, `[ created_at ]` or `[ status, created_at ]`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ status ]`, `[ created_at ]` or `[ status, created_at ]`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ status ]`, `[ created_at ]` or `[ status, created_at ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := email_log.List(files_sdk.EmailLogListParams{})
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

`files_sdk.EmailLog` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Timestamp` | `*time.Time` | `timestamp` | Start Time of Action. Deprecrated: Use created_at. |
| `Message` | `string` | `message` | Log Message |
| `Status` | `string` | `status` | Status of E-Mail delivery |
| `Subject` | `string` | `subject` | Subject line of E-Mail |
| `To` | `string` | `to` | To field of E-Mail |
| `Cc` | `string` | `cc` | CC field of E-Mail |
| `DeliveryMethod` | `string` | `delivery_method` | How was the email delivered? `customer_smtp` or `files.com` |
| `SmtpHostname` | `string` | `smtp_hostname` | Customer SMTP Hostname used. |
| `SmtpIp` | `string` | `smtp_ip` | Customer SMTP IP address as resolved for use (useful for troubleshooting DNS issues with customer SMTP). |
| `CreatedAt` | `*time.Time` | `created_at` | Start Time of Action |
