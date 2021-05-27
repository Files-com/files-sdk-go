package file

import (
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go"
	file_action "github.com/Files-com/files-sdk-go/fileaction"
	"github.com/Files-com/files-sdk-go/folder"
	"github.com/Files-com/files-sdk-go/lib"
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
	absoluteSource, err := filepath.Abs(params.Source)
	if err != nil {
		return []files_sdk.File{}, err
	}
	fi, err := os.Stat(absoluteSource)
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
		var baseDestination string
		if localFolderPath != "." {
			baseDestination = strings.TrimPrefix(path, localFolderPath)
		} else if path != "." {
			baseDestination = path
		}
		baseDestination = strings.TrimLeft(baseDestination, "/")
		baseDestination = strings.TrimPrefix(baseDestination, "/")
		if destinationRootPath == "" {
			destination = baseDestination
		} else {
			destination = filepath.Join(destinationRootPath, baseDestination)
		}

		if destination == "." {
			destination = filename
		}

		file := fileUpload{File: files_sdk.File{Path: destination, Size: info.Size()}, Source: path, Destination: destination}
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
		responseError, ok := (err).(files_sdk.ResponseError)
		if err != nil && ok && responseError.ErrorMessage != "The destination exists." {
			return uploadFiles, err
		}
	}

	batchStatus := UploadBatchStats{LargestSize: int(largestSize), LargestFilePath: largestFilePath, TotalUploads: len(uploadFiles), Size: TotalSize}
	someMapMutex := sync.RWMutex{}
	goc := goccm.New(10)
	for _, uploadFile := range uploadFiles {
		goc.Wait()

		go func(uploadFile fileUpload) {
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

			file, err := c.Upload(localFile, uploadFile.Stat.Size(), files_sdk.FileActionBeginUploadParams{Path: uploadFile.File.Path, MkdirParents: lib.Bool(true)}, &progress)
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
	responseError, ok := (err).(files_sdk.ResponseError)
	if err != nil && ok && responseError.ErrorMessage != "The destination exists." {
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
	beginUpload := files_sdk.FileActionBeginUploadParams{}
	destination := params.Destination
	_, localFileName := filepath.Split(params.Source)
	if params.Destination == "" {
		destination = localFileName
	} else {
		_, err := c.Find(params.Destination)
		responseError, ok := err.(files_sdk.ResponseError)
		if ok && responseError.Type == "bad-request/cannot-download-directory" {
			destination = filepath.Join(params.Destination, localFileName)
		} else if ok && responseError.Type == "not-found" {
			if destination[len(destination)-1:] == "/" {
				destination = filepath.Join(params.Destination, localFileName)
				beginUpload.MkdirParents = lib.Bool(true)
			}
		} else if err != nil {
			return files_sdk.File{}, err
		}
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
			files_sdk.File{Size: fi.Size(), Path: destination, Type: "file"},
			bytesCount,
			UploadBatchStats{Size: fi.Size(), LargestSize: int(fi.Size()), LargestFilePath: len(params.Destination), TotalUploads: 1},
			nil,
		)
	}
	beginUpload.Path = destination
	return c.Upload(localFile, fi.Size(), beginUpload, &progress)
}

func UploadFile(params *UploadParams) (files_sdk.File, error) {
	return (&Client{}).UploadFile(params)
}

func (c *Client) Upload(reader io.ReaderAt, size int64, params files_sdk.FileActionBeginUploadParams, progress *UploadProgress) (files_sdk.File, error) {
	onComplete := make(chan files_sdk.EtagsParam)
	onError := make(chan error)
	bytesWritten := int64(0)
	etags := make([]files_sdk.EtagsParam, 0)
	goc := c.Config.NullConcurrencyManger()
	fileUploadPart, err := c.startUpload(params)
	if err != nil {
		return files_sdk.File{}, err
	}
	if *fileUploadPart.ParallelParts {
		goc = c.Config.ConcurrencyManger()
	}
	partReturnedError := false
	fileUploadPart.Path = params.Path

	byteOffset(
		size,
		fileUploadPart.Partsize,
		func(off int64, len int64) {
			goc.Wait()

			if partReturnedError {
				return
			}
			go func(off int64, len int64, fileUploadPart files_sdk.FileUploadPart) {
				proxyReader := &ProxyReader{
					ReaderAt: reader,
					off:      off,
					len:      len,
					onRead:   progress.AddUploadedBytes,
				}

				etag, bytesRead, err := c.createPart(proxyReader, len, fileUploadPart)
				if err != nil {
					goc.Done()
					onError <- err
					return
				}
				bytesWritten += bytesRead
				goc.Done()
				onComplete <- etag
			}(off, len, fileUploadPart)

			fileUploadPart.PartNumber += 1
		},
	)

	n := int64(0)
	for n < fileUploadPart.PartNumber-1 {
		n++
		select {
		case err := <-onError:
			partReturnedError = true
			return files_sdk.File{}, err
		case etag := <-onComplete:
			etags = append(etags, etag)
		}
	}

	return c.completeUpload(etags, bytesWritten, fileUploadPart.Path, fileUploadPart.Ref)
}

func (c *Client) startUpload(beginUpload files_sdk.FileActionBeginUploadParams) (files_sdk.FileUploadPart, error) {
	fileActionClient := file_action.Client{Config: c.Config}
	uploads, err := fileActionClient.BeginUpload(beginUpload)
	if err != nil {
		return files_sdk.FileUploadPart{}, err
	}
	return uploads[0], err
}

func (c *Client) completeUpload(etags []files_sdk.EtagsParam, bytesWritten int64, path string, ref string) (files_sdk.File, error) {
	return c.Create(files_sdk.FileCreateParams{
		ProvidedMtime: time.Now(),
		EtagsParam:    etags,
		Action:        "end",
		Path:          path,
		Ref:           ref,
		Size:          bytesWritten,
	})
}

func byteOffset(size int64, blockSize int64, callback func(off int64, len int64)) {
	off := int64(0)
	endRange := blockSize
	for {
		if off < size {
			endRange = int64(math.Min(float64(endRange), float64(size)))
			callback(off, endRange-off)
			off = endRange
			endRange = off + blockSize
		} else {
			break
		}
	}
}

func Upload(reader io.ReaderAt, size int64, beginUpload files_sdk.FileActionBeginUploadParams, progress *UploadProgress) (files_sdk.File, error) {
	return (&Client{}).Upload(reader, size, beginUpload, progress)
}

func (c *Client) createPart(reader io.ReadCloser, len int64, fileUploadPart files_sdk.FileUploadPart) (files_sdk.EtagsParam, int64, error) {
	var err error
	if fileUploadPart.PartNumber != 1 {
		fileUploadPart, err = c.startUpload(
			files_sdk.FileActionBeginUploadParams{Path: fileUploadPart.Path, Ref: fileUploadPart.Ref, Part: fileUploadPart.PartNumber},
		)
		if err != nil {
			return files_sdk.EtagsParam{}, int64(0), err
		}
	}

	headers := http.Header{}
	headers.Add("Content-Length", strconv.FormatInt(len, 10))
	res, err := files_sdk.CallRaw(
		&files_sdk.CallParams{
			Method:  fileUploadPart.HttpMethod,
			Config:  c.Config,
			Uri:     fileUploadPart.UploadUri,
			BodyIo:  reader,
			Headers: &headers,
		},
	)
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	if err != nil {
		return files_sdk.EtagsParam{}, len, err
	}

	return files_sdk.EtagsParam{
		Etag: res.Header.Get("Etag"),
		Part: strconv.FormatInt(int64(fileUploadPart.PartNumber), 10),
	}, len, nil
}
