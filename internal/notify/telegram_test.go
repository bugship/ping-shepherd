package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTelegramSend(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/botTEST/sendMessage" {
			t.Fatalf("path %s", r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	tg := Telegram{Token: "TEST", ChatID: "123", Base: srv.URL}
	if err := tg.Send(t.Context(), "DOWN example"); err != nil {
		t.Fatal(err)
	}
	if got["chat_id"] != "123" || got["text"] != "DOWN example" {
		t.Fatalf("payload %#v", got)
	}
}

func TestTelegramRequiresConfig(t *testing.T) {
	tg := Telegram{}
	if err := tg.Send(t.Context(), "hi"); err == nil {
		t.Fatal("expected error")
	}
}
