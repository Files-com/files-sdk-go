package files_sdk

import (
	"errors"
	"net/url"
	"slices"
	"sync"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

type ListParams struct {
	PerPage  int64  `json:"per_page,omitempty" url:"per_page,omitempty" required:"false"`
	Cursor   string `json:"cursor,omitempty" url:"cursor,omitempty" required:"false"`
	MaxPages int64  `json:"-" url:"-"`
}

// ListParamsContainer is a general interface for which all list parameter
// structs should comply. They achieve this by embedding a ListParams struct
// and inheriting its implementation of this interface.
type ListParamsContainer interface {
	GetListParams() *ListParams
}

// GetListParams returns a ListParams struct (itself). It exists because any
// structs that embed ListParams will inherit it, and thus implement the
// ListParamsContainer interface.
func (p *ListParams) GetListParams() *ListParams {
	return p
}

type OnPageError func(error) (*[]interface{}, error)
type Query func(params lib.Values, opts ...RequestResponseOption) (*[]interface{}, string, error)

type IterI interface {
	Next() bool
	Current() interface{}
	Err() error
}

var _ IterI = (*Iter)(nil)

type TypedIterI[T any] interface {
	Next() bool
	Current() interface{}
	Resource() T
	Err() error
}

type IterPagingI interface {
	IterI
	EOFPage() bool
}

var _ IterPagingI = (*Iter)(nil)

type ResourceIterator interface {
	Iterate(interface{}, ...RequestResponseOption) (IterI, error)
}

type ReloadIterator interface {
	Reload(opts ...RequestResponseOption) IterI
}

var _ ReloadIterator = (*Iter)(nil)

type ResourceLoader interface {
	LoadResource(interface{}, ...RequestResponseOption) (interface{}, error)
}

type Identifier interface {
	Identifier() interface{}
}

type Iterable interface {
	Iterable() bool
}

type Iter struct {
	Query
	ListParams   ListParamsContainer
	Params       []interface{}
	CurrentIndex int
	Page         int64
	Values       *[]interface{}
	Cursor       string
	Error        error
	OnPageError
	requestResponseOptions []RequestResponseOption
	// listing holds the listing's parameters from before the first page
	// request. Reload starts from it rather than from parameters that a running
	// traversal is updating.
	listing       *reloadParams
	recordListing sync.Once
}

// reloadParams are the parameters of a reloaded traversal: the query its
// listing sent before paging, with this traversal's own paging state.
type reloadParams struct {
	ListParams
	query url.Values // every exported parameter except per_page and cursor
	err   error      // from exporting the listing's parameters
}

// values copies the recorded query, so changes to one request's values cannot
// reach later requests.
func (p *reloadParams) values() (url.Values, error) {
	if p.err != nil {
		return nil, p.err
	}
	values := make(url.Values, len(p.query))
	for key, value := range p.query {
		values[key] = slices.Clone(value)
	}
	return values, nil
}

// Err returns the error, if any,
// that caused the Iter to stop.
// It must be inspected
// after Next returns false.
func (i *Iter) Err() error {
	return i.Error
}

func (i *Iter) Current() interface{} {
	return (*i.Values)[i.CurrentIndex]
}

func (i *Iter) GetParams() *ListParams {
	return i.ListParams.GetListParams()
}

func (i *Iter) ExportParams() (lib.ExportValues, error) {
	p := lib.Params{Params: i.GetParams()}
	paramValues, err := p.ToValues()
	if err != nil {
		return lib.ExportValues{}, err
	}
	var listParamValues url.Values
	if reload, ok := i.ListParams.(*reloadParams); ok {
		listParamValues, err = reload.values()
	} else {
		listParamValues, err = lib.Params{Params: i.ListParams}.ToValues()
	}
	if err != nil {
		return lib.ExportValues{}, err
	}

	for key, value := range paramValues {
		listParamValues.Set(key, value[0])
	}

	return lib.ExportValues{Values: listParamValues}, nil
}

func (i *Iter) GetPage() bool {
	i.recordedListing() // before the first request changes the cursor
	for {
		if i.GetParams().MaxPages != 0 && i.Page >= i.GetParams().MaxPages {
			return false
		}

		i.CurrentIndex = 0
		i.Page += 1

		if i.Values != nil && i.Cursor == "" {
			return false
		}
		previousCursor := i.GetParams().Cursor
		params, err := i.ExportParams()
		if err != nil {
			i.Error = err
			return false
		}
		i.Values, i.Cursor, i.Error = i.Query(params, i.requestResponseOptions...)
		i.SetCursor(i.Cursor)
		if i.Error != nil && i.OnPageError != nil {
			i.Values, i.Error = i.OnPageError(i.Error)
		}
		if i.Error == nil && len(*i.Values) == 0 && i.Cursor != "" && i.Cursor == previousCursor {
			i.Error = errors.New("pagination cursor did not advance after an empty page")
		}
		if i.Error != nil || len(*i.Values) != 0 || i.Cursor == "" {
			return i.Error == nil && len(*i.Values) != 0
		}
	}
}

func (i *Iter) EOFPage() bool {
	if i.Values == nil {
		return false
	}
	return len(*i.Values) == i.CurrentIndex+1
}

func (i *Iter) Paging() bool {
	return true
}

func (i *Iter) GetCursor() string {
	return i.GetParams().Cursor
}

func (i *Iter) SetCursor(cursor string) {
	i.GetParams().Cursor = cursor
	i.Cursor = cursor
}

// Next iterates the results in i.Current() or i.`ResourceName`().
// It returns true until there are no results remaining.
// To adjust the number of results set ListParams.PerPage.
// To have it auto-paginate set ListParams.MaxPages, default is 1.
//
// To iterate over all results use the following pattern.
//
//	for i.Next() {
//	  i.Current()
//	}
func (i *Iter) Next() bool {
	if i.Values == nil {
		return i.GetPage() && len(*i.Values) > 0
	} else if len(*i.Values) > i.CurrentIndex+1 {
		i.CurrentIndex += 1
		return true
	}

	if i.EOFPage() {
		return i.GetPage()
	}

	return false
}

func (i *Iter) NextPage() bool {
	return i.Cursor != ""
}

// Reload returns a new iterator for the same listing that starts again from
// the first page, even if i started from a cursor. It keeps the listing's
// parameters as they were before i requested its first page, along with its
// page limits, error handler, and request options. opts run after those
// options, so WithContext replaces the listing's context; without it, the
// original context still applies. Reload does not read or change i's
// traversal state, so it can run while another goroutine is still in i.Next,
// as when a caller cancels a traversal and reloads it without waiting.
func (i *Iter) Reload(opts ...RequestResponseOption) IterI {
	params := *i.recordedListing()
	return &Iter{
		Query:                  i.Query,
		ListParams:             &params,
		Params:                 i.Params,
		OnPageError:            i.OnPageError,
		requestResponseOptions: slices.Concat(i.requestResponseOptions, opts),
	}
}

// recordedListing returns the listing's parameters from before the first
// page request, recording them on first use.
func (i *Iter) recordedListing() *reloadParams {
	i.recordListing.Do(func() {
		exported, err := i.ExportParams()
		paging := *i.GetParams()
		paging.Cursor = ""
		delete(exported.Values, "cursor")
		delete(exported.Values, "per_page")
		i.listing = &reloadParams{ListParams: paging, query: exported.Values, err: err}
	})
	return i.listing
}
