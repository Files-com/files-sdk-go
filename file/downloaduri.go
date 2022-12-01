package file

import (
	"net/url"
	"time"

	"github.com/itchyny/timefmt-go"
)

type DownloadUri string

type DownloadUriValid struct {
	*url.URL
	time.Time
	Remaining time.Duration
}

func (d DownloadUri) ToUrl() (u *url.URL, err error) {
	u, err = url.Parse(string(d))
	return
}

func (d DownloadUri) ToTime() (t time.Time, err error) {
	var u *url.URL
	u, err = d.ToUrl()
	if err != nil {
		return t, err
	}
	expires := u.Query().Get("X-Amz-Date")
	t, err = timefmt.Parse(expires, "%Y%m%dT%H%M%SZ")
	if err == nil {
		return t, err
	}
	expires = u.Query().Get("X-Files-Date")
	t, err = timefmt.Parse(expires, "%Y%m%dT%H%M%SZ")
	if err == nil {
		return t, err
	}
	expires = u.Query().Get("X-Goog-Date")
	t, err = timefmt.Parse(expires, "%Y%m%dT%H%M%SZ")
	return
}

func (d DownloadUri) Valid(within time.Duration) (s DownloadUriValid, valid bool, err error) {
	s.URL, err = d.ToUrl()
	s.Time, err = d.ToTime()
	if err != nil {
		return
	}
	s.Remaining = s.Time.Sub(time.Now())
	return s, s.Remaining < within, err
}
