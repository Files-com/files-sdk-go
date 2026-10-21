# InboxUpload

`files_sdk.InboxUpload` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/inboxupload`, named `inbox_upload`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	inbox_upload "github.com/Files-com/files-sdk-go/v3/inboxupload"
)
```

## Operations

### List

`GET /inbox_uploads`

List Inbox Uploads

```go
func (c *Client) List(params files_sdk.InboxUploadListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.InboxUploadListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.InboxUpload` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.InboxUploadListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `folder_behavior_id` or `inbox_registration_id`. Valid field combinations are `[ folder_behavior_id, created_at ]`, `[ inbox_registration_id, created_at ]`, `[ folder_behavior_id, inbox_registration_id ]` or `[ folder_behavior_id, inbox_registration_id, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := inbox_upload.List(files_sdk.InboxUploadListParams{})
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

`files_sdk.InboxUpload` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `InboxRegistration` | `files_sdk.InboxRegistration` | `inbox_registration` |  |
| `Path` | `string` | `path` | Upload path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `CreatedAt` | `*time.Time` | `created_at` | Upload date/time |
