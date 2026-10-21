# BundleDownload

`files_sdk.BundleDownload` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/bundledownload`, named `bundle_download`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	bundle_download "github.com/Files-com/files-sdk-go/v3/bundledownload"
)
```

## Operations

### List

`GET /bundle_downloads`

List Share Link Downloads

```go
func (c *Client) List(params files_sdk.BundleDownloadListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.BundleDownloadListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.BundleDownload` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.BundleDownloadListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `created_at`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |
| `BundleId` | `int64` | `bundle_id` | no | Bundle ID |
| `BundleRegistrationId` | `int64` | `bundle_registration_id` | no | BundleRegistration ID |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := bundle_download.List(files_sdk.BundleDownloadListParams{})
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

`files_sdk.BundleDownload` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `BundleRegistration` | `files_sdk.BundleRegistration` | `bundle_registration` |  |
| `DownloadMethod` | `string` | `download_method` | Download method (file or full_zip) |
| `Path` | `string` | `path` | Download path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
| `CreatedAt` | `*time.Time` | `created_at` | Download date/time |
