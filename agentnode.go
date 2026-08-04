package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type AgentNode struct {
	NodeId                  string     `json:"node_id,omitempty" path:"node_id,omitempty" url:"node_id,omitempty"`
	Name                    string     `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Hostname                string     `json:"hostname,omitempty" path:"hostname,omitempty" url:"hostname,omitempty"`
	AvailabilityRole        string     `json:"availability_role,omitempty" path:"availability_role,omitempty" url:"availability_role,omitempty"`
	ConnectionStatus        string     `json:"connection_status,omitempty" path:"connection_status,omitempty" url:"connection_status,omitempty"`
	IsDefault               *bool      `json:"is_default,omitempty" path:"is_default,omitempty" url:"is_default,omitempty"`
	AgentVersion            string     `json:"agent_version,omitempty" path:"agent_version,omitempty" url:"agent_version,omitempty"`
	DirectTransferAvailable *bool      `json:"direct_transfer_available,omitempty" path:"direct_transfer_available,omitempty" url:"direct_transfer_available,omitempty"`
	LastSeenAt              *time.Time `json:"last_seen_at,omitempty" path:"last_seen_at,omitempty" url:"last_seen_at,omitempty"`
}

// Identifier no path or id

type AgentNodeCollection []AgentNode

func (a *AgentNode) UnmarshalJSON(data []byte) error {
	type agentNode AgentNode
	var v agentNode
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = AgentNode(v)
	return nil
}

func (a *AgentNodeCollection) UnmarshalJSON(data []byte) error {
	type agentNodes AgentNodeCollection
	var v agentNodes
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AgentNodeCollection(v)
	return nil
}

func (a *AgentNodeCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
