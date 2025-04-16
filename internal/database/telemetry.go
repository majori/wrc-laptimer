package database

import (
	"github.com/majori/wrc-laptimer/pkg/telemetry"
)

func (d *Database) AppendTelemetry(t *telemetry.TelemetrySessionUpdate) error {
	// Ignore telemetry if no active session
	sessionID := d.GetActiveSessionID()
	if sessionID == 0 {
		return nil
	}

	return d.appender.AppendRow(
		sessionID,
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

func (d *Database) ClearTelemetry() error {
	_, err := d.exec("DELETE FROM telemetry")
	if err != nil {
		return err
	}
	return nil
}
