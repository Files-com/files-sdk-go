package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// BundleAction is a Files.com API resource.
type BundleAction struct {
	// Type of action.
	Action string `json:"action,omitempty" path:"action,omitempty" url:"action,omitempty"`
	// Object that contains bundle registration information.
	BundleRegistration BundleRegistration `json:"bundle_registration,omitempty" path:"bundle_registration,omitempty" url:"bundle_registration,omitempty"`
	// Action occurrence date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// The destination path for this bundle action, if applicable.
	Destination string `json:"destination,omitempty" path:"destination,omitempty" url:"destination,omitempty"`
	// Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// The source path for this bundle action, if applicable.
	Source string `json:"source,omitempty" path:"source,omitempty" url:"source,omitempty"`
	// Identifier of the associated bundle registration.
	BundleRegistrationId int64 `json:"bundle_registration_id,omitempty" path:"bundle_registration_id,omitempty" url:"bundle_registration_id,omitempty"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Name of the registrant who performed the action.
	BundleRegistrationName string `json:"bundle_registration_name,omitempty" path:"bundle_registration_name,omitempty" url:"bundle_registration_name,omitempty"`
	// Email of the registrant.
	BundleRegistrationEmail string `json:"bundle_registration_email,omitempty" path:"bundle_registration_email,omitempty" url:"bundle_registration_email,omitempty"`
	// IP address of the registrant.
	BundleRegistrationIp string `json:"bundle_registration_ip,omitempty" path:"bundle_registration_ip,omitempty" url:"bundle_registration_ip,omitempty"`
	// Id of associated form field set.
	FormFieldSetId int64 `json:"form_field_set_id,omitempty" path:"form_field_set_id,omitempty" url:"form_field_set_id,omitempty"`
	// Data for form field set with form field ids as keys and user data as values.
	FormFieldData interface{} `json:"form_field_data,omitempty" path:"form_field_data,omitempty" url:"form_field_data,omitempty"`
}

// Identifier returns the resource path.
func (b BundleAction) Identifier() interface{} {
	return b.Path
}

// BundleActionCollection is a list of BundleAction resources.
type BundleActionCollection []BundleAction

// BundleActionListParams contains the request parameters for GET /bundle_actions.
type BundleActionListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `created_at` or `bundle_registration_id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at`, `bundle_id` or `bundle_registration_id`. Valid field combinations are `[ bundle_id, created_at ]`, `[ bundle_registration_id, created_at ]`, `[ bundle_id, bundle_registration_id ]` or `[ bundle_id, bundle_registration_id, created_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// Bundle ID.
	BundleId int64 `url:"bundle_id,omitempty" json:"bundle_id,omitempty" path:"bundle_id"`
	// Bundle Registration ID.
	BundleRegistrationId int64 `url:"bundle_registration_id,omitempty" json:"bundle_registration_id,omitempty" path:"bundle_registration_id"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (b *BundleAction) UnmarshalJSON(data []byte) error {
	type bundleAction BundleAction
	var v bundleAction
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = BundleAction(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (b *BundleActionCollection) UnmarshalJSON(data []byte) error {
	type bundleActions BundleActionCollection
	var v bundleActions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BundleActionCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (b *BundleActionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
