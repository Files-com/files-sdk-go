# FileMigration

`files_sdk.FileMigration` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/filemigration`, named `file_migration`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	file_migration "github.com/Files-com/files-sdk-go/v3/filemigration"
)
```

## Operations

### Find

`GET /file_migrations/{id}`

Show File Migration

```go
func (c *Client) Find(params files_sdk.FileMigrationFindParams, opts ...files_sdk.RequestResponseOption) (fileMigration files_sdk.FileMigration, err error)
func Find(params files_sdk.FileMigrationFindParams, opts ...files_sdk.RequestResponseOption) (fileMigration files_sdk.FileMigration, err error)
```

Parameters, `files_sdk.FileMigrationFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | File Migration ID. |

Example:

```go
result, err := file_migration.Find(files_sdk.FileMigrationFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Other functions and methods

The package also provides these functions and methods:

```go
// Wait polls the migration once per second until it completes, fails, or a request
// returns an error. status is called after each successful poll and must be non-nil.
// WithContext controls cancellation. A failed migration is returned with nil error
// when the status request itself succeeds; inspect the returned Status.
func (c *Client) Wait(fileAction files_sdk.FileAction, status func(files_sdk.FileMigration), opts ...files_sdk.RequestResponseOption) (files_sdk.FileMigration, error)

// LogIterator returns an iterator over the migration log using ctx for cancellation.
func (c *Client) LogIterator(ctx context.Context, f files_sdk.FileMigration) files_sdk.IterI
```

## Model

`files_sdk.FileMigration` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | File migration ID |
| `Path` | `string` | `path` | Source path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `DestPath` | `string` | `dest_path` | Destination path |
| `FailureMessage` | `string` | `failure_message` | Reason for the failure, if applicable. |
| `FilesMoved` | `int64` | `files_moved` | Number of files processed |
| `FilesTotal` | `int64` | `files_total` |  |
| `Operation` | `string` | `operation` | The type of operation |
| `Region` | `string` | `region` | Region |
| `Status` | `string` | `status` | Status |
| `LogUrl` | `string` | `log_url` | Link to download the log file for this migration. |
