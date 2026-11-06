package test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dnaeon/go-vcr/cassette"
)

var recordedSignatureDate = regexp.MustCompile(`([?&]|\\u0026)X-(Amz|Files|Goog)-Date=[0-9]{8}T[0-9]{6}Z`)
var recordedPartExpiry = regexp.MustCompile(`("expires"[[:space:]]*:[[:space:]]*")([^"]+)(")`)

const signatureDateLayout = "20060102T150405Z"

// Old cassettes retain their signed query bytes. Relocate only their known time
// fields for offline replay, using one offset for the complete recorded session.
// The SDK still checks expiration normally; genuinely expired recorded links
// remain expired relative to the session's captured Date.
type replayTimeline struct{ shift time.Duration }

func newReplayTimeline(name string, now time.Time) (*replayTimeline, error) {
	c, err := cassette.Load(name)
	if err != nil {
		return nil, err
	}
	var anchor time.Time
	transfers := false
	for _, interaction := range c.Interactions {
		transfers = transfers || recordedSignatureDate.MatchString(interaction.Request.URL) || recordedSignatureDate.MatchString(interaction.Response.Body)
		for _, value := range interaction.Response.Headers.Values("Date") {
			date, err := http.ParseTime(value)
			if err == nil && (anchor.IsZero() || date.Before(anchor)) {
				anchor = date
			}
		}
	}
	if !transfers {
		return nil, nil
	}
	if anchor.IsZero() {
		return nil, fmt.Errorf("recorded transfer fixture %q has no usable Date header", name)
	}
	return &replayTimeline{shift: now.UTC().Truncate(time.Second).Sub(anchor)}, nil
}

func (t *replayTimeline) normalize(value string) string {
	value = recordedSignatureDate.ReplaceAllStringFunc(value, func(match string) string {
		separator := strings.IndexByte(match, '=') + 1
		date, err := time.Parse(signatureDateLayout, match[separator:])
		if err != nil {
			return match
		}
		return match[:separator] + date.Add(t.shift).UTC().Format(signatureDateLayout)
	})
	return recordedPartExpiry.ReplaceAllStringFunc(value, func(match string) string {
		parts := recordedPartExpiry.FindStringSubmatch(match)
		date, err := time.Parse(time.RFC3339Nano, parts[2])
		if err != nil {
			return match
		}
		return parts[1] + date.Add(t.shift).Format(time.RFC3339Nano) + parts[3]
	})
}

type replayTimeTransport struct {
	base     http.RoundTripper
	timeline *replayTimeline
	apiRoot  *url.URL
}

func (t replayTimeTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := t.base.RoundTrip(request)
	if err != nil || response == nil {
		return response, err
	}
	response.Header = response.Header.Clone()
	for _, header := range []string{"Date", "Retry-After"} {
		for index, value := range response.Header.Values(header) {
			if date, err := http.ParseTime(value); err == nil {
				response.Header[http.CanonicalHeaderKey(header)][index] = date.Add(t.timeline.shift).UTC().Format(http.TimeFormat)
			}
		}
	}
	for index, value := range response.Header.Values("Location") {
		response.Header["Location"][index] = t.timeline.normalize(value)
	}
	// JSON files downloaded from storage are content, not API metadata. Their
	// bytes must not be changed merely because they contain an expiry/date key.
	apiMetadata := strings.EqualFold(request.URL.Scheme, t.apiRoot.Scheme) &&
		strings.EqualFold(request.URL.Host, t.apiRoot.Host) &&
		strings.HasPrefix(request.URL.Path, t.apiRoot.Path+"/")
	if apiMetadata && strings.Contains(response.Header.Get("Content-Type"), "application/json") {
		body, readErr := io.ReadAll(response.Body)
		closeErr := response.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		if closeErr != nil {
			return nil, closeErr
		}
		body = []byte(t.timeline.normalize(string(body)))
		response.Body = io.NopCloser(bytes.NewReader(body))
		response.ContentLength = int64(len(body))
		if response.Header.Get("Content-Length") != "" {
			response.Header.Set("Content-Length", strconv.Itoa(len(body)))
		}
	}
	return response, nil
}
