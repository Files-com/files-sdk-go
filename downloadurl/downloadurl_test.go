package downloadurl

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDownloadUrl_New(t *testing.T) {
	type args struct {
		url func(t *testing.T, ti time.Time) string
	}
	tests := []struct {
		name string
		args
		time.Time
		urlType URLType
	}{
		{
			name: "amazon s3 date",
			args: args{
				url: func(t *testing.T, ti time.Time) string {
					u, err := url.ParseRequestURI(fmt.Sprintf("https://example.com?%v=%v&%v=%v", amazonS3.date, time.Now().UTC().Format(timeDateFormat), amazonS3.expire, 3*60))
					assert.NoError(t, err)
					return u.String()
				},
			},
			Time:    time.Now().Add(time.Minute * 3).UTC(),
			urlType: AmazonS3,
		},
		{
			name: "files date",
			args: args{
				url: func(t *testing.T, ti time.Time) string {
					u, err := url.ParseRequestURI(fmt.Sprintf("https://example.com?%v=%v", filesDate.date, time.Now().Add(time.Minute*3).UTC().Format(timeDateFormat)))
					assert.NoError(t, err)
					return u.String()
				},
			},
			Time:    time.Now().Add(time.Minute * 3).UTC(),
			urlType: Files,
		},
		{
			name: "google date",
			args: args{
				url: func(t *testing.T, ti time.Time) string {
					u, err := url.ParseRequestURI(fmt.Sprintf("https://example.com?%v=%v&%v=%v", googleDate.date, time.Now().UTC().Format(timeDateFormat), googleDate.expire, 3*60))
					assert.NoError(t, err)
					return u.String()
				},
			},
			Time:    time.Now().Add(time.Minute * 3).UTC(),
			urlType: Google,
		},
		{
			name: "azure blob storage",
			args: args{
				url: func(t *testing.T, ti time.Time) string {
					u, err := url.ParseRequestURI(fmt.Sprintf("https://filescomtests.blob.core.windows.net/testazureremote/ntie3buw/file-to-download.txt?sp=se=%v", time.Now().Add(time.Minute*3).UTC().Format(timeDateFormat)))
					assert.NoError(t, err)
					return u.String()
				},
			},
			Time:    time.Now().Add(time.Minute * 3).UTC(),
			urlType: Azure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := New(tt.args.url(t, tt.Time))
			assert.NoError(t, err)
			assert.Equal(t, tt.Time.Truncate(time.Second), d.Time)
			assert.Equal(t, tt.urlType, d.Type)
		})
	}
}

func TestDownloadUrl_NewFromURL(t *testing.T) {
	expiresAt := time.Now().Add(time.Minute * 3).UTC().Truncate(time.Second)
	parsedURL, err := url.ParseRequestURI(fmt.Sprintf("https://example.com?%v=%v", filesDate.date, expiresAt.Format(timeDateFormat)))
	assert.NoError(t, err)

	d, err := NewFromURL(parsedURL)
	assert.NoError(t, err)
	assert.Same(t, parsedURL, d.URL)
	assert.Equal(t, expiresAt, d.Time)
}

func TestDownloadUrl_Init(t *testing.T) {
	expiresAt := time.Now().Add(time.Minute * 3).UTC().Truncate(time.Second)
	parsedURL, err := url.ParseRequestURI(fmt.Sprintf("https://example.com?%v=%v", filesDate.date, expiresAt.Format(timeDateFormat)))
	assert.NoError(t, err)

	d := &URL{}
	err = d.Init(parsedURL)

	assert.NoError(t, err)
	assert.Same(t, parsedURL, d.URL)
	assert.Equal(t, expiresAt, d.Time)
}

func TestDownloadUrl_InitNilURL(t *testing.T) {
	d := &URL{}

	err := d.Init(nil)

	assert.ErrorContains(t, err, "nil URL")
	assert.Nil(t, d.URL)
}

func TestDownloadUrl_InitResetsPreviousState(t *testing.T) {
	expiresAt := time.Now().Add(time.Minute * 3).UTC().Truncate(time.Second)
	parsedURL, err := url.ParseRequestURI(fmt.Sprintf("https://example.com?%v=%v", filesDate.date, expiresAt.Format(timeDateFormat)))
	assert.NoError(t, err)

	d, err := NewFromURL(parsedURL)
	assert.NoError(t, err)
	assert.False(t, d.Time.IsZero())

	invalidURL, err := url.ParseRequestURI("https://example.com?jwt=test")
	assert.NoError(t, err)

	err = d.Init(invalidURL)

	assert.Error(t, err)
	assert.Same(t, invalidURL, d.URL)
	assert.True(t, d.Time.IsZero())
	assert.Empty(t, d.Type)
}

func TestDownloadUrl_Valid(t *testing.T) {
	type args struct {
		within time.Duration
	}
	tests := []struct {
		name string
		args
		*URL
		valid     bool
		remaining time.Duration
	}{
		{
			name:      "Is not within time small difference",
			args:      args{within: time.Millisecond * 500},
			URL:       &URL{Time: time.Now().Add(-time.Millisecond * 500)},
			valid:     false,
			remaining: -time.Millisecond * 500,
		},
		{
			name:      "Is not within time large difference",
			args:      args{within: time.Millisecond * 500},
			URL:       &URL{Time: time.Now().Add(-time.Hour * 24)},
			valid:     false,
			remaining: -time.Hour * 24,
		},
		{
			name:      "Is within time small difference",
			args:      args{within: time.Millisecond * 500},
			URL:       &URL{Time: time.Now().Add(time.Millisecond * 900)},
			valid:     true,
			remaining: time.Millisecond * 900,
		},
		{
			name:      "Is within time large difference",
			args:      args{within: time.Millisecond * 500},
			URL:       &URL{Time: time.Now().Add(time.Hour * 24)},
			valid:     true,
			remaining: time.Hour * 24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remaining, valid := tt.URL.Valid(tt.args.within)
			assert.Equal(t, tt.valid, valid)
			assert.InDelta(t, remaining, tt.remaining, float64(time.Millisecond*100))
		})
	}
}
