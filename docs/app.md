# App

`files_sdk.App` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/app`, named `app`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/app"
)
```

## Operations

### List

`GET /apps`

List Apps

```go
func (c *Client) List(params files_sdk.AppListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.AppListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.App` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.AppListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name` and `app_type`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `name` and `app_type`. Valid field combinations are `[ name, app_type ]`. |
| `FilterPrefix` | `interface{}` | `filter_prefix` | no | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := app.List(files_sdk.AppListParams{})
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

`files_sdk.App` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `AppType` | `string` | `app_type` | The type of the App |
| `DocumentationLinks` | `interface{}` | `documentation_links` | Collection of named links to documentation |
| `ExtendedDescription` | `string` | `extended_description` | Long description for the in-App landing page |
| `ExtendedDescriptionForMarketingSite` | `string` | `extended_description_for_marketing_site` | Long form description of the App |
| `ExternalHomepageUrl` | `string` | `external_homepage_url` | Link to external homepage |
| `Featured` | `*bool` | `featured` | Is featured on the App listing? |
| `FolderBehaviorType` | `string` | `folder_behavior_type` | Associated Folder Behavior type, if any |
| `IconUrl` | `string` | `icon_url` | App icon |
| `LogoThumbnailUrl` | `string` | `logo_thumbnail_url` | Logo thumbnail for the App |
| `LogoUrl` | `string` | `logo_url` | Full size logo for the App |
| `MarketingIntro` | `string` | `marketing_intro` | Marketing introdution of the App |
| `MarketingYoutubeUrl` | `string` | `marketing_youtube_url` | Marketing video page |
| `Name` | `string` | `name` | Name of the App |
| `PackageManagerInstallCommand` | `string` | `package_manager_install_command` | Package manager install command |
| `RemoteServerType` | `string` | `remote_server_type` | Associated Remote Server type, if any |
| `ScreenshotListUrls` | `[]string` | `screenshot_list_urls` | Screenshots of the App |
| `SdkInstallationInstructionsLink` | `string` | `sdk_installation_instructions_link` | Link to SDK installation instructions |
| `ShortDescription` | `string` | `short_description` | Short description of the App |
| `SsoStrategyType` | `string` | `sso_strategy_type` | Associated SSO Strategy type, if any |
| `SiemType` | `string` | `siem_type` | Associated SIEM type, if any |
| `TutorialYoutubeUrl` | `string` | `tutorial_youtube_url` | Tutorial video page |
