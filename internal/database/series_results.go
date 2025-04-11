package database

import (
	"database/sql"
	"fmt"
	"sort"
)

type SeriesResult struct {
	ID           int          `db:"id"`
	UserID       string       `db:"user_id"`
	RaceSeriesID int          `db:"race_series_id"`
	CreatedAt    sql.NullTime `db:"created_at"`
	Points       int          `db:"points"`
	RaceCount    int          `db:"race_count"`
	HCMode       bool         `db:"hc_mode"`
	Position     int          `db:"position"`
	ResultTime   float32      `db:"result_time"`
}

func (d *Database) GetResultsBySeriesID(seriesID int, HCMode bool) ([]Result, error) {
	query := `
        SELECT 
            r.id,
            r.user_id,
            r.race_event_id,
            r.created_at,
            r.points,
            r.hc_mode,
            r.position,
            r.result_time
        FROM 
            results r
        JOIN 
            race_events e ON r.race_event_id = e.id
        WHERE 
            e.race_series_id = ? AND r.hc_mode = ?;
    `
	rows, err := d.db.Query(query, seriesID, HCMode)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("failed to close rows: %v\n", err)
		}
	}()

	var results []Result
	for rows.Next() {
		var result Result
		if err := rows.Scan(
			&result.ID,
			&result.UserID,
			&result.RaceEventID,
			&result.CreatedAt,
			&result.Points,
			&result.HCMode,
			&result.Position,
			&result.ResultTime,
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

func aggregateResultsByUser(seriesID int, eventResults []Result) map[string]*SeriesResult {
	seriesResults := make(map[string]*SeriesResult)
	for _, eventResult := range eventResults {
		if _, exists := seriesResults[eventResult.UserID]; !exists {
			seriesResults[eventResult.UserID] = &SeriesResult{
				UserID:       eventResult.UserID,
				RaceSeriesID: seriesID,
				Points:       0,
				ResultTime:   0,
				RaceCount:    0,
			}
		}
		// Calculate event points (e.g. Power Stage points)
		seriesResults[eventResult.UserID].Points += eventResult.Points
		seriesResults[eventResult.UserID].ResultTime += eventResult.ResultTime
		seriesResults[eventResult.UserID].RaceCount++
	}
	return seriesResults
}

func convertAndSortResults(seriesResults map[string]*SeriesResult) []SeriesResult {
	var aggregatedResults []SeriesResult
	for _, seriesResult := range seriesResults {
		aggregatedResults = append(aggregatedResults, *seriesResult)
	}

	// Sort by RaceCount (descending) and ResultTime (ascending)
	sort.Slice(aggregatedResults, func(i, j int) bool {
		if aggregatedResults[i].RaceCount == aggregatedResults[j].RaceCount {
			return aggregatedResults[i].ResultTime < aggregatedResults[j].ResultTime
		}
		return aggregatedResults[i].RaceCount > aggregatedResults[j].RaceCount
	})

	// Assign positions
	for i := range aggregatedResults {
		aggregatedResults[i].Position = i + 1
	}

	return aggregatedResults
}

func applyPointScale(aggregatedResults []SeriesResult, pointScale []int) []SeriesResult {
	for i := range aggregatedResults {
		if aggregatedResults[i].Position <= len(pointScale) {
			aggregatedResults[i].Points += pointScale[aggregatedResults[i].Position-1]
		}
	}
	return aggregatedResults
}

func reorderByPoints(aggregatedResults []SeriesResult) []SeriesResult {
	sort.Slice(aggregatedResults, func(i, j int) bool {
		if aggregatedResults[i].Points == aggregatedResults[j].Points {
			return aggregatedResults[i].ResultTime < aggregatedResults[j].ResultTime
		}
		return aggregatedResults[i].Points > aggregatedResults[j].Points
	})

	// Reassign positions
	for i := range aggregatedResults {
		aggregatedResults[i].Position = i + 1
	}

	return aggregatedResults
}

func (d *Database) CalculateSeriesResults(seriesID int, eventResults []Result, pointScale []int) ([]SeriesResult, error) {
	// Aggregate results by user
	seriesResults := aggregateResultsByUser(seriesID, eventResults)

	// Convert the map to a sorted slice. This sorts results by fastest time
	aggregatedResults := convertAndSortResults(seriesResults)

	// Apply point scale to the aggregated results. Calculate new points by including points from individual events and series pointscale
	aggregatedResults = applyPointScale(aggregatedResults, pointScale)

	// Reorder by points and resolve ties
	aggregatedResults = reorderByPoints(aggregatedResults)

	return aggregatedResults, nil
}

func (d *Database) CalculateAndStoreSerie(seriesID int) error {
	// Get all results for the series

	serie, err := d.GetSeries(seriesID)
	if err != nil {
		return fmt.Errorf("failed to get series: %w", err)
	}

	pointScale, err := d.parsePointScale(serie.PointScale)
	if err != nil {
		return fmt.Errorf("failed to parse point scale: %w", err)
	}

	bestEventResults, err := d.GetResultsBySeriesID(seriesID, false)
	if err != nil {
		return fmt.Errorf("failed to get best event results: %w", err)
	}

	firstEventResults, err := d.GetResultsBySeriesID(seriesID, true)
	if err != nil {
		return fmt.Errorf("failed to get first event results: %w", err)
	}

	bestEndResults, err := d.CalculateSeriesResults(seriesID, bestEventResults, pointScale)
	if err != nil {
		return fmt.Errorf("failed to calculate best end series results: %w", err)
	}

	firstEndResults, err := d.CalculateSeriesResults(seriesID, firstEventResults, pointScale)
	if err != nil {
		return fmt.Errorf("failed to calculate first end series results: %w", err)
	}

	// Store the calculated results in the database
	for _, result := range bestEndResults {
		result.HCMode = false
		if err := d.StoreSeriesResult(result); err != nil {
			return fmt.Errorf("failed to store best end series result: %w", err)
		}
	}

	for _, result := range firstEndResults {
		result.HCMode = true
		if err := d.StoreSeriesResult(result); err != nil {
			return fmt.Errorf("failed to store first end series result: %w", err)
		}
	}

	return nil
}

func (d *Database) StoreSeriesResult(result SeriesResult) error {
	query := `
        INSERT INTO series_results (
            user_id,
            race_series_id,
            points,
            race_count,
            hc_mode,
            position,
            result_time,
			created_at           
        ) VALUES (?, ?, ?, ?, ?, ?, ?, NOW());
    `
	_, err := d.db.Exec(query,
		result.UserID,
		result.RaceSeriesID,
		result.Points,
		result.RaceCount,
		result.HCMode,
		result.Position,
		result.ResultTime,
	)
	if err != nil {
		return fmt.Errorf("failed to store series result for user %s: %w", result.UserID, err)
	}

	return nil
}
