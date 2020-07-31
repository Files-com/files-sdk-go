package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
	"io"
)

type Behavior struct {
	Id             int             `json:"id,omitempty"`
	Path           string          `json:"path,omitempty"`
	AttachmentUrl  string          `json:"attachment_url,omitempty"`
	Behavior       string          `json:"behavior,omitempty"`
	Value          json.RawMessage `json:"value,omitempty"`
	AttachmentFile io.Reader       `json:"attachment_file,omitempty"`
}

type BehaviorCollection []Behavior

type BehaviorListParams struct {
	Page       int             `url:"page,omitempty"`
	PerPage    int             `url:"per_page,omitempty"`
	Action     string          `url:"action,omitempty"`
	Cursor     string          `url:"cursor,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty"`
	Behavior   string          `url:"behavior,omitempty"`
	lib.ListParams
}

type BehaviorFindParams struct {
	Id int `url:"-,omitempty"`
}

type BehaviorListForParams struct {
	Page       int             `url:"page,omitempty"`
	PerPage    int             `url:"per_page,omitempty"`
	Action     string          `url:"action,omitempty"`
	Cursor     string          `url:"cursor,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty"`
	Path       string          `url:"-,omitempty"`
	Recursive  string          `url:"recursive,omitempty"`
	Behavior   string          `url:"behavior,omitempty"`
	lib.ListParams
}

type BehaviorCreateParams struct {
	Value          string    `url:"value,omitempty"`
	AttachmentFile io.Writer `url:"attachment_file,omitempty"`
	Path           string    `url:"path,omitempty"`
	Behavior       string    `url:"behavior,omitempty"`
}

type BehaviorWebhookTestParams struct {
	Url      string          `url:"url,omitempty"`
	Method   string          `url:"method,omitempty"`
	Encoding string          `url:"encoding,omitempty"`
	Headers  json.RawMessage `url:"headers,omitempty"`
	Body     json.RawMessage `url:"body,omitempty"`
	Action   string          `url:"action,omitempty"`
}

type BehaviorUpdateParams struct {
	Id             int       `url:"-,omitempty"`
	Value          string    `url:"value,omitempty"`
	AttachmentFile io.Writer `url:"attachment_file,omitempty"`
	Behavior       string    `url:"behavior,omitempty"`
	Path           string    `url:"path,omitempty"`
}

type BehaviorDeleteParams struct {
	Id int `url:"-,omitempty"`
}

func (b *Behavior) UnmarshalJSON(data []byte) error {
	type behavior Behavior
	var v behavior
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = Behavior(v)
	return nil
}

func (b *BehaviorCollection) UnmarshalJSON(data []byte) error {
	type behaviors []Behavior
	var v behaviors
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*b = BehaviorCollection(v)
	return nil
}
