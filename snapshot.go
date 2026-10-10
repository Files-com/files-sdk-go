package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Snapshot is a Files.com API resource.
type Snapshot struct {
	Id          int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" path:"expires_at,omitempty" url:"expires_at,omitempty"`
	FinalizedAt *time.Time `json:"finalized_at,omitempty" path:"finalized_at,omitempty" url:"finalized_at,omitempty"`
	Name        string     `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	UserId      int64      `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	BundleId    int64      `json:"bundle_id,omitempty" path:"bundle_id,omitempty" url:"bundle_id,omitempty"`
	WorkspaceId int64      `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	Paths       []string   `json:"paths,omitempty" path:"paths,omitempty" url:"paths,omitempty"`
}

// Identifier returns the resource ID.
func (s Snapshot) Identifier() interface{} {
	return s.Id
}

// SnapshotCollection is a list of Snapshot resources.
type SnapshotCollection []Snapshot

// SnapshotListParams contains the request parameters for GET /snapshots.
type SnapshotListParams struct {
	ListParams
}

// SnapshotFindParams contains the request parameters for GET /snapshots/{id}.
type SnapshotFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// SnapshotCreateParams contains the request parameters for POST /snapshots.
type SnapshotCreateParams struct {
	ExpiresAt   *time.Time `url:"expires_at,omitempty" json:"expires_at,omitempty" path:"expires_at"`
	Name        string     `url:"name,omitempty" json:"name,omitempty" path:"name"`
	Paths       []string   `url:"paths,omitempty" json:"paths,omitempty" path:"paths"`
	WorkspaceId int64      `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// SnapshotFinalizeParams contains the request parameters for POST /snapshots/{id}/finalize.
//
// Finalize Snapshot
type SnapshotFinalizeParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// SnapshotUpdateParams contains the request parameters for PATCH /snapshots/{id}.
type SnapshotUpdateParams struct {
	Id        int64      `url:"-,omitempty" json:"-,omitempty" path:"id"`
	ExpiresAt *time.Time `url:"expires_at,omitempty" json:"expires_at,omitempty" path:"expires_at"`
	Name      string     `url:"name,omitempty" json:"name,omitempty" path:"name"`
	Paths     []string   `url:"paths,omitempty" json:"paths,omitempty" path:"paths"`
}

// SnapshotDeleteParams contains the request parameters for DELETE /snapshots/{id}.
type SnapshotDeleteParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *Snapshot) UnmarshalJSON(data []byte) error {
	type snapshot Snapshot
	var v snapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Snapshot(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SnapshotCollection) UnmarshalJSON(data []byte) error {
	type snapshots SnapshotCollection
	var v snapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SnapshotCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SnapshotCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
