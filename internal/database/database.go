package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/majori/wrc-laptimer/pkg/telemetry"
	"github.com/marcboeker/go-duckdb/v2"
)

type Database struct {
	ctx      context.Context
	db       *sql.DB
	appender *duckdb.Appender
}

//go:embed db.sql
var dbSchema string

func NewDatabase(ctx context.Context, path string) (*Database, error) {
	connector, err := duckdb.NewConnector(path, nil)
	if err != nil {
		return nil, err
	}

	dbConnection, err := connector.Connect(ctx)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(connector)

	_, err = db.ExecContext(ctx, dbSchema)
	if err != nil {
		return nil, err
	}

	appender, err := duckdb.NewAppenderFromConn(dbConnection, "", "telemetry")
	if err != nil {
		return nil, fmt.Errorf("could not create new appender for telemetry: %w", err)
	}

	return &Database{
		ctx:      ctx,
		db:       db,
		appender: appender,
	}, nil
}

func (d *Database) Close() {
	if d.appender != nil {
		//nolint:errcheck
		d.appender.Close()
	}
	if d.db != nil {
		//nolint:errcheck
		d.db.Close()
	}
}

func (d *Database) AppendTelemetry(t *telemetry.TelemetrySessionUpdate) error {
	return d.appender.AppendRow(
		time.Now(),
		t.StageCurrentDistance,
		t.StageCurrentTime,
		t.StagePreviousSplitTime,
		t.StageProgress,
		t.VehicleAccelerationX,
		t.VehicleAccelerationY,
		t.VehicleAccelerationZ,
		t.VehicleBrake,
		t.VehicleBrakeTemperatureBl,
		t.VehicleBrakeTemperatureBr,
		t.VehicleBrakeTemperatureFl,
		t.VehicleBrakeTemperatureFr,
		t.VehicleClutch,
		t.VehicleClusterAbs,
		t.VehicleCpForwardSpeedBl,
		t.VehicleCpForwardSpeedBr,
		t.VehicleCpForwardSpeedFl,
		t.VehicleCpForwardSpeedFr,
		t.VehicleEngineRpmCurrent,
		t.VehicleEngineRpmIdle,
		t.VehicleEngineRpmMax,
		t.VehicleForwardDirectionX,
		t.VehicleForwardDirectionY,
		t.VehicleForwardDirectionZ,
		t.VehicleGearIndex,
		t.VehicleGearIndexNeutral,
		t.VehicleGearIndexReverse,
		t.VehicleGearMaximum,
		t.VehicleHandbrake,
		t.VehicleHubPositionBl,
		t.VehicleHubPositionBr,
		t.VehicleHubPositionFl,
		t.VehicleHubPositionFr,
		t.VehicleHubVelocityBl,
		t.VehicleHubVelocityBr,
		t.VehicleHubVelocityFl,
		t.VehicleHubVelocityFr,
		t.VehicleLeftDirectionX,
		t.VehicleLeftDirectionY,
		t.VehicleLeftDirectionZ,
		t.VehiclePositionX,
		t.VehiclePositionY,
		t.VehiclePositionZ,
		t.VehicleSpeed,
		t.VehicleSteering,
		t.VehicleThrottle,
		t.VehicleTransmissionSpeed,
		t.VehicleTyreStateBl,
		t.VehicleTyreStateBr,
		t.VehicleTyreStateFl,
		t.VehicleTyreStateFr,
		t.VehicleUpDirectionX,
		t.VehicleUpDirectionY,
		t.VehicleUpDirectionZ,
		t.VehicleVelocityX,
		t.VehicleVelocityY,
		t.VehicleVelocityZ,
	)
}

func (d *Database) FlushTelemetry() error {
	return d.appender.Flush()
}

func (d *Database) ListenForUserLogins(cardEvents <-chan string) {
	for id := range cardEvents {
		// Check if the user already exists
		var exists bool
		err := d.db.QueryRowContext(d.ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM users
				WHERE id = ?
			);
		`, id).Scan(&exists)
		if err != nil {
			log.Println("Error checking user existence:", err)
			continue
		}
		if exists {
			continue
		} else {
			// Create the user
			err = d.CreateUser(id)
			if err != nil {
				log.Println("Error creating user:", err)
				continue
			}

		}

		// Insert the user login into the database
		_, err = d.db.ExecContext(d.ctx, `
				INSERT INTO user_logins (user_id) 
				VALUES (?);
			`, id)
		if err != nil {
			log.Println("Error inserting user login:", err)
		}
	}
}

func (d *Database) CreateUser(id string) error {
	name := id[:5] // TODO: Get a better placeholder name
	_, err := d.db.ExecContext(d.ctx, `
		INSERT INTO users (id, name)
		VALUES (?, ?)
		ON CONFLICT (id) DO NOTHING
	`, id, name)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}
	log.Println("User created:", name)
	return nil
}

func (d *Database) GetActiveUserID() (sql.NullString, error) {
	// TODO: Figure out if user has logged in but is not active (too much time from last session)
	var id sql.NullString
	err := d.db.QueryRowContext(d.ctx, `
		SELECT u.id
		FROM user_logins ul
		JOIN users u ON ul.user_id = u.id
		WHERE ul.timestamp = (
			SELECT MAX(timestamp)
			FROM user_logins
		);
	`).Scan(&id)
	if err != nil {
		// No rows means no user is logged in
		if err == sql.ErrNoRows {
			return sql.NullString{}, nil
		}
		return sql.NullString{}, fmt.Errorf("could not get active user: %w", err)
	}
	return id, nil
}

func (d *Database) SaveSession(pkt *telemetry.TelemetrySessionStart) error {
	_, err := d.db.ExecContext(d.ctx, `
		INSERT INTO sessions (
			game_mode,
			location_id,
			route_id,
			stage_length,
			stage_shakedown,
			vehicle_class_id,
			vehicle_id,
			vehicle_manufacturer_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`, pkt.GameMode, pkt.LocationID, pkt.RouteID, pkt.StageLength, pkt.StageShakedown, pkt.VehicleClassID, pkt.VehicleID, pkt.VehicleManufacturerID)
	if err != nil {
		return fmt.Errorf("could not save session: %w", err)
	}
	log.Println("Session saved")
	return nil
}

func (d *Database) FinalizeSession(pkt *telemetry.TelemetrySessionEnd) error {
	userID, err := d.GetActiveUserID()
	if err != nil {
		return err
	}
	_, err = d.db.ExecContext(d.ctx, `
		UPDATE sessions
		SET user_id = ?,
			stage_result_status = ?,
			stage_result_time = ?,
			stage_result_time_penalty = ?
		WHERE id = (
			SELECT MAX(id)
			FROM sessions
		);
	`, userID, pkt.StageResultStatus, pkt.StageResultTime, pkt.StageResultTimePenalty)
	if err != nil {
		return fmt.Errorf("could not finalize session: %w", err)
	}
	log.Println("Session finalized")
	return nil
}
