package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bugship/ping-shepherd/internal/store"
)

type memSender struct {
	msgs []string
}

func (m *memSender) Send(_ context.Context, text string) error {
	m.msgs = append(m.msgs, text)
	return nil
}

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
	if err := Once(t.Context(), mem, time.Second, nil); err != nil {
		t.Fatal(err)
	}
	got, err := mem.LatestCheck(t.Context(), target.ID)
	if err != nil || !got.Up || got.StatusCode != http.StatusOK {
		t.Fatalf("check: %v %#v", err, got)
	}
}

func TestOnceAlertsOnDownAndRecovery(t *testing.T) {
	var code int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(code)
	}))
	t.Cleanup(srv.Close)

	mem := store.NewMemory()
	if _, err := mem.CreateTarget(t.Context(), store.NewTarget{Name: "example", URL: srv.URL}); err != nil {
		t.Fatal(err)
	}
	alerts := &memSender{}

	code = http.StatusInternalServerError
	if err := Once(t.Context(), mem, time.Second, alerts); err != nil {
		t.Fatal(err)
	}
	if len(alerts.msgs) != 1 || !strings.HasPrefix(alerts.msgs[0], "DOWN ") {
		t.Fatalf("down alert: %#v", alerts.msgs)
	}

	code = http.StatusOK
	if err := Once(t.Context(), mem, time.Second, alerts); err != nil {
		t.Fatal(err)
	}
	if len(alerts.msgs) != 2 || !strings.HasPrefix(alerts.msgs[1], "UP ") {
		t.Fatalf("up alert: %#v", alerts.msgs)
	}

	if err := Once(t.Context(), mem, time.Second, alerts); err != nil {
		t.Fatal(err)
	}
	if len(alerts.msgs) != 2 {
		t.Fatalf("should not re-alert: %#v", alerts.msgs)
	}
}
