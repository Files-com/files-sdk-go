package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	goFs "io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/samber/lo"

	"github.com/Files-com/files-sdk-go/v2/lib"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/folder"
)

type FS struct {
	files_sdk.Config
	context.Context
	Root     string
	cache    map[string]*File
	cacheDir map[string][]goFs.DirEntry
	useCache bool
}

func (f *FS) Init(config files_sdk.Config, cache bool) *FS {
	f.Config = config
	f.ClearCache()
	f.useCache = cache
	return f
}

type WithContext interface {
	WithContext(ctx context.Context) interface{}
}

func (f *FS) WithContext(ctx context.Context) interface{} {
	return &FS{Context: ctx, Config: f.Config, cache: f.cache, useCache: f.useCache, cacheDir: f.cacheDir}
}

func (f *FS) ClearCache() {
	f.cache = make(map[string]*File)
	f.cacheDir = make(map[string][]goFs.DirEntry)
}

type File struct {
	*files_sdk.File
	*FS
	io.ReadCloser
	downloadRequestId string
	stat              bool
	fileMutex         *sync.Mutex
	realSize          int64
}

type ReadDirFile struct {
	*File
	count int
}

func (f *File) safeFile() files_sdk.File {
	f.fileMutex.Lock()
	defer f.fileMutex.Unlock()
	return *f.File
}

func (f *File) Init() *File {
	f.fileMutex = &sync.Mutex{}
	return f
}

func (f *File) Name() string {
	return f.safeFile().DisplayName
}

func (f *File) IsDir() bool {
	return f.safeFile().Type == "directory"
}

func (f *File) Type() goFs.FileMode {
	return goFs.ModePerm
}

func (f *File) Info() (goFs.FileInfo, error) {
	return f.Stat()
}

type Info struct {
	files_sdk.File
	realSize int64
}

func (i Info) Name() string {
	return i.File.DisplayName
}

func (i Info) Size() int64 {
	if i.realSize < 0 {
		return i.File.Size
	}
	return i.realSize
}

type UntrustedSize interface {
	UntrustedSize() bool
}

func (i Info) UntrustedSize() bool {
	return i.realSize == -1
}

type PossibleSize interface {
	PossibleSize() int64
}

func (i Info) PossibleSize() int64 {
	return i.File.Size
}

func (i Info) Mode() goFs.FileMode {
	return goFs.ModePerm
}
func (i Info) ModTime() time.Time {
	return *i.File.Mtime
}

func (i Info) IsDir() bool {
	return i.File.Type == "directory"
}
func (i Info) Sys() interface{} {
	return i.File
}

type RealTimeStat interface {
	RealTimeStat() (goFs.FileInfo, error)
}

func (f *File) RealTimeStat() (goFs.FileInfo, error) {
	var err error
	if f.safeFile().Type == "directory" {
		return Info{File: *f.File, realSize: f.realSize}, err
	}
	if f.stat {
		return Info{File: *f.File, realSize: f.realSize}, err
	}

	f.fileMutex.Lock()
	defer f.fileMutex.Unlock()

	statCtx, cancel := context.WithTimeout(f.Context, time.Second*15)
	defer cancel()
	var tempFile files_sdk.File
	tempFile, err = (&Client{Config: f.Config}).FileStats(statCtx, *f.File)
	f.stat = true
	f.realSize = tempFile.Size
	return Info{File: *f.File, realSize: f.realSize}, err
}

func (f *File) Stat() (goFs.FileInfo, error) {
	return Info{File: *f.File, realSize: -2}, nil
}

func (f *File) Read(b []byte) (n int, err error) {
	f.fileMutex.Lock()
	if f.ReadCloser == nil {
		*f.File, err = (&Client{Config: f.Config}).Download(
			f.Context,
			files_sdk.FileDownloadParams{File: *f.File},
			files_sdk.ResponseOption(func(response *http.Response) error {
				if err := lib.ResponseErrors(response, lib.IsStatus(http.StatusForbidden), lib.NotStatus(http.StatusOK)); err != nil {
					return &goFs.PathError{Path: f.File.Path, Err: err, Op: "read"}
				}
				f.downloadRequestId = response.Header.Get("X-Files-Download-Request-Id")
				f.ReadCloser = response.Body
				return nil
			}),
		)
	}

	if downloadRequestExpired(err) {
		f.Config.LogPath(f.File.Path, map[string]interface{}{"message": "downloadRequestExpired", "error": err})
		f.fileMutex.Unlock()
		f.File.DownloadUri = "" // force a new query
		err = f.downloadURI()
		if err != nil {
			return n, err
		}

		return f.Read(b)
	}

	if err != nil {
		return
	}

	defer f.fileMutex.Unlock()
	return f.ReadCloser.Read(b)
}

