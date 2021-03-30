package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIter_Next_MaxPages(t *testing.T) {
	assert := assert.New(t)
	params := ListParams{}
	params.Set(0, 5, "", 2)
	it := Iter{}
	it.ListParams = &params

	it.Query = func() (*[]interface{}, string, error) {
		ret := make([]interface{}, params.PerPage)

		return &ret, "cursor", nil
	}
	recordCount := 0
	for it.Next() {
		recordCount += 1
	}
	assert.Equal(params.PerPage*params.MaxPages, recordCount)
	assert.Equal(nil, it.Err())
	assert.Equal("cursor", it.GetCursor())
}

func TestIter_Next_ZeroMaxPages(t *testing.T) {
	assert := assert.New(t)
	params := ListParams{}
	params.Set(0, 5, "", 0)
	it := Iter{}
	it.ListParams = &params

	it.Query = func() (*[]interface{}, string, error) {
		ret := make([]interface{}, params.PerPage)

		return &ret, "cursor", nil
	}
	recordCount := 0
	for it.Next() {
		recordCount += 1
	}
	assert.Equal(params.PerPage*params.MaxPages, recordCount)
}

func TestIter_Next_PerPage_of_one(t *testing.T) {
	assert := assert.New(t)
	params := ListParams{}
	params.Set(0, 1, "", 2)
	it := Iter{}
	it.ListParams = &params
	var sliceOfSliceInterfaces [2][]interface{}
	sliceOfSliceInterfaces[0] = make([]interface{}, params.PerPage)
	sliceOfSliceInterfaces[1] = make([]interface{}, 0)
	resultCounter := 0
	it.Query = func() (*[]interface{}, string, error) {
		ret := sliceOfSliceInterfaces[resultCounter]
		resultCounter += 1
		return &ret, "cursor", nil
	}
	recordCount := 0
	for it.Next() {
		recordCount += 1
		assert.Equal(Interface(), it.Current())
	}
	assert.Equal(1, recordCount)
}
