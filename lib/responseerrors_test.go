package lib

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorParsersCloseOriginalBodyOnceAndKeepRemainingDiagnostics(t *testing.T) {
	for _, xml := range []bool{false, true} {
		t.Run(map[bool]string{false: "bounded generic response", true: "S3 response"}[xml], func(t *testing.T) {
			payload := strings.Repeat("x", 700)
			header := http.Header{"Content-Type": []string{"text/plain"}}
			if xml {
				payload = `<Error><Code>AccessDenied</Code><Message>denied</Message></Error>`
				header.Set("Content-Type", "application/xml")
			}
			original := &diagnosticResponseBody{Reader: strings.NewReader(payload)}
			response := &http.Response{StatusCode: http.StatusForbidden, Header: header, Body: original, ContentLength: int64(len(payload))}
			var err error
			if xml {
				err = S3XMLError(response)
			} else {
				err = NonOkError(response)
			}
			require.Error(t, err)
			if xml {
				var responseErr S3Error
				require.ErrorAs(t, err, &responseErr)
				assert.Equal(t, "AccessDenied", responseErr.Code)
			} else {
				var responseErr ResponseError
				require.ErrorAs(t, err, &responseErr)
				assert.Same(t, response, responseErr.Response)
				assert.Equal(t, header, responseErr.Response.Header)
				assert.Contains(t, err.Error(), payload[:512])
			}
			var remaining strings.Builder
			_, readErr := io.Copy(&remaining, response.Body)
			require.NoError(t, readErr)
			if xml {
				assert.Empty(t, remaining.String())
			} else {
				assert.Equal(t, payload[512:], remaining.String())
			}
			assert.Equal(t, 1, original.writeTos)
			require.NoError(t, response.Body.Close())
			CloseBody(response)
			assert.Equal(t, 1, original.closes)
		})
	}
}

type diagnosticResponseBody struct {
	*strings.Reader
	closes, writeTos int
}

func (b *diagnosticResponseBody) Close() error { b.closes++; return nil }
func (b *diagnosticResponseBody) WriteTo(w io.Writer) (int64, error) {
	b.writeTos++
	return b.Reader.WriteTo(w)
}
