package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type ExternalEvent struct {
	Id               int64     `json:"id,omitempty"`
	EventType        string    `json:"event_type,omitempty"`
	Status           string    `json:"status,omitempty"`
	Body             string    `json:"body,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	BodyUrl          string    `json:"body_url,omitempty"`
	FolderBehaviorId int64     `json:"folder_behavior_id,omitempty"`
	SuccessfulFiles  int64     `json:"successful_files,omitempty"`
	ErroredFiles     int64     `json:"errored_files,omitempty"`
	BytesSynced      int64     `json:"bytes_synced,omitempty"`
	RemoteServerType string    `json:"remote_server_type,omitempty"`
}

type ExternalEventCollection []ExternalEvent

type ExternalEventStatusEnum string

func (u ExternalEventStatusEnum) String() string {
	return string(u)
}

func (u ExternalEventStatusEnum) Enum() map[string]ExternalEventStatusEnum {
	return map[string]ExternalEventStatusEnum{
		"success":         ExternalEventStatusEnum("success"),
		"error":           ExternalEventStatusEnum("error"),
		"partial_failure": ExternalEventStatusEnum("partial_failure"),
	}
}

type ExternalEventListParams struct {
	Cursor     string          `url:"cursor,omitempty" required:"false"`
	PerPage    int64           `url:"per_page,omitempty" required:"false"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	lib.ListParams
}

type ExternalEventFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type ExternalEventCreateParams struct {
	Status ExternalEventStatusEnum `url:"status,omitempty" required:"true"`
	Body   string                  `url:"body,omitempty" required:"true"`
}

func (e *ExternalEvent) UnmarshalJSON(data []byte) error {
	type externalEvent ExternalEvent
	var v externalEvent
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*e = ExternalEvent(v)
	return nil
}

func (e *ExternalEventCollection) UnmarshalJSON(data []byte) error {
	type externalEvents []ExternalEvent
	var v externalEvents
	if err := json.Unmarshal(data, &v); err != nil {
		return err
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
