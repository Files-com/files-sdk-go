package files_sdk

import (
	"encoding/json"
	"time"
)

type File struct {
	Id               int64     `json:"id,omitempty"`
	Path             string    `json:"path,omitempty"`
	DisplayName      string    `json:"display_name,omitempty"`
	Type             string    `json:"type,omitempty"`
	Size             int       `json:"size,omitempty"`
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
	Length           int       `json:"length,omitempty"`
	MkdirParents     *bool     `json:"mkdir_parents,omitempty"`
	Part             int       `json:"part,omitempty"`
	Parts            int       `json:"parts,omitempty"`
	Ref              string    `json:"ref,omitempty"`
	Restart          int       `json:"restart,omitempty"`
	Structure        string    `json:"structure,omitempty"`
	WithRename       *bool     `json:"with_rename,omitempty"`
}

type FileCollection []File

type EtagsParam struct {
	Etag string `url:"etag,omitempty"`
	Part string `url:"part,omitempty"`
}

type FileDownloadParams struct {
	Path              string `url:"-,omitempty"`
	Action            string `url:"action,omitempty"`
	PreviewSize       string `url:"preview_size,omitempty"`
	WithPreviews      *bool  `url:"with_previews,omitempty"`
	WithPriorityColor *bool  `url:"with_priority_color,omitempty"`
}

type FileCreateParams struct {
	Path          string       `url:"-,omitempty"`
	Action        string       `url:"action,omitempty"`
	EtagsParam    []EtagsParam `url:"etags,omitempty"`
	Length        int          `url:"length,omitempty"`
	MkdirParents  *bool        `url:"mkdir_parents,omitempty"`
	Part          int          `url:"part,omitempty"`
	Parts         int          `url:"parts,omitempty"`
	ProvidedMtime string       `url:"provided_mtime,omitempty"`
	Ref           string       `url:"ref,omitempty"`
	Restart       int          `url:"restart,omitempty"`
	Size          int          `url:"size,omitempty"`
	Structure     string       `url:"structure,omitempty"`
	WithRename    *bool        `url:"with_rename,omitempty"`
}

type FileUpdateParams struct {
	Path          string `url:"-,omitempty"`
	ProvidedMtime string `url:"provided_mtime,omitempty"`
	PriorityColor string `url:"priority_color,omitempty"`
}

type FileDeleteParams struct {
	Path      string `url:"-,omitempty"`
	Recursive *bool  `url:"recursive,omitempty"`
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
