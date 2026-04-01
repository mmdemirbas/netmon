package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func setupTestServer(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	if err := initDatabase(filepath.Join(dir, "test.db")); err != nil {
		t.Fatalf("initDatabase: %v", err)
	}
	return func() { closeDatabase() }
}

func TestHandleMetrics_Empty(t *testing.T) {
	teardown := setupTestServer(t)
	defer teardown()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	handleMetrics(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status %d, want %d", rr.Code, http.StatusOK)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type %q, want application/json", ct)
	}
	var dtos []MetricsDto
	if err := json.Unmarshal(rr.Body.Bytes(), &dtos); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(dtos) != 0 {
		t.Errorf("got %d items, want 0", len(dtos))
	}
}

func TestHandleMetrics_ReturnsOnlyWindowData(t *testing.T) {
	teardown := setupTestServer(t)
	defer teardown()

	now := time.Now()
	rows := []Metrics{
		{Timestamp: now.Add(-2 * time.Hour), NetworkName: "net", Online: false}, // outside default 24h window only when ?since provided
		{Timestamp: now.Add(-30 * time.Minute), NetworkName: "net", Online: false},
	}
	for i := range rows {
		if err := saveMetric(&rows[i]); err != nil {
			t.Fatalf("saveMetric: %v", err)
		}
	}

	// ?since = 1 hour ago → should return only the second row
	sinceMs := now.Add(-1 * time.Hour).UnixMilli()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	q := req.URL.Query()
	q.Set("since", strconv.FormatInt(sinceMs, 10))
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handleMetrics(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status %d, want 200", rr.Code)
	}
	var dtos []MetricsDto
	if err := json.Unmarshal(rr.Body.Bytes(), &dtos); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(dtos) != 1 {
		t.Fatalf("got %d items, want 1", len(dtos))
	}
}

func TestHandleMetrics_InvalidSince(t *testing.T) {
	teardown := setupTestServer(t)
	defer teardown()

	req := httptest.NewRequest(http.MethodGet, "/metrics?since=notanumber", nil)
	rr := httptest.NewRecorder()
	handleMetrics(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status %d, want 400", rr.Code)
	}
}

func TestHandleMetrics_OfflineRow(t *testing.T) {
	teardown := setupTestServer(t)
	defer teardown()

	m := Metrics{Timestamp: time.Now(), NetworkName: "wifi", Online: false}
	if err := saveMetric(&m); err != nil {
		t.Fatalf("saveMetric: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	handleMetrics(rr, req)

	var dtos []MetricsDto
	if err := json.Unmarshal(rr.Body.Bytes(), &dtos); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(dtos) != 1 {
		t.Fatalf("got %d items, want 1", len(dtos))
	}
	if dtos[0].IsOnline {
		t.Errorf("IsOnline = true, want false")
	}
	if dtos[0].NetworkName != "wifi" {
		t.Errorf("NetworkName = %q, want wifi", dtos[0].NetworkName)
	}
}

