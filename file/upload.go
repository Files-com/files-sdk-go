package file

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func (c *Client) UploadFolder(localFolderPath string, destinationRootPath *string, reporters ...func(source string, file files_sdk.File, largestSize int, largestFilePath int, totalUploads int, err error)) ([]fileUpload, error) {
	var uploadFiles []fileUpload
	var largestSize int64
	var largestFilePath int
	directoriesToCreate := make(map[string]fileUpload)
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
		if destinationRootPath == nil {
			destination = baseDestination
		} else {
			destination = filepath.Join(*destinationRootPath, baseDestination)
		}

		if destination == "." {
			destination = filename
		}

		file := fileUpload{File: files_sdk.File{Path: destination}, Source: path, Destination: destination}
		if file.isDir() {
			directoriesToCreate[destination] = file
		} else {
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

	if destinationRootPath != nil {
		folderClient := folder.Client{Config: c.Config}
		_, err := folderClient.Create(files_sdk.FolderCreateParams{Path: filepath.Clean(*destinationRootPath)})
		if err != nil && (err).(*files_sdk.ResponseError).ErrorMessage != "file or folder already exists with that name" {
			return uploadFiles, err
		}
	}
	if c.ConcurrentUploads == 0 {
		c.ConcurrentUploads = 10
	}
	goc := goccm.New(c.Config.ConcurrentUploads)

	fileChannel := make(chan fileUpload)
	for _, dir := range directoriesToCreate {
		go func(file fileUpload) {
			goc.Wait()
			maybeCreateFolder(file)
			fileChannel <- file
			goc.Done()
		}(dir)
	}

	for i := 0; i < len(directoriesToCreate); i++ {
		file := <-fileChannel
		if len(reporters) > 0 {
			reporters[0](file.Source, file.File, int(largestSize), largestFilePath, len(uploadFiles), file.error)
		}
	}
	for _, uploadFile := range uploadFiles {
		go func(uploadFile fileUpload) {
			goc.Wait()
			localFile, err := os.Open(uploadFile.Source)
			defer func() {
				localFile.Close()
				fileChannel <- uploadFile
				goc.Done()
			}()
			if err != nil {
				uploadFile.error = err
				return
			}
			file, err := c.Upload(localFile, uploadFile.File.Path)
			if err != nil {
				uploadFile.error = err
			}
			uploadFile.File = file
		}(uploadFile)
	}
	for i := 0; i < len(uploadFiles); i++ {
		file := <-fileChannel
		if len(reporters) > 0 {
			reporters[0](file.Source, file.File, int(largestSize), largestFilePath, len(uploadFiles), file.error)
		}
	}
	return uploadFiles, err
}

func maybeCreateFolder(file fileUpload) {
	createdFolder, err := folder.Create(files_sdk.FolderCreateParams{Path: file.Destination})

	if err != nil && (err).(*files_sdk.ResponseError).ErrorMessage != "file or folder already exists with that name" {
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

func (c *Client) UploadFile(path string, destination *string) (files_sdk.File, error) {
	if destination == nil {
		_, fileName := filepath.Split(path)
		destination = &fileName
	}
	localFile, err := os.Open(path)
	if err != nil {
		return files_sdk.File{}, err
	}
	defer localFile.Close()

	return c.Upload(localFile, *destination)
}

func UploadFile(path string, destination *string) (files_sdk.File, error) {
	return (&Client{}).UploadFile(path, destination)
}

func (c *Client) Upload(source io.Reader, destination string) (files_sdk.File, error) {
	upload, etags, bytesWritten, err := c.uploadChunks(source, destination)
	if err != nil {
		return files_sdk.File{}, err
	}
	return c.Create(files_sdk.FileCreateParams{
		ProvidedMtime: time.Now(),
		EtagsParam:    etags,
		Action:        "end",
		Size:          bytesWritten,
		Path:          destination,
		Ref:           upload.Ref,
	})
}

func (c *Client) uploadChunks(reader io.Reader, path string) (files_sdk.FileUploadPart, []files_sdk.EtagsParam, int, error) {
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
		headers := http.Header{}
		headers.Add("Content-Length", strconv.FormatInt(int64(bytesRead), 10))
		res, err := files_sdk.CallRaw(upload.HttpMethod, c.Config, upload.UploadUri, nil, &readBytesBuffer, &headers)
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

func Upload(source io.Reader, destination string) (files_sdk.File, error) {
	return (&Client{}).Upload(source, destination)
}
