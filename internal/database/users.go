package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/majori/wrc-laptimer/pkg/username"
)

func (d *Database) ListenForUserLogins(cardEvents <-chan string) {
	for id := range cardEvents {
		// Check if the user already exists
		var exists bool
		err := d.queryRow(`
			SELECT EXISTS (
				SELECT 1
				FROM users
				WHERE id = ?
			)
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

		// Logout previous user
		err = d.LogoutUser()
		if err != nil {
			slog.Error("error when logging out user", "error", err)
			continue
		}

		// Insert the user login into the database
		_, err = d.exec(`
			INSERT INTO user_logins (user_id) 
			VALUES (?)
		`, id)
		if err != nil {
			slog.Error("error inserting user login", "error", err)
			continue
		}

		slog.Info("user logged in", "id", id)
	}
}

func (d *Database) CreateUser(id string) error {
	name := username.GenerateFromSeed(id)
	_, err := d.exec(`
		INSERT INTO users (id, name)
		VALUES (?, ?)
	`, id, name)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}
	return nil
}

func (d *Database) GetActiveUserID() (sql.NullString, error) {
	var id sql.NullString
	err := d.queryRow(`
		SELECT u.id
		FROM user_logins ul
		JOIN users u ON ul.user_id = u.id
		WHERE ul.timestamp = (
			SELECT MAX(timestamp)
			FROM user_logins
			WHERE active IS true
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

func (d *Database) LogoutUser() error {
	_, err := d.exec(`
		UPDATE user_logins
		SET active = false
		WHERE active IS true
	`)

	if err != nil {
		return err
	}

	slog.Info("user logged out (if any)")
	return nil
}
