# RemoteMountBackend

`files_sdk.RemoteMountBackend` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/remotemountbackend`, named `remote_mount_backend`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	remote_mount_backend "github.com/Files-com/files-sdk-go/v3/remotemountbackend"
)
```

## Operations

### List

`GET /remote_mount_backends`

List Remote Mount Backends

```go
func (c *Client) List(params files_sdk.RemoteMountBackendListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.RemoteMountBackendListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.RemoteMountBackend` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.RemoteMountBackendListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `remote_server_mount_id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := remote_mount_backend.List(files_sdk.RemoteMountBackendListParams{})
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

`GET /remote_mount_backends/{id}`

Show Remote Mount Backend

```go
func (c *Client) Find(params files_sdk.RemoteMountBackendFindParams, opts ...files_sdk.RequestResponseOption) (remoteMountBackend files_sdk.RemoteMountBackend, err error)
func Find(params files_sdk.RemoteMountBackendFindParams, opts ...files_sdk.RequestResponseOption) (remoteMountBackend files_sdk.RemoteMountBackend, err error)
```

Parameters, `files_sdk.RemoteMountBackendFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Mount Backend ID. |

Example:

```go
result, err := remote_mount_backend.Find(files_sdk.RemoteMountBackendFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /remote_mount_backends`

Create Remote Mount Backend

```go
func (c *Client) Create(params files_sdk.RemoteMountBackendCreateParams, opts ...files_sdk.RequestResponseOption) (remoteMountBackend files_sdk.RemoteMountBackend, err error)
func Create(params files_sdk.RemoteMountBackendCreateParams, opts ...files_sdk.RequestResponseOption) (remoteMountBackend files_sdk.RemoteMountBackend, err error)
```

Parameters, `files_sdk.RemoteMountBackendCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Enabled` | `*bool` | `enabled` | no | True if this backend is enabled. |
| `Fall` | `int64` | `fall` | no | Number of consecutive failures before considering the backend unhealthy. |
| `HealthCheckEnabled` | `*bool` | `health_check_enabled` | no | True if health checks are enabled for this backend. |
| `HealthCheckType` | `files_sdk.RemoteMountBackendHealthCheckTypeEnum` | `health_check_type` | no | Type of health check to perform. Values: `active`, `passive`. |
| `Interval` | `int64` | `interval` | no | Interval in seconds between health checks. |
| `MinFreeCpu` | `float64` | `min_free_cpu` | no | Minimum free CPU percentage required for this backend to be considered healthy. |
| `MinFreeCpuDecimal` | `*string` | `min_free_cpu` | no | Sends exact decimal text, such as "1.5", as min_free_cpu. Leave MinFreeCpu zero when you set it. |
| `MinFreeMem` | `float64` | `min_free_mem` | no | Minimum free memory percentage required for this backend to be considered healthy. |
| `MinFreeMemDecimal` | `*string` | `min_free_mem` | no | Sends exact decimal text, such as "1.5", as min_free_mem. Leave MinFreeMem zero when you set it. |
| `Priority` | `int64` | `priority` | no | Priority of this backend. |
| `RemotePath` | `string` | `remote_path` | no | Path on the remote server to treat as the root of this mount. |
| `Rise` | `int64` | `rise` | no | Number of consecutive successes before considering the backend healthy. |
| `CanaryFilePath` | `string` | `canary_file_path` | yes | Path to the canary file used for health checks. |
| `RemoteServerMountId` | `int64` | `remote_server_mount_id` | yes | The mount ID of the Remote Server Mount that this backend is associated with. |
| `RemoteServerId` | `int64` | `remote_server_id` | yes | The remote server that this backend is associated with. |

Example:

