package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	env "github.com/caarlos0/env/v6"
	"github.com/majori/wrc-laptimer/internal/database"
	"github.com/majori/wrc-laptimer/internal/events"
	"github.com/majori/wrc-laptimer/internal/http"
	"github.com/majori/wrc-laptimer/internal/nfc"
	"github.com/majori/wrc-laptimer/pkg/telemetry"
)

type Config struct {
	ListenUDP  string `env:"LISTEN_UDP" envDefault:"127.0.0.1:20777"`
	ListenHTTP string `env:"LISTEN_HTTP" envDefault:"127.0.0.1:8080"`
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

	packetCh := make(chan telemetry.TelemetryPacket, 64)
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
	go func() {
		err := nfc.ListenForCardEvents(ctx, cardEvents)
		if err != nil {
			slog.Error("could not start NFC reader", "error", err)
		}
	}()

	go db.ListenForUserLogins(cardEvents)

	go http.StartHTTPServer(db, config.ListenHTTP)

	go events.ProcessTelemetryEvents(ctx, db, packetCh)

	<-ctx.Done()
}
