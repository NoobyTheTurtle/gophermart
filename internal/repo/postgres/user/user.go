package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/jmoiron/sqlx"
)

type UserPostgresRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *UserPostgresRepo {
	return &UserPostgresRepo{db: db}
}

func (r *UserPostgresRepo) CreateUser(ctx context.Context, user *entity.User) error {
	err := r.db.QueryRowContext(ctx, createUserQuery, user.Login, user.PasswordHash, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserPostgresRepo) GetUserByLogin(ctx context.Context, login string) (*entity.User, error) {
	user := &entity.User{}

	err := r.db.GetContext(ctx, user, getUserByLoginQuery, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}

	return user, nil
}

func (r *UserPostgresRepo) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	user := &entity.User{}

	err := r.db.GetContext(ctx, user, getUserByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}
