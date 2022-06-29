package files_sdk

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type File struct {
	Path             string     `json:"path,omitempty"`
	DisplayName      string     `json:"display_name,omitempty"`
	Type             string     `json:"type,omitempty"`
	Size             int64      `json:"size,omitempty"`
	Mtime            *time.Time `json:"mtime,omitempty"`
	ProvidedMtime    *time.Time `json:"provided_mtime,omitempty"`
	Crc32            string     `json:"crc32,omitempty"`
	Md5              string     `json:"md5,omitempty"`
	MimeType         string     `json:"mime_type,omitempty"`
	Region           string     `json:"region,omitempty"`
	Permissions      string     `json:"permissions,omitempty"`
	SubfoldersLocked *bool      `json:"subfolders_locked?,omitempty"`
	DownloadUri      string     `json:"download_uri,omitempty"`
	PriorityColor    string     `json:"priority_color,omitempty"`
	PreviewId        int64      `json:"preview_id,omitempty"`
	Preview          Preview    `json:"preview,omitempty"`
	Action           string     `json:"action,omitempty"`
	Length           int64      `json:"length,omitempty"`
	MkdirParents     *bool      `json:"mkdir_parents,omitempty"`
	Part             int64      `json:"part,omitempty"`
	Parts            int64      `json:"parts,omitempty"`
	Ref              string     `json:"ref,omitempty"`
	Restart          int64      `json:"restart,omitempty"`
	Structure        string     `json:"structure,omitempty"`
	WithRename       *bool      `json:"with_rename,omitempty"`
}

type FileCollection []File

type EtagsParam struct {
	Etag string `url:"etag,omitempty" json:"etag,omitempty"`
	Part string `url:"part,omitempty" json:"part,omitempty"`
}

// Download file
type FileDownloadParams struct {
	Path              string               `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Action            string               `url:"action,omitempty" required:"false" json:"action,omitempty"`
	PreviewSize       string               `url:"preview_size,omitempty" required:"false" json:"preview_size,omitempty"`
	WithPreviews      *bool                `url:"with_previews,omitempty" required:"false" json:"with_previews,omitempty"`
	WithPriorityColor *bool                `url:"with_priority_color,omitempty" required:"false" json:"with_priority_color,omitempty"`
	Writer            io.Writer            `url:"-,omitempty" required:"false" json:"-,omitempty"`
	OnDownload        func(*http.Response) `url:"-,omitempty" required:"false" json:"-,omitempty"`
}

type FileCreateParams struct {
	Path          string       `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Action        string       `url:"action,omitempty" required:"false" json:"action,omitempty"`
	EtagsParam    []EtagsParam `url:"etags,omitempty" required:"false" json:"etags,omitempty"`
	Length        int64        `url:"length,omitempty" required:"false" json:"length,omitempty"`
	MkdirParents  *bool        `url:"mkdir_parents,omitempty" required:"false" json:"mkdir_parents,omitempty"`
	Part          int64        `url:"part,omitempty" required:"false" json:"part,omitempty"`
	Parts         int64        `url:"parts,omitempty" required:"false" json:"parts,omitempty"`
	ProvidedMtime *time.Time   `url:"provided_mtime,omitempty" required:"false" json:"provided_mtime,omitempty"`
	Ref           string       `url:"ref,omitempty" required:"false" json:"ref,omitempty"`
	Restart       int64        `url:"restart,omitempty" required:"false" json:"restart,omitempty"`
	Size          int64        `url:"size,omitempty" required:"false" json:"size,omitempty"`
	Structure     string       `url:"structure,omitempty" required:"false" json:"structure,omitempty"`
	WithRename    *bool        `url:"with_rename,omitempty" required:"false" json:"with_rename,omitempty"`
}

type FileUpdateParams struct {
	Path          string     `url:"-,omitempty" required:"true" json:"-,omitempty"`
	ProvidedMtime *time.Time `url:"provided_mtime,omitempty" required:"false" json:"provided_mtime,omitempty"`
	PriorityColor string     `url:"priority_color,omitempty" required:"false" json:"priority_color,omitempty"`
}

type FileDeleteParams struct {
	Path      string `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Recursive *bool  `url:"recursive,omitempty" required:"false" json:"recursive,omitempty"`
}

type FileFindParams struct {
	Path              string `url:"-,omitempty" required:"true" json:"-,omitempty"`
	PreviewSize       string `url:"preview_size,omitempty" required:"false" json:"preview_size,omitempty"`
	WithPreviews      *bool  `url:"with_previews,omitempty" required:"false" json:"with_previews,omitempty"`
	WithPriorityColor *bool  `url:"with_priority_color,omitempty" required:"false" json:"with_priority_color,omitempty"`
}

// Copy file/folder
type FileCopyParams struct {
	Path        string `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Destination string `url:"destination,omitempty" required:"true" json:"destination,omitempty"`
	Structure   *bool  `url:"structure,omitempty" required:"false" json:"structure,omitempty"`
}

// Move file/folder
type FileMoveParams struct {
	Path        string `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Destination string `url:"destination,omitempty" required:"true" json:"destination,omitempty"`
}

// Begin file upload
type FileBeginUploadParams struct {
	Path         string `url:"-,omitempty" required:"true" json:"-,omitempty"`
	MkdirParents *bool  `url:"mkdir_parents,omitempty" required:"false" json:"mkdir_parents,omitempty"`
	Part         int64  `url:"part,omitempty" required:"false" json:"part,omitempty"`
	Parts        int64  `url:"parts,omitempty" required:"false" json:"parts,omitempty"`
	Ref          string `url:"ref,omitempty" required:"false" json:"ref,omitempty"`
	Restart      int64  `url:"restart,omitempty" required:"false" json:"restart,omitempty"`
	Size         int64  `url:"size,omitempty" required:"false" json:"size,omitempty"`
	WithRename   *bool  `url:"with_rename,omitempty" required:"false" json:"with_rename,omitempty"`
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
