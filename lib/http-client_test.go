package lib

import (
	"net/http"
	"testing"
	"time"
)

func TestDefaultRetryableHttpBackoffHonorsRetryAfterForRetryableStatuses(t *testing.T) {
	client := DefaultRetryableHttp(nil)
	for _, status := range []int{
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			resp := &http.Response{
				StatusCode: status,
				Header:     http.Header{"Retry-After": []string{"7"}},
			}

			if got := client.Backoff(0, 0, 1, resp); got != 7*time.Second {
				t.Fatalf("Backoff() = %v, want 7s", got)
			}
		})
	}
}

func TestRetryAfterDurationParsesHTTPDate(t *testing.T) {
	now := time.Date(2026, 5, 11, 12, 0, 0, 0, time.UTC)
	retryAt := now.Add(9 * time.Second).Format(http.TimeFormat)

	got, ok := retryAfterDuration(retryAt, now)
	if !ok {
		t.Fatal("retryAfterDuration() ok = false, want true")
	}
	if got != 9*time.Second {
		t.Fatalf("retryAfterDuration() = %v, want 9s", got)
	}
}
