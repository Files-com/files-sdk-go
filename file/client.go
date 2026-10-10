// Package file provides the Files.com File API client.
package file

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Files-com/files-sdk-go/v3/downloadurl"
	"github.com/Files-com/files-sdk-go/v3/folder"

	"errors"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Client calls the File API using its embedded Config.
type Client struct {
	files_sdk.Config
}

func (c *Client) Get(Path string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	file := files_sdk.File{}
	path, err := lib.BuildPath("/files/{path}", map[string]string{"path": Path})
	if err != nil {
		return file, err
	}
	data, _, err := files_sdk.Call("GET", c.Config, path, lib.Params{}, opts...)
	if err != nil {
		return file, err
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Get(Path string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	client := Client{}
	return client.Get(Path, opts...)
}

// FileStats checks for current file size and modification time.
// File{}.Size and File{}.Mtime are not always up to date. This calls HEAD on File{}.DownloadUri to get the latest info.
// Some Download URLs won't support HEAD. In this case the size is reported as UntrustedSizeValue. The size can be known post download
// using Client{}.DownloadRequestStatus. This applies to the remote mount types FTP, SFTP, and WebDAV.
func (c *Client) FileStats(file files_sdk.File, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	var err error
	var size int64
	file, err = c.Download(
		files_sdk.FileDownloadParams{File: file},
		append(opts,
			files_sdk.RequestOption(func(req *http.Request) error {
				if req.URL.Host != "s3.amazonaws.com" {
					req.Method = "HEAD"
				}
				return nil
			}),
			files_sdk.ResponseOption(func(response *http.Response) error {
				if response.StatusCode == 422 {
					size = int64(UntrustedSizeValue)
					return nil
				}
				if err := lib.ResponseErrors(response, lib.NonOkError); err != nil {
					return err
				}
				size = response.ContentLength
				if response.Header.Get("Last-Modified") != "" {
					mtime, err := time.Parse(time.RFC1123, response.Header.Get("Last-Modified"))
					if err == nil {
						file.Mtime = &mtime
					}
				}
				return response.Body.Close()
			}),
		)...,
	)
	if err == nil {
		file.Size = size
	}
	return file, err
}

func (c *Client) DownloadRequestStatus(fileDownloadUrl string, downloadRequestId string, opts ...files_sdk.RequestResponseOption) (files_sdk.ResponseError, error) {
	re := files_sdk.ResponseError{}
	uri, err := url.Parse(fileDownloadUrl)
	if err != nil {
		return re, err
	}

	uri = uri.JoinPath(downloadRequestId)

	request, err := http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return re, err
	}

	c.SetHeadersForRequest(request)

	resp, err := files_sdk.WrapRequestOptions(c.Config, request, opts...)
	if err != nil {
		return re, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return re, err
	}
	if lib.IsJSON(resp) {
		err = re.UnmarshalJSON(data)
		if err != nil {
			return re, err
		}
		if re.Type == "" && !re.IsNil() {
			re.Type = "download request status"
		}
	}
	return re, err
}

func (c *Client) DownloadUri(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	var err error
	if params.Path == "" {
		params.Path = params.File.Path
	}
	if !c.Config.DisableDirectTransfers && params.WithDirectConnectionInfo == nil {
		params.WithDirectConnectionInfo = lib.Bool(true)
	}

	if params.File.DownloadUri == "" {
		err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/files/{path}", Params: params, Entity: &params.File}, opts...)
		return params.File, err
	} else {
		url, parseErr := downloadurl.New(params.File.DownloadUri)
		// url is nil when DownloadUri does not parse; Download then returns that parse error from building its request.
		if parseErr == nil {
			if remaining, valid := url.Valid(time.Millisecond * 250); !valid {
				err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/files/{path}", Params: params, Entity: &params.File}, opts...)
				if params.File.DownloadUri == url.URL.String() {
					c.LogPath(params.Path, map[string]interface{}{"message": "URL was expired. Fetched a new URL but it didn't change", "Remaining": remaining, "Time": url.Time})
				} else {
					c.LogPath(params.Path, map[string]interface{}{"message": "URL was expired. Fetched a new URL", "Remaining": remaining, "Time": url.Time})
				}
			}
		}
	}

	return params.File, err
}

func (c *Client) Download(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	if params.Path == "" {
		params.Path = params.File.Path
	}
	var err error

	params.File, err = c.DownloadUri(params, files_sdk.WithContext(files_sdk.ContextOption(opts)))
	if err != nil {
		return params.File, err
	}
	request, err := http.NewRequest("GET", params.File.DownloadUri, nil)
	if err != nil {
		return params.File, err
	}

	c.SetHeadersForRequest(request)

	directSuppressor := directTransferDownloadSuppressorFromContext(files_sdk.ContextOption(opts))
	directAttemptAllowed := !c.Config.DisableDirectTransfers && files_sdk.DirectConnectionInfoPresent(params.File.DirectConnectionInfo)
	if directAttemptAllowed && directSuppressor != nil {
		directAttemptAllowed = directSuppressor.directTransferDownloadAttemptAllowed()
	}
	if params.File.DownloadUri == "" && !directAttemptAllowed {
		return params.File, errors.New("download response did not include a usable download URL")
	}
	if directAttemptAllowed {
		directContext, closeDirectClient := files_sdk.WithDirectTransferClientCache(files_sdk.ContextOption(opts))
		defer closeDirectClient()
		directOptions := append(directTransferDownloadFailureOptions(directSuppressor, opts...), files_sdk.WithContext(directContext))
		for attempt := 1; attempt <= downloadV2RetryAttempts; attempt++ {
			_, err = files_sdk.WrapDirectTransferOptions(
				c.Config,
				params.File.DirectConnectionInfo,
				request,
				directOptions...,
			)
			var responseErr *files_sdk.DirectTransferResponseError
			if !errors.As(err, &responseErr) || responseErr.StatusCode != http.StatusTooManyRequests || attempt == downloadV2RetryAttempts {
				break
			}
			if backpressure := directTransferBackpressureFromContext(directContext); backpressure != nil {
				backpressure.record(responseErr.RetryAfter)
			}
			if responseErr.RetryAfter > 0 {
				select {
				case <-directContext.Done():
					err = directContext.Err()
				case <-time.After(responseErr.RetryAfter):
				}
				if err != nil {
					break
				}
			}
		}
		var responseErr *files_sdk.DirectTransferResponseError
		if errors.As(err, &responseErr) && responseErr.StatusCode == http.StatusTooManyRequests {
			if backpressure := directTransferBackpressureFromContext(directContext); backpressure != nil {
				backpressure.record(responseErr.RetryAfter)
			}
		}
		if err == nil {
			if c.Config.InDebug() {
				c.LogPath(params.Path, map[string]interface{}{"message": "direct download success", "direction": "download"})
			}
			return params.File, nil
		}
		if errors.Is(err, files_sdk.ErrDirectTransferResponseStarted) {
			return params.File, err
		}
		if directSuppressor != nil {
			directSuppressor.disableDirectTransferDownload("direct_request_failed", err)
		}
		if params.File.DownloadUri == "" {
			return params.File, err
		}
		c.LogPath(params.Path, map[string]interface{}{"message": "direct download failed; falling back to proxy URL", "direction": "download", "reason": "direct_request_failed", "error": uploadRetryLogError(err)})
	}

	_, err = files_sdk.WrapRequestOptions(c.Config, request, opts...)

	return params.File, err
}

func Download(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	client := Client{}
	return client.Download(params, opts...)
}

// Create calls POST /files/{path}.
//
// API operation: Upload File.
func (c *Client) Create(params files_sdk.FileCreateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/files/{path}", Params: params, Entity: &file}, opts...)
	return
}

// Create calls Client.Create using the default configuration.
func Create(params files_sdk.FileCreateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	return (&Client{}).Create(params, opts...)
}

// Update calls PATCH /files/{path}.
//
// API operation: Update File/Folder Metadata.
func (c *Client) Update(params files_sdk.FileUpdateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/files/{path}", Params: params, Entity: &file}, opts...)
	return
}

// Update calls Client.Update using the default configuration.
func Update(params files_sdk.FileUpdateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	return (&Client{}).Update(params, opts...)
}

// UpdateWithMap calls PATCH /files/{path} using API parameter names as map keys.
// Include any path parameters in the map. Unlike optional struct fields, explicit
// zero values in the map are included in the request.
func (c *Client) UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "PATCH", Path: "/files/{path}", Params: params, Entity: &file}, opts...)
	return
}

