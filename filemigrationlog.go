package files_sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type FilesMigrationLog struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Status    string    `json:"status"`
	FileType  string    `json:"file_type"`
	Path      string    `json:"path"`
	DestPath  string    `json:"dest_path"`
}

type FilesMigrationLogIter struct {
	context.Context
	Config
	FileMigration
	next func() bool
	more func() bool
	error
	current interface{}
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
		resp, err := l.Config.GetHttpClient().Do(req)
		if err != nil {
			l.error = err
			return &l
		}
		decoder := json.NewDecoder(resp.Body)
		l.more = decoder.More
		l.next = func() bool {
			if l.error != nil {
				return false
			}
			if l.more() {
				var log FilesMigrationLog
				if err := decoder.Decode(&log); err != nil {
					l.error = err
					return false
				}
				l.current = log
				return true
			}
			return false
		}
	}

	return &l
}
