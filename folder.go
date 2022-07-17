package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Folder struct {
	Path             string     `json:"path,omitempty" path:"path"`
	DisplayName      string     `json:"display_name,omitempty" path:"display_name"`
	Type             string     `json:"type,omitempty" path:"type"`
	Size             int64      `json:"size,omitempty" path:"size"`
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
	MkdirParents     *bool      `json:"mkdir_parents,omitempty" path:"mkdir_parents"`
}

type FolderCollection []Folder

type FolderListForParams struct {
	Path              string `url:"-,omitempty" required:"true" json:"-,omitempty" path:"path"`
	Filter            string `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	PreviewSize       string `url:"preview_size,omitempty" required:"false" json:"preview_size,omitempty" path:"preview_size"`
	Search            string `url:"search,omitempty" required:"false" json:"search,omitempty" path:"search"`
	SearchAll         *bool  `url:"search_all,omitempty" required:"false" json:"search_all,omitempty" path:"search_all"`
	WithPreviews      *bool  `url:"with_previews,omitempty" required:"false" json:"with_previews,omitempty" path:"with_previews"`
	WithPriorityColor *bool  `url:"with_priority_color,omitempty" required:"false" json:"with_priority_color,omitempty" path:"with_priority_color"`
	lib.ListParams
}

type FolderCreateParams struct {
	Path         string `url:"-,omitempty" required:"true" json:"-,omitempty" path:"path"`
	MkdirParents *bool  `url:"mkdir_parents,omitempty" required:"false" json:"mkdir_parents,omitempty" path:"mkdir_parents"`
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
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = Folder(v)
	return nil
}

func (f *FolderCollection) UnmarshalJSON(data []byte) error {
	type folders FolderCollection
	var v folders
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
