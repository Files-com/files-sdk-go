package files_sdk

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type File struct {
	Path             string    `json:"path,omitempty"`
	DisplayName      string    `json:"display_name,omitempty"`
	Type             string    `json:"type,omitempty"`
	Size             int64     `json:"size,omitempty"`
	Mtime            time.Time `json:"mtime,omitempty"`
	ProvidedMtime    time.Time `json:"provided_mtime,omitempty"`
	Crc32            string    `json:"crc32,omitempty"`
	Md5              string    `json:"md5,omitempty"`
	MimeType         string    `json:"mime_type,omitempty"`
	Region           string    `json:"region,omitempty"`
	Permissions      string    `json:"permissions,omitempty"`
	SubfoldersLocked *bool     `json:"subfolders_locked?,omitempty"`
	DownloadUri      string    `json:"download_uri,omitempty"`
	PriorityColor    string    `json:"priority_color,omitempty"`
	PreviewId        int64     `json:"preview_id,omitempty"`
	Preview          string    `json:"preview,omitempty"`
	Action           string    `json:"action,omitempty"`
	Length           int64     `json:"length,omitempty"`
	MkdirParents     *bool     `json:"mkdir_parents,omitempty"`
	Part             int64     `json:"part,omitempty"`
	Parts            int64     `json:"parts,omitempty"`
	Ref              string    `json:"ref,omitempty"`
	Restart          int64     `json:"restart,omitempty"`
	Structure        string    `json:"structure,omitempty"`
	WithRename       *bool     `json:"with_rename,omitempty"`
}

type FileCollection []File

type EtagsParam struct {
	Etag string `url:"etag,omitempty"`
	Part string `url:"part,omitempty"`
}

// Download file
type FileDownloadParams struct {
	Path              string               `url:"-,omitempty" required:"true"`
	Action            string               `url:"action,omitempty" required:"false"`
	PreviewSize       string               `url:"preview_size,omitempty" required:"false"`
	WithPreviews      *bool                `url:"with_previews,omitempty" required:"false"`
	WithPriorityColor *bool                `url:"with_priority_color,omitempty" required:"false"`
	Writer            io.Writer            `url:"-,omitempty" required:"false"`
	OnDownload        func(*http.Response) `url:"-,omitempty" required:"false"`
}

type FileCreateParams struct {
	Path          string       `url:"-,omitempty" required:"true"`
	Action        string       `url:"action,omitempty" required:"false"`
	EtagsParam    []EtagsParam `url:"etags,omitempty" required:"false"`
	Length        int64        `url:"length,omitempty" required:"false"`
	MkdirParents  *bool        `url:"mkdir_parents,omitempty" required:"false"`
	Part          int64        `url:"part,omitempty" required:"false"`
	Parts         int64        `url:"parts,omitempty" required:"false"`
	ProvidedMtime time.Time    `url:"provided_mtime,omitempty" required:"false"`
	Ref           string       `url:"ref,omitempty" required:"false"`
	Restart       int64        `url:"restart,omitempty" required:"false"`
	Size          int64        `url:"size,omitempty" required:"false"`
	Structure     string       `url:"structure,omitempty" required:"false"`
	WithRename    *bool        `url:"with_rename,omitempty" required:"false"`
}

type FileUpdateParams struct {
	Path          string    `url:"-,omitempty" required:"true"`
	ProvidedMtime time.Time `url:"provided_mtime,omitempty" required:"false"`
	PriorityColor string    `url:"priority_color,omitempty" required:"false"`
}

type FileDeleteParams struct {
	Path      string `url:"-,omitempty" required:"true"`
	Recursive *bool  `url:"recursive,omitempty" required:"false"`
}

func (f File) ToFolder() (Folder, error) {
	bodyBytes, err := json.Marshal(f)
	if err != nil {
		return Folder{}, err
	}
	folder := Folder{}
	folder.UnmarshalJSON(bodyBytes)
	return folder, nil
}

func (f *File) UnmarshalJSON(data []byte) error {
	type file File
	var v file
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = File(v)
	return nil
}

func (f *FileCollection) UnmarshalJSON(data []byte) error {
	type files []File
	var v files
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FileCollection(v)
	return nil
}

func (f *FileCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
