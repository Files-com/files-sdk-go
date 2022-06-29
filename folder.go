package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Folder struct {
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
	MkdirParents     *bool      `json:"mkdir_parents,omitempty"`
}

type FolderCollection []Folder

type FolderListForParams struct {
	Cursor            string `url:"cursor,omitempty" required:"false" json:"cursor,omitempty"`
	PerPage           int64  `url:"per_page,omitempty" required:"false" json:"per_page,omitempty"`
	Path              string `url:"-,omitempty" required:"true" json:"-,omitempty"`
	Filter            string `url:"filter,omitempty" required:"false" json:"filter,omitempty"`
	PreviewSize       string `url:"preview_size,omitempty" required:"false" json:"preview_size,omitempty"`
	Search            string `url:"search,omitempty" required:"false" json:"search,omitempty"`
	SearchAll         *bool  `url:"search_all,omitempty" required:"false" json:"search_all,omitempty"`
	WithPreviews      *bool  `url:"with_previews,omitempty" required:"false" json:"with_previews,omitempty"`
	WithPriorityColor *bool  `url:"with_priority_color,omitempty" required:"false" json:"with_priority_color,omitempty"`
	lib.ListParams
}

type FolderCreateParams struct {
	Path         string `url:"-,omitempty" required:"true" json:"-,omitempty"`
	MkdirParents *bool  `url:"mkdir_parents,omitempty" required:"false" json:"mkdir_parents,omitempty"`
}

func (f *Folder) ToFile() (File, error) {
	bodyBytes, err := json.Marshal(f)
	if err != nil {
		return File{}, err
	}
	file := File{}
	file.UnmarshalJSON(bodyBytes)
	return file, nil
}

func (f *Folder) UnmarshalJSON(data []byte) error {
	type folder Folder
	var v folder
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = Folder(v)
	return nil
}

func (f *FolderCollection) UnmarshalJSON(data []byte) error {
	type folders []Folder
	var v folders
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FolderCollection(v)
	return nil
}

func (f *FolderCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
