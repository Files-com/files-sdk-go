package files_sdk

import (
	"encoding/json"
	goerrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
    "partial_session_id": "XXX",
	"unknown-key": "unknown-value"
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

var TestStr4 = `
<body></body>
`

var notAuthenticated = `
{
  "error": "The API key or Session token provided could not be used to validate this request.",
  "http-code":401,
  "instance":"5baab0b9dd8b58ffa436cf86498bda05",
  "title":"Invalid Credentials",
  "type":"not-authenticated/not-a-real-auth-type"
}
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
	assert.Equal("unknown-value", subject.RawData["unknown-key"])
	assert.Equal("XXX", subject.Data.PartialSessionId)
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
	assert.False(subject.IsNil())
}

func TestResponseError3_UnmarshalJSON(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte(TestStr3))

	assert.Nil(err)
	assert.Equal("Internal server error, please contact support or the person who created your account.", subject.ErrorMessage)
	assert.Equal(int(500), subject.HttpCode)
	assert.Equal("", subject.Type)
	assert.False(subject.IsNil())
}

func TestResponseError_UnmarshalJSON_Error(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte("{"))

	assert.Equal(err.Error(), "unexpected end of JSON input")
	assert.Equal("", subject.ErrorMessage)
	assert.Equal("", subject.Type)
	assert.True(subject.IsNil(), "Empty ErrorMessage should make IsNil() true")
}

func TestResponseError_UnmarshalJSON_Error2(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte(`{"error": ["error"]}`))

	assert.Error(err)
	assert.Equal("", subject.ErrorMessage)
	assert.Equal("", subject.Type)
	assert.True(subject.IsNil(), "Empty ErrorMessage should make IsNil() true")
}

func TestResponseError_UnmarshalJSON_Error3(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte(`["error"]`))

	assert.Nil(err, "The response is not an error, but a list response.")
	assert.Equal("", subject.ErrorMessage)
	assert.Equal("", subject.Type)
	assert.True(subject.IsNil(), "Empty ErrorMessage should make IsNil() true")
}

func TestResponseError_UnmarshalJSON_Error4(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte(TestStr4))

	assert.Error(err, "\n<body></body>\n")
	assert.Equal("", subject.ErrorMessage)
	assert.Equal("", subject.Type)
	assert.True(subject.IsNil(), "Empty ErrorMessage should make IsNil() true")
}

func TestResponseError_MarshalJSON(t *testing.T) {
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte(TestStr1))
	require.NoError(t, err)
	jsonBytes, err := json.Marshal(subject)
	require.NoError(t, err)
	assert.JSONEq(t, TestStr1, string(jsonBytes))
}

func TestIsNotAuthenticated(t *testing.T) {
	assert := assert.New(t)
	subject := ResponseError{}

	err := subject.UnmarshalJSON([]byte(notAuthenticated))
	require.NoError(t, err)
	assert.True(IsNotAuthenticated(subject))
}

func TestResponseErrorTypeConstants(t *testing.T) {
	assert.Equal(t, ResponseErrorType("processing-failure/destination-exists"), ErrDestinationExists)
	assert.Equal(t, ResponseErrorType("not-authenticated/invalid-credentials"), ErrInvalidCredentials)
	assert.Equal(t, ResponseErrorType("not-authenticated/invalid-username-or-password"), ErrInvalidUsernameOrPassword)
	assert.Equal(t, ResponseErrorType("download_request_expired"), ErrDownloadRequestExpired)
	assert.Equal(t, ResponseErrorType("upload_request_expired"), ErrUploadRequestExpired)
	assert.Equal(t, ResponseErrorGroup("not-found"), ErrNotFound)
	assert.Equal(t, ResponseErrorGroup("not-authenticated"), ErrNotAuthenticated)
	assert.Equal(t, "processing-failure/destination-exists", DestinationExists)
	assert.Equal(t, "download_request_expired", DownloadRequestExpired)
	assert.Equal(t, "upload_request_expired", UploadRequestExpired)
}

func TestResponseErrorTypeOf(t *testing.T) {
	errType, ok := ResponseErrorTypeOf(ResponseError{Type: string(ErrDestinationExists)})
	assert.True(t, ok)
	assert.Equal(t, ErrDestinationExists, errType)

	errType, ok = ResponseErrorTypeOf(&ResponseError{Type: string(ErrFileNotFound)})
	assert.True(t, ok)
	assert.Equal(t, ErrFileNotFound, errType)

	_, ok = ResponseErrorTypeOf(ResponseError{})
	assert.False(t, ok)

	_, ok = ResponseErrorTypeOf(goerrors.New("boom"))
	assert.False(t, ok)
}

func TestResponseErrorStdlibMatching(t *testing.T) {
	err := ResponseError{Type: string(ErrFileNotFound)}

	assert.True(t, goerrors.Is(err, ErrFileNotFound))
	assert.True(t, goerrors.Is(err, ErrNotFound))
	assert.False(t, goerrors.Is(err, ErrNotAuthenticated))

	var responseErr ResponseError
	assert.True(t, goerrors.As(err, &responseErr))
	assert.Equal(t, ErrFileNotFound, ResponseErrorType(responseErr.Type))

	errPtr := &ResponseError{Type: string(ErrDestinationExists)}
	assert.True(t, goerrors.Is(errPtr, ErrDestinationExists))

	var responseErrPtr *ResponseError
	assert.True(t, goerrors.As(errPtr, &responseErrPtr))
	assert.Equal(t, ErrDestinationExists, ResponseErrorType(responseErrPtr.Type))
}

func TestIsErrorType(t *testing.T) {
	err := ResponseError{Type: string(ErrFileNotFound)}
	assert.True(t, IsErrorType(err, ErrFileNotFound))
	assert.False(t, IsErrorType(err, ErrDestinationExists))
}

func TestIsAnyErrorType(t *testing.T) {
	err := ResponseError{Type: string(ErrFileNotFound)}
	assert.True(t, IsAnyErrorType(err, ErrFileNotFound, ErrDestinationExists))
	assert.False(t, IsAnyErrorType(err, ErrDestinationExists))
}

func TestIsErrorGroup(t *testing.T) {
	err := ResponseError{Type: string(ErrFileNotFound)}
	assert.True(t, IsErrorGroup(err, ErrNotFound))
	assert.False(t, IsErrorGroup(err, ErrNotAuthenticated))
	assert.False(t, IsErrorGroup(goerrors.New("boom"), ErrNotFound))
}

func TestGroupHelpers(t *testing.T) {
	assert.True(t, IsNotFound(ResponseError{Type: string(ErrFileNotFound)}))
	assert.True(t, IsAuthenticationError(ResponseError{Type: string(ErrInvalidCredentials)}))
	assert.False(t, IsAuthenticationError(ResponseError{Type: string(ErrFileNotFound)}))
}

func TestLegacyWrapperHelpers(t *testing.T) {
	assert.True(t, IsExist(ResponseError{Type: string(ErrDestinationExists)}))
	assert.False(t, IsExist(ResponseError{Type: string(ErrFileNotFound)}))

	assert.True(t, IsNotExist(ResponseError{Type: string(ErrFileNotFound)}))
	assert.False(t, IsNotExist(ResponseError{Type: string(ErrDestinationExists)}))

	assert.True(t, IsNotAuthenticated(ResponseError{Type: string(ErrInvalidCredentials)}))
	assert.False(t, IsNotAuthenticated(ResponseError{Type: string(ErrFileNotFound)}))

	assert.True(t, IsExpired(ResponseError{Type: string(ErrDownloadRequestExpired)}))
	assert.True(t, IsExpired(ResponseError{Type: string(ErrUploadRequestExpired)}))
	assert.False(t, IsExpired(ResponseError{Type: string(ErrFileNotFound)}))
}
