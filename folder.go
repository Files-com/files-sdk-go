package files_sdk

import (
	"encoding/json"
	lib "github.com/Files-com/files-sdk-go/lib"
	"time"
)

type Folder struct {
	Id               int       `json:"id,omitempty"`
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
	PreviewId        int       `json:"preview_id,omitempty"`
	Preview          string    `json:"preview,omitempty"`
}

type FolderCollection []Folder

type FolderListForParams struct {
	Page              int    `url:"page,omitempty"`
	PerPage           int    `url:"per_page,omitempty"`
	Action            string `url:"action,omitempty"`
	Path              string `url:"-,omitempty"`
	Cursor            string `url:"cursor,omitempty"`
	Filter            string `url:"filter,omitempty"`
	PreviewSize       string `url:"preview_size,omitempty"`
	Search            string `url:"search,omitempty"`
	SearchAll         *bool  `url:"search_all,omitempty"`
	WithPreviews      *bool  `url:"with_previews,omitempty"`
	WithPriorityColor *bool  `url:"with_priority_color,omitempty"`
	lib.ListParams
}

type FolderCreateParams struct {
	Path string `url:"-,omitempty"`
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
