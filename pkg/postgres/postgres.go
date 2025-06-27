package postgres

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func New(ctx context.Context, databaseURI string) (*sqlx.DB, error) {
	if databaseURI == "" {
		return nil, ErrDatabaseURIRequired
	}

	db, err := sqlx.Open("pgx", databaseURI)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func Close(db *sqlx.DB) error {
	if db != nil {
		return db.Close()
	}

	return nil
}
