package files_sdk

import (
	"encoding/json"
	"io"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Behavior struct {
	Id               int64           `json:"id,omitempty"`
	Path             string          `json:"path,omitempty"`
	AttachmentUrl    string          `json:"attachment_url,omitempty"`
	Behavior         string          `json:"behavior,omitempty"`
	Name             string          `json:"name,omitempty"`
	Description      string          `json:"description,omitempty"`
	Value            json.RawMessage `json:"value,omitempty"`
	AttachmentFile   io.Reader       `json:"attachment_file,omitempty"`
	AttachmentDelete *bool           `json:"attachment_delete,omitempty"`
}

type BehaviorCollection []Behavior

type BehaviorListParams struct {
	Cursor     string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage    int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	Behavior   string          `url:"behavior,omitempty" required:"false" json:"behavior,omitempty"`
	lib.ListParams
}

type BehaviorFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
}

type BehaviorListForParams struct {
	Cursor     string          `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage    int64           `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty"`
	Path       string          `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Recursive  string          `url:"recursive,omitempty" required:"false" json:"recursive,omitempty"`
	Behavior   string          `url:"behavior,omitempty" required:"false" json:"behavior,omitempty"`
	lib.ListParams
}

type BehaviorCreateParams struct {
	Value          string    `url:"value,omitempty" required:"false" json:"value,omitempty"`
	AttachmentFile io.Writer `url:"attachment_file,omitempty" required:"false" json:"attachment_file,omitempty"`
	Name           string    `url:"name,omitempty" required:"false" json:"name,omitempty"`
	Description    string    `url:"description,omitempty" required:"false" json:"description,omitempty"`
	Path           string    `url:"path,omitempty" required:"true" json:"path,omitempty"`
	Behavior       string    `url:"behavior,omitempty" required:"true" json:"behavior,omitempty"`
}

type BehaviorWebhookTestParams struct {
	Url      string          `url:"url,omitempty" required:"true" json:"url,omitempty"`
	Method   string          `url:"method,omitempty" required:"false" json:"method,omitempty"`
	Encoding string          `url:"encoding,omitempty" required:"false" json:"encoding,omitempty"`
	Headers  json.RawMessage `url:"headers,omitempty" required:"false" json:"headers,omitempty"`
	Body     json.RawMessage `url:"body,omitempty" required:"false" json:"body,omitempty"`
	Action   string          `url:"action,omitempty" required:"false" json:"action,omitempty"`
}

type BehaviorUpdateParams struct {
	Id               int64     `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Value            string    `url:"value,omitempty" required:"false" json:"value,omitempty"`
	AttachmentFile   io.Writer `url:"attachment_file,omitempty" required:"false" json:"attachment_file,omitempty"`
	Name             string    `url:"name,omitempty" required:"false" json:"name,omitempty"`
	Description      string    `url:"description,omitempty" required:"false" json:"description,omitempty"`
	Behavior         string    `url:"behavior,omitempty" required:"false" json:"behavior,omitempty"`
	Path             string    `url:"path,omitempty" required:"false" json:"path,omitempty"`
	AttachmentDelete *bool     `url:"attachment_delete,omitempty" required:"false" json:"attachment_delete,omitempty"`
}

type BehaviorDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty"`
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

func (b *BehaviorCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
