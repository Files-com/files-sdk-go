# HolidayCalendar

`files_sdk.HolidayCalendar` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/holidaycalendar`, named `holiday_calendar`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	holiday_calendar "github.com/Files-com/files-sdk-go/v3/holidaycalendar"
)
```

## Operations

### List

`GET /holiday_calendars`

List Holiday Calendars

```go
func (c *Client) List(params files_sdk.HolidayCalendarListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.HolidayCalendarListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.HolidayCalendar` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.HolidayCalendarListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are . |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := holiday_calendar.List(files_sdk.HolidayCalendarListParams{})
if err != nil {
	return err
}
for item, err := range it.All() {
	if err != nil {
		return err
	}
	fmt.Println(item)
}
```

### Find

`GET /holiday_calendars/{id}`

Show Holiday Calendar

```go
func (c *Client) Find(params files_sdk.HolidayCalendarFindParams, opts ...files_sdk.RequestResponseOption) (holidayCalendar files_sdk.HolidayCalendar, err error)
func Find(params files_sdk.HolidayCalendarFindParams, opts ...files_sdk.RequestResponseOption) (holidayCalendar files_sdk.HolidayCalendar, err error)
```

Parameters, `files_sdk.HolidayCalendarFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Holiday Calendar ID. |

Example:

```go
result, err := holiday_calendar.Find(files_sdk.HolidayCalendarFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /holiday_calendars`

Create Holiday Calendar

```go
func (c *Client) Create(params files_sdk.HolidayCalendarCreateParams, opts ...files_sdk.RequestResponseOption) (holidayCalendar files_sdk.HolidayCalendar, err error)
func Create(params files_sdk.HolidayCalendarCreateParams, opts ...files_sdk.RequestResponseOption) (holidayCalendar files_sdk.HolidayCalendar, err error)
```

Parameters, `files_sdk.HolidayCalendarCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Definition` | `interface{}` | `definition` | yes | Holiday rules for the calendar. |
| `Name` | `string` | `name` | yes | Holiday Calendar name. |

Example:

```go
// Set this required value before the call.
var definition interface{}
result, err := holiday_calendar.Create(files_sdk.HolidayCalendarCreateParams{
	Definition: definition,
	Name:       "Company Holidays",
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /holiday_calendars/{id}`

Update Holiday Calendar

```go
func (c *Client) Update(params files_sdk.HolidayCalendarUpdateParams, opts ...files_sdk.RequestResponseOption) (holidayCalendar files_sdk.HolidayCalendar, err error)
func Update(params files_sdk.HolidayCalendarUpdateParams, opts ...files_sdk.RequestResponseOption) (holidayCalendar files_sdk.HolidayCalendar, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (holidayCalendar files_sdk.HolidayCalendar, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (holidayCalendar files_sdk.HolidayCalendar, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.HolidayCalendarUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Holiday Calendar ID. |
| `Definition` | `interface{}` | `definition` | no | Holiday rules for the calendar. |
| `Name` | `string` | `name` | no | Holiday Calendar name. |

Example:

```go
result, err := holiday_calendar.Update(files_sdk.HolidayCalendarUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := holiday_calendar.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /holiday_calendars/{id}`

Delete Holiday Calendar

```go
func (c *Client) Delete(params files_sdk.HolidayCalendarDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.HolidayCalendarDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.HolidayCalendarDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Holiday Calendar ID. |

Example:

```go
err := holiday_calendar.Delete(files_sdk.HolidayCalendarDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.HolidayCalendar` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Holiday Calendar ID. Set a scheduled resource's `holiday_region` to `custom_` followed by this ID to make it skip the days in this calendar. |
| `Name` | `string` | `name` | Holiday Calendar name. |
| `Definition` | `interface{}` | `definition` | Holiday rules for the calendar. |
| `CreatedAt` | `*time.Time` | `created_at` | Creation time. |
| `UpdatedAt` | `*time.Time` | `updated_at` | Last update time. |
