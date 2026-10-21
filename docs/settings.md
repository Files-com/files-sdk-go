# Settings

`files_sdk.Settings` contains the API settings used by desktop clients. Its client package,
`github.com/Files-com/files-sdk-go/v3/settings`, named `settings`, retrieves them.

The package function uses the default configuration, including the `FILES_API_KEY` environment variable, and the
`Client` method of the same name uses the `files_sdk.Config` that `Client` embeds. The example runs inside a function
that returns an `error`, with these imports:

```go
import (
	"fmt"

	"github.com/Files-com/files-sdk-go/v3/settings"
)
```

## Operations

### Get

`GET /settings`

```go
func (c Client) Get(opts ...files_sdk.RequestResponseOption) (files_sdk.Settings, error)
func Get(opts ...files_sdk.RequestResponseOption) (files_sdk.Settings, error)
```

`Get` takes only request options: settings are not identified by an ID or path.

Example:

```go
result, err := settings.Get()
if err != nil {
	return err
}
fmt.Println(result)
```

## Model

`files_sdk.Settings` has these fields:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `DesktopDriveMappings` | `map[string]string` | `desktop_drive_mappings` |  |
| `DisableDriveMounting` | `bool` | `disable_drive_mounting` |  |
