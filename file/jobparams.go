package file

import (
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"
)

func SetJobParams(r *status.Job, d direction.Type, params interface{}) {
	r.Params = params
	r.Direction = d
	switch d {
	case direction.DownloadType:
		p := params.(DownloadFolderParams)
		r.SetManager(p.Manager)
		r.SetEventsReporter(p.EventsReporter)
		r.RetryPolicy = string(p.RetryPolicy)
		r.LocalPath = p.LocalPath
		r.RemotePath = p.RemotePath
		r.Sync = p.Sync
	case direction.UploadType:
		p := params.(UploadParams)
		r.SetManager(p.Manager)
		r.SetEventsReporter(p.EventsReporter)
		r.RetryPolicy = string(p.RetryPolicy)
		r.LocalPath = p.LocalPath
		r.RemotePath = p.RemotePath
		r.Sync = p.Sync
	}

}
