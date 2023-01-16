package status

type Status struct {
	Name  string
	Value int
}

var (
	Null        = Status{"", -1}
	Indexed     = Status{"indexed", 0}
	Retrying    = Status{"retrying", 0}
	Queued      = Status{"queued", 1}
	Downloading = Status{"downloading", 2}
	Uploading   = Status{"uploading", 2}
	Skipped     = Status{"skipped", 3}
	Ignored     = Status{"ignored", 3}
	Complete    = Status{"complete", 4}
	Canceled    = Status{"canceled", 5}
	Errored     = Status{"errored", 5}

	Included = []Status{Indexed, Queued, Retrying, Downloading, Uploading, Complete, Canceled, Errored}
	Excluded = []Status{Skipped, Ignored}
	Valid    = []Status{Indexed, Queued, Retrying, Downloading, Uploading, Complete}
	Invalid  = []Status{Null, Canceled, Errored, Skipped, Ignored}
	Running  = []Status{Downloading, Uploading}
	Ended    = []Status{Complete, Canceled, Errored, Skipped, Ignored}
)

func (e Status) String() string {
	return e.Name
}

func (e Status) Has(statuses ...Status) bool {
	return e.Any(statuses...)
}

func (e Status) Is(statuses ...Status) bool {
	return e.Any(statuses...)
}

func (e Status) IsNot(statuses ...Status) bool {
	return !e.Any(statuses...)
}

func (e Status) is(status Status) bool {
	return e.Name == status.Name
}

func (e Status) Any(statuses ...Status) bool {
	if len(statuses) == 0 {
		return true
	}
	for _, status := range statuses {
		if e.is(status) {
			return true
		}
	}
	return false
}

func SetStatus(old Status, new Status, err error) (Status, bool) {
	var setError bool
	if err != nil || new.Is(Retrying) {
		setError = true
	}
	if old.Is(Errored) && new.Is(Running...) {
		new = old
	}

	return new, setError
}
