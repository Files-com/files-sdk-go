package files_sdk

import (
	"encoding/json"
	"io"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Behavior is a Files.com API resource.
type Behavior struct {
	// Folder behavior ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Folder path. Note that Behavior paths cannot be updated once initially set. You will need to remove and re-create the behavior on the new path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// URL for attached file.
	AttachmentUrl string `json:"attachment_url,omitempty" path:"attachment_url,omitempty" url:"attachment_url,omitempty"`
	// Behavior type.
	Behavior string `json:"behavior,omitempty" path:"behavior,omitempty" url:"behavior,omitempty"`
	// Name for this behavior.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Description for this behavior.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Settings for this behavior. See the section above for an example value to provide here. Formatting is different for each Behavior type. May be sent as nested JSON or a single JSON-encoded string. If using XML encoding for the API call, this data must be sent as a JSON-encoded string.
	Value interface{} `json:"value,omitempty" path:"value,omitempty" url:"value,omitempty"`
	// Public URL for this publicly hosted folder when the `Serve Publicly` behavior has a key configured. When a Custom Domain with `public_hosting` destination is attached to this behavior, the URL uses that domain. Otherwise it uses the site's `subdomain.hosted-by-files.com` host, unless `custom_domain_only` is enabled, in which case no URL is returned.
	PublicHostingUrl string `json:"public_hosting_url,omitempty" path:"public_hosting_url,omitempty" url:"public_hosting_url,omitempty"`
	// If true, disables the inherited behavior for this folder and its children. Valid only for behavior types that child folders may override, and requires recursive to be true. Rejected for all other behavior types.
	DisableParentFolderBehavior *bool `json:"disable_parent_folder_behavior,omitempty" path:"disable_parent_folder_behavior,omitempty" url:"disable_parent_folder_behavior,omitempty"`
	// Whether this behavior is recursive for this record. `always` behaviors are always `true`, `never` behaviors are always `false`, and `sometimes` behaviors may be either value.
	Recursive *bool `json:"recursive,omitempty" path:"recursive,omitempty" url:"recursive,omitempty"`
	// If true, this behavior is inherited from a higher scope rather than owned by the requested workspace.
	Inherited *bool `json:"inherited,omitempty" path:"inherited,omitempty" url:"inherited,omitempty"`
	// If true, this behavior is controlled by a parent-site policy and cannot be modified locally.
	Managed *bool `json:"managed,omitempty" path:"managed,omitempty" url:"managed,omitempty"`
	// If true, this behavior may only be modified by a site admin because it is at the site root or disables a root behavior.
	RootBehaviorSiteAdminOnly *bool `json:"root_behavior_site_admin_only,omitempty" path:"root_behavior_site_admin_only,omitempty" url:"root_behavior_site_admin_only,omitempty"`
	// Certain behaviors may require a file, for instance, the `watermark` behavior requires a watermark image. Attach that file here.
	AttachmentFile io.Reader `json:"attachment_file,omitempty" path:"attachment_file,omitempty" url:"attachment_file,omitempty"`
	// If `true`, delete the file stored in `attachment`.
	AttachmentDelete *bool `json:"attachment_delete,omitempty" path:"attachment_delete,omitempty" url:"attachment_delete,omitempty"`
}

// Identifier returns the resource ID.
func (b Behavior) Identifier() interface{} {
	return b.Id
}

// BehaviorCollection is a list of Behavior resources.
type BehaviorCollection []Behavior

// BehaviorListParams contains the request parameters for GET /behaviors.
type BehaviorListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `behavior`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `clickwrap_id`, `form_field_set_id`, `impacts_ui`, `remote_server_id` or `behavior`. Valid field combinations are `[ impacts_ui, behavior ]` and `[ behavior, remote_server_id ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// BehaviorFindParams contains the request parameters for GET /behaviors/{id}.
type BehaviorFindParams struct {
	// Behavior ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// BehaviorListForParams contains the request parameters for GET /behaviors/folders/{path}.
type BehaviorListForParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `behavior`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `impacts_ui` and `behavior`. Valid field combinations are `[ impacts_ui, behavior ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// Path to operate on. Required; sent in the request path.
	Path string `url:"-,omitempty" json:"-,omitempty" path:"path"`
	// If `true`, behaviors above this path are shown.
	AncestorBehaviors *bool `url:"ancestor_behaviors,omitempty" json:"ancestor_behaviors,omitempty" path:"ancestor_behaviors"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// BehaviorCreateParams contains the request parameters for POST /behaviors.
type BehaviorCreateParams struct {
	// This field stores data specific to the type of behavior. See The Behavior Types section for the accepted value for each type of behavior.
	Value interface{} `url:"value,omitempty" json:"value,omitempty" path:"value"`
	// Certain behaviors may require a file, for instance, the `watermark` behavior requires a watermark image. Attach that file here.
	AttachmentFile io.Writer `url:"attachment_file,omitempty" json:"attachment_file,omitempty" path:"attachment_file"`
	// AttachmentFileAttachment sends file content as attachment_file; see lib.JSONAttachment. Leave AttachmentFile nil when you set it.
	AttachmentFileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"attachment_file"`
	// If true, disables the inherited behavior for this folder and its children. Valid only for behavior types that child folders may override, and requires recursive to be true. Rejected for all other behavior types.
	DisableParentFolderBehavior *bool `url:"disable_parent_folder_behavior,omitempty" json:"disable_parent_folder_behavior,omitempty" path:"disable_parent_folder_behavior"`
	// Whether the behavior should apply to child folders. This is only configurable for behavior types whose recursion mode is `sometimes`; `always` behaviors stay recursive and `never` behaviors stay non-recursive.
	Recursive *bool `url:"recursive,omitempty" json:"recursive,omitempty" path:"recursive"`
	// Name for this behavior.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Description for this behavior.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Path where this behavior should apply. Required.
	Path string `url:"path" json:"path" path:"path"`
	// Behavior type. Required.
	Behavior string `url:"behavior" json:"behavior" path:"behavior"`
}

// BehaviorWebhookTestParams contains the request parameters for POST /behaviors/webhook/test.
type BehaviorWebhookTestParams struct {
	// URL for testing the webhook. Required.
	Url string `url:"url" json:"url" path:"url"`
	// HTTP request method (GET or POST).
	Method string `url:"method,omitempty" json:"method,omitempty" path:"method"`
	// Encoding type for the webhook payload. Can be JSON, XML, or RAW (form data).
	Encoding string `url:"encoding,omitempty" json:"encoding,omitempty" path:"encoding"`
	// Additional request headers to send via HTTP.
	Headers interface{} `url:"headers,omitempty" json:"headers,omitempty" path:"headers"`
	// Additional body parameters to include in the webhook payload.
	Body interface{} `url:"body,omitempty" json:"body,omitempty" path:"body"`
	// Action for test body.
	Action string `url:"action,omitempty" json:"action,omitempty" path:"action"`
}

// BehaviorUpdateParams contains the request parameters for PATCH /behaviors/{id}.
type BehaviorUpdateParams struct {
	// Behavior ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// This field stores data specific to the type of behavior. See The Behavior Types section for the accepted value for each type of behavior.
	Value interface{} `url:"value,omitempty" json:"value,omitempty" path:"value"`
	// Certain behaviors may require a file, for instance, the `watermark` behavior requires a watermark image. Attach that file here.
	AttachmentFile io.Writer `url:"attachment_file,omitempty" json:"attachment_file,omitempty" path:"attachment_file"`
	// AttachmentFileAttachment sends file content as attachment_file; see lib.JSONAttachment. Leave AttachmentFile nil when you set it.
	AttachmentFileAttachment *lib.JSONAttachment `json:"-" url:"-" path:"-" attachment:"attachment_file"`
	// If true, disables the inherited behavior for this folder and its children. Valid only for behavior types that child folders may override, and requires recursive to be true. Rejected for all other behavior types.
	DisableParentFolderBehavior *bool `url:"disable_parent_folder_behavior,omitempty" json:"disable_parent_folder_behavior,omitempty" path:"disable_parent_folder_behavior"`
	// Whether the behavior should apply to child folders. This is only configurable for behavior types whose recursion mode is `sometimes`; `always` behaviors stay recursive and `never` behaviors stay non-recursive.
	Recursive *bool `url:"recursive,omitempty" json:"recursive,omitempty" path:"recursive"`
	// Name for this behavior.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Description for this behavior.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// If `true`, delete the file stored in `attachment`.
	AttachmentDelete *bool `url:"attachment_delete,omitempty" json:"attachment_delete,omitempty" path:"attachment_delete"`
}

// BehaviorDeleteParams contains the request parameters for DELETE /behaviors/{id}.
type BehaviorDeleteParams struct {
	// Behavior ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (b *Behavior) UnmarshalJSON(data []byte) error {
	type behavior Behavior
	var v behavior
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*b = Behavior(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (b *BehaviorCollection) UnmarshalJSON(data []byte) error {
	type behaviors BehaviorCollection
	var v behaviors
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*b = BehaviorCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (b *BehaviorCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*b))
	for i, v := range *b {
		ret[i] = v
	}

	return &ret
}
