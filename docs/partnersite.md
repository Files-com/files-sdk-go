# PartnerSite

`files_sdk.PartnerSite` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/partnersite`, named `partner_site`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_site "github.com/Files-com/files-sdk-go/v3/partnersite"
)
```

## Operations

### Delete

`DELETE /partner_sites/{id}`

Delete Partner Site

```go
func (c *Client) Delete(params files_sdk.PartnerSiteDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.PartnerSiteDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.PartnerSiteDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Partner Site ID. |

Example:

```go
err := partner_site.Delete(files_sdk.PartnerSiteDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

The API schema gives `files_sdk.PartnerSite` no fields.
