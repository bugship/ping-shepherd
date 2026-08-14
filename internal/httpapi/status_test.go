package httpapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bugship/ping-shepherd/internal/store"
)

func TestStatusPage(t *testing.T) {
	mem := store.NewMemory()
	target, err := mem.CreateTarget(t.Context(), store.NewTarget{Name: "example", URL: "https://example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := mem.RecordCheck(t.Context(), target.ID, store.Check{Up: true, StatusCode: 200, LatencyMS: 8}); err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(New(mem))
	t.Cleanup(srv.Close)

	res, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	html := string(body)
	if !strings.Contains(html, "example") || !strings.Contains(html, "up") {
		t.Fatalf("page missing target: %s", html)
	}
}
