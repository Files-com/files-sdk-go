package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// PartnerChannelUploadForm is a Files.com API resource.
type PartnerChannelUploadForm struct {
	// Upload destination path relative to the From Partner folder. Must be a configured leaf managed folder path, or an empty string for the From Partner folder when no managed folders are configured. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// ID of the Form Field Set answered for uploads to this destination. Must be in the same Workspace and contain at least one custom field.
	FormFieldSetId int64 `json:"form_field_set_id,omitempty" path:"form_field_set_id,omitempty" url:"form_field_set_id,omitempty"`
}

// Identifier returns the resource path.
func (p PartnerChannelUploadForm) Identifier() interface{} {
	return p.Path
}

// PartnerChannelUploadFormCollection is a list of PartnerChannelUploadForm resources.
type PartnerChannelUploadFormCollection []PartnerChannelUploadForm

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (p *PartnerChannelUploadForm) UnmarshalJSON(data []byte) error {
	type partnerChannelUploadForm PartnerChannelUploadForm
	var v partnerChannelUploadForm
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PartnerChannelUploadForm(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (p *PartnerChannelUploadFormCollection) UnmarshalJSON(data []byte) error {
	type partnerChannelUploadForms PartnerChannelUploadFormCollection
	var v partnerChannelUploadForms
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PartnerChannelUploadFormCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (p *PartnerChannelUploadFormCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
