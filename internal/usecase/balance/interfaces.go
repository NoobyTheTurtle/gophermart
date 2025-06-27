package balance

import (
	"context"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/balance"
)

type BalanceRepo interface {
	CreateBalance(ctx context.Context, userID int) error
	GetBalanceByUserID(ctx context.Context, userID int) (*entity.UserBalance, error)
	UpdateBalance(ctx context.Context, balance *entity.UserBalance) error
	AddAccrual(ctx context.Context, userID int, accrual float64) error
}

var _ BalanceRepo = (*balance.BalancePostgresRepo)(nil)
