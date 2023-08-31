package file

import (
	"bytes"
	"context"
	"io/fs"
	"math"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/directory"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/ignore"
	"github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
)

type MockUploader struct {
	files_sdk.File
	findError   files_sdk.ResponseError
	uploadError error
}

func (m *MockUploader) UploadWithResume(...UploadOption) (UploadResumable, error) {
	return UploadResumable{}, m.uploadError
}

func (m *MockUploader) Find(files_sdk.FileFindParams, ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return m.File, m.findError
}

func Test_skipOrIgnore(t *testing.T) {
	var progressReportError error

	init := func() (*UploadStatus, fstest.MapFS, *status.Job) {
		job := (&status.Job{Logger: (&files_sdk.Config{}).Logger()}).Init()
		job.Ignore, _ = ignore.New()
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
		job.Ignore, _ = ignore.New([]string{"*.css"}...)
		uploadStatus.localPath = "main.css"
		assert.Equal(true, skipOrIgnore(uploadStatus, false))

		uploadStatus.localPath = "main.php"
		assert.Equal(false, skipOrIgnore(uploadStatus, false))

		job.Ignore, _ = ignore.New([]string{"*.css", "*.php"}...)
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
							return client.Uploader(UploaderParams{LocalPath: source, RemotePath: destination, Config: config})
						},
						args: []string{source, destination},
					}
				})
				r.Stop()
			})
		}
	})
}

type ReaderWithOutLen struct {
	buffer *bytes.Buffer
}

func (r ReaderWithOutLen) Read(p []byte) (n int, err error) {
	return r.buffer.Read(p)
}

type ReaderAtWithOutLen struct {
	buffer *bytes.Buffer
}

func (r ReaderAtWithOutLen) ReadAt(p []byte, off int64) (n int, err error) {
	return bytes.NewReader(r.buffer.Bytes()).ReadAt(p, off)
}

