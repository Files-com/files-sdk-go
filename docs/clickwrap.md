# Clickwrap

`files_sdk.Clickwrap` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/clickwrap`, named `clickwrap`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/clickwrap"
)
```

## Operations

### List

`GET /clickwraps`

List Clickwraps

```go
func (c *Client) List(params files_sdk.ClickwrapListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ClickwrapListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Clickwrap` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ClickwrapListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := clickwrap.List(files_sdk.ClickwrapListParams{})
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

`GET /clickwraps/{id}`

Show Clickwrap

```go
func (c *Client) Find(params files_sdk.ClickwrapFindParams, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error)
func Find(params files_sdk.ClickwrapFindParams, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error)
```

Parameters, `files_sdk.ClickwrapFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Clickwrap ID. |

Example:

```go
result, err := clickwrap.Find(files_sdk.ClickwrapFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /clickwraps`

Create Clickwrap

```go
func (c *Client) Create(params files_sdk.ClickwrapCreateParams, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error)
func Create(params files_sdk.ClickwrapCreateParams, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error)
```

Parameters, `files_sdk.ClickwrapCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | no | Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.) |
| `Body` | `string` | `body` | no | Body text of Clickwrap (supports Markdown formatting). |
| `UseWithBundles` | `files_sdk.ClickwrapUseWithBundlesEnum` | `use_with_bundles` | no | Use this Clickwrap for Bundles? Values: `none`, `available`, `require`, `available_to_all_users`. |
| `UseWithInboxes` | `files_sdk.ClickwrapUseWithInboxesEnum` | `use_with_inboxes` | no | Use this Clickwrap for Inboxes? Values: `none`, `available`, `require`, `available_to_all_users`. |
| `UseWithUsers` | `files_sdk.ClickwrapUseWithUsersEnum` | `use_with_users` | no | Use this Clickwrap for Users? Values: `none`, `require` (new user signup via email invitation only), `require_all_users_once` (show to all users at their next web login; once accepted, not shown again), `require_all_users_always` (show to all users on every web login). Values: `none`, `require`, `require_all_users_once`, `require_all_users_always`. |

Example:

```go
result, err := clickwrap.Create(files_sdk.ClickwrapCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /clickwraps/{id}`

Update Clickwrap

```go
func (c *Client) Update(params files_sdk.ClickwrapUpdateParams, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error)
func Update(params files_sdk.ClickwrapUpdateParams, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (clickwrap files_sdk.Clickwrap, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.ClickwrapUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Clickwrap ID. |
| `Name` | `string` | `name` | no | Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.) |
| `Body` | `string` | `body` | no | Body text of Clickwrap (supports Markdown formatting). |
| `UseWithBundles` | `files_sdk.ClickwrapUseWithBundlesEnum` | `use_with_bundles` | no | Use this Clickwrap for Bundles? Values: `none`, `available`, `require`, `available_to_all_users`. |
| `UseWithInboxes` | `files_sdk.ClickwrapUseWithInboxesEnum` | `use_with_inboxes` | no | Use this Clickwrap for Inboxes? Values: `none`, `available`, `require`, `available_to_all_users`. |
| `UseWithUsers` | `files_sdk.ClickwrapUseWithUsersEnum` | `use_with_users` | no | Use this Clickwrap for Users? Values: `none`, `require` (new user signup via email invitation only), `require_all_users_once` (show to all users at their next web login; once accepted, not shown again), `require_all_users_always` (show to all users on every web login). Values: `none`, `require`, `require_all_users_once`, `require_all_users_always`. |

Example:

```go
result, err := clickwrap.Update(files_sdk.ClickwrapUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := clickwrap.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /clickwraps/{id}`

Delete Clickwrap

```go
func (c *Client) Delete(params files_sdk.ClickwrapDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.ClickwrapDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.ClickwrapDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Clickwrap ID. |

Example:

```go
err := clickwrap.Delete(files_sdk.ClickwrapDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Clickwrap` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Clickwrap ID |
| `Name` | `string` | `name` | Name of the Clickwrap agreement (used when selecting from multiple Clickwrap agreements.) |
| `Body` | `string` | `body` | Body text of Clickwrap (supports Markdown formatting). |
| `UseWithUsers` | `string` | `use_with_users` | Use this Clickwrap for Users? Values: `none`, `require` (new user signup via email invitation only), `require_all_users_once` (show to all users at their next web login; once accepted, not shown again), `require_all_users_always` (show to all users on every web login). |
| `UseWithBundles` | `string` | `use_with_bundles` | Use this Clickwrap for Bundles? |
| `UseWithInboxes` | `string` | `use_with_inboxes` | Use this Clickwrap for Inboxes? |
