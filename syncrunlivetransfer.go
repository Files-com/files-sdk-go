package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// SyncRunLiveTransfer is a Files.com API resource.
type SyncRunLiveTransfer struct {
	// The file path being transferred. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Status of this individual transfer.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Bytes transferred so far.
	BytesCopied int64 `json:"bytes_copied,omitempty" path:"bytes_copied,omitempty" url:"bytes_copied,omitempty"`
	// Total bytes of the file being transferred.
	BytesTotal int64 `json:"bytes_total,omitempty" path:"bytes_total,omitempty" url:"bytes_total,omitempty"`
	// Transfer progress from 0.0 to 1.0.
	Percentage float64 `json:"percentage,omitempty" path:"percentage,omitempty" url:"percentage,omitempty"`
	// Estimated time remaining (human-readable).
	Eta string `json:"eta,omitempty" path:"eta,omitempty" url:"eta,omitempty"`
	// When this individual transfer started.
	StartedAt string `json:"started_at,omitempty" path:"started_at,omitempty" url:"started_at,omitempty"`
}

// Identifier returns the resource path.
func (s SyncRunLiveTransfer) Identifier() interface{} {
	return s.Path
}

// SyncRunLiveTransferCollection is a list of SyncRunLiveTransfer resources.
type SyncRunLiveTransferCollection []SyncRunLiveTransfer

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *SyncRunLiveTransfer) UnmarshalJSON(data []byte) error {
	type syncRunLiveTransfer SyncRunLiveTransfer
	var v syncRunLiveTransfer
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SyncRunLiveTransfer(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SyncRunLiveTransferCollection) UnmarshalJSON(data []byte) error {
	type syncRunLiveTransfers SyncRunLiveTransferCollection
	var v syncRunLiveTransfers
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SyncRunLiveTransferCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SyncRunLiveTransferCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
