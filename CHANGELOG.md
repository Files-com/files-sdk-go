# Change Log

All notable changes to this project will be documented in this file.
This project gets auto released on every change to the [Files.com API](https://developers.files.com).
Auto generated releases contain additions and fixes to models and method arguments, theses will not be documented here.

## [1.0.184] - 2021/04/28### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature### Feature
- `file.UploadFile(file.UploadParams{})` and `file.UploadFolder(file.UploadParams{})` now uploads file chunks in parallel. 
  Defaults to 25, but can be changed via `files_sdk.Config{}.SetMaxConcurrentConnections(50)`

### Fix
- Reduce memory usage when not in debug mode.

## [1.0.183] - 2021/04/28### Feature
### Fix
- Race condition: `file.Client#UploadFolder` Uploading nested folders could sometimes skip folders.
- `file.Client#DownloadFolder(files_sdk.FolderListForParams{Path: "documents/report.pdf"}, "local-files")` would result in `local-files/documents/report.pdf`. This is now fixed resulting in `local-files/report.pdf`
- `file.Client#DownloadFolder(files_sdk.FolderListForParams{Path: "documents/report.pdf"}, "report-2020.pdf")` would result in `local-files/documents/report-2020.pdf/report.pdf`. This is now fixed resulting in `local-files/report-2020.pdf`
- Removed `lib.Iter{}.MaxPages` default of 1 allowing for downloaded of folders that contain more than 1000 files/folders.

## [1.0.156] - 2021/02/22
### Fix
- In some cases API errors were resulting in a json.UnmarshalTypeError.
