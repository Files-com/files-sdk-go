# Behavior

`files_sdk.Behavior` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/behavior`, named `behavior`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/behavior"
)
```

## Operations

### List

`GET /behaviors`

List Behaviors

```go
func (c *Client) List(params files_sdk.BehaviorListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.BehaviorListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Behavior` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.BehaviorListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `behavior`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `clickwrap_id`, `form_field_set_id`, `impacts_ui`, `remote_server_id` or `behavior`. Valid field combinations are `[ impacts_ui, behavior ]` and `[ behavior, remote_server_id ]`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := behavior.List(files_sdk.BehaviorListParams{})
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

`GET /behaviors/{id}`

Show Behavior

```go
func (c *Client) Find(params files_sdk.BehaviorFindParams, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error)
func Find(params files_sdk.BehaviorFindParams, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error)
```

Parameters, `files_sdk.BehaviorFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Behavior ID. |

Example:

```go
result, err := behavior.Find(files_sdk.BehaviorFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### ListFor

`GET /behaviors/folders/{path}`

List Behaviors by Path

```go
func (c *Client) ListFor(params files_sdk.BehaviorListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func ListFor(params files_sdk.BehaviorListForParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Behavior` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.BehaviorListForParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `behavior`. |
| `Filter` | `interface{}` | `filter` | no | If set, return records where the specified field is equal to the supplied value. Valid fields are `impacts_ui` and `behavior`. Valid field combinations are `[ impacts_ui, behavior ]`. |
| `Path` | `string` | `path`, in the path | yes | Path to operate on. |
| `AncestorBehaviors` | `*bool` | `ancestor_behaviors` | no | If `true`, behaviors above this path are shown. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := behavior.ListFor(files_sdk.BehaviorListForParams{
	Path: "path",
})
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

### Create

`POST /behaviors`

Create Behavior

```go
func (c *Client) Create(params files_sdk.BehaviorCreateParams, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error)
func Create(params files_sdk.BehaviorCreateParams, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error)
```

Parameters, `files_sdk.BehaviorCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Value` | `interface{}` | `value` | no | This field stores data specific to the type of behavior. See The Behavior Types section for the accepted value for each type of behavior. |
| `AttachmentFile` | `io.Writer` | `attachment_file` | no | Certain behaviors may require a file, for instance, the `watermark` behavior requires a watermark image. Attach that file here. |
| `AttachmentFileAttachment` | `*lib.JSONAttachment` | `attachment_file` | no | Sends file content as attachment_file; see lib.JSONAttachment. Leave AttachmentFile nil when you set it. |
| `DisableParentFolderBehavior` | `*bool` | `disable_parent_folder_behavior` | no | If true, disables the inherited behavior for this folder and its children. Valid only for behavior types that child folders may override, and requires recursive to be true. Rejected for all other behavior types. |
| `Recursive` | `*bool` | `recursive` | no | Whether the behavior should apply to child folders. This is only configurable for behavior types whose recursion mode is `sometimes`; `always` behaviors stay recursive and `never` behaviors stay non-recursive. |
| `Name` | `string` | `name` | no | Name for this behavior. |
| `Description` | `string` | `description` | no | Description for this behavior. |
| `Path` | `string` | `path` | yes | Path where this behavior should apply. |
| `Behavior` | `string` | `behavior` | yes | Behavior type. |

Example:

```go
result, err := behavior.Create(files_sdk.BehaviorCreateParams{
	Path:     "path",
	Behavior: "webhook",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### WebhookTest

`POST /behaviors/webhook/test`

Test Webhook

```go
func (c *Client) WebhookTest(params files_sdk.BehaviorWebhookTestParams, opts ...files_sdk.RequestResponseOption) (err error)
func WebhookTest(params files_sdk.BehaviorWebhookTestParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.BehaviorWebhookTestParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Url` | `string` | `url` | yes | URL for testing the webhook. |
| `Method` | `string` | `method` | no | HTTP request method (GET or POST). |
| `Encoding` | `string` | `encoding` | no | Encoding type for the webhook payload. Can be JSON, XML, or RAW (form data). |
| `Headers` | `interface{}` | `headers` | no | Additional request headers to send via HTTP. |
| `Body` | `interface{}` | `body` | no | Additional body parameters to include in the webhook payload. |
| `Action` | `string` | `action` | no | Action for test body. |

Example:

```go
err := behavior.WebhookTest(files_sdk.BehaviorWebhookTestParams{
	Url: "https://www.site.com/...",
})
if err != nil {
	return err
}
```

### Update

`PATCH /behaviors/{id}`

Update Behavior

```go
func (c *Client) Update(params files_sdk.BehaviorUpdateParams, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error)
func Update(params files_sdk.BehaviorUpdateParams, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (behavior files_sdk.Behavior, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.BehaviorUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Behavior ID. |
| `Value` | `interface{}` | `value` | no | This field stores data specific to the type of behavior. See The Behavior Types section for the accepted value for each type of behavior. |
| `AttachmentFile` | `io.Writer` | `attachment_file` | no | Certain behaviors may require a file, for instance, the `watermark` behavior requires a watermark image. Attach that file here. |
| `AttachmentFileAttachment` | `*lib.JSONAttachment` | `attachment_file` | no | Sends file content as attachment_file; see lib.JSONAttachment. Leave AttachmentFile nil when you set it. |
| `DisableParentFolderBehavior` | `*bool` | `disable_parent_folder_behavior` | no | If true, disables the inherited behavior for this folder and its children. Valid only for behavior types that child folders may override, and requires recursive to be true. Rejected for all other behavior types. |
| `Recursive` | `*bool` | `recursive` | no | Whether the behavior should apply to child folders. This is only configurable for behavior types whose recursion mode is `sometimes`; `always` behaviors stay recursive and `never` behaviors stay non-recursive. |
| `Name` | `string` | `name` | no | Name for this behavior. |
| `Description` | `string` | `description` | no | Description for this behavior. |
| `AttachmentDelete` | `*bool` | `attachment_delete` | no | If `true`, delete the file stored in `attachment`. |

Example:

```go
result, err := behavior.Update(files_sdk.BehaviorUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := behavior.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /behaviors/{id}`

Delete Behavior

```go
func (c *Client) Delete(params files_sdk.BehaviorDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.BehaviorDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.BehaviorDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Behavior ID. |

Example:

```go
err := behavior.Delete(files_sdk.BehaviorDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Behavior` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Folder behavior ID |
| `Path` | `string` | `path` | Folder path. Note that Behavior paths cannot be updated once initially set. You will need to remove and re-create the behavior on the new path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `AttachmentUrl` | `string` | `attachment_url` | URL for attached file |
| `Behavior` | `string` | `behavior` | Behavior type. |
| `Name` | `string` | `name` | Name for this behavior. |
| `Description` | `string` | `description` | Description for this behavior. |
| `Value` | `interface{}` | `value` | Settings for this behavior. See the section above for an example value to provide here. Formatting is different for each Behavior type. May be sent as nested JSON or a single JSON-encoded string. If using XML encoding for the API call, this data must be sent as a JSON-encoded string. |
| `PublicHostingUrl` | `string` | `public_hosting_url` | Public URL for this publicly hosted folder when the `Serve Publicly` behavior has a key configured. When a Custom Domain with `public_hosting` destination is attached to this behavior, the URL uses that domain. Otherwise it uses the site's `subdomain.hosted-by-files.com` host, unless `custom_domain_only` is enabled, in which case no URL is returned. |
| `DisableParentFolderBehavior` | `*bool` | `disable_parent_folder_behavior` | If true, disables the inherited behavior for this folder and its children. Valid only for behavior types that child folders may override, and requires recursive to be true. Rejected for all other behavior types. |
| `Recursive` | `*bool` | `recursive` | Whether this behavior is recursive for this record. `always` behaviors are always `true`, `never` behaviors are always `false`, and `sometimes` behaviors may be either value. |
| `Inherited` | `*bool` | `inherited` | If true, this behavior is inherited from a higher scope rather than owned by the requested workspace. |
| `Managed` | `*bool` | `managed` | If true, this behavior is controlled by a parent-site policy and cannot be modified locally. |
| `RootBehaviorSiteAdminOnly` | `*bool` | `root_behavior_site_admin_only` | If true, this behavior may only be modified by a site admin because it is at the site root or disables a root behavior. |
| `AttachmentFile` | `io.Reader` | `attachment_file` | Certain behaviors may require a file, for instance, the `watermark` behavior requires a watermark image. Attach that file here. |
| `AttachmentDelete` | `*bool` | `attachment_delete` | If `true`, delete the file stored in `attachment`. |
