# Change Log

All notable changes to this project will be documented in this file.

## [1.0.175] - 2021/04/12
### Fix
- Race condition: `file.Client#UploadFolder` Uploading nested folders could sometimes skip folders.
- `file.Client#DownloadFolder(files_sdk.FolderListForParams{Path: "documents/report.pdf"}, "local-files")` would result in `local-files/documents/report.pdf`. This is now fixed resulting in `local-files/report.pdf`
- `file.Client#DownloadFolder(files_sdk.FolderListForParams{Path: "documents/report.pdf"}, "report-2020.pdf")` would result in `local-files/documents/report-2020.pdf/report.pdf`. This is now fixed resulting in `local-files/report-2020.pdf`
- Removed `lib.Iter{}.MaxPages` default of 1 allowing for downloaded of folders that contain more than 1000 files/folders.

## [1.0.156] - 2021/02/22
### Fix
- In some cases API errors were resulting in a json.UnmarshalTypeError.