type ReaderRange interface {
	ReaderRange(off int64, end int64) (io.ReadCloser, error)
}

func (f *File) ReaderRange(off int64, end int64) (r io.ReadCloser, err error) {
	err = f.downloadURI()
	if err != nil {
		return
	}
	headers := &http.Header{}
	headers.Set("Range", fmt.Sprintf("bytes=%v-%v", off, end))
	_, err = (&Client{Config: f.Config}).Download(
		f.Context,
		files_sdk.FileDownloadParams{File: *f.File},
		files_sdk.RequestHeadersOption(headers),
		files_sdk.ResponseOption(func(response *http.Response) error {
			if err := lib.ResponseErrors(response, lib.IsStatus(http.StatusForbidden), lib.NotStatus(http.StatusPartialContent)); err != nil {
				return &goFs.PathError{Path: f.File.Path, Err: err, Op: "ReaderRange"}
			}
			r = response.Body
			return nil
		}),
	)

	if downloadRequestExpired(err) {
		f.Config.LogPath(f.File.Path, map[string]interface{}{"message": "downloadRequestExpired", "error": err})
		f.File.DownloadUri = "" // force a new query
		err = f.downloadURI()
		if err != nil {
			return r, err
		}

		return f.ReaderRange(off, end)
	}
	return r, err
}

func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	err = f.downloadURI()
	if err != nil {
		return
	}
	headers := &http.Header{}
	headers.Set("Range", fmt.Sprintf("bytes=%v-%v", off, int64(len(p))+off-1))
	_, err = (&Client{Config: f.Config}).Download(
		f.Context,
		files_sdk.FileDownloadParams{
			File: *f.File,
		},
		files_sdk.RequestHeadersOption(headers),
		files_sdk.ResponseOption(func(response *http.Response) error {
			if err := lib.ResponseErrors(response, lib.IsStatus(http.StatusForbidden), lib.NotStatus(http.StatusPartialContent)); err != nil {
				return &goFs.PathError{Path: f.File.Path, Err: err, Op: "ReadAt"}
			}
			n, err = io.ReadFull(response.Body, p)
			if err != nil && err != io.EOF {
				return err
			}
			if int64(len(p)) >= response.ContentLength && int64(n) != response.ContentLength {
				return &goFs.PathError{Path: f.File.Path, Err: fmt.Errorf("content-length did not match body"), Op: "ReadAt"}
			}
			return nil
		}),
	)

	if downloadRequestExpired(err) {
		f.Config.LogPath(f.File.Path, map[string]interface{}{"message": "downloadRequestExpired", "error": err})
		f.File.DownloadUri = "" // force a new query
		err = f.downloadURI()
		if err != nil {
			return n, err
		}

		return f.ReadAt(p, off)
	}

	return n, err
}

func downloadRequestExpired(err error) bool {
	if err == nil {
		return false
	}
	responseErr, ok := errors.Unwrap(err).(lib.ResponseError)
	return ok && responseErr.StatusCode == http.StatusForbidden
}

func (f *File) downloadURI() (err error) {
	f.fileMutex.Lock()
	*f.File, err = (&Client{Config: f.Config}).DownloadUri(f.Context, files_sdk.FileDownloadParams{File: *f.File})
	f.fileMutex.Unlock()
	return
}

