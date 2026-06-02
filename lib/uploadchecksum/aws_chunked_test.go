package uploadchecksum

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAWSChunkedEncoderWritesCRC32CTrailer(t *testing.T) {
	encoder, err := NewAWSChunkedEncoder(strings.NewReader("hello"), AWSChunkedEncoderOptions{
		Algorithm: ChecksumAlgorithmCRC32C,
		ChunkSize: 2,
	})
	require.NoError(t, err)

	body, err := io.ReadAll(encoder)
	require.NoError(t, err)
	checksum, err := Checksum([]byte("hello"), ChecksumAlgorithmCRC32C)
	require.NoError(t, err)
	require.Equal(t, "2\r\nhe\r\n2\r\nll\r\n1\r\no\r\n0\r\nx-amz-checksum-crc32c:"+checksum+"\r\n\r\n", string(body))

	encodedLength, err := AWSChunkedEncodedLength(5, 2, ChecksumAlgorithmCRC32C)
	require.NoError(t, err)
	require.Equal(t, int64(len(body)), encodedLength)
}

func TestAWSChunkedEncoderWritesCRC32Trailer(t *testing.T) {
	encoder, err := NewAWSChunkedEncoder(strings.NewReader("hello"), AWSChunkedEncoderOptions{
		Algorithm: ChecksumAlgorithmCRC32,
		ChunkSize: 5,
	})
	require.NoError(t, err)

	body, err := io.ReadAll(encoder)
	require.NoError(t, err)
	checksum, err := Checksum([]byte("hello"), ChecksumAlgorithmCRC32)
	require.NoError(t, err)
	require.Equal(t, "5\r\nhello\r\n0\r\nx-amz-checksum-crc32:"+checksum+"\r\n\r\n", string(body))

	decoder, err := NewAWSChunkedDecoder(strings.NewReader(string(body)), AWSChunkedDecoderOptions{
		Algorithm: ChecksumAlgorithmCRC32,
	})
	require.NoError(t, err)
	decoded, err := io.ReadAll(decoder)
	require.NoError(t, err)
	require.Equal(t, "hello", string(decoded))
	require.Equal(t, checksum, decoder.TrailerValue())
}

func TestAWSChunkedEncoderSupportsTinyReadBuffers(t *testing.T) {
	encoder, err := NewAWSChunkedEncoder(strings.NewReader("hello"), AWSChunkedEncoderOptions{
		Algorithm: ChecksumAlgorithmCRC32,
		ChunkSize: 2,
	})
	require.NoError(t, err)

	var body strings.Builder
	buf := make([]byte, 1)
	for {
		n, err := encoder.Read(buf)
		if n > 0 {
			body.Write(buf[:n])
		}
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}

	checksum, err := Checksum([]byte("hello"), ChecksumAlgorithmCRC32)
	require.NoError(t, err)
	require.Equal(t, "2\r\nhe\r\n2\r\nll\r\n1\r\no\r\n0\r\nx-amz-checksum-crc32:"+checksum+"\r\n\r\n", body.String())
}

func TestAWSChunkedDecoderVerifiesAndReturnsDecodedBody(t *testing.T) {
	encoder, err := NewAWSChunkedEncoder(strings.NewReader("hello world"), AWSChunkedEncoderOptions{
		Algorithm: ChecksumAlgorithmCRC32C,
		ChunkSize: 3,
	})
	require.NoError(t, err)

	decoder, err := NewAWSChunkedDecoder(encoder, AWSChunkedDecoderOptions{
		Algorithm: ChecksumAlgorithmCRC32C,
	})
	require.NoError(t, err)
	body, err := io.ReadAll(decoder)
	require.NoError(t, err)
	require.Equal(t, "hello world", string(body))

	checksum, err := Checksum([]byte("hello world"), ChecksumAlgorithmCRC32C)
	require.NoError(t, err)
	require.Equal(t, checksum, decoder.TrailerValue())
}

func TestAWSChunkedDecoderForHeaders(t *testing.T) {
	encoder, err := NewAWSChunkedEncoder(strings.NewReader("hello"), AWSChunkedEncoderOptions{
		Algorithm: ChecksumAlgorithmCRC32,
		ChunkSize: 5,
	})
	require.NoError(t, err)

	headers := http.Header{}
	headers.Set("x-amz-trailer", "x-amz-checksum-crc32")
	decoder, err := NewAWSChunkedDecoderForHeaders(encoder, headers)
	require.NoError(t, err)

	body, err := io.ReadAll(decoder)
	require.NoError(t, err)
	require.Equal(t, "hello", string(body))
}

func TestAWSChunkedDecoderRejectsChecksumMismatch(t *testing.T) {
	body := "5\r\nhello\r\n0\r\nx-amz-checksum-crc32c:AAAAAA==\r\n\r\n"
	decoder, err := NewAWSChunkedDecoder(strings.NewReader(body), AWSChunkedDecoderOptions{
		Algorithm: ChecksumAlgorithmCRC32C,
	})
	require.NoError(t, err)

	_, err = io.ReadAll(decoder)
	require.ErrorIs(t, err, ErrChecksumMismatch)
}

func TestAWSChunkedDecoderRejectsMissingChecksum(t *testing.T) {
	body := "5\r\nhello\r\n0\r\nx-amz-meta-test:value\r\n\r\n"
	decoder, err := NewAWSChunkedDecoder(strings.NewReader(body), AWSChunkedDecoderOptions{
		Algorithm: ChecksumAlgorithmCRC32C,
	})
	require.NoError(t, err)

	_, err = io.ReadAll(decoder)
	require.ErrorIs(t, err, ErrMissingChecksum)
}

