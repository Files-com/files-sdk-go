package file

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/chilts/sid"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

var serverPort = 50001

type randomReader struct {
	n int
}

func (r *randomReader) Read(p []byte) (int, error) {
	if r.n <= 0 {
		return 0, io.EOF
	}
	if len(p) > r.n {
		p = p[:r.n]
	}

	_, err := rand.Read(p)
	if err != nil {
		return 0, err
	}

	r.n -= len(p)
	return len(p), nil
}

type FakeDownloadServer struct {
	router *gin.Engine
	Addr   string
	*http.Server
	Port      int
	downloads *lib.Map[download]
	MockFiles map[string]mockFile
	*testing.T
	TrackRequest map[string][]string
	traceMutex   *sync.Mutex
}

type download struct {
	Id string
	mockFile
	Requests *lib.Map[files_sdk.ResponseError]
}

func (d download) init() download {
	d.Requests = &lib.Map[files_sdk.ResponseError]{}
	return d
}

type mockFile struct {
	files_sdk.File
	RealSize *int64
	SizeTrust
	ForceRequestStatus  string
	ForceRequestMessage string
	ServerBytesSent     *int64
	MaxConnections      int
	MaxConnectionsMutex *sync.Mutex
}

func (m mockFile) Completed() string {
	if m.ForceRequestStatus != "" {
		return m.ForceRequestStatus
	}
	return "completed"
}

type TestLogger struct {
	*testing.T
}

func (t TestLogger) Printf(format string, args ...any) {
	t.T.Logf(format, args...)
}

func (t TestLogger) Write(p []byte) (n int, err error) {
	t.T.Log(string(p))
	return len(p), nil
}

func (f FakeDownloadServer) Do() FakeDownloadServer {
	f.MockFiles = make(map[string]mockFile)
	f.TrackRequest = make(map[string][]string)
	f.traceMutex = &sync.Mutex{}
	f.downloads = &lib.Map[download]{}
	f.router = gin.New()
	f.router.Use(gin.LoggerWithWriter(TestLogger{f.T}))
	f.Routes()
	serverPort += 1
	f.Port = serverPort
	f.Server = &http.Server{
		Addr:    fmt.Sprintf("localhost:%v", f.Port),
		Handler: f.router,
	}
	go func() {
		var err error
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		for {
			select {
			case <-ctx.Done():
				log.Fatalf("listen: %s\n", err)
				return
			default:
				err = f.Server.ListenAndServe()
				if err != nil && strings.Contains(err.Error(), "address already in use") {
					serverPort += 1
					f.Port = serverPort
					f.Server = &http.Server{
						Addr:    fmt.Sprintf("localhost:%v", f.Port),
						Handler: f.router,
					}
				} else if err != nil && err != http.ErrServerClosed {
					log.Fatalf("listen: %s\n", err)
					return
				} else {
					return
				}
			}
		}
	}()
	return f
}

func (f FakeDownloadServer) Client() *Client {
	client := &Client{}
	httpClient := http.Client{}
	httpClient.Transport = &CustomTransport{Addr: f.Server.Addr}
	client.Config.SetHttpClient(&httpClient)
	client.Config.SetLogger(TestLogger{f.T})
	return client
}

func (f FakeDownloadServer) GetFile(file mockFile) (r io.Reader, contentLengthOk bool, contentLength int64, realSize int64, err error) {
	if file.SizeTrust == NullSizeTrust || file.SizeTrust == TrustedSizeValue {
		contentLengthOk = true
	}

	contentLength = file.File.Size

	if file.RealSize != nil {
		realSize = *file.RealSize
	} else {
		realSize = contentLength
	}
	r = &randomReader{int(realSize)}
	return
}

func (f FakeDownloadServer) trackRequest(c *gin.Context) {
	f.traceMutex.Lock()
	defer f.traceMutex.Unlock()
	f.TrackRequest[c.FullPath()] = append(f.TrackRequest[c.FullPath()], c.Request.URL.String())
}

func (f FakeDownloadServer) GetRouter() *gin.Engine {
	return f.router
}

