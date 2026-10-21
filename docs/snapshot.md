# Snapshot

`files_sdk.Snapshot` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/snapshot`, named `snapshot`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/snapshot"
)
```

## Operations

### List

`GET /snapshots`

List Snapshots

```go
func (c *Client) List(params files_sdk.SnapshotListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.SnapshotListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Snapshot` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.SnapshotListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := snapshot.List(files_sdk.SnapshotListParams{})
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

`GET /snapshots/{id}`

Show Snapshot

```go
func (c *Client) Find(params files_sdk.SnapshotFindParams, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error)
func Find(params files_sdk.SnapshotFindParams, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error)
```

Parameters, `files_sdk.SnapshotFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Snapshot ID. |

Example:

```go
result, err := snapshot.Find(files_sdk.SnapshotFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /snapshots`

Create Snapshot

```go
func (c *Client) Create(params files_sdk.SnapshotCreateParams, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error)
func Create(params files_sdk.SnapshotCreateParams, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error)
```

Parameters, `files_sdk.SnapshotCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `ExpiresAt` | `*time.Time` | `expires_at` | no | When the snapshot expires. |
| `Name` | `string` | `name` | no | A name for the snapshot. |
| `Paths` | `[]string` | `paths` | no | An array of paths to add to the snapshot. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID. `0` means the default workspace. |

Example:

```go
result, err := snapshot.Create(files_sdk.SnapshotCreateParams{})
if err != nil {
	return err
}
fmt.Println(result)
```

### Finalize

`POST /snapshots/{id}/finalize`

Finalize Snapshot

```go
func (c *Client) Finalize(params files_sdk.SnapshotFinalizeParams, opts ...files_sdk.RequestResponseOption) (err error)
func Finalize(params files_sdk.SnapshotFinalizeParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.SnapshotFinalizeParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Snapshot ID. |

Example:

```go
err := snapshot.Finalize(files_sdk.SnapshotFinalizeParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### Update

`PATCH /snapshots/{id}`

Update Snapshot

```go
func (c *Client) Update(params files_sdk.SnapshotUpdateParams, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error)
func Update(params files_sdk.SnapshotUpdateParams, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (snapshot files_sdk.Snapshot, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.SnapshotUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Snapshot ID. |
| `ExpiresAt` | `*time.Time` | `expires_at` | no | When the snapshot expires. |
| `Name` | `string` | `name` | no | A name for the snapshot. |
| `Paths` | `[]string` | `paths` | no | An array of paths to add to the snapshot. |

Example:

```go
result, err := snapshot.Update(files_sdk.SnapshotUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := snapshot.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /snapshots/{id}`

Delete Snapshot

```go
func (c *Client) Delete(params files_sdk.SnapshotDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.SnapshotDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.SnapshotDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Snapshot ID. |

Example:

```go
err := snapshot.Delete(files_sdk.SnapshotDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Snapshot` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | The snapshot's unique ID. |
| `ExpiresAt` | `*time.Time` | `expires_at` | When the snapshot expires. |
| `FinalizedAt` | `*time.Time` | `finalized_at` | When the snapshot was finalized. |
| `Name` | `string` | `name` | A name for the snapshot. |
| `UserId` | `int64` | `user_id` | The user that created this snapshot, if applicable. |
| `BundleId` | `int64` | `bundle_id` | The bundle using this snapshot, if applicable. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID. `0` means the default workspace. |
| `Paths` | `[]string` | `paths` | An array of paths to add to the snapshot. |
