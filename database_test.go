package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// stubLogger satisfies service.Logger with no-ops so tests don't panic
// on the nil global logger.
type stubLogger struct{}

func (stubLogger) Error(v ...any) error                 { return nil }
func (stubLogger) Warning(v ...any) error               { return nil }
func (stubLogger) Info(v ...any) error                  { return nil }
func (stubLogger) Errorf(format string, a ...any) error { return nil }
func (stubLogger) Warningf(format string, a ...any) error { return nil }
func (stubLogger) Infof(format string, a ...any) error  { return nil }

func TestMain(m *testing.M) {
	logger = stubLogger{}
	os.Exit(m.Run())
}

func setupTestDB(t *testing.T) func() {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	if err := initDatabase(path); err != nil {
		t.Fatalf("initDatabase: %v", err)
	}
	return func() {
		if err := closeDatabase(); err != nil {
			t.Errorf("closeDatabase: %v", err)
		}
	}
}

func TestInitDatabase(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "test.db")

	if err := initDatabase(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer closeDatabase()

	if _, err := os.Stat(path); err != nil {
		t.Errorf("database file not created: %v", err)
	}
}

func TestSaveAndQueryMetrics(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	now := time.Now().Truncate(time.Millisecond)
	want := []Metrics{
		{Timestamp: now.Add(-2 * time.Minute), NetworkName: "HomeWifi", Online: false},
		{Timestamp: now.Add(-1 * time.Minute), NetworkName: "HomeWifi", Online: false},
		{Timestamp: now, NetworkName: "HomeWifi", Online: false},
	}

	for i := range want {
		if err := saveMetric(&want[i]); err != nil {
			t.Fatalf("saveMetric[%d]: %v", i, err)
		}
	}

	got, err := getMetricsSince(now.Add(-3 * time.Minute))
	if err != nil {
		t.Fatalf("getMetricsSince: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d rows, want %d", len(got), len(want))
	}
	for i := range want {
		if !got[i].Timestamp.Equal(want[i].Timestamp) {
			t.Errorf("[%d] timestamp: got %v, want %v", i, got[i].Timestamp, want[i].Timestamp)
		}
		if got[i].NetworkName != want[i].NetworkName {
			t.Errorf("[%d] network_name: got %q, want %q", i, got[i].NetworkName, want[i].NetworkName)
		}
		if got[i].Online != want[i].Online {
			t.Errorf("[%d] online: got %v, want %v", i, got[i].Online, want[i].Online)
		}
	}
}

func TestGetMetricsSince_FiltersOldRows(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	now := time.Now().Truncate(time.Millisecond)
	old := Metrics{Timestamp: now.Add(-2 * time.Hour), NetworkName: "net", Online: false}
	recent := Metrics{Timestamp: now.Add(-30 * time.Minute), NetworkName: "net", Online: false}

	for _, m := range []Metrics{old, recent} {
		if err := saveMetric(&m); err != nil {
			t.Fatalf("saveMetric: %v", err)
		}
	}

	got, err := getMetricsSince(now.Add(-1 * time.Hour))
	if err != nil {
		t.Fatalf("getMetricsSince: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d rows, want 1", len(got))
	}
	if !got[0].Timestamp.Equal(recent.Timestamp) {
		t.Errorf("expected recent row, got %v", got[0].Timestamp)
	}
}

func TestGetMetricsSince_Empty(t *testing.T) {
	teardown := setupTestDB(t)
	defer teardown()

	got, err := getMetricsSince(time.Now().Add(-1 * time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("got %d rows, want 0", len(got))
	}
}

func TestCloseDatabase_Idempotent(t *testing.T) {
	teardown := setupTestDB(t)
	teardown() // first close

	// second close should not panic or error
	if err := closeDatabase(); err != nil {
		t.Errorf("second closeDatabase returned error: %v", err)
	}
}
