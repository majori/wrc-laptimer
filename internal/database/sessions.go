package database

import (
	"github.com/majori/wrc-laptimer/pkg/telemetry"
)

var activeSessionID int

func (d *Database) GetActiveSessionID() int {
	return activeSessionID
}

func (d *Database) setActiveSessionID(id int) {
	activeSessionID = id
}

func (d *Database) StartSession(pkt *telemetry.TelemetrySessionStart) error {
	session := d.db.QueryRowContext(d.ctx, `
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
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id;
	`, pkt.GameMode, pkt.LocationID, pkt.RouteID, pkt.StageLength, pkt.StageShakedown, pkt.VehicleClassID, pkt.VehicleID, pkt.VehicleManufacturerID)

	var sessionID int
	err := session.Scan(&sessionID)
	if err != nil {
		return err
	}

	d.setActiveSessionID(sessionID)

	return nil
}

func (d *Database) EndSession(pkt *telemetry.TelemetrySessionEnd) error {
	if activeSessionID == 0 {
		return nil
	}

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
		WHERE id = ?;
	`, userID, pkt.StageResultStatus, pkt.StageResultTime, pkt.StageResultTimePenalty, activeSessionID)
	if err != nil {
		return err
	}

	// Reset active session
	d.setActiveSessionID(0)

	return nil
}
