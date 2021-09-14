package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Priority struct {
	Path  string `json:"path,omitempty"`
	Color string `json:"color,omitempty"`
}

type PriorityCollection []Priority

type PriorityListParams struct {
	Cursor  string `url:"cursor,omitempty" required:"false"`
	PerPage int64  `url:"per_page,omitempty" required:"false"`
	Path    string `url:"path,omitempty" required:"true"`
	lib.ListParams
}

func (p *Priority) UnmarshalJSON(data []byte) error {
	type priority Priority
	var v priority
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = Priority(v)
	return nil
}

func (p *PriorityCollection) UnmarshalJSON(data []byte) error {
	type prioritys []Priority
	var v prioritys
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*p = PriorityCollection(v)
	return nil
}

func (p *PriorityCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
