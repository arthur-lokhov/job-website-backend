// Package postgres provides PostgreSQL database utilities.
package postgres

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("failed to run goose migrations: %w", err)
	}

	return nil
}
