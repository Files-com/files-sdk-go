package file_migration

import (
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/lib/test"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func CreateClient(fixture string) (client *Client, r *recorder.Recorder, err error) {
	client = &Client{}
	client.Config, r, err = test.CreateConfig(fixture)

	return client, r, err
}

func TestClient_Wait(t *testing.T) {
	client, r, err := CreateClient("TestClient_Wait")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)
	var migrationPassedToFunc files_sdk.FileMigration
	statusFunc := func(migration files_sdk.FileMigration) {
		migrationPassedToFunc = migration
	}
	//
	fileAction := files_sdk.FileAction{Status: "in_progress", FileMigrationId: 11}

	migration, err := client.Wait(fileAction, statusFunc)

	assert.Error(err, "Not Found - `Not Found`")
	assert.Equal("", migration.Status)
	assert.Equal(int64(0), migration.Id)
	//
	fileAction = files_sdk.FileAction{Status: "in_progress", FileMigrationId: 12}

	migration, err = client.Wait(fileAction, statusFunc)

	assert.NoError(err)
	assert.Equal("completed", migrationPassedToFunc.Status)
	assert.Equal("completed", migration.Status)
	assert.Equal(int64(12), migration.Id)
	assert.Equal(int64(12), migrationPassedToFunc.Id)
	assert.Equal("video.mp4", migrationPassedToFunc.Path)

	//
	fileAction = files_sdk.FileAction{Status: "completed", FileMigrationId: 13}

	migration, err = client.Wait(fileAction, statusFunc)

	assert.NoError(err)
	assert.Equal("completed", migration.Status)
	assert.Equal(int64(13), migration.Id)

	//
	fileAction = files_sdk.FileAction{FileMigrationId: 14}

	migration, err = client.Wait(fileAction, statusFunc)

	assert.NoError(err)
	assert.Equal("failed", migration.Status)
	assert.Equal(int64(14), migration.Id)

	//
	fileAction = files_sdk.FileAction{FileMigrationId: 15}

	migration, err = client.Wait(fileAction, statusFunc)

	assert.NoError(err)
	assert.Equal("completed", migration.Status)
	assert.Equal(int64(15), migration.Id)
}
