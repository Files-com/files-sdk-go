package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type PartnerConnection struct {
	Id        int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Role      string `json:"role,omitempty" path:"role,omitempty" url:"role,omitempty"`
	SiteId    int64  `json:"site_id,omitempty" path:"site_id,omitempty" url:"site_id,omitempty"`
	SiteName  string `json:"site_name,omitempty" path:"site_name,omitempty" url:"site_name,omitempty"`
	MountPath string `json:"mount_path,omitempty" path:"mount_path,omitempty" url:"mount_path,omitempty"`
}

func (p PartnerConnection) Identifier() interface{} {
	return p.Id
}

type PartnerConnectionCollection []PartnerConnection

func (p *PartnerConnection) UnmarshalJSON(data []byte) error {
	type partnerConnection PartnerConnection
	var v partnerConnection
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PartnerConnection(v)
	return nil
}

func (p *PartnerConnectionCollection) UnmarshalJSON(data []byte) error {
	type partnerConnections PartnerConnectionCollection
	var v partnerConnections
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PartnerConnectionCollection(v)
	return nil
}

func (p *PartnerConnectionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
