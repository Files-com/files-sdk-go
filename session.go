package files_sdk

import (
	"encoding/json"
)

type Session struct {
	Id                  string `json:"id,omitempty"`
	Language            string `json:"language,omitempty"`
	ReadOnly            *bool  `json:"read_only,omitempty"`
	SftpInsecureCiphers *bool  `json:"sftp_insecure_ciphers,omitempty"`
	Username            string `json:"username,omitempty"`
	Password            string `json:"password,omitempty"`
	Otp                 string `json:"otp,omitempty"`
	PartialSessionId    string `json:"partial_session_id,omitempty"`
}

type SessionCollection []Session

type SessionCreateParams struct {
	Username         string `url:"username,omitempty" required:"false" json:"username,omitempty"`
	Password         string `url:"password,omitempty" required:"false" json:"password,omitempty"`
	Otp              string `url:"otp,omitempty" required:"false" json:"otp,omitempty"`
	PartialSessionId string `url:"partial_session_id,omitempty" required:"false" json:"partial_session_id,omitempty"`
}

func (s *Session) UnmarshalJSON(data []byte) error {
	type session Session
	var v session
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = Session(v)
	return nil
}

func (s *SessionCollection) UnmarshalJSON(data []byte) error {
	type sessions []Session
	var v sessions
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*s = SessionCollection(v)
	return nil
}

func (s *SessionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*s))
	for i, v := range *s {
		ret[i] = v
	}

	return &ret
}
