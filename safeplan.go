package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type SafePlan struct {
}

// Identifier no path or id

type SafePlanCollection []SafePlan

func (s *SafePlan) UnmarshalJSON(data []byte) error {
	type safePlan SafePlan
	var v safePlan
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = SafePlan(v)
	return nil
}

func (s *SafePlanCollection) UnmarshalJSON(data []byte) error {
	type safePlans SafePlanCollection
	var v safePlans
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SafePlanCollection(v)
	return nil
}

func (s *SafePlanCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
