# Action

`files_sdk.Action` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/action`, named `action`.

This package has no API operations of its own.

## Model

`files_sdk.Action` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Action ID |
| `Path` | `string` | `path` | Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters. |
| `When` | `*time.Time` | `when` | Action occurrence date/time |
| `Destination` | `string` | `destination` | The destination path for this action, if applicable |
| `Display` | `string` | `display` | Friendly displayed output |
| `Ip` | `string` | `ip` | IP Address that performed this action |
| `Source` | `string` | `source` | The source path for this action, if applicable |
| `Targets` | `interface{}` | `targets` | Targets |
| `UserId` | `int64` | `user_id` | User ID |
| `Username` | `string` | `username` | Username |
| `UserIsFromParentSite` | `*bool` | `user_is_from_parent_site` | true if this change was performed by a user on a parent site. |
| `Action` | `string` | `action` | Type of action |
| `FailureType` | `string` | `failure_type` | Failure type. If action was a user login or session failure, why did it fail? |
| `Interface` | `string` | `interface` | Interface on which this action occurred. |
