package balance

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/jmoiron/sqlx"
)

type BalancePostgresRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *BalancePostgresRepo {
	return &BalancePostgresRepo{db: db}
}

func (r *BalancePostgresRepo) CreateBalance(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx, createBalanceQuery, userID, 0, 0, time.Now())
	if err != nil {
		return fmt.Errorf("balanceRepo.CreateBalance: failed to create balance: %w", err)
	}

	return nil
}

func (r *BalancePostgresRepo) GetBalanceByUserID(ctx context.Context, userID int) (*entity.UserBalance, error) {
	balance := &entity.UserBalance{}

	err := r.db.GetContext(ctx, balance, getBalanceByUserIDQuery, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrBalanceNotFound
		}
		return nil, fmt.Errorf("balanceRepo.GetBalanceByUserID: failed to get balance by user ID: %w", err)
	}

	return balance, nil
}

func (r *BalancePostgresRepo) UpdateBalance(ctx context.Context, balance *entity.UserBalance) error {
	_, err := r.db.ExecContext(ctx, updateBalanceQuery, balance.Current, balance.Withdrawn, time.Now(), balance.UserID)
	if err != nil {
		return fmt.Errorf("balanceRepo.UpdateBalance: failed to update balance: %w", err)
	}

	return nil
}

func (r *BalancePostgresRepo) AddAccrual(ctx context.Context, userID int, accrual float64) error {
	_, err := r.db.ExecContext(ctx, addAccrualQuery, accrual, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("balanceRepo.AddAccrual: failed to add accrual: %w", err)
	}

	return nil
}
