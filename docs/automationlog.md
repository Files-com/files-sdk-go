# AutomationLog

`files_sdk.AutomationLog` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/automationlog`, named `automation_log`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	automation_log "github.com/Files-com/files-sdk-go/v3/automationlog"
)
```

## Operations

### List

`GET /automation_logs`

List Automation Logs

```go
func (c *Client) List(params files_sdk.AutomationLogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.AutomationLogListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.AutomationLog` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.AutomationLogListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` or `created_at`. Valid field combinations are every combination of one or more of the 13 fields `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` and `created_at`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `source_automation_version` and `created_at`. Valid field combinations are every combination of one or more of the 13 fields `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` and `created_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `source_automation_version` and `created_at`. Valid field combinations are every combination of one or more of the 13 fields `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` and `created_at`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `correlation_id`, `error_type`, `node_id`, `node_type`, `output` or `path`. Valid field combinations are every combination of one or more of the 13 fields `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` and `created_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `source_automation_version` and `created_at`. Valid field combinations are every combination of one or more of the 13 fields `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` and `created_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `source_automation_version` and `created_at`. Valid field combinations are every combination of one or more of the 13 fields `automation_id`, `automation_run_id`, `correlation_id`, `error_type`, `node_id`, `node_type`, `operation`, `output`, `path`, `source_automation_id`, `source_automation_version`, `status` and `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := automation_log.List(files_sdk.AutomationLogListParams{})
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

`files_sdk.AutomationLog` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Timestamp` | `*time.Time` | `timestamp` | Start Time of Action. Deprecrated: Use created_at. |
| `AutomationId` | `int64` | `automation_id` | Automation ID |
| `AutomationRunId` | `int64` | `automation_run_id` | Automation Run ID |
| `CorrelationId` | `string` | `correlation_id` | Logical item correlation ID |
| `DestPath` | `string` | `dest_path` | Destination path, for moves and copies |
| `ErrorType` | `string` | `error_type` | Typed error, if applicable |
| `SourceAutomationId` | `int64` | `source_automation_id` | ID of the Automation whose definition owns this node |
| `SourceAutomationVersion` | `int64` | `source_automation_version` | Pinned version of the Automation whose definition owns this node |
| `Message` | `string` | `message` | Message |
| `NodeId` | `string` | `node_id` | Node ID |
| `NodeType` | `string` | `node_type` | Node type |
| `Operation` | `string` | `operation` | Operation type |
| `Output` | `string` | `output` | Node outlet used by the item |
| `Path` | `string` | `path` | File path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `Status` | `string` | `status` | Status |
| `CreatedAt` | `*time.Time` | `created_at` | Start Time of Action |
