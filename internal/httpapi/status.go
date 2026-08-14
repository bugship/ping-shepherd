package httpapi

import (
	"embed"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/bugship/ping-shepherd/internal/store"
)

//go:embed status.html
var statusFS embed.FS

var statusTmpl = template.Must(template.ParseFS(statusFS, "status.html"))

type statusRow struct {
	Name       string
	URL        string
	HasCheck   bool
	Status     string
	Class      string
	StatusCode int
	LatencyMS  int64
	CheckedAt  time.Time
}

type statusPageData struct {
	Now  time.Time
	Rows []statusRow
}

func statusPage(db Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			http.Error(w, "no store", http.StatusServiceUnavailable)
			return
		}
		targets, err := db.ListTargets(r.Context())
		if err != nil {
			http.Error(w, "store error", http.StatusInternalServerError)
			return
		}
		data := statusPageData{Now: time.Now().UTC()}
		for _, t := range targets {
			row := statusRow{Name: t.Name, URL: t.URL, Status: "unknown", Class: "unknown"}
			c, err := db.LatestCheck(r.Context(), t.ID)
			if err == nil {
				row.HasCheck = true
				row.StatusCode = c.StatusCode
				row.LatencyMS = c.LatencyMS
				row.CheckedAt = c.CheckedAt
				if c.Up {
					row.Status = "up"
					row.Class = "up"
				} else {
					row.Status = "down"
					row.Class = "down"
				}
			} else if !errors.Is(err, store.ErrNotFound) {
				http.Error(w, "store error", http.StatusInternalServerError)
				return
			}
			data.Rows = append(data.Rows, row)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := statusTmpl.Execute(w, data); err != nil {
			log.Printf("status page: %v", err)
		}
	}
}
