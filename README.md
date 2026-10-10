# Files.com Go Client

The Files.com Go SDK lets Go applications manage files, folders, users, and other resources through the Files.com API.

Use the examples below to install the SDK, authenticate, and make API requests.
The [developer documentation](https://developers.files.com/go/) covers the complete API.

## Introduction

The Files.com Go SDK lets Go applications manage files, folders, users, automations,
and other API resources. Resource clients accept configuration and request options and use standard Go
error handling.

### Installation

If your project does not already have a `go.mod` file, create a module. Replace
`example.com/myapp` with your module path. Then add the SDK dependency:

```shell
go mod init example.com/myapp
go get github.com/Files-com/files-sdk-go/v3
```

Import the SDK models and the resource clients you need:

``` go
import (
    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/folder"
)
```

The examples show imports followed by code to place inside a function, such as
`main`. Use `go build` to compile your application after adding the SDK dependency.

### Files.com is Committed to Go

Files.com uses Go for internal development. The Go SDK also supports the Files.com
CLI, desktop application, Terraform provider, and Rclone integration.

Explore the [files-sdk-go](https://github.com/Files-com/files-sdk-go) code on GitHub.

### Getting Support

The Files.com Support team provides official support for all of our official Files.com integration tools.

To initiate a support conversation, you can send an [Authenticated Support Request](https://www.files.com/docs/overview/requesting-support) or simply send an E-Mail to support@files.com.

## Authentication

There are two ways to authenticate: API Key authentication and Session-based authentication.

### Authenticate with an API Key

Authenticating with an API key is the recommended authentication method for most scenarios, and is
the method used in the examples on this site.

To use an API Key, first generate an API key from the [web
interface](https://www.files.com/docs/sdk-and-apis/api-keys) or [via the API or an
SDK](https://developers.files.com/go/resources/developers/api-keys).

A user-specific API key uses the user's permissions. [Workspaces](https://developers.files.com/go/overview/workspaces) describes account scope and administrative access.

Set `Config.APIKey` and initialize the configuration before passing it to a resource client.
This configuration applies to that client.

```go title="API Key Authentication"
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/folder"
)

config := files_sdk.Config{APIKey: "YOUR_API_KEY"}.Init()
client := folder.Client{Config: config}
it, err := client.ListFor(files_sdk.FolderListForParams{})
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

for file, err := range it.All() {
    if err != nil {
        var respErr files_sdk.ResponseError
        if errors.As(err, &respErr) {
            fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
        } else {
            fmt.Printf("Request failed: %v\n", err)
        }
        return
    }
    fmt.Println(file.Path)
}
```

You can reuse `files_sdk.GlobalConfig` by passing it to a client, for example
`folder.Client{Config: files_sdk.GlobalConfig}`.

To use the `FILES_API_KEY` environment variable, leave `Config.APIKey` empty.
Package-level resource functions also read this variable:

```go title="API Key From the Environment"
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/folder"
)

it, err := folder.ListFor(files_sdk.FolderListForParams{})
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

for file, err := range it.All() {
    if err != nil {
        var respErr files_sdk.ResponseError
        if errors.As(err, &respErr) {
            fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
        } else {
            fmt.Printf("Request failed: %v\n", err)
        }
        return
    }
    fmt.Println(file.Path)
}
```

Don't forget to replace the placeholder, `YOUR_API_KEY`, with your actual API key.

### Authenticate with a Session

Create a session with the username and password of an active user. The session uses that user's permissions. [Workspaces](https://developers.files.com/go/overview/workspaces) describes account scope and administrative access.

Sessions follow the same timeout settings as web sessions. When a session expires,
create a new session and update the clients that use it. The SDK does not renew
sessions automatically.

#### Logging In

Create a `session.Client` with your Files.com subdomain, then call `Create` with
your username and password. Use the returned session ID as `Config.SessionId`
when creating resource clients.

```go title="Example Request"
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/folder"
    "github.com/Files-com/files-sdk-go/v3/session"
)

sessionClient := session.Client{Config: files_sdk.Config{Subdomain: "MY-SUBDOMAIN"}.Init()}
thisSession, err := sessionClient.Create(files_sdk.SessionCreateParams{Username: "USERNAME", Password: "PASSWORD"})
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

config := files_sdk.Config{Subdomain: "MY-SUBDOMAIN", SessionId: thisSession.Id}.Init()
folderClient := folder.Client{Config: config}

it, err := folderClient.ListFor(files_sdk.FolderListForParams{})
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

for file, err := range it.All() {
    if err != nil {
        var respErr files_sdk.ResponseError
        if errors.As(err, &respErr) {
            fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
        } else {
            fmt.Printf("Request failed: %v\n", err)
        }
        return
    }
    fmt.Println(file.Path)
}
```

#### Using a Session

Set `Config.SessionId` to the session's `Id` and pass the configuration to each
resource client that uses the session.

```go title="Example Request"
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/folder"
)

config := files_sdk.Config{Subdomain: "MY-SUBDOMAIN", SessionId: thisSession.Id}.Init()
folderClient := folder.Client{Config: config}

it, err := folderClient.ListFor(files_sdk.FolderListForParams{})
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

for file, err := range it.All() {
    if err != nil {
        var respErr files_sdk.ResponseError
        if errors.As(err, &respErr) {
            fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
        } else {
            fmt.Printf("Request failed: %v\n", err)
        }
        return
    }
    fmt.Println(file.Path)
}
```

#### Logging Out

User sessions can be ended by calling `Delete()` on the `Session` client.

```go title="Example Request"
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/session"
)

sessionClient := session.Client{Config: files_sdk.Config{SessionId: thisSession.Id}.Init()}
err := sessionClient.Delete()
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}
```

## Configuration

Initialize a `files_sdk.Config` with `Init()` and pass it to each resource client.
The configuration controls authentication, the API endpoint, workspace scoping,
and HTTP behavior for that client.

### Configuration Options

#### Base URL

Set this to the full https:// URL of your Files.com subdomain (e.g. `https://MY-SUBDOMAIN.files.com`).
This is not required in most cases, but one benefit of setting it is that it ensures that authentication failures will be logged to your site's API logs.  Without setting this, we won't know which site to associate the authentication failure with, and it won't be logged to your site's API logs.
This is always required if your site is configured to disable global acceleration.
This can also be set to use a mock server in development or CI.

```go title="Example setting"
import (
    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/file"
)

config := files_sdk.Config{
    EndpointOverride: "https://MY-SUBDOMAIN.files.com",
}.Init()
client := file.Client{Config: config}
```

#### Diagnostic Logging

Set `Config.Debug` to `true`, or set the `FILES_SDK_DEBUG` environment variable
to a non-empty value, to enable detailed SDK request and response diagnostics.
Provide `Config.Logger` to receive the output. The default logger discards it.
These DEBUG diagnostics can include signed URLs, API keys, session IDs, headers,
and payloads.

SDK retry entries at INFO, WARN, and ERROR redact request URLs and credentials,
even when debug mode is enabled. A logger with levels controls which retry DEBUG
entries it displays. For a logger that only implements `Printf`, debug mode
controls whether retry DEBUG entries are emitted.

## Sort and Filter

Several of the Files.com API resources have list operations that return multiple instances of the
resource. The List operations can be sorted and filtered.

### Sorting

To sort the returned data, pass in the ```SortBy``` method argument.

Each resource supports a unique set of valid sort fields and can only be sorted by one field at a
time.

The argument value is a Go ```map[string]interface{}``` map that has a key of the resource field
name to sort on and a value of either ```"asc"``` or ```"desc"``` to specify the sort order.

#### Special note about the List Folder Endpoint

For historical reasons, and to maintain compatibility
with a variety of other cloud-based MFT and EFSS services, Folders will always be listed before Files
when listing a Folder.  This applies regardless of the sorting parameters you provide.  These *will* be
used, after the initial sort application of Folders before Files.

```go title="Sort Example"
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/user"
)

client := user.Client{Config: files_sdk.GlobalConfig}

// users sorted by username
parameters := files_sdk.UserListParams{SortBy: map[string]interface{}{"username":"asc"}}
userIterator, err := client.List(parameters)
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

for user, err := range userIterator.All() {
    if err != nil {
        var respErr files_sdk.ResponseError
        if errors.As(err, &respErr) {
            fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
        } else {
            fmt.Printf("Request failed: %v\n", err)
        }
        return
    }
    fmt.Println(user.Username)
}
```

### Filtering

Filters apply selection criteria to the underlying query that returns the results. They can be
applied individually or combined with other filters, and the resulting data can be sorted by a
single field.

Each resource supports a unique set of valid filter fields, filter combinations, and combinations of
filters and sort fields.

The passed in argument is a Go ```map[string]interface{}``` map that has
a key of the resource field name to filter on and a passed in value to use in the filter comparison.

#### Filter Types

| Filter | Type | Description |
| --------- | --------- | --------- |
| `Filter` | Exact | Find resources that have an exact field value match to a passed in value. (i.e., FIELD_VALUE = PASS_IN_VALUE). |
| `FilterPrefix` | Pattern | Find resources where the specified field is prefixed by the supplied value. This is applicable to values that are strings. |
| `FilterGt` | Range | Find resources that have a field value that is greater than the passed in value.  (i.e., FIELD_VALUE > PASS_IN_VALUE). |
| `FilterGteq` | Range | Find resources that have a field value that is greater than or equal to the passed in value.  (i.e., FIELD_VALUE >=  PASS_IN_VALUE). |
| `FilterLt` | Range | Find resources that have a field value that is less than the passed in value.  (i.e., FIELD_VALUE < PASS_IN_VALUE). |
| `FilterLteq` | Range | Find resources that have a field value that is less than or equal to the passed in value.  (i.e., FIELD_VALUE \<= PASS_IN_VALUE). |

```go title="Exact Filter Example"
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/user"
)

client := user.Client{Config: files_sdk.GlobalConfig}

// non admin users
parameters := files_sdk.UserListParams{
    Filter: map[string]interface{}{"not_site_admin": true},
}
userIterator, err := client.List(parameters)
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

for user, err := range userIterator.All() {
    if err != nil {
        var respErr files_sdk.ResponseError
        if errors.As(err, &respErr) {
            fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
        } else {
            fmt.Printf("Request failed: %v\n", err)
        }
        return
    }
    fmt.Println(user.Username)
}
```

```go title="Range Filter Example"
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/user"
)

client := user.Client{Config: files_sdk.GlobalConfig};

// users who haven't logged in since 2024-01-01
parameters := files_sdk.UserListParams{
    FilterLt: map[string]interface{}{"last_login_at": "2024-01-01"},
}
userIterator, err := client.List(parameters)
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

for user, err := range userIterator.All() {
    if err != nil {
        var respErr files_sdk.ResponseError
        if errors.As(err, &respErr) {
            fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
        } else {
            fmt.Printf("Request failed: %v\n", err)
        }
        return
    }
    fmt.Println(user.Username)
}
```

```go title="Pattern Filter Example"
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/user"
)

client := user.Client{Config: files_sdk.GlobalConfig};

// users whose usernames start with 'test'
parameters := files_sdk.UserListParams{
    FilterPrefix: map[string]interface{}{"username": "test"},
}
userIterator, err := client.List(parameters)
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

for user, err := range userIterator.All() {
    if err != nil {
        var respErr files_sdk.ResponseError
        if errors.As(err, &respErr) {
            fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
        } else {
            fmt.Printf("Request failed: %v\n", err)
        }
        return
    }
    fmt.Println(user.Username)
}
```

```go title="Combination Filter with Sort Example"
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/user"
)

client := user.Client{Config: files_sdk.GlobalConfig};

// users whose usernames start with 'test' and are not admins
parameters := files_sdk.UserListParams{
    FilterPrefix: map[string]interface{}{"username": "test"},
    Filter:       map[string]interface{}{"not_site_admin": true},
    SortBy:       map[string]interface{}{"username": "asc"},
}
userIterator, err := client.List(parameters)
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

for user, err := range userIterator.All() {
    if err != nil {
        var respErr files_sdk.ResponseError
        if errors.As(err, &respErr) {
            fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
        } else {
            fmt.Printf("Request failed: %v\n", err)
        }
        return
    }
    fmt.Println(user.Username)
}
```

## Paths

Files.com preserves the spelling of file and folder paths while comparing them using shared case and Unicode rules. Use the SDK comparison helpers when matching paths locally.
<div></div>

### Capitalization

Files.com uses case-insensitive path matching based on its fixed Unicode comparison map.

For example, the following paths have the same comparison key:

| Path Variant                          | Comparison Key              |
|---------------------------------------|------------------------------|
| `Documents/Reports/Q1.pdf`            | `documents/reports/q1.pdf`  |
| `documents/reports/q1.PDF`            | `documents/reports/q1.pdf`  |
| `DOCUMENTS/REPORTS/Q1.PDF`            | `documents/reports/q1.pdf`  |

This behavior applies across:
- API requests
- Folder and file lookup operations
- Automations and workflows

See also: [Case Sensitivity Documentation](https://www.files.com/docs/files-and-folders/case-sensitivity/)

### Slashes

Use `/` between folder and file names, without leading or trailing slashes. SDK normalization helpers convert backslashes to `/`, remove duplicate separators, and discard exact `.` and `..` components. Discarding `..` leaves the preceding folder name intact.

| Input | Normalized path |
|-------|-----------------|
| `folder/subfolder/file.txt` | `folder/subfolder/file.txt` |
| `/folder/subfolder/file.txt` | `folder/subfolder/file.txt` |
| `folder/subfolder/file.txt/` | `folder/subfolder/file.txt` |
| `//folder//file.txt` | `folder/file.txt` |
| `folder/../file.txt` | `folder/file.txt` |

<div></div>

### Unicode and Path Comparison

Files.com compares paths using a fixed mapping shared by the server and SDKs. It treats case and many accent differences as equivalent: `Résumé.txt` and `resume.txt` identify the same file, as do `q` followed by a combining acute accent and `q`. The mapping also handles other equivalences, such as Hiragana and Katakana. Lowercasing or applying a standard Unicode normalization form alone does not reproduce these rules.

SDK comparison helpers normalize path separators and dot segments, then apply the bundled [versioned comparison map](https://github.com/Files-com/files-sdk-javascript/blob/master/shared/path_comparison.json). The [shared examples](https://github.com/Files-com/files-sdk-javascript/blob/master/shared/comparison_examples.json) give exact comparison results for integrations that implement their own matching. The map uses hexadecimal Unicode scalar values as keys: a missing entry preserves the character, an empty replacement removes it, and other replacements may contain several characters. Apply each replacement once without normalizing or lowercasing the result again.

Use comparison results only for matching. Send the original path spelling in API requests and preserve it for display and local filenames; comparison results can have a different spelling or length.

Trailing whitespace is significant for comparison. `report.txt` and `report.txt ` are different file paths, and SDK helpers preserve spaces, tabs, and newlines. Folder names cannot end in whitespace. See [Unicode Normalization](https://www.files.com/docs/files-and-folders/file-system-semantics/unicode-normalization) for the complete path rules.

<div></div>

## Workspaces

A Workspace groups files, users, groups, Partners, integrations, and workflows within a Files.com Site. An integration can provision a Workspace for a department or project and delegate its operation to a team without making that team Site Administrators. Every Site has a Default Workspace, with ID `0`; additional Workspaces have their own IDs and root folders.

Account membership, request context, and permission grants serve different purposes. Creating an account in a Workspace determines where it belongs. Selecting a Workspace determines which resources a request operates on. A permission grant determines what the caller can do there. Selecting a Workspace never grants access to it.

### Accounts and Administrative Access

A user's or group's `workspace_id` identifies the Workspace the account belongs to. Accounts belonging to a Custom Workspace stay within it. Default Workspace users and groups can receive permissions in one or more Custom Workspaces while keeping their existing accounts in Workspace `0`.

| Account | Workspace Administrator assignment | Scope |
| --- | --- | --- |
| User belonging to a Custom Workspace | Set the user's `workspace_admin` to `true`. | That user's own Custom Workspace. |
| Default Workspace user | Create an `admin` Permission for the user on a Custom Workspace's root folder. | Each Custom Workspace with a root grant. |
| Default Workspace group | Create an `admin` Permission for the group on a Custom Workspace's root folder. | Every member inherits administration of each Workspace with a root grant. |

`workspace_admin` is not a summary of a user's effective administrative access. A Default Workspace user can administer a Custom Workspace through a direct or group root grant while their `workspace_admin` remains `false`. Groups have no `workspace_admin` field. See [Users](https://developers.files.com/go/resources/user-accounts/users) and [Groups](https://developers.files.com/go/resources/user-accounts/groups) for account fields.

An `admin` grant on the **Custom Workspace root** provides full Workspace Administrator authority over its files, users, groups, Partners, workflows, and integrations. An `admin` grant on a subfolder provides Folder Admin authority over that folder and its descendants; it does not provide Workspace administration. Other permission levels provide their corresponding folder access without Workspace administration. [Permissions](https://developers.files.com/go/resources/user-accounts/permissions) defines the levels.

Site Administrators manage cross-Workspace assignments to Default Workspace accounts. Workspace Administrators manage accounts and permissions within their own scope. Site Administrators retain access to every Workspace; adding a Workspace grant does not narrow Site Administrator authority. The [product documentation](https://www.files.com/docs/workspaces/workspace-administrators) explains the administrator's operational scope and site-wide controls.

### Request Context and API Keys

You can include the `X-Files-Workspace-Id` REST header to select a Workspace for a request. SDK request options and CLI configuration send that same selection. When a Workspace is selected, Workspace-scoped resources are listed, created, and changed within that context, and ordinary paths are relative to its root.

A resource's `workspace_id` request field describes the resource's Workspace membership. It is separate from the SDK's Workspace request option or REST header. Creating a Workspace-scoped resource in a Custom Workspace defaults its `workspace_id` to the selected Workspace; a mismatching membership value is rejected with `not-authorized/insufficient-permission-for-params`.

Selecting another Workspace with an API key requires a **Full Access key created in the Default Workspace**. A user key follows that user's current access, including group permissions. A site-wide Full Access key created in the Default Workspace has Site Administrator authority in every Workspace. A Files Only key stays in its creation Workspace, even if its user has cross-Workspace access. Any key created in a Custom Workspace stays within that Workspace. Selecting another context with these confined keys is rejected with `bad-request/invalid-workspace-id-header`.

An account belonging to a Custom Workspace is scoped there when it authenticates normally. For a Default Workspace user, explicitly select the intended Workspace for an integration rather than relying on an interactive login preference. [API Keys](https://developers.files.com/go/resources/developers/api-keys) and [Authentication](https://developers.files.com/go/overview/authentication) cover credentials.

The Files.com Go SDK supports workspace scoping by using the `WorkspaceId` attribute on the `Config` object.

The adjacent scoping example uses a credential authorized for the selected Workspace. A group member uses their own Full Access user key from the Default Workspace; the Site Administrator credential used to assign the grant is not needed for their day-to-day work.

```go title="Example Request"
import (
    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/folder"
)

config := files_sdk.Config{WorkspaceId: 123}.Init()
client := folder.Client{Config: config}

client.ListFor(files_sdk.FolderListForParams{})
```

### Delegating a Workspace to an Existing Group

An operations team already represented by a Default Workspace group can administer a Custom Workspace through one root Permission. The group and its members stay in the Default Workspace, so the same team can receive different access in other Workspaces.

First retrieve the target [Workspace](https://developers.files.com/go/resources/settings/workspaces) and [Group](https://developers.files.com/go/resources/user-accounts/groups) IDs as a Site Administrator in Workspace `0`. The examples use Workspace `123`, group `456`, and member user `789`; replace them with your own IDs. Confirm that the group belongs to Workspace `0` and that the intended user is a member.

Create the Permission using a Default Workspace Full Access site-wide key or a Full Access user key belonging to a Site Administrator. Keep the request context at `0` and use the qualified root path `_/Workspaces/123`. Set `group_id` to the group's ID, `permission` to `admin`, and `recursive` to `true`. Save the returned Permission `id` for later removal. For an individual Default Workspace user, use `user_id` instead of `group_id`.

For a Default Workspace group, a Site Administrator can also select Workspace `123` and use an empty `path` to grant access to its root. The qualified path in Workspace `0` works for both Default Workspace users and groups and keeps the account scope and target Workspace explicit. Appending a subfolder to the path would grant Folder Admin access instead of Workspace Administrator authority.

After the grant, make a request as the member, using their own credential with Workspace `123` selected, such as listing that Workspace's root folder. That member can work with the Workspace's files and perform Workspace Administrator operations, such as managing its users, Partners, and integrations. A Site Administrator's successful request does not establish that the member has the intended access.

```go title="Grant group administration"
import (
    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/lib"
    "github.com/Files-com/files-sdk-go/v3/permission"
)

config := files_sdk.Config{ApiKey: "YOUR_SITE_ADMIN_API_KEY", WorkspaceId: 0}.Init()
permissions := permission.Client{Config: config}
grant, err := permissions.Create(files_sdk.PermissionCreateParams{
    Path: "_/Workspaces/123", GroupId: 456,
    Permission: "admin", Recursive: lib.Bool(true),
})
if err != nil { panic(err) }
println(grant.Id) // Save this Permission ID for removal.
```

### Permission Inspection and Removal

List the member's Permissions with `user_id` and `include_groups=true` to include grants inherited through group membership. Listing only direct user grants can miss the Permission that provides Workspace administration. In Workspace `0`, the Custom Workspace root appears as `_/Workspaces/123`; in Workspace `123`, paths are relative to that root. Inspect the root path and `permission=admin`, rather than treating the user's `workspace_admin` field as their effective administrative access.

Permission lists show individual grants, rather than a single flag for effective administrative access. Membership in several groups combines their access. A Permission using `group_ids` instead of `group_id` requires membership in all the specified groups; it is not a shorthand for assigning the same grant to several independent groups.

Removing a member ends access received through that group. Deleting the root Permission ends the group's Workspace Administrator grant for every member. These changes leave independent direct and other group grants in place, so review all applicable grants when withdrawing access. Default Workspace user API keys follow those permission changes without being recreated.

Group membership maintained through SCIM follows the same rule. A Group Admin allowed to add members can give those users the group's existing Workspace Administrator access. Choose who manages the group with that authority in mind.

Delete the Permission by its returned `id` as the Site Administrator in Workspace `0`. The removal examples use Permission ID `9001`; replace it with the ID returned by your create request. Permissions are created and deleted, rather than updated in place. If narrower folder access is still needed, assign it explicitly; deleting a broad grant does not restore narrower grants it previously replaced.

```go title="Inspect member grants and remove the group grant"
grants, err := permissions.List(files_sdk.PermissionListParams{
    UserId: "789", IncludeGroups: lib.Bool(true),
})
if err != nil { panic(err) }
for grants.Next() {
    item := grants.Permission()
    println(item.Path, item.Permission)
}
if err := grants.Err(); err != nil { panic(err) }
if err := permissions.Delete(files_sdk.PermissionDeleteParams{Id: 9001}); err != nil {
    panic(err)
}
```

## Foreign Language Support

The Files.com Go SDK supports localized responses by using the `Language` attribute on the `Config` struct.
When configured, this guides the API in selecting a preferred language for applicable response content.

Language support currently applies to select human-facing fields only, such as notification messages
and error descriptions.

If the specified language is not supported or the value is omitted, the API defaults to English.

```shell title="Example Request"
import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
)

files_sdk.GlobalConfig.Language = "es";
```

## Errors

SDK methods return an `error` using the standard Go pattern. Check it before
using a returned resource. Listing methods fetch pages during iteration. When
ranging over `All()`, check the loop's error before using the resource. If you use
`Next()`, check the iterator's `Err()` after the loop.

Errors can come from request preparation, HTTP transport, response decoding, or
the Files.com API. A structured API error is a `files_sdk.ResponseError`, which
implements the `error` interface and includes these fields:

- `Type`: the API error identifier.
- `Title`: a short description of the error.
- `ErrorMessage`: additional error details.

Use `errors.As` to inspect a `ResponseError`, including one wrapped by another
error. Use `errors.Is` with the SDK's error constants to match a type or group.

```go title="Example Error Handling"
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/session"
)

_, err := session.Create(files_sdk.SessionCreateParams{Username: "USERNAME", Password: "BADPASSWORD"})
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}
```

### ResponseError Types

Match an exact error type with a `ResponseErrorType` constant, or use a group
helper such as `IsNotFound`. The table below lists API types and their constants.

```go title="Example ResponseError Type Matching"
if err != nil {
    if errors.Is(err, files_sdk.ErrInvalidUsernameOrPassword) {
        fmt.Println("Bad username/password")
    } else if files_sdk.IsNotFound(err) {
        fmt.Println("Resource not found")
    }

    var responseErr files_sdk.ResponseError
    if errors.As(err, &responseErr) {
        fmt.Println(responseErr.Title)
    }
}
```

| Type | Go Error | Title |
| --------- | --------- | --------- |
| `bad-request` | `ErrBadRequest` | Bad Request |
| `bad-request/agent-upgrade-required` | `ErrAgentUpgradeRequired` | Agent Upgrade Required |
| `bad-request/attachment-too-large` | `ErrAttachmentTooLarge` | Attachment Too Large |
| `bad-request/cannot-download-directory` | `ErrCannotDownloadDirectory` | Cannot Download Directory |
| `bad-request/cant-move-with-multiple-locations` | `ErrCantMoveWithMultipleLocations` | Cant Move With Multiple Locations |
| `bad-request/datetime-parse` | `ErrDatetimeParse` | Datetime Parse |
| `bad-request/destination-same` | `ErrDestinationSame` | Destination Same |
| `bad-request/destination-site-mismatch` | `ErrDestinationSiteMismatch` | Destination Site Mismatch |
| `bad-request/does-not-support-sorting` | `ErrDoesNotSupportSorting` | Does Not Support Sorting |
| `bad-request/folder-must-not-be-a-file` | `ErrFolderMustNotBeAFile` | Folder Must Not Be A File |
| `bad-request/folders-not-allowed` | `ErrFoldersNotAllowed` | Folders Not Allowed |
| `bad-request/internal-general-error` | `ErrInternalGeneralError` | Internal General Error |
| `bad-request/invalid-body` | `ErrInvalidBody` | Invalid Body |
| `bad-request/invalid-cursor` | `ErrInvalidCursor` | Invalid Cursor |
| `bad-request/invalid-cursor-type-for-sort` | `ErrInvalidCursorTypeForSort` | Invalid Cursor Type For Sort |
| `bad-request/invalid-etags` | `ErrInvalidEtags` | Invalid Etags |
| `bad-request/invalid-filter-alias-combination` | `ErrInvalidFilterAliasCombination` | Invalid Filter Alias Combination |
| `bad-request/invalid-filter-field` | `ErrInvalidFilterField` | Invalid Filter Field |
| `bad-request/invalid-filter-param` | `ErrInvalidFilterParam` | Invalid Filter Param |
| `bad-request/invalid-filter-param-format` | `ErrInvalidFilterParamFormat` | Invalid Filter Param Format |
| `bad-request/invalid-filter-param-value` | `ErrInvalidFilterParamValue` | Invalid Filter Param Value |
| `bad-request/invalid-input-encoding` | `ErrInvalidInputEncoding` | Invalid Input Encoding |
| `bad-request/invalid-interface` | `ErrInvalidInterface` | Invalid Interface |
| `bad-request/invalid-oauth-provider` | `ErrInvalidOauthProvider` | Invalid Oauth Provider |
| `bad-request/invalid-path` | `ErrInvalidPath` | Invalid Path |
| `bad-request/invalid-return-to-url` | `ErrInvalidReturnToUrl` | Invalid Return To Url |
| `bad-request/invalid-search-query` | `ErrInvalidSearchQuery` | Invalid Search Query |
| `bad-request/invalid-sort-field` | `ErrInvalidSortField` | Invalid Sort Field |
| `bad-request/invalid-sort-filter-combination` | `ErrInvalidSortFilterCombination` | Invalid Sort Filter Combination |
| `bad-request/invalid-upload-offset` | `ErrInvalidUploadOffset` | Invalid Upload Offset |
| `bad-request/invalid-upload-part-gap` | `ErrInvalidUploadPartGap` | Invalid Upload Part Gap |
| `bad-request/invalid-upload-part-size` | `ErrInvalidUploadPartSize` | Invalid Upload Part Size |
| `bad-request/invalid-workspace-id-header` | `ErrInvalidWorkspaceIdHeader` | Invalid Workspace Id Header |
| `bad-request/method-not-allowed` | `ErrMethodNotAllowed` | Method Not Allowed |
| `bad-request/multiple-sort-params-not-allowed` | `ErrMultipleSortParamsNotAllowed` | Multiple Sort Params Not Allowed |
| `bad-request/no-valid-input-params` | `ErrNoValidInputParams` | No Valid Input Params |
| `bad-request/offset-upload-not-allowed-with-malware-scanning` | `ErrOffsetUploadNotAllowedWithMalwareScanning` | Offset Upload Not Allowed With Malware Scanning |
| `bad-request/part-number-too-large` | `ErrPartNumberTooLarge` | Part Number Too Large |
| `bad-request/path-cannot-have-trailing-whitespace` | `ErrPathCannotHaveTrailingWhitespace` | Path Cannot Have Trailing Whitespace |
| `bad-request/reauthentication-needed-fields` | `ErrReauthenticationNeededFields` | Reauthentication Needed Fields |
| `bad-request/request-body-too-large` | `ErrRequestBodyTooLarge` | Request Body Too Large |
| `bad-request/request-params-contain-invalid-character` | `ErrRequestParamsContainInvalidCharacter` | Request Params Contain Invalid Character |
| `bad-request/request-params-invalid` | `ErrRequestParamsInvalid` | Request Params Invalid |
| `bad-request/request-params-required` | `ErrRequestParamsRequired` | Request Params Required |
| `bad-request/search-all-on-child-path` | `ErrSearchAllOnChildPath` | Search All On Child Path |
| `bad-request/unrecognized-sort-index` | `ErrUnrecognizedSortIndex` | Unrecognized Sort Index |
| `bad-request/unsupported-currency` | `ErrUnsupportedCurrency` | Unsupported Currency |
| `bad-request/unsupported-http-response-format` | `ErrUnsupportedHttpResponseFormat` | Unsupported Http Response Format |
| `bad-request/unsupported-media-type` | `ErrUnsupportedMediaType` | Unsupported Media Type |
| `bad-request/user-id-invalid` | `ErrUserIdInvalid` | User Id Invalid |
| `bad-request/user-id-on-user-endpoint` | `ErrUserIdOnUserEndpoint` | User Id On User Endpoint |
| `bad-request/user-required` | `ErrUserRequired` | User Required |
| `not-authenticated/additional-authentication-required` | `ErrAdditionalAuthenticationRequired` | Additional Authentication Required |
| `not-authenticated/api-key-sessions-not-supported` | `ErrApiKeySessionsNotSupported` | Api Key Sessions Not Supported |
| `not-authenticated/authentication-required` | `ErrAuthenticationRequired` | Authentication Required |
| `not-authenticated/bundle-registration-code-failed` | `ErrBundleRegistrationCodeFailed` | Bundle Registration Code Failed |
| `not-authenticated/inbox-registration-code-failed` | `ErrInboxRegistrationCodeFailed` | Inbox Registration Code Failed |
| `not-authenticated/invalid-credentials` | `ErrInvalidCredentials` | Invalid Credentials |
| `not-authenticated/invalid-oauth` | `ErrInvalidOauth` | Invalid Oauth |
| `not-authenticated/invalid-or-expired-code` | `ErrInvalidOrExpiredCode` | Invalid Or Expired Code |
| `not-authenticated/invalid-session` | `ErrInvalidSession` | Invalid Session |
| `not-authenticated/invalid-username-or-password` | `ErrInvalidUsernameOrPassword` | Invalid Username Or Password |
| `not-authenticated/locked-out` | `ErrLockedOut` | Locked Out |
| `not-authenticated/lockout-region-mismatch` | `ErrLockoutRegionMismatch` | Lockout Region Mismatch |
| `not-authenticated/one-time-password-incorrect` | `ErrOneTimePasswordIncorrect` | One Time Password Incorrect |
| `not-authenticated/two-factor-authentication-error` | `ErrTwoFactorAuthenticationError` | Two Factor Authentication Error |
| `not-authenticated/two-factor-authentication-setup-expired` | `ErrTwoFactorAuthenticationSetupExpired` | Two Factor Authentication Setup Expired |
| `not-authorized/api-key-is-disabled` | `ErrApiKeyIsDisabled` | Api Key Is Disabled |
| `not-authorized/api-key-is-path-restricted` | `ErrApiKeyIsPathRestricted` | Api Key Is Path Restricted |
| `not-authorized/api-key-only-for-desktop-app` | `ErrApiKeyOnlyForDesktopApp` | Api Key Only For Desktop App |
| `not-authorized/api-key-only-for-file-operations` | `ErrApiKeyOnlyForFileOperations` | Api Key Only For File Operations |
| `not-authorized/api-key-only-for-mobile-app` | `ErrApiKeyOnlyForMobileApp` | Api Key Only For Mobile App |
| `not-authorized/api-key-only-for-office-integration` | `ErrApiKeyOnlyForOfficeIntegration` | Api Key Only For Office Integration |
| `not-authorized/billing-information-hidden` | `ErrBillingInformationHidden` | Billing Information Hidden |
| `not-authorized/billing-permission-required` | `ErrBillingPermissionRequired` | Billing Permission Required |
| `not-authorized/bundle-maximum-uses-reached` | `ErrBundleMaximumUsesReached` | Bundle Maximum Uses Reached |
| `not-authorized/bundle-permission-required` | `ErrBundlePermissionRequired` | Bundle Permission Required |
| `not-authorized/cannot-administer-higher-level-user` | `ErrCannotAdministerHigherLevelUser` | Cannot Administer Higher Level User |
| `not-authorized/cannot-login-while-using-key` | `ErrCannotLoginWhileUsingKey` | Cannot Login While Using Key |
| `not-authorized/cant-act-for-other-user` | `ErrCantActForOtherUser` | Cant Act For Other User |
| `not-authorized/contact-admin-for-password-change-help` | `ErrContactAdminForPasswordChangeHelp` | Contact Admin For Password Change Help |
| `not-authorized/files-agent-failed-authorization` | `ErrFilesAgentFailedAuthorization` | Files Agent Failed Authorization |
| `not-authorized/folder-admin-or-billing-permission-required` | `ErrFolderAdminOrBillingPermissionRequired` | Folder Admin Or Billing Permission Required |
| `not-authorized/folder-admin-permission-required` | `ErrFolderAdminPermissionRequired` | Folder Admin Permission Required |
| `not-authorized/full-permission-required` | `ErrFullPermissionRequired` | Full Permission Required |
| `not-authorized/history-permission-required` | `ErrHistoryPermissionRequired` | History Permission Required |
| `not-authorized/in-app-ai-assistant-unavailable` | `ErrInAppAiAssistantUnavailable` | In App Ai Assistant Unavailable |
| `not-authorized/insufficient-permission-for-params` | `ErrInsufficientPermissionForParams` | Insufficient Permission For Params |
| `not-authorized/insufficient-permission-for-site` | `ErrInsufficientPermissionForSite` | Insufficient Permission For Site |
| `not-authorized/mover-access-denied` | `ErrMoverAccessDenied` | Mover Access Denied |
| `not-authorized/mover-package-required` | `ErrMoverPackageRequired` | Mover Package Required |
| `not-authorized/must-authenticate-with-api-key` | `ErrMustAuthenticateWithApiKey` | Must Authenticate With Api Key |
| `not-authorized/need-admin-permission-for-inbox` | `ErrNeedAdminPermissionForInbox` | Need Admin Permission For Inbox |
| `not-authorized/non-admins-must-query-by-folder-or-path` | `ErrNonAdminsMustQueryByFolderOrPath` | Non Admins Must Query By Folder Or Path |
| `not-authorized/not-allowed-to-create-bundle` | `ErrNotAllowedToCreateBundle` | Not Allowed To Create Bundle |
| `not-authorized/not-enqueuable-sync` | `ErrNotEnqueuableSync` | Not Enqueuable Sync |
| `not-authorized/password-change-not-required` | `ErrPasswordChangeNotRequired` | Password Change Not Required |
| `not-authorized/password-change-required` | `ErrPasswordChangeRequired` | Password Change Required |
| `not-authorized/payment-method-error` | `ErrPaymentMethodError` | Payment Method Error |
| `not-authorized/preview-only-permission-cannot-download` | `ErrPreviewOnlyPermissionCannotDownload` | Preview Only Permission Cannot Download |
| `not-authorized/read-only-session` | `ErrReadOnlySession` | Read Only Session |
| `not-authorized/read-permission-required` | `ErrReadPermissionRequired` | Read Permission Required |
| `not-authorized/reauthentication-failed` | `ErrReauthenticationFailed` | Reauthentication Failed |
| `not-authorized/reauthentication-failed-final` | `ErrReauthenticationFailedFinal` | Reauthentication Failed Final |
| `not-authorized/reauthentication-needed-action` | `ErrReauthenticationNeededAction` | Reauthentication Needed Action |
| `not-authorized/recaptcha-failed` | `ErrRecaptchaFailed` | Recaptcha Failed |
| `not-authorized/remote-desktop-debug-logging-disabled` | `ErrRemoteDesktopDebugLoggingDisabled` | Remote Desktop Debug Logging Disabled |
| `not-authorized/root-folder-behavior-site-admin-required` | `ErrRootFolderBehaviorSiteAdminRequired` | Root Folder Behavior Site Admin Required |
| `not-authorized/root-folder-behavior-skip-site-admin-required` | `ErrRootFolderBehaviorSkipSiteAdminRequired` | Root Folder Behavior Skip Site Admin Required |
| `not-authorized/self-managed-required` | `ErrSelfManagedRequired` | Self Managed Required |
| `not-authorized/site-admin-or-partner-admin-permission-required` | `ErrSiteAdminOrPartnerAdminPermissionRequired` | Site Admin Or Partner Admin Permission Required |
| `not-authorized/site-admin-or-workspace-admin-or-folder-admin-permission-required` | `ErrSiteAdminOrWorkspaceAdminOrFolderAdminPermissionRequired` | Site Admin Or Workspace Admin Or Folder Admin Permission Required |
| `not-authorized/site-admin-or-workspace-admin-or-partner-admin-or-folder-admin-permission-required` | `ErrSiteAdminOrWorkspaceAdminOrPartnerAdminOrFolderAdminPermissionRequired` | Site Admin Or Workspace Admin Or Partner Admin Or Folder Admin Permission Required |
| `not-authorized/site-admin-or-workspace-admin-or-partner-admin-permission-required` | `ErrSiteAdminOrWorkspaceAdminOrPartnerAdminPermissionRequired` | Site Admin Or Workspace Admin Or Partner Admin Permission Required |
| `not-authorized/site-admin-or-workspace-admin-permission-required` | `ErrSiteAdminOrWorkspaceAdminPermissionRequired` | Site Admin Or Workspace Admin Permission Required |
| `not-authorized/site-admin-required` | `ErrSiteAdminRequired` | Site Admin Required |
| `not-authorized/site-files-are-immutable` | `ErrSiteFilesAreImmutable` | Site Files Are Immutable |
| `not-authorized/two-factor-authentication-required` | `ErrTwoFactorAuthenticationRequired` | Two Factor Authentication Required |
| `not-authorized/user-id-without-site-admin` | `ErrUserIdWithoutSiteAdmin` | User Id Without Site Admin |
| `not-authorized/write-and-bundle-permission-required` | `ErrWriteAndBundlePermissionRequired` | Write And Bundle Permission Required |
| `not-authorized/write-permission-required` | `ErrWritePermissionRequired` | Write Permission Required |
| `not-found` | `ErrNotFound` | Not Found |
| `not-found/api-key-not-found` | `ErrApiKeyNotFound` | Api Key Not Found |
| `not-found/bundle-path-not-found` | `ErrBundlePathNotFound` | Bundle Path Not Found |
| `not-found/bundle-registration-not-found` | `ErrBundleRegistrationNotFound` | Bundle Registration Not Found |
| `not-found/code-not-found` | `ErrCodeNotFound` | Code Not Found |
| `not-found/file-not-found` | `ErrFileNotFound` | File Not Found |
| `not-found/file-upload-not-found` | `ErrFileUploadNotFound` | File Upload Not Found |
| `not-found/group-not-found` | `ErrGroupNotFound` | Group Not Found |
| `not-found/inbox-not-found` | `ErrInboxNotFound` | Inbox Not Found |
| `not-found/nested-not-found` | `ErrNestedNotFound` | Nested Not Found |
| `not-found/plan-not-found` | `ErrPlanNotFound` | Plan Not Found |
| `not-found/site-not-found` | `ErrSiteNotFound` | Site Not Found |
| `not-found/user-not-found` | `ErrUserNotFound` | User Not Found |
| `processing-failure` | `ErrProcessingFailure` | Processing Failure |
| `processing-failure/agent-push-update-blocked` | `ErrAgentPushUpdateBlocked` | Agent Push Update Blocked |
| `processing-failure/agent-unavailable` | `ErrAgentUnavailable` | Agent Unavailable |
| `processing-failure/ai-task-cannot-be-run-manually` | `ErrAiTaskCannotBeRunManually` | Ai Task Cannot Be Run Manually |
| `processing-failure/already-completed` | `ErrAlreadyCompleted` | Already Completed |
| `processing-failure/automation-cannot-be-run-manually` | `ErrAutomationCannotBeRunManually` | Automation Cannot Be Run Manually |
| `processing-failure/behavior-not-allowed-on-remote-server` | `ErrBehaviorNotAllowedOnRemoteServer` | Behavior Not Allowed On Remote Server |
| `processing-failure/buffered-upload-disabled-for-this-destination` | `ErrBufferedUploadDisabledForThisDestination` | Buffered Upload Disabled For This Destination |
| `processing-failure/bundle-only-allows-previews` | `ErrBundleOnlyAllowsPreviews` | Bundle Only Allows Previews |
| `processing-failure/bundle-operation-requires-subfolder` | `ErrBundleOperationRequiresSubfolder` | Bundle Operation Requires Subfolder |
| `processing-failure/configuration-locked-path` | `ErrConfigurationLockedPath` | Configuration Locked Path |
| `processing-failure/could-not-create-parent` | `ErrCouldNotCreateParent` | Could Not Create Parent |
| `processing-failure/destination-exists` | `ErrDestinationExists` | Destination Exists |
| `processing-failure/destination-folder-limited` | `ErrDestinationFolderLimited` | Destination Folder Limited |
| `processing-failure/destination-parent-conflict` | `ErrDestinationParentConflict` | Destination Parent Conflict |
| `processing-failure/destination-parent-does-not-exist` | `ErrDestinationParentDoesNotExist` | Destination Parent Does Not Exist |
| `processing-failure/exceeded-runtime-limit` | `ErrExceededRuntimeLimit` | Exceeded Runtime Limit |
| `processing-failure/expectation-already-has-open-window` | `ErrExpectationAlreadyHasOpenWindow` | Expectation Already Has Open Window |
| `processing-failure/expectation-not-manual-trigger` | `ErrExpectationNotManualTrigger` | Expectation Not Manual Trigger |
| `processing-failure/expired-private-key` | `ErrExpiredPrivateKey` | Expired Private Key |
| `processing-failure/expired-public-key` | `ErrExpiredPublicKey` | Expired Public Key |
| `processing-failure/export-failure` | `ErrExportFailure` | Export Failure |
| `processing-failure/export-not-ready` | `ErrExportNotReady` | Export Not Ready |
| `processing-failure/failed-to-change-password` | `ErrFailedToChangePassword` | Failed To Change Password |
| `processing-failure/file-locked` | `ErrFileLocked` | File Locked |
| `processing-failure/file-not-uploaded` | `ErrFileNotUploaded` | File Not Uploaded |
| `processing-failure/file-pending-processing` | `ErrFilePendingProcessing` | File Pending Processing |
| `processing-failure/file-processing-error` | `ErrFileProcessingError` | File Processing Error |
| `processing-failure/file-too-big-to-decrypt` | `ErrFileTooBigToDecrypt` | File Too Big To Decrypt |
| `processing-failure/file-too-big-to-encrypt` | `ErrFileTooBigToEncrypt` | File Too Big To Encrypt |
| `processing-failure/file-uploaded-to-wrong-region` | `ErrFileUploadedToWrongRegion` | File Uploaded To Wrong Region |
| `processing-failure/filename-too-long` | `ErrFilenameTooLong` | Filename Too Long |
| `processing-failure/folder-locked` | `ErrFolderLocked` | Folder Locked |
| `processing-failure/folder-not-empty` | `ErrFolderNotEmpty` | Folder Not Empty |
| `processing-failure/history-unavailable` | `ErrHistoryUnavailable` | History Unavailable |
| `processing-failure/invalid-bundle-code` | `ErrInvalidBundleCode` | Invalid Bundle Code |
| `processing-failure/invalid-file-type` | `ErrInvalidFileType` | Invalid File Type |
| `processing-failure/invalid-filename` | `ErrInvalidFilename` | Invalid Filename |
| `processing-failure/invalid-priority-color` | `ErrInvalidPriorityColor` | Invalid Priority Color |
| `processing-failure/invalid-range` | `ErrInvalidRange` | Invalid Range |
| `processing-failure/invalid-site` | `ErrInvalidSite` | Invalid Site |
| `processing-failure/invalid-zip-file` | `ErrInvalidZipFile` | Invalid Zip File |
| `processing-failure/metadata-not-supported-on-remotes` | `ErrMetadataNotSupportedOnRemotes` | Metadata Not Supported On Remotes |
| `processing-failure/model-save-error` | `ErrModelSaveError` | Model Save Error |
| `processing-failure/multiple-processing-errors` | `ErrMultipleProcessingErrors` | Multiple Processing Errors |
| `processing-failure/path-too-long` | `ErrPathTooLong` | Path Too Long |
| `processing-failure/recipient-already-shared` | `ErrRecipientAlreadyShared` | Recipient Already Shared |
| `processing-failure/remote-entry-read-only` | `ErrRemoteEntryReadOnly` | Remote Entry Read Only |
| `processing-failure/remote-server-error` | `ErrRemoteServerError` | Remote Server Error |
| `processing-failure/resource-belongs-to-parent-site` | `ErrResourceBelongsToParentSite` | Resource Belongs To Parent Site |
| `processing-failure/resource-locked` | `ErrResourceLocked` | Resource Locked |
| `processing-failure/subfolder-locked` | `ErrSubfolderLocked` | Subfolder Locked |
| `processing-failure/sync-in-progress` | `ErrSyncInProgress` | Sync In Progress |
| `processing-failure/two-factor-authentication-code-already-sent` | `ErrTwoFactorAuthenticationCodeAlreadySent` | Two Factor Authentication Code Already Sent |
| `processing-failure/two-factor-authentication-country-blacklisted` | `ErrTwoFactorAuthenticationCountryBlacklisted` | Two Factor Authentication Country Blacklisted |
| `processing-failure/two-factor-authentication-general-error` | `ErrTwoFactorAuthenticationGeneralError` | Two Factor Authentication General Error |
| `processing-failure/two-factor-authentication-method-unsupported-error` | `ErrTwoFactorAuthenticationMethodUnsupportedError` | Two Factor Authentication Method Unsupported Error |
| `processing-failure/two-factor-authentication-unsubscribed-recipient` | `ErrTwoFactorAuthenticationUnsubscribedRecipient` | Two Factor Authentication Unsubscribed Recipient |
| `processing-failure/updates-not-allowed-for-remotes` | `ErrUpdatesNotAllowedForRemotes` | Updates Not Allowed For Remotes |
| `rate-limited/duplicate-share-recipient` | `ErrDuplicateShareRecipient` | Duplicate Share Recipient |
| `rate-limited/reauthentication-rate-limited` | `ErrReauthenticationRateLimited` | Reauthentication Rate Limited |
| `rate-limited/too-many-concurrent-logins` | `ErrTooManyConcurrentLogins` | Too Many Concurrent Logins |
| `rate-limited/too-many-concurrent-requests` | `ErrTooManyConcurrentRequests` | Too Many Concurrent Requests |
| `rate-limited/too-many-login-attempts` | `ErrTooManyLoginAttempts` | Too Many Login Attempts |
| `rate-limited/too-many-requests` | `ErrTooManyRequests` | Too Many Requests |
| `rate-limited/too-many-shares` | `ErrTooManyShares` | Too Many Shares |
| `service-unavailable/automations-unavailable` | `ErrAutomationsUnavailable` | Automations Unavailable |
| `service-unavailable/lock-operation-busy` | `ErrLockOperationBusy` | Lock Operation Busy |
| `service-unavailable/migration-in-progress` | `ErrMigrationInProgress` | Migration In Progress |
| `service-unavailable/search-unavailable` | `ErrSearchUnavailable` | Search Unavailable |
| `service-unavailable/site-disabled` | `ErrSiteDisabled` | Site Disabled |
| `service-unavailable/uploads-unavailable` | `ErrUploadsUnavailable` | Uploads Unavailable |
| `site-configuration/account-already-exists` | `ErrAccountAlreadyExists` | Account Already Exists |
| `site-configuration/account-overdue` | `ErrAccountOverdue` | Account Overdue |
| `site-configuration/no-account-for-site` | `ErrNoAccountForSite` | No Account For Site |
| `site-configuration/site-was-removed` | `ErrSiteWasRemoved` | Site Was Removed |
| `site-configuration/trial-expired` | `ErrTrialExpired` | Trial Expired |
| `site-configuration/trial-locked` | `ErrTrialLocked` | Trial Locked |
| `site-configuration/user-requests-enabled-required` | `ErrUserRequestsEnabledRequired` | User Requests Enabled Required |

### ResponseError Helpers

Helpers are provided for matching error families.

| Error Type Prefix | Go Error | Helper |
| --------- | --------- | --------- |
| `bad-request` | `ErrBadRequest` | `IsBadRequest` |
| `not-authenticated` | `ErrNotAuthenticated` | `IsAuthenticationError` |
| `not-authorized` | `ErrNotAuthorized` | `IsAuthorizationError` |
| `not-found` | `ErrNotFound` | `IsNotFound` |
| `processing-failure` | `ErrProcessingFailure` | `IsProcessingFailure` |
| `rate-limited` | `ErrRateLimited` | `IsRateLimited` |
| `service-unavailable` | `ErrServiceUnavailable` | `IsServiceUnavailable` |
| `site-configuration` | `ErrSiteConfiguration` | `IsSiteConfiguration` |

## Pagination

Listing methods return lazy iterators. Use `for resource, err := range listing.All()`
to read the results. Calling the listing method or `All()` makes no API request.
The iterator fetches pages during iteration. `All()` does not collect all pages first.

Each resource has a `nil` error. Check the error before using the resource. An
unhandled request, decoding, or cancellation error is yielded once and ends that traversal.
An empty listing, normal completion, or a page limit ends the loop without an
extra resource or error.

Set `ListParams.PerPage` to request a page size. Set `ListParams.MaxPages` to limit
the number of pages requested. Its default value, `0`, allows all pages.

```go title="Range Over Results" hasDataFormatSelector
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/folder"
)

listing, err := folder.ListFor(files_sdk.FolderListForParams{Path: "path"})
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

for file, err := range listing.All() {
    if err != nil {
        var respErr files_sdk.ResponseError
        if errors.As(err, &respErr) {
            fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
        } else {
            fmt.Printf("Request failed: %v\n", err)
        }
        return
    }
    fmt.Println(file.Path)
}
```

### Stop or Resume Iteration

Use `break` or `return` to stop the loop. No further results are yielded or pages
requested by that traversal.

The listing is single-use. `All()` continues from its current position instead
of restarting it. After earlier `Next()` calls or an early `break`, another
`All()` loop or `Next()` call resumes at the next resource. After completion,
reaching `MaxPages`, or a failure, later `All()` calls yield nothing and make no
requests. `Err()` retains a failure, including one encountered through `Next()`.

Call the listing method again to restart. `Reload()` also creates a fresh
iterator, returned as `files_sdk.IterI`. Use a type assertion to the generated
iterator type if you want to call `All()` on that result.

Both `All()` and `Next()` advance the same iterator. Use one traversal at a time. Do not call `Next()` inside an `All()` loop or use
one iterator from several goroutines.

### Use Next and Err

Existing `Next()` and `Err()` code remains supported. Read the typed resource,
such as `File()`, only after `Next()` returns `true`. Check `Err()` after the loop.

```go title="Next and Err"
import (
    "fmt"
    "errors"

    files_sdk "github.com/Files-com/files-sdk-go/v3"
    "github.com/Files-com/files-sdk-go/v3/folder"
)

listing, err := folder.ListFor(files_sdk.FolderListForParams{Path: "path"})
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}

for listing.Next() {
    fmt.Println(listing.File().Path)
}
err = listing.Err()
if err != nil {
    var respErr files_sdk.ResponseError
    if errors.As(err, &respErr) {
        fmt.Printf("API error (%s): %s\n", respErr.Type, respErr.ErrorMessage)
    } else {
        fmt.Printf("Request failed: %v\n", err)
    }
    return
}
```

## Mock Server

Files.com publishes a mock Files.com API server, which is useful for testing your use of the Files.com
SDKs and other direct integrations against the Files.com API in an integration test environment.
It never checks credentials: send any placeholder API key, and never use real Files.com credentials
with it.

The server has two modes, chosen when it starts:

* **Legacy mode** (the default) checks required parameters and parameter types, then returns a fixed
  example response for each API endpoint. It does not maintain state and it does not deeply inspect
  your submissions for correctness, which makes it useful for testing basic network operations and
  JSON encoding for your SDK or API client.
* **Simulation mode** keeps records, files and folders in memory, so a test can create, list, update
  and delete resources, upload a file and download the same bytes, and make chosen requests fail,
  stall or lose their connection on purpose. Requests it does not simulate fail with a clear error
  instead of returning an example response.

Start the server from its source with Ruby and Bundler. `FILES_MOCK_MODE` is read once at startup:
leaving it unset or setting it to `legacy` starts legacy mode, `simulation` starts simulation mode,
and any other value stops startup with an error. Legacy mode listens on port 4041 on all IPv4
interfaces, and simulation mode on `127.0.0.1:4041`.

Simulation mode keeps its state only in the server process. Its control endpoints under
`/__files_mock/v1` report when the server is ready, reset it with your fixtures, add fault rules
and return the journal of the requests it received. It refuses work over its limits instead of
truncating it. The README in the source describes all of these, the operations it simulates and
how to configure its limits.

Download the server as a Docker image via [Docker Hub](https://hub.docker.com/r/filescom/files-mock-server).
The image's `latest` tag moves to whichever server was published last, so it need not include
simulation mode; to run exactly the server the README describes, build the image from the source.

The Source Code is also available on [GitHub](https://github.com/Files-com/files-mock-server).

```shell title="Start the Mock Server"
bundle install

## Legacy mode
bundle exec puma

## Simulation mode
FILES_MOCK_MODE=simulation bundle exec puma
```
