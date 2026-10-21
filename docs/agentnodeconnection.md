# AgentNodeConnection

`files_sdk.AgentNodeConnection` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/agentnodeconnection`, named `agent_node_connection`.

This package has no API operations of its own.

## Model

`files_sdk.AgentNodeConnection` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Mode` | `string` | `mode` | How the Agent process uses this proxy connection |
| `Status` | `string` | `status` | Whether this connection was observed recently and has not disconnected |
| `LastSeenAt` | `*time.Time` | `last_seen_at` | Most recent successful observation for this connection |
