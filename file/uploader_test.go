package file

import (
	"context"
	"io/fs"
	"math"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/Files-com/files-sdk-go/v2/lib"

	"github.com/Files-com/files-sdk-go/v2/directory"

	"github.com/Files-com/files-sdk-go/v2/ignore"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file/status"

	"github.com/stretchr/testify/assert"
)

type MockUploader struct {
	files_sdk.File
	findError   files_sdk.ResponseError
	uploadError error
}

func (m *MockUploader) UploadIO(context.Context, UploadIOParams) (files_sdk.File, files_sdk.FileUploadPart, Parts, []error, error) {
	return files_sdk.File{}, files_sdk.FileUploadPart{}, Parts{}, []error{}, m.uploadError
}

func (m *MockUploader) Find(context.Context, files_sdk.FileFindParams, ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return m.File, m.findError
}

func Test_skipOrIgnore(t *testing.T) {
	var progressReportError error

	init := func() (*UploadStatus, fstest.MapFS, *status.Job) {
		job := status.Job{Logger: (&files_sdk.Config{}).Logger()}.Init()
		job.GitIgnore, _ = ignore.New()
		job.Params = UploaderParams{}
		uploadStatus := &UploadStatus{job: job, Mutex: &sync.RWMutex{}, file: files_sdk.File{Path: "test"}, remotePath: "test"}
		uploadStatus.Job().Add(uploadStatus)
		uploader := &MockUploader{}
		uploadStatus.Uploader = uploader
		uploadStatus.job.EventsReporter = Reporter(func(s status.File) {
			progressReportError = s.Err
		})

		mockFs := make(fstest.MapFS)
		job.RemoteFs = mockFs
		return uploadStatus, mockFs, job
	}
	assert := assert.New(t)
	uploadStatus, mockFs, job := init()

	t.Run("sync not enabled and sizes do match", func(t *testing.T) {
		mockFs["test"] = &fstest.MapFile{
			Sys: files_sdk.File{Size: 10},
		}
		uploadStatus.file.Size = 10
		uploadStatus.Sync = false
		assert.Equal(false, skipOrIgnore(uploadStatus, false))
	})

	t.Run("when sizes don't match", func(t *testing.T) {
		uploadStatus.Sync = true
		mockFs["test"] = &fstest.MapFile{
			Sys: files_sdk.File{Size: 9},
		}
		uploadStatus.file.Size = 10
		assert.Equal(false, skipOrIgnore(uploadStatus, false))
		assert.Equal(nil, progressReportError)
	})

	t.Run("when sizes do match", func(t *testing.T) {
		mockFs["test"] = &fstest.MapFile{
			Sys: files_sdk.File{Size: 10},
		}
		uploadStatus.file.Size = 10
		assert.Equal(true, skipOrIgnore(uploadStatus, false))
	})

	t.Run("There is no server version", func(t *testing.T) {
		delete(mockFs, "test")
		assert.Equal(false, skipOrIgnore(uploadStatus, false))
	})

	t.Run("when sizes do match on a deeply nested path", func(t *testing.T) {
		oldUploadStatus := *uploadStatus
		uploadStatus = &UploadStatus{job: job, Mutex: &sync.RWMutex{}, file: files_sdk.File{Path: "test/path/test"}, remotePath: "test/path/test"}
		uploadStatus.Sync = true
		mockFs["test/path/test"] = &fstest.MapFile{
			Sys: files_sdk.File{Size: 10},
		}
		mockFs["test/path"] = &fstest.MapFile{
			Sys:  files_sdk.File{Size: 10},
			Mode: fs.ModeDir,
		}
		uploadStatus.file.Size = 10
		assert.Equal(true, skipOrIgnore(uploadStatus, false))
		uploadStatus = &oldUploadStatus
	})

	t.Run("when transfer is a single file", func(t *testing.T) {
		uploadStatus.job.Type = directory.File
		uploadStatus.Sync = true
		mockFs["test"] = &fstest.MapFile{
			Sys: files_sdk.File{Size: 10},
		}
		uploadStatus.file.Size = 10
		assert.Equal(true, skipOrIgnore(uploadStatus, false))
		uploadStatus.job.Type = directory.Dir
		uploadStatus.Sync = false
	})

	t.Run("case insensitive", func(t *testing.T) {
		uploadStatus.job.Type = directory.Dir
		uploadStatus.Sync = true
		mockFs["Test"] = &fstest.MapFile{
			Sys: files_sdk.File{Size: 10},
		}
		uploadStatus.file.Size = 10
		assert.Equal(true, skipOrIgnore(uploadStatus, false))
		uploadStatus.Sync = false
	})

	t.Run("Ignore files", func(t *testing.T) {
		job.GitIgnore, _ = ignore.New([]string{"*.css"}...)
		uploadStatus.localPath = "main.css"
		assert.Equal(true, skipOrIgnore(uploadStatus, false))

		uploadStatus.localPath = "main.php"
		assert.Equal(false, skipOrIgnore(uploadStatus, false))

		job.GitIgnore, _ = ignore.New([]string{"*.css", "*.php"}...)
		uploadStatus.localPath = "main.css"
		assert.Equal(true, skipOrIgnore(uploadStatus, false))
	})

	t.Run("FeatureFlag incremental-updates", func(t *testing.T) {
		type Args struct {
			destinationMtime *time.Time
			sourceMtime      time.Time
			sync             bool
		}
		t.Run("when sizes do not match", func(t *testing.T) {
			test := func(t *testing.T, args Args) bool {
				uploadStatus, mockFs, _ := init()
				mockFs["test"] = &fstest.MapFile{
					Sys: files_sdk.File{Size: 11, Mtime: args.destinationMtime},
				}
				uploadStatus.Sync = true
				uploadStatus.file.Size = 10
				uploadStatus.file.Mtime = &args.sourceMtime
				skip := skipOrIgnore(uploadStatus, true)
				if args.destinationMtime != nil {
					diff := args.destinationMtime.Sub(args.sourceMtime)
					diff = time.Duration(math.Abs(float64(diff))).Truncate(time.Millisecond)
					t.Logf("%v - destinationMtime: %v, sourceMtime: %v, diff: %v, sync: %v", t.Name(), args.destinationMtime.Format(time.RFC3339), args.sourceMtime.Format(time.RFC3339), diff, !skip)
				}
				assert.Equal(args.sync, !skip, "should it sync")
				return skip
			}

			t.Run("when destination Mtime is older than source time it should sync", func(t *testing.T) {
				test(
					t,
					Args{
						destinationMtime: lib.Time(time.Now().Add(-time.Hour * 48)),
						sourceMtime:      time.Now().Add(-time.Hour * 24),
						sync:             true,
					},
				)
			})

			t.Run("when destination Mtime is the within the same minute as source Mtime it should not sync", func(t *testing.T) {
				test(
					t,
					Args{
						destinationMtime: lib.Time(
							time.Date(2021, 8, 15, 14, 30, 45, 100, time.Local),
						),
						sourceMtime: time.Date(2021, 8, 15, 14, 30, 00, 100, time.Local),
						sync:        false,
					},
				)
			})

			t.Run("when destination Mtime is the just outside the same minute as source Mtime it should not sync", func(t *testing.T) {
				test(
					t,
					Args{
						destinationMtime: lib.Time(
							time.Date(2021, 8, 15, 14, 30, 45, 100, time.Local),
						),
						sourceMtime: time.Date(2021, 8, 15, 14, 31, 45, 100, time.Local),
						sync:        true,
					},
				)
			})

			t.Run("when destination Mtime is the same source Mtime it should not sync", func(t *testing.T) {
				mtime := time.Now().Add(-time.Hour * 24)
				test(
					t,
					Args{
						destinationMtime: &mtime,
						sourceMtime:      mtime,
						sync:             false,
					},
				)
			})

			t.Run("when destination Mtime is the same source Mtime but in different time zones it should not sync", func(t *testing.T) {
				test(
					t,
					Args{
						destinationMtime: lib.Time(time.Date(2021, 8, 15, 14, 31, 45, 100, time.Local).UTC()),
						sourceMtime:      time.Date(2021, 8, 15, 14, 31, 45, 100, time.Local),
						sync:             false,
					},
				)
			})

			t.Run("when destination Mtime is newer than source Mtime it should not sync", func(t *testing.T) {
				test(
					t,
					Args{
						destinationMtime: lib.Time(time.Date(2021, 8, 16, 14, 31, 45, 100, time.Local)),
						sourceMtime:      time.Date(2021, 8, 15, 14, 31, 45, 100, time.Local),
						sync:             false,
					},
				)
			})

			t.Run("when destination Mtime is nil", func(t *testing.T) {
				test(
					t,
					Args{
						destinationMtime: nil,
						sourceMtime:      time.Now().Add(-time.Hour * 24),
						sync:             true,
					},
				)
			})

			t.Run("when destination Mtime is zero", func(t *testing.T) {
				test(
					t,
					Args{
						destinationMtime: &time.Time{},
						sourceMtime:      time.Now().Add(-time.Hour * 24),
						sync:             true,
					},
				)
			})
		})
	})
}

func TestUploader(t *testing.T) {
	mutex := &sync.Mutex{}
	t.Run("uploader", func(t *testing.T) {
		sourceFs := lib.ReadWriteFs(lib.LocalFileSystem{})
		destinationFs := &FS{Context: context.Background()}
		for _, tt := range lib.PathSpec(sourceFs.PathSeparator(), destinationFs.PathSeparator()) {
			t.Run(tt.Name, func(t *testing.T) {
				client, r, err := CreateClient(t.Name())
				if err != nil {
					t.Fatal(err)
				}
				config := client.Config
				destinationFs = (&FS{Context: context.Background()}).Init(config, true)
				lib.BuildPathSpecTest(t, mutex, tt, sourceFs, destinationFs, func(source, destination string) lib.Cmd {
					return &CmdRunner{
						run: func() *status.Job {
							return client.Uploader(context.Background(), UploaderParams{LocalPath: source, RemotePath: destination, Config: config})
						},
						args: []string{source, destination},
					}
				})
				r.Stop()
			})
		}
	})
}
