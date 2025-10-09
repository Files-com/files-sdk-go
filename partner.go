package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type Partner struct {
	AllowBypassing2faPolicies *bool  `json:"allow_bypassing_2fa_policies,omitempty" path:"allow_bypassing_2fa_policies,omitempty" url:"allow_bypassing_2fa_policies,omitempty"`
	AllowCredentialChanges    *bool  `json:"allow_credential_changes,omitempty" path:"allow_credential_changes,omitempty" url:"allow_credential_changes,omitempty"`
	AllowUserCreation         *bool  `json:"allow_user_creation,omitempty" path:"allow_user_creation,omitempty" url:"allow_user_creation,omitempty"`
	Id                        int64  `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Name                      string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Notes                     string `json:"notes,omitempty" path:"notes,omitempty" url:"notes,omitempty"`
	RootFolder                string `json:"root_folder,omitempty" path:"root_folder,omitempty" url:"root_folder,omitempty"`
}

func (p Partner) Identifier() interface{} {
	return p.Id
}

type PartnerCollection []Partner

type PartnerListParams struct {
	SortBy map[string]interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	ListParams
}

type PartnerFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

type PartnerCreateParams struct {
	AllowBypassing2faPolicies *bool  `url:"allow_bypassing_2fa_policies,omitempty" json:"allow_bypassing_2fa_policies,omitempty" path:"allow_bypassing_2fa_policies"`
	AllowCredentialChanges    *bool  `url:"allow_credential_changes,omitempty" json:"allow_credential_changes,omitempty" path:"allow_credential_changes"`
	AllowUserCreation         *bool  `url:"allow_user_creation,omitempty" json:"allow_user_creation,omitempty" path:"allow_user_creation"`
	Name                      string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	Notes                     string `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	RootFolder                string `url:"root_folder,omitempty" json:"root_folder,omitempty" path:"root_folder"`
}

type PartnerUpdateParams struct {
	Id                        int64  `url:"-,omitempty" json:"-,omitempty" path:"id"`
	AllowBypassing2faPolicies *bool  `url:"allow_bypassing_2fa_policies,omitempty" json:"allow_bypassing_2fa_policies,omitempty" path:"allow_bypassing_2fa_policies"`
	AllowCredentialChanges    *bool  `url:"allow_credential_changes,omitempty" json:"allow_credential_changes,omitempty" path:"allow_credential_changes"`
	AllowUserCreation         *bool  `url:"allow_user_creation,omitempty" json:"allow_user_creation,omitempty" path:"allow_user_creation"`
	Name                      string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	Notes                     string `url:"notes,omitempty" json:"notes,omitempty" path:"notes"`
	RootFolder                string `url:"root_folder,omitempty" json:"root_folder,omitempty" path:"root_folder"`
}

type PartnerDeleteParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

func (p *Partner) UnmarshalJSON(data []byte) error {
	type partner Partner
	var v partner
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = Partner(v)
	return nil
}

func (p *PartnerCollection) UnmarshalJSON(data []byte) error {
	type partners PartnerCollection
	var v partners
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PartnerCollection(v)
	return nil
}

func (p *PartnerCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
