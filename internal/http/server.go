package http

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/majori/wrc-laptimer/internal/database"
	"github.com/majori/wrc-laptimer/web"
)

func StartHTTPServer(db *database.Database, addr string) {
	mux := http.NewServeMux()

	// Add query endpoint
	mux.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("could not read request body", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reqBodyString := string(reqBody)

		result, err := db.ExecuteSelectQuery(reqBodyString)
		if err != nil {
			slog.Error("could not execute select query", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(result))
		if err != nil {
			slog.Error("could not write response", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Add series creation endpoint
	mux.HandleFunc("/api/admin/series/create", CreateSeriesHandler(db))
	mux.HandleFunc("/api/admin/series/{id}/start", StartSeriesHandler(db))
	mux.HandleFunc("/api/admin/series/{id}/end", EndSeriesHandler(db))

	mux.HandleFunc("/api/admin/events/create", CreateEventHandler(db))
	mux.HandleFunc("/api/admin/events/{id}/start", StartEventHandler(db))
	mux.HandleFunc("/api/admin/events/{id}/end", EndEventHandler(db))

	// Serve static files
	staticHandler := http.FileServer(http.FS(web.GetWebFS()))
	mux.Handle("/", staticHandler)

	slog.Info("starting HTTP server", "address", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("HTTP server error", "error", err)
	}
}
