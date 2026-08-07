package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type AgentNodeConnection struct {
	Mode       string     `json:"mode,omitempty" path:"mode,omitempty" url:"mode,omitempty"`
	Status     string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	LastSeenAt *time.Time `json:"last_seen_at,omitempty" path:"last_seen_at,omitempty" url:"last_seen_at,omitempty"`
}

// Identifier no path or id

type AgentNodeConnectionCollection []AgentNodeConnection

func (a *AgentNodeConnection) UnmarshalJSON(data []byte) error {
	type agentNodeConnection AgentNodeConnection
	var v agentNodeConnection
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AgentNodeConnection(v)
	return nil
}

func (a *AgentNodeConnectionCollection) UnmarshalJSON(data []byte) error {
	type agentNodeConnections AgentNodeConnectionCollection
	var v agentNodeConnections
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AgentNodeConnectionCollection(v)
	return nil
}

func (a *AgentNodeConnectionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
