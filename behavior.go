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
	Cursor     string          `url:"cursor,omitempty" required:"false"`
	PerPage    int64           `url:"per_page,omitempty" required:"false"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	Behavior   string          `url:"behavior,omitempty" required:"false"`
	lib.ListParams
}

type BehaviorFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type BehaviorListForParams struct {
	Cursor     string          `url:"cursor,omitempty" required:"false"`
	PerPage    int64           `url:"per_page,omitempty" required:"false"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	Path       string          `url:"-,omitempty" required:"true"`
	Recursive  string          `url:"recursive,omitempty" required:"false"`
	Behavior   string          `url:"behavior,omitempty" required:"false"`
	lib.ListParams
}

type BehaviorCreateParams struct {
	Value          string    `url:"value,omitempty" required:"false"`
	AttachmentFile io.Writer `url:"attachment_file,omitempty" required:"false"`
	Name           string    `url:"name,omitempty" required:"false"`
	Description    string    `url:"description,omitempty" required:"false"`
	Path           string    `url:"path,omitempty" required:"true"`
	Behavior       string    `url:"behavior,omitempty" required:"true"`
}

type BehaviorWebhookTestParams struct {
	Url      string          `url:"url,omitempty" required:"true"`
	Method   string          `url:"method,omitempty" required:"false"`
	Encoding string          `url:"encoding,omitempty" required:"false"`
	Headers  json.RawMessage `url:"headers,omitempty" required:"false"`
	Body     json.RawMessage `url:"body,omitempty" required:"false"`
	Action   string          `url:"action,omitempty" required:"false"`
}

type BehaviorUpdateParams struct {
	Id               int64     `url:"-,omitempty" required:"true"`
	Value            string    `url:"value,omitempty" required:"false"`
	AttachmentFile   io.Writer `url:"attachment_file,omitempty" required:"false"`
	Name             string    `url:"name,omitempty" required:"false"`
	Description      string    `url:"description,omitempty" required:"false"`
	Behavior         string    `url:"behavior,omitempty" required:"false"`
	Path             string    `url:"path,omitempty" required:"false"`
	AttachmentDelete *bool     `url:"attachment_delete,omitempty" required:"false"`
}

type BehaviorDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
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
