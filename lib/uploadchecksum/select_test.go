package uploadchecksum

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBestAlgorithmForPlatformReturnsSupportedS3Algorithm(t *testing.T) {
	algorithm := BestAlgorithmForPlatform()
	require.Contains(t, SupportedS3Algorithms(), algorithm)
}

func TestSelectBestAlgorithmSupportsExplicitCandidates(t *testing.T) {
	algorithm, err := SelectBestAlgorithm(BestAlgorithmOptions{
		Candidates: []Algorithm{ChecksumAlgorithmCRC32},
		SampleSize: 128,
		Iterations: 1,
	})
	require.NoError(t, err)
	require.Equal(t, ChecksumAlgorithmCRC32, algorithm)
}

func TestBenchmarkAlgorithms(t *testing.T) {
	results, err := BenchmarkAlgorithms(BestAlgorithmOptions{
		SampleSize: 128,
		Iterations: 1,
	})
	require.NoError(t, err)
	require.Len(t, results, 2)
	for _, result := range results {
		require.Contains(t, SupportedS3Algorithms(), result.Algorithm)
		require.Equal(t, int64(128), result.BytesProcessed)
		require.Greater(t, result.Duration, time.Duration(0))
	}
}

func TestChecksumSupportsCRC32AndCRC32C(t *testing.T) {
	crc32Value, err := Checksum([]byte("hello"), ChecksumAlgorithmCRC32)
	require.NoError(t, err)
	require.Equal(t, "NhCmhg==", crc32Value)

	crc32cValue, err := Checksum([]byte("hello"), ChecksumAlgorithmCRC32C)
	require.NoError(t, err)
	require.Equal(t, "mnG7TA==", crc32cValue)
}
