# InboundS3Log

`files_sdk.InboundS3Log` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/inbounds3log`, named `inbound_s3_log`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	inbound_s3_log "github.com/Files-com/files-sdk-go/v3/inbounds3log"
)
```

## Operations

### List

`GET /inbound_s3_logs`

List Inbound S3 Logs

```go
func (c *Client) List(params files_sdk.InboundS3LogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.InboundS3LogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.InboundS3Log` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.InboundS3LogListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `operation`, `status`, `path`, `client_ip` or `created_at`. Valid field combinations are `[ operation ]`, `[ status ]`, `[ path ]`, `[ client_ip ]`, `[ created_at ]`, `[ operation, status ]`, `[ operation, path ]`, `[ operation, client_ip ]`, `[ operation, created_at ]`, `[ status, path ]`, `[ status, client_ip ]`, `[ status, created_at ]`, `[ path, client_ip ]`, `[ path, created_at ]`, `[ client_ip, created_at ]`, `[ operation, status, path ]`, `[ operation, status, client_ip ]`, `[ operation, status, created_at ]`, `[ operation, path, client_ip ]`, `[ operation, path, created_at ]`, `[ operation, client_ip, created_at ]`, `[ status, path, client_ip ]`, `[ status, path, created_at ]`, `[ status, client_ip, created_at ]`, `[ path, client_ip, created_at ]`, `[ operation, status, path, client_ip ]`, `[ operation, status, path, created_at ]`, `[ operation, status, client_ip, created_at ]`, `[ operation, path, client_ip, created_at ]` or `[ status, path, client_ip, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ operation ]`, `[ status ]`, `[ path ]`, `[ client_ip ]`, `[ created_at ]`, `[ operation, status ]`, `[ operation, path ]`, `[ operation, client_ip ]`, `[ operation, created_at ]`, `[ status, path ]`, `[ status, client_ip ]`, `[ status, created_at ]`, `[ path, client_ip ]`, `[ path, created_at ]`, `[ client_ip, created_at ]`, `[ operation, status, path ]`, `[ operation, status, client_ip ]`, `[ operation, status, created_at ]`, `[ operation, path, client_ip ]`, `[ operation, path, created_at ]`, `[ operation, client_ip, created_at ]`, `[ status, path, client_ip ]`, `[ status, path, created_at ]`, `[ status, client_ip, created_at ]`, `[ path, client_ip, created_at ]`, `[ operation, status, path, client_ip ]`, `[ operation, status, path, created_at ]`, `[ operation, status, client_ip, created_at ]`, `[ operation, path, client_ip, created_at ]` or `[ status, path, client_ip, created_at ]`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ operation ]`, `[ status ]`, `[ path ]`, `[ client_ip ]`, `[ created_at ]`, `[ operation, status ]`, `[ operation, path ]`, `[ operation, client_ip ]`, `[ operation, created_at ]`, `[ status, path ]`, `[ status, client_ip ]`, `[ status, created_at ]`, `[ path, client_ip ]`, `[ path, created_at ]`, `[ client_ip, created_at ]`, `[ operation, status, path ]`, `[ operation, status, client_ip ]`, `[ operation, status, created_at ]`, `[ operation, path, client_ip ]`, `[ operation, path, created_at ]`, `[ operation, client_ip, created_at ]`, `[ status, path, client_ip ]`, `[ status, path, created_at ]`, `[ status, client_ip, created_at ]`, `[ path, client_ip, created_at ]`, `[ operation, status, path, client_ip ]`, `[ operation, status, path, created_at ]`, `[ operation, status, client_ip, created_at ]`, `[ operation, path, client_ip, created_at ]` or `[ status, path, client_ip, created_at ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `operation`, `status`, `path` or `client_ip`. Valid field combinations are `[ operation ]`, `[ status ]`, `[ path ]`, `[ client_ip ]`, `[ created_at ]`, `[ operation, status ]`, `[ operation, path ]`, `[ operation, client_ip ]`, `[ operation, created_at ]`, `[ status, path ]`, `[ status, client_ip ]`, `[ status, created_at ]`, `[ path, client_ip ]`, `[ path, created_at ]`, `[ client_ip, created_at ]`, `[ operation, status, path ]`, `[ operation, status, client_ip ]`, `[ operation, status, created_at ]`, `[ operation, path, client_ip ]`, `[ operation, path, created_at ]`, `[ operation, client_ip, created_at ]`, `[ status, path, client_ip ]`, `[ status, path, created_at ]`, `[ status, client_ip, created_at ]`, `[ path, client_ip, created_at ]`, `[ operation, status, path, client_ip ]`, `[ operation, status, path, created_at ]`, `[ operation, status, client_ip, created_at ]`, `[ operation, path, client_ip, created_at ]` or `[ status, path, client_ip, created_at ]`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ operation ]`, `[ status ]`, `[ path ]`, `[ client_ip ]`, `[ created_at ]`, `[ operation, status ]`, `[ operation, path ]`, `[ operation, client_ip ]`, `[ operation, created_at ]`, `[ status, path ]`, `[ status, client_ip ]`, `[ status, created_at ]`, `[ path, client_ip ]`, `[ path, created_at ]`, `[ client_ip, created_at ]`, `[ operation, status, path ]`, `[ operation, status, client_ip ]`, `[ operation, status, created_at ]`, `[ operation, path, client_ip ]`, `[ operation, path, created_at ]`, `[ operation, client_ip, created_at ]`, `[ status, path, client_ip ]`, `[ status, path, created_at ]`, `[ status, client_ip, created_at ]`, `[ path, client_ip, created_at ]`, `[ operation, status, path, client_ip ]`, `[ operation, status, path, created_at ]`, `[ operation, status, client_ip, created_at ]`, `[ operation, path, client_ip, created_at ]` or `[ status, path, client_ip, created_at ]`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ operation ]`, `[ status ]`, `[ path ]`, `[ client_ip ]`, `[ created_at ]`, `[ operation, status ]`, `[ operation, path ]`, `[ operation, client_ip ]`, `[ operation, created_at ]`, `[ status, path ]`, `[ status, client_ip ]`, `[ status, created_at ]`, `[ path, client_ip ]`, `[ path, created_at ]`, `[ client_ip, created_at ]`, `[ operation, status, path ]`, `[ operation, status, client_ip ]`, `[ operation, status, created_at ]`, `[ operation, path, client_ip ]`, `[ operation, path, created_at ]`, `[ operation, client_ip, created_at ]`, `[ status, path, client_ip ]`, `[ status, path, created_at ]`, `[ status, client_ip, created_at ]`, `[ path, client_ip, created_at ]`, `[ operation, status, path, client_ip ]`, `[ operation, status, path, created_at ]`, `[ operation, status, client_ip, created_at ]`, `[ operation, path, client_ip, created_at ]` or `[ status, path, client_ip, created_at ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := inbound_s3_log.List(files_sdk.InboundS3LogListParams{})
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

`files_sdk.InboundS3Log` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Path` | `string` | `path` | Request Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `ClientIp` | `string` | `client_ip` | Client IP |
| `Operation` | `string` | `operation` | S3 Operation Type |
| `Status` | `string` | `status` | HTTP Status Code |
| `AwsAccessKey` | `string` | `aws_access_key` | AWS Access Key ID |
| `ErrorMessage` | `string` | `error_message` | Error message, if applicable |
| `ErrorType` | `string` | `error_type` | Error type, if applicable |
| `DurationMs` | `int64` | `duration_ms` | Duration (in milliseconds) |
| `RequestId` | `string` | `request_id` | Request ID |
| `UserAgent` | `string` | `user_agent` | User Agent |
| `CreatedAt` | `*time.Time` | `created_at` | Start Time of Request |
