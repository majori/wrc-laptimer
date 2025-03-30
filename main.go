package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	_ "embed"
	"encoding/base64"
	"log"
	"net"
	"time"

	env "github.com/caarlos0/env/v6"
	"github.com/marcboeker/go-duckdb/v2"
	"github.com/peterhellberg/acr122u"
)

type Config struct {
	ListenUDP string `env:"LISTEN_UDP" envDefault:"0.0.0.0:20777"`
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

	connector, err := duckdb.NewConnector("wrc.db?access_mode=READ_WRITE", nil)
	if err != nil {
		log.Fatalf("could not initialize new connector: %s", err.Error())
	}

	dbConnection, err := connector.Connect(ctx)
	if err != nil {
		log.Fatalf("could not connect: %s", err.Error())
	}

	db := sql.OpenDB(connector)
	defer db.Close()

	initDB(db)

	udpCh := make(chan any, 64)
	go func() {
		for {
			if err := udpReceiver(ctx, udpCh); err != nil {
				log.Print(err)

				// Retry receiving UDP packets after 5 seconds
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}
	}()

	go startCardReader(ctx, db)

	appender, err := duckdb.NewAppenderFromConn(dbConnection, "", "telemetry")
	if err != nil {
		log.Fatalf("could not create new appender for telemetry: %s", err.Error())
	}
	defer appender.Close()

	for pkt := range udpCh {
		switch pkt := pkt.(type) {
		case *TelemetrySessionStart:
			log.Println("Session Start")
			// TODO: Create a new session in the database
		case *TelemetrySessionUpdate:
			log.Println("Session Update")
			appender.AppendRow(
				nil,
				pkt.StageCurrentDistance,
				pkt.StageCurrentTime,
				pkt.StagePreviousSplitTime,
				pkt.StageProgress,
				pkt.VehicleAccelerationX,
				pkt.VehicleAccelerationY,
				pkt.VehicleAccelerationZ,
				pkt.VehicleBrake,
				pkt.VehicleBrakeTemperatureBl,
				pkt.VehicleBrakeTemperatureBr,
				pkt.VehicleBrakeTemperatureFl,
				pkt.VehicleBrakeTemperatureFr,
				pkt.VehicleClutch,
				pkt.VehicleClusterAbs,
				pkt.VehicleCpForwardSpeedBl,
				pkt.VehicleCpForwardSpeedBr,
				pkt.VehicleCpForwardSpeedFl,
				pkt.VehicleCpForwardSpeedFr,
				pkt.VehicleEngineRpmCurrent,
				pkt.VehicleEngineRpmIdle,
				pkt.VehicleEngineRpmMax,
				pkt.VehicleForwardDirectionX,
				pkt.VehicleForwardDirectionY,
				pkt.VehicleForwardDirectionZ,
				pkt.VehicleGearIndex,
				pkt.VehicleGearIndexNeutral,
				pkt.VehicleGearIndexReverse,
				pkt.VehicleGearMaximum,
				pkt.VehicleHandbrake,
				pkt.VehicleHubPositionBl,
				pkt.VehicleHubPositionBr,
				pkt.VehicleHubPositionFl,
				pkt.VehicleHubPositionFr,
				pkt.VehicleHubVelocityBl,
				pkt.VehicleHubVelocityBr,
				pkt.VehicleHubVelocityFl,
				pkt.VehicleHubVelocityFr,
				pkt.VehicleLeftDirectionX,
				pkt.VehicleLeftDirectionY,
				pkt.VehicleLeftDirectionZ,
				pkt.VehiclePositionX,
				pkt.VehiclePositionY,
				pkt.VehiclePositionZ,
				pkt.VehicleSpeed,
				pkt.VehicleSteering,
				pkt.VehicleThrottle,
				pkt.VehicleTransmissionSpeed,
				pkt.VehicleTyreStateBl,
				pkt.VehicleTyreStateBr,
				pkt.VehicleTyreStateFl,
				pkt.VehicleTyreStateFr,
				pkt.VehicleUpDirectionX,
				pkt.VehicleUpDirectionY,
				pkt.VehicleUpDirectionZ,
				pkt.VehicleVelocityX,
				pkt.VehicleVelocityY,
				pkt.VehicleVelocityZ,
			)
		case *TelemetrySessionPause:
			continue
		case *TelemetrySessionResume:
			continue
		case *TelemetrySessionEnd:
			log.Println("Session End")
			// TODO: End the session in the database
		default:
			log.Printf("Unknown packet type: %T", pkt)
		}
	}
}

func startCardReader(ctx context.Context, db *sql.DB) {
	readerCtx, err := acr122u.EstablishContext()
	if err != nil {
		panic(err)
	}

	log.Println("ready for smartcard events")
	readerCtx.ServeFunc(func(c acr122u.Card) {
		hasher := sha256.New()
		hasher.Write(c.UID())
		userID := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		// TODO: Add userID to the database if it doesn't exist

		_, err := db.ExecContext(ctx, "INSERT INTO user_logins (user_id) VALUES (?)", userID)
		if err != nil {
			log.Printf("Error inserting user: %v", err)
		}
	})
}

func udpReceiver(ctx context.Context, ch chan<- any) error {
	// Validate the UDP address
	host, port, err := net.SplitHostPort(config.ListenUDP)
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(host, port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("listen udp:", addr)
	defer log.Println("udp closed:", addr)

	done := make(chan error, 1)
	go func() {
		// Create a buffer large enough
		b := make([]byte, 256)
		for {
			n, _, err := conn.ReadFrom(b)
			if err != nil {
				done <- err
				continue
			}

			// Process only the bytes that were read
			pkt, err := UnmarshalBinary(b[:n])
			if err != nil {
				done <- err
				continue
			}

			ch <- pkt
		}
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
	}
	return nil
}

//go:embed db.sql
var dbSchema string

func initDB(db *sql.DB) {
	_, err := db.Exec(dbSchema)
	if err != nil {
		panic(err)
	}
}
