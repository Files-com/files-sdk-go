package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Announcement struct {
	Id         int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Headline   string     `json:"headline,omitempty" path:"headline,omitempty" url:"headline,omitempty"`
	Body       string     `json:"body,omitempty" path:"body,omitempty" url:"body,omitempty"`
	ButtonText string     `json:"button_text,omitempty" path:"button_text,omitempty" url:"button_text,omitempty"`
	ButtonUrl  string     `json:"button_url,omitempty" path:"button_url,omitempty" url:"button_url,omitempty"`
	HtmlBody   string     `json:"html_body,omitempty" path:"html_body,omitempty" url:"html_body,omitempty"`
	Label      string     `json:"label,omitempty" path:"label,omitempty" url:"label,omitempty"`
	LabelColor string     `json:"label_color,omitempty" path:"label_color,omitempty" url:"label_color,omitempty"`
	PublishAt  *time.Time `json:"publish_at,omitempty" path:"publish_at,omitempty" url:"publish_at,omitempty"`
	Slug       string     `json:"slug,omitempty" path:"slug,omitempty" url:"slug,omitempty"`
}

func (a Announcement) Identifier() interface{} {
	return a.Id
}

type AnnouncementCollection []Announcement

type AnnouncementListParams struct {
	Action string                 `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	SortBy map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	ListParams
}

func (a *Announcement) UnmarshalJSON(data []byte) error {
	type announcement Announcement
	var v announcement
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = Announcement(v)
	return nil
}

func (a *AnnouncementCollection) UnmarshalJSON(data []byte) error {
	type announcements AnnouncementCollection
	var v announcements
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AnnouncementCollection(v)
	return nil
}

func (a *AnnouncementCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
