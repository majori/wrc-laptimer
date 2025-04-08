package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	env "github.com/caarlos0/env/v6"
	"github.com/majori/wrc-laptimer/internal/database"
	"github.com/majori/wrc-laptimer/internal/nfc"
	"github.com/majori/wrc-laptimer/pkg/telemetry"
	"github.com/majori/wrc-laptimer/web"
)

type Config struct {
	ListenUDP  string `env:"LISTEN_UDP" envDefault:"0.0.0.0:20777"`
	ListenHTTP string `env:"LISTEN_HTTP" envDefault:"0.0.0.0:8080"`
	DisableNFC bool   `env:"DISABLE_NFC"`
}

var (
	config Config
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	config = Config{}
	if err := env.Parse(&config); err != nil {
		slog.Error("failed to parse config", "error", err)
		os.Exit(1)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		slog.Debug("received shutdown signal, canceling context")
		cancel()
	}()

	db, err := database.NewDatabase(ctx, "wrc.db?access_mode=READ_WRITE")
	if err != nil {
		slog.Error("could not open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	packetCh := make(chan any, 64)
	go func() {
		for {
			if err := telemetry.StartUDPReceiver(ctx, config.ListenUDP, packetCh); err != nil {
				slog.Error("UDP receiver error", "error", err)

				// Retry receiving UDP packets after 5 seconds
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}
	}()

	cardEvents := make(chan string, 1)
	if !config.DisableNFC {
		go func() {
			err := nfc.ListenForCardEvents(ctx, cardEvents)
			if err != nil {
				slog.Error("could not start NFC reader", "error", err)
			}
		}()
	} else {
		slog.Info("NFC reader disabled")
	}

	go db.ListenForUserLogins(cardEvents)

	// Setup HTTP server
	go func() {
		http.Handle("/", http.FileServer(http.FS(web.GetWebFS())))

		// Setup HTTP server with /api/query endpoint
		http.HandleFunc("/api/query", func(w http.ResponseWriter, r *http.Request) {
			// Parse the query parameter
			query := r.URL.Query().Get("query")
			if query == "" {
				http.Error(w, "missing query parameter", http.StatusBadRequest)
				return
			}

			// Execute the query on the database
			rows, err := db.UnsafeQuery(query)
			if err != nil {
				slog.Error("failed to execute query", "error", err)
				http.Error(w, "failed to execute query", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			// Serialize the rows into JSON
			columns, err := rows.Columns()
			if err != nil {
				slog.Error("failed to get columns", "error", err)
				http.Error(w, "failed to get columns", http.StatusInternalServerError)
				return
			}

			results := []map[string]interface{}{}
			rowCount := 0
			for rows.Next() {
				// Create a slice of interface{} to hold column values
				values := make([]interface{}, len(columns))
				valuePtrs := make([]interface{}, len(columns))
				for i := range values {
					valuePtrs[i] = &values[i]
				}

				// Scan the row into the value pointers
				if err := rows.Scan(valuePtrs...); err != nil {
					slog.Error("failed to scan row", "error", err)
					http.Error(w, "failed to scan row", http.StatusInternalServerError)
					return
				}

				// Create a map for the row
				row := make(map[string]interface{})
				for i, col := range columns {
					row[col] = values[i]
				}
				results = append(results, row)
				rowCount++
			}

			// Write the JSON response
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(results); err != nil {
				slog.Error("failed to encode response", "error", err)
				http.Error(w, "failed to encode response", http.StatusInternalServerError)
			}
		})

		// Start the HTTP server
		slog.Info("starting HTTP server", "address", config.ListenHTTP)
		if err := http.ListenAndServe(config.ListenHTTP, nil); err != nil {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	for {
		select {
		case pkt := <-packetCh:
			switch pkt := pkt.(type) {
			case *telemetry.TelemetrySessionStart:
				err := db.FlushTelemetry()
				if err != nil {
					slog.Error("could not save telemetry", "error", err)
				}

				err = db.StartSession(pkt)
				if err != nil {
					slog.Error("could not save session", "error", err)
				}
				slog.Info("session started")

			case *telemetry.TelemetrySessionUpdate:
				err := db.AppendTelemetry(pkt)
				if err != nil {
					slog.Error("could not create new appender for telemetry", "error", err)
				}

			case *telemetry.TelemetrySessionEnd:
				err := db.FlushTelemetry()
				if err != nil {
					slog.Error("could not save telemetry", "error", err)
				}

				err = db.EndSession(pkt)
				if err != nil {
					slog.Error("could not end session", "error", err)
				}
				slog.Info("session ended")

			case *telemetry.TelemetrySessionPause:
				continue
			case *telemetry.TelemetrySessionResume:
				continue
			default:
				slog.Warn("unknown packet type", "type", pkt)
			}
		case <-ctx.Done():
			slog.Info("exiting...")
			return
		}
	}
}