```go
result, err := remote_mount_backend.Create(files_sdk.RemoteMountBackendCreateParams{
	CanaryFilePath:      "backend1.txt",
	RemoteServerMountId: 1,
	RemoteServerId:      1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### ResetStatus

`POST /remote_mount_backends/{id}/reset_status`

Reset backend status to healthy

```go
func (c *Client) ResetStatus(params files_sdk.RemoteMountBackendResetStatusParams, opts ...files_sdk.RequestResponseOption) (err error)
func ResetStatus(params files_sdk.RemoteMountBackendResetStatusParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.RemoteMountBackendResetStatusParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Mount Backend ID. |

Example:

```go
err := remote_mount_backend.ResetStatus(files_sdk.RemoteMountBackendResetStatusParams{
	Id: 1,
})
if err != nil {
	return err
}
```

### Update

`PATCH /remote_mount_backends/{id}`

Update Remote Mount Backend

```go
func (c *Client) Update(params files_sdk.RemoteMountBackendUpdateParams, opts ...files_sdk.RequestResponseOption) (remoteMountBackend files_sdk.RemoteMountBackend, err error)
func Update(params files_sdk.RemoteMountBackendUpdateParams, opts ...files_sdk.RequestResponseOption) (remoteMountBackend files_sdk.RemoteMountBackend, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (remoteMountBackend files_sdk.RemoteMountBackend, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (remoteMountBackend files_sdk.RemoteMountBackend, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.RemoteMountBackendUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Mount Backend ID. |
| `Enabled` | `*bool` | `enabled` | no | True if this backend is enabled. |
| `Fall` | `int64` | `fall` | no | Number of consecutive failures before considering the backend unhealthy. |
| `HealthCheckEnabled` | `*bool` | `health_check_enabled` | no | True if health checks are enabled for this backend. |
| `HealthCheckType` | `files_sdk.RemoteMountBackendHealthCheckTypeEnum` | `health_check_type` | no | Type of health check to perform. Values: `active`, `passive`. |
| `Interval` | `int64` | `interval` | no | Interval in seconds between health checks. |
| `MinFreeCpu` | `float64` | `min_free_cpu` | no | Minimum free CPU percentage required for this backend to be considered healthy. |
| `MinFreeCpuDecimal` | `*string` | `min_free_cpu` | no | Sends exact decimal text, such as "1.5", as min_free_cpu. Leave MinFreeCpu zero when you set it. |
| `MinFreeMem` | `float64` | `min_free_mem` | no | Minimum free memory percentage required for this backend to be considered healthy. |
| `MinFreeMemDecimal` | `*string` | `min_free_mem` | no | Sends exact decimal text, such as "1.5", as min_free_mem. Leave MinFreeMem zero when you set it. |
| `Priority` | `int64` | `priority` | no | Priority of this backend. |
| `RemotePath` | `string` | `remote_path` | no | Path on the remote server to treat as the root of this mount. |
| `Rise` | `int64` | `rise` | no | Number of consecutive successes before considering the backend healthy. |
| `CanaryFilePath` | `string` | `canary_file_path` | no | Path to the canary file used for health checks. |
| `RemoteServerId` | `int64` | `remote_server_id` | no | The remote server that this backend is associated with. |

Example:

```go
result, err := remote_mount_backend.Update(files_sdk.RemoteMountBackendUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := remote_mount_backend.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /remote_mount_backends/{id}`

Delete Remote Mount Backend

```go
func (c *Client) Delete(params files_sdk.RemoteMountBackendDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.RemoteMountBackendDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.RemoteMountBackendDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Remote Mount Backend ID. |

Example:

```go
err := remote_mount_backend.Delete(files_sdk.RemoteMountBackendDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.RemoteMountBackend` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `CanaryFilePath` | `string` | `canary_file_path` | Path to the canary file used for health checks. |
| `Enabled` | `*bool` | `enabled` | True if this backend is enabled. |
| `Fall` | `int64` | `fall` | Number of consecutive failures before considering the backend unhealthy. |
| `HealthCheckEnabled` | `*bool` | `health_check_enabled` | True if health checks are enabled for this backend. |
| `HealthCheckResults` | `[]map[string]interface{}` | `health_check_results` | Array of recent health check results. |
| `HealthCheckType` | `string` | `health_check_type` | Type of health check to perform. |
| `Id` | `int64` | `id` | Unique identifier for this backend. |
| `Interval` | `int64` | `interval` | Interval in seconds between health checks. |
| `MinFreeCpu` | `string` | `min_free_cpu` | Minimum free CPU percentage required for this backend to be considered healthy. |
| `MinFreeMem` | `string` | `min_free_mem` | Minimum free memory percentage required for this backend to be considered healthy. |
| `Priority` | `int64` | `priority` | Priority of this backend. |
| `RemotePath` | `string` | `remote_path` | Path on the remote server to treat as the root of this mount. |
| `RemoteServerId` | `int64` | `remote_server_id` | The remote server that this backend is associated with. |
| `RemoteServerMountId` | `int64` | `remote_server_mount_id` | The mount ID of the Remote Server Mount that this backend is associated with. |
| `Rise` | `int64` | `rise` | Number of consecutive successes before considering the backend healthy. |
| `Status` | `string` | `status` | Status of this backend. |
| `UndergoingMaintenance` | `*bool` | `undergoing_maintenance` | True if this backend is undergoing maintenance. |
