# SharedWorkflow

`files_sdk.SharedWorkflow` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/sharedworkflow`, named `shared_workflow`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	shared_workflow "github.com/Files-com/files-sdk-go/v3/sharedworkflow"
)
```

## Operations

### List

`GET /shared_workflows`

List Shared Workflows

```go
func (c *Client) List(params files_sdk.SharedWorkflowListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.SharedWorkflowListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.SharedWorkflow` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.SharedWorkflowListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := shared_workflow.List(files_sdk.SharedWorkflowListParams{})
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

`GET /shared_workflows/{id}`

Show Shared Workflow

```go
func (c *Client) Find(params files_sdk.SharedWorkflowFindParams, opts ...files_sdk.RequestResponseOption) (sharedWorkflow files_sdk.SharedWorkflow, err error)
func Find(params files_sdk.SharedWorkflowFindParams, opts ...files_sdk.RequestResponseOption) (sharedWorkflow files_sdk.SharedWorkflow, err error)
```

Parameters, `files_sdk.SharedWorkflowFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Shared Workflow ID. |

Example:

```go
result, err := shared_workflow.Find(files_sdk.SharedWorkflowFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### FindRun

`GET /shared_workflows/{id}/runs/{run_id}`

Show a Shared Workflow run started by the recipient

```go
func (c *Client) FindRun(params files_sdk.SharedWorkflowFindRunParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error)
func FindRun(params files_sdk.SharedWorkflowFindRunParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error)
```

Parameters, `files_sdk.SharedWorkflowFindRunParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Shared Workflow ID. |
| `RunId` | `int64` | `run_id`, in the path | yes | ID of the run returned by the Shared Workflow submission. |

Example:

```go
result, err := shared_workflow.FindRun(files_sdk.SharedWorkflowFindRunParams{
	Id:    1,
	RunId: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### RunWorkflow

`POST /shared_workflows/{id}/run_workflow`

Run a Shared Workflow

```go
func (c *Client) RunWorkflow(params files_sdk.SharedWorkflowRunWorkflowParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error)
func RunWorkflow(params files_sdk.SharedWorkflowRunWorkflowParams, opts ...files_sdk.RequestResponseOption) (automationRun files_sdk.AutomationRun, err error)
```

Parameters, `files_sdk.SharedWorkflowRunWorkflowParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Shared Workflow ID. |
| `SubmissionRevision` | `string` | `submission_revision` | yes | Revision returned when opening the Shared Workflow. |
| `SubmissionId` | `string` | `submission_id` | yes | Client-generated UUID used to retry the same submission without starting another run. |
| `FilePath` | `string` | `file_path` | yes | Path of a file the recipient can read, including a file uploaded to a Snapshot created for this workflow. |

Example:

```go
result, err := shared_workflow.RunWorkflow(files_sdk.SharedWorkflowRunWorkflowParams{
	Id:                 1,
	SubmissionRevision: "submission_revision",
	SubmissionId:       "submission_id",
	FilePath:           "file_path",
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.SharedWorkflow` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Shared Workflow ID. |
| `Name` | `string` | `name` | Name displayed to recipients. |
| `Description` | `string` | `description` | Description displayed to recipients. |
| `SubmissionRevision` | `string` | `submission_revision` | Revision required when submitting a Shared Workflow run. |
| `FormFields` | `[]files_sdk.FormField` | `form_fields` | Current custom questions in display order. |
