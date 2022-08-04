package files_sdk

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type File struct {
	Path             string     `json:"path,omitempty" path:"path"`
	DisplayName      string     `json:"display_name,omitempty" path:"display_name"`
	Type             string     `json:"type,omitempty" path:"type"`
	Size             int64      `json:"size,omitempty" path:"size"`
	CreatedAt        *time.Time `json:"created_at,omitempty" path:"created_at"`
	Mtime            *time.Time `json:"mtime,omitempty" path:"mtime"`
	ProvidedMtime    *time.Time `json:"provided_mtime,omitempty" path:"provided_mtime"`
	Crc32            string     `json:"crc32,omitempty" path:"crc32"`
	Md5              string     `json:"md5,omitempty" path:"md5"`
	MimeType         string     `json:"mime_type,omitempty" path:"mime_type"`
	Region           string     `json:"region,omitempty" path:"region"`
	Permissions      string     `json:"permissions,omitempty" path:"permissions"`
	SubfoldersLocked *bool      `json:"subfolders_locked?,omitempty" path:"subfolders_locked?"`
	DownloadUri      string     `json:"download_uri,omitempty" path:"download_uri"`
	PriorityColor    string     `json:"priority_color,omitempty" path:"priority_color"`
	PreviewId        int64      `json:"preview_id,omitempty" path:"preview_id"`
	Preview          Preview    `json:"preview,omitempty" path:"preview"`
	Action           string     `json:"action,omitempty" path:"action"`
	Length           int64      `json:"length,omitempty" path:"length"`
	MkdirParents     *bool      `json:"mkdir_parents,omitempty" path:"mkdir_parents"`
	Part             int64      `json:"part,omitempty" path:"part"`
	Parts            int64      `json:"parts,omitempty" path:"parts"`
	Ref              string     `json:"ref,omitempty" path:"ref"`
	Restart          int64      `json:"restart,omitempty" path:"restart"`
	Structure        string     `json:"structure,omitempty" path:"structure"`
	WithRename       *bool      `json:"with_rename,omitempty" path:"with_rename"`
}

type FileCollection []File

type EtagsParam struct {
	Etag string `url:"etag,omitempty" json:"etag,omitempty" path:"etag"`
	Part string `url:"part,omitempty" json:"part,omitempty" path:"part"`
}

// Download file
type FileDownloadParams struct {
	Path              string               `url:"-,omitempty" required:"true" json:"-,omitempty" path:"path"`
	Action            string               `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	PreviewSize       string               `url:"preview_size,omitempty" required:"false" json:"preview_size,omitempty" path:"preview_size"`
	WithPreviews      *bool                `url:"with_previews,omitempty" required:"false" json:"with_previews,omitempty" path:"with_previews"`
	WithPriorityColor *bool                `url:"with_priority_color,omitempty" required:"false" json:"with_priority_color,omitempty" path:"with_priority_color"`
	Writer            io.Writer            `url:"-,omitempty" required:"false" json:"-,omitempty"`
	OnDownload        func(*http.Response) `url:"-,omitempty" required:"false" json:"-,omitempty"`
}

type FileCreateParams struct {
	Path          string       `url:"-,omitempty" required:"true" json:"-,omitempty" path:"path"`
	Action        string       `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	EtagsParam    []EtagsParam `url:"etags,omitempty" required:"false" json:"etags,omitempty" path:"etags"`
	Length        int64        `url:"length,omitempty" required:"false" json:"length,omitempty" path:"length"`
	MkdirParents  *bool        `url:"mkdir_parents,omitempty" required:"false" json:"mkdir_parents,omitempty" path:"mkdir_parents"`
	Part          int64        `url:"part,omitempty" required:"false" json:"part,omitempty" path:"part"`
	Parts         int64        `url:"parts,omitempty" required:"false" json:"parts,omitempty" path:"parts"`
	ProvidedMtime *time.Time   `url:"provided_mtime,omitempty" required:"false" json:"provided_mtime,omitempty" path:"provided_mtime"`
	Ref           string       `url:"ref,omitempty" required:"false" json:"ref,omitempty" path:"ref"`
	Restart       int64        `url:"restart,omitempty" required:"false" json:"restart,omitempty" path:"restart"`
	Size          int64        `url:"size,omitempty" required:"false" json:"size,omitempty" path:"size"`
	Structure     string       `url:"structure,omitempty" required:"false" json:"structure,omitempty" path:"structure"`
	WithRename    *bool        `url:"with_rename,omitempty" required:"false" json:"with_rename,omitempty" path:"with_rename"`
}

type FileUpdateParams struct {
	Path          string     `url:"-,omitempty" required:"true" json:"-,omitempty" path:"path"`
	ProvidedMtime *time.Time `url:"provided_mtime,omitempty" required:"false" json:"provided_mtime,omitempty" path:"provided_mtime"`
	PriorityColor string     `url:"priority_color,omitempty" required:"false" json:"priority_color,omitempty" path:"priority_color"`
}

type FileDeleteParams struct {
	Path      string `url:"-,omitempty" required:"true" json:"-,omitempty" path:"path"`
	Recursive *bool  `url:"recursive,omitempty" required:"false" json:"recursive,omitempty" path:"recursive"`
}

type FileFindParams struct {
	Path              string `url:"-,omitempty" required:"true" json:"-,omitempty" path:"path"`
	PreviewSize       string `url:"preview_size,omitempty" required:"false" json:"preview_size,omitempty" path:"preview_size"`
	WithPreviews      *bool  `url:"with_previews,omitempty" required:"false" json:"with_previews,omitempty" path:"with_previews"`
	WithPriorityColor *bool  `url:"with_priority_color,omitempty" required:"false" json:"with_priority_color,omitempty" path:"with_priority_color"`
}

// Copy file/folder
type FileCopyParams struct {
	Path        string `url:"-,omitempty" required:"true" json:"-,omitempty" path:"path"`
	Destination string `url:"destination,omitempty" required:"true" json:"destination,omitempty" path:"destination"`
	Structure   *bool  `url:"structure,omitempty" required:"false" json:"structure,omitempty" path:"structure"`
}

// Move file/folder
type FileMoveParams struct {
	Path        string `url:"-,omitempty" required:"true" json:"-,omitempty" path:"path"`
	Destination string `url:"destination,omitempty" required:"true" json:"destination,omitempty" path:"destination"`
}

// Begin file upload
type FileBeginUploadParams struct {
	Path         string `url:"-,omitempty" required:"true" json:"-,omitempty" path:"path"`
	MkdirParents *bool  `url:"mkdir_parents,omitempty" required:"false" json:"mkdir_parents,omitempty" path:"mkdir_parents"`
	Part         int64  `url:"part,omitempty" required:"false" json:"part,omitempty" path:"part"`
	Parts        int64  `url:"parts,omitempty" required:"false" json:"parts,omitempty" path:"parts"`
	Ref          string `url:"ref,omitempty" required:"false" json:"ref,omitempty" path:"ref"`
	Restart      int64  `url:"restart,omitempty" required:"false" json:"restart,omitempty" path:"restart"`
	Size         int64  `url:"size,omitempty" required:"false" json:"size,omitempty" path:"size"`
	WithRename   *bool  `url:"with_rename,omitempty" required:"false" json:"with_rename,omitempty" path:"with_rename"`
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
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = File(v)
	return nil
}

func (f *FileCollection) UnmarshalJSON(data []byte) error {
	type files FileCollection
	var v files
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
