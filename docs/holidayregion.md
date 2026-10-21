# HolidayRegion

`files_sdk.HolidayRegion` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/holidayregion`, named `holiday_region`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	holiday_region "github.com/Files-com/files-sdk-go/v3/holidayregion"
)
```

## Operations

### GetSupported

`GET /holiday_regions/supported`

List all possible holiday regions

```go
func (c *Client) GetSupported(params files_sdk.HolidayRegionGetSupportedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func GetSupported(params files_sdk.HolidayRegionGetSupportedParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.HolidayRegion` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.HolidayRegionGetSupportedParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := holiday_region.GetSupported(files_sdk.HolidayRegionGetSupportedParams{})
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

## Model

`files_sdk.HolidayRegion` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Code` | `string` | `code` | The code representing a region |
| `Name` | `string` | `name` | The name of the region |
