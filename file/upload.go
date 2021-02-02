package file

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go"
	file_action "github.com/Files-com/files-sdk-go/fileaction"
	folder "github.com/Files-com/files-sdk-go/folder"
	"github.com/zenthangplus/goccm"
)

type fileUpload struct {
	LocalFile   *os.File
	Destination string
	Source      string
	File        files_sdk.File
	Stat        os.FileInfo
	error
}

type UploadParams struct {
	Source           string
	Destination      string
	ProgressReporter func(source string, file files_sdk.File, progressByteCount int64, batchStats UploadBatchStats, err error)
}

type UploadProgress struct {
	Complete        bool
	progressWatcher func(int64)
}

func (u *UploadProgress) AddUploadedBytes(bytes int64) {
	if u.progressWatcher != nil {
		u.progressWatcher(bytes)
	}
}

type UploadBatchStats struct {
	LargestSize     int
	LargestFilePath int
	TotalUploads    int
	Size            int64
}

func (c *Client) UploadFolderOrFile(params *UploadParams) ([]files_sdk.File, error) {
	fi, err := os.Stat(params.Source)
	if err != nil {
		return []files_sdk.File{}, err
	}

	var files []files_sdk.File

	if fi.IsDir() {
		fileUploads, err := c.UploadFolder(params)
		if err != nil {
			return files, err
		}
		for _, file := range fileUploads {
			if file.error != nil {
				return files, file.error
			}
			files = append(files, file.File)
		}
	} else {
		file, err := c.UploadFile(params)
		if err != nil {
			return files, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (c *Client) UploadFolder(params *UploadParams) ([]fileUpload, error) {
	var uploadFiles []fileUpload
	var largestSize int64
	var largestFilePath int
	localFolderPath := params.Source
	destinationRootPath := params.Destination
	directoriesToCreate := make(map[string]fileUpload)
	var TotalSize int64
	addUploads := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		dir, filename := filepath.Split(path)

		if localFolderPath == dir {
			return nil
		}

		if filename == ".DS_Store" {
			return nil
		}

		if info.Size() > largestSize {
			largestSize = info.Size()
		}

		if len(path) > largestFilePath {
			largestFilePath = len(path)
		}

		var destination string
		baseDestination := filepath.Clean(strings.Replace(path, localFolderPath, "", 1))
		baseDestination = strings.TrimLeft(baseDestination, "/")
		if destinationRootPath == "" {
			destination = baseDestination
		} else {
			destination = filepath.Join(destinationRootPath, baseDestination)
		}

		if destination == "." {
			destination = filename
		}

		file := fileUpload{File: files_sdk.File{Path: destination, Size: int(info.Size())}, Source: path, Destination: destination}
		if file.isDir() {
			file.File.Type = "directory"
			directoriesToCreate[destination] = file
		} else {
			TotalSize += info.Size()
			file.File.Type = "file"
			uploadFiles = append(uploadFiles, file)
		}
		return nil
	}
	err := filepath.Walk(localFolderPath, addUploads)

	if err != nil {
		return uploadFiles, err
	}

	if len(uploadFiles) == 0 {
		file := fileUpload{Source: localFolderPath}
		if !file.isDir() {
			if addUploads(localFolderPath, file.Stat, nil) != nil {
				return uploadFiles, err
			}
		}
	}

	if destinationRootPath != "" {
		folderClient := folder.Client{Config: c.Config}
		_, err := folderClient.Create(files_sdk.FolderCreateParams{Path: filepath.Clean(destinationRootPath)})
		if err != nil && (err).(files_sdk.ResponseError).ErrorMessage != "The destination exists." {
			return uploadFiles, err
		}
	}

	batchStatus := UploadBatchStats{LargestSize: int(largestSize), LargestFilePath: largestFilePath, TotalUploads: len(uploadFiles), Size: TotalSize}
	someMapMutex := sync.RWMutex{}
	goc := goccm.New(c.Config.MaxConcurrentConnections())

	for _, uploadFile := range uploadFiles {
		go func(uploadFile fileUpload) {
			goc.Wait()
			progress := UploadProgress{}
			progress.progressWatcher = func(bytesCount int64) {
				if params.ProgressReporter == nil {
					return
				}
				params.ProgressReporter(uploadFile.Source, uploadFile.File, bytesCount, batchStatus, uploadFile.error)
			}
			dir, _ := filepath.Split(uploadFile.File.Path)
			someMapMutex.RLock()
			dirFile, ok := directoriesToCreate[filepath.Clean(dir)]
			someMapMutex.RUnlock()
			progress.progressWatcher(0)
			if ok {
				maybeCreateFolder(dirFile)
				if dirFile.error != nil {
					uploadFile.error = dirFile.error
					progress.progressWatcher(0)
				}
				someMapMutex.Lock()
				delete(directoriesToCreate, filepath.Clean(dir))
				someMapMutex.Unlock()
			}
			localFile, err := os.Open(uploadFile.Source)
			defer func() {
				localFile.Close()
				goc.Done()
			}()
			if err != nil {
				uploadFile.error = err
				progress.progressWatcher(0)
				return
			}

			file, err := c.Upload(localFile, uploadFile.File.Path, &progress)
			if err != nil {
				uploadFile.error = err
				progress.progressWatcher(0)
			}
			uploadFile.File = file
		}(uploadFile)
	}
	goc.WaitAllDone()

	return uploadFiles, err
}

func maybeCreateFolder(file fileUpload) {
	createdFolder, err := folder.Create(files_sdk.FolderCreateParams{Path: file.Destination + "/"})

	if err != nil && (err).(files_sdk.ResponseError).ErrorMessage != "The destination exists." {
		file.error = err
	} else {
		file.File = createdFolder
	}
}

func (u *fileUpload) isDir() bool {
	fi, err := os.Stat(u.Source)
	if err != nil {
		u.error = err
		return false
	}
	u.Stat = fi

	if fi.IsDir() {
		return true
	} else {
		return false
	}
}

func (c *Client) UploadFile(params *UploadParams) (files_sdk.File, error) {
	destination := params.Destination
	if params.Destination == "" {
		_, fileName := filepath.Split(params.Source)
		destination = fileName
	}
	fi, err := os.Stat(params.Source)
	localFile, err := os.Open(params.Source)
	if err != nil {
		return files_sdk.File{}, err
	}
	defer localFile.Close()
	progress := UploadProgress{}
	progress.progressWatcher = func(bytesCount int64) {
		if params.ProgressReporter == nil {
			return
		}
		params.ProgressReporter(
			params.Source,
			files_sdk.File{Size: int(fi.Size()), Path: destination, Type: "file"},
			bytesCount,
			UploadBatchStats{Size: fi.Size(), LargestSize: int(fi.Size()), LargestFilePath: len(params.Destination), TotalUploads: 1},
			nil,
		)
	}
	return c.Upload(localFile, destination, &progress)
}

func UploadFile(params *UploadParams) (files_sdk.File, error) {
	return (&Client{}).UploadFile(params)
}

func (c *Client) Upload(source io.Reader, destination string, progress *UploadProgress) (files_sdk.File, error) {
	upload, etags, bytesWritten, err := c.uploadChunks(source, destination, progress)
	if err != nil {
		return files_sdk.File{}, err
	}
	file, err := c.Create(files_sdk.FileCreateParams{
		ProvidedMtime: time.Now(),
		EtagsParam:    etags,
		Action:        "end",
		Size:          bytesWritten,
		Path:          destination,
		Ref:           upload.Ref,
	})
	if err != nil {
		return file, err
	}
	return file, nil
}

type ProgressReader struct {
	io.Reader
	*UploadProgress
}

func (p *ProgressReader) Read(b []byte) (n int, err error) {
	n, err = p.Reader.Read(b)
	p.UploadProgress.AddUploadedBytes(int64(n))
	return n, err
}

func (c *Client) uploadChunks(reader io.Reader, path string, progress *UploadProgress) (files_sdk.FileUploadPart, []files_sdk.EtagsParam, int, error) {
	bytesWritten := 0
	etags := make([]files_sdk.EtagsParam, 0)
	beginUpload := files_sdk.FileActionBeginUploadParams{}
	upload := files_sdk.FileUploadPart{}
	for {
		beginUpload.Path = path
		beginUpload.Part = upload.PartNumber + 1
		beginUpload.Ref = upload.Ref
		fileActionClient := file_action.Client{Config: c.Config}
		uploads, err := fileActionClient.BeginUpload(beginUpload)
		if err != nil {
			return upload, etags, 0, err
		}
		upload = uploads[0]

		partSizeBuffer := make([]byte, upload.Partsize)
		bytesRead, err := reader.Read(partSizeBuffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return upload, etags, 0, err
		}

		readBytesBuffer := bytes.Buffer{}
		readBytesBuffer.Grow(bytesRead)
		readBytesBuffer.Write(partSizeBuffer[:bytesRead])
		bytesWritten += bytesRead
		progressReader := ProgressReader{Reader: &readBytesBuffer, UploadProgress: progress}
		headers := http.Header{}
		headers.Add("Content-Length", strconv.FormatInt(int64(bytesRead), 10))
		res, err := files_sdk.CallRaw(upload.HttpMethod, c.Config, upload.UploadUri, nil, &progressReader, &headers)
		if err != nil {
			return upload, etags, 0, err
		}
		etags = append(
			etags,
			files_sdk.EtagsParam{
				Etag: res.Header.Get("ETag"),
				Part: strconv.FormatInt(int64(upload.PartNumber), 10),
			},
		)
	}

	return upload, etags, bytesWritten, nil
}

func Upload(source io.Reader, destination string, progress *UploadProgress) (files_sdk.File, error) {
	return (&Client{}).Upload(source, destination, progress)
}
