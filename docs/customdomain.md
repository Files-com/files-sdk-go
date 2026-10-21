# CustomDomain

`files_sdk.CustomDomain` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/customdomain`, named `custom_domain`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	custom_domain "github.com/Files-com/files-sdk-go/v3/customdomain"
)
```

## Operations

### List

`GET /custom_domains`

List Custom Domains

```go
func (c *Client) List(params files_sdk.CustomDomainListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.CustomDomainListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.CustomDomain` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.CustomDomainListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `id` or `available_to_all_workspaces`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := custom_domain.List(files_sdk.CustomDomainListParams{})
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

`GET /custom_domains/{id}`

Show Custom Domain

```go
func (c *Client) Find(params files_sdk.CustomDomainFindParams, opts ...files_sdk.RequestResponseOption) (customDomain files_sdk.CustomDomain, err error)
func Find(params files_sdk.CustomDomainFindParams, opts ...files_sdk.RequestResponseOption) (customDomain files_sdk.CustomDomain, err error)
```

Parameters, `files_sdk.CustomDomainFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Custom Domain ID. |

Example:

```go
result, err := custom_domain.Find(files_sdk.CustomDomainFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### CreateAllocateIp

`POST /custom_domains/{id}/allocate_ips`

Allocate dedicated IP addresses to this Custom Domain

```go
func (c *Client) CreateAllocateIp(params files_sdk.CustomDomainCreateAllocateIpParams, opts ...files_sdk.RequestResponseOption) (customDomain files_sdk.CustomDomain, err error)
func CreateAllocateIp(params files_sdk.CustomDomainCreateAllocateIpParams, opts ...files_sdk.RequestResponseOption) (customDomain files_sdk.CustomDomain, err error)
```

Parameters, `files_sdk.CustomDomainCreateAllocateIpParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Custom Domain ID. |
| `Count` | `int64` | `count` | yes | Number of dedicated IP addresses to allocate. Values: `1`, `2`. |

Example:

```go
result, err := custom_domain.CreateAllocateIp(files_sdk.CustomDomainCreateAllocateIpParams{
	Id:    1,
	Count: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /custom_domains`

Create Custom Domain

```go
func (c *Client) Create(params files_sdk.CustomDomainCreateParams, opts ...files_sdk.RequestResponseOption) (customDomain files_sdk.CustomDomain, err error)
func Create(params files_sdk.CustomDomainCreateParams, opts ...files_sdk.RequestResponseOption) (customDomain files_sdk.CustomDomain, err error)
```

Parameters, `files_sdk.CustomDomainCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `AvailableToAllWorkspaces` | `*bool` | `available_to_all_workspaces` | no | Allow all workspaces to use this default-workspace Custom Domain. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID (0 for the default workspace). |
| `Destination` | `files_sdk.CustomDomainDestinationEnum` | `destination` | no | Where this custom domain routes. Can be `site_alias`, `public_hosting`, `s3_endpoint`, or `unassigned` (not routing traffic). Set to `unassigned` automatically when a bound `public_hosting` folder behavior is deleted, and can be set manually via the API for any reason. Values: `site_alias`, `public_hosting`, `s3_endpoint`, `unassigned`. |
| `FolderBehaviorId` | `int64` | `folder_behavior_id` | no | Public Hosting behavior ID when this domain routes to a specific Public Hosting behavior. Preserved as historical context when `destination` becomes `unassigned`. |
| `SslCertificateId` | `int64` | `ssl_certificate_id` | no | Current SSL certificate ID. |
| `Domain` | `string` | `domain` | yes | Customer-owned domain name. |

Example:

```go
result, err := custom_domain.Create(files_sdk.CustomDomainCreateParams{
	Domain: "files.example.com",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /custom_domains/{id}`

Update Custom Domain

```go
func (c *Client) Update(params files_sdk.CustomDomainUpdateParams, opts ...files_sdk.RequestResponseOption) (customDomain files_sdk.CustomDomain, err error)
func Update(params files_sdk.CustomDomainUpdateParams, opts ...files_sdk.RequestResponseOption) (customDomain files_sdk.CustomDomain, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (customDomain files_sdk.CustomDomain, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (customDomain files_sdk.CustomDomain, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.CustomDomainUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Custom Domain ID. |
| `AvailableToAllWorkspaces` | `*bool` | `available_to_all_workspaces` | no | Allow all workspaces to use this default-workspace Custom Domain. |
| `WorkspaceId` | `int64` | `workspace_id` | no | Workspace ID (0 for the default workspace). |
| `Destination` | `files_sdk.CustomDomainDestinationEnum` | `destination` | no | Where this custom domain routes. Can be `site_alias`, `public_hosting`, `s3_endpoint`, or `unassigned` (not routing traffic). Set to `unassigned` automatically when a bound `public_hosting` folder behavior is deleted, and can be set manually via the API for any reason. Values: `site_alias`, `public_hosting`, `s3_endpoint`, `unassigned`. |
| `FolderBehaviorId` | `int64` | `folder_behavior_id` | no | Public Hosting behavior ID when this domain routes to a specific Public Hosting behavior. Preserved as historical context when `destination` becomes `unassigned`. |
| `SslCertificateId` | `int64` | `ssl_certificate_id` | no | Current SSL certificate ID. |
| `Domain` | `string` | `domain` | no | Customer-owned domain name. |

Example:

```go
result, err := custom_domain.Update(files_sdk.CustomDomainUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := custom_domain.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /custom_domains/{id}`

Delete Custom Domain

```go
func (c *Client) Delete(params files_sdk.CustomDomainDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.CustomDomainDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.CustomDomainDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Custom Domain ID. |

Example:

```go
err := custom_domain.Delete(files_sdk.CustomDomainDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.CustomDomain` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Custom Domain ID. |
| `WorkspaceId` | `int64` | `workspace_id` | Workspace ID (0 for the default workspace). |
| `AvailableToAllWorkspaces` | `*bool` | `available_to_all_workspaces` | Allow all workspaces to use this default-workspace Custom Domain. |
| `OutboundIpAddresses` | `[]string` | `outbound_ip_addresses` | Eligible public IP addresses for Remote Server outbound connections. Empty when this Custom Domain is not eligible for outbound selection. |
| `Domain` | `string` | `domain` | Customer-owned domain name. |
| `Destination` | `string` | `destination` | Where this custom domain routes. Can be `site_alias`, `public_hosting`, `s3_endpoint`, or `unassigned` (not routing traffic). Set to `unassigned` automatically when a bound `public_hosting` folder behavior is deleted, and can be set manually via the API for any reason. |
| `DnsStatus` | `string` | `dns_status` | Current DNS verification status. |
| `SslCertificateId` | `int64` | `ssl_certificate_id` | Current SSL certificate ID. |
| `BrickManaged` | `*bool` | `brick_managed` | Is this domain's SSL certificate automatically managed and renewed by Files.com? |
| `FolderBehaviorId` | `int64` | `folder_behavior_id` | Public Hosting behavior ID when this domain routes to a specific Public Hosting behavior. Preserved as historical context when `destination` becomes `unassigned`. |
| `IpAddresses` | `[]string` | `ip_addresses` | Dedicated public IP addresses allocated to this Custom Domain. |
| `CreatedAt` | `*time.Time` | `created_at` | When this Custom Domain was created. |
| `UpdatedAt` | `*time.Time` | `updated_at` | When this Custom Domain was last updated. |
