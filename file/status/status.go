package status

type Status struct {
	Name  string
	Value int
}

var (
	Queued      = Status{"queued", 0}
	Downloading = Status{"downloading", 1}
	Uploading   = Status{"uploading", 1}
	Skipped     = Status{"skipped", 2}
	Complete    = Status{"complete", 3}
	Canceled    = Status{"canceled", 4}
	Errored     = Status{"errored", 4}
)

func (e Status) String() string {
	return e.Name
}

func (e Status) Queued() bool {
	return e.Name == Queued.Name
}

func (e Status) Downloading() bool {
	return e.Name == Downloading.Name
}

func (e Status) Uploading() bool {
	return e.Name == Uploading.Name
}

func (e Status) Completed() bool {
	return e.Name == Complete.Name
}

func (e Status) Skipped() bool {
	return e.Name == Skipped.Name
}

func (e Status) Errored() bool {
	return e.Name == Errored.Name
}

func (e Status) Running() bool {
	return e.Downloading() || e.Uploading()
}

func (e Status) Invalid() bool {
	return e.Errored() || e.Canceled() || e.Skipped()
}

func (e Status) Valid() bool {
	return e.Queued() || e.Downloading() || e.Uploading() || e.Completed()
}

func (e Status) Canceled() bool {
	return e.Name == Canceled.Name
}

func (e Status) Ended() bool {
	return e.Completed() || e.Errored() || e.Skipped() || e.Canceled()
}

func (e Status) Compare(s Status) bool {
	return e.Name == s.Name
}
