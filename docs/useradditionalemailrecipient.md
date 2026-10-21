# UserAdditionalEmailRecipient

`files_sdk.UserAdditionalEmailRecipient` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/useradditionalemailrecipient`, named `user_additional_email_recipient`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_additional_email_recipient "github.com/Files-com/files-sdk-go/v3/useradditionalemailrecipient"
)
```

## Operations

### List

`GET /user_additional_email_recipients`

List User Additional Email Recipients

```go
func (c *Client) List(params files_sdk.UserAdditionalEmailRecipientListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.UserAdditionalEmailRecipientListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.UserAdditionalEmailRecipient` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.UserAdditionalEmailRecipientListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `email`, `user_id` or `workspace_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `email`, `workspace_id` or `user_id`. Valid field combinations are `[ workspace_id, email ]`, `[ user_id, email ]`, `[ workspace_id, user_id ]` or `[ workspace_id, user_id, email ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `email`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := user_additional_email_recipient.List(files_sdk.UserAdditionalEmailRecipientListParams{})
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

`GET /user_additional_email_recipients/{id}`

Show User Additional Email Recipient

```go
func (c *Client) Find(params files_sdk.UserAdditionalEmailRecipientFindParams, opts ...files_sdk.RequestResponseOption) (userAdditionalEmailRecipient files_sdk.UserAdditionalEmailRecipient, err error)
func Find(params files_sdk.UserAdditionalEmailRecipientFindParams, opts ...files_sdk.RequestResponseOption) (userAdditionalEmailRecipient files_sdk.UserAdditionalEmailRecipient, err error)
```

Parameters, `files_sdk.UserAdditionalEmailRecipientFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User Additional Email Recipient ID. |

Example:

```go
result, err := user_additional_email_recipient.Find(files_sdk.UserAdditionalEmailRecipientFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /user_additional_email_recipients`

Create User Additional Email Recipient

```go
func (c *Client) Create(params files_sdk.UserAdditionalEmailRecipientCreateParams, opts ...files_sdk.RequestResponseOption) (userAdditionalEmailRecipient files_sdk.UserAdditionalEmailRecipient, err error)
func Create(params files_sdk.UserAdditionalEmailRecipientCreateParams, opts ...files_sdk.RequestResponseOption) (userAdditionalEmailRecipient files_sdk.UserAdditionalEmailRecipient, err error)
```

Parameters, `files_sdk.UserAdditionalEmailRecipientCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `UserId` | `int64` | `user_id` | no | User ID. Provide a value of `0` to operate the current session's user. |
| `Email` | `string` | `email` | yes | Additional email recipient address |

Example:

```go
result, err := user_additional_email_recipient.Create(files_sdk.UserAdditionalEmailRecipientCreateParams{
	Email: "user-copy@example.com",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /user_additional_email_recipients/{id}`

Update User Additional Email Recipient

```go
func (c *Client) Update(params files_sdk.UserAdditionalEmailRecipientUpdateParams, opts ...files_sdk.RequestResponseOption) (userAdditionalEmailRecipient files_sdk.UserAdditionalEmailRecipient, err error)
func Update(params files_sdk.UserAdditionalEmailRecipientUpdateParams, opts ...files_sdk.RequestResponseOption) (userAdditionalEmailRecipient files_sdk.UserAdditionalEmailRecipient, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (userAdditionalEmailRecipient files_sdk.UserAdditionalEmailRecipient, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (userAdditionalEmailRecipient files_sdk.UserAdditionalEmailRecipient, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.UserAdditionalEmailRecipientUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User Additional Email Recipient ID. |
| `Email` | `string` | `email` | no | Additional email recipient address |

Example:

```go
result, err := user_additional_email_recipient.Update(files_sdk.UserAdditionalEmailRecipientUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := user_additional_email_recipient.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /user_additional_email_recipients/{id}`

Delete User Additional Email Recipient

```go
func (c *Client) Delete(params files_sdk.UserAdditionalEmailRecipientDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.UserAdditionalEmailRecipientDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.UserAdditionalEmailRecipientDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User Additional Email Recipient ID. |

Example:

```go
err := user_additional_email_recipient.Delete(files_sdk.UserAdditionalEmailRecipientDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.UserAdditionalEmailRecipient` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | User additional email recipient ID |
| `UserId` | `int64` | `user_id` | User ID |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID (0 for default workspace). |
| `Email` | `string` | `email` | Additional email recipient address |
| `CreatedAt` | `*time.Time` | `created_at` | Created at date/time |
