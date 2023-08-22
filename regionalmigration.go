package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type RegionalMigration struct {
}

// Identifier no path or id

type RegionalMigrationCollection []RegionalMigration

func (r *RegionalMigration) UnmarshalJSON(data []byte) error {
	type regionalMigration RegionalMigration
	var v regionalMigration
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = RegionalMigration(v)
	return nil
}

func (r *RegionalMigrationCollection) UnmarshalJSON(data []byte) error {
	type regionalMigrations RegionalMigrationCollection
	var v regionalMigrations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = RegionalMigrationCollection(v)
	return nil
}

func (r *RegionalMigrationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
