package files_sdk

import (
	"encoding/json"
	"io"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Behavior struct {
	Id               int64           `json:"id,omitempty" path:"id"`
	Path             string          `json:"path,omitempty" path:"path"`
	AttachmentUrl    string          `json:"attachment_url,omitempty" path:"attachment_url"`
	Behavior         string          `json:"behavior,omitempty" path:"behavior"`
	Name             string          `json:"name,omitempty" path:"name"`
	Description      string          `json:"description,omitempty" path:"description"`
	Value            json.RawMessage `json:"value,omitempty" path:"value"`
	AttachmentFile   io.Reader       `json:"attachment_file,omitempty" path:"attachment_file"`
	AttachmentDelete *bool           `json:"attachment_delete,omitempty" path:"attachment_delete"`
}

type BehaviorCollection []Behavior

type BehaviorListParams struct {
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty" path:"filter_like"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	Behavior   string          `url:"behavior,omitempty" required:"false" json:"behavior,omitempty" path:"behavior"`
	lib.ListParams
}

type BehaviorFindParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

type BehaviorListForParams struct {
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false" json:"filter_like,omitempty" path:"filter_like"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	Path       string          `url:"-,omitempty" required:"true" json:"-,omitempty" path:"path"`
	Recursive  string          `url:"recursive,omitempty" required:"false" json:"recursive,omitempty" path:"recursive"`
	Behavior   string          `url:"behavior,omitempty" required:"false" json:"behavior,omitempty" path:"behavior"`
	lib.ListParams
}

type BehaviorCreateParams struct {
	Value          string    `url:"value,omitempty" required:"false" json:"value,omitempty" path:"value"`
	AttachmentFile io.Writer `url:"attachment_file,omitempty" required:"false" json:"attachment_file,omitempty" path:"attachment_file"`
	Name           string    `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Description    string    `url:"description,omitempty" required:"false" json:"description,omitempty" path:"description"`
	Path           string    `url:"path,omitempty" required:"true" json:"path,omitempty" path:"path"`
	Behavior       string    `url:"behavior,omitempty" required:"true" json:"behavior,omitempty" path:"behavior"`
}

type BehaviorWebhookTestParams struct {
	Url      string          `url:"url,omitempty" required:"true" json:"url,omitempty" path:"url"`
	Method   string          `url:"method,omitempty" required:"false" json:"method,omitempty" path:"method"`
	Encoding string          `url:"encoding,omitempty" required:"false" json:"encoding,omitempty" path:"encoding"`
	Headers  json.RawMessage `url:"headers,omitempty" required:"false" json:"headers,omitempty" path:"headers"`
	Body     json.RawMessage `url:"body,omitempty" required:"false" json:"body,omitempty" path:"body"`
	Action   string          `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
}

type BehaviorUpdateParams struct {
	Id               int64     `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
	Value            string    `url:"value,omitempty" required:"false" json:"value,omitempty" path:"value"`
	AttachmentFile   io.Writer `url:"attachment_file,omitempty" required:"false" json:"attachment_file,omitempty" path:"attachment_file"`
	Name             string    `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Description      string    `url:"description,omitempty" required:"false" json:"description,omitempty" path:"description"`
	Behavior         string    `url:"behavior,omitempty" required:"false" json:"behavior,omitempty" path:"behavior"`
	Path             string    `url:"path,omitempty" required:"false" json:"path,omitempty" path:"path"`
	AttachmentDelete *bool     `url:"attachment_delete,omitempty" required:"false" json:"attachment_delete,omitempty" path:"attachment_delete"`
}

type BehaviorDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true" json:"-,omitempty" path:"id"`
}

func (b *Behavior) UnmarshalJSON(data []byte) error {
	type behavior Behavior
	var v behavior
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = Behavior(v)
	return nil
}

func (b *BehaviorCollection) UnmarshalJSON(data []byte) error {
	type behaviors BehaviorCollection
	var v behaviors
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
