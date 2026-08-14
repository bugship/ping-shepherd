package store

import "testing"

func TestNormalizeTarget(t *testing.T) {
	_, err := NormalizeTarget(NewTarget{Name: "  ", URL: "https://example.com"})
	if err != ErrInvalidName {
		t.Fatalf("name: %v", err)
	}
	_, err = NormalizeTarget(NewTarget{Name: "example", URL: "ftp://example.com"})
	if err != ErrInvalidURL {
		t.Fatalf("scheme: %v", err)
	}
	got, err := NormalizeTarget(NewTarget{Name: " example ", URL: "https://example.com/status"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "example" || got.URL != "https://example.com/status" {
		t.Fatalf("got %#v", got)
	}
}

func TestMemoryCreateListDelete(t *testing.T) {
	m := NewMemory()
	created, err := m.CreateTarget(t.Context(), NewTarget{Name: "example", URL: "https://example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == "" {
		t.Fatal("missing id")
	}
	list, err := m.ListTargets(t.Context())
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %v %#v", err, list)
	}
	got, err := m.GetTarget(t.Context(), created.ID)
	if err != nil || got.URL != "https://example.com" {
		t.Fatalf("get: %v %#v", err, got)
	}
	if err := m.DeleteTarget(t.Context(), created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := m.GetTarget(t.Context(), created.ID); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}
