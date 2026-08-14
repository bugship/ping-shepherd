// Package httpapi serves the Ping Shepherd HTTP API.
package httpapi

import (
	"context"
	"net/http"

	"github.com/bugship/ping-shepherd/internal/store"
)

// Backend is the storage used by HTTP handlers.
type Backend interface {
	Ping(context.Context) error
	CreateTarget(context.Context, store.NewTarget) (store.Target, error)
	ListTargets(context.Context) ([]store.Target, error)
	GetTarget(context.Context, string) (store.Target, error)
	DeleteTarget(context.Context, string) error
	RecordCheck(context.Context, string, store.Check) (store.Check, error)
	ListChecks(context.Context, string, int) ([]store.Check, error)
}

// New returns the HTTP handler.
func New(db Backend) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("GET /ready", ready(db))
	mux.HandleFunc("POST /targets", createTarget(db))
	mux.HandleFunc("GET /targets", listTargets(db))
	mux.HandleFunc("GET /targets/{id}", getTarget(db))
	mux.HandleFunc("DELETE /targets/{id}", deleteTarget(db))
	mux.HandleFunc("GET /targets/{id}/checks", listChecks(db))
	return mux
}
