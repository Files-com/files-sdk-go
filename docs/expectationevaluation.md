# ExpectationEvaluation

`files_sdk.ExpectationEvaluation` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/expectationevaluation`, named `expectation_evaluation`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	expectation_evaluation "github.com/Files-com/files-sdk-go/v3/expectationevaluation"
)
```

## Operations

### List

`GET /expectation_evaluations`

List Expectation Evaluations

```go
func (c *Client) List(params files_sdk.ExpectationEvaluationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ExpectationEvaluationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.ExpectationEvaluation` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ExpectationEvaluationListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `created_at` or `expectation_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `expectation_id` and `workspace_id`. Valid field combinations are `[ workspace_id, expectation_id ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := expectation_evaluation.List(files_sdk.ExpectationEvaluationListParams{})
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

`GET /expectation_evaluations/{id}`

Show Expectation Evaluation

```go
func (c *Client) Find(params files_sdk.ExpectationEvaluationFindParams, opts ...files_sdk.RequestResponseOption) (expectationEvaluation files_sdk.ExpectationEvaluation, err error)
func Find(params files_sdk.ExpectationEvaluationFindParams, opts ...files_sdk.RequestResponseOption) (expectationEvaluation files_sdk.ExpectationEvaluation, err error)
```

Parameters, `files_sdk.ExpectationEvaluationFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Expectation Evaluation ID. |

Example:

```go
result, err := expectation_evaluation.Find(files_sdk.ExpectationEvaluationFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.ExpectationEvaluation` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | ExpectationEvaluation ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
| `ExpectationId` | `int64` | `expectation_id` | Expectation ID. |
| `Status` | `string` | `status` | Evaluation status. |
| `OpenedVia` | `string` | `opened_via` | How the evaluation window was opened. |
| `OpenedAt` | `*time.Time` | `opened_at` | When the evaluation row was opened. |
| `WindowStartAt` | `*time.Time` | `window_start_at` | Logical start of the candidate window. |
| `WindowEndAt` | `*time.Time` | `window_end_at` | Actual candidate cutoff boundary for the window. |
| `DeadlineAt` | `*time.Time` | `deadline_at` | Logical due boundary for schedule-driven windows. |
| `LateAcceptanceCutoffAt` | `*time.Time` | `late_acceptance_cutoff_at` | Logical cutoff for late acceptance, when present. |
| `HardCloseAt` | `*time.Time` | `hard_close_at` | Hard stop after which the window may not remain open. |
| `ClosedAt` | `*time.Time` | `closed_at` | When the evaluation row was finalized. |
| `MatchedFiles` | `[]map[string]interface{}` | `matched_files` | Captured evidence for files that matched the window. |
| `MissingFiles` | `[]map[string]interface{}` | `missing_files` | Captured evidence for required files that were missing. |
| `CriteriaErrors` | `[]string` | `criteria_errors` | Captured criteria failures for the window, including file-scoped FTS errors and returned validation details. |
| `Summary` | `interface{}` | `summary` | Compact evaluator summary payload. |
| `CreatedAt` | `*time.Time` | `created_at` | Creation time. |
| `UpdatedAt` | `*time.Time` | `updated_at` | Last update time. |