// UpdateWithMap calls Client.UpdateWithMap using the default configuration.
func UpdateWithMap(params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	return (&Client{}).UpdateWithMap(params, opts...)
}

// Delete calls DELETE /files/{path}.
//
// API operation: Delete File/Folder.
func (c *Client) Delete(params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "DELETE", Path: "/files/{path}", Params: params, Entity: nil}, opts...)
	return
}

// Delete calls Client.Delete using the default configuration.
func Delete(params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(params, opts...)
}

// Find calls GET /file_actions/metadata/{path}.
//
// API operation: Find File/Folder by Path.
func (c *Client) Find(params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/file_actions/metadata/{path}", Params: params, Entity: &file}, opts...)
	return
}

// Find calls Client.Find using the default configuration.
func Find(params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	return (&Client{}).Find(params, opts...)
}

// ZipListContents calls GET /file_actions/zip_list/{path}.
//
// API operation: List the contents of a ZIP file.
func (c *Client) ZipListContents(params files_sdk.FileZipListContentsParams, opts ...files_sdk.RequestResponseOption) (zipListEntryCollection files_sdk.ZipListEntryCollection, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "GET", Path: "/file_actions/zip_list/{path}", Params: params, Entity: &zipListEntryCollection}, opts...)
	return
}

// ZipListContents calls Client.ZipListContents using the default configuration.
func ZipListContents(params files_sdk.FileZipListContentsParams, opts ...files_sdk.RequestResponseOption) (zipListEntryCollection files_sdk.ZipListEntryCollection, err error) {
	return (&Client{}).ZipListContents(params, opts...)
}

