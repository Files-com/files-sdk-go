package files_sdk

import (
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
)

func TestIter_Next_MaxPages(t *testing.T) {
	assert := assert.New(t)
	params := ListParams{PerPage: 5, MaxPages: 2}
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
	params := ListParams{PerPage: 2, MaxPages: 0}
	pages := make([][]interface{}, 0)
	pages = append(pages, make([]interface{}, params.PerPage))
	pages = append(pages, make([]interface{}, params.PerPage))
	pages = append(pages, make([]interface{}, 0))
	it := Iter{}
	it.ListParams = &params

	it.Query = func(lib.Values, ...RequestResponseOption) (*[]interface{}, string, error) {
		ret := pages[:1][0]
		pages = pages[1:]
		cursor := "cursor"
		if len(pages) == 0 {
			cursor = ""
		}

		return &ret, cursor, nil
	}
	recordCount := 0
	for it.Next() {
		recordCount += 1
	}
	assert.Equal(4, recordCount)
}

func TestIter_Next_FollowsCursorAfterEmptyPage(t *testing.T) {
	params := ListParams{}
	pages := [][]interface{}{{"first"}, {}, {"last"}}
	cursors := []string{"cursor-1", "cursor-2", ""}
	requests := 0
	it := Iter{ListParams: &params}
	it.Query = func(lib.Values, ...RequestResponseOption) (*[]interface{}, string, error) {
		page := pages[requests]
		cursor := cursors[requests]
		requests++
		return &page, cursor, nil
	}

	var records []interface{}
	for it.Next() {
		records = append(records, it.Current())
	}

	assert.Equal(t, []interface{}{"first", "last"}, records)
	assert.Equal(t, 3, requests)
	assert.NoError(t, it.Err())
}

func TestIter_Next_StopsWhenEmptyPageCursorDoesNotAdvance(t *testing.T) {
	params := ListParams{Cursor: "stalled-cursor"}
	requests := 0
	it := Iter{ListParams: &params}
	it.Query = func(lib.Values, ...RequestResponseOption) (*[]interface{}, string, error) {
		requests++
		page := []interface{}{}
		return &page, "stalled-cursor", nil
	}

	assert.False(t, it.Next())
	assert.Equal(t, 1, requests)
	assert.EqualError(t, it.Err(), "pagination cursor did not advance after an empty page")
}

func TestIter_EOFPageBeforeFetch(t *testing.T) {
	assert.False(t, (&Iter{}).EOFPage())
}

func TestIter_Next_PerPage_of_one(t *testing.T) {
	assert := assert.New(t)
	params := ListParams{PerPage: 1, MaxPages: 2}
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
		assert.Equal(nil, it.Current())
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
		assert.Equal(nil, it.Current())
	}
	assert.Equal(1, recordCount)
}
