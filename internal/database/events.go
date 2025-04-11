package database

import (
	"database/sql"
	"fmt"
	"time"
)

type RaceEvent struct {
	ID             int            `db:"id"`
	Name           string         `db:"name"`
	RaceSeriesID   sql.NullInt32  `db:"race_series_id"`
	LocationID     sql.NullInt16  `db:"location_id"`
	RouteID        sql.NullInt16  `db:"route_id"`
	VehicleClassID sql.NullInt16  `db:"vehicle_class_id"`
	Active         bool           `db:"active"`
	PointScale     sql.NullString `db:"point_scale"`
	CreatedAt      time.Time      `db:"created_at"`
	StartedAt      sql.NullTime   `db:"started_at"`
	EndedAt        sql.NullTime   `db:"ended_at"`
}

var activeEventID sql.NullInt32        // There may be case where no event is running. Using NullInt32 instead of checking 0 value
var activeVehicleClassID sql.NullInt16 // Vehicle class may not be set for all events, so using NullInt16

// TODO Add restrictions also based on location and route
// var activeSeriesLocation sql.NullInt16
// var activeSeriesRoute sql.NullInt16

func (d *Database) GetActiveEventID(eventVehicleClassID int16) (sql.NullInt32, error) {
	// Checks also if the active event has vehicle restrictions
	if activeEventID.Valid && (!activeVehicleClassID.Valid || (activeVehicleClassID.Int16 == eventVehicleClassID)) {
		return activeEventID, nil
	}

	query := `
		SELECT id, vehicle_class_id
		FROM race_events
		WHERE active = TRUE
		ORDER BY started_at DESC LIMIT 1
	`
	var eventID sql.NullInt32
	var vehicleClassID sql.NullInt16
	err := d.db.QueryRow(query).Scan(&eventID, &vehicleClassID)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.NullInt32{}, nil // No active event found
		}
		return sql.NullInt32{}, fmt.Errorf("failed to fetch active event ID: %w", err)
	}

	activeEventID = eventID
	activeVehicleClassID = vehicleClassID
	return activeEventID, nil
}
func (d *Database) CreateEvent(name string, seriesID sql.NullInt32, locationID sql.NullInt16, routeID sql.NullInt16, vehicleClassID sql.NullInt16) (int, error) {
	query := `
		INSERT INTO race_events (name, race_series_id, location_id, route_id, vehicle_class_id, active, created_at)
		VALUES (?, ?, ?, ?, ?, FALSE, CURRENT_TIMESTAMP)
		RETURNING id
	`
	var eventID int
	err := d.db.QueryRow(query, name, seriesID, locationID, routeID, vehicleClassID).Scan(&eventID)
	if err != nil {
		return 0, fmt.Errorf("failed to create event: %w", err)
	}
	return eventID, nil
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

	// Get the vehicle class ID for the active event if exists
	query = "SELECT vehicle_class_id FROM race_events WHERE id = ?"
	err = d.db.QueryRow(query, id).Scan(&activeVehicleClassID)
	if err != nil {
		return fmt.Errorf("failed to fetch vehicle class ID: %w", err)
	}
	activeEventID = sql.NullInt32{Int32: int32(id), Valid: true}
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

	err = d.CalculateAndStoreEventResults(id)
	if err != nil {
		return fmt.Errorf("failed to calculate and store event results: %w", err)
	}

	// TODO Calculate points when event ends
	activeEventID = sql.NullInt32{}
	return nil
}

func (d *Database) GetSeriesEvents(seriesID int) ([]RaceEvent, error) {
	query := `
		SELECT
			id,
			name,
			race_series_id,
			location_id,
			route_id,
			vehicle_class_id,
			point_scale,
			active,
			created_at,
			started_at,
			ended_at
		FROM race_events
		WHERE race_series_id = ?
		ORDER BY created_at DESC
	`
	rows, err := d.db.Query(query, seriesID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch series events: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}()

	var events []RaceEvent
	for rows.Next() {
		var event RaceEvent
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.RaceSeriesID,
			&event.LocationID,
			&event.RouteID,
			&event.VehicleClassID,
			&event.PointScale,
			&event.Active,
			&event.CreatedAt,
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
		SELECT
			id,
			name,
			race_series_id,
			location_id,
			route_id,
			vehicle_class_id,
			point_scale,
			active,
			created_at,
			started_at,
			ended_at
		FROM race_events
		WHERE id = ?
	`
	var event RaceEvent
	err := d.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.Name,
		&event.RaceSeriesID,
		&event.LocationID,
		&event.RouteID,
		&event.VehicleClassID,
		&event.PointScale,
		&event.Active,
		&event.CreatedAt,
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
