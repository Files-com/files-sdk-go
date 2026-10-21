package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Session is a Files.com API resource.
type Session struct {
	// Session ID.
	Id string `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Session language.
	Language string `json:"language,omitempty" path:"language,omitempty" url:"language,omitempty"`
	// Is this session read only?
	ReadOnly *bool `json:"read_only,omitempty" path:"read_only,omitempty" url:"read_only,omitempty"`
	// Are insecure SFTP ciphers allowed for this user? (If this is set to true, the site administrator has signaled that it is ok to use less secure SSH ciphers for this user.).
	SftpInsecureCiphers *bool `json:"sftp_insecure_ciphers,omitempty" path:"sftp_insecure_ciphers,omitempty" url:"sftp_insecure_ciphers,omitempty"`
	// Username to sign in as.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// Password for sign in.
	Password string `json:"password,omitempty" path:"password,omitempty" url:"password,omitempty"`
	// If this user has a 2FA device, provide its OTP or code here.
	Otp string `json:"otp,omitempty" path:"otp,omitempty" url:"otp,omitempty"`
	// Identifier for a partially-completed login.
	PartialSessionId string `json:"partial_session_id,omitempty" path:"partial_session_id,omitempty" url:"partial_session_id,omitempty"`
}

// Identifier returns the resource ID.
func (s Session) Identifier() interface{} {
	return s.Id
}

// SessionCollection is a list of Session resources.
type SessionCollection []Session

// SessionCreateParams contains the request parameters for POST /sessions.
type SessionCreateParams struct {
	// Username to sign in as.
	Username string `url:"username,omitempty" json:"username,omitempty" path:"username"`
	// Password for sign in.
	Password string `url:"password,omitempty" json:"password,omitempty" path:"password"`
	// If this user has a 2FA device, provide its OTP or code here.
	Otp string `url:"otp,omitempty" json:"otp,omitempty" path:"otp"`
	// Identifier for a partially-completed login.
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
