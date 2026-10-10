package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// AgentPushUpdate is a Files.com API resource.
type AgentPushUpdate struct {
	Version        string `json:"version,omitempty" path:"version,omitempty" url:"version,omitempty"`
	Message        string `json:"message,omitempty" path:"message,omitempty" url:"message,omitempty"`
	CurrentVersion string `json:"current_version,omitempty" path:"current_version,omitempty" url:"current_version,omitempty"`
	PendingVersion string `json:"pending_version,omitempty" path:"pending_version,omitempty" url:"pending_version,omitempty"`
	LastError      string `json:"last_error,omitempty" path:"last_error,omitempty" url:"last_error,omitempty"`
	Error          string `json:"error,omitempty" path:"error,omitempty" url:"error,omitempty"`
}

// Identifier no path or id

// AgentPushUpdateCollection is a list of AgentPushUpdate resources.
type AgentPushUpdateCollection []AgentPushUpdate

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *AgentPushUpdate) UnmarshalJSON(data []byte) error {
	type agentPushUpdate AgentPushUpdate
	var v agentPushUpdate
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AgentPushUpdate(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AgentPushUpdateCollection) UnmarshalJSON(data []byte) error {
	type agentPushUpdates AgentPushUpdateCollection
	var v agentPushUpdates
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AgentPushUpdateCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AgentPushUpdateCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
