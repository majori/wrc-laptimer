package main

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"net"
	"time"

	env "github.com/caarlos0/env/v6"
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

	db, err := sql.Open("duckdb", "telemetry.db?access_mode=READ_WRITE")
	if err != nil {
		log.Fatal(err)
	}
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

	saveTelemetryStreamToDB(db, ch)
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

func saveTelemetryStreamToDB(db *sql.DB, ch <-chan *codemasters.PacketEASportsWRC) {
	stmt, err := db.Prepare(`INSERT INTO telemetry (
		packet_uid,
		game_delta_time,
		game_frame_count,
		game_total_time,
		shiftlights_fraction,
		shiftlights_rpm_end,
		shiftlights_rpm_start,
		shiftlights_rpm_valid,
		stage_current_distance,
		stage_current_time,
		stage_length,
		vehicle_acceleration_x,
		vehicle_acceleration_y,
		vehicle_acceleration_z,
		vehicle_brake_temperature_bl,
		vehicle_brake_temperature_br,
		vehicle_brake_temperature_fl,
		vehicle_brake_temperature_fr,
		vehicle_brake,
		vehicle_clutch,
		vehicle_cp_forward_speed_bl,
		vehicle_cp_forward_speed_br,
		vehicle_cp_forward_speed_fl,
		vehicle_cp_forward_speed_fr,
		vehicle_engine_rpm_current,
		vehicle_engine_rpm_idle,
		vehicle_engine_rpm_max,
		vehicle_forward_direction_x,
		vehicle_forward_direction_y,
		vehicle_forward_direction_z,
		vehicle_gear_index,
		vehicle_gear_index_neutral,
		vehicle_gear_index_reverse,
		vehicle_gear_maximum,
		vehicle_handbrake,
		vehicle_hub_position_bl,
		vehicle_hub_position_br,
		vehicle_hub_position_fl,
		vehicle_hub_position_fr,
		vehicle_hub_velocity_bl,
		vehicle_hub_velocity_br,
		vehicle_hub_velocity_fl,
		vehicle_hub_velocity_fr,
		vehicle_left_direction_x,
		vehicle_left_direction_y,
		vehicle_left_direction_z,
		vehicle_position_x,
		vehicle_position_y,
		vehicle_position_z,
		vehicle_speed,
		vehicle_steering,
		vehicle_throttle,
		vehicle_transmission_speed,
		vehicle_up_direction_x,
		vehicle_up_direction_y,
		vehicle_up_direction_z,
		vehicle_velocity_x,
		vehicle_velocity_y,
		vehicle_velocity_z
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for pkt := range ch {
		log.Println(pkt.PacketUid)
		// Save packets to the database asynchronously
		go func() {
			_, err = stmt.Exec(
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
		}()
	}
}
