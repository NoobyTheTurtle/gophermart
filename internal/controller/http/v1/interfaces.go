package v1

import (
	"context"

	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/NoobyTheTurtle/gophermart/internal/usecase/auth"
	"github.com/NoobyTheTurtle/gophermart/internal/usecase/balance"
	"github.com/NoobyTheTurtle/gophermart/internal/usecase/order"
	"github.com/NoobyTheTurtle/gophermart/internal/usecase/withdrawal"
	"github.com/NoobyTheTurtle/gophermart/pkg/logger"
)

type Logger interface {
	Info(message string, args ...any)
	Error(message string, args ...any)
}

var _ Logger = (*logger.ZapLogger)(nil)

type OrderUseCase interface {
	UploadOrder(ctx context.Context, userID int, orderNumber string) error
	GetUserOrders(ctx context.Context, userID int) ([]*entity.Order, error)
}

var _ OrderUseCase = (*order.OrderUseCaseImpl)(nil)

type BalanceUseCase interface {
	GetUserBalance(ctx context.Context, userID int) (*entity.UserBalance, error)
}

var _ BalanceUseCase = (*balance.BalanceUseCaseImpl)(nil)

type AuthUseCase interface {
	Register(ctx context.Context, auth *entity.Auth) (token string, err error)
	Login(ctx context.Context, auth *entity.Auth) (token string, err error)
	ValidateToken(token string) (int, error)
}

var _ AuthUseCase = (*auth.AuthUseCaseImpl)(nil)

type WithdrawalUseCase interface {
	ProcessWithdrawal(ctx context.Context, userID int, orderNumber string, sum float64) error
	GetUserWithdrawals(ctx context.Context, userID int) ([]*entity.Withdrawal, error)
}

var _ WithdrawalUseCase = (*withdrawal.WithdrawalUseCaseImpl)(nil)
