# Change Log

All notable changes to this project will be documented in this file.
This project gets auto released on every change to the [Files.com API](https://developers.files.com).
Auto generated releases contain additions and fixes to models and method arguments, theses will not be documented here.

## [2.0.81-beta] - 2022/06/14
### Change
- `file.Client{}.delete` now returns only `error` instead of `(files_sdk.File, error)`

## [2.0.80-beta] - 2022/06/13
### Fix
- `file.Downloader` and `file.Uploader`
  - Could cause a nil pointer panic when a `file.RetryPolicy` is set.
  - Removing possible blocking while job scanning phase that could block transfers from starting.

## [2.0.38-beta] - 2021/12/16
### Add
- `file_migration.Client.Wait` poll async `file.FileAction` for results. 
- `file_migration.Client.LogIterator` given a completed file migration this returns a detailed log of each file event.

## [2.0.36-beta] - 2021/12/03
### Add
- Alias `folder.Client.ListFor` `file.Client.ListFor`
- `file.Client.ListForRecursive` recursively list folders/files.

## [2.0.30-beta] - 2021/10/27
### Fix
- Models that contained an array if ids is now correctly deserialized from JSON. Fixes error `json: cannot unmarshal number into Go struct field user.admin_group_ids of type string`

## [2.0.28-beta] - 2021/10/25
### Fix
- `file.UploadIO` could cause `panic: runtime error: invalid memory address or nil pointer dereference` when `file.UploadIOParams.FilesUploadParts.ParallelParts` is nil

## [2.0.27-beta] - 2021/10/22
### Improvements
- `file.UploadIO` 
  - Reduce unnecessary HTTP calls when uploading to a remote server.
  - Better invalidation of upload parts when a remote server upload is no longer found.
  - Reduction byte allocations when segmenting file parts.
  - `file.ProxyReader` can be reread after closed to allowing for inline reties of http bodies. 
  - Skip indicating to close connection after part upload for remote mount urls that don't change.
- `status.Job` Less jittery `TransferRate` and `ETA` indicators when uploading to slower remote mounts.

## [2.0.26-beta] - 2021/10/19
### Fix
- `files_sdk.ResponseError.UnmarshalJSON` detects HTML errors and returns an error with it rather than a JSON syntax error.
- `file.Uploader` When uploading to a remote server parts were being incorrectly validating causing `status.File.TransferBytes` to return 0.

## [2.0.25-beta] - 2021/10/19
### Add
- `file.Uploader`/`file.UploadIO` sends known file size in preflight upload request. Allows some remote servers to upload larger files.

### Change
- `file.Downloader`/`file.Uploader` `Sync` param compares file size instead of modified time to match the server sync.

### Fix
- Fix race condition `file.Uploader` where the reported upload bytes were less than the actual.
- All detected race warning have been fixed.

## [2.0.20-beta] - 2021/10/01
### Fix
- `file.Downloader` with sync panics when local file doesn't exist. Reason local error variable was overwritten by another call causing it to be nil when it should have had IsNotExist error object.

## [2.0.19-beta] - 2021/09/30
### Fix
- `file.Uploader` when uploading a single file the `Job.Scanning`/`Job.EndScanning` signals were not called. This caused problems for the CLI expecting that they would always be called.

## [2.0.18-beta] - 2021/09/30
### Change
- `file.RetryByPolicy` takes 4th arg `signalEvents bool` when `true` it resets all events and calls new ones. When `false` fixes the case where events that were only expect to happen once don't get repeated.

## [2.0.17-beta] - 2021/09/30
### Fix
- panic caused by Job.Finish() being called twice after a retry is needed.

## [2.0.16-beta] - 2021/09/30
### Add
- An upload job can be canceled and files will restart where is left off using `file.RetryByPolicy`.

## [2.0.12-beta] - 2021/09/25
### Fix
- Downloading files Error `too many open files`

## [2.0.11-beta] - 2021/09/24
### Fix
- Uploading files Error `too many open files`: solution was to inform the server that the client wants to close any connection after a transaction is complete.

## [2.0.10-beta] - 2021/09/23
### Change
- `status.Job` has subscribable events `Started`,`Finished`,`Canceled`,`Scanning`, and `EndScanning`. Call `Subscribe()` to return `chan time.Time`

### Fix
- `file.Client.Downloader` won't query stats of local/remote file when sync is false.

## [2.0.7-beta] - 2021/09/21
### Changes
- `file.RetryTransfers` changed to `file.RetryByPolicy`.

### Fix
- Retrying Canceled jobs via `file.RetryByPolicy` when files were in progress were not rerun correctly.

## [2.0.4-beta] - 2021/09/17
### Fix
- When using `lib.Iter{}.Next()` if a cursor is not given, in the case of listing a file, it would always return true.

## [2.0.2-beta] - 2021/09/17
### Fix
- `file.Client{}.Upload()` supports up to 4.9 TB files due to improved file chunking.

## [2.0.0-beta] - 2021/09/13
### Changes
- API changes to `file.UploadFolderOrFile` => `file.Uploader`, `file.DownloadFolder` => `file.Downloader`

### Add
- `file.FS{}.Init(context.Background(), files_sdk.Config{})` support for using Files.com as FS implementation.
- `file.Uploader` Concurrent files system scanning improves performance for uploading large numbers of files.
- `status.Job` includes stats `TransferRate`, `ETA`, `ElapsedTime`, and `Percentage`.
- `file.RetryPolicy` for  `file.Uploader`/`file.Downloader`

### Fix
- `lib.Iter{}.Next()` if `PerPage` was not set func didn't not return all results.
- Calling `status.Job.Any()` could cause deadlock.

## [1.2.1146] - 2021/08/03
### Fix
- `file.DownloadFolder` handle concurrent downloads on the same path by incrementing the tmp file name.

### Add
- `file.DownloadRetry(status.File)` and `file.UploadRetry(status.File)` give a interface to retry files that have may of failed.

### Changes
- `status.Report` is removed and replaced with `status.File`
- `file.DownloadFolder` and `file.UploadFolder` no longer return `(*status.Job, error)` instead only `(*status.Job)`. All errors are sent via the Reporter func. 
- When the SDK calls the given `Reporter: func(status.Report, error)` it will block and the user provided function should handle making code inside async. 

## [1.1.1145] - 2021/08/03
### Changes
- Every applicable function now take `Context` as the first parameter. This allows for cancellation of tasks in flight.
- `files_sdk.Config{}` has removed `SetMaxConcurrentConnections`. Instead pass `manager.Manager{}` to `file.DownloadFolder` and `file.UploadFolder`.
- `file.DownloadFolder` and `file.UploadFolder` now take `Reporter: func(status.Report, error)`

## [1.1.1144] - 2021/08/03
### Fix
- `files.Client{}.DownloadFolder()` Fix Windows issue `The process cannot access the file because it is being used by another process.`
- `files.Client{}.DownloadFolder()` in some cases the func hangs after all files are download.
- Enum constants are removed due to issue with duplicates. Use `Enum()["value"]` to validate input.
  
## [1.1.1143] - 2021/08/03
### Add
- Enum constants are available for structs params used as server requests. 

## [1.1.1142] - 2021/08/03### Changes
### Feature
- `file.UploadFile(file.UploadParams{})` and `file.UploadFolder(file.UploadParams{})` now uploads file chunks in parallel. 
  Defaults to 25, but can be changed via `files_sdk.Config{}.SetMaxConcurrentConnections(50)`

### Fix
- Reduce memory usage when not in debug mode.

## [1.0.183] - 2021/04/28### Fix
### Fix
- Race condition: `file.Client#UploadFolder` Uploading nested folders could sometimes skip folders.
- `file.Client#DownloadFolder(files_sdk.FolderListForParams{Path: "documents/report.pdf"}, "local-files")` would result in `local-files/documents/report.pdf`. This is now fixed resulting in `local-files/report.pdf`
- `file.Client#DownloadFolder(files_sdk.FolderListForParams{Path: "documents/report.pdf"}, "report-2020.pdf")` would result in `local-files/documents/report-2020.pdf/report.pdf`. This is now fixed resulting in `local-files/report-2020.pdf`
- Removed `lib.Iter{}.MaxPages` default of 1 allowing for downloaded of folders that contain more than 1000 files/folders.

## [1.0.156] - 2021/02/22### Add
### Fix
- In some cases API errors were resulting in a json.UnmarshalTypeError.
