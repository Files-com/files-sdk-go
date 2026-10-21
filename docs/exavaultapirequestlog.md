# ExavaultApiRequestLog

`files_sdk.ExavaultApiRequestLog` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/exavaultapirequestlog`, named `exavault_api_request_log`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	exavault_api_request_log "github.com/Files-com/files-sdk-go/v3/exavaultapirequestlog"
)
```

## Operations

### List

`GET /exavault_api_request_logs`

List Exavault API Request Logs

```go
func (c *Client) List(params files_sdk.ExavaultApiRequestLogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ExavaultApiRequestLogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.ExavaultApiRequestLog` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ExavaultApiRequestLogListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `request_ip`, `request_method`, `success` or `created_at`. Valid field combinations are `[ request_ip ]`, `[ request_method ]`, `[ success ]`, `[ created_at ]`, `[ request_ip, request_method ]`, `[ request_ip, success ]`, `[ request_ip, created_at ]`, `[ request_method, success ]`, `[ request_method, created_at ]`, `[ success, created_at ]`, `[ request_ip, request_method, success ]`, `[ request_ip, request_method, created_at ]`, `[ request_ip, success, created_at ]`, `[ request_method, success, created_at ]` or `[ request_ip, request_method, success, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ request_ip ]`, `[ request_method ]`, `[ success ]`, `[ created_at ]`, `[ request_ip, request_method ]`, `[ request_ip, success ]`, `[ request_ip, created_at ]`, `[ request_method, success ]`, `[ request_method, created_at ]`, `[ success, created_at ]`, `[ request_ip, request_method, success ]`, `[ request_ip, request_method, created_at ]`, `[ request_ip, success, created_at ]`, `[ request_method, success, created_at ]` or `[ request_ip, request_method, success, created_at ]`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ request_ip ]`, `[ request_method ]`, `[ success ]`, `[ created_at ]`, `[ request_ip, request_method ]`, `[ request_ip, success ]`, `[ request_ip, created_at ]`, `[ request_method, success ]`, `[ request_method, created_at ]`, `[ success, created_at ]`, `[ request_ip, request_method, success ]`, `[ request_ip, request_method, created_at ]`, `[ request_ip, success, created_at ]`, `[ request_method, success, created_at ]` or `[ request_ip, request_method, success, created_at ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `request_ip` and `request_method`. Valid field combinations are `[ request_ip ]`, `[ request_method ]`, `[ success ]`, `[ created_at ]`, `[ request_ip, request_method ]`, `[ request_ip, success ]`, `[ request_ip, created_at ]`, `[ request_method, success ]`, `[ request_method, created_at ]`, `[ success, created_at ]`, `[ request_ip, request_method, success ]`, `[ request_ip, request_method, created_at ]`, `[ request_ip, success, created_at ]`, `[ request_method, success, created_at ]` or `[ request_ip, request_method, success, created_at ]`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. Valid field combinations are `[ request_ip ]`, `[ request_method ]`, `[ success ]`, `[ created_at ]`, `[ request_ip, request_method ]`, `[ request_ip, success ]`, `[ request_ip, created_at ]`, `[ request_method, success ]`, `[ request_method, created_at ]`, `[ success, created_at ]`, `[ request_ip, request_method, success ]`, `[ request_ip, request_method, created_at ]`, `[ request_ip, success, created_at ]`, `[ request_method, success, created_at ]` or `[ request_ip, request_method, success, created_at ]`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. Valid field combinations are `[ request_ip ]`, `[ request_method ]`, `[ success ]`, `[ created_at ]`, `[ request_ip, request_method ]`, `[ request_ip, success ]`, `[ request_ip, created_at ]`, `[ request_method, success ]`, `[ request_method, created_at ]`, `[ success, created_at ]`, `[ request_ip, request_method, success ]`, `[ request_ip, request_method, created_at ]`, `[ request_ip, success, created_at ]`, `[ request_method, success, created_at ]` or `[ request_ip, request_method, success, created_at ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := exavault_api_request_log.List(files_sdk.ExavaultApiRequestLogListParams{})
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

`files_sdk.ExavaultApiRequestLog` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Timestamp` | `*time.Time` | `timestamp` | Start Time of Action. Deprecrated: Use created_at. |
| `Endpoint` | `string` | `endpoint` | Name of API Endpoint |
| `Version` | `float64` | `version` | Exavault API Version |
| `RequestIp` | `string` | `request_ip` | IP of requesting client |
| `RequestMethod` | `string` | `request_method` | HTTP Method |
| `ErrorType` | `string` | `error_type` | Error type, if applicable |
| `ErrorMessage` | `string` | `error_message` | Error message, if applicable |
| `UserAgent` | `string` | `user_agent` | User-Agent |
| `ResponseCode` | `int64` | `response_code` | HTTP Response Code |
| `Success` | `*bool` | `success` | `false` if HTTP Response Code is 4xx or 5xx |
| `DurationMs` | `int64` | `duration_ms` | Duration (in milliseconds) |
| `CreatedAt` | `*time.Time` | `created_at` | Start Time of Action |
