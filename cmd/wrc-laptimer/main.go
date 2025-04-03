package main

import (
	"context"

	"log"
	"time"

	env "github.com/caarlos0/env/v6"
	"github.com/majori/wrc-laptimer/internal/database"
	"github.com/majori/wrc-laptimer/internal/nfc"
	"github.com/majori/wrc-laptimer/pkg/telemetry"
)

type Config struct {
	ListenUDP  string `env:"LISTEN_UDP" envDefault:"0.0.0.0:20777"`
	DisableNFC bool   `env:"DISABLE_NFC"`
}

var (
	config Config
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	config = Config{}
	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.NewDatabase(ctx, "wrc.db?access_mode=READ_WRITE")
	if err != nil {
		log.Fatalf("could not open database: %s", err.Error())
	}
	defer db.Close()

	packetCh := make(chan any, 64)
	go func() {
		for {
			if err := telemetry.StartUDPReceiver(ctx, config.ListenUDP, packetCh); err != nil {
				log.Print(err)

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
				log.Printf("could not start NFC reader: %s", err.Error())
			}
		}()
	} else {
		log.Println("NFC reader disabled")
	}

	go db.ListenForUserLogins(cardEvents)

	for pkt := range packetCh {
		switch pkt := pkt.(type) {
		case *telemetry.TelemetrySessionStart:
			log.Println("Session Start")
			err = db.FlushTelemetry()
			if err != nil {
				log.Printf("could not save telemetry: %s", err.Error())
			}

			err := db.SaveSession(pkt)
			if err != nil {
				log.Printf("could not save session: %s", err.Error())
			}
			log.Println("Session saved")

		case *telemetry.TelemetrySessionUpdate:
			// log.Println("Session Update")
			err = db.AppendTelemetry(pkt)
			if err != nil {
				log.Printf("could not create new appender for telemetry: %s", err.Error())
			}

		case *telemetry.TelemetrySessionEnd:
			log.Println("Session End")
			err = db.FlushTelemetry()
			if err != nil {
				log.Printf("could not save telemetry: %s", err.Error())
			}

			err := db.FinalizeSession(pkt)
			if err != nil {
				log.Printf("could not finalize session: %s", err.Error())
			}
			log.Println("Session finalized")

		case *telemetry.TelemetrySessionPause:
			continue
		case *telemetry.TelemetrySessionResume:
			continue
		default:
			log.Printf("Unknown packet type: %T", pkt)
		}
	}
}
