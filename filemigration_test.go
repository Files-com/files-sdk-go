package files_sdk

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrationLogIteratorReleasesResponse(t *testing.T) {
	for _, tc := range []struct {
		name      string
		payload   string
		earlyStop bool
		malformed bool
	}{
		{name: "exhausted", payload: `{"type":"file"}`},
		{name: "malformed", payload: `{"type":`, malformed: true},
		{name: "early stop", payload: `{"type":"file"}` + "\n" + `{"path":"` + strings.Repeat("x", 4096) + `"}`, earlyStop: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &optionResponseBody{Reader: strings.NewReader(tc.payload)}
			config := Config{}.Init().SetCustomClient(&http.Client{Transport: responseOptionTransport(func(req *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusOK, Body: body, Header: make(http.Header), Request: req}, nil
			})})
			it := (FilesMigrationLogIter{Context: context.Background(), Config: config, FileMigration: FileMigration{LogUrl: "https://example.com/log"}}).Init()
			closer, ok := any(it).(io.Closer)
			if tc.earlyStop {
				require.True(t, it.Next())
				require.True(t, ok, "migration iterators must support closing an unfinished stream")
				require.NoError(t, closer.Close())
				assert.Less(t, body.readBytes, len(tc.payload), "early close must not consume the remaining log")
			} else {
				for it.Next() {
				}
			}
			if tc.malformed {
				require.ErrorIs(t, it.Err(), io.ErrUnexpectedEOF)
			} else {
				require.NoError(t, it.Err())
			}
			assert.Equal(t, 1, body.closes)
			require.True(t, ok)
			require.NoError(t, closer.Close())
			assert.Equal(t, 1, body.closes, "closing an exhausted or stopped iterator is idempotent")
			assert.False(t, it.Next())
		})
	}
}
