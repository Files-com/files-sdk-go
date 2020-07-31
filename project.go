package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Project struct {
	Id           int    `json:"id,omitempty"`
	GlobalAccess string `json:"global_access,omitempty"`
}

type ProjectCollection []Project

type ProjectListParams struct {
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Action  string `url:"action,omitempty"`
	lib.ListParams
}

type ProjectFindParams struct {
	Id int `url:"-,omitempty"`
}

type ProjectCreateParams struct {
	GlobalAccess string `url:"global_access,omitempty"`
}

type ProjectUpdateParams struct {
	Id           int    `url:"-,omitempty"`
	GlobalAccess string `url:"global_access,omitempty"`
}

type ProjectDeleteParams struct {
	Id int `url:"-,omitempty"`
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
