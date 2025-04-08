package database

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func (d *Database) ListenForUserLogins(cardEvents <-chan string) {
	for id := range cardEvents {
		// Check if the user already exists
		var exists bool
		err := d.db.QueryRowContext(d.ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM users
				WHERE id = ?
			);
		`, id).Scan(&exists)
		if err != nil {
			slog.Error("error checking user existence", "error", err)
			continue
		}
		if !exists {
			// Create the user
			err = d.CreateUser(id)
			if err != nil {
				slog.Error("error creating user", "error", err)
				continue
			}
			slog.Info("user created", "id", id)
		}

		// Insert the user login into the database
		_, err = d.db.ExecContext(d.ctx, `
				INSERT INTO user_logins (user_id) 
				VALUES (?);
			`, id)
		if err != nil {
			slog.Error("error inserting user login", "error", err)
		}

		slog.Info("user logged in", "id", id)
	}
}

func (d *Database) CreateUser(id string) error {
	name := id[:5] // TODO: Get a better placeholder name
	_, err := d.db.ExecContext(d.ctx, `
		INSERT INTO users (id, name)
		VALUES (?, ?)
	`, id, name)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}
	return nil
}

func (d *Database) GetActiveUserID() (sql.NullString, error) {
	// TODO: Figure out if user has logged in but is not active (too much time from last session)
	var id sql.NullString
	err := d.db.QueryRowContext(d.ctx, `
		SELECT u.id
		FROM user_logins ul
		JOIN users u ON ul.user_id = u.id
		WHERE ul.timestamp = (
			SELECT MAX(timestamp)
			FROM user_logins
		);
	`).Scan(&id)
	if err != nil {
		// No rows means no user is logged in
		if err == sql.ErrNoRows {
			return sql.NullString{}, nil
		}
		return sql.NullString{}, fmt.Errorf("could not get active user: %w", err)
	}
	return id, nil
}
