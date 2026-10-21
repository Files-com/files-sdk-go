# InboxRegistration

`files_sdk.InboxRegistration` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/inboxregistration`, named `inbox_registration`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	inbox_registration "github.com/Files-com/files-sdk-go/v3/inboxregistration"
)
```

## Operations

### List

`GET /inbox_registrations`

List Inbox Registrations

```go
func (c *Client) List(params files_sdk.InboxRegistrationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.InboxRegistrationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.InboxRegistration` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.InboxRegistrationListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `FolderBehaviorId` | `int64` | `folder_behavior_id` | no | ID of the associated Inbox. This is required if the user is not a site admin. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := inbox_registration.List(files_sdk.InboxRegistrationListParams{})
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

`files_sdk.InboxRegistration` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Code` | `string` | `code` | Registration cookie code |
| `Name` | `string` | `name` | Registrant name |
| `Company` | `string` | `company` | Registrant company name |
| `Email` | `string` | `email` | Registrant email address |
| `Ip` | `string` | `ip` | Registrant IP Address |
| `ClickwrapBody` | `string` | `clickwrap_body` | Clickwrap text that was shown to the registrant |
| `FormFieldSetId` | `int64` | `form_field_set_id` | Id of associated form field set |
| `FormFieldData` | `interface{}` | `form_field_data` | Data for form field set with form field ids as keys and user data as values |
| `InboxId` | `int64` | `inbox_id` | Id of associated inbox |
| `InboxRecipientId` | `int64` | `inbox_recipient_id` | Id of associated inbox recipient |
| `InboxTitle` | `string` | `inbox_title` | Title of associated inbox |
| `CreatedAt` | `*time.Time` | `created_at` | Registration creation date/time |
