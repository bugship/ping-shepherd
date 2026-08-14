package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bugship/ping-shepherd/internal/store"
)

func TestTargetsCRUD(t *testing.T) {
	srv := httptest.NewServer(New(store.NewMemory()))
	t.Cleanup(srv.Close)

	body := bytes.NewBufferString(`{"name":"example","url":"https://example.com"}`)
	res, err := http.Post(srv.URL+"/targets", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create status %d", res.StatusCode)
	}
	var created store.Target
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.Name != "example" {
		t.Fatalf("created %#v", created)
	}

	res, err = http.Get(srv.URL + "/targets")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var list []store.Target
	if err := json.NewDecoder(res.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("list %#v", list)
	}

	res, err = http.Get(srv.URL + "/targets/" + created.ID)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("get status %d", res.StatusCode)
	}

	req, err := http.NewRequest(http.MethodDelete, srv.URL+"/targets/"+created.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("delete status %d", res.StatusCode)
	}
}

func TestCreateTargetRejectsBadURL(t *testing.T) {
	srv := httptest.NewServer(New(store.NewMemory()))
	t.Cleanup(srv.Close)

	body := bytes.NewBufferString(`{"name":"bad","url":"not-a-url"}`)
	res, err := http.Post(srv.URL+"/targets", "application/json", body)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status %d", res.StatusCode)
	}
}
