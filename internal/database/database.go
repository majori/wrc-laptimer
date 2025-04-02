package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/majori/wrc-laptimer/pkg/telemetry"
	"github.com/marcboeker/go-duckdb/v2"
)

type Database struct {
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

	_, err = db.Exec(dbSchema)
	if err != nil {
		return nil, err
	}

	appender, err := duckdb.NewAppenderFromConn(dbConnection, "", "telemetry")
	if err != nil {
		return nil, fmt.Errorf("could not create new appender for telemetry: %w", err)
	}

	return &Database{
		db:       db,
		appender: appender,
	}, nil
}

func (d *Database) Close() {
	if d.appender != nil {
		d.appender.Close()
	}
	if d.db != nil {
		d.db.Close()
	}
}

func (d *Database) AppendTelemetry(t *telemetry.TelemetrySessionUpdate) error {
	return d.appender.AppendRow(
		nil,
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

func (d *Database) FlushTelemetry() {
	d.appender.Flush()
}
