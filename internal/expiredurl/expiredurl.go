// Package expiredurl recognizes a transfer response that says the signed URL
// it was sent to has expired, before the SDK's response parsers consume its
// body. files_sdk.IsExpired reports what it recognizes.
package expiredurl

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

// scanBytes is how much of a refusal's body is searched for an expiry answer,
// which is far smaller. One byte more is read, to tell a body that ends within
// it from a longer one. It bounds the search, not what the response parsers
// read afterwards.
const scanBytes = 4096

// What storage services answer, in their XML error, when a signed URL has
// expired, matched without regard to case: Amazon S3 ("Request has expired"),
// Amazon and Google (ExpiredToken) and Azure.
var storageAnswers = [][]byte{
	[]byte("request has expired"),
	[]byte("<code>expiredtoken</code>"),
	[]byte("signature not valid in the specified time frame"),
}

// The Files.com error types that say a transfer URL has expired, matched exactly.
var filesTypes = []string{"download_request_expired", "upload_request_expired"}

// says reports whether a response with status whose body begins with read,
// at most scanBytes+1 bytes of it, says that its URL has expired: an HTTP 400
// or 403 whose whole body, ending within scanBytes, is a JSON object of one of
// filesTypes, or a body that is not JSON with one of storageAnswers in its
// first scanBytes. A body that starts as JSON does (an object, array or
// string, even one cut short or malformed) is never searched for words, and
// no other status counts.
func says(status int, read []byte) bool {
	if status != http.StatusBadRequest && status != http.StatusForbidden {
		return false
	}
	whole := len(read) <= scanBytes
	prefix := read[:min(len(read), scanBytes)]
	if trimmed := bytes.TrimSpace(prefix); len(trimmed) > 0 && bytes.IndexByte([]byte(`{["`), trimmed[0]) >= 0 {
		var body struct {
			Type string `json:"type"`
		}
		return whole && json.Unmarshal(trimmed, &body) == nil && slices.Contains(filesTypes, body.Type)
	}
	lower := bytes.ToLower(prefix)
	for _, answer := range storageAnswers {
		if bytes.Contains(lower, answer) {
			return true
		}
	}
	return false
}

func candidate(res *http.Response) bool {
	return res != nil && res.Body != nil && (res.StatusCode == http.StatusBadRequest || res.StatusCode == http.StatusForbidden)
}

// ResponseErrors runs handlers on res as lib.ResponseErrors does. When res is
// an HTTP 400 or 403 whose body says its URL has expired, the error they
// return is marked so that Marked reports it; its text and the errors it wraps
// are unchanged. For those statuses only, up to scanBytes+1 bytes of the body
// are read first and given back to the handlers ahead of the rest, followed by
// any read error, so the handlers read and close the body as they would have.
// A body that fails to read within those bytes is not taken as an expiry
// answer.
func ResponseErrors(res *http.Response, handlers ...func(*http.Response) error) error {
	expired := false
	if candidate(res) {
		original := res.Body
		read, readErr := io.ReadAll(io.LimitReader(original, scanBytes+1))
		var rest io.Reader = original
		if readErr != nil {
			rest = failedRead{readErr}
		}
		res.Body = replayedBody{Reader: io.MultiReader(bytes.NewReader(read), rest), Closer: original}
		expired = readErr == nil && says(res.StatusCode, read)
	}
	err := lib.ResponseErrors(res, handlers...)
	if err != nil && expired {
		return marked{err}
	}
	return err
}

// Read reports whether res, an unsuccessful response whose body the caller
// closes next without reading it, says its URL has expired. It reads at most
// scanBytes+1 bytes of the body, and only for an HTTP 400 or 403.
func Read(res *http.Response) bool {
	if !candidate(res) {
		return false
	}
	read, err := io.ReadAll(io.LimitReader(res.Body, scanBytes+1))
	return err == nil && says(res.StatusCode, read)
}

// Marked reports whether err, or an error it wraps, came from ResponseErrors
// for a response that says its URL has expired.
func Marked(err error) bool {
	var expired marked
	return errors.As(err, &expired)
}

type marked struct{ error }

func (m marked) Unwrap() error { return m.error }

type failedRead struct{ err error }

func (r failedRead) Read([]byte) (int, error) { return 0, r.err }

type replayedBody struct {
	io.Reader
	io.Closer
}
