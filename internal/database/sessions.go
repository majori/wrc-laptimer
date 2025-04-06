package database

import (
	"database/sql"
	"fmt"

	"github.com/majori/wrc-laptimer/pkg/telemetry"
)

var activeSessionID sql.NullInt32

func (d *Database) StartSession(pkt *telemetry.TelemetrySessionStart) error {
	result, err := d.db.ExecContext(d.ctx, `
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

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("could not get session ID: %w", err)
	}
	activeSessionID = sql.NullInt32{
		Int32: int32(id),
		Valid: true,
	}

	return nil
}

func (d *Database) EndSession(pkt *telemetry.TelemetrySessionEnd) error {
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

	// Reset active session
	activeSessionID = sql.NullInt32{}

	return nil
}
