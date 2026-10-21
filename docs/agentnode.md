# AgentNode

`files_sdk.AgentNode` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/agentnode`, named `agent_node`.

This package has no API operations of its own.

## Model

`files_sdk.AgentNode` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `NodeId` | `string` | `node_id` | Stable Agent installation ID |
| `Name` | `string` | `name` | Customer-configured Agent node name |
| `Hostname` | `string` | `hostname` | Hostname reported by the Agent |
| `AvailabilityRole` | `string` | `availability_role` | Configured traffic preference |
| `Status` | `string` | `status` | Whether this node currently has an available Agent instance |
| `IsDefault` | `*bool` | `is_default` | Whether this node is the current default route for new unscoped work |
| `DirectTransferAvailable` | `*bool` | `direct_transfer_available` | Whether the proxy recently validated a direct connection to this Agent node. False means direct transfers are enabled but not currently available; null means disabled or unsupported. |
| `LastSeenAt` | `*time.Time` | `last_seen_at` | Most recent successful node observation |
| `Instances` | `[]files_sdk.AgentNodeInstance` | `instances` | Current Agent processes for this node |
