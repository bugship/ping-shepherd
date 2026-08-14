package store

import (
	"context"
	"testing"
)

func TestOpenEmptyURL(t *testing.T) {
	_, err := Open(context.Background(), "")
	if err == nil {
		t.fatal("expected error for empty DATABASE_URL")
	}
}
