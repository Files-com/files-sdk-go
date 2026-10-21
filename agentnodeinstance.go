package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// AgentNodeInstance is a Files.com API resource.
type AgentNodeInstance struct {
	// Ephemeral ID for this running Agent process.
	InstanceId string `json:"instance_id,omitempty" path:"instance_id,omitempty" url:"instance_id,omitempty"`
	// Role of this process during an Agent update.
	ProcessState string `json:"process_state,omitempty" path:"process_state,omitempty" url:"process_state,omitempty"`
	// Whether this process has an available proxy connection.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Whether this process receives new unscoped work for its node.
	IsDefault *bool `json:"is_default,omitempty" path:"is_default,omitempty" url:"is_default,omitempty"`
	// Agent version reported by this process.
	AgentVersion string `json:"agent_version,omitempty" path:"agent_version,omitempty" url:"agent_version,omitempty"`
	// Most recent successful observation for this process.
	LastSeenAt *time.Time `json:"last_seen_at,omitempty" path:"last_seen_at,omitempty" url:"last_seen_at,omitempty"`
	// Proxy connections observed for this process.
	Connections []AgentNodeConnection `json:"connections,omitempty" path:"connections,omitempty" url:"connections,omitempty"`
}

// Identifier no path or id

// AgentNodeInstanceCollection is a list of AgentNodeInstance resources.
type AgentNodeInstanceCollection []AgentNodeInstance

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *AgentNodeInstance) UnmarshalJSON(data []byte) error {
	type agentNodeInstance AgentNodeInstance
	var v agentNodeInstance
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AgentNodeInstance(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AgentNodeInstanceCollection) UnmarshalJSON(data []byte) error {
	type agentNodeInstances AgentNodeInstanceCollection
	var v agentNodeInstances
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AgentNodeInstanceCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AgentNodeInstanceCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
