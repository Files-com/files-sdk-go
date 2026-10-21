# AutomationRun

`files_sdk.AutomationRun` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/automationrun`, named `automation_run`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	automation_run "github.com/Files-com/files-sdk-go/v3/automationrun"
)
```

## Operations

### List

`GET /automation_runs`

List Automation Runs

```go
func (c *Client) List(params files_sdk.AutomationRunListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.AutomationRunListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.AutomationRun` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.AutomationRunListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `automation_id`, `created_at` or `status`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `status`, `workspace_id` or `automation_id`. Valid field combinations are `[ workspace_id, status ]`, `[ automation_id, status ]`, `[ workspace_id, automation_id ]` or `[ workspace_id, automation_id, status ]`. |
| `AutomationId` | `int64` | `automation_id` | yes | ID of the associated Automation. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := automation_run.List(files_sdk.AutomationRunListParams{
	AutomationId: 1,
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

### Find

`GET /automation_runs/{id}`

Show Automation Run

```go
func (c *Client) Find(params files_sdk.AutomationRunFindParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error)
func Find(params files_sdk.AutomationRunFindParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error)
```

Parameters, `files_sdk.AutomationRunFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Automation Run ID. |

Example:

```go
result, err := automation_run.Find(files_sdk.AutomationRunFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### FindNode

`GET /automation_runs/{id}/node`

Show Automation Run Node

```go
func (c *Client) FindNode(params files_sdk.AutomationRunFindNodeParams, opts ...files_sdk.RequestResponseOption) (automationExecutionNode files_sdk.AutomationExecutionNode, err error)
func FindNode(params files_sdk.AutomationRunFindNodeParams, opts ...files_sdk.RequestResponseOption) (automationExecutionNode files_sdk.AutomationExecutionNode, err error)
```

Parameters, `files_sdk.AutomationRunFindNodeParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Automation Run ID. |
| `NodeId` | `string` | `node_id` | yes | Node ID from the pinned Automation definition. |

Example:

```go
result, err := automation_run.FindNode(files_sdk.AutomationRunFindNodeParams{
	Id:     1,
	NodeId: "node_id",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Cancel

`POST /automation_runs/{id}/cancel`

Cancel Automation Run

```go
func (c *Client) Cancel(params files_sdk.AutomationRunCancelParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error)
func Cancel(params files_sdk.AutomationRunCancelParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error)
```

Parameters, `files_sdk.AutomationRunCancelParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Automation Run ID. |

Example:

```go
result, err := automation_run.Cancel(files_sdk.AutomationRunCancelParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Rerun

`POST /automation_runs/{id}/rerun`

Re-run Automation from Node

```go
func (c *Client) Rerun(params files_sdk.AutomationRunRerunParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error)
func Rerun(params files_sdk.AutomationRunRerunParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error)
```

Parameters, `files_sdk.AutomationRunRerunParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Automation Run ID. |
| `NodeId` | `string` | `node_id` | yes | Node ID at which execution should resume. |

Example:

```go
result, err := automation_run.Rerun(files_sdk.AutomationRunRerunParams{
	Id:     1,
	NodeId: "node_id",
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.AutomationRun` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | ID. |
| `AutomationId` | `int64` | `automation_id` | ID of the associated Automation. |
| `AutomationVersionId` | `int64` | `automation_version_id` | ID of the immutable Automation version pinned by this run. |
| `Version` | `int64` | `version` | Pinned Automation v2 definition version. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. |
| `CancelRequestedAt` | `*time.Time` | `cancel_requested_at` | Date/time at which cancellation was requested. |
| `CompletedAt` | `*time.Time` | `completed_at` | Automation run completion/failure date/time. |
| `CreatedAt` | `*time.Time` | `created_at` | Automation run start date/time. |
| `RetryAt` | `*time.Time` | `retry_at` | If set, this automation will be retried at this date/time due to `failure` or `partial_failure`. |
| `RetriedAt` | `*time.Time` | `retried_at` | If set, this Automation run was retried due to `failure` or `partial_failure`. |
| `RetriedInRunId` | `int64` | `retried_in_run_id` | ID of the run that is or will be retrying this run. |
| `RetryOfRunId` | `int64` | `retry_of_run_id` | ID of the original run that this run is retrying. |
| `RerunOfRunId` | `int64` | `rerun_of_run_id` | ID of the run whose persisted node outputs this run reused. |
| `RerunFromNodeId` | `string` | `rerun_from_node_id` | Node at which this run resumed execution. |
| `Runtime` | `float64` | `runtime` | Automation run runtime. |
| `Status` | `string` | `status` | The status of the AutomationRun. One of `queued`, `running`, `success`, `partial_failure`, `failure`, `skipped`, or `canceled`. |
| `SuccessfulOperations` | `int64` | `successful_operations` | Count of successful operations. |
| `FailedOperations` | `int64` | `failed_operations` | Count of failed operations. |
| `Definition` | `interface{}` | `definition` | Automation definition snapshot pinned by this run. For performance reasons, this is not provided when listing Automation runs. |
| `NodeStates` | `interface{}` | `node_states` | Status and execution stage for each node in this run. For performance reasons, this is not provided when listing Automation runs. |
| `ExecutionNodes` | `[]files_sdk.AutomationExecutionNode` | `execution_nodes` | Execution status, timing, and bounded output summaries for each node. For performance reasons, this is not provided when listing Automation runs. |
| `JournalUrl` | `string` | `journal_url` | Link to the run journal artifact. |
| `StatusMessagesUrl` | `string` | `status_messages_url` | Link to status messages log file. |
