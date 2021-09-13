package direction

type Type struct {
	name   string
	Symbol string
}

var (
	DownloadType = Type{name: "download", Symbol: "⬆"}
	UploadType   = Type{name: "upload", Symbol: "⬆"}
)

func (t *Type) Name() string {
	return t.name
}
