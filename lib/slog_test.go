package lib

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/require"
)

func slogTestRecords(t *testing.T, output []byte) []map[string]any {
	t.Helper()
	var records []map[string]any
	decoder := json.NewDecoder(bytes.NewReader(output))
	for {
		var record map[string]any
		err := decoder.Decode(&record)
		if err == io.EOF {
			return records
		}
		require.NoError(t, err)
		records = append(records, record)
	}
}

func TestSlogLoggerUsesHandlerLevelsAndFormatsSDKMessages(t *testing.T) {
	var output bytes.Buffer
	var level slog.LevelVar
	logger := NewSlogLogger(slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: &level})).With("component", "sdk"))
	leveled := NewLeveledLogger(logger)
	logger.Printf("processed %d%% of %s", 25, "report.csv")
	leveled.Info("attempt %d", 1)
	leveled.Debug("hidden %d", 1)
	level.Set(slog.LevelDebug)
	leveled.Debug("retry %d", 2)
	leveled.Trace("read %d bytes", 32)
	level.Set(slog.LevelWarn)
	leveled.Info("hidden %d", 2)
	leveled.Warn("wait %d seconds", 7)
	leveled.Error("failed %s", "report.csv")

	records := slogTestRecords(t, output.Bytes())
	require.Len(t, records, 6)
	for i, want := range []struct{ level, message string }{
		{"INFO", "processed 25% of report.csv"},
		{"INFO", "attempt 1"},
		{"DEBUG", "retry 2"},
		{"DEBUG", "read 32 bytes"},
		{"WARN", "wait 7 seconds"},
		{"ERROR", "failed report.csv"},
	} {
		require.Equal(t, want.level, records[i]["level"])
		require.Equal(t, want.message, records[i]["msg"])
		require.Equal(t, "sdk", records[i]["component"], "the caller's logger attributes are preserved")
		require.NotContains(t, records[i], "!BADKEY")
	}
}

func TestDefaultRetryableHttpKeepsSlogDependencyAttributes(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = io.WriteString(w, "retried")
	}))
	defer server.Close()

	var output bytes.Buffer
	logger := NewSlogLogger(slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})).With("component", "retry"))
	client := DefaultRetryableHttp(logger)
	request, err := retryablehttp.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	response, err := client.Do(request)
	require.NoError(t, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, "retried", string(body))
	require.Equal(t, int32(2), attempts.Load())

	performing := false
	for _, record := range slogTestRecords(t, output.Bytes()) {
		require.Equal(t, "retry", record["component"])
		require.NotContains(t, record, "!BADKEY")
		require.NotContains(t, record["msg"], "%!")
		if record["msg"] == "performing request" {
			performing = true
			require.Equal(t, "GET", record["method"])
			require.Equal(t, server.URL, record["url"])
		}
	}
	require.True(t, performing, "dependency key/value fields remain structured slog attributes")
}
