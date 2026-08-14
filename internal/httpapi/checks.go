package httpapi

import (
	"net/http"
	"strconv"
)

func listChecks(db Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "no_store"})
			return
		}
		limit := 50
		if raw := r.URL.Query().Get("limit"); raw != "" {
			n, err := strconv.Atoi(raw)
			if err != nil || n <= 0 {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_limit"})
				return
			}
			limit = n
		}
		list, err := db.ListChecks(r.Context(), r.PathValue("id"), limit)
		if writeStoreErr(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, list)
	}
}
