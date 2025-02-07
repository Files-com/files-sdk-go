package file

import (
	"encoding/json"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/stretchr/testify/assert"
)

func TestJobFile_JSONMarshal(t *testing.T) {
	assert := assert.New(t)
	file := JobFile{
		Status: status.Errored,
		Err:    MashableError{files_sdk.ResponseError{Title: "error"}},
	}
	data, err := json.Marshal(file)
	assert.Nil(err)
	assert.Contains(string(data), "error")
}
