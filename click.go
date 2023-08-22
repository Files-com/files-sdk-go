package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Click struct {
	Code      string `json:"code,omitempty" path:"code,omitempty" url:"code,omitempty"`
	PartnerId int64  `json:"partner_id,omitempty" path:"partner_id,omitempty" url:"partner_id,omitempty"`
	Subid     string `json:"subid,omitempty" path:"subid,omitempty" url:"subid,omitempty"`
	Ip        string `json:"ip,omitempty" path:"ip,omitempty" url:"ip,omitempty"`
	Referer   string `json:"referer,omitempty" path:"referer,omitempty" url:"referer,omitempty"`
}

// Identifier no path or id

type ClickCollection []Click

type ClickCreateParams struct {
	PartnerId int64  `url:"partner_id,omitempty" required:"true" json:"partner_id,omitempty" path:"partner_id"`
	Subid     string `url:"subid,omitempty" required:"false" json:"subid,omitempty" path:"subid"`
	Ip        string `url:"ip,omitempty" required:"false" json:"ip,omitempty" path:"ip"`
	Referer   string `url:"referer,omitempty" required:"false" json:"referer,omitempty" path:"referer"`
}

func (c *Click) UnmarshalJSON(data []byte) error {
	type click Click
	var v click
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*c = Click(v)
	return nil
}

func (c *ClickCollection) UnmarshalJSON(data []byte) error {
	type clicks ClickCollection
	var v clicks
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*c = ClickCollection(v)
	return nil
}

func (c *ClickCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*c))
	for i, v := range *c {
		ret[i] = v
	}

	return &ret
}
