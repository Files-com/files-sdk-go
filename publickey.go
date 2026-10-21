package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// PublicKey is a Files.com API resource.
type PublicKey struct {
	// Public key ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Workspace ID (0 for default workspace).
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Public key title.
	Title string `json:"title,omitempty" path:"title,omitempty" url:"title,omitempty"`
	// Public key created at date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Public key expiration date/time.
	ExpiresAt *time.Time `json:"expires_at,omitempty" path:"expires_at,omitempty" url:"expires_at,omitempty"`
	// Is this public key expired?
	Expired *bool `json:"expired,omitempty" path:"expired,omitempty" url:"expired,omitempty"`
	// Public key fingerprint (MD5).
	Fingerprint string `json:"fingerprint,omitempty" path:"fingerprint,omitempty" url:"fingerprint,omitempty"`
	// Public key fingerprint (SHA256).
	FingerprintSha256 string `json:"fingerprint_sha256,omitempty" path:"fingerprint_sha256,omitempty" url:"fingerprint_sha256,omitempty"`
	// Only returned when generating keys. Can be invalid, not_generated, generating, complete.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Key's most recent login time via SFTP.
	LastLoginAt *time.Time `json:"last_login_at,omitempty" path:"last_login_at,omitempty" url:"last_login_at,omitempty"`
	// Only returned when generating keys. Private key generated for the user.
	GeneratedPrivateKey string `json:"generated_private_key,omitempty" path:"generated_private_key,omitempty" url:"generated_private_key,omitempty"`
	// Only returned when generating keys. Public key generated for the user.
	GeneratedPublicKey string `json:"generated_public_key,omitempty" path:"generated_public_key,omitempty" url:"generated_public_key,omitempty"`
	// Username of the user this public key is associated with.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// User ID this public key is associated with.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Actual contents of SSH key.
	PublicKey string `json:"public_key,omitempty" path:"public_key,omitempty" url:"public_key,omitempty"`
	// If true, generate a new SSH key pair. Can not be used with `public_key`.
	GenerateKeypair *bool `json:"generate_keypair,omitempty" path:"generate_keypair,omitempty" url:"generate_keypair,omitempty"`
	// Password for the private key. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
	GeneratePrivateKeyPassword string `json:"generate_private_key_password,omitempty" path:"generate_private_key_password,omitempty" url:"generate_private_key_password,omitempty"`
	// Type of key to generate. One of rsa, dsa, ecdsa, ed25519. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
	GenerateAlgorithm string `json:"generate_algorithm,omitempty" path:"generate_algorithm,omitempty" url:"generate_algorithm,omitempty"`
	// Length of key to generate. If algorithm is ecdsa, this is the signature size. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
	GenerateLength int64 `json:"generate_length,omitempty" path:"generate_length,omitempty" url:"generate_length,omitempty"`
}

// Identifier returns the resource ID.
func (p PublicKey) Identifier() interface{} {
	return p.Id
}

// PublicKeyCollection is a list of PublicKey resources.
type PublicKeyCollection []PublicKey

// PublicKeyListParams contains the request parameters for GET /public_keys.
type PublicKeyListParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `user_id`, `title` or `created_at`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `created_at` and `workspace_id`. Valid field combinations are `[ workspace_id, created_at ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is greater than the supplied value. Valid fields are `created_at`.
	FilterGt interface{} `url:"filter_gt,omitempty" json:"filter_gt,omitempty" path:"filter_gt"`
	// If set, return records where the specified field is greater than or equal the supplied value. Valid fields are `created_at`.
	FilterGteq interface{} `url:"filter_gteq,omitempty" json:"filter_gteq,omitempty" path:"filter_gteq"`
	// If set, return records where the specified field is less than the supplied value. Valid fields are `created_at`.
	FilterLt interface{} `url:"filter_lt,omitempty" json:"filter_lt,omitempty" path:"filter_lt"`
	// If set, return records where the specified field is less than or equal the supplied value. Valid fields are `created_at`.
	FilterLteq interface{} `url:"filter_lteq,omitempty" json:"filter_lteq,omitempty" path:"filter_lteq"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// PublicKeyFindParams contains the request parameters for GET /public_keys/{id}.
type PublicKeyFindParams struct {
	// Public Key ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// PublicKeyCreateParams contains the request parameters for POST /public_keys.
type PublicKeyCreateParams struct {
	// User ID. Provide a value of `0` to operate the current session's user.
	UserId int64 `url:"user_id,omitempty" json:"user_id,omitempty" path:"user_id"`
	// Internal reference for key. Required.
	Title string `url:"title" json:"title" path:"title"`
	// Actual contents of SSH key.
	PublicKey string `url:"public_key,omitempty" json:"public_key,omitempty" path:"public_key"`
	// If true, generate a new SSH key pair. Can not be used with `public_key`.
	GenerateKeypair *bool `url:"generate_keypair,omitempty" json:"generate_keypair,omitempty" path:"generate_keypair"`
	// Password for the private key. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
	GeneratePrivateKeyPassword string `url:"generate_private_key_password,omitempty" json:"generate_private_key_password,omitempty" path:"generate_private_key_password"`
	// Type of key to generate. One of rsa, dsa, ecdsa, ed25519. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
	GenerateAlgorithm string `url:"generate_algorithm,omitempty" json:"generate_algorithm,omitempty" path:"generate_algorithm"`
	// Length of key to generate. If algorithm is ecdsa, this is the signature size. Used for the generation of the key. Will be ignored if `generate_keypair` is false.
	GenerateLength int64 `url:"generate_length,omitempty" json:"generate_length,omitempty" path:"generate_length"`
}

// PublicKeyUpdateParams contains the request parameters for PATCH /public_keys/{id}.
type PublicKeyUpdateParams struct {
	// Public Key ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// Internal reference for key. Required.
	Title string `url:"title" json:"title" path:"title"`
}

// PublicKeyDeleteParams contains the request parameters for DELETE /public_keys/{id}.
type PublicKeyDeleteParams struct {
	// Public Key ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (p *PublicKey) UnmarshalJSON(data []byte) error {
	type publicKey PublicKey
	var v publicKey
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*p = PublicKey(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (p *PublicKeyCollection) UnmarshalJSON(data []byte) error {
	type publicKeys PublicKeyCollection
	var v publicKeys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*p = PublicKeyCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (p *PublicKeyCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*p))
	for i, v := range *p {
		ret[i] = v
	}

	return &ret
}
