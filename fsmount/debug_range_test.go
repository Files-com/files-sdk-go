//go:build (linux || windows) && filescomfs_debug

package fsmount

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDebugRangeStatsReportsModeAndResetsDiagnostics(t *testing.T) {
	remote := &RemoteFs{sparseRangeReads: true}
	coordinator := remote.rangeDownloads()
	originalStart := time.Unix(1, 0).UTC()
	coordinator.mu.Lock()
	coordinator.diagnostics.StartedAt = originalStart
	coordinator.diagnostics.RequestsStarted = 3
	coordinator.mu.Unlock()
	registry := &mountRegistry{hosts: map[string]*Host{
		"X:": {fs: &Filescomfs{remote: remote}},
	}}

	getSnapshot := requestRangeStats(t, registry, http.MethodGet, http.StatusOK)
	if !getSnapshot.SparseRangeReadsEnabled || getSnapshot.RequestsStarted != 3 || !getSnapshot.StartedAt.Equal(originalStart) {
		t.Fatalf("GET snapshot = %#v", getSnapshot)
	}

	resetSnapshot := requestRangeStats(t, registry, http.MethodPost, http.StatusOK)
	if !resetSnapshot.SparseRangeReadsEnabled || resetSnapshot.RequestsStarted != 0 || !resetSnapshot.StartedAt.After(originalStart) {
		t.Fatalf("POST snapshot = %#v", resetSnapshot)
	}
}

func TestDebugRangeStatsRejectsResetWhileRequestIsActive(t *testing.T) {
	remote := &RemoteFs{}
	coordinator := remote.rangeDownloads()
	coordinator.mu.Lock()
	coordinator.diagnostics.ActiveRequests = 1
	coordinator.mu.Unlock()
	registry := &mountRegistry{hosts: map[string]*Host{
		"X:": {fs: &Filescomfs{remote: remote}},
	}}

	request := httptest.NewRequest(http.MethodPost, "/debug/ranges?mnt=X%3A", nil)
	response := httptest.NewRecorder()
	registry.handleDebugRangeStats(response, request)
	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusConflict)
	}
}

func TestDebugRangeStatsHandlesClosedCoordinator(t *testing.T) {
	remote := &RemoteFs{}
	remote.closeRangeDownloads()
	registry := &mountRegistry{hosts: map[string]*Host{
		"X:": {fs: &Filescomfs{remote: remote}},
	}}

	request := httptest.NewRequest(http.MethodGet, "/debug/ranges?mnt=X%3A", nil)
	response := httptest.NewRecorder()
	registry.handleDebugRangeStats(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}
}

func requestRangeStats(t *testing.T, registry *mountRegistry, method string, wantStatus int) rangeStreamDiagnostics {
	t.Helper()
	request := httptest.NewRequest(method, "/debug/ranges?mnt=X%3A", nil)
	response := httptest.NewRecorder()
	registry.handleDebugRangeStats(response, request)
	if response.Code != wantStatus {
		t.Fatalf("%s status = %d, want %d", method, response.Code, wantStatus)
	}

	var snapshot rangeStreamDiagnostics
	if err := json.NewDecoder(response.Body).Decode(&snapshot); err != nil {
		t.Fatalf("%s response decode failed: %v", method, err)
	}
	return snapshot
}
