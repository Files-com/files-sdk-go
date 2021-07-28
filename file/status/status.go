package status

type Status string

const (
	Queued      = Status("queued")
	Downloading = Status("downloading")
	Uploading   = Status("uploading")
	Skipped     = Status("skipped")
	Complete    = Status("complete")
	Errored     = Status("errored")
	Canceled    = Status("canceled")
)

func (e Status) String() string {
	return string(e)
}

func (e Status) Queued() bool {
	return e == Queued
}

func (e Status) Downloading() bool {
	return e == Downloading
}

func (e Status) Uploading() bool {
	return e == Uploading
}

func (e Status) Completed() bool {
	return e == Complete
}

func (e Status) Skipped() bool {
	return e == Skipped
}

func (e Status) Errored() bool {
	return e == Errored
}

func (e Status) Invalid() bool {
	return e.Errored() || e.Canceled() || e.Skipped()
}

func (e Status) Valid() bool {
	return e.Queued() || e.Downloading() || e.Uploading() || e.Completed()
}

func (e Status) Canceled() bool {
	return e == Canceled
}

func (e Status) Ended() bool {
	return e.Completed() || e.Errored() || e.Skipped() || e.Canceled()
}

func (e Status) Compare(s Status) bool {
	return e == s
}

func (e Status) Type() Status {
	return e
}

type IStatus interface {
	String() string
	Queued() bool
	Downloading() bool
	Uploading() bool
	Completed() bool
	Skipped() bool
	Errored() bool
	Invalid() bool
	Valid() bool
	Canceled() bool
	Ended() bool
	Compare(Status) bool
	Type() Status
}
