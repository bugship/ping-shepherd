package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bugship/ping-shepherd/internal/store"
)

func TestListChecks(t *testing.T) {
	mem := store.NewMemory()
	target, err := mem.CreateTarget(t.Context(), store.NewTarget{Name: "example", URL: "https://example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := mem.RecordCheck(t.Context(), target.ID, store.Check{Up: true, StatusCode: 200}); err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(New(mem))
	t.Cleanup(srv.Close)

	res, err := http.Get(srv.URL + "/targets/" + target.ID + "/checks")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status %d", res.StatusCode)
	}
	var list []store.Check
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || !list[0].Up {
		t.Fatalf("list %#v", list)
	}
}
