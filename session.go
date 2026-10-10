package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Session is a Files.com API resource.
type Session struct {
	Id                  string `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	Language            string `json:"language,omitempty" path:"language,omitempty" url:"language,omitempty"`
	ReadOnly            *bool  `json:"read_only,omitempty" path:"read_only,omitempty" url:"read_only,omitempty"`
	SftpInsecureCiphers *bool  `json:"sftp_insecure_ciphers,omitempty" path:"sftp_insecure_ciphers,omitempty" url:"sftp_insecure_ciphers,omitempty"`
	Username            string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	Password            string `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
	Otp                 string `json:"otp,omitempty" path:"otp,omitempty" url:"otp,omitempty"`
	PartialSessionId    string `json:"partial_session_id,omitempty" path:"partial_session_id,omitempty" url:"partial_session_id,omitempty"`
}

// Identifier returns the resource ID.
func (s Session) Identifier() interface{} {
	return s.Id
}

// SessionCollection is a list of Session resources.
type SessionCollection []Session

// SessionCreateParams contains the request parameters for POST /sessions.
type SessionCreateParams struct {
	Username         string `url:"username,omitempty" json:"username,omitempty" path:"username"`
	Password         string `url:"password,omitempty" json:"password,omitempty" path:"password"`
	Otp              string `url:"otp,omitempty" json:"otp,omitempty" path:"otp"`
	PartialSessionId string `url:"partial_session_id,omitempty" json:"partial_session_id,omitempty" path:"partial_session_id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (s *Session) UnmarshalJSON(data []byte) error {
	type session Session
	var v session
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*s = Session(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (s *SessionCollection) UnmarshalJSON(data []byte) error {
	type sessions SessionCollection
	var v sessions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*s = SessionCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (s *SessionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
