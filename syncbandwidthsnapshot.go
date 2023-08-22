package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type SyncBandwidthSnapshot struct {
	Id                int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	SiteId            int64      `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	SyncBytesReceived string     `json:"sync_bytes_received,omitempty" path:"sync_bytes_received,omitempty" url:"sync_bytes_received,omitempty"`
	SyncBytesSent     string     `json:"sync_bytes_sent,omitempty" path:"sync_bytes_sent,omitempty" url:"sync_bytes_sent,omitempty"`
	CreatedAt         *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	RemoteServerId    int64      `json:"remote_server_id,omitempty" path:"remote_server_id,omitempty" url:"remote_server_id,omitempty"`
}

func (s SyncBandwidthSnapshot) Identifier() interface{} {
	return s.Id
}

type SyncBandwidthSnapshotCollection []SyncBandwidthSnapshot

type SyncBandwidthSnapshotCreateParams struct {
	RemoteServerId    int64 `url:"remote_server_id,omitempty" required:"true" json:"remote_server_id,omitempty" path:"remote_server_id"`
	SyncBytesSent     int64 `url:"sync_bytes_sent,omitempty" required:"true" json:"sync_bytes_sent,omitempty" path:"sync_bytes_sent"`
	SyncBytesReceived int64 `url:"sync_bytes_received,omitempty" required:"true" json:"sync_bytes_received,omitempty" path:"sync_bytes_received"`
}

func (s *SyncBandwidthSnapshot) UnmarshalJSON(data []byte) error {
	type syncBandwidthSnapshot SyncBandwidthSnapshot
	var v syncBandwidthSnapshot
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SyncBandwidthSnapshot(v)
	return nil
}

func (s *SyncBandwidthSnapshotCollection) UnmarshalJSON(data []byte) error {
	type syncBandwidthSnapshots SyncBandwidthSnapshotCollection
	var v syncBandwidthSnapshots
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SyncBandwidthSnapshotCollection(v)
	return nil
}

func (s *SyncBandwidthSnapshotCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
