# AgentPushUpdate

`files_sdk.AgentPushUpdate` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/agentpushupdate`, named `agent_push_update`.

This package has no API operations of its own.

## Model

`files_sdk.AgentPushUpdate` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Version` | `string` | `version` | Pushed agent version |
| `Message` | `string` | `message` | Update accepted or reason |
| `CurrentVersion` | `string` | `current_version` | Installed agent version |
| `PendingVersion` | `string` | `pending_version` | Pending agent version or null |
| `LastError` | `string` | `last_error` | Last error message or null |
| `Error` | `string` | `error` | Error code |
