package downloadurl

import (
	"net/url"
	"strconv"
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

	return
}

func (d *URL) parseUrl(urlStr string) (err error) {
	d.URL, err = url.Parse(urlStr)
	if err != nil {
		return
	}
	return
}

type expire struct {
	date       string
	expire     string
	expireDate string
}

var (
	amazonS3       = expire{"X-Amz-Date", "X-Amz-Expires", ""}
	filesDate      = expire{"X-Files-Date", "X-Files-Expires", ""}
	googleDate     = expire{"X-Goog-Date", "X-Goog-Expires", ""}
	azureBlob      = expire{"sp", "", "se"}
	timeDateFormat = "%Y%m%dT%H%M%SZ"
)

var Dates = []expire{amazonS3, filesDate, googleDate, azureBlob}

func (d *URL) parseTime() (err error) {
	for _, parser := range Dates {
		date := d.URL.Query().Get(parser.date)
		if parser.expireDate != "" {
			query := (&url.URL{RawQuery: date}).Query()
			if len(query[parser.expireDate]) == 1 {
				d.Time, err = timefmt.Parse(query[parser.expireDate][0], timeDateFormat)
			} else {
				continue
			}
		} else {
			d.Time, err = timefmt.Parse(date, timeDateFormat)
		}

		if parser.expire != "" {
			duration, err := strconv.Atoi(d.URL.Query().Get(parser.expire))
			if err == nil {
				t := time.Second * time.Duration(duration)
				d.Time = d.Time.Add(t)
			}
		}
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
