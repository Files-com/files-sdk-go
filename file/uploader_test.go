package file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"math"
	"net/http"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/directory"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/ignore"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/gin-gonic/gin"
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

func (m *MockUploader) CreateFolder(files_sdk.FolderCreateParams, ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return files_sdk.File{}, nil
}

func Test_excludeFile(t *testing.T) {
	var progressReportError error

	init := func() (*UploadStatus, fstest.MapFS, *Job) {
		job := (&Job{Logger: files_sdk.Config{}.Init().Logger}).Init()
		job.Ignore, _ = ignore.New()
		job.Params = UploaderParams{}
		uploadStatus := &UploadStatus{job: job, Mutex: &sync.RWMutex{}, file: files_sdk.File{Path: "test"}, remotePath: "test"}
		uploadStatus.Job().Add(uploadStatus)
		uploader := &MockUploader{}
		uploadStatus.Uploader = uploader
		uploadStatus.job.EventsReporter = CreateReporter(func(s JobFile) {
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
		assert.False(excludeFile(uploadStatus, false))
	})

	t.Run("when sizes don't match", func(t *testing.T) {
		uploadStatus.Sync = true
		mockFs["test"] = &fstest.MapFile{
			Sys: files_sdk.File{Size: 9},
		}
		uploadStatus.file.Size = 10
		assert.False(excludeFile(uploadStatus, false))
		assert.Equal(nil, progressReportError)
	})

	t.Run("when sizes do match", func(t *testing.T) {
		mockFs["test"] = &fstest.MapFile{
			Sys: files_sdk.File{Size: 10},
		}
		uploadStatus.file.Size = 10
		assert.True(excludeFile(uploadStatus, false))
	})

	t.Run("There is no server version", func(t *testing.T) {
		delete(mockFs, "test")
		assert.False(excludeFile(uploadStatus, false))
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
		assert.True(excludeFile(uploadStatus, false))
		uploadStatus = &oldUploadStatus
	})

	t.Run("when transfer is a single file", func(t *testing.T) {
		uploadStatus.job.Type = directory.File
		uploadStatus.Sync = true
		mockFs["test"] = &fstest.MapFile{
			Sys: files_sdk.File{Size: 10},
		}
		uploadStatus.file.Size = 10
		assert.True(excludeFile(uploadStatus, false))
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
		assert.True(excludeFile(uploadStatus, false))
		uploadStatus.Sync = false
	})

	t.Run("Ignore files", func(t *testing.T) {
		job.Ignore, _ = ignore.New([]string{"*.css"}...)
		uploadStatus.localPath = "main.css"
		assert.True(excludeFile(uploadStatus, false))

		uploadStatus.localPath = "main.php"
		assert.False(excludeFile(uploadStatus, false))

		job.Ignore, _ = ignore.New([]string{"*.css", "*.php"}...)
		uploadStatus.localPath = "main.css"
		assert.True(excludeFile(uploadStatus, false))
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
				skip := excludeFile(uploadStatus, true)
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

	t.Run("No overwrite file exists", func(t *testing.T) {
		uploadStatus, mockFs, _ := init()
		mockFs["test"] = &fstest.MapFile{
			Sys: files_sdk.File{},
		}
		uploadStatus.NoOverwrite = true
		assert.True(excludeFile(uploadStatus, false))
		assert.Equal(status.FileExists, uploadStatus.Status())
	})

	t.Run("No overwrite file does not exists", func(t *testing.T) {
		uploadStatus, _, _ := init()
		uploadStatus.NoOverwrite = true
		uploadStatus.status = status.Queued
		assert.False(excludeFile(uploadStatus, false))
		assert.Equal(status.Queued, uploadStatus.Status())
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
				lib.BuildPathSpecTest(t, mutex, tt, sourceFs, destinationFs, func(args lib.PathSpecArgs) lib.Cmd {
					return &CmdRunner{
						run: func() *Job {
							return client.Uploader(UploaderParams{LocalPath: args.Src, RemotePath: args.Dest, config: config, PreserveTimes: args.PreserveTimes})
						},
						args: []string{args.Src, args.Dest, "--times", fmt.Sprintf("%v", args.PreserveTimes)},
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
		server := (&MockAPIServer{T: t}).Do()
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
		server := (&MockAPIServer{T: t}).Do()
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
		server := (&MockAPIServer{T: t}).Do()
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
		server := (&MockAPIServer{T: t}).Do()
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
		firstTry := func(t *testing.T, filename string) (context.Context, context.CancelFunc, *Client, *MockAPIServer, UploadResumable) {
			ctx, cancel := context.WithCancel(context.Background())
			server := (&MockAPIServer{T: t}).Do()
			client := server.Client()

			server.MockFiles[filename] = mockFile{File: files_sdk.File{Size: 10}}
			progressMutex := sync.Mutex{}
			bytesUploaded := []int64{0, 0}

			firstTry, err := client.UploadWithResume(
				func(io uploadIO) (uploadIO, error) {
					io.PartSizes = []int64{2, 4, 8, 16, 32}
					return io, nil
				},
				UploadWithReaderAt(bytes.NewReader(bytes.NewBufferString("0123456789").Bytes())),
				UploadWithDestinationPath(filename),
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

			assert.Len(t, firstTry.Parts, 3)
			assert.Equal(t, int64(2), firstTry.Parts.SuccessfulBytes())
			assert.Equal(t, int64(0), bytesUploaded[0])
			assert.Equal(t, filename, firstTry.FileUploadPart.Path)
			server.CloseClientConnections()
			ctx, cancel = context.WithCancel(context.Background())
			// The last request might still be processing
			time.Sleep(10 * time.Millisecond)
			return ctx, func() {
					server.CloseClientConnections()
					cancel()
				}, client,
				server,
				firstTry
		}
		// Retry
		retry := func(ctx context.Context, u UploadResumable, client *Client) (UploadResumable, error) {
			return client.UploadWithResume(
				func(io uploadIO) (uploadIO, error) {
					io.PartSizes = []int64{2, 4, 8, 16, 32}
					return io, nil
				},
				UploadWithReaderAt(bytes.NewReader(bytes.NewBufferString("0123456789").Bytes())),
				UploadWithDestinationPath(u.FileUploadPart.Path),
				UploadWithSize(10),
				UploadWithManager(lib.NewConstrainedWorkGroup(2)),
				UploadWithContext(ctx),
				UploadWithResume(u),
			)
		}

		t.Run("native", func(t *testing.T) {
			ctx, cancel, client, server, resume := firstTry(t, "native-file")
			defer cancel()

			server.traceMutex.Lock()
			server.TrackRequest = make(map[string][]string)
			server.traceMutex.Unlock()
			var beginUploadRequests []files_sdk.FileBeginUploadParams

			server.MockRoute("/api/rest/v1/file_actions/begin_upload/native-file", func(ctx *gin.Context, model interface{}) bool {
				beginUploadRequests = append(beginUploadRequests, model.(files_sdk.FileBeginUploadParams))
				return false
			})

			u, err := retry(ctx, resume, client)
			require.NoError(t, err)

			assert.ElementsMatch(t, []string{"/api/rest/v1/file_actions/begin_upload/native-file"}, server.TrackRequest["/api/rest/v1/file_actions/begin_upload/*path"], "only requests part 3")
			assert.ElementsMatch(t, []string{"/upload/native-file?part_number=2", "/upload/native-file?part_number=3"}, server.TrackRequest["/upload/*path"], "1 already succeed rest are uploaded")
			assert.Equal(t, "native-file", u.File.Path)
			assert.Equal(t, int64(10), u.Size)
			assert.Len(t, u.Parts, 3)
			assert.Equal(t, "native-file", u.FileUploadPart.Path)
		})

		t.Run("remote_mount", func(t *testing.T) {
			ctx, cancel, client, server, resume := firstTry(t, "remote_mount-file")
			defer cancel()

			resume.FileUploadPart.ParallelParts = lib.Bool(false)
			server.traceMutex.Lock()
			server.TrackRequest = make(map[string][]string)
			server.traceMutex.Unlock()
			server.MockRoute("/api/rest/v1/file_actions/begin_upload/remote_mount-file", func(ctx *gin.Context, model interface{}) bool {
				file := model.(files_sdk.FileBeginUploadParams)
				if file.Part == 0 {
					file.Part = 1
				}
				path := strings.TrimPrefix(ctx.Param("path"), "/")
				ctx.JSON(http.StatusOK, files_sdk.FileUploadPartCollection{
					files_sdk.FileUploadPart{
						HttpMethod:    "POST",
						Path:          path,
						UploadUri:     fmt.Sprintf("%v?part_number=%v", lib.UrlJoinNoEscape(server.URL, "upload", path), file.Part),
						ParallelParts: lib.Bool(false),
						Expires:       time.Now().Add(time.Hour).Format(time.RFC3339),
						PartNumber:    file.Part,
					},
				})
				return true
			})

			u, err := retry(ctx, resume, client)
			require.NoError(t, err)

			assert.ElementsMatch(t, []string{"/api/rest/v1/file_actions/begin_upload/remote_mount-file"}, server.TrackRequest["/api/rest/v1/file_actions/begin_upload/*path"], "upload is invalided because of ParallelParts")
			assert.ElementsMatch(t, []string{"/upload/remote_mount-file?part_number=1", "/upload/remote_mount-file?part_number=2", "/upload/remote_mount-file?part_number=3"}, server.TrackRequest["/upload/*path"], "all parts are uploaded")
			assert.Equal(t, "remote_mount-file", u.File.Path)
			assert.Equal(t, int64(10), u.Size)
			assert.Len(t, u.Parts, 3)
			assert.Equal(t, int64(10), u.Parts.SuccessfulBytes())
			assert.Equal(t, "remote_mount-file", u.FileUploadPart.Path)
		})
	})

	t.Run("missing UploadWithDestinationPath", func(t *testing.T) {
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()

		err := client.Upload(
			UploadWithReader(ReaderWithOutLen{buffer: bytes.NewBufferString("Hello")}),
		)
		assert.Equal(t, "UploadWithDestinationPath is required", err.Error())
	})

	t.Run("missing reader", func(t *testing.T) {
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		client := server.Client()

		err := client.Upload(
			UploadWithDestinationPath("reader-no-size.txt"),
		)
		assert.Equal(t, "UploadWithReader or UploadWithReaderAt required", err.Error())
	})

	accessDenied := CustomResponse{
		Status:      http.StatusForbidden,
		ContentType: "application/xml",
		Body:        []byte(`<?xml version="1.0" encoding="UTF-8"?> <Error><Code>AccessDenied</Code><Message>Request has expired</Message><X-Amz-Expires>900</X-Amz-Expires><Expires>2023-09-06T05:27:01Z</Expires><ServerTime>2023-09-06T05:27:21Z</ServerTime><RequestId>DCZ7NV6P08Y6SKY2</RequestId><HostId>a0ww8xPnO34ZC2to9wizy501VJcZicTFKdohzq5P7SArZuXJ7cCo6GpJbUXITjkyFHNPla8Sd1U=</HostId></Error>`),
	}

	t.Run("socket connection error", func(t *testing.T) {
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()

		r, w := io.Pipe()
		defer w.Close()
		defer r.Close()
		go func() {
			for {
				// Simulate a slow disk access
				w.Write([]byte("1"))
				time.Sleep(time.Millisecond * 10)
			}
		}()

		client := server.Client()
		u, err := client.UploadWithResume(
			UploadWithSize(1024),
			UploadWithReader(r),
			UploadWithDestinationPath("file.bak"),
			UploadWithManager(lib.NewConstrainedWorkGroup(1)),
		)

		assert.Equal(t, lib.S3Error{Message: "Your socket connection to the server was not read from or written to within the timeout period. Idle connections will be closed.", Code: "RequestTimeout"}.Error(), err.Error())
		assert.Equal(t, "file.bak", u.FileUploadPart.Path)
	})

	t.Run("socket connection error recovers after retry", func(t *testing.T) {
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()

		client := server.Client()
		u, err := client.UploadWithResume(
			UploadWithSize(1024),
			UploadWithReaderAt(bytes.NewReader(make([]byte, 1024))),
			UploadWithDestinationPath("file.bak"),
			UploadWithManager(lib.NewConstrainedWorkGroup(1)),
		)

		assert.NoError(t, err)
		assert.Equal(t, "file.bak", u.FileUploadPart.Path)
	})

	t.Run("request expired error", func(t *testing.T) {
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		server.MockRoute("/upload/file.bak", func(c *gin.Context, _ interface{}) bool {
			c.Data(accessDenied.Status, accessDenied.ContentType, accessDenied.Body)
			return true
		})

		client := server.Client()
		u, err := client.UploadWithResume(
			UploadWithReader(bytes.NewBufferString("0123456789")),
			UploadWithDestinationPath("file.bak"),
			UploadWithManager(lib.NewConstrainedWorkGroup(1)),
		)

		assert.Equal(t, lib.S3Error{Message: "Request has expired", Code: "AccessDenied"}.Error(), err.Error())
		assert.Equal(t, "file.bak", u.FileUploadPart.Path)
	})

	t.Run("File Upload Not Found", func(t *testing.T) {
		server := (&MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		var partCount int
		server.MockRoute("/api/rest/v1/file_actions/begin_upload/file.bak", func(c *gin.Context, model interface{}) bool {
			partCount += 1
			if beginUpload, ok := model.(files_sdk.FileBeginUploadParams); ok {
				if beginUpload.Part == 0 {
					return false
				}
			}

			c.Data(http.StatusNotFound, "text", []byte("File Upload Not Found"))
			return true
		})

		client := server.Client()
		u, err := client.UploadWithResume(
			func(io uploadIO) (uploadIO, error) {
				io.PartSizes = []int64{2, 4, 8, 16, 32}
				return io, nil
			},
			UploadWithReader(bytes.NewBufferString("0123456789")),
			UploadWithDestinationPath("file.bak"),
			UploadWithManager(lib.NewConstrainedWorkGroup(1)),
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "File Upload Not Found", "it invalidates any resuming")
		assert.Len(t, u.Parts, 0)
		assert.Equal(t, partCount, 3)
	})
}
