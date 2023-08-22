package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Revision struct {
	Current  string `json:"current,omitempty" path:"current,omitempty" url:"current,omitempty"`
	Prior    string `json:"prior,omitempty" path:"prior,omitempty" url:"prior,omitempty"`
	Revision string `json:"revision,omitempty" path:"revision,omitempty" url:"revision,omitempty"`
	UpToDate *bool  `json:"up_to_date,omitempty" path:"up_to_date,omitempty" url:"up_to_date,omitempty"`
}

// Identifier no path or id

type RevisionCollection []Revision

type RevisionListParams struct {
	Action string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	ListParams
}

func (r *Revision) UnmarshalJSON(data []byte) error {
	type revision Revision
	var v revision
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = Revision(v)
	return nil
}

func (r *RevisionCollection) UnmarshalJSON(data []byte) error {
	type revisions RevisionCollection
	var v revisions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = RevisionCollection(v)
	return nil
}

func (r *RevisionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
