package files_sdk

import (
	"testing"

	"github.com/Files-com/files-sdk-go/v2/lib"

	"github.com/stretchr/testify/assert"
)

func TestIter_Next_MaxPages(t *testing.T) {
	assert := assert.New(t)
	params := ListParams{Page: 0, PerPage: 5, MaxPages: 2}
	it := Iter{}
	it.ListParams = &params

	it.Query = func(lib.Values, ...RequestResponseOption) (*[]interface{}, string, error) {
		ret := make([]interface{}, params.PerPage)

		return &ret, "cursor", nil
	}
	recordCount := 0
	for it.Next() {
		recordCount += 1
	}
	assert.Equal(int(params.PerPage*params.MaxPages), recordCount)
	assert.Equal(nil, it.Err())
	assert.Equal("cursor", it.GetCursor())
}

func TestIter_Next_ZeroMaxPages(t *testing.T) {
	assert := assert.New(t)
	params := ListParams{Page: 0, PerPage: 2, MaxPages: 0}
	pages := make([][]interface{}, 0)
	pages = append(pages, make([]interface{}, params.PerPage))
	pages = append(pages, make([]interface{}, params.PerPage))
	pages = append(pages, make([]interface{}, 0))
	it := Iter{}
	it.ListParams = &params

	it.Query = func(lib.Values, ...RequestResponseOption) (*[]interface{}, string, error) {
		ret := pages[:1][0]
		pages = pages[1:]

		return &ret, "cursor", nil
	}
	recordCount := 0
	for it.Next() {
		recordCount += 1
	}
	assert.Equal(4, recordCount)
}

func TestIter_Next_PerPage_of_one(t *testing.T) {
	assert := assert.New(t)
	params := ListParams{Page: 0, PerPage: 1, MaxPages: 2}
	it := Iter{}
	it.ListParams = &params
	var sliceOfSliceInterfaces [2][]interface{}
	sliceOfSliceInterfaces[0] = make([]interface{}, params.PerPage)
	sliceOfSliceInterfaces[1] = make([]interface{}, 0)
	resultCounter := 0
	it.Query = func(lib.Values, ...RequestResponseOption) (*[]interface{}, string, error) {
		ret := sliceOfSliceInterfaces[resultCounter]
		resultCounter += 1
		return &ret, "cursor", nil
	}
	recordCount := 0
	for it.Next() {
		recordCount += 1
		assert.Equal(lib.Interface(), it.Current())
	}
	assert.Equal(1, recordCount)
}

func TestIter_Next_No_Cursor(t *testing.T) {
	assert := assert.New(t)
	params := ListParams{}
	it := Iter{}
	it.ListParams = &params
	resultCounter := 0
	it.Query = func(lib.Values, ...RequestResponseOption) (*[]interface{}, string, error) {
		ret := make([]interface{}, 1)
		resultCounter += 1
		return &ret, "", nil
	}
	recordCount := 0
	for it.Next() {
		recordCount += 1
		assert.Equal(lib.Interface(), it.Current())
	}
	assert.Equal(1, recordCount)
}
