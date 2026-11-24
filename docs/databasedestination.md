# DatabaseDestination

`files_sdk.DatabaseDestination` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/databasedestination`, named `database_destination`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	database_destination "github.com/Files-com/files-sdk-go/v3/databasedestination"
)
```

## Operations

### List

`GET /database_destinations`

List Database Destinations

```go
func (c *Client) List(params files_sdk.DatabaseDestinationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.DatabaseDestinationListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.DatabaseDestination` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.DatabaseDestinationListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id` and `name`. Valid field combinations are `[ workspace_id, name ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := database_destination.List(files_sdk.DatabaseDestinationListParams{})
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

`GET /database_destinations/{id}`

Show Database Destination

```go
func (c *Client) Find(params files_sdk.DatabaseDestinationFindParams, opts ...files_sdk.RequestResponseOption) (databaseDestination files_sdk.DatabaseDestination, err error)
func Find(params files_sdk.DatabaseDestinationFindParams, opts ...files_sdk.RequestResponseOption) (databaseDestination files_sdk.DatabaseDestination, err error)
```

Parameters, `files_sdk.DatabaseDestinationFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Database Destination ID. |

Example:

```go
result, err := database_destination.Find(files_sdk.DatabaseDestinationFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /database_destinations`

Create Database Destination

```go
func (c *Client) Create(params files_sdk.DatabaseDestinationCreateParams, opts ...files_sdk.RequestResponseOption) (databaseDestination files_sdk.DatabaseDestination, err error)
func Create(params files_sdk.DatabaseDestinationCreateParams, opts ...files_sdk.RequestResponseOption) (databaseDestination files_sdk.DatabaseDestination, err error)
```

Parameters, `files_sdk.DatabaseDestinationCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | yes | Name of the destination. |
| `Provider` | `files_sdk.DatabaseDestinationProviderEnum` | `provider` | yes | Database provider. Values: `snowflake`, `bigquery`, `supabase`, `postgresql`, `mysql`. |
| `Config` | `interface{}` | `config` | yes | Provider settings for the destination. Use the schema matching provider. The connection address and table are fixed after creation. |
| `SecretId` | `int64` | `secret_id` | yes | ID of a Secret in the same Workspace. Use a Token Secret for Snowflake, a service account Key/Value Secret for BigQuery, or a Username/Password Secret for a database login. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. 0 means the default workspace. |

Example:

```go
// Set this required value before the call.
var config interface{}
result, err := database_destination.Create(files_sdk.DatabaseDestinationCreateParams{
	Name:     "Orders",
	Provider: "snowflake",
	Config:   config,
	SecretId: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /database_destinations/{id}`

Update Database Destination

```go
func (c *Client) Update(params files_sdk.DatabaseDestinationUpdateParams, opts ...files_sdk.RequestResponseOption) (databaseDestination files_sdk.DatabaseDestination, err error)
func Update(params files_sdk.DatabaseDestinationUpdateParams, opts ...files_sdk.RequestResponseOption) (databaseDestination files_sdk.DatabaseDestination, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (databaseDestination files_sdk.DatabaseDestination, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (databaseDestination files_sdk.DatabaseDestination, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.DatabaseDestinationUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Database Destination ID. |
| `Name` | `string` | `name` | no | Name of the destination. |
| `SecretId` | `int64` | `secret_id` | no | ID of a Secret in the same Workspace. Use a Token Secret for Snowflake, a service account Key/Value Secret for BigQuery, or a Username/Password Secret for a database login. |
| `Config` | `interface{}` | `config` | no | Provider settings for the destination. Use the schema matching provider. The connection address and table are fixed after creation. |

Example:

```go
result, err := database_destination.Update(files_sdk.DatabaseDestinationUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := database_destination.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /database_destinations/{id}`

Delete Database Destination

```go
func (c *Client) Delete(params files_sdk.DatabaseDestinationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.DatabaseDestinationDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.DatabaseDestinationDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Database Destination ID. |

Example:

```go
err := database_destination.Delete(files_sdk.DatabaseDestinationDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.DatabaseDestination` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Database Destination ID. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. 0 means the default workspace. |
| `Name` | `string` | `name` | Name of the destination. |
| `Provider` | `string` | `provider` | Database provider. |
| `Config` | `interface{}` | `config` | Provider settings for the destination. Use the schema matching provider. The connection address and table are fixed after creation. |
| `SecretId` | `int64` | `secret_id` | ID of a Secret in the same Workspace. Use a Token Secret for Snowflake, a service account Key/Value Secret for BigQuery, or a Username/Password Secret for a database login. |
| `CreatedAt` | `*time.Time` | `created_at` | When the destination was created. |
| `UpdatedAt` | `*time.Time` | `updated_at` | When the destination was last updated. |
