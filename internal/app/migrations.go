package app

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func runMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("app - runMigrations: failed to set dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("app - runMigrations: failed to run migrations: %w", err)
	}

	return nil
}
