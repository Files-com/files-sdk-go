package files_sdk

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallRawSetsConfiguredUserAgent(t *testing.T) {
	var gotUserAgent, gotAPIKey string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserAgent = r.UserAgent()
		gotAPIKey = r.Header.Get("X-FilesAPI-Key")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := Config{
		APIKey:    "secret",
		UserAgent: "Files.com Desktop Helper 1.2.3",
	}.Init()

	res, err := CallRaw(&CallParams{
		Method:  http.MethodGet,
		Config:  config,
		Uri:     server.URL,
		Headers: &http.Header{},
	})
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, "Files.com Desktop Helper 1.2.3", gotUserAgent)
	assert.Empty(t, gotAPIKey)
}

func TestCallRawPreservesExplicitUserAgentHeader(t *testing.T) {
	var gotUserAgent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserAgent = r.UserAgent()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := Config{UserAgent: "Files.com Go SDK test"}.Init()
	headers := http.Header{}
	headers.Set("User-Agent", "Custom Transfer Client")

	res, err := CallRaw(&CallParams{
		Method:  http.MethodGet,
		Config:  config,
		Uri:     server.URL,
		Headers: &headers,
	})
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, "Custom Transfer Client", gotUserAgent)
}
