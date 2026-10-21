# Secret

`files_sdk.Secret` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/secret`, named `secret`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/secret"
)
```

## Operations

### List

`GET /secrets`

List Secrets

```go
func (c *Client) List(params files_sdk.SecretListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.SecretListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Secret` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.SecretListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name` or `secret_type`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`, `name` or `secret_type`. Valid field combinations are `[ workspace_id, name ]`, `[ workspace_id, secret_type ]`, `[ secret_type, name ]` or `[ workspace_id, secret_type, name ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := secret.List(files_sdk.SecretListParams{})
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

`GET /secrets/{id}`

Show Secret

```go
func (c *Client) Find(params files_sdk.SecretFindParams, opts ...files_sdk.RequestResponseOption) (secret files_sdk.Secret, err error)
func Find(params files_sdk.SecretFindParams, opts ...files_sdk.RequestResponseOption) (secret files_sdk.Secret, err error)
```

Parameters, `files_sdk.SecretFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Secret ID. |

Example:

```go
result, err := secret.Find(files_sdk.SecretFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /secrets`

Create Secret

```go
func (c *Client) Create(params files_sdk.SecretCreateParams, opts ...files_sdk.RequestResponseOption) (secret files_sdk.Secret, err error)
func Create(params files_sdk.SecretCreateParams, opts ...files_sdk.RequestResponseOption) (secret files_sdk.Secret, err error)
```

Parameters, `files_sdk.SecretCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | yes | Secret name. |
| `Description` | `string` | `description` | no | Internal description for your reference. |
| `SecretType` | `files_sdk.SecretSecretTypeEnum` | `secret_type` | yes | Secret type. Values: `basic`, `token`, `headers`, `certificate`, `key_value`. |
| `Metadata` | `interface{}` | `metadata` | no | Non-secret metadata for the Secret type. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. 0 means the default workspace. |

Example:

```go
result, err := secret.Create(files_sdk.SecretCreateParams{
	Name:       "Production API token",
	SecretType: "token",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /secrets/{id}`

Update Secret

```go
func (c *Client) Update(params files_sdk.SecretUpdateParams, opts ...files_sdk.RequestResponseOption) (secret files_sdk.Secret, err error)
func Update(params files_sdk.SecretUpdateParams, opts ...files_sdk.RequestResponseOption) (secret files_sdk.Secret, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (secret files_sdk.Secret, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (secret files_sdk.Secret, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.SecretUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Secret ID. |
| `Name` | `string` | `name` | no | Secret name. |
| `Description` | `string` | `description` | no | Internal description for your reference. |
| `SecretType` | `files_sdk.SecretSecretTypeEnum` | `secret_type` | no | Secret type. Values: `basic`, `token`, `headers`, `certificate`, `key_value`. |
| `Metadata` | `interface{}` | `metadata` | no | Non-secret metadata for the Secret type. |

Example:

```go
result, err := secret.Update(files_sdk.SecretUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := secret.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /secrets/{id}`

Delete Secret

```go
func (c *Client) Delete(params files_sdk.SecretDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.SecretDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.SecretDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Secret ID. |

Example:

```go
err := secret.Delete(files_sdk.SecretDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Secret` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Secret ID. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. 0 means the default workspace. |
| `Name` | `string` | `name` | Secret name. |
| `Description` | `string` | `description` | Internal description for your reference. |
| `SecretType` | `string` | `secret_type` | Secret type. |
| `Metadata` | `interface{}` | `metadata` | Non-secret metadata for the Secret type. |
| `ValueFieldNames` | `[]string` | `value_field_names` | Names of configured secret value fields. Secret values are never returned. |
| `CreatedAt` | `*time.Time` | `created_at` | Secret create date/time. |
| `UpdatedAt` | `*time.Time` | `updated_at` | Secret update date/time. |
