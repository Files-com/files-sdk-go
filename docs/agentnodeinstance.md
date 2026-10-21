# AgentNodeInstance

`files_sdk.AgentNodeInstance` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/agentnodeinstance`, named `agent_node_instance`.

This package has no API operations of its own.

## Model

`files_sdk.AgentNodeInstance` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `InstanceId` | `string` | `instance_id` | Ephemeral ID for this running Agent process |
| `ProcessState` | `string` | `process_state` | Role of this process during an Agent update |
| `Status` | `string` | `status` | Whether this process has an available proxy connection |
| `IsDefault` | `*bool` | `is_default` | Whether this process receives new unscoped work for its node |
| `AgentVersion` | `string` | `agent_version` | Agent version reported by this process |
| `LastSeenAt` | `*time.Time` | `last_seen_at` | Most recent successful observation for this process |
| `Connections` | `[]files_sdk.AgentNodeConnection` | `connections` | Proxy connections observed for this process |
