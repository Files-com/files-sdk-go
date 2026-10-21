package files_sdk

import (
	"encoding/json"
	"net/url"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// RemoteMountBackend is a Files.com API resource.
type RemoteMountBackend struct {
	// Path to the canary file used for health checks.
	CanaryFilePath string `json:"canary_file_path,omitempty" path:"canary_file_path,omitempty" url:"canary_file_path,omitempty"`
	// True if this backend is enabled.
	Enabled *bool `json:"enabled,omitempty" path:"enabled,omitempty" url:"enabled,omitempty"`
	// Number of consecutive failures before considering the backend unhealthy.
	Fall int64 `json:"fall,omitempty" path:"fall,omitempty" url:"fall,omitempty"`
	// True if health checks are enabled for this backend.
	HealthCheckEnabled *bool `json:"health_check_enabled,omitempty" path:"health_check_enabled,omitempty" url:"health_check_enabled,omitempty"`
	// Array of recent health check results.
	HealthCheckResults []map[string]interface{} `json:"health_check_results,omitempty" path:"health_check_results,omitempty" url:"health_check_results,omitempty"`
	// Type of health check to perform.
	HealthCheckType string `json:"health_check_type,omitempty" path:"health_check_type,omitempty" url:"health_check_type,omitempty"`
	// Unique identifier for this backend.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Interval in seconds between health checks.
	Interval int64 `json:"interval,omitempty" path:"interval,omitempty" url:"interval,omitempty"`
	// Minimum free CPU percentage required for this backend to be considered healthy.
	MinFreeCpu string `json:"min_free_cpu,omitempty" path:"min_free_cpu,omitempty" url:"min_free_cpu,omitempty"`
	// Minimum free memory percentage required for this backend to be considered healthy.
	MinFreeMem string `json:"min_free_mem,omitempty" path:"min_free_mem,omitempty" url:"min_free_mem,omitempty"`
	// Priority of this backend.
	Priority int64 `json:"priority,omitempty" path:"priority,omitempty" url:"priority,omitempty"`
	// Path on the remote server to treat as the root of this mount.
	RemotePath string `json:"remote_path,omitempty" path:"remote_path,omitempty" url:"remote_path,omitempty"`
	// The remote server that this backend is associated with.
	RemoteServerId int64 `json:"remote_server_id,omitempty" path:"remote_server_id,omitempty" url:"remote_server_id,omitempty"`
	// The mount ID of the Remote Server Mount that this backend is associated with.
	RemoteServerMountId int64 `json:"remote_server_mount_id,omitempty" path:"remote_server_mount_id,omitempty" url:"remote_server_mount_id,omitempty"`
	// Number of consecutive successes before considering the backend healthy.
	Rise int64 `json:"rise,omitempty" path:"rise,omitempty" url:"rise,omitempty"`
	// Status of this backend.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// True if this backend is undergoing maintenance.
	UndergoingMaintenance *bool `json:"undergoing_maintenance,omitempty" path:"undergoing_maintenance,omitempty" url:"undergoing_maintenance,omitempty"`
}

// Identifier returns the resource ID.
func (r RemoteMountBackend) Identifier() interface{} {
	return r.Id
}

// RemoteMountBackendCollection is a list of RemoteMountBackend resources.
type RemoteMountBackendCollection []RemoteMountBackend

// RemoteMountBackendHealthCheckTypeEnum is a string value for health_check_type.
// Enum lists the values documented by the API.
type RemoteMountBackendHealthCheckTypeEnum string

// String returns the API parameter value.
func (u RemoteMountBackendHealthCheckTypeEnum) String() string {
	return string(u)
}

// Enum returns the documented values keyed by their API strings.
func (u RemoteMountBackendHealthCheckTypeEnum) Enum() map[string]RemoteMountBackendHealthCheckTypeEnum {
	return map[string]RemoteMountBackendHealthCheckTypeEnum{
		"active":  RemoteMountBackendHealthCheckTypeEnum("active"),
		"passive": RemoteMountBackendHealthCheckTypeEnum("passive"),
	}
}

// RemoteMountBackendListParams contains the request parameters for GET /remote_mount_backends.
type RemoteMountBackendListParams struct {
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `remote_server_mount_id`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// RemoteMountBackendFindParams contains the request parameters for GET /remote_mount_backends/{id}.
type RemoteMountBackendFindParams struct {
	// Remote Mount Backend ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// RemoteMountBackendCreateParams contains the request parameters for POST /remote_mount_backends.
type RemoteMountBackendCreateParams struct {
	// True if this backend is enabled.
	Enabled *bool `url:"enabled,omitempty" json:"enabled,omitempty" path:"enabled"`
	// Number of consecutive failures before considering the backend unhealthy.
	Fall int64 `url:"fall,omitempty" json:"fall,omitempty" path:"fall"`
	// True if health checks are enabled for this backend.
	HealthCheckEnabled *bool `url:"health_check_enabled,omitempty" json:"health_check_enabled,omitempty" path:"health_check_enabled"`
	// Type of health check to perform. Values: active, passive.
	HealthCheckType RemoteMountBackendHealthCheckTypeEnum `url:"health_check_type,omitempty" json:"health_check_type,omitempty" path:"health_check_type"`
	// Interval in seconds between health checks.
	Interval int64 `url:"interval,omitempty" json:"interval,omitempty" path:"interval"`
	// Minimum free CPU percentage required for this backend to be considered healthy.
	MinFreeCpu float64 `url:"min_free_cpu,omitempty" json:"min_free_cpu,omitempty" path:"min_free_cpu"`
	// MinFreeCpuDecimal sends exact decimal text, such as "1.5", as min_free_cpu. Leave MinFreeCpu zero when you set it.
	MinFreeCpuDecimal *string `json:"-" url:"-" path:"-"`
	// Minimum free memory percentage required for this backend to be considered healthy.
	MinFreeMem float64 `url:"min_free_mem,omitempty" json:"min_free_mem,omitempty" path:"min_free_mem"`
	// MinFreeMemDecimal sends exact decimal text, such as "1.5", as min_free_mem. Leave MinFreeMem zero when you set it.
	MinFreeMemDecimal *string `json:"-" url:"-" path:"-"`
	// Priority of this backend.
	Priority int64 `url:"priority,omitempty" json:"priority,omitempty" path:"priority"`
	// Path on the remote server to treat as the root of this mount.
	RemotePath string `url:"remote_path,omitempty" json:"remote_path,omitempty" path:"remote_path"`
	// Number of consecutive successes before considering the backend healthy.
	Rise int64 `url:"rise,omitempty" json:"rise,omitempty" path:"rise"`
	// Path to the canary file used for health checks. Required.
	CanaryFilePath string `url:"canary_file_path" json:"canary_file_path" path:"canary_file_path"`
	// The mount ID of the Remote Server Mount that this backend is associated with. Required.
	RemoteServerMountId int64 `url:"remote_server_mount_id" json:"remote_server_mount_id" path:"remote_server_mount_id"`
	// The remote server that this backend is associated with. Required.
	RemoteServerId int64 `url:"remote_server_id" json:"remote_server_id" path:"remote_server_id"`
}

func (p *RemoteMountBackendCreateParams) decimalOverrides() []lib.DecimalOverride {
	return []lib.DecimalOverride{
		{Key: "min_free_cpu", Field: "MinFreeCpu", Float: &p.MinFreeCpu, Decimal: &p.MinFreeCpuDecimal},
		{Key: "min_free_mem", Field: "MinFreeMem", Float: &p.MinFreeMem, Decimal: &p.MinFreeMemDecimal},
	}
}

// MarshalJSON sends each set MinFreeCpuDecimal or MinFreeMemDecimal as a JSON string.
func (p RemoteMountBackendCreateParams) MarshalJSON() ([]byte, error) {
	type remoteMountBackendCreateParams RemoteMountBackendCreateParams
	return lib.MarshalDecimalOverrides(remoteMountBackendCreateParams(p), p.decimalOverrides()...)
}

// UnmarshalJSON decodes JSON strings for min_free_cpu and min_free_mem as exact decimal text.
func (p *RemoteMountBackendCreateParams) UnmarshalJSON(data []byte) error {
	type fields RemoteMountBackendCreateParams
	// Each override decodes its key in the same pass as the other fields, so
	// key matching and repeated keys work as in plain decoding.
	type remoteMountBackendCreateParams struct {
		*fields
		MinFreeCpu lib.DecimalOverride `json:"min_free_cpu"`
		MinFreeMem lib.DecimalOverride `json:"min_free_mem"`
	}
	decoded := *p
	overrides := decoded.decimalOverrides()
	target := remoteMountBackendCreateParams{fields: (*fields)(&decoded), MinFreeCpu: overrides[0], MinFreeMem: overrides[1]}
	// A failed decode leaves p unchanged.
	if err := json.Unmarshal(data, &target); err != nil {
		return err
	}
	*p = decoded
	return nil
}

// ToValues encodes p as query values, sending each set MinFreeCpuDecimal or MinFreeMemDecimal as its text.
func (p RemoteMountBackendCreateParams) ToValues() (url.Values, error) {
	type remoteMountBackendCreateParams RemoteMountBackendCreateParams
	return lib.DecimalOverrideValues(remoteMountBackendCreateParams(p), p.decimalOverrides()...)
}

// RemoteMountBackendResetStatusParams contains the request parameters for POST /remote_mount_backends/{id}/reset_status.
//
// Reset backend status to healthy
type RemoteMountBackendResetStatusParams struct {
	// Remote Mount Backend ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// RemoteMountBackendUpdateParams contains the request parameters for PATCH /remote_mount_backends/{id}.
type RemoteMountBackendUpdateParams struct {
	// Remote Mount Backend ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// True if this backend is enabled.
	Enabled *bool `url:"enabled,omitempty" json:"enabled,omitempty" path:"enabled"`
	// Number of consecutive failures before considering the backend unhealthy.
	Fall int64 `url:"fall,omitempty" json:"fall,omitempty" path:"fall"`
	// True if health checks are enabled for this backend.
	HealthCheckEnabled *bool `url:"health_check_enabled,omitempty" json:"health_check_enabled,omitempty" path:"health_check_enabled"`
	// Type of health check to perform. Values: active, passive.
	HealthCheckType RemoteMountBackendHealthCheckTypeEnum `url:"health_check_type,omitempty" json:"health_check_type,omitempty" path:"health_check_type"`
	// Interval in seconds between health checks.
	Interval int64 `url:"interval,omitempty" json:"interval,omitempty" path:"interval"`
	// Minimum free CPU percentage required for this backend to be considered healthy.
	MinFreeCpu float64 `url:"min_free_cpu,omitempty" json:"min_free_cpu,omitempty" path:"min_free_cpu"`
	// MinFreeCpuDecimal sends exact decimal text, such as "1.5", as min_free_cpu. Leave MinFreeCpu zero when you set it.
	MinFreeCpuDecimal *string `json:"-" url:"-" path:"-"`
	// Minimum free memory percentage required for this backend to be considered healthy.
	MinFreeMem float64 `url:"min_free_mem,omitempty" json:"min_free_mem,omitempty" path:"min_free_mem"`
	// MinFreeMemDecimal sends exact decimal text, such as "1.5", as min_free_mem. Leave MinFreeMem zero when you set it.
	MinFreeMemDecimal *string `json:"-" url:"-" path:"-"`
	// Priority of this backend.
	Priority int64 `url:"priority,omitempty" json:"priority,omitempty" path:"priority"`
	// Path on the remote server to treat as the root of this mount.
	RemotePath string `url:"remote_path,omitempty" json:"remote_path,omitempty" path:"remote_path"`
	// Number of consecutive successes before considering the backend healthy.
	Rise int64 `url:"rise,omitempty" json:"rise,omitempty" path:"rise"`
	// Path to the canary file used for health checks.
	CanaryFilePath string `url:"canary_file_path,omitempty" json:"canary_file_path,omitempty" path:"canary_file_path"`
	// The remote server that this backend is associated with.
	RemoteServerId int64 `url:"remote_server_id,omitempty" json:"remote_server_id,omitempty" path:"remote_server_id"`
}

func (p *RemoteMountBackendUpdateParams) decimalOverrides() []lib.DecimalOverride {
	return []lib.DecimalOverride{
		{Key: "min_free_cpu", Field: "MinFreeCpu", Float: &p.MinFreeCpu, Decimal: &p.MinFreeCpuDecimal},
		{Key: "min_free_mem", Field: "MinFreeMem", Float: &p.MinFreeMem, Decimal: &p.MinFreeMemDecimal},
	}
}

// MarshalJSON sends each set MinFreeCpuDecimal or MinFreeMemDecimal as a JSON string.
func (p RemoteMountBackendUpdateParams) MarshalJSON() ([]byte, error) {
	type remoteMountBackendUpdateParams RemoteMountBackendUpdateParams
	return lib.MarshalDecimalOverrides(remoteMountBackendUpdateParams(p), p.decimalOverrides()...)
}

// UnmarshalJSON decodes JSON strings for min_free_cpu and min_free_mem as exact decimal text.
func (p *RemoteMountBackendUpdateParams) UnmarshalJSON(data []byte) error {
	type fields RemoteMountBackendUpdateParams
	// Each override decodes its key in the same pass as the other fields, so
	// key matching and repeated keys work as in plain decoding.
	type remoteMountBackendUpdateParams struct {
		*fields
		MinFreeCpu lib.DecimalOverride `json:"min_free_cpu"`
		MinFreeMem lib.DecimalOverride `json:"min_free_mem"`
	}
	decoded := *p
	overrides := decoded.decimalOverrides()
	target := remoteMountBackendUpdateParams{fields: (*fields)(&decoded), MinFreeCpu: overrides[0], MinFreeMem: overrides[1]}
	// A failed decode leaves p unchanged.
	if err := json.Unmarshal(data, &target); err != nil {
		return err
	}
	*p = decoded
	return nil
}

// ToValues encodes p as query values, sending each set MinFreeCpuDecimal or MinFreeMemDecimal as its text.
func (p RemoteMountBackendUpdateParams) ToValues() (url.Values, error) {
	type remoteMountBackendUpdateParams RemoteMountBackendUpdateParams
	return lib.DecimalOverrideValues(remoteMountBackendUpdateParams(p), p.decimalOverrides()...)
}

// RemoteMountBackendDeleteParams contains the request parameters for DELETE /remote_mount_backends/{id}.
type RemoteMountBackendDeleteParams struct {
	// Remote Mount Backend ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (r *RemoteMountBackend) UnmarshalJSON(data []byte) error {
	type remoteMountBackend RemoteMountBackend
	var v remoteMountBackend
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*r = RemoteMountBackend(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (r *RemoteMountBackendCollection) UnmarshalJSON(data []byte) error {
	type remoteMountBackends RemoteMountBackendCollection
	var v remoteMountBackends
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*r = RemoteMountBackendCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (r *RemoteMountBackendCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*r))
	for i, v := range *r {
		ret[i] = v
	}

	return &ret
}
