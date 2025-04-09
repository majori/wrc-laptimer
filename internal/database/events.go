package database

import (
	"database/sql"
	"fmt"
)

type RaceEvent struct {
	ID             int    `db:"id"`
	RaceSeriesID   int    `db:"race_series_id"`
	LocationID     uint16 `db:"location_id"`
	Name           string `db:"name"`
	CreatedAt      string `db:"created_at"`
	VehicleID      uint16 `db:"vehicle_id"`
	VehicleClassID uint16 `db:"vehicle_class_id"`
	Active         bool   `db:"active"`
	StartedAt      string `db:"started_at"`
	EndedAt        string `db:"ended_at"`
}

var activeEventID sql.NullInt64 // Zero means no active event

func (d *Database) GetActiveEventID() (sql.NullInt64, error) {
	if activeEventID.Valid {
		return activeEventID, nil
	}

	var eventID sql.NullInt64
	query := "SELECT id FROM race_events WHERE active = TRUE ORDER BY started_at DESC LIMIT 1"
	err := d.db.QueryRow(query).Scan(&eventID)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.NullInt64{}, nil // No active event found
		}
		return sql.NullInt64{}, fmt.Errorf("failed to fetch active event ID: %w", err)
	}

	activeEventID = eventID
	return activeEventID, nil
}
func (d *Database) CreateEvent(name string, seriesID int, locationID sql.NullInt16, vehicleID sql.NullInt16, vehicleClassID sql.NullInt16) error {
	query := `
		INSERT INTO race_events (name, race_series_id, location_id, vehicle_id, vehicle_class_id, active, created_at)
		VALUES (?, ?, ?, ?, ?, FALSE, CURRENT_TIMESTAMP)
	`
	_, err := d.db.Exec(query, name, seriesID, locationID, vehicleID, vehicleClassID)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

func (d *Database) StartEvent(id int) error {
	query := `
			UPDATE race_events
			SET active = TRUE, started_at = CURRENT_TIMESTAMP
			WHERE id = ? AND active = FALSE
		`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to start event: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no event was started, possibly the event is already active or does not exist")
	}

	activeEventID = sql.NullInt64{Int64: int64(id), Valid: true}
	return nil
}

func (d *Database) EndEvent(id int) error {
	query := `
		UPDATE race_events
		SET active = FALSE, ended_at = CURRENT_TIMESTAMP
		WHERE id = ? AND active = TRUE
	`
	result, err := d.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to end event: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no event was ended, possibly the event is already inactive or does not exist")
	}

	// TODO Calculate points when event ends
	activeEventID = sql.NullInt64{}
	return nil
}

func (d *Database) GetSeriesEvents(seriesID int) ([]RaceEvent, error) {
	query := `
		SELECT id, race_series_id, location_id, name, created_at, vehicle_id, vehicle_class_id, active, started_at, ended_at
		FROM race_events
		WHERE race_series_id = ?
		ORDER BY created_at DESC
	`
	rows, err := d.db.Query(query, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch series events: %w", err)
	}
	defer rows.Close()

	var events []RaceEvent
	for rows.Next() {
		var event RaceEvent
		err := rows.Scan(
			&event.ID,
			&event.RaceSeriesID,
			&event.LocationID,
			&event.Name,
			&event.CreatedAt,
			&event.VehicleID,
			&event.VehicleClassID,
			&event.Active,
			&event.StartedAt,
			&event.EndedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return events, nil
}

func (d *Database) GetEvent(id int) (*RaceEvent, error) {
	query := `
		SELECT id, race_series_id, location_id, name, created_at, vehicle_id, vehicle_class_id, active, started_at, ended_at
		FROM race_events
		WHERE id = ?
	`
	var event RaceEvent
	err := d.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.RaceSeriesID,
		&event.LocationID,
		&event.Name,
		&event.CreatedAt,
		&event.VehicleID,
		&event.VehicleClassID,
		&event.Active,
		&event.StartedAt,
		&event.EndedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No event found
		}
		return nil, fmt.Errorf("failed to fetch event: %w", err)
	}

	return &event, nil
}
