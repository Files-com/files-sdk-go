package files_sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestStr1 = `
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

var TestStr2 = `
{
  "error": "Hidden reason can't be blank",
  "http-code": 422,
  "instance": "9a0165ca-cbe4-480d-b2f3-376c5a3c5ff6",
  "model_errors": {
    "hidden_reason": [
      "Hidden reason can't be blank"
    ]
  },
  "model_error_keys": {
    "hidden_reason": [
      "blank"
    ]
  },
  "errors": [
    "Hidden reason can't be blank2"
  ],
  "title": "Model Save Error",
  "type": "processing-failure/model-save-error"
}
`

var TestStr3 = `
{"error":"Internal server error, please contact support or the person who created your account.","http-code":"500"}
`

func TestResponseError1_UnmarshalJSON(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte(TestStr1))

	assert.Nil(err)
	assert.Equal("2FA Authenication error: Insert your U2F/FIDO key and press its button.", subject.ErrorMessage)
	assert.Equal("not-authenticated/two-factor-authentication-error", subject.Type)
	assert.Equal("Two Factor Authentication Error", subject.Title)
	assert.Equal("Two Factor Authentication Error", subject.Errors[0].Title)
	assert.Equal(false, subject.IsNil())
}

func TestResponseError2_UnmarshalJSON(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte(TestStr2))

	assert.Nil(err)
	assert.Equal("Hidden reason can't be blank", subject.ErrorMessage)
	assert.Equal("processing-failure/model-save-error", subject.Type)
	assert.Equal("Model Save Error", subject.Title)
	assert.Equal("Hidden reason can't be blank2", subject.Errors[0].ErrorMessage)
	assert.Equal(false, subject.IsNil())
}

func TestResponseError3_UnmarshalJSON(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte(TestStr3))

	assert.Nil(err)
	assert.Equal("Internal server error, please contact support or the person who created your account.", subject.ErrorMessage)
	assert.Equal(int(500), subject.HttpCode)
	assert.Equal("", subject.Type)
	assert.Equal(false, subject.IsNil())
}

func TestResponseError_UnmarshalJSON_Error(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte("{"))

	assert.Equal(err.Error(), "unexpected end of JSON input")
	assert.Equal("", subject.ErrorMessage)
	assert.Equal("", subject.Type)
	assert.Equal(true, subject.IsNil())
}

func TestResponseError_UnmarshalJSON_Error2(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte(`{"error": ["error"]}`))

	assert.Error(err)
	assert.Equal("", subject.ErrorMessage)
	assert.Equal("", subject.Type)
	assert.Equal(true, subject.IsNil())
}

func TestResponseError_UnmarshalJSON_Error3(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte(`["error"]`))

	assert.Nil(err, "The response is not an error, but a list response.")
	assert.Equal("", subject.ErrorMessage)
	assert.Equal("", subject.Type)
	assert.Equal(true, subject.IsNil())
}
