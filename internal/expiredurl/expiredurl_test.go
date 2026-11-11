package expiredurl

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/require"
)

type trackedBody struct {
	io.Reader
	read, closes int
}

func (b *trackedBody) Read(p []byte) (int, error) {
	n, err := b.Reader.Read(p)
	b.read += n
	return n, err
}

func (b *trackedBody) Close() error { b.closes++; return nil }

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

// Storage's expiry answer counts only within the first scanBytes of a 400 or 403 body that reads without error, even
// when the body is far longer; an expired type counts only in JSON that ends within them. The parsers still read,
// report and close that body exactly as they would have, and a success's body is not touched. Read, for a body the
// caller closes unread, decides the same.
func TestResponseErrorsMarksOnlyAnExpiryAnswerItReadsAndLeavesTheBodyToTheParsers(t *testing.T) {
	answer := `<Error><Code>ExpiredToken</Code><Message>The provided token has expired.</Message></Error>`
	readErr := errors.New("connection reset")
	typed := `{"type":"download_request_expired","error":"expired"}`
	typedToTheLimit := typed + strings.Repeat(" ", scanBytes-len(typed))
	for name, test := range map[string]struct {
		status  int
		body    string
		readErr error
		want    bool
	}{
		"answer at the start of a long body":     {http.StatusBadRequest, answer + strings.Repeat(" ", 1<<20), nil, true},
		"answer beyond the searched prefix":      {http.StatusBadRequest, strings.Repeat(" ", scanBytes) + answer, nil, false},
		"answer with another status":             {http.StatusInternalServerError, answer, nil, false},
		"answer in a body that fails to be read": {http.StatusForbidden, answer, readErr, false},
		"success that looks like an answer":      {http.StatusOK, answer, nil, false},
		"expired type in JSON that ends":         {http.StatusBadRequest, typedToTheLimit, nil, true},
		"expired type in JSON that goes on":      {http.StatusBadRequest, typedToTheLimit + `}`, nil, false},
	} {
		t.Run(name, func(t *testing.T) {
			response := func() (*http.Response, *trackedBody) {
				var reader io.Reader = strings.NewReader(test.body)
				if test.readErr != nil {
					reader = io.MultiReader(strings.NewReader(test.body[:20]), failingReader{test.readErr})
				}
				body := &trackedBody{Reader: reader}
				return &http.Response{StatusCode: test.status, Header: http.Header{"Content-Type": {"application/xml"}}, ContentLength: -1, Body: body}, body
			}
			plain, plainBody := response()
			want := lib.ResponseErrors(plain, lib.S3XMLError, lib.NonOkError)
			unread, _ := response()
			require.Equal(t, test.want, Read(unread))
			res, body := response()

			err := ResponseErrors(res, lib.S3XMLError, lib.NonOkError)

			require.Equal(t, test.want, Marked(err))
			if want == nil {
				require.NoError(t, err)
				require.Same(t, body, res.Body)
				require.Zero(t, body.read)
				return
			}
			require.EqualError(t, err, want.Error())
			require.Equal(t, plainBody.read, body.read)
			require.Equal(t, 1, body.closes)
			if test.readErr != nil {
				require.ErrorIs(t, err, test.readErr)
			}
		})
	}
}
