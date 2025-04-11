package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Result struct {
	ID          int          `db:"id"`
	UserID      string       `db:"user_id"`
	RaceEventID int          `db:"race_event_id"`
	CreatedAt   sql.NullTime `db:"created_at"`
	Points      int          `db:"points"`
	HCMode      bool         `db:"hc_mode"`
	Position    int          `db:"position"`
	ResultTime  float32      `db:"result_time"`
}

type SessionResult struct {
	UserID           string  `db:"user_id"`
	EventID          int     `db:"race_event_id"`
	StageTotalResult float32 `db:"stage_total_result"`
}

/* type EventResult struct {
	eventID      int
	bestResults  []Result
	firstResults []Result
} */

func (d *Database) QuerySessionResults(eventID int, query string) ([]SessionResult, error) {

	rows, err := d.db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}()

	var results []SessionResult
	for rows.Next() {
		var result SessionResult
		if err := rows.Scan(
			&result.UserID,
			&result.EventID,
			&result.StageTotalResult,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return results, nil
}

func (d *Database) GetBestSessionsByEventID(eventID int) ([]SessionResult, error) {
	// SQL query to get the first session for each user in the event, sorted by time
	query := `
		SELECT 
			user_id, race_event_id, MIN(stage_result_time + stage_result_time_penalty) as stage_total_result,
		FROM 
			sessions
		WHERE 
			stage_result_status = 1 AND
			race_event_id = ?
		GROUP BY
			user_id, race_event_id
		ORDER BY 
			stage_total_result ASC;
		`

	return d.QuerySessionResults(eventID, query)
}

func (d *Database) GetFirstSessionsByEventID(eventID int) ([]SessionResult, error) {
	// SQL query to get the first session for each user in the event, sorted by time
	query := `
		SELECT 
			user_id, race_event_id, stage_result_time + stage_result_time_penalty as stage_total_result,
		FROM 
			sessions
		WHERE 
			id IN (
				SELECT MIN(id) FROM sessions
				WHERE
					stage_result_status = 1 
					AND race_event_id = ?
				GROUP BY user_id
			)
		ORDER BY 
			stage_total_result ASC;
		`
	return d.QuerySessionResults(eventID, query)
}

func calculatePoints(sessions []SessionResult, pointScale []int, HCMode bool) []Result {
	var results []Result
	for i, session := range sessions {
		pointsValue := 0
		if i < len(pointScale) {
			pointsValue = pointScale[i]
		}
		result := Result{
			UserID:      session.UserID,
			RaceEventID: session.EventID,
			ResultTime:  session.StageTotalResult,
			Points:      pointsValue,
			Position:    i + 1,
			HCMode:      HCMode,
		}
		results = append(results, result)
	}
	return results
}

// Store Event results
func (d *Database) StoreResults(results []Result) error {
	query := `
		INSERT INTO results (user_id, race_event_id, created_at, points, hc_mode, position, result_time)
		VALUES (?, ?, NOW(), ?, ?, ?, ?)
	`
	for _, result := range results {
		_, err := d.db.Exec(query, result.UserID, result.RaceEventID, result.Points, result.HCMode, result.Position, result.ResultTime)
		if err != nil {
			return fmt.Errorf("failed to store result for user %s: %w", result.UserID, err)
		}
	}
	return nil
}

func (d *Database) parsePointScale(pointScaleString sql.NullString) ([]int, error) {
	if !pointScaleString.Valid {
		return []int{}, nil
	}

	pointScaleParts := strings.Split(pointScaleString.String, "-")
	var pointScale []int
	for _, part := range pointScaleParts {
		point, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid point scale value: %w", err)
		}
		pointScale = append(pointScale, point)
	}
	return pointScale, nil
}

func (d *Database) CalculateAndStoreEventResults(eventID int) error {

	var event, err = d.GetEvent(eventID)
	if err != nil {
		return fmt.Errorf("failed to get event: %w", err)
	}

	pointScale, err := d.parsePointScale(event.PointScale)
	if err != nil {
		return fmt.Errorf("failed to parse point scale: %w", err)
	}

	bestSessions, err := d.GetBestSessionsByEventID(eventID)
	if err != nil {
		return fmt.Errorf("failed to get best sessions: %w", err)
	}
	firstSessions, err := d.GetFirstSessionsByEventID(eventID)
	if err != nil {
		return fmt.Errorf("failed to get first sessions: %w", err)
	}

	// Calculate both best and first session points
	var bestSessionPoints = calculatePoints(bestSessions, pointScale, false)
	var firstSessionPoints = calculatePoints((firstSessions), pointScale, true)

	if err := d.StoreResults(bestSessionPoints); err != nil {
		return fmt.Errorf("failed to store best session points: %w", err)
	}

	if err := d.StoreResults(firstSessionPoints); err != nil {
		return fmt.Errorf("failed to store first session points: %w", err)
	}
	return nil
}

func (d *Database) GetPointsByEventID(eventID int, HCMode bool) ([]Result, error) {
	query := `
		SELECT 
			id, user_id, race_event_id, created_at, points, hc_mode, position, result_time
		FROM 
			points
		WHERE 
			race_event_id = ? AND hc_mode = ?
		ORDER BY 
			position ASC
	`
	rows, err := d.db.Query(query, eventID, HCMode)
	if err != nil {
		return nil, fmt.Errorf("failed to query points: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}()

	var points []Result
	for rows.Next() {
		var point Result
		if err := rows.Scan(
			&point.ID,
			&point.UserID,
			&point.RaceEventID,
			&point.CreatedAt,
			&point.Points,
			&point.HCMode,
			&point.Position,
			&point.ResultTime,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		points = append(points, point)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return points, nil
}
