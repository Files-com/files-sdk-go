package files_sdk

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// FileMigration is a Files.com API resource.
type FileMigration struct {
	// File migration ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Source path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Destination path.
	DestPath string `json:"dest_path,omitempty" path:"dest_path,omitempty" url:"dest_path,omitempty"`
	// Reason for the failure, if applicable.
	FailureMessage string `json:"failure_message,omitempty" path:"failure_message,omitempty" url:"failure_message,omitempty"`
	// Number of files processed.
	FilesMoved int64 `json:"files_moved,omitempty" path:"files_moved,omitempty" url:"files_moved,omitempty"`
	FilesTotal int64 `json:"files_total,omitempty" path:"files_total,omitempty" url:"files_total,omitempty"`
	// The type of operation.
	Operation string `json:"operation,omitempty" path:"operation,omitempty" url:"operation,omitempty"`
	// Region.
	Region string `json:"region,omitempty" path:"region,omitempty" url:"region,omitempty"`
	// Status.
	Status string `json:"status,omitempty" path:"status,omitempty" url:"status,omitempty"`
	// Link to download the log file for this migration.
	LogUrl string `json:"log_url,omitempty" path:"log_url,omitempty" url:"log_url,omitempty"`
}

// Identifier returns the resource ID.
func (f FileMigration) Identifier() interface{} {
	return f.Id
}

// FileMigrationCollection is a list of FileMigration resources.
type FileMigrationCollection []FileMigration

// this code is merged into the built version of filemigration.go

type FilesMigrationLog struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Status    string    `json:"status"`
	FileType  string    `json:"file_type"`
	Path      string    `json:"path"`
	DestPath  string    `json:"dest_path"`
}

// FilesMigrationLogIter Transforms migrations into a log iterator
type FilesMigrationLogIter struct {
	context.Context
	Config
	FileMigration
	next func() bool
	more func() bool
	error
	current interface{}
	body    io.ReadCloser
}

func (l *FilesMigrationLogIter) Err() error {
	return l.error
}

func (l *FilesMigrationLogIter) Current() interface{} {
	return l.current
}

func (l *FilesMigrationLogIter) Next() bool {
	return l.next()
}

// Close releases the log response and stops iteration. Call it when stopping
// before Next returns false; exhaustion and decoding failures close automatically.
func (l *FilesMigrationLogIter) Close() error {
	if l.body == nil {
		return nil
	}
	body := l.body
	l.body = nil
	l.next = func() bool { return false }
	return body.Close()
}

func (l FilesMigrationLogIter) Init() *FilesMigrationLogIter {
	l.next = func() bool {
		return false
	}
	if l.FileMigration.LogUrl != "" {
		req, err := http.NewRequestWithContext(l.Context, "GET", l.FileMigration.LogUrl, nil)
		if err != nil {
			l.error = err
			return &l
		}
		resp, err := l.Config.Do(req)
		if err != nil {
			l.error = err
			return &l
		}
		decoder := json.NewDecoder(resp.Body)
		l.body = resp.Body
		l.more = decoder.More
		l.next = func() bool {
			if l.error != nil {
				return false
			}
			if l.more() {
				var log FilesMigrationLog
				if err := decoder.Decode(&log); err != nil {
					l.error = err
					l.Close()
					return false
				}
				l.current = log
				return true
			}
			l.error = l.Close()
			return false
		}
	}

	return &l
}

// FileMigrationFindParams contains the request parameters for GET /file_migrations/{id}.
type FileMigrationFindParams struct {
	// File Migration ID. Required; sent in the request path.
	Id int64 `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (f *FileMigration) UnmarshalJSON(data []byte) error {
	type fileMigration FileMigration
	var v fileMigration
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*f = FileMigration(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (f *FileMigrationCollection) UnmarshalJSON(data []byte) error {
	type fileMigrations FileMigrationCollection
	var v fileMigrations
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*f = FileMigrationCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (f *FileMigrationCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*f))
	for i, v := range *f {
		ret[i] = v
	}

	return &ret
}
