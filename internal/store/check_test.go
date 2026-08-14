package store

import "testing"

func TestMemoryRecordAndListChecks(t *testing.T) {
	m := NewMemory()
	target, err := m.CreateTarget(t.Context(), NewTarget{Name: "example", URL: "https://example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := m.RecordCheck(t.Context(), target.ID, Check{Up: true, StatusCode: 200, LatencyMS: 12}); err != nil {
		t.Fatal(err)
	}
	if _, err := m.RecordCheck(t.Context(), target.ID, Check{Up: false, StatusCode: 500, LatencyMS: 40}); err != nil {
		t.Fatal(err)
	}
	latest, err := m.LatestCheck(t.Context(), target.ID)
	if err != nil || latest.Up || latest.StatusCode != 500 {
		t.Fatalf("latest: %v %#v", err, latest)
	}
	list, err := m.ListChecks(t.Context(), target.ID, 10)
	if err != nil || len(list) != 2 {
		t.Fatalf("list: %v %#v", err, list)
	}
	if _, err := m.RecordCheck(t.Context(), "missing", Check{Up: true}); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}
