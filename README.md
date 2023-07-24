# Files.com Go Client

The Files.com Go client library provides convenient access to the Files.com API from applications written in the Go language.

## Installation

Make sure your project is using Go Modules (it will have a `go.mod` file in its
root if it already is):

``` sh
go mod init
```

Then, reference files-sdk-go in a Go program with `import`:

``` go
import (
    "github.com/Files-com/files-sdk-go/v2"
    "github.com/Files-com/files-sdk-go/v2/folder"
)
```

Run any of the normal `go` commands (`build`/`install`/`test`). The Go
toolchain will resolve and fetch the files module automatically.

## Documentation

### Setting API Key

#### Setting by ENV 

``` sh
export FILES_API_KEY="XXXX-XXXX..."
```

#### Set Global Variable

```go 
import (
    "github.com/Files-com/files-sdk-go/v2"
)

 files_sdk.APIKey = "XXXX-XXXX..."
```

#### Set Per Client

```go 
import (
    "github.com/Files-com/files-sdk-go/v2"
    "github.com/Files-com/files-sdk-go/v2/file"
)

config := files_sdk.Config{APIKey: "XXXX-XXXX..."}
client := file.Client{Config: config}
```

### List

```go 
import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	folder "github.com/Files-com/files-sdk-go/v2/folder"
    "fmt"
)

func main() {
    it, err := folder.ListFor(files_sdk.FolderListForParams{})

    if err != nil {
        // deal with error
    }

    for it.Next() {
        entry := it.Folder()
        fmt.Println(entry.Path)
    }
}

```

### Upload a File
```go 
import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	file "github.com/Files-com/files-sdk-go/v2/file"
)

func main() {
    client := file.Client{}
    uploadPath := "file-to-upload.txt"
    destinationPath := "file-to-upload.txt"
    err := client.Upload(UploadWithFile(uploadPath), UploadWithDestinationPath(destinationPath))
    if err != nil {
        panic(err)
    }
}
```

#### Via io.Reader

```go 
import file "github.com/Files-com/files-sdk-go/v2/file"

func main() {
    client := file.Client{}
    io := strings.NewReader("my file contents")
    destinationPath := "my-file.txt"
    err := client.Upload(UploadWithReader(io), UploadWithDestinationPath(destinationPath))
    if err != nil {
        panic(err)
    }
}
```

### Download a File
```go 
import file "github.com/Files-com/files-sdk-go/v2/file"

func main() {
    client := file.Client{}
    downloadPath := "file-to-download.txt"
    fileEntry, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: "file-to-download.txt"}, downloadPath)
    if err != nil {
        panic(err)
    }
}
```

### Docker

```shell
docker build . --tag files-sdk-go:latest
docker run --workdir /app --volume ${PWD}:/app -it files-sdk-go
```
