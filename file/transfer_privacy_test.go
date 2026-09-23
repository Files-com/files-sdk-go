package file

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/require"
)

type failingTransferTransport struct{ err error }

func (t failingTransferTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, t.err
}

func TestUploadV2NetworkErrorHidesSignedURL(t *testing.T) {
	for _, cause := range []error{os.ErrDeadlineExceeded, context.Canceled} {
		t.Run(cause.Error(), func(t *testing.T) {
			var logs bytes.Buffer
			config := files_sdk.Config{Logger: log.New(&logs, "", 0)}.Init()
			config.Client.RetryMax = 0
			config.Client.HTTPClient = &http.Client{Transport: failingTransferTransport{err: cause}}
			u := &uploadIO{Client: &Client{Config: config}}
			part := &uploadV2Part{}
			params := &files_sdk.CallParams{
				Method:  http.MethodPut,
				Config:  config,
				Uri:     "https://transfer.example.test/private-part?X-Amz-Credential=credential&X-Amz-Signature=signature",
				BodyIo:  io.NopCloser(bytes.NewReader([]byte("part"))),
				Context: context.Background(),
			}

			_, err := (&uploadV2Engine{u: u}).callUploadV2Part(context.Background(), part, params, func() bool { return true }, nil)

			require.ErrorIs(t, err, cause)
			var urlErr *url.Error
			require.ErrorAs(t, err, &urlErr)
			require.Equal(t, cause == os.ErrDeadlineExceeded, urlErr.Timeout())
			for _, secret := range []string{"transfer.example.test", "private-part", "credential", "signature"} {
				require.NotContains(t, err.Error(), secret)
				require.NotContains(t, logs.String(), secret)
			}
		})
	}
}
