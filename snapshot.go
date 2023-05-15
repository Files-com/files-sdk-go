package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Snapshot struct {
	ExpiresAt   *time.Time `json:"expires_at,omitempty" path:"expires_at"`
	FinalizedAt *time.Time `json:"finalized_at,omitempty" path:"finalized_at"`
	Name        string     `json:"name,omitempty" path:"name"`
	UserId      int64      `json:"user_id,omitempty" path:"user_id"`
	BundleId    int64      `json:"bundle_id,omitempty" path:"bundle_id"`
	Paths       []string   `json:"paths,omitempty" path:"paths"`
	Id          int64      `json:"id,omitempty" path:"id"`
}

func (s Snapshot) Identifier() interface{} {
	return s.Id
}

type SnapshotCollection []Snapshot

type SnapshotListParams struct {
	ListParams
}

type SnapshotFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type SnapshotCreateParams struct {
	ExpiresAt *time.Time `url:"expires_at,omitempty" required:"false" json:"expires_at,omitempty" path:"expires_at"`
	Name      string     `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Paths     []string   `url:"paths,omitempty" required:"false" json:"paths,omitempty" path:"paths"`
}

type SnapshotUpdateParams struct {
	Id        int64      `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	ExpiresAt *time.Time `url:"expires_at,omitempty" required:"false" json:"expires_at,omitempty" path:"expires_at"`
	Name      string     `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Paths     []string   `url:"paths,omitempty" required:"false" json:"paths,omitempty" path:"paths"`
}

type SnapshotDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (s *Snapshot) UnmarshalJSON(data []byte) error {
	type snapshot Snapshot
	var v snapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Snapshot(v)
	return nil
}

func (s *SnapshotCollection) UnmarshalJSON(data []byte) error {
	type snapshots SnapshotCollection
	var v snapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SnapshotCollection(v)
	return nil
}

func (s *SnapshotCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
