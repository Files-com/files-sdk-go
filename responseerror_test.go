package files_sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestStr = `
{
  "data": {
    "u2f_sign_requests": [
      {
        "app_id": "https://dustin.files.com",
        "challenge": "XXXX",
        "sign_request": {
          "version": "U2F_V2",
          "keyHandle": "XXXX"
        }
      }
    ],
    "two_factor_authentication_methods": [
      "u2f"
    ],
    "u2f_redirect": "https://dustin.files.com",
    "partial_session_id": "XXX"
  },
  "error": "2FA Authenication error: Insert your U2F/FIDO key and press its button.",
  "http-code": 401,
  "instance": "XXX",
  "errors": [
    {
      "data": {
        "u2f_sign_requests": [
          {
            "app_id": "https://dustin.files.com",
            "challenge": "XXX",
            "sign_request": {
              "version": "U2F_V2",
              "keyHandle": "XXXX"
            }
          }
        ],
        "two_factor_authentication_methods": [
          "u2f"
        ],
        "u2f_redirect": "https://dustin.files.com",
        "partial_session_id": "XXX"
      },
      "error": "2FA Authenication error: Insert your U2F/FIDO key and press its button.",
      "http-code": 401,
      "instance": "XXX",
      "title": "Two Factor Authentication Error",
      "type": "401-two-factor-authentication-error"
    }
  ],
  "title": "Two Factor Authentication Error",
  "type": "not-authenticated/two-factor-authentication-error"
}
`

func TestResponseError_UnmarshalJSON(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	subject.UnmarshalJSON([]byte(TestStr))

	assert.Equal("2FA Authenication error: Insert your U2F/FIDO key and press its button.", subject.ErrorMessage)
	assert.Equal("not-authenticated/two-factor-authentication-error", subject.Type)
	assert.Equal("Two Factor Authentication Error", subject.Title)
	assert.Equal(false, subject.IsNil())
}

func TestResponseError_UnmarshalJSON_Error(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	subject.UnmarshalJSON([]byte("{"))

	assert.Equal("", subject.ErrorMessage)
	assert.Equal("", subject.Type)
	assert.Equal(true, subject.IsNil())
}
