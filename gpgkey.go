package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// GpgKey is a Files.com API resource.
type GpgKey struct {
	// GPG key ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID (0 for default workspace).
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// GPG key expiration date.
	ExpiresAt *time.Time `json:"expires_at,omitempty" path:"expires_at,omitempty" url:"expires_at,omitempty"`
	// GPG key name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Partner ID who owns this GPG Key, if applicable.
	PartnerId int64 `json:"partner_id,omitempty" path:"partner_id,omitempty" url:"partner_id,omitempty"`
	// Name of the Partner who owns this GPG Key, if applicable.
	PartnerName string `json:"partner_name,omitempty" path:"partner_name,omitempty" url:"partner_name,omitempty"`
	// User ID who owns this GPG Key, if applicable.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// MD5 hash of the GPG public key.
	PublicKeyMd5 string `json:"public_key_md5,omitempty" path:"public_key_md5,omitempty" url:"public_key_md5,omitempty"`
	// MD5 hash of the GPG private key.
	PrivateKeyMd5 string `json:"private_key_md5,omitempty" path:"private_key_md5,omitempty" url:"private_key_md5,omitempty"`
	// GPG public key.
	GeneratedPublicKey string `json:"generated_public_key,omitempty" path:"generated_public_key,omitempty" url:"generated_public_key,omitempty"`
	// GPG private key.
	GeneratedPrivateKey string `json:"generated_private_key,omitempty" path:"generated_private_key,omitempty" url:"generated_private_key,omitempty"`
	// GPG private key password. Only required for password protected keys.
	PrivateKeyPasswordMd5 string `json:"private_key_password_md5,omitempty" path:"private_key_password_md5,omitempty" url:"private_key_password_md5,omitempty"`
	// The GPG public key.
	PublicKey string `json:"public_key,omitempty" path:"public_key,omitempty" url:"public_key,omitempty"`
	// The GPG private key.
	PrivateKey string `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
	// The GPG private key password.
	PrivateKeyPassword string `json:"private_key_password,omitempty" path:"private_key_password,omitempty" url:"private_key_password,omitempty"`
	// Expiration date of the key. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
	GenerateExpiresAt string `json:"generate_expires_at,omitempty" path:"generate_expires_at,omitempty" url:"generate_expires_at,omitempty"`
	// If true, generate a new GPG key pair. Can not be used with `public_key`/`private_key`.
	GenerateKeypair *bool `json:"generate_keypair,omitempty" path:"generate_keypair,omitempty" url:"generate_keypair,omitempty"`
	// Full name of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
	GenerateFullName string `json:"generate_full_name,omitempty" path:"generate_full_name,omitempty" url:"generate_full_name,omitempty"`
	// Email address of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
	GenerateEmail string `json:"generate_email,omitempty" path:"generate_email,omitempty" url:"generate_email,omitempty"`
}

// Identifier returns the resource ID.
func (g GpgKey) Identifier() interface{} {
	return g.Id
}

// GpgKeyCollection is a list of GpgKey resources.
type GpgKeyCollection []GpgKey

// GpgKeyListParams contains the request parameters for GET /gpg_keys.
type GpgKeyListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name` or `expires_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`, `partner_id` or `expires_at`. Valid field combinations are `[ workspace_id, partner_id ]`, `[ workspace_id, expires_at ]`, `[ partner_id, expires_at ]` or `[ workspace_id, partner_id, expires_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `expires_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `expires_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `expires_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `expires_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// GpgKeyFindParams contains the request parameters for GET /gpg_keys/{id}.
type GpgKeyFindParams struct {
	// Gpg Key ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// GpgKeyCreateParams contains the request parameters for POST /gpg_keys.
type GpgKeyCreateParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// Partner ID who owns this GPG Key, if applicable.
	PartnerId int64 `url:"partner_id,omitempty" json:"partner_id,omitempty" path:"partner_id"`
	// The GPG public key.
	PublicKey string `url:"public_key,omitempty" json:"public_key,omitempty" path:"public_key"`
	// The GPG private key.
	PrivateKey string `url:"private_key,omitempty" json:"private_key,omitempty" path:"private_key"`
	// The GPG private key password.
	PrivateKeyPassword string `url:"private_key_password,omitempty" json:"private_key_password,omitempty" path:"private_key_password"`
	// GPG key name. Required.
	Name string `url:"name" json:"name" path:"name"`
	// Workspace ID (0 for default workspace).
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Expiration date of the key. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
	GenerateExpiresAt *time.Time `url:"generate_expires_at,omitempty" json:"generate_expires_at,omitempty" path:"generate_expires_at"`
	// If true, generate a new GPG key pair. Can not be used with `public_key`/`private_key`.
	GenerateKeypair *bool `url:"generate_keypair,omitempty" json:"generate_keypair,omitempty" path:"generate_keypair"`
	// Full name of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
	GenerateFullName string `url:"generate_full_name,omitempty" json:"generate_full_name,omitempty" path:"generate_full_name"`
	// Email address of the key owner. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
	GenerateEmail string `url:"generate_email,omitempty" json:"generate_email,omitempty" path:"generate_email"`
}

// GpgKeyUpdateParams contains the request parameters for PATCH /gpg_keys/{id}.
type GpgKeyUpdateParams struct {
	// Gpg Key ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Partner ID who owns this GPG Key, if applicable.
	PartnerId int64 `url:"partner_id,omitempty" json:"partner_id,omitempty" path:"partner_id"`
	// The GPG public key.
	PublicKey string `url:"public_key,omitempty" json:"public_key,omitempty" path:"public_key"`
	// The GPG private key.
	PrivateKey string `url:"private_key,omitempty" json:"private_key,omitempty" path:"private_key"`
	// The GPG private key password.
	PrivateKeyPassword string `url:"private_key_password,omitempty" json:"private_key_password,omitempty" path:"private_key_password"`
	// GPG key name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
}

// GpgKeyDeleteParams contains the request parameters for DELETE /gpg_keys/{id}.
type GpgKeyDeleteParams struct {
	// Gpg Key ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (g *GpgKey) UnmarshalJSON(data []byte) error {
	type gpgKey GpgKey
	var v gpgKey
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*g = GpgKey(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (g *GpgKeyCollection) UnmarshalJSON(data []byte) error {
	type gpgKeys GpgKeyCollection
	var v gpgKeys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*g = GpgKeyCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (g *GpgKeyCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*g))
	for i, v := range *g {
		ret[i] = v
	}

	return &ret
}
