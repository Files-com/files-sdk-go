# SiteSubdomainRedirect

`files_sdk.SiteSubdomainRedirect` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/sitesubdomainredirect`, named `site_subdomain_redirect`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	site_subdomain_redirect "github.com/Files-com/files-sdk-go/v3/sitesubdomainredirect"
)
```

## Operations

### List

`GET /site_subdomain_redirects`

List Site Subdomain Redirects

```go
func (c *Client) List(params files_sdk.SiteSubdomainRedirectListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.SiteSubdomainRedirectListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.SiteSubdomainRedirect` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.SiteSubdomainRedirectListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `id`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := site_subdomain_redirect.List(files_sdk.SiteSubdomainRedirectListParams{})
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

`GET /site_subdomain_redirects/{id}`

Show Site Subdomain Redirect

```go
func (c *Client) Find(params files_sdk.SiteSubdomainRedirectFindParams, opts ...files_sdk.RequestResponseOption) (siteSubdomainRedirect files_sdk.SiteSubdomainRedirect, err error)
func Find(params files_sdk.SiteSubdomainRedirectFindParams, opts ...files_sdk.RequestResponseOption) (siteSubdomainRedirect files_sdk.SiteSubdomainRedirect, err error)
```

Parameters, `files_sdk.SiteSubdomainRedirectFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Site Subdomain Redirect ID. |

Example:

```go
result, err := site_subdomain_redirect.Find(files_sdk.SiteSubdomainRedirectFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /site_subdomain_redirects/{id}`

Delete Site Subdomain Redirect

```go
func (c *Client) Delete(params files_sdk.SiteSubdomainRedirectDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.SiteSubdomainRedirectDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.SiteSubdomainRedirectDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Site Subdomain Redirect ID. |

Example:

```go
err := site_subdomain_redirect.Delete(files_sdk.SiteSubdomainRedirectDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.SiteSubdomainRedirect` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Site subdomain redirect ID. |
| `Subdomain` | `string` | `subdomain` | Files.com subdomain that continues to route to the current site subdomain. |
| `CreatedAt` | `*time.Time` | `created_at` | When this redirect was created. |
| `UpdatedAt` | `*time.Time` | `updated_at` | When this redirect was last updated. |
