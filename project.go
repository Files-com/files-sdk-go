package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Project struct {
	Id           int64  `json:"id,omitempty"`
	GlobalAccess string `json:"global_access,omitempty"`
}

type ProjectCollection []Project

type ProjectListParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int64  `url:"per_page,omitempty" required:"false"`
	lib.ListParams
}

type ProjectFindParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

type ProjectCreateParams struct {
	GlobalAccess string `url:"global_access,omitempty" required:"true"`
}

type ProjectUpdateParams struct {
	Id           int64  `url:"-,omitempty" required:"true"`
	GlobalAccess string `url:"global_access,omitempty" required:"true"`
}

type ProjectDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"true"`
}

func (p *Project) UnmarshalJSON(data []byte) error {
	type project Project
	var v project
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = Project(v)
	return nil
}

func (p *ProjectCollection) UnmarshalJSON(data []byte) error {
	type projects []Project
	var v projects
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = ProjectCollection(v)
	return nil
}

func (p *ProjectCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
