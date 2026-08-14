package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bugship/ping-shepherd/internal/store"
)

type targetRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func createTarget(db Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "no_store"})
			return
		}
		var req targetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_json"})
			return
		}
		t, err := db.CreateTarget(r.Context(), store.NewTarget{Name: req.Name, URL: req.URL})
		if writeStoreErr(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, t)
	}
}

func listTargets(db Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "no_store"})
			return
		}
		list, err := db.ListTargets(r.Context())
		if writeStoreErr(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, list)
	}
}

func getTarget(db Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "no_store"})
			return
		}
		t, err := db.GetTarget(r.Context(), r.PathValue("id"))
		if writeStoreErr(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, t)
	}
}

func deleteTarget(db Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "no_store"})
			return
		}
		err := db.DeleteTarget(r.Context(), r.PathValue("id"))
		if writeStoreErr(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func writeStoreErr(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, store.ErrInvalidName), errors.Is(err, store.ErrInvalidURL):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, store.ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "store_error"})
	}
	return true
}