func (f *File) Close() error {
	if f.ReadCloser == nil {
		return nil
	}
	defer func() { f.ReadCloser = nil }()
	err := f.ReadCloser.Close()

	if err != nil {
		return err
	}

	info, err := f.Info()
	if err == nil && info.(UntrustedSize).UntrustedSize() {
		status, err := (&Client{Config: f.Config}).DownloadRequestStatus(f.Context, f.File.DownloadUri, f.downloadRequestId)
		if err != nil {
			return err
		}
		if !status.IsNil() {
			return status
		}

		if status.Data.Status == "completed" {
			f.realSize = status.Data.BytesTransferred
			if dataBytes, err := json.Marshal(status.Data); err == nil {
				dataMap := make(map[string]interface{})
				if err = json.Unmarshal(dataBytes, &dataMap); err == nil {
					f.Config.LogPath(info.Name(), lo.Assign(dataMap, map[string]interface{}{"message": "download request server status"}))
				}
			}
		} else {
			return fmt.Errorf("server reported transfer '%v'", status.Data.Status)
		}
	}
	return err
}
func (f *File) WithContext(ctx context.Context) interface{} {
	newF := *f
	fs := *newF.FS
	newF.FS = fs.WithContext(ctx).(*FS)
	return &newF
}

func (f *FS) Open(name string) (goFs.File, error) {
	if name == "." {
		name = ""
	}
	file, ok := f.cache[name]
	if ok {
		if file.IsDir() {
			return &ReadDirFile{File: file}, nil
		}
		return file, nil
	}
	path := lib.UrlJoinNoEscape(f.Root, name)
	var err error
	var fileInfo files_sdk.File
	if path == "" { // skip call on root path
		fileInfo = files_sdk.File{Type: "directory"}
	} else {
		fileInfo, err = (&Client{Config: f.Config}).Find(f.Context, files_sdk.FileFindParams{Path: path})
		if err != nil {
			return &File{}, &goFs.PathError{Path: fileInfo.Path, Err: err, Op: "open"}
		}
	}

	file = (&File{File: &fileInfo, FS: f}).Init()
	if f.useCache {
		f.cache[name] = file
	}
	if fileInfo.Type == "directory" {
		return &ReadDirFile{File: file}, nil
	} else {
		return file, nil
	}
}

func (f *FS) ReadDir(name string) ([]goFs.DirEntry, error) {
	if name == "." {
		name = ""
	}
	dirs, ok := f.cacheDir[name]
	if ok {
		return dirs, nil
	}

	dirs, err := ReadDirFile{File: (&File{File: &files_sdk.File{Path: name}, FS: f}).Init()}.ReadDir(0)
	if err != nil {
		return dirs, err
	}
	if f.useCache {
		f.cacheDir[name] = dirs
	}
	return dirs, nil
}

func (f ReadDirFile) ReadDir(n int) ([]goFs.DirEntry, error) {
	var files []goFs.DirEntry
	if f.Context != nil && f.Context.Err() != nil {
		return files, &goFs.PathError{Path: f.Path, Err: f.Context.Err(), Op: "readdir"}
	}
	folderClient := folder.Client{Config: f.Config}
	it, err := folderClient.ListFor(f.Context, files_sdk.FolderListForParams{Path: f.Path})
	if err != nil {
		return files, &goFs.PathError{Path: f.Path, Err: err, Op: "readdir"}
	}
	if f.count > 0 {
		return files, io.EOF
	}
	for it.Next() && (n <= 0 || n > 0 && n >= f.count) {
		fi := it.File()
		if err != nil {
			return files, &goFs.PathError{Path: f.Path, Err: err, Op: "readdir"}
		}
		parts := strings.Split(fi.Path, "/")
		dir := strings.Join(parts[0:len(parts)-1], "/")
		if dir == strings.TrimSuffix(f.Path, "/") {
			// There is a bug in the API that it could return a nested file not in the current directory.
			file := (&File{File: &fi, FS: f.FS}).Init()
			if f.useCache {
				f.cache[fi.Path] = file
			}
			files = append(files, file)
		}

		f.count += 1
	}

	if it.Err() != nil {
		return files, &goFs.PathError{Path: f.Path, Err: it.Err(), Op: "readdir"}
	}
	return files, nil
}

func (f *FS) MkdirAll(dir string, _ goFs.FileMode) error {
	var parentPath string
	for _, dirPath := range strings.Split(dir, "/") {
		if dirPath == "" {
			break
		}
		folderClient := folder.Client{Config: f.Config}
		_, err := folderClient.Create(f.Context, files_sdk.FolderCreateParams{Path: lib.UrlJoinNoEscape(parentPath, dirPath)})
		rErr, ok := err.(files_sdk.ResponseError)
		if err != nil && ok && rErr.Type != "processing-failure/destination-exists" {
			return err
		}

		parentPath = lib.UrlJoinNoEscape(parentPath, dirPath)
	}
	return nil
}
