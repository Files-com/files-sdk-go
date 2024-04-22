package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type File struct {
	Path             string     `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	DisplayName      string     `json:"display_name,omitempty" path:"display_name,omitempty" url:"display_name,omitempty"`
	Type             string     `json:"type,omitempty" path:"type,omitempty" url:"type,omitempty"`
	Size             int64      `json:"size,omitempty" path:"size,omitempty" url:"size,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	Mtime            *time.Time `json:"mtime,omitempty" path:"mtime,omitempty" url:"mtime,omitempty"`
	ProvidedMtime    *time.Time `json:"provided_mtime,omitempty" path:"provided_mtime,omitempty" url:"provided_mtime,omitempty"`
	Crc32            string     `json:"crc32,omitempty" path:"crc32,omitempty" url:"crc32,omitempty"`
	Md5              string     `json:"md5,omitempty" path:"md5,omitempty" url:"md5,omitempty"`
	MimeType         string     `json:"mime_type,omitempty" path:"mime_type,omitempty" url:"mime_type,omitempty"`
	Region           string     `json:"region,omitempty" path:"region,omitempty" url:"region,omitempty"`
	Permissions      string     `json:"permissions,omitempty" path:"permissions,omitempty" url:"permissions,omitempty"`
	SubfoldersLocked *bool      `json:"subfolders_locked?,omitempty" path:"subfolders_locked?,omitempty" url:"subfolders_locked?,omitempty"`
	IsLocked         *bool      `json:"is_locked,omitempty" path:"is_locked,omitempty" url:"is_locked,omitempty"`
	DownloadUri      string     `json:"download_uri,omitempty" path:"download_uri,omitempty" url:"download_uri,omitempty"`
	PriorityColor    string     `json:"priority_color,omitempty" path:"priority_color,omitempty" url:"priority_color,omitempty"`
	PreviewId        int64      `json:"preview_id,omitempty" path:"preview_id,omitempty" url:"preview_id,omitempty"`
	Preview          Preview    `json:"preview,omitempty" path:"preview,omitempty" url:"preview,omitempty"`
	Action           string     `json:"action,omitempty" path:"action,omitempty" url:"action,omitempty"`
	Length           int64      `json:"length,omitempty" path:"length,omitempty" url:"length,omitempty"`
	MkdirParents     *bool      `json:"mkdir_parents,omitempty" path:"mkdir_parents,omitempty" url:"mkdir_parents,omitempty"`
	Part             int64      `json:"part,omitempty" path:"part,omitempty" url:"part,omitempty"`
	Parts            int64      `json:"parts,omitempty" path:"parts,omitempty" url:"parts,omitempty"`
	Ref              string     `json:"ref,omitempty" path:"ref,omitempty" url:"ref,omitempty"`
	Restart          int64      `json:"restart,omitempty" path:"restart,omitempty" url:"restart,omitempty"`
	Structure        string     `json:"structure,omitempty" path:"structure,omitempty" url:"structure,omitempty"`
	WithRename       *bool      `json:"with_rename,omitempty" path:"with_rename,omitempty" url:"with_rename,omitempty"`
}

func (f File) Identifier() interface{} {
	return f.Path
}

type FileCollection []File

type EtagsParam struct {
	Etag string `url:"etag,omitempty" json:"etag,omitempty" path:"etag"`
	Part string `url:"part,omitempty" json:"part,omitempty" path:"part"`
}

// Download file
type FileDownloadParams struct {
	Path              string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	Action            string `url:"action,omitempty" required:"false" json:"action,omitempty" path:"action"`
	PreviewSize       string `url:"preview_size,omitempty" required:"false" json:"preview_size,omitempty" path:"preview_size"`
	WithPreviews      *bool  `url:"with_previews,omitempty" required:"false" json:"with_previews,omitempty" path:"with_previews"`
	WithPriorityColor *bool  `url:"with_priority_color,omitempty" required:"false" json:"with_priority_color,omitempty" path:"with_priority_color"`
	File              File   `url:"-,omitempty" required:"false" json:"-,omitempty"`
}

type FileCreateParams struct {
	Path          string       `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
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
	Path          string     `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	ProvidedMtime *time.Time `url:"provided_mtime,omitempty" required:"false" json:"provided_mtime,omitempty" path:"provided_mtime"`
	PriorityColor string     `url:"priority_color,omitempty" required:"false" json:"priority_color,omitempty" path:"priority_color"`
}

type FileDeleteParams struct {
	Path      string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	Recursive *bool  `url:"recursive,omitempty" required:"false" json:"recursive,omitempty" path:"recursive"`
}

type FileFindParams struct {
	Path              string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	PreviewSize       string `url:"preview_size,omitempty" required:"false" json:"preview_size,omitempty" path:"preview_size"`
	WithPreviews      *bool  `url:"with_previews,omitempty" required:"false" json:"with_previews,omitempty" path:"with_previews"`
	WithPriorityColor *bool  `url:"with_priority_color,omitempty" required:"false" json:"with_priority_color,omitempty" path:"with_priority_color"`
}

// Copy file/folder
type FileCopyParams struct {
	Path        string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	Destination string `url:"destination,omitempty" required:"true" json:"destination,omitempty" path:"destination"`
	Structure   *bool  `url:"structure,omitempty" required:"false" json:"structure,omitempty" path:"structure"`
}

// Move file/folder
type FileMoveParams struct {
	Path        string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
	Destination string `url:"destination,omitempty" required:"true" json:"destination,omitempty" path:"destination"`
}

// Begin file upload
type FileBeginUploadParams struct {
	Path         string `url:"-,omitempty" required:"false" json:"-,omitempty" path:"path"`
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

func (f File) String() string {
	return f.Path
}

func (f File) Iterable() bool {
	return f.IsDir()
}

func (f File) IsDir() bool {
	return f.Type == "directory"
}

func (f File) ModTime() time.Time {
	if f.ProvidedMtime != nil {
		return *f.ProvidedMtime
	}
	if f.Mtime != nil {
		return *f.Mtime
	}
	return time.Time{}
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
