package lib

import (
	"testing"

	"github.com/appscode/go-querystring/query"
	"github.com/stretchr/testify/assert"
)

type A struct {
	Name   string `url:"name" required:"true"`
	Age    int    `url:"age" required:"false"`
	Ignore int    `url:"-" required:"false"`
}

func TestCheckRequired_Valid(t *testing.T) {
	assert := assert.New(t)
	a := A{Name: "Dustin", Age: 90, Ignore: 50}
	values, _ := query.Values(a)

	err := CheckRequired(a, &values)
	assert.Equal(err, nil)
}

func TestCheckRequired_Invalid(t *testing.T) {
	assert := assert.New(t)
	a := A{Age: 90}
	values, _ := query.Values(a)

	err := CheckRequired(a, &values)
	assert.EqualError(err, "missing required field: A{}.Name")
}
