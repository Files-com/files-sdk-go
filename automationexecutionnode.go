package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// AutomationExecutionNode is a Files.com API resource.
type AutomationExecutionNode struct {
	// Node ID from the pinned Automation definition.
	NodeId string `json:"node_id,omitempty" path:"node_id,omitempty" url:"node_id,omitempty"`
	// Node type.
	NodeType string `json:"node_type,omitempty" path:"node_type,omitempty" url:"node_type,omitempty"`
	// Node status.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Current node execution stage.
	RunStage string `json:"run_stage,omitempty" path:"run_stage,omitempty" url:"run_stage,omitempty"`
	// Whether this node reused persisted output from an earlier run.
	Reused *bool `json:"reused,omitempty" path:"reused,omitempty" url:"reused,omitempty"`
	// Count of successful operations in this node.
	SuccessfulOperations int64 `json:"successful_operations,omitempty" path:"successful_operations,omitempty" url:"successful_operations,omitempty"`
	// Count of failed operations in this node.
	FailedOperations int64 `json:"failed_operations,omitempty" path:"failed_operations,omitempty" url:"failed_operations,omitempty"`
	// When this node started.
	StartedAt *time.Time `json:"started_at,omitempty" path:"started_at,omitempty" url:"started_at,omitempty"`
	// When this node completed.
	CompletedAt *time.Time `json:"completed_at,omitempty" path:"completed_at,omitempty" url:"completed_at,omitempty"`
	// Node runtime in milliseconds.
	DurationMs int64 `json:"duration_ms,omitempty" path:"duration_ms,omitempty" url:"duration_ms,omitempty"`
	// Ordered inbound edge references.
	Inputs []map[string]interface{} `json:"inputs,omitempty" path:"inputs,omitempty" url:"inputs,omitempty"`
	// Output counts, item kinds, and typed-error summaries by outlet.
	Outputs interface{} `json:"outputs,omitempty" path:"outputs,omitempty" url:"outputs,omitempty"`
	// Materialized input items currently available, grouped by inlet.
	InputItems interface{} `json:"input_items,omitempty" path:"input_items,omitempty" url:"input_items,omitempty"`
	// Materialized output items grouped by outlet.
	OutputItems interface{} `json:"output_items,omitempty" path:"output_items,omitempty" url:"output_items,omitempty"`
}

// Identifier no path or id

// AutomationExecutionNodeCollection is a list of AutomationExecutionNode resources.
type AutomationExecutionNodeCollection []AutomationExecutionNode

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *AutomationExecutionNode) UnmarshalJSON(data []byte) error {
	type automationExecutionNode AutomationExecutionNode
	var v automationExecutionNode
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AutomationExecutionNode(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AutomationExecutionNodeCollection) UnmarshalJSON(data []byte) error {
	type automationExecutionNodes AutomationExecutionNodeCollection
	var v automationExecutionNodes
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AutomationExecutionNodeCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AutomationExecutionNodeCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
