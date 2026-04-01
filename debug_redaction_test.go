package files_sdk

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeDebugRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "https://app.files.com/api/rest/v1/folders", nil)
	require.NoError(t, err)

	req.Header.Set("X-FilesAPI-Key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	req.Header.Set("X-FilesAPI-Auth", "session-token")
	req.Header.Set("X-Files-Reauthentication", "password")

	sanitized := sanitizeDebugRequest(req)

	require.Equal(t, "0123****************", sanitized.Header.Get("X-FilesAPI-Key"))
	require.Equal(t, "<redacted>", sanitized.Header.Get("X-FilesAPI-Auth"))
	require.Equal(t, "<redacted>", sanitized.Header.Get("X-Files-Reauthentication"))
	require.Equal(t, "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", req.Header.Get("X-FilesAPI-Key"))
	require.Equal(t, "session-token", req.Header.Get("X-FilesAPI-Auth"))
}
