package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type IpAbuseEntry struct {
	Ip       string `json:"ip,omitempty" path:"ip,omitempty" url:"ip,omitempty"`
	List     string `json:"list,omitempty" path:"list,omitempty" url:"list,omitempty"`
	Hostname string `json:"hostname,omitempty" path:"hostname,omitempty" url:"hostname,omitempty"`
	Reason   string `json:"reason,omitempty" path:"reason,omitempty" url:"reason,omitempty"`
}

// Identifier no path or id

type IpAbuseEntryCollection []IpAbuseEntry

type IpAbuseEntryCreateParams struct {
	Ip       string `url:"ip,omitempty" required:"true" json:"ip,omitempty" path:"ip"`
	List     string `url:"list,omitempty" required:"true" json:"list,omitempty" path:"list"`
	Hostname string `url:"hostname,omitempty" required:"false" json:"hostname,omitempty" path:"hostname"`
	Reason   string `url:"reason,omitempty" required:"false" json:"reason,omitempty" path:"reason"`
}

func (i *IpAbuseEntry) UnmarshalJSON(data []byte) error {
	type ipAbuseEntry IpAbuseEntry
	var v ipAbuseEntry
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*i = IpAbuseEntry(v)
	return nil
}

func (i *IpAbuseEntryCollection) UnmarshalJSON(data []byte) error {
	type ipAbuseEntrys IpAbuseEntryCollection
	var v ipAbuseEntrys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*i = IpAbuseEntryCollection(v)
	return nil
}

func (i *IpAbuseEntryCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*i))
	for i, v := range *i {
		ret[i] = v
	}

	return &ret
}
