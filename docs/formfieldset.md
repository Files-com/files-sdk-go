# FormFieldSet

`files_sdk.FormFieldSet` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/formfieldset`, named `form_field_set`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	form_field_set "github.com/Files-com/files-sdk-go/v3/formfieldset"
)
```

## Operations

### List

`GET /form_field_sets`

List Form Field Sets

```go
func (c *Client) List(params files_sdk.FormFieldSetListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.FormFieldSetListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.FormFieldSet` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.FormFieldSetListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := form_field_set.List(files_sdk.FormFieldSetListParams{})
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

`GET /form_field_sets/{id}`

Show Form Field Set

```go
func (c *Client) Find(params files_sdk.FormFieldSetFindParams, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error)
func Find(params files_sdk.FormFieldSetFindParams, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error)
```

Parameters, `files_sdk.FormFieldSetFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Form Field Set ID. |

Example:

```go
result, err := form_field_set.Find(files_sdk.FormFieldSetFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /form_field_sets`

Create Form Field Set

```go
func (c *Client) Create(params files_sdk.FormFieldSetCreateParams, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error)
func Create(params files_sdk.FormFieldSetCreateParams, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error)
```

Parameters, `files_sdk.FormFieldSetCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `Title` | `string` | `title` | no | Title to be displayed |
| `RequireCompleteMetadataCapture` | `*bool` | `require_complete_metadata_capture` | no | Reject answers over 1024 characters or more than 32 answered custom fields to ensure complete file metadata capture |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID |
| `SkipEmail` | `*bool` | `skip_email` | no | Skip validating form email |
| `SkipName` | `*bool` | `skip_name` | no | Skip validating form name |
| `SkipCompany` | `*bool` | `skip_company` | no | Skip validating company |
| `FormFields` | `[]map[string]interface{}` | `form_fields` | no |  |

Example:

```go
result, err := form_field_set.Create(files_sdk.FormFieldSetCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /form_field_sets/{id}`

Update Form Field Set

```go
func (c *Client) Update(params files_sdk.FormFieldSetUpdateParams, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error)
func Update(params files_sdk.FormFieldSetUpdateParams, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (formFieldSet files_sdk.FormFieldSet, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.FormFieldSetUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Form Field Set ID. |
| `Title` | `string` | `title` | no | Title to be displayed |
| `RequireCompleteMetadataCapture` | `*bool` | `require_complete_metadata_capture` | no | Reject answers over 1024 characters or more than 32 answered custom fields to ensure complete file metadata capture |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID |
| `SkipEmail` | `*bool` | `skip_email` | no | Skip validating form email |
| `SkipName` | `*bool` | `skip_name` | no | Skip validating form name |
| `SkipCompany` | `*bool` | `skip_company` | no | Skip validating company |
| `FormFields` | `[]map[string]interface{}` | `form_fields` | no |  |

Example:

```go
result, err := form_field_set.Update(files_sdk.FormFieldSetUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := form_field_set.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /form_field_sets/{id}`

Delete Form Field Set

```go
func (c *Client) Delete(params files_sdk.FormFieldSetDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.FormFieldSetDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.FormFieldSetDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Form Field Set ID. |

Example:

```go
err := form_field_set.Delete(files_sdk.FormFieldSetDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.FormFieldSet` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Form field set id |
| `Title` | `string` | `title` | Title to be displayed |
| `RequireCompleteMetadataCapture` | `*bool` | `require_complete_metadata_capture` | Reject answers over 1024 characters or more than 32 answered custom fields to ensure complete file metadata capture |
| `FormLayout` | `[]int64` | `form_layout` | Layout of the form |
| `FormFields` | `[]files_sdk.FormField` | `form_fields` | Associated form field definitions; authenticated form field set responses include historical definitions, while form_layout identifies current fields |
| `SkipName` | `*bool` | `skip_name` | Any associated InboxRegistrations or BundleRegistrations can be saved without providing name |
| `SkipEmail` | `*bool` | `skip_email` | Any associated InboxRegistrations or BundleRegistrations can be saved without providing email |
| `SkipCompany` | `*bool` | `skip_company` | Any associated InboxRegistrations or BundleRegistrations can be saved without providing company |
| `InUse` | `*bool` | `in_use` | Form Field Set is in use by an active Inbox / Bundle / Inbox Registration / Bundle Registration |
| `UserId` | `int64` | `user_id` | User ID. Provide a value of `0` to operate the current session's user. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID |
