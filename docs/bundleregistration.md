# BundleRegistration

`files_sdk.BundleRegistration` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/bundleregistration`, named `bundle_registration`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	bundle_registration "github.com/Files-com/files-sdk-go/v3/bundleregistration"
)
```

## Operations

### List

`GET /bundle_registrations`

List Share Link Registrations

```go
func (c *Client) List(params files_sdk.BundleRegistrationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.BundleRegistrationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.BundleRegistration` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.BundleRegistrationListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `bundle_id` or `created_at`. |
| `BundleId` | `int64` | `bundle_id` | no | ID of the associated Bundle |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := bundle_registration.List(files_sdk.BundleRegistrationListParams{})
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

`files_sdk.BundleRegistration` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Code` | `string` | `code` | Registration cookie code |
| `Name` | `string` | `name` | Registrant name |
| `Company` | `string` | `company` | Registrant company name |
| `Email` | `string` | `email` | Registrant email address |
| `Ip` | `string` | `ip` | Registrant IP Address |
| `InboxCode` | `string` | `inbox_code` | InboxRegistration cookie code, if there is an associated InboxRegistration |
| `ClickwrapBody` | `string` | `clickwrap_body` | Clickwrap text that was shown to the registrant |
| `FormFieldSetId` | `int64` | `form_field_set_id` | Id of associated form field set |
| `FormFieldData` | `interface{}` | `form_field_data` | Data for form field set with form field ids as keys and user data as values |
| `BundleCode` | `string` | `bundle_code` | Bundle URL code |
| `BundleId` | `int64` | `bundle_id` | Id of associated bundle |
| `BundleRecipientId` | `int64` | `bundle_recipient_id` | Id of associated bundle recipient |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
| `CreatedAt` | `*time.Time` | `created_at` | Registration creation date/time |
