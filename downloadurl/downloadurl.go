package downloadurl

import (
	"net/url"
	"time"

	"github.com/itchyny/timefmt-go"
)

type URL struct {
	*url.URL
	time.Time
}

func New(urlStr string) (d *URL, err error) {
	d = &URL{}
	err = d.parseUrl(urlStr)
	if err != nil {
		return
	}

	err = d.parseTime()

	if err != nil {
		return
	}

	err = d.parseTime()

	if err != nil {
		return
	}
	return
}

func (d *URL) parseUrl(urlStr string) (err error) {
	d.URL, err = url.Parse(urlStr)
	if err != nil {
		return
	}
	return
}

const (
	AmazonS3Date   = "X-Amz-Date"
	FilesDate      = "X-Files-Date"
	GoogleDate     = "X-Goog-Date"
	TimeDateFormat = "%Y%m%dT%H%M%SZ"
)

var Dates = []string{AmazonS3Date, FilesDate, GoogleDate}

func (d *URL) parseTime() (err error) {
	for _, date := range Dates {
		expires := d.URL.Query().Get(date)
		d.Time, err = timefmt.Parse(expires, TimeDateFormat)
		if err == nil {
			break
		}
	}
	return err
}

func (d *URL) Valid(within time.Duration) (remaining time.Duration, valid bool) {
	remaining = d.Time.Sub(time.Now())
	valid = remaining >= within
	return
}
