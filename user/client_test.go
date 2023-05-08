package user

import (
	"context"
	"fmt"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/Files-com/files-sdk-go/v2/lib/test"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
)

func createClient(fixture string) (client *Client, r *recorder.Recorder, err error) {
	client = &Client{}
	client.Config, r, err = test.CreateConfig(fixture)

	return client, r, err
}

func TestClient_Create(t *testing.T) {
	client, r, err := createClient("TestClient_Create")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)

	user, err := findOrCreateUser(client, files_sdk.UserCreateParams{Username: "TestMo"})
	assert.NoError(err)

	assert.Equal("TestMo", user.Username)
}

func TestClient_Update(t *testing.T) {
	client, r, err := createClient("TestClient_Update")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)

	user, err := findOrCreateUser(client, files_sdk.UserCreateParams{Username: "TestMo"})
	assert.NoError(err)

	assert.Equal(lib.Bool(true), user.SftpPermission)

	user, err = client.Update(
		context.Background(),
		files_sdk.UserUpdateParams{
			Id:             user.Id,
			SftpPermission: lib.Bool(false),
		},
	)

	assert.NoError(err)

	assert.Equal(lib.Bool(false), user.SftpPermission)
}

func TestClient_List(t *testing.T) {
	client, r, err := createClient("TestClient_List")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	assert := assert.New(t)

	_, err = findOrCreateUser(client, files_sdk.UserCreateParams{Username: "test-list-user"})
	assert.NoError(err)

	it, err := client.List(context.Background(), files_sdk.UserListParams{})
	assert.NoError(err)
	var users []files_sdk.User
	for it.Next() {
		users = append(users, it.User())
		loaderUser, err := it.LoadResource(it.User().Identifier())
		assert.NoError(err)
		assert.Equal(loaderUser, it.User())
	}
	assert.NoError(it.Err())
	assert.Len(users, 1)
}

func findOrCreateUser(client *Client, params files_sdk.UserCreateParams) (files_sdk.User, error) {
	user, err := findUser(client, params)
	if err != nil && err.Error() == "user not found" {
		return client.Create(
			context.Background(),
			params,
		)
	}
	return user, err
}

func findUser(client *Client, params files_sdk.UserCreateParams) (files_sdk.User, error) {
	it, err := client.List(
		context.Background(),
		files_sdk.UserListParams{},
	)

	if err != nil {
		return files_sdk.User{}, err
	}
	var user *files_sdk.User
	for it.Next() {
		if it.User().Username == params.Username {
			u := it.User()
			user = &u
			continue
		}
	}
	if user == nil {
		return files_sdk.User{}, fmt.Errorf("user not found")
	}

	return *user, nil
}
