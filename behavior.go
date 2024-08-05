package files_sdk

import (
	"encoding/json"
	"io"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type Behavior struct {
	Id                          int64                  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Path                        string                 `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	AttachmentUrl               string                 `json:"attachment_url,omitempty" path:"attachment_url,omitempty" url:"attachment_url,omitempty"`
	Behavior                    string                 `json:"behavior,omitempty" path:"behavior,omitempty" url:"behavior,omitempty"`
	Name                        string                 `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Description                 string                 `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	Value                       map[string]interface{} `json:"value,omitempty" path:"value,omitempty" url:"value,omitempty"`
	DisableParentFolderBehavior *bool                  `json:"disable_parent_folder_behavior,omitempty" path:"disable_parent_folder_behavior,omitempty" url:"disable_parent_folder_behavior,omitempty"`
	Recursive                   *bool                  `json:"recursive,omitempty" path:"recursive,omitempty" url:"recursive,omitempty"`
	AttachmentFile              io.Reader              `json:"attachment_file,omitempty" path:"attachment_file,omitempty" url:"attachment_file,omitempty"`
	AttachmentDelete            *bool                  `json:"attachment_delete,omitempty" path:"attachment_delete,omitempty" url:"attachment_delete,omitempty"`
}

func (b Behavior) Identifier() interface{} {
	return b.Id
}

type BehaviorCollection []Behavior

type BehaviorListParams struct {
	SortBy       map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter       Behavior               `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix map[string]interface{} `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

type BehaviorFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type BehaviorListForParams struct {
	SortBy            map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter            Behavior               `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix      map[string]interface{} `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	Path              string                 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	AncestorBehaviors *bool                  `url:"ancestor_behaviors,omitempty" required:"false" json:"ancestor_behaviors,omitempty" path:"ancestor_behaviors"`
	ListParams
}

type BehaviorCreateParams struct {
	Value                       string    `url:"value,omitempty" required:"false" json:"value,omitempty" path:"value"`
	AttachmentFile              io.Writer `url:"attachment_file,omitempty" required:"false" json:"attachment_file,omitempty" path:"attachment_file"`
	DisableParentFolderBehavior *bool     `url:"disable_parent_folder_behavior,omitempty" required:"false" json:"disable_parent_folder_behavior,omitempty" path:"disable_parent_folder_behavior"`
	Recursive                   *bool     `url:"recursive,omitempty" required:"false" json:"recursive,omitempty" path:"recursive"`
	Name                        string    `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Description                 string    `url:"description,omitempty" required:"false" json:"description,omitempty" path:"description"`
	Path                        string    `url:"path,omitempty" required:"true" json:"path,omitempty" path:"path"`
	Behavior                    string    `url:"behavior,omitempty" required:"true" json:"behavior,omitempty" path:"behavior"`
}

type BehaviorWebhookTestParams struct {
	Url      string                 `url:"url,omitempty" required:"true" json:"url,omitempty" path:"url"`
	Method   string                 `url:"method,omitempty" required:"false" json:"method,omitempty" path:"method"`
	Encoding string                 `url:"encoding,omitempty" required:"false" json:"encoding,omitempty" path:"encoding"`
	Headers  map[string]interface{} `url:"headers,omitempty" required:"false" json:"headers,omitempty" path:"headers"`
	Body     map[string]interface{} `url:"body,omitempty" required:"false" json:"body,omitempty" path:"body"`
	Action   string                 `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
}

type BehaviorUpdateParams struct {
	Id                          int64     `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	Value                       string    `url:"value,omitempty" required:"false" json:"value,omitempty" path:"value"`
	AttachmentFile              io.Writer `url:"attachment_file,omitempty" required:"false" json:"attachment_file,omitempty" path:"attachment_file"`
	DisableParentFolderBehavior *bool     `url:"disable_parent_folder_behavior,omitempty" required:"false" json:"disable_parent_folder_behavior,omitempty" path:"disable_parent_folder_behavior"`
	Recursive                   *bool     `url:"recursive,omitempty" required:"false" json:"recursive,omitempty" path:"recursive"`
	Name                        string    `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
	Description                 string    `url:"description,omitempty" required:"false" json:"description,omitempty" path:"description"`
	AttachmentDelete            *bool     `url:"attachment_delete,omitempty" required:"false" json:"attachment_delete,omitempty" path:"attachment_delete"`
}

type BehaviorDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
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
