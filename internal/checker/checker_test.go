package checker

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bugship/ping-shepherd/internal/store"
)

func TestOnceRecordsProbe(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	mem := store.NewMemory()
	target, err := mem.CreateTarget(t.Context(), store.NewTarget{Name: "example", URL: srv.URL})
	if err != nil {
		t.Fatal(err)
	}
	if err := Once(t.Context(), mem, time.Second); err != nil {
		t.Fatal(err)
	}
	got, err := mem.LatestCheck(t.Context(), target.ID)
	if err != nil || !got.Up || got.StatusCode != http.StatusOK {
		t.Fatalf("check: %v %#v", err, got)
	}
}
