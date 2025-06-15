package withdrawal

import (
	"context"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/balance"
	"github.com/NoobyTheTurtle/gophermart/internal/repo/postgres/withdrawal"
	"github.com/NoobyTheTurtle/gophermart/pkg/luhn"
)

type WithdrawalRepo interface {
	CreateWithdrawal(ctx context.Context, withdrawal *entity.Withdrawal) error
	GetWithdrawalsByUserID(ctx context.Context, userID int) ([]*entity.Withdrawal, error)
}

var _ WithdrawalRepo = (*withdrawal.WithdrawalPostgresRepo)(nil)

type BalanceRepo interface {
	GetBalanceByUserID(ctx context.Context, userID int) (*entity.UserBalance, error)
	UpdateBalance(ctx context.Context, balance *entity.UserBalance) error
}

var _ BalanceRepo = (*balance.BalancePostgresRepo)(nil)

type LuhnService interface {
	IsValid(number string) bool
}

var _ LuhnService = (*luhn.LuhnService)(nil)