func TestUploadReader(t *testing.T) {
	t.Run("reader with nil size", func(t *testing.T) {
		server := (&FakeDownloadServer{T: t}).Do()
		defer server.Shutdown()
		server.MockFiles["reader-no-size.txt"] = mockFile{File: files_sdk.File{Size: 5}}
		client := server.Client()
		u, err := client.UploadWithResume(
			func(io uploadIO) (uploadIO, error) {
				io.PartSizes = []int64{2, 4, 8, 16, 32}
				return io, nil
			},
			UploadWithReader(ReaderWithOutLen{buffer: bytes.NewBufferString("Hello")}),
			UploadWithDestinationPath("reader-no-size.txt"),
			UploadWithManager(lib.NewConstrainedWorkGroup(2)),
		)

		require.NoError(t, err)

		assert.Equal(t, server.TrackRequest["/upload/*path"], []string{"/upload/reader-no-size.txt?part_number=1", "/upload/reader-no-size.txt?part_number=2"})
		assert.Equal(t, "reader-no-size.txt", u.File.Path)
		assert.Equal(t, int64(5), u.Size)
		assert.Len(t, u.Parts, 0, "individual parts are not retryable with nil size")
		assert.Equal(t, "reader-no-size.txt", u.FileUploadPart.Path)
	})

	t.Run("reader with size present", func(t *testing.T) {
		server := (&FakeDownloadServer{T: t}).Do()
		defer server.Shutdown()
		server.MockFiles["reader-size.txt"] = mockFile{File: files_sdk.File{Size: 10}}
		client := server.Client()
		u, err := client.UploadWithResume(
			func(io uploadIO) (uploadIO, error) {
				io.PartSizes = []int64{2, 4, 8, 16, 32}
				return io, nil
			},
			UploadWithReader(bytes.NewBufferString("0123456789")),
			UploadWithDestinationPath("reader-size.txt"),
			UploadWithManager(lib.NewConstrainedWorkGroup(2)),
		)

		require.NoError(t, err)

		assert.Equal(t, server.TrackRequest["/upload/*path"], []string{"/upload/reader-size.txt?part_number=1", "/upload/reader-size.txt?part_number=2", "/upload/reader-size.txt?part_number=3"})
		assert.Equal(t, "reader-size.txt", u.File.Path)
		assert.Equal(t, int64(10), u.Size)
		assert.Len(t, u.Parts, 0, "individual parts are not retryable with nil size")
		assert.Equal(t, "reader-size.txt", u.FileUploadPart.Path)
	})

	t.Run("io.ReaderAt and no size", func(t *testing.T) {
		server := (&FakeDownloadServer{T: t}).Do()
		defer server.Shutdown()
		server.MockFiles["reader-at_no-size.txt"] = mockFile{File: files_sdk.File{Size: 10}}
		client := server.Client()
		u, err := client.UploadWithResume(
			func(io uploadIO) (uploadIO, error) {
				io.PartSizes = []int64{2, 4, 8, 16, 32}
				return io, nil
			},
			UploadWithReaderAt(ReaderAtWithOutLen{buffer: bytes.NewBufferString("0123456789")}),
			UploadWithDestinationPath("reader-at_no-size.txt"),
			UploadWithManager(lib.NewConstrainedWorkGroup(2)),
		)

		require.NoError(t, err)

		assert.Equal(t, server.TrackRequest["/upload/*path"], []string{"/upload/reader-at_no-size.txt?part_number=1", "/upload/reader-at_no-size.txt?part_number=2", "/upload/reader-at_no-size.txt?part_number=3"})
		assert.Equal(t, "reader-at_no-size.txt", u.File.Path)
		assert.Equal(t, int64(10), u.Size)
		assert.Len(t, u.Parts, 0, "individual parts are not retryable with nil size")
		assert.Equal(t, "reader-at_no-size.txt", u.FileUploadPart.Path)
	})

	t.Run("io.ReaderAt and size", func(t *testing.T) {
		server := (&FakeDownloadServer{T: t}).Do()
		defer server.Shutdown()
		server.MockFiles["reader-at-size.txt"] = mockFile{File: files_sdk.File{Size: 10}}
		client := server.Client()
		u, err := client.UploadWithResume(
			func(io uploadIO) (uploadIO, error) {
				io.PartSizes = []int64{2, 4, 8, 16, 32}
				return io, nil
			},
			UploadWithReaderAt(bytes.NewReader(bytes.NewBufferString("0123456789").Bytes())),
			UploadWithDestinationPath("reader-at-size.txt"),
			UploadWithSize(10),
			UploadWithManager(lib.NewConstrainedWorkGroup(2)),
		)

		require.NoError(t, err)
		expectation := []string{"/upload/reader-at-size.txt?part_number=1", "/upload/reader-at-size.txt?part_number=2", "/upload/reader-at-size.txt?part_number=3"}
		slices.Sort(expectation)
		slices.Sort(server.TrackRequest["/upload/*path"])
		assert.Equal(t, expectation, server.TrackRequest["/upload/*path"])
		assert.Equal(t, "reader-at-size.txt", u.File.Path)
		assert.Equal(t, int64(10), u.Size)
		assert.Len(t, u.Parts, 3, "individual parts are not retryable with nil size")
		assert.Equal(t, "reader-at-size.txt", u.FileUploadPart.Path)
	})

	t.Run("io.ReaderAt and size with resume", func(t *testing.T) {
		server := (&FakeDownloadServer{T: t}).Do()
		defer server.Shutdown()
		server.MockFiles["reader-at-size.txt"] = mockFile{File: files_sdk.File{Size: 10}}
		progressMutex := sync.Mutex{}
		bytesUploaded := []int64{0, 0}
		ctx, cancel := context.WithCancel(context.Background())
		client := server.Client()
		u, err := client.UploadWithResume(
			func(io uploadIO) (uploadIO, error) {
				io.PartSizes = []int64{2, 4, 8, 16, 32}
				return io, nil
			},
			UploadWithReaderAt(bytes.NewReader(bytes.NewBufferString("0123456789").Bytes())),
			UploadWithDestinationPath("reader-at-size.txt"),
			UploadWithSize(10),
			UploadWithManager(lib.NewConstrainedWorkGroup(1)),
			UploadWithContext(ctx),
			UploadRewindAllProgressOnFailure(),
			UploadWithProgress(func(i int64) {
				progressMutex.Lock()
				defer progressMutex.Unlock()
				bytesUploaded[0] += i

				if bytesUploaded[0] > 5 {
					cancel()
				}
			}),
		)
		cancel()
		require.ErrorIs(t, err, context.Canceled)

		assert.Len(t, u.Parts, 3)
		assert.Equal(t, int64(0), bytesUploaded[0])
		assert.Equal(t, "reader-at-size.txt", u.FileUploadPart.Path)

		// Retry
		ctx, cancel = context.WithCancel(context.Background())
		u, err = client.UploadWithResume(
			func(io uploadIO) (uploadIO, error) {
				io.PartSizes = []int64{2, 4, 8, 16, 32}
				return io, nil
			},
			UploadWithReaderAt(bytes.NewReader(bytes.NewBufferString("0123456789").Bytes())),
			UploadWithDestinationPath("reader-at-size.txt"),
			UploadWithSize(10),
			UploadWithManager(lib.NewConstrainedWorkGroup(2)),
			UploadWithContext(ctx),
			UploadWithProgress(func(i int64) {
				progressMutex.Lock()
				defer progressMutex.Unlock()
				bytesUploaded[1] += i
			}),
			UploadWithResume(u),
		)

		assert.ElementsMatch(t, []string{"/upload/reader-at-size.txt?part_number=1", "/upload/reader-at-size.txt?part_number=2", "/upload/reader-at-size.txt?part_number=2", "/upload/reader-at-size.txt?part_number=3"}, server.TrackRequest["/upload/*path"])
		assert.Equal(t, "reader-at-size.txt", u.File.Path)
		assert.Equal(t, int64(10), u.Size)
		assert.Len(t, u.Parts, 3)
		assert.Equal(t, bytesUploaded[1], u.Parts.SuccessfulBytes())
		assert.Equal(t, bytesUploaded[0]+bytesUploaded[1], u.Size)
		assert.Equal(t, "reader-at-size.txt", u.FileUploadPart.Path)
	})

	t.Run("missing UploadWithDestinationPath", func(t *testing.T) {
		server := (&FakeDownloadServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()

		err := client.Upload(
			UploadWithReader(ReaderWithOutLen{buffer: bytes.NewBufferString("Hello")}),
		)
		assert.Equal(t, "UploadWithDestinationPath is required", err.Error())
	})

	t.Run("missing reader", func(t *testing.T) {
		server := (&FakeDownloadServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()

		err := client.Upload(
			UploadWithDestinationPath("reader-no-size.txt"),
		)
		assert.Equal(t, "UploadWithReader or UploadWithReaderAt required", err.Error())
	})
}
