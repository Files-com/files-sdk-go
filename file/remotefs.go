package file

import (
	"context"
	"fmt"
	"io"
	goFs "io/fs"
	"net/http"
	"path/filepath"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/folder"
)

type FS struct {
	files_sdk.Config
	context.Context
	Root  string
	cache map[string]File
}

func (f FS) Init(config files_sdk.Config) FS {
	return FS{Config: config, cache: make(map[string]File)}
}

func (f FS) WithContext(ctx context.Context) FS {
	return FS{Context: ctx, Config: f.Config, cache: f.cache}
}

type File struct {
	*files_sdk.File
	FS
	io.ReadCloser
}

type ReadDirFile struct {
	File
	count int
}

func (f File) Name() string {
	return f.File.DisplayName
}

func (f File) IsDir() bool {
	return f.File.Type == "directory"
}

func (f File) Type() goFs.FileMode {
	return goFs.ModePerm
}

func (f File) Info() (goFs.FileInfo, error) {
	return Info{File: f.File}, nil
}

type Info struct {
	*files_sdk.File
}

func (i Info) Name() string {
	return i.File.DisplayName
}

func (i Info) Size() int64 {
	return i.File.Size
}

func (i Info) Mode() goFs.FileMode {
	return goFs.ModePerm
}
func (i Info) ModTime() time.Time {
	return i.File.Mtime
}

func (i Info) IsDir() bool {
	return i.File.Type == "directory"
}
func (i Info) Sys() interface{} {
	return *i.File
}

func (f File) Stat() (goFs.FileInfo, error) {
	return Info{File: f.File}, nil
}

func (f *File) Read(b []byte) (int, error) {
	pathErr := f.load()
	if pathErr != nil {
		return 0, pathErr
	}
	return f.ReadCloser.Read(b)
}

func (f *File) load() error {
	if f.ReadCloser != nil {
		return nil
	}
	f1, err := f.Reload()
	if err != nil {
		return &goFs.PathError{Path: f.File.Path, Err: err, Op: "read"}
	}
	request, err := http.NewRequestWithContext(f.Context, "GET", f1.File.DownloadUri, nil)
	if err != nil {
		return &goFs.PathError{Path: f1.File.Path, Err: err, Op: "read"}
	}
	resp, err := f.Config.GetHttpClient().Do(request)
	if err != nil {
		return &goFs.PathError{Path: f1.File.Path, Err: err, Op: "read"}
	}
	if resp.StatusCode != 200 {
		body := make([]byte, resp.ContentLength)
		resp.Body.Read(body)
		return &goFs.PathError{Path: f1.File.Path, Err: fmt.Errorf(string(body)), Op: "read"}
	}
	f.ReadCloser = resp.Body
	f.Size = resp.ContentLength
	return nil
}

func (f *File) Reload() (File, error) {
	fileInfo, err := (&Client{Config: f.Config}).Get(f.Context, f.Path)
	if err != nil {
		return File{File: &fileInfo, FS: f.FS}, err
	}
	f.File = &fileInfo
	return File{File: &fileInfo, FS: f.FS}, nil
}

func (f *File) Close() error {
	if f.ReadCloser != nil {
		return f.ReadCloser.Close()
	}
	return nil
}

func (f FS) Open(name string) (goFs.File, error) {
	file, ok := f.cache[name]
	if ok {
		if file.IsDir() {
			return &ReadDirFile{File: file}, nil
		}
		return &file, nil
	}
	fileInfo, err := (&Client{Config: f.Config}).Get(f.Context, filepath.Join(f.Root, name))
	responseError, ok := err.(files_sdk.ResponseError)

	if err != nil {
		if ok && responseError.Type == "bad-request/cannot-download-directory" {
			fileInfo.Type = "directory"
			fileInfo.DisplayName = name
			fileInfo.Path = filepath.Join(f.Root, name)
			f.cache[name] = File{File: &fileInfo, FS: f}
			return &ReadDirFile{File: File{File: &fileInfo, FS: f}}, nil
		} else {
			return &File{File: &fileInfo, FS: f}, &goFs.PathError{Path: fileInfo.Path, Err: err, Op: "open"}
		}
	}
	return &File{File: &fileInfo, FS: f}, nil
}

func (f ReadDirFile) ReadDir(n int) ([]goFs.DirEntry, error) {
	var files []goFs.DirEntry
	if f.Context.Err() != nil {
		return files, &goFs.PathError{Path: f.Path, Err: f.Context.Err(), Op: "readdir"}
	}
	folderClient := folder.Client{Config: f.Config}
	it, err := folderClient.ListFor(f.Context, files_sdk.FolderListForParams{Path: f.File.Path})
	if err != nil {
		return files, &goFs.PathError{Path: f.Path, Err: err, Op: "readdir"}
	}
	if f.count > 0 {
		return files, io.EOF
	}
	for it.Next() && (n <= 0 || n > 0 && n >= f.count) {
		fl := it.Folder()
		fi, err := fl.ToFile()
		if err != nil {
			return files, &goFs.PathError{Path: f.Path, Err: err, Op: "readdir"}
		}
		f.cache[fi.Path] = File{File: &fi, FS: f.FS}
		files = append(files, File{File: &fi, FS: f.FS})
		f.count += 1
	}

	if it.Err() != nil {
		return files, &goFs.PathError{Path: f.Path, Err: err, Op: "readdir"}
	}
	return files, nil
}
