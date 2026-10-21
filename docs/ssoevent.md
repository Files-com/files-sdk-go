# SsoEvent

`files_sdk.SsoEvent` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/ssoevent`, named `sso_event`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	sso_event "github.com/Files-com/files-sdk-go/v3/ssoevent"
)
```

## Operations

### List

`GET /sso_events`

List SSO Events

```go
func (c *Client) List(params files_sdk.SsoEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.SsoEventListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.SsoEvent` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.SsoEventListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `created_at`, `event_type`, `status` or `user_id`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `event_type`, `idp_uid`, `ip`, `provider`, `status`, `user_id` or `username`. Valid field combinations are `[ event_type, created_at ]`, `[ idp_uid, created_at ]`, `[ ip, created_at ]`, `[ provider, created_at ]`, `[ status, created_at ]`, `[ user_id, created_at ]`, `[ username, created_at ]`, `[ event_type, status ]`, `[ idp_uid, status ]`, `[ ip, status ]`, `[ provider, status ]`, `[ user_id, status ]`, `[ username, status ]`, `[ event_type, status, created_at ]`, `[ idp_uid, status, created_at ]`, `[ ip, status, created_at ]`, `[ provider, status, created_at ]`, `[ user_id, status, created_at ]` or `[ username, status, created_at ]`. |
| `FilterGt` | `interface{}` | `filter_gt` | no | If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`. |
| `FilterGteq` | `interface{}` | `filter_gteq` | no | If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `idp_uid`, `ip`, `provider` or `username`. |
| `FilterLt` | `interface{}` | `filter_lt` | no | If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`. |
| `FilterLteq` | `interface{}` | `filter_lteq` | no | If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := sso_event.List(files_sdk.SsoEventListParams{})
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

`GET /sso_events/{id}`

Show SSO Event

```go
func (c *Client) Find(params files_sdk.SsoEventFindParams, opts ...files_sdk.RequestResponseOption) (ssoEvent files_sdk.SsoEvent, err error)
func Find(params files_sdk.SsoEventFindParams, opts ...files_sdk.RequestResponseOption) (ssoEvent files_sdk.SsoEvent, err error)
```

Parameters, `files_sdk.SsoEventFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Sso Event ID. |

Example:

```go
result, err := sso_event.Find(files_sdk.SsoEventFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.SsoEvent` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Event ID |
| `EventType` | `string` | `event_type` | Type of SSO event being recorded. |
| `Status` | `string` | `status` | Status of event. |
| `Body` | `string` | `body` | Event body. |
| `EventErrors` | `[]string` | `event_errors` | Event errors. |
| `CreatedAt` | `*time.Time` | `created_at` | Event create date/time. |
| `BodyUrl` | `string` | `body_url` | Link to log file. |
| `UserId` | `int64` | `user_id` | User ID. |
| `Username` | `string` | `username` | Username on Files.com for the SSO login attempt. |
| `IdpUid` | `string` | `idp_uid` | Identity Provider UID for the SSO login attempt. |
| `Provider` | `string` | `provider` | SSO provider for the SSO login attempt. |
| `ProviderLabel` | `string` | `provider_label` | SSO provider label for the SSO login attempt. |
| `Ip` | `string` | `ip` | IP address for the SSO login attempt. |
| `Region` | `string` | `region` | Region for the SSO login attempt. |
