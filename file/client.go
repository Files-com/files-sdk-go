package file

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/downloadurl"
	"github.com/Files-com/files-sdk-go/v2/folder"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Get(ctx context.Context, Path string) (files_sdk.File, error) {
	file := files_sdk.File{}
	path, err := lib.BuildPath("/files/{path}", map[string]string{"path": Path})
	if err != nil {
		return file, err
	}
	data, _, err := files_sdk.Call(ctx, "GET", c.Config, path, lib.Params{Params: lib.Interface()})
	if err != nil {
		return file, err
	}
	if err := file.UnmarshalJSON(*data); err != nil {
		return file, err
	}

	return file, nil
}

func Get(ctx context.Context, Path string) (files_sdk.File, error) {
	client := Client{}
	return client.Get(ctx, Path)
}

// File{}.Size and File{}.Mtime are not always up to date. This calls HEAD on File{}.DownloadUri to get the latest info.
// Some Download URLs won't support HEAD. In this case the size is reported as -1. The size can be known post download
// using Client{}.DownloadRequestStatus. This applies to the remote mount types FTP, SFTP, and WebDAV.
func (c *Client) FileStats(ctx context.Context, file files_sdk.File) (files_sdk.File, error) {
	var err error
	var size int64
	file, err = c.Download(
		ctx,
		files_sdk.FileDownloadParams{File: file},
		files_sdk.RequestOption(func(req *http.Request) error {
			if req.URL.Host != "s3.amazonaws.com" {
				req.Method = "HEAD"
			}
			return nil
		}),
		files_sdk.ResponseOption(func(response *http.Response) error {
			if response.StatusCode == 422 {
				size = -1 // Size is unknown
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
	)
	if err == nil {
		file.Size = size
	}
	return file, err
}

func (c *Client) DownloadRequestStatus(ctx context.Context, fileDownloadUrl string, downloadRequestId string, opts ...files_sdk.RequestResponseOption) (files_sdk.ResponseError, error) {
	re := files_sdk.ResponseError{}
	uri, err := url.Parse(fileDownloadUrl)
	if err != nil {
		return re, err
	}

	uri = uri.JoinPath(downloadRequestId)

	request, err := http.NewRequestWithContext(ctx, "GET", uri.String(), nil)
	if err != nil {
		return re, err
	}
	resp, err := files_sdk.WrapRequestOptions(c.Config.GetHttpClient(), request, opts...)
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
		re.Errors = append(re.Errors, files_sdk.ResponseError{Type: "download request status"})
	}
	return re, err
}

func (c *Client) DownloadUri(ctx context.Context, params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	var err error
	if params.Path == "" {
		params.Path = params.File.Path
	}

	if params.File.DownloadUri == "" {
		err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/files/{path}", Params: params, Entity: &params.File}, opts...)
		return params.File, err
	} else {
		url, parseErr := downloadurl.New(params.File.DownloadUri)
		remaining, valid := url.Valid(time.Millisecond * 100)
		if parseErr == nil {
			if !valid {
				err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/files/{path}", Params: params, Entity: &params.File})
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

func (c *Client) Download(ctx context.Context, params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	if params.Path == "" {
		params.Path = params.File.Path
	}
	var err error

	params.File, err = c.DownloadUri(ctx, params)
	if err != nil {
		return params.File, err
	}
	request, err := http.NewRequestWithContext(ctx, "GET", params.File.DownloadUri, nil)
	if err != nil {
		return params.File, err
	}

	_, err = files_sdk.WrapRequestOptions(c.Config.GetHttpClient(), request, opts...)

	return params.File, err
}

func Download(ctx context.Context, params files_sdk.FileDownloadParams) (files_sdk.File, error) {
	client := Client{}
	return client.Download(ctx, params)
}

func (c *Client) Create(ctx context.Context, params files_sdk.FileCreateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/files/{path}", Params: params, Entity: &file}, opts...)
	return
}

func Create(ctx context.Context, params files_sdk.FileCreateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	return (&Client{}).Create(ctx, params, opts...)
}

func (c *Client) Update(ctx context.Context, params files_sdk.FileUpdateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/files/{path}", Params: params, Entity: &file}, opts...)
	return
}

func Update(ctx context.Context, params files_sdk.FileUpdateParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	return (&Client{}).Update(ctx, params, opts...)
}

func (c *Client) UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "PATCH", Path: "/files/{path}", Params: params, Entity: &file}, opts...)
	return
}

func UpdateWithMap(ctx context.Context, params map[string]interface{}, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	return (&Client{}).UpdateWithMap(ctx, params, opts...)
}

func (c *Client) Delete(ctx context.Context, params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "DELETE", Path: "/files/{path}", Params: params, Entity: nil}, opts...)
	return
}

func Delete(ctx context.Context, params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) (err error) {
	return (&Client{}).Delete(ctx, params, opts...)
}

func (c *Client) Find(ctx context.Context, params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "GET", Path: "/file_actions/metadata/{path}", Params: params, Entity: &file}, opts...)
	return
}

func Find(ctx context.Context, params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (file files_sdk.File, err error) {
	return (&Client{}).Find(ctx, params, opts...)
}

func (c *Client) Copy(ctx context.Context, params files_sdk.FileCopyParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/file_actions/copy/{path}", Params: params, Entity: &fileAction}, opts...)
	return
}

func Copy(ctx context.Context, params files_sdk.FileCopyParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	return (&Client{}).Copy(ctx, params, opts...)
}

func (c *Client) Move(ctx context.Context, params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/file_actions/move/{path}", Params: params, Entity: &fileAction}, opts...)
	return
}

func Move(ctx context.Context, params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (fileAction files_sdk.FileAction, err error) {
	return (&Client{}).Move(ctx, params, opts...)
}

func (c *Client) BeginUpload(ctx context.Context, params files_sdk.FileBeginUploadParams, opts ...files_sdk.RequestResponseOption) (fileUploadPartCollection files_sdk.FileUploadPartCollection, err error) {
	err = files_sdk.Resource(ctx, c.Config, lib.Resource{Method: "POST", Path: "/file_actions/begin_upload/{path}", Params: params, Entity: &fileUploadPartCollection}, opts...)
	return
}

func BeginUpload(ctx context.Context, params files_sdk.FileBeginUploadParams, opts ...files_sdk.RequestResponseOption) (fileUploadPartCollection files_sdk.FileUploadPartCollection, err error) {
	return (&Client{}).BeginUpload(ctx, params, opts...)
}

func (c *Client) ListFor(ctx context.Context, params files_sdk.FolderListForParams) (*folder.Iter, error) {
	client := folder.Client{Config: c.Config}
	return client.ListFor(ctx, params)
}

func (c *Client) ListForRecursive(ctx context.Context, params files_sdk.FolderListForParams) (lib.IterI, error) {
	it := lib.IterChan{}.Init()

	go func(params files_sdk.FolderListForParams) {
		f := (&FS{}).Init(c.Config, true).WithContext(ctx).(*FS)
		err := fs.WalkDir(f, params.Path, func(path string, d fs.DirEntry, err error) error {
			if path == "" && err == nil {
				return nil // Skip root directory
			}

			if err == nil {
				info, _ := d.Info()
				it.Send <- info.Sys()
			} else {
				it.SendError <- err
			}
			return err
		})
		if err != nil {
			it.Error.Store(err)
		}
		it.Stop <- true
	}(params)
	return it, nil
}