func TestAWSChunkedDecoderAcceptsChunkExtensions(t *testing.T) {
	checksum, err := Checksum([]byte("hello"), ChecksumAlgorithmCRC32C)
	require.NoError(t, err)
	body := "5;chunk-signature=unused\r\nhello\r\n0\r\nx-amz-checksum-crc32c:" + checksum + "\r\n\r\n"
	decoder, err := NewAWSChunkedDecoder(strings.NewReader(body), AWSChunkedDecoderOptions{
		Algorithm: ChecksumAlgorithmCRC32C,
	})
	require.NoError(t, err)

	decoded, err := io.ReadAll(decoder)
	require.NoError(t, err)
	require.Equal(t, "hello", string(decoded))
}

func TestS3TrailerHeaders(t *testing.T) {
	headers, encodedLength, err := S3TrailerHeaders(S3TrailerHeadersOptions{
		Algorithm:     ChecksumAlgorithmCRC32C,
		DecodedLength: 5,
		ChunkSize:     2,
	})
	require.NoError(t, err)
	require.Equal(t, "aws-chunked", headers.Get("Content-Encoding"))
	require.Equal(t, "5", headers.Get("x-amz-decoded-content-length"))
	require.Equal(t, "CRC32C", headers.Get("x-amz-sdk-checksum-algorithm"))
	require.Equal(t, "STREAMING-UNSIGNED-PAYLOAD-TRAILER", headers.Get("x-amz-content-sha256"))
	require.Equal(t, "x-amz-checksum-crc32c", headers.Get("x-amz-trailer"))
	require.Equal(t, int64(57), encodedLength)
	require.Equal(t, "57", headers.Get("Content-Length"))
}

func TestS3TrailerHeadersDefaultToBestPlatformAlgorithm(t *testing.T) {
	headers, _, err := S3TrailerHeaders(S3TrailerHeadersOptions{
		DecodedLength: 5,
		ChunkSize:     5,
	})
	require.NoError(t, err)
	require.Contains(t, []string{"CRC32", "CRC32C"}, headers.Get("x-amz-sdk-checksum-algorithm"))
	require.Contains(t, []string{"x-amz-checksum-crc32", "x-amz-checksum-crc32c"}, headers.Get("x-amz-trailer"))
}

func TestIsAWSChunked(t *testing.T) {
	headers := http.Header{}
	headers.Set("Content-Encoding", "gzip, aws-chunked")
	require.True(t, IsAWSChunked(headers))
}

func TestAlgorithmFromHeaders(t *testing.T) {
	headers := http.Header{}
	headers.Set("x-amz-trailer", "x-amz-checksum-crc32")
	headers.Set("x-amz-sdk-checksum-algorithm", "CRC32C")

	algorithm, err := AlgorithmFromHeaders(headers)
	require.NoError(t, err)
	require.Equal(t, ChecksumAlgorithmCRC32, algorithm)

	headers = http.Header{}
	headers.Set("x-amz-sdk-checksum-algorithm", "crc32c")
	algorithm, err = AlgorithmFromHeaders(headers)
	require.NoError(t, err)
	require.Equal(t, ChecksumAlgorithmCRC32C, algorithm)

	headers = http.Header{}
	_, err = AlgorithmFromHeaders(headers)
	require.ErrorIs(t, err, ErrUnsupportedAlgorithm)
}

func TestUnsupportedAlgorithm(t *testing.T) {
	_, err := NewAWSChunkedEncoder(strings.NewReader("hello"), AWSChunkedEncoderOptions{
		Algorithm: Algorithm("SHA256"),
	})
	require.ErrorIs(t, err, ErrUnsupportedAlgorithm)

	_, err = AWSChunkedEncodedLength(1, 1, Algorithm("SHA256"))
	require.ErrorIs(t, err, ErrUnsupportedAlgorithm)
}

func TestDecoderMalformedBody(t *testing.T) {
	decoder, err := NewAWSChunkedDecoder(strings.NewReader("z\r\n"), AWSChunkedDecoderOptions{
		Algorithm: ChecksumAlgorithmCRC32C,
	})
	require.NoError(t, err)

	_, err = io.ReadAll(decoder)
	require.True(t, errors.Is(err, ErrMalformedAWSChunked), err)
}

func BenchmarkAWSChunkedEncoder64MiB(b *testing.B) {
	const size = int64(64 * 1024 * 1024)
	buffer := make([]byte, 32*1024)
	b.SetBytes(size)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		encoder, err := NewAWSChunkedEncoder(&zeroReader{remaining: size}, AWSChunkedEncoderOptions{
			Algorithm: ChecksumAlgorithmCRC32C,
			ChunkSize: DefaultChunkSize,
		})
		require.NoError(b, err)
		n, err := io.CopyBuffer(io.Discard, encoder, buffer)
		require.NoError(b, err)
		if n <= size {
			b.Fatalf("encoded bytes = %d, want more than decoded size %d", n, size)
		}
	}
}

type zeroReader struct {
	remaining int64
}

func (r *zeroReader) Read(p []byte) (int, error) {
	if r.remaining <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > r.remaining {
		p = p[:r.remaining]
	}
	clear(p)
	r.remaining -= int64(len(p))
	return len(p), nil
}
