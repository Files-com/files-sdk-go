package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"time"

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

var (
	zeroTimeString = time.Time{}.Format(time.RFC3339)
)

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
	b, err = sanitizeJSON(b)
	return bytes.NewBuffer(b), err
}

func (p Params) ToValues() (url.Values, error) {
	v, err := p.values()
	if err != nil {
		return url.Values{}, err
	}

	if err := CheckRequired(p.Params); err != nil {
		return url.Values{}, err
	}

	return removeDash(v), nil
}

// valuesEncoder is implemented by request types whose query values need more
// than their struct tags, such as those with DecimalOverride fields.
type valuesEncoder interface {
	ToValues() (url.Values, error)
}

func (p Params) values() (url.Values, error) {
	// query.Values only consults encoders on fields, not on the root value.
	if encoder, ok := p.Params.(valuesEncoder); ok && !isNilPointer(p.Params) {
		return encoder.ToValues()
	}
	return query.Values(p.Params)
}

func isNilPointer(v interface{}) bool {
	value := reflect.ValueOf(v)
	return value.Kind() == reflect.Pointer && value.IsNil()
}

func sanitizeJSON(b []byte) ([]byte, error) {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return b, err
	}

	for key, value := range m {
		switch {
		case key == "-":
			delete(m, key)
		case value == zeroTimeString:
			m[key] = nil
		}
	}

	return json.Marshal(m)
}

func removeDash(params url.Values) url.Values {
	for key := range params {
		if string(key[0]) == "-" {
			params.Del(key)
		}
	}

	return params
}

type UnmarshalJSON interface {
	UnmarshalJSON(data []byte) error
}

type Resource struct {
	Path   string
	Params interface{}
	Method string
	Entity UnmarshalJSON
}

func (r Resource) Out() (ResourceOut, error) {
	path, err := BuildPath(r.Path, r.Params)
	if err != nil {
		return ResourceOut{}, err
	}
	return ResourceOut{
		Resource: Resource{
			Path:   path,
			Method: r.Method,
			Entity: r.Entity,
		},
		Values: Params{Params: r.Params},
	}, nil
}

type ResourceOut struct {
	Resource
	Values
}