// Copy calls POST /file_actions/copy/{path}.
//
// API operation: Copy File/Folder.
func (c *Client) Copy(params files_sdk.FileCopyParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/file_actions/copy/{path}", Params: params, Entity: &fileAction}, opts...)
	return
}

// Copy calls Client.Copy using the default configuration.
func Copy(params files_sdk.FileCopyParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	return (&Client{}).Copy(params, opts...)
}

// Move calls POST /file_actions/move/{path}.
//
// API operation: Move File/Folder.
func (c *Client) Move(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/file_actions/move/{path}", Params: params, Entity: &fileAction}, opts...)
	return
}

// Move calls Client.Move using the default configuration.
func Move(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	return (&Client{}).Move(params, opts...)
}

// Transform calls POST /file_actions/transform/{path}.
//
// API operation: Transform a file and save the output to a destination path.
func (c *Client) Transform(params files_sdk.FileTransformParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/file_actions/transform/{path}", Params: params, Entity: &fileAction}, opts...)
	return
}

// Transform calls Client.Transform using the default configuration.
func Transform(params files_sdk.FileTransformParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	return (&Client{}).Transform(params, opts...)
}

// GpgDecrypt calls POST /file_actions/gpg_decrypt/{path}.
//
// API operation: Decrypt a GPG-encrypted file and save it to a destination path.
func (c *Client) GpgDecrypt(params files_sdk.FileGpgDecryptParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/file_actions/gpg_decrypt/{path}", Params: params, Entity: &fileAction}, opts...)
	return
}

// GpgDecrypt calls Client.GpgDecrypt using the default configuration.
func GpgDecrypt(params files_sdk.FileGpgDecryptParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	return (&Client{}).GpgDecrypt(params, opts...)
}

// GpgEncrypt calls POST /file_actions/gpg_encrypt/{path}.
//
// API operation: Encrypt a file with GPG and save it to a destination path.
func (c *Client) GpgEncrypt(params files_sdk.FileGpgEncryptParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/file_actions/gpg_encrypt/{path}", Params: params, Entity: &fileAction}, opts...)
	return
}

// GpgEncrypt calls Client.GpgEncrypt using the default configuration.
func GpgEncrypt(params files_sdk.FileGpgEncryptParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	return (&Client{}).GpgEncrypt(params, opts...)
}

// Unzip calls POST /file_actions/unzip.
//
// API operation: Extract a ZIP file to a destination folder.
func (c *Client) Unzip(params files_sdk.FileUnzipParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/file_actions/unzip", Params: params, Entity: &fileAction}, opts...)
	return
}

// Unzip calls Client.Unzip using the default configuration.
func Unzip(params files_sdk.FileUnzipParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	return (&Client{}).Unzip(params, opts...)
}

