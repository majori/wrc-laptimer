package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"log"
	"net"
	"time"

	env "github.com/caarlos0/env/v6"
	"github.com/marcboeker/go-duckdb/v2"
	_ "github.com/marcboeker/go-duckdb/v2"
	"github.com/nobonobo/obs-codemasters-telemetry/codemasters"
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

	connector, err := duckdb.NewConnector("telemetry.db?access_mode=READ_WRITE", nil)
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

	ch := make(chan *codemasters.PacketEASportsWRC, 64)
	go func() {
		for {
			if err := udpReceiver(ctx, ch); err != nil {
				log.Print(err)

				// Retry receiving UDP packets after 5 seconds
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}
	}()

	saveTelemetryStreamToDB(dbConnection, ch)
}

func udpReceiver(ctx context.Context, ch chan<- *codemasters.PacketEASportsWRC) error {
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
		b := make([]byte, codemasters.PacketEASportsWRCLength)
		for {
			_, _, err := conn.ReadFrom(b)
			if err != nil {
				done <- err
			}

			pkt := &codemasters.PacketEASportsWRC{}
			if err := pkt.UnmarshalBinary(b); err != nil {
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

func saveTelemetryStreamToDB(con driver.Conn, ch <-chan *codemasters.PacketEASportsWRC) {
	appender, err := duckdb.NewAppenderFromConn(con, "", "telemetry")
	if err != nil {
		log.Fatalf("could not create new appender for telemetry: %s", err.Error())
	}
	defer appender.Close()

	for pkt := range ch {
		err := appender.AppendRow(
			pkt.PacketUid,
			pkt.GameDeltaTime,
			pkt.GameFrameCount,
			pkt.GameTotalTime,
			pkt.ShiftlightsFraction,
			pkt.ShiftlightsRpmEnd,
			pkt.ShiftlightsRpmStart,
			pkt.ShiftlightsRpmValid,
			pkt.StageCurrentDistance,
			pkt.StageCurrentTime,
			pkt.StageLength,
			pkt.VehicleAccelerationX,
			pkt.VehicleAccelerationY,
			pkt.VehicleAccelerationZ,
			pkt.VehicleBrake,
			pkt.VehicleBrakeTemperatureBl,
			pkt.VehicleBrakeTemperatureBr,
			pkt.VehicleBrakeTemperatureFl,
			pkt.VehicleBrakeTemperatureFr,
			pkt.VehicleClutch,
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
			pkt.VehicleUpDirectionX,
			pkt.VehicleUpDirectionY,
			pkt.VehicleUpDirectionZ,
			pkt.VehicleVelocityX,
			pkt.VehicleVelocityY,
			pkt.VehicleVelocityZ,
		)
		if err != nil {
			log.Printf("Error inserting telemetry data: %v", err)
		}
	}
}
