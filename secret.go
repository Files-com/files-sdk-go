package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Secret is a Files.com API resource.
type Secret struct {
	Id              int64       `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	WorkspaceId     int64       `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	Name            string      `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	Description     string      `json:"description,omitempty" path:"description,omitempty" url:"description,omitempty"`
	SecretType      string      `json:"secret_type,omitempty" path:"secret_type,omitempty" url:"secret_type,omitempty"`
	Metadata        interface{} `json:"metadata,omitempty" path:"metadata,omitempty" url:"metadata,omitempty"`
	ValueFieldNames []string    `json:"value_field_names,omitempty" path:"value_field_names,omitempty" url:"value_field_names,omitempty"`
	CreatedAt       *time.Time  `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt       *time.Time  `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
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
	SortBy       interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	Filter       interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

// SecretFindParams contains the request parameters for GET /secrets/{id}.
type SecretFindParams struct {
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// SecretCreateParams contains the request parameters for POST /secrets.
type SecretCreateParams struct {
	Name        string               `url:"name" json:"name" path:"name"`
	Description string               `url:"description,omitempty" json:"description,omitempty" path:"description"`
	SecretType  SecretSecretTypeEnum `url:"secret_type" json:"secret_type" path:"secret_type"`
	Metadata    interface{}          `url:"metadata,omitempty" json:"metadata,omitempty" path:"metadata"`
	WorkspaceId int64                `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// SecretUpdateParams contains the request parameters for PATCH /secrets/{id}.
type SecretUpdateParams struct {
	Id          int64                `url:"-,omitempty" json:"-,omitempty" path:"id"`
	Name        string               `url:"name,omitempty" json:"name,omitempty" path:"name"`
	Description string               `url:"description,omitempty" json:"description,omitempty" path:"description"`
	SecretType  SecretSecretTypeEnum `url:"secret_type,omitempty" json:"secret_type,omitempty" path:"secret_type"`
	Metadata    interface{}          `url:"metadata,omitempty" json:"metadata,omitempty" path:"metadata"`
}

// SecretDeleteParams contains the request parameters for DELETE /secrets/{id}.
type SecretDeleteParams struct {
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
