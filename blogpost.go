package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type BlogPost struct {
	Id          int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Title       string     `json:"title,omitempty" path:"title,omitempty" url:"title,omitempty"`
	Content     string     `json:"content,omitempty" path:"content,omitempty" url:"content,omitempty"`
	Link        *time.Time `json:"link,omitempty" path:"link,omitempty" url:"link,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty" path:"published_at,omitempty" url:"published_at,omitempty"`
}

func (b BlogPost) Identifier() interface{} {
	return b.Id
}

type BlogPostCollection []BlogPost

type BlogPostListParams struct {
	Action     string                 `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	SortBy     map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter     BlogPost               `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterGt   map[string]interface{} `url:"filter_gt,omitempty" required:"false" json:"filter_gt,omitempty" path:"filter_gt"`
	FilterGteq map[string]interface{} `url:"filter_gteq,omitempty" required:"false" json:"filter_gteq,omitempty" path:"filter_gteq"`
	FilterLt   map[string]interface{} `url:"filter_lt,omitempty" required:"false" json:"filter_lt,omitempty" path:"filter_lt"`
	FilterLteq map[string]interface{} `url:"filter_lteq,omitempty" required:"false" json:"filter_lteq,omitempty" path:"filter_lteq"`
	ListParams
}

func (b *BlogPost) UnmarshalJSON(data []byte) error {
	type blogPost BlogPost
	var v blogPost
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BlogPost(v)
	return nil
}

func (b *BlogPostCollection) UnmarshalJSON(data []byte) error {
	type blogPosts BlogPostCollection
	var v blogPosts
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BlogPostCollection(v)
	return nil
}

func (b *BlogPostCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
