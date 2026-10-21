# SettingsChange

`files_sdk.SettingsChange` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/settingschange`, named `settings_change`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	settings_change "github.com/Files-com/files-sdk-go/v3/settingschange"
)
```

## Operations

### List

`GET /settings_changes`

List Settings Changes

```go
func (c *Client) List(params files_sdk.SettingsChangeListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.SettingsChangeListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.SettingsChange` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.SettingsChangeListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `api_key_id` or `user_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `api_key_id` and `user_id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := settings_change.List(files_sdk.SettingsChangeListParams{})
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

`files_sdk.SettingsChange` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `ApiKeyId` | `int64` | `api_key_id` | The API key id responsible for this change. |
| `Changes` | `[]string` | `changes` | Markdown-formatted change messages. |
| `CreatedAt` | `*time.Time` | `created_at` | The time this change was made. |
| `UserId` | `int64` | `user_id` | The user id responsible for this change. |
| `UserIsFilesSupport` | `*bool` | `user_is_files_support` | true if this change was performed by Files.com support. |
| `UserIsFromParentSite` | `*bool` | `user_is_from_parent_site` | true if this change was performed by a user on a parent site. |
| `Username` | `string` | `username` | The username of the user responsible for this change. |
