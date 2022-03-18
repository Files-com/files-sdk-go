package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/appscode/go-querystring/query"
)

type Params struct {
	Params interface{}
}

type Values interface {
	ToValues() (url.Values, error)
	ToJSON() (io.Reader, error)
}

type ExportValues struct {
	url.Values
}

func (m ExportValues) ToValues() (url.Values, error) {
	return m.Values, nil
}

func (m ExportValues) ToJSON() (io.Reader, error) {
	return nil, fmt.Errorf("not Implemented")
}

func (p Params) ToJSON() (io.Reader, error) {
	_, err := p.ToValues()
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(p.Params)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

func (p Params) ToValues() (url.Values, error) {
	v, err := query.Values(p.Params)
	if err != nil {
		return url.Values{}, err
	}

	if err := CheckRequired(p.Params); err != nil {
		return url.Values{}, err
	}

	return removeDash(v), nil
}

func removeDash(params url.Values) url.Values {
	for key := range params {
		if string(key[0]) == "-" {
			params.Del(key)
		}
	}

	return params
}
