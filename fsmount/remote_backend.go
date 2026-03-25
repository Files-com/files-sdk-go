package fsmount

import (
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	api_key "github.com/Files-com/files-sdk-go/v3/apikey"
	"github.com/Files-com/files-sdk-go/v3/file"
	file_migration "github.com/Files-com/files-sdk-go/v3/filemigration"
	"github.com/Files-com/files-sdk-go/v3/lock"
)

type remoteFileIter interface {
	Next() bool
	File() files_sdk.File
	Err() error
}

type remoteLockIter interface {
	Next() bool
	Lock() files_sdk.Lock
	Err() error
}

type remoteBackend interface {
	findCurrent(opts ...files_sdk.RequestResponseOption) (files_sdk.ApiKey, error)
	find(params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)
	listFor(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (remoteFileIter, error)
	createFolder(params files_sdk.FolderCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)
	move(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error)
	update(params files_sdk.FileUpdateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)
	uploadWithResume(opts ...file.UploadOption) (file.UploadResumable, error)
	upload(opts ...file.UploadOption) error
	downloadToFile(params files_sdk.FileDownloadParams, filePath string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)
	download(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error)
	createLock(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.Lock, error)
	deleteLock(params files_sdk.LockDeleteParams, opts ...files_sdk.RequestResponseOption) error
	listLocksFor(params files_sdk.LockListForParams, opts ...files_sdk.RequestResponseOption) (remoteLockIter, error)
	delete(params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) error
	wait(action files_sdk.FileAction, status func(files_sdk.FileMigration), opts ...files_sdk.RequestResponseOption) (files_sdk.FileMigration, error)
}

type sdkRemoteBackend struct {
	fileClient      *file.Client
	lockClient      *lock.Client
	apiKeyClient    *api_key.Client
	migrationClient *file_migration.Client
}

func (b *sdkRemoteBackend) findCurrent(opts ...files_sdk.RequestResponseOption) (files_sdk.ApiKey, error) {
	return b.apiKeyClient.FindCurrent(opts...)
}

func (b *sdkRemoteBackend) find(params files_sdk.FileFindParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return b.fileClient.Find(params, opts...)
}

func (b *sdkRemoteBackend) listFor(params files_sdk.FolderListForParams, opts ...files_sdk.RequestResponseOption) (remoteFileIter, error) {
	return b.fileClient.ListFor(params, opts...)
}

func (b *sdkRemoteBackend) createFolder(params files_sdk.FolderCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return b.fileClient.CreateFolder(params, opts...)
}

func (b *sdkRemoteBackend) move(params files_sdk.FileMoveParams, opts ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
	return b.fileClient.Move(params, opts...)
}

func (b *sdkRemoteBackend) update(params files_sdk.FileUpdateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return b.fileClient.Update(params, opts...)
}

func (b *sdkRemoteBackend) uploadWithResume(opts ...file.UploadOption) (file.UploadResumable, error) {
	return b.fileClient.UploadWithResume(opts...)
}

func (b *sdkRemoteBackend) upload(opts ...file.UploadOption) error {
	return b.fileClient.Upload(opts...)
}

func (b *sdkRemoteBackend) downloadToFile(params files_sdk.FileDownloadParams, filePath string, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return b.fileClient.DownloadToFile(params, filePath, opts...)
}

func (b *sdkRemoteBackend) download(params files_sdk.FileDownloadParams, opts ...files_sdk.RequestResponseOption) (files_sdk.File, error) {
	return b.fileClient.Download(params, opts...)
}

func (b *sdkRemoteBackend) createLock(params files_sdk.LockCreateParams, opts ...files_sdk.RequestResponseOption) (files_sdk.Lock, error) {
	return b.lockClient.Create(params, opts...)
}

func (b *sdkRemoteBackend) deleteLock(params files_sdk.LockDeleteParams, opts ...files_sdk.RequestResponseOption) error {
	return b.lockClient.Delete(params, opts...)
}

func (b *sdkRemoteBackend) listLocksFor(params files_sdk.LockListForParams, opts ...files_sdk.RequestResponseOption) (remoteLockIter, error) {
	return b.lockClient.ListFor(params, opts...)
}

func (b *sdkRemoteBackend) delete(params files_sdk.FileDeleteParams, opts ...files_sdk.RequestResponseOption) error {
	return b.fileClient.Delete(params, opts...)
}

func (b *sdkRemoteBackend) wait(action files_sdk.FileAction, status func(files_sdk.FileMigration), opts ...files_sdk.RequestResponseOption) (files_sdk.FileMigration, error) {
	return b.migrationClient.Wait(action, status, opts...)
}
