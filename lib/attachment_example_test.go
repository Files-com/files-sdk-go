package lib_test

import (
	"context"
	"log"
	"net/http"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/user"
)

// A file's size is known, so NewJSONAttachment checks the file has exactly that
// many bytes when the request reads it.
func ExampleNewJSONAttachment() {
	ctx := context.Background()
	client := user.Client{Config: files_sdk.Config{APIKey: os.Getenv("FILES_API_KEY")}.Init()}

	avatar, err := os.Open("avatar.png")
	if err != nil {
		log.Fatal(err)
	}
	// The SDK reads the file but leaves closing it to you.
	defer avatar.Close()
	info, err := avatar.Stat()
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Update(files_sdk.UserUpdateParams{
		Id:                   123,
		Name:                 "Avatar Owner",
		AvatarFileAttachment: lib.NewJSONAttachment(avatar, info.Size(), "avatar.png", "image/png"),
	}, files_sdk.WithContext(ctx))
	if err != nil {
		log.Fatal(err)
	}
}

// A response body's length is not known in advance, so NewJSONAttachmentUpTo
// sets the most it may have.
func ExampleNewJSONAttachmentUpTo() {
	ctx := context.Background()
	client := user.Client{Config: files_sdk.Config{APIKey: os.Getenv("FILES_API_KEY")}.Init()}

	response, err := http.Get("https://example.com/avatar.png")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	_, err = client.Update(files_sdk.UserUpdateParams{
		Id:                   123,
		AvatarFileAttachment: lib.NewJSONAttachmentUpTo(response.Body, 5<<20, "avatar.png", response.Header.Get("Content-Type")),
	}, files_sdk.WithContext(ctx))
	if err != nil {
		log.Fatal(err)
	}
}
