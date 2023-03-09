package downloadurl

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/itchyny/timefmt-go"
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
	}{
		{
			name: "amazon s3 date",
			args: args{
				url: func(t *testing.T, ti time.Time) string {
					u, err := url.ParseRequestURI(fmt.Sprintf("https://example.com?%v=%v&%v=%v", amazonS3.date, timefmt.Format(time.Now().UTC(), timeDateFormat), amazonS3.expire, 3*60))
					assert.NoError(t, err)
					return u.String()
				},
			},
			Time: time.Now().Add(time.Minute * 3).UTC(),
		},
		{
			name: "files date",
			args: args{
				url: func(t *testing.T, ti time.Time) string {
					u, err := url.ParseRequestURI(fmt.Sprintf("https://example.com?%v=%v", filesDate.date, timefmt.Format(time.Now().Add(time.Minute*3).UTC(), timeDateFormat)))
					assert.NoError(t, err)
					return u.String()
				},
			},
			Time: time.Now().Add(time.Minute * 3).UTC(),
		},
		{
			name: "google date",
			args: args{
				url: func(t *testing.T, ti time.Time) string {
					u, err := url.ParseRequestURI(fmt.Sprintf("https://example.com?%v=%v&%v=%v", googleDate.date, timefmt.Format(time.Now().UTC(), timeDateFormat), googleDate.expire, 3*60))
					assert.NoError(t, err)
					return u.String()
				},
			},
			Time: time.Now().Add(time.Minute * 3).UTC(),
		},
		{
			name: "azure blob storage",
			args: args{
				url: func(t *testing.T, ti time.Time) string {
					u, err := url.ParseRequestURI(fmt.Sprintf("https://filescomtests.blob.core.windows.net/testazureremote/ntie3buw/file-to-download.txt?sp=se=%v", timefmt.Format(time.Now().Add(time.Minute*3).UTC(), timeDateFormat)))
					assert.NoError(t, err)
					return u.String()
				},
			},
			Time: time.Now().Add(time.Minute * 3).UTC(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := New(tt.args.url(t, tt.Time))
			assert.NoError(t, err)
			assert.Equal(t, tt.Time.Truncate(time.Second), d.Time)
		})
	}
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
