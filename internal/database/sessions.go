package database

import (
	"fmt"

	"github.com/majori/wrc-laptimer/pkg/telemetry"
)

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

	return nil
}
