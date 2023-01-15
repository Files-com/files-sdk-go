package file

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Files-com/files-sdk-go/v2/lib"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/directory"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/Files-com/files-sdk-go/v2/lib/direction"
)

// DeleteSource files after a sync
//
//	job.RegisterFileEvent(func(file status.File) {
//			log, err := file.DeleteSource{Direction: f.Direction, Config: config}.Call(ctx, f)
//	}, status.Complete, status.Skipped)
type DeleteSource struct {
	direction.Direction
	Config files_sdk.Config
}

func (ad DeleteSource) Call(ctx context.Context, f status.File) (status.Log, error) {
	switch f.Direction {
	case direction.UploadType:
		return status.Log{Path: f.LocalPath, Action: "delete source"}, os.Remove(f.LocalPath)
	case direction.DownloadType:
		client := Client{Config: ad.Config}
		err := client.Delete(ctx, files_sdk.FileDeleteParams{Path: f.RemotePath})
		return status.Log{Path: f.RemotePath, Action: "delete source"}, err
	default:
		panic(fmt.Sprintf("unknown direction %v", f.Direction))
	}
}

// MoveSource files after a sync
//
//	job.RegisterFileEvent(func(file status.File) {
//			log, err := file.MoveSource{Direction: f.Direction, Config: config}.Call(ctx, f)
//	}, status.Complete, status.Skipped)
type MoveSource struct {
	direction.Direction
	Path   string
	Config files_sdk.Config
}

func (am MoveSource) Call(ctx context.Context, f status.File) (status.Log, error) {
	var err error
	log := status.Log{Action: "move source"}
	log.Path = am.movePath(f)

	switch f.Direction {
	case direction.UploadType:
		dir, _ := filepath.Split(log.Path)
		err = os.MkdirAll(dir, 0755)
		if err != nil && !errors.Is(err, syscall.EEXIST) {
			return log, err
		}
		err = os.Rename(f.LocalPath, log.Path)
		if err != nil && errors.Is(err, syscall.EEXIST) {
			err = os.Remove(log.Path)
			if err != nil {
				return log, err
			}
			return am.Call(ctx, f)
		}
		return log, err
	case direction.DownloadType:
		client := &Client{Config: am.Config}
		_, err := client.Move(
			ctx,
			files_sdk.FileMoveParams{Path: f.RemotePath, Destination: lib.Path{Path: log.Path}.NormalizePathSystemForAPI().String()},
		)
		rErr, ok := err.(files_sdk.ResponseError)
		if ok && rErr.Type == "processing-failure/destination-parent-does-not-exist" {
			err := (&FS{}).Init(am.Config, true).WithContext(ctx).(*FS).MkdirAll(filepath.Dir(log.Path), 0755)
			if err != nil {
				return log, err
			}
			return am.Call(ctx, f)
		}
		if ok && rErr.Type == "processing-failure/destination-exists" {
			err := client.Delete(ctx, files_sdk.FileDeleteParams{Path: lib.Path{Path: log.Path}.NormalizePathSystemForAPI().String()})
			if err != nil {
				return log, err
			}
			return am.Call(ctx, f)
		}
		return log, err
	default:
		panic(fmt.Sprintf("unknown direction %v", f.Direction))
	}
}

func (am MoveSource) movePath(f status.File) string {
	switch f.Job.Type {
	case directory.Dir:
		return filepath.Join(
			append([]string{am.Path}, strings.Split(strings.TrimPrefix(f.RemotePath, f.Job.RemotePath), "/")...)...,
		)
	case directory.File:
		return am.Path
	default:
		panic("")
	}
}
