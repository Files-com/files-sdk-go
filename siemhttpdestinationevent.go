package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// SiemHttpDestinationEvent is a Files.com API resource.
type SiemHttpDestinationEvent struct {
	Id                    int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	EventType             string     `json:"event_type,omitempty" path:"event_type,omitempty" url:"event_type,omitempty"`
	Status                string     `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	Body                  string     `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	EventErrors           []string   `json:"event_errors,omitempty" path:"event_errors,omitempty" url:"event_errors,omitempty"`
	CreatedAt             *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	BodyUrl               string     `json:"body_url,omitempty" path:"body_url,omitempty" url:"body_url,omitempty"`
	SiemHttpDestinationId int64      `json:"siem_http_destination_id,omitempty" path:"siem_http_destination_id,omitempty" url:"siem_http_destination_id,omitempty"`
}

// Identifier returns the resource ID.
func (s SiemHttpDestinationEvent) Identifier() interface{} {
	return s.Id
}

// SiemHttpDestinationEventCollection is a list of SiemHttpDestinationEvent resources.
type SiemHttpDestinationEventCollection []SiemHttpDestinationEvent

// SiemHttpDestinationEventListParams contains the request parameters for GET /siem_http_destination_events.
type SiemHttpDestinationEventListParams struct {
	SortBy     interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter     interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	FilterGt   interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt   interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	ListParams
}

// SiemHttpDestinationEventFindParams contains the request parameters for GET /siem_http_destination_events/{id}.
type SiemHttpDestinationEventFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *SiemHttpDestinationEvent) UnmarshalJSON(data []byte) error {
	type siemHttpDestinationEvent SiemHttpDestinationEvent
	var v siemHttpDestinationEvent
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SiemHttpDestinationEvent(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SiemHttpDestinationEventCollection) UnmarshalJSON(data []byte) error {
	type siemHttpDestinationEvents SiemHttpDestinationEventCollection
	var v siemHttpDestinationEvents
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SiemHttpDestinationEventCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SiemHttpDestinationEventCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