func (f FakeDownloadServer) Routes() {
	f.router.GET("/api/rest/v1/files/*path", func(c *gin.Context) {
		f.trackRequest(c)
		path := strings.TrimPrefix(c.Param("path"), "/")

		file, ok := f.MockFiles[path]
		if ok {
			if file.Path == "" {
				file.Path = path
				file.DisplayName = filepath.Base(path)
			}
			downloadId := sid.IdHex()
			f.downloads.Store(downloadId, download{Id: downloadId, mockFile: file}.init())
			file.DownloadUri = lib.UrlJoinNoEscape("http://localhost:8080/download", downloadId)

			if file.MaxConnections != 0 {
				file.MaxConnectionsMutex = &sync.Mutex{}
			}

			c.JSON(http.StatusOK, file.File)
		} else {
			c.JSON(http.StatusNotFound, nil)
		}
	})
	f.router.GET("/api/rest/v1/folders/*path", func(c *gin.Context) {
		f.trackRequest(c)
		path := strings.TrimPrefix(c.Param("path"), "/")

		var files []files_sdk.File
		for k, v := range f.MockFiles {
			dir, _ := filepath.Split(k)
			if lib.NormalizeForComparison(filepath.Clean(path)) == lib.NormalizeForComparison(filepath.Clean(dir)) {
				if v.Path == "" {
					v.Path = k
					v.DisplayName = filepath.Base(k)
				}
				files = append(files, v.File)
			}
		}

		if len(files) > 0 {
			c.JSON(http.StatusOK, files)
		} else {
			c.JSON(http.StatusNotFound, nil)
		}
	})
	f.router.GET("/api/rest/v1/file_actions/metadata/*path", func(c *gin.Context) {
		f.trackRequest(c)
		path := strings.TrimPrefix(c.Param("path"), "/")
		file, ok := f.MockFiles[path]
		if ok {
			if file.Path == "" {
				file.Path = path
				file.DisplayName = filepath.Base(path)
			}
			c.JSON(http.StatusOK, file.File)
		} else {
			c.JSON(http.StatusNotFound, nil)
		}

	})
	f.router.GET("/download/:download_id/:download_request_id", func(c *gin.Context) {
		f.trackRequest(c)
		downloadId := c.Param("download_id")
		downloadJob, downloadOk := f.downloads.Load(downloadId)
		if !downloadOk {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		downloadRequestJob, requestOk := downloadJob.Requests.Load(c.Param("download_request_id"))
		if requestOk {
			c.JSON(http.StatusOK, downloadRequestJob)
		} else {
			c.JSON(http.StatusNotFound, nil)
		}

	})
	f.router.GET("/download/:download_id", func(c *gin.Context) {
		f.trackRequest(c)
		downloadJob, ok := f.downloads.Load(c.Param("download_id"))
		if !ok {
			c.JSON(http.StatusNotFound, nil)
			return
		}

		if downloadJob.mockFile.MaxConnectionsMutex != nil {
			downloadJob.mockFile.MaxConnectionsMutex.Lock()
		}

		start, end, okRange := rangeValue(c.Request.Header)

		reader, contentLengthOk, contentLength, realSize, err := f.GetFile(downloadJob.mockFile)
		if err != nil {
			panic(err)
		}
		status := http.StatusOK
		if okRange {
			if realSize < int64(start) {
				reader = &randomReader{0}
			} else {
				reader = &randomReader{(lo.Min[int]([]int{int(realSize - 1), end}) - start) + 1}
			}
			status = http.StatusPartialContent
		}
		downloadRequestId := sid.IdHex()
		if downloadJob.mockFile.MaxConnections == 0 {
			c.Header("X-Files-Max-Connections", "*")
		} else {
			c.Header("X-Files-Max-Connections", fmt.Sprintf("%v", downloadJob.mockFile.MaxConnections))
		}

		c.Header("X-Files-Download-Request-Id", downloadRequestId)
		responseError := files_sdk.ResponseError{ErrorMessage: downloadJob.ForceRequestMessage}
		extraHeaders := map[string]string{}
		if contentLengthOk {
			if okRange && contentLength < int64(end) {
				c.Status(http.StatusBadRequest)
			}

			if okRange {
				extraHeaders["Content-Range"] = fmt.Sprintf("%v-%v/%v", start, end, contentLength)
				contentLength = int64(end-start) + 1
			}

			c.DataFromReader(status, contentLength, "application/zip, application/octet-stream", reader, extraHeaders)
			downloadJob.Requests.Store(downloadRequestId, responseError)
			if downloadJob.mockFile.MaxConnectionsMutex != nil {
				downloadJob.mockFile.MaxConnectionsMutex.Unlock()
			}
		} else {
			finish := func() {
				if downloadJob.ServerBytesSent != nil {
					responseError.Data.BytesTransferred = *downloadJob.ServerBytesSent
				}
				downloadJob.Requests.Store(downloadRequestId, responseError)
				if downloadJob.mockFile.MaxConnectionsMutex != nil {
					downloadJob.mockFile.MaxConnectionsMutex.Unlock()
				}
			}
			if okRange {
				c.Header("Content-Range", fmt.Sprintf("%v-%v/*", start, end))
			}
			c.Status(status)
			c.Stream(func(w io.Writer) bool {
				buf := make([]byte, 1024*1024)

				n, err := reader.Read(buf)
				if err == io.EOF {
					responseError.Data.Status = downloadJob.Completed()
					finish()
					return false
				}
				if err != nil && err != io.EOF {
					responseError.Data.Status = "errored"
					finish()
					return false
				}

				wn, err := w.Write(buf[:n])
				if err != nil {
					responseError.Data.Status = "errored"
					finish()
					return false
				}

				responseError.Data.BytesTransferred += int64(wn)

				if err == io.EOF {
					responseError.Data.Status = "errored"
					finish()
					return false
				}

				return true
			})
		}
	})
	f.router.HEAD("/download/:download_id", func(c *gin.Context) {
		f.trackRequest(c)
		downloadJob, ok := f.downloads.Load(c.Param("download_id"))
		if !ok {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		_, contentLengthOk, contentLength, _, err := f.GetFile(downloadJob.mockFile)
		if err != nil {
			panic(err)
		}
		if contentLengthOk {
			c.Header("Content-Length", fmt.Sprintf("%v", contentLength))
		}
		if downloadJob.mockFile.MaxConnections == 0 {
			c.Header("X-Files-Max-Connections", "*")
		} else {
			c.Header("X-Files-Max-Connections", fmt.Sprintf("%v", downloadJob.mockFile.MaxConnections))
		}
		c.Status(http.StatusOK)
	})
}

func (f FakeDownloadServer) Shutdown() error {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return f.Server.Shutdown(ctx)
}

type CustomTransport struct {
	http.Transport
	Addr string
}

func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Host = t.Addr
	req.URL.Scheme = "http"

	return t.Transport.RoundTrip(req)
}

func rangeValue(header http.Header) (start, end int, ok bool) {
	r := header.Get("Range")
	if r == "" {
		return
	}

	r = strings.SplitN(r, "=", 2)[1]
	ranges := strings.Split(r, "-")
	var err error
	start, err = strconv.Atoi(ranges[0])
	if err != nil {
		return
	}
	end, err = strconv.Atoi(ranges[1])
	if err != nil {
		return
	}

	ok = true

	return
}
