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
    "github.com/Files-com/files-sdk-go"
    "github.com/Files-com/files-sdk-go/folder"
)
```

Run any of the normal `go` commands (`build`/`install`/`test`). The Go
toolchain will resolve and fetch the files module automatically.

## Documentation

### Setting API Key

#### Setting by env 

``` sh
FILES_API_KEY="XXXX-XXXX..."
```

#### Set Global Variable

```go 
import (
    "github.com/Files-com/files-sdk-go"
)

 files_sdk.APIKey = "XXXX-XXXX..."
```

#### Set Per Client

```go 
import (
    "github.com/Files-com/files-sdk-go"
    "github.com/Files-com/files-sdk-go/file"
)

config :=  files_sdk.Config{APIKey: "XXXX-XXXX..."}
client := file.Client{Config: config}
```