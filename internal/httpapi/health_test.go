package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	srv := httptest.NewServer(New(nil))
	t.Cleanup(srv.Close)

	res, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status %d", res.StatusCode)
	}
	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.fatal(err)
	}
	if body["status"] != "ok" || body["service"] != "canary-coop" {
		t.Fatalf("body %#v", body)
	}
}

func TestReadyWithoutDB(t *testing.T) {
	srv := httptest.NewServer(New(nil))
	t.Cleanup(srv.Close)

	res, err := http.Get(srv.URL + "/ready")
	if err != nil {
		t.fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status %d", res.StatusCode)
	}
}
