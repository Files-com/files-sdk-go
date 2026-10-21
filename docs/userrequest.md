# UserRequest

`files_sdk.UserRequest` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/userrequest`, named `user_request`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	user_request "github.com/Files-com/files-sdk-go/v3/userrequest"
)
```

## Operations

### List

`GET /user_requests`

List User Requests

```go
func (c *Client) List(params files_sdk.UserRequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.UserRequestListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.UserRequest` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.UserRequestListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := user_request.List(files_sdk.UserRequestListParams{})
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

`GET /user_requests/{id}`

Show User Request

```go
func (c *Client) Find(params files_sdk.UserRequestFindParams, opts ...files_sdk.RequestResponseOption) (userRequest files_sdk.UserRequest, err error)
func Find(params files_sdk.UserRequestFindParams, opts ...files_sdk.RequestResponseOption) (userRequest files_sdk.UserRequest, err error)
```

Parameters, `files_sdk.UserRequestFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User Request ID. |

Example:

```go
result, err := user_request.Find(files_sdk.UserRequestFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /user_requests`

Create User Request

```go
func (c *Client) Create(params files_sdk.UserRequestCreateParams, opts ...files_sdk.RequestResponseOption) (userRequest files_sdk.UserRequest, err error)
func Create(params files_sdk.UserRequestCreateParams, opts ...files_sdk.RequestResponseOption) (userRequest files_sdk.UserRequest, err error)
```

Parameters, `files_sdk.UserRequestCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | yes | Name of user requested |
| `Email` | `string` | `email` | yes | Email of user requested |
| `Details` | `string` | `details` | yes | Details of the user request |
| `Company` | `string` | `company` | no | Company of the user requested |

Example:

```go
result, err := user_request.Create(files_sdk.UserRequestCreateParams{
	Name:    "name",
	Email:   "email",
	Details: "details",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /user_requests/{id}`

Delete User Request

```go
func (c *Client) Delete(params files_sdk.UserRequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.UserRequestDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.UserRequestDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | User Request ID. |

Example:

```go
err := user_request.Delete(files_sdk.UserRequestDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.UserRequest` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | ID |
| `Name` | `string` | `name` | User's full name |
| `Email` | `string` | `email` | User email address |
| `Details` | `string` | `details` | Details of the user's request |
| `Company` | `string` | `company` | User's company name |
