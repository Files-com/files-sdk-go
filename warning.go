package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Warning struct {
	Warnings []string `json:"warnings,omitempty" path:"warnings,omitempty" url:"warnings,omitempty"`
}

// Identifier no path or id

type WarningCollection []Warning

type WarningListParams struct {
	Action string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	ListParams
}

func (w *Warning) UnmarshalJSON(data []byte) error {
	type warning Warning
	var v warning
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*w = Warning(v)
	return nil
}

func (w *WarningCollection) UnmarshalJSON(data []byte) error {
	type warnings WarningCollection
	var v warnings
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*w = WarningCollection(v)
	return nil
}

func (w *WarningCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*w))
	for i, v := range *w {
		ret[i] = v
	}

	return &ret
}
