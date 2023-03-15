package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type ExternalEvent struct {
	Id               int64      `json:"id,omitempty" path:"id"`
	EventType        string     `json:"event_type,omitempty" path:"event_type"`
	Status           string     `json:"status,omitempty" path:"status"`
	Body             string     `json:"body,omitempty" path:"body"`
	CreatedAt        *time.Time `json:"created_at,omitempty" path:"created_at"`
	BodyUrl          string     `json:"body_url,omitempty" path:"body_url"`
	FolderBehaviorId int64      `json:"folder_behavior_id,omitempty" path:"folder_behavior_id"`
	SuccessfulFiles  int64      `json:"successful_files,omitempty" path:"successful_files"`
	ErroredFiles     int64      `json:"errored_files,omitempty" path:"errored_files"`
	BytesSynced      int64      `json:"bytes_synced,omitempty" path:"bytes_synced"`
	RemoteServerType string     `json:"remote_server_type,omitempty" path:"remote_server_type"`
}

type ExternalEventCollection []ExternalEvent

type ExternalEventStatusEnum string

func (u ExternalEventStatusEnum) String() string {
	return string(u)
}

func (u ExternalEventStatusEnum) Enum() map[string]ExternalEventStatusEnum {
	return map[string]ExternalEventStatusEnum{
		"success":         ExternalEventStatusEnum("success"),
		"failure":         ExternalEventStatusEnum("failure"),
		"partial_failure": ExternalEventStatusEnum("partial_failure"),
		"in_progress":     ExternalEventStatusEnum("in_progress"),
		"skipped":         ExternalEventStatusEnum("skipped"),
	}
}

type ExternalEventListParams struct {
	SortBy       json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter       json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt     json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq   json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterPrefix json.RawMessage `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	FilterLt     json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq   json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	lib.ListParams
}

type ExternalEventFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

type ExternalEventCreateParams struct {
	Status ExternalEventStatusEnum `url:"status,omitempty" required:"true" json:"status,omitempty" path:"status"`
	Body   string                  `url:"body,omitempty" required:"true" json:"body,omitempty" path:"body"`
}

func (e *ExternalEvent) UnmarshalJSON(data []byte) error {
	type externalEvent ExternalEvent
	var v externalEvent
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*e = ExternalEvent(v)
	return nil
}

func (e *ExternalEventCollection) UnmarshalJSON(data []byte) error {
	type externalEvents ExternalEventCollection
	var v externalEvents
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*e = ExternalEventCollection(v)
	return nil
}

func (e *ExternalEventCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*e))
	for i, v := range *e {
		ret[i] = v
	}

	return &ret
}
