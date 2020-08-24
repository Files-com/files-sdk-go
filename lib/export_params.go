package lib

import (
	"net/url"

	"github.com/appscode/go-querystring/query"
)

func ExportParams(i interface{}) url.Values {
	v, err := query.Values(i)
	if err != nil {
		panic(err)
	}
	return v
}
