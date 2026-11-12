package files_sdk

import (
	"errors"
	"fmt"
	"iter"
	"net/url"
	"slices"
	"sync"

	"github.com/Files-com/files-sdk-go/v3/lib"
)

// ListParams controls cursor-based pagination. A zero MaxPages allows all pages;
// a positive value limits the number of pages requested.
type ListParams struct {
	// PerPage requests the number of results per page. Zero uses the API default.
	PerPage int64 `json:"per_page,omitempty" url:"per_page,omitempty" required:"false"`
	// Cursor starts the listing at an API cursor. Empty starts at the first page.
	Cursor   string `json:"cursor,omitempty" url:"cursor,omitempty" required:"false"`
	MaxPages int64  `json:"-" url:"-"`
}

// ListParamsContainer exposes paging settings. Generated list parameter
// structs implement it by embedding ListParams.
type ListParamsContainer interface {
	GetListParams() *ListParams
}

// GetListParams returns p. Structs embedding ListParams inherit this method.
func (p *ListParams) GetListParams() *ListParams {
	return p
}

// OnPageError handles a page-fetch error. It can return replacement results and
// a nil error to continue iteration, or an error to stop.
type OnPageError func(error) (*[]interface{}, error)

// Query fetches one page and returns its resources, the next cursor, and any error.
type Query func(params lib.Values, opts ...RequestResponseOption) (*[]interface{}, string, error)

// IterI traverses resources. Call Current only after Next returns true, and
// check Err after Next returns false.
type IterI interface {
	Next() bool
	Current() interface{}
	Err() error
}

var _ IterI = (*Iter)(nil)

// TypedIterI adds typed resource access to an iterator.
type TypedIterI[T any] interface {
	Next() bool
	Current() interface{}
	Resource() T
	Err() error
}

// IterPagingI reports whether the current resource is the last one on its page.
type IterPagingI interface {
	IterI
	EOFPage() bool
}

var _ IterPagingI = (*Iter)(nil)

// ResourceIterator creates a listing for a resource identifier.
type ResourceIterator interface {
	Iterate(interface{}, ...RequestResponseOption) (IterI, error)
}

// ReloadIterator starts a fresh traversal of the same listing.
type ReloadIterator interface {
	Reload(opts ...RequestResponseOption) IterI
}

var _ ReloadIterator = (*Iter)(nil)

// ResourceLoader fetches a resource by its identifier.
type ResourceLoader interface {
	LoadResource(interface{}, ...RequestResponseOption) (interface{}, error)
}

// Identifier exposes an API resource ID or path.
type Identifier interface {
	Identifier() interface{}
}

// Iterable reports whether a resource can be listed for child resources.
type Iterable interface {
	Iterable() bool
}

// Iter traverses API results, fetching additional pages on demand. Use the
// generated client listing methods to create an iterator.
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
	// ended records that the last GetPage call ended the traversal, so IterAll
	// yields nothing more.
	ended bool
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

// Err returns the error that stopped iteration, or nil if none occurred.
// Check it after Next returns false.
func (i *Iter) Err() error {
	return i.Error
}

// Current returns the resource selected by the last successful Next call.
// Call it only after Next returns true.
func (i *Iter) Current() interface{} {
	return (*i.Values)[i.CurrentIndex]
}

// GetParams returns the paging settings used by the iterator.
func (i *Iter) GetParams() *ListParams {
	return i.ListParams.GetListParams()
}

// ExportParams combines the listing filters and current paging settings as query values.
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

// GetPage fetches the next nonempty page within MaxPages. It returns false when
// no results remain, the page limit is reached, or an error occurs. Check Err
// after a false return.
func (i *Iter) GetPage() bool {
	more := i.getPage()
	i.ended = !more
	return more
}

func (i *Iter) getPage() bool {
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

// EOFPage reports whether the current resource is the last one on the loaded page.
// It does not report whether the listing has more pages.
func (i *Iter) EOFPage() bool {
	if i.Values == nil {
		return false
	}
	return len(*i.Values) == i.CurrentIndex+1
}

// Paging reports that this iterator supports pagination.
func (i *Iter) Paging() bool {
	return true
}

// GetCursor returns the cursor that will be sent with the next page request.
func (i *Iter) GetCursor() string {
	return i.GetParams().Cursor
}

// SetCursor sets the cursor for the next page request.
func (i *Iter) SetCursor(cursor string) {
	i.GetParams().Cursor = cursor
	i.Cursor = cursor
}

// Next advances to the next resource, fetching pages as needed. It returns false
// when the listing ends, MaxPages is reached, or an error occurs. Check Err
// after it returns false. PerPage controls page size; MaxPages defaults to zero,
// which allows all pages.
//
//	for i.Next() {
//		// Process i.Current().
//	}
//	if err := i.Err(); err != nil {
//		// Handle the listing error.
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

// IterAll returns an iterator over the resources that later calls to i.Next
// would return, as values of type T. Generated listing iterators provide it as
// their All method:
//
//	for file, err := range it.All() {
//		if err != nil {
//			// Handle the listing error.
//			return err
//		}
//		// Process file.
//	}
//
// Creating the iterator requests nothing. Each page is requested when the loop
// reaches it, and breaking out of the loop requests nothing more. Each resource
// is yielded with a nil error. If a page request fails, the loop ends with one
// final yield of the zero value and the error, which Err also returns. The end
// of the listing, MaxPages, or an empty listing ends the loop without an extra
// yield. A resource that is not a T ends the loop the same way, with an error
// that Err also returns.
//
// The iterator shares i's position: after a break or calls to Next, a new loop
// continues with the next resource. Once the listing has ended or failed, it
// yields nothing; call the listing method again, or Reload, to start over. Do
// not call i.Next during the loop.
func IterAll[T any](i *Iter) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		if i.ended {
			return
		}
		var zero T
		for i.Next() {
			resource, ok := i.Current().(T)
			if !ok {
				i.Error = fmt.Errorf("files_sdk: listing resource is %T, not %T", i.Current(), zero)
				i.ended = true
				yield(zero, i.Error)
				return
			}
			if !yield(resource, nil) {
				return
			}
		}
		if err := i.Err(); err != nil {
			yield(zero, err)
		}
	}
}

// NextPage reports whether the API supplied a cursor for another page.
// It does not fetch that page or account for MaxPages.
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
