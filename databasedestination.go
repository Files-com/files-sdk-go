package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// DatabaseDestination is a Files.com API resource.
type DatabaseDestination struct {
	// Database Destination ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID. 0 means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Name of the destination.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Database provider.
	Provider string `json:"provider,omitempty" path:"provider,omitempty" url:"provider,omitempty"`
	// Provider settings for the destination. Use the schema matching provider. The connection address and table are fixed after creation.
	Config interface{} `json:"config,omitempty" path:"config,omitempty" url:"config,omitempty"`
	// ID of a Secret in the same Workspace. Use a Token Secret for Snowflake, a service account Key/Value Secret for BigQuery, or a Username/Password Secret for a database login.
	SecretId int64 `json:"secret_id,omitempty" path:"secret_id,omitempty" url:"secret_id,omitempty"`
	// When the destination was created.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// When the destination was last updated.
	UpdatedAt *time.Time `json:"updated_at,omitempty" path:"updated_at,omitempty" url:"updated_at,omitempty"`
}

// Identifier returns the resource ID.
func (d DatabaseDestination) Identifier() interface{} {
	return d.Id
}

// DatabaseDestinationCollection is a list of DatabaseDestination resources.
type DatabaseDestinationCollection []DatabaseDestination

// DatabaseDestinationProviderEnum is a string value for provider.
// Enum lists the values documented by the API.
type DatabaseDestinationProviderEnum string

// String returns the API parameter value.
func (u DatabaseDestinationProviderEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u DatabaseDestinationProviderEnum) Enum() map[string]DatabaseDestinationProviderEnum {
	return map[string]DatabaseDestinationProviderEnum{
		"snowflake":  DatabaseDestinationProviderEnum("snowflake"),
		"bigquery":   DatabaseDestinationProviderEnum("bigquery"),
		"supabase":   DatabaseDestinationProviderEnum("supabase"),
		"postgresql": DatabaseDestinationProviderEnum("postgresql"),
		"mysql":      DatabaseDestinationProviderEnum("mysql"),
	}
}

// DatabaseDestinationListParams contains the request parameters for GET /database_destinations.
type DatabaseDestinationListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id` and `name`. Valid field combinations are `[ workspace_id, name ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// DatabaseDestinationFindParams contains the request parameters for GET /database_destinations/{id}.
type DatabaseDestinationFindParams struct {
	// Database Destination ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// DatabaseDestinationCreateParams contains the request parameters for POST /database_destinations.
type DatabaseDestinationCreateParams struct {
	// Name of the destination. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Database provider. Values: snowflake, bigquery, supabase, postgresql, mysql. Required.
	Provider DatabaseDestinationProviderEnum `url:"provider" json:"provider" path:"provider"`
	// Provider settings for the destination. Use the schema matching provider. The connection address and table are fixed after creation. Required.
	Config interface{} `url:"config" json:"config" path:"config"`
	// ID of a Secret in the same Workspace. Use a Token Secret for Snowflake, a service account Key/Value Secret for BigQuery, or a Username/Password Secret for a database login. Required.
	SecretId int64 `url:"secret_id" json:"secret_id" path:"secret_id"`
	// Workspace ID. 0 means the default workspace.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
}

// DatabaseDestinationUpdateParams contains the request parameters for PATCH /database_destinations/{id}.
type DatabaseDestinationUpdateParams struct {
	// Database Destination ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Name of the destination.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// ID of a Secret in the same Workspace. Use a Token Secret for Snowflake, a service account Key/Value Secret for BigQuery, or a Username/Password Secret for a database login.
	SecretId int64 `url:"secret_id,omitempty" json:"secret_id,omitempty" path:"secret_id"`
	// Provider settings for the destination. Use the schema matching provider. The connection address and table are fixed after creation.
	Config interface{} `url:"config,omitempty" json:"config,omitempty" path:"config"`
}

// DatabaseDestinationDeleteParams contains the request parameters for DELETE /database_destinations/{id}.
type DatabaseDestinationDeleteParams struct {
	// Database Destination ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (d *DatabaseDestination) UnmarshalJSON(data []byte) error {
	type databaseDestination DatabaseDestination
	var v databaseDestination
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*d = DatabaseDestination(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (d *DatabaseDestinationCollection) UnmarshalJSON(data []byte) error {
	type databaseDestinations DatabaseDestinationCollection
	var v databaseDestinations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*d = DatabaseDestinationCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (d *DatabaseDestinationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*d))
	for i, v := range *d {
		ret[i] = v
	}

	return &ret
}
