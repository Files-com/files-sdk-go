package file

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go"
	file_action "github.com/Files-com/files-sdk-go/fileaction"
	lib "github.com/Files-com/files-sdk-go/lib"
)

type Client struct {
	files_sdk.Config
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

func Upload(source io.Reader, destination string) (files_sdk.File, error) {
	return (&Client{}).Upload(source, destination)
}

func (c *Client) DownloadToFile(params files_sdk.FileDownloadParams, filePath string) (files_sdk.File, error) {
	out, err := os.Create(filePath)
	if err != nil {
		return files_sdk.File{}, err
	}
	params.Writer = out
	return c.Download(params)
}

func DownloadToFile(params files_sdk.FileDownloadParams, filePath string) (files_sdk.File, error) {
	return (&Client{}).DownloadToFile(params, filePath)
}

func (c *Client) Download(params files_sdk.FileDownloadParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := "/files/" + lib.QueryEscape(params.Path) + ""
	exportParams, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, _, err := files_sdk.Call("GET", c.Config, path, exportParams)
	if err != nil {
		return file, err
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	resp, err := c.Config.GetHttpClient().Get(file.DownloadUri)
	if err != nil {
		return file, err
	}
	_, err = io.Copy(params.Writer, resp.Body)
	if err != nil {
		return file, err
	}

	return file, nil
}

func Download(params files_sdk.FileDownloadParams) (files_sdk.File, error) {
	client := Client{}
	return client.Download(params)
}

func (c *Client) Create(params files_sdk.FileCreateParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := "/files/" + lib.QueryEscape(params.Path) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, res, err := files_sdk.Call("POST", c.Config, path, exportedParms)
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Create(params files_sdk.FileCreateParams) (files_sdk.File, error) {
	return (&Client{}).Create(params)
}

func (c *Client) Update(params files_sdk.FileUpdateParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := "/files/" + lib.QueryEscape(params.Path) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, res, err := files_sdk.Call("PATCH", c.Config, path, exportedParms)
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Update(params files_sdk.FileUpdateParams) (files_sdk.File, error) {
	return (&Client{}).Update(params)
}

func (c *Client) Delete(params files_sdk.FileDeleteParams) (files_sdk.File, error) {
	file := files_sdk.File{}
	path := "/files/" + lib.QueryEscape(params.Path) + ""
	exportedParms, err := lib.ExportParams(params)
	if err != nil {
		return file, err
	}
	data, res, err := files_sdk.Call("DELETE", c.Config, path, exportedParms)
	if err != nil {
		return file, err
	}
	if res.StatusCode == 204 {
		return file, nil
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Delete(params files_sdk.FileDeleteParams) (files_sdk.File, error) {
	return (&Client{}).Delete(params)
}