// Zip calls POST /file_actions/zip.
//
// API operation: Create a ZIP from one or more paths and save it to a destination path.
func (c *Client) Zip(params files_sdk.FileZipParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/file_actions/zip", Params: params, Entity: &fileAction}, opts...)
	return
}

// Zip calls Client.Zip using the default configuration.
func Zip(params files_sdk.FileZipParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	return (&Client{}).Zip(params, opts...)
}

// BeginUpload calls POST /file_actions/begin_upload/{path}.
//
// API operation: Begin File Upload.
func (c *Client) BeginUpload(params files_sdk.FileBeginUploadParams, opts ...files_sdk.RequestResponseOption) (fileUploadPartCollection files_sdk.FileUploadPartCollection, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/file_actions/begin_upload/{path}", Params: params, Entity: &fileUploadPartCollection}, opts...)
	return
}

// BeginUpload calls Client.BeginUpload using the default configuration.
func BeginUpload(params files_sdk.FileBeginUploadParams, opts ...files_sdk.RequestResponseOption) (fileUploadPartCollection files_sdk.FileUploadPartCollection, err error) {
	return (&Client{}).BeginUpload(params, opts...)
}

type Iter struct {
	*folder.Iter
}

var _ files_sdk.ResourceIterator = Iter{}
var _ files_sdk.ResourceLoader = Iter{}

func (i Iter) LoadResource(identifier interface{}, opts ...files_sdk.RequestResponseOption) (interface{}, error) {
	params := files_sdk.FileFindParams{}
	if path, ok := identifier.(string); ok {
		params.Path = path
	}

	return (&Client{Config: i.Config}).Find(params, opts...)
}

func (i Iter) Iterate(identifier interface{}, opts ...files_sdk.RequestResponseOption) (files_sdk.IterI, error) {
	it, err := i.Iter.Iterate(identifier, opts...)
	return Iter{Iter: it.(*folder.Iter)}, err
}

func (c *Client) ListFor(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (Iter, error) {
	it, err := (&folder.Client{Config: c.Config}).ListFor(params, opts...)
	return Iter{Iter: it}, err
}

func (c *Client) CreateFolder(params files_sdk.FolderCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return (&folder.Client{Config: c.Config}).Create(params, opts...)
}

type RecursiveItem struct {
	files_sdk.File
	error `json:"error"`
}

func (r RecursiveItem) Err() error {
	return r.error
}

type recursiveIter struct {
	Iter
	root    *files_sdk.File
	current RecursiveItem
}

func (i *recursiveIter) Next() bool {
	if i.root != nil {
		i.current = RecursiveItem{File: *i.root}
		i.root = nil
		return true
	}
	if !i.Iter.Next() {
		return false
	}
	i.current = RecursiveItem{File: i.Iter.Current().(files_sdk.File)}
	return true
}

func (i *recursiveIter) Current() interface{} {
	return i.current
}

func (i *recursiveIter) Resource() RecursiveItem {
	return i.current
}

func (c *Client) ListForRecursive(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (files_sdk.TypedIterI[RecursiveItem], error) {
	if params.SortBy != nil || params.Search != "" || params.SearchAll != nil || params.SearchCustomMetadataKey != "" || params.ModifiedAtDatetime != nil {
		return nil, errors.New("recursive listings do not support sort_by, search, search_all, search_custom_metadata_key, or modified_at_datetime")
	}

	root := files_sdk.File{Type: "directory"}
	if params.Path == "." || lib.NormalizeForComparison(params.Path) == "" {
		params.Path = ""
	} else {
		var err error
		root, err = c.Find(files_sdk.FileFindParams{Path: params.Path}, opts...)
		if err != nil {
			return nil, err
		}
		if root.Type == "directory" && root.Path != params.Path {
			params.Path = root.Path
		}
	}

	opts = append(opts, files_sdk.RequestOption(func(request *http.Request) error {
		query := request.URL.Query()
		query.Set("recursive", "true")
		request.URL.RawQuery = query.Encode()
		return nil
	}))
	it, err := c.ListFor(params, opts...)
	if err != nil {
		return nil, err
	}
	recursive := &recursiveIter{Iter: it}
	if params.Cursor == "" && params.Type != "file" && root.Type == "directory" && root.Path != "" {
		recursive.root = &root
	}
	return recursive, nil
}
