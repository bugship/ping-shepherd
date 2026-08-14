package httpapi

import (
	"net/http"

	"github.com/bugship/canary-coop/internal/store"
)

func New(db *store.DB) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("GET /ready", ready(db))
	return mux
}
