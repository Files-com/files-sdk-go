package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// AgentNodeConnection is a Files.com API resource.
type AgentNodeConnection struct {
	Mode       string     `json:"mode,omitempty" path:"mode,omitempty" url:"mode,omitempty"`
	Status     string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	LastSeenAt *time.Time `json:"last_seen_at,omitempty" path:"last_seen_at,omitempty" url:"last_seen_at,omitempty"`
}

// Identifier no path or id

// AgentNodeConnectionCollection is a list of AgentNodeConnection resources.
type AgentNodeConnectionCollection []AgentNodeConnection

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *AgentNodeConnection) UnmarshalJSON(data []byte) error {
	type agentNodeConnection AgentNodeConnection
	var v agentNodeConnection
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AgentNodeConnection(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AgentNodeConnectionCollection) UnmarshalJSON(data []byte) error {
	type agentNodeConnections AgentNodeConnectionCollection
	var v agentNodeConnections
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AgentNodeConnectionCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AgentNodeConnectionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
