package database

import (
	"database/sql"
	"fmt"
)

type RaceSerie struct {
	ID             int           `db:"id"`
	Name           string        `db:"name"`
	VehicleClassID sql.NullInt16 `db:"vehicle_class_id"`
	Active         bool          `db:"active"`
	CreatedAt      sql.NullTime  `db:"created_at"`
	StartedAt      sql.NullTime  `db:"started_at"`
	EndedAt        sql.NullTime  `db:"ended_at"`
}

// Cache the activeSeriedID for quicker usage
var activeSeriesID sql.NullInt32

func (d *Database) GetActiveSeriesID() (sql.NullInt32, error) {
	if !activeSeriesID.Valid {
		var id int
		err := d.db.QueryRowContext(d.ctx, `
			SELECT id
			FROM race_series
			WHERE active = true
			ORDER BY started_at DESC
			LIMIT 1;
		`).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				return sql.NullInt32{}, nil // No active series
			}
			return sql.NullInt32{}, fmt.Errorf("could not get active series: %w", err)
		}
		activeSeriesID = sql.NullInt32{Int32: int32(id), Valid: true}
	}
	return activeSeriesID, nil
}

func (d *Database) CreateSeries(name string, vehicleClassID sql.NullInt16) (int, error) {
	var id int
	err := d.db.QueryRowContext(d.ctx, `
		INSERT INTO race_series (name, vehicle_class_id)
		VALUES (?, ?)
		RETURNING id;
	`, name, vehicleClassID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("could not create series: %w", err)
	}
	return id, nil
}

func (d *Database) StartSeries(id int) error {
	_, err := d.db.ExecContext(d.ctx, `
		UPDATE race_series
		SET active = true,
			started_at = CURRENT_TIMESTAMP
		WHERE id = ?;
	`, id)
	if err != nil {
		return fmt.Errorf("could not start series: %w", err)
	}
	activeSeriesID = sql.NullInt64{Int64: int64(id), Valid: true}
	return nil
}

func (d *Database) EndSeries(id int) error {
	_, err := d.db.ExecContext(d.ctx, `
		UPDATE race_series
		SET active = false,
			ended_at = CURRENT_TIMESTAMP
		WHERE id = ?;
	`, id)
	if err != nil {
		return fmt.Errorf("could not end series: %w", err)
	}
	// TODO Loop through all series Events and close those
	activeSeriesID = sql.NullInt64{}
	return nil
}

func (d *Database) GetSeries(id int) (*RaceSerie, error) {
	var series RaceSerie
	err := d.db.QueryRowContext(d.ctx, `
		SELECT id, name, vehicle_class_id, active, created_at, started_at, ended_at
		FROM race_series
		WHERE id = ?;
	`, id).Scan(
		&series.ID,
		&series.Name,
		&series.VehicleClassID,
		&series.Active,
		&series.CreatedAt,
		&series.StartedAt,
		&series.EndedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No series found
		}
		return nil, fmt.Errorf("could not get series: %w", err)
	}
	return &series, nil
}

func (d *Database) GetAllSeries() ([]*RaceSerie, error) {
	var series []*RaceSerie
	rows, err := d.db.QueryContext(d.ctx, `
		SELECT id, name, vehicle_class_id, active, created_at, started_at, ended_at
		FROM race_series
	`)
	if err != nil {
		return nil, fmt.Errorf("could not get series: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("error closing rows: %v\n", err)
		}
	}()

	for rows.Next() {
		var s RaceSerie
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.VehicleClassID,
			&s.Active,
			&s.CreatedAt,
			&s.StartedAt,
			&s.EndedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan series: %w", err)
		}
		series = append(series, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return series, nil
}
