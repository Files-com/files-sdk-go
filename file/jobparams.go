package file

import (
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"
	"github.com/hashicorp/go-retryablehttp"
)

func SetJobParams(r *status.Job, d direction.Direction, params interface{}, logger retryablehttp.Logger) {
	r.Params = params
	r.Direction = d
	r.Logger = logger
	switch d {
	case direction.DownloadType:
		p := params.(DownloaderParams)
		r.SetManager(p.Manager)
		r.SetEventsReporter(p.EventsReporter)
		r.RetryPolicy = p.RetryPolicy
		r.LocalPath = p.LocalPath
		r.RemotePath = p.RemotePath
		r.Sync = p.Sync
	case direction.UploadType:
		p := params.(UploaderParams)
		r.SetManager(p.Manager)
		r.SetEventsReporter(p.EventsReporter)
		r.RetryPolicy = p.RetryPolicy
		r.LocalPath = p.LocalPath
		r.RemotePath = p.RemotePath
		r.Sync = p.Sync
	}

}
