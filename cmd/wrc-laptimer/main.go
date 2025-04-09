package main

import (
	"context"
	"io"
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
	ListenUDP  string `env:"LISTEN_UDP" envDefault:"127.0.0.1:20777"`
	ListenHTTP string `env:"LISTEN_HTTP" envDefault:"127.0.0.1:8080"`
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
		mux := http.NewServeMux()

		// Add query endpoint
		mux.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
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

		// Serve static files
		staticHandler := http.FileServer(http.FS(web.GetWebFS()))
		mux.Handle("/", staticHandler)

		slog.Info("starting HTTP server", "address", config.ListenHTTP)
		if err := http.ListenAndServe(config.ListenHTTP, mux); err != nil {
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
