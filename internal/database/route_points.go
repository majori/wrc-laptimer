package database

import "fmt"

type RoutePoint struct {
	id              int
	routeID         int
	stageDistance   int
	x               float64
	y               float64
	z               float64
	avgVelocity     float64
	avgAcceleration float64
	avgThrottle     float64
	avgRpm          float64
	avgBreaking     float64
	avgHandbreak    float64
	avgGear         float64
}

func (d *Database) GetRoutePoints(routeID int) ([]*RoutePoint, error) {
	stmt, err := d.db.PrepareContext(d.ctx, `
		SELECT
			id,
			route_id,
			stage_distance,
			x,
			y,
			z,
			avg_velocity,
			avg_acceleration,
			avg_throttle,
			avg_rpm,
			avg_breaking,
			avg_handbreak,
			avg_gear
  		FROM route_points WHERE route_id = ?
	`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(d.ctx, routeID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	routePoints := make([]*RoutePoint, 0)
	for rows.Next() {
		rp := new(RoutePoint)
		err := rows.Scan(
			&rp.id, &rp.routeID, &rp.stageDistance,
			&rp.x, &rp.y, &rp.z,
			&rp.avgVelocity, &rp.avgAcceleration, &rp.avgThrottle,
			&rp.avgRpm, &rp.avgBreaking, &rp.avgHandbreak, &rp.avgGear,
		)

		if err != nil {
			return nil, err
		}

		routePoints = append(routePoints, rp)
	}

	return routePoints, nil
}

// UpdateRoutePoints updates the route points for a given route ID.
// The update is done in place, so the number of route points must be
// the same as the number of route points in the database.
func (d *Database) UpdateRoutePoints(routeID int, routePoints []*RoutePoint) error {
	stmt, err := d.db.PrepareContext(d.ctx, `
		UPDATE route_points
		SET
			x = ?,
			y = ?,
			z = ?,
			avg_velocity = ?,
			avg_acceleration = ?,
			avg_throttle = ?,
			avg_rpm = ?,
			avg_breaking = ?,
			avg_handbreak = ?,
			avg_gear = ?
		WHERE
			route_id = ? AND
			stage_distance = ?
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, rp := range routePoints {
		_, err := stmt.ExecContext(
			d.ctx,
			rp.x, rp.y, rp.z,
			rp.avgVelocity, rp.avgAcceleration, rp.avgThrottle,
			rp.avgRpm, rp.avgBreaking, rp.avgHandbreak, rp.avgGear,
			rp.routeID, rp.stageDistance,
		)
		if err != nil {
			return fmt.Errorf("could not update a route point: %w", err)
		}
	}

	return nil
}

// ReplaceRoutePoints replaces the route points for a given route ID.
// The old route points are deleted and the new ones are inserted.
func (d *Database) ReplaceRoutePoints(routeID int, routePoints []*RoutePoint) error {
	tx, err := d.db.BeginTx(d.ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(d.ctx, `
		INSERT INTO route_points (
			route_id,
			stage_distance,
			x,
			y,
			z,
			avg_velocity,
			avg_acceleration,
			avg_throttle,
			avg_rpm,
			avg_breaking,
			avg_handbreak,
			avg_gear
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = tx.ExecContext(d.ctx, `
		DELETE
		FROM route_points
		WHERE route_id = ?
	`, routeID)

	for _, rp := range routePoints {
		_, err := stmt.ExecContext(
			d.ctx,
			rp.routeID, rp.stageDistance,
			rp.x, rp.y, rp.z,
			rp.avgVelocity, rp.avgAcceleration, rp.avgThrottle,
			rp.avgRpm, rp.avgBreaking, rp.avgHandbreak, rp.avgGear,
		)
		if err != nil {
			return fmt.Errorf("could not create a route point: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// SetRoutePoints sets the route points for a given route ID.
// Uses either UpdateRoutePoints or ReplaceRoutePoints depending on the number of route points.
func (d *Database) SetRoutePoints(routeID int, routePoints []*RoutePoint) error {
	// Get the number of current route points
	var count int
	err := d.db.QueryRowContext(d.ctx, `
		SELECT COUNT(*)
		FROM route_points
		WHERE route_id = ?
	`, routeID).Scan(&count)
	if err != nil {
		return err
	}

	if (count == len(routePoints)) {
		return d.UpdateRoutePoints(routeID, routePoints)
	}

	return d.ReplaceRoutePoints(routeID, routePoints)
}
