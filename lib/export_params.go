package lib

import (
	"net/url"

	"github.com/google/go-querystring/query"
)

func ExportParams(i interface{}) (url.Values, error) {
	v, err := query.Values(i)
	if err != nil {
		panic(err)
	}

	return v, CheckRequired(i, &v)
}
