package session

import (
	"os"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/lib/test"
	recorder "github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func CreateClient(fixture string) (client *Client, r *recorder.Recorder, err error) {
	client = &Client{}
	client.Config, r, err = test.CreateConfig(fixture)

	return client, r, err
}

func TestClient_Delete(t *testing.T) {
	client, r, err := CreateClient("Delete")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	assert := assert.New(t)
	os.Unsetenv("FILES_API_KEY")
	client.Config.SessionId = "9f799aff7f518514a0b6b5cfd1047e73dddd5cf5"
	err = client.Delete()
	assert.Nil(err, "logout returns success")
}
