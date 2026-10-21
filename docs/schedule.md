# Schedule

`files_sdk.Schedule` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/schedule`, named `schedule`.

Each operation has a package function, which uses the default configuration, including the `FILES_API_KEY`
environment variable, and a `Client` method of the same name, which uses the client's embedded `Config`. The
examples run inside a function that returns an `error`, with these imports:

```go
import (
	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/schedule"
)
```

## Operations

### List

`GET /schedules`

List Schedules

```go
func (c *Client) List(params files_sdk.ScheduleListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
func List(params files_sdk.ScheduleListParams, opts ...files_sdk.RequestResponseOption) (*Iter, error)
```

Returns this package's `*Iter`, which requests pages as iteration reaches them. `Iter.All()` yields `files_sdk.Schedule` values, each with a nil error, followed by the error that stops the listing, if any.

Parameters, `files_sdk.ScheduleListParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `SortBy` | `interface{}` | `sort_by` | no | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name`. |
| `PerPage` | `int64` | `per_page` | no | From the embedded `files_sdk.ListParams`: requests the number of results per page. Zero uses the API default. |
| `Cursor` | `string` | `cursor` | no | From the embedded `files_sdk.ListParams`: starts the listing at an API cursor. Empty starts at the first page. |
| `MaxPages` | `int64` | not sent | no | From the embedded `files_sdk.ListParams`. Its tags send nothing to the API. |

Example:

```go
it, err := schedule.List(files_sdk.ScheduleListParams{})
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

`GET /schedules/{id}`

Show Schedule

```go
func (c *Client) Find(params files_sdk.ScheduleFindParams, opts ...files_sdk.RequestResponseOption) (schedule files_sdk.Schedule, err error)
func Find(params files_sdk.ScheduleFindParams, opts ...files_sdk.RequestResponseOption) (schedule files_sdk.Schedule, err error)
```

Parameters, `files_sdk.ScheduleFindParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Schedule ID. |

Example:

```go
result, err := schedule.Find(files_sdk.ScheduleFindParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Create

`POST /schedules`

Create Schedule

```go
func (c *Client) Create(params files_sdk.ScheduleCreateParams, opts ...files_sdk.RequestResponseOption) (schedule files_sdk.Schedule, err error)
func Create(params files_sdk.ScheduleCreateParams, opts ...files_sdk.RequestResponseOption) (schedule files_sdk.Schedule, err error)
```

Parameters, `files_sdk.ScheduleCreateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Name` | `string` | `name` | yes | Schedule name. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | yes | 0-based weekdays used by the Schedule. 0 is Sunday. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | yes | Times of day in HH:MM format (24-hour). |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | no | Time zone for scheduled times. If not set, times are interpreted as UTC. |
| `HolidayRegion` | `string` | `holiday_region` | no | Optional holiday region on which linked resources do not run. |

Example:

```go
result, err := schedule.Create(files_sdk.ScheduleCreateParams{
	Name:               "Weekday overnight",
	ScheduleDaysOfWeek: []int64{1, 2, 3, 4, 5},
	ScheduleTimesOfDay: []string{"01:00"},
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Update

`PATCH /schedules/{id}`

Update Schedule

```go
func (c *Client) Update(params files_sdk.ScheduleUpdateParams, opts ...files_sdk.RequestResponseOption) (schedule files_sdk.Schedule, err error)
func Update(params files_sdk.ScheduleUpdateParams, opts ...files_sdk.RequestResponseOption) (schedule files_sdk.Schedule, err error)
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (schedule files_sdk.Schedule, err error)
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (schedule files_sdk.Schedule, err error)
```

`UpdateWithMap` takes the parameters as a map keyed by their API names, path parameters
included. Unlike optional struct fields, explicit zero values in the map are sent.

Parameters, `files_sdk.ScheduleUpdateParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Schedule ID. |
| `Name` | `string` | `name` | no | Schedule name. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | no | 0-based weekdays used by the Schedule. 0 is Sunday. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | no | Times of day in HH:MM format (24-hour). |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | no | Time zone for scheduled times. If not set, times are interpreted as UTC. |
| `HolidayRegion` | `string` | `holiday_region` | no | Optional holiday region on which linked resources do not run. |

Example:

```go
result, err := schedule.Update(files_sdk.ScheduleUpdateParams{
	Id: 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

With a map:

```go
result, err := schedule.UpdateWithMap(map[string]interface{}{
	"id": 1,
})
if err != nil {
	return err
}
fmt.Println(result)
```

### Delete

`DELETE /schedules/{id}`

Delete Schedule

```go
func (c *Client) Delete(params files_sdk.ScheduleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
func Delete(params files_sdk.ScheduleDeleteParams, opts ...files_sdk.RequestResponseOption) (err error)
```

Parameters, `files_sdk.ScheduleDeleteParams`:

| Field | Go type | API name | Required | Description |
| --- | --- | --- | --- | --- |
| `Id` | `int64` | `id`, in the path | yes | Schedule ID. |

Example:

```go
err := schedule.Delete(files_sdk.ScheduleDeleteParams{
	Id: 1,
})
if err != nil {
	return err
}
```

## Model

`files_sdk.Schedule` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Schedule ID. |
| `Name` | `string` | `name` | Schedule name. |
| `ScheduleDaysOfWeek` | `[]int64` | `schedule_days_of_week` | 0-based weekdays used by the Schedule. 0 is Sunday. |
| `ScheduleTimesOfDay` | `[]string` | `schedule_times_of_day` | Times of day in HH:MM format (24-hour). |
| `ScheduleTimeZone` | `string` | `schedule_time_zone` | Time zone for scheduled times. If not set, times are interpreted as UTC. |
| `HolidayRegion` | `string` | `holiday_region` | Optional holiday region on which linked resources do not run. |
| `HumanReadableSchedule` | `string` | `human_readable_schedule` | Human-readable Schedule description. |
| `CreatedAt` | `*time.Time` | `created_at` | Creation time. |
| `UpdatedAt` | `*time.Time` | `updated_at` | Last update time. |
