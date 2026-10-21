package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// As2Station is a Files.com API resource.
type As2Station struct {
	// Id of the AS2 Station.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// ID of the Workspace associated with this AS2 Station.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// The station's formal AS2 name.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Public URI for sending AS2 message to.
	Uri string `json:"uri,omitempty" path:"uri,omitempty" url:"uri,omitempty"`
	// The station's AS2 domain name.
	Domain string `json:"domain,omitempty" path:"domain,omitempty" url:"domain,omitempty"`
	// Serial of public certificate used for message security in hex format.
	HexPublicCertificateSerial string `json:"hex_public_certificate_serial,omitempty" path:"hex_public_certificate_serial,omitempty" url:"hex_public_certificate_serial,omitempty"`
	// MD5 hash of public certificate used for message security.
	PublicCertificateMd5 string `json:"public_certificate_md5,omitempty" path:"public_certificate_md5,omitempty" url:"public_certificate_md5,omitempty"`
	// Public certificate used for message security.
	PublicCertificate string `json:"public_certificate,omitempty" path:"public_certificate,omitempty" url:"public_certificate,omitempty"`
	// MD5 hash of private key used for message security.
	PrivateKeyMd5 string `json:"private_key_md5,omitempty" path:"private_key_md5,omitempty" url:"private_key_md5,omitempty"`
	// Subject of public certificate used for message security.
	PublicCertificateSubject string `json:"public_certificate_subject,omitempty" path:"public_certificate_subject,omitempty" url:"public_certificate_subject,omitempty"`
	// Issuer of public certificate used for message security.
	PublicCertificateIssuer string `json:"public_certificate_issuer,omitempty" path:"public_certificate_issuer,omitempty" url:"public_certificate_issuer,omitempty"`
	// Serial of public certificate used for message security.
	PublicCertificateSerial string `json:"public_certificate_serial,omitempty" path:"public_certificate_serial,omitempty" url:"public_certificate_serial,omitempty"`
	// Not before value of public certificate used for message security.
	PublicCertificateNotBefore string `json:"public_certificate_not_before,omitempty" path:"public_certificate_not_before,omitempty" url:"public_certificate_not_before,omitempty"`
	// Not after value of public certificate used for message security.
	PublicCertificateNotAfter string `json:"public_certificate_not_after,omitempty" path:"public_certificate_not_after,omitempty" url:"public_certificate_not_after,omitempty"`
	// MD5 hash of private key password used for message security.
	PrivateKeyPasswordMd5 string `json:"private_key_password_md5,omitempty" path:"private_key_password_md5,omitempty" url:"private_key_password_md5,omitempty"`
	// PEM-encoded private key matching public_certificate.
	PrivateKey string `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
	// Password for the PEM-encoded private key.
	PrivateKeyPassword string `json:"private_key_password,omitempty" path:"private_key_password,omitempty" url:"private_key_password,omitempty"`
	// Base64-encoded PKCS#12 (.pfx or .p12) file containing the identity's certificate and private key. Provide this instead of public_certificate and private_key.
	Pkcs12 string `json:"pkcs12,omitempty" path:"pkcs12,omitempty" url:"pkcs12,omitempty"`
	// Password for pkcs12. The file and password are used only for import; the extracted certificate and private key are stored as PEM.
	Pkcs12Password string `json:"pkcs12_password,omitempty" path:"pkcs12_password,omitempty" url:"pkcs12_password,omitempty"`
}

// Identifier returns the resource ID.
func (a As2Station) Identifier() interface{} {
	return a.Id
}

// As2StationCollection is a list of As2Station resources.
type As2StationCollection []As2Station

// As2StationListParams contains the request parameters for GET /as2_stations.
type As2StationListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `name`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// As2StationFindParams contains the request parameters for GET /as2_stations/{id}.
type As2StationFindParams struct {
	// As2 Station ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// As2StationCreateParams contains the request parameters for POST /as2_stations.
type As2StationCreateParams struct {
	// The station's formal AS2 name. Required.
	Name string `url:"name" json:"name" path:"name"`
	// ID of the Workspace associated with this AS2 Station.
	WorkspaceId int64 `url:"workspace_id,omitempty" json:"workspace_id,omitempty" path:"workspace_id"`
	// Public certificate used for message security.
	PublicCertificate string `url:"public_certificate,omitempty" json:"public_certificate,omitempty" path:"public_certificate"`
	// PEM-encoded private key matching public_certificate.
	PrivateKey string `url:"private_key,omitempty" json:"private_key,omitempty" path:"private_key"`
	// Password for the PEM-encoded private key.
	PrivateKeyPassword string `url:"private_key_password,omitempty" json:"private_key_password,omitempty" path:"private_key_password"`
	// Base64-encoded PKCS#12 (.pfx or .p12) file containing the identity's certificate and private key. Provide this instead of public_certificate and private_key.
	Pkcs12 string `url:"pkcs12,omitempty" json:"pkcs12,omitempty" path:"pkcs12"`
	// Password for pkcs12. The file and password are used only for import; the extracted certificate and private key are stored as PEM.
	Pkcs12Password string `url:"pkcs12_password,omitempty" json:"pkcs12_password,omitempty" path:"pkcs12_password"`
}

// As2StationUpdateParams contains the request parameters for PATCH /as2_stations/{id}.
type As2StationUpdateParams struct {
	// As2 Station ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
	// The station's formal AS2 name.
	Name string `url:"name,omitempty" json:"name,omitempty" path:"name"`
	// Public certificate used for message security.
	PublicCertificate string `url:"public_certificate,omitempty" json:"public_certificate,omitempty" path:"public_certificate"`
	// PEM-encoded private key matching public_certificate.
	PrivateKey string `url:"private_key,omitempty" json:"private_key,omitempty" path:"private_key"`
	// Password for the PEM-encoded private key.
	PrivateKeyPassword string `url:"private_key_password,omitempty" json:"private_key_password,omitempty" path:"private_key_password"`
	// Base64-encoded PKCS#12 (.pfx or .p12) file containing the identity's replacement certificate and private key. Provide this instead of public_certificate and private_key.
	Pkcs12 string `url:"pkcs12,omitempty" json:"pkcs12,omitempty" path:"pkcs12"`
	// Password for pkcs12. The file and password are used only for import; the extracted certificate and private key are stored as PEM.
	Pkcs12Password string `url:"pkcs12_password,omitempty" json:"pkcs12_password,omitempty" path:"pkcs12_password"`
}

// As2StationDeleteParams contains the request parameters for DELETE /as2_stations/{id}.
type As2StationDeleteParams struct {
	// As2 Station ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *As2Station) UnmarshalJSON(data []byte) error {
	type as2Station As2Station
	var v as2Station
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = As2Station(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *As2StationCollection) UnmarshalJSON(data []byte) error {
	type as2Stations As2StationCollection
	var v as2Stations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = As2StationCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *As2StationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
