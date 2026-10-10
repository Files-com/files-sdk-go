package files_sdk

import (
	"errors"
	"fmt"
	"iter"
	"net/url"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestIter_Reload_RequestsFirstPageAfterFinishedTraversal(t *testing.T) {
	listing := &pagedQuery{pages: map[string]queryPage{
		"saved-cursor": {values: []interface{}{"old"}},
	}}
	it := &Iter{Query: listing.Query, ListParams: &ListParams{Cursor: "saved-cursor", PerPage: 1}}
	assert.Equal(t, []interface{}{"old"}, collect(it))

	listing.pages[""] = queryPage{values: []interface{}{"new"}}
	reloaded := it.Reload()

	assert.Equal(t, []interface{}{"new"}, collect(reloaded))
	assert.False(t, it.Next(), "the original traversal stays finished")
	require.Equal(t, []string{"saved-cursor", ""}, listing.requestedCursors())
	assert.Equal(t, "1", listing.requests[1].Get("per_page"))
}

func TestIter_Reload_DoesNotServeRowsLeftOnSourcePage(t *testing.T) {
	listing := &pagedQuery{pages: map[string]queryPage{
		"": {values: []interface{}{"a1", "a2"}, next: "page-2"},
	}}
	it := &Iter{Query: listing.Query, ListParams: &ListParams{MaxPages: 1}}
	assert.True(t, it.Next())

	listing.pages[""] = queryPage{values: []interface{}{"b1", "b2"}, next: "page-2"}

	assert.Equal(t, []interface{}{"b1", "b2"}, collect(it.Reload()))
	assert.Equal(t, []interface{}{"a2"}, collect(it), "the original continues from its own page")
	assert.Equal(t, []string{"", ""}, listing.requestedCursors(), "MaxPages stops both traversals after one page")
}

func TestIter_Reload_ClearsSourceErrorAndKeepsOnPageError(t *testing.T) {
	listing := &pagedQuery{pages: map[string]queryPage{
		"": {err: errors.New("service unavailable")},
	}}
	it := &Iter{
		Query:      listing.Query,
		ListParams: &ListParams{},
		OnPageError: func(err error) (*[]interface{}, error) {
			return &[]interface{}{}, fmt.Errorf("listing users: %w", err)
		},
	}
	assert.Empty(t, collect(it))
	assert.EqualError(t, it.Err(), "listing users: service unavailable")

	listing.pages[""] = queryPage{values: []interface{}{"u1"}, next: "page-2"}
	listing.pages["page-2"] = queryPage{err: errors.New("timeout")}
	reloaded := it.Reload()

	assert.Equal(t, []interface{}{"u1"}, collect(reloaded))
	assert.EqualError(t, reloaded.Err(), "listing users: timeout")
	assert.EqualError(t, it.Err(), "listing users: service unavailable")
}

func TestIter_Reload_KeepsSourceCursorWithEmbeddedListParamsPointer(t *testing.T) {
	type searchParams struct {
		*ListParams
		Search string `url:"search,omitempty"`
	}
	listing := &pagedQuery{pages: map[string]queryPage{
		"": {values: []interface{}{"first"}, next: "page-2"},
	}}
	source := &Iter{
		Query:      listing.Query,
		ListParams: &searchParams{ListParams: &ListParams{Cursor: "resume", PerPage: 7}, Search: "ops"},
	}

	reloaded := source.Reload()
	assert.True(t, reloaded.Next())

	assert.Equal(t, "resume", source.GetCursor(), "the reload keeps its cursor to itself")
	require.Len(t, listing.requests, 1)
	request := listing.requests[0]
	assert.Empty(t, request.Get("cursor"))
	assert.Equal(t, "7", request.Get("per_page"))
	assert.Equal(t, "ops", request.Get("search"))
}

func TestIter_ParameterExportErrorStopsBeforeRequest(t *testing.T) {
	type requiredSearchParams struct {
		Search string `url:"search" required:"true"`
		ListParams
	}
	listing := &pagedQuery{pages: map[string]queryPage{
		"": {values: []interface{}{"unfiltered"}},
	}}
	source := &Iter{Query: listing.Query, ListParams: &requiredSearchParams{}}
	reloaded := source.Reload()

	for _, it := range []IterI{source, reloaded} {
		assert.False(t, it.Next())
		assert.ErrorContains(t, it.Err(), "missing required field")
	}
	assert.Empty(t, listing.requests, "nothing is listed without the listing's parameters")
}

func TestIterAll(t *testing.T) {
	t.Run("requests pages lazily, in order", func(t *testing.T) {
		listing := &pagedQuery{pages: map[string]queryPage{
			"":       {values: []interface{}{"a", "b"}, next: "page-2"},
			"page-2": {values: []interface{}{"c"}},
		}}
		it := &Iter{Query: listing.Query, ListParams: &ListParams{}}

		resources := IterAll[string](it)
		assert.Empty(t, listing.requests, "creating the iterator requests nothing")
		values, errs := collectAll(resources)

		assert.Equal(t, []string{"a", "b", "c"}, values)
		assert.Empty(t, errs)
		assert.Equal(t, []string{"", "page-2"}, listing.requestedCursors())
	})

	t.Run("a break requests nothing more, and the next loop continues", func(t *testing.T) {
		listing := &pagedQuery{pages: map[string]queryPage{
			"":       {values: []interface{}{"a", "b"}, next: "page-2"},
			"page-2": {values: []interface{}{"c"}},
		}}
		it := &Iter{Query: listing.Query, ListParams: &ListParams{}}

		var first []string
		for value, err := range IterAll[string](it) {
			require.NoError(t, err)
			first = append(first, value)
			if len(first) == 1 {
				break
			}
		}
		assert.Equal(t, []string{"a"}, first)
		assert.Len(t, listing.requests, 1)

		values, errs := collectAll(IterAll[string](it))
		assert.Equal(t, []string{"b", "c"}, values)
		assert.Empty(t, errs)

		values, errs = collectAll(IterAll[string](it))
		assert.Empty(t, values, "an ended listing yields nothing and is not restarted")
		assert.Empty(t, errs)
		assert.Len(t, listing.requests, 2)
	})

	t.Run("continues after Next, and yields nothing once Next reached the end", func(t *testing.T) {
		listing := &pagedQuery{pages: map[string]queryPage{"": {values: []interface{}{"a", "b", "c"}}}}
		it := &Iter{Query: listing.Query, ListParams: &ListParams{}}
		require.True(t, it.Next())

		values, _ := collectAll(IterAll[string](it))
		assert.Equal(t, []string{"b", "c"}, values)

		ended := &Iter{Query: listing.Query, ListParams: &ListParams{}}
		for ended.Next() {
		}
		values, errs := collectAll(IterAll[string](ended))
		assert.Empty(t, values)
		assert.Empty(t, errs)
	})

	t.Run("yields a page error once and keeps it in Err", func(t *testing.T) {
		pageErr := fmt.Errorf("page 2: %w", errors.New("service unavailable"))
		listing := &pagedQuery{pages: map[string]queryPage{
			"":       {values: []interface{}{"a"}, next: "page-2"},
			"page-2": {err: pageErr},
		}}
		it := &Iter{Query: listing.Query, ListParams: &ListParams{}}

		var values []string
		var errs []error
		for value, err := range IterAll[string](it) {
			if err != nil {
				assert.Empty(t, value)
				errs = append(errs, err)
				continue
			}
			values = append(values, value)
		}

		assert.Equal(t, []string{"a"}, values)
		require.Len(t, errs, 1)
		assert.Same(t, pageErr, errs[0])
		assert.Same(t, pageErr, it.Err())
		values, errs = collectAll(IterAll[string](it))
		assert.Empty(t, values)
		assert.Empty(t, errs, "the error is yielded only by the loop that hit it")
		assert.Len(t, listing.requests, 2)
	})

	t.Run("MaxPages and an empty listing end without an extra yield", func(t *testing.T) {
		listing := &pagedQuery{pages: map[string]queryPage{
			"":       {values: []interface{}{"a"}, next: "page-2"},
			"page-2": {values: []interface{}{"b"}},
		}}
		limited := &Iter{Query: listing.Query, ListParams: &ListParams{MaxPages: 1}}
		values, errs := collectAll(IterAll[string](limited))
		assert.Equal(t, []string{"a"}, values)
		assert.Empty(t, errs)

		empty := &Iter{Query: (&pagedQuery{pages: map[string]queryPage{"": {}}}).Query, ListParams: &ListParams{}}
		values, errs = collectAll(IterAll[string](empty))
		assert.Empty(t, values)
		assert.Empty(t, errs)
	})

	t.Run("a resource of another type ends the loop with an error", func(t *testing.T) {
		listing := &pagedQuery{pages: map[string]queryPage{"": {values: []interface{}{"a", "b"}}}}
		it := &Iter{Query: listing.Query, ListParams: &ListParams{}}

		values, errs := collectAll(IterAll[int](it))
		assert.Empty(t, values)
		require.Len(t, errs, 1)
		assert.ErrorContains(t, errs[0], "string, not int")
		typeErr := errs[0]
		assert.Same(t, typeErr, it.Err())

		values, errs = collectAll(IterAll[int](it))
		assert.Empty(t, values)
		assert.Empty(t, errs, "the error is yielded only by the loop that hit it")
		assert.Same(t, typeErr, it.Err())
		assert.Len(t, listing.requests, 1)
	})
}

// collectAll ranges over resources, collecting the values yielded with a nil
// error and the errors.
func collectAll[T any](resources iter.Seq2[T, error]) (values []T, errs []error) {
	for value, err := range resources {
		if err != nil {
			errs = append(errs, err)
			continue
		}
		values = append(values, value)
	}
	return values, errs
}

// pagedQuery serves pages by requested cursor and records every request.
type pagedQuery struct {
	pages    map[string]queryPage
	requests []url.Values
}

type queryPage struct {
	values []interface{}
	next   string
	err    error
}

func (q *pagedQuery) Query(params lib.Values, _ ...RequestResponseOption) (*[]interface{}, string, error) {
	values, err := params.ToValues()
	if err != nil {
		return &[]interface{}{}, "", err
	}
	q.requests = append(q.requests, values)
	page := q.pages[values.Get("cursor")]
	return &page.values, page.next, page.err
}

func (q *pagedQuery) requestedCursors() (cursors []string) {
	for _, request := range q.requests {
		cursors = append(cursors, request.Get("cursor"))
	}
	return cursors
}

func collect(it IterI) (values []interface{}) {
	for it.Next() {
		values = append(values, it.Current())
	}
	return values
}
