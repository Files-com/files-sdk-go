package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Secret is a Files.com API resource.
type Secret struct {
	// Secret ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID. 0 means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Secret name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Internal description for your reference.
	Description string `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	// Secret type.
	SecretType string `json:"secret_type,omitempty" path:"secret_type,omitempty" url:"secret_type,omitempty"`
	// Non-secret metadata for the Secret type.
	Metadata interface{} `json:"metadata,omitempty" path:"metadata,omitempty" url:"metadata,omitempty"`
	// Names of configured secret value fields. Secret values are never returned.
	ValueFieldNames []string `json:"value_field_names,omitempty" path:"value_field_names,omitempty" url:"value_field_names,omitempty"`
	// Secret create date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Secret update date/time.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (s Secret) Identifier() interface{} {
	return s.Id
}

// SecretCollection is a list of Secret resources.
type SecretCollection []Secret

// SecretSecretTypeEnum is a string value for secret_type.
// Enum lists the values documented by the API.
type SecretSecretTypeEnum string

// String returns the API parameter value.
func (u SecretSecretTypeEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u SecretSecretTypeEnum) Enum() map[string]SecretSecretTypeEnum {
	return map[string]SecretSecretTypeEnum{
		"basic":       SecretSecretTypeEnum("basic"),
		"token":       SecretSecretTypeEnum("token"),
		"headers":     SecretSecretTypeEnum("headers"),
		"certificate": SecretSecretTypeEnum("certificate"),
		"key_value":   SecretSecretTypeEnum("key_value"),
	}
}

// SecretListParams contains the request parameters for GET /secrets.
type SecretListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name` or `secret_type`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`, `name` or `secret_type`. Valid field combinations are `[ workspace_id, name ]`, `[ workspace_id, secret_type ]`, `[ secret_type, name ]` or `[ workspace_id, secret_type, name ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// SecretFindParams contains the request parameters for GET /secrets/{id}.
type SecretFindParams struct {
	// Secret ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// SecretCreateParams contains the request parameters for POST /secrets.
type SecretCreateParams struct {
	// Secret name. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Internal description for your reference.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Secret type. Values: basic, token, headers, certificate, key_value. Required.
	SecretType SecretSecretTypeEnum `url:"secret_type" json:"secret_type" path:"secret_type"`
	// Non-secret metadata for the Secret type.
	Metadata interface{} `url:"metadata,omitempty" json:"metadata,omitempty" path:"metadata"`
	// Workspace ID. 0 means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// SecretUpdateParams contains the request parameters for PATCH /secrets/{id}.
type SecretUpdateParams struct {
	// Secret ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Secret name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Internal description for your reference.
	Description string `url:"description,omitempty" json:"description,omitempty" path:"description"`
	// Secret type. Values: basic, token, headers, certificate, key_value.
	SecretType SecretSecretTypeEnum `url:"secret_type,omitempty" json:"secret_type,omitempty" path:"secret_type"`
	// Non-secret metadata for the Secret type.
	Metadata interface{} `url:"metadata,omitempty" json:"metadata,omitempty" path:"metadata"`
}

// SecretDeleteParams contains the request parameters for DELETE /secrets/{id}.
type SecretDeleteParams struct {
	// Secret ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *Secret) UnmarshalJSON(data []byte) error {
	type secret Secret
	var v secret
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Secret(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SecretCollection) UnmarshalJSON(data []byte) error {
	type secrets SecretCollection
	var v secrets
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SecretCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SecretCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
