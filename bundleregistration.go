package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// BundleRegistration is a Files.com API resource.
type BundleRegistration struct {
	// Registration cookie code.
	Code string `json:"code,omitempty" path:"code,omitempty" url:"code,omitempty"`
	// Registrant name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Registrant company name.
	Company string `json:"company,omitempty" path:"company,omitempty" url:"company,omitempty"`
	// Registrant email address.
	Email string `json:"email,omitempty" path:"email,omitempty" url:"email,omitempty"`
	// Registrant IP Address.
	Ip string `json:"ip,omitempty" path:"ip,omitempty" url:"ip,omitempty"`
	// InboxRegistration cookie code, if there is an associated InboxRegistration.
	InboxCode string `json:"inbox_code,omitempty" path:"inbox_code,omitempty" url:"inbox_code,omitempty"`
	// Clickwrap text that was shown to the registrant.
	ClickwrapBody string `json:"clickwrap_body,omitempty" path:"clickwrap_body,omitempty" url:"clickwrap_body,omitempty"`
	// Id of associated form field set.
	FormFieldSetId int64 `json:"form_field_set_id,omitempty" path:"form_field_set_id,omitempty" url:"form_field_set_id,omitempty"`
	// Data for form field set with form field ids as keys and user data as values.
	FormFieldData interface{} `json:"form_field_data,omitempty" path:"form_field_data,omitempty" url:"form_field_data,omitempty"`
	// Bundle URL code.
	BundleCode string `json:"bundle_code,omitempty" path:"bundle_code,omitempty" url:"bundle_code,omitempty"`
	// Id of associated bundle.
	BundleId int64 `json:"bundle_id,omitempty" path:"bundle_id,omitempty" url:"bundle_id,omitempty"`
	// Id of associated bundle recipient.
	BundleRecipientId int64 `json:"bundle_recipient_id,omitempty" path:"bundle_recipient_id,omitempty" url:"bundle_recipient_id,omitempty"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Registration creation date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
}

// Identifier no path or id

// BundleRegistrationCollection is a list of BundleRegistration resources.
type BundleRegistrationCollection []BundleRegistration

// BundleRegistrationListParams contains the request parameters for GET /bundle_registrations.
type BundleRegistrationListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `bundle_id` or `created_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// ID of the associated Bundle.
	BundleId int64 `url:"bundle_id,omitempty" json:"bundle_id,omitempty" path:"bundle_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (b *BundleRegistration) UnmarshalJSON(data []byte) error {
	type bundleRegistration BundleRegistration
	var v bundleRegistration
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BundleRegistration(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (b *BundleRegistrationCollection) UnmarshalJSON(data []byte) error {
	type bundleRegistrations BundleRegistrationCollection
	var v bundleRegistrations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BundleRegistrationCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (b *BundleRegistrationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
