package lib

import (
	"net/url"
)

type ListParams struct {
	Page     int64  `json:"page,omitempty" url:"page,omitempty" required:"false"`
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

func (p *ListParams) Set(page int64, perPage int64, cursor string, maxPages int64) {
	p.Page = page
	p.PerPage = perPage
	p.Cursor = cursor
	p.MaxPages = maxPages
}

type Query func() (*[]interface{}, string, error)

type Iter struct {
	Query
	ListParams   ListParamsContainer
	Params       []interface{}
	CurrentIndex int
	Values       *[]interface{}
	Cursor       string
	Error        error
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

func (i *Iter) ExportParams() (url.Values, error) {
	paramValues, err := ExportParams(i.GetParams())
	if err != nil {
		return paramValues, err
	}
	listParamValues, _ := ExportParams(i.ListParams)

	for key, value := range paramValues {
		listParamValues.Set(key, value[0])
	}

	if i.GetCursor() != "" {
		listParamValues.Del("page")
	}

	return listParamValues, nil
}

func (i *Iter) GetPage() bool {
	if i.GetParams().MaxPages != 0 && i.GetParams().Page == i.GetParams().MaxPages {
		return false
	}

	i.GetParams().Page += 1
	if i.GetParams().Page == 2 && i.Cursor == "" {
		return false
	}
	i.Values, i.Cursor, i.Error = i.Query()
	i.SetCursor(i.Cursor)
	return i.Error == nil && len(*i.Values) != 0
}

func (i *Iter) EOFPage() bool {
	return len(*i.Values) == i.CurrentIndex+1
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
//   for i.Next() {
//     i.Current()
//   }
func (i *Iter) Next() bool {
	if i.Values == nil {
		return i.GetPage() && len(*i.Values) > 0
	} else if len(*i.Values) > i.CurrentIndex+1 {
		i.CurrentIndex += 1
		return true
	}

	if i.EOFPage() {
		i.CurrentIndex = 0
		return i.GetPage()
	}

	return false
}
