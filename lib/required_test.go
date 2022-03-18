package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
	Name string `url:"name" required:"true"`
	Age  int    `url:"age" required:"false"`
}

func TestCheckRequired_Valid(t *testing.T) {
	assert := assert.New(t)
	a := A{Name: "Dustin", Age: 90}

	err := CheckRequired(a)
	assert.Equal(err, nil)
}

func TestCheckRequired_Invalid(t *testing.T) {
	assert := assert.New(t)
	a := A{Age: 90}

	err := CheckRequired(a)
	assert.EqualError(err, "missing required field: A{}.Name")
}
