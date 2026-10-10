package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// HolidayRegion is a Files.com API resource.
type HolidayRegion struct {
	Code string `json:"code,omitempty" path:"code,omitempty" url:"code,omitempty"`
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
}

// Identifier no path or id

// HolidayRegionCollection is a list of HolidayRegion resources.
type HolidayRegionCollection []HolidayRegion

// HolidayRegionGetSupportedParams contains the request parameters for GET /holiday_regions/supported.
type HolidayRegionGetSupportedParams struct {
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (h *HolidayRegion) UnmarshalJSON(data []byte) error {
	type holidayRegion HolidayRegion
	var v holidayRegion
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*h = HolidayRegion(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (h *HolidayRegionCollection) UnmarshalJSON(data []byte) error {
	type holidayRegions HolidayRegionCollection
	var v holidayRegions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*h = HolidayRegionCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (h *HolidayRegionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*h))
	for i, v := range *h {
		ret[i] = v
	}

	return &ret
}
