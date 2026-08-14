package probe

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestURLUp(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	got := URL(t.Context(), srv.URL, time.Second)
	if !got.Up || got.StatusCode != http.StatusOK {
		t.Fatalf("got %#v", got)
	}
}

func TestURLDownOnServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	got := URL(t.Context(), srv.URL, time.Second)
	if got.Up || got.StatusCode != http.StatusInternalServerError {
		t.Fatalf("got %#v", got)
	}
}

func TestURLDownOnUnreachable(t *testing.T) {
	got := URL(t.Context(), "http://127.0.0.1:1", 200*time.Millisecond)
	if got.Up || got.Err == "" {
		t.Fatalf("got %#v", got)
	}
}
