# AutomationExecutionNode

`files_sdk.AutomationExecutionNode` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/automationexecutionnode`, named `automation_execution_node`.

This package has no API operations of its own.

## Model

`files_sdk.AutomationExecutionNode` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `NodeId` | `string` | `node_id` | Node ID from the pinned Automation definition. |
| `NodeType` | `string` | `node_type` | Node type. |
| `Status` | `string` | `status` | Node status. |
| `RunStage` | `string` | `run_stage` | Current node execution stage. |
| `Reused` | `*bool` | `reused` | Whether this node reused persisted output from an earlier run. |
| `SuccessfulOperations` | `int64` | `successful_operations` | Count of successful operations in this node. |
| `FailedOperations` | `int64` | `failed_operations` | Count of failed operations in this node. |
| `StartedAt` | `*time.Time` | `started_at` | When this node started. |
| `CompletedAt` | `*time.Time` | `completed_at` | When this node completed. |
| `DurationMs` | `int64` | `duration_ms` | Node runtime in milliseconds. |
| `Inputs` | `[]map[string]interface{}` | `inputs` | Ordered inbound edge references. |
| `Outputs` | `interface{}` | `outputs` | Output counts, item kinds, and typed-error summaries by outlet. |
| `InputItems` | `interface{}` | `input_items` | Available input items grouped by inlet. Resolve large values with item_refs and item_values. |
| `OutputItems` | `interface{}` | `output_items` | Output items grouped by outlet. Resolve large values with item_refs and item_values. |
| `ItemRefs` | `[]map[string]interface{}` | `item_refs` | References with a path from the response root and a digest. Replace the null at each path with the value for that digest. Path segments are literal object keys or zero-based array indexes. |
| `ItemValues` | `[]map[string]interface{}` | `item_values` | Ordered entries with a digest, a value, and optional item_refs relative to that value. References inside entries point to earlier entries in this array. |
